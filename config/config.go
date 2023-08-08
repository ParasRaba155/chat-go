// config for handling app config
package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Redis      Redis
	Database   Database
	Server     Server
	Session    Session
	LogFileLoc string
}

// Server struct contains all the configuration related to Server in application
type Server struct {
	JwtSecretKey     string
	Port             int
	ReadTimeOut      int
	WriteTimeOut     int
	JWTExpireMinutes int
}

// Database struct contains all the configuration related to Database  in application
type Database struct {
	Host        string
	User        string
	Password    string
	DefaultDB   string
	SSLMODE     string
	Port        int
	MaxOpenConn int
	TimeoutSecs int
	LogQueries  bool
}

type Redis struct {
	Address     string
	Password    string
	Username    string
	MaxRetries  int
	PoolSize    int
	DB          int
	PoolTimeout time.Duration
}

type Session struct {
	Name       string
	CookieName string
	MaxAgeSec  int
}

// mustSetEnv will read the config file and set the env variable
// in the app and it will panic on any error
func mustSetEnv(filename string) error {
	configFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer configFile.Close()

	sc := bufio.NewScanner(configFile)
	for sc.Scan() {
		line := sc.Text()
		// ignore empty lines and commented lines
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2) // split the line into 2 parts

		if len(parts) != 2 {
			return fmt.Errorf("could not split the %s line", line)
		}

		key, value := parts[0], parts[1]
		os.Setenv(key, value)
	}

	return nil
}

// LoadConfig takes the filename of env file and returns the
// or it panics on any error
func LoadConfig(filename string) (Config, error) {
	if err := mustSetEnv(filename); err != nil {
		return Config{}, err
	}

	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		return Config{}, fmt.Errorf("could no get SERVER_PORT :%w", err)
	}

	serverReadTimeOut, err := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	if err != nil {
		return Config{}, fmt.Errorf("could no get SERVER_READ_TIMEOUT :%w", err)
	}

	serverWriteTimeOut, err := strconv.Atoi(os.Getenv("SERVER_WRITE_TIMEOUT"))
	if err != nil {
		return Config{}, fmt.Errorf("could no get SERVER_WRITE_TIMEOUT :%w", err)
	}
	serverJwtExpMins, err := strconv.Atoi(os.Getenv("SERVER_JWT_ETMIN"))
	if err != nil {
		return Config{}, fmt.Errorf("could no get SERVER_JWT_ETMIN :%w", err)
	}

	dbPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		return Config{}, fmt.Errorf("could no get POSTGRES_PORT :%w", err)
	}

	dbMaxOpenConn, err := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNECTIONS"))
	if err != nil {
		return Config{}, fmt.Errorf("could no get DB_MAX_OPEN_CONNECTIONS :%w", err)
	}

	dbTimeOutSec, err := strconv.Atoi(os.Getenv("DB_TIMEOUT_SECONDS"))
	if err != nil {
		return Config{}, fmt.Errorf("could no get DB_TIMEOUT_SECONDS :%w", err)
	}

	dbLogQueries, err := strconv.ParseBool(os.Getenv("DB_LOG_QUERIES"))
	if err != nil {
		return Config{}, fmt.Errorf("could no get DB_LOG_QUERIES :%w", err)
	}

	sessionMaxAge, err := strconv.Atoi(os.Getenv("SESSION_COOKIE_MAXAGE"))
	if err != nil {
		sessionMaxAge = 15
	}

	redisMaxRetries, err := strconv.Atoi(os.Getenv("REDIS_MAXRETRIES"))
	if err != nil {
		return Config{}, fmt.Errorf("could not get REDIS_MAXRETRIES :%w", err)
	}

	redisPoolSize, err := strconv.Atoi(os.Getenv("REDIS_POOLSIZE"))
	if err != nil {
		return Config{}, fmt.Errorf("could not get REDIS_POOLSIZE:%w", err)
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDISDB"))
	if err != nil {
		return Config{}, fmt.Errorf("could not get REDISDB :%w", err)
	}

	redisPoolTimeout, err := strconv.Atoi(os.Getenv("REDIS_POOLTIMEOUT"))
	if err != nil {
		return Config{}, fmt.Errorf("could not get REDIS_POOLTIMEOUT :%w", err)
	}

	return Config{
		Redis: Redis{
			Address:     os.Getenv("REDIS_ADDRESS"),
			MaxRetries:  redisMaxRetries,
			PoolSize:    redisPoolSize,
			Password:    os.Getenv("REDIS_PASSWORD"),
			DB:          redisDB,
			Username:    os.Getenv("REDIS_USERNAME"),
			PoolTimeout: time.Second * time.Duration(redisPoolTimeout),
		},
		Server: Server{
			JwtSecretKey:     os.Getenv("SERVER_JWT_SECRETKEY"),
			Port:             serverPort,
			ReadTimeOut:      serverReadTimeOut,
			WriteTimeOut:     serverWriteTimeOut,
			JWTExpireMinutes: serverJwtExpMins,
		},
		Database: Database{
			Host:        os.Getenv("POSTGRES_HOST"),
			User:        os.Getenv("POSTGRES_USER"),
			Password:    os.Getenv("POSTGRES_PASSWORD"),
			DefaultDB:   os.Getenv("POSTGRES_DB"),
			Port:        dbPort,
			MaxOpenConn: dbMaxOpenConn,
			TimeoutSecs: dbTimeOutSec,
			SSLMODE:     os.Getenv("POSTGRES_SSLMODE"),
			LogQueries:  dbLogQueries,
		},
		Session: Session{
			Name:       os.Getenv("SESSION_NAME"),
			MaxAgeSec:  sessionMaxAge,
			CookieName: os.Getenv("SESSION_COOKIE_NAME"),
		},
		LogFileLoc: os.Getenv("LOG_DIR"),
	}, nil
}
