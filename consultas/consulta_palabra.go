package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Palabra struct {
	Palabra       string `json:"palabra"`
	SignificadoEs string `json:"significado_es"`
	SignificadoEn string `json:"significado_en"`
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/api/palabra/", GetPalabra)

	fmt.Println("Servidor escuchando en http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

func GetPalabra(w http.ResponseWriter, r *http.Request) {

	palabra := r.URL.Path[len("/api/palabra/"):]

	db, err := sql.Open("mysql", "usuario:contraseña@tcp(localhost:3306)/basededatos")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT palabra, significado_es, significado_en FROM palabras WHERE palabra = '%s'", palabra)

	// Ejecutar la consulta
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var resultado Palabra

	for rows.Next() {
		err := rows.Scan(&resultado.Palabra, &resultado.SignificadoEs, &resultado.SignificadoEn)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Verificar si se encontró la palabra
	if resultado.Palabra == "" {
		http.Error(w, "Palabra no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado)
}
