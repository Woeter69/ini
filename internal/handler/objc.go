package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("objc", &ObjCHandler{}) }

type ObjCHandler struct{}

func (o *ObjCHandler) Name() string { return "Objective-C" }
func (o *ObjCHandler) Validate() error { return nil }

func (o *ObjCHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	os.MkdirAll(filepath.Join(config.Path, "src"), 0o755)
	os.MkdirAll(filepath.Join(config.Path, "build"), 0o755)

	fmt.Printf("  %s Creating main.m and Makefile...\n", ui.Arrow)

	mainM := fmt.Sprintf(`#import <Foundation/Foundation.h>

int main(int argc, const char * argv[]) {
    @autoreleasepool {
        NSLog(@"Hello from %s!");
    }
    return 0;
}
`, config.Name)
	if err := os.WriteFile(filepath.Join(config.Path, "src", "main.m"), []byte(mainM), 0o644); err != nil {
		return err
	}

	makefile := fmt.Sprintf(`CC = clang
CFLAGS = -Wall -framework Foundation

SRC_DIR = src
BUILD_DIR = build
TARGET = $(BUILD_DIR)/%s

SRCS = $(wildcard $(SRC_DIR)/*.m)
OBJS = $(patsubst $(SRC_DIR)/%%.m,$(BUILD_DIR)/%%.o,$(SRCS))

all: $(TARGET)

$(TARGET): $(OBJS)
	$(CC) $(CFLAGS) -o $@ $^

$(BUILD_DIR)/%%.o: $(SRC_DIR)/%%.m
	@mkdir -p $(BUILD_DIR)
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -rf $(BUILD_DIR)/*.o $(TARGET)
`, config.Name)
	if err := os.WriteFile(filepath.Join(config.Path, "Makefile"), []byte(makefile), 0o644); err != nil {
		return err
	}
	fmt.Printf("  %s Scaffolding complete\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git { if err := scaffold.InitGit(config.Path); err != nil { return err } }

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your Objective-C project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  make\n")
	s.WriteString(fmt.Sprintf("  ./build/%s\n", config.Name))
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
