package bootstrap

import (
	log "unknwon.dev/clog/v2"

	"github.com/1076736617/Cardinal/internal/asteroid"
	"github.com/1076736617/Cardinal/internal/conf"
	"github.com/1076736617/Cardinal/internal/dbold"
	"github.com/1076736617/Cardinal/internal/dynamic_config"
	"github.com/1076736617/Cardinal/internal/game"
	"github.com/1076736617/Cardinal/internal/install"
	"github.com/1076736617/Cardinal/internal/livelog"
	"github.com/1076736617/Cardinal/internal/misc"
	"github.com/1076736617/Cardinal/internal/misc/webhook"
	"github.com/1076736617/Cardinal/internal/route"
	"github.com/1076736617/Cardinal/internal/store"
	"github.com/1076736617/Cardinal/internal/timer"
)

func init() {
	// Init log
	_ = log.NewConsole(100)
}

// LinkStart starts the Cardinal.
func LinkStart() {
	// Install
	install.Init()

	// Config
	if err := conf.Init("./conf/Cardinal.toml"); err != nil {
		log.Fatal("Failed to load configuration file: %v", err)
	}

	// Check version
	misc.CheckVersion()

	// Sentry
	misc.Sentry()

	// Init MySQL database.
	dbold.InitMySQL()

	// Check manager
	install.InitManager()

	// Refresh the dynamic config from the database.
	dynamic_config.Init()

	// Check if the database need update.
	misc.CheckDatabaseVersion()

	// Game timer.
	GameToTimerBridge()
	timer.Init()

	// Cache
	store.Init()
	webhook.RefreshWebHookStore()

	// Unity3D Asteroid
	asteroid.Init(game.AsteroidGreetData)

	// Live log
	livelog.Init()

	// Web router.
	router := route.Init()

	log.Fatal("Failed to start web server: %v", router.Run(conf.App.HTTPAddr))
}
