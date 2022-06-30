package jinshuju

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Conf struct {
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
}

type Client struct {
	Conf
}

type FieldData struct {
	Name        string     `json:"name"`
	Token       string     `json:"token"`
	Description string     `json:"description"`
	Fields      []FieldMap `json:"fields"`
}

type EntriesData struct {
	Total int     `json:"total"`
	Count int     `json:"count"`
	Data  []Entry `json:"data"`
	Next  int     `json:"next"`
}

func NewClient(conf Conf) *Client {
	return &Client{
		Conf: conf,
	}
}

func (c *Client) GetFormFields(formName string) ([]FieldMap, error) {
	api := "https://jinshuju.net/api/v1/forms/" + formName
	req, err := http.NewRequest("GET", api, nil)
	req.SetBasicAuth(c.AppKey, c.AppSecret)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var form FieldData

	err = json.Unmarshal(body, &form)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return form.Fields, nil
}

func (c *Client) GetFormEntries(formName string) ([]Entry, error) {

	form, error := c.GetFormFields(formName)
	if error != nil {
		log.Error(error)
		return nil, error
	}

	var allEntries []Entry
	var next = 1

	for next != 0 {
		api := "https://jinshuju.net/api/v1/forms/" + formName + "/entries?next=" + strconv.Itoa(next)
		req, err := http.NewRequest("GET", api, nil)
		req.SetBasicAuth(c.AppKey, c.AppSecret)

		if err != nil {
			log.Error(err)
			return nil, err
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		var jinshujuData map[string]interface{}
		var entries EntriesData

		err = json.Unmarshal(body, &jinshujuData)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		err = json.Unmarshal(body, &entries)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		// 无data的时候，返回空
		if jinshujuData["data"] == nil {
			next = 0
			continue
		}
		var data = jinshujuData["data"].([]interface{})
		for i, entry := range entries.Data {
			if data[i] == nil {
				continue
			}
			jinshujuFieldMap := data[i].(map[string]interface{})
			// range form
			for _, fieldMap := range form {
				// get field key and value
				var fieldKey = ""
				var field Field
				for k, v := range fieldMap {
					fieldKey = k
					field = v
				}
				value, ok := jinshujuFieldMap[fieldKey]
				if !ok || value == nil {
					// 无此字段
					if field.isSectionBreak() {
						continue
					} else {
						log.Info(`entry 无此字段 ` + fieldKey + ` ` + field.Label)
					}
				}
				var choose []Choice
				if field.isInput() {
					choose = append(choose, Choice{Value: value.(string)})
				}
				if field.isSingleChoice() {
					// 单选
					for _, choice := range field.Choices {
						if choice.Name == value.(string) {
							choose = append(choose, choice)
						}
					}
					// 看下是否可输入,空值不考虑
					if len(choose) == 0 && value.(string) != "" {
						if field.isAllowOther() && strings.HasPrefix(value.(string), "其他") {
							choose = append(choose, Choice{OtherValue: value.(string), IsOther: true})
						} else {
							log.Info(`单选：无此选项 fieldKey: %v, fieldLabel: %v, value, %v`, fieldKey, field.Label, value.(string))
						}
					}
				}
				if field.isMultipleChoice() {
					// 多选
					choiceValueMap := make(map[string]bool)
					for _, choice := range field.Choices {
						choiceValueMap[choice.Name] = true
					}
					for _, v := range value.([]interface{}) {
						if _, ok := v.(string); ok {
							choose = append(choose, Choice{Value: v.(string)})
						} else {
							// 看下是否可输入
							if v.(string) != "" {
								if field.isAllowOther() && strings.HasPrefix(v.(string), "其他") {
									choose = append(choose, Choice{OtherValue: v.(string), IsOther: true})
								} else {
									log.Info(`多选：无此选项 fieldKey: %v, fieldLabel: %v, value, %v`, fieldKey, field.Label, v.(string))
								}
							}
						}
					}
				}

				field.ChooseChoices = choose
				fieldMap[fieldKey] = field
				entry.Fields = append(entry.Fields, fieldMap)
			}
			allEntries = append(allEntries, entry)
		}
		next = entries.Next
	}
	return allEntries, nil

}
