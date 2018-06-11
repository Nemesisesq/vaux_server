package chat

import (
	"log"

	"github.com/gobuffalo/pop"

	"encoding/json"

	"github.com/nemesisesq/vaux_server/models"
	"github.com/gobuffalo/uuid"
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
	//defer tx.Close()
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

func (c *Client) broadcastThreads() {

	tx, err := pop.Connect("development")
	if err != nil {
		log.Panic(err)
	}
	threads := models.Threads{}

	err = tx.Eager().Where("owner_id  =  ? ", c.user.ID).All(&threads)

	if err != nil {
		log.Panic(err)
	}
	tu := models.User{}
	err = tx.Eager().Where("id = ?", c.user.ID).First(&tu)

	if err != nil {
		log.Panic(err)
	}

	//ids := getIDs(threads)

	threads = append(tu.OwnedThreads, tu.JoinedThreads...)

	data := Data{
		"threads",
		threads,
		nil,
	}
	out, err := json.Marshal(data)

	if err != nil {
		errors := Data{
			"errors",
			"There was a problem marshaling threads",
			nil,
		}
		errOut, _ := json.Marshal(errors)
		c.out <- errOut
	}
	c.out <- out

}

func getIDs(t models.Threads) []uuid.UUID {
	ids := []uuid.UUID{}
	for _, v := range t {
		ids = append(ids, v.ID)
	}
	return ids
}
