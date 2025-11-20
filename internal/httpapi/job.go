package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/tomgooo/JobStream/internal/common/db"
	"github.com/tomgooo/JobStream/internal/model"
)

// 创建职位请求体
type CreateJobRequest struct {
	Title     string `json:"title"`
	Company   string `json:"company"`
	City      string `json:"city"`
	SalaryMin int    `json:"salary_min"`
	SalaryMax int    `json:"salary_max"`
	TechStack string `json:"tech_stack"`
	Source    string `json:"source"`
	URL       string `json:"url"`
}

// HandleCreateJob 处理创建职位：POST /api/v1/jobs
func HandleCreateJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Company = strings.TrimSpace(req.Company)
	req.City = strings.TrimSpace(req.City)

	if req.Title == "" || req.Company == "" {
		WriteError(w, http.StatusBadRequest, "title and company are required")
		return
	}

	job := &model.Job{
		Title:     req.Title,
		Company:   req.Company,
		City:      req.City,
		SalaryMin: req.SalaryMin,
		SalaryMax: req.SalaryMax,
		TechStack: req.TechStack,
		Source:    req.Source,
		URL:       req.URL,
	}

	if err := db.DB.Create(job).Error; err != nil {
		WriteError(w, http.StatusInternalServerError, "failed to create job")
		return
	}

	WriteJSON(w, http.StatusCreated, baseResponse{
		Code:    0,
		Message: "job created",
		Data:    job,
	})
}

// HandleListJobs 列出职位：GET /api/v1/jobs?city=xxx&keyword=go&limit=20
func HandleListJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	city := strings.TrimSpace(r.URL.Query().Get("city"))
	keyword := strings.TrimSpace(r.URL.Query().Get("keyword"))
	limitStr := r.URL.Query().Get("limit")

	limit := 20
	if limitStr != "" {
		if n, err := strconv.Atoi(limitStr); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}

	var jobs []model.Job
	q := db.DB.Model(&model.Job{})

	if city != "" {
		q = q.Where("city = ?", city)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("title LIKE ? OR company LIKE ? OR tech_stack LIKE ?", like, like, like)
	}

	if err := q.Order("id DESC").Limit(limit).Find(&jobs).Error; err != nil {
		WriteError(w, http.StatusInternalServerError, "db error")
		return
	}

	WriteJSON(w, http.StatusOK, baseResponse{
		Code:    0,
		Message: "ok",
		Data:    jobs,
	})
}
