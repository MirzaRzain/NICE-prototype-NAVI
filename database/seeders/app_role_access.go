package seeder

import (
	"log"

	"github.com/kristijorgji/goseeder"
)

type RoleAccessResponse struct {
	RoleID           int
	MenuPermissionID int
}

var roleAccesses = []RoleAccessResponse{
	{RoleID: 1, MenuPermissionID: 101},

	{RoleID: 1, MenuPermissionID: 201},
	{RoleID: 1, MenuPermissionID: 202},
	{RoleID: 1, MenuPermissionID: 203},
	{RoleID: 1, MenuPermissionID: 204},

	{RoleID: 1, MenuPermissionID: 301},
	{RoleID: 1, MenuPermissionID: 302},
	{RoleID: 1, MenuPermissionID: 303},
	{RoleID: 1, MenuPermissionID: 304},

	{RoleID: 1, MenuPermissionID: 401},
	{RoleID: 1, MenuPermissionID: 402},
	{RoleID: 1, MenuPermissionID: 403},
	{RoleID: 1, MenuPermissionID: 404},

	{RoleID: 2, MenuPermissionID: 101},
	{RoleID: 2, MenuPermissionID: 201},
}

func roleAccessSeeder(s goseeder.Seeder) {
	const statement = `
	INSERT INTO app_role_access(
		role_id, menu_permission_id
	)
	VALUES($1, $2)
	ON CONFLICT (role_id, menu_permission_id) DO NOTHING;
	`
	for i := range roleAccesses {
		_, err := s.DB.Exec(statement,
			roleAccesses[i].RoleID,
			roleAccesses[i].MenuPermissionID,
		)
		if err != nil {
			log.Fatalf("❌ ERROR execute role_access seeder : %v", err.Error())
		}
	}

}
