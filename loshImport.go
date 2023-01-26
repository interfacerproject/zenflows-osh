package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

type Config struct {
	host string
	port int
}

func loadEnvConfig() Config {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	return Config{
		host: os.Getenv("HOST"),
		port: port,
	}
}

func cloneAndAnalyze(w http.ResponseWriter, repoURL string) {
	tmpDir, err := ioutil.TempDir("", "repo-")
	if err != nil {
		fmt.Fprintf(w, "Error creating temporary directory: %s", err)
		return
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("git", "clone", repoURL, tmpDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(w, "Error cloning repository: %s", err)
		return
	}

	oshCmd := exec.Command("./osh", "-C", tmpDir, "check", "--json")
	oshOut, err := oshCmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(w, "Error running osh tool: %s, Output: %s", err, oshOut)
		return
	}

	fmt.Fprintf(w, "Successfully cloned and analyzed repository %s. Output: %s", repoURL, oshOut)
}

func cloneHandler(w http.ResponseWriter, r *http.Request) {
	repo := r.URL.Query().Get("repo")
	if repo == "" {
		fmt.Fprintf(w, "Please provide a repository URL to clone")
		return
	}
	log.Printf("Cloning repository %s\n", repo)

	cloneAndAnalyze(w, repo)
}

func main() {
	config := loadEnvConfig()
	log.Printf("Starting service on %s:%d\n", config.host, config.port)

	mux := http.NewServeMux()
	mux.HandleFunc("/clone", cloneHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", config.host, config.port), mux))
}
