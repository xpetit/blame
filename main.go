package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/xpetit/blame/system"
	. "github.com/xpetit/x/v2"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func atoi(s string) int { return C2(strconv.Atoi(s)) }

func lines(command ...string) []string {
	return strings.Split(
		strings.TrimSpace(string(C2(Output(exec.Command(command[0], command[1:]...))))),
		"\n",
	)
}

func main() {
	var minCPU float64
	var minMem int
	{
		const threshold = 0.001 // 0.1 %
		meminfo := string(C2(os.ReadFile("/proc/meminfo")))
		totalMem := float64(1024 * atoi(regexp.MustCompile(`MemTotal: *(\d+) kB`).FindStringSubmatch(meminfo)[1]))
		defaultMem := int(math.Round(totalMem * threshold))

		flag.Float64Var(&minCPU, "min-cpu", 100*threshold, "Minimum CPU percentage")
		flag.IntVar(&minMem, "min-mem", defaultMem, fmt.Sprintf("Minimum memory bytes (%s)", FormatByte(defaultMem)))
		flag.Parse()
	}

	programs := map[int]string{}
	{
		for _, line := range lines("wmctrl", "-lpx") {
			fields := strings.Fields(line)

			if desktopID := fields[1]; desktopID == "-1" {
				continue
			}

			pidString := fields[2]
			if pidString == "0" {
				log.Println("PID not reported by", fields[3], fields[5:])
				continue
			}

			pid := atoi(pidString)
			class := fields[3]
			if p, ok := programs[pid]; ok {
				Assert(p == class)
			} else {
				programs[pid] = class
			}
		}
	}

	processes, uptime := system.Status()
	pids := map[int]*system.Process{}
	for _, p := range processes {
		pids[p.PID] = p
	}
	var children []int
	for _, p := range processes {
		child := p
		for ; p.PPID != 0 && !HasKey(programs, p.PID); p = pids[p.PPID] {
		}
		if p != child && HasKey(programs, p.PID) {
			p.Mem += child.Mem
			p.CPUTime += child.CPUTime
			children = append(children, child.PID)
		}
	}
	for _, child := range children {
		delete(pids, child)
	}

	var totalMem int
	var totalPercent float64
	percentages := map[int]float64{}
	for pid, p := range pids {
		totalMem += p.Mem
		elapsed := uptime - float64(p.StartTime)/system.Tick
		percent := 100 * float64(p.CPUTime) / system.Tick / elapsed
		totalPercent += percent
		if percent > minCPU {
			percentages[pid] = percent
		}
	}

	var cpuNameWidth int
	for pid := range percentages {
		if l := len(pids[pid].Name); l > cpuNameWidth {
			cpuNameWidth = l
		}
	}

	byMem := maps.Keys(pids)
	byCPU := slices.Clone(byMem)

	slices.SortFunc(byMem, func(a, b int) bool { return pids[a].Mem > pids[b].Mem })
	slices.SortFunc(byCPU, func(a, b int) bool { return percentages[a] > percentages[b] })

	cpuWidth := len(fmt.Sprintf("%.1f", percentages[byCPU[0]]))
	for i := 0; i < len(pids); i++ {
		{ // left column
			pid := byCPU[i]
			percent := percentages[pid]
			if percent > 0 {
				barWidth := int(math.Round(12 * 8 * percent / totalPercent))
				fmt.Printf("%*.1f%% [%-*s] %-*s",
					cpuWidth, percent,
					12, UnicodeBar(barWidth),
					cpuNameWidth, pids[pid].Name,
				)
			} else {
				if pids[byMem[i]].Mem < minMem {
					break
				}
				// padding with spaces
				fmt.Printf("%*s %*s",
					cpuWidth+1, "",
					cpuNameWidth+len(`[            ] `), "",
				)
			}
		}
		fmt.Print("   ")
		{ // right column
			if p := pids[byMem[i]]; p.Mem < minMem {
				fmt.Println()
			} else {
				barWidth := int(math.Round(12 * 8 * float64(p.Mem) / float64(totalMem)))
				fmt.Printf("%6s [%-*s] %s\n",
					FormatByte(p.Mem),
					12, UnicodeBar(barWidth),
					p.Name,
				)
			}
		}
	}
}
