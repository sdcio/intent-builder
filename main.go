package main

import (
	"context"
	"os"
	"time"

	"github.com/sdcio/data-server/pkg/config"
	schemaClient "github.com/sdcio/data-server/pkg/datastore/clients/schema"
	"github.com/sdcio/data-server/pkg/schema"
	"github.com/sdcio/data-server/pkg/tree"
	"github.com/sdcio/intent-builder/pkg/clients"
	"github.com/sdcio/intent-builder/pkg/data"
	"github.com/sdcio/intent-builder/pkg/prompt"
	"github.com/sdcio/intent-builder/pkg/prompt/prompts"
	"github.com/sirupsen/logrus"
)

func main() {
	err := run()
	if err != nil {
		logrus.Panic(err)
	}
}

func run() error {
	ctx := context.Background()

	ds, err := clients.NewDataServerClient(os.Args[1])
	if err != nil {
		return err
	}

	// setup grpc connection
	conn, err := clients.NewGrpcClient(os.Args[1])
	if err != nil {
		return err
	}

	scb := schemaClient.NewSchemaClientBound(&config.SchemaConfig{Vendor: "srl.nokia.sdcio.dev", Version: "25.7.1"}, schema.NewRemoteClient(conn, &config.RemoteSchemaCache{Capacity: 1, WithDescription: true, TTL: time.Second}))

	ssc, err := clients.NewSchemaServerClient(scb)
	if err != nil {
		return err
	}

	backend, err := data.NewBackend(ctx, ds, ssc, tree.NewTreeContext(scb, "New"))
	if err != nil {
		return err
	}

	pc := prompt.NewPrompter(ctx)

	pc.AddPrompt(prompts.NewPromptsExit(), [][]string{{"quit"}, {"exit"}})
	pc.AddPrompt(data.NewDataPrompt(backend), [][]string{{"set"}, {"delete"}, {"show"}})

	pc.Run()
	return nil
}
