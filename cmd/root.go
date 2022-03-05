package cmd

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/LucasRoesler/openfaas-loki/pkg"
	"github.com/LucasRoesler/openfaas-loki/pkg/faas"
	"github.com/LucasRoesler/openfaas-loki/pkg/handlers"
	"github.com/LucasRoesler/openfaas-loki/pkg/loki"

	"github.com/openfaas/faas-provider/logs"
)

//nolint:gochecknoinits // cobra is initialized in init()
func init() {
	rootCmd.Flags().String("log-level", "INFO", "Logging level")
	rootCmd.Flags().Int("port", 9191, "address the HTTP server will be listening to")
	rootCmd.Flags().Duration("timeout", 30*time.Second, "log request timeout")
	rootCmd.Flags().String("url", "", "base url of the Loki API")

	viper.SetEnvPrefix("OF_LOKI")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	_ = viper.BindPFlags(rootCmd.Flags())
}

var rootCmd = &cobra.Command{
	Use:     "openfaas-loki",
	Short:   "openfaas-loki is a log provider for openfaas, powered by loki",
	Version: pkg.Version,
	Run: func(cmd *cobra.Command, args []string) {
		configureLogging()
		log.WithFields(log.Fields(viper.AllSettings())).Debug("configuration")

		client := loki.New(viper.GetString("url"))
		requester := faas.New(client)

		routes := chi.NewRouter()
		routes.Use(middleware.Recoverer)
		routes.Use(middleware.Heartbeat("/-/health"))
		routes.Get("/-/config", handlers.ConfigHandlerFunc)
		routes.Get("/system/logs", logs.NewLogHandlerFunc(requester, viper.GetDuration("timeout")))

		srv := http.Server{
			Addr:    ":" + viper.GetString("port"),
			Handler: routes,
		}

		idleConnsClosed := make(chan struct{})
		go func() {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			<-sigint

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// We received an interrupt signal, shut down.
			err := srv.Shutdown(ctx)
			if err != nil {
				// Error from closing listeners, or context timeout:
				log.Printf("server Shutdown: %v\n", err)
			}
			close(idleConnsClosed)
		}()

		log.Printf("starting server at %v\n", srv.Addr)
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			// Error starting or closing listener:
			log.Error(err)
		}

		<-idleConnsClosed
		log.Print("server Stopped")
	},
}

func configureLogging() {
	lvl, err := log.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(lvl)
}

// Execute starts the openfaas-loki server
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
