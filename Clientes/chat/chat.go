package chat

import (
	context "context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
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

//OrdenarPyme is para recibir las ordenes de las pymes y llamar la funcion que las guarda
func (s *Server) OrdenarPyme(ctx context.Context, message *Orden) (*Message, error) {
	code := " "
	if strings.Compare(message.Prioritario, "0") == 0 {
		code = GuardarOrden(message.Id, message.Producto, message.Valor, message.Tienda, message.Destino, "Normal")
	} else {
		code = GuardarOrden(message.Id, message.Producto, message.Valor, message.Tienda, message.Destino, "Prioritario")
	}
	return &Message{Body: code}, nil
}

//OrdenarRetail is para recibir las ordenes de retail y llamar la funcion que las guarda
func (s *Server) OrdenarRetail(ctx context.Context, message *Orden) (*Message, error) {
	code := GuardarOrden(message.Id, message.Producto, message.Valor, message.Tienda, message.Destino, "Retail")
	return &Message{Body: code}, nil
}

//GuardarOrden is para guardar las ordenes en un array y en memoria
func GuardarOrden(id string, producto string, valor string, tienda string, destino string, tipo string) string {
	//Registro en memoria
	csvfile, err := os.OpenFile("logistica.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	csvwriter := csv.NewWriter(csvfile)

	//Hora Actual
	t := time.Now()
	timestamp := fmt.Sprintf("%02d-%02d-%d %02d:%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute())

	//Codigo Seguimiento
	code := id + "177013"

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
				Id:          prioritario[1].id,
				Seguimiento: prioritario[1].seguimiento,
				Tipo:        prioritario[1].tipo,
				Valor:       prioritario[1].valor,
				Intentos:    0,
				Estado:      "En Camino",
			}

		} else if len(noprioritario) > 0 {
			pac = MPaquete{
				Id:          noprioritario[1].id,
				Seguimiento: noprioritario[1].seguimiento,
				Tipo:        noprioritario[1].tipo,
				Valor:       noprioritario[1].valor,
				Intentos:    0,
				Estado:      "En Camino",
			}
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
	pac = MPaquete{
		Id:          "NOHAY",
		Seguimiento: "NOHAY",
		Tipo:        "NOHAY",
		Valor:       "177013 pero no entro a la otra wea si",
		Intentos:    0,
		Estado:      "NOHAY",
	}
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
