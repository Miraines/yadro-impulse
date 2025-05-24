package errors

import "fmt"

type BiathlonError struct {
	Code    string
	Message string
	Cause   error
}

func (e *BiathlonError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

const (
	ErrCodeInvalidConfig      = "INVALID_CONFIG"
	ErrCodeInvalidEvent       = "INVALID_EVENT"
	ErrCodeCompetitorNotFound = "COMPETITOR_NOT_FOUND"
	ErrCodeInvalidTime        = "INVALID_TIME"
	ErrCodeFileNotFound       = "FILE_NOT_FOUND"
	ErrCodeProcessingFailed   = "PROCESSING_FAILED"
)

func NewInvalidConfigError(msg string, cause error) *BiathlonError {
	return &BiathlonError{ErrCodeInvalidConfig, msg, cause}
}

func NewInvalidEventError(msg string, cause error) *BiathlonError {
	return &BiathlonError{ErrCodeInvalidEvent, msg, cause}
}

func NewCompetitorNotFoundError(competitorID int) *BiathlonError {
	return &BiathlonError{ErrCodeCompetitorNotFound, fmt.Sprintf("competitor %d not found", competitorID), nil}
}

func NewInvalidTimeError(msg string, cause error) *BiathlonError {
	return &BiathlonError{ErrCodeInvalidTime, msg, cause}
}

func NewFileNotFoundError(filename string, cause error) *BiathlonError {
	return &BiathlonError{ErrCodeFileNotFound, fmt.Sprintf("file %s not found", filename), cause}
}

func NewProcessingFailedError(msg string, cause error) *BiathlonError {
	return &BiathlonError{ErrCodeProcessingFailed, msg, cause}
}
