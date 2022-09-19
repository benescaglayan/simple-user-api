package error

import "errors"

var EmailAlreadyInUseError = errors.New("a user with that email already exists")

var BadRequestError = errors.New("bad request")

var NotFoundError = errors.New("user with that id does not exist")

var ServerError = errors.New("server error")
