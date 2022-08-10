# Go To-do App

I built a simple CLI To-do app in Go while streaming live on 
[YouTube](https://youtu.be/FwZKh3Fq8t0)!

## Usage

To list the to-do's run:
```sh
go run todo.go
```

To create a to-do run the following command and follow it's instructions:
```sh
go run todo.go create
```

To mark a to-do as done run the following command with the name of the task to 
complete:
```sh
go run todo.go done name
```
Replace spaces in the name with dashes. It will also search for all tasks that have 
a name that start with whatever you entered for the name. 
