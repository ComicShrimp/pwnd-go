package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/term"
)

func main() {
	fmt.Printf("Digite a senha que deseja consultar: ")

	userPassword, error := term.ReadPassword(int(os.Stdin.Fd()))
	if error != nil {
		log.Fatalln("Erro, é necessário passar uma senha: ", error)
	}

	h := sha1.New()
	h.Write(userPassword)

	hash := h.Sum(nil)

	hashedPassword := strings.ToUpper(hex.EncodeToString(hash))

	prefix := hashedPassword[:5]
	suffix := hashedPassword[5:]

	baseURL := "https://api.pwnedpasswords.com/range"
	fullURL := fmt.Sprintf("%s/%s", baseURL, prefix)

	resp, err := http.Get(fullURL)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, suffix) {
			count := strings.Split(line, ":")[1]
			fmt.Printf("\nOps, sua senha foi vazada %s vezes\n", count)
			return
		}
	}

	fmt.Println("Password is safe!")
}
