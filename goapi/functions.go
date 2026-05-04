package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// ! <needs_testing>
func connectToDB(host string, port int, user string, password string, dbname string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	dbConn, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("goAPI failed to connect to DB: %v", err)
	}

	// проверка соединения
	if err = dbConn.Ping(); err != nil {
		dbConn.Close()
		return nil, fmt.Errorf("goAPI failed to ping to DB: %v", err)
	}

	fmt.Println("[GO API] Successfully connected to the database!")
	return dbConn, nil
}

// ! <draft>
// ! be aware of SQL Injection
func addNewHostToDB(dbConn *sql.DB, hostgroup string, hostIP string, tableName string) error {
	query := fmt.Sprintf("INSERT INTO %s(group, ip) VALUES ($1, $2);", tableName)

	_, err := dbConn.Exec(query, hostIP, hostgroup)
	if err != nil {
		return fmt.Errorf("failed to insert row into %s: %v", tableName, err)
	}

	fmt.Printf("[GO API] Successfully inserted new value %s in %s\n", hostIP, tableName)
	return nil
}

// ! <draft>
func deleteHostFromDB(dbConn *sql.DB, hostIP string, tableName string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE ip='$1';", tableName)

	_, err := dbConn.Exec(query, hostIP)
	if err != nil {
		return fmt.Errorf("failed to delete row in %s: %v", tableName, err)
	}

	fmt.Printf("[GO API] Successfully deleted host with IP(%s) in %s\n", hostIP, tableName)
	return nil
}
