package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	grpcserver "desafio-clean-architecture/internal/transport/grpc"
	graphqltransport "desafio-clean-architecture/internal/transport/graphql"
	resttransport "desafio-clean-architecture/internal/transport/rest"
	postgresrepo "desafio-clean-architecture/internal/repository/postgres"
	"desafio-clean-architecture/internal/infrastructure/db"
	orderspb "desafio-clean-architecture/internal/proto/orderspb"
	"desafio-clean-architecture/internal/usecase"
	"google.golang.org/grpc"
)

type config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	RESTAddr   string
	GRPCAddr   string
	GraphQLAddr string
}

func main() {
	cfg := loadConfig()

	conn, err := waitForDatabase(db.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Name:     cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	})
	if err != nil {
		log.Fatalf("database init: %v", err)
	}
	defer conn.Close()

	migrationsDir := os.Getenv("MIGRATIONS_PATH")
	if migrationsDir == "" {
		migrationsDir = filepath.Join(".", "internal", "db", "migrations")
	}
	if err := db.RunMigrations(conn, migrationsDir); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	repo := postgresrepo.NewOrderRepository(conn)
	listUseCase := usecase.NewListOrdersUseCase(repo)
	createUseCase := usecase.NewCreateOrderUseCase(repo)

	go startREST(cfg.RESTAddr, resttransport.NewHandler(listUseCase, createUseCase))
	go startGRPC(cfg.GRPCAddr, listUseCase)
	go startGraphQL(cfg.GraphQLAddr, listUseCase)

	log.Printf("application started")
	select {}
}

func loadConfig() config {
	return config{
		DBHost:      getenv("DB_HOST", "localhost"),
		DBPort:      getenv("DB_PORT", "5432"),
		DBUser:      getenv("DB_USER", "orders"),
		DBPassword:  getenv("DB_PASSWORD", "orders"),
		DBName:      getenv("DB_NAME", "orders"),
		DBSSLMode:   getenv("DB_SSLMODE", "disable"),
		RESTAddr:    getenv("REST_ADDR", ":8080"),
		GRPCAddr:    getenv("GRPC_ADDR", ":50051"),
		GraphQLAddr: getenv("GRAPHQL_ADDR", ":8081"),
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func waitForDatabase(cfg db.Config) (*sql.DB, error) {
	var conn *sql.DB
	var err error

	for i := 0; i < 30; i++ {
		conn, err = db.Connect(cfg)
		if err == nil {
			return conn, nil
		}
		log.Printf("waiting for database: %v", err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("database unavailable: %w", err)
}

func startREST(addr string, handler http.Handler) {
	log.Printf("REST server listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("rest server: %v", err)
	}
}

func startGRPC(addr string, listUseCase *usecase.ListOrdersUseCase) {
	grpcServer := grpc.NewServer()
	orderspb.RegisterOrderServiceServer(grpcServer, grpcserver.NewServer(listUseCase))

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("grpc listen: %v", err)
	}

	log.Printf("gRPC server listening on %s", addr)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("grpc server: %v", err)
	}
}

func startGraphQL(addr string, listUseCase *usecase.ListOrdersUseCase) {
	mux := http.NewServeMux()
	mux.Handle("/graphql", graphqltransport.NewHandler(listUseCase))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	log.Printf("GraphQL server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("graphql server: %v", err)
	}
}
