package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type Message struct {
	ID        uuid.UUID    `json:"id" db:"id"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
	Text      nulls.String `json:"text" db:"text"`
	Image     nulls.String `json:"image" db:"image"`
	Sounds    Sounds       `json:"sound" db:"-" many_to_many:"message_sounds"`
	User      User         `json:"user" db:"-" belongs_to:"user"`
	UserID    uuid.UUID    `json:"user_id" db:"user_id"`
	Thread    Thread       `json:"thread" db:"-" belongs_to:"thread"`
	ThreadID  uuid.UUID    `json:"thread_id" db:"thread_id"`
}

// String is not required by pop and may be deleted
func (m Message) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Messages is not required by pop and may be deleted
type Messages []Message

// String is not required by pop and may be deleted
func (m Messages) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *Message) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *Message) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *Message) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
