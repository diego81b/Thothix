package vault

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/vault/api"
)

type Client struct {
	client *api.Client
	mount  string
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type ClerkConfig struct {
	SecretKey      string `json:"secret_key"`
	WebhookSecret  string `json:"webhook_secret"`
	PublishableKey string `json:"publishable_key"`
}

type AppConfig struct {
	JWTSecret     string `json:"jwt_secret"`
	EncryptionKey string `json:"encryption_key"`
	Environment   string `json:"environment"`
}

func NewVaultClient() (*Client, error) {
	config := api.DefaultConfig()
	config.Address = os.Getenv("VAULT_ADDR")
	if config.Address == "" {
		config.Address = "http://localhost:8200"
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("VAULT_TOKEN environment variable not set")
	}
	client.SetToken(token)

	mount := os.Getenv("VAULT_MOUNT")
	if mount == "" {
		mount = "thothix"
	}

	return &Client{
		client: client,
		mount:  mount,
	}, nil
}

func (v *Client) GetDatabaseConfig() (*DatabaseConfig, error) {
	secret, err := v.client.Logical().Read(fmt.Sprintf("%s/data/database", v.mount))
	if err != nil {
		return nil, fmt.Errorf("failed to read database secrets: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no database secrets found")
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid secret format")
	}

	return &DatabaseConfig{
		Host:     data["host"].(string),
		Port:     data["port"].(string),
		Username: data["username"].(string),
		Password: data["password"].(string),
		Database: data["database"].(string),
	}, nil
}

func (v *Client) GetClerkConfig() (*ClerkConfig, error) {
	secret, err := v.client.Logical().Read(fmt.Sprintf("%s/data/clerk", v.mount))
	if err != nil {
		return nil, fmt.Errorf("failed to read clerk secrets: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no clerk secrets found")
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid secret format")
	}

	return &ClerkConfig{
		SecretKey:      data["secret_key"].(string),
		WebhookSecret:  data["webhook_secret"].(string),
		PublishableKey: data["publishable_key"].(string),
	}, nil
}

func (v *Client) GetAppConfig() (*AppConfig, error) {
	secret, err := v.client.Logical().Read(fmt.Sprintf("%s/data/app", v.mount))
	if err != nil {
		return nil, fmt.Errorf("failed to read app secrets: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no app secrets found")
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid secret format")
	}

	return &AppConfig{
		JWTSecret:     data["jwt_secret"].(string),
		EncryptionKey: data["encryption_key"].(string),
		Environment:   data["environment"].(string),
	}, nil
}

// LoadConfigFromVault loads all configuration from Vault
func LoadConfigFromVault() error {
	useVault := os.Getenv("USE_VAULT")
	if useVault != "true" {
		log.Println("Vault integration disabled, using environment variables")
		return nil
	}

	log.Println("🔐 Loading configuration from Vault...")

	vaultClient, err := NewVaultClient()
	if err != nil {
		return fmt.Errorf("failed to initialize vault client: %w", err)
	}

	// Load database config
	dbConfig, err := vaultClient.GetDatabaseConfig()
	if err != nil {
		return fmt.Errorf("failed to load database config: %w", err)
	}

	// Set environment variables from Vault
	os.Setenv("DB_HOST", dbConfig.Host)
	os.Setenv("DB_PORT", dbConfig.Port)
	os.Setenv("DB_USER", dbConfig.Username)
	os.Setenv("DB_PASSWORD", dbConfig.Password)
	os.Setenv("DB_NAME", dbConfig.Database)

	// Load Clerk config
	clerkConfig, err := vaultClient.GetClerkConfig()
	if err != nil {
		return fmt.Errorf("failed to load clerk config: %w", err)
	}

	os.Setenv("CLERK_SECRET_KEY", clerkConfig.SecretKey)
	os.Setenv("CLERK_WEBHOOK_SECRET", clerkConfig.WebhookSecret)

	// Load app config
	appConfig, err := vaultClient.GetAppConfig()
	if err != nil {
		return fmt.Errorf("failed to load app config: %w", err)
	}

	os.Setenv("JWT_SECRET", appConfig.JWTSecret)
	os.Setenv("ENVIRONMENT", appConfig.Environment)

	log.Println("✅ Configuration loaded from Vault successfully")
	return nil
}
