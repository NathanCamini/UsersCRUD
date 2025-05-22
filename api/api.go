package api

import (
	"UsersCRUD/models"
	"UsersCRUD/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

func NewHandler(db map[models.ID]models.User) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	// create a new user
	r.Post("/api/users", handleNewUser(db))

	// r.Get("/api/users", handlePost(db)) return all users slice (arr)
	// r.Get("/api/users/{id}", handlePost(db)) return the user with the id

	// r.Delete("/api/users", handlePost(db)) delete user with the id

	// r.Put("/api/users", handlePost(db)) update the id user with the body of request

	return r
}

func handleNewUser(db map[models.ID]models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type dataResponse struct {
			Id   string      `json:"id"`
			User models.User `json:"user"`
		}

		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			utils.SendJson(w, utils.Response{Error: "invalid body"}, http.StatusUnprocessableEntity)
			return
		}

		if err := models.Validate.Struct(user); err != nil {
			utils.SendJson(w, utils.Response{Error: "please provide FirstName, LastName and bio for the user"}, http.StatusBadRequest)
			return
		}

		uuid, err := uuid.NewUUID()
		if err != nil {
			utils.SendJson(w, utils.Response{Error: "error generating uuid"}, http.StatusInternalServerError)
			return
		}

		db[models.ID(uuid)] = user

		data := dataResponse{Id: uuid.String(), User: user}

		utils.SendJson(w, utils.Response{Data: data}, http.StatusCreated)
	}
}
