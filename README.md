This Project is a simple Server implementation for a MessageProvider called Chirpy (Twitter clone). It provides the basic ablitiy to Add, Update and Delete a User, aswell for Messages.
This was a Guided exercise from boot.dev. It was a pleasure to learn how the basic aspects of a server are implemented in Go.
    
How to install and run your project: 
You need a Postgres Database. 
Create a .env File with the following Params
    - PORT= portnumber 
    - DB_URL= url to your database
    - PLATFORM= Dev for enabling reset of users via url link
    - SECRET= Secret Hash for Token generation
    - POLKA_KEY=Key for "POLKA_Service". Can be left empty

--Enpoints of this API. Avaiavle under .../admin/doc--

Entpoint Nr. 1
Name: /admin/metrics,
Method: GET,
QueryParams: [],
Authentification Required: false,
Input Body: null,
Output Body: null,
ErrorCodes: [],
Description: Enpoint for Server metrics

Entpoint Nr. 2
Name: /admin/reset,
Method: POST,
QueryParams: [],
Authentification Required: false,
Input Body: null,
Output Body: null,
ErrorCodes: [],
Description: Resets Counter for Server metrics. Deletes all users (only in dev)

Entpoint Nr. 3
Name: /api/chirps,
Method: POST,
QueryParams: [],
Authentification Required: true,
Input Body: {"body":"Text"},
Output Body: {"id":"00000000-0000-0000-0000-000000000000","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","body":"","user_id":"00000000-0000-0000-0000-000000000000"},
ErrorCodes: [500 400 201],
Description: Endpoint for Sending Message to Chirpy

Entpoint Nr. 4
Name: /api/chirps/,
Method: GET,
QueryParams: [author_id=UUID sort=asc|desc],
Authentification Required: false,
Input Body: null,
Output Body: [{"id":"00000000-0000-0000-0000-000000000000","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","body":"Text","user_id":"00000000-0000-0000-0000-000000000000"}],
ErrorCodes: [400 200 500],
Description: Endpoint for Getting Messages from User by Id

Entpoint Nr. 5
Name: /api/chirps/{chirpID},
Method: GET,
QueryParams: [chirpID=UUID],
Authentification Required: false,
Input Body: null,
Output Body: {"id":"00000000-0000-0000-0000-000000000000","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","body":"Text","user_id":"00000000-0000-0000-0000-000000000000"},
ErrorCodes: [400 500 404 200],
Description: Endpoint for Getting Messages by Chirp Id

Entpoint Nr. 6
Name: /api/chirps/{chirpID},
Method: DELETE,
QueryParams: [chirpID=UUID],
Authentification Required: true,
Input Body: null,
Output Body: null,
ErrorCodes: [400 500 404 403 204],
Description: Endpoint for Deöeting Messages by Chirp Id

Entpoint Nr. 7
Name: /api/refresh,
Method: POST,
QueryParams: [],
Authentification Required: true,
Input Body: null,
Output Body: {"token":""},
ErrorCodes: [401 500 200],
Description: Enpoint to refresh Access Token for a User

Entpoint Nr. 8
Name: /api/revoke,
Method: POST,
QueryParams: [],
Authentification Required: true,
Input Body: null,
Output Body: null,
ErrorCodes: [401 500 204],
Description: Enpoint to revoke provided refresh Token for a User

Entpoint Nr. 9
Name: /api/users,
Method: POST,
QueryParams: [],
Authentification Required: false,
Input Body: {"password":"","email":""},
Output Body: {"id":"00000000-0000-0000-0000-000000000000","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","email":"","is_chirpy_red":false},
ErrorCodes: [400 500],
Description: Enpoint to register a User

Entpoint Nr. 10
Name: /api/login,
Method: POST,
QueryParams: [],
Authentification Required: false,
Input Body: {"password":"","email":""},
Output Body: {"id":"00000000-0000-0000-0000-000000000000","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","email":"","token":"","refresh_token":"","is_chirpy_red":false},
ErrorCodes: [400 401 500 200],
Description: Entpoint for User to log in

Entpoint Nr. 11
Name: /api/users,
Method: PUT,
QueryParams: [],
Authentification Required: true,
Input Body: {"password":"","email":""},
Output Body: {"id":"00000000-0000-0000-0000-000000000000","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","email":"","is_chirpy_red":false},
ErrorCodes: [400 500 200],
Description: Entpoint for updating Users email and password

Entpoint Nr. 12
Name: /api/users,
Method: DELETE,
QueryParams: [],
Authentification Required: true,
Input Body: null,
Output Body: null,
ErrorCodes: [404 204],
Description: Entpoint for Deleting a User by its id

Entpoint Nr. 13
Name: /api/polka/webhooks,
Method: POST,
QueryParams: [],
Authentification Required: true,
Input Body: {"event":"","data":{"user_Id":"00000000-0000-0000-0000-000000000000"}},
Output Body: null,
ErrorCodes: [400 500 200],
Description: Webhook for changing Users subscription status