package main

import (
	"log"
	"strconv"
	"sync"

	"example.com/lab_07/models" 
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mailru/easyjson"
)

var (
	notes     = make(map[int]models.Note)
	idCounter = 1
	mu        sync.Mutex
	validate  = validator.New()
)

func main() {
	app := fiber.New()

	app.Get("/notes", getNotes)
	app.Get("/notes/:id", getNote)
	app.Post("/notes", createNote)
	app.Put("/notes/:id", updateNote)
	app.Delete("/notes/:id", deleteNote)

	log.Println("Сервер запущено на http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

func sendError(c *fiber.Ctx, status int, message string) error {
	resp := models.ErrorResponse{Error: message}
	data, _ := easyjson.Marshal(&resp)
	c.Set("Content-Type", "application/json")
	return c.Status(status).SendString(string(data) + "\n")
}

func getNotes(c *fiber.Ctx) error {
	mu.Lock()
	defer mu.Unlock()

	list := make(models.NoteList, 0, len(notes))
	for _, n := range notes {
		list = append(list, n)
	}

	data, err := easyjson.Marshal(&list)
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Помилка кодування JSON")
	}

	c.Set("Content-Type", "application/json")
	return c.SendString(string(data) + "\n")
}

func getNote(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Неправильний формат ID")
	}

	mu.Lock()
	note, exists := notes[id]
	mu.Unlock()

	if !exists {
		return sendError(c, fiber.StatusNotFound, "Нотатку не знайдено")
	}

	data, err := easyjson.Marshal(&note)
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Помилка кодування JSON")
	}

	c.Set("Content-Type", "application/json")
	return c.SendString(string(data) + "\n")
}

func createNote(c *fiber.Ctx) error {
	var note models.Note

	if err := easyjson.Unmarshal(c.Body(), &note); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Неправильний формат JSON")
	}

	if err := validate.Struct(&note); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Помилка валідації: title має бути >= 3 символів, content є обов'язковим")
	}

	mu.Lock()
	note.ID = idCounter
	notes[idCounter] = note
	idCounter++
	mu.Unlock()

	data, _ := easyjson.Marshal(&note)
	c.Set("Content-Type", "application/json")
	return c.Status(fiber.StatusCreated).SendString(string(data) + "\n")
}

func updateNote(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Неправильний формат ID")
	}

	var updatedNote models.Note
	if err := easyjson.Unmarshal(c.Body(), &updatedNote); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Неправильний формат JSON")
	}

	if err := validate.Struct(&updatedNote); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Помилка валідації: title має бути >= 3 символів, content є обов'язковим")
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := notes[id]; !exists {
		return sendError(c, fiber.StatusNotFound, "Нотатку не знайдено")
	}

	updatedNote.ID = id
	notes[id] = updatedNote

	data, _ := easyjson.Marshal(&updatedNote)
	c.Set("Content-Type", "application/json")
	return c.SendString(string(data) + "\n")
}

func deleteNote(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Неправильний формат ID")
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := notes[id]; !exists {
		return sendError(c, fiber.StatusNotFound, "Нотатку не знайдено")
	}

	delete(notes, id)
	return c.SendStatus(fiber.StatusNoContent)
}