# 4600project1


**NOTE: EDITED main.go FILE IS IN THE PROJECT1 FOLDER! **
# Process Scheduling Algorithms

This program carries out different scheduling algorithms, including First Come First Serve (FCFS), Shortest Job First (SJF), Shortest Job First with Priority (SJF Priority), Round-Robin (RR)

## Getting Started

### Prerequisites

- Go (Golang) must be installed on system

### Installation

1. Clone repository to local machine:

```
git clone https://github.com/yourusername/scheduling-algorithms
```

2. Navigate to project directory:

```
cd scheduling-algorithms
```

### Usage

1. Compile Go code:

```
go build
```

2. Run the program while passing a command-line argument for an input file. 'input.csv' should be replaced with input file of choice:

```
./scheduling-algorithms input.csv
```

As an alternative, you can use this command to run the application without compiling.


```
go run main.go input.csv
```

### Input Format

The input file should be in CSV format, with each line representing a process. Each line should include the following fields:

```
<ProcessID>,<Burst Duration>,<Arrival Time>,<Priority>
```

- `<ProcessID>`: This field is used to uniquely identify each process.
- `<Burst Duration>`: This is an integer value that represents the time needed to finish executing the process.
- `<Arrival Time>`: An integer value indicating the point in time when the process enters the system.
- `<Priority>`: This field, if provided, is an integer between 1 and 50 and represents the priority of the process. Note that it's only applicable for SJF Priority scheduling.

### Output

The program will produce a GANTT chart and a summary table displaying the scheduling outcomes. The results will encompass calculations for the average turnaround time, average waiting time, and average throughput.


