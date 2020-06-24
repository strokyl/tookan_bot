package googlecloudfunctions

import (
	"fmt"
	"net/http"
	"github.com/strokyl/tookan_bot/pkg/generate_hash"
)

func CancelTask(w http.ResponseWriter, r *http.Request) {
	jobId, err := getJobId(w, r)
	if err != nil {
		return
	}
	hashSalt := generate_hash.Secret(secret)
	if !checkSecret(w, r, string(generate_hash.GenerateHash(hashSalt, string(jobId)))) {
		return
	}
	err = client.MarkTaskFailed(jobId, "Tâche annulé par le cuisto")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "Ton Panier à bien été annulé. On est au regret de ne pas avoir pu répandre ta générosité dans nos rues #PourEux. On espère que tu redéposeras très vite un Panier et que cette fois on honorera celà comme il se doit !")
}
