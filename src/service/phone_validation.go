package service

import (
	"encoding/json"
	"fmt"
	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type PhoneNumberValidateService struct{}

type ApiResponseData struct {
	Valid               bool   `json:"valid"`
	Number              string `json:"number"`
	LocalFormat         string `json:"local_format"`
	InternationalFormat string `json:"international_format"`
	CountryPrefix       string `json:"country_prefix"`
	CountryCode         string `json:"country_code"`
	CountryName         string `json:"country_name"`
	Location            string `json:"location"`
	Carrier             string `json:"carrier"`
	LineType            string `json:"line_type"`
}

func (p *PhoneNumberValidateService) Validate(phoneNumber string, countryCode string) (valid bool, err error) {
	checkValidPhoneAPI := fmt.Sprintf("http://apilayer.net/api/validate?access_key=%s&number=%s&country_code=%s&format=1", config.Params.NumValidAccessKey, phoneNumber, countryCode)
	statusCode, body, errs := fiber.
		Post(checkValidPhoneAPI).String()
	if statusCode != fiber.StatusOK {
		log.Error(errs)
		valid = false
		err = errors.New("Something is wrong while phone number validation check")
		return
	}
	var resData ApiResponseData
	_ = json.Unmarshal([]byte(body), &resData)
	valid = resData.Valid
	return
}
