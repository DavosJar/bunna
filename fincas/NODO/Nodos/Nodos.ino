#include "esp_camera.h"
#include <WiFi.h>
#include <HTTPClient.h>
#include <WiFiClientSecure.h>

// ==========================================
// PARÁMETROS CONFIGURABLES — solo tocar aquí
// ==========================================
const char* WIFI_SSID        = "CELERITY_GUICE";
const char* WIFI_PASSWORD    = "CECI.@1954m";
const char* API_ENDPOINT = "https://i2u6hsbhf1.execute-api.us-east-2.amazonaws.com";
const char* NODE_API_KEY     = "bunna-fincaPrueba";
const int   CAPTURE_INTERVAL = 10000;  // ms — 10s para pruebas

// ==========================================
// PINES AI-THINKER — no tocar
// ==========================================
#define PWDN_GPIO_NUM     32
#define RESET_GPIO_NUM    -1
#define XCLK_GPIO_NUM      0
#define SIOD_GPIO_NUM     26
#define SIOC_GPIO_NUM     27
#define Y9_GPIO_NUM       35
#define Y8_GPIO_NUM       34
#define Y7_GPIO_NUM       39
#define Y6_GPIO_NUM       36
#define Y5_GPIO_NUM       21
#define Y4_GPIO_NUM       19
#define Y3_GPIO_NUM       18
#define Y2_GPIO_NUM        5
#define VSYNC_GPIO_NUM    25
#define HREF_GPIO_NUM     23
#define PCLK_GPIO_NUM     22

void initCamera() {
  camera_config_t config;
  config.ledc_channel = LEDC_CHANNEL_0;
  config.ledc_timer   = LEDC_TIMER_0;
  config.pin_d0       = Y2_GPIO_NUM;
  config.pin_d1       = Y3_GPIO_NUM;
  config.pin_d2       = Y4_GPIO_NUM;
  config.pin_d3       = Y5_GPIO_NUM;
  config.pin_d4       = Y6_GPIO_NUM;
  config.pin_d5       = Y7_GPIO_NUM;
  config.pin_d6       = Y8_GPIO_NUM;
  config.pin_d7       = Y9_GPIO_NUM;
  config.pin_xclk     = XCLK_GPIO_NUM;
  config.pin_pclk     = PCLK_GPIO_NUM;
  config.pin_vsync    = VSYNC_GPIO_NUM;
  config.pin_href     = HREF_GPIO_NUM;
  config.pin_sscb_sda = SIOD_GPIO_NUM;
  config.pin_sscb_scl = SIOC_GPIO_NUM;
  config.pin_pwdn     = PWDN_GPIO_NUM;
  config.pin_reset    = RESET_GPIO_NUM;
  config.xclk_freq_hz = 20000000;
  config.pixel_format = PIXFORMAT_JPEG;
  config.frame_size   = FRAMESIZE_VGA;  // 640x480 — compatible con YOLO
  config.jpeg_quality = 12;             // 0-63, menor = mejor calidad
  config.fb_count     = 1;

  esp_err_t err = esp_camera_init(&config);
  if (err != ESP_OK) {
    Serial.printf("[CAMARA] Error al iniciar: 0x%x\n", err);
    ESP.restart();
  }
  Serial.println("[CAMARA] Iniciada correctamente");
}

void connectWiFi() {
  Serial.printf("[WIFI] Conectando a %s", WIFI_SSID);
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.printf("\n[WIFI] Conectado — IP: %s\n", WiFi.localIP().toString().c_str());
}

void captureAndSend() {
  // 1. Capturar frame
  camera_fb_t* fb = esp_camera_fb_get();
  if (!fb) {
    Serial.println("[CAMARA] Error al capturar imagen");
    return;
  }
  Serial.printf("[CAMARA] Imagen capturada — %d bytes\n", fb->len);

  // 2. Verificar WiFi
  if (WiFi.status() != WL_CONNECTED) {
    Serial.println("[WIFI] Sin conexión — descartando frame");
    esp_camera_fb_return(fb);
    return;
  }

  // 3. Construir multipart/form-data manualmente
  String boundary = "----BunnaEdge";
  String bodyStart = "--" + boundary + "\r\n"
                   + "Content-Disposition: form-data; name=\"archivo\"; filename=\"foto.jpg\"\r\n"
                   + "Content-Type: image/jpeg\r\n\r\n";
  String bodyEnd = "\r\n--" + boundary + "--\r\n";

  int totalLen = bodyStart.length() + fb->len + bodyEnd.length();

  // 4. Configurar el Cliente Seguro (Vital para HTTPS)
  WiFiClientSecure client;
  client.setInsecure(); // <-- Ignora la validación estricta del certificado de DuckDNS

  // 5. Enviar HTTP POST
  HTTPClient http;
  http.begin(client, API_ENDPOINT); // <-- Usamos el cliente seguro
  http.addHeader("Content-Type", "multipart/form-data; boundary=" + boundary);
  http.addHeader("X-Node-Key", NODE_API_KEY);
  http.setTimeout(30000);

  // 6. Usar PSRAM para no ahogar la RAM interna del ESP32
  uint8_t* payload = (uint8_t*)ps_malloc(totalLen);
  if (!payload) {
    Serial.println("[HTTP] Error al reservar memoria PSRAM para payload");
    esp_camera_fb_return(fb);
    return;
  }

  memcpy(payload, bodyStart.c_str(), bodyStart.length());
  memcpy(payload + bodyStart.length(), fb->buf, fb->len);
  memcpy(payload + bodyStart.length() + fb->len, bodyEnd.c_str(), bodyEnd.length());

  esp_camera_fb_return(fb);  // liberar buffer de cámara antes de enviar

  // 7. POST
  Serial.println("[HTTP] Enviando imagen...");
  int httpCode = http.POST(payload, totalLen);
  free(payload);

  if (httpCode == 200) {
    Serial.println("[HTTP] Diagnóstico recibido correctamente");
  } else {
    Serial.printf("[HTTP] Error — código: %d\n", httpCode);
  }

  http.end();
}
void setup() {
  Serial.begin(115200);
  initCamera();
  connectWiFi();
}

void loop() {
  captureAndSend();
  delay(CAPTURE_INTERVAL);
}