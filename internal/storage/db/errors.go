package db

import "fmt"

func ErrUserNotFound(id int) error {
	return fmt.Errorf("user %d not found", id)
}
