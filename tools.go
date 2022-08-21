//go:build tools

package tools

import (
	_ "github.com/golang/mock/mockgen"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/rubenv/sql-migrate/sql-migrate"
	_ "golang.org/x/tools/cmd/goimports"
)
