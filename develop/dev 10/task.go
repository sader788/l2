package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Config struct {
	timeout time.Duration
	host    string
	port    string
}

func parseConfig() (Config, error) {
	cfg := Config{}
	cfg.timeout = time.Duration(10 * time.Second)

	for i := 1; i < len(os.Args); i++ {
		if strings.HasPrefix(os.Args[i], "--timeout=") {
			timeout, err := parseTimeout(os.Args[i])
			if err != nil {
				return Config{}, err
			}
			cfg.timeout = timeout
		} else if len(cfg.host) == 0 {
			cfg.host = os.Args[i]
		} else if len(cfg.port) == 0 {
			cfg.port = os.Args[i]
		} else {
			break
		}
	}

	if len(cfg.host) == 0 || len(cfg.port) == 0 {
		return Config{}, errors.New("bad args")
	}

	return Config{}, nil
}

func parseTimeout(str string) (time.Duration, error) {
	timeStr := str[len("--timeout="):]

	value, err := strconv.Atoi(timeStr)
	if err != nil {
		return 10 * time.Second, err
	}

	switch timeStr[len(timeStr)-1] {
	case 's':
		return time.Duration(value) * time.Second, nil
	case 'm':
		return time.Duration(value) * time.Minute, nil
	case 'h':
		return time.Duration(value) * time.Hour, nil
	default:
		return 10 * time.Second, errors.New("bad timeout")
	}
}

func telnet(cfg Config) error {
	notify := make(chan os.Signal)

	signal.Notify(notify, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	_, err := strconv.Atoi(cfg.port)
	if err != nil {
		return errors.New("bad port")
	}

	conn, err := net.DialTimeout("tcp", cfg.host+":"+cfg.port, cfg.timeout)
	if err != nil {
		return err
	}

	go func(conn net.Conn) {
		<-notify
		err := conn.Close()
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}(conn)

	_ = conn.SetReadDeadline(time.Now().Add(time.Duration(2) * time.Second))
	for {
		fmt.Print(">")

		_, err = io.Copy(conn, os.Stdin)

		_ = conn.SetReadDeadline(time.Now().Add(time.Duration(700) * time.Millisecond))

		_, err := io.Copy(os.Stdout, conn)

		if errors.Is(err, net.ErrClosed) {
			os.Exit(0)
		}

		fmt.Println()
	}
}

func main() {
	cfg, err := parseConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = telnet(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
