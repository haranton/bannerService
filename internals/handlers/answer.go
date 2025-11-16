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
	QuestionIdErrIsRequired = "question id is required"
)

func (h *Handler) CreateAnswer(w http.ResponseWriter, r *http.Request) {

	idQuestionStr := r.PathValue("id")
	idQuestion, err := strconv.Atoi(idQuestionStr)
	if err != nil {
		http.Error(w, "invalid answer id", http.StatusBadRequest)
		return
	}

	var answerRequest dto.AnswerCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&answerRequest); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if answerRequest.Text == "" {
		http.Error(w, TextErrIsRequired, http.StatusBadRequest)
		return
	}

	if answerRequest.UserID == "" {
		http.Error(w, UserIdErrIsRequired, http.StatusBadRequest)
		return
	}

	answer := models.Answer{
		Text:       answerRequest.Text,
		UserID:     answerRequest.UserID,
		QuestionID: idQuestion,
	}

	answerCreated, err := h.service.SrvAnswer.CreateAnswer(r.Context(), &answer)
	if err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) {
			http.Error(w, "question not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(answerCreated); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) GetAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid answer id", http.StatusBadRequest)
		return
	}

	answer, err := h.service.SrvAnswer.Answer(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrAnswerNotFound) {
			http.Error(w, "question not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(answer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		if errors.Is(err, service.ErrAnswerNotFound) {
			http.Error(w, "question not found", http.StatusNotFound)
			return
		}
		http.Error(w, "invalid answer id", http.StatusBadRequest)
		return
	}

	if err := h.service.SrvAnswer.DeleteAnswer(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
