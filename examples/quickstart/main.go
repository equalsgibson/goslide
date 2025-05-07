package main

import (
	"context"
	"fmt"
	"os"

	"github.com/equalsgibson/goslide"
)

/**
This example performs the following steps using the slide API:

- Checks for an API Token variable, and if not found, requests the user set this and then exits.
- Queries the Slide API for a list of Devices
- Presents the list of Devices to the user

This example could be extended to perform other checks (smart pagination, nicer TUI, check if devices listed have any Alerts or warnings etc.)
*/

func main() {
	// * NOTE:
	// Ensure you set your environment variable before using this example, or enter it when prompted
	// export SLIDE_AUTH_TOKEN=xxxabc123
	slideAuthToken := os.Getenv("SLIDE_AUTH_TOKEN")
	if slideAuthToken == "" {
		fmt.Println("Did not detect SLIDE_AUTH_TOKEN environment variable. Please set this, and then re-run the example.")

		os.Exit(1)
	}

	// * NOTE:
	// If you do not want to make actual network requests, include a custom roundtripper, similar to the example below
	// slideService := goslide.NewService(strings.TrimSuffix(slideAuthToken, "\n"), goslide.WithCustomRoundtripper(
	// 	roundtripper.MockNetworkQueue(
	// 		[]roundtripper.MockRoundTripFunc{
	// 			roundtripper.Serve(&roundtripper.MockResponseFile{
	// 				FilePath:   "./mock_network_responses/device_list_page1_200.json",
	// 				StatusCode: http.StatusOK,
	// 			}),
	// 			roundtripper.Serve(&roundtripper.MockResponseFile{
	// 				FilePath:   "./mock_network_responses/device_list_page2_200.json",
	// 				StatusCode: http.StatusOK,
	// 			}),
	// 		},
	// 	),
	// ))

	// Create the slide service by calling goslide.NewService
	slideService := goslide.NewService(slideAuthToken)

	fmt.Println("Querying Slide API for devices...")

	ctx := context.Background()

	slideDevices := []goslide.Device{}
	if err := slideService.Devices().List(ctx, func(response goslide.ListResponse[goslide.Device]) error {
		slideDevices = append(slideDevices, response.Data...)

		return nil
	}); err != nil {
		fmt.Printf("Encountered error while querying devices from Slide API: %s\n", err.Error())

		os.Exit(1)
	}

	fmt.Printf("Found %d devices\n", len(slideDevices))
	fmt.Printf("Slide Device details:\n")
	for _, device := range slideDevices {
		fmt.Printf("\t%s - %s\n", device.DeviceID, device.DisplayName)
	}

	fmt.Printf("\nGoodbye!\n")
}
