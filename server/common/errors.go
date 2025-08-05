package common

import "fmt"

var ErrFileTooLarge = fmt.Errorf("file too large")
var ErrFileContentType = fmt.Errorf("invalid file type")
var ErrReadFileFromFormFailed = fmt.Errorf("failed to read file from form")
var ErrSaveFileFail = fmt.Errorf("failed to save file")
var ErrInvalidJSON = fmt.Errorf("invalid JSON format")
var ErrInvalidRequest = fmt.Errorf("invalid request")
var ErrInvalidCredentials = fmt.Errorf("invalid credentials")

var ErrUserExists = fmt.Errorf("user already exists")

var ErrUserCreationFailed = fmt.Errorf("failed to create user")
