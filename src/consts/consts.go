package constants

import(
	"net/rpc"
	"errors"
	"time"
	"fmt"
)

const RpcWaitMillis = 1000 // Milliseconds to wait when an RPC fails
const RpcTestMillis = 2000 // Milliseconds to wait when Testing Connection fails
const RpcTries = 5

const SAMPLE_STORAGEHP = "localhost:9321"

func DialToServer(targetHP string) (error, *rpc.Client) {
	// Keep Dialing until connected!
	count := 0
	for count < RpcTries {
		conn, err := rpc.DialHTTP("tcp", targetHP)
		if err == nil {
			return nil, conn
		}
		time.Sleep(time.Duration(RpcWaitMillis) * time.Millisecond)
		fmt.Println("Connecting.... trying again")
		count++
	}
	if count == RpcTries {
		return errors.New("Too many tries to connect to server"), nil
	}
	return nil, nil
}