# go-jwt

JWT issuer, validator, token-revocation cache and context helpers for Go projects. Tokens are signed with Ed25519, valid or revoked token IDs can be tracked in memory, Redis or SQLite, and invalidation updates can be distributed over RabbitMQ. Requires Go 1.25 (per `go.mod`).

## Installation

```bash
go get github.com/ralvarezdev/go-jwt
```

Direct dependencies include `golang-jwt/jwt/v5`, `gin`, `amqp091-go`, `go-redis/v9`, `golang.org/x/crypto`, and `go-cache`, `go-databases`, `go-flags`, `go-strings` from `github.com/ralvarezdev`.

## Packages

- **`gojwt`** (root) — `BearerPrefix`, context keys `CtxTokenClaimsKey` / `CtxTokenKey`, claim names `IDClaim` (`jti`), `IsRefreshTokenClaim` (`irt`), `SubjectClaim` (`sub`), shared errors.
- **`token`** — `Token` type with `RefreshToken` and `AccessToken` kinds.
- **`token/issuer`** — `Issuer` interface and `NewEd25519Issuer(privateKey []byte)`.
- **`token/validator`** — `Validator` interface (`GetToken`, `GetClaims`, `ValidateClaims(ctx, rawToken, token)`) and `NewEd25519Validator(...)`.
- **`token/claims`** — `ClaimsValidator` and `TokenValidator` interfaces (`AddRefreshToken`, `AddAccessToken`, `RevokeToken`, `IsTokenValid`), `NewDefaultClaimsValidator`.
- **`token/claims/{cache,redis,sqlite}`** — in-memory, Redis-backed and SQLite-backed `NewTokenValidator(...)`.
- **`sync`, `sync/sqlite`** — `Service` interface tracking the last token-sync timestamp, with a SQLite implementation.
- **`rabbitmq`, `rabbitmq/publisher`, `rabbitmq/consumer`** — publish and consume token messages (`NewDefaultPublisher`, `NewDefaultConsumer`, `NewDefaultService`, `DeclareTokensMessageQueue`).
- **`flags`** — registers `-private-key` and `-public-key` (paths to PEM key files).
- **`gin`, `grpc`, `net/http`** — store and read the token and claims in a `gin.Context`, a gRPC `context.Context` (with subject and JWT ID getters) or an `*http.Request`.

## Usage (sketch)

```go
iss, err := issuer.NewEd25519Issuer(privateKeyBytes)
raw, err := iss.IssueToken(claims) // claims implements jwt.Claims
```

For validation, combine a `token/validator` Ed25519 validator with one of the claims validators and call `ValidateClaims(ctx, rawToken, gojwttoken.AccessToken)`. See the source for constructor parameters.

## Development

```bash
go build ./...
go vet ./...
```

There are no tests.

## License

GNU General Public License v3.0 (see [LICENSE](LICENSE)).
