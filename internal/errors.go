package metrics

import "errors"

var ErrorMetricNotFound error = errors.New("no such metric")
var ErrorWrongStringConvertion error = errors.New("string convertion went wrong")
