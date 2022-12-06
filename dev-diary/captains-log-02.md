# Dev Diary 02

## Problem

For handling changes to user following/followers, do we want to use a Method on the User struct, or a function in the DB?

## Pretext

Each user record maintains a list of followers (users who follow them) and following (users they themselves follow) as a slice of pointers; the rationale for this decision was to make better use of memory than having duplicate strings existing everywhere (while also offering some extensibility later if you liked e.g. extending the User record with more fields would implicitly extend user data retrieval functionality) - every user will exist in-memory anyway, so it seems more efficient to use this existing  data structure than try to engineer some "simple" solution which ends up degrading performance over time.

This does lead to an interesting edge case however; the structure of the `db` package implies that it will be this which manages the data - a request to follow a user goes via the DB, to establish a new link. However, a method on the User struct would allow user records to be relatively self-maintaining in a simple and clear way.

## Solution

It does make sense to have a method on the User struct which handles adding and removing followers; since the User holds a record of pointers, it makes sense for that record to update itself with assistance from the `db`'s operations queue. The DB is used to order queries chronologically, as the channel-based map structure means reading data is instant and free, while adjusting data is done in the correct order.