package draw

import (
	"math/rand"

	"github.com/olekukonko/tablewriter"
)

type Draw struct {
	Users         map[string]*User
	PossibleUsers []*User
}

type User struct {
	ID         string
	SpouseID   string
	Drawn      bool
	Selected   bool
	LastYearID string
	Pick       string
}

func NewDraw() *Draw {
	return &Draw{}
}

func (d *Draw) StartDraw(table *tablewriter.Table) {
	for _, u := range d.Users {
		d.PrepPossibleUsers(u)
		if len(d.PossibleUsers) > 0 {
			d.doDraw(u)
			table.Append([]string{u.ID, u.Pick})
		}
	}
}

func (d *Draw) doDraw(u *User) {
	var picked *User

	if len(d.PossibleUsers) > 1 {
		pickedIdx := rand.Intn(len(d.PossibleUsers) - 1)
		picked = d.PossibleUsers[pickedIdx]
	} else {
		picked = d.PossibleUsers[0]
	}

	u.Pick = picked.ID
	d.Users[picked.ID].Selected = true
}

func (d *Draw) SetUsers() {
	d.Users = map[string]*User{
		"kirsten": {"kirsten", "adrian", false, false, "rudi", ""},
		"adrian":  {"adrian", "kirsten", false, false, "kelita", ""},
		"rudi":    {"rudi", "liane", false, false, "bill", ""},
		"liane":   {"liane", "rudi", false, false, "adrian", ""},
		"kelita":  {"kelita", "bill", false, false, "liane", ""},
		"bill":    {"bill", "kelita", false, false, "kirsten", ""},
	}
}

func (d *Draw) PrepPossibleUsers(u *User) {
	d.PossibleUsers = []*User{}

	for _, possible := range d.Users {
		if !possible.Selected {
			if u.ID != possible.ID && u.SpouseID != possible.ID && possible.ID != u.LastYearID {
				d.PossibleUsers = append(d.PossibleUsers, possible)
			}
		}
	}
}
