package nat

import (
	"os/exec"
	"strings"
)

func getGateway() (string, error) {
	out, err := exec.Command("bash", "-c", "netstat -rn | grep default | awk '{print $2}'").Output()
	if err != nil {
		return "", err
	}
	return strings.Split(string(out), "\n")[0], nil
}
