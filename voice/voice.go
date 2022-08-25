package voice

import (
	//	"gorm.io/gorm"

	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	soundcloudapi "github.com/zackradisic/soundcloud-api"
	"github.com/zmb3/spotify"
	//	"gorm.io/driver/sqlite3"
)

type userQueries struct {
	trackURL []string
	ID       string
}

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
	server string
	queue  []track

	user []users
}

type ServersStruct struct {
	servers []serverStruct
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

func fetchSong(songURI string, z ServersStruct, i *discordgo.InteractionCreate) (bool, *discordgo.InteractionResponse, []string) {

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

func structBuilder(trackURL []string, i *discordgo.InteractionCreate, z ServersStruct, id string) {
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
			server: i.Member.GuildID,
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

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, c spotify.Client, z ServersStruct){
		"cinplay": func(s *discordgo.Session, i *discordgo.InteractionCreate, c spotify.Client, z ServersStruct) {

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

				fmt.Println(i.Interaction.ID, "last ", lastInteract.ID)

			}

		},
		"cinlatch": func(s *discordgo.Session, i *discordgo.InteractionCreate, c spotify.Client, z ServersStruct) {

		},
		"cinstream": func(s *discordgo.Session, i *discordgo.InteractionCreate, c spotify.Client, z ServersStruct) {

		},
	}
)

func OnInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate, z ServersStruct, ctrl chan bool) {
	if i.Interaction.Type == 3 {

		defer func() {
			if err := recover(); err != nil {
				log.Println("panic occurred:", err)
			}
		}()
		interactionResponse := i.MessageComponentData()

		if interactionResponse.CustomID == "track" {
			fmt.Print(interactionResponse.Values[0])
			defer func() {
				if err := recover(); err != nil {
					log.Println("panic occurred:", err)
				}
			}()
			trackNumber, _ := strconv.Atoi(interactionResponse.Values[0])
			var whereServerAt int
			var whereUserAt int

			for x := 0; x < len(z.servers); x++ {
				if z.servers[x].server == i.Member.GuildID {
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

			defer func() {
				if err := recover(); err != nil {
					log.Println("panic occurred:", err)
				}
			}()

			voiceControl(s, i, z, whereServerAt, ctrl)
		}

	}
}

func voiceControl(s *discordgo.Session, i *discordgo.InteractionCreate, z ServersStruct, whereServerAt int, ctrl chan bool) {

	selfInVoice, userInVoice, userChannelID, selfChannelID := checkVC(s, i)
	var dgv *discordgo.VoiceConnection
	play := false
	if userInVoice {
		if selfInVoice {
			if userChannelID != selfChannelID {
				//todo insert handling for bot in use in different channel
			}
		} else {
			dgv = voiceConnect(s, userChannelID, i.Member.GuildID)
			play = true
		}
	} else {
		//todo insert handling for user not in voicechannel
	}
	if play {
		defer playQueue(dgv, whereServerAt, z, ctrl)
	}
}

func playQueue(dgv *discordgo.VoiceConnection, whereServerAt int, z ServersStruct, ctrl chan bool) {

	for x := 0; x < len(z.servers[whereServerAt].queue); x++ {

		if x == 0 {
			downloadQueue(z, whereServerAt, x)
			downloadQueue(z, whereServerAt, x+1)
		} else {
			oldFilename := z.servers[whereServerAt].queue[x-1].uuid
			deleteOldFile(oldFilename, whereServerAt, x-1)
			downloadQueue(z, whereServerAt, x+1)
		}
		filename := z.servers[whereServerAt].queue[x].uuid

		os.Open("./tmp/" + filename)
		dgvoice.PlayAudioFile(dgv, "./tmp/"+filename, ctrl)

	}

}

func deleteOldFile(filename string, whereServerAt int, whereTrackAt int) {

	os.Remove("./tmp/" + filename)
}
func downloadQueue(z ServersStruct, whereServerAt int, whereTrackAt int) {
	sc, err := soundcloudapi.New(soundcloudapi.APIOptions{})
	if err != nil {
		log.Panic(err)
	}

	f, _ := os.Open("./tmp/" + z.servers[whereServerAt].queue[whereTrackAt].uuid)

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

func queueTrack(trak int, whereServerAt int, whereUserAt int, id string, userID string, z ServersStruct) {
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
