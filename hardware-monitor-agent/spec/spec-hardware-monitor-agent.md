---
title: Hardware Monitor Agent — Especificación de Agente Sidecar para Monitoreo de Hardware en Docker Swarm
version: 1.0
date_created: 2026-06-19
owner: Bunna Infrastructure Team
tags: rust, monitoring, docker-swarm, kafka, hardware, sidecar
---

# Introducción

Especificación del agente de monitoreo de hardware diseñado para correr como sidecar en contenedores Docker dentro de un Swarm. Recolecta métricas del host anfitrión, las evalúa contra umbrales configurables y publica tanto métricas periódicas como alertas inmediatas a Apache Kafka.

## 1. Propósito y Alcance

**Propósito**: Proveer monitoreo de hardware en tiempo real para nodos Docker Swarm mediante un agente ligero que lee métricas del host, las agrega, las evalúa contra umbrales y las publica en Kafka para consumo por sistemas externos de visualización y alertamiento.

**Alcance**:
- Recolección de métricas de CPU, RAM, disco, red y contenedores desde `/host/proc/`
- Agregación periódica en snapshots unificados
- Evaluación contra umbrales configurables (WARN y CRITICAL)
- Publicación a Kafka de métricas (cada intervalo) y alertas (inmediatas)
- Control de inundación de alertas con detección de recovery
- Ejecución como sidecar en Docker Swarm

**Fuera de alcance**:
- Visualización de métricas (se consume desde otra fuente)
- Persistencia local de métricas históricas
- Orquestación de los tópicos de Kafka
- Autenticación/autorización de clientes Kafka

**Audiencia**: Equipos de infraestructura y desarrollo que despliegan y operan el agente.

## 2. Definiciones

| Término | Definición |
|---------|-----------|
| **Sidecar** | Contenedor auxiliar que acompaña al contenedor principal en un mismo pod/task de Swarm |
| **Snapshot** | Conjunto completo de métricas de todas las fuentes en un instante de tiempo |
| **Intervalo** | Período de recolección configurable (default 10 s) |
| **WARN** | Nivel de alerta que excede un umbral menor; solo genera log local |
| **CRITICAL** | Nivel de alerta que excede un umbral mayor; publica a Kafka inmediatamente |
| **Recovery** | Evento que indica que una métrica ha vuelto a niveles normales tras haber estado en CRITICAL |
| **Inundación de alertas** | Publicación repetitiva de la misma alerta mientras la métrica se mantiene sobre el umbral |
| **node_id** | Identificador único del nodo Swarm, inyectado vía variable de entorno |
| **cgroup** | Mecanismo del kernel Linux para limitar y aislar recursos de procesos/contenedores |
| **/host/proc/** | Punto de montaje del filesystem /proc del host dentro del contenedor sidecar |

## 3. Requisitos, Restricciones y Directrices

### Requisitos Funcionales

- **REQ-FUNC-001**: El agente debe recolectar CPU, RAM, disco, red y contenedores desde `/host/proc/`
- **REQ-FUNC-002**: La recolección debe ser paralela usando tokio async tasks
- **REQ-FUNC-003**: El intervalo de recolección debe ser configurable, con default de 10 s
- **REQ-FUNC-004**: Todas las métricas de un intervalo deben consolidarse en un snapshot unificado
- **REQ-FUNC-005**: El snapshot completo debe publicarse en el tópico `hardware.metrics` cada intervalo
- **REQ-FUNC-006**: Las alertas CRITICAL deben publicarse en el tópico `hardware.alerts` inmediatamente
- **REQ-FUNC-007**: Las alertas WARN deben solo registrar en log local, sin publicar a Kafka
- **REQ-FUNC-008**: El agente debe implementar control de inundación: una alerta CRITICAL se publica una sola vez hasta que la métrica retorne a nivel normal y luego la supere nuevamente
- **REQ-FUNC-009**: El agente debe publicar un evento de recovery cuando una métrica en CRITICAL retorna a niveles normales
- **REQ-FUNC-010**: Los thresholds deben ser configurables independientemente para CPU, RAM y disco
- **REQ-FUNC-011**: El node_id debe poderse inyectar vía variable de entorno

### Requisitos No Funcionales

- **REQ-NF-001**: El agente debe ser un binario Rust ligero y estático
- **REQ-NF-002**: La imagen Docker debe construirse con Dockerfile multistage para minimizar tamaño
- **REQ-NF-003**: La comunicación entre componentes internos debe usar `tokio::sync::mpsc`
- **REQ-NF-004**: El agente debe manejar errores de lectura de `/host/proc/` sin colapsar (graceful degradation)
- **REQ-NF-005**: El agente debe publicar mensajes JSON a Kafka con esquema consistente

### Restricciones

- **CON-001**: El agente corre como sidecar; debe acceder a `/host/proc/` del host mediante bind mount
- **CON-002**: No se debe almacenar estado en disco local
- **CON-003**: No se debe exponer API HTTP ni otro servicio de red (solo Kafka como salida)
- **CON-004**: Los tópicos de Kafka deben llamarse `hardware.metrics` y `hardware.alerts`

### Directrices

- **GUD-001**: Cada fuente de métrica debe ser un módulo independiente dentro de `collectors/`
- **GUD-002**: El agregador debe esperar a todos los collectors antes de construir el snapshot
- **GUD-003**: El guardian debe mantener un mapa de estado por métrica para control de inundación
- **GUD-004**: Los errores transitorios de lectura (e.g., archivo temporalmente no disponible) deben logging, no panic

## 4. Interfaces y Contratos de Datos

### 4.1 Estructura del Snapshot (tópico `hardware.metrics`)

```json
{
  "node_id": "swarm-node-01",
  "timestamp": "2026-06-19T14:30:00Z",
  "interval_ms": 10000,
  "cpu": {
    "usage_percent": 45.2,
    "cores": 8
  },
  "ram": {
    "total_mb": 16384,
    "used_mb": 8192,
    "available_mb": 8192,
    "usage_percent": 50.0
  },
  "disks": [
    {
      "mount": "/",
      "total_gb": 256.0,
      "used_gb": 128.0,
      "available_gb": 128.0,
      "usage_percent": 50.0
    }
  ],
  "net": {
    "interfaces": [
      {
        "name": "eth0",
        "received_bytes": 1024000,
        "transmitted_bytes": 512000,
        "received_bytes_per_sec": 102400.0,
        "transmitted_bytes_per_sec": 51200.0
      }
    ]
  },
  "containers": [
    {
      "container_id": "abc123",
      "cpu_shares": 512,
      "memory_limit_mb": 1024
    }
  ]
}
```

### 4.2 Estructura de Alerta (tópico `hardware.alerts`)

```json
{
  "node_id": "swarm-node-01",
  "timestamp": "2026-06-19T14:30:05Z",
  "metric": "cpu",
  "severity": "CRITICAL",
  "value": 92.5,
  "threshold": 90.0,
  "message": "CPU usage at 92.5% exceeds CRITICAL threshold of 90.0%",
  "previous_state": "normal",
  "event_type": "alert"
}
```

### 4.3 Estructura de Recovery (tópico `hardware.alerts`)

```json
{
  "node_id": "swarm-node-01",
  "timestamp": "2026-06-19T14:35:00Z",
  "metric": "cpu",
  "severity": "INFO",
  "value": 45.2,
  "threshold": 90.0,
  "message": "CPU usage returned to normal: 45.2%",
  "previous_state": "critical",
  "event_type": "recovery"
}
```

### 4.4 Canales Internos (tokio::sync::mpsc)

| Canal | Tipo | Emisor | Receptor | Capacidad |
|-------|------|--------|----------|-----------|
| metrics_tx | `mpsc::Sender<HashMap<MetricType, MetricData>>` | Cada collector | Aggregator | 32 |
| snapshot_tx | `mpsc::Sender<Snapshot>` | Aggregator | Guardian | 16 |
| alert_tx | `mpsc::Sender<Alert>` | Guardian | Publisher | 64 |

### 4.5 Configuración (variables de entorno)

| Variable | Ejemplo | Descripción |
|----------|---------|-------------|
| `NODE_ID` | `swarm-node-01` | Identificador único del nodo |
| `KAFKA_BROKERS` | `kafka:9092` | Lista de brokers Kafka |
| `INTERVAL_MS` | `10000` | Intervalo de recolección en ms |
| `CPU_WARN_PERCENT` | `80.0` | Umbral WARN de CPU (%) |
| `CPU_CRITICAL_PERCENT` | `90.0` | Umbral CRITICAL de CPU (%) |
| `RAM_WARN_PERCENT` | `85.0` | Umbral WARN de RAM (%) |
| `RAM_CRITICAL_PERCENT` | `95.0` | Umbral CRITICAL de RAM (%) |
| `DISK_WARN_PERCENT` | `85.0` | Umbral WARN de disco (%) |
| `DISK_CRITICAL_PERCENT` | `95.0` | Umbral CRITICAL de disco (%) |
| `RUST_LOG` | `info` | Nivel de log de tracing |

## 5. Criterios de Aceptación

- **AC-001**: Dado un nodo Swarm con el sidecar desplegado, cuando transcurre un intervalo, entonces se publica un snapshot JSON en `hardware.metrics`
- **AC-002**: Dada una métrica que supera el umbral CRITICAL, cuando persiste por un tick, entonces se publica una alerta en `hardware.alerts` una sola vez
- **AC-003**: Dada una métrica en estado CRITICAL, cuando retorna a nivel normal, entonces se publica un evento de recovery en `hardware.alerts`
- **AC-004**: Dada una métrica que supera el umbral WARN pero no CRITICAL, entonces solo se registra en log local sin publicación a Kafka
- **AC-005**: Dado un error de lectura de `/host/proc/stat`, cuando el archivo no existe o está corrupto, entonces el agente continúa operando con las demás métricas
- **AC-006**: Dado el contenedor sin acceso a `/host/proc/`, cuando se inicia el agente, entonces falla con un mensaje de error claro

## 6. Estrategia de Pruebas

- **Nivel Unitario**: Cada collector con datos simulados (archivos proc dummy)
- **Nivel Integración**: Agregador + Guardian con canales mpsc en memoria
- **Nivel Integración**: publicador Kafka contra mock o container local
- **Cobertura mínima**: 80% en módulos core (guardian, aggregator, tipos)
- **CI/CD**: Pruebas en GitHub Actions con servicio Kafka usando github-action kafka container

## 7. Justificación y Contexto

- Se elige `tokio::sync::mpsc` como mecanismo de comunicación interna por su integración nativa con async Rust y backpressure controlada
- El control de inundación se implementa con un `HashMap<MetricType, AlertState>` para evitar reintentos innecesarios y reducir carga en Kafka
- No se expone API HTTP porque la visualización se realiza externamente (Grafana, dashboard propio, etc.)
- Los thresholds se configuran por variable de entorno para facilitar el despliegue en Swarm sin archivos de configuración montados
- Se usa `/host/proc/` en lugar de `/proc` para cumplir con el aislamiento de contenedores: el host monta `/proc` en `/host/proc/` dentro del sidecar

## 8. Dependencias e Integraciones Externas

### Sistemas externos
- **Apache Kafka**: Sistema de mensajería para publicación de métricas y alertas

### Dependencias de infraestructura
- **Docker Swarm**: Orquestador de contenedores donde se despliega el sidecar
- **/host/proc/**: Bind mount del filesystem proc del host

### Dependencias de plataforma
- **Rust edition 2024**: Versión del lenguaje
- **tokio v1**: Runtime asíncrono
- **rdkafka v0.36**: Cliente Kafka con soporte tokio

## 9. Ejemplos y Casos Borde

### Ciclo de alerta completo

```
Estado normal (CPU 30%) → sin acción
Estado normal (CPU 95%) → supera CRITICAL → publica alerta → estado: alert_sent
Estado alert_sent (CPU 96%) → sigue sobre umbral → NO publica (flood control)
Estado alert_sent (CPU 50%) → baja del umbral → publica recovery → estado: normal
Estado normal (CPU 92%) → supera CRITICAL → publica alerta → ciclo se repite
```

### Error transitorio de lectura

Si `/host/proc/stat` falla una vez, el collector de CPU devuelve `None` o un error. El agregador omite esa métrica en el snapshot parcial pero no bloquea el ciclo. En el siguiente intervalo se reintenta.

## 10. Criterios de Validación

- [ ] Todos los tests unitarios pasan
- [ ] El binario compila en modo release sin warnings
- [ ] La imagen Docker ocupa menos de 50 MB
- [ ] El agente inicia y publica métricas a Kafka sin errores
- [ ] El control de inundación envía una sola alerta por evento
- [ ] El recovery se publica correctamente al normalizarse la métrica
- [ ] Un error en un collector no detiene el ciclo principal

## 11. Especificaciones Relacionadas

- N/A — especificación inicial del proyecto
