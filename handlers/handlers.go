package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"
)

func GetList(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("call method", r.Method),
			slog.String("request ID", middleware.GetReqID(r.Context())),
			slog.String("request IP", r.RemoteAddr),
		)

		exec.Command("git", "fetch")
		out, err := exec.Command("git", "branch").Output()

		if err != nil {
			log.Error(err.Error())
		}

		list := strings.Fields(string(out))
		branches := make([]string, 0, len(list)-1)

		for _, branch := range list {
			if branch != "*" {
				branches = append(branches, branch)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		err = encoder.Encode(branches)

		if err != nil {
			return
		}

		log.Info("request completed")
	}
}

func SetBranch(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("call method", r.Method),
			slog.String("request ID", middleware.GetReqID(r.Context())),
			slog.String("request IP", r.RemoteAddr),
		)

		var req SetBranchRequest

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to parse request body")
			return
		}

		out, err := exec.Command("git", "checkout", req.Branch).Output()
		if err != nil {
			log.Error(err.Error())
		} else {
			log.Info("cool", out)
		}
		//exec.Command("git", "pull")
		//exec.Command("make", "dev-reset")

		log.Info("request completed", slog.String("branch", req.Branch))
	}
}

type SetBranchRequest struct {
	Branch string `json:"branch"`
}
