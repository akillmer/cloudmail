package cloudmail

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"strings"
	"time"
)

type Message struct {
	Name      string `json:"name"`
	ReplyTo   string `json:"replyTo"`
	Message   string `json:"message"`
	Recaptcha string `json:"recaptcha"`
}

type RecaptchaResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

var (
	recaptchaSecret = os.Getenv("RECAPTCHA_SECRET")
	smtpUser        = os.Getenv("SMTP_USER")
	smtpPass        = os.Getenv("SMTP_PW")
	smtpAddr        = os.Getenv("SMTP_ADDR")
	smtpPort        = os.Getenv("SMTP_PORT")
	mailTo          = os.Getenv("MAIL_TO")
)

// SendMessage is the function exposed to the Cloud
func SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// axios sends the payload in the request's body as a JSON string
	var msg Message
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&msg); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	/* the message payload isn't verified here, it should be vetted on the client side.
	validating the client's Recaptcha response should help confirm that the data is still clean */

	if code, err := verifyRecaptcha(msg.Recaptcha); err != nil {
		respondWithError(w, code, err.Error())
		return
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpAddr)
	err := smtp.SendMail(smtpAddr+":"+smtpPort, auth, smtpUser, []string{mailTo}, msg.RFC822())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("Sent message from %s\n", msg.ReplyTo)
}

func respondWithError(w http.ResponseWriter, code int, errMsg string) {
	w.WriteHeader(code)
	w.Write([]byte(errMsg))
	log.Printf("Error %d: %v", code, errMsg)
}

func verifyRecaptcha(clientResponse string) (int, error) {
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{"secret": {recaptchaSecret}, "response": {clientResponse}})

	if err != nil {
		return resp.StatusCode, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return http.StatusNoContent, err
	}

	var recaptchaResponse RecaptchaResponse

	if err := json.Unmarshal(body, &recaptchaResponse); err != nil {
		return http.StatusInternalServerError, err
	}

	if !recaptchaResponse.Success {
		errorCodes := strings.Join(recaptchaResponse.ErrorCodes, ", ")
		return http.StatusForbidden, errors.New(errorCodes)
	}

	return http.StatusOK, nil
}

// RFC822 is an old and obsoleted format for emails, but is referred to in Go's SendMail function
func (m *Message) RFC822() []byte {
	var body = fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Contact Form\r\n\r\nSubmitted by: %s [%s]:\r\n\r\n%s",
		smtpUser, mailTo, m.Name, m.ReplyTo, m.Message)
	return []byte(body)
}
