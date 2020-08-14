package permissions

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testsRoleChecker = []struct {
	testName string

	actingUserRole Role
	permittedRoles []Role

	expected bool
}{
	{
		testName: "be sure the role doesn't have an access to the alien resource by default",

		actingUserRole: RoleManager,
		permittedRoles: []Role{},

		expected: false,
	},
	{
		testName: "be sure the role has an access to the resource permitted only for this role",

		actingUserRole: RoleManager,
		permittedRoles: []Role{RoleManager},

		expected: true,
	},
	{
		testName: "be sure the role has an access to the resource permitted for this and other roles",

		actingUserRole: RoleManager,
		permittedRoles: []Role{RoleClient, RoleManager},

		expected: true,
	},
}

// TestRoleChecker тестирует проверяльщик прав доступа, разрешающий доступ указанным ролям.
func TestRoleChecker(t *testing.T) {
	for _, tt := range testsRoleChecker {
		tt := tt
		t.Run(tt.testName, func(t *testing.T) {
			c := NewRoleChecker(tt.permittedRoles)

			req, err := http.NewRequest("GET", "test.url", nil)
			if err != nil {
				assert.NoError(t, err)
			}
			req.Header.Set(HeaderUserRole, string(tt.actingUserRole))

			actual := c.Check(req)
			assert.Exactly(t, tt.expected, actual)
		})
	}
}
