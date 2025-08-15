// Package cmd provides commands for RSS feed aggregation.
package cmd

import (
	"os"
	"context"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"gator/rss"
	"github.com/google/uuid"
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


func HandlerAddFeed(state *State, cmd Command, user *database.User) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("missing feed url oR name")
		os.Exit(1)
	}
	feedName := cmd.Arguments[0]
	feedURL := cmd.Arguments[1]
	
	_, err := state.DB.CreateFeed(
		context.Background(), 
		database.CreateFeedParams{
			Name: feedName,
			Url: feedURL,
			UserID: user.ID,
		},
	)
	if err != nil {
		fmt.Printf("error creating feed: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Feed created: %s\n", feedName)
	
	feed := state.GetFeedByURL(feedURL)
	_ = state.CreateFeedFollow(user.ID, feed.ID)

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

func HandlerFollow(state *State, cmd Command, user *database.User) error {
	
	feed := state.GetFeedByURL(cmd.Arguments[0])
	
	follows := state.CreateFeedFollow(user.ID, feed.ID)
	fmt.Println(follows)
	return nil
}

func HandlerListUserFollows(state *State, cmd Command, user *database.User) error {
	follows, err := state.DB.GetFeedFollowsForUser(
		context.Background(),
		user.ID,
	)
	if err != nil { 
		fmt.Printf("error listing user follows: %v", err)
		os.Exit(1)
	}
	for _, follow := range follows {
		fmt.Printf("* %s %s\n", follow.FeedName, follow.UserName)
	}
	return nil
}


func HandlerUnfollow(state *State, cmd Command, user *database.User) error {
	feedURL := cmd.Arguments[0]
	feed := state.GetFeedByURL(feedURL)

	err := state.DB.UnfollowFeed(
		context.Background(),
		database.UnfollowFeedParams{
			UserID: user.ID,
			FeedID: feed.ID,
		},
	)
	if err != nil { 
		fmt.Printf("error unfollowing feed: %v", err)
		os.Exit(1)
	}
	return nil
}


// helpers

func MiddlewareLoggedIn(handler func(state *State, cmd Command, user *database.User) error) func(*State, Command) error {
	return func(state *State, cmd Command) error {
		currentUser := state.GetCurrentUser()
		return handler(state, cmd, currentUser)
	}
}


func (state *State) GetCurrentUser() *database.User {
	currentUser, err := state.DB.GetUser(
		context.Background(), state.Config.CurrentUser)
	if err != nil {
		fmt.Printf("error getting current user: %v", err)
		os.Exit(1)
	}	
	return &currentUser
}

func (state *State) GetFeedByURL(url string) *database.Feed {
	feed, err := state.DB.GetFeed(context.Background(), url)
	if err != nil { 
		fmt.Printf("error getting feed: %v", err)
		os.Exit(1)
	}
	return &feed
}

func (state *State) CreateFeedFollow(userID uuid.UUID, feedID uuid.UUID) *database.CreateFeedFollowRow {
	follows, err := state.DB.CreateFeedFollow(
		context.Background(), 
		database.CreateFeedFollowParams{
			UserID: userID,
			FeedID: feedID,
		},
	)
	if err != nil {
		fmt.Printf("error creating feed follow: %v", err)
		os.Exit(1)
	}
	return &follows
}

