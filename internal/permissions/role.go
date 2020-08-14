package permissions

import (
	"net/http"
)

// IRoleChecker — это интерфейс проверяльщика доступа к ресурсу по роли.
type IRoleChecker interface {
	Check(r *http.Request) bool
}

// RoleChecker — это проверяльщик, разрешающий доступ к ресурсу перечисленным ролям.
type RoleChecker struct {
	permittedRoles []Role
}

// NewRoleChecker создаёт проверяльщик, разрешающий доступ к ресурсу перечисленным ролям.
func NewRoleChecker(permittedRoles []Role) RoleChecker {
	return RoleChecker{
		permittedRoles: permittedRoles,
	}
}

// Check сообщает, разрешён ли доступ к ресурсу.
func (c RoleChecker) Check(r *http.Request) bool {
	permittedRoles := map[string]struct{}{}
	for _, r := range c.permittedRoles {
		permittedRoles[string(r)] = struct{}{}
	}

	roleHeader := r.Header.Get(HeaderUserRole)

	if _, ok := permittedRoles[roleHeader]; ok {
		return true
	}
	return false
}
