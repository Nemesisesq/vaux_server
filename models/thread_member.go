package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type ThreadMember struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	ThreadID  uuid.UUID `json:"thread_id" db:"thread_id"`
	MemberID  uuid.UUID `json:"user_id" db:"user_id"`
	Active    bool      `json:"active" db:"active"`
	Member    User      `belongs_to:"users" db:"-"`
	Thread    Thread    `belongs_to:"threads" db:"-"`
}

// String is not required by pop and may be deleted
func (t ThreadMember) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// ThreadMembers is not required by pop and may be deleted
type ThreadMembers []ThreadMember

// String is not required by pop and may be deleted
func (t ThreadMembers) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *ThreadMember) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *ThreadMember) ValidateCreate(tx *pop.Connection) (errors *validate.Errors, err error) {
	return validate.NewErrors(), nil

}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *ThreadMember) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
