# Python Requirements Generator

The Python Requirements Generator is a Go-based command-line tool that automates the creation of `requirements.txt` files for your Python projects. It scans your project, identifies imported modules, and then matches those modules with the exact versions of packages installed in your current environment using `pip freeze`.

---
## 🚀 Features

* **Recursive scanning**: Automatically finds all `.py` files in a directory and its subdirectories.
* **Smart import detection**: Extracts both `import module` and `from module import` statements.
* **Version matching**: Matches detected modules with installed package versions using `pip freeze`.
* **Flexible output**: Customize the output file name and location.
* **Cross-platform**: Works on Windows, macOS, and Linux.
* **Error handling**: Provides clear warnings and error messages.
* **Fast execution**: Written in Go for optimal performance.

---
## 📋 Prerequisites

* **Go 1.16 or higher**
* **Python with pip** installed and accessible from your command line.
* **Python packages** installed in the current environment that you want to generate requirements for.

---
## 🛠️ Installation

You have a couple of options to get the tool set up:

### Option 1: Build from Source

```bash
# Clone or download the source code
git clone https://github.com/LaamiriOuail/go-pyreqs.git
cd go-pyreqs

# Build the executable
go build -o py-requirements-gen main.go
```

### Option 2: Direct Compilation

If you have the `main.go` file directly, you can compile it:

```go
go build -o py-requirements-gen main.go
```

---
## 📖 Usage

### Basic Usage

```bash
# Scan the current directory and generate requirements.txt
./py-requirements-gen

# Scan a specific directory
./py-requirements-gen /path/to/your/python/project
```

### Advanced Usage

```bash
# Custom output file
./py-requirements-gen --output my-requirements.txt

# Scan a specific directory with custom output
./py-requirements-gen --output deps.txt /path/to/project

# Show help message
./py-requirements-gen -h
```

### Command Line Options

| Option      | Description                    | Default            |
| :---------- | :----------------------------- | :----------------- |
| `--output`  | Specify the output file name   | `requirements.txt` |
| `-h`        | Show help message              | -                  |

---
## 📁 Example

Given a Python project with the following structure:

```plaintext
my-project/
├── main.py
├── utils/
│   ├── helper.py
│   └── api.py
└── tests/
└── test_main.py
```

Where `main.py` contains:

```python
import requests
import pandas as pd
from flask import Flask
import numpy as np
```

Running:

```bash
./py-requirements-gen my-project/
```

Will generate `requirements.txt` with content similar to this (versions depend on your installed packages):

Flask==2.3.2
numpy==1.24.3
pandas==2.0.2
requests==2.31.0


---
## 🔍 How It Works

1.  **Directory Scanning**: Recursively walks through the target directory to find all `.py` files.
2.  **Import Extraction**: Uses regex patterns to identify `import module_name` and `from module_name import something` statements.
3.  **Module Normalization**: Handles package name variations (e.g., hyphens vs. underscores, case differences) for accurate matching.
4.  **Version Matching**: Executes `pip freeze` to get a list of all installed Python package versions in the current environment.
5.  **Requirements Generation**: Matches the detected modules from your code with the installed packages and outputs them in a sorted `requirements.txt` format.

---
## ⚠️ Limitations

* **Regex-based parsing**: The tool uses regex instead of Abstract Syntax Tree (AST) parsing, which means it might miss complex or dynamic import patterns.
* **`pip` dependency**: Requires `pip` to be available in your system's `PATH`.
* **Environment-specific**: Only detects packages installed in the current Python environment where the tool is run.
* **No version constraints**: Only outputs exact versions (`==`), not ranges or minimum versions.

---
## 🐛 Troubleshooting

### Common Issues

* **Error: "pip command not found"**
    * Ensure Python and `pip` are installed and accessible from your command line.
    * Try running `python -m pip freeze` manually to verify `pip` works.

* **Warning: "Could not parse file.py"**
    * The file may have syntax errors or use complex import patterns that the regex cannot handle.
    * The tool will continue processing other files even if this warning appears.

* **Empty `requirements.txt`**
    * No matching packages were found between your imports and your installed packages.
    * Verify that the packages are actually installed in your current environment by running `pip freeze`.
    * Check if the import names in your Python files correctly match the package names.

### Debug Tips

* **Verify `pip` works**: Run `pip freeze` manually in your terminal to see what packages are installed.
* **Check Python files**: Ensure your `.py` files contain standard `import` statements.
* **Environment check**: Make sure you're running the tool in the correct Python environment (e.g., virtual environment, Conda environment) where your project's dependencies are installed.

---
## 🔧 Development

### Building

```bash
# Build for current platform
go build -o py-requirements-gen main.go

# Build for different platforms
GOOS=windows GOARCH=amd64 go build -o py-requirements-gen.exe main.go
GOOS=linux GOARCH=amd64 go build -o py-requirements-gen-linux main.go
GOOS=darwin GOARCH=amd64 go build -o py-requirements-gen-mac main.go
```

### Testing

```bash
# Test with a sample Python project
mkdir test-project
echo "import requests" > test-project/main.py
echo "from flask import Flask" > test-project/app.py
./py-requirements-gen test-project/
```

---
## 📝 License

This project is licensed under the MIT License - see the `LICENSE` file for details.