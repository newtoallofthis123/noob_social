package utils

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"strings"

	"github.com/anthonynsimon/bild/transform"
	"github.com/o1egl/govatar"
)

// FILEPATH is the path to the profile pictures
const FILEPATH string = "./static/assets/images/"

// STATICPATH is the path to the static folder from the web server
const STATICPATH string = "/static/assets/images/"

// Uses the govatar library to generate an avatar
// for the user
// Since the gender of the user is not known, we
// use MALE as a placeholder default
func GetAvatar(username string) (bytes.Buffer, error) {
	img, err := govatar.GenerateForUsername(govatar.MALE, username)
	if err != nil {
		return bytes.Buffer{}, err
	}

	var imgBuff bytes.Buffer

	if png.Encode(&imgBuff, img) != nil {
		return bytes.Buffer{}, err
	}

	return imgBuff, nil
}

func CheckPicture(name string, resize bool) (string, error) {
	// read from /static/assets/profile_pics
	file, err := os.Open(FILEPATH + name)
	if err != nil {
		return "", err
	}

	imageBuff := bytes.Buffer{}

	// convert file to image.Image
	fileImage, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	resizedImage := fileImage
	if resize {
		resizedImage = transform.Resize(fileImage, 128, 128, transform.Linear)
	}
	if png.Encode(&imageBuff, resizedImage) != nil {
		return "", err
	}

	newName := GenerateRandomString(8) + ".png"

	// delete file
	err = os.Remove(FILEPATH + name)
	if err != nil {
		return "", err
	}

	// write imageBuff to file
	err = os.WriteFile(FILEPATH+newName, imageBuff.Bytes(), 0644)
	if err != nil {
		return "", err
	}

	return newName, nil
}

func GetImage(name string) (bytes.Buffer, error) {
	file, err := os.Open(FILEPATH + name)
	if err != nil {
		return bytes.Buffer{}, err
	}

	imageBuff := bytes.Buffer{}

	// convert file to image.Image
	fileImage, _, err := image.Decode(file)
	if err != nil {
		return bytes.Buffer{}, err
	}

	if png.Encode(&imageBuff, fileImage) != nil {
		return bytes.Buffer{}, err
	}

	return imageBuff, nil
}

func GetImageFile(name string) (image.Image, error) {
	file, err := os.Open(FILEPATH + name)
	if err != nil {
		return nil, err
	}

	// convert file to image.Image
	fileImage, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return fileImage, nil
}
func DeleteUnused(used []string) error {
	files, err := os.ReadDir(FILEPATH)
	if err != nil {
		return err
	}

	for _, file := range files {
		// check if file is in used
		// if not, delete it
		if !strings.Contains(strings.Join(used, " "), file.Name()) {
			err = os.Remove(FILEPATH + file.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
