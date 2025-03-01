package seeders

import (
	"docker-control-go/src/configs"
	"docker-control-go/src/database/models"
	logger "docker-control-go/src/log"
	"log"
)

// SeedPermissions akan menambahkan role dan permission default ke Casbin
func SeedPermissions() {
	// Tambahkan roles
	roles := []string{"admin", "user"}

	// Tambahkan permission default
	permissions := []models.Permission{
		{Role: "admin", Obj: "resource:user", Act: "read"},
		{Role: "admin", Obj: "resource:user", Act: "write"},
		{Role: "user", Obj: "resource:user", Act: "read"},
	}

	// Tambahkan role ke Casbin
	for _, role := range roles {
		_, err := configs.Enforcer.AddGroupingPolicy(role, role) // Pastikan role tersedia
		if err != nil {
			log.Printf("Failed to add role %s: %v", role, err)
			logger.Log.Errorf("Failed to add role %s: %v", role, err)
		}
	}

	// Tambahkan permission ke Casbin
	for _, perm := range permissions {
		_, err := configs.Enforcer.AddPolicy(perm.Role, perm.Obj, perm.Act)
		if err != nil {
			log.Printf("Failed to add permission %s -> %s:%s: %v", perm.Role, perm.Obj, perm.Act, err)
			logger.Log.Errorf("Failed to add permission %s -> %s:%s: %v", perm.Role, perm.Obj, perm.Act, err)
		}
	}

	log.Println("Seeding completed successfully!")
	logger.Log.Info("Seeding completed successfully!")
}

// ResetCasbinData akan menghapus semua policy dan grouping policy di Casbin
func ResetCasbinData() {
	var err error
	// Hapus semua policy
	configs.Enforcer.ClearPolicy() // Tidak perlu menangkap nilai kembalian
	log.Println("✅ Casbin policies cleared!")
	logger.Log.Info("✅ Casbin policies cleared!")

	// Simpan perubahan
	err = configs.Enforcer.SavePolicy()
	if err != nil {
		log.Fatalf("Failed to save cleared Casbin policies: %v", err)
		logger.Log.Fatalf("Failed to save cleared Casbin policies: %v", err)
	}
	log.Println("✅ Casbin policies saved after reset!")
	logger.Log.Info("✅ Casbin policies saved after reset!")
}
