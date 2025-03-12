package proto

var (
	// General.
	ErrTODO           = WebRPCError{Code: 1000, Name: "NotImplemented", Message: "TODO: Not implemented", HTTPStatus: 500}
	ErrDeprecated     = WebRPCError{Code: 1001, Name: "Deprecated", Message: "Endpoint is deprecated", HTTPStatus: 410}
	ErrRateLimited    = WebRPCError{Code: 1002, Name: "RateLimited", Message: "Rate limited. Please, slow down", HTTPStatus: 429}
	ErrInvalidRequest = WebRPCError{Code: 1003, Name: "InvalidRequest", Message: "Invalid request", HTTPStatus: 400}
	ErrUnexpected     = WebRPCError{Code: 1004, Name: "Unexpected", Message: "Unexpected server error", HTTPStatus: 500}

	// Pets.
	ErrPetNotFound = WebRPCError{Code: 2000, Name: "PetNotFound", Message: "Pet not found", HTTPStatus: 400}
)
