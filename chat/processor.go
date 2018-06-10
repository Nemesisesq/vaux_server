package chat

import (
	"github.com/gobuffalo/pop"
	"log"

	"github.com/nemesisesq/vaux_server/models"
	"encoding/json"
)

func (c *Client) processData(d Data) {
	switch d.Type {
	case "SET_USER":
		SetUser(d, c)
	case "CREATE_THREAD":
		tx, err := pop.Connect("development")

		if err != nil {
			log.Panic(err)
		}
		tmp, err := json.Marshal(d.Paylaod)
		// Allocate an empty User
		thread := &models.Thread{}

		err = json.Unmarshal(tmp, &thread)

		thread.Owner = c.user
		thread.OwnerID = c.user.ID

		err = tx.Create(thread)
		if err != nil {
			log.Panic(err)
		}
	}
}

func SetUser(d Data, c *Client) {
	tx, err := pop.Connect("development")
	defer tx.Close()
	if err != nil {
		log.Panic(err)
	}
	// Allocate an empty User
	user := &models.User{}
	err = tx.Where("email = ?", d.Paylaod).First(user)
	if err != nil {
		log.Panic(err)
	}
	c.user = *user
	c.broadcastThreads()
}
