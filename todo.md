    # BookVault → Production-Grade Backend: TODO

Reviewed: `dev` branch, commit at time of review. Every item below is traced to an actual
file/line in the repo — this isn't generic advice.

Work top to bottom. Phase 0 is "must fix before you show this to anyone." Phase 1 is your
LocalStack→AWS migration. The rest is what turns a working app into an operable one.

---

## Phase 0 — Correctness bugs (fix today, these are actively broken)

- [ ] **Cart timestamps are always zero.** `internal/services/cart_service.go:134-135` —
  `convertToCartResponse` builds `dto.CartItemResponse{CreatedAt: cartItems[i].CreatedAt, ...}`
  but `cartItems` is the *destination* slice being built (`make([]dto.CartItemResponse, ...)`),
  not the source. It should read from `cart.CartItems[i].CreatedAt` / `.UpdatedAt`. Every cart
  item in every API/GraphQL response currently shows `0001-01-01T00:00:00Z`.

- [ ] **Refresh tokens silently die after ~24h instead of 30 days.**
  `internal/utils/jwt.go:47` — the refresh JWT's `ExpiresAt` claim is set with `cfg.ExpiresIn`
  (the *access* token TTL), not `cfg.RefreshTokenExpires`. The DB row for the refresh token is
  valid 30 days, but the JWT itself is dead in 24h, so `ValidateToken` will reject it long
  before the DB check ever runs. Refresh-token flow is effectively non-functional past day 1.

- [ ] **`orderRoutes.POST("/orders", ...)` doubles the path.**
  `internal/server/server.go:128-131` — the group is already mounted at `/api/v1/orders`, so
  this registers `POST /api/v1/orders/orders`, not `/api/v1/orders` (contradicts the Swagger
  comment on the handler and every other route in that group). Change to `orderRoutes.POST("/", s.createOrder)`.

- [ ] **`deleteProduct` is missing a `return` after a bad ID.**
  `internal/server/product_handlers.go:237-241` — on `ParseUint` failure it calls
  `BadRequestResponse` but doesn't `return`, so execution falls through and calls
  `DeleteProduct(uint(id))` with `id=0` anyway, then tries to write a second response.

- [ ] **Cart writes fail silently.** `internal/services/cart_service.go` — in `AddToCart`,
  both `s.db.Create(&cartItem)` and `s.db.Save(&cartItem)` have their returned `error`
  discarded. A failed insert/update currently looks identical to a successful one to the caller.

- [ ] **Auth is coupled to your event bus's uptime.** `internal/services/auth_service.go` —
  `Register`/`Login` call `eventPublisher.Publish("USER_LOGGED_IN", ...)` synchronously and
  return an error (failing the whole login/register) if the SQS publish fails. A blip in SQS
  should never lock users out of logging in. Publish in a goroutine (or fire-and-forget with
  logging) instead of blocking the auth path on it.

- [ ] **Cart-creation failure during registration is swallowed.**
  `internal/services/auth_service.go` — `fmt.Println("unable to creatr cart")` (typo, and wrong
  logger — you already have zerolog) hides the error entirely. User account is created with no
  cart and nobody finds out. Wrap `Register` (create user + create cart) in a single DB
  transaction so it's all-or-nothing.

---

## Phase 1 — LocalStack → AWS migration (your main ask)

The current code doesn't have an AWS "mode" — it has a LocalStack mode that happens to also
work against AWS by accident for S3, and will actively misbehave for SQS. Fix the design, not
just the `.env` values.

- [ ] **Stop hardcoding `test`/`test` credentials.**
  `internal/providers/aws.go` — right now, *any time* `S3Endpoint` is non-empty, it force-sets
  static creds `"test"/"test"` (a LocalStack-only assumption baked into shared code). Split this
  into two paths:
  - **Local/dev** (endpoint set, e.g. LocalStack): keep static test creds + path-style addressing.
  - **AWS** (no custom endpoint): call `config.LoadDefaultConfig(ctx, config.WithRegion(region))`
    with **no static credentials at all**, so the SDK's default credential chain picks up an
    IAM role (EC2 instance profile / ECS task role / EKS IRSA). Never put real AWS access keys
    in `.env` for a running service — that's what roles are for.

- [ ] **Separate the S3 endpoint from the SQS endpoint — don't reuse one config value for both.**
  `internal/events/watermill.go` — the SQS client is built with `providers.CreateAWSConfig(ctx, cfg.S3Endpoint, cfg.Region)`.
  This only "works" because LocalStack multiplexes every service behind `:4566`. Real AWS has
  distinct endpoints per service and per region (`sqs.<region>.amazonaws.com`, `s3.<region>.amazonaws.com`).
  Add a separate `cfg.SQSEndpoint` (empty in AWS, `http://localhost:4566` in LocalStack) and stop
  passing S3's endpoint into the SQS client constructor. Same fix needed anywhere else `S3Endpoint`
  is reused for a non-S3 client.

- [ ] **Turn off forced path-style addressing for real AWS.**
  `internal/providers/s3.go` — `UsePathStyle = true` is set whenever an endpoint is configured.
  Path-style is a LocalStack/MinIO requirement; real S3 should use virtual-hosted style
  (`UsePathStyle = false`, and just don't set a custom endpoint at all in AWS mode).

- [ ] **Add IAM least-privilege policies** for whatever role/user runs in AWS: scope S3 access
  to `s3:PutObject`/`s3:GetObject`/`s3:DeleteObject` on the specific bucket+prefix only, and SQS
  to `sqs:SendMessage`/`sqs:ReceiveMessage`/`sqs:DeleteMessage` on the specific queue ARN only.
  Don't reuse one broad policy across both.

- [ ] **Real S3 bucket setup that `init-localstack.sh` currently fakes for you:**
  bucket versioning (recommended), default encryption (SSE-S3 or SSE-KMS), a bucket policy that
  blocks public access unless you intentionally want public reads, and lifecycle rules if
  uploads should expire/transition to cheaper storage classes.

- [ ] **Real SQS queue setup that `init-localstack.sh` fakes for you:** a redrive policy pointing
  at a dead-letter queue (so a poison message doesn't loop forever), and a visibility timeout
  tuned to your notifier's processing time.

- [ ] **Write this as Infrastructure-as-Code** (Terraform or CDK) rather than shell scripts —
  you'll want the same bucket/queue defs reproducible across dev/staging/prod AWS accounts, not
  just LocalStack. This is also where bucket policy, encryption, and DLQ config above should live.

- [ ] **Fix `.env.example` inconsistency**: `AWS_S3_ENDPOINT=http://localhost:9000` (port 9000
  is MinIO's default, not LocalStack's `4566`) — pick one and make docker-compose match it.

- [ ] **Config validation on startup**: right now missing/blank env vars silently fall back to
  defaults (including a placeholder JWT secret in `internal/config/config.go`). Add an explicit
  `Validate()` step that hard-fails startup if `JWT_SECRET`, `DB_PASSWORD`, or (in AWS mode)
  `AWS_S3_BUCKET`/`SQS` queue URL are empty, so a misconfigured prod deploy fails at boot, not
  at 2am when the first request comes in.

---

## Phase 2 — Data integrity & concurrency

- [ ] **Prevent overselling stock.** `AddToCart` / `UpdateCartItem`
  (`internal/services/cart_service.go`) and the stock check in
  `internal/services/order_service.go` read `Product.Stock`, then act on it later with no
  transaction/row lock. Two concurrent requests can both pass the check. Fix with either:
  - an atomic `UPDATE products SET stock = stock - ? WHERE id = ? AND stock >= ?` and checking
    `RowsAffected`, or
  - `SELECT ... FOR UPDATE` (GORM: `clause.Locking{Strength: "UPDATE"}`) inside a transaction.

- [ ] **`RemoveFromCart` doesn't check `RowsAffected`.** A delete that matches zero rows (item
  doesn't exist / doesn't belong to this user) currently still returns success to the client.

- [ ] **Fix the GORM tag that doesn't match your real schema.**
  `internal/models/order.go` — `CartItem.CartID` is tagged `gorm:"uniqueIndex;not null"`. Your
  actual migration (`db/migrations/007_create_cart_items_table.up.sql`) correctly has
  `UNIQUE(cart_id, product_id)`. If anyone ever runs `AutoMigrate` against this model, it will
  try to create a *single-column* unique index on `cart_id` — which would mean **one cart item
  per cart, ever**. Fix the tag to a composite unique index, or explicitly document that
  AutoMigrate must never run in this project (migrations are the source of truth) and remove
  the misleading tag.

- [ ] **Fix invalid GORM tag syntax**: `internal/models/user.go` (and `RefreshToken`) use
  `gorm:"primary key"` (with a space) instead of `gorm:"primaryKey"`. It currently "works" only
  because GORM's naming convention auto-detects a field named `ID` as PK regardless of the tag
  — but it's a landmine for the next person who renames that field or copies the pattern
  elsewhere. Fix for consistency with `Product`/`Category`, which use it correctly.

- [ ] **Money fields are `float64`.** `internal/models/product.go`, `internal/dto/product.go`,
  `internal/dto/order.go` (`Price`, `TotalAmount`, etc.) use `float64` for currency. Floating
  point cannot represent money exactly and will eventually produce off-by-a-cent totals. Move to
  integer minor units (`int64` cents) or `shopspring/decimal` end-to-end (model → DTO → DB column).

- [ ] **Fix the validator tag typo.** `internal/dto/product.go` has `binding:"required, min=1"`
  — note the space after the comma. `go-playground/validator` tag parsing is comma-delimited
  with no spaces; this can silently break the `min=1` rule. Should be `binding:"required,min=1"`.

---

## Phase 3 — Security hardening

- [ ] **Pin the JWT signing algorithm on verification.** `internal/utils/jwt.go`'s
  `ValidateToken` doesn't restrict which algorithm is accepted. Add
  `jwt.WithValidMethods([]string{"HS256"})` (or check `token.Method` in the keyfunc) so a
  crafted token can't try to switch algorithms.

- [ ] **Don't leak internal error strings to API clients.** Several handlers pass raw
  `err.Error()` (which can be a raw GORM/Postgres error) straight into the JSON response body
  via `BadRequestResponse`/`InternalServerErrorResponse`. Log the full error server-side (with
  request ID), return a generic client-safe message.

- [ ] **Gate the GraphQL Playground and introspection behind an env check.**
  `internal/server/graphql.go` / `server.go` — `/playground`, `/playground/public`,
  `/playground/protected` are always mounted. Only register them when
  `cfg.Env != "production"`, and consider disabling introspection in prod too.

- [ ] **Enforce upload size limits server-side.** `MAX_UPLOAD_SIZE` exists in `.env.example` but
  isn't read anywhere in `internal/services/upload_service.go` or the handler. Wire it into
  `router.MaxMultipartMemory` and/or check `file.Size` before accepting the upload.

- [ ] **Validate file content, not just extension.** `upload_service.go` whitelists extensions
  (good) but doesn't check magic bytes/MIME sniffing — a file can be renamed to bypass the
  extension check. Low severity here since nothing executes the file, but cheap to add
  (`http.DetectContentType` on the first 512 bytes).

- [ ] **Tighten CORS for production.** `internal/server/server.go` sets
  `Access-Control-Allow-Origin: *` unconditionally. Make it a configured allow-list of real
  frontend origins in prod; also fix the trailing comma in the `Access-Control-Allow-Headers`
  value.

---

## Phase 4 — Architecture consistency

- [ ] **Pick one data-access pattern and use it everywhere.** `UserRepository`/`CartRepository`
  go through the `repository` interfaces; `ProductService`/`OrderService` hit `*gorm.DB`
  directly. Both are legitimate patterns on their own, but mixing them means testing and
  transaction handling can't be done consistently. Recommend: standardize on the repository
  interface pattern (you've already built it) so every service is unit-testable against a mock.

- [ ] **Thread `context.Context` through repositories and services.** None of the repository
  interfaces in `internal/repository/interfaces.go` take a `ctx`. Add it to every method
  signature (`GetByID(ctx context.Context, id uint) (...)`) so request cancellation, timeouts,
  and tracing propagate down to the DB layer — needed for #5 below too.

- [ ] **Move DB connection tuning into `database.New`.** `internal/database/database.go` opens
  a connection with no `SetMaxOpenConns`/`SetMaxIdleConns`/`SetConnMaxLifetime`, and logs every
  query at `logger.Info` unconditionally (fine for dev, noisy/slow in prod). Make the GORM log
  level driven by `cfg.Env`, and set sane pool limits.

- [ ] **Rename `abc` in `database.New`.** Cosmetic but real — it's the DSN string; call it `dsn`.

---

## Phase 5 — Observability & operability

- [ ] **Use the zerolog logger you already built, everywhere.** Several places
  (`auth_service.go`'s `fmt.Println`, `log.Println` for refresh-token cleanup errors) bypass the
  structured logger entirely, and `gin.Logger()` (plain text) is used for HTTP request logging
  instead of a zerolog-based middleware. Structured, leveled logs everywhere or you'll fight your
  own log pipeline in prod.

- [ ] **Make `/health` actually check dependencies.** Currently it likely returns a static OK.
  Have it ping the DB (and optionally SQS) with a short timeout, and separate `/healthz`
  (liveness) from `/readyz` (readiness, checks DB) if you're deploying to k8s/ECS.

- [ ] **Add a request-ID / correlation-ID middleware** and put it in the logger context, so you
  can trace one request across API logs, GORM query logs, and the notifier's SQS consumer logs.

- [ ] **Add basic metrics** (request count/latency by route+status, DB pool stats, SQS
  publish/consume success-failure counts) — Prometheus client + `/metrics` endpoint is the
  standard here given your stack.

---

## Phase 6 — Testing

- [ ] **There are currently zero `_test.go` files in the repo.** For a "production-grade"
  bar, at minimum:
  - Unit tests for `services/*` (mock the repository interfaces from Phase 4).
  - Unit tests for `utils/jwt.go` and `utils/password.go` (cheap, high value, would have caught
    the refresh-token bug immediately).
  - Integration tests for `order_service.go`'s transactional order-creation path, including a
    concurrent-request test that proves the stock-locking fix in Phase 2 actually prevents
    overselling.
  - A basic handler-level test suite using `httptest` for the REST routes, to catch route
    registration mistakes like the `/orders/orders` bug before they ship.

---

## Phase 7 — CI/CD & deployment

- [ ] **No `.github/workflows` exists.** Add a CI pipeline that runs on every PR: `go build ./...`,
  `go vet ./...`, `golangci-lint run` (you already reference it in the `Makefile`), and
  `go test ./... -race -cover`.
- [ ] **Add a `go.sum`/dependency vulnerability check** (`govulncheck`) to CI.
- [ ] **Terraform/CDK for the AWS resources** from Phase 1, applied via CI/CD on merge to main,
  so infra changes are reviewed the same way code changes are.
- [ ] Dockerfile itself is in decent shape already (multi-stage, non-root user, healthcheck) —
  just make sure the AWS-mode env vars (no static creds) are what actually gets injected in the
  deployed task definition/pod spec, via a role, not baked into the image or a checked-in `.env`.

---

## Your DTO question — short answer

**Keep the DTO layer as-is, structurally.** It's genuinely one of the stronger parts of this
codebase: request/response DTOs are cleanly separated from GORM models, and — notably — the
same DTOs are reused by both the REST handlers (`internal/server/*_handlers.go`) and the
GraphQL resolvers (`graph/resolver/schema.resolvers.go`), so business-shape logic isn't
duplicated across your two API surfaces. Pointer fields for optional booleans (`IsActive *bool`)
are used correctly to distinguish "not provided" from "false."

The one change needed isn't structural — it's the underlying type: switch money fields from
`float64` to integer minor-units or `decimal.Decimal` (Phase 2 item above). Do that in the model
and DTO together so there's no float↔decimal conversion drift at the boundary.