package main

import (
	"fmt"
	
	nodospostgres "github.com/davosjar/bunna/services/fincas/internal/nodos/infrastructure/persistence/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost port=5432 user=fincas_user password=fincas_pass_dev dbname=fincas_db sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	var nodos []nodospostgres.NodoModel
	db.Find(&nodos)
	fmt.Printf("NODOS EN DB: %+v\n", nodos)
}
