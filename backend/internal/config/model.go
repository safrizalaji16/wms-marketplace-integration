package config

type Config struct {
	Server      Server
	Database    Database
	Jwt         Jwt
	Marketplace Marketplace
}

type Server struct {
	Host string
	Port string
}

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Tz       string
}

type Jwt struct {
	Key string
	Exp int
}

type Marketplace struct {
	BaseURL     string
	PartnerID   string
	PartnerKey  string
	ShopID      string
	RedirectURL string
}
