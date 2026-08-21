package sanitizer

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/bogem/id3v2/v2"
)

const maxSize = 10 * 1024 * 1024

func SanitizeFile(file multipart.File) ([]byte, error) {
	limitedReader := io.LimitReader(file, maxSize+1)

	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, err
	}

	if len(data) > maxSize {
		return nil, fmt.Errorf("file exceeds maximum size of 10 MB")
	}

	contentType := http.DetectContentType(data)

	switch contentType {
	case "image/jpeg":
		return SanitizeJPEG(data)
	case "image/png":
		return SanitizePNG(data)
	case "audio/mpeg":
		return SanitizeMP3(data)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", contentType)
	}
}

func SanitizeJPEG(data []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode jpeg image: %w", err)
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85}); err != nil {
		return nil, fmt.Errorf("failed to re-encode jpeg image: %w", err)
	}

	return buf.Bytes(), nil
}

func SanitizePNG(data []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode png image: %w", err)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to re-encode png image: %w", err)
	}

	return buf.Bytes(), nil
}

// Manually strip ID3V1, use lib to strip ID3V2
func SanitizeMP3(data []byte) ([]byte, error) {
	// Excise ID3V1 first, stored in 128 bytes at the end, ensure that payload exists
	if len(data) >= 128 && string(data[len(data)-128:len(data)-125]) == "TAG" {
		data = data[:len(data)-128]
	}

	// Excise ID3V2, variable & located at the start of the file, if present
	tag, err := id3v2.ParseReader(bytes.NewReader(data), id3v2.Options{Parse: true})
	if err != nil {
		return nil, err
	}

	id3v2size := tag.Size()

	data = data[id3v2size:]

	return data, nil
}
