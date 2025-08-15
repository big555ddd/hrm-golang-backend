package model

import (
	"github.com/uptrace/bun"
)

type UserForgot struct {
	bun.BaseModel `bun:"table:user_forgots"`

	ID      string `bun:",default:gen_random_uuid(),pk"`
	UserID  string `bun:"user_id,notnull"`
	Ref     string `bun:"ref,notnull"`
	Otp     string `bun:"otp,notnull"`
	Expires int64  `bun:"expires,notnull"`
	Used    bool   `bun:"used,default:false"`
	CreateUpdateUnixTimestamp
}
