// validator package contains the various functions to for validations
package validator

import (
	"net/mail"
	"regexp"
)

const (
	minPasswordLength int = 8
)

// IsValidEmail checks if given email is valid or not
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err != nil
}

// IsValidPassword will check if
//  1. Password contains at-least one upper case letter
//  2. Password contains at-least one lower case letter
//  3. Password contains at-least one digit
//  4. Password contains at-least one special character from "@$!%*#?&^_-"
//  5. Password contains at-least 8 characters
func IsValidPassword(s string) bool {
	uppercaseRegex := regexp.MustCompile(`[A-Z]+`)
	hasUppercase := uppercaseRegex.MatchString(s)

	lowercaseRegex := regexp.MustCompile(`[a-z]+`)
	hasLowercase := lowercaseRegex.MatchString(s)

	digitRegex := regexp.MustCompile(`\d+`)
	hasDigits := digitRegex.MatchString(s)

	specialCharacterRegex := regexp.MustCompile(`[@$!%*#?&^_-]+`)
	hasSpecialCharacter := specialCharacterRegex.MatchString(s)

	hasLengthGTMinPasswordLength := len(s) >= minPasswordLength

	return hasUppercase && hasLowercase && hasDigits && hasSpecialCharacter && hasLengthGTMinPasswordLength
}
