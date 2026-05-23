package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	postgresdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTestDB initializes test database connection directly (avoiding import cycle)
func setupTestDB(t *testing.T) *gorm.DB {
	// Load configuration from environment
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "identidad_user"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "identidad_pass_dev"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "identidad_db"
	}
	dbSSLMode := os.Getenv("DB_SSLMODE")
	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	// Build DSN
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s client_encoding=UTF8",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	// Initialize database connection
	db, err := gorm.Open(postgresdriver.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Create uuid-ossp extension if it doesn't exist
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		t.Fatalf("Failed to create uuid-ossp extension: %v", err)
	}
	
	if err := db.AutoMigrate(&CredencialesModel{}); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

// createTestCredenciales creates a new test credential with specified usuarioID
func createTestCredenciales(usuarioID string) *domain.CredencialesUsuario {
	return domain.NuevaCredencialesUsuario(usuarioID, "test_hash_"+usuarioID)
}

// TestCrearCredenciales tests basic creation of credentials
func TestCrearCredenciales(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Arrange
	cred := createTestCredenciales("user_create_1")

	// Act
	created, err := repo.Crear(ctx, cred)

	// Assert
	if err != nil {
		t.Fatalf("Crear() error = %v", err)
	}
	if created == nil {
		t.Fatal("Crear() returned nil")
	}
	if created.UsuarioID() != "user_create_1" {
		t.Errorf("UsuarioID() = %q, want %q", created.UsuarioID(), "user_create_1")
	}
	if !created.Activo() {
		t.Error("Activo() should be true by default")
	}
	if created.PasswordHash() != "test_hash_user_create_1" {
		t.Errorf("PasswordHash mismatch")
	}
}

// TestCrearMultipleCredenciales tests creating multiple credentials
func TestCrearMultipleCredenciales(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Arrange
	userIDs := []string{"user1", "user2", "user3"}

	// Act & Assert
	for _, id := range userIDs {
		cred := createTestCredenciales(id)
		created, err := repo.Crear(ctx, cred)
		if err != nil {
			t.Fatalf("Crear() for user %s error = %v", id, err)
		}
		if created.UsuarioID() != id {
			t.Errorf("UsuarioID mismatch for %s", id)
		}
	}

	// Verify all were created
	all, err := repo.Find(ctx, domain.EspecificacionCredenciales{}, shareddomain.Paginacion{
		Pagina:       1,
		TamanoPagina: 100,
	})
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	if len(all) != 3 {
		t.Errorf("Find() returned %d credentials, want 3", len(all))
	}
}

// TestActualizarCredenciales tests updating credential fields
func TestActualizarCredenciales(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Arrange - Create initial credential
	cred := createTestCredenciales("user_update_1")
	created, err := repo.Crear(ctx, cred)
	if err != nil {
		t.Fatalf("Crear() error = %v", err)
	}

	// Act - Update fields
	created.VerificarCorreo()
	created.Desactivar()
	ahora := time.Now()
	created.MarcarIntentoFallido(ahora)

	updated, err := repo.Actualizar(ctx, created)

	// Assert
	if err != nil {
		t.Fatalf("Actualizar() error = %v", err)
	}
	if !updated.CorreoVerificado() {
		t.Error("CorreoVerificado should be true after update")
	}
	if updated.Activo() {
		t.Error("Activo should be false after update")
	}
	if updated.IntentosFallidos() != 1 {
		t.Errorf("IntentosFallidos = %d, want 1", updated.IntentosFallidos())
	}
}

// TestActualizarMultipleFields tests updating multiple credential fields
func TestActualizarMultipleFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Arrange
	cred := createTestCredenciales("user_update_multi")
	created, _ := repo.Crear(ctx, cred)

	// Act - Perform multiple updates
	created.VerificarCorreo()
	for i := 0; i < 3; i++ {
		created.MarcarIntentoFallido(time.Now())
	}
	updated1, _ := repo.Actualizar(ctx, created)

	// Act again - more updates
	updated1.Desactivar()
	updated2, _ := repo.Actualizar(ctx, updated1)

	// Assert - verify cumulative updates
	if !updated2.CorreoVerificado() {
		t.Error("CorreoVerificado should still be true")
	}
	if updated2.Activo() {
		t.Error("Activo should now be false")
	}
	if updated2.IntentosFallidos() != 3 {
		t.Errorf("IntentosFallidos = %d, want 3", updated2.IntentosFallidos())
	}
}

// TestActualizarNoExistente tests updating a credential that doesn't exist
func TestActualizarNoExistente(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Arrange
	cred := createTestCredenciales("user_nonexistent")

	// Act
	_, err := repo.Actualizar(ctx, cred)

	// Assert
	if err == nil {
		t.Error("Actualizar() should return error for nonexistent credential")
	}
}

// TestObtenerPorUsuarioID tests retrieving credential by usuarioID
func TestObtenerPorUsuarioID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Arrange
	cred := createTestCredenciales("user_get_1")
	created, _ := repo.Crear(ctx, cred)

	// Act
	retrieved, err := repo.ObtenerPorUsuarioID(ctx, "user_get_1")

	// Assert
	if err != nil {
		t.Fatalf("ObtenerPorUsuarioID() error = %v", err)
	}
	if retrieved.UsuarioID() != created.UsuarioID() {
		t.Errorf("UsuarioID mismatch")
	}
	if retrieved.PasswordHash() != created.PasswordHash() {
		t.Errorf("PasswordHash mismatch")
	}
}

// TestObtenerPorUsuarioIDNoExistente tests retrieving a credential that doesn't exist
func TestObtenerPorUsuarioIDNoExistente(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Act
	_, err := repo.ObtenerPorUsuarioID(ctx, "nonexistent_user")

	// Assert
	if err == nil {
		t.Error("ObtenerPorUsuarioID() should return error for nonexistent credential")
	}
}

// TestEliminarCredenciales tests deleting an existing credential
func TestEliminarCredenciales(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Arrange
	cred := createTestCredenciales("user_delete_1")
	repo.Crear(ctx, cred)

	// Act
	err := repo.Eliminar(ctx, "user_delete_1")

	// Assert
	if err != nil {
		t.Fatalf("Eliminar() error = %v", err)
	}

	// Verify it's gone
	_, err = repo.ObtenerPorUsuarioID(ctx, "user_delete_1")
	if err == nil {
		t.Error("Credential should not exist after deletion")
	}
}

// TestEliminarNoExistente tests deleting a credential that doesn't exist
func TestEliminarNoExistente(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Act
	err := repo.Eliminar(ctx, "nonexistent_user")

	// Assert
	if err == nil {
		t.Error("Eliminar() should return error for nonexistent credential")
	}
}

// TestFindSinFiltros tests Find() with no filters returns all credentials
func TestFindSinFiltros(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Arrange - Create 5 credentials
	for i := 1; i <= 5; i++ {
		id := "user_find_" + string(rune(48+i))
		cred := createTestCredenciales(id)
		repo.Crear(ctx, cred)
	}

	// Act
	result, err := repo.Find(ctx, domain.EspecificacionCredenciales{}, shareddomain.Paginacion{
		Pagina:       1,
		TamanoPagina: 100,
	})

	// Assert
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	if len(result) != 5 {
		t.Errorf("Find() returned %d credentials, want 5", len(result))
	}
}

// TestFindConFiltroActivo tests Find() with activo = true filter
func TestFindConFiltroActivo(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Arrange - Create mix of active and inactive credentials
	for i := 1; i <= 5; i++ {
		id := "user_active_" + string(rune(48+i))
		cred := createTestCredenciales(id)
		created, _ := repo.Crear(ctx, cred)

		// Deactivate every other one
		if i%2 == 0 {
			created.Desactivar()
			repo.Actualizar(ctx, created)
		}
	}

	// Act - Filter for active = true
	result, err := repo.Find(ctx, domain.EspecificacionCredenciales{
		ListaFiltros: []shareddomain.CriterioFiltro{
			{Campo: "activo", Operador: "=", Valor: true},
		},
	}, shareddomain.Paginacion{
		Pagina:       1,
		TamanoPagina: 100,
	})

	// Assert
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	if len(result) != 3 {
		t.Errorf("Find(activo=true) returned %d, want 3", len(result))
	}
	for _, cred := range result {
		if !cred.Activo() {
			t.Errorf("Found inactive credential when filtering for active=true")
		}
	}
}

// TestFindConFiltroInactivo tests Find() with activo = false filter
func TestFindConFiltroInactivo(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Arrange - Create mix of active and inactive
	for i := 1; i <= 6; i++ {
		id := "user_inactive_" + string(rune(48+i))
		cred := createTestCredenciales(id)
		created, _ := repo.Crear(ctx, cred)

		// Deactivate some
		if i > 3 {
			created.Desactivar()
			repo.Actualizar(ctx, created)
		}
	}

	// Act
	result, err := repo.Find(ctx, domain.EspecificacionCredenciales{
		ListaFiltros: []shareddomain.CriterioFiltro{
			{Campo: "activo", Operador: "=", Valor: false},
		},
	}, shareddomain.Paginacion{
		Pagina:       1,
		TamanoPagina: 100,
	})

	// Assert
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	if len(result) != 3 {
		t.Errorf("Find(activo=false) returned %d, want 3", len(result))
	}
}

// TestFindConFiltroNegacion tests Find() with != operator
func TestFindConFiltroNegacion(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Arrange
	for i := 1; i <= 5; i++ {
		id := "user_negation_" + string(rune(48+i))
		cred := createTestCredenciales(id)
		created, _ := repo.Crear(ctx, cred)
		if i <= 2 {
			created.Desactivar()
			repo.Actualizar(ctx, created)
		}
	}

	// Act - Find all that are NOT inactive (i.e., active)
	result, err := repo.Find(ctx, domain.EspecificacionCredenciales{
		ListaFiltros: []shareddomain.CriterioFiltro{
			{Campo: "activo", Operador: "!=", Valor: false},
		},
	}, shareddomain.Paginacion{
		Pagina:       1,
		TamanoPagina: 100,
	})

	// Assert
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	if len(result) != 3 {
		t.Errorf("Find(activo != false) returned %d, want 3", len(result))
	}
}

// TestFindOrdenacionASC tests Find() with ASC ordering
func TestFindOrdenacionASC(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Arrange - Create credentials with different usuarioIDs
	ids := []string{"user_z", "user_a", "user_m", "user_b"}
	for _, id := range ids {
		cred := createTestCredenciales(id)
		repo.Crear(ctx, cred)
	}

	// Act - Find with ASC ordering on usuarioID
	result, err := repo.Find(ctx, domain.EspecificacionCredenciales{}, shareddomain.Paginacion{
		Pagina:       1,
		TamanoPagina: 100,
		Ordenaciones: []shareddomain.Ordenacion{
			{Campo: "usuarioID", Tipo: shareddomain.ASC},
		},
	})

	// Assert
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	if len(result) != 4 {
		t.Errorf("Find() returned %d, want 4", len(result))
	}
	// Verify ascending order
	expected := []string{"user_a", "user_b", "user_m", "user_z"}
	for i, exp := range expected {
		if result[i].UsuarioID() != exp {
			t.Errorf("Order mismatch at position %d: got %q, want %q",
				i, result[i].UsuarioID(), exp)
		}
	}
}

// TestFindOrdenacionDESC tests Find() with DESC ordering
func TestFindOrdenacionDESC(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Arrange
	ids := []string{"user_z", "user_a", "user_m", "user_b"}
	for _, id := range ids {
		cred := createTestCredenciales(id)
		repo.Crear(ctx, cred)
	}

	// Act
	result, err := repo.Find(ctx, domain.EspecificacionCredenciales{}, shareddomain.Paginacion{
		Pagina:       1,
		TamanoPagina: 100,
		Ordenaciones: []shareddomain.Ordenacion{
			{Campo: "usuarioID", Tipo: shareddomain.DESC},
		},
	})

	// Assert
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	// Verify descending order
	expected := []string{"user_z", "user_m", "user_b", "user_a"}
	for i, exp := range expected {
		if result[i].UsuarioID() != exp {
			t.Errorf("Order mismatch at position %d: got %q, want %q",
				i, result[i].UsuarioID(), exp)
		}
	}
}

// TestFindPaginacionPrimeraPagee tests Find() with pagination - first page
func TestFindPaginacionPrimeraPage(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Arrange - Create 10 credentials
	for i := 1; i <= 10; i++ {
		id := "user_page_" + string(rune(47+i))
		cred := createTestCredenciales(id)
		repo.Crear(ctx, cred)
	}

	// Act - First page with size 3
	result, err := repo.Find(ctx, domain.EspecificacionCredenciales{}, shareddomain.Paginacion{
		Pagina:       1,
		TamanoPagina: 3,
	})

	// Assert
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	if len(result) != 3 {
		t.Errorf("Find() returned %d, want 3", len(result))
	}
}

// TestFindPaginacionSegundaPage tests Find() with pagination - second page
func TestFindPaginacionSegundaPage(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Arrange - Create 10 credentials with predictable IDs
	for i := 1; i <= 10; i++ {
		id := "user_" + string(rune(47+i))
		cred := createTestCredenciales(id)
		repo.Crear(ctx, cred)
	}

	// Act - Get second page with size 3, ordered by usuarioID
	result, err := repo.Find(ctx, domain.EspecificacionCredenciales{}, shareddomain.Paginacion{
		Pagina:       2,
		TamanoPagina: 3,
		Ordenaciones: []shareddomain.Ordenacion{
			{Campo: "usuarioID", Tipo: shareddomain.ASC},
		},
	})

	// Assert
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	if len(result) != 3 {
		t.Errorf("Find() page 2 returned %d, want 3", len(result))
	}
}

// TestFindPaginacionUltimaPage tests Find() with pagination - last page with fewer items
func TestFindPaginacionUltimaPage(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Arrange - Create 7 credentials
	for i := 1; i <= 7; i++ {
		id := "user_last_" + string(rune(47+i))
		cred := createTestCredenciales(id)
		repo.Crear(ctx, cred)
	}

	// Act - Last page with size 3 (should have 1 item)
	result, err := repo.Find(ctx, domain.EspecificacionCredenciales{}, shareddomain.Paginacion{
		Pagina:       3,
		TamanoPagina: 3,
	})

	// Assert
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	if len(result) != 1 {
		t.Errorf("Find() last page returned %d, want 1", len(result))
	}
}

// TestFindCompoundFiltersAndOrdering tests Find() with multiple filters and ordering
func TestFindCompoundFiltersAndOrdering(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Arrange - Create varied credentials
	testCases := []struct {
		id     string
		activo bool
		fallos int
	}{
		{"user_compound_a", true, 0},
		{"user_compound_b", true, 2},
		{"user_compound_c", false, 0},
		{"user_compound_d", true, 1},
		{"user_compound_e", false, 5},
	}

	for _, tc := range testCases {
		cred := createTestCredenciales(tc.id)
		created, _ := repo.Crear(ctx, cred)

		if !tc.activo {
			created.Desactivar()
		}
		for i := 0; i < tc.fallos; i++ {
			created.MarcarIntentoFallido(time.Now())
		}
		repo.Actualizar(ctx, created)
	}

	// Act - Find active credentials, ordered by usuarioID DESC
	result, err := repo.Find(ctx, domain.EspecificacionCredenciales{
		ListaFiltros: []shareddomain.CriterioFiltro{
			{Campo: "activo", Operador: "=", Valor: true},
		},
	}, shareddomain.Paginacion{
		Pagina:       1,
		TamanoPagina: 100,
		Ordenaciones: []shareddomain.Ordenacion{
			{Campo: "usuarioID", Tipo: shareddomain.DESC},
		},
	})

	// Assert
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	if len(result) != 3 {
		t.Errorf("Find() returned %d, want 3", len(result))
	}
	// Verify DESC ordering
	expected := []string{"user_compound_d", "user_compound_b", "user_compound_a"}
	for i, exp := range expected {
		if result[i].UsuarioID() != exp {
			t.Errorf("Order mismatch at position %d: got %q, want %q",
				i, result[i].UsuarioID(), exp)
		}
	}
}

// TestFindConIntentosFallidos tests Find() filtering by number of failed attempts
func TestFindConIntentosFallidos(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Arrange
	for i := 1; i <= 5; i++ {
		id := "user_fallos_" + string(rune(47+i))
		cred := createTestCredenciales(id)
		created, _ := repo.Crear(ctx, cred)

		// Add some failed attempts
		for j := 0; j < i; j++ {
			created.MarcarIntentoFallido(time.Now())
		}
		repo.Actualizar(ctx, created)
	}

	// Act - Find credentials with intentosFallidos = 3
	result, err := repo.Find(ctx, domain.EspecificacionCredenciales{
		ListaFiltros: []shareddomain.CriterioFiltro{
			{Campo: "intentosFallidos", Operador: "=", Valor: 3},
		},
	}, shareddomain.Paginacion{
		Pagina:       1,
		TamanoPagina: 100,
	})

	// Assert
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	if len(result) != 1 {
		t.Errorf("Find(intentosFallidos=3) returned %d, want 1", len(result))
	}
	if result[0].IntentosFallidos() != 3 {
		t.Errorf("Found credential with %d failures, want 3", result[0].IntentosFallidos())
	}
}

// TestFindConCorreoVerificado tests Find() filtering by email verification status
func TestFindConCorreoVerificado(t *testing.T) {
	db := setupTestDB(t)
	defer db.Migrator().DropTable(&CredencialesModel{})

	repo := NewCredencialesRepositorio(db)
	ctx := context.Background()

	// Arrange - Create credentials, some with verified email
	for i := 1; i <= 5; i++ {
		id := "user_verified_" + string(rune(47+i))
		cred := createTestCredenciales(id)
		created, _ := repo.Crear(ctx, cred)

		// Verify email for every other one
		if i%2 == 0 {
			created.VerificarCorreo()
			repo.Actualizar(ctx, created)
		}
	}

	// Act
	result, err := repo.Find(ctx, domain.EspecificacionCredenciales{
		ListaFiltros: []shareddomain.CriterioFiltro{
			{Campo: "correoVerificado", Operador: "=", Valor: true},
		},
	}, shareddomain.Paginacion{
		Pagina:       1,
		TamanoPagina: 100,
	})

	// Assert
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Find(correoVerificado=true) returned %d, want 2", len(result))
	}
	for _, cred := range result {
		if !cred.CorreoVerificado() {
			t.Errorf("Found credential with unverified email")
		}
	}
}
