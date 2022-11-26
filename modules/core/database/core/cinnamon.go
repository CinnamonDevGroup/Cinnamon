package coredb

func Init() {

}

type Cinnamon struct {
	BotID           string         `gorm:"primaryKey"`
	TotalUsers      []UserStats    `gorm:"serializer:json"`
	TotalServers    []ServerStats  `gorm:"serializer:json"`
	TotalMessages   []MessageStats `gorm:"serializer:json"`
	TotalXP         []XPStats      `gorm:"serializer:json"`
	Uptime          int64
	UpSince         int64
	TotalUptime     int64
	TotalDowntime   int64
	DowntimePercent float64
	PastUptime      []PastUptime `gorm:"serializer:json"`
}

type UserStats struct {
	UserCount int
	TimeCount int64
}

type XPStats struct {
	XPCount   int
	TimeCount int64
}

type ServerStats struct {
	ServerCount int
	TimeCount   int64
}

type MessageStats struct {
	TimeCount    int64
	MessageCount int
}

type PastUptime struct {
	Uptime  int64
	UpSince int64
}
