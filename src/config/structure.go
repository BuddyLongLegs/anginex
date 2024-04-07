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
*/

type Config struct {
	Routes Route `mapstructure:"routes"`
}

type Route map[string]RouteConfig

type RouteConfig struct {
	ProxyPass       string   `mapstructure:"proxy_pass"`
	ProxySetHeader  []Header `mapstructure:"proxy_set_header"`
	ProxyHideHeader []string `mapstructure:"proxy_hide_header"`
	PorxySetCookie  []Cookie `mapstructure:"proxy_set_cookie"`
	ProxyHideCookie []string `mapstructure:"proxy_hide_cookie"`
	TrimPrefix      string   `mapstructure:"trim_prefix"`
	TrimSuffix      string   `mapstructure:"trim_suffix"`
}

type Header struct {
	Key   string `mapstructure:"key"`
	Value string `mapstructure:"value"`
}

type Cookie struct {
	Key   string `mapstructure:"key"`
	Value string `mapstructure:"value"`
}