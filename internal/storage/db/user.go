package db

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go-effective-mobile/internal/models"
)

//go:embed queries/insert_user.sql
var insertUser string

func NewUser(u *models.User) (int, error) {
	var userID int

	err := client.Pool.QueryRow(client.Ctx, insertUser,
		u.Surname, u.Name, u.Patronymic, u.Address, u.Passport).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

//go:embed queries/get_user.sql
var getUser string

func GetUser(id int) (*models.User, error) {
	var user models.User

	err := client.Pool.QueryRow(client.Ctx, getUser, id).Scan(
		&user.Surname,
		&user.Name,
		&user.Patronymic,
		&user.Address,
		&user.Passport,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound(id)
		}
		return nil, err
	}

	return &user, nil
}

func ErrUserNotFound(id int) error {
	return fmt.Errorf("user %d not found", id)
}
