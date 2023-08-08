package queries

import "github.com/jackc/pgx/v5/pgtype"

// RegisterUserRequest is the user form struct
type RegisterUserRequest struct {
	Intro    *string `json:"intro"`
	Profile  *string `json:"profile"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	FullName string  `json:"full_name"`
}

// ToRegisterUserParams handles the data transformation to ToRegisterUserParams
func (req RegisterUserRequest) ToRegisterUserParams(hashedPass string) RegisterUserParams {
	var q RegisterUserParams

	q.Email = req.Email
	q.FullName = req.FullName
	q.PasswordHash = hashedPass

	if req.Profile != nil {
		q.Profile = pgtype.Text{
			String: *req.Profile,
			Valid:  true,
		}
	}

	if req.Intro != nil {
		q.Intro = pgtype.Text{
			String: *req.Intro,
			Valid:  true,
		}
	}

	q.Role = RoleUser

	return q
}

type UserWithToken struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Token    string `json:"token"`
}

func (u User) ToUserWithToken(token string) UserWithToken {
	return UserWithToken{
		Email:    u.Email,
		FullName: u.FullName,
		Token:    token,
	}
}
