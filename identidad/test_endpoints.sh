#!/bin/bash
BASE_URL="http://localhost:8080"
LOG_FILE="api_tests_happy_path.log"
echo "Iniciando pruebas de la API (Ruta Feliz)..." | tee "$LOG_FILE"
echo "=====================================" >> "$LOG_FILE"
# Datos de prueba
TEST_EMAIL="test_$(date +%s)@bunna.com"
TEST_PASS="Password123!"
FAKE_PASSWORD="WrongPass123!"
# Función de ayuda
test_endpoint() {
    local method=$1
    local path=$2
    local body=$3
    local token=$4
    
    echo "=====================================" >> "$LOG_FILE"
    echo "⚙️ TEST: $method $path" | tee -a "$LOG_FILE"
    
    local response
    local http_code
    
    if [ "$method" = "GET" ] || [ "$method" = "DELETE" ]; then
        response=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X "$method" "$BASE_URL$path" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $token")
    else
        response=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X "$method" "$BASE_URL$path" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $token" \
            -d "$body")
    fi
    
    http_code=$(echo "$response" | grep "HTTP_CODE:" | cut -d':' -f2)
    body_resp=$(echo "$response" | sed '/HTTP_CODE:/d')
    
    if [ "$http_code" -ge 400 ]; then
        echo "⚠️ FALLO (Status: $http_code)" | tee -a "$LOG_FILE"
    else
        echo "✅ ÉXITO (Status: $http_code)" | tee -a "$LOG_FILE"
    fi
    echo "$body_resp" >> "$LOG_FILE"
}
# 1. Health
test_endpoint "GET" "/health" "" ""
# 2. Registrar Usuario
REGISTER_BODY="{\"correo\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASS\",\"nombre\":\"Test\",\"apellido\":\"User\"}"
test_endpoint "POST" "/api/v1/auth/register" "$REGISTER_BODY" ""
# 3. Login para obtener Token
LOGIN_BODY="{\"correo\":\"$TEST_EMAIL\",\"password\":\"$FAKE_PASSWORD\"}"
echo "⚙️ Obteniendo Token JWT..."
LOGIN_RESP=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" -H "Content-Type: application/json" -d "$LOGIN_BODY")
# Extraer el access token (requiere jq, si no tienes jq puedes instalarlo con sudo apt install jq)
REAL_TOKEN=$(echo "$LOGIN_RESP" | grep -o '"access_token":"[^"]*' | grep -o '[^"]*$')
if [ -z "$REAL_TOKEN" ]; then
    echo "❌ No se pudo obtener el token. ¿Está bien el formato del login?"
    echo "Respuesta del Login: $LOGIN_RESP"
    exit 1
fi
echo "✅ Token obtenido correctamente."
# 4. Probar endpoints protegidos usando el token real
# Ahora enviamos JSONs simulados con datos básicos para evitar los Errores 422
test_endpoint "GET" "/api/v1/mi-perfil" "" "$REAL_TOKEN"
test_endpoint "PUT" "/api/v1/mi-perfil" "{\"nombre\":\"Nuevo Nombre\"}" "$REAL_TOKEN"
test_endpoint "GET" "/api/v1/sesiones" "" "$REAL_TOKEN"
test_endpoint "GET" "/api/v1/roles" "" "$REAL_TOKEN"
# Algunos de estos podrían seguir fallando si tu lógica de negocio exige cosas muy específicas 
# (ej: si el usuario recién creado no tiene rol de admin, le dará un 403 Forbidden al intentar ver todos los usuarios)
test_endpoint "GET" "/api/v1/usuarios" "" "$REAL_TOKEN"
echo "====================================="
echo "✅ Pruebas finalizadas. Revisa '$LOG_FILE' para ver el detalle completo."