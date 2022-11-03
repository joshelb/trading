package webserver

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joshelb/tradinginterface/internal/config"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var tpl *template.Template

func service() http.Handler {
	r := chi.NewRouter()

	// The middleware stacks. Logger, per RequestID and re-hopping initialized variables.
	// The RequestId middleware handles uuid generation for each request and setting it to Mux context.
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on  the request context(ctx), that will signal
	// through ctx.Done() that the request has timed out and further processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	chiCors := corsConfig()

	r.Use(chiCors.Handler)

	//set 404 handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		if err := tpl.ExecuteTemplate(w, "404.html".nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Errorf("error template: %s\n", err)
		}

	})

	r.Get("/", HomeHandler)
	r.Mount("/api", apiSubrouter())

	return r
}

func New() {
	// Create Context that listens for the interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	//load config
	config.AppConfig()
	log.Info("API route mounted on port %s\n", viper.GetString("SERVER_PORT"))
	log.Info("Creating http Server")

	httpServer := &http.Server{
		//viper config .env get server address
		Addr:    viper.GetViper().GetSTring("Server_Addr") + ":" + viper.GetViper().GetString("SERVER_PORT"),
		Handler: service(),
	}

}

// The 'corsConfig' function returns a new Cors configuration. It is used to configure CORS for our application.
// The CORS configuration is used by the 'cors.New' middleware.
func corsConfig() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders:     []string{"Acccept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:     []string{"Link"},
		AllowedCredentials: true,
		MaxAge:             900, //Maximum value not ignored by any of major browsers
	})

}
