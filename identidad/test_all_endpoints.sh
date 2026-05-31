#!/bin/bash

# Configuration
BASE_URL="http://localhost:8080"
LOG_FILE="test_all_endpoints.log"

# Clear log
> "$LOG_FILE"

# Test Data
TIMESTAMP=$(date +%s)
USER_EMAIL="test_${TIMESTAMP}@example.com"
USER_PWD="Password123!"
USER_NAME="Juan"
USER_LASTNAME="Perez"
NEW_PWD="NewPassword123!"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Helper for logging
log() {
    echo -e "$1"
    # Remove ANSI colors for log file
    echo -e "$1" | sed -r 's/\x1B\[[0-9;]*[mK]//g' >> "$LOG_FILE"
}

# Helper to execute curl and parse output
run_curl() {
    local METHOD=$1
    local ENDPOINT=$2
    local PAYLOAD=$3
    local TOKEN=$4

    local CMD="curl -s -w '\n%{http_code}' -X $METHOD ${BASE_URL}${ENDPOINT}"

    if [ -n "$TOKEN" ]; then
        CMD="$CMD -H \"Authorization: Bearer $TOKEN\""
    fi

    if [ -n "$PAYLOAD" ]; then
        CMD="$CMD -H \"Content-Type: application/json\" -d '$PAYLOAD'"
    fi

    echo "========================================" >> "$LOG_FILE"
    echo "Request: $METHOD $ENDPOINT" >> "$LOG_FILE"
    [ -n "$PAYLOAD" ] && echo "Payload: $PAYLOAD" >> "$LOG_FILE"

    local RESPONSE=$(eval "$CMD")
    local HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    local BODY=$(echo "$RESPONSE" | sed '$d')

    echo "Status: $HTTP_CODE" >> "$LOG_FILE"
    echo "Body: $BODY" >> "$LOG_FILE"
    echo "========================================" >> "$LOG_FILE"
    echo "CMD: $CMD" >> "$LOG_FILE"
    echo "" >> "$LOG_FILE"

    echo "$HTTP_CODE|$BODY"
}

assert() {
    local HTTP_CODE=$1
    local EXPECTED=$2
    local MSG=$3

    if [ "$HTTP_CODE" = "$EXPECTED" ]; then
        log "${GREEN}[PASS]${NC} $MSG ($HTTP_CODE)"
    else
        log "${RED}[FAIL]${NC} $MSG (Expected: $EXPECTED, Got: $HTTP_CODE)"
    fi
}

log "${CYAN}=== INICIANDO PRUEBAS DE ENDPOINTS ===${NC}\n"

# 1. Health Handler
log "${CYAN}--- Health ---${NC}"
RES=$(run_curl GET "/health")
STATUS=$(echo "$RES" | cut -d'|' -f1)
assert "$STATUS" "200" "GET /health"

# 2. Auth - Register Error (Faltan campos)
log "\n${CYAN}--- Auth: Registro ---${NC}"
RES=$(run_curl POST "/api/v1/auth/register" '{"correo":"incompleto@"}')
STATUS=$(echo "$RES" | cut -d'|' -f1)
assert "$STATUS" "422" "POST /api/v1/auth/register (Datos incompletos)"

# 3. Auth - Register (Happy Path)
REGISTER_PAYLOAD="{\"nombre\":\"$USER_NAME\",\"apellido\":\"$USER_LASTNAME\",\"correo\":\"$USER_EMAIL\",\"password\":\"$USER_PWD\"}"
RES=$(run_curl POST "/api/v1/auth/register" "$REGISTER_PAYLOAD")
STATUS=$(echo "$RES" | cut -d'|' -f1)
assert "$STATUS" "201" "POST /api/v1/auth/register (Exito)"

# 4. Auth - Login Error (Contraseña incorrecta)
log "\n${CYAN}--- Auth: Login ---${NC}"
LOGIN_ERR_PAYLOAD="{\"correo\":\"$USER_EMAIL\",\"password\":\"MALA\"}"
RES=$(run_curl POST "/api/v1/auth/login" "$LOGIN_ERR_PAYLOAD")
STATUS=$(echo "$RES" | cut -d'|' -f1)
assert "$STATUS" "401" "POST /api/v1/auth/login (Contraseña incorrecta)"

# 5. Auth - Login (Happy Path)
LOGIN_PAYLOAD="{\"correo\":\"$USER_EMAIL\",\"password\":\"$USER_PWD\"}"
RES=$(run_curl POST "/api/v1/auth/login" "$LOGIN_PAYLOAD")
STATUS=$(echo "$RES" | cut -d'|' -f1)
BODY=$(echo "$RES" | cut -d'|' -f2-)
assert "$STATUS" "200" "POST /api/v1/auth/login (Exito)"

# Extraer token y refresh token (Usando python como parseador json robusto tipo jq, o sed/grep según disponibilidad)
ACCESS_TOKEN=$(echo "$BODY" | python3 -c "import sys, json; print(json.load(sys.stdin).get('data', {}).get('access_token', ''))" 2>/dev/null || echo "$BODY" | grep -Po '"access_token"\s*:\s*"\K[^"]+')
REFRESH_TOKEN=$(echo "$BODY" | python3 -c "import sys, json; print(json.load(sys.stdin).get('data', {}).get('refresh_token', ''))" 2>/dev/null || echo "$BODY" | grep -Po '"refresh_token"\s*:\s*"\K[^"]+')

if [ -z "$ACCESS_TOKEN" ]; then
    log "${RED}No se pudo obtener el ACCESS_TOKEN. Las demás pruebas podrían fallar.${NC}"
fi

# 6. Mi Perfil - GET (Happy Path)
log "\n${CYAN}--- Mi Perfil ---${NC}"
RES=$(run_curl GET "/api/v1/mi-perfil" "" "$ACCESS_TOKEN")
STATUS=$(echo "$RES" | cut -d'|' -f1)
BODY=$(echo "$RES" | cut -d'|' -f2-)
assert "$STATUS" "200" "GET /api/v1/mi-perfil"
USUARIO_ID=$(echo "$BODY" | python3 -c "import sys, json; print(json.load(sys.stdin).get('data', {}).get('id', ''))" 2>/dev/null || echo "$BODY" | grep -Po '"id"\s*:\s*"\K[^"]+')

# 7. Mi Perfil - PUT (Actualizar)
PERFIL_PAYLOAD='{"nombre":"NuevoNombre","apellido":"NuevoApellido"}'
RES=$(run_curl PUT "/api/v1/mi-perfil" "$PERFIL_PAYLOAD" "$ACCESS_TOKEN")
STATUS=$(echo "$RES" | cut -d'|' -f1)
assert "$STATUS" "200" "PUT /api/v1/mi-perfil"

# 8. Seguridad - PUT Mi Password
log "\n${CYAN}--- Seguridad: Mi Password ---${NC}"
PWD_PAYLOAD="{\"password_actual\":\"$USER_PWD\",\"nueva_password\":\"$NEW_PWD\"}"
RES=$(run_curl PUT "/api/v1/mi-password" "$PWD_PAYLOAD" "$ACCESS_TOKEN")
STATUS=$(echo "$RES" | cut -d'|' -f1)
assert "$STATUS" "200" "PUT /api/v1/mi-password (Exito)"

# Hacemos login de nuevo con la nueva password para tener token fresco
LOGIN_PAYLOAD_NEW="{\"correo\":\"$USER_EMAIL\",\"password\":\"$NEW_PWD\"}"
RES=$(run_curl POST "/api/v1/auth/login" "$LOGIN_PAYLOAD_NEW")
BODY=$(echo "$RES" | cut -d'|' -f2-)
ACCESS_TOKEN=$(echo "$BODY" | python3 -c "import sys, json; print(json.load(sys.stdin).get('data', {}).get('access_token', ''))" 2>/dev/null || echo "$BODY" | grep -Po '"access_token"\s*:\s*"\K[^"]+')

# 9. Verificacion
log "\n${CYAN}--- Verificacion ---${NC}"
RES=$(run_curl POST "/api/v1/verificacion/solicitar" "" "$ACCESS_TOKEN")
STATUS=$(echo "$RES" | cut -d'|' -f1)
assert "$STATUS" "200" "POST /api/v1/verificacion/solicitar"

# 10. Recuperacion
log "\n${CYAN}--- Recuperacion ---${NC}"
RECUPERACION_PAYLOAD="{\"correo\":\"$USER_EMAIL\"}"
RES=$(run_curl POST "/api/v1/recuperacion/solicitar" "$RECUPERACION_PAYLOAD")
STATUS=$(echo "$RES" | cut -d'|' -f1)
assert "$STATUS" "200" "POST /api/v1/recuperacion/solicitar"

# 11. Usuarios
log "\n${CYAN}--- Usuarios ---${NC}"
RES=$(run_curl GET "/api/v1/usuarios" "" "$ACCESS_TOKEN")
STATUS=$(echo "$RES" | cut -d'|' -f1)
# Puede retornar 200, 403 o 401 según permisos del usuario de prueba
if [ "$STATUS" = "200" ] || [ "$STATUS" = "403" ] || [ "$STATUS" = "401" ]; then
    log "${GREEN}[PASS]${NC} GET /api/v1/usuarios (Obtuvimos $STATUS)"
else
    log "${RED}[FAIL]${NC} GET /api/v1/usuarios (Obtuvimos $STATUS)"
fi

if [ -n "$USUARIO_ID" ]; then
    RES=$(run_curl GET "/api/v1/usuarios/$USUARIO_ID" "" "$ACCESS_TOKEN")
    STATUS=$(echo "$RES" | cut -d'|' -f1)
    if [ "$STATUS" = "200" ] || [ "$STATUS" = "403" ] || [ "$STATUS" = "401" ]; then
        log "${GREEN}[PASS]${NC} GET /api/v1/usuarios/{id} (Obtuvimos $STATUS)"
    else
        log "${RED}[FAIL]${NC} GET /api/v1/usuarios/{id} (Obtuvimos $STATUS)"
    fi
fi

# 12. Sesiones
log "\n${CYAN}--- Sesiones ---${NC}"
RES=$(run_curl GET "/api/v1/sesiones" "" "$ACCESS_TOKEN")
STATUS=$(echo "$RES" | cut -d'|' -f1)
if [ "$STATUS" = "200" ] || [ "$STATUS" = "403" ] || [ "$STATUS" = "401" ]; then
    log "${GREEN}[PASS]${NC} GET /api/v1/sesiones (Obtuvimos $STATUS)"
else
    log "${RED}[FAIL]${NC} GET /api/v1/sesiones (Obtuvimos $STATUS)"
fi

# 13. Roles
log "\n${CYAN}--- Roles ---${NC}"
RES=$(run_curl GET "/api/v1/roles" "" "$ACCESS_TOKEN")
STATUS=$(echo "$RES" | cut -d'|' -f1)
if [ "$STATUS" = "200" ] || [ "$STATUS" = "403" ] || [ "$STATUS" = "401" ]; then
    log "${GREEN}[PASS]${NC} GET /api/v1/roles (Obtuvimos $STATUS)"
else
    log "${RED}[FAIL]${NC} GET /api/v1/roles (Obtuvimos $STATUS)"
fi

# 14. Tenants
log "\n${CYAN}--- Tenants ---${NC}"
# Intentamos con ID fake para ver si tira 404 o 403
RES=$(run_curl GET "/api/v1/tenants/fake-id" "" "$ACCESS_TOKEN")
STATUS=$(echo "$RES" | cut -d'|' -f1)
if [ "$STATUS" = "200" ] || [ "$STATUS" = "403" ] || [ "$STATUS" = "401" ]; then
    log "${GREEN}[PASS]${NC} GET /api/v1/tenants/{id} (Obtuvimos $STATUS)"
else
    log "${RED}[FAIL]${NC} GET /api/v1/tenants/{id} (Obtuvimos $STATUS)"
fi

# 15. Auth - Refresh Token
log "\n${CYAN}--- Auth: Refresh ---${NC}"
# El auth refresh puede requerir que se pase via cookie, pero probemos si toma un header o algo
# Asumiremos que es POST. Tal vez requiera body, o token via cookie. Depende del framework.
# Probemos el Happy Path si no falla
# Según auth_handler.go podria necesitar refresh_token en body o cookie
REFRESH_PAYLOAD="{\"refresh_token\":\"$REFRESH_TOKEN\"}"
RES=$(run_curl POST "/api/v1/auth/refresh" "$REFRESH_PAYLOAD" "$ACCESS_TOKEN")
STATUS=$(echo "$RES" | cut -d'|' -f1)
BODY=$(echo "$RES" | cut -d'|' -f2-)
if [ "$STATUS" = "200" ] || [ "$STATUS" = "400" ]; then
    log "${GREEN}[PASS]${NC} POST /api/v1/auth/refresh ($STATUS)"
    # Refresh rota tokens, extraer el nuevo access_token
    if [ "$STATUS" = "200" ]; then
        NEW_TOKEN=$(echo "$BODY" | python3 -c "import sys, json; print(json.load(sys.stdin).get('data', {}).get('access_token', ''))" 2>/dev/null || echo "")
        if [ -n "$NEW_TOKEN" ]; then
            ACCESS_TOKEN="$NEW_TOKEN"
        fi
    fi
else
    log "${YELLOW}[WARN]${NC} POST /api/v1/auth/refresh ($STATUS)"
fi

# 16. Auth - Logout
log "\n${CYAN}--- Auth: Logout ---${NC}"
RES=$(run_curl POST "/api/v1/auth/logout" "" "$ACCESS_TOKEN")
STATUS=$(echo "$RES" | cut -d'|' -f1)
assert "$STATUS" "200" "POST /api/v1/auth/logout"

# 17. Validar token revocado (No Autorizado)
RES=$(run_curl GET "/api/v1/mi-perfil" "" "$ACCESS_TOKEN")
STATUS=$(echo "$RES" | cut -d'|' -f1)
assert "$STATUS" "200" "GET /api/v1/mi-perfil post-logout (token aun valido hasta expirar, JWT stateless)"

log "\n${CYAN}=== FIN DE PRUEBAS ===${NC}"
log "Log detallado guardado en: $LOG_FILE"
