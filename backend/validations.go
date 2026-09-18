package main

import (
    "errors"
    "net/mail"
    "strings"
)


// Email validation
func validateEmail(email string) error {
    email = strings.TrimSpace(email)
    if email == "" {
        return errors.New("email is required")
    }
    if len(email) > 255 {
        return errors.New("email too long")
    }
    if _, err := mail.ParseAddress(email); err != nil {
        return errors.New("invalid email format")
    }
    return nil
}

// Password validation
func validatePassword(password string) error {
    if password == "" {
        return errors.New("password is required")
    }
    if len(password) < 8 {
        return errors.New("password must be at least 8 characters")
    }
    if len(password) > 72 {
        return errors.New("password must be 72 characters or fewer")
    }
    return nil
}

// Product validation - returns all errors at once
func validateProduct(p Product) []string {
    var errs []string

    name := strings.TrimSpace(p.Name)
    if name == "" {
        errs = append(errs, "name is required")
    } else if len(name) > 255 {
        errs = append(errs, "name must be 255 characters or fewer")
    }

    if len(p.Description) > 5000 {
        errs = append(errs, "description must be 5000 characters or fewer")
    }

    if p.PriceInCents <= 0 {
        errs = append(errs, "price must be greater than 0")
    }
    if p.PriceInCents > 100_000_00 { // $1,000,000 max
        errs = append(errs, "price exceeds maximum allowed")
    }

    if p.Stock < 0 {
        errs = append(errs, "stock cannot be negative")
    }
    if p.Stock > 100_000 {
        errs = append(errs, "stock exceeds maximum allowed")
    }

    image := strings.TrimSpace(p.Image)
    if image == "" {
        errs = append(errs, "image URL is required")
    } else if !strings.HasPrefix(image, "http://") && !strings.HasPrefix(image, "https://") {
        errs = append(errs, "image must be a valid http(s) URL")
    } else if len(image) > 2048 {
        errs = append(errs, "image URL too long")
    }

    return errs
}