package templates

import "fmt"

// Readme generates a README.md for the given project name and language.
func Readme(projectName, language string) string {
	switch language {
	case "python":
		return fmt.Sprintf(`# %s

A Python project.

## Getting Started

`+"```bash"+`
# Activate virtual environment
source .venv/bin/activate  # Linux/macOS
# .venv\Scripts\activate   # Windows

# Install dependencies
uv pip install -r requirements.txt

# Run
python main.py
`+"```"+`
`, projectName)
	case "go":
		return fmt.Sprintf(`# %s

A Go project.

## Getting Started

`+"```bash"+`
# Run
go run .

# Build
go build -o %s .

# Test
go test ./...
`+"```"+`
`, projectName, projectName)
	case "fortran":
		return fmt.Sprintf(`# %s

A Fortran project built with fpm.

## Getting Started

`+"```bash"+`
# Run
fpm run

# Build
fpm build

# Test
fpm test
`+"```"+`
`, projectName)
	case "assembly":
		return fmt.Sprintf(`# %s

An x86-64 Assembly project (NASM + ld).

## Getting Started

`+"```bash"+`
# Build
make

# Run
./build/%s

# Clean
make clean
`+"```"+`
`, projectName, projectName)
	case "c":
		return fmt.Sprintf(`# %s

A C project.

## Getting Started

`+"```bash"+`
# Build
make

# Run
./build/%s

# Clean
make clean
`+"```"+`
`, projectName, projectName)
	case "cpp":
		return fmt.Sprintf(`# %s

A C++ project.

## Getting Started

`+"```bash"+`
# Build
make

# Run
./build/%s

# Clean
make clean
`+"```"+`
`, projectName, projectName)
	case "rust":
		return fmt.Sprintf(`# %s

A Rust project.

## Getting Started

`+"```bash"+`
# Run
cargo run

# Build (release)
cargo build --release

# Test
cargo test
`+"```"+`
`, projectName)
	case "zig":
		return fmt.Sprintf(`# %s

A Zig project.

## Getting Started

`+"```bash"+`
# Build & Run
zig build run

# Build only
zig build

# Test
zig build test
`+"```"+`
`, projectName)
	case "shell":
		return fmt.Sprintf(`# %s

A shell script project.

## Getting Started

`+"```bash"+`
# Run
./%s.sh

# Or with bash explicitly
bash %s.sh
`+"```"+`
`, projectName, projectName, projectName)
	case "java":
		return fmt.Sprintf(`# %s

A Java project built with Gradle.

## Getting Started

`+"```bash"+`
# Run
gradle run

# Build
gradle build

# Test
gradle test
`+"```"+`
`, projectName)
	case "csharp":
		return fmt.Sprintf(`# %s

A C# project.

## Getting Started

`+"```bash"+`
# Run
dotnet run

# Build
dotnet build

# Test
dotnet test
`+"```"+`
`, projectName)
	case "kotlin":
		return fmt.Sprintf(`# %s

A Kotlin project built with Gradle.

## Getting Started

`+"```bash"+`
# Run
gradle run

# Build
gradle build

# Test
gradle test
`+"```"+`
`, projectName)
	case "flutter":
		return fmt.Sprintf(`# %s

A Flutter project.

## Getting Started

`+"```bash"+`
# Run (debug)
flutter run

# Build
flutter build apk

# Test
flutter test
`+"```"+`
`, projectName)
	case "swift":
		return fmt.Sprintf(`# %s

A Swift project.

## Getting Started

`+"```bash"+`
# Run
swift run

# Build
swift build

# Test
swift test
`+"```"+`
`, projectName)
	case "ruby":
		return fmt.Sprintf(`# %s

A Ruby project.

## Getting Started

`+"```bash"+`
# Install dependencies
bundle install

# Run
ruby main.rb
`+"```"+`
`, projectName)
	case "r":
		return fmt.Sprintf(`# %s

An R project.

## Getting Started

`+"```bash"+`
# Run
Rscript main.R

# Interactive
R
`+"```"+`
`, projectName)
	case "scala":
		return fmt.Sprintf(`# %s

A Scala project using scala-cli.

## Getting Started

`+"```bash"+`
# Run
scala-cli run .

# Compile
scala-cli compile .

# Test
scala-cli test .
`+"```"+`
`, projectName)
	case "perl":
		return fmt.Sprintf(`# %s

A Perl project.

## Getting Started

`+"```bash"+`
# Run
perl main.pl

# Install dependencies
cpanm --installdeps .
`+"```"+`
`, projectName)
	case "php":
		return fmt.Sprintf(`# %s

A PHP project.

## Getting Started

`+"```bash"+`
# Run
php main.php

# Install dependencies
composer install

# Dev server
php -S localhost:8000
`+"```"+`
`, projectName)
	case "erlang":
		return fmt.Sprintf(`# %s

An Erlang project built with rebar3.

## Getting Started

`+"```bash"+`
# Compile
rebar3 compile

# Run shell
rebar3 shell

# Test
rebar3 eunit
`+"```"+`
`, projectName)
	case "elixir":
		return fmt.Sprintf(`# %s

An Elixir project.

## Getting Started

`+"```bash"+`
# Run
mix run

# Interactive
iex -S mix

# Test
mix test
`+"```"+`
`, projectName)
	case "haskell":
		return fmt.Sprintf(`# %s

A Haskell project built with Cabal.

## Getting Started

`+"```bash"+`
# Run
cabal run

# Build
cabal build

# Test
cabal test
`+"```"+`
`, projectName)
	case "lua":
		return fmt.Sprintf(`# %s

A Lua project.

## Getting Started

`+"```bash"+`
# Run
lua main.lua
`+"```"+`
`, projectName)
	case "julia":
		return fmt.Sprintf(`# %s

A Julia project.

## Getting Started

`+"```bash"+`
# Run
julia main.jl
`+"```"+`
`, projectName)
	case "nim":
		return fmt.Sprintf(`# %s

A Nim project.

## Getting Started

`+"```bash"+`
# Run
nim c -r main.nim

# Build (release)
nim c -d:release main.nim
`+"```"+`
`, projectName)
	case "ocaml":
		return fmt.Sprintf(`# %s

An OCaml project built with Dune.

## Getting Started

`+"```bash"+`
# Run
dune exec ./bin/main.exe

# Build
dune build

# Test
dune runtest
`+"```"+`
`, projectName)
	case "clojure":
		return fmt.Sprintf(`# %s

A Clojure project.

## Getting Started

`+"```bash"+`
# Run
clj -M -m main
`+"```"+`
`, projectName)
	case "d":
		return fmt.Sprintf(`# %s

A D language project built with DUB.

## Getting Started

`+"```bash"+`
# Run
dub run

# Build
dub build

# Test
dub test
`+"```"+`
`, projectName)
	case "v":
		return fmt.Sprintf(`# %s

A V language project.

## Getting Started

`+"```bash"+`
# Run
v run .

# Build
v .

# Test
v test .
`+"```"+`
`, projectName)
	case "crystal":
		return fmt.Sprintf(`# %s

A Crystal project.

## Getting Started

`+"```bash"+`
# Run
crystal run src/%s.cr

# Build
crystal build src/%s.cr

# Test
crystal spec
`+"```"+`
`, projectName, projectName, projectName)
	case "fsharp":
		return fmt.Sprintf(`# %s

An F# project built with dotnet.

## Getting Started

`+"```bash"+`
# Run
dotnet run

# Build
dotnet build

# Test
dotnet test
`+"```"+`
`, projectName)
	case "objc":
		return fmt.Sprintf(`# %s

An Objective-C project.

## Getting Started

`+"```bash"+`
# Build
make

# Run
./build/%s

# Clean
make clean
`+"```"+`
`, projectName, projectName)
	case "groovy":
		return fmt.Sprintf(`# %s

A Groovy project.

## Getting Started

`+"```bash"+`
# Run
groovy main.groovy
`+"```"+`
`, projectName)
	case "pascal":
		return fmt.Sprintf(`# %s

A Pascal project (Free Pascal).

## Getting Started

`+"```bash"+`
# Compile
fpc main.pas

# Run
./main
`+"```"+`
`, projectName)
	case "ada":
		return fmt.Sprintf(`# %s

An Ada project.

## Getting Started

`+"```bash"+`
# Compile and Build
gnatmake main.adb

# Run
./main
`+"```"+`
`, projectName)
	case "cobol":
		return fmt.Sprintf(`# %s

A COBOL project (GnuCOBOL).

## Getting Started

`+"```bash"+`
# Compile
cobc -x -free -o main main.cbl

# Run
./main
`+"```"+`
`, projectName)
	case "bun":
		return fmt.Sprintf(`# %s

A modern JavaScript/TypeScript project built with [Bun](https://bun.sh/).

## Getting Started

To install dependencies:
`+"```bash"+`
bun install
`+"```"+`

To run the project:
`+"```bash"+`
bun run index.ts
`+"```"+`
`, projectName)
	default:
		return fmt.Sprintf("# %s\n\nA new project.\n", projectName)
	}
}
