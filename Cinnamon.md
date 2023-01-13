└── modules
   ├── core
   │   ├── basic_handlers
   │   │   └── [core_handlers]
   │   │        └── core_handlers.go
   │   │            ├── func regenAuthKey(string) (string)
   │   │            ├── func updateDBAuthKeys(GID string, AuthKey string, OldKey string) ()
   │   │            ├── func OnServerJoin(z *discordgo.GuildCreate) (discord_handler) ()
   │   │            └── var
   │   │                ├── admin
   │   │                ├── Commands
   │   │                │   ├── regenauthkey
   │   │                │   └── authkey
   │   │                └── map[string]func(i *discordgo.InteractionCreate) ()
   │   │                    ├── "regenauthkey"
   │   │                    └── "authkey"
   │   ├── common
   │   │   └── [common]
   │   │        ├── common_structs.go
   │   │        │   ├── type Data struct│
   │   │        │   │   ├── Token           string
   │   │        │   │   ├── AdminServer     string
   │   │        │   │   ├── AdminChannel    string
   │   │        │   │   ├── InfoChannel     string
   │   │        │   │   ├── WarnChannel     string
   │   │        │   │   ├── ErrChannel      string
   │   │        │   │   ├── UpdateChannel   string
   │   │        │   │   └── FeedbackChannel string
   │   │        │   └── type LogEntry struct
   │   │        │       ├── Time    time.Time
   │   │        │       ├── Message string
   │   │        │       └── Level   string
   │   │        └── common_utils.go
   │   │            ├── func BabbleWords()                          (string)
   │   │            ├── func AddGuildToDB(GID string)           ()
   │   │            ├── const
   │   │            │   ├── LogError
   │   │            │   ├── LogWarning
   │   │            │   ├── LogInfo
   │   │            │   ├── LogUpdate
   │   │            │   └── LogFeedback
   │   │            ├── func GetGuildName(GID string)               (string)
   │   │            ├── func GetGuildOwnerName(GID string)          (string)
   │   │            ├── func LogEvent(message string, level string) ()
   │   │            ├── func RecoverPanic(channelID string)         ()
   │   │            └── func CheckGuildExists(GID string)           ()
   │   ├── database
   │   │   ├── core_models
   │   │   │   └── [core_models]
   │   │   │        ├── cinnamon_db.go
   │   │   │        │   ├── type Cinnamon struct
   │   │   │        │   │   ├── BotID               string
   │   │   │        │   │   ├── TotalUsers          []UserStats
   │   │   │        │   │   ├── TotalServers        []ServerStats
   │   │   │        │   │   ├── Uptime              int64
   │   │   │        │   │   ├── CreationDate        int64
   │   │   │        │   │   ├── UpSince             int64
   │   │   │        │   │   ├── TotalUptime         int64
   │   │   │        │   │   ├── TotalDowntime       int64
   │   │   │        │   │   ├── DowntimePercent     float64
   │   │   │        │   │   └── PastUptime          []PastUptime
   │   │   │        │   ├── type UserStats struct
   │   │   │        │   │   ├── UserCount int
   │   │   │        │   │   └── TimeCount int64
   │   │   │        │   ├── type ServerStats struct
   │   │   │        │   │   ├── ServerCount int
   │   │   │        │   │   └── TimeCount   int64
   │   │   │        │   ├── type MessageStats struct
   │   │   │        │   │   ├── TimeCount    int64
   │   │   │        │   │   └── MessageCount int
   │   │   │        │   ├── type PastUptime struct
   │   │   │        │   │   ├── Uptime  int64
   │   │   │        │   │   └── UpSince int64
   │   │   │        │   ├── type UserModule struct
   │   │   │        │   │   ├── Service string
   │   │   │        │   │   ├── UUID    string
   │   │   │        │   │   ├── Data    json.RawMessage
   │   │   │        │   │   ├── UID     string
   │   │   │        │   │   └── AuthKey string
   │   │   │        │   └── type ServerModule struct
   │   │   │        │       ├── Service string
   │   │   │        │       ├── UUID    string
   │   │   │        │       ├── Data    json.RawMessage
   │   │   │        │       ├── UID     string
   │   │   │        │       └── AuthKey string
   │   │   │        ├── guild_db.go
   │   │   │        │   ├── type Guild struct
   │   │   │        │   │   ├── GID     string
   │   │   │        │   │   ├── Joined  int64
   │   │   │        │   │   ├── AuthKey string
   │   │   │        │   │   ├── About   Information
   │   │   │        │   │   ├── Config  Config
   │   │   │        │   │   └── Modules []ServerModule
   │   │   │        │   ├── type Config struct
   │   │   │        │   │   └── DefaultAdminChannel string
   │   │   │        │   ├── type Information struct
   │   │   │        │   │   ├── JoinedAt        []int64
   │   │   │        │   │   ├── LeftAt          []int64
   │   │   │        │   │   ├── UserAmount      int
   │   │   │        │   │   ├── UserInfo        []GuildUser
   │   │   │        │   │   └── MessageCount    []Message
   │   │   │        │   ├── type Message struct
   │   │   │        │   │   ├── MessageCount int
   │   │   │        │   │   └── TimeCount    int64
   │   │   │        │   ├── type GuildUser struct
   │   │   │        │   │   ├── UID             string
   │   │   │        │   │   ├── UUID            string
   │   │   │        │   │   ├── JoinedPositions []int
   │   │   │        │   │   ├── LeftPositions   []int
   │   │   │        │   │   └── Modules         []UserModule
   │   │   │        └── user_db.go
   │   │   │            ├── type User struct
   │   │   │            │   ├── UID             string
   │   │   │            │   ├── CurrentServers  []string
   │   │   │            │   └── UserModules     []UserModule
   │   │   └── [database]
   │   │        └── database.go
   │   │            └── func InitDB() (*gorm.DB)
   │   ├── discord
   │   │   └── [discord_session]
   │   │        └──discord_session.go
   │   │           └── func InitSession() (*discordgo.Session) 
   │   └── websocket
   │       └── [websocket]
   │            ├── client.go
   │            │   ├── func (c *Client) readPump()     ()
   │            │   ├── func (c *Client) writePump()    ()
   │            │   ├── func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) ()
   │            │   └── func GenUserID() (string)
   │            ├── hub.go
   │            │   ├── type Hub struct
   │            │   │   ├── Clients     map[*Client]bool
   │            │   │   ├── Broadcast   chan    IncomingData
   │            │   │   ├── Client      chan    *Client
   │            │   │   ├── Register    chan    *Client
   │            │   │   ├── Unregister  chan    *Client
   │            │   ├── var WriteToWebsocket = make(chan ClientCache)
   │            │   ├── func newHub() (*Hub)
   │            │   ├── var WebsocketHandlers map[string]func(receivedData IncomingData, h *Hub)
   │            │   ├── type ClientCache struct
   │            │   │   ├── Client          *Client
   │            │   │   ├── UUID            string
   │            │   │   ├── OutboundData    OutboundData
   │            │   │   ├── GID             string
   │            │   │   ├── DefaultChannel  string
   │            │   │   ├── AuthKey         string
   │            │   │   └── Service         string
   │            │   ├── var ClientsCache []ClientCache
   │            │   ├── func (h *Hub) run() ()
   │            ├── models.go
   │            │   ├── type Client struct
   │            │   │   ├── Hub         *Hub
   │            │   │   ├── Conn        *websocket.Conn
   │            │   │   ├── Send        chan    []byte
   │            │   │   ├── Authenticated       bool
   │            │   │   └── User
   │            │   ├── const APINegotiate
   │            │   ├── const maxAPI
   │            │   ├── type NegotiateAPI struct
   │            │   │   └── APIVersion  int
   │            │   ├── type ConnectionStatus struct
   │            │   │   ├── AuthKey string
   │            │   │   ├── GID     string
   │            │   │   └── Status  int
   │            │   ├── type IncomingData struct
   │            │   │   ├── DataType    string
   │            │   │   ├── RawData     json.RawMessage
   │            │   │   ├── APIVersion  int
   │            │   │   └── Client      *Client
   │            │   ├── type OutboundData struct
   │            │   │   ├── DataType    string
   │            │   │   ├── RawData     json.RawMessage
   │            │   │   ├── APIVersion  int
   │            │   │   └── UUID        string
   │            │   ├── type User struct
   │            │   │   ├── AuthKey     string
   │            │   │   ├── APIVersion  int
   │            │   │   ├── Service     string
   │            │   │   ├── Addr        string
   │            │   │   ├── EnterAt     time.Time
   │            │   │   └── UUID        string
   │            │   ├── const
   │            │   │   ├── writeWait
   │            │   │   ├── pongWait
   │            │   │   └── pingPeriod
   │            │   ├── var newline
   │            │   └── var upgrader websocket.Upgrader
   │            └── websocket_server.go
   │                └── func InitWebsocket() ()
   └── integrations
       ├── affection
       │   └── database
       │       └── [affection_db]
       │            └── affection_db.go
       │                ├── type PerUserPreferences struct
       │                │   ├── UID                 string
       │                │   ├── OverrideServerAffct []string
       │                │   └── AffectionAllowed    bool
       │                ├── type PerServerPreferences struct
       │                │   ├── AffectionAllowed    bool
       │                │   └── GID                 string
       │                └── type AffectData struct
       │                    ├── DefaultAffect       bool
       │                    ├── UserPreferences     []PerUserPreferences
       │                    └── ServerPreferences   []PerServerPreferences
       ├── experience
       │   └── database
       │       └── [experience_db]
       │            └── experience_db.go
       │                 type GlobalXP struct 
       │                │   ├── TotalXP     int
       │                │   └── PerServerXP []GuildXP
       │                ├── type GuildXP struct
       │                │   ├── GID string
       │                │   └── XP  int
       │                └── type TotalXP struct
       │                    └── TotalXP int64
       ├── minecraft
       │   ├── api
       │   │   └── v1
       │   │       └── [minecraft_api_v1]
       │   │            ├── inbound_models.go
       │   │            │   ├── type ChatMessage struct
       │   │            │   │   ├── UUID    string
       │   │            │   │   ├── Message string
       │   │            │   │   ├── Mention string
       │   │            │   │   └── Channel string
       │   │            │   ├── type PlayerJoin struct
       │   │            │   │   └── UUID    string
       │   │            │   ├── type Data struct
       │   │            │   │   ├── DataType    string
       │   │            │   │   ├── RawData     json.RawMessage
       │   │            │   │   ├── AuthKey     string
       │   │            │   │   └── GID         string
       │   │            │   ├── type Authenticate struct
       │   │            │   │   ├── AuthKey         string
       │   │            │   │   ├── DefaultChannel  string
       │   │            │   │   └── GuildID         string
       │   │            │   ├── type ConnectionStatus struct
       │   │            │   │   ├── AuthKey string
       │   │            │   │   ├── GID     string
       │   │            │   │   └── Status  int
       │   │            │   └── type DiscordMessage struct
       │   │            │       ├── User    string
       │   │            │       ├── Mention string
       │   │            │       ├── Message string
       │   │            │       └── Channel string
       │   │            └── outbound_models.go
       │   │                ├── const AuthKickEvent
       │   │                ├── type KickForAuth struct
       │   │                │   ├── UUID    string
       │   │                │   └── AuthKey string
       │   │                ├── const PlayerAuthEvent
       │   │                ├── type PlayerAuthSuccessful struct
       │   │                │   ├── UUID        string
       │   │                │   └── Username    string
       │   │                ├── const NotFoundEvent
       │   │                ├── type KickForNotOnServer struct
       │   │                │   └── UUID    string
       │   │                └── type OutboundData struct
       │   │                    ├── DataType    string
       │   │                    ├── RawData     json.RawMessage
       │   │                    └── API         int
       │   ├── database
       │   │   └── [minecraft_db]
       │   │        └── minecraft_db.go
       │   │            ├── type Minecraft struct
       │   │            │   ├── AuthKey         string
       │   │            │   ├── GID             string
       │   │            │   ├── DefaultChannel  string
       │   │            │   └── Active          bool
       │   │            ├── type User struct
       │   │            │   ├── MCUUID      string
       │   │            │   ├── MCUsername  string
       │   │            │   ├── MCPFP       string
       │   │            │   ├── CurrentID   ID
       │   │            │   └── OldIDs      []ID
       │   │            └── type ID struct
       │   │                ├── UID     string
       │   │                └── UUID    string
       │   ├── discord
       │   │   └── [minecraft_discord]
       │   │        └── minecraft_discord.go
       │   │            ├── func checkGuildExists(GID string) (bool)
       │   │            ├── func addGuildToAdd(GID string) ()
       │   │            ├── func unlinkServer(GID string) ()
       │   │            ├── func RegenAuthKeys(GID string, AuthKey string, OldKey string) ()
       │   │            └── func deleteGuildData(GID string) ()
       │   └── websocket
       │       └── [minecraft_websocket]
       │            └── minecraft_websocket.go
       │                ├── var WebsocketHandler = map[string]func(data websocket.IncomingData, h *websocket.Hub)
       │                │   └── minecraft
       │                ├── func onPlayerMessage(data minecraft_api_v1.Data)
       │                ├── func clientAuthenticate(client *websocket.Client, h *websocket.Hub, responseData json.RawMessage)
       │                ├── func checkPlayerAuth(UUID string, cache websocket.ClientCache) (bool, coredb.Service)
       │                ├── func KickNewPlayerAuth(UUID string, cache websocket.ClientCache)
       │                ├── func kickPlayerAuth(UUID string, Cache websocket.ClientCache)
       │                ├── func DecideAuth(UUID string, Cache websocket.ClientCache)
       │                ├── func authFail(UUID string, Cache websocket.ClientCache)
       │                ├── func authSuccess(UUID string, username string, Cache websocket.ClientCache)
       │                └── func onPlayerJoin(m []byte, Cache websocket.ClientCache)
       ├── moderation
       │   └── database
       │       └── [moderation_db]
       │            ├── type UserModeration struct
       │            │   ├── Warnings    []Warning
       │            │   ├── Mutes       []Mute
       │            │   ├── Kicks       []Kick
       │            │   ├── Bans        []Ban
       │            │   ├── IsMember    bool
       │            │   ├── ActiveBan   bool
       │            │   └── ActiveMute  bool
       │            ├── type Warning struct
       │            │   ├── WarningReason   string
       │            │   └── WarnedAt        int64 
       │            ├── type Ban struct
       │            │   ├── BanReason   string
       │            │   └── BannedOn    int64
       │            ├── type Kick
       │            │   ├── KickReason  string
       │            │   └── KickedOn    string
       │            └── type Mute struct
       │                ├── MuteReason      string
       │                ├── MuteDuration    int64
       │                ├── MuteStartedAt   int64
       │                └── MuteExpired     bool
       └── personalization
           └── database
               └── [personalization_db]
                    └── personalization_db.go
                        ├── type PerServerPreferences struct 
                        │   ├── Pronouns    Pronouns
                        │   ├── GID         string
                        │   └── Nickname    string
                        ├── type PerUserPreferences struct
                        │   ├── UID             string
                        │   ├── Pronouns        Pronouns
                        │   ├── Nickname        string
                        │   └── OverrideServer  []string
                        ├── type Pronouns struct
                        │   ├── Nominative  string
                        │   ├── Objective   string
                        │   └── Possessive  string
                        └── type PersonalizationData struct
                            ├── ServerPreferences   []PerServerPreferences
                            ├── DefaultPronouns     []Pronouns
                            └── DefaultNickname     string
