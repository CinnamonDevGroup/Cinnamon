package minecraftdb

type Minecraft struct {
	AuthKey string `gorm:"primaryKey"`
	GID     string `gorm:"primaryKey"`
	Users   []user `gorm:"serializer:json"`
}

type user struct {
	MCUUID     string
	MCUsername string
	MCPFP      string
	UID        string
	AuthKey    string
	UUID       []string
}
