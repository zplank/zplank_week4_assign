package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type housesInput struct {
	value    int
	income   float64
	age      int
	rooms    int
	bedrooms int
	pop      int
	hh       int
}

func main() {
	// Open the CSV file
	file, err := os.Open("housesInput.csv")
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	// Process each record
	var houses []housesInput
	for _, record := range records {
		housesinput := housesInput{
			value:    convertToInt(record[0]),
			income:   convertToFloat64(record[1]),
			age:      convertToInt(record[2]),
			rooms:    convertToInt(record[3]),
			bedrooms: convertToInt(record[4]),
			pop:      convertToInt(record[5]),
			hh:       convertToInt(record[6]),
		}
		houses = append(houses, housesinput)
	}

	// Print the processed data
	for _, housesinput := range houses {
		fmt.Printf("value: %d, income: %f, age: %d, rooms: %d, bedrooms: %d, pop: %d, hh: %d\n", housesinput.value, housesinput.income, housesinput.age, housesinput.rooms, housesinput.bedrooms, housesinput.pop, housesinput.hh)
	}
}

// Helper function to convert string to int
func convertToInt(s string) int {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		fmt.Println("Error converting to int:", err)
	}
	return result
}

// Helper function to convert string to float64
func convertToFloat64(s string) float64 {
	var result float64
	_, err := fmt.Sscanf(s, "%f", &result)
	if err != nil {
		fmt.Println("Error converting to float64:", err)
	}
	return result
}
