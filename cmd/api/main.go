package main

import (
	"log"
	"time"

	"github.com/szoumoc/social/internal/auth"
	"github.com/szoumoc/social/internal/db"
	"github.com/szoumoc/social/internal/env"
	"github.com/szoumoc/social/internal/mailer"
	"github.com/szoumoc/social/internal/ratelimiter"
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
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgresql://neondb_owner:npg_hipBck5aqYG7@ep-patient-queen-a1omvmm2-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"),
			maxIdleConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, //3 days
			fromEmail: env.GetString("FROM_EMAIL", "hello@demomailtrap.co"),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			mailTrap: mailTrapConfig{
				apiKey: env.GetString("MAILTRAP_API_KEY", "8344d666a1d199a460b907383d17e4e4"),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
				exp:    time.Hour * 24 * 3, //3days
				iss:    "szoumosocial",
			},
		},
		rateLimiter: ratelimiter.Config{
			RequestPerTimeFrame: env.GetInt("RATELIMITER_REQUEST_COUNT", 20),
			TimeFrame:           time.Second * 5,
			Enabled:             env.GetBool("RATE_LIMITER_ENABLED", true),
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

	ratelimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	// mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	mailTrap, err := mailer.NewMailTrapClient(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)

	app := &application{
		config:        cfg,
		store:         store,
		logger:        logger,
		mailer:        mailTrap,
		authenticator: jwtAuthenticator,
		rateLimiter:   ratelimiter,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
