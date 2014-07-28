package auth

var (
	mock_auth_config *Config = &Config{
		Type:              BUILTIN,
		AllowRegistration: true,
		Username:          "admin",
		Password:          "admin",
	}
)
