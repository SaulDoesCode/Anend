package backend

import (
	"errors"
	"io/ioutil"

	tr "github.com/SaulDoesCode/transplacer"
)

// Errors for days
var (
	// ErrBadDBConnection bad database connection, try different details
	ErrBadDBConnection = errors.New("bad database connection error, try different details")
	// ErrInvalidEmail bad email
	ErrInvalidEmail = errors.New(`invalid email address`)
	// ErrIncompleteWrit someone probably tried to mutate a writ that is invalid or non existing (in db)
	ErrIncompleteWrit = errors.New(`attempted to modify either an invalid writ, or one not in the db`)
	// ErrMissingTags writ is missing tags
	ErrMissingTags = errors.New(`writ doesn't have any tags, add some`)
	// ErrAuthorIsNoUser writ's author is persona non grata
	ErrAuthorIsNoUser = errors.New(`writ author is not a registered user`)
	// UnauthorizedError unauthorized request, cannot proceed
	UnauthorizedError = MakeCodedResponse(403, obj{"err": "unauthorized request, cannot proceed"})
	// InvalidDetailsError invalid details, could not authorize user
	InvalidDetailsError = MakeCodedResponse(401, obj{"err": "invalid details, could not authorize user"})
	// BadUsernameError invalid username, could not authorize user
	BadUsernameError = MakeCodedResponse(401, obj{"err": "invalid username, could not authorize user"})
	// BadEmailError invalid email, could not authorize user
	BadEmailError = MakeCodedResponse(401, obj{"err": "invalid email, could not authorize user"})
	// BadRequestError bad request, check details and try again
	BadRequestError = MakeCodedResponse(400, obj{"err": "bad request, check details and try again"})
	// ServerDecodeError ran into trouble decoding your request
	ServerDecodeError = MakeCodedResponse(400, obj{"err": "ran into trouble decoding your request"})
	// ServerDBError server error, could not complete your request
	ServerDBError = MakeCodedResponse(500, obj{"err": "server error, could not complete your request"})
	// AlreadyLoggedIn user is logged in, but they tried to login again.
	AlreadyLoggedIn = MakeCodedResponse(203, obj{"msg": "You're already logged in :D"})
	// RateLimitingError somebody probably sent too many emails
	RateLimitingError = MakeCodedResponse(429, obj{"err": "ratelimiting: too many auth requests/emails, wait a bit and try again"})
	// NoSuchWrit could not find a writ matching the query/request
	NoSuchWrit = MakeCodedResponse(404, obj{"err": "couldn't find a writ like that"})
	// RequestQueryOverLimitMembers a member is requesting too many things at once
	RequestQueryOverLimitMembers = MakeCodedResponse(403, obj{"err": "requesting too many items at once, members >= 200"})
	// RequestQueryOverLimit a viewer (non-member-user) is requesting too many things at once
	RequestQueryOverLimit = MakeCodedResponse(403, obj{"err": "requesting too many items at once, non-members >= 100"})
	// SuccessMsg send a success response
	SuccessMsg = MakeCodedResponse(203, obj{"msg": "success!"})
	// DeleteWritError there was trouble when attempting to delete a writ, prolly database/bad-input related
	DeleteWritError = MakeCodedResponse(500, obj{"err": "could not delete writ, maybe it didn't exist in the first place"})
)

// PageError implements error but can send an .html file as a response
type PageError struct {
	Code    int
	Path    string
	Content []byte
	Value   string
}

func (pe *PageError) Error() string {
	return pe.Value
}

// Send responds to a request with an .html error page
func (pe *PageError) Send(c ctx) error {
	// See https://github.com/golang/go/issues/27139
	// We can't set the status code to the error's
	// Because http.ServeContent does this with no
	// way to change it :(
	return c.HTMLBlob(pe.Code, pe.Content)
}

// MakePageErr generates a new *PageError
func MakePageErr(code int, value, path string) *PageError {
	path = tr.PrepPath(Conf.Assets, path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic("PageError needs a valid path to a servable error page: " + err.Error())
	}
	return &PageError{Code: code, Path: path, Content: content, Value: value}
}

// Err404NotFound is the standard 404 error response returned by Anend
var Err404NotFound *PageError
