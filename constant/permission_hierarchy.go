package constant

// MenuGroup is a logical grouping for menu keys (useful for "folder-like" hierarchy without introducing new menu keys).
type MenuGroup string

const (
	MenuGroupSettings  MenuGroup = "settings"
	MenuGroupDashboard MenuGroup = "dashboard"
)

// MenuGroupHierarchy defines the hierarchy for menu keys at the group level.
var MenuGroupHierarchy = map[MenuGroup][]MenuKey{
	MenuGroupSettings: {
		MenuSettingsMenu,
		MenuSettingsRole,
		MenuSettingsUser,
	},
	MenuGroupDashboard: {
		MenuDashboard,
	},
}

// MenuKeyPermissions defines which permission codes are valid/expected per menu key.
var MenuKeyPermissions = map[MenuKey][]PermissionCode{
	MenuSettingsMenu: {
		ActionRead,
		ActionCreate,
		ActionUpdate,
		ActionDelete,
		ActionReadMenuPermission,
		ActionCreateMenuPermission,
		ActionUpdateMenuPermission,
		ActionDeleteMenuPermission,
	},
	MenuSettingsRole: {
		ActionRead,
		ActionCreate,
		ActionUpdate,
		ActionDelete,
	},
	MenuSettingsUser: {
		ActionRead,
		ActionCreate,
		ActionUpdate,
		ActionDelete,
	},
	MenuDashboard: {
		ActionRead,
	},
}

// PermissionImplications defines inherited permissions.
// Example: if a user can Update, they also implicitly can Read.
var PermissionImplications = map[PermissionCode][]PermissionCode{
	ActionCreate: {ActionRead},
	ActionUpdate: {ActionRead},
	ActionDelete: {ActionRead},

	ActionCreateMenuPermission: {ActionReadMenuPermission},
	ActionUpdateMenuPermission: {ActionReadMenuPermission},
	ActionDeleteMenuPermission: {ActionReadMenuPermission},
}

// InheritedPermissionCodes defines alternative permission codes that are also accepted for authorization checks.
// Example: for menu-permission actions, the generic CRUD codes can be treated as a fallback.
var InheritedPermissionCodes = map[PermissionCode][]PermissionCode{
	ActionReadMenuPermission:   {ActionReadMenuPermission, ActionRead},
	ActionCreateMenuPermission: {ActionCreateMenuPermission, ActionCreate},
	ActionUpdateMenuPermission: {ActionUpdateMenuPermission, ActionUpdate},
	ActionDeleteMenuPermission: {ActionDeleteMenuPermission, ActionDelete},
}

// GetInheritedPermissions returns permission codes (as strings) that are acceptable for a given permission check.
// If no inheritance rule exists, it returns the code itself.
func GetInheritedPermissions(code PermissionCode) []string {
	if codes, found := InheritedPermissionCodes[code]; found {
		out := make([]string, 0, len(codes))
		for _, c := range codes {
			out = append(out, c.String())
		}
		return out
	}

	return []string{code.String()}
}

// ExpandPermissionCodes returns the closure of codes + implied codes, de-duplicated.
func ExpandPermissionCodes(codes []PermissionCode) []PermissionCode {
	seen := make(map[PermissionCode]struct{}, len(codes))
	out := make([]PermissionCode, 0, len(codes))

	var add func(code PermissionCode)
	add = func(code PermissionCode) {
		if _, ok := seen[code]; ok {
			return
		}
		seen[code] = struct{}{}
		out = append(out, code)
		for _, implied := range PermissionImplications[code] {
			add(implied)
		}
	}

	for _, code := range codes {
		add(code)
	}

	return out
}
