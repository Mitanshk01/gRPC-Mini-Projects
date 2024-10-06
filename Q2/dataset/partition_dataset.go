package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
)

const (
    numServers = 5
    inputFile  = "dataset.txt"
    outputFilePattern = "dataset_server_%d.txt"
)

func partitionDataset() error {
    file, err := os.Open(inputFile)
    if err != nil {
        return fmt.Errorf("failed to open dataset: %v", err)
    }
    defer file.Close()

    var dataset []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        dataset = append(dataset, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        return fmt.Errorf("error reading dataset: %v", err)
    }

    totalDataPoints := len(dataset)
    partitionSize := totalDataPoints / numServers
    extra := totalDataPoints % numServers

    start := 0
    for i := 0; i < numServers; i++ {
        end := start + partitionSize
        if i < extra {
            end++
        }

        outputFileName := fmt.Sprintf(outputFilePattern, i+1)
        if err := writePartitionToFile(outputFileName, dataset[start:end]); err != nil {
            return fmt.Errorf("failed to write partition: %v", err)
        }

        start = end
    }

    return nil
}

func writePartitionToFile(filename string, data []string) error {
    outFile, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("failed to create file: %v", err)
    }
    defer outFile.Close()

    writer := bufio.NewWriter(outFile)
    for _, line := range data {
        if _, err := writer.WriteString(line + "\n"); err != nil {
            return fmt.Errorf("failed to write data: %v", err)
        }
    }

    if err := writer.Flush(); err != nil {
        return fmt.Errorf("failed to flush data: %v", err)
    }

    log.Printf("Partition written to %s", filename)
    return nil
}

func main() {
    err := partitionDataset()
    if err != nil {
        log.Fatalf("Error partitioning dataset: %v", err)
    }
}
