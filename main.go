package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xuri/excelize/v2"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
)

type DataChart struct {
	Tanggal time.Time
	Kode    string
	Data    string
}

func createRedisClient(ctx context.Context) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return client
}

func main() {
	ctx := context.Background()
	redisClient := createRedisClient(ctx)
	defer redisClient.Close()

	f, err := excelize.OpenFile("tes.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	tanggal, err := f.GetCellValue("HSL PACK", "E2")
	fmt.Println("Nilai Tanggal New :", tanggal)

	// kode, err := f.GetCellValue("HSL PACK", "A2")
	// // fmt.Println("Nilai A2:", kode)

	// data, err := f.GetCellValue("HSL PACK", "I5")
	// // fmt.Println("Nilai I5:", data)

	// tanggalTime, err := time.Parse("02/01/2006", tanggal)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// SaveDataToRedis(redisClient, tanggalTime, kode, data)

	startDateStr := "25-11-2023"
	endDateStr := "28-11-2023"

	startDate, err := time.Parse("02-01-2006", startDateStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	endDate, err := time.Parse("02-01-2006", endDateStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ReadDataFromRedis(redisClient, startDate, endDate)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func SaveDataToRedis(client *redis.Client, tanggal time.Time, kode string, data string) {
	ctx := context.Background()
	key := fmt.Sprintf("%s", tanggal.Format("20060102"))

	client.HSet(ctx, key, "kode", kode)
	client.HSet(ctx, key, "data", data)
}

func ReadDataFromRedis(redisClient *redis.Client, startDate, endDate time.Time) error {

	// f := excelize.NewFile()

	// Menginisialisasi header untuk tanggal
	// header := []interface{}{nil}

	// Iterasi melalui rentang tanggal dan menambahkan setiap tanggal ke dalam header
	currentDate := startDate
	// for currentDate.Before(endDate) || currentDate.Equal(endDate) {
	// 	header = append(header, currentDate.Format("02/01/2006"))
	// 	currentDate = currentDate.Add(24 * time.Hour)
	// }

	// // Menulis header ke file Excel
	// if err := f.SetSheetRow("Sheet1", "A1", &header); err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	// currentDate := startDate
	// rowIndex := 2
	dataMap := make(map[string]map[string]string)

	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		key := fmt.Sprintf("%s", currentDate.Format("20060102"))

		kode, err := redisClient.HGet(ctx, key, "kode").Result()
		if err != nil {
			// fmt.Printf("Error membaca data untuk %s: %v\n", currentDate.Format("2006-01-02"), err)
			return excelize.ErrAddVBAProject
		} else {
			data, err := redisClient.HGet(ctx, key, "data").Result()
			if err != nil {
				// fmt.Printf("Error membaca data untuk %s: %v\n", currentDate.Format("2006-01-02"), err)
				return err
			} else {
				dataMap[key] = map[string]string{
					"Kode": kode,
					"Data": data,
				}
				fmt.Printf("Key : %s\n", currentDate.Format("02-01-2006"))
				fmt.Printf("Kode: %s\n", kode)
				fmt.Printf("Data: %s\n", data)
			}

		}
		currentDate = currentDate.Add(24 * time.Hour)
	}
	f := excelize.NewFile()

	sheetName := "DataSheet"
	f.NewSheet(sheetName)

	highestRow := 1 // Default to first row
	if rows, err := f.GetRows(sheetName); err == nil {
		highestRow = len(rows) + 1
	}

	for key, value := range dataMap {
		values := []interface{}{key, value["Kode"], value["Data"]}
		startCell, err := excelize.JoinCellName("A", highestRow)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if err := f.SetSheetRow(sheetName, startCell, &values); err != nil {
			fmt.Println(err)
			return err
		}
		highestRow++
	}

	f.SaveAs("laporan.xlsx")

	// for key, _ := range dataMap {
	// 	data := [][]interface{}{
	// 		{nil, key},
	// 		// {subKey, subValue},
	// 	}
	// 	fmt.Println(data[0]...)

	// 	for i, row := range data[0]... {
	// 		startCell, err := excelize.JoinCellName("A", i+1)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return err
	// 		}
	// 		if err := f.SetSheetRow(sheetName, startCell, &row); err != nil {
	// 			fmt.Println(err)
	// 			return err
	// 		}
	// 	}

	// 	if err := f.SaveAs("test.xlsx"); err != nil {
	// 		fmt.Println(err)
	// 		return err
	// 	}
	// }
	return nil
}
