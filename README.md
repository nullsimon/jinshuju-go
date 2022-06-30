# Jinshuju-go

Jinshuju-go is a Go library for the Jinshuju API.

This is the official docs from [jinshuju](https://jinshuju.net/help/articles/api-intro)

## Supported API
* 获取单个表单的表头字段和字段选项
* 获取单个表单的所有的数据条目，并且根据表头组合成一个对象

## Usage

完整示例代码见examples/main.go

```go
    // 获取config
    var conf jinshuju.Conf
	raw, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(raw, &conf)
    
	client := jinshuju.NewClient(conf)
	// 获取单个表单的表头字段和字段选项
	form, err := client.GetFormFields("test")
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Info(len(form))
	// 获取单个表单的所有的数据条目，并且根据表头组合成一个对象
	entries, err := client.GetFormEntries("test")
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Info(len(entries))
```

## Contributing

目前仅支持部分API，欢迎来补全整个金数据的API

