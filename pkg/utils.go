package pkg

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

// FloatToNumeric converts a float64 to a pgtype.Numeric
func FloatToNumeric(f float64) pgtype.Numeric {
	n := pgtype.Numeric{}
	_ = n.Scan(fmt.Sprintf("%f", f))
	return n
}

// NumericToFloat safely converts a pgtype.Numeric to float64
func NumericToFloat(n pgtype.Numeric) float64 {
	f, _ := n.Float64Value()
	return f.Float64
}

// StringToText converts a string to a valid pgtype.Text
func StringToText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

// TextToString safely extracts a string from pgtype.Text
func TextToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}
