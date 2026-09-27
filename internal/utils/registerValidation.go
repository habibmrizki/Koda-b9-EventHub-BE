package utils

import (
	"errors"
	"regexp"
	"strings"

	"github.com/habibmrizki/BE-EventHub/internal/dto"
)

var (
	emailRegex   = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	regexLower   = regexp.MustCompile(`[a-z]`)
	regexUpper   = regexp.MustCompile(`[A-Z]`)
	regexSpecial = regexp.MustCompile(`[!@#$%^&*/><]`)
)

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email tidak boleh kosong")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("format email tidak valid")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password minimal 8 karakter")
	}
	if !regexLower.MatchString(password) {
		return errors.New("password harus mengandung huruf kecil")
	}
	if !regexUpper.MatchString(password) {
		return errors.New("password harus mengandung huruf besar")
	}
	if !regexSpecial.MatchString(password) {
		return errors.New("password harus mengandung karakter spesial (!@#$%^&*/<>)")
	}
	return nil
}

func RegisterValidation(body dto.RegisterRequest) error {
	if strings.TrimSpace(body.FullName) == "" {
		return errors.New("nama lengkap wajib diisi")
	}
	if len(strings.TrimSpace(body.FullName)) < 3 {
		return errors.New("nama lengkap minimal 3 karakter")
	}
	if err := ValidateEmail(body.Email); err != nil {
		return err
	}
	if err := ValidatePassword(body.Password); err != nil {
		return err
	}
	return nil
}

func ResetPasswordValidation(body dto.ResetPasswordRequest) error {
	if err := ValidateEmail(body.Email); err != nil {
		return err
	}
	if body.NewPassword != body.ConfirmPassword {
		return errors.New("konfirmasi password tidak cocok")
	}
	if err := ValidatePassword(body.NewPassword); err != nil {
		return err
	}
	return nil
}
