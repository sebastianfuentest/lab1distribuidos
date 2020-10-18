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

	fmt.Println(prioritario)

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
