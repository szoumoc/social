package main

import (
	"log"

	"github.com/szoumoc/social/internal/db"
	"github.com/szoumoc/social/internal/env"
	"github.com/szoumoc/social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgresql://neondb_owner:npg_hipBck5aqYG7@ep-patient-queen-a1omvmm2-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require")
	conn, err := db.New(addr, 30, 30, "15m")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close()
	store := store.NewStorage(conn) // Assuming db is your *sql.DB instance
	db.Seed(store)
}
