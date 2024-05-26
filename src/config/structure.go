package config

/*
config file is a yaml file with the following structure:

routes:
	/test:
		proxy_pass: http://localhost:8000
		proxy_set_header:
			- X-Forwarded-Host: $host
		proxy_hide_header:
			- Server
		proxy_hide_query:
			- key
		proxy_set_cookie:
			- key: value
		proxy_hide_cookie:
			- key
		trim_prefix: /test
		trim_suffix: .html

analytics:
	username: admin
	password: admin
	enabled: true
	port: 4000
	extra_queries:
		key: value
*/

type Config struct {
	Routes    Route     `mapstructure:"routes"`
	Analytics Analytics `mapstructure:"analytics"`
}

type Route []RouteConfig

type RouteConfig struct {
	Location        string   `mapstructure:"location"`
	ProxyPass       string   `mapstructure:"proxy_pass"`
	ProxySetHeader  []Header `mapstructure:"proxy_set_header"`
	ProxyHideHeader []string `mapstructure:"proxy_hide_header"`
	PorxySetCookie  []Cookie `mapstructure:"proxy_set_cookie"`
	ProxyHideCookie []string `mapstructure:"proxy_hide_cookie"`
	TrimPrefix      string   `mapstructure:"trim_prefix"`
	TrimSuffix      string   `mapstructure:"trim_suffix"`
}

type Header map[string]string

type Cookie map[string]string
type Analytics struct {
	Username             string `mapstructure:"username"`
	Password             string `mapstructure:"password"`
	Enabled              bool   `mapstructure:"enabled"`
	Port                 int    `mapstructure:"port"`
	DisableSystemMetrics bool   `mapstructure:"disable_system_metrics"`
}
