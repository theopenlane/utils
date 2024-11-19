package testutils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"entgo.io/ent/dialect"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

const (
	defaultMaxConn       = 100
	defaultExpiryMinutes = 5
)

// Option is a functional option for configuring the DBConfig
type Option func(*DBConfig)

// TestFixture is a struct that holds the pool, resource and URI for the database
type TestFixture struct {
	Pool     *dockertest.Pool
	resource *dockertest.Resource
	URI      string
	Dialect  string
}

// DBConfig is a struct that holds the configuration for the database
type DBConfig struct {
	// image is the docker image to use
	image string

	// expiry is the duration after which the container will be killed
	expiry time.Duration

	// maxConn is the maximum number of connections to the database
	// defaults to 100
	maxConn int
}

// WithImage sets the image used for docker tests
func WithImage(image string) Option {
	return func(c *DBConfig) {
		c.image = image
	}
}

// WithExpiryMinutes sets the expiration of the container, defaults to 10 minutes
func WithExpiryMinutes(expiry int) Option {
	return func(c *DBConfig) {
		c.expiry = time.Duration(expiry) * time.Minute
	}
}

// WithMaxConn sets the max connections to the db container, defaults to 100
func WithMaxConn(maxConn int) Option {
	return func(c *DBConfig) {
		c.maxConn = maxConn
	}
}

func TeardownFixture(tf *TestFixture) {
	if tf.Pool != nil {
		if err := tf.Pool.Purge(tf.resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}
}
func GetPostgresDockerTest(image string, expiry time.Duration, maxConn int) (*TestFixture, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, err
	}

	defaultImg := "postgres"
	imgTag := "alpine"

	if strings.Contains(image, ":") {
		p := strings.SplitN(image, ":", 2) //nolint:mnd
		imgTag = p[1]
	}

	password := "password"
	user := "postgres"
	dbName := "postgres"

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: defaultImg,
			Tag:        imgTag,
			Cmd:        []string{"postgres", "-c", fmt.Sprintf("max_connections=%d", maxConn)},
			Env: []string{
				fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
				fmt.Sprintf("POSTGRES_USER=%s", user),
				fmt.Sprintf("POSTGRES_DB=%s", dbName),
				"listen_addresses='*'",
			},
		}, func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		})
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}

	port := resource.GetPort("5432/tcp")

	// when running locally, the host is localhost
	// however, when running in a CI environment, and using docker-in-docker
	// the host is the docker host network
	// - `host.docker.internal` on mac
	// - `172.17.0.1` on linux
	host := os.Getenv("TEST_DB_HOST")
	if host == "" {
		host = "localhost"
	}

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

	if err := resource.Expire(uint(expiry.Seconds())); err != nil {
		return nil, err
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		db, err := sql.Open("postgres", databaseURL)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("unable to connect to database: %s", err)
	}

	return &TestFixture{
		Pool:     pool,
		resource: resource,
		URI:      databaseURL,
		Dialect:  dialect.Postgres,
	}, nil
}

func getTestDB(u string, expiry time.Duration, maxConn int) (*TestFixture, error) {
	switch {
	case strings.HasPrefix(u, "postgres"):
		return GetPostgresDockerTest(u, expiry, maxConn)
	default:
		return nil, newURIError(u)
	}
}

// GetTestURI returns the dialect, connection string and if used a testcontainer for database connectivity in tests
func GetTestURI(opts ...Option) *TestFixture {
	config := &DBConfig{
		maxConn: defaultMaxConn,                                    // default to 100
		expiry:  time.Duration(defaultExpiryMinutes) * time.Minute, // default to 5 minutes
	}

	for _, opt := range opts {
		opt(config)
	}

	switch {
	case strings.HasPrefix(config.image, "sqlite://"):
		// return dialect.SQLite, strings.TrimPrefix(u, "sqlite://")
		return &TestFixture{Dialect: dialect.SQLite, URI: strings.TrimPrefix(config.image, "sqlite://")}
	case strings.HasPrefix(config.image, "libsql://"):
		// return dialect.SQLite, strings.TrimPrefix(u, "libsql://")
		return &TestFixture{Dialect: "libsql", URI: strings.TrimPrefix(config.image, "libsql://")}
	case strings.HasPrefix(config.image, "postgres://"), strings.HasPrefix(config.image, "postgresql://"):
		// return dialect.Postgres, u
		return &TestFixture{Dialect: dialect.Postgres, URI: config.image}
	case strings.HasPrefix(config.image, "docker://"):
		tf, err := getTestDB(strings.TrimPrefix(config.image, "docker://"), config.expiry, config.maxConn)
		if err != nil {
			panic(err)
		}

		tf.Dialect = dialect.Postgres

		return tf
	default:
		panic("invalid DB URI, uri: " + config.image)
	}
}
