package constant

const (
	InternalErrorCode = 1
	InternalError     = "It seems to meet an internal error, please wait for a while."
)

const (
	GateWayParamsErrorCode = 2
	GateWayParamsError     = "The gateway cannot response for the prameters. Please check your parameters."
)

const (
	UserExistsErrorCode = 3
	UserExistsError     = "User exists, choose another username."
)

const (
	UserNotExistsErrorCode = 4
	UserNotExistsError     = "User not exists!"
)

const (
	LoginFailErrorCode = 5
	LoginFailError     = "Login failed. Please check your password and username."
)

const (
	UnauthorizedErrorCode = 6
	UnauthorizedError     = "Visiting is not authorized. Please login."
)
