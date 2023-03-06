package http

import (
	"encoding/json"
	"github.com/iamnator/movie-api/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/iamnator/movie-api/model"
)

type handlers struct {
	service service.IServices
}

func NewHandlers(srv service.IServices) handlers {
	return handlers{
		service: srv,
	}
}

func (h handlers) getMoviesHandler(w http.ResponseWriter, r *http.Request) {

	movieList, err := h.service.GetMovies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(movieList)
}

func (h handlers) handleCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID, err := strconv.Atoi(vars["movie"])
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	if r.Method == "POST" {
		var comment model.Comment
		err := json.NewDecoder(r.Body).Decode(&comment)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		comment.MovieID = movieID
		comment.IPAddress = r.RemoteAddr
		comment.CreatedAt = time.Now().UTC()

		if err = h.service.SaveComment(movieID, comment); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	} else if r.Method == "GET" {

		comments, err := h.service.GetComment(movieID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(comments)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (h handlers) getMovieCharacterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID, err := strconv.Atoi(vars["movie_id"])
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	characterList, err := h.service.GetCharactersByMovieID(movieID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(characterList)
}
