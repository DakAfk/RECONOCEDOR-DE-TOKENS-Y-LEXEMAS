package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	ruta := "entrada.txt" // Puedes cambiar esto por un argumento dinámico
	if !archivoExiste(ruta) {
		fmt.Println("❌ Error: El archivo no existe o la ruta es incorrecta.")
		return
	}

	fmt.Println("Archivo encontrado. Leyendo contenido...")
	leerArchivo(ruta)
}

func archivoExiste(ruta string) bool {
	absPath, err := filepath.Abs(ruta)
	if err != nil {
		return false
	}
	_, err = os.Stat(absPath)
	return !os.IsNotExist(err)
}
func leerArchivo(ruta string) {
	archivo, err := os.Open(ruta)
	if err != nil {
		fmt.Println("❌ No se pudo abrir el archivo:", err)
		return
	}
	defer archivo.Close()

	scanner := bufio.NewScanner(archivo)
	linea := 1
	for scanner.Scan() {
		texto := scanner.Text()
		fmt.Printf("Línea %d: %s\n", linea, texto)
		linea++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("❌ Error al leer el archivo:", err)
	}
}
