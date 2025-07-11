package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/a-novel-kit/configurator/chans"
	"github.com/a-novel-kit/configurator/utilstest"

	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/internal/api/apiclient"
	"github.com/a-novel/service-story-schematics/internal/api/codegen"
)

var logs *chans.MultiChan[string]

func getServerClient() (*codegen.Client, *apiclient.SecuritySource, error) {
	security := apiclient.NewSecuritySource()

	client, err := codegen.NewClient(fmt.Sprintf("http://127.0.0.1:%v/v1", config.API.Port), security)
	if err != nil {
		return nil, nil, fmt.Errorf("create client: %w", err)
	}

	start := time.Now()
	_, err = client.Ping(context.Background())

	for time.Since(start) < 16*time.Second && err != nil {
		_, err = client.Ping(context.Background())
	}

	if err != nil {
		return nil, nil, fmt.Errorf("ping server: %w", err)
	}

	return client, security, nil
}

func _patchSTD() {
	patchedStd, _, err := utilstest.MonkeyPatchStderr()
	if err != nil {
		panic(err)
	}

	logs, _, err = utilstest.CaptureSTD(patchedStd)
	if err != nil {
		panic(err)
	}

	go func() {
		listener := logs.Register()
		for msg := range listener {
			// Forward logs to default system outputs, in case we need them for debugging.
			log.Println(msg)
		}
	}()
}

// Create a separate database to run integration tests.
func init() {
	_patchSTD()

	go func() {
		main()
	}()

	_, _, err := getServerClient()
	if err != nil {
		panic(err)
	}
}
