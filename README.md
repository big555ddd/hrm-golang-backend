# HRM Golang Backend

This repository is a Human Resource Management (HRM) backend service written in Go. Below is the project structure and a brief description of each main directory and file:

## Project Structure

```
docker-compose.yml      # Docker Compose configuration for multi-container setup
Dockerfile              # Dockerfile for building the application image
go.mod, go.sum          # Go modules dependencies
main.go                 # Application entry point
Makefile                # Build and management commands

app/
  console/              # Console commands and kernel logic
  enum/                 # Enumerations (day, document type, gender, status, etc.)
  helper/               # Helper functions (email, time, etc.)
  messsage/             # Message definitions
  middleware/           # Middleware (auth, logger, permission)
  model/                # Data models (user, document, attendance, etc.)
  modules/              # Business modules (attendance, auth, branch, department, etc.)
    activitylog/        # Activity log module
    attendance/         # Attendance module
    ...                 # Other modules
  response/             # Response formatting
  routes/               # API route definitions
  util/                 # Utility functions (hashing, jwt)

config/                 # Configuration files (database, mail, redis, etc.)
database/
  migrations/           # Database migration files
  seeds/                # Database seed files
internal/
  interface.go          # Internal interfaces
  cmd/                  # Internal command logic
  database/             # Internal database logic
  logger/               # Internal logger
```

## Getting Started

## Main Components

### Models

- User, Document, Attendance, Role, Permission, Department, Organization, WorkShift, Holiday, Notification, Leave, ActivityLog, DocumentLeave, DocumentOvertime, HolidayOrganization, LeaveOrganization, UserDepartment, UserForgot, UserRole, UserWorkShift, ShiftSchedule
- Common fields: `CreateUpdateUnixTimestamp`, `SoftDelete` for audit and soft delete support

### Enums

- Day: Sunday-Saturday
- DocumentType: leave, overtime, addTime
- Gender: Unknown, Female, Male
- StatusDocument: pending, approved, rejected
- Status: active, inactive
- OverTimeType: dayWork, dayOfWork, holiday

### Helpers

- Email: OTP/REF code generation, email masking
- Time: Unix to Day conversion, schedule lookup, distance calculation, month range
- General: Get user from JWT token

### Middleware

- Auth: JWT authentication for API endpoints
- Logger: Logs requests/responses, creates activity logs
- Permission: Checks user permissions for protected routes

### Configuration

- config.go: Loads environment variables, sets defaults
- database.go: Database connection and management
- mail.go: Email service setup and sending
- redis.go: Redis client setup

### Messages

- Standardized error and status messages (success, internal-server-error, forbidden, unauthorized, user/role/department/branch/org/holiday/email errors, OTP errors)

## License

MIT
