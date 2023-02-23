package cmd

import (
	"errors"
	"strconv"
)

func parseInt(input string) (float64, error) {
	i, err := strconv.Atoi(input)
	if err != nil {
		return 0, errors.New("please provide an integer")
	}

	return float64(i), nil
}

func parseFloat(input string) (float64, error) {
	f, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, errors.New("please provide a decimal number")
	}

	return f, nil
}
