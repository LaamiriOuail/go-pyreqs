package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type RequirementsGenerator struct {
	targetDir    string
	outputFile   string
	foundModules map[string]bool
}

func main() {
	var outputFile string
	flag.StringVar(&outputFile, "output", "requirements.txt", "Output file for requirements")
	flag.Parse()

	// Get target directory (default to current directory)
	targetDir := "."
	if flag.NArg() > 0 {
		targetDir = flag.Arg(0)
	}

	generator := &RequirementsGenerator{
		targetDir:    targetDir,
		outputFile:   outputFile,
		foundModules: make(map[string]bool),
	}

	if err := generator.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (rg *RequirementsGenerator) run() error {
	// Check if target directory exists
	if _, err := os.Stat(rg.targetDir); os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' not found", rg.targetDir)
	}

	fmt.Printf("Scanning directory '%s' for Python files...\n", rg.targetDir)

	// Find and process all Python files
	if err := rg.findAndProcessPythonFiles(); err != nil {
		return fmt.Errorf("failed to process Python files: %v", err)
	}

	// Get installed packages
	installedPackages, err := rg.getInstalledPackages()
	if err != nil {
		return fmt.Errorf("failed to get installed packages: %v", err)
	}

	// Generate requirements
	requirements := rg.generateRequirements(installedPackages)

	// Write to output file
	if err := rg.writeRequirements(requirements); err != nil {
		return fmt.Errorf("failed to write requirements: %v", err)
	}

	rg.printResults(requirements)
	return nil
}

func (rg *RequirementsGenerator) findAndProcessPythonFiles() error {
	return filepath.Walk(rg.targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".py") {
			if err := rg.extractModulesFromFile(path); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Could not parse %s: %v\n", path, err)
			}
		}
		return nil
	})
}

func (rg *RequirementsGenerator) extractModulesFromFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Parse Python imports using regex (since we're in Go, we can't use Python's ast)
	imports := rg.extractImportsFromPythonCode(string(content))
	
	for _, module := range imports {
		rg.foundModules[module] = true
	}

	return nil
}

func (rg *RequirementsGenerator) extractImportsFromPythonCode(content string) []string {
	var modules []string
	
	// Regex patterns for Python imports
	importRegex := regexp.MustCompile(`(?m)^import\s+([a-zA-Z_][a-zA-Z0-9_]*(?:\.[a-zA-Z_][a-zA-Z0-9_]*)*)`)
	fromImportRegex := regexp.MustCompile(`(?m)^from\s+([a-zA-Z_][a-zA-Z0-9_]*(?:\.[a-zA-Z_][a-zA-Z0-9_]*)*)\s+import`)
	
	// Find "import module" statements
	matches := importRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			// Get top-level module (e.g., "requests" from "requests.auth")
			topLevel := strings.Split(match[1], ".")[0]
			modules = append(modules, topLevel)
		}
	}
	
	// Find "from module import" statements
	matches = fromImportRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			// Get top-level module
			topLevel := strings.Split(match[1], ".")[0]
			modules = append(modules, topLevel)
		}
	}
	
	return modules
}

func (rg *RequirementsGenerator) getInstalledPackages() (map[string]string, error) {
	cmd := exec.Command("pip", "freeze")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run 'pip freeze': %v", err)
	}

	packages := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "==") {
			parts := strings.Split(line, "==")
			if len(parts) >= 2 {
				name := strings.ToLower(parts[0])
				packages[name] = line
			}
		}
	}
	
	return packages, scanner.Err()
}

func (rg *RequirementsGenerator) generateRequirements(installedPackages map[string]string) []string {
	var requirements []string
	normalizedFound := make(map[string]bool)
	
	// Normalize found module names
	for module := range rg.foundModules {
		normalized := strings.ToLower(strings.ReplaceAll(module, "-", "_"))
		normalizedFound[normalized] = true
	}
	
	// Match installed packages with found modules
	var packageNames []string
	for pkgName := range installedPackages {
		packageNames = append(packageNames, pkgName)
	}
	sort.Strings(packageNames) // Sort for consistent output
	
	for _, pkgName := range packageNames {
		normalizedPkg := strings.ToLower(strings.ReplaceAll(pkgName, "-", "_"))
		if normalizedFound[normalizedPkg] {
			requirements = append(requirements, installedPackages[pkgName])
		}
	}
	
	return requirements
}

func (rg *RequirementsGenerator) writeRequirements(requirements []string) error {
	file, err := os.Create(rg.outputFile)
	if err != nil {
		return err
	}
	defer file.Close()
	
	writer := bufio.NewWriter(file)
	for _, req := range requirements {
		fmt.Fprintln(writer, req)
	}
	
	return writer.Flush()
}

func (rg *RequirementsGenerator) printResults(requirements []string) {
	if len(requirements) > 0 {
		fmt.Printf("Successfully generated '%s' with detected Python modules and their versions.\n", rg.outputFile)
		fmt.Printf("Contents of '%s':\n", rg.outputFile)
		for _, req := range requirements {
			fmt.Println(req)
		}
	} else {
		fmt.Println("No external Python modules with installed versions were found.")
	}
}