package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	imageDir := flag.String("dir", "./images", "Directorio con imágenes")
	mqttBroker := flag.String("broker", "tcp://localhost:1883", "URL del broker MQTT")
	maxWorkers := flag.Int("workers", 4, "Número máximo de workers concurrentes")
	flag.Parse()

	log.Printf("Iniciando Image Service")
	log.Printf("Directorio: %s", *imageDir)
	log.Printf("Broker MQTT: %s", *mqttBroker)
	log.Printf("Workers: %d", *maxWorkers)

	// Inicializar cliente MQTT
	publisher, err := NewMQTTPublisher(*mqttBroker)
	if err != nil {
		log.Fatalf("Error conectando a MQTT: %v", err)
	}

	// Canal para capturar señales de terminación
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Procesar imágenes iniciales
	processImages(*imageDir, publisher, *maxWorkers)

	// Mantener servicio activo esperando nuevas imágenes
	log.Println("Servicio activo. MQTT broker disponible. Ctrl+C para salir.")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			log.Println("Señal de terminación recibida")
			publisher.Close()
			os.Exit(0)
		case <-ticker.C:
			// Verificar periódicamente por nuevas imágenes
			images, err := GetImages(*imageDir)
			if err == nil && len(images) > 0 {
				processImages(*imageDir, publisher, *maxWorkers)
			}
		}
	}
}

func processImages(imageDir string, publisher *MQTTPublisher, maxWorkers int) {
	// Obtener lista de imágenes
	images, err := GetImages(imageDir)
	if err != nil {
		log.Printf("Error leyendo directorio: %v", err)
		return
	}

	if len(images) == 0 {
		return
	}

	log.Printf("Procesando %d imágenes", len(images))

	// Procesar imágenes concurrentemente
	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, maxWorkers)

	for _, imgPath := range images {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire

		go func(path string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release

			if err := ProcessAndPublish(path, publisher); err != nil {
				log.Printf("Error procesando %s: %v", path, err)
			}
		}(imgPath)
	}

	wg.Wait()
	log.Println("Lote procesado")
}

func ProcessAndPublish(imagePath string, publisher *MQTTPublisher) error {
	// Procesar imagen
	processedData, err := ProcessImage(imagePath, 640)
	if err != nil {
		return err
	}

	// Publicar en MQTT
	filename := getFilename(imagePath)
	topic := fmt.Sprintf("images/processed/%s", filename)
	return publisher.Publish(topic, processedData)
}
