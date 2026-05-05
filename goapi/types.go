package main

type Config struct {
	DBConf databaseConfig `toml:"database"`
	GoAPI  goapiConfig    `toml:"goapi"`
}

type databaseConfig struct {
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	User         string `toml:"username"`
	Password     string `toml:"password"`
	DatabaseName string `toml:"database_name"`
}

type goapiConfig struct {
	Token string `toml:"token"`
	Port  int64  `toml:"port"`
}

// --------------- HTTP handling ---------------

type HostGroup int

const (
	DATAPK HostGroup = iota
	ITM
)

var hostGroups = map[HostGroup]int{
	DATAPK: 0,
	ITM:    1,
}

func (hg HostGroup) Int() int {
	return hostGroups[hg]
}

type GoAPIRequest struct {
	HostIP    string    `json:"hostip"`
	Hostname  string    `json:"hostname"`
	HostGroup HostGroup `json:"hostgrp"`
}

type GoAPIResponse struct {
	Result string `json:"result"`
	Msg    string `json:"msg"`
}
