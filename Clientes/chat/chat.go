package chat

import (
	context "context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

//Server is
type Server struct {
	mute sync.Mutex
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

var retail []Paquete
var prioritario []Paquete
var noprioritario []Paquete
var listapaquetes []Paquete

//remove is
func remove(slice []Paquete, p int) []Paquete {
	return append(slice[:p], slice[p+1:]...)
}

//OrdenarPyme is
func (s *Server) OrdenarPyme(ctx context.Context, message *Orden) (*Message, error) {
	code := " "
	if strings.Compare(message.Prioritario, "0") == 0 {
		code = GuardarOrden(message.Id, message.Producto, message.Valor, message.Tienda, message.Destino, "Normal")
	} else {
		code = GuardarOrden(message.Id, message.Producto, message.Valor, message.Tienda, message.Destino, "Prioritario")
	}
	return &Message{Body: code}, nil
}

//OrdenarRetail is
func (s *Server) OrdenarRetail(ctx context.Context, message *Orden) (*Message, error) {
	code := GuardarOrden(message.Id, message.Producto, message.Valor, message.Tienda, message.Destino, "Retail")
	return &Message{Body: code}, nil
}

//GuardarOrden is
func GuardarOrden(id string, producto string, valor string, tienda string, destino string, tipo string) string {
	//Registro en Memoria
	csvfile, err := os.OpenFile("dblogistica.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	csvwriter := csv.NewWriter(csvfile)

	//Hora Actual
	t := time.Now()
	timestamp := fmt.Sprintf("%02d-%02d-%d %02d:%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute())

	//Codigo Seguimiento
	var banana int
	banana = len(listapaquetes)
	var strbanana string
	strbanana = strconv.Itoa(banana)
	code := id + strbanana //+"177013"

	//Registro en struct
	registro := []string{timestamp, id, tipo, producto, valor, tienda, destino, code}

	//Agregar Paquete a la cola

	appPaquete := Paquete{
		id:          id,
		seguimiento: code,
		tipo:        tipo,
		valor:       valor,
		intentos:    0,
		estado:      "En Cola",
	}

	if strings.Compare(tipo, "Normal") == 0 {
		noprioritario = append(noprioritario, appPaquete)
	} else if (strings.Compare(tipo, "Prioritario")) == 0 {
		prioritario = append(prioritario, appPaquete)
	} else {
		retail = append(retail, appPaquete)
	}

	//fmt.Println(prioritario)
	listapaquetes = append(listapaquetes, appPaquete)

	//Agrego al csv
	csvwriter.Write(registro)
	csvwriter.Flush()
	csvfile.Close()

	return code
}

//RecibirPaquete is Funcion para darle paquetes al camion
func (s *Server) RecibirPaquete(ctx context.Context, message *Message) (*MPaquete, error) {
	var pac MPaquete
	s.mute.Lock()
	if message.GetBody() == "normal" {
		if len(prioritario) > 0 {
			pac = MPaquete{
				Id:          prioritario[0].id,
				Seguimiento: prioritario[0].seguimiento,
				Tipo:        prioritario[0].tipo,
				Valor:       prioritario[0].valor,
				Intentos:    0,
				Estado:      "En Camino",
			}
			prioritario = remove(prioritario, 0)

		} else if len(noprioritario) > 0 {
			pac = MPaquete{
				Id:          noprioritario[0].id,
				Seguimiento: noprioritario[0].seguimiento,
				Tipo:        noprioritario[0].tipo,
				Valor:       noprioritario[0].valor,
				Intentos:    0,
				Estado:      "En Camino",
			}
			noprioritario = remove(noprioritario, 0)
		} else {
			pac = MPaquete{
				Id:          "NOHAY",
				Seguimiento: "NOHAY",
				Tipo:        "NOHAY",
				Valor:       "177013 xd",
				Intentos:    0,
				Estado:      "NOHAY",
			}
		}
		s.mute.Unlock()
		return &pac, nil

	}
	if message.GetBody() == "retail" {
		if len(retail) > 0 {
			pac = MPaquete{
				Id:          retail[0].id,
				Seguimiento: retail[0].seguimiento,
				Tipo:        retail[0].tipo,
				Valor:       retail[0].valor,
				Intentos:    0,
				Estado:      "En Camino",
			}
			retail = remove(retail, 0)
		} else {
			pac = MPaquete{
				Id:          "NOHAY",
				Seguimiento: "NOHAY",
				Tipo:        "NOHAY",
				Valor:       "177013",
				Intentos:    0,
				Estado:      "NOHAY",
			}
		}

		s.mute.Unlock()
		return &pac, nil
	}
	pac = MPaquete{
		Id:          "NOHAYNADADENADA",
		Seguimiento: "NOHAY",
		Tipo:        "NOHAY",
		Valor:       "177013",
		Intentos:    0,
		Estado:      "NOHAY",
	}

	s.mute.Unlock()
	return &pac, nil
}

//SeguimientoPaquete is para recibir el estado actual del paquete
func (s *Server) SeguimientoPaquete(ctx context.Context, message *Message) (*Message, error) {

	//Recorro la lista de paquetes
	for i := len(listapaquetes) - 1; i >= 0; i-- {
		//Obtengo el paquete en la posicion i y comparo su codigo de seguimiento con el mensaje
		log.Print(message.Body)
		if strings.Compare(message.Body, listapaquetes[i].seguimiento) == 0 {
			return &Message{Body: listapaquetes[i].estado}, nil
		}
	}
	return &Message{Body: "No se ha encontrado el paquete"}, nil
}

//CambiarEstado is Para actualizar el estado de los paquetes en la ListaPaquetes y en el csv
func (s *Server) CambiarEstado(ctx context.Context, message *NuevoEstado) (*Message, error) {

	//Recorro la lista de paquetes
	for i := len(listapaquetes) - 1; i >= 0; i-- {
		//Obtengo el paquete en la posicion i y comparo su codigo de seguimiento con el mensaje
		if strings.Compare(message.Seguimiento, listapaquetes[i].seguimiento) == 0 {
			listapaquetes[i].estado = message.Nuevoestado
			return &Message{Body: "Cambio realizado con exito"}, nil
		}
	}

	return &Message{Body: "No se ha encontrado el paquete"}, nil

}
