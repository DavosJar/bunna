package registry

import (
	diagnosticopostgres "github.com/davosjar/bunna/services/fincas/internal/diagnostico/infrastructure/persistence/postgres"
	fincaspostgres "github.com/davosjar/bunna/services/fincas/internal/fincas/infrastructure/persistence/postgres"
	"gorm.io/gorm"
)

// runAutoMigrate ejecuta AutoMigrate para todos los modelos del microservicio.
func runAutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&fincaspostgres.FincaModel{},
		&fincaspostgres.LoteModel{},
		&diagnosticopostgres.MuestraModel{},
		&diagnosticopostgres.DiagnosticoModel{},
		&diagnosticopostgres.CandidatoModel{},
	)
}
