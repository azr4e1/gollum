package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/azr4e1/gollum"
	"github.com/azr4e1/gollum/message"
	// "github.com/joho/godotenv"
)

// var resMess string

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	key := os.Getenv("OPENAI_API_KEY")
	client, _ := gollum.NewClient(gollum.WithAPIKey(key), gollum.WithProvider(gollum.OPENAI))
	// client.EnableStream(streamingFunc)
	// client.Timeout = time.Second * 40

	chat := message.NewChat()
	chat.SetSystemMessage("You are a Yakuza member. Act like it! Feel free to take the lead and ask the user what they want only when necessary. When you are finished reasoning or talking and you want the user to prompt you with something, terminate your sentence with this word: DONE!!!")
	tools := map[string]func(json.RawMessage) (string, error){
		ReadFileSchema.Name():  ReadFile,
		ListFilesSchema.Name(): ListFiles,
		EditFileSchema.Name():  EditFile,
	}
	// chat.Add(message.UserMessage("read the content of the file `secret-file.txt`"))
	// req, res, _ := client.Complete(gollum.WithChat(chat), gollum.WithModel("gpt-4o"), gollum.WithTool(ReadFileSchema))
	// v, _ := json.MarshalIndent(req, "", "  ")
	// fmt.Println(string(v))
	// v, _ = json.MarshalIndent(res, "", "  ")
	// fmt.Println(string(v))
	// t, _ := res.Tool()
	// fmt.Println(t)

	// args := t.Arguments
	// // print(string(args))
	// readFileInput := ReadFileInput{}
	// err = json.Unmarshal(args, &readFileInput)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(readFileInput)

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
		_, res, err := client.Complete(gollum.WithChat(chat), gollum.WithModel("gpt-4o"), gollum.WithTool(ReadFileSchema, ListFilesSchema, EditFileSchema))
		if err != nil {
			log.Fatal(err)
		}
		switch res.Type {
		case gollum.Text:
			resMess := strings.TrimSpace(res.Content())
			if strings.ToUpper(resMess[len(resMess)-7:]) == "DONE!!!" {
				userInput = true
				fmt.Printf("\u001b[93mBot\u001b[0m: ")
				fmt.Println(resMess[:len(resMess)-7])
				chat.Add(message.AssistantMessage(resMess[:len(resMess)-7]))
				continue
			}
			fmt.Printf("\u001b[93mBot\u001b[0m: ")
			fmt.Println(resMess)
			chat.Add(message.AssistantMessage(resMess))
			userInput = false
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

type ListFilesInput struct {
	Path string `json:"path,omitempty" jsonschema_description:"Optional relative path to list files from. Defaults to current directory if not provided."`
}

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
