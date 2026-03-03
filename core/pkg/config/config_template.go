package config

type Config struct {
	ContainerName string
	Adminer       string
	Hertz         HertzConfig
	PgSQL         PgSQLConfig
	OAuth         OAuthConfig
}

type HertzConfig struct {
	ListenAddr  string
	ListenPort  string
	MonitorPort string
}

type PgSQLConfig struct {
	Addr   string
	Port   string
	User   string
	Passwd string
	DBName string
}

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectUrl  string
}
