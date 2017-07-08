# chatServerAsync
A simple async chatServer with limited API
The chat server will run on http://localhost:8081 and will support the following REST API: 

 1. GET /messages

     list 100 most recent messages, sorted by 'timestamp' posted to the chat server.

     example:
     ========

     ```
     curl -H "Content-Type: application/json" http://localhost:8081/messages

     {
       "messages: [
         {"timestamp": 1491345710.18, "user": "superman", "text": "hello"},
         {"timestamp": 1491345713.18, "user": "batman", "text": "hello"}
       ]
     }

     ```

 2. POST /message 

     a request to post the given message. 
     when the message is processed by the server a unix timestamp is recorded with each message.

     example:
     ========

     ```
     curl -X POST -H "Content-Type: application/json" --data '{"user":"superman", "text":"hello"}' http://localhost:8081/message

     {
       "ok": true
     }
     ```

 3. GET /users

     a request to return a set of users that have posted messages so far.

     example:
     ========

     ```
     curl -H "Content-Type: application/json" http://localhost:8081/users

     {
       "users": [
         "superman", "batman"
       ]
     }
     ```

the server should respond with 404 to all other requests not listed above
