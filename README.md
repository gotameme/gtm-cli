# Go! Tame! Me! CLI (gtm-cli)

Welcome to the official command-line interface (CLI) for Go! Tame! Me!, an ambitious simulation project inspired by the
classic game Ant!Me!. The `gtm-cli` tool is designed to streamline the development, testing, and execution of custom ant
behavior simulations. This CLI leverages the core functionalities of Go! Tame! Me!, allowing users to compile and run
simulation plugins directly from the command line.

## Features

- **Compile and Run Simulation Plugins:** Quickly compile and run your custom Go! Tame! Me! simulation plugins with a
  simple command.
- **Customizable Simulation Parameters:** Launch simulations with customizable settings such as immediate start,
  headless mode, and desired sugar cone count.
- **Open-Source Development:** Built with the open-source community in mind, `gtm-cli` supports collaborative
  development and experimentation.

## Installation

Before you can use `gtm-cli`, make sure you have Go installed on your machine (version 1.22 or newer is recommended).

To install `gtm-cli`, run the following command:

```bash
go install github.com/gotameme/gtm-cli@latest
```

This command downloads and installs the `gtm-cli` tool, making it ready to use.

## Usage

`gtm-cli` currently supports two primary commands:

### 1. Root Command

The base command for `gtm-cli`, which can be executed with just `gtm-cli`, provides a brief description of the
application and available commands.

### 2. Run Command

The `run` command compiles and executes a Go! Tame! Me! simulation plugin.

```bash
gtm-cli run /path/to/your/plugin.go
```

#### Flags

- `--startImmediately` or `-i`: Start the game immediately after the plugin has been compiled.
- `--headless`: Run the game in headless mode for simulations without a graphical interface.
- `--sugar`: Set the desired amount of sugar cones in the simulation environment.

## Contributing

Contributions are welcome! Whether you're fixing bugs, adding new features, or improving documentation, your help is
invaluable. For more details on how to contribute, please refer to our `CONTRIBUTING.md` file.

## Support and Communication

For questions, suggestions, and discussions, please use the GitHub Issue Tracker. We strive for a supportive and
harassment-free communication environment. English is the preferred language for project communication to include as
many contributors as possible.

Thank you for your interest in contributing to Go! Tame! Me! CLI. Together, we're building an engaging and educational
simulation platform.

## License

`gtm-cli` is licensed under the MIT License. See the `LICENSE` file for more details.
