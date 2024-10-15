module test

go 1.23.0

require (
	github.com/fl0mb/mockoidc v1.1.2
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292
)

require (
	github.com/go-jose/go-jose/v3 v3.0.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/text v0.3.6 // indirect
)

replace github.com/fl0mb/mockoidc v1.1.2 => ./mockoidc
