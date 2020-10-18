package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"papa.com/Clientes/chat"
)

//Camion is
type Camion struct {
	Tipo   string
	cargo1 string
	cargo2 string
}

//Mandar is
func Mandar(camion Camion) {
	var conn *grpc.ClientConn
	mensaje := chat.Message{
		Body: "normal",
	}
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	c := chat.NewChatServiceClient(conn)
	response, err := c.RecibirPaquete(context.Background(), &mensaje)
	if err != nil {
		log.Fatalf("We couldn't say hello: %s", err)
	}
	fmt.Printf("Su id es: %s", response.Id)
}

func main() {

	CNormal := Camion{
		Tipo: "Normal",
	}
	/*CRetail1 := Camion{
		Tipo: "Retail",
	}*/
	/*CRetail2:=Camion{
		Tipo:"Retail"
	}*/

	//reader := bufio.NewReader(os.Stdin)
	/*fmt.Println("Tiempo de espera para los camiones?")
	espera, _ := reader.ReadString("\n")
	fmt.Println("Tiempo de envio Paquete")
	tenvio, _ := reader.ReadString('\n')*/
	go Mandar(CNormal)
	time.Sleep(time.Second * 10)

}
