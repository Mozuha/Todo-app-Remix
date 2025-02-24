package utils

import "errors"

var MsgResourceNotFound = "Resource not found"
var MsgUIDNotFoundInCtx = "UserID not found in context"
var MsgInternalServerErr = "The server encountered unexpected error"
var MsgInvalidReq = "Invalid request"
var MsgInvalidEmailOrPswd = "Invalid email or password"

var ErrUIDNotFoundInCtx = errors.New("userID not found in context")
var ErrNoRowsMatchedSQLC = errors.New("no rows in result set")
var ErrInvalidEmailOrPswd = errors.New("invalid email or password")
var ErrInvalidUID = errors.New("invalid userID")
