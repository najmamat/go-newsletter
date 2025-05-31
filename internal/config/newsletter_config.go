package config

// NewsletterConfig contains all business rules and constants for the newsletter domain
type NewsletterConfig struct {
	// Name constraints
	MaxNameLength        int
	MinNameLength        int
	RequiredNameMessage  string
	TooLongNameMessage   string
	TooShortNameMessage  string
	EmptyNameMessage     string
	DuplicateNameMessage string

	// Description constraints
	MaxDescriptionLength int
	TooLongDescMessage   string

	// ID constraints
	InvalidIDMessage string
}

// DefaultNewsletterConfig returns the default configuration for newsletters
func DefaultNewsletterConfig() *NewsletterConfig {
	return &NewsletterConfig{
		// Name constraints
		MaxNameLength:        100,
		MinNameLength:        1,
		RequiredNameMessage:  "Newsletter name is required",
		TooLongNameMessage:   "Newsletter name must be less than 100 characters",
		TooShortNameMessage:  "Newsletter name must be at least 1 character",
		EmptyNameMessage:     "Newsletter name cannot be empty",
		DuplicateNameMessage: "You already have a newsletter with this name",

		// Description constraints
		MaxDescriptionLength: 500,
		TooLongDescMessage:   "Description must be less than 500 characters",

		// ID constraints
		InvalidIDMessage: "Invalid newsletter ID format",
	}
}
