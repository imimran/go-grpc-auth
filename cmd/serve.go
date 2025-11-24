package cmd

import (
	"log"
	"net"

	"github.com/imimran/go-grpc-auth/config"
	"github.com/imimran/go-grpc-auth/infrastructure"
	pb "github.com/imimran/go-grpc-auth/proto"
	grpcDelivery "github.com/imimran/go-grpc-auth/user/delivery/grpc"
	"github.com/imimran/go-grpc-auth/user/repository"
	"github.com/imimran/go-grpc-auth/user/usecase"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start gRPC server",
	Run: func(cmd *cobra.Command, args []string) {

		// Load full config using Viper
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("Config load error: %v", err)
		}

		// Build DB + connect
		postgresDB, err := db.NewPostgresDB(cfg.Database)
		if err != nil {
			log.Fatalf("Database connection failed: %v", err)
		}

		// Run AutoMigration
		if err := AutoMigrate(postgresDB); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

		// Layer wiring (Repository → Usecase → Handler)
		userRepo := repository.NewUserRepository(postgresDB)
		userUsecase := usecase.NewUserUsecase(userRepo, []byte(cfg.JWT.Secret))
		userHandler := grpcDelivery.NewUserHandler(userUsecase)

		// Start gRPC listener
		lis, err := net.Listen("tcp", cfg.Server.Port)
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterUserServiceServer(grpcServer, userHandler)

		log.Printf("gRPC server running on port %s", cfg.Server.Port)

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	},
}

func init() {
	serveCmd.Flags().String("config", "config.yaml", "config file path")
	_ = viper.BindPFlag("config", serveCmd.Flags().Lookup("config"))

	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(serveCmd)
}

func initConfig() {
	configPath := viper.GetString("config")

	viper.SetConfigFile(configPath) // REQUIRED

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
}
