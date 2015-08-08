package processing

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Identify(url string) ([]string, float64) {
	// download it, save to tmp
	fileName := download(url)

	// process it
	out, err := exec.Command("./OverFeat/src/overfeat", "-n 1", fileName).CombinedOutput()

	if err != nil {
		log.Println(out, err)
	}

	split := strings.Split(string(out), " ")
	confidence, _ := strconv.ParseFloat(strings.Trim(split[len(split)-1], "\n"), 64)
	classifier := strings.Split(strings.Trim(strings.Join(split[:len(split)-1], " "), ","), ", ")

	// delete from /tmp

	return classifier, confidence
}

func download(url string) string {
	tokens := strings.Split(url, "/")
	fileName := "/tmp/" + tokens[len(tokens)-1]

	output, err := os.Create(fileName)

	if err != nil {
		log.Println("Error while creating", fileName, " - ", err)
	}

	defer output.Close()

	response, err := http.Get(url)

	if err != nil {
		log.Println("Error while downloading", url, " - ", err)
	}

	defer response.Body.Close()

	_, er := io.Copy(output, response.Body)

	if er != nil {
		log.Println("Error while saving", url, "-", er)
	}

	return fileName
}
