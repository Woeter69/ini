package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() {
	Register("assembly", &AssemblyHandler{})
}

// AssemblyHandler scaffolds x86-64 Assembly projects with NASM and a Makefile.
type AssemblyHandler struct{}

func (a *AssemblyHandler) Name() string {
	return "Assembly"
}

func (a *AssemblyHandler) Validate() error {
	// No strict toolchain requirement — just create the project.
	// nasm/ld will be needed at build time, not at scaffold time.
	return nil
}

func (a *AssemblyHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	// 1. Create project directory
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 2. Create src/ directory
	srcDir := filepath.Join(projectDir, "src")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		return fmt.Errorf("failed to create src/: %w", err)
	}

	// 3. Create main.asm (x86-64 Linux hello world)
	fmt.Printf("  %s Creating src/main.asm...\n", ui.Arrow)
	mainAsm := fmt.Sprintf(`; %s — x86-64 Assembly (NASM, Linux)

section .data
    msg db "Hello from %s!", 10
    msg_len equ $ - msg

section .text
    global _start

_start:
    ; write(stdout, msg, msg_len)
    mov rax, 1          ; syscall: write
    mov rdi, 1          ; fd: stdout
    lea rsi, [rel msg]  ; buffer
    mov rdx, msg_len    ; count
    syscall

    ; exit(0)
    mov rax, 60         ; syscall: exit
    xor rdi, rdi        ; status: 0
    syscall
`, config.Name, config.Name)
	if err := os.WriteFile(filepath.Join(srcDir, "main.asm"), []byte(mainAsm), 0o644); err != nil {
		return fmt.Errorf("failed to create main.asm: %w", err)
	}
	fmt.Printf("  %s src/main.asm created\n", ui.CheckMark)

	// 4. Create Makefile
	fmt.Printf("  %s Creating Makefile...\n", ui.Arrow)
	makefile := fmt.Sprintf(`NAME = %s
SRC  = src/main.asm
BUILD_DIR = build

all: $(BUILD_DIR)/$(NAME)

$(BUILD_DIR)/$(NAME): $(SRC) | $(BUILD_DIR)
	nasm -f elf64 -o $(BUILD_DIR)/main.o $(SRC)
	ld -o $(BUILD_DIR)/$(NAME) $(BUILD_DIR)/main.o

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

run: $(BUILD_DIR)/$(NAME)
	./$(BUILD_DIR)/$(NAME)

clean:
	rm -rf $(BUILD_DIR)

.PHONY: all run clean
`, config.Name)
	if err := os.WriteFile(filepath.Join(projectDir, "Makefile"), []byte(makefile), 0o644); err != nil {
		return fmt.Errorf("failed to create Makefile: %w", err)
	}
	fmt.Printf("  %s Makefile created\n", ui.CheckMark)

	// 5. Create .gitignore (shared)
	if err := scaffold.WriteGitignore(projectDir, config.Language); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	// 6. Create README.md (shared)
	if err := scaffold.WriteReadme(projectDir, config.Name, config.Language); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	// 7. Initialize git repo (if --git flag is set)
	if config.Git {
		if err := scaffold.InitGit(projectDir); err != nil {
			return err
		}
	}

	// Print summary
	fmt.Println()
	relPath, _ := filepath.Rel(".", projectDir)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render("🚀 Your Assembly project is ready!"))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  make && make run\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()

	return nil
}
