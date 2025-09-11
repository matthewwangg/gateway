package logger

import (
	"context"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/matthewwangg/gateway/internal/pb"
)

type RemoteLogger struct {
	Source   string
	Hostname string
}

func NewRemoteLogger(source string, hostname string) *RemoteLogger {
	return &RemoteLogger{
		Source:   source,
		Hostname: hostname,
	}
}

func (rl *RemoteLogger) LogWithLevel(message string, level string) {
	connection, err := grpc.NewClient(os.Getenv("GCP_LOGGER"), grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		return
	}
	defer func() {
		if err := connection.Close(); err != nil {
			return
		}
	}()

	client := pb.NewLogServiceClient(connection)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	apiKey := os.Getenv("LOG_SERVICE_API_KEY")
	if apiKey != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", apiKey)
	}

	_, err = client.SendLog(ctx, &pb.LogEntry{
		Timestamp: timestamppb.New(time.Now()),
		Message:   message,
		Source:    rl.Source,
		Level:     level,
		Hostname:  rl.Hostname,
	})
	if err != nil {
		return
	}
}

func (rl *RemoteLogger) Debug(message string) {
	rl.LogWithLevel(message, "DEBUG")
}

func (rl *RemoteLogger) Info(message string) {
	rl.LogWithLevel(message, "INFO")
}

func (rl *RemoteLogger) Warn(message string) {
	rl.LogWithLevel(message, "WARN")
}

func (rl *RemoteLogger) Error(message string) {
	rl.LogWithLevel(message, "ERROR")
}
