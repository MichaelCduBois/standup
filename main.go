package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type StandupItem struct {
	id        int8
	blocker   bool
	item      string
	yesterday bool
}

func main() {
	// Check or Create .standup Directory
	homeDir, err := os.UserHomeDir()
	checkErr(err)
	standupDir := homeDir + "/.standup"
	err = os.Mkdir(standupDir, 0700)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

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

	// Create Database Connection
	db, err := sql.Open("sqlite3", standupDir+"/data.db")
	checkErr(err)
	defer db.Close()

	sqlQuery := `
    CREATE TABLE IF NOT EXISTS notes (
      id        INTEGER PRIMARY KEY,
      blocker   BOOLEAN,
      item      STRING  NOT NULL,
      yesterday BOOLEAN
    )
  `
	_, err = db.Exec(sqlQuery)
	checkErr(err)

	switch {
	case *add:
		standupItem := StandupItem{}
		if *blocker {
			standupItem.blocker = true
		}
		if *yesterday {
			standupItem.yesterday = true
		}
		standupItem.item = item

		sqlQuery = `
      INSERT INTO notes (
        id,
        blocker,
        item,
        yesterday
      ) VALUES (
        ?,
        ?,
        ?,
        ?
      )
    `
		sqlStmt, _ := db.Prepare(sqlQuery)
		sqlStmt.Exec(
			nil,
			standupItem.blocker,
			standupItem.item,
			standupItem.yesterday,
		)
		defer sqlStmt.Close()
		return

	case *del:
		sqlQuery = "DELETE FROM notes WHERE id=?"
		sqlStmt, err := db.Prepare(sqlQuery)
		checkErr(err)
		defer sqlStmt.Close()
		_, err = sqlStmt.Exec(item)
		checkErr(err)
		return

	case *age:
		sqlQuery = "DELETE FROM notes WHERE yesterday=1"
		sqlStmt, err := db.Prepare(sqlQuery)
		checkErr(err)
		defer sqlStmt.Close()
		_, err = sqlStmt.Exec()
		checkErr(err)
		sqlQuery = "UPDATE notes SET yesterday=1"
		sqlStmt, err = db.Prepare(sqlQuery)
		checkErr(err)
		defer sqlStmt.Close()
		_, err = sqlStmt.Exec()
		checkErr(err)
		return

	case *list:
		rows, err := db.Query("SELECT id, item FROM notes")
		defer rows.Close()
		err = rows.Err()
		checkErr(err)
		items := make([]StandupItem, 0)
		for rows.Next() {
			item := StandupItem{}
			err = rows.Scan(&item.id, &item.item)
			checkErr(err)
			items = append(items, item)
		}
		err = rows.Err()
		checkErr(err)
		if len(items) > 0 {
			fmt.Println("# -- Standup Item")
			fmt.Println("=================")
			for _, item := range items {
				fmt.Printf("%v -- %v\n", item.id, item.item)
			}
		} else {
			fmt.Println("No Standup Items. Please use the '-add' flag.")
		}
		return

	case *reset:
		db.Exec("DELETE FROM notes")
		return
	}

	// Generate Standup Notes
	rows, err := db.Query("SELECT * FROM notes")
	defer rows.Close()
	err = rows.Err()
	checkErr(err)
	blockerItems := make([]StandupItem, 0)
	items := make([]StandupItem, 0)
	yesterdayItems := make([]StandupItem, 0)
	for rows.Next() {
		item := StandupItem{}
		err = rows.Scan(&item.id, &item.blocker, &item.item, &item.yesterday)
		checkErr(err)
		switch {
		case item.blocker:
			blockerItems = append(blockerItems, item)
			break
		case item.yesterday:
			yesterdayItems = append(yesterdayItems, item)
			break
		case !item.blocker && !item.yesterday:
			items = append(items, item)
			break
		}
	}
	err = rows.Err()
	checkErr(err)
	if len(blockerItems) > 0 {
		fmt.Println("##### Blockers #####")
		for _, blocker := range blockerItems {
			fmt.Printf("- %v\n", blocker.item)
		}
		if len(yesterdayItems) > 0 || len(items) > 0 {
			fmt.Println()
		}
	}
	if len(yesterdayItems) > 0 {
		fmt.Println("##### Yesterday #####")
		for _, yesterday := range yesterdayItems {
			fmt.Printf("- %v\n", yesterday.item)
		}
		if len(items) > 0 {
			fmt.Println()
		}
	}
	if len(items) > 0 {
		fmt.Println("##### Today #####")
		for _, item := range items {
			fmt.Printf("- %v\n", item.item)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
