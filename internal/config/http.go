package config

import (
	"net"
	"os"
)

type HttpConfig struct {
	host string
	port string

	musicInfoServiceURL string
}

func (c *HttpConfig) MusicInfoServiceURL() string {
	return c.musicInfoServiceURL
}

func NewHttpConfig() *HttpConfig {
	host := os.Getenv("HTTP_HOST")
	if host == "" {
		panic("HTTP_HOST environment variable is empty")
	}
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		panic("HTTP_PORT environment variable is empty")
	}

	musicInfoServiceURL := os.Getenv("MUSIC_INFO_SERVICE_URL")
	if musicInfoServiceURL == "" {
		panic("MUSIC_INFO_SERVICE_URL environment variable is empty")
	}

	return &HttpConfig{
		host:                host,
		port:                port,
		musicInfoServiceURL: musicInfoServiceURL,
	}
}

func (c *HttpConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}
