package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func main() {
	// CLI args
	f, closeFile, err := openProcessingFile(os.Args...)
	if err != nil {
		log.Fatal(err)
	}
	defer closeFile()

	// Load and parse processes
	processes, err := loadProcesses(f)
	if err != nil {
		log.Fatal(err)
	}

	// First-come, first-serve scheduling
	FCFSSchedule(os.Stdout, "First-come, first-serve", processes)

	//SJFSchedule(os.Stdout, "Shortest-job-first", processes)
	//
	//SJFPrioritySchedule(os.Stdout, "Priority", processes)
	//
	//RRSchedule(os.Stdout, "Round-robin", processes)
}

func openProcessingFile(args ...string) (*os.File, func(), error) {
	if len(args) != 2 {
		return nil, nil, fmt.Errorf("%w: must give a scheduling file to process", ErrInvalidArgs)
	}
	// Read in CSV process CSV file
	f, err := os.Open(args[1])
	if err != nil {
		return nil, nil, fmt.Errorf("%v: error opening scheduling file", err)
	}
	closeFn := func() {
		if err := f.Close(); err != nil {
			log.Fatalf("%v: error closing scheduling file", err)
		}
	}

	return f, closeFn, nil
}

type (
	Process struct {
		ProcessID     int64
		ArrivalTime   int64
		BurstDuration int64
		Priority      int64
	}
	TimeSlice struct {
		PID   int64
		Start int64
		Stop  int64
	}
)

//region Schedulers

// FCFSSchedule outputs a schedule of processes in a GANTT chart and a table of timing given:
// • an output writer
// • a title for the chart
// • a slice of processes
func FCFSSchedule(w io.Writer, title string, processes []Process) {
	var (
		serviceTime     int64
		totalWait       float64
		totalTurnaround float64
		lastCompletion  float64
		waitingTime     int64
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
	)
	for i := range processes {
		if processes[i].ArrivalTime > 0 {
			waitingTime = serviceTime - processes[i].ArrivalTime
		}
		totalWait += float64(waitingTime)

		start := waitingTime + processes[i].ArrivalTime

		turnaround := processes[i].BurstDuration + waitingTime
		totalTurnaround += float64(turnaround)

		completion := processes[i].BurstDuration + processes[i].ArrivalTime + waitingTime
		lastCompletion = float64(completion)

		schedule[i] = []string{
			fmt.Sprint(processes[i].ProcessID),
			fmt.Sprint(processes[i].Priority),
			fmt.Sprint(processes[i].BurstDuration),
			fmt.Sprint(processes[i].ArrivalTime),
			fmt.Sprint(waitingTime),
			fmt.Sprint(turnaround),
			fmt.Sprint(completion),
		}
		serviceTime += processes[i].BurstDuration

		gantt = append(gantt, TimeSlice{
			PID:   processes[i].ProcessID,
			Start: start,
			Stop:  serviceTime,
		})
	}

	count := float64(len(processes))
	aveWait := totalWait / count
	aveTurnaround := totalTurnaround / count
	aveThroughput := count / lastCompletion

	outputTitle(w, title)
	outputGantt(w, gantt)
	outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}

func SJFPrioritySchedule(w io.Writer, title string, processes []Process) {
    sort.Slice(processes, func(i, j int) bool {
        if processes[i].BurstDuration == processes[j].BurstDuration {
            return processes[i].Priority < processes[j].Priority
        }
        return processes[i].BurstDuration < processes[j].BurstDuration
    })

    var (
        serviceTime     int64
        totalWait       float64
        totalTurnaround float64
        lastCompletion  float64
        schedule        = make([][]string, len(processes))
        gantt           = make([]TimeSlice, 0)
    )

    for i := range processes {
        if processes[i].ArrivalTime > serviceTime {
            serviceTime = processes[i].ArrivalTime
        }
        waitingTime := serviceTime - processes[i].ArrivalTime
        totalWait += float64(waitingTime)

        start := serviceTime
        turnaround := waitingTime + processes[i].BurstDuration
        totalTurnaround += float64(turnaround)

        completion := serviceTime + turnaround
        lastCompletion = float64(completion)

        schedule[i] = []string{
            fmt.Sprint(processes[i].ProcessID),
            fmt.Sprint(processes[i].Priority),
            fmt.Sprint(processes[i].BurstDuration),
            fmt.Sprint(processes[i].ArrivalTime),
            fmt.Sprint(waitingTime),
            fmt.Sprint(turnaround),
            fmt.Sprint(completion),
        }

        gantt = append(gantt, TimeSlice{
            PID:   processes[i].ProcessID,
            Start: start,
            Stop:  start + turnaround,
        })

        serviceTime += processes[i].BurstDuration
    }

    count := float64(len(processes))
    aveWait := totalWait / count
    aveTurnaround := totalTurnaround / count
    aveThroughput := count / lastCompletion

    outputTitle(w, title)
    outputGantt(w, gantt)
    outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}

func SJFSchedule(w io.Writer, title string, processes []Process) {
    sort.Slice(processes, func(i, j int) bool {
        return processes[i].BurstDuration < processes[j].BurstDuration
    })

    var (
        serviceTime     int64
        totalWait       float64
        totalTurnaround float64
        lastCompletion  float64
        schedule        = make([][]string, len(processes))
        gantt           = make([]TimeSlice, 0)
    )

    for i := range processes {
        if processes[i].ArrivalTime > serviceTime {
            serviceTime = processes[i].ArrivalTime
        }
        waitingTime := serviceTime - processes[i].ArrivalTime
        totalWait += float64(waitingTime)

        start := serviceTime
        turnaround := waitingTime + processes[i].BurstDuration
        totalTurnaround += float64(turnaround)

        completion := serviceTime + turnaround
        lastCompletion = float64(completion)

        schedule[i] = []string{
            fmt.Sprint(processes[i].ProcessID),
            fmt.Sprint(processes[i].Priority),
            fmt.Sprint(processes[i].BurstDuration),
            fmt.Sprint(processes[i].ArrivalTime),
            fmt.Sprint(waitingTime),
            fmt.Sprint(turnaround),
            fmt.Sprint(completion),
        }

        gantt = append(gantt, TimeSlice{
            PID:   processes[i].ProcessID,
            Start: start,
            Stop:  start + turnaround,
        })

        serviceTime += processes[i].BurstDuration
    }

    count := float64(len(processes))
    aveWait := totalWait / count
    aveTurnaround := totalTurnaround / count
    aveThroughput := count / lastCompletion

    outputTitle(w, title)
    outputGantt(w, gantt)
    outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}

func RRSchedule(w io.Writer, title string, processes []Process, timeQuantum int64) {
    var (
        time    int64
        queue   []Process
        gantt   []TimeSlice
        schedule = make([][]string, 0)
    )

    for len(processes) > 0 || len(queue) > 0 {
        if len(queue) == 0 {
            queue = append(queue, processes[0])
            processes = processes[1:]
        }
        process := queue[0]
        queue = queue[1:]

        start := time
        if process.BurstDuration <= timeQuantum {
            // Process completes within time quantum
            time += process.BurstDuration
            process.BurstDuration = 0
        } else {
            time += timeQuantum
            process.BurstDuration -= timeQuantum
            queue = append(queue, process)
        }

  
        waitingTime := start - process.ArrivalTime
        turnaround := waitingTime + process.BurstDuration

        schedule = append(schedule, []string{
            fmt.Sprint(process.ProcessID),
            fmt.Sprint(process.Priority),
            fmt.Sprint(process.BurstDuration),
            fmt.Sprint(process.ArrivalTime),
            fmt.Sprint(waitingTime),
            fmt.Sprint(turnaround),
            fmt.Sprint(time),
        })

        gantt = append(gantt, TimeSlice{
            PID:   process.ProcessID,
            Start: start,
            Stop:  time,
        })
    }

    count := float64(len(schedule))
    totalWait := 0.0
    totalTurnaround := 0.0
    lastCompletion := 0.0

    for _, entry := range schedule {
        totalWait += strToFloat(entry[4])
        totalTurnaround += strToFloat(entry[5])
        lastCompletion = strToFloat(entry[6])
    }

    aveWait := totalWait / count
    aveTurnaround := totalTurnaround / count
    aveThroughput := float64(len(schedule)) / lastCompletion

    outputTitle(w, title)
    outputGantt(w, gantt)
    outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}

func strToFloat(s string) float64 {
    f, err := strconv.ParseFloat(s, 64)
    if err != nil {
        log.Fatalf("Error converting string to float: %v", err)
    }
    return f
}


func outputTitle(w io.Writer, title string) {
	_, _ = fmt.Fprintln(w, strings.Repeat("-", len(title)*2))
	_, _ = fmt.Fprintln(w, strings.Repeat(" ", len(title)/2), title)
	_, _ = fmt.Fprintln(w, strings.Repeat("-", len(title)*2))
}

func outputGantt(w io.Writer, gantt []TimeSlice) {
	_, _ = fmt.Fprintln(w, "Gantt schedule")
	_, _ = fmt.Fprint(w, "|")
	for i := range gantt {
		pid := fmt.Sprint(gantt[i].PID)
		padding := strings.Repeat(" ", (8-len(pid))/2)
		_, _ = fmt.Fprint(w, padding, pid, padding, "|")
	}
	_, _ = fmt.Fprintln(w)
	for i := range gantt {
		_, _ = fmt.Fprint(w, fmt.Sprint(gantt[i].Start), "\t")
		if len(gantt)-1 == i {
			_, _ = fmt.Fprint(w, fmt.Sprint(gantt[i].Stop))
		}
	}
	_, _ = fmt.Fprintf(w, "\n\n")
}

func outputSchedule(w io.Writer, rows [][]string, wait, turnaround, throughput float64) {
	_, _ = fmt.Fprintln(w, "Schedule table")
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"ID", "Priority", "Burst", "Arrival", "Wait", "Turnaround", "Exit"})
	table.AppendBulk(rows)
	table.SetFooter([]string{"", "", "", "",
		fmt.Sprintf("Average\n%.2f", wait),
		fmt.Sprintf("Average\n%.2f", turnaround),
		fmt.Sprintf("Throughput\n%.2f/t", throughput)})
	table.Render()
}

//endregion

//region Loading processes.

var ErrInvalidArgs = errors.New("invalid args")

func loadProcesses(r io.Reader) ([]Process, error) {
	rows, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%w: reading CSV", err)
	}

	processes := make([]Process, len(rows))
	for i := range rows {
		processes[i].ProcessID = mustStrToInt(rows[i][0])
		processes[i].BurstDuration = mustStrToInt(rows[i][1])
		processes[i].ArrivalTime = mustStrToInt(rows[i][2])
		if len(rows[i]) == 4 {
			processes[i].Priority = mustStrToInt(rows[i][3])
		}
	}

	return processes, nil
}

func mustStrToInt(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return i
}

//endregion
