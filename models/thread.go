package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"log"
	"github.com/gobuffalo/envy"
)

type Thread struct {
	ID            uuid.UUID `json:"id" db:"id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	ImageUrl      string    `json:"image_url" db:"image_url"`
	LastMessageAt time.Time `json:"last_message_at" db:"last_message_at"`
	Name          string    `json:"name" db:"name"`
	New           bool      `json:"new" db:"new"`
	Messages      Messages  `json:"messages" has_many:"messages" order_by:"created_at desc"`
	Owner         User      `json:"user" db:"-" belongs_to:"user"`
	OwnerID       uuid.UUID `json:"owner_id" db:"owner_id"`
	Private       bool      `json:"private" db:"private"`
	Active        bool      `json:"active" db:"active"`
	Members       Users     `json:"members" many_to_many:"thread_members" db:"-"`
}

func (thread Thread) Join(u User) error {
	tx, err := pop.Connect("database")
	if err != nil {
		return errors.WithStack(err)
	}

	thread.Members = append(thread.Members, u)

	tx.Eager().Update(&thread)
	if err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a success message
	// and redirect to the threads index page
	return nil
}

// String is not required by pop and may be deleted
func (t Thread) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Threads is not required by pop and may be deleted
type Threads []Thread

// String is not required by pop and may be deleted
func (t Threads) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Thread) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.ImageUrl, Name: "ImageUrl"},
		&validators.TimeIsPresent{Field: t.LastMessageAt, Name: "LastMessageAt"},
		&validators.StringIsPresent{Field: t.Name, Name: "Name"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Thread) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Thread) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}



func GetThread(id string) Thread {
	tx, err := pop.Connect(envy.Get("GO_ENV", "development"))
	if err != nil {
		log.Panic(err)
	}

	thread := Thread{}

	err = tx.Find(&thread, id)

	if err != nil {
		log.Panic(err)
	}
	return thread
}