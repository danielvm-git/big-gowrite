package permissions

import "testing"

func TestRoleInvariantsFromRolesPHP(t *testing.T) {
	set := func(ss []string) map[string]struct{} {
		m := map[string]struct{}{}
		for _, s := range ss {
			m[s] = struct{}{}
		}
		return m
	}
	eq := func(a, b map[string]struct{}) bool {
		if len(a) != len(b) {
			return false
		}
		for k := range a {
			if _, ok := b[k]; !ok {
				return false
			}
		}
		return true
	}

	// admin ≡ developer
	if !eq(set(roleScopes[RoleAdmin]), set(roleScopes[RoleDeveloper])) {
		t.Fatal("admin and developer scope sets must be identical (roles.php)")
	}
	// owner == member(users) ∪ admins
	union := set(roleScopes[RoleUsers])
	for _, s := range roleScopes[RoleAdmin] {
		union[s] = struct{}{}
	}
	if !eq(set(roleScopes[RoleOwner]), union) {
		t.Fatal("owner scope set must equal users ∪ admin (roles.php)")
	}
	// keys is the tiny application scope set
	if len(roleScopes[RoleKeys]) != 3 {
		t.Fatalf("keys should have 3 scopes, got %d", len(roleScopes[RoleKeys]))
	}
}

func TestCanScope(t *testing.T) {
	owner := NewPrincipal("u1", []string{RoleOwner}, nil)
	if !owner.Can("projects.write") {
		t.Fatal("owner should have projects.write")
	}
	member := NewPrincipal("u2", []string{RoleUsers}, nil)
	if member.Can("projects.write") {
		t.Fatal("member should NOT have projects.write (admin-only)")
	}
	if !member.Can("teams.write") {
		t.Fatal("member should have teams.write")
	}
	guest := NewPrincipal("", []string{RoleGuests}, nil)
	if guest.Can("users.write") {
		t.Fatal("guest should not have users.write")
	}
}

func TestParsePermission(t *testing.T) {
	cases := map[string][2]string{
		`read("any")`:            {"read", "any"},
		`update("user:abc")`:     {"update", "user:abc"},
		`write("team:t1/owner")`: {"write", "team:t1/owner"},
		` delete( "users" ) `:    {"delete", "users"},
	}
	for in, want := range cases {
		op, role, ok := parsePermission(in)
		if !ok || op != want[0] || role != want[1] {
			t.Fatalf("parse %q = (%q,%q,%v), want (%q,%q,true)", in, op, role, ok, want[0], want[1])
		}
	}
	if _, _, ok := parsePermission("garbage"); ok {
		t.Fatal("garbage should not parse")
	}
}

func TestDocumentACLRead(t *testing.T) {
	doc := Document{ID: "d1", Permissions: []string{`read("user:owner1")`, `read("team:t1")`}}

	owner := NewPrincipal("owner1", nil, []string{"user:owner1", "users"})
	if !owner.Allows(doc, OpRead) {
		t.Fatal("owner should read own doc")
	}
	teamMember := NewPrincipal("m1", nil, []string{"user:m1", "team:t1/member"})
	if !teamMember.Allows(doc, OpRead) {
		t.Fatal("team member should read team doc (team:t1/member satisfies team:t1)")
	}
	stranger := NewPrincipal("x", nil, []string{"user:x", "users"})
	if stranger.Allows(doc, OpRead) {
		t.Fatal("stranger must NOT read the doc")
	}
}

func TestWriteGrantsCUD(t *testing.T) {
	doc := Document{ID: "d", Permissions: []string{`write("user:editor")`}}
	editor := NewPrincipal("editor", nil, []string{"user:editor"})
	for _, op := range []Op{OpCreate, OpUpdate, OpDelete} {
		if !editor.Allows(doc, op) {
			t.Fatalf("write permission should grant %s", op)
		}
	}
	// but not read, unless a read perm exists
	if editor.Allows(doc, OpRead) {
		t.Fatal("write should not imply read")
	}
}

func TestFilterReadable(t *testing.T) {
	docs := []Document{
		{ID: "a", Permissions: []string{`read("any")`}},
		{ID: "b", Permissions: []string{`read("user:me")`}},
		{ID: "c", Permissions: []string{`read("user:someone-else")`}},
	}
	me := NewPrincipal("me", nil, []string{"user:me", "any", "users"})
	got := me.FilterReadable(docs)
	if len(got) != 2 {
		t.Fatalf("readable = %d docs, want 2 (a via any, b via user:me)", len(got))
	}
	for _, d := range got {
		if d.ID == "c" {
			t.Fatal("doc c must be filtered out")
		}
	}
}
