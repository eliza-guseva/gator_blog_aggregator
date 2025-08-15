Assignment

    Create a feed_follows table with a new migration. It should:
        Have an id column that is a primary key.
        Have created_at and updated_at columns.
        Have user_id and feed_id foreign key columns. Feed follows should auto delete when a user or feed is deleted.
        Add a unique constraint on user/feed pairs - we don't want duplicate follow records.
    Add a CreateFeedFollow query. It will be a deceptively complex SQL query. It should insert a feed follow record, but then return all the fields from the feed follow as well as the names of the linked user and feed. I'll add a tip at the bottom of this lesson if you need it.
    Add a follow command. It takes a single url argument and creates a new feed follow record for the current user. It should print the name of the feed and the current user once the record is created (which the query we just made should support). You'll need a query to look up feeds by URL.
    Add a GetFeedFollowsForUser query. It should return all the feed follows for a given user, and include the names of the feeds and user in the result.
    Add a following command. It should print all the names of the feeds the current user is following.
    Enhance the addfeed command. It should now automatically create a feed follow record for the current user when they add a feed.

Run and submit the CLI tests.
