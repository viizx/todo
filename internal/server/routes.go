package server

import (
	"context"
	"encoding/json"
	"evolutio_to-do/internal/database"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func sendErrorResponse(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(CORS)
	r.Use(JSONResponseMiddleware)
	r.Get("/api/hello", s.HelloWorldHandler)

	r.Get("/api/health", s.healthHandler)

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "index.html")
	})
	r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "openapi.yaml")
	})

	r.HandleFunc("/api/todos", s.toDosHandler)
	r.HandleFunc("/api/todos/{id}", s.toDoHandler)
	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if health, err := s.db.Health(); err != nil {
		sendErrorResponse(w, "Database health check failed", http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(health)
	}
}

func (s *Server) toDoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		s.getTodoHandler(w, r)
		return
	} else if r.Method == "PATCH" {
		s.updateTodoHandler(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (s *Server) toDosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.getAllToDos(w, r)
	case "POST":
		s.createToDoHandler(w, r)
	default:
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) getAllToDos(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	if sort != "ASC" && sort != "DESC" {
		sort = "DESC"
	}

	todos, err := s.db.GetAllToDos(sort)
	if err != nil {
		sendErrorResponse(w, "Failed to retrieve todos", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(todos)
}

func (s *Server) getTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	todo, err := s.db.GetToDo(id)
	if err != nil {
		sendErrorResponse(w, "Todo not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(todo)
}

func (s *Server) updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req database.ToDoRequest

	todo, err := s.db.GetToDo(id)
	if err != nil {
		sendErrorResponse(w, "Todo not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := s.db.UpdateToDo(id, req); err != nil {
		sendErrorResponse(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	if req.Done != nil && *req.Done && todo.Done != *req.Done {
		smsReceiver := os.Getenv("SMS_RECIEVER")
		message := fmt.Sprintf("ToDo [%s] completed", todo.ID)
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if err := s.sms.SendSMS(ctx, smsReceiver, message); err != nil {
			sendErrorResponse(w, "Failed to send SMS notification", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) createToDoHandler(w http.ResponseWriter, r *http.Request) {
	var req database.ToDoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := s.db.CreateToDo(req); err != nil {
		sendErrorResponse(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
