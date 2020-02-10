package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexeldeib/iomond/iostat"
	"github.com/alexeldeib/iomond/limits"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(8 * time.Second)
	for {
		select {
		case <-signals:
			return
		case t := <-ticker.C:
			fmt.Println("Tick at", t)
			monitor()
		}
	}
}

func monitor() {
	limit, err := limits.New()
	if err != nil {
		log.Fatal(err)
	}

	stats, err := iostat.New()
	if err != nil {
		log.Fatal(err)
	}

	threshold := .8

	for disk, sku := range limit.Individual {
		diskStats, ok := stats[disk]
		if !ok {
			log.Fatal(fmt.Errorf("failed to find disk %s in iostat output", disk))
		}
		iops := diskStats.RS + diskStats.WS
		throughputMB := (diskStats.RkBS + diskStats.WkBS) / 1024
		if iops >= sku.Limit.IOPS*threshold {
			fmt.Printf("disk /dev/%s above threshold %.2f of iops limit\ncurrent: %.2f \nlimit: %.2f\n\n", disk, threshold, iops, sku.Limit.IOPS)
		} else {
			fmt.Printf("disk /dev/%s below threshold %.2f of iops limit\ncurrent: %.2f \nlimit: %.2f\n\n", disk, threshold, iops, sku.Limit.IOPS)
		}
		if throughputMB >= sku.Limit.Throughput*threshold {
			fmt.Printf("disk /dev/%s above threshold %.2f of throughput limit\ncurrent: %.2f MB/s\nlimit: %.2f\n\n", disk, threshold, throughputMB, sku.Limit.Throughput)
		} else {
			fmt.Printf("disk /dev/%s below threshold %.2f of throughput limit\ncurrent: %.2f MB/s\nlimit: %.2f\n\n", disk, threshold, throughputMB, sku.Limit.Throughput)
		}
	}

	var totalIOPS, totalThroughput float64
	for _, diskStats := range stats {
		totalIOPS = totalIOPS + diskStats.RS + diskStats.WS
		totalThroughput = totalThroughput + (diskStats.RkBS+diskStats.WkBS)/1024
	}
	if totalIOPS >= limit.Total.Uncached.IOPS*threshold {
		fmt.Printf("vm above threshold %.2f of iops limit\ncurrent: %.2f \nlimit: %.2f\n\n", threshold, totalIOPS, limit.Total.Uncached.IOPS)
	} else {
		fmt.Printf("vm below threshold %.2f of iops limit\ncurrent: %.2f \nlimit: %.2f\n\n", threshold, totalIOPS, limit.Total.Uncached.IOPS)
	}
	if totalThroughput >= limit.Total.Uncached.Throughput*threshold {
		fmt.Printf("vm above threshold %.2f of throughput limit\ncurrent: %.2f MB/s\nlimit: %.2f\n\n", threshold, totalThroughput, limit.Total.Uncached.Throughput)
	} else {
		fmt.Printf("vm below threshold %.2f of throughput limit\ncurrent: %.2f MB/s\nlimit: %.2f\n\n", threshold, totalThroughput, limit.Total.Uncached.Throughput)
	}
}
