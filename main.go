package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Weather struct {
	Status struct {
		Water int `json:"water"`
		Wind  int `json:"wind"`
	} `json:"status"`
}

var (
	fi  []byte
	di  Weather
	do  []byte
	err error
)

func main() {
	go func() {
		for {
			<-time.After(15 * time.Second)
			func() {
				fi, err = os.ReadFile("status.json")
				Check(err)

				// json --> struct
				err = json.Unmarshal(fi, &di)
				Check(err)

				di.Randomize()

				// struct --> json
				do, err = json.MarshalIndent(di, "", "  ")
				Check(err)

				// write output file with 52 bytes
				err = os.WriteFile("status.json", do, 0644)
				Check(err)
			}()

			fmt.Printf(
				"water : %d m   - %s \nwind  : %d m/s - %s \n",
				di.Status.Water,
				StatusWater(di.Status.Water),
				di.Status.Wind,
				StatusWind(di.Status.Wind),
			)
		}
	}()

	routers := gin.Default()
	routers.LoadHTMLGlob("templates/*.html")
	routers.Static("/static", "./templates/static")
	routers.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	routers.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"water":     di.Status.Water,
			"wind":      di.Status.Wind,
			"waterStat": StatusWater(di.Status.Water),
			"windStat":  StatusWind(di.Status.Wind),
			"name":      "szczynk",
		})
	})
	routers.Run(":8080")
}

func (w *Weather) Randomize() {
	rand.Seed(time.Now().UTC().UnixNano())
	min, max := 1, 100
	w.Status.Water = rand.Intn(max-min+1) + min
	w.Status.Wind = rand.Intn(max-min+1) + min
}

func StatusWater(water int) string {
	var waterStat string
	switch {
	case (water < 6):
		waterStat = "Aman"
	case (water >= 6) && (water <= 8):
		waterStat = "Siaga"
	default:
		waterStat = "Bahaya"
	}
	return waterStat
}

func StatusWind(wind int) string {
	var windStat string
	switch {
	case (wind < 7):
		windStat = "Aman"
	case (wind >= 7) && (wind <= 15):
		windStat = "Siaga"
	default:
		windStat = "Bahaya"
	}
	return windStat
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
