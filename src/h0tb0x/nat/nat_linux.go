package nat

import (
	"os/exec"
	"strings"
)

func getGateway() (string, error) {
	out, err := exec.Command("bash", "-c", "ip route | grep ^default | awk '{printf $3}'").Output()
	if err != nil {
		return "", err
	}
	return strings.Split(string(out), "\n")[0], nil
}
