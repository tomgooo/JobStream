package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tomgooo/JobStream/internal/common/db"
	"github.com/tomgooo/JobStream/internal/httpapi"
)

func main() {
	// 初始化数据库
	if err := db.Init(); err != nil {
		log.Fatalf("init db failed: %v", err)
	}

	mux := http.NewServeMux()

	// 健康检查
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		now := time.Now().Format(time.RFC3339)
		resp := fmt.Sprintf(`{"status":"ok","time":"%s"}`, now)
		_, _ = w.Write([]byte(resp))
	})

	// 可以保留之前的 debug/db 接口，也可以先注释掉
	// mux.HandleFunc("/debug/db", ...)

	// 用户相关接口
	mux.HandleFunc("/api/v1/register", httpapi.HandleRegister)
	mux.HandleFunc("/api/v1/login", httpapi.HandleLogin)

	// 职位相关接口
	mux.HandleFunc("/api/v1/jobs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			httpapi.HandleListJobs(w, r)
		case http.MethodPost:
			httpapi.HandleCreateJob(w, r)
		default:
			httpapi.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	addr := ":8080"
	log.Printf("JobStream service listening on %s ...", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
