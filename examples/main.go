package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/azr4e1/gollum"
	"github.com/azr4e1/gollum/message"
	// "github.com/joho/godotenv"
)

const Keyword = "***CONTINUE***"

// var resMess string

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	key := os.Getenv("OPENAI_API_KEY")
	client, err := gollum.NewClient(gollum.WithAPIKey(key), gollum.WithProvider(gollum.OPENAI))
	if err != nil {
		panic(err)
	}
	// client.EnableStream(streamingFunc)
	// client.Timeout = time.Second * 40

	chat := message.NewChat()
	chat.SetSystemMessage("You are a helpful assistant that always thinks step by step. Think step by step and don't ask for user input or confirmation until you have completely satisfied the user's request. Feel free to take the lead and ask the user what they want only when necessary. If you need more turns to completely answer the user question, end your turn with the keyword ***CONTINUE***")
	tools := map[string]func(json.RawMessage) (string, error){
		ReadFileSchema.Name():      ReadFile,
		ListFilesSchema.Name():     ListFiles,
		EditFileSchema.Name():      EditFile,
		ExecuteScriptSchema.Name(): ExecuteScript,
	}

	userInput := true
	for {
		if userInput {
			fmt.Print("\u001b[94mYou\u001b[0m: ")
			input, ok := getInput()
			if !ok {
				break
			}
			chat.Add(message.UserMessage(input))
		}
		_, res, err := client.Complete(gollum.WithChat(chat), gollum.WithModel("gpt-4o"), gollum.WithTool(ReadFileSchema, ListFilesSchema, EditFileSchema, ExecuteScriptSchema))
		if err != nil {
			log.Fatal(err)
		}
		switch res.Type {
		case gollum.Text:
			resMess := strings.TrimSpace(res.Content())
			if strings.ToUpper(resMess[len(resMess)-len(Keyword):]) == Keyword {
				userInput = false
				fmt.Printf("\u001b[93mBot\u001b[0m: ")
				fmt.Println(resMess)
				chat.Add(message.AssistantMessage(resMess[:len(resMess)-len(Keyword)]))
				continue
			}
			fmt.Printf("\u001b[93mBot\u001b[0m: ")
			fmt.Println(resMess)
			chat.Add(message.AssistantMessage(resMess))
			userInput = true
		case gollum.ToolCall:
			t, err := res.Tool()
			fmt.Printf("System: using tool: '%s'\n", t.Name)
			if err != nil {
				panic(err)
			}
			// fmt.Println(t.Name)
			useTool, ok := tools[t.Name]
			if !ok {
				resMess := fmt.Sprintf("You used tool: %s; Tool %s is not available", t.Name, t.Name)
				fmt.Printf("System: %s", resMess)
				chat.Add(message.UserMessage(resMess))
				userInput = false
				continue
			}
			res, err := useTool(t.Arguments)
			if err != nil {
				resMess := fmt.Sprintf("You used tool: %s; Error %s", t.Name, err)
				fmt.Printf("System: %s", resMess)
				chat.Add(message.UserMessage(resMess))
				userInput = true
				continue
			}
			chat.Add(message.UserMessage(fmt.Sprintf("You used tool: %s; This is the output: %s", t.Name, res)))
			userInput = false
		}
	}
}

func getInput() (string, bool) {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return "", false
	}
	return scanner.Text(), true
}

type ReadFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory"`
}

type EditFileInput struct {
	Path   string `json:"path" jsonschema_description:"The path to the file"`
	OldStr string `json:"old_str" jsonschema_description:"Text to search for - must match exactly and must only have one match exactly"`
	NewStr string `json:"new_str" jsonschema_description:"Text to replace old_str with"`
}

type ExecuteScriptInput struct {
	Executable string `json:"executable" jsonschema_description:"The name of the executable that runs the script, e.g. 'python'"`
	Path       string `json:"path" jsonschema_description:"The path to the script"`
}

type ListFilesInput struct {
	Path string `json:"path,omitempty" jsonschema_description:"Optional relative path to list files from. Defaults to current directory if not provided."`
}

var ExecuteScriptSchema = gollum.NewTool("execute_script",
	"Execute a script with the given executable and script path",
	gollum.GenerateArguments[ExecuteScriptInput](),
	nil,
)

var ListFilesSchema = gollum.NewTool("list_files",
	"List files and directories at a given path. If no path is provided, lists files in the current directory.",
	gollum.GenerateArguments[ListFilesInput](),
	nil,
)

var ReadFileSchema = gollum.NewTool("read_file",
	"Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names",
	gollum.GenerateArguments[ReadFileInput](),
	nil)

var EditFileSchema = gollum.NewTool("edit_file",
	`Make edits to a text file.

Replaces 'old_str' with 'new_str' in the given file. 'old_str' and 'new_str' MUST be different from each other.

If the file specified with path doesn't exist, it will be created.`,
	gollum.GenerateArguments[EditFileInput](),
	nil,
)

func ReadFile(input json.RawMessage) (string, error) {
	readFileInput := ReadFileInput{}
	err := json.Unmarshal([]byte(input), &readFileInput)
	if err != nil {
		panic(err)
	}

	content, err := os.ReadFile(readFileInput.Path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func ListFiles(input json.RawMessage) (string, error) {
	listFilesInput := ListFilesInput{}
	err := json.Unmarshal(input, &listFilesInput)
	if err != nil {
		panic(err)
	}

	dir := "."
	if listFilesInput.Path != "" {
		dir = listFilesInput.Path
	}

	var files []string
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		if relPath != "." {
			if info.IsDir() {
				files = append(files, relPath+"/")
			} else {
				files = append(files, relPath)
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	result, err := json.Marshal(files)
	if err != nil {
		return "", err
	}

	fmt.Println(string(result))
	return string(result), nil
}

func EditFile(input json.RawMessage) (string, error) {
	editFileInput := EditFileInput{}
	err := json.Unmarshal(input, &editFileInput)
	if err != nil {
		return "", err
	}

	if editFileInput.Path == "" || editFileInput.OldStr == editFileInput.NewStr {
		return "", fmt.Errorf("invalid input parameters")
	}

	content, err := os.ReadFile(editFileInput.Path)
	if err != nil {
		if os.IsNotExist(err) && editFileInput.OldStr == "" {
			return createNewFile(editFileInput.Path, editFileInput.NewStr)
		}
		return "", err
	}

	oldContent := string(content)
	newContent := strings.Replace(oldContent, editFileInput.OldStr, editFileInput.NewStr, -1)

	if oldContent == newContent && editFileInput.OldStr != "" {
		return "", fmt.Errorf("old_str not found in file")
	}

	err = os.WriteFile(editFileInput.Path, []byte(newContent), 0644)
	if err != nil {
		return "", err
	}

	return "OK", nil
}

func createNewFile(filePath, content string) (string, error) {
	dir := path.Dir(filePath)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
	}

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	return fmt.Sprintf("Successfully created file %s", filePath), nil
}

func ExecuteScript(input json.RawMessage) (string, error) {
	executeScriptInput := ExecuteScriptInput{}
	err := json.Unmarshal(input, &executeScriptInput)
	if err != nil {
		return "", err
	}

	if executeScriptInput.Executable == "" || executeScriptInput.Path == "" {
		return "", errors.New("invalid input parameters")
	}
	executable := executeScriptInput.Executable
	filePath := executeScriptInput.Path

	cmd := exec.Command(executable, filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
