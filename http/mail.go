package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jingyan/justcoin/models"
	"github.com/open-falcon/mail-provider/config"
	"github.com/toolkits/web/param"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func configProcRoutes() {

	http.HandleFunc("/sender/mail", func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Config()
		token := param.String(r, "token", "")
		if cfg.Http.Token != token {
			http.Error(w, "no privilege", http.StatusForbidden)
			return
		}

		tos := param.MustString(r, "tos")
		subject := param.MustString(r, "subject")
		content := param.MustString(r, "content")
		tos = strings.Replace(tos, ",", ";", -1)

		fmt.Println("subject", subject, "content", content)
		errCode := SendMail(tos, subject, content)
		if errCode != models.No_Error {
			fmt.Println("send email err")
			return
		}

		/*s := smtp.New(cfg.Smtp.Addr, cfg.Smtp.Username, cfg.Smtp.Password)
		err := s.SendMail(cfg.Smtp.From, tos, subject, content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "success", http.StatusOK)
		}
		*/
	})

}

func SendMail(tos string, subject, content string) int32 {
	cfg := config.Config()
	apiUser := cfg.Sendcloud.ApiUser //"justcoin"
	apiKey := cfg.Sendcloud.ApiKey   //"pdSiJgtO4hCNnHqN"
	from := cfg.Sendcloud.From

	requestURI := "http://api.sendcloud.net/apiv2/mail/send"

	/*xsmtpapi := struct {
		To  []string            `json:"to"`
		Sub map[string][]string `json:"sub"`
	}{
		To:  []string{tos},
		Sub: map[string][]string{},
	}
	xsmtpapi_str, _ := json.Marshal(&xsmtpapi)
	*/
	postParams := url.Values{
		"apiUser": {apiUser},
		"apiKey":  {apiKey},
		"from":    {from},
		//"xsmtpapi": {string(xsmtpapi_str)},
		"to":             {tos},
		"subject":        {subject},
		"html":           {content},
		"useAddressList": {"true"},
	}
	postBody := bytes.NewBufferString(postParams.Encode())
	fmt.Println("req:", models.Json2String(&postParams))
	responseHandler, err := http.Post(requestURI, "application/x-www-form-urlencoded", postBody)
	if err != nil {
		fmt.Println("http.Post err", err.Error())
		return models.Err_Auth_Send_Mail_Err
	}
	defer responseHandler.Body.Close()
	BodyByte, err := ioutil.ReadAll(responseHandler.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll err", err.Error())
		return models.Err_Auth_Send_Mail_Err
	}
	fmt.Println("resp", string(BodyByte))
	resp := &MessageResponse{}
	json.Unmarshal(BodyByte, resp)
	if resp.StatusCode != 200 {
		return models.Err_Auth_Send_Mail_Err
	}
	return models.No_Error
}

type MessageResponse struct {
	Result     bool   `json:"result"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}
