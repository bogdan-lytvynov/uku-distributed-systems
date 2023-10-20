### Second iteraction of the module 1

In this iteraction the leader wait only for W ACKs received from leader + replicas.
This implementation also respects order of messages and includes order nubmer in the RPC call from leader to replicas.
Leader is responisble for order number increment.

If you want to quickly test this code just run ./test.sh
