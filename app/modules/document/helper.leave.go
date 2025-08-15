package document

import (
	"app/app/enum"
	"app/app/helper"
	message "app/app/messsage"
	attendancedto "app/app/modules/attendance/dto"
	documentdto "app/app/modules/document/dto"
	workshiftdto "app/app/modules/workshift/dto"
	"app/internal/logger"
	"context"
	"errors"
	"time"
)

func (s *Service) calculateLeaveQuotaAndHours(ctx context.Context, leave *documentdto.LeaveReq, userID string) (float64, float64, error) {
	user, err := s.user.Svc.Get(ctx, userID)
	if err != nil {
		return 0, 0, err
	}
	if !user.IsActive {
		return 0, 0, errors.New(message.UserNotFound)
	}
	// Check if the leave exists
	ex, err := s.leave.Svc.ExistOnOrg(ctx, user.OrganizationID, leave.LeaveID)
	if err != nil {
		return 0, 0, err
	}
	if !ex {
		return 0, 0, errors.New(message.LeaveNotFound)
	}
	//calculate leave quota and hours with unix timestamps
	startDate := leave.StartDate
	endDate := leave.EndDate
	if startDate >= endDate {
		return 0, 0, errors.New(message.InvalidRequest)
	}
	leaveHours := (endDate - startDate) / 3600 // Convert seconds to hours
	logger.Infof("Leave hours: %d", leaveHours)
	logger.Infof("Leave start date: %d, end date: %d", startDate, endDate)
	if leaveHours <= 0 {
		logger.Err("Leave hours must be greater than 0", "leaveHours ", leaveHours)
		return 0, 0, errors.New(message.InvalidRequest)
	}
	// Get the leave quota for the user get workshift first
	workShift, err := s.workshift.Svc.Get(ctx, user.WorkShiftID)
	if err != nil {
		return 0, 0, err
	}

	// Calculate working hours between start and end date
	_, actualLeaveHours := s.calculateWorkingHours(startDate, endDate, workShift)

	// Calculate quota used based on 8-hour standard working day
	// Each full working day counts as 1 quota regardless of actual working hours
	quotaUsed := float64(actualLeaveHours) / 8.0 // Store as integer with 3 decimal places (1.375 = 1375)

	return quotaUsed, actualLeaveHours, nil
}

// calculateWorkingHours calculates working hours between two Unix timestamps
func (s *Service) calculateWorkingHours(startUnix, endUnix int64, workShift *workshiftdto.ListWorkShiftResponse) (float64, float64) {
	startTime := time.Unix(startUnix, 0).In(time.FixedZone("Asia/Bangkok", 7*3600))
	endTime := time.Unix(endUnix, 0).In(time.FixedZone("Asia/Bangkok", 7*3600))

	totalWorkingHours := float64(0)
	actualLeaveHours := float64(0)
	currentTime := startTime

	// If it's the same day
	if startTime.Format("2006-01-02") == endTime.Format("2006-01-02") {
		dayOfWeek := helper.UnixToDay(currentTime.Unix())
		if s.isWorkingDay(dayOfWeek, workShift) {
			// Get the work schedule for this day
			schedule := helper.GetScheduleForDay(workShift.Schedules, dayOfWeek)
			if schedule != nil {
				workStartHour := float64(schedule.StartTime)
				workEndHour := float64(schedule.EndTime)

				// Convert leave times to hours of day
				leaveStartHour := float64(startTime.Hour()) + float64(startTime.Minute())/60.0
				leaveEndHour := float64(endTime.Hour()) + float64(endTime.Minute())/60.0

				// Ensure leave times are within working hours
				if leaveStartHour < workStartHour {
					leaveStartHour = workStartHour
				}
				if leaveEndHour > workEndHour {
					leaveEndHour = workEndHour
				}

				// Calculate actual leave hours for this day
				if leaveEndHour > leaveStartHour {
					actualLeaveHours = float64((leaveEndHour-leaveStartHour)*100) / 100 // Round to 2 decimal places
					totalWorkingHours = float64(workEndHour - workStartHour)
				}
			}
		}
		return totalWorkingHours, actualLeaveHours
	}

	// Multi-day leave calculation
	for currentTime.Before(endTime) || currentTime.Format("2006-01-02") == endTime.Format("2006-01-02") {
		dayOfWeek := helper.UnixToDay(currentTime.Unix())
		nextDay := currentTime.AddDate(0, 0, 1)
		nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, nextDay.Location())

		if s.isWorkingDay(dayOfWeek, workShift) {
			schedule := helper.GetScheduleForDay(workShift.Schedules, dayOfWeek)
			if schedule != nil {
				workStartHour := float64(schedule.StartTime)
				workEndHour := float64(schedule.EndTime)
				dailyWorkingHours := (workEndHour - workStartHour)

				// First day - from leave start time to end of work day
				if currentTime.Format("2006-01-02") == startTime.Format("2006-01-02") {
					leaveStartHour := float64(startTime.Hour()) + float64(startTime.Minute())/60.0
					if leaveStartHour < workStartHour {
						leaveStartHour = workStartHour
					}

					leaveHoursThisDay := ((workEndHour - leaveStartHour) * 100) / 100
					if leaveHoursThisDay > 0 {
						actualLeaveHours += leaveHoursThisDay
						totalWorkingHours += dailyWorkingHours
					}
				} else if currentTime.Format("2006-01-02") == endTime.Format("2006-01-02") {
					// Last day - from start of work day to leave end time
					leaveEndHour := float64(endTime.Hour()) + float64(endTime.Minute())/60.0
					if leaveEndHour > workEndHour {
						leaveEndHour = workEndHour
					}

					leaveHoursThisDay := ((leaveEndHour - workStartHour) * 100) / 100
					if leaveHoursThisDay > 0 {
						actualLeaveHours += leaveHoursThisDay
						totalWorkingHours += dailyWorkingHours
					}
				} else {
					// Middle days - full day leave
					actualLeaveHours += dailyWorkingHours
					totalWorkingHours += dailyWorkingHours
				}
			}
		}

		// Move to next day
		currentTime = nextDay
		if currentTime.After(endTime) {
			break
		}
	}

	return totalWorkingHours, actualLeaveHours
}

// getDailyWorkingHours returns working hours per day from workshift
func (s *Service) getDailyWorkingHours(workShift *workshiftdto.ListWorkShiftResponse, day enum.Day) int64 {

	dayString := day.String()

	for _, schedule := range workShift.Schedules {

		if string(schedule.Day) == dayString {
			start := schedule.StartTime
			end := schedule.EndTime

			workingHours := int64(end - start)
			if workingHours > 0 {
				return workingHours
			}
		}
	}

	return 0
}

// isWorkingDay checks if a specific day is a working day based on workshift schedule
func (s *Service) isWorkingDay(day enum.Day, workShift *workshiftdto.ListWorkShiftResponse) bool {
	// Get working hours for this day
	workingHours := s.getDailyWorkingHours(workShift, day)
	// If working hours > 0, it's a working day
	return workingHours > 0
}

func (s *Service) createLeaveRecord(ctx context.Context, documentID string) error {
	// Calculate leave quota and hours
	doc, err := s.Get(ctx, documentID)
	if err != nil {
		return err
	}

	if (doc.Type != enum.DOCUMENT_TYPE_LEAVE) || (doc.LeaveDetails == nil) {
		return errors.New(message.InvalidRequest)
	}

	start := doc.LeaveDetails.StartDate
	end := doc.LeaveDetails.EndDate
	//get workshift for the user
	user, err := s.user.Svc.Get(ctx, doc.UserID)
	if err != nil {
		return err
	}

	// Get workshift details for working day validation
	workShift, err := s.workshift.Svc.Get(ctx, user.WorkShiftID)
	if err != nil {
		return err
	}

	// Convert timestamps to time objects with timezone
	startTime := time.Unix(start, 0).In(time.FixedZone("Asia/Bangkok", 7*3600))
	endTime := time.Unix(end, 0).In(time.FixedZone("Asia/Bangkok", 7*3600))

	// Create attendance records for each day in the leave period
	currentTime := startTime
	for currentTime.Before(endTime) || currentTime.Format("2006-01-02") == endTime.Format("2006-01-02") {
		dayOfWeek := helper.UnixToDay(currentTime.Unix())

		// Only create attendance record for working days
		if s.isWorkingDay(dayOfWeek, workShift) {
			// Create attendance record for this day
			dayStartUnix := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).Unix()

			m := &attendancedto.CreateAttendance{
				UserID:      doc.UserID,
				WorkShiftID: user.WorkShiftID,
				CheckIn:     0,
				CheckOut:    0,
				Date:        dayStartUnix,
				IsOnTime:    true,
				IsLate:      false,
				IsLeave:     true,
			}

			_, err = s.attendance.Svc.Create(ctx, m)
			if err != nil {
				return err
			}
		}

		// Move to next day
		currentTime = currentTime.AddDate(0, 0, 1)
		if currentTime.After(endTime) {
			break
		}
	}

	return nil
}
