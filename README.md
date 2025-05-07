<!-- markdownlint-configure-file { "MD004": { "style": "consistent" } } -->
<!-- markdownlint-disable MD033 -->

#

<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://equalsgibson.github.io/goslide/assets/goslide-dark.png">
    <source media="(prefers-color-scheme: light)" srcset="https://equalsgibson.github.io/goslide/assets/goslide-light.png">
    <img src="https://equalsgibson.github.io/goslide/assets/goslide-light.png" width="400" height="175" alt="Go Slide repository logo">
  </picture>
  <br>
  <strong>Interact with your Slide devices, agents and restores!</strong>
</p>

<!-- markdownlint-enable MD033 -->

<div align="right">

[![Slide API Version][slideapiversion]][slideapiversion-url]
[![Go][golang]][golang-url]
[![Code Coverage][coverage]][coverage-url]
[![Go Reference][goref]][goref-url]
[![Go Report Card][goreport]][goreport-url]

</div>

## Getting Started

This is a library that allows developers to interact with the publicly accessible REST API for [Slide](https://slide.tech), using Go. 

> **Info**  
> Slide is a modern backup solution, purpose built for MSPs (Managed Service Providers)

For a full API reference, see the official [Slide API Documentation](https://docs.slide.tech/)

### Prerequisites

  - Download and install Go, version 1.23+, from the [official Go website](https://go.dev/doc/install).
  - Obtain your Slide API Token and review the **[Authentication](https://docs.slide.tech/api/#description/authentication)** section.
  - Open your editor of choice, create a new directory and initialize your Go project.
  - Open a terminal and navigate to the directory of your project, then run: `go get github.com/equalsgibson/goslide@latest`

### Quickstart
After following the above steps, you could create a simple `main.go` file and include the following to list all your current Slide Devices:


```golang
...
	// Create the slide service by calling goslide.NewService
	slide := goslide.NewService(slideAuthToken)

	fmt.Println("Querying Slide API for devices...")

	ctx := context.Background()
	slideDevices := []goslide.Device{}
	if err := slide.Devices().List(ctx, func(response goslide.ListResponse[goslide.Device]) error {
		slideDevices = append(slideDevices, response.Data...)

		return nil
	}); err != nil {
		fmt.Printf("Encountered error while querying devices from Slide API: %s\n", err.Error())

		os.Exit(1)
	}
...
```

![Quickstart Demo](/ops/docs/assets/quickstart.gif)

### More Examples

To see examples on how this library could be used to create a basic CLI tool, checkout the [examples](/examples/) directory. 

<!-- CONTRIBUTING -->

## Contributing

Contributions are what make the open source community such an amazing place to learn, get inspired, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<!-- LICENSE -->

## License

Distributed under the GPL-3.0 License. See [LICENSE](/LICENSE) for more information.

<!-- CONTACT -->

## Contact

[Chris Gibson (@equalsgibson)](https://github.com/equalsgibson)

Project Link: [https://github.com/equalsgibson/slide](https://github.com/equalsgibson/slide)


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[slideapiversion]: https://img.shields.io/badge/v1.10.2-000?labelColor=000&color=001021&logo=data:image/svg+xml;base64,PHN2ZyBpZD0ibG9nbyIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB2aWV3Qm94PSIwIDAgNjkuNDc5IDY5LjQ3OSI+DQogICAgICAgICAgICAgICAgICA8c3R5bGU+DQogICAgICAgICAgICAgICAgICAgICAgLmRvdCB7DQogICAgICAgICAgICAgICAgICAgICAgICAgIGZpbGw6ICNmZmY7DQogICAgICAgICAgICAgICAgICAgICAgfQ0KICAgICAgICAgICAgICAgICAgICAgIC5kb3QtbyB7DQogICAgICAgICAgICAgICAgICAgICAgICAgIG9wYWNpdHk6IDAuMjc4MDU3Ow0KICAgICAgICAgICAgICAgICAgICAgIH0NCiAgICAgICAgICAgICAgICAgIDwvc3R5bGU+DQogICAgICAgICAgICAgICAgICA8ZyB0cmFuc2Zvcm09InJvdGF0ZSgxNDcuMzc3IDExMi45OTIgNDMuMTE1KXNjYWxlKDEuOTkzNikiIHN0eWxlPSJkaXNwbGF5OiBpbmxpbmUiPg0KICAgICAgICAgICAgICAgICAgICAgIDxyZWN0IGlkPSJkb3Qtby0xIiByeT0iMy41ODEiIHk9IjQyLjc0NSIgeD0iNzUuNTAzIiBoZWlnaHQ9IjcuMTYzIiB3aWR0aD0iNy4xOTMiIGNsYXNzPSJkb3QgZG90LW8iPjwvcmVjdD4NCiAgICAgICAgICAgICAgICAgICAgICA8cmVjdCBpZD0iZG90LW8tMiIgcnk9IjMuNTgxIiB5PSIzNC4zODgiIHg9Ijc1LjUwMyIgaGVpZ2h0PSI3LjE2MyIgd2lkdGg9IjcuMTkzIiBjbGFzcz0iZG90IGRvdC1vIj48L3JlY3Q+DQogICAgICAgICAgICAgICAgICAgICAgPHJlY3QgaWQ9ImRvdC1vLTMiIHJ5PSIzLjU4MSIgeT0iNDIuNzQ1IiB4PSI5Mi4yMTciIGhlaWdodD0iNy4xNjMiIHdpZHRoPSI3LjE5MyIgY2xhc3M9ImRvdCBkb3QtbyI+PC9yZWN0Pg0KICAgICAgICAgICAgICAgICAgICAgIDxyZWN0IGlkPSJkb3Qtby00IiByeT0iMy41ODEiIHk9IjUxLjEwMiIgeD0iOTIuMjQ2IiBoZWlnaHQ9IjcuMTYzIiB3aWR0aD0iNy4xOTMiIGNsYXNzPSJkb3QgZG90LW8iPjwvcmVjdD4NCiAgICAgICAgICAgICAgICAgICAgICA8cmVjdCBpZD0iZG90LXMtMSIgcnk9IjMuNTgxIiB5PSI1MS4xMDIiIHg9Ijc1LjUwMyIgaGVpZ2h0PSI3LjE2MyIgd2lkdGg9IjcuMTkzIiBjbGFzcz0iZG90IGRvdC1zIj48L3JlY3Q+DQogICAgICAgICAgICAgICAgICAgICAgPHJlY3QgaWQ9ImRvdC1zLTIiIHJ5PSIzLjU4MSIgeT0iNTEuMTAyIiB4PSI4My44NiIgaGVpZ2h0PSI3LjE2MyIgd2lkdGg9IjcuMTkzIiBjbGFzcz0iZG90IGRvdC1zIj48L3JlY3Q+DQogICAgICAgICAgICAgICAgICAgICAgPHJlY3QgaWQ9ImRvdC1zLTMiIHJ5PSIzLjU4MSIgeT0iNDIuNzQ1IiB4PSI4My44OSIgaGVpZ2h0PSI3LjE2MyIgd2lkdGg9IjcuMTkzIiBjbGFzcz0iZG90IGRvdC1zIj48L3JlY3Q+DQogICAgICAgICAgICAgICAgICAgICAgPHJlY3QgaWQ9ImRvdC1zLTQiIHJ5PSIzLjU4MSIgeT0iMzQuMzg4IiB4PSI4My44OSIgaGVpZ2h0PSI3LjE2MyIgd2lkdGg9IjcuMTkzIiBjbGFzcz0iZG90IGRvdC1zIj48L3JlY3Q+DQogICAgICAgICAgICAgICAgICAgICAgPHJlY3QgaWQ9ImRvdC1zLTUiIHJ5PSIzLjU4MSIgeT0iMzQuMzg4IiB4PSI5Mi4yNDYiIGhlaWdodD0iNy4xNjMiIHdpZHRoPSI3LjE5MyIgY2xhc3M9ImRvdCBkb3QtcyI+PC9yZWN0Pg0KICAgICAgICAgICAgICAgICAgPC9nPg0KICAgICAgICAgICAgICA8L3N2Zz4=
[slideapiversion-url]: https://docs.slide.tech/releases/#api-v1102
[golang]: https://img.shields.io/badge/v1.23-000?logo=go&logoColor=fff&labelColor=444&color=%2300ADD8
[golang-url]: https://go.dev/
[coverage]: https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fequalsgibson.github.io%2Fgoslide%2Fcoverage.json&query=%24.total&label=Coverage
[coverage-url]: https://equalsgibson.github.io/goslide/coverage.html
[goaction]: https://github.com/equalsgibson/goslide/actions/workflows/go.yml/badge.svg?branch=main
[goaction-url]: https://github.com/equalsgibson/goslide/actions/workflows/go.yml
[goref]: https://pkg.go.dev/badge/github.com/equalsgibson/goslide.svg
[goref-url]: https://pkg.go.dev/github.com/equalsgibson/goslide
[goreport]: https://goreportcard.com/badge/github.com/equalsgibson/goslide
[goreport-url]: https://goreportcard.com/report/github.com/equalsgibson/goslide
