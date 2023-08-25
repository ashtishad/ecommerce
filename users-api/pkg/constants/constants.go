package constants

const (
	EmailRegex       = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	FullNameRegex    = `^[a-zA-Z\s]+$`
	PhoneRegex       = `^\d{10,15}$`
	TimezoneRegex    = `^(UTC|[A-Za-z]+(?:/[A-Za-z_]+)+)$`
	SignUpOptGoogle  = "google"
	SignupOptGeneral = "general"
	UUIDRegex        = `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[1-5][a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`
	UserStatusActive = "active"

	// UserStatusInactive = "inactive"
	// UserStatusDeleted  = "deleted"
)
