# Image Service

Microservicio Go para procesar y publicar imágenes vía MQTT.

## Características

- Lectura concurrente de imágenes desde directorio
- Redimensionamiento a 640px de ancho (aspecto preservado)
- Publicación en MQTT
- Control de concurrencia configurable
- **MQTT broker permanentemente activo** (no se desconecta al terminar de procesar)
- Monitorea periódicamente nuevas imágenes

## Uso

```bash
# Compilar
go build -o image-service

# Ejecutar
./image-service -dir ./images -broker tcp://localhost:1883 -workers 4
```

### Flags

- `-dir` - Directorio con imágenes (default: `./images`)
- `-broker` - URL del broker MQTT (default: `tcp://localhost:1883`)
- `-workers` - Máximo de workers concurrentes (default: 4)

## Ejemplo

```bash
# Con broker MQTT en localhost
mkdir -p images
# Copiar imágenes a ./images
./image-service

# Con broker remoto y 8 workers
./image-service -dir /path/to/images -broker tcp://mqtt.example.com:1883 -workers 8
```

## MQTT Topics

Las imágenes procesadas se publican en:
```
images/processed/{nombre_archivo}
```

## Dependencias

```
github.com/eclipse/paho.mqtt.golang - Cliente MQTT
github.com/nfnt/resize - Redimensionamiento de imágenes
```
