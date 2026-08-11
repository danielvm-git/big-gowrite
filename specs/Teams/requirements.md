# Requirements: Teams

> Identifier: `TMS-COLLAB`
> Date: 2026-08-04
> Source: Reverse engineering from Appwrite legacy codebase
> Confidence: 🟢 CONFIRMED · 🟡 INFERRED · 🔴 GAP

## 1. Executive summary

The Teams module provides collaboration primitives: **teams** (user groups with roles), **memberships** (user–team links with roles), **team preferences** (arbitrary JSON), and the invitation/confirmation flow (email, SMS, or direct userId). Teams live in the per-project database and are the primary access-control grouping mechanism.

## 2. Context from legacy

| Source | Relevant excerpt | Confidence |
|--------|------------------|------------|
| `specs/inventory.md#Teams` | Module purpose | 🟢 |
| `specs/Teams/legacy-mapping.md` | Endpoints, business rules, DI | 🟢 |
| `specs/domain.md#Team` | Team/membership model | 🟢 |
| `specs/security.md#RBAC` | Roles, privacy settings | 🟢 |

## 3. Personas and use scenarios

| Persona | Goal | Key scenario |
|---------|------|-------------|
| App user | Collaborate | Create team; invite members |
| App user | Manage roles | Update membership roles |
| App user | Join team | Confirm email invite |
| App developer | Access control | Team-scoped permissions |

## 4. Business rules

1. **BR-01:** Team creation in user mode auto-adds creator as owner (`confirm=true`); app/server mode skips (total starts 0). 🟢
2. **BR-02:** Membership creation in app/privileged mode adds user directly (confirmed, `joined` immediately) + increments `total`; user mode creates unconfirmed + sends email/SMS invite. 🟢
3. **BR-03:** Only owner, privileged user, or app key may send invitations. 🟢
4. **BR-04:** SMTP must be configured (`_APP_SMTP_HOST`) for non-privileged/non-app email invites. 🟢
5. **BR-05:** SMS invite requires `_APP_SMS_PROVIDER`. 🟢
6. **BR-06:** `createTeamMembership` accepts `userId`, `email`, or `phone` (≥1); `userId` must match email/phone values when co-provided. 🟢
7. **BR-07:** No matching user → stub user created with random password hash, `passwordUpdate=null`. 🟢
8. **BR-08:** Email validation rules (disposable, canonical, free, corporate) enforced on new-user creation during invite, gated by project auths + plan. 🟢
9. **BR-09:** Status confirmation validates invite secret via `proofForToken`; mismatch → `TEAM_INVALID_SECRET`. 🟢
10. **BR-10:** Not-logged-in invitee confirming → new session set as cookie (fallback header for cross-domain). 🟢
11. **BR-11:** `membershipsPrivacy` settings (UserName/Email/Id/Phone/Mfa) control exposed fields; privileged + app keys see all. 🟢
12. **BR-12:** Console project: deleting team or removing last owner blocked → `MEMBERSHIP_DELETION_PROHIBITED` / `MEMBERSHIP_DOWNGRADE_PROHIBITED`. 🟢
13. **BR-13:** Removing primary team user (`team.userId`) promotes next confirmed member. 🟢
14. **BR-14:** Team delete: memberships sync-deleted; team-owned projects async-deleted (`DELETE_TYPE_TEAM_PROJECTS`) only on console project. 🟢
15. **BR-15:** `total` counts confirmed members only; decremented on confirmed-member removal. 🟡

## 5. Functional requirements

| ID | Requirement | Priority | Acceptance criteria | Conf |
|----|-------------|----------|---------------------|------|
| RF-01 | Team CRUD + name update | Must | Creator ownership; cascade delete | 🟢 |
| RF-02 | Team preferences get/replace | Must | Arbitrary JSON; full replace | 🟢 |
| RF-03 | Membership create (invite/direct) | Must | email/phone/userId; invite vs direct; stub user | 🟢 |
| RF-04 | Membership get/list | Must | Privacy settings honored | 🟢 |
| RF-05 | Membership role update | Must | Roles replaced; owner guards | 🟢 |
| RF-06 | Membership delete | Must | Primary promotion; last-owner guard | 🟢 |
| RF-07 | Membership status confirm | Must | Secret validation; session creation | 🟢 |

## 6. Non-functional requirements

| Type | Requirement | Evidence | Conf |
|------|-------------|----------|------|
| Security | Invite secret verification | proofForToken | 🟢 |
| Privacy | Membership field gating | membershipsPrivacy settings | 🟢 |
| Governança | Email validation on invite | plan-gated disposable/canonical/free/corporate | 🟢 |
| Reliability | Primary-member promotion | team.userId reassignment | 🟢 |

## 7. Acceptance criteria

```gherkin
Scenario: User-mode team creation
  Given a session-authenticated user
  When POST /v1/teams is called
  Then the creator is added as confirmed owner and total=1

Scenario: Confirm invite with wrong secret
  Given an invite with a mismatched secret
  When PATCH .../memberships/:membershipId/status is called
  Then TEAM_INVALID_SECRET is returned

Scenario: Last owner removal blocked
  Given a console project team with one owner
  When removing that owner is attempted
  Then MEMBERSHIP_DELETION_PROHIBITED is returned

Scenario: Privacy-gated membership list
  Given membershipsPrivacy hiding emails
  When an end user lists memberships
  Then email fields are omitted
```

## 8. MoSCoW priority

| Item | MoSCoW | Justification |
|------|--------|---------------|
| RF-01..07 (core) | Must | Collaboration lifecycle |

## 9. Clarifications

> No clarification session recorded.

## 10. Gaps

- 🟡 [INFERRED] `total` counter semantics for edge cases.
- 🔴 [GAP] Exact invite email/SMS template flow.

## 11. Change history

| Date | Change | Author |
|------|--------|--------|
| 2026-08-04 | Initial version from legacy reverse engineering | reversa-writer |
