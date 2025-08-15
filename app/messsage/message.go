package message

var (
	Success             = "success"
	InternalServerError = "internal-server-error"
	BadRequest          = "bad-request"
	Forbidden           = "forbidden"
	Unauthorized        = "unauthorized"
	InvalidRequest      = "invalid-request-form"

	UserAlreadyExists = "user-already-exists"
	UserNotFound      = "user-not-found"
	UserIsInUse       = "user-is-in-use"

	RoleNotFound = "role-not-found"
	RoleInUse    = "role-in-use"

	EmailAlreadyExists = "email-already-exists"
	EmailNotFound      = "email-not-found"

	DepartmentNotMatched = "department-not-matched"
	DepartmentNotFound   = "department-not-found"
	DepartmentInUse      = "department-in-use"

	BranchNotFound   = "branch-not-found"
	BranchNotMatched = "branch-not-matched"
	BranchInUse      = "branch-in-use"

	OrganizationNotFound   = "organization-not-found"
	OrganizationNotMatched = "organization-not-matched"
	OrganizationInUse      = "organization-in-use"

	OTPCodeInvalid = "otp-code-invalid"
	OTPCodeUsed    = "otp-code-used"
	OTPCodeExpired = "otp-code-expired"

	HolidayNotFound      = "holiday-not-found"
	HolidayInUse         = "holiday-in-use"
	HolidayAlreadyExists = "holiday-already-exists"
	HolidayInPast        = "holiday-in-past"

	WorkShiftNotFound      = "work-shift-not-found"
	WorkShiftInUse         = "work-shift-in-use"
	WorkShiftAlreadyExists = "work-shift-already-exists"

	LeaveNotFound = "leave-not-found"
	LeaveInUse    = "leave-in-use"

	DocumentNotFound        = "document-not-found"
	DocumentInUse           = "document-in-use"
	DocumentAlreadyExists   = "document-already-exists"
	DocumentAlreadyApproved = "document-already-approved"

	CheckInLocationNotAllowed = "check-in-location-not-allowed"

// ClientInvalid    = "Client invalid"
// GrantTypeInvalid = "Grant type invalid"

// RefreshTokenInvalid = "Refresh token invalid"
// AccessTokenInvalid  = "Access token invalid"
// TokenInvalid        = "Token invalid"

// UserInvalid         = "User invalid"
// UserNotFound        = "User not found"
// LoginFailed         = "Login failed"
// LoginPasswordFailed = "username or password is wrong"
)
