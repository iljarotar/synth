package utils

import (
	"errors"
	"strconv"
)

func ParseInt(input string) (int, error) {
	i, err := strconv.Atoi(input)
	if err != nil {
		return 0, errors.New("please provide an integer")
	}

	return int(i), nil
}

func ParseFloat(input string) (float64, error) {
	f, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, errors.New("please provide a decimal number")
	}

	return f, nil
}
