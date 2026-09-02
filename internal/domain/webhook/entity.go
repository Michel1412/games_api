package webhook

type Webhook struct {
	ID    int    `json:"id" firestore:"id"`
	URL   string `json:"url" firestore:"url"`
	Ativo bool   `json:"ativo" firestore:"ativo"`
}
