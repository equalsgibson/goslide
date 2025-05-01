<!-- markdownlint-configure-file { "MD004": { "style": "consistent" } } -->
<!-- markdownlint-disable MD033 -->

#

<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://equalsgibson.github.io/slide/assets/goslide-dark.png">
    <source media="(prefers-color-scheme: light)" srcset="https://equalsgibson.github.io/slide/assets/goslide-light.png">
    <img src="https://equalsgibson.github.io/slide/assets/goslide-light.png" width="400" height="175" alt="Go Slide repository logo">
  </picture>
  <br>
  <strong>Interact with your Slide devices, agents and restores!</strong>
</p>

<!-- markdownlint-enable MD033 -->

<div align="right">

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
  - Open a terminal and navigate to the directory of your project, then run: `go get github.com/equalsgibson/slide@latest`

### Quickstart
After following the above steps, you could create a simple `main.go` file and include the following to list all your current Slide Devices:
  
![Quickstart Demo](/ops/assets/quickstart.gif)

```golang
...
	// Create the slide service by calling slide.NewService
	slideService := slide.NewService(slideAuthToken)

	fmt.Println("Querying Slide API for devices...")

	ctx := context.Background()
	slideDevices := []slide.Device{}
	if err := slideService.Devices().List(ctx, func(response slide.ListResponse[slide.Device]) error {
		slideDevices = append(slideDevices, response.Data...)

		return nil
	}); err != nil {
		fmt.Printf("Encountered error while querying devices from Slide API: %s\n", err.Error())

		os.Exit(1)
	}
...
```

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

[golang]: https://img.shields.io/badge/v1.23-000?logo=go&logoColor=fff&labelColor=444&color=%2300ADD8
[golang-url]: https://go.dev/
[coverage]: https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fequalsgibson.github.io%2Fslide%2Fcoverage.json&query=%24.total&label=Coverage
[coverage-url]: https://equalsgibson.github.io/slide/coverage.html
[goaction]: https://github.com/equalsgibson/slide/actions/workflows/go.yml/badge.svg?branch=main
[goaction-url]: https://github.com/equalsgibson/slide/actions/workflows/go.yml
[goref]: https://pkg.go.dev/badge/github.com/equalsgibson/slide.svg
[goref-url]: https://pkg.go.dev/github.com/equalsgibson/slide
[goreport]: https://goreportcard.com/badge/github.com/equalsgibson/slide
[goreport-url]: https://goreportcard.com/report/github.com/equalsgibson/slide
