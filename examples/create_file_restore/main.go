package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/equalsgibson/slide"
)

/**
This example performs the following steps using the slide API:

- Checks for an API Token variable, and if not found, prompts the user for a token
- Queries the Slide API for a list of Agents
- Presents the list of Agents to the user, and requests the index of the agent that
	the restore should be created for
- Queries the Slide API for a list of snapshots for the provided Agent
- Creates a File Restore for the given Agent and Snapshot
- Provides the download URL for the File Restore.

This example could be extended to perform other checks (smart pagination, nicer TUI, check agents listed have valid snapshots etc.)
*/

func main() {

	buf := bufio.NewReader(os.Stdin)

	// * NOTE:
	// Ensure you set your environment variable before using this example, or enter it when prompted
	// export SLIDE_AUTH_TOKEN=xxxabc123
	slideAuthToken := os.Getenv("SLIDE_AUTH_TOKEN")
	if slideAuthToken == "" {
		fmt.Println("Did not detect 'SLIDE_AUTH_TOKEN' environment variable. Please provide this or set the environment variable and restart the example.")

		var err error
		slideAuthToken, err = buf.ReadString('\n')
		if err != nil {
			fmt.Printf("Encountered error while reading Slide Auth Token from os.Stdin: %s\n", err.Error())

			os.Exit(1)
		}
	}

	// * NOTE:
	// If you do not want to make actual network requests, include a custom roundtripper, similar to the example below
	// slideService := slide.NewService(strings.TrimSuffix(slideAuthToken, "\n"), slide.WithCustomRoundtripper(
	// 	roundtripper.MockNetworkQueue(
	// 		[]roundtripper.MockRoundTripFunc{
	// 			roundtripper.Serve(&roundtripper.MockResponseFile{
	// 				FilePath:   "./mock_network_responses/agent_list_page1_200.json",
	// 				StatusCode: http.StatusOK,
	// 			}),
	// 			roundtripper.Serve(&roundtripper.MockResponseFile{
	// 				FilePath:   "./mock_network_responses/agent_list_page2_200.json",
	// 				StatusCode: http.StatusOK,
	// 			}),
	// 			roundtripper.Serve(&roundtripper.MockResponseFile{
	// 				FilePath:   "./mock_network_responses/snapshot_list_page1_200.json",
	// 				StatusCode: http.StatusOK,
	// 			}),
	// 			roundtripper.Serve(&roundtripper.MockResponseFile{
	// 				FilePath:   "./mock_network_responses/snapshot_list_page2_200.json",
	// 				StatusCode: http.StatusOK,
	// 			}),
	// 			roundtripper.Serve(&roundtripper.MockResponseFile{
	// 				FilePath:   "./mock_network_responses/restore_file_create_201.json",
	// 				StatusCode: http.StatusCreated,
	// 			}),
	// 		},
	// 	),
	// ))

	slideService := slide.NewService(strings.TrimSuffix(slideAuthToken, "\n"))

	ctx := context.Background()

	fmt.Println("Querying Slide API for agents...")

	agents := []slide.Agent{}
	if err := slideService.Agents().List(ctx, func(response slide.ListResponse[slide.Agent]) error {
		agents = append(agents, response.Data...)

		return nil
	}); err != nil {
		fmt.Printf("Encountered error while querying agents from Slide API: %s\n", err.Error())

		os.Exit(1)
	}

	fmt.Printf("Type the number of the agent that you want to create a File Restore for and then press [ENTER]:\n\n")
	for index, agent := range agents {
		fmt.Printf("\t[%d] %s (ID: %s)\n", index, agent.DisplayName, agent.AgentID)
	}

	agentID, err := buf.ReadString('\n')
	if err != nil {
		fmt.Printf("Encountered error while reading Agent Index from os.Stdin: %s\n", err.Error())

		os.Exit(1)
	}

	agentIDInt, err := strconv.ParseInt(strings.TrimSuffix(agentID, "\n"), 10, 64)
	if err != nil {
		fmt.Printf("Encountered error while converting Agent Index to int: %s\n", err.Error())

		os.Exit(1)
	}

	if agentIDInt > int64(len(agents)) || agentIDInt < 0 {
		fmt.Printf("Agent index provided (%d) is out of range. It needs to be greater than or equal to 0, and less than or equal to %d\n", agentIDInt, len(agents))

		os.Exit(1)
	}

	snapshots := []slide.Snapshot{}
	if err := slideService.Snapshots().ListWithQueryParameters(
		ctx,
		func(response slide.ListResponse[slide.Snapshot]) error {
			snapshots = append(snapshots, response.Data...)

			return nil
		},
		slide.WithAgentID(agents[agentIDInt].AgentID),
	); err != nil {
		fmt.Printf("Encountered error while querying snapshots for agent (%s) from Slide API: %s\n", agents[agentIDInt].AgentID, err.Error())

		os.Exit(1)
	}

	fmt.Printf("Type the number of the snapshot that you want to create a File Restore for and then press [ENTER]:\n\n")
	for index, snapshot := range snapshots {
		fmt.Printf("\t[%d] %s (Started At: %s)\n", index, snapshot.SnapshotID, snapshot.BackupStartedAt)
	}

	snapshotID, err := buf.ReadString('\n')
	if err != nil {
		fmt.Printf("Encountered error while reading Agent Index from os.Stdin: %s\n", err.Error())

		os.Exit(1)
	}

	snapshotIDInt, err := strconv.ParseInt(strings.TrimSuffix(snapshotID, "\n"), 10, 64)
	if err != nil {
		fmt.Printf("Encountered error while converting Snapshot Index to int: %s\n", err.Error())

		os.Exit(1)
	}

	if snapshotIDInt > int64(len(snapshots)) || snapshotIDInt < 0 {
		fmt.Printf("Snapshot index provided (%d) is out of range. It needs to be greater than or equal to 0, and less than or equal to %d\n", snapshotIDInt, len(snapshots))

		os.Exit(1)
	}

	fmt.Printf("Creating a File Restore for the following Agent and Snapshot:\n\n\tAgent:\t%s\n\tSnapshot:\t%s\n\tDevice:\t%s\n", agents[agentIDInt].AgentID, snapshots[snapshotIDInt].SnapshotID, agents[agentIDInt].DeviceID)

	fileRestore, err := slideService.FileRestores().Create(ctx, slide.FileRestorePayload{
		DeviceID:   agents[agentIDInt].DeviceID,
		SnapshotID: snapshots[snapshotIDInt].SnapshotID,
	})

	if err != nil {
		fmt.Printf("Encountered error while creating File Restore: %s\n", err.Error())

		os.Exit(1)
	}

	fmt.Printf("\nFile Restore successfully created! Please see details below:\n\n")

	fileRestoreBytes, err := json.MarshalIndent(fileRestore, "", "	")
	if err != nil {
		fmt.Printf("Encountered error while formatting File Restore details: %s\n", err.Error())

		os.Exit(1)
	}

	fmt.Println(string(fileRestoreBytes))

	fmt.Printf("\nGoodbye!\n")
}
