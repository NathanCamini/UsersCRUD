package api

import (
	"UsersCRUD/models"
	"UsersCRUD/utils"
	"UsersCRUD/utils/users"
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

	// get all users
	r.Get("/api/users", handleGetUsers(db))

	// get specified user by ID
	r.Get("/api/users/{id}", handleGetUserByID(db))

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
			utils.SendJson(w, utils.Response{Error: "Invalid body"}, http.StatusUnprocessableEntity)
			return
		}

		if err := models.Validate.Struct(user); err != nil {
			utils.SendJson(w, utils.Response{Error: "Please provide FirstName, LastName and Bio for the user"}, http.StatusBadRequest)
			return
		}

		uuid, err := uuid.NewUUID()
		if err != nil {
			utils.SendJson(w, utils.Response{Error: "There was an error while saving the user to the database"}, http.StatusInternalServerError)
			return
		}

		db[models.ID(uuid)] = user

		data := dataResponse{Id: uuid.String(), User: user}

		utils.SendJson(w, utils.Response{Data: data}, http.StatusCreated)
	}
}

func handleGetUsers(db map[models.ID]models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := users.FindAll(db)
		utils.SendJson(w, utils.Response{Data: data}, http.StatusOK)
	}
}

func handleGetUserByID(db map[models.ID]models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		parsedID, err := uuid.Parse(id)
		if err != nil {
			utils.SendJson(w, utils.Response{Error: "The user information could not be retrieved"}, http.StatusInternalServerError)
			return
		}

		user, err := users.FindByID(db, parsedID)
		if err != nil {
			utils.SendJson(w, utils.Response{Error: "The user with the specified ID does not exist"}, http.StatusNotFound)
			return
		}

		utils.SendJson(w, utils.Response{Data: user}, http.StatusOK)
	}
}
