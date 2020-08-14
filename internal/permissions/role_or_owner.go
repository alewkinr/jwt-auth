package permissions

import (
	"net/http"
	"strconv"
)

// IRoleOrOwnerChecker — это интерфейс проверяльщика доступа к ресурсу по роли.
type IRoleOrOwnerChecker interface {
	Check(r *http.Request, ownerUserID int64) bool
}

// RoleOrOwnerChecker — это проверяльщик, разрешающий доступ к ресурсу владельцу и перечисленным ролям.
type RoleOrOwnerChecker struct {
	permittedRoles []Role
}

// NewRoleOrOwnerChecker создаёт проверяльщик, разрешающий доступ к ресурсу владельцу и перечисленным ролям.
func NewRoleOrOwnerChecker(permittedRoles []Role) RoleOrOwnerChecker {
	return RoleOrOwnerChecker{
		permittedRoles: permittedRoles,
	}
}

// Check сообщает, разрешён ли доступ к ресурсу.
func (c RoleOrOwnerChecker) Check(r *http.Request, ownerUserID int64) bool {
	var userID int64
	permittedRoles := map[string]struct{}{}
	for _, r := range c.permittedRoles {
		permittedRoles[string(r)] = struct{}{}
	}

	clientIDHeader := r.Header.Get(HeaderUserID)
	roleHeader := r.Header.Get(HeaderUserRole)

	if userID, _ = strconv.ParseInt(clientIDHeader, 10, 0); userID == ownerUserID {
		return true
	}
	if _, ok := permittedRoles[roleHeader]; ok {
		return true
	}
	return false
}
