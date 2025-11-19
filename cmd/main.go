package main
import (
	"fmt"
	"os"
	"slices"
	"github.com/fatih/color"
	"github.com/yorks/youdaoweb"
)


func main(){
   if len(os.Args) < 2 {
	fmt.Printf("Usage: %s word\n", os.Args[0])
	return 
   }
   doc, err := youdaoweb.Qword(os.Args[1])
   if err != nil {
	fmt.Println(err)
   }else{
	phonetics := doc.GetPhonetic()
	for lang, ph := range phonetics {
		fmt.Print(lang)
		fmt.Print(color.CyanString(ph)+" ")
	}
	fmt.Println()
	trans := doc.GetTrans()
	for _, tran := range trans {
		fmt.Print(tran + " ")
	}
	fmt.Println()

	err = doc.ParseJsCode()
	if err != nil {
		fmt.Println(err)
		return
	}
	dicts := doc.GetDicts()
	if slices.Contains(dicts, "collins"){
		collins := doc.GetCollins()
		if len(collins.Trans) > 0 {
			color.New(color.BgGreen).Add(color.Underline, color.Bold).Println("\t  柯林斯 ")
		}
		for _, col := range collins.Trans {
			fmt.Print("["+col.Pos+"]")
			fmt.Println(youdaoweb.RenderHTML(col.Tran))
			for _, sent := range col.ExamSents {
				color.White(youdaoweb.RenderHTML(sent.Eng))
				color.White(sent.Chn)
			}
		}
	}
	if slices.Contains(dicts, "web_trans"){
		webs := doc.GetWebs()
		if len(webs.Trans) > 0 {
			color.New(color.BgGreen).Add(color.Underline, color.Bold).Println("\t  网络释义 ")
		}
		isShort := false
		shorted := false
		for _, wt := range webs.Trans {
			if isShort && !shorted {
				color.New(color.BgGreen).Add(color.Underline, color.Bold).Println("\t  短语 ")
				shorted = true
			}
			fmt.Print(color.YellowString(wt.Key)+" ")
			for _, w := range wt.Trans {
				if len(w.Cls)>0 {
					fmt.Print(w.Cls)
				}else{
					isShort = true
				}
				fmt.Print(w.Value)
				for _, summ := range w.Summarys {
					fmt.Print(youdaoweb.RenderHTML(summ)+" ")
				}
			}
			fmt.Println()
		}
	}

	if slices.Contains(dicts, "blng_sents_part"){
		bl := doc.GetBlngSents()
		if len(bl.Trans) > 0 {
			color.New(color.BgGreen).Add(color.Underline, color.Bold).Println("\t  双语例句 ")
		}
		for _, bs := range bl.Trans {
			color.Yellow(youdaoweb.RenderHTML(bs.Sentence))
			color.White(bs.Source + bs.Chn)
		}
	}
	if slices.Contains(dicts, "ee"){
		ees := doc.GetEEs()
		if len(ees.Trans) > 0 {
			color.New(color.BgGreen).Add(color.Underline, color.Bold).Println("\t  英英释义 ")
		}
		for _, ee := range ees.Trans {
			fmt.Print(ee.Pos)
			for _, t := range ee.Trans {
				color.Yellow(t.Tran)
				if len(t.Swords) > 0 {
					fmt.Print(color.WhiteString("同义词: "))
					for _, w := range t.Swords {
						fmt.Print(color.CyanString(w) + " ")
					}
					fmt.Println()
				}
			}
		}
	}
	if slices.Contains(dicts, "phrs"){
		phrs := doc.GetPhRs()
		if len(phrs.Trans) > 0 {
			color.New(color.BgGreen).Add(color.Underline, color.Bold).Println("\t  词典短语 ")
		}
		for _, phr := range phrs.Trans {
			fmt.Print(color.CyanString(phr.Word))
			color.White(phr.Tran)
		}
	}
	if slices.Contains(dicts, "syno"){
		synos := doc.GetSynos()
		if len(synos.Trans) > 0 {
			color.New(color.BgGreen).Add(color.Underline, color.Bold).Println("\t  同近义词 ")
		}
		for _, syno := range synos.Trans {
			fmt.Print(color.CyanString(syno.Pos))
			fmt.Println(syno.Tran)
			for i, w := range syno.WS {
				fmt.Print( color.CyanString(w) )
				if i < len(syno.WS)-1 {
					fmt.Print( " / ")
				}
			}
			if len(syno.WS) > 0 {
				fmt.Println()
			}
		}
	}

	if slices.Contains(dicts, "rel_word"){
		rels := doc.GetRelWords()
		if len(rels.Trans) > 0 {
			color.New(color.BgGreen).Add(color.Underline, color.Bold).Println("\t  同根词 ")
		}
		for _, rel := range rels.Trans {
			fmt.Print(color.WhiteString(rel.Pos)+" ")
			fmt.Print(color.CyanString(rel.Word)+" ")
			color.White(rel.Tran)
		}
	}
	if slices.Contains(dicts, "etym"){
		etyms := doc.GetETYMs()
		if len(etyms.Trans) > 0 {
			color.New(color.BgGreen).Add(color.Underline, color.Bold).Println("\t  词源 ")
		}
		for _, etym := range etyms.Trans {
			fmt.Print(color.CyanString(etym.Word)+" ")
			fmt.Print(color.WhiteString(etym.Desc)+" ")
			color.White(etym.Value)
		}
	}
	if slices.Contains(dicts, "media_sents_part"){
		sents := doc.GetMediaSents()
		if len(sents.Trans) > 0 {
			color.New(color.BgGreen).Add(color.Underline, color.Bold).Println("\t  原声例句 ")
		}
		for _, sent := range sents.Trans {
			fmt.Print(color.YellowString(sent.Eng))
			fmt.Print(color.WhiteString(sent.Source+" "+sent.Name+" "))
			color.White(sent.StreamUrl)
		}
	}
	if slices.Contains(dicts, "auth_sents_part"){
		sents := doc.GetAuthSents()
		if len(sents.Trans) > 0 {
			color.New(color.BgGreen).Add(color.Underline, color.Bold).Println("\t  权威例句 ")
		}
		for _, sent := range sents.Trans {
			fmt.Print(color.YellowString(youdaoweb.RenderHTML(sent.Foreign)))
			fmt.Print(color.WhiteString(youdaoweb.RenderHTML(sent.Source)))
			fmt.Println()
		}
	}
	if slices.Contains(dicts, "wikipedia_digest"){
		wikis := doc.GetWIKIs()
		if len(wikis.Trans) > 0 {
			color.New(color.BgGreen).Add(color.Underline, color.Bold).Println("\t  WIKI ")
		}
		for _, w := range wikis.Trans {
			fmt.Print(color.CyanString(w.Key)+" ")
			color.White(w.Url)
			fmt.Print(color.WhiteString(youdaoweb.RenderHTML(w.Summary)))
			fmt.Println()
		}
	}

	//fmt.Println(doc.GetEC())
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
}
