# psw
Два сервиса producer(загружает сообщения в очередь) и consumer(получает их через GET запрос) общаются через NATS JetStream

## How to run ##
    nast-server -js
    
    cd producer && go run .
    
    cd ../../consumer/cmd && go run .

## Queue
Забрать сообщения из очереди можно так: 
    `curl 'http://127.0.0.1:8081/quotes/v1/all'`
