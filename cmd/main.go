package cmd

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HoneySinghDev/go-templ-htmx-template/internal/router"
	"github.com/HoneySinghDev/go-templ-htmx-template/pkg/sb"
	"github.com/HoneySinghDev/go-templ-htmx-template/pkg/server"
	"github.com/HoneySinghDev/go-templ-htmx-template/pkl/pklgen/environment"
	"github.com/labstack/echo/v4"

	"github.com/rs/zerolog/log"

	"github.com/HoneySinghDev/go-templ-htmx-template/internal/config"
	"github.com/rs/zerolog"
)

const (
	timeOutValue = 10
)

func App(static embed.FS) {
	c := config.DefaultServiceConfig()

	// Setup logger
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(c.LogLevelFromString(c.Logger.Level))

	if c.Logger.PreetyPrintConsole {
		log.Logger = log.Output(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = "15:04:05"
		}))
	}

	s := server.NewServer(c)

	if err := s.InitDB(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}

	// Init SupaBase Client
	s.SB = sb.InitSB(c.Supabase.ApiUrl, c.Supabase.SecretKey)

	router.Init(s)

	// Serve static files
	assetHandler := http.FileServer(getFileSystem(c.Env == environment.Prod, static))
	s.Echo.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))

	go func() {
		if err := s.Start(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Info().Msg("Server closed")
			} else {
				log.Fatal().Err(err).Msg("Failed to start server")
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOutValue)*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Panic().Err(err).Msg("Failed to gracefully shut down server")
	}
}

func getFileSystem(useOS bool, static embed.FS) http.FileSystem {
	if useOS {
		log.Debug().Msg("using OS mode")
		return http.FS(os.DirFS("static"))
	}

	log.Debug().Msg("using embed mode")
	fsys, err := fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
