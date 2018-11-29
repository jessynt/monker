package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type GlobalConfig struct {
	LogLevel   string `envconfig:"LOG_LEVEL"`
	Kafka      KafkaConfig
	Worker     WorkerConfig
	Nsq        NsqConfig
	StorageDSN string `envconfig:"STORAGE_DSN"`
}

type KafkaConfig struct {
	Brokers []string `envconfig:"KAFKA_BROKERS"`
	Topic   string   `envconfig:"KAFKA_TOPIC"`
}

type NsqConfig struct {
	Lookupds    []string `envconfig:"NSQ_LOOKUPDS"`
	Topic       string   `envconfig:"NSQ_TOPIC"`
	Channel     string   `envconfig:"NSQ_CHANNEL"`
	Concurrency int      `envconfig:"NSQ_CONCURRENCY"`
}

type WorkerConfig struct {
	//  worker 每个循环休眠的时间，以避免 CPU 过载
	CycleTimeout time.Duration `envconfig:"WORKER_CYCLE_TIMEOUT"`
	// cache 里面最多缓存多少消息，超出所有消息将立即被推送到 Kafka
	CacheSize int `envconfig:"WORKER_CACHE_SIZE"`
	// 消息最多在 cache 里面保存多久，超时所有消息将立即被推送到 Kafka
	CacheFlushTimeout time.Duration `envconfig:"WORKER_CACHE_FLUSH_TIMEOUT"`
	// 从仓库读取消息的间隔时间，必须比 CycleTimeout 大两倍以上
	StorageReadTimeout time.Duration `envconfig:"WORKER_STORAGE_READ_TIMEOUT"`
}

func init() {
	viper.SetDefault("logLevel", "info")
	viper.SetDefault("worker.cycleTimeout", time.Duration(2)*time.Second)
	viper.SetDefault("worker.cacheSize", 10)
	viper.SetDefault("worker.cacheFlushTimeout", time.Duration(5)*time.Second)
	viper.SetDefault("worker.storageReadTimeout", time.Duration(10)*time.Second)
	viper.SetDefault("nsq.concurrency", 1)
}

func Load(path string) (*GlobalConfig, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Warn("No config file found, loading config from environment variables")
		return LoadConfigFromEnv()
	}
	log.WithField("path", viper.ConfigFileUsed()).Info("Config loaded from file")

	var instance GlobalConfig
	if err := viper.Unmarshal(&instance); err != nil {
		return nil, err
	}

	return &instance, nil
}

func LoadConfigFromEnv() (*GlobalConfig, error) {
	var instance GlobalConfig

	if err := viper.Unmarshal(&instance); err != nil {
		return nil, err
	}

	err := envconfig.Process("", &instance)
	if err != nil {
		return &instance, err
	}

	return &instance, nil
}
