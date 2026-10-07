<!--
Sync Impact Report:
- Version change: none → 1.0.0
- Modified principles: none (initial generation)
- Added sections: Core Principles (I–XVI), Technology Stack, Development Workflow, Governance
- Removed sections: Speckit template placeholders
- Follow-up TODOs:
  - TODO(PRODUCT_SPEC): no product-spec repository or pinned version is declared for this fork
  - TODO(CHANGELOG_FORMAT): CHANGELOG.md and WENI-CHANGELOG.md are not Keep a Changelog
  - TODO(STRUCTURED_LOGS): cmd/rp-indexer/main.go uses logrus TextFormatter
  - TODO(SENTRY_DSN_LOG): an invalid Sentry DSN is interpolated into the fatal log in cmd/rp-indexer/main.go
  - TODO(RETRY_5XX): ShouldRetry in indexer.go retries HTTP 429 and body-read failures, not HTTP 5xx
  - TODO(POLL_LITERAL): the main loop sleeps a fixed 5s and does not use config.Poll
  - TODO(PEAK_LOAD): peak indexing throughput is not declared in an engineering spec
  - TODO(BRANCH_PROTECTION): main-branch protection cannot be verified from the repository

Provenance:
- Source: weni-ai/vtex-cx-engineering-constitutions (main)
- Domains: backend
-->

# rp-indexer Constitution

## Core Principles

### I. Version Control and Review

All code MUST enter the default branch through a pull request. A merge MUST require at least one approved review and a green run of `.github/workflows/ci.yml`. Direct pushes to the default branch MUST be blocked by platform branch protection.

**Rationale:** the policy is only real when the platform enforces it. Peer review and a protected default branch keep the indexer history auditable and stop unreviewed indexing changes from reaching a deployment.

### II. Security and Secrets

Secrets MUST never be committed. `INDEXER_DB`, `INDEXER_ELASTIC_URL`, and `INDEXER_SENTRY_DSN` MUST be injected at runtime through the environment (or a deployment secret store) and MUST NOT be written to `indexer.toml` in this repository. Access to PostgreSQL and Elasticsearch MUST follow least privilege. Module dependencies MUST come from the Go checksum database recorded in `go.sum` and MUST be checked for known vulnerabilities.

Error logs and Sentry events MUST NOT contain connection strings, DSNs, or other secrets.

**Rationale:** leaked database URLs and error-reporting credentials are enough to expose RapidPro contact data. Prevention is far cheaper than remediation.

### III. Observability

Logs MUST be structured key/value records on stdout and MUST NOT contain secrets or contact personal data (names, e-mail addresses, phone numbers, or government identifiers). Indexing errors MUST be traceable across a cycle through a correlation identifier for that cycle, together with the physical index name.

The process MUST expose Prometheus metrics on `GET /metrics` (`INDEXER_METRICS_PORT`, default `8070`), including contact batch outcomes and database and Elasticsearch timings. When `INDEXER_SENTRY_DSN` is set, error-level failures MUST be reported to Sentry.

**Rationale:** structured, privacy-safe telemetry is what makes a stalled or partial index diagnosable without copying contact records into the log stream.

### IV. Versioned Contracts

Any change to a public interface MUST follow SemVer. The public interfaces of this service are:

- the contact document and index mapping in `contacts/index_settings.json` (embedded at build time) and the alias named by `INDEXER_INDEX` (default `contacts`);
- the configuration surface loaded by `ezconf` in `cmd/rp-indexer/main.go` (`indexer.toml`, `INDEXER_*` environment variables, and CLI flags);
- the `rp-indexer` binary and container image published for a tag.

Changes MUST be backward compatible or MUST ship with an announced deprecation path. A silent breaking change to indexed fields, the alias swap, or configuration MUST NOT be introduced. A mapping change that drops or retypes a field consumers query MUST be a major version.

**Rationale:** RapidPro search depends on a stable contact index. Explicit versioning gives consumers a predictable path to adapt without query outages.

### V. Specification Traceability

Every engineering spec MUST derive from exactly one approved product spec and MUST reference it through an immutable, pinned version (commit or tag). A mutable URL or ID alone MUST NOT be used. The product spec MUST exist and be tagged before its engineering spec is created. An engineering spec MUST NOT redefine the "what" it inherits: problem, scope, success criteria, and binding decisions belong to the product spec. A technical architecture document SHOULD be produced for non-trivial features; when it exists it MUST be linked from the engineering spec, also pinned by commit or tag, but its absence MUST NOT block the engineering spec.

Every engineering spec under `specs/` MUST open with an inheritance section in exactly this format:

```
## Inheritance from Product Spec
- Product Spec: <title> — <URL>
- Pinned version: <commit/tag>
- Architecture doc: <none | URL + commit/tag>
- Inherited binding decisions: <short list>
- Scope of this spec: <slice implemented by this repo>
- Divergences: <none | link to amendment>
```

**Rationale:** traceability from product intent to the indexing behavior keeps decisions auditable. Pinning the version is what guarantees every deployment implements the same feature instead of divergent readings of a spec that changed mid-flight.

### VI. No Silent Divergence

When a technical need contradicts something inherited from the product spec — scope, success criteria, or a binding decision — the divergence MUST NOT be implemented silently in `indexer.go` or `cmd/rp-indexer`. It MUST be raised as an amendment in the product repository and recorded in the `Divergences` field of the engineering spec's inheritance section, linking to that amendment. Once the amendment is approved and produces a new tag, the engineering spec's `Pinned version` MUST be updated to it. A technical difference that contradicts nothing inherited is an implementation decision and MUST live in the engineering spec.

**Rationale:** a silent code deviation makes the contact-index contract and the product intent drift apart with no audit trail. Amendments keep the spec authoritative.

### VII. Commit Messages

Commits MUST follow Conventional Commits: `<type>: <description>`. Allowed types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`. The description MUST be imperative, specific, and no longer than 50 characters. Commits MUST be atomic: one logical change per commit.

**Rationale:** conventional commits enable changelog generation and semantic versioning. Atomic commits simplify bisecting, reverting, and reviewing indexing changes.

### VIII. Changelog Maintenance

This repository publishes a versioned `rp-indexer` binary (GoReleaser) and container image. Every user-facing change MUST appear in `WENI-CHANGELOG.md` under Added, Changed, Deprecated, Removed, Fixed, or Security, following Keep a Changelog. Version bumps MUST follow SemVer. `CHANGELOG.md` is the upstream history and MUST NOT be rewritten to simulate that format.

**Rationale:** operators of the indexer need release notes that state impact. SemVer alignment keeps upgrades predictable. The upstream file stays intact so this fork does not rewrite history it did not author.

### IX. Never Trust the Client

Every external input MUST be validated for type, format, range, and business rules at the boundary before use. For this service those inputs are PostgreSQL rows from `contacts_contact`, Elasticsearch responses, and configuration from `ezconf`. Contact fields MUST be validated before they are written into the index. Authorization for contact data MUST be enforced by database and Elasticsearch credentials, not by a caller the indexer does not have.

The only inbound HTTP surface is `GET /metrics`. It MUST NOT return contact documents or personal data. If it is left unauthenticated, it MUST be reachable only from the metrics-scraper network.

**Rationale:** rows, search responses, and environment values can be incomplete or hostile. Validating them at the boundary is what prevents a bad record from corrupting the contact index. Client-side checks do not exist on this path, so the server boundary is the only check.

### X. Fail Gracefully and Predictably

Calls to PostgreSQL (`INDEXER_DB`) and Elasticsearch (`INDEXER_ELASTIC_URL`) MUST have explicit timeouts and MUST NOT block indefinitely. Failures MUST be handled explicitly. In the default polling mode, a cycle error MUST be logged and the process MUST keep running. In `--rebuild`, a cycle error MUST fail the process rather than swap the alias on a partial index. Error text returned or logged MUST NOT include secrets or raw contact personal data.

**Rationale:** Elasticsearch and PostgreSQL will fail. Containing that failure keeps one bad cycle from either killing a healthy poller or promoting an incomplete rebuild.

### XI. Bounded Retry Over REST

Propagation of contact batches to Elasticsearch MUST be retried rather than dropped when the failure could plausibly succeed on another attempt: connection errors, timeouts, HTTP 5xx, or HTTP 429. A 4xx that reflects a defect in the request MUST NOT be retried. Bulk index operations MUST stay idempotent (stable document ids, deduplicated batches). The retry policy MUST live in `ElasticRetries` / `retryConfig` in `indexer.go`, MUST define a maximum number of attempts and a backoff, and MUST NOT retry without bound. When attempts are exhausted, the failure MUST be logged and the batch MUST remain recoverable on a later cycle. It MUST NOT be silently discarded.

**Rationale:** transient Elasticsearch failures are common; dropping a batch makes the index drift from RapidPro. Unbounded or non-idempotent retries turn a degraded cluster into an outage and duplicate writes.

### XII. Scalability and Peak Load

Durable indexing progress MUST NOT be kept only in process memory or on local disk. The cursor MUST be the last-modified contact read from the Elasticsearch index, and the source of truth MUST remain PostgreSQL.

**Exception:** this service MUST run as a single instance per deployment and index alias. The incremental poll and the alias swap in rebuild mode are not safe to run concurrently, as documented in `README.md`. Horizontal scale-out MUST NOT be used as the capacity strategy. Peak load MUST still be declared in the engineering spec as peak indexing throughput (contacts per cycle, batch size, poll interval), not as a replica count.

**Rationale:** capacity is a design input. Stateless durable progress is what lets a restarted process resume from Elasticsearch. A second instance would race the same alias and cursor, so the single-instance limit is an explicit exception to scale-out, not permission to hide progress on local disk.

### XIII. Diagnosable Errors

Every error reported to Sentry MUST carry enough context to locate and filter it without reproducing it: at minimum a project identifier, the index alias, the physical index name, and the cycle correlation identifier. When the failure is tied to an organization or contact, the report MUST include only opaque organization and contact identifiers. Names, e-mail addresses, phone numbers, and government identifiers MUST NOT be attached to an error report.

**Rationale:** an error without identifying context can be counted but not investigated. Opaque identifiers support filtering without putting personal data in Sentry, which is what Observability requires.

### XIV. Tests Exercise Flows

Every indexing flow MUST have at least one test that runs the use case from database input to the resulting Elasticsearch effect. `TestIndexing` in `indexer_test.go`, backed by `testdb.sql` and `testdb_update.sql`, is the flow test for the contact index and MUST be extended when the indexed document or the rebuild/delete behavior changes. Method-level tests MAY cover edge cases, and `TestRetryServer` MUST cover the retry decision, but they MUST NOT be the only coverage of a flow. Success and failure paths MUST both be exercised. Tests MUST be runnable with `go test ./... -p=1` against PostgreSQL database `elastic_test` and Elasticsearch, matching `.github/workflows/ci.yml`.

**Rationale:** isolated tests can stay green while the composition of query, batch, and index update is broken. The flow test is what proves a contact written in PostgreSQL becomes the document search expects, including the failure paths that production hits first.

### XV. Explicit Over Clever

What a piece of code does MUST be evident where it happens. Hidden side effects and implicit control flow MUST NOT be introduced to save lines. Any literal that carries meaning — batch size, poll interval, overlap window, retry count, backoff, metrics port — MUST be a named constant or a configuration field (`config.Poll`, `INDEXER_METRICS_PORT`, `batchSize`, `ElasticRetries`), not an inline value. A literal with no meaning beyond itself, such as an index of 0, is exempt. Comments MUST explain why a decision was made. A comment that restates the code SHOULD be removed by rewriting the code.

**Rationale:** the indexer is read by people who do not have the original RapidPro context. An unexplained timeout or overlap window cannot be reviewed because its safe range is invisible.

### XVI. Contained Changes

A change MUST be limited to the context it was asked to address. Refactoring, renaming, reformatting, or behavior adjustments outside that context MUST NOT ride along; each belongs to its own change. A change that stays within scope MAY still span several atomic commits.

**Rationale:** a diff that reaches past its stated scope hides the intended fix, makes review expensive, and turns a revert into a choice between losing the fix and keeping an unrelated regression.

## Technology Stack

The service MUST stay a single Go module, `github.com/nyaruka/rp-indexer`, on the Go version declared in `go.mod` (1.23) and in `docker/Dockerfile`. The binary entrypoint MUST remain `cmd/rp-indexer`. Indexing logic MUST remain in the root `indexer` package.

Runtime dependencies MUST stay explicit:

- PostgreSQL through `github.com/lib/pq` and `INDEXER_DB`
- Elasticsearch 7 through the HTTP index API and `INDEXER_ELASTIC_URL`, with `github.com/olivere/elastic/v7` used by tests
- configuration through `github.com/nyaruka/ezconf`
- logs through `github.com/sirupsen/logrus` and Sentry through `github.com/evalphobia/logrus_sentry`
- metrics through `github.com/prometheus/client_golang` and `github.com/go-chi/chi` in `server.go`

The container image MUST be built from `docker/Dockerfile` as a non-root user. Local composition lives in `docker/docker-compose.yml`. Index settings MUST be the embedded document `contacts/index_settings.json`.

## Development Workflow

Work MUST land through a reviewed pull request whose CI is green. The required check is `.github/workflows/ci.yml`: `go test -p=1` with coverage, PostgreSQL 12 and 13, and Elasticsearch 7.10.1. A tag that starts a release MUST pass that test job before GoReleaser publishes binaries (`goreleaser.yml`). Image publish workflows under `.github/workflows/` MUST NOT replace the test job as the merge gate.

Plans MUST include a Constitution Check against this document before implementation. A conflict with a MUST found by speckit analyze MUST be treated as CRITICAL and resolved before implementation, unless this constitution records an explicit exception (see Scalability and Peak Load).

## Governance

This constitution supersedes informal practice for `rp-indexer`. Amendments MUST be made in `.specify/memory/constitution.md`, MUST state why, and MUST bump the version below. SemVer for this document: MAJOR when a principle is removed or redefined, MINOR when a principle or section is added, PATCH when wording is clarified without changing a requirement. The ratification date MUST NOT change on amendment.

Engineering specs, plans, and tasks MUST comply with these principles. Complexity beyond a single indexing cycle, one batch writer, and the existing PostgreSQL and Elasticsearch boundaries MUST be justified in the spec. Runtime guidance for contributors is this file plus `README.md`.

**Version**: 1.0.0 | **Ratified**: 2026-10-06 | **Last Amended**: 2026-10-06
