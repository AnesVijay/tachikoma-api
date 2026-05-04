package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

func main() {
	configDir := os.Getenv("CONFIG_DIR")

	var databaseInfo DatabaseInfo

	if _, err := toml.Decode(fmt.Sprintf("%s/base.toml", configDir), &databaseInfo); err != nil {
		fmt.Printf("error decoding base.toml: %v", err)
	}

	if os.Getenv("APP_ENVIRONMENT") == "production" {
		if _, err := toml.Decode(fmt.Sprintf("%s/production.toml", configDir), &databaseInfo); err != nil {
			fmt.Printf("error decoding production.toml: %v", err)
		}
	} else if os.Getenv("APP_ENVIRONMENT") == "local" {
		if _, err := toml.Decode(fmt.Sprintf("%s/local.toml", configDir), &databaseInfo); err != nil {
			fmt.Printf("error decoding local.toml: %v", err)
		}
	}

	// ! <draft>
	dbConnection, err := connectToDB(
		databaseInfo.Config.Host,
		databaseInfo.Config.Port,
		databaseInfo.Config.User,
		databaseInfo.Config.Password,
		databaseInfo.Config.DatabaseName)
	if err != nil {
		log.Fatal("[GO API] Connection to DB failed: %v", err)
	}

	

}
