from fastapi import FastAPI
from routers import diagnostico

app = FastAPI(
    title="Modelo YOLO API",
    description="API para diagnóstico usando modelo YOLO",
    version="1.0.0",
)

app.include_router(diagnostico.router, prefix="/api/v1")


@app.get("/health")
async def salud():
    return {"estado": "ok"}
