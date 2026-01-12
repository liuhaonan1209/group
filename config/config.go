package config

type Config struct {
	Mysql Mysql
	Redis Redis
}
type Mysql struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}
type Redis struct {
	Password string
	Host     string
	Port     int
}
