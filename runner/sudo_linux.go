package runner

import "os/exec"

var queries = []string{
	"gksu",
	"gksudo",
	"kdesu",
	"kdesudo",
}

func withSudo(args []string) []string {
	var sudo = "sudo"
	for _, name := range queries {
		if _, err := exec.LookPath(name); err == nil {
			sudo = name
			break
		}
	}
	var ans []string
	ans = append(ans, sudo)
	ans = append(ans, args...)
	return ans
}
