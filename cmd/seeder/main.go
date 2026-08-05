package main

import (
	"log"

	"project-2026-06-misoastory-be-go/internal/config"
	"project-2026-06-misoastory-be-go/internal/common/models"
)

func main() {
	config.Load()

	db, err := config.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 1. Get the Admin position
	var admin models.Position
	if err := db.Where("name = ?", "Admin").First(&admin).Error; err != nil {
		log.Fatalf("Admin position not found. Did you register a user yet? Error: %v", err)
	}

	// 2. Define the exact permissions we need
	permissions := []models.Permission{
		{Name: "Add Category", Resource: "CATEGORY", Action: "ADD"},
		{Name: "Update Category", Resource: "CATEGORY", Action: "UPDATE"},
		{Name: "Delete Category", Resource: "CATEGORY", Action: "DELETE"},
		{Name: "Add Location", Resource: "LOCATION", Action: "ADD"},
		{Name: "Update Location", Resource: "LOCATION", Action: "UPDATE"},
		{Name: "Delete Location", Resource: "LOCATION", Action: "DELETE"},
		{Name: "View User", Resource: "USER", Action: "VIEW"},
	}

	// 3. Create permissions and link them to Admin
	for _, p := range permissions {
		var perm models.Permission

		// Create or find permission
		if err := db.Where("resource = ? AND action = ?", p.Resource, p.Action).FirstOrCreate(&perm, p).Error; err != nil {
			log.Printf("Failed to create permission %s: %v", p.Name, err)
			continue
		}

		// Link to Admin
		pp := models.PositionPermission{
			PositionID:   admin.ID,
			PermissionID: perm.ID,
		}
		if err := db.Where("position_id = ? AND permission_id = ?", admin.ID, perm.ID).FirstOrCreate(&pp, pp).Error; err != nil {
			log.Printf("Failed to link permission %s to Admin: %v", p.Name, err)
		} else {
			log.Printf("Granted %s:%s to Admin", p.Resource, p.Action)
		}
	}

	log.Println("Seeding complete! Admin has full access.")
}
