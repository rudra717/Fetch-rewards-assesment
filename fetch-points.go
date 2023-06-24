package main

import (
	"net/http"
	"strconv"
	"strings"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Receipt struct {
	Retailer     string  `json:"retailer"`
	Total        string  `json:"total"`
	Items        []Item  `json:"items"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
}

type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            string `json:"price"`
}

type ReceiptsMap map[string]int

var receipts ReceiptsMap

func main() {
	receipts = make(ReceiptsMap)
	router := gin.Default()
	router.POST("/receipts/process", processReceipts)
	router.GET("/receipts/:receipt_id/points", getPoints)
	router.Run(":8080")
}

func processReceipts(c *gin.Context) {
	var receipt Receipt
	err := c.ShouldBindJSON(&receipt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if receipt.Items == nil{
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt cannot be made without Items"})
		return
	}
	if receipt.PurchaseTime == ""{
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt should have a purchase time"})
		return
	}
	if receipt.PurchaseDate == ""{
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt should have a purchase date"})
		return
	}
	if receipt.Retailer == ""{
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt should have a retailer"})
		return
	}
	if receipt.Total == ""{
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt should have a Total"})
		return
	}
	receiptID := uuid.New().String()
	points := calculatePoints(receipt)
	receipts[receiptID] = points

	c.JSON(http.StatusOK, gin.H{"id": receiptID})
}

func getPoints(c *gin.Context) {
	receiptID := c.Param("receipt_id")
	points, ok := receipts[receiptID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"points": points})
}

func calculatePoints(receipt Receipt) int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name
	points += countAlphanumeric(receipt.Retailer)

	// Rule 2: 50 points if the total is a round dollar amount
	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil && total == float64(int(total)) {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if math.Mod(total*100, 25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt
	if len(receipt.Items) > 0 {
		points += len(receipt.Items) / 2 * 5
	} else {
		points = 0 // Set points to zero if there are no items
	}

	// Rule 5: Multiply the price by 0.2 and round up to the nearest integer if the trimmed length of
	// the item description is a multiple of 3. The result is the number of points earned.
	for _, item := range receipt.Items {
		trimmedLength := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLength%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				points += int(math.Ceil(price * 0.2))
			}
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd
	day, err := strconv.Atoi(strings.Split(receipt.PurchaseDate, "-")[2])
	if err == nil && day%2 != 0 {
		points += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm
	hour, err := strconv.Atoi(strings.Split(receipt.PurchaseTime, ":")[0])
	if err == nil && hour >= 14 && hour < 16 {
		points += 10
	}

	return points
}

func countAlphanumeric(s string) int {
	count := 0
	for _, ch := range s {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
			count++
		}
	}
	return count
}
