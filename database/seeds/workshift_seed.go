package seeds

import (
	"app/app/enum"
	"app/app/model"
	"context"

	"github.com/uptrace/bun"
)

func workShiftSeed(db *bun.DB) error {
	workShift := model.WorkShift{

		Name:          "Default Work Shift",
		WorkLocationX: 0,
		WorkLocationY: 0,
		LateMinutes:   15,
		Description:   "This is the default work shift created during seed.",
	}
	_, err := db.NewInsert().Model(&workShift).Exec(context.Background())
	if err != nil {
		return err
	}

	shiftSchedule := []model.ShiftSchedule{
		{
			WorkShiftID: workShift.ID,
			Day:         enum.DAY_MONDAY,
			StartTime:   9,
			EndTime:     17,
		},
		{
			WorkShiftID: workShift.ID,
			Day:         enum.DAY_TUESDAY,
			StartTime:   9,
			EndTime:     17,
		},
		{
			WorkShiftID: workShift.ID,
			Day:         enum.DAY_WEDNESDAY,
			StartTime:   9,
			EndTime:     17,
		},
		{
			WorkShiftID: workShift.ID,
			Day:         enum.DAY_THURSDAY,
			StartTime:   9,
			EndTime:     17,
		},
		{
			WorkShiftID: workShift.ID,
			Day:         enum.DAY_FRIDAY,
			StartTime:   9,
			EndTime:     17,
		},
	}

	_, err = db.NewInsert().Model(&shiftSchedule).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}
