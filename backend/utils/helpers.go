package utils

import (
	"fmt"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
    "io"
	"github.com/google/uuid"
)

// GenerateUniqueFileName genera un nombre de archivo único
func GenerateUniqueFileName(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	return fmt.Sprintf("%s%s", uuid.New().String(), ext)
}

// GenerateIDTransaccion genera un ID de transacción único para cada método de pago
func GenerateIDTransaccion(metodo string) string {
	timestamp := time.Now().Unix()
	randomPart := rand.Intn(100000)

	// Prefijo según el método de pago
	prefix := "txn"
	switch metodo {
	case "tarjeta":
		prefix = "card"
	case "paypal":
		prefix = "pp"
	case "coinbase":
		prefix = "cb"
	case "transferencia":
		prefix = "trf"
	case "dev":
		prefix = "dev"
	case "stripe":
		prefix = "ch"
	case "mercadopago":
		prefix = "mp"
	}

	return fmt.Sprintf("%s_%d_%d", prefix, timestamp, randomPart)
}

// CalculateDateRanges calcula los rangos de fechas para estadísticas según el período seleccionado
func CalculateDateRanges(period string) (startDate, endDate, prevStartDate, prevEndDate time.Time) {
	now := time.Now()

	switch period {
	case "week":
		// Período actual: últimos 7 días
		endDate = now
		startDate = now.AddDate(0, 0, -7)
		// Período anterior: 7 días antes del período actual
		prevEndDate = startDate
		prevStartDate = startDate.AddDate(0, 0, -7)
	case "month":
		// Período actual: últimos 30 días
		endDate = now
		startDate = now.AddDate(0, 0, -30)
		// Período anterior: 30 días antes del período actual
		prevEndDate = startDate
		prevStartDate = startDate.AddDate(0, 0, -30)
	case "year":
		// Período actual: último año
		endDate = now
		startDate = now.AddDate(-1, 0, 0)
		// Período anterior: año anterior al período actual
		prevEndDate = startDate
		prevStartDate = startDate.AddDate(-1, 0, 0)
	default:
		// Por defecto, usar vista mensual
		endDate = now
		startDate = now.AddDate(0, 0, -30)
		prevEndDate = startDate
		prevStartDate = startDate.AddDate(0, 0, -30)
	}

	return startDate, endDate, prevStartDate, prevEndDate
}

// InitRand inicializa el generador de números aleatorios
func InitRand() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt genera un número aleatorio entre 0 y max-1
func RandomInt(max int) int {
	return rand.Intn(max)
}

// StringToInt convierte un string a int, retornando un valor por defecto si falla
func StringToInt(s string, defaultValue int) (int, error) {
	if s == "" {
		return defaultValue, nil
	}
	
	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue, err
	}
	
	return i, nil
}

// PathMatchesPattern verifica si un path coincide con un patrón
func PathMatchesPattern(path, pattern string) bool {
	// Implementación simple: verifica si el path comienza con el patrón
	return strings.HasPrefix(path, pattern)
}

// GetEnv obtiene una variable de entorno o retorna un valor por defecto


// SaveUploadedFile guarda un archivo subido en el sistema de archivos
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Crear directorio si no existe
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// GenerateUUID genera un UUID v4
func GenerateUUID() string {
	return uuid.New().String()
}

// FileExists verifica si un archivo existe
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// EnsureDir asegura que un directorio exista

// RemoveFile elimina un archivo del sistema
func RemoveFile(path string) error {
	return os.Remove(path)
}

// GetFileExtension obtiene la extensión de un archivo
func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// IsValidImageExtension verifica si la extensión es de imagen válida
func IsValidImageExtension(ext string) bool {
	valid := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	ext = strings.ToLower(ext)
	for _, v := range valid {
		if ext == v {
			return true
		}
	}
	return false
}

// GenerateRandomString genera una cadena aleatoria de longitud n
func GenerateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// ParseTime parsea un string a time.Time
func ParseTime(timeStr string, layout string) (time.Time, error) {
	return time.Parse(layout, timeStr)
}

// FormatTime formatea un time.Time a string
func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// CurrentTimestamp obtiene el timestamp actual en segundos
func CurrentTimestamp() int64 {
	return time.Now().Unix()
}

// AddDaysToTime añade días a una fecha
func AddDaysToTime(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// DaysBetween calcula los días entre dos fechas
func DaysBetween(start, end time.Time) int {
	return int(end.Sub(start).Hours() / 24)
}