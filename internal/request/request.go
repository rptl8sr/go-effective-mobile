package request

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go-effective-mobile/internal/models"
	"go-effective-mobile/internal/storage/db"
)

var InfoChan = make(chan int)

type userResponse struct {
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	Patronymic string `json:"patronymic,omitempty"`
}

func GetUserInfo(ctx context.Context, ext string, userId int) (*models.User, error) {
	user, err := db.GetUser(userId)
	if err != nil {
		return nil, err
	}

	pn := strings.Split(user.PassportNumber, " ")
	if len(pn) != 2 {
		return nil, fmt.Errorf("invalid user passport number")
	}

	s := pn[0]
	n := pn[1]

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ext, http.NoBody)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("passportSerie", s)
	query.Add("passportNumber", n)
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API request failed with status code: %d", resp.StatusCode)
	}

	var userInfo userResponse
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	user.Name = userInfo.Name
	user.Surname = userInfo.Surname
	user.Patronymic = userInfo.Patronymic
	user.Address = userInfo.Address

	err = db.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
