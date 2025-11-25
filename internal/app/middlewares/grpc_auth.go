package middlewares

import (
	"context"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/sviatilnik/url-shortener/internal/app/config"
	"github.com/sviatilnik/url-shortener/internal/app/models"
)

// GRPCAuthInterceptor обеспечивает авторизацию клиентов gRPC по metadata authorization.
type GRPCAuthInterceptor struct {
	config *config.Config
	logger *zap.SugaredLogger
}

func NewGRPCAuthInterceptor(cfg *config.Config, logger *zap.SugaredLogger) *GRPCAuthInterceptor {
	return &GRPCAuthInterceptor{
		config: cfg,
		logger: logger,
	}
}

func (i *GRPCAuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		key := i.config.AuthSecret

		token := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get("authorization")
			if len(values) > 0 {
				token = values[0]
			}
		}

		token = strings.TrimSpace(token)
		userID := ""
		if token != "" {
			userID = getUserID(key, token)
		}

		if strings.TrimSpace(userID) == "" {
			var err error
			userID, err = generateUserID()
			if err != nil {
				i.logger.Errorw("Failed to generate user ID", "error", err)
				return nil, status.Error(codes.Internal, "internal error")
			}

			signed := signUserID(key, userID)
			if strings.TrimSpace(signed) != "" {
				if err := grpc.SetHeader(ctx, metadata.Pairs("authorization", signed)); err != nil {
					i.logger.Warnw("Failed to set authorization metadata", "error", err)
				}
			}
		}

		ctx = context.WithValue(ctx, models.ContextUserID, userID)

		return handler(ctx, req)
	}
}
