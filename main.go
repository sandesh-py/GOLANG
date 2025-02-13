package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Car struct {
	ID       string `json:"id"`
	Brand    string `json:"brand"`
	Number   string `json:"number"`
	Type     string `json:"type"`
	Incoming string `json:"incoming_time"`
	Outgoing string `json:"outgoing_time"`
	Slot     string `json:"parking_slot"`
}

var dataFile = "data.json"

func loadCars() ([]Car, error) {
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var cars []Car
	err = json.NewDecoder(file).Decode(&cars)
	if err != nil {
		return nil, err
	}
	return cars, nil
}

func saveCars(cars []Car) error {
	file, err := json.MarshalIndent(cars, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dataFile, file, 0644)
}

func createCar(c *gin.Context) {
	var newCar Car
	if err := c.ShouldBindJSON(&newCar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Load existing cars
	cars, _ := loadCars()

	// Determine the next available ID
	newID := "1" // Default ID if the file is empty
	if len(cars) > 0 {
		lastCar := cars[len(cars)-1] // Get last car in the list
		lastID, err := strconv.Atoi(lastCar.ID)
		if err == nil {
			newID = strconv.Itoa(lastID + 1) // Increment ID
		}
	}

	// Assign the new ID
	newCar.ID = newID

	// Append and save
	cars = append(cars, newCar)
	saveCars(cars)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Car added successfully!",
		"car":     newCar,
	})
}

func getCars(c *gin.Context) {
	cars, err := loadCars()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading cars"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Cars retrieved successfully",
		"cars":    cars,
	})
}

func getCar(c *gin.Context) {
	id := c.Param("id")
	cars, _ := loadCars()
	for _, car := range cars {
		if car.ID == id {
			c.JSON(http.StatusOK, gin.H{
				"message": "Car retrieved successfully",
				"car":     car,
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
}

func updateCar(c *gin.Context) {
	id := c.Param("id")
	var updatedCar Car
	if err := c.ShouldBindJSON(&updatedCar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	cars, _ := loadCars()
	for i, car := range cars {
		if car.ID == id {
			cars[i] = updatedCar
			saveCars(cars)
			c.JSON(http.StatusOK, gin.H{
				"message": "Car updated successfully!",
				"car":     updatedCar,
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
}

func deleteCar(c *gin.Context) {
	id := c.Param("id")
	cars, _ := loadCars()
	for i, car := range cars {
		if car.ID == id {
			cars = append(cars[:i], cars[i+1:]...)
			saveCars(cars)
			c.JSON(http.StatusOK, gin.H{"message": "Car deleted successfully!"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
}

func main() {
	r := gin.Default()

	// Enable CORS Middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // React Frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/cars", getCars)
	r.GET("/cars/:id", getCar)
	r.POST("/cars", createCar)
	r.PUT("/cars/:id", updateCar)
	r.DELETE("/cars/:id", deleteCar)

	port := "8080"
	fmt.Printf("Server is running at: http://localhost:%s\n", port)
	r.Run(":" + port)
}
