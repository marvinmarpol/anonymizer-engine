package app

import (
	"errors"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBAddress        string `env:"DATA_PROTECTION_DB_ADDRESS" env-default:"localhost:5432"`
	DBUser           string `env:"DATA_PROTECTION_DB_USER" env-default:"postgres"`
	DBPassword       string `env:"DATA_PROTECTION_DB_PASSWORD" env-default:"password"`
	DBName           string `env:"DATA_PROTECTION_DB_NAME" env-default:"anonymizer-engine"`
	DBPoolSize       int    `env:"DATA_PROTECTION_DB_POOL_SIZE" env-default:"10"`
	DBMaxRetries     int    `env:"DATA_PROTECTION_DB_MAX_RETRIES" env-default:"3"`
	DBRetryDelay     int    `env:"DATA_PROTECTION_DB_RETRY_DELAY" env-default:"3"`
	DBIdleTimeout    int    `env:"DATA_PROTECTION_DB_IDLE_TIMEOUT" env-default:"30"`
	DBWriteTimeout   int    `env:"DATA_PROTECTION_DB_WRITE_TIMEOUT" env-default:"30"`
	DBPoolTimeout    int    `env:"DATA_PROTECTION_DB_POOL_TIMEOUT" env-default:"30"`
	DBReadTimeout    int    `env:"DATA_PROTECTION_DB_READ_TIMEOUT" env-default:"30"`
	ServiceAddress   string `env:"DATA_PROTECTION_SERVICE_ADDRESS" env-default:":8080"`
	NewrelicAPIKey   string `env:"DATA_PROTECTION_NEWRELIC_API_KEY" env-default:"secret"`
	NewrelicAPMName  string `env:"DATA_PROTECTION_NEWRELIC_APM_NAME" env-default:"ANONYMIZER-ENGINE[LOCAL]"`
	LogLevel         int    `env:"DATA_PROTECTION_LOG_LEVEL" env-default:"4"`
	PublicKeyPath    string `env:"DATA_PROTECTION_PUBLIC_KEY_PATH" env-default:"infrastructures/kek/public_key.pem"`
	PrivateKeyPath   string `env:"DATA_PROTECTION_PRIVATE_KEY_PATH" env-default:"infrastructures/kek/private_key.pem"`
	Realm            string `env:"DATA_PROTECTION_REALM" env-default:"local"`
	BasicUsername    string `env:"DATA_PROTECTION_BASIC_USERNAME" env-default:"admin"`
	BasicPassword    string `env:"DATA_PROTECTION_BASIC_PASSWORD" env-default:"password"`
	GoogleCredential string `env:"GOOGLE_APPLICATION_CREDENTIALS" env-default:"infrastructures/gcpserviceaccount/dev-pubsub.json"`
	GoogleProjectID  string `env:"GOOGLE_PROJECT_ID" env-default:"pubsub-development"`
	PubsubMaxRetry   int    `env:"DATA_PROTECTION_PUBSUB_MAX_RETRY" env-default:"3"`
	PubsubRetryMS    int    `env:"DATA_PROTECTION_PUBSUB_RETRY_MS" env-default:"200"`
}

func loadConfig() (Config, error) {
	var config Config
	err := cleanenv.ReadConfig(".env", &config)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return config, err
	}

	if errors.Is(err, os.ErrNotExist) {
		err := cleanenv.ReadEnv(&config)
		if err != nil {
			return config, err
		}
	}

	return config, nil
}
