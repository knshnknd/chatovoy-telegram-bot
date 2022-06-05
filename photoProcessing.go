package main

import (
	"fmt"
	"math/rand"
)

func generatePhotoName() string {
	return fmt.Sprintf("kuzya%d", rand.Intn(numberOfKuzyasPictures))
}

func makePhotoPath(photoName string) string {
	return fmt.Sprintf("resources/%s.jpg", photoName)
}
