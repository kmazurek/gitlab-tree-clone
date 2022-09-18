package infra

import (
	"log"
	"os/exec"
	"strings"
)

func CloneProject(path string, projectName string, url string) error {
	log.Println("Cloning project: " + path + "/" + projectName)

	out, err := exec.Command("git", "clone", url, path+"/"+projectName).CombinedOutput()
	if err != nil {
		output := string(out)
		if strings.Contains(output, "already exists and is not an empty directory") {
			return nil
		}

		log.Fatal(string(out))
		return err
	}

	return nil
}
