package grpcserver

import (
	"context"
	"errors"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/sviatilnik/url-shortener/internal/app/models"
	"github.com/sviatilnik/url-shortener/internal/app/shortener"
	"github.com/sviatilnik/url-shortener/internal/app/storages"
	pb "github.com/sviatilnik/url-shortener/proto"
)

// ShortenerServer реализует gRPC сервис сокращения ссылок.
type ShortenerServer struct {
	pb.UnimplementedShortenerServiceServer
	shortener *shortener.Shortener
	logger    *zap.SugaredLogger
}

func NewShortenerServer(shorter *shortener.Shortener, logger *zap.SugaredLogger) *ShortenerServer {
	return &ShortenerServer{
		shortener: shorter,
		logger:    logger,
	}
}

func (s *ShortenerServer) ShortenURL(ctx context.Context, req *pb.URLShortenRequest) (*pb.URLShortenResponse, error) {
	url := strings.TrimSpace(req.GetUrl())
	if url == "" {
		return nil, status.Error(codes.InvalidArgument, "url is required")
	}

	result, err := s.shortener.GenerateShortLink(ctx, url)
	if err != nil {
		switch {
		case errors.Is(err, shortener.ErrInvalidURL):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, shortener.ErrLinkConflict):
			// Возвращаем уже существующую ссылку без ошибки, как и HTTP-обработчик.
			return &pb.URLShortenResponse{Result: result}, nil
		default:
			s.logger.Errorw("Failed to shorten URL", "error", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &pb.URLShortenResponse{Result: result}, nil
}

func (s *ShortenerServer) ExpandURL(ctx context.Context, req *pb.URLExpandRequest) (*pb.URLExpandResponse, error) {
	id := strings.TrimSpace(req.GetId())
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	link, err := s.shortener.GetFullLinkByShortCode(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, shortener.ErrIDIsRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, storages.ErrKeyNotFound):
			return nil, status.Error(codes.NotFound, "link not found")
		default:
			s.logger.Errorw("Failed to expand URL", "error", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &pb.URLExpandResponse{Result: link.OriginalURL}, nil
}

func (s *ShortenerServer) ListUserURLs(ctx context.Context, _ *emptypb.Empty) (*pb.UserURLsResponse, error) {
	userID, _ := ctx.Value(models.ContextUserID).(string)
	if strings.TrimSpace(userID) == "" {
		return nil, status.Error(codes.Unauthenticated, "authorization required")
	}

	links, err := s.shortener.GetUserLinks(ctx, userID)
	if err != nil {
		s.logger.Errorw("Failed to list user URLs", "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	resp := &pb.UserURLsResponse{}
	for _, link := range links {
		resp.Url = append(resp.Url, &pb.URLData{
			ShortUrl:    link.ShortURL,
			OriginalUrl: link.OriginalURL,
		})
	}

	return resp, nil
}
