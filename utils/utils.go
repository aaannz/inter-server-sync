package utils

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

//ReverseArray reverses the array
func ReverseArray(slice interface{}) {
	size := reflect.ValueOf(slice).Len()
	swap := reflect.Swapper(slice)
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

// Contains is a helper method to check if a string element exist in the string slice
func Contains(slice []string, elementToFind string) bool {
	for _, element := range slice {
		if elementToFind == element {
			return true
		}
	}
	return false
}

func GetAbsPath(path string) string {
	result := path
	if filepath.IsAbs(path) {
		result, _ = filepath.Abs(path)
	} else {
		homedir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Msg("Couldn't determine the home directory")
		}
		if strings.HasPrefix(path, "~") {
			result = strings.Replace(path, "~", homedir, -1)
		}
	}
	return result
}

func GetCurrentServerVersion() (string, string) {
	var version, product string
	rhndefault := "/etc/rhn/rhn.conf"
	webpath := "/usr/share/rhn/config-defaults/rhn_web.conf"
	altpath := "/usr/share/rhn/config-defaults/rhn.conf"

	files := []string{rhndefault, webpath, altpath}
	property := []string{"product_name", "web.product_name"}
	product = "SUSE Manager"
	for _, path:= range files {
		for _, search := range property {
			p, err := ScannerFunc(path, search)
			if err == nil {
				product = p
				break
			} else if err != nil {
				continue
			}
		}
	}

	if product != "SUSE Manager" {
		product = "uyuni"
		v, err := ScannerFunc(rhndefault, "web.version.uyuni")
		if err != nil {
			v, err = ScannerFunc(webpath, "web.version.uyuni")
			if err != nil {
				log.Fatal().Msgf("No version found for product %s", product)
			}
		}
		version = v
	} else if product == "SUSE Manager" {
		v, err := ScannerFunc(rhndefault, "web.version")
		if err != nil {
			v, err = ScannerFunc(webpath, "web.version")
		}
		if err != nil {
			log.Fatal().Msgf("No version found for product %s", product)
		}
		version = v
	}
	return version, product
}

func ScannerFunc(path string, search string) (string, error) {
	var output string
	f, err := os.Open(path)
	if err != nil {
		log.Fatal().Msg("Couldn't open file")
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), search) {
			splits := strings.Split(scanner.Text(), "= ")
			output = splits[1]
			return output, nil
		}
	}
	return "", fmt.Errorf("String not found!")
}
