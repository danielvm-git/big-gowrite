package errors

import "net/http"

// catalogEntry is a known error code's default HTTP status and public message.
type catalogEntry struct {
	status  int
	message string
}

// catalog maps the subset of legacy Appwrite error types the migrated contexts
// use to their HTTP status + default message. It preserves the legacy `type`
// strings so the console SPA's error handling keeps working (DEV-004).
//
// Source: legacy app/config/errors.php (typed subset for the 11 contexts).
var catalog = map[string]catalogEntry{
	// general
	"general_unknown":             {http.StatusInternalServerError, "An unknown error occurred."},
	"general_rate_limit_exceeded": {http.StatusTooManyRequests, "Rate limit exceeded. Try again later."},
	"general_unauthorized_scope":  {http.StatusUnauthorized, "App or user is not authorized to take this action."},
	"general_access_forbidden":    {http.StatusForbidden, "Access to this resource is forbidden."},
	"general_argument_invalid":    {http.StatusBadRequest, "Invalid argument provided."},
	"general_not_found":           {http.StatusNotFound, "The requested resource could not be found."},

	// user / identity
	"user_not_found":             {http.StatusNotFound, "User with the requested ID could not be found."},
	"user_already_exists":        {http.StatusConflict, "A user with the same id, email, or phone already exists."},
	"user_invalid_credentials":   {http.StatusUnauthorized, "Invalid credentials. Check the email and password."},
	"user_blocked":               {http.StatusUnauthorized, "The current user has been blocked."},
	"user_more_factors_required": {http.StatusUnauthorized, "More factors are required to complete the sign in process."},

	// team / organization
	"team_not_found":       {http.StatusNotFound, "Team with the requested ID could not be found."},
	"membership_not_found": {http.StatusNotFound, "Membership with the requested ID could not be found."},
	"team_owner_required":  {http.StatusBadRequest, "A team must have at least one owner."},

	// project
	"project_not_found": {http.StatusNotFound, "Project with the requested ID could not be found."},

	// storage / document
	"document_not_found":       {http.StatusNotFound, "Document with the requested ID could not be found."},
	"document_update_conflict": {http.StatusConflict, "Remote document is newer than local."},
	"storage_file_not_found":   {http.StatusNotFound, "The requested file could not be found."},
	"storage_device_not_found": {http.StatusNotFound, "The requested storage device could not be found."},
}

// Codes returns all known catalog codes (test/introspection helper).
func Codes() []string {
	out := make([]string, 0, len(catalog))
	for k := range catalog {
		out = append(out, k)
	}
	return out
}
