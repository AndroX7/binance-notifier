package utility

import (
	"fmt"
	"strconv"
)

func ParseFloat(v interface{}) (float64, error) {
	switch t := v.(type) {
	case string:
		return strconv.ParseFloat(t, 64)
	case float64:
		return t, nil
	default:
		return 0, fmt.Errorf("unknown type for parseFloat")
	}
}
