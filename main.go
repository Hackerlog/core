package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/akamensky/argparse"
)

var (
	xHeader = "X-Hackerlog-EditorToken"

	// All of the params that our CLI can accept
	pAPIURL      = []string{"u", "api-url"}
	pEditor      = []string{"e", "editor"}
	pProjectName = []string{"p", "project-name"}
	pFileName    = []string{"f", "file-name"}
	pLocWritten  = []string{"w", "loc-written"}
	pLocDeleted  = []string{"d", "loc-deleted"}
	pStartedAt   = []string{"s", "started-at"}
	pStoppedAt   = []string{"x", "stopped-at"}
)

// Unit This is the data type that we need to send to the API
type Unit struct {
	EditorType  *string
	ProjectName *string
	FileName    *string
	LocWritten  *int
	LocDeleted  *int
	Os          *string
	StartedAt   *string
	StoppedAt   *string
}

func main() {
	parser := argparse.NewParser("hackerlog", "Collects coding stats and submits them to the API.")
	apiUrl := parser.String(pAPIURL[0], pAPIURL[1], &argparse.Options{
		Required: true,
		Help:     "The URL of the API to send the request",
	})
	editorType := parser.String(pEditor[0], pEditor[1], &argparse.Options{
		Required: true,
		Help:     "The editor that is being used",
	})
	projectName := parser.String(pProjectName[0], pProjectName[1], &argparse.Options{
		Required: true,
		Help:     "The name of the project associated with the unit of work.",
	})
	fileName := parser.String(pFileName[0], pFileName[1], &argparse.Options{
		Required: true,
		Help:     "The file name that was edited",
	})
	locWritten := parser.Int(pLocWritten[0], pLocWritten[1], &argparse.Options{
		Required: false,
		Help:     "The amount of lines of code that has been written.",
	})
	locDeleted := parser.Int(pLocDeleted[0], pLocDeleted[1], &argparse.Options{
		Required: false,
		Help:     "The amount of lines of code that has been deleted.",
	})
	operatingSystem := runtime.GOOS
	startedAt := parser.String(pStartedAt[0], pStartedAt[1], &argparse.Options{
		Required: true,
		Help:     "When did the file start being edited",
	})
	stoppedAt := parser.String(pStoppedAt[0], pStoppedAt[1], &argparse.Options{
		Required: true,
		Help:     "When did the file stop being edited",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Errorf("There is an error: %s", err)
	}

	unit := Unit{
		EditorType:  editorType,
		ProjectName: projectName,
		FileName:    fileName,
		LocWritten:  locWritten,
		LocDeleted:  locDeleted,
		Os:          &operatingSystem,
		StartedAt:   startedAt,
		StoppedAt:   stoppedAt,
	}

	if err := sendUnit(apiUrl, unit); err != nil {
		// We need reason codes here so we know what to do next
	}
}

func sendUnit(url *string, u Unit) error {
	// TODO: Finish making HTTP request...
	resp, err := http.Post(url)
}
