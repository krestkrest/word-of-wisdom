package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/krestkrest/word-of-wisdom/internal/network/tcp"
	"github.com/krestkrest/word-of-wisdom/internal/service"
	"github.com/krestkrest/word-of-wisdom/internal/storage/file"
)

const (
	address = "127.0.0.1"
)

var (
	fileName   = flag.String("f", "", "file containing quotes, required")
	complexity = flag.Uint("c", 20, "complexity of POW task")
	port       = flag.Uint("p", 3434, "TCP port for incoming connections")
	level      = flag.String("l", "INFO", "log level")
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	flag.Parse()

	if len(*fileName) == 0 {
		log.Fatal("File with quotes should be defined")
	}
	if *complexity == 0 {
		log.Fatal("Complexity should not be zero")
	}

	logLevel, err := log.ParseLevel(*level)
	if err != nil {
		log.Fatalf("Incorrect log level %s", *level)
	}
	log.SetLevel(logLevel)

	ctx, cancel := context.WithCancel(context.Background())
	logger := log.WithContext(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigChan
		cancel()
	}()

	storage := file.NewStorage(*fileName)
	if err := storage.Start(); err != nil {
		logger.Fatal("Failed to start storage", err)
	}

	fullAddress := address + ":" + strconv.Itoa(int(*port))

	s := service.NewService(storage, uint8(*complexity))
	server := tcp.NewServer(fullAddress, s)
	if err := server.Start(ctx); err != nil {
		logger.Fatal("Error executing server", err)
	}
	logger.Info("Server stopped")
}
