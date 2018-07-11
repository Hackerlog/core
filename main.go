package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/akamensky/argparse"
	raven "github.com/getsentry/raven-go"
)

var (
	isProd = os.Getenv("ENV") == "production"

	// Use this for build information
	version = ""
	commit  = ""
	date    = ""

	xHeader = "X-Hackerlog-EditorToken"

	// All of the params that our CLI can accept
	pVersion     = []string{"v", "version"}
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

func init() {
	if isProd {
		raven.SetDSN("https://dd748bbffcf04dca9bfb64b312efbcd8:90583a39dbe344678fbc813dd138af56@sentry.io/1240770")
		raven.SetRelease(commit)
	}
}

func main() {
	parser := argparse.NewParser("hackerlog", "Collects coding stats and submits them to the API.")

	// Print out the current version
	printVersion := parser.Flag(pVersion[0], pVersion[1], &argparse.Options{
		Required: false,
		Help:     "Prints the current version of the executable",
	})

	// This allows for running the command like: ./core send --foo bar --baz true, etc.
	send := parser.NewCommand("send", "Sends a unit of work to the API")

	apiUrl := send.String(pAPIURL[0], pAPIURL[1], &argparse.Options{
		Required: true,
		Help:     "The URL of the API to send the request",
	})

	editorToken := send.String(pEditorToken[0], pEditorToken[1], &argparse.Options{
		Required: true,
		Help:     "The editor token associated with a user",
	})

	editorType := send.String(pEditorType[0], pEditorType[1], &argparse.Options{
		Required: true,
		Help:     "The editor that is being used",
	})

	projectName := send.String(pProjectName[0], pProjectName[1], &argparse.Options{
		Required: true,
		Help:     "The name of the project associated with the unit of work.",
	})

	fileName := send.String(pFileName[0], pFileName[1], &argparse.Options{
		Required: true,
		Help:     "The file name that was edited",
	})

	locWritten := send.Int(pLocWritten[0], pLocWritten[1], &argparse.Options{
		Required: false,
		Help:     "The amount of lines of code that has been written.",
	})

	locDeleted := send.Int(pLocDeleted[0], pLocDeleted[1], &argparse.Options{
		Required: false,
		Help:     "The amount of lines of code that has been deleted.",
	})

	startedAt := send.String(pStartedAt[0], pStartedAt[1], &argparse.Options{
		Required: true,
		Help:     "When did the file start being edited",
	})

	stoppedAt := send.String(pStoppedAt[0], pStoppedAt[1], &argparse.Options{
		Required: true,
		Help:     "When did the file stop being edited",
	})

	operatingSystem := runtime.GOOS

	err := parser.Parse(os.Args)
	if err != nil {
		if isProd {
			raven.CaptureErrorAndWait(err, map[string]string{"args": strings.Join(os.Args, " | ")})
		}
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

	if *printVersion {
		fmt.Println("v" + version)
	} else {
		// TODO: Need to solve offline cases here
		if err := sendUnit(*apiUrl, unit, *editorToken); err != nil {
			if isProd {
				raven.CaptureErrorAndWait(err, map[string]string{
					"editorType": unit.EditorType,
					"os":         unit.Os,
					"arch":       runtime.GOARCH,
				})
			}
			// We need reason codes here so we know what to do next
			log.Fatal(err)
		}
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
