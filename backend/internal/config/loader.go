package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Get() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error when loading file configuration: ", err.Error())
	}

	expInt, _ := strconv.Atoi(os.Getenv("JWT_EXP"))

	return &Config{
		Server: Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		Database: Database{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Tz:       os.Getenv("DB_TZ"),
		},
		Jwt: Jwt{
			Key: os.Getenv("JWT_KEY"),
			Exp: expInt,
		},
		Marketplace: Marketplace{
			BaseURL:     os.Getenv("MARKETPLACE_BASE_URL"),
			PartnerID:   os.Getenv("MARKETPLACE_PARTNER_ID"),
			PartnerKey:  os.Getenv("MARKETPLACE_PARTNER_KEY"),
			ShopID:      os.Getenv("MARKETPLACE_SHOP_ID"),
			RedirectURL: os.Getenv("MARKETPLACE_REDIRECT_URL"),
		},
	}
}
