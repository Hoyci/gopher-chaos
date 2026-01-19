package interceptors

import (
	"google.golang.org/grpc"
)

type wrappedStream struct {
	grpc.ServerStream
	interceptor *Interceptor
}

func (w *wrappedStream) RecvMsg(m any) error {
	if err := w.interceptor.inject(w.Context()); err != nil {
		return err
	}
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m any) error {
	if err := w.interceptor.inject(w.Context()); err != nil {
		return err
	}
	return w.ServerStream.SendMsg(m)
}

func (I *Interceptor) StreamInterceptor(
	srv any,
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	if err := I.inject(ss.Context()); err != nil {
		return err
	}

	wrapper := &wrappedStream{
		ServerStream: ss,
		interceptor:  I,
	}

	return handler(srv, wrapper)
}
