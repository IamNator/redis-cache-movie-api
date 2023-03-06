package http

import (
	"github.com/iamnator/movie-api/service"
	"net/http"

	"github.com/gorilla/mux"
)

func Run(r *mux.Router, srv service.IServices) error {

	handler := NewHandlers(srv)

	r.HandleFunc("/movies", handler.getMoviesHandler).Methods("GET")

	r.HandleFunc("/comments/{movie_id}", handler.handleCommentHandler).Methods("GET", "POST")

	r.HandleFunc("/characters/{movie_id}", handler.getMovieCharacterHandler).Methods("GET")

	return http.ListenAndServe(":8000", r)

}
