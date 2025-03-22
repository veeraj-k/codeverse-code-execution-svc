package scheduler

import (
	"context"
	"fmt"
	"lc-code-execution-service/types"
	"log"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func SpinUpContainer(job types.Job, path string, image string) error {

	apiclient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	defer apiclient.Close()
	ctx := context.Background()
	// timeout := 10 * time.Second
	// ctx, cancel := context.WithTimeout(context.Background(), timeout)
	// defer cancel()
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}

	resp, err := apiclient.ContainerCreate(ctx, &container.Config{
		Image: image,
	}, &container.HostConfig{

		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: absolutePath,
				Target: "/app",
			},
		},
		Resources: container.Resources{
			Memory:   256 * 1024 * 1024, // 256MB
			NanoCPUs: 250000000,         // 0.5 CPU (500ms CPU time per second)
		},
		AutoRemove: true,
	}, nil, nil, "submission_"+job.JobID)

	fmt.Println("Container ID: ", resp.ID)
	if err != nil {
		return err
	}

	if err := apiclient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}
	fmt.Println("Container Started")
	waitCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	waitC, errC := apiclient.ContainerWait(waitCtx, resp.ID, "")

	// fmt.Println("Container Wait Result: ", waitRes.StatusCode)
	// fmt.Println("Container Wait Err: ", waitRes.Error.Message)
	select {
	case val := <-waitC:
		if val.StatusCode == 0 {
			fmt.Println("Container Exited Successfully")
			return nil
		} else if val.StatusCode == 137 {
			return fmt.Errorf("resource Limit Exceeded")
		}

	case err := <-errC:
		if err != nil {
			fmt.Println("ErrC received:", err)
			return err
		}
	}

	// if waitRes := <-waitC; waitRes.StatusCode == 0 {
	// 	return nil
	// }

	// if err := <-errC; err != nil {
	// 	fmt.Printf("Error: %v", err)
	// 	fmt.Println("Stopping Container")
	// 	apiclient.ContainerStop(ctx, resp.ID, container.StopOptions{})

	// 	// log.Fatal(err)

	// 	return err
	// }
	fmt.Println("Container Stopped")

	return nil
}
