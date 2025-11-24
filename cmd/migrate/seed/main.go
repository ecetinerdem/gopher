package main

import (
	"log"

	"github.com/ecetinerdem/gopherSocial/internal/db"
	"github.com/ecetinerdem/gopherSocial/internal/env"
	"github.com/ecetinerdem/gopherSocial/internal/store"
)

func main() {

	addr := env.GetString("DB_ADDR", "postgres://postgres:postgres@localhost:5432/gopher?sslmode=disable")
	log.Println("Connecting to:", addr)

	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	store := store.NewStore(conn)

	db.Seed(store)
}
