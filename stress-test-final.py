#!/usr/bin/env python3
"""1000 requests, 20/s, 1000 usuarios random, 1 password"""
import urllib.request, time, sys, json, threading
from concurrent.futures import ThreadPoolExecutor, as_completed

URL = "https://i2u6hsbhf1.execute-api.us-east-2.amazonaws.com/api/v1/identidad/auth/login"
TOTAL = 1000
POR_SEGUNDO = 20
CONCURRENCIA = 20
PASSWORD = "As654321!"

results = []
lock = threading.Lock()
start_global = time.time()

def do_request(i):
    correo = f"stress-final-{i:04d}@bunna.app"
    body = json.dumps({"correo": correo, "password": PASSWORD}).encode()
    t0 = time.time()
    try:
        req = urllib.request.Request(URL, data=body, method="POST")
        req.add_header("Content-Type", "application/json")
        with urllib.request.urlopen(req, timeout=10) as r:
            code = r.status
    except urllib.error.HTTPError as e:
        code = e.code
    except:
        code = 0
    elapsed = time.time() - t0
    with lock:
        results.append((code, elapsed, i, correo))

print(f"\n{'='*50}")
print(f" Stress Test - 1000 usuarios @ 20 req/s")
print(f"{'='*50}")
print(f" URL:         {URL}")
print(f" Total:       {TOTAL}")
print(f" Rate:        {POR_SEGUNDO}/s")
print(f" Password:    {PASSWORD}")
print(f" Usuarios:    stress-final-0000 .. stress-final-0999")
print(f"{'='*50}\n")

pool = ThreadPoolExecutor(max_workers=CONCURRENCIA)
enviados = 0
t0_lote = time.time()

while enviados < TOTAL:
    ahora = time.time()
    pasaron = ahora - t0_lote
    deben_ir = min(int(pasaron * POR_SEGUNDO), TOTAL - enviados)
    
    if deben_ir > 0:
        for _ in range(deben_ir):
            enviados += 1
            pool.submit(do_request, enviados)
        t0_lote = ahora
    
    bar = '█' * (enviados * 50 // TOTAL) + '░' * (50 - enviados * 50 // TOTAL)
    sys.stdout.write(f"\r Enviados: {enviados:4d}/{TOTAL}  [{bar}]  {time.time()-start_global:.0f}s")
    sys.stdout.flush()
    time.sleep(0.05)

pool.shutdown(wait=True)

total_time = time.time() - start_global
codes = {}
times = [r[1] for r in results]
times.sort()

errores_422 = [r for r in results if r[0]==422]
errores_401 = [r for r in results if r[0]==401]

print(f"\n\n{'='*50}")
print(f" RESULTADOS")
print(f"{'='*50}")
print(f" Total:       {len(results)}")
print(f" Duración:    {total_time:.1f}s")
print(f" Rate real:   {len(results)/total_time:.1f} req/s")
for c in sorted(set(r[0] for r in results)):
    print(f"   HTTP {c}: {sum(1 for r in results if r[0]==c)}")

print(f"\n Tiempos:")
print(f"   Promedio:  {sum(times)/len(times)*1000:.1f}ms")
print(f"   Mínimo:    {times[0]*1000:.1f}ms")
print(f"   Máximo:    {times[-1]*1000:.1f}ms")
print(f"   P50:       {times[len(times)//2]*1000:.1f}ms")
print(f"   P90:       {times[int(len(times)*.9)]*1000:.1f}ms")
print(f"   P99:       {times[int(len(times)*.99)]*1000:.1f}ms")

print(f"\n 422 vs 401:")
print(f"   422 (cuenta no existe): {len(errores_422)}")
print(f"   401 (password malo):    {len(errores_401)}")
print(f"{'='*50}\n")
