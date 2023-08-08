package webhooks

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/tjandrayana/ee/config"
	"github.com/tjandrayana/ee/grace"
)

type Config struct {
	HTTP HTTP `yaml:"http"`
}

// HTTP defines server config for http server
type HTTP struct {
	Port         string        `yaml:"port"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
}

func Run() {

	// Get Flag parameter from user
	var configFile string
	flag.StringVar(&configFile, "config.file", "etc/ee/config.yaml", "Path of config location")

	// Parse input flag
	flag.Parse()

	// Initialize Config
	var conf Config
	err := config.Read(&conf, configFile)
	if err != nil {
		log.Fatalf("Failed to read config in %s. Error: %s", configFile, err)
	}

	fmt.Printf("%v", conf)

	httpRouter := chi.NewRouter()

	checker := systemCheck{}
	httpRouter.Get("/ping", checker.ping)
	httpRouter.Get("/health", checker.health)

	srv := http.Server{
		ReadTimeout:  conf.HTTP.ReadTimeout * time.Second,
		WriteTimeout: conf.HTTP.WriteTimeout * time.Second,
		Handler:      httpRouter,
	}
	err = grace.ServeHTTP(&srv, conf.HTTP.Port)
	if err != nil {
		log.Fatal("Failed to start HTTP Server gracefully. Error: ", err)
	}

}

//-----------[ Pinger ]-----------------

type Tester interface {
	Ping() error
}

type systemCheck struct {
	pinger map[string]Tester
}

func (sys *systemCheck) ping(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("pong"))
}

func (sys *systemCheck) health(w http.ResponseWriter, r *http.Request) {
	var str string
	for k, v := range sys.pinger {
		start := time.Now()
		status := "Success"
		message := "successful"
		if err := v.Ping(); err != nil {
			status = "Error"
			message = err.Error()
		}
		duration := time.Since(start).Nanoseconds()
		str = fmt.Sprintf("%s%s | %s | %s | %dms\n", str, k, status, message, duration)
	}
	_, _ = w.Write([]byte(str))
}
