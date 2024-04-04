# Behavior design of Chat Engine

**App Init - Sequence Diagram:**
1. App stablishes a WebSocket connection to the ChatEngine.
2. Sends a msg to the ChatEngine to retrieve the list of chat rooms the user has joined.
3. ChatEngine responds with the list of chat rooms.
4. App checks the cached data for the chat rooms.
    - If New Rooms Found: App sends msg to retrieve the messages for the new chat rooms.
    - If Rooms Left: App removes the chat rooms that the user has left from the cached data.
[for each existing chat room]
5. App sends a msg to retrieve the latest messages by providing the last message timestamp.
    - If no timestamp, server returns the last X messages.
    - If has timestamp, server returns the messages after that timestamp.
6. Server responds with the latest messages.
[end for]
7. App listens for new messages and updates the chat rooms in parallel with other app activities.

* The chat room data should be in memory until the user closes the app or logs out.

**Driver Creates Route - Sequence Diagram:**
1. Driver creates a route, App requests RouteAPI to create a route.
2. RouteAPI creates a route and returns the route ID.
3. RouteAPI requests ChatEngine to create a chat room for the route with the route ID and the driver's user ID.
4. ChatEngine acknowledges RouteAPIs message.
    - If the ACK fails the RouteAPI should retry the request.
6. ChatEngine creates a chat room for the route
7. ChatEngine joins the driver to the chat room 
8. ChatEngine sends a msg through WS to notify the driver joined a chat room.
9. App adds the room to the cached data.

**User Joins Route - Sequence Diagram:**
1. User joins a route, requests RouteAPI to join the route.
2. RouteAPI joins the user to the route and returns the route ID.
3. RouteAPI requests ChatEngine to join the user to the chat room for the route with the route ID and the user's user ID.
4. ChatEngine acknowledges RouteAPIs message.
    - If the ACK fails the RouteAPI should retry the request.
5. ChatEngine joins the user to the chat room.
6. ChatEngine sends a msg through WS to notify the user joined a chat room.
7. App adds the room to the cached data.

**User Leaves Route - Sequence Diagram:**
1. User leaves a route, requests RouteAPI to leave the route.
2. RouteAPI leaves the user from the route.
3. RouteAPI requests ChatEngine to leave the user from the chat room for the route with the route ID and the user's user ID.
4. ChatEngine acknowledges RouteAPIs message.
    - If the ACK fails the RouteAPI should retry the request.
5. ChatEngine leaves the user from the chat room.
6. ChatEngine sends a msg through WS to notify the user left the chat room.
7. App removes the room from the cached data.

Task Inception:

### Phase 1: Basic Chat Application of One Room
[BACK] 1.0 Set up the project
[BACK] 1.1 Implement basic WebSocket server
[BACK] 1.2 Implement basic message broadcasting
[BACK] 1.3 Implement tests the basic chat functionality

### Phase 2: Multi-room & 1 to 1 Chats
[BACK] 2.0 Implement room support
[BACK] 2.1 Implement room message broadcast
[BACK] 2.2 Implement 1 to 1 chats
[BACK] 2.3 Test multi-room and 1 to 1 Chats

### Phase 3: Using Redis Pub/Sub for Scalability
[BACK] 3.0 Set up Redis
[BACK] 3.1 Implement persistance layer
[BACK] 3.2 Implement Redis Pub/Sub handling

### Phase 4: Route Service Integration
[BACK] 4.0 Setup Swagger docs
[BACK] 4.0 Implement Create room API endpoint
[BACK] 4.1 Implement Join room API endpoint
[BACK] 4.2 Implement Leave room API endpoint

### Phase 5: Authentication
TBD 