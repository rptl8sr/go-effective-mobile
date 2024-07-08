package db

import (
	_ "embed"
	"errors"
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

//go:embed queries/update_user.sql
var updateUser string

func UpdateUser(u *models.User) error {
	saved, err := GetUser(u.ID)
	if err != nil {
		return err
	}

	if u.Name != "" && u.Name != saved.Name {
		saved.Name = u.Name
	}

	if u.Surname != "" && u.Surname != saved.Surname {
		saved.Surname = u.Surname
	}

	if u.Patronymic != "" && u.Patronymic != saved.Patronymic {
		saved.Patronymic = u.Patronymic
	}

	if u.Address != "" && u.Address != saved.Address {
		saved.Address = u.Address
	}

	if u.Passport != "" && u.Passport != saved.Passport {
		saved.Passport = u.Passport
	}

	_, err = client.Pool.Exec(client.Ctx, updateUser,
		u.ID,
		saved.Surname,
		saved.Name,
		saved.Patronymic,
		saved.Address,
		saved.Passport,
	)

	if err != nil {
		return err
	}

	return nil
}

//go:embed queries/delete_user.sql
var deleteUser string

func DeleteUser(id int) (int64, error) {
	res, err := client.Pool.Exec(client.Ctx, deleteUser, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected(), nil
}
