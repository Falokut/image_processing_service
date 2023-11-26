package image_processing

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"

	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
)

func EncodeImage(img image.Image, Extension string) ([]byte, error) {
	buf := new(bytes.Buffer)
	var err error

	format, err := imaging.FormatFromExtension(Extension)
	if err != nil {
		return []byte{}, err
	}
	err = imaging.Encode(buf, img, format, imaging.PNGCompressionLevel(png.BestSpeed), imaging.JPEGQuality(70))

	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func DecodeImage(img []byte) (decoded image.Image, mimeType *mimetype.MIME, err error) {
	mimeType = mimetype.Detect(img)
	buf := bytes.NewBuffer(img)
	decoded, err = imaging.Decode(buf)
	if err != nil {
		return nil, mimeType, err
	}
	return decoded, mimeType, nil
}

func ConvertToBase64(img []byte) []byte {
	base64Encoded := make([]byte, base64.StdEncoding.EncodedLen(len(img)))
	base64.StdEncoding.Encode(base64Encoded, img)
	return base64Encoded
}

func GetMimeTypeExt(img []byte) string {
	return mimetype.Detect(img).Extension()
}
