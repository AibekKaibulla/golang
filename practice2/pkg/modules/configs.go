package modules
import "time"

type PostgresConfig struct {
	Host string
	Port string
	Username string
	Password string
	DBName string
	SSLMode string
	ExecTimeout time.Duration
}