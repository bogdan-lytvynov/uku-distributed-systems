curl -H "Content-Type: application/json" -XPOST -d '{"message":"first"}' localhost:8080/message
curl -H "Content-Type: application/json" -XPOST -d '{"message":"second"}' localhost:8080/message
curl -H "Content-Type: application/json" -XPOST -d '{"message":"third"}' localhost:8080/message

echo "Check replica 1\n"
curl localhost:3000/logs
echo "\n"

echo "Check replica 2\n"
curl localhost:3001/logs
echo "\n"

echo "Check replica 3 with delay\n"
curl localhost:3002/logs
echo "\n"

sleep 20

echo "Check replica 3 after delay\n"
curl localhost:3002/logs


