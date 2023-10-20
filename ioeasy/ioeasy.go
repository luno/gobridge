package ioeasy

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

func CreateFileFromPath(path string) error {
	pCon := strings.Split(path, "/")
	reconPath := "."
	for index, item := range pCon {
		if item == "." {
			continue
		}

		reconPath += "/" + item

		if index == (len(pCon) - 1) {
			err := CreateFileIfNotExists(reconPath)
			if err != nil {
				return err
			}
		} else {
			err := CreateDirIfNotExists(reconPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CreateFileIfNotExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = ioutil.WriteFile(path, []byte{}, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func CreateDirIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return errors.New("unable to create directory")
		}
	}

	return nil
}
