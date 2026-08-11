# Design: Teams

> Unit: `Teams`
> Date: 2026-08-04
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Architecture overview

The Teams module is an **HTTP collaboration service** over the project DB (`teams`, `memberships`, `users`, `sessions`). Invitations fan out through mail/messaging queues; confirmation creates sessions; deletion cascades synchronously (memberships) and asynchronously (team projects on console).

```
Client
  │ scopes: teams.*
  ▼
app/http.php → Teams\Http\** (13 actions)
  │
  ├── dbForProject (teams, memberships, users, sessions)
  ├── publisherForMails (email invites)
  ├── publisherForMessaging (SMS invites)
  ├── publisherForDeletes (async project deletion)
  ├── proofForToken / proofForPassword
  └── store / geoRecord (session creation)
```

## 2. Route-to-action map

| Method | Path | Action ID | Scope |
|--------|------|-----------|-------|
| POST | `/v1/teams` | `createTeam` | teams.write |
| GET | `/v1/teams/:teamId` | `getTeam` | teams.read |
| GET | `/v1/teams` | `listTeams` | teams.read |
| DELETE | `/v1/teams/:teamId` | `deleteTeam` | teams.write |
| PUT | `/v1/teams/:teamId` | `updateTeamName` | teams.write |
| GET | `/v1/teams/:teamId/prefs` | `getTeamPreferences` | teams.read |
| PUT | `/v1/teams/:teamId/prefs` | `updateTeamPreferences` | teams.write |
| POST | `/v1/teams/:teamId/memberships` | `createTeamMembership` | teams.write |
| GET | `/v1/teams/:teamId/memberships/:membershipId` | `getTeamMembership` | teams.read |
| GET | `/v1/teams/:teamId/memberships` | `listTeamMemberships` | teams.read |
| PATCH | `/v1/teams/:teamId/memberships/:membershipId` | `updateTeamMembership` | teams.write |
| DELETE | `/v1/teams/:teamId/memberships/:membershipId` | `deleteTeamMembership` | teams.write |
| PATCH | `/v1/teams/:teamId/memberships/:membershipId/status` | `updateTeamMembershipStatus` | teams.write |

## 3. Data flows — key actions

### 3.1 Create team

```
1. Validate name
2. Create team doc (total=0)
3. User mode → create owner membership (confirm=true), total=1
4. Emit teams.[teamId].create
```

### 3.2 Membership create (invite/direct)

```
1. Require owner/privileged/app-key
2. Resolve target: userId | email | phone (consistency check)
3. Existing user? else create stub (random hash, passwordUpdate=null) with email validation gates
4. Privileged/app mode → confirmed membership + total++
5. User mode → unconfirmed + email/SMS invite (SMTP/SMS provider required)
6. Emit memberships.create
```

### 3.3 Status confirmation

```
1. Validate invite secret (proofForToken) → TEAM_INVALID_SECRET on mismatch
2. If not logged in → create session (cookie / cross-domain header)
3. Mark confirmed; total++; emit memberships.update.status
```

### 3.4 Membership delete / role update

```
1. Owner guards (console: last owner prohibited)
2. Delete → if primary removed → promote next confirmed member
3. Update → replace roles; owner downgrade guard
```

### 3.5 Team delete

```
1. Sync-delete memberships
2. Console project → enqueue DELETE_TYPE_TEAM_PROJECTS
3. Emit teams.[teamId].delete
```

## 4. Dependency injection

| Injection | Role |
|-----------|------|
| `dbForProject` | teams, memberships, users, sessions |
| `getProjectDB` | Delete cascade |
| `queueForEvents` | Events |
| `authorization` | Role/scope checks |
| `user`, `project` | Context |
| `publisherForMails` | Email invites |
| `publisherForMessaging` | SMS invites |
| `publisherForDeletes` | Async project deletion |
| `locale`, `platform` | Email i18n |
| `timelimit` | Rate limiting |
| `usage` | Metrics |
| `plan` | Email validation flags |
| `proofForPassword`, `proofForToken` | Hashing |
| `store`, `geoRecord` | Sessions |

## 5. Events emitted

- `teams.[teamId].{create,update,update.prefs,delete}`
- `teams.[teamId].memberships.[membershipId].{create,update,update.status,delete}`

## 6. Error handling

| Error | HTTP | Trigger |
|-------|------|---------|
| `TEAM_INVALID_SECRET` | 400 | Invite secret mismatch |
| `MEMBERSHIP_DELETION_PROHIBITED` | 403 | Last owner removal (console) |
| `MEMBERSHIP_DOWNGRADE_PROHIBITED` | 403 | Last owner downgrade (console) |
| SMTP/SMS missing | 400 | Invite channel unconfigured |

## 7. Security

- Owner/privileged gating on invitations.
- Privacy settings on membership fields.
- Last-owner protection in console.
