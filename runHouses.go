package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

func main() {
	N := 100 //set iterations

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

func printTableHeader(outputFile *os.File, columnNames []string) {
	fmt.Fprint(outputFile, "Statistic\t")
	for _, col := range columnNames {
		fmt.Fprintf(outputFile, "%-15s\t", col)
	}
	fmt.Fprintln(outputFile)
}

// print summary stats into output file
func printTableRow(outputFile *os.File, statType string, stats map[string]columnStats, columnNames []string) {
	fmt.Fprintf(outputFile, "%-15s\t", statType)
	for _, col := range columnNames {
		colStat := stats[col]
		var statValue float64

		switch statType {
		case "Min":
			statValue = colStat.Min
		case "Max":
			statValue = colStat.Max
		case "Mean":
			statValue = colStat.Sum / float64(colStat.Count)
		case "Median":
			statValue = calculateMedian(colStat.Values)
		case "1st Quartile":
			statValue = calculateQuartile(colStat.Values, 0.25)
		case "3rd Quartile":
			statValue = calculateQuartile(colStat.Values, 0.75)
		}

		fmt.Fprintf(outputFile, "%15.4f\t", statValue)
	}
	fmt.Fprintln(outputFile)
}

// build median calculation
func calculateMedian(data []float64) float64 {
	sort.Float64s(data)
	n := len(data)
	if n%2 == 0 {
		mid := n / 2
		return (data[mid-1] + data[mid]) / 2.0
	}
	return data[n/2]
}

// build quartile function
func calculateQuartile(data []float64, percentile float64) float64 {
	sort.Float64s(data)
	n := len(data)
	index := int(percentile * float64(n-1))
	return data[index]
}

type columnStats struct {
	Min    float64
	Max    float64
	Sum    float64
	Count  int
	Values []float64
}
