# youdaoweb

get https://dict.youdao.com result from your terminal



# install

```bash
go get github.com/yorks/youdaoweb
```


# Usage
```bash
cd cmd && go run main.go $word
```
or

```golang
doc, _ := youdaoweb.Qword(os.Args[1])
if doc.ParseJsCode() == nil {
	fmt.Println(doc.GetPhonetic())
	//fmt.Println(doc.GetWebs())
	//fmt.Println(doc.GetSpecs())
	//fmt.Println(doc.GetEEs())
	//fmt.Println(doc.GetExpandECs())
	//fmt.Println(doc.GetBlngSents())
	//fmt.Println(doc.GetMediaSents())
	//fmt.Println(doc.GetAuthSents())
	//fmt.Println(doc.GetPhRs())
	//fmt.Println(doc.GetSynos())
	//fmt.Println(doc.GetRelWords())
	//fmt.Println(doc.GetETYMs())
	//fmt.Println(doc.GetWIKIs())

}

