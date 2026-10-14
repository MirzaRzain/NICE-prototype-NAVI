package seeder

import (
	"log"
	"os"

	"github.com/kristijorgji/goseeder"
)

type MenuPermission struct {
	ID         int
	MenuID     int
	ActionName string
	Code       string
}

var menuPermissions = []MenuPermission{
	{ID: 101, MenuID: 1, ActionName: "read", Code: "R"},

	{ID: 301, MenuID: 3, ActionName: "create", Code: "C"},
	{ID: 302, MenuID: 3, ActionName: "read", Code: "R"},
	{ID: 303, MenuID: 3, ActionName: "update", Code: "U"},
	{ID: 304, MenuID: 3, ActionName: "delete", Code: "D"},

	{ID: 401, MenuID: 4, ActionName: "create", Code: "C"},
	{ID: 402, MenuID: 4, ActionName: "read", Code: "R"},
	{ID: 403, MenuID: 4, ActionName: "update", Code: "U"},
	{ID: 404, MenuID: 4, ActionName: "delete", Code: "D"},

	{ID: 201, MenuID: 2, ActionName: "read", Code: "R"},
	{ID: 202, MenuID: 2, ActionName: "create", Code: "C"},
	{ID: 203, MenuID: 2, ActionName: "update", Code: "U"},
	{ID: 204, MenuID: 2, ActionName: "delete", Code: "D"},
}

func menuPermissionSeeder(s goseeder.Seeder) {
	const statement = `
	INSERT INTO app_menu_permission(
		id, menu_id, action_name, code, created_by, created_at
	)
	VALUES($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
	ON CONFLICT (id) DO NOTHING;
	`

	const createdBy = "fcdb2142-4731-470e-8a1b-8d0037665fb2"

	for _, mp := range menuPermissions {
		_, err := s.DB.Exec(statement,
			mp.ID,
			mp.MenuID,
			mp.ActionName,
			mp.Code,
			createdBy,
		)
		if err != nil {
			log.Fatalf("❌ ERROR execute menu_permission seeder : %v", err.Error())
		}
	}

	resetSeq := `
        SELECT setval(
            pg_get_serial_sequence('app_menu_permission', 'id'),
            COALESCE((SELECT MAX(id) FROM app_menu_permission), 1)
        );
    `

	if os.Getenv("DB_DRIVER") == "sqlite" {
		resetSeq = "SELECT 1"
	}
	if _, err := s.DB.Exec(resetSeq); err != nil {
		log.Fatalf("❌ ERROR execute menu_permission reset sequence : %v", err.Error())
	}
}
