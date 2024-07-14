package db

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"strings"

	"go-effective-mobile/internal/models"
)

//go:embed queries/insert_user.sql
var insertUser string

func NewUser(u *models.User) (int, error) {
	var userID int

	err := client.Pool.QueryRow(client.Ctx, insertUser,
		u.Surname, u.Name, u.Patronymic, u.Address, u.PassportNumber).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

//go:embed queries/get_user.sql
var getUser string

func GetUser(id int) (*models.User, error) {
	var user models.User

	row := client.Pool.QueryRow(client.Ctx, getUser, id)
	err := row.Scan(
		&user.ID,
		&user.Surname,
		&user.Name,
		&user.Patronymic,
		&user.Address,
		&user.PassportNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
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

	if u.PassportNumber != "" && u.PassportNumber != saved.PassportNumber {
		saved.PassportNumber = u.PassportNumber
	}

	_, err = client.Pool.Exec(client.Ctx, updateUser,
		u.ID,
		saved.Surname,
		saved.Name,
		saved.Patronymic,
		saved.Address,
		saved.PassportNumber,
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

var DefaultLimit = 10

func GetUsers(ctx context.Context, filter models.UserFilter) ([]*models.User, error) {
	var users []*models.User
	var sb strings.Builder

	sb.WriteString(
		"SELECT id, surname, name, patronymic, address, passport, created_at, updated_at FROM users WHERE true ")

	queryArgs := []interface{}{}

	if filter.Surname != nil {
		sb.WriteString(fmt.Sprintf(" AND LOWER(surname) LIKE LOWER($%d) ", len(queryArgs)+1))
		queryArgs = append(queryArgs, "%"+*filter.Surname+"%")
	}

	if filter.Name != nil {
		sb.WriteString(fmt.Sprintf(" AND LOWER(name) LIKE LOWER($%d) ", len(queryArgs)+1))
		queryArgs = append(queryArgs, "%"+*filter.Name+"%")
	}

	if filter.Patronymic != nil {
		sb.WriteString(fmt.Sprintf(" AND LOWER(patronymic) LIKE LOWER($%d) ", len(queryArgs)+1))
		queryArgs = append(queryArgs, "%"+*filter.Patronymic+"%")
	}

	if filter.Address != nil {
		sb.WriteString(fmt.Sprintf(" AND LOWER(address) LIKE LOWER($%d) ", len(queryArgs)+1))
		queryArgs = append(queryArgs, "%"+*filter.Address+"%")
	}

	if filter.PassportNumber != nil {
		sb.WriteString(fmt.Sprintf(" AND LOWER(passport) LIKE LOWER($%d) ", len(queryArgs)+1))
		queryArgs = append(queryArgs, "%"+*filter.PassportNumber+"%")
	}

	if filter.MinDate != nil {
		sb.WriteString(fmt.Sprintf(" AND created_at >= $%d ", len(queryArgs)+1))
		queryArgs = append(queryArgs, filter.MinDate.Format("2006-01-02"))
	}

	if filter.MaxDate != nil {
		sb.WriteString(fmt.Sprintf(" AND created_at <= $%d ", len(queryArgs)+1))
		queryArgs = append(queryArgs, filter.MaxDate.Format("2006-01-02"))
	}

	limit := filter.Pagination.Limit
	if limit == 0 {
		limit = DefaultLimit
	}

	offset := filter.Pagination.Offset

	sb.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d ", len(queryArgs)+1, len(queryArgs)+2))
	queryArgs = append(queryArgs, limit, offset)

	rows, err := client.Pool.Query(ctx, sb.String(), queryArgs...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.ID,
			&user.Surname,
			&user.Name,
			&user.Patronymic,
			&user.Address,
			&user.PassportNumber,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}
