package main

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTPublisher struct {
	client mqtt.Client
}

// NewMQTTPublisher crea un nuevo cliente MQTT
func NewMQTTPublisher(brokerURL string) (*MQTTPublisher, error) {
	opts := mqtt.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID("image-service").
		SetAutoReconnect(true)

	client := mqtt.NewClient(opts)
	token := client.Connect()

	if !token.WaitTimeout(5000 * 1000000) { // 5 segundos
		return nil, fmt.Errorf("timeout conectando a MQTT")
	}

	if token.Error() != nil {
		return nil, fmt.Errorf("error conectando a MQTT: %v", token.Error())
	}

	log.Println("Conectado a MQTT broker")

	return &MQTTPublisher{client: client}, nil
}

// Publish publica datos en un topic MQTT
func (p *MQTTPublisher) Publish(topic string, payload []byte) error {
	token := p.client.Publish(topic, 1, false, payload)

	if !token.WaitTimeout(5000 * 1000000) { // 5 segundos
		return fmt.Errorf("timeout publicando en MQTT")
	}

	if token.Error() != nil {
		return fmt.Errorf("error publicando: %v", token.Error())
	}

	log.Printf("Publicado en %s (%d bytes)", topic, len(payload))
	return nil
}

// Close cierra la conexión MQTT
func (p *MQTTPublisher) Close() {
	p.client.Disconnect(250)
	log.Println("MQTT desconectado")
}
