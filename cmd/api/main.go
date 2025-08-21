package main

import (
	"log"
	"time"

	"github.com/szoumoc/social/internal/db"
	"github.com/szoumoc/social/internal/env"
	"github.com/szoumoc/social/internal/mailer"
	"github.com/szoumoc/social/internal/store"
	"go.uber.org/zap"
)

const version = " 0.0.1"

//	@title			social API
//	@description	API for the social application, a social network for sharing text posts.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				API key for authorization

func main() {
	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgresql://neondb_owner:npg_hipBck5aqYG7@ep-patient-queen-a1omvmm2-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"),
			maxIdleConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "production"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, //3 days
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", "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"),
			},
		},
	}
	log.Printf("Connecting to DB: %s", cfg.db.addr)

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync() // flushes buffer, if any

	// Database

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Info("db connected")
	store := store.NewStorage(db)

	mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mailer,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
