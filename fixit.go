package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const inputDir = "IMFDB_FINAL"
const outputDir = "IMFDB_FIXED"

const dsStore = ".DS_Store"

func main() {
	people, err := readDirNames(inputDir)
	if err != nil {
		panic(err)
	}

	for _, person := range people {
		if person == dsStore {
			continue
		}

		handlePerson(person)
	}
}

func handlePerson(person string) {
	movies, err := readDirNames(filepath.Join(inputDir, person))
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(filepath.Join(outputDir, person), os.ModeDir|0666)
	if err != nil {
		panic(err)
	}

	personAttributesFile, err := os.OpenFile(filepath.Join(outputDir, person, "attributes.txt"), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer personAttributesFile.Close()

	var id int

	for _, movie := range movies {
		if movie == dsStore {
			continue
		}

		handleMovie(person, movie, personAttributesFile, &id)
	}
}

func handleMovie(person, movie string, personAttributesFile *os.File, id *int) {
	movieFiles, err := readDirNames(filepath.Join(inputDir, person, movie))
	if err != nil {
		panic(err)
	}

	var movieAttributesFile *string
	for _, movieFile := range movieFiles {
		if strings.HasSuffix(movieFile, ".txt") {
			movieAttributesFile = &movieFile
		}
	}

	if movieAttributesFile == nil {
		fmt.Printf("> Ignoring %s/%s because it has no attributes file\n", person, movie)
	}

	images, err := readDirNames(filepath.Join(inputDir, person, movie, "images"))
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		err = copyFile(filepath.Join(inputDir, person, movie, "images", image), filepath.Join(outputDir, person, strconv.Itoa(*id)+".jpg"))
		(*id)++
	}
}

func readDirNames(dirPath string) ([]string, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	return dir.Readdirnames(0)
}

func copyFile(src, dest string) error {
	fromFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	toFile, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	_, err = io.Copy(toFile, fromFile)
	if err != nil {
		return err
	}

	return nil
}
