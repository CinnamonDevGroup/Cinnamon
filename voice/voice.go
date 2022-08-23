package voice

import (
	//	"gorm.io/gorm"

	"log"
	"os"
	"regexp"
	"strings"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	soundcloudapi "github.com/zackradisic/soundcloud-api"
	"github.com/zmb3/spotify"
	//	"gorm.io/driver/sqlite3"
)

type deets struct {
	gID  string
	cID  string
	mute bool
	deaf bool
}

type SpotifyClient func()

var (
	Commands = []discordgo.ApplicationCommand{
		{
			Name:        "play",
			Description: "Play a song",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "songuri",
					Description: "Name/URL",
					Required:    true,
				},
			},
		},
		{
			Name:        "latch",
			Description: "Latch onto your rich presence status status.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "userlatch",
					Description: "User",
					Required:    false,
				},
			},
		},
		{
			Name:        "stream",
			Description: "Stream audio from your computer.",
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, c spotify.Client){
		"play": func(s *discordgo.Session, i *discordgo.InteractionCreate, c spotify.Client) {

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

			shouldStream := false

			if option, ok := optionMap["songuri"]; ok {

				songURI = option.StringValue()

			}
			sc, err := soundcloudapi.New(soundcloudapi.APIOptions{})
			if err != nil {
				log.Panic(err)
			}

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

				if isSoundcloud {

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
					}

					shouldStream = true

				} else if isSpotify {
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

					f, err := os.Create("./tmp/dat")
					if err != nil {
						log.Panic(err)
					}

					defer f.Close()
					err = sc.DownloadTrack(tracks[0].Media.Transcodings[0], f)

					if err != nil {
						log.Panic(err)
					}

					shouldStream = true

				}

				if shouldStream {
					channelID := ""
					guild, err := s.State.Guild(i.GuildID)
					if err != nil {
						log.Panic(err)
					}

					userInVoice := false

					for _, key := range guild.VoiceStates {
						if key.UserID == i.Member.User.ID {
							channelID = key.ChannelID
							userInVoice = true
						}
					}
					if userInVoice {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Searched.",
							},
						})
						dgv, err := s.ChannelVoiceJoin(i.GuildID, channelID, false, true)
						if err != nil {
							log.Panic(err)
						}
						if err != nil {
							log.Panic(err)
						}

						s.UpdateListeningStatus("uwu!~")
						dgvoice.PlayAudioFile(dgv, "./tmp/dat", make(chan bool))
						dgv.Close()
					} else {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Please join a Voice Channel to play music!",
							},
						})
					}

				}

			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "URL format not recognized. Supported song services are: Youtube, Spotify, and Soundcloud.",
					},
				})
			}

		},
		"latch": func(s *discordgo.Session, i *discordgo.InteractionCreate, c spotify.Client) {

		},
		"stream": func(s *discordgo.Session, i *discordgo.InteractionCreate, c spotify.Client) {

		},
	}
)

func VoiceConnect(s *discordgo.Session, deets deets) {
	s.ChannelVoiceJoin(deets.gID, deets.cID, deets.mute, deets.deaf)

}
