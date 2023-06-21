package errors

import "errors"

var (
	ErrorMetricNotFound         error = errors.New("no such metric")
	ErrorWrongStringConvertion  error = errors.New("string convertion went wrong")
	ErrorWithEnvConfig          error = errors.New("error with env config occured")
	ErrorWithIntervalConvertion error = errors.New("error converting Intervals to int64")
)
