package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func getNTPTime(server string) (time.Time, error) {
	t, err := ntp.Time(server)
	if err != nil {
		return time.Time{}, errors.Join(err, errors.New("cannot get time from NTP"))
	}
	return t, nil
}
func main() {

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	server := flag.String("server", "pool.ntp.org", "NTP server to fetch time from")
	flag.Parse()

	currentTime, err := getNTPTime(*server)
	if err != nil {
		logger.Error("Failed to fetch NTP time", "error", err, "server", *server)
		os.Exit(1)
	}

	fmt.Println("Current time:", currentTime.Format("2006-01-02 15:04:05 MST"))
}
