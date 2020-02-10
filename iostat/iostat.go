package iostat

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type Disk struct {
	DiskDevice string  `json:"disk_device"`
	RS         float64 `json:"r/s"`
	WS         float64 `json:"w/s"`
	RkBS       float64 `json:"rkB/s"`
	WkBS       float64 `json:"wkB/s"`
	RrqmS      float64 `json:"rrqm/s"`
	WrqmS      float64 `json:"wrqm/s"`
	Rrqm       float64 `json:"rrqm"`
	Wrqm       float64 `json:"wrqm"`
	RAwait     float64 `json:"r_await"`
	WAwait     float64 `json:"w_await"`
	AquSz      float64 `json:"aqu-sz"`
	RareqSz    float64 `json:"rareq-sz"`
	WareqSz    float64 `json:"wareq-sz"`
	Svctm      float64 `json:"svctm"`
	Util       float64 `json:"util"`
}

const iostatBinary = "iostat"

func New() (map[string]Disk, error) {
	if _, err := exec.LookPath(iostatBinary); err != nil {
		return nil, err
	}

	cmd := exec.Command(iostatBinary, []string{"-xty", "1", "1"}...)
	cmd.Env = []string{"S_TIME_FORMAT=ISO"}

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(stdout), "\n")
	if len(lines) < 8 {
		return nil, fmt.Errorf("too few lines to parse to iostat output")
	}

	stats := map[string]Disk{}
	for _, line := range lines[7 : len(lines)-2] {
		parts := strings.Split(line, " ")
		numerics := make([]float64, 16)
		n := 1
		for _, part := range parts[1:] {
			if part == "" {
				continue
			}
			numeric, err := strconv.ParseFloat(part, 64)
			if err != nil {
				return nil, err
			}
			numerics[n] = numeric
			n++
		}

		if len(numerics) != 16 {
			log.Fatal(fmt.Errorf("expected 16 parts per line but found %d", len(numerics)))
		}

		datum := Disk{
			parts[0],
			numerics[1],
			numerics[2],
			numerics[3],
			numerics[4],
			numerics[5],
			numerics[6],
			numerics[7],
			numerics[8],
			numerics[9],
			numerics[10],
			numerics[11],
			numerics[12],
			numerics[13],
			numerics[14],
			numerics[15],
		}
		stats[datum.DiskDevice] = datum
	}

	return stats, nil
}
