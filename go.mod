module test

go 1.23.0

require (
	github.com/fl0mb/mockoidc v1.1.2
	golang.org/x/crypto v0.17.0
)

require (
	github.com/go-jose/go-jose/v3 v3.0.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

replace github.com/fl0mb/mockoidc v1.1.2 => ./mockoidc
