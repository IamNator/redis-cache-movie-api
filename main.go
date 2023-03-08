package main

import (
	"log"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/iamnator/movie-api/adapter/cache"
	"github.com/iamnator/movie-api/adapter/repository"
	"github.com/iamnator/movie-api/docs"
	"github.com/iamnator/movie-api/env"
	"github.com/iamnator/movie-api/handler/http"
	"github.com/iamnator/movie-api/service"
	"github.com/iamnator/movie-api/thirdparty/swapi"
)

//	@title			Busha Movie API documentation
//	@version		1.0.0
//	@description	This documents all rest endpoints exposed by this application.

//	@contact.name	Nator Verinumbe
//	@contact.email	natorverinumbe@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:9500
// @BasePath	/
func main() {

	if er := env.Init(); er != nil {
		panic(er)
	}

	//programmatically set swagger info
	docs.SwaggerInfo.Title = "Busha Movie API"
	docs.SwaggerInfo.Description = "This is a sample server for a movie API."
	docs.SwaggerInfo.Version = "1.0"

	docs.SwaggerInfo.Host = env.Get().HOST_MACHINE
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := mux.NewRouter()

	redisCache, err := cache.NewRedisCache(env.Get().REDIS_URL) //
	if err != nil {
		panic(err)
	}
	log.Println("Connected to redis")

	commentRepo, err := repository.NewPgxCommentRepository(env.Get().POSTGRES_URL)
	if err != nil {
		panic(err)
	}
	log.Println("Connected to postgres")

	restyClient := resty.New()
	restyClient.SetTimeout(120 * time.Second)
	restyClient.SetRetryCount(3)
	swapiClient, err := swapi.NewSwapi(restyClient.GetClient())
	if err != nil {
		panic(err)
	}

	srv := service.NewServices(redisCache, commentRepo, swapiClient)

	log.Println("Starting server on port ", env.Get().PORT)

	log.Fatal(http.Run(env.Get().PORT, r, srv))
}
