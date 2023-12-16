package default_processing

import (
	"context"
	"image"
	"strings"

	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/transform"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type ImageProcessing struct {
	logger *logrus.Logger
}

func NewImageProcessingService(logger *logrus.Logger) *ImageProcessing {
	return &ImageProcessing{
		logger: logger,
	}
}

func (p ImageProcessing) Resize(ctx context.Context, img image.Image, width, height int, resampleMethod string) image.Image {
	span, _ := opentracing.StartSpanFromContext(ctx, "ImageProcessing.Resize")
	defer span.Finish()
	p.logger.Info("Resizing started")

	return transform.Resize(img, width, height, resolveFiter(resampleMethod))
}

func (p ImageProcessing) Crop(ctx context.Context, img image.Image, x0, y0, x1, y1 int) image.Image {
	span, _ := opentracing.StartSpanFromContext(ctx, "ImageProcessing.Crop")
	defer span.Finish()
	return transform.Crop(img, image.Rectangle{
		Min: image.Point{X: x0, Y: y0},
		Max: image.Point{X: x1, Y: y1}})
}

func (p ImageProcessing) Desaturate(ctx context.Context, img image.Image) image.Image {
	span, _ := opentracing.StartSpanFromContext(ctx, "ImageProcessing.Desaturate")
	defer span.Finish()

	return effect.Grayscale(img)
}

func (p ImageProcessing) Hue(ctx context.Context, img image.Image, hue int) image.Image {
	span, _ := opentracing.StartSpanFromContext(ctx, "ImageProcessing.Hue")
	defer span.Finish()
	hue %= 360
	return adjust.Hue(img, hue)
}

type BlurMethod string

const (
	Box      BlurMethod = "Box"
	Gaussian BlurMethod = "Gaussian"
)

func (p ImageProcessing) Blur(ctx context.Context, img image.Image, radius float64, method string) image.Image {
	span, _ := opentracing.StartSpanFromContext(ctx, "ImageProcessing.Hue")
	defer span.Finish()
	p.logger.Info("Blur started")

	blurMethod := resolveBlurMethod(method)
	switch blurMethod {
	default:
		return blur.Box(img, radius)
	case Gaussian:
		return blur.Gaussian(img, radius)
	}
}

func resolveFiter(filter string) transform.ResampleFilter {
	switch filter {
	case "Lanczos":
		return transform.Lanczos
	case "CatmullRom":
		return transform.CatmullRom
	case "MitchellNetravali":
		return transform.MitchellNetravali
	case "Linear":
		return transform.Linear
	case "Box":
		return transform.Box
	case "NearestNeighbor":
		return transform.NearestNeighbor
	default:
		return transform.NearestNeighbor
	}
}

func resolveBlurMethod(method string) BlurMethod {
	method = strings.ToUpper(method)
	switch method {
	case "BOX":
		return Box
	case "GAUSSIAN":
		return Gaussian
	default:
		return Box
	}
}
