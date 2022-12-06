# README

## API Specification

### POST /CreateNewUser

``` POST /createNewUser/{username} ```

Creates a user with a specific username

### GET /GetFollowers

``` GET /getFollowers/{username} ```

Returns a list of followers of a given user

### GET /GetFollowing

``` GET /getFollowing/{username} ```

Returns a list of users the given user is following

### PATCH /Follow

``` PATCH /follow -d '{ "username": <following-user>, "target": <target-user> }' ```

Follows the target-user as the following-user

#### Example

##### Request

```json
{
    "username": "deanfoley",
    "target": "domgreen",
}
```

##### Response

```"User deanfoley is now following domgreen!"```

### PATCH /Unfollow

``` PATCH /unfollow -d '{ "username": <following-user>, "target": <target-user> }' ```

Unfollows the target-user as the following-user

#### Example

##### Request

```json
{
    "username": "deanfoley",
    "target": "domgreen",
}
```

##### Response

```"User deanfoley is no longer following domgreen!"```

