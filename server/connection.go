package main

import (
	c "github.com/ostafen/clover/v2"
)

// It connects to the database.
func ConnectClover() *c.DB {
	db, _ := c.Open("clover-db")
	return db
}
