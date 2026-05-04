package main

type DatabaseInfo struct {
	Config databaseConfig `toml:"database"`
}

type databaseConfig struct {
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	User         string `toml:"username"`
	Password     string `toml:"password"`
	DatabaseName string `toml:"database_name"`
}
