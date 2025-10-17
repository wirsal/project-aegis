package config

import (
	"github.com/spf13/viper"
)

// Config menampung semua konfigurasi aplikasi.
type Config struct {
	Gateway      GatewayConfig      `mapstructure:"gateway"`
	RuleEngine   RuleEngineConfig   `mapstructure:"rule_engine"`
	Persistence  PersistenceConfig  `mapstructure:"persistence"`
	Notification NotificationConfig `mapstructure:"notification"`
	Database     DatabaseConfig     `mapstructure:"database"`
}

type GatewayConfig struct {
	TCPPort           string `mapstructure:"tcp_port"`
	RuleEngineAddress string `mapstructure:"rule_engine_address"`
}

type RuleEngineConfig struct {
	GRPCPort string `mapstructure:"grpc_port"`
}

// Struct baru untuk konfigurasi persistence
type PersistenceConfig struct {
	GRPCPort     string `mapstructure:"grpc_port"`
	GRPCPAddress string `mapstructure:"grpc_address"`
}

type NotificationConfig struct {
	GRPCPort        string `mapstructure:"grpc_port"`
	GRPCPAddress    string `mapstructure:"grpc_address"`
	FCMgatewayURL   string `mapstructure:"fcmgateway_url"`
	SMSGatewayURL   string `mapstructure:"smsgateway_url"`
	EmailGatewayURL string `mapstructure:"emailgateway_url"`
	WAGatewayURL    string `mapstructure:"wagateway_url"`
}

// Struct untuk konfigurasi database
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

// Fungsi LoadConfig tidak perlu diubah.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
