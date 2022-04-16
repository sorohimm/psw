nats-server -js &

cd producer 

go run . &

cd ../consumer/cmd 

go run . &
