package utils

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func GetToken(url, username, password string) (*string, error) {
	url = url + "/oauth/authorize"
	newReq, _ := http.NewRequest("GET", url, nil)
	userAndPasswd := username + ":" + password
	basicAuth := base64.StdEncoding.EncodeToString([]byte(userAndPasswd))

	newReq.Header.Set("Authorization", "Basic "+basicAuth)
	newReq.Header.Set("X-CSRF-Token", "1")
	newReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	q := newReq.URL.Query()
	q.Add("client_id", "openshift-challenging-client")
	q.Add("response_type", "token")
	newReq.URL.RawQuery = q.Encode()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {

		return errors.New("redirect in implicit grant flow, and return 302")
	}
	res, err := client.Do(newReq)

	if err != nil && res.StatusCode == 302 { //status code 302
		u, err := res.Location()
		if err != nil {
			return nil, err
		}
		// https://docs.openshift.org/latest/architecture/additional_concepts/authentication.html
		s := u.String()
		l := strings.Index(s, "access_token=")
		// no access_token in return.
		if l == -1 {
			return nil, errors.New("no access_token return")
		}
		token := s[l+13 : l+56]
		return &token, nil
	} else {
		return nil, fmt.Errorf("%v %v : %v", newReq.Method, newReq.URL.String(), err)
	}
	return nil, fmt.Errorf("没有获取到token")
}

// token通过base64编码
func TokenToBase64Encode(token string) string {
	encodeStr := base64.URLEncoding.EncodeToString([]byte(token))
	//去掉最后生成的2个==
	eocodeSubStr := encodeStr[0 : len(encodeStr)-2]
	return eocodeSubStr
}
