package http

import (
	"encoding/json"
	"github.com/go-openapi/runtime/middleware"
	"github.com/iamnator/movie-api/docs"
	"github.com/iamnator/movie-api/env"
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

func Run(port string, r *mux.Router, srv service.IServices) error {

	handler := NewHandlers(srv)

	//programmatically set swagger info
	docs.SwaggerInfo.Title = "Busha Movie API"
	docs.SwaggerInfo.Description = "This is a sample server for a movie API."
	docs.SwaggerInfo.Version = "1.0"

	if env.Get().HOST_MACHINE != "" {
		docs.SwaggerInfo.Host = env.Get().HOST_MACHINE
	}

	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// documentation for developers
	opts := middleware.SwaggerUIOpts{SpecURL: "./docs/swagger.json"}
	sh := middleware.SwaggerUI(opts, nil)
	r.Handle("/docs", sh)

	r.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))))

	r.Path("/swagger.yaml").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.yaml")
	})

	r.HandleFunc("/movies", handler.getMoviesHandler).Methods(http.MethodGet)
	r.HandleFunc("/characters/{movie_id}", handler.getMovieCharacterHandler).Methods(http.MethodGet)

	r.HandleFunc("/comments/{movie_id}", handler.addCommentHandler).Methods(http.MethodPost)
	r.HandleFunc("/comments/{movie_id}", handler.getCommentHandler).Methods(http.MethodGet)

	return http.ListenAndServe(":"+port, r)

}

func respondWithError(w http.ResponseWriter, code int, msg string, err ...interface{}) {
	_ = json.NewEncoder(w).Encode(model.GenericResponse{
		Error:   err,
		Data:    nil,
		Code:    code,
		Message: msg,
	})
}

func respondWithSuccess(w http.ResponseWriter, code int, msg string, count int64, data interface{}) {
	_ = json.NewEncoder(w).Encode(model.GenericResponse{
		Error:   nil,
		Code:    code,
		Message: msg,
		Data:    data,
		Count:   count,
	})
}

// getMoviesHandler handles the request to get all movies
// @Summary Get all movies
// @Description Get all movies
// @Tags Movies
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} model.GenericResponse{data=[]model.Movie}
// @Failure 400,502 {object} model.GenericResponse{error=string}
// @Router /movies [get]
func (h handlers) getMoviesHandler(w http.ResponseWriter, r *http.Request) {

	//get page and page size from query params
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		pageSize = 10
	}

	movieList, _, err := h.service.GetMovies(page, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting movies", err)
		return
	}

	_ = json.NewEncoder(w).Encode(movieList)
}

// getMovieCharacterHandler handles the request to get all characters in a movie
// @Summary Get all characters in a movie
// @Description Get all characters in a movie
// @Tags Characters
// @Param movie_id path int true "Movie ID"
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} model.GenericResponse{data=[]model.Character}
// @Failure 400,502 {object} model.GenericResponse{error=string}
// @Router /characters/{movie_id} [get]
func (h handlers) getMovieCharacterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID, err := strconv.Atoi(vars["movie_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid movie id", err)
		return
	}

	//get page and page size from query params
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		pageSize = 10
	}

	characterList, count, err := h.service.GetCharactersByMovieID(movieID, page, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting characters", err)
		return
	}

	respondWithSuccess(w, 200, "Success", count, characterList)
}

// getCommentHandler handles the request to get all comments for a movie
// @Summary Get all comments for a movie
// @Description Get all comments for a movie
// @Tags Comments
// @Param movie_id path int true "Movie ID"
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} model.GenericResponse{data=[]model.Comment}
// @Failure 400,502 {object} model.GenericResponse{error=string}
// @Router /comments/{movie_id} [get]
func (h handlers) getCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID, err := strconv.Atoi(vars["movie"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid movie id", err)
		return
	}

	//get page and page size from query params
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		pageSize = 10
	}

	comments, count, err := h.service.GetComment(movieID, page, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting comments", err)
		return
	}

	respondWithSuccess(w, 200, "Success", count, comments)
}

// addCommentHandler handles the request to add a comment to a movie
// @Summary Add a comment to a movie
// @Description Add a comment to a movie
// @Tags Comments
// @Accept json
// @Param movie_id path int true "Movie ID"
// @Param comment body model.Comment true "Comment"
// @Success 201 {object} model.GenericResponse{data=model.Comment}
// @Failure 400,502 {object} model.GenericResponse{error=string}
// @Router /comments/{movie_id} [post]
func (h handlers) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID, err := strconv.Atoi(vars["movie"])
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	var comment model.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	comment.SwapiMovieID = movieID
	comment.IPv4Addr = r.RemoteAddr
	comment.CreatedAt = time.Now().UTC()

	if err = h.service.SaveComment(movieID, comment); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error saving comment", err)
		return
	}

	respondWithSuccess(w, 201, "Comment added successfully", 0, nil)

}
