package cmd

import (
	"log"

	"github.com/imimran/go-grpc-auth/internal/delivery/grpc"
	"github.com/imimran/go-grpc-auth/internal/infrastructure"
	"github.com/imimran/go-grpc-auth/internal/repository"
	"github.com/imimran/go-grpc-auth/internal/usecase"

	"net"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		dbDsn := viper.GetString("database.dsn")
		port := viper.GetString("server.port")
		jwtSecret := []byte(viper.GetString("jwt.secret"))

		db := infrastructure.NewPostgresDB(dbDsn)

		userRepo := repository.NewUserRepository(db)
		userUsecase := usecase.NewUserUsecase(userRepo, jwtSecret)

		handler := grpc.NewUserHandler(userUsecase)

		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterUserServiceServer(grpcServer, handler)

		log.Printf("gRPC server listening on %s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}

func init() {
	serveCmd.Flags().String("config", "configs/config.yaml", "config file")
	_ = viper.BindPFlag("config", serveCmd.Flags().Lookup("config"))
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	configFile := viper.GetString("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}
}
