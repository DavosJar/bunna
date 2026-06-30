#!/usr/bin/env python3
"""1000 intentos de romper aislamiento de tenant - solo Alexis"""
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

# Targets de OTROS tenants (conocidos de la DB)
targets = [
    ("GET", f"{BASE}/api/v1/identidad/usuarios/019f117a-2671-700e-941c-d7bada0cf123", "Ivan user"),
    ("GET", f"{BASE}/api/v1/identidad/usuarios/019f1219-4fd5-73af-81a3-56195a305147", "Pepe user"),
    ("GET", f"{BASE}/api/v1/identidad/tenants/019f1142-f701-7ab6-8e00-6c21ff30d578", "Cesar tenant"),
    ("GET", f"{BASE}/api/v1/identidad/tenants/019f117a-26b5-723c-b678-28777432bc96", "Ivan tenant"),
    ("GET", f"{BASE}/api/v1/identidad/usuarios?pagina=1&tamano=50", "todos los usuarios"),
]

TOTAL = 1000
CONCURRENCIA = 20
results = []
leaks = []

print(f"\n{'='*55}")
print(f" STRESS TEST DE AISLAMIENTO - 1000 intentos")
print(f"{'='*55}")
print(f" Tenant: {mi_tenant}")
print(f" Targets: 5 endpoints de OTROS tenants")
print(f" Requests: {TOTAL} ({TOTAL//len(targets)} c/u)")
print(f" Concurrentes: {CONCURRENCIA}")
print(f"{'='*55}\n")

def try_fuga(metodo, url, desc):
    t0 = time.time()
    try:
        req = urllib.request.Request(url, method=metodo, headers=headers)
        with urllib.request.urlopen(req, timeout=10) as r:
            code = r.status
            data = json.loads(r.read())
    except urllib.error.HTTPError as e:
        code = e.code
        data = {"code": e.code}
    except:
        code = 0
        data = {}
    t = (time.time()-t0)*1000
    return {"code": code, "ms": t, "desc": desc, "url": url}

pool = ThreadPoolExecutor(max_workers=CONCURRENCIA)
futuros = []
for i in range(TOTAL):
    metodo, url, desc = targets[i % len(targets)]
    futuros.append(pool.submit(try_fuga, metodo, url, desc))

for i, f in enumerate(as_completed(futuros), 1):
    r = f.result()
    results.append(r)
    if r["code"] == 200:
        leaks.append(r)
    bar = '█' * (i * 50 // TOTAL) + '░' * (50 - i * 50 // TOTAL)
    sys.stdout.write(f"\r {i:5d}/{TOTAL}  [{bar}]  leaks:{len(leaks)}  {r['ms']:.0f}ms")
    sys.stdout.flush()

pool.shutdown()

total_time = sum(r["ms"] for r in results) / 1000
codes = {}
for r in results:
    codes[r["code"]] = codes.get(r["code"], 0) + 1

print(f"\n\n{'='*55}")
print(f" RESULTADOS")
print(f"{'='*55}")
print(f" Total intentos: {len(results)}")
print(f" Tiempo total:   {time.time()-time.mktime(time.localtime()):.1f}s (estimado)")
print(f"{'='*55}")
print(f" CODIGOS HTTP:")
for c in sorted(codes):
    pct = codes[c]*100/TOTAL
    nombre = {200: "✅ DATOS FILTRADOS", 403: "🚫 BLOQUEADO", 404: "✅ NO ENCONTRADO", 0: "💥 ERROR"}.get(c, f"OTRO {c}")
    print(f"  HTTP {c}: {codes[c]:4d} ({pct:.1f}%)  {nombre}")

print(f"\n{'='*55}")
if leaks:
    print(f" ❌ FUGA DE DATOS DETECTADA!")
    print(f"    {len(leaks)} requests devolvieron datos de otros tenants")
    for l in leaks[:5]:
        print(f"     - {l['desc']}: {l['code']} ({l['ms']:.0f}ms)")
else:
    print(f" ✅ CERO fugas en {TOTAL} intentos agresivos")
    print(f"    El aislamiento aguanta bajo estres")
print(f"{'='*55}\n")
