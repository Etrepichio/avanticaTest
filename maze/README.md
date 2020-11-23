# Synack Challenge

## Maze CRUD App

This little application provides an API Rest server with CRUD operations on Spots, Paths and Quadrants. 
Considerations:
- A path distance depends on the points it connects. Even if I'm saving the distance attribute as requested,
Im recalculating it everytime it's required, since one of the points may have its value changed during
that period of time.
- Each quadrant depends on the origin of the coordinates. That's why I've made CRUD operations around the origin. 
Once we have the origin, we can, for example, calculate which spots belong to a determined quadrant.

## Deployment
- The application is deployed locally using docker compose. 
Run docker-compose build and then docker-compose up to run both the application and database together.

## API

### Spots
Create Spot - POST
- Endpoint: /spot
- Payload:
{
    "x_coordinate": -4,
    "y_coordinate": -5,
    "name": "treasure",
    "number": 600
}

Get Single Spot - GET
- Endpoint: /spot/{id}

Modify Single Spot - PUT
- Endpoint: /spot/{id}
- Payload:
{
    "x_coordinate": -4,
    "y_coordinate": -5,
    "name": "treasure",
    "number": 600
}

Get Spots - GET
- Endpoint: /spots 

Delete Spot - DELETE
- Endpoint: /spot/{id}

### Paths
Create Path - POST
- Endpoint: /spot
- Payload:
{
    "point_a": "5fbb3712e3c84f4e02ff4e31",
    "point_b": "5fbb4b798edc5836096f87ea"
}

Get Single Path - GET
- Endpoint: /path/{id}

Modify Single Path - PUT
- Endpoint: /path/{id}
- Payload:
{
    "point_a": "5fbb3712e3c84f4e02ff4e31",
    "point_b": "5fbb4b798edc5836096f87ea"
}

Get Paths - GET
- Endpoint: /paths 

Delete Path - DELETE
- Endpoint: /path/{id}

### Origin
Create Origin - POST (Only one)
- Endpoint: /origin
- Payload:
{
    "x_origin": 4,
    "y_origin": 5
}

Get Origin - GET
- Endpoint: /origin

Modify Origin - PUT
- Endpoint: /origin
- Payload:
{
    "x_origin": 4,
    "y_origin": 5
} 

Delete Origin - DELETE
- Endpoint: /origin


Quadrant Spots - POST
- Endpoint: /quadrantSpots
- Payload:
{
    "name": "bottom_left"
}


# Improvements
- This whole application could be improved in several ways. For example, I'd change the quadrantSpots method to a GET
and add the name as a query parameter possibly,since it's a POST, and it isn't creating anything, but I made it like this 
because it needed to send some data.
- The codebase could be much better. I've done this quickly and even delayed a bit in the deployment construction.
(something as small as forgetting to connect the mongodb to the docker network). 
I've made some simple handlers to manage the different endpoints, but I don't really like to have transport mixed
with functionality. I'd separate the logic from the transport first and foremost, ideally using Go-Kit to have 
endpoints separated from the transport, decoding and encoding stuff, and all that aside from the actual logic. 
- Another improvement would be to define a proper db interface, so I don't reuse so much all the db interactions. 
- Needed to add logs (the logger is there)
- In short, the code is too tight, it would be better to separate the concerns to improve its flexibility looking forward
to future changes, also that would improve legibility(sorry!)
- Most of these things to improve are there because I'd need more time. I just wanted to develop something in a couple hours,
but the deployment part delayed me a bit more than expected.

###Thanks!
- Thanks for the opportunity! I really appreciate your time reviewing this!