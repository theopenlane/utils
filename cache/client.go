package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config for the redis client used to store key-value pairs
type Config struct {
	// Enabled to enable redis client in the server
	Enabled bool `json:"enabled" koanf:"enabled" default:"true"`
	// Address is the host:port to connect to redis
	Address string `json:"address" koanf:"address" default:"localhost:6379"`
	// Name of the connecting client
	Name string `json:"name" koanf:"name" default:""`
	// Username to connect to redis
	Username string `json:"username" koanf:"username"`
	// Password, must match the password specified in the server configuration
	Password string `json:"password" koanf:"password"`
	// DB to be selected after connecting to the server, 0 uses the default
	DB int `json:"db" koanf:"db" default:"0"`
	// Dial timeout for establishing new connections, defaults to 5s
	DialTimeout time.Duration `json:"dialtimeout" koanf:"dialtimeout" default:"5s"`
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Supported values:
	//   - `0` - default timeout (3 seconds).
	//   - `-1` - no timeout (block indefinitely).
	//   - `-2` - disables SetReadDeadline calls completely.
	ReadTimeout time.Duration `json:"readtimeout" koanf:"readtimeout" default:"0"`
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.  Supported values:
	//   - `0` - default timeout (3 seconds).
	//   - `-1` - no timeout (block indefinitely).
	//   - `-2` - disables SetWriteDeadline calls completely.
	WriteTimeout time.Duration `json:"writetimeout" koanf:"writetimeout" default:"0"`
	// MaxRetries before giving up.
	// Default is 3 retries; -1 (not 0) disables retries.
	MaxRetries int `json:"maxretries" koanf:"maxretries" default:"3"`
	// MinIdleConns is useful when establishing new connection is slow.
	// Default is 0. the idle connections are not closed by default.
	MinIdleConns int `json:"minidleconns" koanf:"minidleconns" default:"0"`
	// Maximum number of idle connections.
	// Default is 0. the idle connections are not closed by default.
	MaxIdleConns int `json:"maxidleconns" koanf:"maxidleconns" default:"0"`
	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActiveConns int `json:"maxactiveconns" koanf:"maxactiveconns" default:"0"`
}

// New returns a new redis client based on the configuration settings
func New(c Config) *redis.Client {
	opts := &redis.Options{
		Addr:            c.Address,
		DB:              c.DB,
		DialTimeout:     c.DialTimeout,
		ReadTimeout:     c.ReadTimeout,
		WriteTimeout:    c.WriteTimeout,
		MaxRetries:      c.MaxRetries,
		MinIdleConns:    c.MinIdleConns,
		MaxIdleConns:    c.MaxIdleConns,
		MaxActiveConns:  c.MaxActiveConns,
		DisableIdentity: true,
	}

	// optional fields
	if c.Name != "" {
		opts.ClientName = c.Name
	}

	if c.Username != "" {
		opts.Username = c.Username
	}

	if c.Password != "" {
		opts.Password = c.Password
	}

	return redis.NewClient(opts)
}

// Healthcheck pings the client to check if the connection is working
func Healthcheck(c *redis.Client) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		// check if its alive
		if err := c.Ping(ctx).Err(); err != nil {
			return err
		}

		return nil
	}
}
