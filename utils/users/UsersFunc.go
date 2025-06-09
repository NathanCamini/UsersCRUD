package users

import (
	"UsersCRUD/models"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type UserWithID struct {
	ID   string      `json:"id"`
	User models.User `json:"user"`
}

func FindAll(db map[models.ID]models.User) []UserWithID {
	var data []UserWithID

	for i, v := range db {
		data = append(data, UserWithID{ID: uuid.UUID(i).String(), User: v})
	}

	return data
}

func FindByID(db map[models.ID]models.User, id uuid.UUID) (*UserWithID, int, error) {
	user, ok := db[models.ID(id)]
	if !ok {
		return nil, http.StatusNotFound, fmt.Errorf("user not found")
	}

	return &UserWithID{
		ID:   id.String(),
		User: user,
	}, http.StatusOK, nil
}

func InsertNewUser(db map[models.ID]models.User, newUser models.User) (*UserWithID, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	db[models.ID(uuid)] = newUser

	return &UserWithID{
		ID:   uuid.String(),
		User: newUser,
	}, nil
}

func UpdateUser(db map[models.ID]models.User, id uuid.UUID, updatedUser models.User) *UserWithID {

	db[models.ID(id)] = updatedUser

	return &UserWithID{
		ID:   id.String(),
		User: updatedUser,
	}
}

func DeleteUser(db map[models.ID]models.User, id uuid.UUID) (int, error) {

	_, ok := db[models.ID(id)]
	if !ok {
		return http.StatusNotFound, fmt.Errorf("user not found")
	}

	delete(db, models.ID(id))

	return http.StatusNoContent, nil
}
