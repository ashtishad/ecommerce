package constants

const (
	EmailRegex        = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	FullNameRegex     = `^[a-zA-Z\s]+$`
	PhoneRegex        = `^\d{10,15}$`
	TimezoneRegex     = `^(UTC|[A-Za-z]+(?:/[A-Za-z_]+)+)$`
	UUIDRegex         = `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[1-5][a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`
	StatusRegex       = `^(active|inactive|deleted)$`
	SignUpOptionRegex = `^(google|general)$`
	SignUpOptGoogle   = "google"
	SignupOptGeneral  = "general"
	UserStatusActive  = "active"

	UserStatusInactive = "inactive"
	UserStatusDeleted  = "deleted"

	DefaultPageSizeString = "20"
)
