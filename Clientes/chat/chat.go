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
var hicelacosa int = 1

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

	//fmt.Println(prioritario)

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
	if hicelacosa == 1 {
		prioritario = remove(prioritario, 0)
		hicelacosa = 0
	}
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

		} else if len(noprioritario) > 1 {
			pac = MPaquete{
				Id:          noprioritario[1].id,
				Seguimiento: noprioritario[1].seguimiento,
				Tipo:        noprioritario[1].tipo,
				Valor:       noprioritario[1].valor,
				Intentos:    0,
				Estado:      "En Camino",
			}
			noprioritario = remove(noprioritario, 1)
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
				Id:          retail[1].id,
				Seguimiento: retail[1].seguimiento,
				Tipo:        retail[1].tipo,
				Valor:       retail[1].valor,
				Intentos:    0,
				Estado:      "En Camino",
			}
			retail = remove(retail, 1)
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
