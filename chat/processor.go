package chat

import (
	"log"

	"github.com/gobuffalo/pop"

	"encoding/json"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/uuid"
	"github.com/nemesisesq/vaux_server/models"
	"hash/fnv"
	"sort"
	"github.com/mitchellh/hashstructure"
	"strconv"
)

func (c *Client) processData(d Data) {
	switch d.Type {
	case "SET_USER":
		SetUser(d, c)
	case "CREATE_THREAD":
		createThread(d, c)
	case "ADD_MESSAGE":
		message, err := json.Marshal(d)
		if err != nil {
			log.Panic(err)
		}
		c.in <- message
	}
}

func createThread(d Data, c *Client) {
	tx, err := pop.Connect(envy.Get("GO_ENV", "development"))
	defer tx.Close()
	if err != nil {
		log.Panic(err)
	}
	tmp, err := json.Marshal(d.Paylaod)
	if err != nil {
		log.Panic(err)
	}
	// Allocate an empty User
	thread := &models.Thread{}
	err = json.Unmarshal(tmp, &thread)
	if err != nil {
		log.Panic(err)
	}
	thread.OwnerID = thread.Owner.ID

	uqk, err := hashMembers(thread)
	if err != nil {
		panic(err)
	}

	thread.UniqueKey = strconv.Itoa(int(uqk))
	verrs, err := tx.ValidateAndCreate(thread)

	if err != nil {
		log.Panic(err)
	}

	if verrs.HasAny() {

		data := Data{
			"validation_errors",
			verrs,
			nil,
		}
		out, err := json.Marshal(data)

		if err != nil {
			panic(err)
		}

		c.out <- out

		return

	}

	for _, v := range thread.Members {
		threadMember := &models.ThreadMember{ThreadID: thread.ID, MemberID: v.ID}

		verrs, err := tx.ValidateAndCreate(threadMember)
		if err != nil {
			log.Panic(err)
		}

		if verrs.HasAny() {
			// Make the errors available inside the html template

			// Render again the new.html template that the user can
			// correct the input.

			data := Data{
				"validation_errors",
				verrs,
				nil,
			}
			out, err := json.Marshal(data)

			if err != nil {
				errors := Data{
					"errors",
					"There was a problem marshaling users",
					nil,
				}
				errOut, _ := json.Marshal(errors)
				c.out <- errOut
			}

			c.out <- out
			return
		}
	}

	c.broadcastThreads()
}

func hashMembers(thread *models.Thread) (uint64, error) {
	hashes := []int{}
	for _, v := range thread.Members {
		hashes = append(hashes, int(hash(v.ID.String())))
	}
	sort.Ints(hashes)
	superKey, err := hashstructure.Hash(hashes, nil)
	return superKey, err
}

func SetUser(d Data, c *Client) {
	tx, err := pop.Connect(envy.Get("GO_ENV", "development"))
	defer tx.Close()
	if err != nil {
		log.Panic(err)
	}
	// Allocate an empty User
	user := &models.User{}
	err = tx.Where("email = ?", d.Paylaod).First(user)
	if err != nil {
		data := Data{
			"bad_user",
			"There is no user please sign back in.",
			nil,
		}

		out, _ := json.Marshal(data)

		c.out <- out
	}
	c.user = *user
	c.broadcastThreads()
	c.broadcastUsers()
	//go c.Subscribe()
	//go c.Publish()
}
func (c *Client) broadcastUsers() {

	users, err := GetAllUsers(c)

	data := Data{
		"users",
		users,
		nil,
	}
	out, err := json.Marshal(data)

	if err != nil {
		errors := Data{
			"errors",
			"There was a problem marshaling users",
			nil,
		}
		errOut, _ := json.Marshal(errors)
		c.out <- errOut
	}
	c.out <- out

}

func (c *Client) broadcastThreads() {

	threads, err := GetAllThreads(c)

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
func GetAllThreads(c *Client) (models.Threads, error) {
	tx, err := pop.Connect(envy.Get("GO_ENV", "development"))
	defer tx.Close()
	if err != nil {
		log.Panic(err)
	}
	threads := models.Threads{}
	tx.Eager().All(&threads)

	return threads, nil
}
func GetAllUsers(c *Client) (models.Users, error) {
	tx, err := pop.Connect(envy.Get("GO_ENV", "development"))
	defer tx.Close()
	if err != nil {
		log.Panic(err)
	}
	users := models.Users{}
	tx.Eager().All(&users)

	return users, nil
}

func GetThreads(c *Client) (models.Threads, error) {

	tx, err := pop.Connect(envy.Get("GO_ENV", "development"))
	defer tx.Close()
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
	return threads, err
}

func getIDs(t models.Threads) []uuid.UUID {
	ids := []uuid.UUID{}
	for _, v := range t {
		ids = append(ids, v.ID)
	}
	return ids
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
