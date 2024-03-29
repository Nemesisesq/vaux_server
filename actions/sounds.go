package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/nemesisesq/vaux_server/models"
	"github.com/pkg/errors"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Sound)
// DB Table: Plural (sounds)
// Resource: Plural (Sounds)
// Path: Plural (/sounds)
// View Template Folder: Plural (/templates/sounds/)

// SoundsResource is the resource for the Sound model
type SoundsResource struct {
	buffalo.Resource
}

// List gets all Sounds. This function is mapped to the path
// GET /sounds
func (v SoundsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	sounds := &models.Sounds{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Sounds from the DB
	if err := q.All(sounds); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, sounds))
}

// Show gets the data for one Sound. This function is mapped to
// the path GET /sounds/{sound_id}
func (v SoundsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Sound
	sound := &models.Sound{}

	// To find the Sound the parameter sound_id is used.
	if err := tx.Find(sound, c.Param("sound_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, sound))
}

// New renders the form for creating a new Sound.
// This function is mapped to the path GET /sounds/new
func (v SoundsResource) New(c buffalo.Context) error {
	return c.Render(200, r.Auto(c, &models.Sound{}))
}

// Create adds a Sound to the DB. This function is mapped to the
// path POST /sounds
func (v SoundsResource) Create(c buffalo.Context) error {
	// Allocate an empty Sound
	sound := &models.Sound{}

	// Bind sound to the html form elements
	if err := c.Bind(sound); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(sound)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, sound))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Sound was created successfully")

	// and redirect to the sounds index page
	return c.Render(201, r.Auto(c, sound))
}

// Edit renders a edit form for a Sound. This function is
// mapped to the path GET /sounds/{sound_id}/edit
func (v SoundsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Sound
	sound := &models.Sound{}

	if err := tx.Find(sound, c.Param("sound_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, sound))
}

// Update changes a Sound in the DB. This function is mapped to
// the path PUT /sounds/{sound_id}
func (v SoundsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Sound
	sound := &models.Sound{}

	if err := tx.Find(sound, c.Param("sound_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Sound to the html form elements
	if err := c.Bind(sound); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(sound)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, sound))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Sound was updated successfully")

	// and redirect to the sounds index page
	return c.Render(200, r.Auto(c, sound))
}

// Destroy deletes a Sound from the DB. This function is mapped
// to the path DELETE /sounds/{sound_id}
func (v SoundsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Sound
	sound := &models.Sound{}

	// To find the Sound the parameter sound_id is used.
	if err := tx.Find(sound, c.Param("sound_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(sound); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "Sound was destroyed successfully")

	// Redirect to the sounds index page
	return c.Render(200, r.Auto(c, sound))
}
