package user

import (
	"context"
	"github.com/sorawaslocked/ap2final_base/pkg/logger"
	"github.com/sorawaslocked/ap2final_user_service/internal/app"
	"github.com/sorawaslocked/ap2final_user_service/internal/config"
)

func main() {
	ctx := context.Background()

	// TODO: load config
	cfg := config.MustLoad()

	// TODO: setup logger
	log := logger.SetupLogger(cfg.Env)

	// TODO: initialize app and run it
	log.Info("initializing application")
	application, err := app.New(ctx, cfg, log)
	if err != nil {
		log.Error("failed to initialize application")

		return
	}

	log.Info("running application")
	application.Run()
}
