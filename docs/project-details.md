# 📦 Models

There are 3 models used in this project:

- user: See [user.go](../models/user.go)
- post: See [post.go](../models/post.go)
- comment: See [comment.go](../models/comment.go)

Each of them also inherit from the [base model](../models/base.go) which contains 3 base attributes:

```py
- ID          # PrimaryKey
- CreatedAt   # Time that entry was created
- UpdatedAt   # Time that entry was updated
```

## Relational diagram

![relational-diagram](images/relational-diagram.png)

<br>

# 🛣️ API Endpoints

API endpoints can be categorized into 2 access-levels:

1. `public`:

   - Does not require user authentication to access
   - Routes in this category are initialized in [public.go](../routes/public.go)

2. `protected`:
   - Requires user authentication for access (JWT token)
   - Routes in this category are initialized in [protected.go](../routes/protected.go)

There are 4 `domains` in this project which define all the available API endpoints.

These domains mirror the 4 [features](https://github.com/mfjkri/OneNUS/blob/master/docs/project-details.md#-features) in our frontend.

- [auth](../controllers/auth/)
- [posts](../controllers/posts/)
- [comments](../controllers/comments/)
- [users](../controllers/users/)

Below is a quick reference to the access level of each domain and the API endpoints they define:

- `auth`:

  ```py
  auth (public)
  ├── login       # Login of existing account
  ├── register    # Registration of new account
  └── me          # Authenticating an existing session using JWT token
  ```

- `posts`:

  ```py
  posts (protected)
  ├── get         # Fetches a list of posts based on given params
  ├── getbyid     # Fetches a single post based on ID (if any)
  ├── create      # Creates a new post
  ├── updatetext  # Updates an existing post text
  └── delete      # Deletes an existing post
  ```

- `comments`:

  ```py
  comments (protected)
  ├── get         # Fetches a list of comments from given postID
  ├── create      # Creates a new comment
  ├── updatetext  # Updates an existing comment text
  └── delete      # Deletes an existing comment
  ```

- `users`:

  ```py
  users (protected)
  ├── getbyid     # Fetches user details based on ID (if any)
  ├── updatebio   # Update user bio
  └── delete      # Deletes user account
  ```

<br>

# 🎮 Controllers

Each domain has its own controller with the following directory structure:

```sh
src
├── routes.go   # Define all the endpoints in the domain
├── handler.go  # Handler functions for each endpoint
└── helper.go   # Helper functions used by handler functions
```

## Handler functions

Each handler function follows a consistent style that is easy to follow when adding more endpoints.

```go
type ExpectedRequestTypeForHandler {
   RequestParam1  paramType `uri:"requestParam1" json:"requestParam1" binding:"required"`
   // ...
}

func Handler(c *gin.Context) {
   // ...
}
```

An example from `CreatePost` handler in [posts.go](../controllers/posts/handler.go):

```go
type CreatePostRequest struct {
  Title string `json:"title" binding:"required"`
  Tag   string `json:"tag" binding:"required"`
  Text  string `json:"text" binding:"required"`
}

func CreatePost(c *gin.Context) {
	// Check that RequestUser is authenticated
  user, found := VerifyAuth(c)
  if found == false {
    return
  }

  // Parse RequestBody
  var json CreatePostRequest
  if err := c.ShouldBindJSON(&json); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
    return
  }

  // Request params are now accessible through:
  // json.Title, json.Tag, json.Tag

  // Rest of Handler logic
  // ...
}
```

### JSON vs URI

- All request params for `POST` requests are included in the `JSON` request body.

- Meanwhile, all request params for `GET` and `DELETE` requests are included in the URI instead.

  `ginContext.ShouldBindJSON` and `ginContext.ShouldBindUri` are used to bind and verify the request params accordingly.
