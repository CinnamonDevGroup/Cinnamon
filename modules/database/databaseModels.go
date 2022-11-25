package databaseModels

type Server struct {
	GID    string
	Joined string
	Config
}

type Config struct {
	Minecraft Minecraft
}

type Minecraft struct {
	AuthKey string
}
