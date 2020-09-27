package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-homedir"
	"github.com/pinheirolucas/shortinho/database"
	"github.com/pinheirolucas/shortinho/shortener"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

var cfgFile string

func main() {
	logLevelStr := viper.GetString("log.level")
	if logLevelStr == "" {
		fmt.Println("empty log level")
		os.Exit(1)
	}

	logLevel, err := zerolog.ParseLevel(logLevelStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	disableLogColors := viper.GetBool("log.disable-colors")

	zerolog.SetGlobalLevel(logLevel)
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			NoColor:    disableLogColors,
			TimeFormat: time.RFC3339,
		},
	)

	databaseEngine, err := database.EngineFromString(viper.GetString("database.engine"))
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	databaseURI := viper.GetString("database.uri")
	if databaseURI == "" {
		log.Fatal().Msg("empty database URI")
	}

	log.Info().
		Interface("engine", databaseEngine).
		Msg("connecting to database")
	closeDatabase, err := database.Connect(databaseEngine, databaseURI)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	defer closeDatabase()
	log.Info().
		Msg("connected to database")

	host := viper.GetString("server.host")
	if host == "" {
		log.Fatal().Msg("empty http server host")
	}

	port := viper.GetInt("server.port")
	if port == 0 {
		log.Fatal().Msg("empty http server port")
	}

	address := fmt.Sprintf("%s:%d", host, port)

	allowedOrigins := viper.GetStringSlice("server.headers.allowed-origins")
	if len(allowedOrigins) == 0 {
		log.Fatal().Msg("no allowed origins provided")
	}

	router := gin.New()

	router.Use(logger.SetLogger(logger.Config{
		Logger: &log.Logger,
		UTC:    true,
	}))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(gin.Recovery())

	shortenerGroup := router.Group("/")
	shortenerHandlers, err := shortener.NewHandlers()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	shortenerHandlers.Register(shortenerGroup)

	log.Info().Str("address", address).Msg("listening for HTTP requests")
	router.Run(address)
}

func init() {
	pflag.String(
		"host",
		"0.0.0.0",
		"host to bind shortinhos HTTP server, by defaults it binds to the entire interface at the given port",
	)
	pflag.Int("port", 3000, "port to bind shortinhos HTTP server")
	pflag.StringSlice("allowed-origins", []string{"*"}, "the origins that the HTTP server will accept requests from")
	pflag.String("database-engine", "postgres", "the database engine, it suports: postgres and mongodb")
	pflag.String("database-uri", "", "the URI where the database is hosted")
	pflag.String("log-level", "debug", "set the application log level (default is debug)")
	pflag.Bool("no-log-colors", false, "disable log colors")
	pflag.Int("slug-size", 4, "slug size for generated URL's (default and minimum is 4)")
	pflag.StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.tagsrv.yaml)")

	pflag.Parse()

	viper.BindPFlag("server.host", pflag.Lookup("host"))
	viper.BindPFlag("server.port", pflag.Lookup("port"))
	viper.BindPFlag("server.headers.allowed-origins", pflag.Lookup("allowed-origins"))
	viper.BindPFlag("log.level", pflag.Lookup("log-level"))
	viper.BindPFlag("log.disable-colors", pflag.Lookup("no-log-colors"))
	viper.BindPFlag("database.engine", pflag.Lookup("database-engine"))
	viper.BindPFlag("database.uri", pflag.Lookup("database-uri"))
	viper.BindPFlag("app.slug-size", pflag.Lookup("slug-size"))

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.SetConfigName(".shortinho")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(cwd)
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/shortinho")
	}

	viper.SetEnvPrefix("SHORTINHO")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
