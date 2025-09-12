// Analizador léxico en Go para reconocer tokens y lexemas
// Identifica palabras reservadas, variables, operadores, signos de agrupación, números y errores
// Muestra línea y columna de cada token y una tabla de conteo al final
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	// Ruta del archivo de entrada (puedes cambiarlo por argumento dinámico si lo deseas)
	ruta := "entrada.txt"
	if !archivoExiste(ruta) {
		fmt.Println("❌ Error: El archivo no existe o la ruta es incorrecta.")
		return
	}

	fmt.Println("Archivo encontrado. Analizando contenido...")
	leerArchivo(ruta)
}

// Verifica si el archivo existe en la ruta dada
func archivoExiste(ruta string) bool {
	absPath, err := filepath.Abs(ruta)
	if err != nil {
		return false
	}
	_, err = os.Stat(absPath)
	return !os.IsNotExist(err)
}

// Tipos de token reconocidos por el analizador
type TokenType string

const (
	RESERVADA  TokenType = "Palabra Reservada"
	VARIABLE   TokenType = "Variable"
	OPERADOR   TokenType = "Operador"
	AGRUPACION TokenType = "Signo de Agrupación"
	NUMERO     TokenType = "Número"
	ERROR      TokenType = "Error"
)

// Palabras reservadas del lenguaje
var palabrasReservadas = map[string]bool{
	"if": true, "else": true, "for": true, "while": true, "return": true, "func": true, "var": true, "package": true, "import": true,
}

// Operadores soportados
var operadores = map[string]bool{
	"+": true, "-": true, "*": true, "/": true, "=": true, "==": true, "!=": true, "<": true, ">": true, "<=": true, ">=": true,
}

// Signos de agrupación
var agrupaciones = map[rune]bool{
	'(': true, ')': true, '{': true, '}': true, '[': true, ']': true,
	';': true, // punto y coma
	'"': true, // comillas dobles
}

// Estructura que representa un token identificado
type Token struct {
	Tipo    TokenType // Tipo de token
	Lexema  string    // Texto del token
	Linea   int       // Línea donde se encontró
	Columna int       // Columna donde inicia
}

// Retorna true si el caracter es una letra o guion bajo
func esLetra(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

// Retorna true si el caracter es un dígito
func esDigito(r rune) bool {
	return r >= '0' && r <= '9'
}

// Analiza una línea de texto y retorna los tokens encontrados
func analizarLinea(texto string, linea int) []Token {
	tokens := []Token{}
	runes := []rune(texto)
	i := 0
	for i < len(runes) {
		r := runes[i]
		// Ignorar espacios y tabulaciones
		if r == ' ' || r == '\t' {
			i++
			continue
		}
		col := i + 1 // Columna donde inicia el token
		// Palabra reservada o variable
		if esLetra(r) {
			ini := i
			i++
			for i < len(runes) && (esLetra(runes[i]) || esDigito(runes[i])) {
				i++
			}
			lex := string(runes[ini:i])
			if palabrasReservadas[lex] {
				tokens = append(tokens, Token{RESERVADA, lex, linea, col})
			} else {
				tokens = append(tokens, Token{VARIABLE, lex, linea, col})
			}
			continue
		}
		// Número
		if esDigito(r) {
			ini := i
			i++
			for i < len(runes) && esDigito(runes[i]) {
				i++
			}
			tokens = append(tokens, Token{NUMERO, string(runes[ini:i]), linea, col})
			continue
		}
		// Operadores de dos caracteres (ej: ==, !=, <=, >=)
		if i+1 < len(runes) {
			op2 := string(runes[i : i+2])
			if operadores[op2] {
				tokens = append(tokens, Token{OPERADOR, op2, linea, col})
				i += 2
				continue
			}
		}
		// Operador de un caracter
		if operadores[string(r)] {
			tokens = append(tokens, Token{OPERADOR, string(r), linea, col})
			i++
			continue
		}
		// Signos de agrupación
		if agrupaciones[r] {
			tokens = append(tokens, Token{AGRUPACION, string(r), linea, col})
			i++
			continue
		}
		// Si no es nada válido, es error léxico
		tokens = append(tokens, Token{ERROR, string(r), linea, col})
		i++
	}
	return tokens
}

// Lee el archivo línea por línea, analiza y muestra los tokens y la tabla de conteo
func leerArchivo(ruta string) {
	archivo, err := os.Open(ruta)
	if err != nil {
		fmt.Println("❌ No se pudo abrir el archivo:", err)
		return
	}
	defer archivo.Close()

	scanner := bufio.NewScanner(archivo)
	linea := 1
	todosTokens := []Token{}
	for scanner.Scan() {
		texto := scanner.Text()
		// Analiza la línea y obtiene los tokens
		tokens := analizarLinea(texto, linea)
		for _, t := range tokens {
			fmt.Printf("[Línea %d, Col %d] %-22s : %s\n", t.Linea, t.Columna, t.Tipo, t.Lexema)
		}
		todosTokens = append(todosTokens, tokens...)
		linea++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("❌ Error al leer el archivo:", err)
		return
	}

	// Tabla de conteo por tipo de token (ordenada para mejor visualización)
	conteo := map[TokenType]int{}
	for _, t := range todosTokens {
		conteo[t.Tipo]++
	}
	tipos := []TokenType{}
	for tipo := range conteo {
		tipos = append(tipos, tipo)
	}
	sort.Slice(tipos, func(i, j int) bool { return string(tipos[i]) < string(tipos[j]) })
	fmt.Println("\n--- Tabla de conteo por tipo de token ---")
	for _, tipo := range tipos {
		fmt.Printf("%-22s : %d\n", tipo, conteo[tipo])
	}
}
