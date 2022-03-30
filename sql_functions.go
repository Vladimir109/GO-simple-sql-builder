package sql_builder

import "strings"

func COUNT(alias, field string) string {
	if alias != "" { alias += "." }
	return "COUNT(" + alias + field + ")"
}

func CONCAT(fieldVals ...string) string {
	return "CONCAT(" + strings.Join(fieldVals, ", ") + ")"
}
func CONCAT_WS(delimiter string, fieldVals ...string) string {
	if delimiter != "?" {
		delimiter = "'"+delimiter+"'"
	}
	return "CONCAT_WS(" + delimiter + "," + strings.Join(fieldVals, ", ") + ")"
}

func GROUP_CONCAT(separator string, fieldVals ...string) string {
	if separator != "" { separator = "SEPARATOR '"+separator+"'" }
	return "GROUP_CONCAT(" + strings.Join(fieldVals, ", ") + separator + ")"
}
