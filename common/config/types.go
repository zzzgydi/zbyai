package config

type AppConfig struct {
	HttpPort      int            `yaml:"http_port"`
	JwtSecret     string         `yaml:"jwt_secret"`
	Search        SearchConfig   `yaml:"search"`
	Gpt           GPTConfig      `yaml:"gpt"`
	Redis         RedisConfig    `yaml:"redis"`
	Postgres      PostgresConfig `yaml:"postgres"`
	Mysql         MysqlConfig    `yaml:"mysql"`
	Supabase      SupabaseConfig `yaml:"supabase"`
	ChatModels    []ModelConfig  `yaml:"chat_models"`
	RewriteModels []ModelConfig  `yaml:"rewrite_models"`
}

type GPTConfig struct {
	Endpoint string `yaml:"endpoint"`
	ApiKey   string `yaml:"api_key"`
	ProxyURL string `yaml:"proxy_url,omitempty"`
}

type SearchConfig struct {
	Github  *GithubSearchConfig `yaml:"github,omitempty"`
	Serper  *SearchItemConfig   `yaml:"serper,omitempty"`
	Bing    *SearchItemConfig   `yaml:"bing,omitempty"`
	Searxng *SearchItemConfig   `yaml:"searxng,omitempty"`
}

type GithubSearchConfig struct {
	Tokens []string `yaml:"tokens"`
}

type SearchItemConfig struct {
	Endpoint string `yaml:"endpoint"`
	ApiKey   string `yaml:"api_key"`
}

type RedisConfig struct {
	Url      string `yaml:"url"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

type PostgresConfig struct {
	Host          string `yaml:"host"`
	Port          string `yaml:"port"`
	DbName        string `yaml:"db_name"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	SSLMode       string `yaml:"sslmode"`
	SlowThreshold int    `yaml:"slow_threshold,omitempty"` // ms
	AutoMigrate   bool   `yaml:"auto_migrate"`
}

type MysqlConfig struct {
	Host          string `yaml:"host"`
	Port          string `yaml:"port"`
	DbName        string `yaml:"db_name"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	SlowThreshold int    `yaml:"slow_threshold,omitempty"` // ms
	AutoMigrate   bool   `yaml:"auto_migrate"`
}

type SupabaseConfig struct {
	Url       string `yaml:"url"`
	Key       string `yaml:"key"`
	JwtSecret string `yaml:"jwt_secret"`
}

type ModelConfig struct {
	Model   string `yaml:"model"`             // model value for provider
	Display string `yaml:"display,omitempty"` // display model name
	Weight  int    `yaml:"weight,omitempty"`
}
