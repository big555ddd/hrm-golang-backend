package enum

type Day string

const (
	DAY_SUNDAY    Day = "Sunday"
	DAY_MONDAY    Day = "Monday"
	DAY_TUESDAY   Day = "Tuesday"
	DAY_WEDNESDAY Day = "Wednesday"
	DAY_THURSDAY  Day = "Thursday"
	DAY_FRIDAY    Day = "Friday"
	DAY_SATURDAY  Day = "Saturday"
)

func GetDay(t Day) Day {
	switch t {
	case DAY_SUNDAY:
		return DAY_SUNDAY
	case DAY_MONDAY:
		return DAY_MONDAY
	case DAY_TUESDAY:
		return DAY_TUESDAY
	case DAY_WEDNESDAY:
		return DAY_WEDNESDAY
	case DAY_THURSDAY:
		return DAY_THURSDAY
	case DAY_FRIDAY:
		return DAY_FRIDAY
	case DAY_SATURDAY:
		return DAY_SATURDAY
	default:
		return DAY_SUNDAY
	}
}

func (d Day) String() string {
	switch d {
	case DAY_SUNDAY:
		return string(DAY_SUNDAY)
	case DAY_MONDAY:
		return string(DAY_MONDAY)
	case DAY_TUESDAY:
		return string(DAY_TUESDAY)
	case DAY_WEDNESDAY:
		return string(DAY_WEDNESDAY)
	case DAY_THURSDAY:
		return string(DAY_THURSDAY)
	case DAY_FRIDAY:
		return string(DAY_FRIDAY)
	case DAY_SATURDAY:
		return string(DAY_SATURDAY)
	default:
		return "Unknown Day"
	}
}
