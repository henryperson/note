package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"time"
)

var args = map[string]func(){
	"new":  createNote,
	"open": openInVSCode,
	"help": printHelp,
	"":     printHelp,
}

var usageString = `Usage: "note new" will create a note with today's date in a subdirectory within $HOME/notes/. It will be stored in the directory YYYY/MM.
"note open" will open today's note.`

func main() {
	flag.Parse()
	cmd := flag.Arg(0)
	if _, ok := args[cmd]; !ok {
		log.Fatalf("Invalid command: %v is not a valid command.\n", cmd)
	} else if len(flag.Args()) > 1 {
		log.Fatalf("Too many commands. Expected 1, got %v\n", len(flag.Args()))
	}
	args[cmd]()
}

func createNote() {
	now := time.Now()
	yearMonth := now.Format("2006/01")
	usr, err := user.Current()
	// TODO: Make an option to supply the notes directory as a flag.
	if err != nil {
		log.Fatal("Could not get the user (used to find the home directory): %w\n", err)
	}
	absPath := fmt.Sprintf("%v/notes/%v", usr.HomeDir, yearMonth)
	err = os.MkdirAll(absPath, 0755)
	if err != nil {
		log.Fatalf("Could not create directory %v: %v\n", absPath, err)
	}

	today := now.Format("02")
	fileStr := fmt.Sprintf("%v/%v.md", absPath, today)
	_, err = os.Stat(fileStr)
	if os.IsNotExist(err) {
		f, err := os.Create(fileStr)
		defer f.Close()
		if err != nil {
			log.Fatalf("Could not create file %v: %v\n", fileStr, err)
		}
		_, err = f.Write([]byte(fmt.Sprintf("# %v/%v\n", yearMonth, today)))
		if err != nil {
			log.Fatalf("Could not write to file %v: %v\n", fileStr, err)
		}
		fmt.Printf("Created file at %v/%v.md\n", yearMonth, today)
	} else if err != nil {
		log.Fatalf("Could not read status of file %v: %v\n", fileStr, err)
	} else {
		fmt.Printf("File exists for %v/%v\n", yearMonth, today)
	}
}

func openInVSCode() {
	now := time.Now()
	ymd := now.Format("2006/01/02")
	usr, err := user.Current()
	if err != nil {
		log.Fatal("Could not get the user (used to find the home directory).\n")
	}

	notesPath := fmt.Sprintf("%v/notes", usr.HomeDir)
	filePath := fmt.Sprintf("%v/%v.md", notesPath, ymd)
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("File does not exist, please create it with \"note new\"")
		return
	}

	cmd := exec.Command("code", notesPath, filePath, "-n")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Could not run command to open file: %v\n", err)
	}
}

func printHelp() {
	fmt.Println(usageString)
}
