package voice

import (
	//	"gorm.io/gorm"

	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	soundcloudapi "github.com/zackradisic/soundcloud-api"
	"github.com/zmb3/spotify"
	"layeh.com/gopus"
	//	"gorm.io/driver/sqlite3"
)

type userQueries struct {
	trackURL []string
	ID       string
}

const (
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // max size of opus data
)

type users struct {
	userID  string
	queries []userQueries
}

type track struct {
	url  string
	uuid string
	user string
}

type serverStruct struct {
	server      string
	queue       []track
	stop        chan bool
	keepAlive   bool
	wg          sync.WaitGroup
	cancelTimer bool
	user        []users
}

type ServersStruct struct {
	servers []serverStruct
}

var (
	opusEncoder *gopus.Encoder
)

var (
	Commands = []discordgo.ApplicationCommand{
		{
			Name:        "cinplay",
			Description: "Play a song",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "song",
					Description: "Name/URL",
					Required:    true,
				},
			},
		},
		{
			Name:        "cinlatch",
			Description: "Latch onto your rich presence status status.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User",
					Required:    false,
				},
			},
		},
		{
			Name:        "cintream",
			Description: "Stream audio from your computer.",
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, c spotify.Client, z *ServersStruct){
		"cinplay": func(s *discordgo.Session, i *discordgo.InteractionCreate, c spotify.Client, z *ServersStruct) {

			defer func() {
				if err := recover(); err != nil {
					log.Println("panic occurred:", err)
				}
			}()

			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			// This example stores the provided arguments in an []interface{}
			// which will be used to format the bot's response

			// Get the value from the option map.
			// When the option exists, ok = true
			var songURI string

			if option, ok := optionMap["song"]; ok {

				songURI = option.StringValue()

			}

			defer func() {
				if err := recover(); err != nil {
					log.Println("panic occurred:", err)
				}
			}()
			regURL, _ := regexp.Compile(`com/|https://|www.a-z.`)
			shouldContinue := false

			var response *discordgo.InteractionResponse
			var trackURL []string

			if regURL.MatchString(songURI) {
				defer func() {
					if err := recover(); err != nil {
						log.Println("panic occurred:", err)
					}
				}()
				shouldContinue, response = fetchSongURL(songURI, c, s, i)
			} else {
				defer func() {
					if err := recover(); err != nil {
						log.Println("panic occurred:", err)
					}
				}()
				shouldContinue, response, trackURL = fetchSong(songURI, z, i)
			}

			if shouldContinue {
				defer func() {
					if err := recover(); err != nil {
						log.Println("panic occurred:", err)
					}
				}()
				_ = s.InteractionRespond(i.Interaction, response)

				defer func() {
					if err := recover(); err != nil {
						log.Println("panic occurred:", err)
					}
				}()
				lastInteract, _ := s.InteractionResponse(i.Interaction)

				structBuilder(trackURL, i, z, lastInteract.ID)

			}

		},
		"cinlatch": func(s *discordgo.Session, i *discordgo.InteractionCreate, c spotify.Client, z *ServersStruct) {

		},
		"cinstream": func(s *discordgo.Session, i *discordgo.InteractionCreate, c spotify.Client, z *ServersStruct) {

		},
	}
)

func OnInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate, z *ServersStruct) {

	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()
	interactionResponse := i.MessageComponentData()

	if interactionResponse.CustomID == "track" {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic occurred:", err)
			}
		}()
		trackNumber, _ := strconv.Atoi(interactionResponse.Values[0])
		trackNumber = trackNumber - 1
		var whereServerAt int
		var whereUserAt int

		for x := 0; x < len(z.servers); x++ {
			if z.servers[x].server == i.GuildID {
				whereServerAt = x
			}
		}
		for x := 0; x < len(z.servers[whereServerAt].user); x++ {
			if z.servers[whereServerAt].user[x].userID == i.Member.User.ID {
				whereUserAt = x
			}
		}

		defer func() {
			if err := recover(); err != nil {
				log.Println("panic occurred:", err)
			}
		}()

		queueTrack(trackNumber, whereServerAt, whereUserAt, i.Message.ID, i.Member.User.ID, z)

		defer voiceControl(s, i, z, whereServerAt, z.servers[whereServerAt].stop)
	}

}

func voiceControl(s *discordgo.Session, i *discordgo.InteractionCreate, z *ServersStruct, whereServerAt int, ctrl chan bool) {

	selfInVoice, userInVoice, userChannelID, selfChannelID := checkVC(s, i)
	var dgv *discordgo.VoiceConnection
	play := false
	if userInVoice {
		if selfInVoice {
			if userChannelID != selfChannelID {
				response := &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "⚠️Bot is already in use by a different Voice Channel⚠️",
					},
				}
				s.InteractionRespond(i.Interaction, response)
			}
		} else {
			dgv = voiceConnect(s, userChannelID, i.GuildID)
			play = true
		}
	} else {
		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "⚠️Please join a Voice Channel to use bot ⚠️",
			},
		}
		s.InteractionRespond(i.Interaction, response)
	}
	if play {
		connectionTimer(z, whereServerAt, dgv, s, i)
	}
}
func SendPCM(v *discordgo.VoiceConnection, pcm <-chan []int16) {
	if pcm == nil {
		return
	}

	var err error

	opusEncoder, err = gopus.NewEncoder(frameRate, channels, gopus.Audio)

	if err != nil {
		return
	}

	for {

		// read pcm from chan, exit if channel is closed.
		recv, ok := <-pcm
		if !ok {
			return
		}

		// try encoding pcm frame with Opus
		opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
		if err != nil {
			return
		}

		if !v.Ready || v.OpusSend == nil {
			// OnError(fmt.Sprintf("Discordgo not ready for opus packets. %+v : %+v", v.Ready, v.OpusSend), nil)
			// Sending errors here might not be suited
			return
		}
		// send encoded opus data to the sendOpus channel
		v.OpusSend <- opus
	}
}
func playFile(dgv *discordgo.VoiceConnection, whereServerAt int, z *ServersStruct) bool {

	filename := z.servers[whereServerAt].queue[0].uuid

	//Open File
	f, _ := os.Open("./tmp/" + filename)

	//Close file after function ends
	defer f.Close()

	run := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")

	ffmpegOut, _ := run.StdoutPipe()
	ffmpegIn, _ := run.StdinPipe()

	ffmpegOutBuffer := bufio.NewReaderSize(ffmpegOut, 16384)
	ffmpegInBuffer := bufio.NewWriterSize(f, 16384)

	run.Start()
	go io.Copy(ffmpegIn, f)

	defer run.Process.Kill()

	go func() {
		<-z.servers[whereServerAt].stop
		run.Process.Kill()
	}()

	dgv.Speaking(true)

	defer func() {
		dgv.Speaking(false)
	}()

	send := make(chan []int16, 2)

	defer close(send)

	close := make(chan bool)
	go func() {
		SendPCM(dgv, send)

		close <- true
	}()

	defer func() {
		fmt.Print("owoooowoowo")
		fmt.Print("uwu")
		z.servers[whereServerAt].wg.Done()
	}()

	go func() {
		for {
			audioBuffer := make([]int16, frameSize*channels)

			binary.Write(ffmpegInBuffer, binary.LittleEndian, &audioBuffer)
		}

	}()
	var uwu sync.WaitGroup

	uwu.Add(1)

	go func() {

		go func() {
			for {

				audioBuffer := make([]int16, frameSize*channels)

				err := binary.Read(ffmpegOutBuffer, binary.LittleEndian, &audioBuffer)

				if err == io.EOF || err == io.ErrUnexpectedEOF {
					return
				}

				select {
				case send <- audioBuffer:
				case <-close:
					return
				}
			}
		}()

		<-close
		uwu.Done()
		return

	}()
	uwu.Wait()

	return true

}

func playQueue(dgv *discordgo.VoiceConnection, whereServerAt int, z *ServersStruct) int64 {

	uwu := z.servers[whereServerAt]

	defer func() {
		fmt.Print("uwuwuuu")
		fmt.Print("uwu")
		z.servers[whereServerAt].wg.Done()
	}()
	if len(uwu.queue) >= 1 {
		defer func() { z.servers[whereServerAt].queue = RemoveIndex(z.servers[whereServerAt].queue, 0) }()
	}
	for len(z.servers[whereServerAt].queue) >= 1 {

		func() {

			z.servers[whereServerAt].cancelTimer = true
			downloadQueue(z, whereServerAt, 0)

			filename := z.servers[whereServerAt].queue[0].uuid

			_ = playFile(dgv, whereServerAt, z)

			deleteOldFile(filename)

		}()
	}

	return time.Now().Unix()

}

func connectionTimer(z *ServersStruct, whereServerAt int, dgv *discordgo.VoiceConnection, s *discordgo.Session, i *discordgo.InteractionCreate) {
	timerStartTime := time.Now().Unix()

	fmt.Print(timerStartTime)

	timeout := 30

	defer func() { z.servers[whereServerAt].queue = nil }()

	for {
		currentTime := time.Now().Unix()

		timeUntilTimeout := currentTime - timerStartTime
		fmt.Println(timeUntilTimeout)

		if timeUntilTimeout >= int64(timeout) {
			z.servers[whereServerAt].keepAlive = false
			s.ChannelMessageSend(i.ChannelID, "Bot session timed out due to no queued tracks activity.")
			dgv.Disconnect()
			return
		} else {
			if len(z.servers[whereServerAt].queue) >= 1 {
				z.servers[whereServerAt].wg.Add(1)
				go playQueue(dgv, whereServerAt, z)
				fmt.Print("uwo")
				z.servers[whereServerAt].wg.Wait()
				fmt.Print("owo")
				timerStartTime = time.Now().Unix()

			}

		}

	}

}

func RemoveIndex(s []track, index int) []track {
	return append(s[:index], s[index+1:]...)
}

func deleteOldFile(filename string) {

	os.Remove("./tmp/" + filename)
}
func downloadQueue(z *ServersStruct, whereServerAt int, whereTrackAt int) {
	sc, err := soundcloudapi.New(soundcloudapi.APIOptions{})
	if err != nil {
		log.Panic(err)
	}

	f, _ := os.Create("./tmp/" + z.servers[whereServerAt].queue[whereTrackAt].uuid)

	track := soundcloudapi.GetTrackInfoOptions{
		URL: z.servers[whereServerAt].queue[whereTrackAt].url,
	}
	trackInfo, _ := sc.GetTrackInfo(track)

	sc.DownloadTrack(trackInfo[0].Media.Transcodings[0], f)

	f.Close()

}

func checkVC(s *discordgo.Session, i *discordgo.InteractionCreate) (bool, bool, string, string) {
	selfChannelID := ""
	userChannelID := ""
	guild, err := s.State.Guild(i.GuildID)
	if err != nil {
		log.Panic(err)
	}

	userInVoice := false

	for _, key := range guild.VoiceStates {
		if key.UserID == i.Member.User.ID {
			userChannelID = key.ChannelID
			userInVoice = true
		}
	}

	selfInVoice := false

	for _, key := range guild.VoiceStates {
		if key.UserID == s.State.User.ID {
			selfChannelID = key.ChannelID
			selfInVoice = true
		}
	}
	return selfInVoice, userInVoice, userChannelID, selfChannelID
}

func queueTrack(trak int, whereServerAt int, whereUserAt int, id string, userID string, z *ServersStruct) {
	var trackURL string
	for x := 0; x < len(z.servers[whereServerAt].user[whereUserAt].queries); x++ {
		if z.servers[whereServerAt].user[whereUserAt].queries[x].ID == id {
			trackURL = z.servers[whereServerAt].user[whereUserAt].queries[x].trackURL[trak]
		}
	}

	trackStruct := track{
		url:  trackURL,
		uuid: uuid.NewString(),
		user: userID,
	}
	z.servers[whereServerAt].queue = append(z.servers[whereServerAt].queue, trackStruct)
}

func voiceConnect(s *discordgo.Session, cid string, gid string) *discordgo.VoiceConnection {

	dgv, err := s.ChannelVoiceJoin(gid, cid, false, true)
	if err != nil {
		log.Panic(err)
	}
	if err != nil {
		log.Panic(err)
	}

	return dgv
}
func fetchSoundcloudURL(songURI string, sc *soundcloudapi.API) (bool, string) {

	shouldContinue := false
	trackURL := songURI
	hasHTTPS := strings.Contains(trackURL, "https://")
	hasWWW := strings.Contains(trackURL, "www")
	if !hasHTTPS {

		if hasWWW {
			tmpString := strings.TrimPrefix(trackURL, "www.")
			trackURL = tmpString
		}

		tmpString := "https://" + trackURL
		trackURL = tmpString
	} else if hasWWW {
		tmpString := strings.TrimPrefix(trackURL, "https://www.")
		trackURL = "https://" + tmpString
	}

	tracks, err := sc.GetTrackInfo(soundcloudapi.GetTrackInfoOptions{
		URL: trackURL,
	})
	if err != nil {
		log.Panic(err)
	}

	f, err := os.Create("./tmp/dat2")
	if err != nil {
		log.Panic(err)
	}

	defer f.Close()
	err = sc.DownloadTrack(tracks[0].Media.Transcodings[0], f)
	if err != nil {
		log.Panic(err)
	} else {
		shouldContinue = true
	}

	f.Close()
	return shouldContinue, tracks[0].Title + tracks[0].User.Username
}

func FetchSpotifyURL(songURI string, c spotify.Client, sc *soundcloudapi.API) (bool, string) {

	shouldContinue := false
	suffix := strings.Index(songURI, "?")
	prefix := strings.Index(songURI, "track/")

	songID := spotify.ID(songURI[prefix+6 : suffix])

	track, err := c.GetTrack(songID)
	if err != nil {
		log.Panic(err)
	}

	songName := track.Name
	songArtists := track.Artists
	songArtist := songArtists[0]
	songArtistName := songArtist.Name

	searchQuery := songName + " " + songArtistName

	var options soundcloudapi.SearchOptions
	options.Query = searchQuery

	searchResults, err := sc.Search(options)

	if err != nil {
		log.Panic(err)
	}

	searchTracks, err := searchResults.GetTracks()
	if err != nil {
		log.Print(err)
	}

	trackResult := searchTracks[0]

	trackURL := trackResult.PermalinkURL
	tracks, err := sc.GetTrackInfo(soundcloudapi.GetTrackInfoOptions{
		URL: trackURL,
	})
	if err != nil {
		log.Panic(err)
	}

	f, err := os.Create("./tmp/dat2")
	if err != nil {
		log.Panic(err)
	}

	defer f.Close()
	err = sc.DownloadTrack(tracks[0].Media.Transcodings[0], f)
	if err != nil {
		log.Panic(err)
	} else {
		shouldContinue = true
	}

	f.Close()
	return shouldContinue, tracks[0].Title + tracks[0].User.Username
}

func fetchSong(songURI string, z *ServersStruct, i *discordgo.InteractionCreate) (bool, *discordgo.InteractionResponse, []string) {

	var response *discordgo.InteractionResponse
	sc, err := soundcloudapi.New(soundcloudapi.APIOptions{})
	if err != nil {
		log.Panic(err)
	}

	searchQuery := songURI

	var options soundcloudapi.SearchOptions
	options.Query = searchQuery

	searchResults, err := sc.Search(options)

	if err != nil {
		log.Panic(err)
	}

	searchTracks, err := searchResults.GetTracks()

	var trackURL []string
	for x := 0; x < len(searchTracks); x++ {
		trackURL = append(trackURL, searchTracks[x].PermalinkURL)
	}

	if err != nil {
		log.Print(err)
	}

	response = compileTracks(searchTracks)

	return true, response, trackURL

}

func structBuilder(trackURL []string, i *discordgo.InteractionCreate, z *ServersStruct, id string) {
	var exists bool
	var whereServerAt int
	var whereUserAt int

	for x := 0; x < len(z.servers); x++ {
		if z.servers[x].server == i.Interaction.GuildID {
			exists = true
			whereServerAt = x
		}
	}

	if exists {
		exists = false
		for x := 0; x < len(z.servers[whereServerAt].user); x++ {
			if z.servers[whereServerAt].user[x].userID == i.Member.User.ID {
				whereUserAt = x
				exists = true
			}
		}
		if exists {

			userQueriesStruct := userQueries{
				trackURL: trackURL,
				ID:       id,
			}

			z.servers[whereServerAt].user[whereUserAt].queries = append(z.servers[whereServerAt].user[whereUserAt].queries, userQueriesStruct)

		} else {
			userQueriesStruct := userQueries{
				trackURL: trackURL,
				ID:       id,
			}
			userQueriesIndexedStruct := []userQueries{
				userQueriesStruct,
			}

			user := users{
				userID:  i.Member.User.ID,
				queries: userQueriesIndexedStruct,
			}
			z.servers[whereServerAt].user = append(z.servers[whereServerAt].user, user)
		}
	} else {

		userQueriesStruct := userQueries{
			trackURL: trackURL,
			ID:       id,
		}
		userQueriesIndexedStruct := []userQueries{
			userQueriesStruct,
		}
		userStruct := users{
			userID:  i.Member.User.ID,
			queries: userQueriesIndexedStruct,
		}
		userIndexedStruct := []users{
			userStruct,
		}
		serverStruct := serverStruct{
			server: i.GuildID,
			user:   userIndexedStruct,
		}
		z.servers = append(z.servers, serverStruct)
	}
}

func fetchSongURL(songURI string, c spotify.Client, s *discordgo.Session, i *discordgo.InteractionCreate) (bool, *discordgo.InteractionResponse) {

	var response *discordgo.InteractionResponse
	sc, err := soundcloudapi.New(soundcloudapi.APIOptions{})
	if err != nil {
		log.Panic(err)
	}

	shouldContinue := false
	regURL, err := regexp.Compile(`spotify\.com|soundcloud\.com`)

	if err != nil {
		log.Panic(err)
	}

	isURL := regURL.MatchString(songURI)

	if isURL {
		regSoundcloud, err := regexp.Compile(`soundcloud\.com`)
		if err != nil {
			log.Panic(err)
		}

		isSoundcloud := regSoundcloud.MatchString(songURI)

		regSpotify, err := regexp.Compile(`spotify\.com`)

		if err != nil {
			log.Panic(err)
		}

		isSpotify := regSpotify.MatchString(songURI)

		var trackName string

		if isSoundcloud {

			shouldContinue, trackName = fetchSoundcloudURL(songURI, sc)
			response = &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Play track " + trackName + "?",
				},
			}

		} else if isSpotify {

			shouldContinue, trackName = FetchSpotifyURL(songURI, c, sc)
			response = &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Play track " + trackName + "?",
				},
			}

		} else {
			response = &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "URL format not recognized. Supported song services are: Youtube, Spotify, and Soundcloud.",
				},
			}
		}
	}
	return shouldContinue, response
}

func compileTracks(tracks []soundcloudapi.Track) *discordgo.InteractionResponse {
	var trackNames []string

	for x := 0; x < len(tracks); x++ {
		track := tracks[x]
		trackString := track.Title + " - " + track.User.Username
		trackNames = append(trackNames, trackString)
	}

	var selectionOptions []discordgo.SelectMenuOption

	var messageComponent []discordgo.MessageComponent

	for x := 0; x < len(trackNames); x++ {

		if len(trackNames[x]) > 99 {
			tmpVar := trackNames[x]

			trackNames[x] = tmpVar[:99]
		}
		option := discordgo.SelectMenuOption{
			Label: trackNames[x],
			Value: fmt.Sprint(x + 1),
		}
		selectionOptions = append(selectionOptions, option)

	}
	selection := discordgo.SelectMenu{
		CustomID:    "track",
		Placeholder: "Select the track.",
		Options:     selectionOptions,
	}

	messageComponent = append(messageComponent, selection)
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "The following tracks were found, please select the track.`",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: messageComponent,
				},
			},
		},
	}

	return response

}
