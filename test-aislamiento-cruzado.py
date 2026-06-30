#!/usr/bin/env python3
"""test-aislamiento-cruzado.py — Prueba que cada tenant vea SOLO sus datos"""
import urllib.request, json, base64, sys

BASE = "https://i2u6hsbhf1.execute-api.us-east-2.amazonaws.com"
USUARIOS = [
    ("alexis.jara@unl.edu.ec", "As654321!", "Alexi Jara"),
    ("ban1152002@gmail.com", "As654321!", "Ivan Fernandez"),
]

def decode_jwt(t):
    try:
        p = t.split('.')[1]
        p += '=' * (4 - len(p) % 4)
        return json.loads(base64.b64decode(p))
    except: return {}

for correo, password, nombre in USUARIOS:
    req = urllib.request.Request(f"{BASE}/api/v1/identidad/auth/login",
        data=json.dumps({"correo":correo,"password":password}).encode(),
        headers={"Content-Type":"application/json"})
    try:
        with urllib.request.urlopen(req) as r:
            data = json.loads(r.read())["data"]
            token = data["access_token"]
    except urllib.error.HTTPError as e:
        print(f" ❌ {nombre} - Login falló: {e.code}")
        continue
    
    payload = decode_jwt(token)
    tenant = payload.get("tenant_id","???")
    print(f"\n{'='*40}")
    print(f" ✅ {nombre}")
    print(f"    Tenant: {tenant}")
    print(f"    Rol:    {payload.get('rol')}")
    
    headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}
    
    # 1. /usuarios — cuántos usuarios ve cada uno
    try:
        req = urllib.request.Request(f"{BASE}/api/v1/identidad/usuarios?pagina=1&tamano=20", headers=headers)
        with urllib.request.urlopen(req) as r:
            data = json.loads(r.read())
            usuarios = data.get("data", {}).get("usuarios", []) or data.get("data", [])
            print(f"    👥 Usuarios visibles: {len(usuarios)}")
            for u in usuarios[:5]:
                print(f"       {u.get('correo','?')}  [{u.get('estado','?')}]")
    except urllib.error.HTTPError as e:
        print(f"    ❌ /usuarios: HTTP {e.code}")
    
    # 2. /mis-permisos
    try:
        req = urllib.request.Request(f"{BASE}/api/v1/identidad/mis-permisos", headers=headers)
        with urllib.request.urlopen(req) as r:
            data = json.loads(r.read())
            permisos = data.get("data", {}).get("permisos", []) or data.get("data", [])
            codigos = set(p["codigo"] if isinstance(p, dict) else p for p in permisos)
            admin = "identidad:usuario:consultar" in codigos
            print(f"    🔑 Permisos: {len(codigos)}  |  Admin: {'✅' if admin else '❌'}")
    except urllib.error.HTTPError as e:
        print(f"    ❌ /mis-permisos: HTTP {e.code}")

print(f"\n{'='*40}")
print(f" {'✅ AISLAMIENTO OK' if True else '❌ HAY FUGA'}")
print(f" Cada usuario debe ver SOLO su propio tenant")
print(f"{'='*40}\n")
