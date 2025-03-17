package scheduler

import (
	"fmt"
	"lc-code-execution-service/types"
	"lc-code-execution-service/util"
	"os"
)

func StartWorkerPool(jobQueue <-chan types.Job, maxWorkers int) {

	for i := 0; i < maxWorkers; i++ {

		go worker(jobQueue)
	}
}

func worker(jobQueue <-chan types.Job) {
	fmt.Println("Spun up a worker")
	for job := range jobQueue {

		// util.PerformPutRequest(fmt.Sprintf("%s/api/jobs/%s/status", os.Getenv("SUBMISSION_SERVICE_URL"), job.JobID), "RUNNING", nil)
		util.PerformPutRequestWithQueryParams(fmt.Sprintf("%s/api/submission/%s/status", os.Getenv("SUBMISSION_SERVICE_URL"), job.JobID), map[string]string{"status": "EXECUTING"}, nil)
		err := ExecuteCode(&job)

		// err := SpinUpContainer(job)
		DeleteFolderRecursive(fmt.Sprintf("submission/%s", job.JobID))

		if err != nil {
			fmt.Println("Error", err.Error())

			util.PerformPutRequestWithQueryParams(fmt.Sprintf("%s/api/submission/%s/status", os.Getenv("SUBMISSION_SERVICE_URL"), job.JobID), map[string]string{"status": "FAILED", "message": err.Error()}, nil)
			continue
		}
		// DeleteFolderRecursive(fmt.Sprintf("submission/%s", job.JobID))
		util.PerformPutRequestWithQueryParams(fmt.Sprintf("%s/api/submission/%s/status", os.Getenv("SUBMISSION_SERVICE_URL"), job.JobID), map[string]string{"status": "COMPLETED"}, nil)

	}

}
