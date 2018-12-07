package main
import(
  "testing"
  pb "package1"
  "golang.org/x/net/context"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"

)
func TestMessages(t *testing.T) {

    // Set up a connection to the Server.
    const address = "localhost:50051"
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        t.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewGreeterClient(conn)

    // Test SayHello
    t.Run("SayHello", func(t *testing.T) {
        name := "world"
        r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
        if err != nil {
            t.Fatalf("could not greet: %v", err)
        }
        t.Logf("Greeting: %s", r.Message)
        if r.Message != "Hello "+name {
            t.Error("Expected 'Hello world', got ", r.Message)
        }

    })
}
