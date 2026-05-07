package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"lab_08/models"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

func main() {
	var err error
	//connStr := "postgres://postgres:postgres@localhost:5432/contacts_db"
	connStr := "postgres://postgres:postgres@db:5432/contacts_db"
	db, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal("не вдалося підключитися до БД:", err)
	}
	defer db.Close(context.Background())

	migration, err := os.ReadFile("migrations/001_create_contacts.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(context.Background(), string(migration))
	if err != nil {
		log.Fatal("помилка міграції:", err)
	}

	r := chi.NewRouter()
	r.Get("/contacts", getContacts)
	r.Get("/contacts/{id}", getContact)
	r.Post("/contacts", createContact)
	r.Put("/contacts/{id}", updateContact)
	r.Delete("/contacts/{id}", deleteContact)

	fmt.Println("Сервер запущено на :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getContacts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(context.Background(), "SELECT id, name, phone FROM contacts")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	contacts := []models.Contact{}
	for rows.Next() {
		var c models.Contact
		if err := rows.Scan(&c.ID, &c.Name, &c.Phone); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		contacts = append(contacts, c)
	}

	w.Header().Set("Content-Type", "application/json")
	result := make([]byte, 0)
	result = append(result, '[')
	for i, c := range contacts {
		data, _ := c.MarshalJSON()
		result = append(result, data...)
		if i < len(contacts)-1 {
			result = append(result, ',')
		}
	}
	result = append(result, ']')
	w.Write(result)
}

func getContact(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "невірний id", http.StatusBadRequest)
		return
	}

	var c models.Contact
	err = db.QueryRow(context.Background(),
			  "SELECT id, name, phone FROM contacts WHERE id = $1", id).
	Scan(&c.ID, &c.Name, &c.Phone)
	if err != nil {
		http.Error(w, "контакт не знайдено", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	data, _ := c.MarshalJSON()
	w.Write(data)
}

func createContact(w http.ResponseWriter, r *http.Request) {
	var c models.Contact
	if err := c.UnmarshalJSON(readBody(r)); err != nil {
		http.Error(w, "невірний формат даних", http.StatusBadRequest)
		return
	}

	err := db.QueryRow(context.Background(),
			   "INSERT INTO contacts (name, phone) VALUES ($1, $2) RETURNING id",
			   c.Name, c.Phone).Scan(&c.ID)
			   if err != nil {
				   http.Error(w, err.Error(), http.StatusInternalServerError)
				   return
			   }

			   w.Header().Set("Content-Type", "application/json")
			   w.WriteHeader(http.StatusCreated)
			   data, _ := c.MarshalJSON()
			   w.Write(data)
}

func updateContact(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "невірний id", http.StatusBadRequest)
		return
	}

	var c models.Contact
	if err := c.UnmarshalJSON(readBody(r)); err != nil {
		http.Error(w, "невірний формат даних", http.StatusBadRequest)
		return
	}
	c.ID = id

	_, err = db.Exec(context.Background(),
			 "UPDATE contacts SET name = $1, phone = $2 WHERE id = $3",
		  c.Name, c.Phone, c.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	data, _ := c.MarshalJSON()
	w.Write(data)
}

func deleteContact(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "невірний id", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(context.Background(),
			 "DELETE FROM contacts WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func readBody(r *http.Request) []byte {
	buf := make([]byte, r.ContentLength)
	r.Body.Read(buf)
	return buf
}
