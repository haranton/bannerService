package handlers

import (
	"bannerService/internals/dto"
	"bannerService/internals/mapper"
	"bannerService/internals/models"
	"bannerService/internals/storage"
	"encoding/json"
	"errors"
	"fmt"
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

	tagIdHeader := r.URL.Query().Get("tag_id")
	if tagIdHeader == "" {
		writeJSONError(w, http.StatusBadRequest, TextErrTagIsRequired)
	}

	tagId, err := strconv.Atoi(tagIdHeader)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, TextErrTagIdValidate)
	}

	featureIdHeader := r.URL.Query().Get("feature_id")
	if featureIdHeader == "" {
		writeJSONError(w, http.StatusBadRequest, TextErrFeatureIdIsRequired)
	}

	featureId, err := strconv.Atoi(featureIdHeader)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, TextErrFeatureIdValidate)
	}

	var UseLastRevision bool
	UseLastRevisionHeader := r.URL.Query().Get("use_last_revision")
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
			writeJSONError(w, http.StatusNotFound, err.Error())
			return
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

func (h *Handler) Banners(w http.ResponseWriter, r *http.Request) {

	var err error

	tagIdHeader := r.URL.Query().Get("tag_id")
	var tagId int
	if tagIdHeader != "" {
		tagId, err = strconv.Atoi(tagIdHeader)
		if err != nil {
			tagId = 0
		}
	}

	featureIdHeader := r.URL.Query().Get("feature_id")
	var featureId int
	if featureIdHeader != "" {
		featureId, err = strconv.Atoi(featureIdHeader)
		if err != nil {
			featureId = 0
		}
	}

	limitHeader := r.URL.Query().Get("limit")
	var limit int
	if limitHeader == "" {
		limit, err = strconv.Atoi(limitHeader)
		if err != nil {
			limit = 0
		}
	}

	limitOffsetHeader := r.URL.Query().Get("offset")
	var offset int
	if limitOffsetHeader == "" {
		offset, err = strconv.Atoi(limitOffsetHeader)
		if err != nil {
			offset = 0
		}
	}

	params := dto.BannersQuery{
		Feature_id: featureId,
		Tag_id:     tagId,
		Offset:     offset,
		Limit:      limit,
	}

	banners, err := h.service.SrvBanner.Banners(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(banners); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) CreateBanner(w http.ResponseWriter, r *http.Request) {

	var bannerRequest dto.BannerCreateUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&bannerRequest); err != nil {
		errMessage := fmt.Sprintf("failed decode body request, err: %s", err.Error())
		writeJSONError(w, http.StatusBadRequest, errMessage)
		return
	}

	if string(bannerRequest.Content) == "" {
		http.Error(w, TextErrIsRequired, http.StatusBadRequest)
		return
	}

	if len(bannerRequest.TagIds) == 0 {
		writeJSONError(w, http.StatusBadRequest, TextErrTagIsRequired)
		return
	}

	if bannerRequest.FeatureId == 0 {
		errMessage := fmt.Sprintf("FeatureId is reqired")
		writeJSONError(w, http.StatusBadRequest, errMessage)
		return
	}

	banner := models.Banner{
		Content:  bannerRequest.Content,
		IsActive: bannerRequest.IsActive,
	}

	featureTags := dto.FeatureTags{
		TagIds:    bannerRequest.TagIds,
		FeatureId: bannerRequest.FeatureId,
	}

	bannerCreated, err := h.service.SrvBanner.CreateBanner(r.Context(), &banner, &featureTags)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicateFeatureTag) {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(bannerCreated); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) PatchBanner(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid question id", http.StatusBadRequest)
		return
	}

	var bannerRequest dto.BannerCreateUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&bannerRequest); err != nil {
		errMessage := fmt.Sprintf("failed decode body request, err: %s", err.Error())
		writeJSONError(w, http.StatusBadRequest, errMessage)
		return
	}

	if string(bannerRequest.Content) == "" {
		http.Error(w, TextErrIsRequired, http.StatusBadRequest)
		return
	}

	if len(bannerRequest.TagIds) == 0 {
		errMessage := fmt.Sprintf("tags is reqired")
		writeJSONError(w, http.StatusBadRequest, errMessage)
		return
	}

	if bannerRequest.FeatureId == 0 {
		errMessage := fmt.Sprintf("FeatureId is reqired")
		writeJSONError(w, http.StatusBadRequest, errMessage)
		return
	}

	banner := models.Banner{
		ID:       id,
		Content:  bannerRequest.Content,
		IsActive: bannerRequest.IsActive,
	}

	featureTags := dto.FeatureTags{
		TagIds:    bannerRequest.TagIds,
		FeatureId: bannerRequest.FeatureId,
	}

	featureTagsModel := mapper.FeatureTagsBanner(&featureTags, id)

	err = h.service.SrvBanner.UpdateBanner(r.Context(), &banner, featureTagsModel)
	if err != nil {
		if errors.Is(err, storage.ErrBannerNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errMsg := fmt.Sprintf("invalid question id err: %s", err.Error())
		writeJSONError(w, http.StatusBadRequest, errMsg)
		return
	}

	if err := h.service.SrvBanner.DeleteBanner(r.Context(), id); err != nil {
		if errors.Is(err, storage.ErrBannerNotFound) {
			writeJSONError(w, http.StatusNotFound, "")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
