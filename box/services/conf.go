package services

import (
	"bufio"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"strings"
)

func UpdateLitConfField(filePath, field, newValue string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, field+"=") {
			lines = append(lines, field+"="+"LNT."+newValue)
			found = true
		} else {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if !found {
		lines = append(lines, field+"="+newValue)
	}

	output := strings.Join(lines, "\n")
	err = os.WriteFile(filePath, []byte(output), 0644)
	if err != nil {
		return errors.Wrap(err, "os.WriteFile")
	}

	Cmd := exec.Command("systemctl", "restart", "litd")
	_, err = Cmd.Output()
	if err != nil {
		return errors.Wrap(err, "exec.Command systemctl restart litd")
	}

	return nil
}
