package service

import (
	"errors"
	"fmt"

	"github.com/Falokut/grpc_errors"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInternal        = errors.New("internal")
	ErrImageTooLarge   = errors.New("image too large")
	ErrImageTooSmall   = errors.New("image too small")
	ErrNoInstructions  = errors.New("no instructions received")
	ErrInvalidArgument = errors.New("invalid received data")
)

var errorCodes = map[error]codes.Code{
	ErrInternal:        codes.Internal,
	ErrImageTooLarge:   codes.InvalidArgument,
	ErrImageTooSmall:   codes.InvalidArgument,
	ErrNoInstructions:  codes.InvalidArgument,
	ErrInvalidArgument: codes.InvalidArgument,
}

type errorHandler struct {
	logger *logrus.Logger
}

func newErrorHandler(logger *logrus.Logger) errorHandler {
	return errorHandler{
		logger: logger,
	}
}

func (e *errorHandler) createErrorResponceWithSpan(span opentracing.Span, err error, developerMessage string) error {
	if err == nil {
		return nil
	}

	span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))
	ext.LogError(span, err)
	return e.createErrorResponce(err, developerMessage)
}

func (e *errorHandler) createErrorResponce(err error, developerMessage string) error {
	var msg string
	if len(developerMessage) == 0 {
		msg = err.Error()
	} else {
		msg = fmt.Sprintf("%s. error: %v", developerMessage, err)
	}

	err = status.Error(grpc_errors.GetGrpcCode(err), msg)
	e.logger.Error(err)
	return err
}

func init() {
	grpc_errors.RegisterErrors(errorCodes)
}
