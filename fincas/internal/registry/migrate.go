package registry

import (
	diagnosticopostgres "github.com/davosjar/bunna/services/fincas/internal/diagnostico/infrastructure/persistence/postgres"
	fincaspostgres "github.com/davosjar/bunna/services/fincas/internal/fincas/infrastructure/persistence/postgres"
	nodospostgres "github.com/davosjar/bunna/services/fincas/internal/nodos/infrastructure/persistence/postgres"
	iampostgres "github.com/davosjar/bunna/services/fincas/internal/infrastructure/security/iam/postgres"
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
		&nodospostgres.NodoModel{},
		&iampostgres.IamRolPermisosModel{},
	)
}
