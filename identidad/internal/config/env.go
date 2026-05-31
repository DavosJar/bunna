package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config contiene toda la configuración de la aplicación cargada desde variables de entorno.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	Port        string
	FrontendURL string
	CORSOrigins string

	APIGatewayEnabled bool

	BcryptCost int

	JWTSecret            string
	JWTIssuer            string
	JWTAccessExpiracion  time.Duration
	JWTRefreshExpiracion time.Duration

	SesionTimeoutInactividad time.Duration
	SesionTimeoutAbsoluto    time.Duration
	SesionMaxRefrescos       int

	BloqueoIPMaxIntentos int
	BloqueoIPVentana     time.Duration
	BloqueoIPDuracion    time.Duration

	RateLimitMaxRequests int
	RateLimitVentana     time.Duration

	CuentaBloqueoMaxIntentos int
	CuentaBloqueoDuracion    time.Duration

	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string

	VerificacionTokenExpiracion time.Duration
	VerificacionMaxReenvios     int
	VerificacionVentanaReenvios time.Duration

	RecuperacionTokenExpiracion     time.Duration
	RecuperacionRateLimitIPMax      int
	RecuperacionRateLimitUsuarioMax int
	RecuperacionRateLimitVentana    time.Duration
}

// LoadConfig carga y valida todas las variables de entorno.
func LoadConfig() (*Config, error) {
	bcryptCost, err := parsarEntero("BCRYPT_COST", "12", 10, 14)
	if err != nil {
		return nil, err
	}
	jwtAccessExpiracion, err := parsarDuracion("JWT_ACCESS_EXPIRACION", "15m")
	if err != nil {
		return nil, err
	}
	jwtRefreshExpiracion, err := parsarDuracion("JWT_REFRESH_EXPIRACION", "24h")
	if err != nil {
		return nil, err
	}
	sesionTimeoutInactividad, err := parsarDuracion("SESION_TIMEOUT_INACTIVIDAD", "30m")
	if err != nil {
		return nil, err
	}
	sesionTimeoutAbsoluto, err := parsarDuracion("SESION_TIMEOUT_ABSOLUTO", "168h")
	if err != nil {
		return nil, err
	}
	sesionMaxRefrescos, err := parsarEnteroSinRango("SESION_MAX_REFRESCOS", "0")
	if err != nil {
		return nil, err
	}
	bloqueoIPMaxIntentos, err := parsarEnteroSinRango("BLOQUEO_IP_MAX_INTENTOS", "20")
	if err != nil {
		return nil, err
	}
	bloqueoIPVentana, err := parsarDuracion("BLOQUEO_IP_VENTANA", "15m")
	if err != nil {
		return nil, err
	}
	bloqueoIPDuracion, err := parsarDuracion("BLOQUEO_IP_DURACION", "30m")
	if err != nil {
		return nil, err
	}
	rateLimitMaxRequests, err := parsarEnteroSinRango("RATE_LIMIT_MAX_REQUESTS", "10")
	if err != nil {
		return nil, err
	}
	rateLimitVentana, err := parsarDuracion("RATE_LIMIT_VENTANA", "1m")
	if err != nil {
		return nil, err
	}

	verificacionTokenExpiracion, err := parsarDuracion("VERIFICACION_TOKEN_EXPIRACION", "24h")
	if err != nil {
		return nil, err
	}
	verificacionMaxReenvios, err := parsarEnteroSinRango("VERIFICACION_MAX_REENVIOS", "5")
	if err != nil {
		return nil, err
	}
	verificacionVentanaReenvios, err := parsarDuracion("VERIFICACION_VENTANA_REENVIOS", "24h")
	if err != nil {
		return nil, err
	}

	recuperacionTokenExpiracion, err := parsarDuracion("RECUPERACION_TOKEN_EXPIRACION", "1h")
	if err != nil {
		return nil, err
	}
	recuperacionRateLimitIPMax, err := parsarEnteroSinRango("RECUPERACION_RATE_LIMIT_IP_MAX", "3")
	if err != nil {
		return nil, err
	}
	recuperacionRateLimitUsuarioMax, err := parsarEnteroSinRango("RECUPERACION_RATE_LIMIT_USUARIO_MAX", "1")
	if err != nil {
		return nil, err
	}
	recuperacionRateLimitVentana, err := parsarDuracion("RECUPERACION_RATE_LIMIT_VENTANA", "15m")
	if err != nil {
		return nil, err
	}

	apiGatewayEnabled, err := strconv.ParseBool(getEnv("API_GATEWAY_ENABLED", "false"))
	if err != nil {
		return nil, fmt.Errorf("API_GATEWAY_ENABLED debe ser un booleano válido: %w", err)
	}

	cuentaBloqueoMaxIntentos, err := parsarEnteroSinRango("CUENTA_BLOQUEO_MAX_INTENTOS", "5")
	if err != nil {
		return nil, err
	}
	cuentaBloqueoDuracion, err := parsarDuracion("CUENTA_BLOQUEO_DURACION", "15m")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "identidad_user"),
		DBPassword: getEnv("DB_PASSWORD", "identidad_pass_dev"),
		DBName:     getEnv("DB_NAME", "identidad_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		Port:        getEnv("PORT", "8080"),
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
		CORSOrigins: getEnv("CORS_ORIGINS", "*"),

		APIGatewayEnabled: apiGatewayEnabled,

		BcryptCost: bcryptCost,

		JWTSecret:            getEnv("JWT_SECRET", ""),
		JWTIssuer:            getEnv("JWT_ISSUER", "ServicioIdentidad"),
		JWTAccessExpiracion:  jwtAccessExpiracion,
		JWTRefreshExpiracion: jwtRefreshExpiracion,

		SesionTimeoutInactividad: sesionTimeoutInactividad,
		SesionTimeoutAbsoluto:    sesionTimeoutAbsoluto,
		SesionMaxRefrescos:       sesionMaxRefrescos,

		BloqueoIPMaxIntentos: bloqueoIPMaxIntentos,
		BloqueoIPVentana:     bloqueoIPVentana,
		BloqueoIPDuracion:    bloqueoIPDuracion,

		RateLimitMaxRequests: rateLimitMaxRequests,
		RateLimitVentana:     rateLimitVentana,

		CuentaBloqueoMaxIntentos: cuentaBloqueoMaxIntentos,
		CuentaBloqueoDuracion:    cuentaBloqueoDuracion,

		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", ""),

		VerificacionTokenExpiracion: verificacionTokenExpiracion,
		VerificacionMaxReenvios:     verificacionMaxReenvios,
		VerificacionVentanaReenvios: verificacionVentanaReenvios,

		RecuperacionTokenExpiracion:     recuperacionTokenExpiracion,
		RecuperacionRateLimitIPMax:      recuperacionRateLimitIPMax,
		RecuperacionRateLimitUsuarioMax: recuperacionRateLimitUsuarioMax,
		RecuperacionRateLimitVentana:    recuperacionRateLimitVentana,
	}

	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD es requerido")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET es requerido")
	}

	return cfg, nil
}

// GetDSN retorna el Data Source Name para conectarse a PostgreSQL.
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s client_encoding=UTF8",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

// getEnv retorna el valor de la variable de entorno o el valor por defecto.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parsarEntero parsea una variable de entorno como entero con rango válido.
func parsarEntero(key, defaultValue string, min, max int) (int, error) {
	val, err := strconv.Atoi(getEnv(key, defaultValue))
	if err != nil {
		return 0, fmt.Errorf("%s debe ser un número válido: %w", key, err)
	}
	if val < min || val > max {
		return 0, fmt.Errorf("%s debe estar entre %d y %d, se obtuvo: %d", key, min, max, val)
	}
	return val, nil
}

// parsarEnteroSinRango parsea una variable de entorno como entero sin validación de rango.
func parsarEnteroSinRango(key, defaultValue string) (int, error) {
	val, err := strconv.Atoi(getEnv(key, defaultValue))
	if err != nil {
		return 0, fmt.Errorf("%s debe ser un número válido: %w", key, err)
	}
	return val, nil
}

// parsarDuracion parsea una variable de entorno como time.Duration.
func parsarDuracion(key, defaultValue string) (time.Duration, error) {
	val, err := time.ParseDuration(getEnv(key, defaultValue))
	if err != nil {
		return 0, fmt.Errorf("%s debe ser una duración válida (ej: 15m, 24h): %w", key, err)
	}
	return val, nil
}
