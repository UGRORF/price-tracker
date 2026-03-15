package domain

type UserRole string

const (
	RoleUser      UserRole = "user"
	RoleAdmin     UserRole = "admin"
	RoleModerator UserRole = "moderator"
)

type User struct {
	ID       int64    `db:"id" json:"id"`
	Username string   `db:"username" json:"username"`
	Password string   `db:"password" json:"-"`
	Role     UserRole `db:"role" json:"role"`
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsModerator() bool {
	return u.Role == RoleModerator
}

func (u *User) CanModerate() bool {
	return u.Role == RoleAdmin || u.Role == RoleModerator
}

func (u *User) IsValidRole() bool {
	return u.Role == RoleUser || u.Role == RoleAdmin || u.Role == RoleModerator
}
