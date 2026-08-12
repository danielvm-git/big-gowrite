// Package permissions implements gowrite's authorization model: console-role
// scope checks (from roles.php, BR-MIGRAR-008) and document-level ACL
// (Role::any/users/guests/user/team, BR-MIGRAR-007).
package permissions

import "strings"

// Principal is the resolved caller identity used for authorization.
type Principal struct {
	// UserID is the authenticated user's ID, empty for guests.
	UserID string
	// Roles are ACL roles the principal holds, e.g. "any", "users",
	// "user:<id>", "team:<id>", "team:<id>/<role>".
	Roles []string
	// Scopes is the effective console scope set (union of the principal's
	// console roles' scopes, plus any granted directly).
	Scopes map[string]struct{}
}

// ScopesForRole returns the scope set for a console role (nil if unknown).
func ScopesForRole(role string) []string {
	return roleScopes[role]
}

// NewPrincipal builds a Principal whose Scopes are the union of the given
// console roles' scope sets.
func NewPrincipal(userID string, consoleRoles []string, aclRoles []string) Principal {
	scopes := map[string]struct{}{}
	for _, r := range consoleRoles {
		for _, s := range roleScopes[r] {
			scopes[s] = struct{}{}
		}
	}
	return Principal{UserID: userID, Roles: aclRoles, Scopes: scopes}
}

// Can reports whether the principal holds the given console scope
// (e.g. "documents.write"). Route handlers gate on their scope label.
func (p Principal) Can(scope string) bool {
	_, ok := p.Scopes[scope]
	return ok
}

// HasRole reports whether the principal holds an exact ACL role.
func (p Principal) HasRole(role string) bool {
	for _, r := range p.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// Op is a document access operation.
type Op string

const (
	OpRead   Op = "read"
	OpCreate Op = "create"
	OpUpdate Op = "update"
	OpDelete Op = "delete"
)

// Document is the minimal shape needed for ACL checks: its per-op permission
// strings (e.g. "read(\"any\")", "update(\"user:abc\")").
type Document struct {
	ID          string
	Permissions []string
}

// permission parses one permission string into (op, role).
// "read(\"team:abc/owner\")" -> ("read", "team:abc/owner").
func parsePermission(s string) (op string, role string, ok bool) {
	s = strings.TrimSpace(s)
	open := strings.IndexByte(s, '(')
	if open < 0 || !strings.HasSuffix(s, ")") {
		return "", "", false
	}
	op = strings.TrimSpace(s[:open])
	inner := s[open+1 : len(s)-1]
	inner = strings.Trim(strings.TrimSpace(inner), `"`)
	if op == "" || inner == "" {
		return "", "", false
	}
	return op, inner, true
}

// roleMatches reports whether a principal ACL role satisfies a permission's
// required role. A team membership role "team:abc/owner" satisfies both
// "team:abc/owner" and the broader "team:abc".
func (p Principal) roleMatches(required string) bool {
	for _, held := range p.Roles {
		if held == required {
			return true
		}
		// "team:abc/owner" (held) satisfies "team:abc" (required)
		if slash := strings.IndexByte(held, '/'); slash > 0 {
			if held[:slash] == required {
				return true
			}
		}
	}
	return false
}

// Allows reports whether the principal may perform op on doc per its ACL.
// "write" permissions grant create/update/delete (legacy aggregate), in
// addition to any op-specific permission.
func (p Principal) Allows(doc Document, op Op) bool {
	want := string(op)
	for _, perm := range doc.Permissions {
		pop, role, ok := parsePermission(perm)
		if !ok {
			continue
		}
		if pop != want && !(pop == "write" && (op == OpCreate || op == OpUpdate || op == OpDelete)) {
			continue
		}
		if p.roleMatches(role) {
			return true
		}
	}
	return false
}

// FilterReadable returns the subset of docs the principal may read — the ACL
// query filter applied to a result set (BR-MIGRAR-007).
func (p Principal) FilterReadable(docs []Document) []Document {
	out := make([]Document, 0, len(docs))
	for _, d := range docs {
		if p.Allows(d, OpRead) {
			out = append(out, d)
		}
	}
	return out
}
