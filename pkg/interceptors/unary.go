package interceptors

import (
	"context"

	"google.golang.org/grpc"
)

func (I *Interceptor) UnaryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	if err := I.inject(ctx); err != nil {
		return nil, err
	}

	return handler(ctx, req)
}
