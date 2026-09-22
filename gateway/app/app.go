package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/olahol/melody"
	"github.com/smtdfc/nagare/core/logger"
	config2 "github.com/smtdfc/nagare/gateway/common/config"
	"github.com/smtdfc/nagare/gateway/common/websocket"
)

type App struct {
	fiberApp      *fiber.App
	melody        *melody.Melody
	config        *config2.Config
	wsCoordinator *websocket.Coordinator
	logger        *logger.BaseLogger
}

// @Injectable
func NewApp(
	config *config2.Config,
	wsCoordinator *websocket.Coordinator,
	logger *logger.BaseLogger,
) *App {
	m := melody.New()
	m.Config.MaxMessageSize = 10 * 1024 * 1024 // 10MB
	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	fiberApp.Use(cors.New())
	return &App{
		fiberApp:      fiberApp,
		melody:        m,
		config:        config,
		wsCoordinator: wsCoordinator,
		logger:        logger.With("module", "gateway:app"),
	}
}
