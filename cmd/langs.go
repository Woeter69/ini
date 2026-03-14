package cmd

func init() {
	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "java", Lang: "java", DisplayName: "Java",
		Short: "Initialize a new Java project",
		Long: `Create a new Java project using Gradle. Scaffolds src/, build.gradle,
.gitignore, and README.md.

Examples:
  ini java my-project
  ini java --git my-project`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "csharp", Aliases: []string{"cs", "c#", "dotnet"},
		Lang: "csharp", DisplayName: "C#",
		Short: "Initialize a new C# project",
		Long: `Create a new C# console project using dotnet CLI.

Examples:
  ini csharp my-project
  ini cs my-app
  ini dotnet my-api
  ini csharp --git my-project`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "kotlin", Aliases: []string{"kt"},
		Lang: "kotlin", DisplayName: "Kotlin",
		Short: "Initialize a new Kotlin project",
		Long: `Create a new Kotlin project using Gradle.

Examples:
  ini kotlin my-project
  ini kt my-app
  ini kotlin --git my-project`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "flutter", Aliases: []string{"dart"},
		Lang: "flutter", DisplayName: "Flutter",
		Short: "Initialize a new Flutter project",
		Long: `Create a new Flutter project using flutter create.

Examples:
  ini flutter my-project
  ini dart my-app
  ini flutter --git my-project`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "swift", Lang: "swift", DisplayName: "Swift",
		Short: "Initialize a new Swift project",
		Long: `Create a new Swift executable package using Swift Package Manager.

Examples:
  ini swift my-project
  ini swift --git my-project`,
		Placeholder: "my-tool",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "ruby", Aliases: []string{"rb"},
		Lang: "ruby", DisplayName: "Ruby",
		Short: "Initialize a new Ruby project",
		Long: `Create a new Ruby project with main.rb, Gemfile, .gitignore,
and README.md.

Examples:
  ini ruby my-project
  ini rb my-app
  ini ruby --git my-project`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "r", Aliases: []string{"rlang"},
		Lang: "r", DisplayName: "R",
		Short: "Initialize a new R project",
		Long: `Create a new R project with main.R, R/ directory, .Rprofile,
.gitignore, and README.md.

Examples:
  ini r my-analysis
  ini rlang my-project
  ini r --git my-project`,
		Placeholder: "my-analysis",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "scala", Lang: "scala", DisplayName: "Scala",
		Short: "Initialize a new Scala project",
		Long: `Create a new Scala 3 project using scala-cli.

Examples:
  ini scala my-project
  ini scala --git my-project`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "perl", Aliases: []string{"pl"},
		Lang: "perl", DisplayName: "Perl",
		Short: "Initialize a new Perl project",
		Long: `Create a new Perl project with main.pl, cpanfile, lib/,
.gitignore, and README.md.

Examples:
  ini perl my-project
  ini pl my-script
  ini perl --git my-project`,
		Placeholder: "my-script",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "php", Lang: "php", DisplayName: "PHP",
		Short: "Initialize a new PHP project",
		Long: `Create a new PHP project with main.php, composer.json, src/,
.gitignore, and README.md. Requires PHP 8.2+.

Examples:
  ini php my-project
  ini php --git my-project`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "erlang", Aliases: []string{"erl"},
		Lang: "erlang", DisplayName: "Erlang",
		Short: "Initialize a new Erlang project",
		Long: `Create a new Erlang OTP application using rebar3.

Examples:
  ini erlang my-project
  ini erl my-app
  ini erlang --git my-project`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "elixir", Aliases: []string{"ex", "exs"},
		Lang: "elixir", DisplayName: "Elixir",
		Short: "Initialize a new Elixir project",
		Long: `Create a new Elixir project using mix.

Examples:
  ini elixir my-project
  ini ex my-app
  ini elixir --git my-project`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "haskell", Aliases: []string{"hs"},
		Lang: "haskell", DisplayName: "Haskell",
		Short: "Initialize a new Haskell project",
		Long: `Create a new Haskell project using Cabal.

Examples:
  ini haskell my-project
  ini hs my-app
  ini haskell --git my-project`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "python", Aliases: []string{"py"},
		Lang: "python", DisplayName: "Python",
		Short: "Initialize a new Python project",
		Long: `Create a new Python project using uv. Scaffolds venv, requirements.txt,
main.py, .gitignore, and README.md.

Examples:
  ini python my-project
  ini py my-api`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "go", Aliases: []string{"golang"},
		Lang: "go", DisplayName: "Go",
		Short: "Initialize a new Go project",
		Long: `Create a new Go project using go mod init.

Examples:
  ini go my-project
  ini golang github.com/user/repo`,
		Placeholder: "my-project",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "fortran", Aliases: []string{"f90", "f95", "f03"},
		Lang: "fortran", DisplayName: "Fortran",
		Short: "Initialize a new Fortran project",
		Long: `Create a new Fortran project using fpm.

Examples:
  ini fortran my-project
  ini f90 my-app`,
		Placeholder: "my-simulation",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "asm", Aliases: []string{"assembly", "nasm"},
		Lang: "assembly", DisplayName: "Assembly",
		Short: "Initialize a new Assembly project",
		Long: `Create a new x86-64 Assembly project with NASM and a Makefile.

Examples:
  ini asm my-project
  ini nasm my-kernel`,
		Placeholder: "my-kernel",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "c",
		Lang: "c", DisplayName: "C",
		Short: "Initialize a new C project",
		Long: `Create a new C project with gcc and a Makefile.

Examples:
  ini c my-project`,
		Placeholder: "my-project",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "cpp", Aliases: []string{"c++", "cxx"},
		Lang: "cpp", DisplayName: "C++",
		Short: "Initialize a new C++ project",
		Long: `Create a new C++ project with g++ and a Makefile.

Examples:
  ini cpp my-project
  ini cxx my-engine`,
		Placeholder: "my-engine",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "rust", Aliases: []string{"rs"},
		Lang: "rust", DisplayName: "Rust",
		Short: "Initialize a new Rust project",
		Long: `Create a new Rust project using cargo.

Examples:
  ini rust my-project
  ini rs my-cli`,
		Placeholder: "my-cli",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "zig",
		Lang: "zig", DisplayName: "Zig",
		Short: "Initialize a new Zig project",
		Long: `Create a new Zig project using zig init.

Examples:
  ini zig my-project`,
		Placeholder: "my-project",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "sh", Aliases: []string{"bash", "shell"},
		Lang: "shell", DisplayName: "Shell",
		Short: "Initialize a new shell script project",
		Long: `Create a new shell script project with strict mode.

Examples:
  ini sh my-project
  ini bash my-script`,
		Placeholder: "my-scripts",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "lua",
		Lang: "lua", DisplayName: "Lua",
		Short: "Initialize a new Lua project",
		Long: `Create a new Lua project.

Examples:
  ini lua my-project`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "julia", Aliases: []string{"jl"},
		Lang: "julia", DisplayName: "Julia",
		Short: "Initialize a new Julia project",
		Long: `Create a new Julia project.

Examples:
  ini julia my-project
  ini jl my-app`,
		Placeholder: "my-project",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "nim",
		Lang: "nim", DisplayName: "Nim",
		Short: "Initialize a new Nim project",
		Long: `Create a new Nim project using nimble.

Examples:
  ini nim my-project`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "ocaml", Aliases: []string{"ml"},
		Lang: "ocaml", DisplayName: "OCaml",
		Short: "Initialize a new OCaml project",
		Long: `Create a new OCaml project using dune.

Examples:
  ini ocaml my-project
  ini ml my-app`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "clojure", Aliases: []string{"clj"},
		Lang: "clojure", DisplayName: "Clojure",
		Short: "Initialize a new Clojure project",
		Long: `Create a new Clojure project using tools.deps.

Examples:
  ini clojure my-project
  ini clj my-app`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "d", Aliases: []string{"dlang"},
		Lang: "d", DisplayName: "D",
		Short: "Initialize a new D project",
		Long: `Create a new D project using dub.

Examples:
  ini d my-project
  ini dlang my-app`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "v", Aliases: []string{"vlang"},
		Lang: "v", DisplayName: "V",
		Short: "Initialize a new V project",
		Long: `Create a new V project using v init.

Examples:
  ini v my-project
  ini vlang my-app`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "crystal", Aliases: []string{"cr"},
		Lang: "crystal", DisplayName: "Crystal",
		Short: "Initialize a new Crystal project",
		Long: `Create a new Crystal project using crystal init.

Examples:
  ini crystal my-project
  ini cr my-app`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "fsharp", Aliases: []string{"fs", "f#"},
		Lang: "fsharp", DisplayName: "F#",
		Short: "Initialize a new F# project",
		Long: `Create a new F# project using dotnet new.

Examples:
  ini fsharp my-project
  ini fs my-app`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "objc", Aliases: []string{"objective-c", "objectivec"},
		Lang: "objc", DisplayName: "Objective-C",
		Short: "Initialize a new Objective-C project",
		Long: `Create a new Objective-C project.

Examples:
  ini objc my-project`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "groovy",
		Lang: "groovy", DisplayName: "Groovy",
		Short: "Initialize a new Groovy project",
		Long: `Create a new Groovy project.

Examples:
  ini groovy my-project`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "pascal",
		Lang: "pascal", DisplayName: "Pascal",
		Short: "Initialize a new Pascal project",
		Long: `Create a new Pascal project.

Examples:
  ini pascal my-project`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "ada",
		Lang: "ada", DisplayName: "Ada",
		Short: "Initialize a new Ada project",
		Long: `Create a new Ada project.

Examples:
  ini ada my-project`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "cobol",
		Lang: "cobol", DisplayName: "COBOL",
		Short: "Initialize a new COBOL project",
		Long: `Create a new COBOL project.

Examples:
  ini cobol my-project`,
		Placeholder: "my-app",
	}))

	rootCmd.AddCommand(makeLangCmd(langCmdConfig{
		Use: "bun", Aliases: []string{"js", "ts", "javascript", "typescript", "node"},
		Lang: "bun", DisplayName: "JavaScript/TypeScript",
		Short: "Initialize a new JS/TS project (via Bun)",
		Long: `Create a new JavaScript or TypeScript project using Bun for blazing fast execution.

Examples:
  ini js my-project
  ini ts my-api
  ini bun my-app`,
		Placeholder: "my-app",
	}))
}
