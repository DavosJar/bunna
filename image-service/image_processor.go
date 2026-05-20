package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

// ProcessImage redimensiona una imagen a 640px de ancho manteniendo aspecto
func ProcessImage(imagePath string, maxWidth uint) ([]byte, error) {
	// Leer archivo
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("error abriendo imagen: %v", err)
	}
	defer file.Close()

	// Decodificar imagen según extensión
	var img image.Image
	ext := strings.ToLower(filepath.Ext(imagePath))

	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return nil, fmt.Errorf("formato no soportado: %s", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("error decodificando imagen: %v", err)
	}

	// Redimensionar a 640px de ancho
	resized := resize.Resize(maxWidth, 0, img, resize.Lanczos3)

	// Codificar a bytes
	processedPath := imagePath + ".resized"
	outFile, err := os.Create(processedPath)
	if err != nil {
		return nil, fmt.Errorf("error creando archivo temporal: %v", err)
	}
	defer outFile.Close()

	if ext == ".png" {
		err = png.Encode(outFile, resized)
	} else {
		err = jpeg.Encode(outFile, resized, &jpeg.Options{Quality: 90})
	}

	if err != nil {
		return nil, fmt.Errorf("error codificando imagen: %v", err)
	}

	// Leer bytes procesados
	data, err := ioutil.ReadFile(processedPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo procesado: %v", err)
	}

	// Limpiar archivo temporal
	os.Remove(processedPath)

	return data, nil
}

// GetImages obtiene lista de archivos de imagen del directorio
func GetImages(dirPath string) ([]string, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var images []string
	validExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}

	for _, f := range files {
		if !f.IsDir() {
			ext := strings.ToLower(filepath.Ext(f.Name()))
			if validExts[ext] {
				fullPath := filepath.Join(dirPath, f.Name())
				images = append(images, fullPath)
			}
		}
	}

	return images, nil
}

func getFilename(path string) string {
	return filepath.Base(path)
}
