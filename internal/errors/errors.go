package errors

import "errors"

var ErrorMetricNotFound error = errors.New("no such metric")
var ErrorWrongStringConvertion error = errors.New("string convertion went wrong")
var ErrorWithEnvConfig error = errors.New("error with env config occured")
var ErrorWithIntervalConvertion error = errors.New("error converting Intervals to int64")
