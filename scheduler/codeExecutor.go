package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lc-code-execution-service/types"
	"lc-code-execution-service/util"
	"net/http"
	"os"
	"strconv"
)

func ExecuteCode(job *types.Job) error {
	resp, err := http.Get(os.Getenv("PMS_SERVICE_URL") + "/api/problems/" + strconv.Itoa(int(job.ProblemId)) + "/metadata/" + job.Language)
	if err != nil {
		return err

	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("problem not found")
	}

	type ProblemMetadata struct {
		ProblemId int    `json:"problemId"`
		Language  string `json:"language"`
		Runner    string `json:"runner"`
		Solution  string `json:"evaluator"`
	}

	var problemMetadata ProblemMetadata
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&problemMetadata); err != nil {
		return err
	}
	fmt.Println("Creating directories")

	var fileEx, image string

	if job.Language == "python" {
		fileEx = ".py"
		image = "py-img-lc"
	} else if job.Language == "java" {
		fileEx = ".java"
		image = "java-img-lc"
	} else if job.Language == "c" {
		fileEx = ".c"
		image = "c-img-lc"
	}

	submissionDir := fmt.Sprintf("submission/%s", job.JobID)
	if err := os.MkdirAll(submissionDir, os.ModePerm); err != nil {
		return err
	}

	runnerFilePath := fmt.Sprintf("%s/Evaluator%s", submissionDir, fileEx)
	runnerFile, err := os.Create(runnerFilePath)
	if err != nil {
		return err
	}
	defer runnerFile.Close()

	if _, err := runnerFile.WriteString(problemMetadata.Runner); err != nil {
		return err
	}

	solutionFilePath := fmt.Sprintf("%s/Solution%s", submissionDir, fileEx)
	solutionFile, err := os.Create(solutionFilePath)
	if err != nil {
		return err
	}
	defer solutionFile.Close()

	if _, err := solutionFile.WriteString(job.Code); err != nil {
		return err
	}

	// resp1, err := http.Get(os.Getenv("PMS_SERVICE_URL") + "/api/problems/" + strconv.Itoa(int(job.ProblemId)) + "/testcases")
	resp1, err := util.PerformGetRequest(os.Getenv("PMS_SERVICE_URL") + "/api/problems/" + strconv.Itoa(int(job.ProblemId)) + "/testcases")

	if err != nil {
		return err
	}

	if resp1.StatusCode != 200 {
		return errors.New("testcases not found")
	}

	testCaseFilePath := fmt.Sprintf("%s/test_cases.json", submissionDir)
	testCasesFile, err := os.Create(testCaseFilePath)
	if err != nil {
		return err
	}
	defer testCasesFile.Close()

	bodyBytes, err := io.ReadAll(resp1.Body)
	if err != nil {
		return err
	}
	if _, err := testCasesFile.Write(bodyBytes); err != nil {
		return err
	}

	// time.Sleep(1000 * time.Millisecond)
	fmt.Println("Spinning up container")

	err = SpinUpContainer(*job, submissionDir, image)
	fmt.Println("Done spinning up container")

	if err != nil {
		return err
	}

	// time.Sleep(5000 * time.Millisecond)
	outputFilePath := fmt.Sprintf("%s/output.json", submissionDir)

	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		logs, _ := os.ReadFile(fmt.Sprintf("%s/logs.txt", submissionDir))
		// logs := "hello world"
		fmt.Println("No output file found")
		return errors.New(string(logs))
	}

	outputFile, err := os.Open(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	outputFileBytes, err := io.ReadAll(outputFile)
	if err != nil {
		return err
	}

	// testCaseUpdateResp, err := http.NewRequest(http.MethodPut, os.Getenv("SUBMISSION_SERVICE_URL")+"/submission/"+job.JobID+"/testcases", bytes.NewReader(outputFileBytes))
	// if err != nil {
	// 	return err
	// }
	// client := &http.Client{}
	// resp2, err := client.Do(testCaseUpdateResp)
	// if err != nil {
	// 	return err
	// }
	// defer resp2.Body.Close()

	testCaseUpdateResp, err := util.PerformPutRequest(os.Getenv("SUBMISSION_SERVICE_URL")+"/api/submission/"+job.JobID+"/testcases", outputFileBytes)
	if err != nil {
		return err
	}

	if testCaseUpdateResp.StatusCode != http.StatusCreated {
		return errors.New("failed to update test cases")
	}
	fmt.Println("Test cases updated")

	return nil

}

func DeleteFolderRecursive(folderPath string) error {
	err := os.RemoveAll(folderPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
