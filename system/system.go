package system

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	. "github.com/xpetit/x/v4"
)

type Process struct {
	Name      string
	PID       int
	PPID      int
	Mem       int
	CPUTime   int // in clock ticks, divide by Tick
	StartTime int // in clock ticks, divide by Tick
}

var Tick = float64(atoi(strings.TrimSpace(string(Must(exec.Command("getconf", "CLK_TCK").Output())))))

func atoi(s string) int { return Must(strconv.Atoi(s)) }

func Status() (processes []*Process, uptime float64) {
	{
		s, _, _ := strings.Cut(string(Must(os.ReadFile("/proc/uptime"))), " ")
		uptime = Must(strconv.ParseFloat(s, 64))
	}
	proc := Must(os.Open("/proc"))
	defer Closing(proc)

	for _, dir := range Must(proc.Readdirnames(-1)) {
		if dir[0] < '0' || dir[0] > '9' {
			continue
		}
		PID := dir
		p := Process{PID: atoi(dir)}
		{ // read process name
			b, err := os.ReadFile("/proc/" + PID + "/comm")
			if err != nil {
				continue
			}
			comm := strings.TrimSpace(string(b))
			exe, err := filepath.EvalSymlinks("/proc/" + PID + "/exe")
			if err != nil {
				p.Name = comm
			} else {
				exe = filepath.Base(exe)
				if exe != comm && !strings.HasPrefix(strings.ToLower(exe), strings.ToLower(comm)) {
					exe += " (" + comm + ")"
				}
				p.Name = exe
			}
		}
		{ // read process CPU values
			b, err := os.ReadFile("/proc/" + PID + "/stat")
			if err != nil {
				continue
			}
			s := string(b)
			_, s, _ = strings.Cut(s, ") ") // skips "pid (comm) "
			stats := strings.Split(s, " ")
			p.PPID = atoi(stats[1])       // ppid
			p.CPUTime += atoi(stats[11])  // utime
			p.CPUTime += atoi(stats[12])  // stime
			p.StartTime = atoi(stats[19]) // starttime
		}
		{ // read process memory values
			b, err := os.ReadFile("/proc/" + PID + "/status")
			if err != nil {
				continue
			}
			if _, rss, ok := strings.Cut(string(b), "RssAnon:\t"); ok {
				rss = strings.TrimLeft(rss, " ")
				rss, _, _ = strings.Cut(rss, " ")
				p.Mem = atoi(rss) * 1024
			}
		}
		processes = append(processes, &p)
	}
	return
}
