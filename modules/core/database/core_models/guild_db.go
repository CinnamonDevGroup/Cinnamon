package core_models

type Guild struct {
	GID     string `gorm:"primaryKey"`
	Joined  int64
	AuthKey string
	About   Information    `gorm:"serializer:json"`
	Config  Config         `gorm:"serializer:json"`
	Modules []ServerModule `gorm:"serializer:json"`
}

type Config struct {
	DefaultAdminChannel  string
	DefaultUpdateChannel string
}

type Information struct {
	JoinedAt     []int64
	LeftAt       []int64
	UserAmount   int
	UserInfo     []GuildUser
	MessageCount []Message
}

type Message struct {
	MessageCount int
	TimeCount    int64
}

type GuildUser struct {
	UID             string
	UUID            string
	JoinedPositions []int
	LeftPositions   []int
	Modules         []UserModule
}
