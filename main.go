package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var mds = `# header
|Sample《サンプル》 text.
[link](http://example.com)

` + "```" + `
$ sudo rm --no-preserve-root -rf /
` + "```" + `

**asf**
~~**asf~~

## あｓｄｆasdf
|Sample《サンプル》 text.
[link](http://example.com)
### あｓｄｆasdf
|Sample《サンプル》 text.
[link](http://example.com)
#### asdf
|Sample《サンプル》 text.
[link](http://example.com)
##### asdf
|Sample《サンプル》 text.
[link](http://example.com)
##### asdf
|Sample《サンプル》 text.
[link](http://example.com)
###### asdf
|Sample《サンプル》 text.
[link](http://example.com)
####### asdf
|Sample《サンプル》 text.
[link](http://example.com)


　私が先生と知り合いになったのは<ruby>鎌倉<rt>かまくら</rt></ruby>である。その時私はまだ若々しい書生であった。暑中休暇を利用して海水浴に行った友達からぜひ来いという<ruby>端書<rt>はがき</rt></ruby>を受け取ったので、私は多少の金を<ruby>工面<rt>くめん</rt></ruby>して、出掛ける事にした。私は金の工面に<ruby>二<rt>に</rt></ruby>、<ruby>三日<rt>さんち</rt></ruby>を費やした。ところが私が鎌倉に着いて三日と経たたないうちに、私を呼び寄せた友達は、急に国元から帰れという電報を受け取った。電報には母が病気だからと断ってあったけれども友達はそれを信じなかった。友達はかねてから国元にいる親たちに勧すすまない結婚を強しいられていた。彼は現代の習慣からいうと結婚するにはあまり年が若過ぎた。それに肝心かんじんの当人が気に入らなかった。それで夏休みに当然帰るべきところを、わざと避けて東京の近くで遊んでいたのである。彼は電報を私に見せてどうしようと相談をした。私にはどうしていいか分らなかった。けれども実際彼の母が病気であるとすれば彼は固もとより帰るべきはずであった。それで彼はとうとう帰る事になった。せっかく来た私は一人取り残された。
`

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	return markdown.Render(doc, renderer)
}

func divMd() {
}

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile) // ログの出力書式を設定する

	var input string
	var sc = bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		input += sc.Text()
	}

	md := []byte(input)
	//md := []byte(mds)
	html := mdToHTML(md)

	/*
		前処理する。
		改行を消す。
		別にhtmlタグを削除した文字列を生成して1ページに表示する文字列をカウントし区切る。
	*/
	output := ""
	buf := bytes.NewBufferString(string(html))
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		t := scanner.Text()
		t = strings.Replace(t, "\n", "", -1)

		for { // 青空文庫形式のルビを処理 1行に複数個ある場合を想定して1つづつ処理する。
			e := strings.Contains(t, "|") && strings.Contains(t, "《") && strings.Contains(t, "》")
			if e {
				t = strings.Replace(t, "|", "<ruby>", 1)
				t = strings.Replace(t, "《", "<rt>", 1)
				t = strings.Replace(t, "》", "</rt></ruby>", 1)
			} else {
				break
			}
		}
		// 読まない句読点などにこっそりルビを振る。
		str := "。、「」（）"
		for _, r := range str {
			c := string(r)
			if strings.Contains(t, c) {
				t = strings.Replace(t, c, "<ruby>"+c+"<rt>&nbsp;</rt></ruby>", -1)
			}
		}
		// 12月31日のような日付文字列の数字を１文字の幅に2文字入れる。
		//t = rmonth.ReplaceAllString(t, "")
		output += t
	}
	/*
		// htmlタグを除去する
		notag := func(in string) string {
			// htmlタグを除去する前にhtmlの中で改行されるタグに改行代わりに
			// 他で絶対使われない文字を入れる。﹆白ゴマらしい
			t := strings.Replace(in, "</", "﹆</", -1)
			// htmlタグを除去する
			t = bluemonday.StripTagsPolicy().Sanitize(t)
			// 改行代わりの文字が複数個並んでいるところを１つにする
			for {
				t = strings.Replace(t, "﹆﹆", "﹆", -1)
				if !strings.Contains(t, "﹆﹆") {
					break
				}
			}
			// 改行代わりの文字を改行にする
			return strings.Replace(t, "﹆", "\n", -1)
		}(output)

		// 一行ごとに文字列の表示長さを計算する。
		s := strings.Split(notag, "\n")
		notag = ""
		for _, line := range s {
			notag += runewidth.Wrap(line, 100) + "\n"
		}
	*/

	//fmt.Printf("--- Markdown:\n%s\n\n--- HTML:\n%s\n", md, output)
	//fmt.Printf("--- Markdown:\n%s\n\n--- HTML:\n%s\n\n--- no tag:\n%s\n", md, output, notag)
	fmt.Printf("%s\n", output)
}
