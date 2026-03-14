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
	Register("c", &CHandler{})
}

// CHandler scaffolds C projects with gcc and a Makefile.
type CHandler struct{}

func (c *CHandler) Name() string {
	return "C"
}

func (c *CHandler) Validate() error {
	return nil
}

func (c *CHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	// 1. Create project directory + src/ and include/
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	srcDir := filepath.Join(projectDir, "src")
	includeDir := filepath.Join(projectDir, "include")
	os.MkdirAll(srcDir, 0o755)
	os.MkdirAll(includeDir, 0o755)

	// 2. Create src/main.c
	fmt.Printf("  %s Creating src/main.c...\n", ui.Arrow)
	mainC := fmt.Sprintf(`#include <stdio.h>

int main(void) {
    printf("Hello from %s!\n");
    return 0;
}
`, config.Name)
	if err := os.WriteFile(filepath.Join(srcDir, "main.c"), []byte(mainC), 0o644); err != nil {
		return fmt.Errorf("failed to create main.c: %w", err)
	}
	fmt.Printf("  %s src/main.c created\n", ui.CheckMark)

	// 3. Create Makefile
	fmt.Printf("  %s Creating Makefile...\n", ui.Arrow)
	makefile := fmt.Sprintf(`CC        = gcc
CFLAGS    = -Wall -Wextra -std=c17 -O2
INCLUDES  = -Iinclude
SRC_DIR   = src
BUILD_DIR = build
NAME      = %s

SRCS = $(wildcard $(SRC_DIR)/*.c)
OBJS = $(SRCS:$(SRC_DIR)/%%.c=$(BUILD_DIR)/%%.o)

all: $(BUILD_DIR)/$(NAME)

$(BUILD_DIR)/$(NAME): $(OBJS)
	$(CC) $(CFLAGS) -o $@ $^

$(BUILD_DIR)/%%.o: $(SRC_DIR)/%%.c | $(BUILD_DIR)
	$(CC) $(CFLAGS) $(INCLUDES) -c -o $@ $<

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

	// 4. Create .gitignore (shared)
	if err := scaffold.WriteGitignore(projectDir, config.Language); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	// 5. Create README.md (shared)
	if err := scaffold.WriteReadme(projectDir, config.Name, config.Language); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	// 6. Initialize git repo (if --git flag is set)
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
	summary.WriteString(ui.SuccessStyle.Render("🚀 Your C project is ready!"))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  make && make run\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()

	return nil
}
