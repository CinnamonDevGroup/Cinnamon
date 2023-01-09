package minecraft_db

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
	CurrentID  ID
	OldIDs     []ID
}

type ID struct {
	UID  string
	UUID string
}
