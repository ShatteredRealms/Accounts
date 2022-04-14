package model

// ErrorModel A model representing all information pertaining to an individual error
type ErrorModel struct {
	// ErrorCode Error code for the error
	ErrorCode int `json:"error_code"`

	// Text The error text message
	Text string `json:"text"`

	// Hints Hints that may help resolve the error
	Hints string `json:"hints"`

	// Info Information about the error
	Info string `json:"info"`
}

// HTTP errors
var (
	// UnsupportedMediaError The base information for an HTTP Unsupported Media rror (415).
	UnsupportedMediaError = ErrorModel{
		ErrorCode: 4151,
		Text:      "Unsupported Media Type",
		Hints:     "Set the content-type in the header to 'application/json'",
		Info:      "Only Accept JSON Data type",
	}

	// NotFoundError The base information for an HTTP Not Found error (404)
	NotFoundError = ErrorModel{
		ErrorCode: 4041,
		Text:      "Not Found",
		Hints:     "Ensure you're using the correct version path for this release",
		Info:      "Route does not exist",
	}

	// BadRequestError The base information for an HTTP Bad Request error (400)
	BadRequestError = ErrorModel{
		ErrorCode: 4001,
		Text:      "Bad Request",
		Hints:     "Ensure you're using the correct payload",
		Info:      "",
	}

	// InternalServError The base information for an HTTP Internal Server error (500)
	InternalServError = ErrorModel{
		ErrorCode: 5001,
		Text:      "Internal Server Error",
		Hints:     "Contact the support team if the issue persists",
		Info:      "Please try again in a few moments",
	}

	// UnauthorizedError The base information for an HTTP Unauthorized error (401)
	UnauthorizedError = ErrorModel{
		ErrorCode: 4011,
		Text:      "Unauthorized",
		Hints:     "If you're logging in, your credentials was denied",
		Info:      "You do not have access to this resources",
	}
)

// Application errors
var (

	// RegistrationFailedError The base information resulting for a registration failed error
	RegistrationFailedError = ErrorModel{
		ErrorCode: 100,
		Text:      "Registration Failed",
		//Hints:     "A user with the given information may already exist",
		Info: "The supplied registration information is insufficient",
	}
)
