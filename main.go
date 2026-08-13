package main

import (
	"fmt"
	"log"
	"stayt/internal/db"
)

func main() {
	pool, ctx, err := db.DbConnect()
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	fmt.Println("Connection with database sucessfully established!")
}
