package utils

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func ParseIDParameter(r *http.Request, parameterName string) (string, error) {
	ID := r.URL.Query().Get(parameterName)

	if len(ID) == 0 {
		return "", errors.New(fmt.Sprintf("'%s' parameter must be supplied.", parameterName))
	}

	_, err := uuid.Parse(ID)
	if err != nil {
		return "", errors.New(fmt.Sprintf("'%s' parameter must be valid UUID.", parameterName))
	}

	return ID, nil
}
