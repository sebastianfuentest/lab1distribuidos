package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"papa.com/Clientes/chat"
)

//Camion is la estructura de los camiones
type Camion struct {
	nombre     string
	conpaquete int
	Tipo       string
	tEspera    int
	tEntre     int
	cargo1     Paquete
	cargo2     Paquete
}

//Paquete is la estructura de los paquetes
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

	for {
		if camion.conpaquete == 0 {
			camion.conpaquete = 1
			log.Printf("Mandando camion %s \n", camion.nombre)
			camion.conpaquete = Mandar(camion)

		}

	}
}

//Mandar is la funcion encargada de enviar los paquetes
func Mandar(camion Camion) (ret int) {
	var i int
	camion.cargo1.id = "nada"
	camion.cargo2.id = "nada"
	var conn *grpc.ClientConn
	mensaje := chat.Message{
		Body: camion.Tipo,
	}
	conn, err := grpc.Dial("dist49:9000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	c := chat.NewChatServiceClient(conn)
	response, err := c.RecibirPaquete(context.Background(), &mensaje)
	if err != nil {
		log.Fatalf("We couldn't say hello: %s", err)
	} //espero el tiempo antes de pedir el segundo pedido
	time.Sleep(time.Duration(camion.tEspera) * time.Second)
	response2, err := c.RecibirPaquete(context.Background(), &mensaje)
	if err != nil {
		log.Fatalf("We couldn't say hello: %s", err)
	}

	//Registro en Memoria
	csvfile, err := os.OpenFile(camion.nombre+".csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	csvwriter := csv.NewWriter(csvfile)
	//

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
					fmt.Println("el camion ", camion.nombre, "entrego con exito luego de ", camion.cargo1.intentos, " intentos el paquete con seguimiento: ", camion.cargo1.seguimiento, " \n")
					message := chat.NuevoEstado{
						Seguimiento: camion.cargo1.seguimiento,
						Nuevoestado: "Entrega Exitosa",
					}
					c.CambiarEstado(context.Background(), &message)

					//Escribir en el registro
					t := time.Now()
					timestamp := fmt.Sprintf("%02d-%02d-%d %02d:%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute())
					intenstring := strconv.Itoa(int(camion.cargo1.intentos) + 1)
					registro := []string{camion.cargo1.id, camion.cargo1.tipo, camion.cargo1.valor, intenstring, timestamp}
					csvwriter.Write(registro)
					break
				}
				time.Sleep(time.Duration(camion.tEntre) * time.Second)
			}
			if camion.cargo1.intentos == 3 {
				//fallo entrega cargo1
				message := chat.NuevoEstado{
					Seguimiento: camion.cargo1.seguimiento,
					Nuevoestado: "Entrega Fallida",
				}
				c.CambiarEstado(context.Background(), &message)

				//Escribir en el registro
				intenstring := strconv.Itoa(int(camion.cargo1.intentos))
				registro := []string{camion.cargo1.id, camion.cargo1.tipo, camion.cargo1.valor, intenstring, "0"}
				csvwriter.Write(registro)
				//---------------------------------------------------

			}
			for camion.cargo2.intentos < 3 {
				i = rand.Intn(100)
				if i >= 80 {
					camion.cargo2.intentos = camion.cargo2.intentos + 1
					fmt.Println("el camion ", camion.nombre, " lleva ", camion.cargo2.intentos, " intentos fallidos en la entrega de seguimiento: ", camion.cargo2.seguimiento, " \n")

				} else {
					fmt.Println("el camion ", camion.nombre, "entrego con exito luego de ", camion.cargo2.intentos, " intentos el paquete con seguimiento: ", camion.cargo2.seguimiento, " \n")
					message2 := chat.NuevoEstado{
						Seguimiento: camion.cargo2.seguimiento,
						Nuevoestado: "Entrega Exitosa",
					}
					c.CambiarEstado(context.Background(), &message2)
					//Escribir en el registro
					t := time.Now()
					timestamp := fmt.Sprintf("%02d-%02d-%d %02d:%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute())
					intenstring := strconv.Itoa(int(camion.cargo2.intentos) + 1)
					registro := []string{camion.cargo2.id, camion.cargo2.tipo, camion.cargo2.valor, intenstring, timestamp}
					csvwriter.Write(registro)
					break

				}
				time.Sleep(time.Duration(camion.tEntre) * time.Second)
			}
			if camion.cargo2.intentos == 3 {
				//fallo entrega cargo 2
				message2 := chat.NuevoEstado{
					Seguimiento: camion.cargo2.seguimiento,
					Nuevoestado: "Entrega Fallida",
				}
				c.CambiarEstado(context.Background(), &message2)
				//Escribir en el registro
				intenstring := strconv.Itoa(int(camion.cargo2.intentos))
				registro := []string{camion.cargo2.id, camion.cargo2.tipo, camion.cargo2.valor, intenstring, "0"}
				csvwriter.Write(registro)
			}
		} else {
			for camion.cargo2.intentos < 3 {
				i = rand.Intn(100)
				if i >= 80 {

					camion.cargo2.intentos = camion.cargo2.intentos + 1
					fmt.Println("el camion ", camion.nombre, " lleva ", camion.cargo2.intentos, " intentos fallidos en la entrega de seguimiento: ", camion.cargo2.seguimiento, " \n")

				} else {
					fmt.Println("el camion ", camion.nombre, "entrego con exito luego de ", camion.cargo2.intentos, " intentos el paquete con seguimiento: ", camion.cargo2.seguimiento, " \n")
					message2 := chat.NuevoEstado{
						Seguimiento: camion.cargo2.seguimiento,
						Nuevoestado: "Entrega Exitosa",
					}
					c.CambiarEstado(context.Background(), &message2)
					//Escribir en el registro
					t := time.Now()
					timestamp := fmt.Sprintf("%02d-%02d-%d %02d:%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute())
					intenstring := strconv.Itoa(int(camion.cargo2.intentos) + 1)
					registro := []string{camion.cargo2.id, camion.cargo2.tipo, camion.cargo2.valor, intenstring, timestamp}
					csvwriter.Write(registro)
					break
				}
				time.Sleep(time.Duration(camion.tEntre) * time.Second)
			}
			if camion.cargo2.intentos == 3 {
				//fallo entrega cargo 2
				message2 := chat.NuevoEstado{
					Seguimiento: camion.cargo2.seguimiento,
					Nuevoestado: "Entrega Fallida",
				}
				c.CambiarEstado(context.Background(), &message2)
				//Escribir en el registro
				intenstring := strconv.Itoa(int(camion.cargo2.intentos))
				registro := []string{camion.cargo2.id, camion.cargo2.tipo, camion.cargo2.valor, intenstring, "0"}
				csvwriter.Write(registro)
				//---------------
			}
			for camion.cargo1.intentos < 3 {
				i = rand.Intn(100)
				if i >= 80 {
					camion.cargo2.intentos = camion.cargo2.intentos + 1
					fmt.Println("el camion ", camion.nombre, " lleva ", camion.cargo1.intentos, " intentos fallidos en la entrega de seguimiento: ", camion.cargo1.seguimiento, " \n")

				} else {
					fmt.Println("el camion ", camion.nombre, "entrego con exito luego de ", camion.cargo1.intentos, " intentos el paquete con seguimiento: ", camion.cargo1.seguimiento, " \n")
					message := chat.NuevoEstado{
						Seguimiento: camion.cargo1.seguimiento,
						Nuevoestado: "Entrega Exitosa",
					}
					c.CambiarEstado(context.Background(), &message)
					//Escribir en el registro
					t := time.Now()
					timestamp := fmt.Sprintf("%02d-%02d-%d %02d:%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute())
					intenstring := strconv.Itoa(int(camion.cargo1.intentos) + 1)
					registro := []string{camion.cargo1.id, camion.cargo1.tipo, camion.cargo1.valor, intenstring, timestamp}
					csvwriter.Write(registro)
					break
				}
			}
			if camion.cargo1.intentos == 3 {
				//fallo entrega cargo1
				message := chat.NuevoEstado{
					Seguimiento: camion.cargo1.seguimiento,
					Nuevoestado: "Entrega Fallida",
				}
				c.CambiarEstado(context.Background(), &message)
				//Escribir en el registro
				intenstring := strconv.Itoa(int(camion.cargo1.intentos))
				registro := []string{camion.cargo1.id, camion.cargo1.tipo, camion.cargo1.valor, intenstring, "0"}
				csvwriter.Write(registro)
				//------------------
				time.Sleep(time.Duration(camion.tEntre) * time.Second)
			}

		}

	} else if camion.cargo1.id != "NOHAY" { //si solo hay un paquete
		for camion.cargo1.intentos < 3 {
			i = rand.Intn(100)
			if i >= 80 {
				camion.cargo1.intentos = camion.cargo1.intentos + 1
				log.Println("el camion %s lleva %b intentos fallidos en la entrega de seguimiento: %s \n", camion.nombre, camion.cargo1.intentos, camion.cargo1.seguimiento)

			} else {
				log.Println("el camion %s entrego con exito luego de %b intentos el paquete con seguimiento: %s \n", camion.nombre, camion.cargo1.valor, camion.cargo1.seguimiento)
				message := chat.NuevoEstado{
					Seguimiento: camion.cargo1.seguimiento,
					Nuevoestado: "Entrega Exitosa",
				}
				c.CambiarEstado(context.Background(), &message)
				//Escribir en el registro
				t := time.Now()
				timestamp := fmt.Sprintf("%02d-%02d-%d %02d:%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute())
				intenstring := strconv.Itoa(int(camion.cargo1.intentos) + 1)
				registro := []string{camion.cargo1.id, camion.cargo1.tipo, camion.cargo1.valor, intenstring, timestamp}
				csvwriter.Write(registro)
				break
			}
			time.Sleep(time.Duration(camion.tEntre) * time.Second)
		}
		if camion.cargo1.intentos == 3 {
			message := chat.NuevoEstado{
				Seguimiento: camion.cargo1.seguimiento,
				Nuevoestado: "Entrega Fallida",
			}
			c.CambiarEstado(context.Background(), &message)
			//Escribir en el registro
			intenstring := strconv.Itoa(int(camion.cargo1.intentos))
			registro := []string{camion.cargo1.id, camion.cargo1.tipo, camion.cargo1.valor, intenstring, "0"}
			csvwriter.Write(registro)
		}
	}

	//Cerrar los archivos
	csvwriter.Flush()
	csvfile.Close()
	//-------

	return 0
}

func main() {

	CNormal := Camion{
		nombre: "normal",
		Tipo:   "normal",
	}
	CRetail1 := Camion{
		nombre: "retail1",
		Tipo:   "retail",
	}
	CRetail2 := Camion{
		nombre: "retail2",
		Tipo:   "retail",
	}
	var espera int
	var esperaEntrega int
	fmt.Println("Ingrese Tiempo que esperaran los camiones para recibir un segundo paquete")
	_, err2 := fmt.Scanf("%d\n", &espera)
	if err2 != nil {
		fmt.Println(err2)
	}
	CNormal.tEspera = espera
	CRetail1.tEspera = espera
	CRetail2.tEspera = espera

	fmt.Println("Ingrese Tiempo de demora para la entrega de los camiones ")
	_, err1 := fmt.Scanf("%d\n", &esperaEntrega)
	if err1 != nil {
		fmt.Println(err1)
	}
	CNormal.tEntre = esperaEntrega
	go EnRuta(CNormal)
	go EnRuta(CRetail1)
	go EnRuta(CRetail2)
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
