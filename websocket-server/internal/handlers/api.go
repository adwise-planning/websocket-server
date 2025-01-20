package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	"websocket-server/internal/services"
	"websocket-server/internal/models"
	"websocket-server/internal/logging"
	"encoding/json"
)

type APIHandler struct {
	UserService    services.UserService
	MessageService services.MessageService
	EventService   services.EventService
}

func NewAPIHandler(userService services.UserService, messageService services.MessageService, eventService services.EventService) *APIHandler {
	return &APIHandler{
		UserService:    userService,
		MessageService: messageService,
		EventService:   eventService,
	}
}

func (h *APIHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/users", h.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/api/users/{id}", h.GetUser).Methods(http.MethodGet)
	r.HandleFunc("/api/messages", h.SendMessage).Methods(http.MethodPost)
	r.HandleFunc("/api/events", h.LogEvent).Methods(http.MethodPost)
}

func (h *APIHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.UserService.CreateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *APIHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	user, err := h.UserService.GetUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *APIHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var message models.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.MessageService.SendMessage(&message); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func (h *APIHandler) LogEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.EventService.LogEvent(&event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}