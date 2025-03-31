package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time" // time.Duration için eklendi

	"github.com/spf13/viper"
)

// DatabaseConfig veritabanı bağlantı parametrelerini tutar
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`             // YAML'daki 'host' anahtarına karşılık gelir
	Port            int           `mapstructure:"port"`             // YAML'daki 'port' anahtarına karşılık gelir
	User            string        `mapstructure:"user"`             // YAML'daki 'user' anahtarına karşılık gelir
	Password        string        `mapstructure:"password"`         // YAML'daki 'password' anahtarına karşılık gelir
	DBName          string        `mapstructure:"dbname"`           // YAML'daki 'dbname' anahtarına karşılık gelir
	SSLMode         string        `mapstructure:"sslmode"`          // YAML'daki 'sslmode' anahtarına karşılık gelir
	MaxOpenConns    int           `mapstructure:"max_open_conns"`   // YAML'daki 'max_open_conns' anahtarına karşılık gelir
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`   // YAML'daki 'max_idle_conns' anahtarına karşılık gelir
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`// YAML'daki 'conn_max_lifetime' anahtarına karşılık gelir (örn: "5m")
}

// SMTPConfig SMTP sunucu detaylarını tutar
type SMTPConfig struct {
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	User   string `mapstructure:"user"`
	Pass   string `mapstructure:"pass"`
	From   string `mapstructure:"from"`
	Secure bool   `mapstructure:"secure"`
}

// ServerConfig sunucu konfigürasyonunu tutar
type ServerConfig struct {
	Port        string `mapstructure:"port"`
	MetricsPort string `mapstructure:"metrics_port"`
}

// Config genel uygulama konfigürasyonunu tutar
type Config struct {
	SMTP     SMTPConfig     `mapstructure:"smtp"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"` // <-- YENİ EKLENEN ALAN
}

// GlobalConfig global konfigürasyon değişkeni
var GlobalConfig Config

// LoadConfig konfigürasyon dosyasını okur ve struct'a yükler
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // config.yaml veya config.json vb. aranır
	viper.SetConfigType("yaml")   // Dosya tipini belirtir

	// Config dosyasının konumunu bul
	configPath := findConfigPath()
	if configPath == "" {
		// Belirtilen yollarda config.yaml bulunamadıysa, proje kök dizinini de deneyebiliriz
        // Veya doğrudan bir hata dönebiliriz
		// Proje kök dizinini bulmak biraz daha karmaşık olabilir,
		// bu yüzden şimdilik sadece belirli yolları arıyoruz.
		// Gerekirse Viper'ın AddConfigPath fonksiyonu da kullanılabilir.
		// viper.AddConfigPath(".") // Çalışma dizini
        // viper.AddConfigPath("./config") // ./config dizini
		return nil, fmt.Errorf("config.yaml file not found in any of the specified search paths")
	}
    fmt.Printf("Using config file: %s\n", configPath) // Hangi config dosyasının kullanıldığını logla
	viper.SetConfigFile(configPath) // Bulunan config dosyasını kullan

	// Config dosyasını oku
	if err := viper.ReadInConfig(); err != nil {
		// Dosyanın bulunamaması durumunu özel olarak kontrol et
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
             return nil, fmt.Errorf("config file not found at %s: %w", configPath, err)
        }
		return nil, fmt.Errorf("error reading config file '%s': %w", configPath, err)
	}

	// Ortam değişkenlerini de oku (isteğe bağlı, önceliklendirme için kullanışlı)
	viper.AutomaticEnv()
	// Örnek: DATABASE_HOST ortam değişkeni config.database.host'u override edebilir
	// viper.SetEnvPrefix("APP") // Ortam değişkenleri için önek belirleyebilirsiniz (örn: APP_DATABASE_HOST)
	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // config.database.host -> DATABASE_HOST dönüşümü için

	// Config'i struct'a dönüştür
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return nil, fmt.Errorf("error unmarshaling config into struct: %w", err)
	}

	return &GlobalConfig, nil
}

// GetConfig yüklenmiş global konfigürasyonu döndürür
func GetConfig() *Config {
	return &GlobalConfig
}

// findConfigPath config.yaml dosyasını çeşitli konumlarda arar
func findConfigPath() string {
	searchPaths := []string{
		"./config", // Proje içindeki /config klasörü
		".",        // Proje kök dizini / çalışma dizini
		// Gerekirse diğer yaygın konumlar eklenebilir
		// "/etc/myapp/",
		// filepath.Join(os.Getenv("HOME"), ".myapp"),
	}

	for _, path := range searchPaths {
		fullPath := filepath.Join(path, "config.yaml")
		if _, err := os.Stat(fullPath); err == nil {
			// Dosya bulundu ve erişilebilir
			return fullPath
		}
	}

	// Hiçbir yerde bulunamadı
	return ""
}