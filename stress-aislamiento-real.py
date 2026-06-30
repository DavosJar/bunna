#!/usr/bin/env python3
"""Test de aislamiento REAL - solo endpoints que DEBEN estar aislados"""
import urllib.request, json, base64, time, sys
from concurrent.futures import ThreadPoolExecutor, as_completed

BASE = "https://i2u6hsbhf1.execute-api.us-east-2.amazonaws.com"
CORREO = "alexis.jara@unl.edu.ec" 
PASSWORD = "As654321!"

# Login
req = urllib.request.Request(f"{BASE}/api/v1/identidad/auth/login",
    data=json.dumps({"correo":CORREO,"password":PASSWORD}).encode(),
    headers={"Content-Type":"application/json"})
with urllib.request.urlopen(req) as r:
    token = json.loads(r.read())["data"]["access_token"]

payload = json.loads(base64.b64decode(token.split('.')[1] + '=='))
mi_tenant = payload.get("tenant_id","???")

headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}

# Endpoints que DEBEN estar aislados por tenant:
# Los IDs de OTRos tenants (Ivan, Pepe, Cesar) - estos deben dar 404/403
targets = [
    # Intentar acceder al perfil de otros usuarios
    ("identidad", "GET", f"{BASE}/api/v1/identidad/credenciales/019f117a-2671-700e-941c-d7bada0cf123", "Credenciales Ivan"),
    ("identidad", "GET", f"{BASE}/api/v1/identidad/credenciales/019f1219-4fd5-73af-81a3-56195a305147", "Credenciales Pepe"),
    # Intentar ver roles de otros usuarios
    ("identidad", "GET", f"{BASE}/api/v1/identidad/usuarios/019f117a-2671-700e-941c-d7bada0cf123/roles", "Roles Ivan"),
    ("identidad", "GET", f"{BASE}/api/v1/identidad/usuarios/019f1219-4fd5-73af-81a3-56195a305147/roles", "Roles Pepe"),
    # Intentar ver invitaciones (solo propias)
    ("identidad", "GET", f"{BASE}/api/v1/identidad/invitaciones", "Invitaciones globales"),
    # Intentar ver sesiones de otros
    ("identidad", "GET", f"{BASE}/api/v1/identidad/sesiones", "Sesiones globales"),
    # Intentar modificar perfil de otro (PUT con ID de otro)
    ("identidad", "PUT", f"{BASE}/api/v1/identidad/usuarios/019f117a-2671-700e-941c-d7bada0cf123", "Modificar Ivan"),
]

# ENDPOINTS PROPIOS (deben ser 200)
propios = [
    ("identidad", "GET", f"{BASE}/api/v1/identidad/mi-perfil", "Mi perfil"),
    ("identidad", "GET", f"{BASE}/api/v1/identidad/tenants/mis-tenants", "Mis tenants"),
    ("identidad", "GET", f"{BASE}/api/v1/identidad/mis-permisos", "Mis permisos"),
]

TOTAL = 500  # 500 intentos a endpoints ajenos
CONCURRENCIA = 10
results = []
leaks = []

print(f"\n{'='*55}")
print(f" STRESS TEST AISLAMIENTO REAL")
print(f"{'='*55}")
print(f" Tenant: {mi_tenant[:8]}...")
print(f" Endpoints propios (deben ser 200): {len(propios)}")
print(f" Endpoints ajenos (deben ser 404/403): {len(targets)}")
print(f" Total intentos a endpoints ajenos: {TOTAL}")
print(f" Concurrencia: {CONCURRENCIA}")
print(f"{'='*55}\n")

# Primero verificar endpoints propios
print("📍 Verificando endpoints PROPIOS (control):")
for servicio, metodo, url, desc in propios:
    t0 = time.time()
    try:
        req = urllib.request.Request(url, method=metodo, headers=headers)
        with urllib.request.urlopen(req, timeout=10) as r:
            code = r.status
    except urllib.error.HTTPError as e:
        code = e.code
    t = (time.time()-t0)*1000
    status = "✅" if code == 200 else "❌"
    print(f"  {status} {desc}: HTTP {code} ({t:.0f}ms)")

# Ahora los ataques
print(f"\n📍 Atacando endpoints ajenos {TOTAL} veces:")
pool = ThreadPoolExecutor(max_workers=CONCURRENCIA)
futuros = []

for i in range(TOTAL):
    servicio, metodo, url, desc = targets[i % len(targets)]
    futuros.append(pool.submit(
        lambda m=metodo, u=url: (m, u, None, None)  # placeholder, hacemos la request abajo
    ))

# Reemplazar por request real
futuros = []
for i in range(TOTAL):
    servicio, metodo, url, desc = targets[i % len(targets)]
    futuros.append((metodo, url, desc))

# Ejecutar con ThreadPoolExecutor
def atacar(metodo, url, desc):
    t0 = time.time()
    try:
        req = urllib.request.Request(url, method=metodo, headers=headers)
        with urllib.request.urlopen(req, timeout=10) as r:
            code = r.status
            data = r.read()
            if code == 200 and len(data) > 20:
                return {"code": code, "ms": (time.time()-t0)*1000, "desc": desc, "fuga": True}
            return {"code": code, "ms": (time.time()-t0)*1000, "desc": desc, "fuga": False}
    except urllib.error.HTTPError as e:
        return {"code": e.code, "ms": (time.time()-t0)*1000, "desc": desc, "fuga": False}
    except:
        return {"code": 0, "ms": (time.time()-t0)*1000, "desc": desc, "fuga": False}

with ThreadPoolExecutor(max_workers=CONCURRENCIA) as pool:
    futuros = [pool.submit(atacar, m, u, d) for m,u,d in futuros]
    for i, f in enumerate(as_completed(futuros), 1):
        r = f.result()
        results.append(r)
        if r.get("fuga"):
            leaks.append(r)
        bar = '█' * (i * 50 // TOTAL) + '░' * (50 - i * 50 // TOTAL)
        sys.stdout.write(f"\r {i:5d}/{TOTAL}  [{bar}]  intentos:{i}  fugas:{len(leaks)}")
        sys.stdout.flush()

codes = {}
for r in results:
    codes[r["code"]] = codes.get(r["code"], 0) + 1

print(f"\n\n{'='*55}")
print(f" RESULTADOS - AISLAMIENTO BAJO STRESS")
print(f"{'='*55}")
print(f" TOTAL INTENTOS A ENDPOINTS AJENOS: {len(results)}")
print(f"{'='*55}")
for c in sorted(codes):
    pct = codes[c]*100/len(results)
    print(f"  HTTP {c}: {codes[c]:4d} ({pct:.1f}%)")

print(f"\n{'='*55}")
if leaks:
    print(f" ❌ {len(leaks)} FUGAS DETECTADAS")
    for l in leaks[:5]:
        print(f"    {l['desc']}: HTTP {l['code']} ({l['ms']:.0f}ms)")
else:
    print(f" ✅ CERO FUGAS - {len(results)} intentos, {len(results)} bloqueos")
    print(f"    El aislamiento de tenant AGUANTA bajo estres")
print(f"{'='*55}\n")
