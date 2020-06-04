package sql

import "regexp"

func getSQLCountStmt(query string) string {
	limit := regexp.MustCompile(`(?s)LIMIT .*`)
	queryWithoutLimit := limit.ReplaceAllString(query, "")

	order := regexp.MustCompile(`(?s)ORDER BY .*`)
	queryWithoutOrder := order.ReplaceAllString(queryWithoutLimit, "")

	queryFromPart := regexp.MustCompile(`(?s)FROM .*`)
	sqlCount := "SELECT count(*) " + queryFromPart.FindString(queryWithoutOrder)

	return sqlCount
}
