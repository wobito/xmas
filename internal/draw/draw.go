package draw

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"

	"gopkg.in/yaml.v3"
)

type Draw struct {
	Users         map[string]*User
	PossibleUsers []*User
}

type DrawResult struct {
	Giver    string
	Receiver string
}

var (
	ErrNoValidDraw   = errors.New("no valid draw found - impossible to satisfy all constraints")
	ErrDuplicateUser = errors.New("duplicate user ID found in configuration")
	ErrInvalidSpouse = errors.New("spouse ID not found in user list")
	ErrSelfSpouse    = errors.New("user cannot be their own spouse")
)

type User struct {
	ID         string
	SpouseID   string
	Selected   bool
	LastYearID string
	Pick       string
}

type Config struct {
	Users []ConfigUser `yaml:"users"`
}

type ConfigUser struct {
	ID       string `yaml:"id"`
	Spouse   string `yaml:"spouse"`
	LastYear string `yaml:"lastYear"`
}

func NewDraw() (*Draw, error) {
	return NewDrawFromFile("config.yaml")
}

func NewDrawFromFile(configPath string) (*Draw, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	users := make(map[string]*User, len(config.Users))
	for _, cu := range config.Users {
		if _, exists := users[cu.ID]; exists {
			return nil, fmt.Errorf("%w: %s", ErrDuplicateUser, cu.ID)
		}
		users[cu.ID] = &User{
			ID:         cu.ID,
			SpouseID:   cu.Spouse,
			LastYearID: cu.LastYear,
		}
	}

	// Validate spouse relationships
	for id, user := range users {
		if user.SpouseID != "" {
			if user.SpouseID == id {
				return nil, fmt.Errorf("%w: %s", ErrSelfSpouse, id)
			}
			if _, exists := users[user.SpouseID]; !exists {
				return nil, fmt.Errorf("%w: %s references %s", ErrInvalidSpouse, id, user.SpouseID)
			}
		}
	}

	return &Draw{
		Users: users,
	}, nil
}

func (d *Draw) StartDraw() ([]DrawResult, error) {
	userList := make([]*User, 0, len(d.Users))
	for _, u := range d.Users {
		userList = append(userList, u)
	}

	if !d.backtrackDraw(userList, 0) {
		return nil, ErrNoValidDraw
	}

	results := make([]DrawResult, 0, len(d.Users))
	for _, u := range d.Users {
		if u.Pick != "" {
			results = append(results, DrawResult{
				Giver:    u.ID,
				Receiver: u.Pick,
			})
		}
	}

	return results, nil
}

func (d *Draw) backtrackDraw(users []*User, index int) bool {
	if index == len(users) {
		return true // All users assigned
	}

	user := users[index]
	d.PrepPossibleUsers(user)

	// Shuffle possible users for randomness
	rand.Shuffle(len(d.PossibleUsers), func(i, j int) {
		d.PossibleUsers[i], d.PossibleUsers[j] = d.PossibleUsers[j], d.PossibleUsers[i]
	})

	for _, possible := range d.PossibleUsers {
		// Try this assignment
		user.Pick = possible.ID
		d.Users[possible.ID].Selected = true

		// Recursively try to assign remaining users
		if d.backtrackDraw(users, index+1) {
			return true
		}

		// Backtrack
		user.Pick = ""
		d.Users[possible.ID].Selected = false
	}

	return false // No valid assignment found
}

func (d *Draw) PrepPossibleUsers(u *User) {
	d.PossibleUsers = make([]*User, 0, len(d.Users))

	for _, possible := range d.Users {
		if !possible.Selected {
			if u.ID != possible.ID && u.SpouseID != possible.ID && possible.ID != u.LastYearID {
				d.PossibleUsers = append(d.PossibleUsers, possible)
			}
		}
	}
}
