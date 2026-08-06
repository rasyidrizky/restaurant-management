package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	"project-2026-06-misoastory-be-go/internal/common/models"
	"project-2026-06-misoastory-be-go/internal/common/utils"
	"project-2026-06-misoastory-be-go/internal/config"

	"gorm.io/gorm"
)

func main() {
	config.Load()

	db, err := config.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("🌱 Starting database seeding...")

	db.AutoMigrate(&models.Product{}, &models.ProductLocation{})
	
	cleanDatabase(db)

	// 1. Seed Positions
	adminPos := seedPosition(db, "Administrator", "Full system access with all permissions")
	adminMisoaPos := seedPosition(db, "Admin Misoa", "Full system access for Admin Misoa")
	memberPos := seedPosition(db, "Member", "Standard user with limited permissions")
	staffPos := seedPosition(db, "Staff Outlet", "Staff at the store level")
	log.Println("✅ Positions seeded")

	// 2. Seed Permissions
	permissions := seedPermissions(db)

	// 3. Assign Permissions to Positions
	assignPermissions(db, permissions, adminPos, adminMisoaPos, memberPos, staffPos)

	// 4. Seed Users
	seedUsers(db, adminPos, adminMisoaPos, memberPos, staffPos)

	// 5. Seed Locations & Categories
	seedRestaurantData(db)

	log.Println("\n✨ Database seeding completed!")
}

func cleanDatabase(db *gorm.DB) {
	log.Println("🧹 Cleaning database...")
	err := db.Exec(`
		TRUNCATE TABLE 
			users, 
			position_permissions, 
			permissions, 
			positions, 
			categories, 
			locations,
			products,
			product_locations
		RESTART IDENTITY CASCADE;
	`).Error
	if err != nil {
		log.Printf("⚠️ Failed to clean database: %v\n", err)
	} else {
		log.Println("✅ Database cleaned")
	}
}

func seedPosition(db *gorm.DB, name, description string) models.Position {
	var pos models.Position
	err := db.Where("name = ?", name).FirstOrCreate(&pos, models.Position{
		Name:        name,
		Description: utils.StringPtr(description),
	}).Error
	if err != nil {
		log.Fatalf("Failed to seed position %s: %v", name, err)
	}
	return pos
}

func seedPermissions(db *gorm.DB) []models.Permission {
	log.Println("🔐 Seeding permissions...")
	permissionsData := []models.Permission{
		// User Management
		{Name: "VIEW_USER", Resource: "USER", Action: "VIEW"},
		{Name: "ADD_USER", Resource: "USER", Action: "ADD"},
		{Name: "UPDATE_USER", Resource: "USER", Action: "UPDATE"},
		{Name: "DELETE_USER", Resource: "USER", Action: "DELETE"},
		{Name: "MANAGE_USER_PERMISSION", Resource: "USER", Action: "MANAGE_PERMISSION"},
		{Name: "CHANGE_USER_POSITION", Resource: "USER", Action: "CHANGE_POSITION"},

		// Position Management
		{Name: "VIEW_POSITION", Resource: "POSITION", Action: "VIEW"},
		{Name: "ADD_POSITION", Resource: "POSITION", Action: "ADD"},
		{Name: "UPDATE_POSITION", Resource: "POSITION", Action: "UPDATE"},
		{Name: "DELETE_POSITION", Resource: "POSITION", Action: "DELETE"},

		// Permission Management
		{Name: "VIEW_PERMISSION", Resource: "PERMISSION", Action: "VIEW"},
		{Name: "ADD_PERMISSION", Resource: "PERMISSION", Action: "ADD"},
		{Name: "UPDATE_PERMISSION", Resource: "PERMISSION", Action: "UPDATE"},
		{Name: "DELETE_PERMISSION", Resource: "PERMISSION", Action: "DELETE"},

		// Category Management
		{Name: "VIEW_CATEGORY", Resource: "CATEGORY", Action: "VIEW"},
		{Name: "ADD_CATEGORY", Resource: "CATEGORY", Action: "ADD"},
		{Name: "UPDATE_CATEGORY", Resource: "CATEGORY", Action: "UPDATE"},
		{Name: "DELETE_CATEGORY", Resource: "CATEGORY", Action: "DELETE"},

		// Product Management
		{Name: "VIEW_PRODUCT", Resource: "PRODUCT", Action: "VIEW"},
		{Name: "ADD_PRODUCT", Resource: "PRODUCT", Action: "ADD"},
		{Name: "UPDATE_PRODUCT", Resource: "PRODUCT", Action: "UPDATE"},
		{Name: "DELETE_PRODUCT", Resource: "PRODUCT", Action: "DELETE"},



		// Location Management
		{Name: "VIEW_LOCATION", Resource: "LOCATION", Action: "VIEW"},
		{Name: "ADD_LOCATION", Resource: "LOCATION", Action: "ADD"},
		{Name: "UPDATE_LOCATION", Resource: "LOCATION", Action: "UPDATE"},
		{Name: "DELETE_LOCATION", Resource: "LOCATION", Action: "DELETE"},
	}

	var seeded []models.Permission
	for _, p := range permissionsData {
		var perm models.Permission
		if err := db.Where("resource = ? AND action = ?", p.Resource, p.Action).FirstOrCreate(&perm, p).Error; err != nil {
			log.Printf("Failed to create permission %s: %v", p.Name, err)
			continue
		}
		seeded = append(seeded, perm)
	}
	log.Printf("✅ %d permissions seeded\n", len(seeded))
	return seeded
}

func assignPermissions(db *gorm.DB, permissions []models.Permission, adminPos, adminMisoaPos, memberPos, staffPos models.Position) {
	log.Println("🔗 Assigning permissions to positions...")
	for _, p := range permissions {
		// Admin & Admin Misoa get all
		db.FirstOrCreate(&models.PositionPermission{}, models.PositionPermission{PositionID: adminPos.ID, PermissionID: p.ID})
		db.FirstOrCreate(&models.PositionPermission{}, models.PositionPermission{PositionID: adminMisoaPos.ID, PermissionID: p.ID})

		// Member gets VIEW_USER
		if p.Name == "VIEW_USER" {
			db.FirstOrCreate(&models.PositionPermission{}, models.PositionPermission{PositionID: memberPos.ID, PermissionID: p.ID})
		}

		// Staff gets all EXCEPT USER, POSITION, PERMISSION
		if p.Resource != "USER" && p.Resource != "POSITION" && p.Resource != "PERMISSION" {
			db.FirstOrCreate(&models.PositionPermission{}, models.PositionPermission{PositionID: staffPos.ID, PermissionID: p.ID})
		}
	}
	log.Println("✅ Permissions assigned to positions")
}

func seedUsers(db *gorm.DB, adminPos, adminMisoaPos, memberPos, staffPos models.Position) {
	log.Println("👥 Seeding users...")
	hashBytes, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	hashedPassword := string(hashBytes)

	users := []models.User{
		{Email: "admin@kulidigital.com", Password: hashedPassword, FirstName: "Admin", LastName: "Kuli Digital", PositionID: adminPos.ID, IsActive: true},
		{Email: "adminmisoa@kulidigital.com", Password: hashedPassword, FirstName: "Admin", LastName: "Misoa", PositionID: adminMisoaPos.ID, IsActive: true},
		{Email: "member@kulidigital.com", Password: hashedPassword, FirstName: "Member", LastName: "User", PositionID: memberPos.ID, IsActive: true},
		{Email: "staffmisoa@kulidigital.com", Password: hashedPassword, FirstName: "Staff", LastName: "Misoa Outlet", PositionID: staffPos.ID, IsActive: true},
	}

	for _, u := range users {
		db.Where("email = ?", u.Email).FirstOrCreate(&u, u)
	}
	log.Println("✅ Users seeded")
}

func seedRestaurantData(db *gorm.DB) {
	log.Println("🍜 Seeding restaurant data (Categories, Locations, Products)...")

	// Locations
	bandung := models.Location{Name: "Misoa Bandung", Slug: "misoa-bandung", Address: "Jl. Dago No. 1", City: "Bandung", Phone: utils.StringPtr("022-1234567"), IsActive: true, IsOpen24Hours: false, HasDineIn: true, SupportsHomeService: true, SupportsEvents: true}
	jakarta := models.Location{Name: "Misoa Jakarta", Slug: "misoa-jakarta", Address: "Jl. Sudirman No. 10", City: "Jakarta", Phone: utils.StringPtr("021-7654321"), IsActive: true, IsOpen24Hours: true, HasDineIn: true, SupportsHomeService: false, SupportsEvents: false}
	db.Where("slug = ?", bandung.Slug).FirstOrCreate(&bandung, bandung)
	db.Where("slug = ?", jakarta.Slug).FirstOrCreate(&jakarta, jakarta)

	// Categories
	snacks := models.Category{Name: "Snacks", Slug: "snacks", DisplayOrder: 1, IsActive: true}
	drinks := models.Category{Name: "Drinks", Slug: "drinks", DisplayOrder: 2, IsActive: true}
	db.Where("slug = ?", snacks.Slug).FirstOrCreate(&snacks, snacks)
	db.Where("slug = ?", drinks.Slug).FirstOrCreate(&drinks, drinks)

	// Products
	misoGoreng := models.Product{Name: "Miso Goreng Original", Slug: "miso-goreng-original", Description: utils.StringPtr("Mie soa goreng khas Misoa"), Price: 18000, IsBestSeller: true, CategoryID: snacks.ID, IsAvailable: true}
	esTehManis := models.Product{Name: "Es Teh Manis", Slug: "es-teh-manis", Description: utils.StringPtr("Teh manis dingin segar"), Price: 5000, CategoryID: drinks.ID, IsAvailable: true}
	db.Where("slug = ?", misoGoreng.Slug).FirstOrCreate(&misoGoreng, misoGoreng)
	db.Where("slug = ?", esTehManis.Slug).FirstOrCreate(&esTehManis, esTehManis)

	// Product Locations
	db.FirstOrCreate(&models.ProductLocation{}, models.ProductLocation{ProductID: misoGoreng.ID, LocationID: bandung.ID, IsAvailable: true})
	db.FirstOrCreate(&models.ProductLocation{}, models.ProductLocation{ProductID: misoGoreng.ID, LocationID: jakarta.ID, IsAvailable: true})

	log.Println("✅ Restaurant data seeded")
}
