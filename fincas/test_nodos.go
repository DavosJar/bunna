package main

import (
	"context"
	"fmt"
	"github.com/davosjar/bunna/services/fincas/internal/fincas/infrastructure/persistence/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost port=5432 user=fincas_user password=fincas_pass_dev dbname=bunna_fincas sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	var nodos []postgres.NodoModel
	db.Find(&nodos)
	fmt.Printf("NODOS EN DB: %+v\n", nodos)
}
