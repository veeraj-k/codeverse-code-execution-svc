package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
)

func LoadEnv() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
func LoadImages() {
	apiclient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	defer apiclient.Close()

	_, err = apiclient.ImagePull(context.Background(), os.Getenv("PY_IMG"), image.PullOptions{})

	if err != nil {
		panic(err)
	}
	_, err = apiclient.ImagePull(context.Background(), os.Getenv("JAVA_IMG"), image.PullOptions{})

	if err != nil {
		panic(err)
	}
	_, err = apiclient.ImagePull(context.Background(), os.Getenv("C_IMG"), image.PullOptions{})

	if err != nil {
		panic(err)
	}
	fmt.Println("Images loaded")
}
