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
	contactHTTP "dev.acorello.it/go/contacts/contact/http"
	"dev.acorello.it/go/contacts/public_assets"
	"github.com/acorello/uttpil"
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
	mux.Handle(publicRootPath, http.StripPrefix(publicRootPath, public_assets.FileServer()))

	contactResourcePaths := contactHTTP.Paths{
		Root:  "/contact/",
		Form:  "/contact/form",
		List:  "/contact/list",
		Email: "/contact/email",
	}

	if validatedPaths, err := contactResourcePaths.Validated(); err != nil {
		log.Fatal(err)
	} else {
		contactHTTP.RegisterHandlers(mux, validatedPaths, &repo)
		homeRedirect := http.RedirectHandler(validatedPaths.List.String(), http.StatusFound)
		mux.Handle("/", homeRedirect)
	}

	mux.HandleFunc(healthCheckPath, healthcheck)
	var srv = http.Server{
		Addr:    bindAddress(),
		Handler: uttpil.LoggingHandler(mux),
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
