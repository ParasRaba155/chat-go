// hasher package handles all the hash related operations
package hasher

import "golang.org/x/crypto/bcrypt"

// HashAndSalt will hash and salt the password with bcrypt package's default cost
func HashAndSalt(pwd string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPass), nil
}

// MatchHashedPassword will match the hashedPass against the plain text pass
func MatchHashedPassword(hashedPass, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(pass))
}
