package src

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"sync"

	_ "image/gif"
	_ "image/jpeg"

	"golang.org/x/image/draw"
)

// imageCache stores processed image sequences by URL and dimensions
var imageCache = struct {
	sync.RWMutex
	cache map[string]string
}{
	cache: make(map[string]string),
}

// getCacheKey creates a unique key for the image cache
func getCacheKey(url string, width, height int) string {
	return fmt.Sprintf("%s_%dx%d", url, width, height)
}

// RenderKittyImageFromURL fetches the provided image URL, resizes it to the
// requested pixel dimensions (always 200x300 in our usage), encodes it as PNG
// and returns the Kitty graphics protocol escape sequence as a string. The
// caller can render that string directly into the terminal output.
func RenderKittyImageFromURL(url string, width, height int) (string, error) {
	if url == "" {
		return "", nil
	}

	// Check cache first
	cacheKey := getCacheKey(url, width, height)
	imageCache.RLock()
	if cached, ok := imageCache.cache[cacheKey]; ok {
		imageCache.RUnlock()
		return cached, nil
	}
	imageCache.RUnlock()

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status fetching image: %s", resp.Status)
	}

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Always resize to provided width/height (200x300 in our usage)
	img = resizeImage(img, width, height)

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", fmt.Errorf("failed to encode png: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	sequence := buildKittySequence(encoded)

	// Cache the processed image sequence
	imageCache.Lock()
	imageCache.cache[cacheKey] = sequence
	imageCache.Unlock()

	return sequence, nil
}

func resizeImage(img image.Image, width, height int) image.Image {
	bounds := img.Bounds()
	if width <= 0 || height <= 0 {
		return img
	}

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)
	return dst
}

// buildKittySequence builds the kitty graphics escape sequence and returns it
// as a string. We return the full sequence (possibly split into chunks) so
// callers can place it inside their TUI output rather than printing it.
func buildKittySequence(base64Data string) string {
	chunkSize := 4096
	dataLen := len(base64Data)
	var out bytes.Buffer

	// Save cursor position
	out.WriteString("\x1b[s")

	// Delete all images at the current cursor position
	out.WriteString("\x1b_Ga=d\x1b\\")

	for i := 0; i < dataLen; i += chunkSize {
		end := i + chunkSize
		if end > dataLen {
			end = dataLen
		}

		chunk := base64Data[i:end]

		if i == 0 {
			// Place new image, using:
			// f=100: PNG format
			// a=T: transmit data as base64
			// z=1: put in background
			// C=1: allow cursor movement
			if end >= dataLen {
				out.WriteString(fmt.Sprintf("\x1b_Gf=100,a=T,z=1,C=1;%s\x1b\\", chunk))
			} else {
				out.WriteString(fmt.Sprintf("\x1b_Gf=100,a=T,z=1,C=1,m=1;%s\x1b\\", chunk))
			}
		} else if end >= dataLen {
			out.WriteString(fmt.Sprintf("\x1b_Gm=0;%s\x1b\\", chunk))
		} else {
			out.WriteString(fmt.Sprintf("\x1b_Gm=1;%s\x1b\\", chunk))
		}
	}

	// Restore cursor position
	out.WriteString("\x1b[u")

	return out.String()
}
