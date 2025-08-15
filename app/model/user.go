package model

import (
	"context"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID        string  `bun:",default:gen_random_uuid(),pk"`
	FirstName string  `bun:"first_name,notnull"`
	LastName  string  `bun:"last_name,notnull"`
	Email     string  `bun:"email,unique,notnull"`
	Password  string  `bun:"password,notnull"`
	EmpCode   string  `bun:"emp_code,unique,notnull"`
	Phone     string  `bun:"phone,notnull"`
	Salary    float64 `bun:"salary,notnull"`
	IsActive  bool    `bun:"is_active,notnull"`

	CreateUpdateUnixTimestamp
	SoftDelete
}

func (u *User) Exist(ctx context.Context, db *bun.DB, id string) (bool, error) {
	ex, err := db.NewSelect().Model(u).Where("id = ?", id).Exists(ctx)

	return ex, err
}
