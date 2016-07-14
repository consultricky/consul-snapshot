package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var hostname string

// Config is a struct to hold the backup configuration
type Config struct {
	S3Bucket       string
	S3Region       string
	S3AccessKey    string
	S3SecretKey    string
	Hostname       string
	BackupInterval time.Duration
	TmpDir         string
	Acceptance     bool
	Version        string
	Encryption     string
}

// When starting, just set the hostname
func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		log.Fatalf("Unable to determine hostname: %v", err)
	}
}

// This checks a slice to see if anything is empty
func checkEmpty(s []string) bool {
	for _, item := range s {
		if len(item) < 1 {
			return false
		}
	}
	return true
}

// Set the environment variables that are required
func setEnvVars(conf *Config, tests bool) error {
	conf.S3Bucket = os.Getenv("S3BUCKET")
	conf.S3Region = os.Getenv("S3REGION")
	conf.S3AccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	conf.S3SecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	backupInterval := os.Getenv("BACKUPINTERVAL")
	conf.TmpDir = os.Getenv("SNAPSHOT_TMP_DIR")
	acceptanceTest := os.Getenv("ACCEPTANCE_TEST")
	conf.Encryption = os.Getenv("CRYPTO_PASSWORD")

	// if the environment variable isn't set, just set the dir to /tmp
	if conf.TmpDir == "" {
		conf.TmpDir = "/tmp"
	}

	// if the environment variable isn't set, require specific env vars
	if acceptanceTest == "" {
		conf.Acceptance = false
		if tests {
			log.Println("Running tests, skipping ENV var requirements")
		} else {
			envChecks := []string{conf.S3Bucket, conf.S3Region, backupInterval}
			if checkEmpty(envChecks) == false {
				log.Fatal("[ERR] Required env var missing, exiting")
			}
		}
		backupStrToInt, err := strconv.Atoi(backupInterval)
		if err != nil {
			return fmt.Errorf("Unable to convert BACKUPINTERVAL environment var to integer: %v", err)
		}
		backupTimeDuration := time.Duration(backupStrToInt) * time.Second
		conf.BackupInterval = backupTimeDuration
	} else {
		conf.Acceptance = true
		conf.BackupInterval = 60 * time.Second
	}

	return nil
}

// ParseConfig parses the config and returns it
func ParseConfig(tests bool) *Config {
	// Set some defaults
	conf := &Config{}

	err := setEnvVars(conf, tests)
	if err != nil {
		log.Fatalf("[ERR] %v", err)
	}

	conf.Hostname = hostname
	return conf
}
