package http

import (
	"encoding/json"
	"errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/iamnator/movie-api/service"
	"github.com/rs/zerolog/log"
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

	// documentation for developers
	opts := middleware.SwaggerUIOpts{SpecURL: "./docs/swagger.json"}
	sh := middleware.SwaggerUI(opts, nil)
	r.Handle("/docs", sh)

	r.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))))

	r.Path("/swagger.yaml").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.yaml")
	})

	//add logging middleware
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			elapsed := time.Since(start)
			log.Info().Msgf("%s %s in %s", r.Method, r.RequestURI, elapsed)
		})
	}

	r.Use(loggingMiddleware)

	//add health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.HandleFunc("/movies", handler.getMoviesHandler).Methods(http.MethodGet)
	r.HandleFunc("/movies/{movie_id}", handler.getMovieHandler).Methods(http.MethodGet)
	r.HandleFunc("/characters/{movie_id}", handler.getMovieCharacterHandler).Methods(http.MethodGet)

	r.HandleFunc("/comments/{movie_id}", handler.addCommentHandler).Methods(http.MethodPost)
	r.HandleFunc("/comments/{movie_id}", handler.getCommentHandler).Methods(http.MethodGet)

	return http.ListenAndServe(":"+port, r)

}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err == nil {
		err = errors.New(msg)
	}
	_ = json.NewEncoder(w).Encode(model.GenericResponse{
		Error:   err.Error(),
		Data:    nil,
		Code:    code,
		Message: msg,
	})
}

func respondWithSuccess(w http.ResponseWriter, code int, msg string, count int64, data interface{}) {
	_ = json.NewEncoder(w).Encode(model.GenericResponse{
		Error:   "",
		Code:    code,
		Message: msg,
		Data:    data,
		Count:   count,
	})
}

// getMoviesHandler handles the request to get all movies
//
//	@Summary		Get all movies
//	@Description	Get all movies
//	@Tags			Movies
//	@Param			page		query		int	false	"Page number"
//	@Param			pageSize	query		int	false	"Page size"
//	@Success		200			{object}	model.GenericResponse{data=[]model.Movie, count=int64, message=string}
//	@Failure		400,502		{object}	model.GenericResponse{error=string, message=string}
//	@Router			/movies [get]
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

	movieList, count, err := h.service.GetMovies(page, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting movies", err)
		return
	}

	respondWithSuccess(w, http.StatusOK, "Success", count, movieList)
}

// getMoviesHandler handles the request to get a movie
//
//	@Summary		Get a movie
//	@Description	Get a movie
//	@Tags			Movies
//	@Param			movie_id	path		int	true	"Movie ID"
//	@Success		200			{object}	model.GenericResponse{data=model.Movie}
//	@Failure		400,502		{object}	model.GenericResponse{error=string}
//	@Router			/movies/{movie_id} [get]
func (h handlers) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID, err := strconv.Atoi(vars["movie_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	movie, err := h.service.GetMovieByID(movieID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting movie", err)
		return //
	}

	respondWithSuccess(w, http.StatusOK, "Success", 1, movie)
}

// getMovieCharacterHandler handles the request to get all characters in a movie
//
//	@Summary		Get all characters in a movie
//	@Description	Get all characters in a movie
//	@Tags			Characters
//	@Param			movie_id	path		int		true	"Movie ID"
//	@Param			page		query		int		false	"Page number"
//	@Param			pageSize	query		int		false	"Page size"
//	@Param			sortKey		query		string	false	"Sort key (name | gender | height)"
//	@Param			sortOrder	query		string	false	"Sort order (asc | desc)"
//	@Param			gender		query		string	false	" ' Gender (female' | 'male')"
//	@Success		200			{object}	model.GenericResponse{data=model.CharacterList}
//	@Failure		400,502		{object}	model.GenericResponse{error=string}
//	@Router			/characters/{movie_id} [get]
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

	sortKey := r.URL.Query().Get("sortKey")

	if sortKey != "" {
		switch sortKey {
		case "name", "height":

		default:
			respondWithError(w, http.StatusBadRequest, "Invalid sort key e.g ['name']", err)
			return
		}
	}

	sortOrder := r.URL.Query().Get("sortOrder")
	if sortOrder != "" {
		switch sortOrder {
		case "asc", "desc":

		default:
			respondWithError(w, http.StatusBadRequest, "Invalid sort order e.g ['asc', 'desc']", nil)
			return
		}
	}

	gender := r.URL.Query().Get("gender")
	if gender != "" {
		switch gender {
		case "male", "female": //yeah, I know, but it's just a demo

		default:
			respondWithError(w, http.StatusBadRequest, "Invalid gender e.g ['female', 'male']", nil)
			return

		}
	}

	var arg = model.GetCharactersByMovieIDArgs{
		MovieID:   movieID,
		Page:      page,
		PageSize:  pageSize,
		SortKey:   sortKey,
		SortOrder: sortOrder,
		Gender:    gender,
	}

	characterList, count, err := h.service.GetCharactersByMovieID(arg)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting characters", err)
		return
	}

	respondWithSuccess(w, 200, "Success", count, characterList)
}

// getCommentHandler handles the request to get all comments for a movie
//
//	@Summary		Get all comments for a movie
//	@Description	Get all comments for a movie
//	@Tags			Comments
//	@Param			movie_id	path		int	true	"Movie ID"
//	@Param			page		query		int	false	"Page number"
//	@Param			pageSize	query		int	false	"Page size"
//	@Success		200			{object}	model.GenericResponse{data=[]model.Comment, count=int64, message=string}
//	@Failure		400,502		{object}	model.GenericResponse{error=string}
//	@Router			/comments/{movie_id} [get]
func (h handlers) getCommentHandler(w http.ResponseWriter, r *http.Request) {
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

	comments, count, err := h.service.GetComment(movieID, page, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting comments", err)
		return
	}

	respondWithSuccess(w, 200, "Success", count, comments)
}

// addCommentHandler handles the request to add a comment to a movie
//
//	@Summary		Add a comment to a movie
//	@Description	Add a comment to a movie
//	@Tags			Comments
//	@Accept			json
//	@Param			movie_id	path		int						true	"Movie ID"
//	@Param			comment		body		model.AddCommentRequest	true	"Comment"
//	@Success		201			{object}	model.GenericResponse{message=string}
//	@Failure		400,502		{object}	model.GenericResponse{error=string}
//	@Router			/comments/{movie_id} [post]
func (h handlers) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieID, err := strconv.Atoi(vars["movie_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid movie id", err)
		return
	}

	var req model.AddCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if err := req.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	comment := req.ToComment()

	comment.SwapiMovieID = movieID
	comment.IPv4Addr = r.RemoteAddr
	comment.CreatedAt = time.Now().UTC()

	if err = h.service.SaveComment(movieID, comment); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error saving comment", err)
		return
	}

	respondWithSuccess(w, 201, "Comment added successfully", 0, nil)

}
