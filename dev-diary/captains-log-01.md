# Dev Diary 01

## Problem

What is the best way to maintain lists of Followers & Following in respective User records?

## Pretext

I have decided to go with an in-memory map for storage of users and records, using a custom type as defined in `internal/data`; in-memory maps are free to access records from and extremely fast when storing records using a channel-based concurrent architecture.

## Solutions

But what would be the best form for a user to have its Following & Followers lists stored?:

- A list of strings per record, listing all the users usernames (or IDs if you were so inclined but I, dear reader, am NOT) but this could be extremely expensive and hard to maintain links between User records with.
- Using pointers, while somewhat complicated to implement and requiring an additional layer to process the data returned from the database (de-referencing each user record and extracting the usernames for return to the calling user) is *somewhat* cumbersome, but I think the performance trade-off of returning a list of users is worth making said list of users easier to manage internally. This also further de-couples internal processing services (the database is solely concerned with returning its data; other services do what they like with it).

## Decision

In case it wasn't obvious by my labour of love written about the pointer-based option, I've gone for the pointer-based option. Users are stored in memory with an address, so each User can store a list of User pointers: every record will get updated implicitly when a change is made to any linked records, it encourages you to further de-couple the layers in your application code, and the performance knock is (hopefully) somewhat trivial.