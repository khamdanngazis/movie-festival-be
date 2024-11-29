## Admin APIs
### Create and Upload Movies

Endpoint: POST /admin/movies
Input: JSON (title, description, duration, artists, genres, watch_url)
Logic:
Validate input.
Save to movies table.
Controller: controllers/admin_movies.go

### Update Movie

Endpoint: PUT /admin/movies/:id
Input: JSON (same as create).
Logic:
Find movie by ID, update fields.
Controller: controllers/admin_movies.go

### See Most Viewed Movie and Genre

Endpoint: GET /admin/reports/views
Output: Most viewed movie, genre.
Logic:
Aggregate view counts, group by genres.
Controller: controllers/reports.go

## All Users APIs
### List Movies with Pagination

Endpoint: GET /movies
Query: page, limit.
Logic:
Fetch paginated movies from DB.
Controller: controllers/movies.go

### Search Movies

Endpoint: GET /movies/search
Query: query (matches title, description, artists, genres).
Logic:
Use SQL LIKE queries to fetch matches.
Controller: controllers/movies.go

### Track Viewership
Endpoint: POST /movies/:id/view
Logic:
Increment view count in movies if wact duration is > 60  menutes.
Optionally track viewerships.
Controller: controllers/movies.go

## Bonus Features
## Vote System
### Vote a Movie
Endpoint: POST /movies/:id/vote
Logic:
Ensure user hasn't voted before for the same movie.
Add entry in votes.
Controller: controllers/movies.go

### Unvote a Movie
Endpoint: DELETE /movies/:id/vote
Logic:
Remove vote from votes.

### List User's Voted Movies
Endpoint: GET /users/me/votes
Logic:
Fetch movies by user_id from votes.

## User Authentication

### Register
Endpoint: POST /auth/register
Logic:
Validate input, hash password, save to users.
Controller: controllers/auth.go

### Login
Endpoint: POST /auth/login
Logic:
Validate credentials, return JWT token.

### Logout
Endpoint: POST /auth/logout
Logic:
Invalidate JWT.
