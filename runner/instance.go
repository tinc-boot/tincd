package runner

import (
	"bufio"
	"context"
	"github.com/tinc-boot/tincd/utils"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var (
	addSubnetPattern = regexp.MustCompile(`ADD_SUBNET\s+from\s+([^\s]+)\s+\(([^\s]+)\s+port\s+(\d+)\)\:\s+\d+\s+[\w\d]+\s+([^\s]+)\s+([^#]+)`)
	delSubnetPattern = regexp.MustCompile(`DEL_SUBNET\s+[^:]+:\s+\d+\s+[\w\d]+\s+([^\s]+)\s+([^#]+)`)
)

//Sending DEL_SUBNET to everyone (BROADCAST): 11 3f17d1ce hubreddecnet_PEN005 6e:6a:5e:26:39:d2#10
func fromLine(line string) *SubnetEvent {
	if match := addSubnetPattern.FindAllStringSubmatch(line, -1); len(match) > 0 {
		groups := match[0]
		if len(groups) != 6 {
			return nil
		}
		var event SubnetEvent
		event.Add = true
		event.Advertising.Node = groups[1]
		event.Advertising.Host = groups[2]
		event.Advertising.Port = groups[3]
		event.Peer.Node = groups[4]
		event.Peer.Subnet = groups[5]
		return &event
	} else if match := delSubnetPattern.FindAllStringSubmatch(line, -1); len(match) > 0 {
		groups := match[0]
		if len(groups) != 3 {
			return nil
		}
		var event SubnetEvent
		event.Add = false
		event.Peer.Node = groups[1]
		event.Peer.Subnet = groups[2]
		return &event
	}
	return nil
}

type SubnetEvent struct {
	Add         bool
	Advertising struct {
		Node string
		Host string
		Port string
	}
	Peer struct {
		Node   string
		Subnet string
	}
}

func makeArgs(tincBin string, dir string) []string {
	return []string{tincBin, "-D", "-d", "-d", "-d", "-d",
		"--pidfile", filepath.Join(dir, "pid.run"),
		"-c", dir}
}

// Run tinc application and scan output for events
func RunTinc(global context.Context, askSudo bool, tincBin string, dir string) <-chan SubnetEvent {

	var events = make(chan SubnetEvent)

	reader, writer := io.Pipe()
	scanner := bufio.NewScanner(reader)
	args := makeArgs(tincBin, dir)
	if askSudo {
		args = withSudo(args)
	}

	logfile, err := os.Create(filepath.Join(dir, "log.txt"))
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stderr = io.MultiWriter(writer, logfile)
	utils.SetCmdAttrs(cmd)
	cmd.Stdout = io.MultiWriter(writer, logfile)

	child, cancel := context.WithCancel(global)
	go func() {
		// kill process when context canceled
		defer cancel()
		<-child.Done()
		killProcess(cmd)
	}()

	go func() {
		// run process, cancel context after
		defer writer.Close()
		defer logfile.Close()
		defer cancel()
		err := cmd.Run()
		if err != nil {
			log.Println("run tincd:", err)
		}

	}()

	go func() {
		// read events from stdout, stderr
		defer close(events)
		defer cancel()
		for scanner.Scan() {
			if event := fromLine(scanner.Text()); event != nil {
				select {
				case events <- *event:
				case <-global.Done():
					return
				}
			}
		}
	}()

	return events
}
