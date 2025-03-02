package common

import "math"

const (
	TOLERANCE = 0.001
)

func FloatsEqual(a float64, b float64) bool {
	return math.Abs(a-b) <= TOLERANCE

}
