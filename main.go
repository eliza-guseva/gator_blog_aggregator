package main

import _ "github.com/lib/pq"

import (
	"fmt"
	"os"
	"database/sql"
	"gator/internal/database"
	"gator/internal/config"
	"gator/internal/cmd"
)

func main() {
	myConfig, err := config.Read()
	if err != nil { panic(err) }
	
	db, err := sql.Open("postgres", myConfig.DBUrl)
	if err != nil { panic(err) }
	defer db.Close()

	dbQueries := database.New(db)
	state := cmd.State{
		Config: myConfig,
		DB: dbQueries,
	}

	commands := cmd.Commands{
		Commands: map[string]func(*cmd.State, cmd.Command) error{
			"login": cmd.HandlerLogin,
			"register": cmd.HandlerRegister,
			"reset": cmd.HandlerReset,
			"users": cmd.HandlerListUsers,
			"agg": cmd.HandlerAgg,
			"addfeed": cmd.HandlerAddFeed,
			"feeds": cmd.HandlerListFeeds,
		},
	}
	args := make([]string,0)
	args = append(args, os.Args[1:]...)
	if len(args) < 1 {
		fmt.Printf("You need to provide at least one command and an argument")
		os.Exit(1)
	}
	err = commands.Run(&state, cmd.Command{Name: args[0], Arguments: args[1:]})
	if err != nil { 
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Done\n")
}
