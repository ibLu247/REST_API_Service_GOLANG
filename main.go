package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Car struct {
	ID             int64  `json:"id"`
	Brand          string `json:"brand"`
	Model          string `json:"model"`
	Mileage        string `json:"mileage"`
	NumberOfOwners uint8  `json:"numberOfOwners"`
}

var cars []Car

// Функция загружает данные из файла "db.json"
func loadData() {
	file, err := os.Open("db.json")
	if err != nil {
		fmt.Println("Не удалось открыть файл", err)
	}

	writtenData := json.NewDecoder(file)
	writtenData.Decode(&cars)
}

// Функция сохраняет полученные из запроса данные в файл "db.json"
func saveData() {
	file, err := os.Create("db.json")
	if err != nil {
		fmt.Println("Не удалось создать файл", err)
		return
	}
	defer file.Close()

	recievedValue := json.NewEncoder(file)
	recievedValue.SetIndent("", " ")
	recievedValue.Encode(cars)
}

// Добавить автомобиль
func addCar(c *gin.Context) {
	var newCar Car

	if err := c.BindJSON(&newCar); err != nil {
		c.IndentedJSON(400, gin.H{"Ошибка": "Не удалось добавить автомобиль"})
		return
	}

	cars = append(cars, newCar)
	saveData()
	c.IndentedJSON(http.StatusCreated, newCar)
}

// Получить список всех автомобилей
func getCars(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, cars)
}

// Получить авто по id
func getCar(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"Ошибка": "Неверный ID"})
	}

	for _, car := range cars {
		if car.ID == id {
			c.IndentedJSON(http.StatusOK, car)
		}
	}
}

// Обновить все поля авто
func fullUpdateCar(c *gin.Context) {
	var updatedCar Car
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"Ошибка": "Неверный ID"})
	}

	if err := c.BindJSON(&updatedCar); err != nil {
		c.IndentedJSON(400, gin.H{"Ошибка": "Не удалось обновить автомобиль"})
		return
	}

	for i, car := range cars {
		if car.ID == id {
			cars[i] = updatedCar
			saveData()
			c.IndentedJSON(http.StatusOK, updatedCar)
			return
		}
	}
}

// Обновить авто частично
func updateCar(c *gin.Context) {
	var updatedCar Car
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"Ошибка": "Неверный ID"})
	}

	if err := c.BindJSON(&updatedCar); err != nil {
		c.IndentedJSON(400, gin.H{"Ошибка": "Не удалось обновить автомобиль"})
		return
	}

	for i, car := range cars {
		if car.ID == id {
			if updatedCar.ID != 0 {
				cars[i].ID = updatedCar.ID
			}
			if updatedCar.Brand != "" {
				cars[i].Brand = updatedCar.Brand
			}
			if updatedCar.Model != "" {
				cars[i].Model = updatedCar.Model
			}
			if updatedCar.Mileage != "" {
				cars[i].Mileage = updatedCar.Mileage
			}
			if updatedCar.NumberOfOwners != 0 {
				cars[i].NumberOfOwners = updatedCar.NumberOfOwners
			}
			saveData()
			c.Status(204)
			return
		}
	}

}

// Удалить авто
func deleteCar(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"Ошибка": "Неверный ID"})
	}

	for i, car := range cars {
		if car.ID == id {
			cars = append(cars[:i], cars[i+1:]...)
			saveData()
			c.Status(204)
			return
		}
	}
}

func main() {
	// Функция loadData() выполняется самая первая, чтобы не перезаписались ранее созданные записи
	loadData()

	router := gin.Default()

	router.POST("/cars", addCar)
	router.GET("/cars", getCars)
	router.GET("/cars/:id", getCar)
	router.PUT("/cars/:id", fullUpdateCar)
	router.PATCH("/cars/:id", updateCar)
	router.DELETE("/cars/:id", deleteCar)

	router.Run("localhost:8080")
}
