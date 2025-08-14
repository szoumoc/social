package main

import (
	"log"

	"github.com/szoumoc/social/internal/db"
	"github.com/szoumoc/social/internal/env"
	"github.com/szoumoc/social/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgresql://neondb_owner:npg_hipBck5aqYG7@ep-patient-queen-a1omvmm2-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"),
			maxIdleConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Printf("db connected")
	store := store.NewStorage(db)
	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
