package core_models

type User struct {
	UID            string       `gorm:"primaryKey"`
	CurrentServers []string     `gorm:"serializer:json"`
	Modules        []UserModule `gorm:"serializer:json"`
}
