package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func FileCreatesIfNotExists(filePath string) error {
	if !FileExists(filePath) {
		if err := os.MkdirAll(path.Dir(filePath), 0755); err != nil {
			return err
		}
		fp, err := os.OpenFile(filePath, os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer fp.Close()
	}
	return nil
}

type WithInputCheck func(text string) error

func EmptyNotAllowed() WithInputCheck {
	return func(text string) error {
		if strings.TrimSpace(text) == "" {
			return errors.New("empty text")
		}
		return nil
	}
}

func PositiveIntegerOnly() WithInputCheck {
	return func(text string) error {
		number, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return err
		}
		if number < 0 {
			return errors.New("negative number")
		}
		return nil
	}
}

func Input(prompt string, inputChecks ...WithInputCheck) string {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(prompt)
		if reader.Scan() {
			text := strings.TrimSpace(reader.Text())
			isValid := true
			for _, inputCheck := range inputChecks {
				if err := inputCheck(text); err != nil {
					isValid = false
					fmt.Printf("Error: %s\n", err)
					break
				}
			}
			if isValid {
				return text
			}
		}
	}
}
