# ppf-chat-engine
The PPF chat engine, from scratch

1. Service Starts
2. Route is created at RouteAPI
3. RouteAPI sends a request to the ChatEngine API to create a room (attaches the driver id)
4. ChatEngine requests the driver's information from UserAPI and stores it
5. ChatEngine creates a room with the driver's information