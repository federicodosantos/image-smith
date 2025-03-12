package regex

import (
	"errors"
	"regexp"
)

func Password(pass string) error {
	if len(pass) < 8 {
		return errors.New("Password must be at least 8 characters long")
	}

	if !regexp.MustCompile(`[a-z]`).MatchString(pass) {
		return errors.New("Password must contain at least one lowercase letter")
	}

	if !regexp.MustCompile(`[A-Z]`).MatchString(pass) {
		return errors.New("Password must contain at least one uppercase letter")
	}

	if !regexp.MustCompile(`\d`).MatchString(pass) {
		return errors.New("Password must contain at least one number")
	}

	if !regexp.MustCompile(`[@$!%*?&#]`).MatchString(pass) {
		return errors.New("Password must contain at least one special character")
	}

	return nil
}
