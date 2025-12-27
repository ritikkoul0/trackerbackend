package config

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB     *mongo.Database
	Client *mongo.Client
)

func ConnectDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable not set")
	}

	// Load system CA certs
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal("Failed to load system certs:", err)
	}

	tlsConfig := &tls.Config{
		RootCAs:    rootCAs,
		MinVersion: tls.VersionTLS12,
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	clientOpts := options.Client().
		ApplyURI(mongoURI).
		SetServerAPIOptions(serverAPI).
		SetTLSConfig(tlsConfig).
		SetAppName("Cluster0").
		SetConnectTimeout(20 * time.Second)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	log.Println("MongoDB connected successfully")

	DB = client.Database("investmentdb")
	Client = client
}

func DisconnectDatabase() {
	if Client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.Disconnect(ctx); err != nil {
		log.Fatal("Failed to disconnect MongoDB:", err)
	}

	log.Println("MongoDB disconnected")
}
