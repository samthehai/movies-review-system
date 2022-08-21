//go:build tools

package tools

import (
	_ "github.com/golang/mock/mockgen"
	_ "github.com/rubenv/sql-migrate/sql-migrate"
	_ "github.com/swaggo/swag/cmd/swag"
)
