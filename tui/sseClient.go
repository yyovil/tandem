package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func CreateSSEClient(endpoint string) {
	runRequest := RunRequest{
		SessionId: "string",
		UserId:    "string",
		Model:     ModelGemini25FlashPreview0417,
		Stream:    true,
		Message:   "update me about meta's anti-trust trials",
	}

	serialisedReqBody, err := json.Marshal(runRequest)
	if err != nil {
		log.Println("Couldn't seralise the body: ", err.Error())
	}

	// i'm not gonna need cancel unless i stitch it with the tui.
	ctx, _ := context.WithCancel(context.Background())

	if endpoint == "" {
		log.Println("ENDPOINT var is not defined or exported in the current env.")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(serialisedReqBody))

	if err != nil {
		log.Println("Couldn't create a new HTTP POST req with ctx: ", err.Error())
	}

	client := http.Client{}
	fmt.Println("making the request.")
	res, err := client.Do(req)
	if err != nil {
		log.Println("Couldn't send a HTTP POST req: ", err.Error())
	}
	defer res.Body.Close()

	scanner := bufio.NewScanner(req.Body)
	// read until io.EOF is raised.
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
}

/*
TODO: create an SSE client.
1. read the whole response till io.EOF.
2. be able to stop streaming from nowhere in the middle.
*/
