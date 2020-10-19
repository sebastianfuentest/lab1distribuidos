package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"papa.com/Clientes/chat"
)

//Camion is
type Camion struct {
	nombre     string
	conpaquete int
	Tipo       string
	tEspera    int
	tEntre     int
	cargo1     Paquete
	cargo2     Paquete
}

//Paquete is
type Paquete struct {
	id          string
	seguimiento string
	tipo        string
	valor       string
	intentos    int32
	estado      string
}

//EnRuta funcion auxiliar para que los camiones busquen paquetes en un loop
func EnRuta(camion Camion) {
	log.Printf("entro a la")
	for {
		if camion.conpaquete == 0 {
			camion.conpaquete = 1
			log.Printf("Mandando camion %s \n", camion.nombre)
			camion.conpaquete = Mandar(camion)
			time.Sleep(time.Duration(5) * time.Second)
			log.Printf("conpaquete igual a %b ", camion.conpaquete)

		}

	}
}

//Mandar is
func Mandar(camion Camion) (ret int) {
	var i int
	camion.cargo1.id = "nada"
	camion.cargo2.id = "nada"
	var conn *grpc.ClientConn
	mensaje := chat.Message{
		Body: camion.Tipo,
	}
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	c := chat.NewChatServiceClient(conn)
	//pido 2 paquetese altiro
	response, err := c.RecibirPaquete(context.Background(), &mensaje)
	if err != nil {
		log.Fatalf("We couldn't say hello: %s", err)
	}
	time.Sleep(time.Duration(camion.tEspera) * time.Second)
	response2, err := c.RecibirPaquete(context.Background(), &mensaje)
	if err != nil {
		log.Fatalf("We couldn't say hello: %s", err)
	}

	appPaquete := Paquete{
		id:          response.Id,
		seguimiento: response.Seguimiento,
		tipo:        response.Tipo,
		valor:       response.Valor,
		intentos:    response.Intentos,
		estado:      response.Estado,
	}
	appPaquete2 := Paquete{
		id:          response2.Id,
		seguimiento: response2.Seguimiento,
		tipo:        response2.Tipo,
		valor:       response2.Valor,
		intentos:    response2.Intentos,
		estado:      response2.Estado,
	}

	if camion.cargo1.id == "nada" {
		camion.cargo1 = appPaquete
		fmt.Printf("camion tipo %s lleva cargo %s\n ", camion.Tipo, camion.cargo1.id)
	}

	if camion.cargo2.id == "nada" {

		camion.cargo2 = appPaquete2
		fmt.Printf("camion tipo %s lleva cargo %s\n", camion.Tipo, camion.cargo2.id)
	}

	if camion.cargo2.id != "NOHAY" {

		if camion.cargo1.valor > camion.cargo2.valor {
			for camion.cargo1.intentos < 3 {
				i = rand.Intn(100)
				if i >= 80 {
					camion.cargo2.intentos = camion.cargo2.intentos + 1
					fmt.Println("el camion ", camion.nombre, " lleva ", camion.cargo2.intentos, " intentos fallidos en la entrega de seguimiento: ", camion.cargo1.seguimiento, " \n")

				} else {
					fmt.Println("el camion ", camion.nombre, "entrego con exito el paquete con seguimiento: ", camion.cargo1.seguimiento, " \n")
					break
				}
				time.Sleep(time.Duration(camion.tEntre) * time.Second)
			}
			for camion.cargo2.intentos < 3 {
				i = rand.Intn(100)
				if i >= 80 {
					camion.cargo2.intentos = camion.cargo2.intentos + 1
					fmt.Println("el camion ", camion.nombre, " lleva ", camion.cargo2.intentos, " intentos fallidos en la entrega de seguimiento: ", camion.cargo2.seguimiento, " \n")

				} else {
					fmt.Println("el camion ", camion.nombre, "entrego con exito el paquete con seguimiento: ", camion.cargo2.seguimiento, " \n")
					break
				}
				time.Sleep(time.Duration(camion.tEntre) * time.Second)
			}
		} else {
			for camion.cargo2.intentos < 3 {
				i = rand.Intn(100)
				if i >= 80 {

					camion.cargo2.intentos = camion.cargo2.intentos + 1
					fmt.Println("el camion ", camion.nombre, " lleva ", camion.cargo2.intentos, " intentos fallidos en la entrega de seguimiento: ", camion.cargo2.seguimiento, " \n")

				} else {
					fmt.Println("el camion ", camion.nombre, "entrego con exito el paquete con seguimiento: ", camion.cargo2.seguimiento, " \n")
					break
				}
				time.Sleep(time.Duration(camion.tEntre) * time.Second)
			}
			for camion.cargo1.intentos < 3 {
				i = rand.Intn(100)
				if i >= 80 {
					camion.cargo2.intentos = camion.cargo2.intentos + 1
					fmt.Println("el camion ", camion.nombre, " lleva ", camion.cargo2.intentos, " intentos fallidos en la entrega de seguimiento: ", camion.cargo1.seguimiento, " \n")

				} else {
					fmt.Println("el camion ", camion.nombre, "entrego con exito el paquete con seguimiento: ", camion.cargo1.seguimiento, " \n")
					break
				}
			}
			time.Sleep(time.Duration(camion.tEntre) * time.Second)

		}

	} else if camion.cargo1.id != "NOHAY" { //si solo hay un paquete
		for camion.cargo1.intentos < 3 {
			i = rand.Intn(100)
			if i >= 80 {
				camion.cargo1.intentos = camion.cargo1.intentos + 1
				log.Println("el camion %s lleva %b intentos fallidos en la entrega de seguimiento: %s \n", camion.nombre, camion.cargo1.intentos, camion.cargo1.seguimiento)

			} else {
				log.Println("el camion %s entrego con exito el paquete con seguimiento: %s \n", camion.nombre, camion.cargo1.seguimiento)
				break
			}
			time.Sleep(time.Duration(camion.tEntre) * time.Second)
		}
	}

	return 0
}

func main() {

	CNormal := Camion{
		nombre: "normal",
		Tipo:   "normal",
	}
	/*CRetail1 := Camion{
		nombre: "retail1",
		Tipo:   "retail",
	}
	CRetail2 := Camion{
		nombre: "retail2",
		Tipo:   "retail",
	}*/
	var espera int
	var esperaEntrega int
	fmt.Println("Ingrese Tiempo que esperaran los camiones para recibir un segundo paquete")
	_, err2 := fmt.Scanf("%d\n", &espera)
	if err2 != nil {
		fmt.Println(err2)
	}
	CNormal.tEspera = espera
	//CRetail1.tEspera = espera
	//CRetail2.tEspera = espera
	/*
		fmt.Println("Tiempo de envio Paquete")
		tenvio, _ := reader.ReadString('\n')*/

	fmt.Println("Ingrese Tiempo de demora para la entrega de los camiones ")
	_, err1 := fmt.Scanf("%d\n", &esperaEntrega)
	if err1 != nil {
		fmt.Println(err1)
	}
	CNormal.tEntre = esperaEntrega
	go EnRuta(CNormal)
	//go Mandar(CRetail1)
	//go Mandar(CRetail2)
	reader := bufio.NewReader(os.Stdin)
	log.Printf("ingrese 1 para terminar el programa")
	for {
		text, _ := reader.ReadString('\n')
		text = strings.ToLower(strings.Trim(text, " \r\n"))

		if strings.Compare(text, "1") == 0 {
			break
		}

	}

}
