package conf

type BusinessConfig struct {
	Domain string `json:"domain"`
	OldDomain string `json:"oldDomain"`
	OldMongodbPassword string `json:"oldMongodbPassword"`
}
