package config

type Config struct {
	Server struct {
		Port int    `envconfig:"PORT"`
		Host string `envconfig:"SERVER_HOST"`
	}
	Telegram struct {
		ApiKey  string `envconfig:"TELEGRAM_API_KEY"`
		BotName string `envconfig:"TELEGRAM_BOT_NAME"`
	}
	Oauth struct {
		Domain       string `envconfig:"OAUTH_DOMAIN"`
		ClientID     string `envconfig:"OAUTH_CLIENT_ID"`
		ClientSecret string `envconfig:"OAUTH_CLIENT_SECRET"`
		Secret       string `envconfig:"SECRET"`
	}
	MyshareAPI struct {
		BaseUrl  string `envconfig:"MYSHARE_BASE_URL"`
		Audience string `envconfig:"MYSHARE_AUDIENCE"`
	}
}
