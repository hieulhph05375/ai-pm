package entity

type RoleWithPermissions struct {
	Role
	Permissions []Permission `json:"permissions"`
}
