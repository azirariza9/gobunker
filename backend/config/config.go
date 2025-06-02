package config

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type TokenConfig struct {
	ApplicationName     string
	JWTSignatureKey     []byte
	JWTSigningMethod    *jwt.SigningMethodHMAC
	AccessTokenLifeTime time.Duration
}

type DBConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Driver   string
}

type APIConfig struct {
	ApiPort string
}

type Config struct {
	DBConfig
	APIConfig
	TokenConfig
}

func (c *Config) readConfig() error {
	// Set default config file
	viper.SetConfigName(".env") // Nama file .env
	viper.SetConfigType("env")  // Tipe file .env
	viper.AddConfigPath(".")    // Menambahkan path direktori tempat file .env berada

	// Membaca konfigurasi dari file
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading .env file: %v", err)
	}

	// Mengambil konfigurasi DB
	c.DBConfig = DBConfig{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		Database: viper.GetString("DB_NAME"),
		Username: viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		Driver:   viper.GetString("DB_DRIVER"),
	}

	// Mengambil konfigurasi API
	c.APIConfig = APIConfig{
		ApiPort: viper.GetString("API_PORT"),
	}

	c.TokenConfig = TokenConfig{
		ApplicationName:     viper.GetString("APPLICATION_NAME"),
		JWTSignatureKey:     []byte(viper.GetString("JWT_SECRET")),
		JWTSigningMethod:    jwt.SigningMethodHS256,
		AccessTokenLifeTime: time.Duration(1) * time.Hour, // 1 jam
	}

	// Validasi apakah konfigurasi yang diperlukan sudah ada
	if c.DBConfig.Host == "" || c.DBConfig.Port == "" || c.DBConfig.Username == "" || c.DBConfig.Password == "" || c.APIConfig.ApiPort == "" {
		return fmt.Errorf("required config missing")
	}

	return nil
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cfg.readConfig()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
