//client de Clientes
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"papa.com/Clientes/chat"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("dist49:9000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ingrese una opción \n1.Ver estado Finanzas\n2.Salir")
	log.Printf("---------------------")

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.ToLower(strings.Trim(text, " \r\n"))

		if strings.Compare(text, "1") == 0 {
			message := chat.Message{
				Body: "Holi",
			}
			response, err := c.ConsultarFinanzas(context.Background(), &message)
			if err != nil {
				log.Fatalf("We couldn't say hello: %s", err)
			}
			//log.Printf(response.Body)

			lines := strings.Split(response.Body, "\n")

			//Escribe la respuesta en un csv
			csvfile, err := os.OpenFile("finanzas.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				log.Fatal(err)
			}
			csvwriter := csv.NewWriter(csvfile)

			for writethis := 0; writethis != len(lines); writethis++ {
				csvline := strings.Split(lines[writethis], ",")
				csvwriter.Write(csvline)
			}

			csvwriter.Flush()
			csvfile.Close()

		}
		//Termina de leerlo

		if strings.Compare(text, "2") == 0 {
			break
		}

	}

}
