package services

// Pembersih HTML pesan korespondensi (tiket 56) - port SATU-SATU `frontend/components/htmlKorespondensi.ts`
// (`bersihkanHtml`): hanya tag di tagHTMLKorespondensi, TANPA atribut apa pun; isi `script` / `style` dibuang utuh (tanpa
// penutup = sisa dibuang); tag lain dibuang tetapi teksnya dipertahankan; komentar dibuang; teks di-escape (`&` kecuali
// entitas sah `&name;` `&#n;` `&#xh;`, `<`, `>`). Hasil hanya tersusun dari tag polos yang diizinkan + teks ter-escape,
// jadi aman sebagai innerHTML dan badan email. Dijalankan server SEBELUM simpan dan kirim.

import (
	"regexp"
	"strings"
)

// tagHTMLKorespondensi - `TAG_HTML_KORESPONDENSI` frontend.
var tagHTMLKorespondensi = map[string]bool{"b": true, "strong": true, "i": true, "em": true, "u": true, "p": true, "br": true,
	"div": true, "ul": true, "ol": true, "li": true}

var (
	isiDibuangKorespondensi = map[string]bool{"script": true, "style": true}
	polaEntitasHTML         = regexp.MustCompile(`^&(?:[a-zA-Z][a-zA-Z0-9]{1,31}|#\d{1,7}|#[xX][0-9a-fA-F]{1,6});`)
	polaTagHTML             = regexp.MustCompile(`<(/?)([a-zA-Z][a-zA-Z0-9]*)\b[^>]*>|<!--[\s\S]*?-->`)
)

// escapeTeksHTML - `escapeTeks` frontend.
func escapeTeksHTML(teks string) string {
	var b strings.Builder
	for i := 0; i < len(teks); i++ {
		switch c := teks[i]; c {
		case '&':
			if m := polaEntitasHTML.FindString(teks[i:]); m != "" {
				b.WriteString(m)
				i += len(m) - 1
			} else {
				b.WriteString("&amp;")
			}
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}

// BersihkanHTMLKorespondensi - `bersihkanHtml` frontend.
func BersihkanHTMLKorespondensi(html string) string {
	var hasil strings.Builder
	buang := ""
	akhir := 0
	for _, m := range polaTagHTML.FindAllStringSubmatchIndex(html, -1) {
		if buang == "" {
			hasil.WriteString(escapeTeksHTML(html[akhir:m[0]]))
		}
		akhir = m[1]
		if m[4] < 0 {
			continue // komentar
		}
		tutup := m[3] > m[2]
		nama := strings.ToLower(html[m[4]:m[5]])
		if buang != "" {
			if tutup && nama == buang {
				buang = ""
			}
			continue
		}
		if isiDibuangKorespondensi[nama] {
			if !tutup {
				buang = nama
			}
			continue
		}
		if !tagHTMLKorespondensi[nama] {
			continue
		}
		switch {
		case nama == "br":
			if !tutup {
				hasil.WriteString("<br>")
			}
		case tutup:
			hasil.WriteString("</" + nama + ">")
		default:
			hasil.WriteString("<" + nama + ">")
		}
	}
	if buang == "" {
		hasil.WriteString(escapeTeksHTML(html[akhir:]))
	}
	return hasil.String()
}
