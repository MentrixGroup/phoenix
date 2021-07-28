package main

import (
	"fmt"
	"strconv"
	"strings"
)

// sum10 returns the weighted sum of the provided ISBN-10 string. It is used
// to calculate the ISBN-10 check digit or to validate an ISBN-10.
func sum10(isbn string) (int, error) {
	s := 0
	w := 10
	for k, v := range isbn {
		if k == 9 && v == 88 {
			// Handle "X" as the digit.
			s += 10
		} else {
			n, err := strconv.Atoi(string(v))
			if err != nil {
				return -1, fmt.Errorf("Failed to convert ISBN-10 character to int: %s", string(v))
			}
			s += n * w
		}
		w--
	}
	return s, nil
}

// sum13 returns the weighted sum of the provided ISBN-13 string. It is used
// to calculate the ISBN-13 check digit or to validate an ISBN-13.
func sum13(isbn string) (int, error) {
	s := 0
	w := 1

	for _, v := range isbn {
		n, err := strconv.Atoi(string(v))
		if err != nil {
			return -1, fmt.Errorf("Failed to convert ISBN-13 character to int: %s", string(v))
		}
		s += n * w
		if w == 1 {
			w = 3
		} else {
			w = 1
		}
	}
	return s, nil
}

// CheckDigit10 returns the check digit for an ISBN-10.
func checkDigit10(isbn10 string) (string, error) {
	if len(isbn10) != 9 && len(isbn10) != 10 {
		return "", fmt.Errorf("A string of length 9 or 10 is required to calculate the ISBN-10 check digit. Provided was: %s", isbn10)
	}
	s, err := sum10(isbn10[:9])
	if err != nil {
		return "", err
	}
	d := (11 - (s % 11)) % 11
	if d == 10 {
		return "X", nil
	}
	return strconv.Itoa(d), nil
}

// CheckDigit13 returns the check digit for an ISBN-13.
func checkDigit13(isbn13 string) (string, error) {
	if len(isbn13) != 12 && len(isbn13) != 13 {
		return "", fmt.Errorf("A string of length 12 or 13 is required to calculate the ISBN-13 check digit. Provided was: %s", isbn13)
	}
	s, err := sum13(isbn13[:12])
	if err != nil {
		return "", err
	}
	d := 10 - (s % 10)
	if d == 10 {
		return "0", nil
	}
	return strconv.Itoa(d), nil
}

// Validate returns true if the provided string is a valid ISBN-10 or ISBN-13.
func validate(isbn string) bool {
	str := strings.ReplaceAll(isbn, "-", "")
	// fmt.Println(strings.ReplaceAll(isbn, "-", ""))
	// fmt.Println(len(strings.ReplaceAll(isbn, "-", "")))
	switch len(str) {
	case 10:
		return validate10(str)
	case 13:
		return validate13(str)
	}
	return false
}

// Validate10 returns true if the provided string is a valid ISBN-10.
func validate10(isbn10 string) bool {
	if len(isbn10) == 10 {
		s, _ := sum10(isbn10)
		return s%11 == 0
	}
	return false
}

// Validate13 returns true if the provided string is a valid ISBN-13.
func validate13(isbn13 string) bool {

	fmt.Println(len(isbn13))
	fmt.Println(isbn13)
	if len(isbn13) == 13 {
		s, _ := sum13(isbn13)
		fmt.Println(s % 10)
		return s%10 == 0
	}
	return false
}

// To13 converts an ISBN-10 to an ISBN-13.
func convertTo13(isbn10 string) string {
	str := strings.ReplaceAll(isbn10, "-", "")

	isbn13 := "978" + str[:9]
	d, err := checkDigit13(isbn13)

	if err != nil {
		return ""
	}

	return isbn13 + d
}

func getSourceId(isbn string) string {
	str := strings.ReplaceAll(isbn, "-", "")

	switch len(str) {
	case 10:
		return fmt.Sprintf("/books/%s.json", convertTo13(str))
	case 13:
		return fmt.Sprintf("/books/%s.json", str)
	}

	return str
}
