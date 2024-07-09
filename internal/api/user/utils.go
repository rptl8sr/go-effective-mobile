package user

import (
	"fmt"
	"regexp"
	"strings"
)

func requiredFieldsNewUser(u *userBody) []string {
	var emptyFields []string

	rfs := map[string]string{
		"PassportNumber": *u.PassportNumber,
	}

	for fieldName, fieldValue := range rfs {
		if fieldValue == "" {
			emptyFields = append(emptyFields, fieldName)
		}
	}

	return emptyFields
}

func validatePassport(s string) (string, error) {
	re := regexp.MustCompile(`^(\d{4})\s?(\d{6})$`)
	matches := re.FindStringSubmatch(s)

	if len(matches) != 3 {
		return "", fmt.Errorf("invalid passport number")
	}

	pn := strings.Join(matches[1:], " ")

	return pn, nil
}

func safeDereference(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
