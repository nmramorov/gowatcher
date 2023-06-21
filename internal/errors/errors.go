package errors

import "errors"

var (
	ErrorMetricNotFound         = errors.New("no such metric")
	ErrorWrongStringConvertion  = errors.New("string conversion went wrong")
	ErrorWithEnvConfig          = errors.New("error with env config occurred")
	ErrorWithIntervalConvertion = errors.New("error converting Intervals to int64")
)
