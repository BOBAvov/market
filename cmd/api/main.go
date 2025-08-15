package main

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"market/internal/config"
	"market/internal/db"
	"market/internal/handler"
	"market/internal/logger"
	"market/internal/middleware"
	"market/internal/repository"
	"market/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}
	zl, err := logger.New(cfg.Logger.Level)
	if err != nil {
		log.Fatalf("logger: %v", err)
	}
	defer func(zl *zap.Logger) {
		err := zl.Sync()
		if err != nil {
			return
		}
	}(zl)
	z := zl.Sugar()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := db.NewPool(ctx, db.Config{
		DSN:             cfg.DB.DSN,
		MaxConns:        cfg.DB.MaxConns,
		MinConns:        cfg.DB.MinConns,
		MaxConnLifetime: cfg.DB.MaxConnLifetime,
		MaxConnIdleTime: cfg.DB.MaxConnIdleTime,
	})
	if err != nil {
		z.Fatalw("db connect failed", "err", err)
	}
	defer pool.Close()

	ctxPing, cancelPing := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelPing()
	if err := pool.Ping(ctxPing); err != nil {
		z.Fatalw("db ping failed", "err", err)
	}
	z.Infow("db connected")

	app := fiber.New(fiber.Config{
		Prefork:      cfg.Server.Prefork,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		AppName:      "Marketplace API",
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			var fe *fiber.Error
			if errors.As(e, &fe) {
				return c.Status(fe.Code).JSON(fiber.Map{"error": fe.Message})
			}
			z.Errorw("unhandled error", "err", e)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		},
	})
	app.Use(recover.New())
	app.Use(flogger.New())
	// ...
	// Repos
	userRepo := repository.NewUserRepository(pool)
	productRepo := repository.NewProductRepository(pool)
	pictureRepo := repository.NewPictureRepository(pool) // NEW

	// Services
	authSvc := service.NewAuthService(userRepo, cfg.Auth.JWTSecret, cfg.Auth.AccessTTL)
	productSvc := service.NewProductService(productRepo)
	pictureSvc := service.NewPictureService(productRepo, pictureRepo) // NEW

	// Handlers
	authH := handler.NewAuthHandler(authSvc)
	prodH := handler.NewProductHandler(productSvc)
	picH := handler.NewPictureHandler(pictureSvc) // NEW

	// Routes
	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register", authH.Register)
	auth.Post("/login", authH.Login)

	products := api.Group("/products")
	products.Get("/", prodH.List)            // public
	products.Get("/:id", prodH.Get)          // public
	products.Get("/:id/pictures", picH.List) // public
	api.Get("/pictures/:id", picH.Download)  // public

	// seller-only
	secured := products.Use(middleware.AuthRequired(middleware.AuthConfig{JWTSecret: cfg.Auth.JWTSecret}))
	secured.Use(middleware.RequireSeller())
	secured.Post("/", prodH.Create)
	secured.Put("/:id", prodH.Update)
	secured.Delete("/:id", prodH.Delete)

	// pictures (seller only)
	secured.Post("/:id/pictures", picH.Upload)
	secured.Delete("/:id/pictures/:pid", picH.Delete)
	secured.Put("/:id/cover/:pid", picH.SetCover)
	// Graceful shutdown
	go func() {
		if err := app.Listen(cfg.Server.Addr); err != nil {
			z.Fatalw("server error", "err", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	z.Infow("shutting down...")
	_ = app.Shutdown()
}
