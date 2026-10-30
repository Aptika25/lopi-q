package http

import (
	"context"
	"net/http"
	"strings"

	"api-gateway/internal/client"
	userProto "proto/user"
)

var profileLookup = &client.UserClientDirectStub{}

func requestRole(r *http.Request) string {
	role, _ := r.Context().Value("role").(string)
	return role
}

func isSuperAdmin(r *http.Request) bool {
	return requestRole(r) == "superadmin"
}

func isInternRole(role string) bool {
	role = strings.ToLower(strings.TrimSpace(role))
	return role == "intern" || role == "peserta"
}

// userRoleByID returns the stored role of an account, or ok=false when it does not exist.
func userRoleByID(ctx context.Context, id int32) (string, bool) {
	res, _ := profileLookup.GetProfile(ctx, &userProto.GetProfileRequest{UserId: id})
	if res == nil || !res.Success || res.User == nil {
		return "", false
	}
	return res.User.Role, true
}

// logAdminAction records an audit entry attributed to the authenticated caller.
func logAdminAction(r *http.Request, action, details string) {
	userID, _ := r.Context().Value("user_id").(int)
	nip, _ := r.Context().Value("nip").(string)
	name := ""
	if res, _ := profileLookup.GetProfile(r.Context(), &userProto.GetProfileRequest{UserId: int32(userID)}); res != nil && res.Success && res.User != nil {
		name = res.User.Name
	}
	client.RecordActivityLog(userID, nip, name, action, details, client.GetClientIP(r), r.UserAgent())
}
