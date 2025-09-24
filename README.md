# SNetT: Secure/Simple Net Tools

SNetT is a secure, cross-platform desktop application for managing files over a network built on top of [SNetT Engine](https://github.com/owbird/SNetT-Engine).

## Features

- **File Sharing**: Securely share files between devices using a simple code. Based on [Magic Wormhole](https://www.lothar.com/~warner/MagicWormhole-PyCon2016.pdf)
- **File Server**: Host a directory of files with a web-based UI for browsing, downloading, and uploading.

## Getting Started

### Prerequisites

- Go 1.22 or later
- Fyne dependencies (see [Fyne's documentation](https://developer.fyne.io/started/))

### Installation & Running

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/Owbird/snett.git
    cd snett
    ```

2.  **Update the engine submodule:**
    ```sh
    git submodule update --init --recursive
    ```

3.  **Run the application:**
    ```sh
    go run .
    ```

For development with live reload, you can use `air`:
```sh
make dev
```

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
