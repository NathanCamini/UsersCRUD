package users

import (
	"UsersCRUD/models"
	"fmt"

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

func FindByID(db map[models.ID]models.User, id uuid.UUID) (*UserWithID, error) {
	user, ok := db[models.ID(id)]
	if !ok {
		return nil, fmt.Errorf("The user with the specified ID does not exist")
	}

	return &UserWithID{
		ID:   id.String(),
		User: user,
	}, nil
}
