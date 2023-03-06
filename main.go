package main

import (
	"github.com/iamnator/movie-api/service"
	"github.com/iamnator/movie-api/thirdparty/swapi"
	"log"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/iamnator/movie-api/adapter/cache"
	"github.com/iamnator/movie-api/adapter/repository"
	"github.com/iamnator/movie-api/env"
	"github.com/iamnator/movie-api/handler/http"
)

// @title Busha Movie API documentation
// @version 1.0.0
// @description This documents all rest endpoints exposed by this application.

// @contact.name Busha Support
// @contact.email natorverinumbe@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8500
// @BasePath /api
func main() {

	if er := env.Init(); er != nil {
		panic(er)
	}

	r := mux.NewRouter()

	redisCache, err := cache.NewRedisCache("redis://busha_redis:6380")
	if err != nil {
		panic(err)
	}
	log.Println("Connected to redis")

	commentRepo, err := repository.NewPgxCommentRepository(env.Get().POSTGRES_URL)
	if err != nil {
		panic(err)
	}
	log.Println("Connected to postgres")

	swapiClient, err := swapi.NewSwapi()
	if err != nil {
		panic(err)
	}

	srv := service.NewServices(redisCache, commentRepo, swapiClient)

	log.Println("Starting server on port ", env.Get().PORT)

	log.Fatal(http.Run(env.Get().PORT, r, srv))
}
