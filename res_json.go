package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func generateErrorJson(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with %v error: %s\n", code, msg)
	}

	type errorJson struct {
		Error string `json:"error"`
	}

	res := errorJson{
		Error: msg,
	}

	respondWithJson(w, code, res)
}

func respondWithJson(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	rawData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError) //--> 500
		return
	}
	w.WriteHeader(code) //--> send the httpStatus code
	w.Write(rawData)
}

// payload any também poderia ser : payload interface{}
// Ao definir payload any (ou interface{}), a função pode ser usada com qualquer struct, mapa ou até mesmo um slice. Isso permite maior reutilização do código.
// Se usassemos payload responseValues, então payload apenas podia ser de type responseValues

// 3. Exemplo de Flexibilidade
// Com interface{}, posso chamar a função com diferentes tipos de dados:

// respondWithJSON(w, http.StatusOK, responseValues{Valid: true})
// respondWithJSON(w, http.StatusOK, map[string]string{"message": "ok"})
// respondWithJSON(w, http.StatusOK, struct {
//     Status  string `json:"status"`
//     Message string `json:"message"`
// }{"success", "Operação concluída"})
// Tudo isso funcionará porque json.Marshal(payload) consegue converter diferentes tipos para JSON.
