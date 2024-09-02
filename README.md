
# Ajtacka Grilovacky Demo

This project is a Go-based web application that provides functionalities related to managing and tracking drink consumption (possibly alcohol-related), including statistical analysis. The project utilizes websockets for real-time updates and follows a modular structure.

## Table of Contents

1. [Introduction](#introduction)
2. [Project Structure](#project-structure)
3. [Requirements](#requirements)
4. [Installation](#installation)
5. [Usage](#usage)
6. [Contributing](#contributing)
7. [License](#license)

## Introduction

Ajtacka Grilovacky Demo is designed to help track the consumption of drinks at gatherings, providing users with real-time statistics and updates. The application is structured in a modular way to ensure scalability and ease of maintenance.

## Project Structure

The project is structured as follows:

- **main.go**: The entry point of the application.
- **go.mod & go.sum**: Go modules files for dependency management.
- **docker/**: Contains Docker-related configurations.
- **internal/**: This directory holds the core business logic of the application.
  - **handlers/**: Contains the HTTP handlers for different routes.
    - `alcoholik.go`: Handler related to alcohol consumption.
    - `drink.go`: Handler for general drink management.
    - `stats.go`: Handler for statistical data.
  - **models/**: Contains the data models used in the application.
    - `alkoholik.go`: Model for alcohol-related data.
    - `db.go`: Database-related models and operations.
    - `qr.go`: QR code related models.
- **pkg/**: Contains packages that can be reused in other projects.
  - **websocket/**: Implements websocket functionality.
    - `client.go`: Manages individual websocket connections.
    - `hub.go`: Manages the central hub for websocket communication.
- **web/**: Contains the frontend part of the application.
  - **templates/**: HTML templates for rendering the UI.
    - `alcoholik.html`: Template for alcohol-related data display.
    - `drink.html`: Template for general drink management.
    - `embed.go`: Go code to embed templates into the binary.
    - `layout.html`: Base layout template.
    - `stats.html`: Template for statistics display.

## Requirements

- Go 1.16 or later
- Docker (optional, for containerized deployment)

## Installation

To set up the project locally, follow these steps:

1. Clone the repository:
    ```bash
    git clone <repository-url>
    cd ajtacka-grilovacky-demo
    ```

2. Install the dependencies:
    ```bash
    go mod tidy
    ```

3. Run the application:
    ```bash
    go run main.go
    ```

4. Open your browser and navigate to `http://localhost:8080` to see the application in action.

## Usage

The application provides several functionalities:

- **Track Alcohol Consumption**: Users can track their alcohol intake in real-time.
- **Drink Management**: Add or remove drinks from the list.
- **Statistics**: View real-time statistics based on the data entered.
- **Websocket Communication**: The application supports real-time updates using websockets.

## Contributing

Contributions are welcome! If you find any bugs or have suggestions for improvements, feel free to create an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.
