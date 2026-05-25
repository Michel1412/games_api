package jogo

type Jogo struct {
	ID     int    `json:"id" firestore:"id"`
	Nome   string `json:"nome" firestore:"nome"`
	Tipo   string `json:"tipo" firestore:"tipo"`
	Nota   int    `json:"nota" firestore:"nota"`
	Review string `json:"review" firestore:"review"`
}
