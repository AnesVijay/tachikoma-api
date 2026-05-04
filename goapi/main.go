package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
)

var DBConnection *sql.DB
var AppConfig Config

func main() {
	configDir := os.Getenv("CONFIG_DIR")

	if _, err := toml.Decode(fmt.Sprintf("%s/base.toml", configDir), &AppConfig); err != nil {
		fmt.Printf("[GO API] error decoding base.toml: %v\n", err)
	}

	if os.Getenv("APP_ENVIRONMENT") == "production" {
		if _, err := toml.Decode(fmt.Sprintf("%s/production.toml", configDir), &AppConfig); err != nil {
			fmt.Printf("[GO API] error decoding production.toml: %v\n", err)
		}
	} else if os.Getenv("APP_ENVIRONMENT") == "local" {
		if _, err := toml.Decode(fmt.Sprintf("%s/local.toml", configDir), &AppConfig); err != nil {
			fmt.Printf("[GO API] error decoding local.toml: %v\n", err)
		}
	}

	// ! <draft>
	dbConnection, err := connectToDB(
		AppConfig.DBConf.Host,
		AppConfig.DBConf.Port,
		AppConfig.DBConf.User,
		AppConfig.DBConf.Password,
		AppConfig.DBConf.DatabaseName)
	if err != nil {
		log.Fatalf("[GO API] Connection to DB failed: %v\n", err)
	}
	DBConnection = dbConnection

	mux := http.NewServeMux()
	mux.HandleFunc("POST /addhost", addhost)
	mux.HandleFunc("POST /removehost", removehost)

	fmt.Printf("[GO API] Listening on port %s\n", AppConfig.goAPI.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", AppConfig.goAPI.port), mux))
}

func addhost(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	if token == AppConfig.goAPI.token {

		var req GoAPIRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		var result string

		e := addNewHostToDB(DBConnection, req.HostGroup.String(), req.HostIP, "hosts")
		if e != nil {
			result = fmt.Sprintf("Failed to add host <%s> to DB: %v", req.HostIP, e)
		} else {
			result = fmt.Sprintf("Successfully added a new host <%s> to DB", req.HostIP)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GoAPIResponse{Result: result})
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GoAPIResponse{Result: "Invalid token"})
	}
}

func removehost(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	if token == AppConfig.goAPI.token {

		var req GoAPIRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		var result string

		e := deleteHostFromDB(DBConnection, req.HostIP, "hosts")
		if e != nil {
			result = fmt.Sprintf("Failed to remove host <%s> from DB: %v", req.HostIP, e)
		} else {
			result = fmt.Sprintf("Successfully removed a host <%s> from DB", req.HostIP)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GoAPIResponse{Result: result})
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GoAPIResponse{Result: "Invalid token"})
	}
}
