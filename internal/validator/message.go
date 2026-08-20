package validator

var (
	//	Built-in validation rules message's
	Required   = "%s field is required"
	Max        = "%s field can be at most %s characters long"
	Min        = "%s field must be at least %s characters long"
	EqField    = "%s field must be equal to %s field"
	AlphaSpace = "%s field can contain only alphabetic characters and spaces"
	Longitude  = "%s field can contain only longitude format"
	Latitude   = "%s field can contain only latitude format"

	//	Custom
	Login = "%s field must be a valid email address or a username, allowed following characters: ._-"

	//	Default validation message
	Default = "%s field failed on %s validation rule"
)
