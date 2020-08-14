package permissions

// TODO casbin
// По возможности отказаться от своей реализации в пользу https://github.com/casbin/casbin.

// Role — это роль пользователя в системе example.comмунд Онлайн.
type Role string

const (
	// Роли пользователей.
	RoleAdmin        Role = "admin"        // Админ.
	RoleManager      Role = "manager"      // Оператор.
	RoleClient       Role = "client"       // Клиент.
	RolePsychologist Role = "psychologist" // Психолог.

	HeaderUserID   = "X-Zig-User-ID"   // HTTP-заголовок в запросе, содержащий идентификатор пользователя.
	HeaderUserRole = "X-Zig-User-Role" // HTTP-заголовок в запросе, содержащий роль пользователя.
)
