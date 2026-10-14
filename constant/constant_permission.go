package constant

type MenuKey string
type PermissionCode string

const (
	ActionRead                 PermissionCode = "R"
	ActionCreate               PermissionCode = "C"
	ActionUpdate               PermissionCode = "U"
	ActionDelete               PermissionCode = "D"
	ActionReadMenuPermission   PermissionCode = "RMP"
	ActionCreateMenuPermission PermissionCode = "CMP"
	ActionUpdateMenuPermission PermissionCode = "UMP"
	ActionDeleteMenuPermission PermissionCode = "DMP"
)

const (
	MenuSettingsMenu MenuKey = "/settings/menu"
	MenuSettingsRole MenuKey = "/settings/role"
	MenuSettingsUser MenuKey = "/settings/user"
	MenuDashboard    MenuKey = "/dashboard"
)

func (a PermissionCode) String() string {
	return string(a)
}

func (m MenuKey) String() string {
	return string(m)
}
