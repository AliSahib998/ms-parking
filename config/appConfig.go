package config

import (
	"github.com/alexflint/go-arg"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

const RootPath = "/v1/post"
const PrivateRootPath = "/v1/private/post"

type args struct {
	LogLevel            string   `arg:"env:LOG_LEVEL"`
	Port                int      `arg:"env:PORT"`
	RedisURL            []string `arg:"env:REDIS_SERVER_URLS"`
	DBB2BPostURL        string   `arg:"env:DB_B2B_POST_URL"`
	DBB2BPostUser       string   `arg:"env:DB_B2B_POST_USER"`
	DBB2BPostPass       string   `arg:"env:DB_B2B_POST_PASS"`
	DBB2BPostName       string   `arg:"env:DB_B2B_POST_NAME"`
	DBPoolSize          int      `arg:"env:DB_POOL_SIZE"`
	DBConnectionTimeout string   `arg:"env:DB_CONNECTION_TIMEOUT"`
	CustomerEndpoint    string   `arg:"env:MS_CUSTOMER_ENDPOINT,required"`
	FileStorageEndpoint string   `arg:"env:MS_FILE_STORAGE_ENDPOINT,required"`
	MaxFileUploadSize   int      `arg:"env:MAX_FILE_SIZE,required"`
}

// Props is global variable for environment variables usage
var Props args
var opts struct {
	Profile string `arg:"-p" default:"default"`
}

// LoadConfig loads configuration when project run
func LoadConfig() {
	arg.MustParse(&opts)
	initLogger()
	log.Info("Application is starting with profile: ", opts.Profile)

	initEnvVars()
	_ = arg.Parse(&Props)
	applyLoggerLevel()
}

func initEnvVars() {
	if godotenv.Load("config/profiles/default.env") != nil {
		log.Fatal("Error in loading environment variables from: profiles/default.env")
	} else {
		log.Info("Environment variables loaded from: profiles/default.env")
	}
	//http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func initLogger() {
	log.SetLevel(log.InfoLevel)
	if opts.Profile == "default" {
		log.SetFormatter(&log.JSONFormatter{})
	}
}

func applyLoggerLevel() {
	loglevel, err := log.ParseLevel(Props.LogLevel)
	if err != nil {
		loglevel = log.InfoLevel
	}
	log.SetLevel(loglevel)
}
