package image_processing

import (
	"context"
	"errors"
	"image"
)

var (
	ErrUnsupported = errors.New("unsupported")
	ErrInternal    = errors.New("internal")
)

type ImagesProcessing interface {
	Resize(ctx context.Context, img image.Image, width, height int, resampleMethod string) image.Image
	Crop(ctx context.Context, img image.Image, x0, y0, x1, y1 int) image.Image
	Desaturate(ctx context.Context, img image.Image) image.Image
	Hue(ctx context.Context, img image.Image, hue int) image.Image
	Blur(ctx context.Context, img image.Image, radius float64, method string) image.Image
}
