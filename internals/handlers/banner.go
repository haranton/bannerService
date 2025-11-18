package handlers

import (
	"bannerService/internals/dto"
	"bannerService/internals/models"
	"bannerService/internals/service"
	"bannerService/internals/storage"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

var (
	TextErrIsRequired          = "text is required"
	TextErrTagIsRequired       = "tag id is required"
	TextErrFeatureIdIsRequired = "feature id is required"
	TextErrFeatureIdValidate   = "feature id must be int"
	TextErrTagIdValidate       = "tag id must be int"
)

func (h *Handler) Banner(w http.ResponseWriter, r *http.Request) {

	tagIdHeader := r.Header.Get("tag_id")
	if tagIdHeader == "" {
		writeJSONError(w, http.StatusBadRequest, TextErrTagIsRequired)
	}

	tagId, err := strconv.Atoi(tagIdHeader)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, TextErrTagIdValidate)
	}

	featureIdHeader := r.Header.Get("feature_id")
	if featureIdHeader == "" {
		writeJSONError(w, http.StatusBadRequest, TextErrFeatureIdIsRequired)
	}

	featureId, err := strconv.Atoi(featureIdHeader)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, TextErrFeatureIdValidate)
	}

	var UseLastRevision bool
	UseLastRevisionHeader := r.Header.Get("use_last_revision")
	if UseLastRevisionHeader == "true" {
		UseLastRevision = true
	}

	params := dto.BannerQuery{
		UseLastRevision: UseLastRevision,
		Feature_id:      featureId,
		Tag_id:          tagId,
	}

	banner, err := h.service.SrvBanner.Banner(r.Context(), params)
	if err != nil {
		if errors.Is(err, storage.ErrBannerNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(banner); err != nil {
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
