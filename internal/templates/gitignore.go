package templates

// Gitignore patterns per language. Each handler just picks its key.
var Gitignore = map[string]string{
	"python": `# Byte-compiled / optimized / DLL files
__pycache__/
*.py[cod]
*$py.class

# Virtual environments
.venv/
venv/
ENV/
env/

# Distribution / packaging
build/
dist/
*.egg-info/
*.egg
.eggs/

# Installer logs
pip-log.txt
pip-delete-this-directory.txt

# Unit test / coverage
htmlcov/
.tox/
.nox/
.coverage
.coverage.*
.cache
nosetests.xml
coverage.xml
*.cover
*.py,cover
.hypothesis/
.pytest_cache/
pytestdebug.log

# mypy
.mypy_cache/

# ruff
.ruff_cache/

# Environments
.env
.env.local
*.env

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# uv
uv.lock
`,

	"go": `# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output of go coverage
*.out

# Go workspace
go.work
go.work.sum

# Vendor directory
vendor/

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Build output
bin/
dist/

# Environment
.env
.env.local
`,

	"fortran": `# Build output
build/
*.o
*.mod
*.smod
*.a
*.so
*.dylib
*.exe

# fpm
build/

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db
`,

	"assembly": `# Object files
*.o
*.obj

# Binaries
*.bin
*.exe
*.out
*.elf

# Listings
*.lst

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Build output
build/
`,

	"c": `# Object files
*.o
*.obj

# Binaries
*.out
*.exe
*.so
*.dylib
*.a

# Debug
*.dSYM/

# Build output
build/

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db
`,

	"cpp": `# Object files
*.o
*.obj

# Binaries
*.out
*.exe
*.so
*.dylib
*.a

# Debug
*.dSYM/

# Build output
build/

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db
`,

	"rust": `# Build output
target/

# Cargo.lock for libraries (keep for binaries)
# Cargo.lock

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Environment
.env
.env.local
`,

	"zig": `# Build output
zig-out/
zig-cache/
.zig-cache/

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db
`,

	"shell": `# Logs
*.log

# Environment
.env
.env.local

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db
`,

	"java": `# Compiled
*.class
*.jar
*.war
*.ear

# Build
build/
.gradle/
out/

# IDE
.idea/
*.iml
.vscode/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"csharp": `# Build
bin/
obj/
out/

# User-specific
*.user
*.suo
*.userprefs

# IDE
.vs/
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"kotlin": `# Compiled
*.class
*.jar

# Build
build/
.gradle/
out/

# IDE
.idea/
*.iml
.vscode/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"flutter": `# Flutter/Dart
.dart_tool/
.packages
build/
.flutter-plugins
.flutter-plugins-dependencies

# IDE
.idea/
*.iml
.vscode/
*.swp
*~

# OS
.DS_Store
Thumbs.db

# Pub
.pub-cache/
.pub/
pubspec.lock
`,

	"dart": `# Dart
.dart_tool/
.packages
build/

# IDE
.idea/
*.iml
.vscode/
*.swp
*~

# OS
.DS_Store
Thumbs.db

# Pub
.pub-cache/
.pub/
pubspec.lock
`,

	"swift": `# Build
.build/
.swiftpm/
Packages/
xcuserdata/
*.xcodeproj/

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"ruby": `# Gems
vendor/bundle/
.bundle/
*.gem

# Environment
.env
.env.local

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"r": `# R history and workspace
.Rhistory
.Rdata
.RData
.Ruserdata

# R environment
.Renviron
renv/library/
renv/staging/

# IDE
.Rproj.user/
.vscode/
*.swp
*~

# OS
.DS_Store
Thumbs.db

# Output
*.pdf
plots/
`,

	"scala": `# Build
target/
.bsp/
.metals/
.bloop/
project/target/
project/project/

# IDE
.idea/
*.iml
.vscode/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"perl": `# Build
blib/
Makefile
Makefile.old
pm_to_blib
MYMETA.*

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"php": `# Dependencies
vendor/

# Environment
.env
.env.local
.env.*.local

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db

# Cache
*.cache
`,

	"erlang": `# Build
_build/
ebin/
*.beam
*.plt
*.d

# Deps
deps/

# Release
rel/

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"elixir": `# Build
_build/
deps/

# Coverage
cover/

# Dialyzer
priv/plts/

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db

# Elixir
*.ez
erl_crash.dump
`,

	"haskell": `# Build
dist/
dist-newstyle/
.cabal-sandbox/
cabal.sandbox.config
*.o
*.hi
*.dyn_o
*.dyn_hi
*.prof
*.hp
*.tix

# Stack
.stack-work/

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"lua": `# Compiled/Luac
*.luac

# Modules/Rocks
luarocks_modules/
.luarocks/

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"julia": `# Build
Project.toml.backup
Manifest.toml

# Coverage
*.jl.cov
*.jl.*.cov

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"nim": `# Nimcache
nimcache/

# Binaries
*.exe
*.out

# Dependencies
nimble/

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"ocaml": `# Compiled and Build
_build/
*.class
*.cmi
*.cmo
*.cmt
*.cmti
*.cmx
*.cma
*.cmxa
*.a
*.o
*.so

# Opam
_opam/

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"clojure": `# Leiningen
/target/
/pom.xml
/.lein-env
/.lein-failures
/.lein-plugins
/.lein-repl-history

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"d": `# Build
*.o
*.obj
*.exe
*.out
.dub/
docs/

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"v": `# Compiled
*.exe
*.out

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"crystal": `# Build
lib/
bin/
.shards/
*.dwarf

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"fsharp": `# Build
bin/
obj/
out/

# IDE
.ionide/
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"objc": `# Build
build/
*.o
*.obj
*.out
*.exe
*.dSYM/
*.a
*.dylib

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"groovy": `# Build
build/
out/
.gradle/

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"pascal": `# Compiled
*.o
*.ppu
*.a
*.exe
*.out
backup/
*.bak

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"ada": `# Build
obj/
*.o
*.ali
*.exe
*.out

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"cobol": `# Compiled
*.o
*.so
*.exe
*.out
*.cbl.c
*.cbl.h
*.cbl.i

# IDE
.vscode/
.idea/
*.swp
*~

# OS
.DS_Store
Thumbs.db
`,

	"bun": `# Logs
logs
*.log
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# Runtime data
pids
*.pid
*.seed
*.pid.lock

# Directory for instrumented libs generated by jscoverage/JSCover
lib-cov

# Coverage directory used by tools like istanbul
coverage

# nyc test coverage
.nyc_output

# Grunt intermediate storage (http://gruntjs.com/creating-plugins#storing-task-files)
.grunt

# Bower dependency directory (https://bower.io/)
bower_components

# node-waf configuration
.lock-wscript

# Compiled binary addons (https://nodejs.org/api/addons.html)
build/Release

# Dependency directories
node_modules/
jspm_packages/

# TypeScript v1 declaration files
typings/

# Optional npm cache directory
.npm

# Optional eslint cache
.eslintcache

# Optional REPL history
.node_repl_history

# Output of 'npm pack'
*.tgz

# Yarn Integrity file
.yarn-integrity

# dotenv environment variables file
.env
.env.test

# parcel-bundler cache (https://parceljs.org/)
.cache

# next.js build output
.next

# nuxt.js build output
.nuxt

# vuepress build output
.vuepress/dist

# serverless directories
.serverless/

# FuseBox cache
.fusebox/

# DynamoDB Local files
.dynamodb/

# OS
.DS_Store
Thumbs.db
`,
}
