package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type task struct {
	Title       string
	Description string
}

func main() {
	command := getCommand(os.Args)

	fmt.Println("Hello!")
	switch command {
	case "list":
		commandList()
	case "create":
		commandCreate(os.Stdin)
	case "done":
		commandDone(os.Args)
	default:
		fmt.Println("Sorry, command not recognized.")
	}
}

func getCommand(args []string) string {
	if len(args) == 1 {
		return "list"
	} else {
		return args[1]
	}
}

func commandCreate(reader io.Reader) {
	title, description := getNewTaskInfo(reader)
	newTask := createNewTask(title, description)
	saveTask(newTask)
	fmt.Printf("\nCreated %s", title)
}

func createNewTask(title string, description string) task {
	return task{Title: title, Description: description}
}

func getNewTaskInfo(rd io.Reader) (string, string) {
	reader := bufio.NewReader(rd)
	fmt.Print("Please enter the title: ")
	title, _ := reader.ReadString('\n')

	fmt.Print("Please enter the description: ")
	description, _ := reader.ReadString('\n')

	return strings.Replace(title, "\n", "", 1), strings.Replace(description, "\n", "", 1)
}

func saveTask(t task) {
	content, _ := json.Marshal(t)
	ioutil.WriteFile("tasks/"+getTaskSlug(t.Title)+".json", content, 0644)
}

func loadTask(slug string) task {
	jsonFile, err := os.Open("tasks/" + slug + ".json")
	if err != nil {
		fmt.Println(err)
		panic("Bad")
	}
	defer jsonFile.Close()

	jsonContent, _ := ioutil.ReadAll(jsonFile)

	t := task{}
	json.Unmarshal(jsonContent, &t)
	return t
}

func getTaskSlug(title string) string {
	return strings.ToLower(strings.Replace(title, " ", "-", -1))
}

func commandDone(args []string) {
	if len(args) < 3 {
		fmt.Println("You must name a task")
		return
	}

	tasks := findTasks(args[2])
	switch len(tasks) {
	case 0:
		fmt.Println("No tasks matched your search")
	case 1:
		removeTask(tasks[0])
	default:
		taskSlug, err := getTaskToRemove(tasks)
		if err != nil {
			fmt.Print(err)
			fmt.Println()
		} else {
			removeTask(taskSlug)
		}
	}
}

func getTaskToRemove(taskSlugs []string) (string, error) {
	for index, taskSlug := range taskSlugs {
		task := loadTask(taskSlug)
		fmt.Printf("%3d. %v\n", index, task.Title)
	}
	var taskIndex int
	_, err := fmt.Scanf("%d", &taskIndex)
	if err != nil || taskIndex >= len(taskSlugs) || taskIndex < 0 {
		return "", errors.New("Invalid task selection")
	}
	return taskSlugs[taskIndex], nil
}

func removeTask(slug string) {
	task := loadTask(slug)
	os.Remove("tasks/" + slug + ".json")
	fmt.Printf("Marked %v as done", task.Title)
}

func findTasks(search string) []string {
	slug := getTaskSlug(search)
	taskSlugs := getTaskNames()
	var tasks []string
	for _, taskSlug := range taskSlugs {
		if strings.HasPrefix(taskSlug, slug) {
			tasks = append(tasks, taskSlug)
		}
	}
	return tasks
}

func commandList() {
	fmt.Println("Here are your tasks:")
	taskNames := getTaskNames()
	for _, name := range taskNames {
		t := loadTask(name)
		fmt.Println("  " + t.Title)
		fmt.Println("  " + t.Description)
		fmt.Println()
	}
}

func getTaskNames() []string {
	var tasks []string

	files, _ := ioutil.ReadDir("tasks")
	for _, file := range files {
		if file.Name()[0] != '.' {
			tasks = append(tasks, strings.TrimSuffix(file.Name(), ".json"))
		}
	}

	return tasks
}
