package youdaoweb
import (
	"fmt"
	"net/http"
	"time"
	"strings"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"errors"
)

type WordResult struct {
	Html string
	Code int 
	Msg  string
}
type MyDoc struct {
	*goquery.Document
	JSCode string
	WordData map[string]interface{}
}
// 简明 
type EC struct {
	ExamTypes	[]string `json:"exam_type"`
	Trans		[]string `json:"trans"`
	UKPhone		string	 `json:"ukphone"`
	USPhone		string	 `json:"usphone"`
}
// 柯林斯
type Collins struct {
	Trans		[]CoTran `json:"collins_entries"`
}
type CoTran struct {
	Pos 	string	`json:"pos"`
	Tran 	string	`json:"tran"`
	ExamSents []ExamSent `json:"sent"`
}
type ExamSent struct {
	Chn   string `json:"chn_sent"`
	Eng   string `json:"eng_sent"`
}
// 网络释义
type Webs struct {
	Trans []Web `json:"web_trans"`
}
type Web struct {
	Key	string   `json:"key"`
	Trans []WebTrans `json:"trans"`
}
type WebTrans struct {
	Cls	[]string  `json:"cl"`
	Summarys []string `json:"line"`
	Value	string    `json:"value"`
}
// 专业释义
type Specs struct {
	Trans []Spec	`json:"entries"`	
}
type Spec struct {
	Major string `json:"major"`
	Trans []SpecTran `json:"trs"`
}
type SpecTran struct {
	Chn	string `json:"chnSent"`	
	Eng	string `json:"engSent"`	
	Tran	string `json:"nat"`	
}
//英英
type EEs struct {
	Trans []EE `json:"ee.word"`
}
type EE struct {
	Pos string `json:"pos"`
	Trans []EETran `json:"tr"`
}
type EETran struct {
	Tran string `json:"tran"`
	Swords []string `json:"similar-words"`
}

// 双语例句
type BlngSents struct {
	Trans []BlngSent `json:"双语例句"`
}
type BlngSent struct {
	Sentence  string `json:"sentence"`
	Eng       string `json:"sentence-eng"`
	Chn       string `json:"sentence-translation"`
	Speech    string `json:"sentence-speech"`
	Source    string `json:"source"`
}

// 原声例句
type MediaSents struct {
	Trans []MediaSent  `json:"原声例句"`
}
type MediaSent struct {
	Eng	string `json:"eng"`
	Name	string `json:"name"`
	Source	string `json:"source"`
	StreamUrl	string `json:"streamUrl"`
}
// 权威例句
type AuthSents struct {
	Trans []AuthSent `json:"权威例句"`
}
type AuthSent struct {
	Foreign	string `json:"foreign"`
	Source	string `json:"source"`
	Speech	string `json:"speech"`
	// score 
}
// 词典短语
type PhRs struct {
	Trans []PhR  `json:"词典短语"`
}
type PhR struct {
	Word string `json:"headword"`
	Tran string `json:"translation"`
}
// 同近义词
type Synos struct {
	Trans []Syno `json:"同近义词"`
}
type Syno struct {
	Pos	string	`json:"pos"`
	Tran	string	`json:"tran"`
	WS	[]string `json:"ws"`
}
// 同根词
type RelWords struct {
	Trans []RelWord `json:"同根词"`
}
type RelWord struct {
	Pos	string	`json:"pos"`
	Tran	string	`json:"tran"`
	Word	string 	`json:"word"`
}

// 词源
type ETYMs struct {
	Trans []ETYM  `json:"词源"`
}

type ETYM struct {
	Desc string `json:"desc"`
	Value string `json:"value"`
	Word string `json:"word"`
}
type ETYMRoot struct {
	ETYM map[string][]ETYM `json:"etyms"`
}
// 百科
type WIKIs struct {
	Trans []Wiki `json:"百科"`
}
type Wiki struct {
	Name string `json:"name"`
	Url string `json:"url"`
	Key string `json:"key"`
	Summary string `json:"summary"`
}

//
type ExpandECs struct {
	Trans []ExpandEC `json:"expan_ec.word.transList"`
}
type ExpandEC struct {
	SentOrig  string `json:"sentOrig"`
	SentTrans string `json:"sentTrans"`
	SentSpeech string `json:"sentSpeech"`
	Source     string `json:"source"`
}



func Qword(word string) (*MyDoc, error){
	url := fmt.Sprintf("https://dict.youdao.com/result?word=%s&lang=en", word)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Host", "dict.youdao.com")
	req.Header.Set("User-Agent", GetRandomUA())

	cli := &http.Client{Timeout: 5 * time.Second}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	mydoc := &MyDoc{doc, "",  data}
	return mydoc, nil
}

func (doc *MyDoc) GetPhonetic() map[string]string {
/*
<div data-v-39fab836="" class="phone_con">
  <div data-v-39fab836="" class="per-phone">
    <span data-v-39fab836="">英</span>
    <span data-v-39fab836="" class="phonetic">/ ˈheləʊ /</span>
    <div class="phraseSpeech phonetic-speech">
      <a title="点击发音" href="javascript:;" class="pronounce"></a>
    </div>
  </div>
  <div data-v-39fab836="" class="per-phone">
    <span data-v-39fab836="">美</span>
    <span data-v-39fab836="" class="phonetic">/ ˈheloʊ; ˈhiːloʊ /</span>
    <div class="phraseSpeech phonetic-speech">
      <a title="点击发音" href="javascript:;" class="pronounce"></a>
    </div>
  </div>
</div>
*/
	
	result := make(map[string]string)
	doc.Find(".per-phone").Each(func(i int, s *goquery.Selection) {
		label := s.Find("span").First().Text()
		phonetic := s.Find(".phonetic").Text()
		//fmt.Println(label, phonetic)
		result[label] = phonetic
	})
	return result
}
func (doc *MyDoc) GetDictsHtml() []string {
	var result  []string
	doc.Find(".dict-tabs .tabs-con .tab-item").Each(func(i int, s *goquery.Selection){
		fmt.Println(i)
		d := s.Text()
		fmt.Println(d)
		result = append(result, d)
	})
	return result
}
func (doc *MyDoc) GetTrans() []string {
	var result  []string
	doc.Find(".dict-book .trans-container li.word-exp").Each(func(i int, s *goquery.Selection){
		label := s.Find("span.pos").First().Text()
		trans := s.Find("span.trans").Text()
		result = append(result, label+" "+trans)
		//fmt.Println(label, trans)
	})
	/*
	doc.Find(".dict-book .trans-container .exam_type-value").Each(func(i int, s *goquery.Selection){
		exam := s.Text()
		fmt.Println(exam)
	})
	*/
	return result
}
func (doc *MyDoc) GetJsCode() error {
	var jsCode string
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "window.__NUXT__") {
			jsCode = s.Text()
		}
	})
	if len(jsCode) < 1 {
		return errors.New("Cannot find window.__NUXT__ in youdao response!")
	}
	doc.JSCode = jsCode
	//fmt.Println("✅ 成功提取 __NUXT__ 数据！")
	return nil
}
func (doc *MyDoc) ParseJsCode() error {
	if len(doc.JSCode) < 1 {
		err := doc.GetJsCode()
		if err != nil {
			return err
		}
	}
	vm := goja.New()
	_, err := vm.RunString("var window = {};")
	if err != nil {
		return err
	}
	_, err = vm.RunString(doc.JSCode)
	if err != nil {
		return err
	}
	
	windowObj := vm.Get("window").ToObject(vm)
	nuxtData := windowObj.Get("__NUXT__")
	if nuxtData == nil || goja.IsUndefined(nuxtData) || goja.IsNull(nuxtData) {
		return errors.New("cannot get nextData by running window.__NUXT__ js code.")
	}
	//fmt.Println("✅ 成功解析 __NUXT__ 数据！")
	var data map[string]interface{}
	err = vm.ExportTo(nuxtData, &data)
	if err != nil {
		fmt.Println("")
		return errors.New("nextData cannot export to golang interface.")
	}
	// data["data"][0].(map)["wordData"].(map)["collins"]
	if dataArr, ok := data["data"].([]interface{}); ok && len(dataArr) > 0 {
		if dataZero, ok := dataArr[0].(map[string]interface{}); ok {
			if wordData, ok := dataZero["wordData"].(map[string]interface{}); ok {
				doc.WordData = wordData
				return nil
			}
		}
	}
	//prettyJSON, _:= json.MarshalIndent(data, "", "  ")
	//fmt.Println(prettyJSON)
	return errors.New("nextData missing data key??")
}
func (doc *MyDoc) GetDicts() []string {
	var list []string
	type Meta struct {
		Dicts []string `json:"dicts"`
		
	}
	var m Meta
	if meta, ok := doc.WordData["meta"].(map[string]interface{}); ok {
		jsonBytes, _ := json.Marshal(meta)
		json.Unmarshal(jsonBytes, &m)
		return m.Dicts
		
	}
/*	

	if meta, ok := doc.WordData["meta"].(map[string]interface{}); ok {
		//fmt.Println(meta)
		if dicts, ok := meta["dicts"].([]interface{}); ok {
			for _, dict := range dicts {
				if dictStr, ok := dict.(string); ok {
					list=append(list, dictStr)
				}
			}
		}
			
	}
*/
	return list
}
// --- 提取简明释义 ---
// 路径: wordData["ec"]["word"]{}
func (doc *MyDoc) GetEC() EC {
	var t EC
	if ec, ok := doc.WordData["ec"].(map[string]interface{}); ok {
		// ec{"exam_type":[]}
		type Exam struct {
			Types []string `json:"exam_type"`
		}
		var exam Exam
		jsonBytes, _ := json.Marshal(ec)
		json.Unmarshal(jsonBytes, &exam)
		t.ExamTypes = exam.Types

		// ec{"word":[]}
		if word, ok := ec["word"].(map[string]interface{}); ok {
			//fmt.Println(word)
			if usphone, ok := word["usphone"].(string); ok {
				t.USPhone = usphone
			}
			if ukphone, ok := word["ukphone"].(string); ok {
				t.UKPhone = ukphone
			}
			if trs, ok := word["trs"].([]interface{}); ok {
				for _, tr := range trs {
					if trMap, ok := tr.(map[string]interface{}); ok {
						type Trans struct {
							Pos  string `json:"pos"`
							Tran  string `json:"tran"`
						}
						var trans Trans
						jsonBytes, _ := json.Marshal(trMap)
						json.Unmarshal(jsonBytes, &trans)
						tran := trans.Pos + " " + trans.Tran
						t.Trans = append(t.Trans, tran)
						/*
						if tran, ok := trMap["tran"].(string); ok {
							t.Trans = append(t.Trans, tran)
						}
						*/
					}
				}
			}
		}
	}
	return t
}
// 路径: wordData["collins"]["collins_entries"][0]["entries"]["entry"][0]["tran_entry"][0]["tran"]
func (doc *MyDoc) GetCollins() Collins {
	var c Collins
	if collins, ok := doc.WordData["collins"].(map[string]interface{}); ok {
		collinsEntries := collins["collins_entries"].([]interface{})
		entry0 := collinsEntries[0].(map[string]interface{})
		entries := entry0["entries"].(map[string]interface{})
		entryList := entries["entry"].([]interface{})
		for _, entry := range entryList {
			//fmt.Println(entry)
			entryObj := entry.(map[string]interface{})
			tranEntrys := entryObj["tran_entry"].([]interface{})
			for _, te := range tranEntrys {
				var ct CoTran

				teObj := te.(map[string]interface{})
				if tran, ok := teObj["tran"].(string); ok {
					//fmt.Println(tran)
					ct.Tran = tran
				}
				if posEntry, ok := teObj["pos_entry"].(map[string]interface{}); ok {
					pos := posEntry["pos"].(string)
					//fmt.Println(pos)
					ct.Pos = pos
				}
				if examSents, ok := teObj["exam_sents"].(map[string]interface{}); ok {
					//sents := examSents["sent"].([]interface{})
					jsonBytes, _ := json.Marshal(examSents)
					//fmt.Println(string(jsonBytes))
					json.Unmarshal(jsonBytes, &ct)
				}
				c.Trans = append(c.Trans, ct)
			}
		}
	}
	return c
}
// 路径: wordData["web_trans"]["web-translation"][0]["trans"][0] value cls summary[line]
func (doc *MyDoc) GetWebs() Webs {
	var ws Webs
	
	if webTrans, ok := doc.WordData["web_trans"].(map[string]interface{}); ok {
		webTranslations := webTrans["web-translation"].([]interface{})
		for _, wt := range webTranslations {
			wtObj := wt.(map[string]interface{})
			key := wtObj["key"].(string)
			var web Web
			web.Key = key
			trans := wtObj["trans"].([]interface{})
			for _, tran := range trans {
				var webT WebTrans
				tranObj := tran.(map[string]interface{})
				value := tranObj["value"].(string)
				//fmt.Println(value)
				webT.Value = value
				
				if cls, ok := tranObj["cls"].(map[string]interface{}); ok {
					jsonBytes, _ := json.Marshal(cls)
					//fmt.Println(string(jsonBytes))
					json.Unmarshal(jsonBytes, &webT)
					//fmt.Println(webT)
				}
				if summary, ok := tranObj["summary"].(map[string]interface{}); ok {
					jsonBytes, _ := json.Marshal(summary)
					//fmt.Println(string(jsonBytes))
					json.Unmarshal(jsonBytes, &webT)
					//fmt.Println(webT)
				}
				web.Trans = append(web.Trans, webT)
			}
			ws.Trans = append(ws.Trans, web)
			
		}
		
	}
	return ws
	
	
}

// 路径: wordData["spceial"]["entries"][0]["entry"] major
func (doc *MyDoc) GetSpecs() Specs {
	var ss Specs
	if special, ok := doc.WordData["special"].(map[string]interface{}); ok {
		entries := special["entries"].([]interface{})
		for _, entry := range entries {
			entryObj := entry.(map[string]interface{})
			//fmt.Println(`"wordData["web_trans"]["web-translation"]"`)
			var s Spec
			if one, ok := entryObj["entry"].(map[string]interface{}); ok {
				major := one["major"].(string)
				//fmt.Println(major)
				s.Major = major
				trs := one["trs"].([]interface{})
				for _, tr := range trs {
					trObj := tr.(map[string]interface{})
					if t, ok := trObj["tr"].(map[string]interface{}); ok {
						var st SpecTran
						jsonBytes, _ := json.Marshal(t)
						fmt.Println(string(jsonBytes))
						json.Unmarshal(jsonBytes, &st)
						fmt.Println(st)
						s.Trans = append(s.Trans, st)
					}
				}
				
			}
			ss.Trans = append(ss.Trans, s)
		}
	}
	return ss
}
// ee->word->trs[]->tr
func (doc *MyDoc) GetEEs() EEs {
	var ees EEs
	if eeDict, ok := doc.WordData["ee"].(map[string]interface{}); ok {
		wordDict := eeDict["word"].(map[string]interface{})
		trs := wordDict["trs"].([]interface{})
		for _, tr := range trs {
			trDict := tr.(map[string]interface{})
			var ee EE
			ee.Pos = trDict["pos"].(string)

			if trList, ok := trDict["tr"].([]interface{}); ok {
				for _, t := range trList{
					var eet EETran
					jsonBytes, _ := json.Marshal(t)
					json.Unmarshal(jsonBytes, &eet)
					ee.Trans = append(ee.Trans, eet)
				}
			}
			ees.Trans = append(ees.Trans, ee)
			
		}
	}
	return ees
	
}

func (doc *MyDoc) GetBlngSents() BlngSents {
	var sents BlngSents
	if blngDict, ok := doc.WordData["blng_sents_part"].(map[string]interface{}); ok {
		pairDict := blngDict["sentence-pair"].([]interface{})
		for _, pair := range pairDict {
			var sent BlngSent
			jsonBytes, _ := json.Marshal(pair)
			json.Unmarshal(jsonBytes, &sent)
			sents.Trans = append(sents.Trans, sent)
		}
	}
	return sents
	
}
func (doc *MyDoc) GetMediaSents() MediaSents {
	var sents MediaSents
	if mediaDict, ok := doc.WordData["media_sents_part"].(map[string]interface{}); ok {
		sentDicts := mediaDict["sent"].([]interface{})
		for _, sentDict := range sentDicts {
			sDict := sentDict.(map[string]interface{})
			var sent MediaSent
			if eng, ok := sDict["eng"].(string); ok {
				sent.Eng = eng
			}
			snippetsDict := sDict["snippets"].(map[string]interface{})
			snippetsDicts := snippetsDict["snippet"].([]interface{})
			for _, snippet := range snippetsDicts {
				jsonBytes, _ := json.Marshal(snippet)
				//fmt.Println(string(jsonBytes))
				json.Unmarshal(jsonBytes, &sent)
				break
			}
			sents.Trans = append(sents.Trans, sent)
			
		}
	}
	return sents
}

func (doc *MyDoc) GetAuthSents() AuthSents {
	var sents AuthSents
	if authDict, ok := doc.WordData["auth_sents_part"].(map[string]interface{}); ok {
		sentDicts := authDict["sent"].([]interface{})
		for _, sentDict := range sentDicts {
			var sent AuthSent
			jsonBytes, _ := json.Marshal(sentDict)
			json.Unmarshal(jsonBytes, &sent)
			sents.Trans = append(sents.Trans, sent)
		}
	}
	return sents
}
func (doc *MyDoc) GetPhRs() PhRs {
	var sents PhRs
	if phrsDict, ok := doc.WordData["phrs"].(map[string]interface{}); ok {
		phrsDicts := phrsDict["phrs"].([]interface{})
		for _, ph := range phrsDicts {
			var sent PhR
			jsonBytes, _ := json.Marshal(ph)
			json.Unmarshal(jsonBytes, &sent)
			sents.Trans = append(sents.Trans, sent)
		}
	}
	return sents
}
func (doc *MyDoc) GetSynos() Synos {
	var sents Synos
	if synoDict, ok := doc.WordData["syno"].(map[string]interface{}); ok {
		synosDicts := synoDict["synos"].([]interface{})
		for _, sy := range synosDicts {
			var sent Syno
			jsonBytes, _ := json.Marshal(sy)
			json.Unmarshal(jsonBytes, &sent)
			sents.Trans = append(sents.Trans, sent)
		}
	}
	return sents
}
func (doc *MyDoc) GetRelWords() RelWords {
	var sents RelWords
	if relDict, ok := doc.WordData["rel_word"].(map[string]interface{}); ok {
		relsDicts := relDict["rels"].([]interface{})
		for _, rels := range relsDicts {
			relsDict := rels.(map[string]interface{})
			var sent RelWord
			if rel, ok := relsDict["rel"].(map[string]interface{}); ok {
				if pos, ok := rel["pos"].(string); ok {
					sent.Pos = pos
				}
				if words, ok := rel["words"].([]interface{}); ok {
					for _, word := range words {
						jsonBytes, _ := json.Marshal(word)
						json.Unmarshal(jsonBytes, &sent)
						break // FIXME
					}
				}
				
			}
			sents.Trans = append(sents.Trans, sent)
		}
	}
	return sents
}

func (doc *MyDoc) GetETYMs() ETYMs {
	var sents ETYMs
	if relDict, ok := doc.WordData["etym"].(map[string]interface{}); ok {
		//relsDicts := relDict["etyms"].([]interface{})
		var root ETYMRoot 
		jsonBytes, _ := json.Marshal(relDict)
		json.Unmarshal(jsonBytes, &root)
		//fmt.Println(root)
		for lang, items := range root.ETYM {
			fmt.Println(lang, items)
			sents.Trans = append(sents.Trans, items...)
		}
	}
	return sents
}

func (doc *MyDoc) GetWIKIs() WIKIs {
	var sents WIKIs
	if wikiRoot, ok := doc.WordData["wikipedia_digest"].(map[string]interface{}); ok {
		source := wikiRoot["source"].(map[string]interface{})
		var wiki Wiki
		jsonBytes, _ := json.Marshal(source)
		json.Unmarshal(jsonBytes, &wiki)
		//fmt.Println(wiki)
		summarys := wikiRoot["summarys"].([]interface{})
		for _, summary := range summarys {
			jsonBytes, _ := json.Marshal(summary)
			json.Unmarshal(jsonBytes, &wiki)
			break // FIXME
		}
		sents.Trans = append(sents.Trans, wiki)
		
	}
	return sents
}








func (doc *MyDoc) GetExpandECs() ExpandECs {
	var expand_ecs ExpandECs
	if expandECDict, ok := doc.WordData["expand_ec"].(map[string]interface{}); ok {
		wordDictList := expandECDict["word"].([]interface{})
		for _, word := range wordDictList {
			wordDict := word.(map[string]interface{})
			if transList, ok := wordDict["transList"].([]interface{}); ok {
				
				for _, tran := range transList {
					tranDict := tran.(map[string]interface{})
					contentDict := tranDict["content"].(map[string]interface{})
					if sendsList, ok := contentDict["sents"].([]interface{}); ok {
						for _, send := range sendsList {
							var expandec  ExpandEC
							jsonBytes, _ := json.Marshal(send)
							json.Unmarshal(jsonBytes, &expandec)
							expand_ecs.Trans = append(expand_ecs.Trans, expandec)
						}
					}
				}
				
			}
			
		}
	}
	return expand_ecs
	
}
