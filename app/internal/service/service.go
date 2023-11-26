package service

import (
	"context"
	"fmt"
	"image"
	"strings"

	"github.com/Falokut/grpc_errors"
	"github.com/Falokut/image_processing_service/internal/image_processing"
	image_service "github.com/Falokut/image_processing_service/pkg/image_processing_service/v1/protos"
	"github.com/gabriel-vasile/mimetype"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/api/httpbody"
)

type ImageProcessingService struct {
	image_service.UnimplementedImageProcessingServiceV1Server
	logger           *logrus.Logger
	errorHandler     errorHandler
	imagesProcessing image_processing.ImagesProcessing
}

func NewImagesProcessingService(logger *logrus.Logger,
	imagesProcessing image_processing.ImagesProcessing) *ImageProcessingService {
	errorHandler := newErrorHandler(logger)
	return &ImageProcessingService{
		logger:           logger,
		errorHandler:     errorHandler,
		imagesProcessing: imagesProcessing,
	}
}

func (s *ImageProcessingService) Crop(ctx context.Context, in *image_service.CropRequest) (*httpbody.HttpBody, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImageProcessingService.Crop")
	defer span.Finish()

	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))
	img, Type, err := s.decodeImage(ctx, in.Image.Image)
	if err != nil {
		return nil, err
	}

	encoded, err := s.encodeImage(ctx, s.imagesProcessing.Crop(ctx, img,
		int(in.StartX), int(in.StartY), int(in.EndX), int(in.EndY)), Type)
	if err != nil {
		return nil, err
	}

	return &httpbody.HttpBody{
		ContentType: Type.String(),
		Data:        encoded,
	}, nil
}

func (s *ImageProcessingService) Resize(ctx context.Context,
	in *image_service.ResizeRequest) (*httpbody.HttpBody, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImageProcessingService.Resize")
	defer span.Finish()

	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	img, Type, err := s.decodeImage(ctx, in.Image.Image)
	if err != nil {
		return nil, err
	}

	encoded, err := s.encodeImage(ctx, s.imagesProcessing.Resize(ctx, img,
		int(in.Width), int(in.Height), in.ResampleFilter.String()), Type)
	if err != nil {
		return nil, err
	}

	return &httpbody.HttpBody{
		ContentType: Type.String(),
		Data:        encoded,
	}, nil
}

func (s *ImageProcessingService) Validate(ctx context.Context,
	in *image_service.ValidateRequest) (*image_service.ValidateResponce, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImageProcessingService.Validate")
	defer span.Finish()

	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))
	checked := false
	if len(in.Image.Image) == 0 {
		msg := "invalid image, received image has zero size"
		err = s.errorHandler.createErrorResponce(ErrImageTooSmall, msg)
		return &image_service.ValidateResponce{ImageValid: false, Details: &msg}, err
	}

	img, format, err := s.decodeImage(ctx, in.Image.Image)
	if err != nil {
		msg := "can't decode image"
		err = s.errorHandler.createErrorResponce(ErrInternal, msg)
		return &image_service.ValidateResponce{ImageValid: false, Details: &msg}, err
	}

	if len(in.SupportedTypes) != 0 {
		supported := false
		checked = true
		for _, Type := range in.SupportedTypes {
			if strings.EqualFold(Type, format.String()) {
				supported = true
				break
			}
		}

		if !supported {
			msg := fmt.Sprintf("image has unsupported type, supported types: %s, image has: %s",
				strings.Join(in.SupportedTypes, ","), format)
			return &image_service.ValidateResponce{ImageValid: true, Details: &msg}, nil
		}
	}

	if in.MaxHeight != nil {
		checked = true
		if img.Bounds().Dy() > int(*in.MaxHeight) {
			msg := fmt.Sprintf("image has height bigger than max height, image height: %d max height: %d",
				img.Bounds().Dy(), *in.MaxHeight)
			return &image_service.ValidateResponce{ImageValid: false, Details: &msg}, nil
		}
	}

	if in.MaxWidth != nil {
		checked = true
		if img.Bounds().Dx() > int(*in.MaxWidth) {
			msg := fmt.Sprintf("image has width bigger than max width, image width: %d max width: %d",
				img.Bounds().Dx(), *in.MaxWidth)
			return &image_service.ValidateResponce{ImageValid: false, Details: &msg}, nil
		}
	}

	if in.MinHeight != nil {
		checked = true
		if img.Bounds().Dy() < int(*in.MinHeight) {
			msg := fmt.Sprintf("image has height less than min height, image height: %d min height: %d",
				img.Bounds().Dy(), *in.MinHeight)
			return &image_service.ValidateResponce{ImageValid: false, Details: &msg}, nil
		}
	}

	if in.MinWidth != nil {
		checked = true
		if img.Bounds().Dx() < int(*in.MinWidth) {
			msg := fmt.Sprintf("image has width less than min width, image width: %d min width: %d",
				img.Bounds().Dx(), *in.MinWidth)
			return &image_service.ValidateResponce{ImageValid: false, Details: &msg}, nil
		}
	}

	if !checked {
		msg := "no received instructions for image validation"
		err = s.errorHandler.createErrorResponce(ErrNoInstructions, msg)
		return &image_service.ValidateResponce{ImageValid: true, Details: &msg}, err
	}

	return &image_service.ValidateResponce{ImageValid: true}, nil
}

func (s *ImageProcessingService) Hue(ctx context.Context, in *image_service.HueRequest) (*httpbody.HttpBody, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImageProcessingService.Hue")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	img, Type, err := s.decodeImage(ctx, in.Image.Image)
	if err != nil {
		return nil, err
	}

	encoded, err := s.encodeImage(ctx, s.imagesProcessing.Hue(ctx, img, int(in.Hue)), Type)
	if err != nil {
		return nil, err
	}

	return &httpbody.HttpBody{
		ContentType: Type.String(),
		Data:        encoded,
	}, nil
}

func (s *ImageProcessingService) Desaturate(ctx context.Context, in *image_service.Image) (*httpbody.HttpBody, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImageProcessingService.Desaturate")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	img, Type, err := s.decodeImage(ctx, in.Image)
	if err != nil {
		return nil, err
	}

	encoded, err := s.encodeImage(ctx, s.imagesProcessing.Desaturate(ctx, img), Type)
	if err != nil {
		return nil, err
	}

	return &httpbody.HttpBody{
		ContentType: Type.String(),
		Data:        encoded,
	}, nil
}

func (s *ImageProcessingService) Blur(ctx context.Context, in *image_service.BlurRequest) (*httpbody.HttpBody, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImageProcessingService.Blur")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	img, Type, err := s.decodeImage(ctx, in.Image.Image)
	if err != nil {
		return nil, err
	}

	encoded, err := s.encodeImage(ctx, s.imagesProcessing.Blur(ctx, img, in.BlurRadius, in.Method.String()), Type)
	if err != nil {
		return nil, err
	}

	return &httpbody.HttpBody{
		ContentType: Type.String(),
		Data:        encoded,
	}, nil
}

func (s *ImageProcessingService) decodeImage(ctx context.Context, img []byte) (image.Image, *mimetype.MIME, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ImageProcessingService.decodeImage")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	decoded, Type, err := image_processing.DecodeImage(img)
	if err != nil {
		err = s.errorHandler.createErrorResponce(ErrInvalidArgument, err.Error())
		return nil, nil, err
	}

	return decoded, Type, nil
}

func (s *ImageProcessingService) encodeImage(ctx context.Context, img image.Image, Type *mimetype.MIME) ([]byte, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ImageProcessingService.encodeImage")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	encoded, err := image_processing.EncodeImage(img, Type.Extension())
	if err != nil {
		err = s.errorHandler.createErrorResponce(ErrInvalidArgument, err.Error())
		return []byte{}, err
	}

	return encoded, nil
}
