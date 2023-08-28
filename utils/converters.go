package utils

import "github.com/jackc/pgtype"

// StringToNumeric convert  string to pgtype numeric
func StringToNumeric(s string) (pgtype.Numeric, error) {
	var numeric pgtype.Numeric
	err := numeric.Set(s)
	if err != nil {
		return numeric, err
	}
	return numeric, nil
}
