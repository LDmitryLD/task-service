package models

import userpb "github.com/LDmitryLD/protos2/gen/user"

type User struct {
	FirstName string         `db:"first_name"`
	LastName  string         `db:"last_name"`
	Email     string         `db:"email"`
	Tasks     []*userpb.Task `json:"tasks"`
}
