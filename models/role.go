package models

type Role string

const (
	RoleAdmin        Role = "admin"
	RoleChatOwner    Role = "owner"
	RoleChatAdmin    Role = "admin"
	RoleChatMember   Role = "member"
	RoleDeelnemer    Role = "Deelnemer"
	RoleBegeleider   Role = "Begeleider"
	RoleVrijwilliger Role = "Vrijwilliger"
	RoleStaff        Role = "staff"
)
