package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type MessageSound struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	MessageID uuid.UUID `json:"message_id" db:"message_id"`
	SoundID   uuid.UUID `json:"sound_id" db:"sound_id"`
	Message   Message   `json:"message" db:"-" belongs_to:"messages"`
	Sound     Sound     `json:"sound" db:"-" belongs_to:"sounds"`
}

// String is not required by pop and may be deleted
func (m MessageSound) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// MessageSounds is not required by pop and may be deleted
type MessageSounds []MessageSound

// String is not required by pop and may be deleted
func (m MessageSounds) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *MessageSound) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *MessageSound) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *MessageSound) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
