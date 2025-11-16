package handlers

import (
	"bannerService/internals/dto"
	"bannerService/internals/models"
	"bannerService/internals/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

var (
	TextErrIsRequired = "text is required"
)

func (h *Handler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := h.service.SrvQuestion.Questions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(questions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var questionRequest dto.QuestionCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&questionRequest); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if questionRequest.Text == "" {
		http.Error(w, TextErrIsRequired, http.StatusBadRequest)
		return
	}

	question := models.Question{Text: questionRequest.Text}

	questionCreated, err := h.service.SrvQuestion.CreateQuestion(r.Context(), &question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(questionCreated); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) GetQuestionWithAnswers(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid question id", http.StatusBadRequest)
		return
	}

	questionWithAnswers, err := h.service.SrvQuestion.QuestionWithAnswers(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(questionWithAnswers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid question id", http.StatusBadRequest)
		return
	}

	if err := h.service.SrvQuestion.DeleteQuestion(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
