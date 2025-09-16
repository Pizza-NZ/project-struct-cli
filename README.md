# Project Structure CLI
A simple command-line tool to scan a source code directory and compile all file contents into a single document, formatted for different use cases like code review or AI model ingestion.
### What it Does
This tool walks through a specified directory, reads the files, respects .gitignore rules, and then uses templates to generate a single output file. This is useful for:
- Creating a single-file context of a project for Large Language Models (LLMs).
- Packaging up all changes for a manual code review.
- Generating simple project documentation.
### Installation
To install the command-line tool, you can use `go install`:
```bash 
go install pizza-nz/project-struct-cli@latest
```
Alternatively, you can clone the repository and build it manually.
```bash
git clone https://github.com/pizza-nz/project-struct-cli.git
cd project-struct-cli
make build
```
This will create the binary at `./out/project-struct-cli`.
### Usage
Run the tool against a source directory. The output file will be created in your current directory.
```bash
# Scan the current directory and create 'project_structure.md'
./out/project-struct-cli -src .

# Specify a different source directory and output file
./out/project-struct-cli -src ../my-other-project -out my-project.md

# Generate output specifically for an LLM
./out/project-struct-cli -src . -format llm -out for-llm.txt
```
### Command-Line Flags
- `-src`:The source directory to scan (default: `.` a.k.a. the current directory).
- `-out`: The name of the output document (default: `project_structure.md`).
- `-format`: The output format template to use. Options are `default`, `review`, `llm` (default: `default`).
##### Contributing
Contributions are welcome! Please feel free to open an issue or submit a pull request.

##### License
This project is licensed under the MIT License.