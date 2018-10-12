package imageupload

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"

	"github.com/nfnt/resize"
)

var extMap map[string]int

var fs fileSystem = osFS{}

var ErrFileNotSupported = errors.New("file is not an image")

func init() {
	initExtMap()
}

// UploadFile function is a simple helper function that uploads and saves an image on the server
// Params:
// r: to get the picture from multi-part form using key "get_picture"
// location: the path on server where you wish to save file. Ex: /users/images/
// ID: unique string ID for the image
// size: to resize the image, the function will keep the aspect ratio intact
func UploadFile(r *http.Request, location string, ID string, size uint) (string, error) {
	var path string
	file, hdr, err := r.FormFile("get_picture")
	if err != nil {
		return path, nil
	}

	ext := getExt(hdr.Filename)
	defer file.Close()

	path, err = saveFile(file, location, ID, ext, size)
	if err != nil {
		return "", err
	}

	return path, nil
}

// SaveFile function helps in uploading the profile picture of user
func saveFile(src io.Reader, location, ID, ext string, size uint) (string, error) {
	name := ID + ".jpg"
	path := "." + location + name
	var img image.Image
	var op jpeg.Options
	var err error
	op.Quality = 50

	e, ok := extMap[ext]
	if !ok {
		return "", ErrFileNotSupported
	}

	switch e {
	case JPG:
		img, err = decodeJPG(src, size)
		if err != nil {
			return "", err
		}
	case PNG:
		img, err = decodePNG(src, size)
		if err != nil {
			return "", err
		}
	case GIF:
		img, err = decodeGIF(src, size)
		if err != nil {
			return "", err
		}
	case BMP:
		img, err = DecodeBMP(src, size)
		if err != nil {
			return "", err
		}
	case TIFF:
		img, err = DecodeTIFF(src, size)
		if err != nil {
			return "", err
		}
	}

	dst, err := fs.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if err := jpeg.Encode(dst, img, &op); err != nil {
		return "", err
	}

	path = location + name
	return path, err
}

// DecodeTIFF function decodes TIFF image
func DecodeTIFF(src io.Reader, size uint) (image.Image, error) {
	img, err := tiff.Decode(src)
	if err != nil {
		return nil, err
	}
	return resize.Resize(size, 0, img, resize.Lanczos3), nil
}

// DecodeBMP function decodes BMP image
func DecodeBMP(src io.Reader, size uint) (image.Image, error) {
	img, err := bmp.Decode(src)
	if err != nil {
		return nil, err
	}
	return resize.Resize(size, 0, img, resize.Lanczos3), nil
}

// DecodeJPG function decodes JPG image
func decodeJPG(src io.Reader, size uint) (image.Image, error) {
	img, err := jpeg.Decode(src)
	if err != nil {
		return nil, err
	}
	return resize.Resize(size, 0, img, resize.Lanczos3), nil
}

// DecodePNG function decodes PNG image
func decodePNG(src io.Reader, size uint) (image.Image, error) {
	img, err := png.Decode(src)
	if err != nil {
		return nil, err
	}
	return resize.Resize(size, 0, img, resize.Lanczos3), nil
}

// DecodeGIF function decodes GIF image
func decodeGIF(src io.Reader, size uint) (image.Image, error) {
	img, err := gif.Decode(src)
	if err != nil {
		return nil, err
	}
	return resize.Resize(size, 0, img, resize.Lanczos3), nil
}

func initExtMap() {
	extMap = make(map[string]int)

	extMap["jpeg"] = JPG
	extMap["jpg"] = JPG
	extMap["JPG"] = JPG
	extMap["JPEG"] = JPG

	extMap["png"] = PNG
	extMap["PNG"] = PNG

	extMap["gif"] = GIF
	extMap["GIF"] = GIF

	extMap["BMP"] = BMP
	extMap["bmp"] = BMP

	extMap["TIFF"] = TIFF
	extMap["TIF"] = TIFF
	extMap["tiff"] = TIFF
	extMap["tif"] = TIFF
}

func getExt(filename string) string {
	revName := reverse(filename)
	revExt := strings.Split(revName, ".")[0]
	return reverse(revExt)
}

func reverse(txt string) string {
	data := []rune(txt)
	var result []rune

	for i := len(data) - 1; i >= 0; i-- {
		result = append(result, data[i])
	}

	return string(result)
}

type fileSystem interface {
	Create(name string) (file, error)
}

type file interface {
	io.Closer
	io.Writer
}

// osFS implements fileSystem using the local disk.
type osFS struct{}

func (osFS) Create(name string) (file, error) { return os.Create(name) }
