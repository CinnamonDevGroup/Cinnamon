package cinnamonModel

type Cinnamon struct {
	TotalUsers      []userStats    `gorm:"serializer:json"`
	TotalServers    []serverStats  `gorm:"serializer:json"`
	TotalMessages   []messageStats `gorm:"serializer:json"`
	TotalXP         []xpStats      `gorm:"serializer:json"`
	Uptime          int
	UpSince         int
	TotalUptime     int
	TotalDowntime   int
	DowntimePercent int
	PastUptime      []pastUptime `gorm:"serializer:json"`
}

type userStats struct {
	UserCount int
	TimeCount int
}

type xpStats struct {
	XPCount   int
	TimeCount int
}

type serverStats struct {
	ServerCount int
	TimeCount   int
}

type messageStats struct {
	TimeCount    int
	MessageCount int
}

type pastUptime struct {
	Uptime  int
	UpSince int
}
