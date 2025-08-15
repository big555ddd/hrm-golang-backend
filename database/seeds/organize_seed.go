package seeds

import (
	"app/app/model"
	"context"

	"github.com/uptrace/bun"
)

func organizeSeed(db *bun.DB) error {
	organize := model.Organization{

		Name:        "Default Organization",
		Description: "This is the default organization created during seed.",
	}
	_, err := db.NewInsert().Model(&organize).Exec(context.Background())
	if err != nil {
		return err
	}

	branches := model.Branch{
		Name:           "Default Branch",
		Description:    "This is the default branch created during seed.",
		OrganizationID: organize.ID,
	}

	_, err = db.NewInsert().Model(&branches).Exec(context.Background())
	if err != nil {
		return err
	}

	departments := model.Department{
		Name:        "Default Department",
		Description: "This is the default department created during seed.",
		BranchID:    branches.ID,
	}
	_, err = db.NewInsert().Model(&departments).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}
