package enum

type OverTimeType string

const (
	STATUS_OVERTIME_DAY_WORK    OverTimeType = "dayWork"
	STATUS_OVERTIME_DAY_OF_WORK OverTimeType = "dayOfWork"
	STATUS_OVERTIME_HOLIDAY     OverTimeType = "holiday"
)

func GetOverTimeType(t OverTimeType) OverTimeType {
	switch t {
	case STATUS_OVERTIME_DAY_WORK:
		return STATUS_OVERTIME_DAY_WORK
	case STATUS_OVERTIME_DAY_OF_WORK:
		return STATUS_OVERTIME_DAY_OF_WORK
	case STATUS_OVERTIME_HOLIDAY:
		return STATUS_OVERTIME_HOLIDAY
	default:
		return STATUS_OVERTIME_DAY_WORK
	}
}

// ConvertToFloat converts OverTimeType string to float64
func (o OverTimeType) ConvertToFloat() float64 {
	switch o {
	case STATUS_OVERTIME_DAY_WORK:
		return 1.5
	case STATUS_OVERTIME_DAY_OF_WORK:
		return 2.0
	case STATUS_OVERTIME_HOLIDAY:
		return 3.0
	default:
		return 1.5
	}
}
