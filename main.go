package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"dev.acorello.it/go/contacts/contact"
	http_contact "dev.acorello.it/go/contacts/contact/http"
	"dev.acorello.it/go/contacts/public_assets"
)

var repo = contact.NewPopulatedInMemoryContactRepository()

var CommitHash = func() string {
	if sha, found := os.LookupEnv("GITHUB_SHA"); found {
		return sha
	} else if buildInfo, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range buildInfo.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}
	return "🤷"
}()

func main() {
	mux := http.NewServeMux()
	const publicRootPath = "/public/"
	mux.HandleFunc(publicRootPath,
		LoggingHandler(http.StripPrefix(publicRootPath, public_assets.FileServer())))

	contactResourcePaths := http_contact.ResourcePaths{
		Root:  "/contact/",
		Form:  "/contact/form",
		List:  "/contact/list",
		Email: "/contact/email",
	}

	if validatedPaths, err := contactResourcePaths.Validated(); err != nil {
		log.Fatal(err)
	} else {
		contactHandler := http_contact.NewContactHandler(validatedPaths, &repo)
		mux.HandleFunc(validatedPaths.Root.String(), LoggingHandler(contactHandler))
		homeRedirect := http.RedirectHandler(validatedPaths.List.String(), http.StatusFound)
		mux.HandleFunc("/", LoggingHandler(homeRedirect))
	}

	mux.HandleFunc(healthCheckPath, healthcheck)
	var srv = http.Server{
		Addr:    bindAddress(),
		Handler: mux,
	}

	shutdownDone := make(chan struct{})
	go waitShutdownSignal(&srv, shutdownDone)

	log.Printf("Starting server at %q", srv.Addr)
	if err := srv.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
		<-shutdownDone
		log.Printf("Bye.")
	} else {
		log.Fatal(err)
	}
}

func waitShutdownSignal(srv *http.Server, done chan<- struct{}) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	signal := <-signals
	log.Printf("Received shutdown signal %q", signal)
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("Shutdown error: %v", err)
	} else {
		log.Printf("Shutdown.")
	}
	close(done)
}

func bindAddress() string {
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	return host + ":8080"
}

func LoggingHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf(`serving '%s %s'`, r.Method, r.URL)
		h.ServeHTTP(w, r)
	}
}

const healthCheckPath = "/healthcheck"

func healthcheck(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format(time.RFC1123Z)
	_, err := fmt.Fprintf(w, "Commit: %s\nTime: %s\n", CommitHash, now)
	if err != nil {
		log.Printf("error reporting %s: %v", healthCheckPath, err)
	} else {
		log.Printf("%s reported\n", healthCheckPath)
	}
}
