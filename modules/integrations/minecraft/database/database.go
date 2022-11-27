package minecraftdb

type Minecraft struct {
	AuthKey        string `gorm:"primaryKey"`
	GID            string `gorm:"primaryKey"`
	DefaultChannel string
	Active         bool
}

type User struct {
	MCUUID     string
	MCUsername string
	MCPFP      string
	UID        []string
	AuthKey    string
}
