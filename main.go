package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var cacheDir = "/var/cache/nginx/"

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "Service is up and running",
		})
	})

	r.GET("/purge", func(c *gin.Context) {
		path, status := c.GetQuery("path")
		if !status {
			c.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "Path is empty",
			})
			return
		}
		go clearNginxCache(path)
		c.JSON(http.StatusOK, Response{
			Success: true,
			Message: "Cache clearing request successfully submitted.",
		})
	})
	r.Run("0.0.0.0:4040")
}

func clearNginxCache(path string) {
	key := fmt.Sprintf("KEY: %s", path)

	fmt.Printf("Key is: %s\n", key)

	cmd := exec.Command("grep", "-r", "-l", key, cacheDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("No cache found.\n")
		return
	}

	type cache struct {
		Key     string `json:"key"`
		Address string `json:"address"`
	}

	var cacheArray []cache

	temp := string(output)
	temp2 := strings.Split(temp, "\n")
	temp3 := temp2[:len(temp2)-1]
	for _, v := range temp3 {
		address := v

		sed_cmd := exec.Command("sed", "-n", "2p", address)
		sed_output, sed_err := sed_cmd.CombinedOutput()
		if sed_err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		key := strings.Split(string(sed_output), "KEY:")[1]

		cacheArray = append(cacheArray, cache{
			Key:     key,
			Address: address,
		})

		rm_cmd := exec.Command("rm", address)
		_, rm_err := rm_cmd.CombinedOutput()
		if rm_err != nil {
			fmt.Printf("Error in deleting cache: %s\n", err)
			return
		}

	}

	fmt.Printf("%d caches deleted. \n", len(cacheArray))
}
