// Package cmd provides commands for RSS feed aggregation.
package cmd

import (
	"os"
	"context"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"gator/rss"
)

type State struct {
	Config *config.Config
	DB *database.Queries
}

type Commands struct {
	Commands map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	function, ok := c.Commands[cmd.Name]
	if !ok { return fmt.Errorf("unknown command: %s", cmd.Name) }
	return function(s, cmd)
}

func (c *Commands) Register(name string, function func(*State, Command) error) {
	c.Commands[name] = function
}


type Command struct {
	Name string
	Arguments []string
}

func HandlerLogin(state *State, cmd Command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("missing username")
	}
	_, err := state.DB.GetUser(context.Background(), cmd.Arguments[0])
	if err != nil { 
		fmt.Printf("error getting user: %v", err) 
		os.Exit(1)
	}

		err = state.Config.SetUser(cmd.Arguments[0])
	if err != nil { 
		fmt.Printf("error setting user: %v", err) 
	os.Exit(1)
	}
	fmt.Printf("User set to %s\n", cmd.Arguments[0])
	return nil
}

func HandlerRegister(state *State, cmd Command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("missing username")
	}

	user, err := state.DB.CreateUser(
		context.Background(), 	
		cmd.Arguments[0],
	)
	if err != nil { 
		fmt.Printf("error creating user: %v", err) 
		os.Exit(1)
	}
	err = state.Config.SetUser(user.Name)
	if err != nil { 
		fmt.Printf("error setting user: %v", err) 
		os.Exit(1)
	}
	fmt.Printf("User created: %s\n", user)
	return nil
}


func HandlerReset(state *State, cmd Command) error {
	err := state.DB.TruncateUsers(context.Background())
	if err != nil { 
		fmt.Printf("error truncating users: %v", err) 
		os.Exit(1)
	}
	fmt.Printf("Users truncated\n")
	return nil
}


func HandlerListUsers(state *State, cmd Command) error {
	users, err := state.DB.GetUsers(context.Background())
	if err != nil { 
		fmt.Printf("error listing users: %v", err)
		os.Exit(1)
	}
	currentUser := state.Config.CurrentUser
	for _, user := range users {
		if user == currentUser {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}


func HandlerAgg(state *State, cmd Command) error {
	// if len(cmd.Arguments) < 1 {
	// 	return fmt.Errorf("missing feed url")
	// }
	_, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil { 
		fmt.Printf("error fetching feed: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Feed fetched\n")
	return nil
}


func HandlerAddFeed(state *State, cmd Command) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("missing feed url or name")
		os.Exit(1)
	}
	feedName := cmd.Arguments[0]
	feedURL := cmd.Arguments[1]
	currentUser, err := state.GetCurrentUser()
	if err != nil { 
		fmt.Printf("error getting current user: %v", err)
		os.Exit(1)
	}
	_, err = state.DB.CreateFeed(
		context.Background(), 
		database.CreateFeedParams{
			Name: feedName,
			Url: feedURL,
			UserID: currentUser.ID,
		},
	)
	if err != nil {
		fmt.Printf("error creating feed: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Feed created: %s\n", feedName)
	return nil
}	


func HandlerListFeeds(state *State, cmd Command) error {
	feeds, err := state.DB.GetFeeds(context.Background())
	if err != nil { 
		fmt.Printf("error listing feeds: %v", err)
		os.Exit(1)
	}
	for _, feed := range feeds {
		fmt.Printf("* %s %s %s\n", feed.Name, feed.Url, feed.UserName)
	}
	return nil
}


func (state *State) GetCurrentUser() (*database.User, error) {
	currentUser, err := state.DB.GetUser(
		context.Background(), state.Config.CurrentUser)
	if err != nil {
		return nil, fmt.Errorf("error getting current user: %v", err)
	}	
	return &currentUser, nil
}
