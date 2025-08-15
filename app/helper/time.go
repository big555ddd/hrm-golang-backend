package helper

import (
	"app/app/enum"
	workshiftdto "app/app/modules/workshift/dto"
	"math"
	"time"
)

func UnixToDay(unixTimestamp int64) enum.Day {
	t := time.Unix(unixTimestamp, 0)
	weekday := t.Weekday()

	switch weekday {
	case time.Sunday:
		return enum.DAY_SUNDAY
	case time.Monday:
		return enum.DAY_MONDAY
	case time.Tuesday:
		return enum.DAY_TUESDAY
	case time.Wednesday:
		return enum.DAY_WEDNESDAY
	case time.Thursday:
		return enum.DAY_THURSDAY
	case time.Friday:
		return enum.DAY_FRIDAY
	case time.Saturday:
		return enum.DAY_SATURDAY
	default:
		return enum.DAY_SUNDAY
	}
}

func GetScheduleForDay(Schedules []workshiftdto.ShiftSchedule, day enum.Day) *workshiftdto.ShiftSchedule {
	dayString := day.String()

	for _, schedule := range Schedules {
		if string(schedule.Day) == dayString {
			return &schedule
		}
	}

	return nil
}

func CalculateDistance(reqx float64, reqy float64, workx float64, worky float64) bool {
	distance := math.Sqrt(math.Pow(reqx-workx, 2) + math.Pow(reqy-worky, 2))
	return distance <= 100
}

func GetMonthRange(month int, year int) (int64, int64) {
	// Create start of month (first day at 00:00:00)
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.FixedZone("Asia/Bangkok", 7*3600))

	// Create end of month (first day of next month minus 1 nanosecond)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return startOfMonth.Unix(), endOfMonth.Unix()
}
