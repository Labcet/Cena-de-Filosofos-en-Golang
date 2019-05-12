package main

import (
	"fmt"
	"math/rand"
	"time"
)

type filosofo struct {
	nombre      string
	cubierto chan bool
	vecino  *filosofo
}

func hacerfilosofo(nombre string, vecino *filosofo) *filosofo {
	filo := &filosofo{nombre, make(chan bool, 1), vecino}
	filo.cubierto <- true
	return filo
}

func (filo *filosofo) pensar() {
	fmt.Printf("%v esta pensando.\n", filo.nombre)
	time.Sleep(time.Duration(rand.Int63n(1e9)))
}

func (filo *filosofo) comer() {
	fmt.Printf("%v esta comiendo.\n", filo.nombre)
	time.Sleep(time.Duration(rand.Int63n(1e9)))
}

func (filo *filosofo) obtenerCubiertos() {
	tiempo_fuera := make(chan bool, 1)
	go func() { time.Sleep(1e9); tiempo_fuera <- true }()
	<-filo.cubierto
	fmt.Printf("%v tiene su cubierto.\n", filo.nombre)
	select {
	case <-filo.vecino.cubierto:
		fmt.Printf("%v coge el cubierto que le falta.\n", filo.nombre)
		fmt.Printf("%v tiene 2 cubiertos.\n", filo.nombre)
		return // REGRESA A COMER
	case <-tiempo_fuera:
		filo.cubierto <- true
		filo.pensar()
		filo.obtenerCubiertos()
	}
}

func (filo *filosofo) retornarCubiertos() {
	filo.cubierto <- true
	filo.vecino.cubierto <- true
}

func (filo *filosofo) cenar(anuncia chan *filosofo) {
	filo.pensar()
	filo.obtenerCubiertos()
	filo.comer()
	filo.retornarCubiertos()
	anuncia <- filo
}

func main() {
	nombres := []string{"Aristoteles", "Socrates", "Platon", "Anaxagoras", "Arquimedes"}
	filosofos := make([]*filosofo, len(nombres))
	var filo *filosofo
	for i, nombre := range nombres {
		filo = hacerfilosofo(nombre, filo)
		filosofos[i] = filo
	}
	filosofos[0].vecino = filo
	fmt.Printf("Hay %v filsofos sentados en la mesa.\n", len(filosofos))
	fmt.Println("Cada uno tiene un cubierto, y debe prestarse otro de su vecino para comer.\n")
	anuncia := make(chan *filosofo)

	for a:=0; a < 3; a++{
		for _, filo := range filosofos {
			go filo.cenar(anuncia)
		}

		for i := 0; i < len(nombres); i++ {
			filo := <-anuncia
			fmt.Printf("%v esta terminando de comer.\n", filo.nombre)
		}
	}
}