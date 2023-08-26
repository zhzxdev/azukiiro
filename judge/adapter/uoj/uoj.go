package uoj

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/zhzxdev/azukiiro/azukiiro/client"
	"github.com/zhzxdev/azukiiro/azukiiro/common"
	"github.com/zhzxdev/azukiiro/azukiiro/storage"
)

type UojAdapter struct{}

func (u *UojAdapter) Name() string {
	return "uoj"
}

func Unzip(source string, target string) (string, error) {
	log.Printf("Unzipping %s to %s\n", source, target)
	dir, err := storage.MkdirTemp(target)
	if err != nil {
		return dir, err
	}
	err = exec.Command("unzip", source, "-d", dir).Run()
	if err != nil {
		os.RemoveAll(dir)
		return dir, err
	}
	// remove all symlinks to avoid security issues
	err = exec.Command("find", dir, "-type", "l", "-delete").Run()
	if err != nil {
		os.RemoveAll(dir)
		return dir, err
	}
	return dir, nil
}

type Test struct {
	Num    int    `xml:"num,attr"`
	Score  int    `xml:"score,attr"`
	Info   string `xml:"info,attr"`
	Time   int    `xml:"time,attr"`
	Memory int    `xml:"memory,attr"`
	In     string `xml:"in"`
	Out    string `xml:"out"`
	Res    string `xml:"res"`
}

type Details struct {
	Tests []Test `xml:"test"`
	Error string `xml:"error"`
}

type Result struct {
	XMLName xml.Name `xml:"result"`
	Score   int      `xml:"score"`
	Time    int      `xml:"time"`
	Memory  int      `xml:"memory"`
	Error   string   `xml:"error"`
	Details Details  `xml:"details"`
}

func toCodeBlock(v interface{}) string {
	return fmt.Sprintf("```\n%s\n```", v)
}

func ReadResult(resultPath string) (client.PatchSolutionTaskRequest, common.SolutionDetails, error) {
	// read result
	resultFile, err := os.ReadFile(resultPath)
	if err != nil {
		return client.PatchSolutionTaskRequest{
				Score: 0,
				Metrics: &map[string]float64{
					"cpu": 0,
					"mem": 0,
				},
				Status:  "Unknown Error",
				Message: "An unknown error occurred when reading the result",
			}, common.SolutionDetails{
				Jobs:    nil,
				Summary: "Unknown Error",
			}, err
	}

	// unmarshal XML
	var result Result
	if err := xml.Unmarshal(resultFile, &result); err != nil {
		return client.PatchSolutionTaskRequest{
				Score: 0,
				Metrics: &map[string]float64{
					"cpu": 0,
					"mem": 0,
				},
				Status:  "Unknown Error",
				Message: "An unknown error occurred when unmarshaling the result",
			}, common.SolutionDetails{
				Jobs:    nil,
				Summary: "Unknown Error",
			}, err
	}

	// Result -> common.SolutionDetails
	var testsResult []*common.SolutionDetailsTest
	status := "Accepted"
	if result.Error != "" {
		status = result.Error
	} else {
		for _, r := range result.Details.Tests {
			if status == "Accepted" && r.Info != "Accepted" {
				status = r.Info
			}
			testsResult = append(testsResult, &common.SolutionDetailsTest{
				Name:    "Test " + fmt.Sprint(r.Num),
				Score:   float64(r.Score),
				Status:  r.Info,
				Summary: "Time: `" + fmt.Sprint(r.Time) + "`\tMemory: `" + fmt.Sprint(r.Memory) + "`\n\nInput:\n\n" + toCodeBlock(r.In) + "\n\nOutput:\n\n" + toCodeBlock(r.Out) + "\n\nResult:\n\n" + toCodeBlock(r.Res),
			})
		}
	}
	return client.PatchSolutionTaskRequest{
			Score: result.Score,
			Metrics: &map[string]float64{
				"cpu": float64(result.Time),
				"mem": float64(result.Memory),
			},
			Status:  status,
			Message: "UOJ Judger finished with exit code 0",
		}, common.SolutionDetails{
			Jobs: []*common.SolutionDetailsJob{
				{
					Name:       "default",
					Score:      float64(result.Score),
					ScoreScale: 100,
					Status:     status,
					Tests:      testsResult,
					Summary:    "The default subtask",
				},
			},
			Summary: fmt.Sprintf("Error:\n\n%s", toCodeBlock(result.Details.Error)),
		}, nil
}

type UOJAdapterConfig struct {
	SandboxMode string `json:"sandbox_mode"`
}

type SolutionMetadata struct {
	Language string `json:"language"`
}

func (u *UojAdapter) Judge(ctx context.Context, config common.ProblemConfig, problemData string, solutionData string) error {
	adapterConfig := UOJAdapterConfig{
		SandboxMode: "bwrap",
	}
	json.Unmarshal([]byte(config.Judge.Config), &adapterConfig)

	judgerPath := "/opt/uoj_judger"

	// unzip data
	problemDir, err := Unzip(problemData, "problem")
	if err != nil {
		return err
	}
	defer os.RemoveAll(problemDir)
	solutionDir, err := Unzip(solutionData, "solution")
	if err != nil {
		return err
	}
	defer os.RemoveAll(solutionDir)

	language := "C++11"
	if content, err := os.ReadFile(solutionDir + "/.metadata.json"); err != nil {
		var metadata SolutionMetadata
		json.Unmarshal(content, &metadata)
		if metadata.Language == "C++" || metadata.Language == "C++11" {
			language = metadata.Language
		}
	}

	if err := os.WriteFile(solutionDir+"/submission.conf", []byte("answer_language "+language), 0666); err != nil {
		return err
	}

	// run judger
	var cmd *exec.Cmd
	switch adapterConfig.SandboxMode {
	case "none":
		cmd = exec.Command(judgerPath+"/main_judger", solutionDir, problemDir)
	case "bwrap":
		cmd = exec.Command("bwrap",
			"--dir", "/tmp",
			"--dir", "/var",
			"--bind", solutionDir, "/tmp/solution",
			"--ro-bind", problemDir, "/tmp/problem",
			"--ro-bind", judgerPath, "/opt/uoj_judger",
			"--bind", judgerPath+"/result", "/opt/uoj_judger/result",
			"--bind", judgerPath+"/work", "/opt/uoj_judger/work",
			"--ro-bind", "/usr", "/usr",
			"--symlink", "../tmp", "var/tmp",
			"--proc", "/proc",
			"--dev", "/dev",
			"--ro-bind", "/etc/resolv.conf", "/etc/resolv.conf",
			"--symlink", "usr/lib", "/lib",
			"--symlink", "usr/lib64", "/lib64",
			"--symlink", "usr/bin", "/bin",
			"--symlink", "usr/sbin", "/sbin",
			"--chdir", "/opt/uoj_judger/main_judger",
			"--unshare-all",
			"--die-with-parent",
			"/opt/uoj_judger/main_judger", "/tmp/solution", "/tmp/problem")
	default:
		return fmt.Errorf("unknown sandbox mode: %s", adapterConfig.SandboxMode)
	}
	cmd.Dir = judgerPath
	log.Printf("Running %s\n", cmd)
	if err := cmd.Run(); err != nil {
		return err
	}

	// read & report result
	result, resultDetails, _ := ReadResult(judgerPath + "/result/result.txt")
	client.PatchSolutionTask(ctx, &result)
	client.SaveSolutionDetails(ctx, &resultDetails)

	return nil
}
