package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"testing"
)

func BenchmarkImportData(b *testing.B) {
	N := 100
	for i := 0; i < b.N; i++ {
		runImportDataBenchmark(N)
	}
}

func runImportDataBenchmark(N int) {
	//create output file
	outputFile, err := os.Create("housesOutputGo.txt")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	//input data
	inputFile, err := os.Open("housesInput.csv")
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	//read in data
	reader := csv.NewReader(inputFile)
	columnNames, err := reader.Read()
	if err != nil {
		panic(err)
	}

	//initiate loop for iterations
	for i := 0; i < N; i++ {
		_, err := inputFile.Seek(0, 0)
		if err != nil {
			panic(err)
		}

		//ignore header row
		_, err = reader.Read()
		if err != nil {
			panic(err)
		}

		//get summary stats for each column
		stats := make(map[string]columnStats)
		for _, col := range columnNames {
			stats[col] = columnStats{Min: math.MaxFloat64, Max: -math.MaxFloat64, Sum: 0.0, Count: 0, Values: []float64{}}
		}

		for {
			record, err := reader.Read()
			if err != nil {
				break
			}

			for colIdx, colValue := range record {
				if colStat, ok := stats[columnNames[colIdx]]; ok {
					num, err := strconv.ParseFloat(colValue, 64)
					if err == nil {
						if num < colStat.Min {
							colStat.Min = num
						}
						if num > colStat.Max {
							colStat.Max = num
						}
						colStat.Sum += num
						colStat.Count++
						colStat.Values = append(colStat.Values, num)
					}
					stats[columnNames[colIdx]] = colStat
				}
			}
		}

		// set column headers
		printTableHeader(outputFile, columnNames)

		//set summary stats
		statTypes := []string{"Min", "Max", "Mean", "Median", "1st Quartile", "3rd Quartile"}
		for _, statType := range statTypes {
			printTableRow(outputFile, statType, stats, columnNames)
		}

		fmt.Fprintln(outputFile, "---------------------------------------------")
	}
}
