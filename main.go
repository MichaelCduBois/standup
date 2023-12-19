package main

import (
	"flag"
	"fmt"
)

func main() {
	// Setup Flags
	add := flag.Bool("add", false, "Add item to standup notes")
	age := flag.Bool("age", false, "Age standup items creating new day")
	del := flag.Bool("delete", false, "Delete item from standup notes")
	list := flag.Bool("list", false, "List all standup items")
	reset := flag.Bool("reset", false, "Delete all standup items")
	// Initialize Add Options
	blocker := flag.Bool("blocker", false, "Mark item as Blocker")
	yesterday := flag.Bool("yesterday", false, "Add item to Standup Notes for previous day")
	// Parse Flags
	flag.Parse()

	item := flag.Arg(0)

	switch {
	case *add:
		if *blocker {
			fmt.Printf("Adding %v as Blocker\n", item)
		}
		if *yesterday {
			fmt.Printf("Adding %v for Yesterday\n", item)
		}
		fmt.Printf("Adding %v\n", item)
		return
	case *del:
		fmt.Printf("Deleting %v\n", item)
		return
	case *age:
		fmt.Println("Aging all standup items")
		return
	case *list:
		fmt.Println("Listing all standup items")
		return
	case *reset:
		fmt.Println("Deleting all standup items")
		return
	}
	fmt.Println("Generating Standup Notes")
}
