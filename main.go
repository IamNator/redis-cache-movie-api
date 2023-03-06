package main

import (
	"github.com/iamnator/movie-api/service"
	"log"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/iamnator/movie-api/adapter/cache"
	"github.com/iamnator/movie-api/adapter/repository"
	"github.com/iamnator/movie-api/env"
	"github.com/iamnator/movie-api/handler/http"
)

func main() {

	if er := env.Init(); er != nil {
		panic(er)
	}

	r := mux.NewRouter()

	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: env.Get().REDIS_URL,
	//})
	//
	//db, err := sql.Open("postgres", env.Get().POSTGRES_URL)
	//if err != nil {
	//	panic(err)
	//}

	//gorm.Op

	redisCache := cache.RedisCache{}
	commentRepo := repository.PgxCommentRepoImpl{}

	srv := service.NewServices(redisCache, commentRepo)

	log.Println("Starting server on port ", env.Get().PORT)

	log.Fatal(http.Run(r, srv))
}
