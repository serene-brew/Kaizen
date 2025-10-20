package src

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"strconv"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	//	"github.com/mattn/go-sixel"
	"golang.org/x/image/draw"
)

func main() {
	if len(os.Args) < 1 {
		return
	}

	imageURL := "https://s4.anilist.co/file/anilistcdn/media/anime/cover/large/bx127230-DdP4vAdssLoz.png"
	var targetWidth, targetHeight int
	var err error

	if len(os.Args) >= 2 {
		targetWidth, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid width\n")
			return
		}
	}

	if len(os.Args) >= 3 {
		targetHeight, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid height\n")
			return
		}
	}

	resp, err := http.Get(imageURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching image: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Error: received status code %d\n", resp.StatusCode)
		return
	}

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading image: %v\n", err)
		return
	}

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding image: %v\n", err)
		return
	}

	if targetWidth > 0 || targetHeight > 0 {
		img = sixelResizeImage(img, targetWidth, targetHeight)
	}

	if err := displaySixelImage(img); err != nil {
		fmt.Fprintf(os.Stderr, "Error displaying image: %v\n", err)
		return
	}
}

func sixelResizeImage(img image.Image, width, height int) image.Image {
	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	if width > 0 && height == 0 {
		aspectRatio := float64(origHeight) / float64(origWidth)
		height = int(float64(width) * aspectRatio)
	} else if height > 0 && width == 0 {
		aspectRatio := float64(origWidth) / float64(origHeight)
		width = int(float64(height) * aspectRatio)
	}

	if width <= 0 || height <= 0 {
		return img
	}

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	return dst
}

func displaySixelImage(img image.Image) error {
	enc := NewEncoder(os.Stdout)

	if err := enc.Encode(img); err != nil {
		return fmt.Errorf("failed to encode sixel: %w", err)
	}

	return nil
}
