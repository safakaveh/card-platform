package uploadcsv

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

const maxUploadSize = int64(2 << 30)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	multipartReader, err := r.MultipartReader()
	if err != nil {
		writeError(w, http.StatusBadRequest, "درخواست باید از نوع multipart/form-data باشد")
		return
	}

	var orderName string
	var fileName string
	var filePart *multipart.Part

	for {
		part, nextErr := multipartReader.NextPart()
		if errors.Is(nextErr, io.EOF) {
			break
		}
		if nextErr != nil {
			writeError(w, http.StatusBadRequest, "خواندن اطلاعات آپلود ناموفق بود")
			return
		}

		switch part.FormName() {
		case "order_name":
			value, readErr := io.ReadAll(io.LimitReader(part, 4*1024))
			_ = part.Close()
			if readErr != nil {
				writeError(w, http.StatusBadRequest, "خواندن نام سفارش ناموفق بود")
				return
			}
			orderName = string(value)
		case "file":
			filePart = part
			fileName = part.FileName()
			goto importFile
		default:
			_ = part.Close()
		}
	}

importFile:
	if filePart == nil {
		writeError(w, http.StatusBadRequest, ErrMissingFile.Error())
		return
	}
	defer filePart.Close()

	if !strings.HasSuffix(strings.ToLower(fileName), ".csv") {
		writeError(w, http.StatusBadRequest, "فقط فایل با پسوند CSV قابل قبول است")
		return
	}

	result, err := h.service.Import(r.Context(), orderName, fileName, filePart)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCSV),
			errors.Is(err, ErrNoMappedColumns),
			errors.Is(err, ErrDuplicateColumn),
			errors.Is(err, ErrEmptyOrderName):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrOrderNameConflict):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "ذخیره فایل ناموفق بود")
		}
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	imports, err := h.service.List(r.Context(), limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "دریافت فایل‌ها ناموفق بود")
		return
	}
	writeJSON(w, http.StatusOK, imports)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	details, err := h.service.Get(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "فایل موردنظر پیدا نشد")
			return
		}
		writeError(w, http.StatusInternalServerError, "دریافت اطلاعات فایل ناموفق بود")
		return
	}
	writeJSON(w, http.StatusOK, details)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(r.Context(), chi.URLParam(r, "id")); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "فایل موردنظر پیدا نشد")
			return
		}
		writeError(w, http.StatusInternalServerError, "حذف فایل ناموفق بود")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, statusCode int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]string{
		"error": fmt.Sprintf("%s", message),
	})
}
