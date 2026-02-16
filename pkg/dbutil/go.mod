module github.com/bernardinorafael/gogem/pkg/dbutil

go 1.24.1

require (
	github.com/bernardinorafael/gogem/fault v0.1.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.11.2
)

replace github.com/bernardinorafael/gogem/fault => ../fault
