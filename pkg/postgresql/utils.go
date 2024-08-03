package postgresql

import (
	"fmt"
	"net/url"
	"strings"
)

func replaceDbName(dbUrl, dbName string) string {
	parsed, err := url.Parse(dbUrl)
	if err != nil {
		return dbUrl
	}

	parsed.Path = "/" + dbName

	return parsed.String()
}

func parseDatabaseName(dbUrl string) (string, error) {
	parsed, err := url.Parse(dbUrl)
	if err != nil {
		return "", err
	}

	path := strings.TrimPrefix(parsed.Path, "/")

	if path == "" {
		return "", fmt.Errorf("empty db name")
	}

	return path, nil
}
