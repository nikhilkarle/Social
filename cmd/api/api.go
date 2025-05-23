package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/nikhilkarle/social/docs"
	"github.com/nikhilkarle/social/internal/store"
	httpSwagger "github.com/swaggo/http-swagger"
)

type application struct{
	config config
	store  store.Storage
	logger *zap.SugaredLogger
}

type config struct{
	addr string
	db dbConfig
	env string
	apiURL string
}

type dbConfig struct{
	addr string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime string
}

func (app *application) mount() http.Handler{
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)


	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router){
		r.Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router){
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}",  func(r chi.Router){
				r.Use(app.postContextMiddleware)

				r.Get("/", app.getPostHandler)
				r.Patch("/", app.updatePostHandler)
				r.Delete("/", app.deletePostHandler)
			})
		})

		r.Route("/users", func(r chi.Router){
			r.Route("/{userID}",  func(r chi.Router){
				r.Use(app.userContextMiddleware)

				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router){
				r.Get("/feed", app.getUserFeedHandler)
			})
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error{
	//Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"


	srv := http.Server{
		Addr: app.config.addr,
		Handler: mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second*10,
		IdleTimeout: time.Minute,
	}

	app.logger.Infow("server has started at %s", app.config.addr)

	return srv.ListenAndServe()

}