package user

type Role int

const (
	RoleUnidentified Role = iota
	RoleNormal            = 1
	RoleAdmin             = 9
)
