curl -H "Content-Type: application/json" -XPOST -d '{"message":"first", "w":1}' localhost:8080/message
curl -H "Content-Type: application/json" -XPOST -d '{"message":"second", "w":2}' localhost:8080/message
curl -H "Content-Type: application/json" -XPOST -d '{"message":"third", "w":3}' localhost:8080/message

echo "Check replica 1\n"
curl localhost:3000/log
echo "\n"

echo "Check replica 2\n"
curl localhost:3001/log
echo "\n"

echo "Check replica 3 with delay\n"
curl localhost:3002/log
echo "\n"

sleep 20

echo "Check replica 3 after delay\n"
curl localhost:3002/log


