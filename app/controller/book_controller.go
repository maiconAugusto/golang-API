package controller

import (
	"encoding/json"
	"io"
	"service/app/model"
	"service/database"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"
)

func CreateBook(response http.ResponseWriter, request *http.Request) {
	body, _ := io.ReadAll(request.Body)

	book := model.Book{}

	if err := json.Unmarshal(body, &book); err != nil {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "error, parser body.",
		})
		return
	}

	db, err := database.DatabaseConnection()

	if err != nil {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "error open database.",
		})
		return
	}

	defer db.Close()

	statement, err := db.Prepare("insert into book (name, author) values (?, ?)")

	if err != nil {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "error prepare database.",
		})
		return
	}

	res, err := statement.Exec(book.Name, book.Author)

	if err != nil {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "error insert book",
		})
		return
	}

	statement.Close()

	id, err := res.LastInsertId()

	if err != nil {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "error get insert id book",
		})
		return
	}

	book.Id = id

	response.Header().Set("Content=Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(book)
}

func GetBookById(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id := params["id"]

	idParsed, err := strconv.ParseUint(id, 10, 32)

	if err != nil {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "error in parser uuid",
		})
		return
	}

	db, err := database.DatabaseConnection()

	if err != nil {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "open database error",
		})
		return
	}

	defer db.Close()

	res, err := db.Query("SELECT * FROM book WHERE id = ?", idParsed)

	if err != nil {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "show book by id erro.",
		})
		return
	}

	book := model.Book{}

	if res.Next() {
		if err := res.Scan(&book.Name, &book.Author, &book.Id); err != nil {
			response.Header().Set("Content=Type", "application/json")
			response.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(response).Encode(map[string]string{
				"error": "scan book error.",
			})
		}
	}

	if book.Id == 0 {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusAccepted)
		json.NewEncoder(response).Encode(nil)
		return
	}

	response.Header().Set("Content=Type", "application/json")
	response.WriteHeader(http.StatusAccepted)
	json.NewEncoder(response).Encode(book)
}

func UpdateBookById(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id := params["id"]

	idParsed, err := strconv.ParseUint(id, 10, 32)

	if err != nil {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "parser id error",
		})
		return
	}

	book := model.Book{}
	book.Id = int64(idParsed)

	body, err := io.ReadAll(request.Body)

	if err != nil {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "body empty",
		})
		return
	}

	if err := json.Unmarshal(body, &book); err != nil {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "body parser error",
		})
		return
	}

	db, err := database.DatabaseConnection()

	if err != nil {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "database error.",
		})
		return
	}

	defer db.Close()

	statement, err := db.Prepare("UPDATE book SET name = ?, author = ? WHERE id = ?")

	if err != nil {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "prepare book error.",
		})
		return
	}

	if _, err := statement.Exec(book.Name, book.Author, idParsed); err != nil {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "update book error.",
		})
		return
	}

	statement.Close()

	response.Header().Set("Content=Type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(book)
}

func DeleteBookById(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id := params["id"]

	idParsed, err := strconv.ParseUint(id, 10, 32)

	if err != nil {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "parser id error",
		})
		return
	}

	db, err := database.DatabaseConnection()

	if err != nil {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "error open database",
		})
		return
	}

	defer db.Close()

	if err != nil {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "error parser id",
		})
		return
	}

	statement, err := db.Prepare("DELETE  FROM book WHERE id = ?")

	if err != nil {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "prepare delete error.",
		})
		return
	}

	if _, err := statement.Exec(idParsed); err != nil {
		response.Header().Set("Content=Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(map[string]string{
			"error": "delete book error.",
		})
		return
	}
	statement.Close()

	response.Header().Set("Content=Type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(nil)

}
