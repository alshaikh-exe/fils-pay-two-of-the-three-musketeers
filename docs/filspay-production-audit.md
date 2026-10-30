# FilsPay Production Audit Report

## Section A: Executive Summary

### Overall Application Health
Filspay is a functional monorepo payment platform with a TypeScript/Fastify API, a Go settlement service, TypeScript admin/web frontends, and shared packages for money and events. The codebase is well-structured with clear module boundaries, integer-based money handling, and comprehensive test coverage. However, several confirmed bugs and deployment blockers prevent production readiness on Dilmune Cloud.

### Main Defects Discovered

| ID | Severity | Description |
|----|----------|-------------|
| BUG-001 | High | `parsePromotionalDiscountFils` uses `parseFloat` with x100 multiplier, causing 10x error on BHD discount amounts |
| BUG-002 | High | `NetPayableFils` uses `float32` arithmetic on monetary values, risking rounding errors on net payout |
| BUG-003 | High | `inboundVerify.ts` uses `JSON.stringify` instead of `canonicalJson`, causing potential signature mismatches with acquirer callbacks |
| BUG-004 | Medium | `dayStartIso` uses UTC day boundaries instead of Bahrain, making transaction date filters 3 hours off |
| BUG-005 | Medium | `RebuildAllBalances` doesn't lock `merchant_balances` rows, risking cache corruption from concurrent CloseBatch |
| BUG-006 | Medium | `refund.succeeded` event enqueued by `processRefund.ts` but not handled by `buildEnvelope` — dispatcher crash risk when refunds ship |
| BUG-007 | High | No Dockerfiles in repository — `ops/docker-compose.yml` references `ops/Dockerfile.node` and `ops/Dockerfile.go` which don't exist |
| BUG-008 | High | `WEBHOOK_SECRET_ENCRYPTION_KEY` required for production boot but not provided in compose config |

### Changes Implemented
None yet — fixes are being implemented in Phase 5.

### Current Deployment-Readiness Assessment
**Not ready.** Critical deployment blockers (missing Dockerfiles, missing production encryption key) and confirmed bugs (money calculation with floats, signature verification mismatch, discount parsing error) must be resolved first.

### Major Unresolved Risks
- Missing Dockerfiles prevent containerized deployment
- `WEBHOOK_SECRET_ENCRYPTION_KEY` must be provisioned for production
- Latent crash risk when `refund.succeeded` webhooks ship in v1.9 without `buildEnvelope` support

## Section B: Documentation Inventory

| File | Reviewed | Key Requirements/Rules | Conflicts/Ambiguities |
|------|----------|----------------------|----------------------|
| `docs/ARCHITECTURE.md` | Yes | Money = integer fils, settlement day = Asia/Bahrain, outbox pattern | "Code wins over docs" |
| `docs/HANDOVER.md` | Yes | Refunds stubbed at 501, migration gaps 0039/0041 abandoned | `amount_3dp` said gone but migration 0038 exists |
| `docs/ADMIN.md` | Yes | Admin roles support < finance < super_admin | Adjust balance journaled claim vs code |
| `docs/RUNBOOK.md` | Yes | readyz boot-only, `settlement rebuild-balances` for wrong balances | FP-215 readyz change unclear |
| `docs/adr/0001-integer-fils.md` | Yes | Integer fils everywhere, `parseBhdToFils` for all parsing | N/A |
| `docs/adr/0002-scheme-residue-at-export.md` | Yes | Residue posted per line at export | Auto-posting not done yet |
| `docs/adr/0003-report-date-windows.md` | Yes | Daily = Bahrain, monthly = UTC calendar month | YA disagrees with UTC months |
| `docs/adr/0004-balance-rebuild-source.md` | Yes | Console uses settled data, CLI uses ledger walk | Two paths may disagree |
| `docs/adr/0005-settlement-v2-shadow-ledger.md` | Yes | Paused, migrations 0054-0056 reserved | N/A |
| `docs/adr/0006-display-adjustments.md` | Yes | Adjustments are cache-only | N/A |
| `AGENTS.md` | Yes | Don't modify `packages/money/`, use `parseBhdToFils` | N/A |
| `CLAUDE.md` | Yes | Go settlement is reference for fees | N/A |
| `.cursorrules` | Yes | packages/money frozen, float32 in batch.go intentional | payable.go also uses float32 |
| `CHANGELOG.md` | Yes | v1.9.2 latest | N/A |
| `ops/README.md` | Yes | Docker compose stack, WEBHOOK_SECRET_ENCRYPTION_KEY needed | Dockerfiles missing |
| `ops/.env.docker` | Yes | NODE_ENV=production, sandbox secret set | Contains secret value |

### Files that could not be read
None.