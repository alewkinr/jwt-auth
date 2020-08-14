package permissions

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testsRoleOrOwner = []struct {
	testName string

	actingUserID    int64
	actingUserRole  Role
	resourceOwnerID int64
	permittedRoles  []Role

	expected bool
}{
	{
		testName: "be sure the admin role doesn't have an access to the alien resource by default",

		actingUserID:    321,
		actingUserRole:  RoleAdmin,
		resourceOwnerID: 123,
		permittedRoles:  []Role{},

		expected: false,
	},
	{
		testName: "be sure the manager role doesn't have an access to the alien resource by default",

		actingUserID:    321,
		actingUserRole:  RoleManager,
		resourceOwnerID: 123,
		permittedRoles:  []Role{},

		expected: false,
	},
	{
		testName: "be sure the client role doesn't have an access to the alien resource by default",

		actingUserID:    321,
		actingUserRole:  RoleClient,
		resourceOwnerID: 123,
		permittedRoles:  []Role{},

		expected: false,
	},
	{
		testName: "be sure the psychologist role doesn't have an access to the alien resource by default",

		actingUserID:    321,
		actingUserRole:  RolePsychologist,
		resourceOwnerID: 123,
		permittedRoles:  []Role{},

		expected: false,
	},
	{
		testName: "be sure owner has an access to the resource",

		actingUserID:    123,
		actingUserRole:  RoleClient,
		resourceOwnerID: 123,
		permittedRoles:  []Role{},

		expected: true,
	},
	{
		testName: "be sure a role permission works for the manager",

		actingUserID:    321,
		actingUserRole:  RoleManager,
		resourceOwnerID: 123,
		permittedRoles:  []Role{RoleManager},

		expected: true,
	},
	{
		testName: "be sure a role permission works for the admin",

		actingUserID:    321,
		actingUserRole:  RoleAdmin,
		resourceOwnerID: 123,
		permittedRoles:  []Role{RoleAdmin},

		expected: true,
	},
	{
		testName: "be sure a role permission works for the client",

		actingUserID:    321,
		actingUserRole:  RoleClient,
		resourceOwnerID: 123,
		permittedRoles:  []Role{RoleClient},

		expected: true,
	},
	{
		testName: "be sure a role permission works for the psychologist",

		actingUserID:    321,
		actingUserRole:  RolePsychologist,
		resourceOwnerID: 123,
		permittedRoles:  []Role{RolePsychologist},

		expected: true,
	},
	{
		testName: "be sure multiple role permissions work",

		actingUserID:    321,
		actingUserRole:  RolePsychologist,
		resourceOwnerID: 123,
		permittedRoles:  []Role{RoleClient, RolePsychologist},

		expected: true,
	},
}

// TestRoleOrOwnerChecker тестирует проверяльщик прав доступа, разрешающий доступ владельцу ресурса или группам.
func TestRoleOrOwnerChecker(t *testing.T) {
	for _, tt := range testsRoleOrOwner {
		tt := tt
		t.Run(tt.testName, func(t *testing.T) {
			c := NewRoleOrOwnerChecker(tt.permittedRoles)

			req, err := http.NewRequest("GET", "test.url", nil)
			if err != nil {
				assert.NoError(t, err)
			}
			req.Header.Set(HeaderUserID, strconv.FormatInt(tt.actingUserID, 10))
			req.Header.Set(HeaderUserRole, string(tt.actingUserRole))

			actual := c.Check(req, 123)
			assert.Exactly(t, tt.expected, actual)
		})
	}
}
