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

	// delete the user by ID
	r.Delete("/api/users/{id}", handleDeleteUser(db))

	// update the user
	r.Put("/api/users/{id}", handleUpdateUser(db))

	return r
}

func handleNewUser(db map[models.ID]models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			utils.SendJson(w, utils.Response{Error: "invalid body"}, http.StatusUnprocessableEntity)
			return
		}

		if err := models.Validate.Struct(user); err != nil {
			utils.SendJson(w, utils.Response{Error: "please provide FirstName, LastName and Bio for the user"}, http.StatusBadRequest)
			return
		}

		userCreated, err := users.InsertNewUser(db, user)
		if err != nil {
			utils.SendJson(w, utils.Response{Error: err.Error()}, http.StatusInternalServerError)
			return
		}

		utils.SendJson(w, utils.Response{Data: userCreated}, http.StatusCreated)
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
			utils.SendJson(w, utils.Response{Error: "invalid user ID"}, http.StatusBadRequest)
			return
		}

		user, statusCodeErr, err := users.FindByID(db, parsedID)
		if err != nil {
			utils.SendJson(w, utils.Response{Error: err.Error()}, statusCodeErr)
			return
		}

		utils.SendJson(w, utils.Response{Data: user}, http.StatusOK)
	}
}

func handleUpdateUser(db map[models.ID]models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUser models.User
		id := chi.URLParam(r, "id")

		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			utils.SendJson(w, utils.Response{Error: "invalid body"}, http.StatusUnprocessableEntity)
			return
		}

		if err := models.Validate.Struct(newUser); err != nil {
			utils.SendJson(w, utils.Response{Error: "please provide FirstName, LastName and Bio for the user"}, http.StatusBadRequest)
			return
		}

		parsedID, err := uuid.Parse(id)
		if err != nil {
			utils.SendJson(w, utils.Response{Error: "invalid user ID"}, http.StatusBadRequest)
			return
		}

		_, statusCodeErr, err := users.FindByID(db, parsedID)
		if err != nil {
			utils.SendJson(w, utils.Response{Error: err.Error()}, statusCodeErr)
			return
		}

		users.UpdateUser(db, parsedID, newUser)

		utils.SendJson(w, utils.Response{Data: newUser}, http.StatusOK)
	}
}

func handleDeleteUser(db map[models.ID]models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		parsedID, err := uuid.Parse(id)
		if err != nil {
			utils.SendJson(w, utils.Response{Error: "invalid user ID"}, http.StatusBadRequest)
			return
		}

		statusCodeErr, err := users.DeleteUser(db, parsedID)
		if err != nil {
			utils.SendJson(w, utils.Response{Error: err.Error()}, statusCodeErr)
			return
		}

		utils.SendJson(w, utils.Response{}, http.StatusNoContent)
	}
}
