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
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)
	log.Printf("Ingrese una opción \n1. Subir archivo pymes.csv\n2.Subir archivo retail.csv")
	log.Printf("---------------------")

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.ToLower(strings.Trim(text, " \r\n"))

		if strings.Compare(text, "1") == 0 {
			//Toda esta parte es para leer los archivos pymes.csv
			csvFile, err := os.Open("../pymes.csv")
			if err != nil {
				fmt.Println(err)
			}
			defer csvFile.Close()

			csvLines, err := csv.NewReader(csvFile).ReadAll()
			if err != nil {
				fmt.Println(err)
			}
			for _, line := range csvLines {
				message := chat.Orden{
					Id:          line[0],
					Producto:    line[1],
					Valor:       line[2],
					Tienda:      line[3],
					Destino:     line[4],
					Prioritario: line[5],
				}
				response, err := c.OrdenarPyme(context.Background(), &message)
				if err != nil {
					log.Fatalf("We couldn't say hello: %s", err)
				}

				log.Printf("Su codigo de seguimiento es: %s", response.Body)
			}
			//Termina de leerlo
		}
		if strings.Compare(text, "2") == 0 {
			//Toda esta parte es para leer los archivos retail.csv
			csvRetail, err := os.Open("../retail.csv")
			if err != nil {
				fmt.Println(err)
			}
			defer csvRetail.Close()

			RetailLines, err := csv.NewReader(csvRetail).ReadAll()
			if err != nil {
				fmt.Println(err)
			}
			for _, line := range RetailLines {
				messageRetail := chat.Orden{
					Id:          line[0],
					Producto:    line[1],
					Valor:       line[2],
					Tienda:      line[3],
					Destino:     line[4],
					Prioritario: "0",
				}
				response, err := c.OrdenarRetail(context.Background(), &messageRetail)
				if err != nil {
					log.Fatalf("We couldn't say hello: %s", err)
				}

				log.Printf("Su codigo de seguimiento es: %s", response.Body)
			}
			//Termina de leerlo
		}

	}

}
