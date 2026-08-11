# Contracts: Teams

> Unit: `Teams`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## HTTP contract

Base: `/v1/teams`. Scopes: `teams.*` (session, privileged, or API key).

| Method | Path | Key params | Success | Errors |
|--------|------|-----------|---------|--------|
| POST | `/teams` | `teamId`, `name`, `roles` (user mode) | 201 team | GENERAL_ARGUMENT_INVALID |
| GET | `/teams/:teamId` | — | 200 team | TEAM_NOT_FOUND |
| GET | `/teams` | `search`, queries | 200 teamList | — |
| PUT | `/teams/:teamId` | `name` | 200 team | TEAM_NOT_FOUND |
| DELETE | `/teams/:teamId` | — | 204 | TEAM_NOT_FOUND |
| GET | `/teams/:teamId/prefs` | — | 200 prefs | TEAM_NOT_FOUND |
| PUT | `/teams/:teamId/prefs` | `prefs` (arbitrary JSON) | 200 prefs | TEAM_NOT_FOUND |
| POST | `/teams/:teamId/memberships` | `roles`, `userId` XOR `email` XOR `phone`, `url` | 201 membership | SMTP/SMS missing, email validation gates |
| GET | `/teams/:teamId/memberships/:membershipId` | — | 200 membership | 404 |
| GET | `/teams/:teamId/memberships` | `search`, queries | 200 membershipList | — |
| PATCH | `/teams/:teamId/memberships/:membershipId` | `roles` | 200 membership | MEMBERSHIP_DOWNGRADE_PROHIBITED |
| DELETE | `/teams/:teamId/memberships/:membershipId` | — | 204 | MEMBERSHIP_DELETION_PROHIBITED |
| PATCH | `/teams/:teamId/memberships/:membershipId/status` | `secret` | 200 membership | TEAM_INVALID_SECRET |

## Invitation contract

| Aspect | Detail |
|--------|--------|
| Direct add (privileged/app) | Membership confirmed, `joined` set, `total` incremented |
| Invite (user mode) | Unconfirmed membership + email/SMS invite carrying `secret` |
| Stub user | Random password hash, `passwordUpdate=null`, email validation gates (disposable/canonical/free/corporate) |
| Confirmation | `proofForToken` secret validation; mismatch → `TEAM_INVALID_SECRET` |

## Session contract (status confirmation)

- Not-logged-in invitee → new session created via `store`/`geoRecord`; set as cookie, fallback header for cross-domain.

## Queue contracts

| Channel | Direction | Payload |
|---------|-----------|---------|
| `publisherForMails` | Http → worker-mails | Email invite: recipient, team name, invite URL + secret |
| `publisherForMessaging` | Http → worker-messaging | SMS invite: phone, team name, short URL |
| `publisherForDeletes` | Team delete → worker-deletes | `DELETE_TYPE_TEAM_PROJECTS` (console project only) |

## Event contract

- `teams.[teamId].{create,update,update.prefs,delete}`
- `teams.[teamId].memberships.[membershipId].{create,update,update.status,delete}`

## Privacy contract

| Setting | Effect |
|---------|--------|
| `membershipsPrivacy.userName/userEmail/userId/userPhone/mfa` | Fields omitted from end-user responses; privileged + API key see all |

## Owner guard contract

- Console project: removing or downgrading the **last owner** → `MEMBERSHIP_DELETION_PROHIBITED` / `MEMBERSHIP_DOWNGRADE_PROHIBITED`.
- Removing primary team user (`team.userId`) promotes the next confirmed member; `total` decrements on confirmed-member removal.
