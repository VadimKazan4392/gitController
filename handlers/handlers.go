package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"
)

func GetList(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(slog.String("request IP", r.RemoteAddr))

		execCommand(log, "git", "fetch", "--prune")

		out := execCommand(log, "git", "branch", "-r")
		list := strings.Fields(string(out))
		branches := make([]string, 0, len(list))
		result := make(map[string][]string)

		for _, branch := range list {
			if branch != "*" {
				branches = append(branches, branch)
			}
		}

		curStr := execCommand(log, "git", "branch")
		curList := strings.Fields(string(curStr))
		var cur string
		for i, c := range curList {
			if c == "*" {
				cur = curList[i+1]
			}
		}

		if cur != "" {
			result["currentBranch"] = []string{cur}
		} else {
			result["currentBranch"] = []string{}
		}

		result["branches"] = branches

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		err := encoder.Encode(result)

		if err != nil {
			log.Error("Error encoding result", slog.String("error", err.Error()))
			return
		}

		log.Info("request completed")
	}
}

func SetBranch(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(slog.String("request IP", r.RemoteAddr))

		branch := chi.URLParam(r, "branch")

		execCommand(log, "git", "checkout", branch)

		execCommand(log, "git", "pull")
		log.Info("success pulling branch", slog.String("branch", branch))

		execCommand(log, "make", "dev-set-branch")
		log.Info("cache clear success")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		err := encoder.Encode([]string{"ok"})

		if err != nil {
			log.Error("Error encoding result", slog.String("error", err.Error()))
			return
		}

		log.Info("request completed", slog.String("branchName", branch))
	}
}

func UpdateBranch(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(slog.String("request IP", r.RemoteAddr))

		execCommand(log, "git", "pull")
		execCommand(log, "make", "dev-set-branch")
		log.Info("cache clear success")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")

		encoder := json.NewEncoder(w)
		err := encoder.Encode([]string{"ok"})

		if err != nil {
			log.Error("Error encoding result", slog.String("error", err.Error()))
			return
		}
	}
}

func execCommand(log *slog.Logger, command string, args ...string) []byte {
	result, err := exec.Command(command, args...).Output()

	if err != nil {
		log.Error("error for: ", slog.String("command", command), slog.String("error", err.Error()))
	}

	return result
}

type SetBranchRequest struct {
	Branch string `json:"branch"`
}
