package config

type Config struct {
	Server struct {
		Port       int    `envconfig:"PORT"`
		Host       string `envconfig:"SERVER_HOST"`
		GcpProject string `envconfig:"GCP_PROJECT"`
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
		BaseUrl1 string `envconfig:"MYSHARE_BASE_URL1"`
		BaseUrl2 string `envconfig:"MYSHARE_BASE_URL2"`
		Audience string `envconfig:"MYSHARE_AUDIENCE"`
	}
	ContributionsAPI struct {
		BaseUrl  string `envconfig:"CONTRIBUTIONS_BASE_URL"`
		Audience string `envconfig:"CONTRIBUTIONS_AUDIENCE"`
	}
	Ngrok struct {
		AuthToken string `envconfig:"NGROK_AUTHTOKEN"`
		Domain    string `envconfig:"NGROK_DOMAIN"`
	}
}
