package seeds

import (
	"app/app/model"
	"context"

	"github.com/uptrace/bun"
)

func roleSeed(db *bun.DB) error {
	roles := []model.Role{
		{
			Name:        "Admin",
			Description: "Administrator with full access",
		},
		{
			Name:        "User",
			Description: "Regular user with limited access",
		},
	}

	_, err := db.NewInsert().Model(&roles).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}
