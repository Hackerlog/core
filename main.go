package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/akamensky/argparse"
)

var (
	xHeader = "X-Hackerlog-EditorToken"

	// All of the params that our CLI can accept
	pAPIURL      = []string{"u", "api-url"}
	pEditorToken = []string{"t", "editor-token"}
	pEditorType  = []string{"e", "editor-type"}
	pProjectName = []string{"p", "project-name"}
	pFileName    = []string{"f", "file-name"}
	pLocWritten  = []string{"w", "loc-written"}
	pLocDeleted  = []string{"d", "loc-deleted"}
	pStartedAt   = []string{"s", "started-at"}
	pStoppedAt   = []string{"x", "stopped-at"}
)

// Unit This is the data type that we need to send to the API
type Unit struct {
	EditorType  string `json:"editor_type"`
	ProjectName string `json:"project_name"`
	FileName    string `json:"file_name"`
	LocWritten  int    `json:"loc_written"`
	LocDeleted  int    `json:"loc_deleted"`
	Os          string `json:"os"`
	StartedAt   string `json:"started_at"`
	StoppedAt   string `json:"stopped_at"`
}

func main() {
	parser := argparse.NewParser("hackerlog", "Collects coding stats and submits them to the API.")
	apiUrl := parser.String(pAPIURL[0], pAPIURL[1], &argparse.Options{
		Required: true,
		Help:     "The URL of the API to send the request",
	})
	editorToken := parser.String(pEditorToken[0], pEditorToken[1], &argparse.Options{
		Required: true,
		Help:     "The editor token associated with a user",
	})
	editorType := parser.String(pEditorType[0], pEditorType[1], &argparse.Options{
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
	startedAt := parser.String(pStartedAt[0], pStartedAt[1], &argparse.Options{
		Required: true,
		Help:     "When did the file start being edited",
	})
	stoppedAt := parser.String(pStoppedAt[0], pStoppedAt[1], &argparse.Options{
		Required: true,
		Help:     "When did the file stop being edited",
	})
	operatingSystem := runtime.GOOS

	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatalf("There is an error: %s", err)
	}

	unit := Unit{
		EditorType:  *editorType,
		ProjectName: *projectName,
		FileName:    *fileName,
		LocWritten:  *locWritten,
		LocDeleted:  *locDeleted,
		Os:          operatingSystem,
		StartedAt:   *startedAt,
		StoppedAt:   *stoppedAt,
	}

	// TODO: Need to solve offline cases here
	if err := sendUnit(*apiUrl, unit, *editorToken); err != nil {
		// We need reason codes here so we know what to do next
		log.Fatal(err)
	}
}

func sendUnit(url string, u Unit, et string) error {
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set(xHeader, et)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return errors.New("The request did not go through")
	}

	return nil
}
