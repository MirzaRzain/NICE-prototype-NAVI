package seeder

import (
	"database/sql"
	"log"
	"os"

	"github.com/kristijorgji/goseeder"
)

type Menu struct {
	ID        int
	ParentID  sql.NullInt32
	Group     string
	Name      string
	URL       string
	Icon      string
	Sort      int
	CreatedBy string
}

var menus = []Menu{
	{
		ID:        1,
		ParentID:  sql.NullInt32{Valid: false},
		Group:     "main",
		Name:      "Dashboard",
		URL:       "/dashboard",
		Icon:      "ri-dashboard-fill",
		Sort:      1,
		CreatedBy: "fcdb2142-4731-470e-8a1b-8d0037665fb2",
	},
	{
		ID:        2,
		ParentID:  sql.NullInt32{Valid: false},
		Group:     "system",
		Name:      "Settings",
		Icon:      "ri-settings-3-fill",
		Sort:      1,
		CreatedBy: "fcdb2142-4731-470e-8a1b-8d0037665fb2",
	},
	{
		ID:        3,
		ParentID:  sql.NullInt32{Int32: 2, Valid: true},
		Group:     "system",
		Name:      "Menu",
		URL:       "/settings/menu",
		Icon:      "ri-menu-line",
		Sort:      1,
		CreatedBy: "fcdb2142-4731-470e-8a1b-8d0037665fb2",
	},
	{
		ID:        4,
		ParentID:  sql.NullInt32{Int32: 2, Valid: true},
		Group:     "system",
		Name:      "Role",
		URL:       "/settings/role",
		Icon:      "ri-menu-line",
		Sort:      2,
		CreatedBy: "fcdb2142-4731-470e-8a1b-8d0037665fb2",
	},
}

func menuSeeder(s goseeder.Seeder) {
	const statement = `
	INSERT INTO app_menu(
		id, parent_id, name, url, icon, sort, created_by, created_at
	)
	VALUES($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
	ON CONFLICT (id) DO UPDATE
	SET name = EXCLUDED.name, url = EXCLUDED.url, parent_id = EXCLUDED.parent_id;
	`
	for i := range menus {
		_, err := s.DB.Exec(statement,
			menus[i].ID,
			menus[i].ParentID,
			menus[i].Name,
			menus[i].URL,
			menus[i].Icon,
			menus[i].Sort,
			menus[i].CreatedBy,
		)

		if err != nil {
			log.Fatalf("❌ ERROR execute menu seeder : %v", err.Error())
		}
	}
	resetSeq := `
        SELECT setval(
            pg_get_serial_sequence('app_menu', 'id'),
            COALESCE((SELECT MAX(id) FROM app_menu), 1)
        );
    `
	if os.Getenv("DB_DRIVER") == "sqlite" {
		resetSeq = "SELECT 1"
	}
	if _, err := s.DB.Exec(resetSeq); err != nil {
		log.Fatalf("❌ ERROR execute menu seeder : %v", err.Error())
	}

}
