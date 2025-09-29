package config

import (
	"github.com/spf13/viper"
)

// Config menampung semua konfigurasi aplikasi.
// Field-field di struct ini harus cocok dengan key di config.yaml.
type Config struct {
	Gateway    GatewayConfig    `mapstructure:"gateway"`
	RuleEngine RuleEngineConfig `mapstructure:"rule_engine"`
}

type GatewayConfig struct {
	TCPPort           string `mapstructure:"tcp_port"`
	RuleEngineAddress string `mapstructure:"rule_engine_address"`
}

type RuleEngineConfig struct {
	GRPCPort string `mapstructure:"grpc_port"`
}

// LoadConfig membaca konfigurasi dari file atau environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)     // Path ke folder config
	viper.SetConfigName("config") // Nama file tanpa ekstensi
	viper.SetConfigType("yaml")

	viper.AutomaticEnv() // Baca juga dari environment variables jika ada

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
