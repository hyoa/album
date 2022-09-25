package user

type role int

const (
	RoleUnidentified role = iota
	RoleNormal            = 1
	RoleAdmin             = 9
)
