package provider

import "time"

type Provider struct {
	Data         `bson:",inline"`
	ProviderType string    `bson:"type"`
	LastUpdated  time.Time `bson:"lastUpdated"`
}

type Data struct {
	Id   string `bson:"id"`
	Name string `bson:"name"`
}
