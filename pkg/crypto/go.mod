module github.com/bernardinorafael/gogem/pkg/crypto

go 1.24.1

require (
	github.com/bernardinorafael/gogem/uid v0.1.0
	github.com/golang-jwt/jwt/v5 v5.3.1
	golang.org/x/crypto v0.48.0
)

replace github.com/bernardinorafael/gogem/uid => ../uid
