package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Lenguaje struct {
	ID          int    `json:"id"`
	PalabraEsp  string `json:"palabra_esp"`
	PalabraQeq  string `json:"palabra_qeq"`
	Descripcion string `json:"descripcion"`
}

func main() {

	db, err := sql.Open("mysql", "usuario:contraseña@tcp(127.0.0.1:3306)/bd_traductor")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	// Manejador para la ruta /traducir
	router.HandleFunc("/traducir", func(w http.ResponseWriter, r *http.Request) {
		// Obtener los parámetros de la solicitud
		idioma, _ := strconv.Atoi(r.URL.Query().Get("idioma"))
		palabra := r.URL.Query().Get("palabra")

		// Realizar la consulta a la base de datos
		var lenguaje Lenguaje
		var consulta string
		if idioma == 1 {
			consulta = fmt.Sprintf("SELECT palabra_qeq FROM tb_lengua WHERE palabra_esp = '%s'", palabra)
		} else if idioma == 2 {
			consulta = fmt.Sprintf("SELECT palabra_esp FROM tb_lengua WHERE palabra_qeq = '%s'", palabra)
		} else {
			http.Error(w, "Idioma no válido", http.StatusBadRequest)
			return
		}

		err := db.QueryRow(consulta).Scan(&lenguaje.PalabraQeq)
		if err != nil {
			http.Error(w, "Palabra no encontrada", http.StatusNotFound)
			return
		}

		// Codificar la respuesta como JSON y enviarla
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lenguaje)
	}).Methods("GET")

	// Configurar el servidor HTTP
	fmt.Println("Servidor en ejecución en el puerto 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
