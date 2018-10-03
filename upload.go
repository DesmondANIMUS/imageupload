package imageupload

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/nfnt/resize"
)

var extMap map[string]int

// UploadFile function is a simple helper function that uploads and saves an image on the server
// Params:
// r: to get the picture from multi-part form using key "get_picture"
// location: the path on server where you wish to save file. Ex: /users/images/
// ID: unique string ID for the image
// size: to resize the image, the function will keep the aspect ratio intact
func UploadFile(r *http.Request, location string, ID string, size uint) (string, error) {
	initExtMap()
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
func saveFile(src multipart.File, location, id, ext string, size uint) (string, error) {
	name := id + ".jpg"
	path := "." + location + name
	var img image.Image
	var op jpeg.Options
	op.Quality = 50

	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	switch extMap[ext] {
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
	default:
		return "", errors.New("File is not an image")
	}

	if err := jpeg.Encode(dst, img, &op); err != nil {
		return "", err
	}

	path = location + name
	return path, err
}

// DecodeJPG function decodes JPG image
func decodeJPG(src multipart.File, size uint) (image.Image, error) {
	img, err := jpeg.Decode(src)
	img = resize.Resize(size, 0, img, resize.Lanczos3)
	return img, err
}

// DecodePNG function decodes PNG image
func decodePNG(src multipart.File, size uint) (image.Image, error) {
	img, err := png.Decode(src)
	img = resize.Resize(size, 0, img, resize.Lanczos3)
	return img, err
}

// DecodeGIF function decodes GIF image
func decodeGIF(src multipart.File, size uint) (image.Image, error) {
	img, err := gif.Decode(src)
	img = resize.Resize(size, 0, img, resize.Lanczos3)
	return img, err
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
