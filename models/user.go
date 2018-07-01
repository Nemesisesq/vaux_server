package models

import (
	"encoding/json"
	"time"

	"strings"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserForm struct {
	User
	Password             string `json:"password" db:"-"`
	PasswordConfirmation string `json:"password_confirmation" db:"-"`
}

type User struct {
	ID            uuid.UUID    `json:"id" db:"id"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at" db:"updated_at"`
	Name          string       `json:"name" db:"name"`
	Email         string       `json:"email" db:"email"`
	PasswordHash  string       `json:"-" db:"password_hash"`
	RefreshToken  nulls.String `json:"-" db:"refresh_token"`
	Admin         bool         `json:"is_admin" db:"is_admin"`
	Avatar        string       `json:"avatar" db:"avatar"`
	OwnedThreads  Threads      `json:"threads" has_many:"threads" fk_id:"owner_id" order_by:"updated_at asc"`
	JoinedThreads Threads      `json:"joined_threads" many_to_many:"thread_members" db:"-"`
	Profile       Profile      `json:"profile,omitempty" has_one:"profile"`
	ProfileID     uuid.UUID    `db:"profile_id" json:"-"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		//&validators.StringIsPresent{Field: u.Name, Name: "Name"},
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
		//&validators.StringIsPresent{Field: u.Avatar, Name: "Avatar"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateAndCrete wraps up the pattern of encrypting the password and
// running validations. Useful when writing tests.
func (u *UserForm) ValidateAndCreate(tx *pop.Connection) (*validate.Errors, error) {
	verrs := validate.NewErrors()

	u.Email = strings.ToLower(u.Email)
	ph, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return verrs, errors.WithStack(err)
	}
	u.PasswordHash = string(ph)

	verrs.Append(validate.Validate(
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&validators.StringsMatch{Name: "Password", Field: u.Password, Field2: u.PasswordConfirmation, Message: "Password does not match confirmation"},
	))
	uv, err := u.User.Validate(tx)
	verrs.Append(uv)

	if err != nil || verrs.HasAny() {
		return verrs, errors.WithStack(err)
	}
	err = tx.Create(&u.User)
	if err != nil {
		return verrs, errors.WithStack(err)
	}
	u.Profile.UserID = u.User.ID
	pv, err := u.Profile.Validate(tx)
	verrs.Append(pv)
	if err != nil || verrs.HasAny() {
		return verrs, errors.WithStack(err)
	}
	err = tx.Eager().Create(&u.User.Profile)
	if err != nil {
		return verrs, err
	}
	u.User.ProfileID = u.Profile.ID
	return verrs, tx.Update(&u.User)
}
