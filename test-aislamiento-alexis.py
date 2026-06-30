#!/usr/bin/env python3
"""Test de aislamiento - Alexis intenta ver datos de OTROS tenants"""
import urllib.request, json, base64, time

BASE = "https://i2u6hsbhf1.execute-api.us-east-2.amazonaws.com"
CORREO = "alexis.jara@unl.edu.ec"
PASSWORD = "As654321!"

# Login
req = urllib.request.Request(f"{BASE}/api/v1/identidad/auth/login",
    data=json.dumps({"correo":CORREO,"password":PASSWORD}).encode(),
    headers={"Content-Type":"application/json"})
with urllib.request.urlopen(req) as r:
    data = json.loads(r.read())["data"]
    token = data["access_token"]

# Decodificar JWT
payload = json.loads(base64.b64decode(token.split('.')[1] + '=='))
mi_tenant = payload.get("tenant_id","???")
print(f"\n{'='*55}")
print(f" TEST DE AISLAMIENTO - {CORREO}")
print(f"{'='*55}")
print(f" Mi tenant ID: {mi_tenant}")
print(f" Mi rol:       {payload.get('rol')}")
print(f"{'='*55}\n")

headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}

def test(nombre, metodo, url, esperado="200"):
    t0 = time.time()
    try:
        req = urllib.request.Request(url, method=metodo, headers=headers)
        with urllib.request.urlopen(req, timeout=10) as r:
            data = json.loads(r.read())
            t = (time.time()-t0)*1000
            return True, data, t
    except urllib.error.HTTPError as e:
        t = (time.time()-t0)*1000
        return False, {"code": e.code}, t

# 1. Ver los usuarios que devuelve el sistema
print("1️⃣  /usuarios — ¿cuántos usuarios ve Alexis?")
ok, data, t = test("usuarios", "GET", f"{BASE}/api/v1/identidad/usuarios?pagina=1&tamano=50")
if ok:
    usuarios = data.get("data", {}).get("usuarios", []) or data.get("data", [])
    print(f"    Responde en {t:.0f}ms")
    print(f"    Total usuarios visibles: {len(usuarios)}")
    for u in usuarios:
        print(f"      - {u.get('correo','?'):35s}  tenant: {u.get('tenant_id','?')[:8] if u.get('tenant_id') else 'N/A'}")
    # Verificar si hay usuarios de OTROS tenants
    otros = [u for u in usuarios if u.get('tenant_id') and u['tenant_id'] != mi_tenant]
    if otros:
        print(f"\n    ⚠️  VIÓ {len(otros)} USUARIOS DE OTROS TENANTS!")
        for u in otros:
            print(f"       {u.get('correo','?')}  -> tenant {u['tenant_id'][:8]}...")
    else:
        print(f"\n    ✅ SOLO usuarios de mi tenant")
else:
    print(f"    ❌ HTTP {data['code']} ({t:.0f}ms)")

# 2. Intentar acceder a tenants CONOCIDOS de otros usuarios
print(f"\n2️⃣  Intento acceder a tenants de OTROS usuarios:")
otros_tenants = [
    ("Cesar Lopez", "019f1142-f701-7ab6-8e00-6c21ff30d578"),
    ("ivan Fernandez", "019f117a-26b5-723c-b678-28777432bc96"),
    ("PEpe Armijos", "019f1219-4fd5-73b2-a195-fa0ca0913157"),
]
for nombre, tid in otros_tenants:
    # Intentar GET a config tenant (si existe endpoint)
    ok, data, t = test(f"configurar-{nombre}", "GET", f"{BASE}/api/v1/identidad/tenants/{tid}")
    if ok:
        print(f"    ⚠️  PUDO ACCEDER a tenant de {nombre} (HTTP 200)")
    else:
        print(f"    ✅ Bloqueado al acceder a {nombre}: HTTP {data['code']} ({t:.0f}ms)")

# 3. Ver usuarios de otros tenants por ID directo
print(f"\n3️⃣  Intento ver usuarios de otros tenants por ID:")
# IDs de usuarios conocidos:
otros_usuarios = [
    ("Ivan", "019f117a-2671-700e-941c-d7bada0cf123"),
    ("Pepe", "019f1219-4fd5-73af-81a3-56195a305147"),
]
for nombre, uid in otros_usuarios:
    ok, data, t = test(f"usuario-{nombre}", "GET", f"{BASE}/api/v1/identidad/usuarios/{uid}")
    if ok:
        print(f"    ⚠️  PUDO VER usuario de {nombre}")
    else:
        print(f"    ✅ Bloqueado al ver {nombre}: HTTP {data['code']} ({t:.0f}ms)")

print(f"\n{'='*55}")
print(f" CONCLUSION")
print(f"{'='*55}")
print(f" Alexis (tenant: {mi_tenant[:8]}...)")
print(f" {'✅ AISLAMIENTO FUNCIONA' if False else '❌ REVISAR RESULTADOS'}")
print(f"{'='*55}\n")
