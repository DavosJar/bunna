#!/bin/bash

# Colores para la salida
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Iniciando configuración del stack de monitoreo...${NC}"

# 1. Crear tópicos en Kafka
echo -e "${BLUE}1/2 Creando tópicos en Kafka...${NC}"
topics=("hardware.metrics" "hardware.alerts" "telemetry")

for topic in "${topics[@]}"; do
    sudo docker exec bunna-kafka kafka-topics --create --bootstrap-server localhost:9092 --topic "$topic" --partitions 1 --replication-factor 1 2>/dev/null
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}  ✔ Tópico '$topic' creado.${NC}"
    else
        echo -e "  ℹ Tópico '$topic' ya existe o hubo un aviso."
    fi
done

# 2. Inicializar ClickHouse
echo -e "${BLUE}2/2 Inicializando tablas en ClickHouse...${NC}"
if [ -f "servicio-monitoreo/config/clickhouse/init/001_tablas.sql" ]; then
    sudo docker exec -i bunna-clickhouse clickhouse-client < servicio-monitoreo/config/clickhouse/init/001_tablas.sql
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}  ✔ Tablas creadas en ClickHouse (BD: bunna_monitoreo).${NC}"
    else
        echo -e "  ✖ Error al crear las tablas en ClickHouse."
    fi
else
    echo -e "  ✖ No se encontró el archivo de esquema SQL."
fi

echo -e "${BLUE}Reiniciando Ingestor para aplicar cambios...${NC}"
sudo docker compose -f docker-compose.monitoring.yml up -d servicio-monitoreo

echo -e "${GREEN}¡Configuración completada! Verifica los logs con: sudo docker logs bunna-ingestor${NC}"
