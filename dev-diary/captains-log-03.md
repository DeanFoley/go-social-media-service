# Dev Diary 3

## Problem

Should requests for data (GetFollowers/GetFollowing) use the message queue?

## Pretext

The in-memory data store is a map of Users where each User is stored under a key value of its username.

Maps are inexpensive to query. You can request a value from a map with a given key and it is instantaneous in Go.

However, in the way my program is structed, requests for data are still appended to the queue of requests. This does make the database idempotent, but also, you could make the argument that queries for data are unnecesarily queued behind queries for data *manipulation*.

## Decision

One of the motivators behind structuring the DB in the way that it is is that this provides a graceful shutdown mechanism for the in-memory store. This is important for any application hosted in Cloud servers; an application may need to scale, either horizontally  or vertically, at a moments notice to meet an increase in demand. If this is done in an ungraceful way (the server just shutting down to scale), some API requests might be lost in the chaos.

Using a message queue with the channels-based initialisation architecture means that, when a signal is sent to terminate the process, the process will not terminate until the queue has finished processing all of its messages. As  such, the server will only shut down once all current users have been served. This might even occur while additional infrastructure is being prepared behind the scenes, reducing the perceptible "downtime" of the application.

This is, as with all things, a trade-off; you might wish for an application to make read operations instantaneous, "jumping the queue" ahead of the requests which are changing the data. It is, in my opinion, a worthwhile trade-off.