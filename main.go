package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

type Word struct {
	Datouno       string `json:"datouno"`
	Datodos       string `json:"datodos"`
	Description   string `json:"description"`
	Pronunciacion string `json:"pronunciacion"`
}

func initDB() {
	dsn := "balcore64:Djandelo30096446#@tcp(122.8.182.16:3306)/corpo_recibos"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error verificando la conexión a la base de datos: %v", err)
	}
	fmt.Println("Conexión exitosa a la base de datos")
}

func getWordsByLetter(letter string) ([]Word, error) {
	query := "SELECT datouno, datodos, descripcion, pronunciacion FROM tb_palabras WHERE datouno LIKE ?"
	rows, err := db.Query(query, letter+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		var word Word
		if err := rows.Scan(&word.Datouno, &word.Datodos, &word.Description, &word.Pronunciacion); err != nil {
			return nil, err
		}
		words = append(words, word)
	}
	return words, nil
}

func wordsHandler(w http.ResponseWriter, r *http.Request) {
	letter := r.URL.Query().Get("letter")
	if letter == "" {
		http.Error(w, "Falta la letra", http.StatusBadRequest)
		return
	}

	words, err := getWordsByLetter(letter)
	if err != nil {
		http.Error(w, "Error obteniendo las palabras", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(words)
}

func serveHTML(templateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, templateName)
	}
}

func traducirPalabra(palabra, direccion string) (string, error) {
	query := ""
	if direccion == "esp-qeqchi" {
		query = "SELECT datodos FROM tb_palabras WHERE datouno = ?"
	} else if direccion == "qeqchi-esp" {
		query = "SELECT datouno FROM tb_palabras WHERE datodos = ?"
	} else {
		return "", fmt.Errorf("dirección de traducción no válida")
	}

	row := db.QueryRow(query, palabra)

	var traduccion string
	err := row.Scan(&traduccion)
	if err != nil {
		return "", err
	}

	return traduccion, nil
}

func traducirHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	palabraOriginal := r.FormValue("textoInput")
	direccion := r.FormValue("opciones")

	palabraTraducida, err := traducirPalabra(palabraOriginal, direccion)
	if err != nil {
		http.Error(w, "Error en la traducción", http.StatusInternalServerError)
		return
	}

	data := struct {
		DatoUno    string
		Traduccion string
	}{
		DatoUno:    palabraOriginal,
		Traduccion: palabraTraducida,
	}

	tmpl, err := template.ParseFiles("static/templates/index.html")
	if err != nil {
		http.Error(w, "Error cargando la plantilla", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error ejecutando la plantilla", http.StatusInternalServerError)
		return
	}
}

func main() {
	initDB()
	defer db.Close()

	// Servir archivos estáticos
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handlers para las páginas HTML
	http.HandleFunc("/", serveHTML("static/templates/index.html"))
	http.HandleFunc("/diccionario.html", serveHTML("static/templates/diccionario.html"))

	// Handlers para las funcionalidades adicionales
	http.HandleFunc("/words", wordsHandler)
	http.HandleFunc("/traducir", traducirHandler)

	fmt.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
