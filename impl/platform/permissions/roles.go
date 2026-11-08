package permissions

// Console role scope sets, generated verbatim from the legacy
// app/config/roles.php (BR-MIGRAR-008). Invariants preserved from source:
//   admin == developer  ·  owner == users(member) ∪ admin
// Do not hand-edit; regenerate from roles.php.

// Console role identifiers.
const (
	RoleGuests    = "guests"
	RoleUsers     = "users"
	RoleAdmin     = "admin"
	RoleDeveloper = "developer"
	RoleOwner     = "owner"
	RoleKeys      = "keys"
)

// roleScopes maps each console role to its exact scope set.
var roleScopes = map[string][]string{
	RoleGuests: {
		"global", "public", "home", "console", "graphql", "sessions.write",
		"documents.read", "documents.write", "rows.read", "rows.write", "files.read", "files.write",
		"locale.read", "avatars.read", "executions.write",
	},
	RoleUsers: {
		"global", "public", "home", "console", "graphql", "sessions.write",
		"account", "teams.read", "teams.write", "presences.read", "presences.write", "documents.read",
		"documents.write", "rows.read", "rows.write", "files.read", "files.write", "projects.read",
		"locale.read", "avatars.read", "executions.read", "executions.write", "targets.read", "targets.write",
		"subscribers.write", "subscribers.read", "assistant.read", "rules.read",
	},
	RoleAdmin: {
		"global", "graphql", "sessions.write", "teams.read", "teams.write", "documents.read",
		"documents.write", "rows.read", "rows.write", "files.read", "files.write", "buckets.read",
		"buckets.write", "users.read", "users.write", "presences.read", "presences.write", "databases.read",
		"databases.write", "collections.read", "collections.write", "tables.read", "tables.write", "platforms.read",
		"platforms.write", "mocks.read", "mocks.write", "project.policies.read", "project.policies.write", "project.oauth2.read",
		"project.oauth2.write", "templates.read", "templates.write", "projects.write", "keys.read", "keys.write",
		"devKeys.read", "devKeys.write", "webhooks.read", "webhooks.write", "project.read", "project.write",
		"locale.read", "avatars.read", "health.read", "functions.read", "functions.write", "sites.read",
		"sites.write", "log.read", "log.write", "executions.read", "executions.write", "rules.read",
		"rules.write", "migrations.read", "migrations.write", "vcs.read", "vcs.write", "targets.read",
		"targets.write", "providers.write", "providers.read", "messages.write", "messages.read", "topics.write",
		"topics.read", "subscribers.write", "subscribers.read", "tokens.read", "tokens.write", "schedules.read",
		"schedules.write", "stages.read", "stages.write", "insights.read", "insights.write", "reports.read",
		"reports.write",
	},
	RoleDeveloper: {
		"global", "graphql", "sessions.write", "teams.read", "teams.write", "documents.read",
		"documents.write", "rows.read", "rows.write", "files.read", "files.write", "buckets.read",
		"buckets.write", "users.read", "users.write", "presences.read", "presences.write", "databases.read",
		"databases.write", "collections.read", "collections.write", "tables.read", "tables.write", "platforms.read",
		"platforms.write", "mocks.read", "mocks.write", "project.policies.read", "project.policies.write", "project.oauth2.read",
		"project.oauth2.write", "templates.read", "templates.write", "projects.write", "keys.read", "keys.write",
		"devKeys.read", "devKeys.write", "webhooks.read", "webhooks.write", "project.read", "project.write",
		"locale.read", "avatars.read", "health.read", "functions.read", "functions.write", "sites.read",
		"sites.write", "log.read", "log.write", "executions.read", "executions.write", "rules.read",
		"rules.write", "migrations.read", "migrations.write", "vcs.read", "vcs.write", "targets.read",
		"targets.write", "providers.write", "providers.read", "messages.write", "messages.read", "topics.write",
		"topics.read", "subscribers.write", "subscribers.read", "tokens.read", "tokens.write", "schedules.read",
		"schedules.write", "stages.read", "stages.write", "insights.read", "insights.write", "reports.read",
		"reports.write",
	},
	RoleOwner: {
		"global", "public", "home", "console", "graphql", "sessions.write",
		"account", "teams.read", "teams.write", "presences.read", "presences.write", "documents.read",
		"documents.write", "rows.read", "rows.write", "files.read", "files.write", "projects.read",
		"locale.read", "avatars.read", "executions.read", "executions.write", "targets.read", "targets.write",
		"subscribers.write", "subscribers.read", "assistant.read", "rules.read", "buckets.read", "buckets.write",
		"users.read", "users.write", "databases.read", "databases.write", "collections.read", "collections.write",
		"tables.read", "tables.write", "platforms.read", "platforms.write", "mocks.read", "mocks.write",
		"project.policies.read", "project.policies.write", "project.oauth2.read", "project.oauth2.write", "templates.read", "templates.write",
		"projects.write", "keys.read", "keys.write", "devKeys.read", "devKeys.write", "webhooks.read",
		"webhooks.write", "project.read", "project.write", "health.read", "functions.read", "functions.write",
		"sites.read", "sites.write", "log.read", "log.write", "rules.write", "migrations.read",
		"migrations.write", "vcs.read", "vcs.write", "providers.write", "providers.read", "messages.write",
		"messages.read", "topics.write", "topics.read", "tokens.read", "tokens.write", "schedules.read",
		"schedules.write", "stages.read", "stages.write", "insights.read", "insights.write", "reports.read",
		"reports.write",
	},
	RoleKeys: {
		"global", "health.read", "graphql",
	},
}
