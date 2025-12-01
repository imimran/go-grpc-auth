package cmd

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/imimran/go-grpc-auth/config"
	"github.com/imimran/go-grpc-auth/infrastructure"
	pb "github.com/imimran/go-grpc-auth/proto"
	grpcDelivery "github.com/imimran/go-grpc-auth/user/delivery/grpc"
	"github.com/imimran/go-grpc-auth/user/repository"
	"github.com/imimran/go-grpc-auth/user/usecase"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start gRPC server with grpc-gateway",
	Run:   serve,
}

func serve(cmd *cobra.Command, args []string) {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config load error: %v", err)
	}

	// Connect DB (GORM DB)
	postgresDB, err := db.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Run migration
	if err := AutoMigrate(postgresDB); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Setup user repository, usecase, handler
	userRepo := repository.NewUserRepository(postgresDB)
	userUsecase := usecase.NewUserUsecase(userRepo, []byte(cfg.JWT.Secret))
	userHandler := grpcDelivery.NewUserHandler(userUsecase)

	grpcPort := cfg.Server.GRPCPort   // e.g. ":50051"
	httpPort := cfg.Server.HTTPPort         // grpc-gateway port

	// gRPC server listener
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userHandler)
	reflection.Register(grpcServer) // optional

	// Graceful shutdown channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run gRPC server goroutine
	go func() {
		log.Printf("Starting gRPC server on %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	// Setup grpc-gateway
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	grpcEndpoint := grpcPort
	if strings.HasPrefix(grpcEndpoint, ":") {
		grpcEndpoint = "localhost" + grpcEndpoint
	}

	if err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts); err != nil {
		log.Fatalf("Failed to start HTTP gateway: %v", err)
	}

	httpServer := &http.Server{
		Addr:    httpPort,
		Handler: mux,
	}

	// Run HTTP server goroutine
	go func() {
		log.Printf("Starting HTTP gateway on %s", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP gateway failed: %v", err)
		}
	}()

	// Wait for stop signal
	<-stop
	log.Println("Shutting down servers...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server Shutdown error: %v", err)
	}

	grpcServer.GracefulStop()

	log.Println("Servers stopped gracefully")
}


