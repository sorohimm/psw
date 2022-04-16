Markdown Cheatsheet<a name="TOP"></a>
===================

- - - - 
# Heading 1 #

    Markup :  # Heading 1 #

    -OR-

    Markup :  ============= (below H1 text)

## Heading 2 ##

    Markup :  ## Heading 2 ##

    -OR-

    Markup: --------------- (below H2 text)
    
# psw
Два сервиса producer(загружает сообщения в очередь) и consumer(получает их через GET запрос) общаются через NATS JetStream

## How to run ##
nast-server -js 
cd ../../producer && go run .
cd consumer/cmd && go run .

## Queue
Забрать сообщения из очереди можно так: curl 'http://127.0.0.1:8081/quotes/v1/all'
