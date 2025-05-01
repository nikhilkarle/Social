package main

import (
	"log"

	"github.com/nikhilkarle/social/internal/db"
	"github.com/nikhilkarle/social/internal/store"
)

func main(){
	addr := "postgres://admin:adminpassword@localhost:5433/social?sslmode=disable"
	conn, err := db.New(addr, 30, 30, "15m")

	if err != nil{
		log.Fatal(err)
	}
	
	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store)

}