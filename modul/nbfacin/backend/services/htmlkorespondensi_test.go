package services

import "testing"

// TestBersihkanHTMLKorespondensiParitas - INSTRUMEN: jawaban dihasilkan oleh pembersih FRONTEND
// (`frontend/components/htmlKorespondensi.ts` `bersihkanHtml`, dijalankan Node 24 06-10-2026) atas 22 masukan, termasuk XSS
// (img onerror, javascript:, <SCRIPT>, svg onload, iframe, tag bersarang, script tanpa penutup). Port Go wajib sama persis.
func TestBersihkanHTMLKorespondensiParitas(t *testing.T) {
	for _, k := range []struct{ masuk, mau string }{
		{"Halo <b>dunia</b>", "Halo <b>dunia</b>"},
		{"<p style=\"color:red\" onclick=\"x()\">teks</p>", "<p>teks</p>"},
		{"<img src=x onerror=alert(1)>gambar", "gambar"},
		{"<a href=\"javascript:alert(1)\">klik</a>", "klik"},
		{"<SCRIPT>alert(1)</SCRIPT>aman", "aman"},
		{"<script>alert(1)", ""},
		{"<style>p{color:red}</style><p>a</p>", "<p>a</p>"},
		{"<!-- rahasia --><i>x</i>", "<i>x</i>"},
		{"a & b &amp; c &#169; &#xA9; &bogus &nbsp; < > \"kutip\"", "a &amp; b &amp; c &#169; &#xA9; &amp;bogus &nbsp; &lt; &gt; \"kutip\""},
		{"<br><br/><BR />baris", "<br><br><br>baris"},
		{"<div><ul><li>satu</li><li>dua</li></ul><ol><li>tiga</li></ol></div>", "<div><ul><li>satu</li><li>dua</li></ul><ol><li>tiga</li></ol></div>"},
		{"<iframe src=\"//jahat\"></iframe>teks", "teks"},
		{"<svg onload=alert(1)><circle/></svg>", ""},
		{"<p>tanpa tutup", "<p>tanpa tutup"},
		{"<STRONG>besar</STRONG><EM>miring</EM><u>garis</u>", "<strong>besar</strong><em>miring</em><u>garis</u>"},
		{"x < y > z", "x &lt; y &gt; z"},
		{"<scr<script>ipt>alert(1)</script>", "ipt&gt;alert(1)"},
		{"<p\nclass=\"a\">baris baru di tag</p>", "<p>baris baru di tag</p>"},
		{"&#1234567; &#12345678; &abcdefghijklmnopqrstuvwxyz0123456789ab;", "&#1234567; &amp;#12345678; &amp;abcdefghijklmnopqrstuvwxyz0123456789ab;"},
		{"teks biasa saja", "teks biasa saja"},
		{"", ""},
		{"<b onmouseover=\"x\">tebal</b><span>span dibuang</span>", "<b>tebal</b>span dibuang"},
	} {
		if got := BersihkanHTMLKorespondensi(k.masuk); got != k.mau {
			t.Errorf("%q:\n got %q\n mau %q", k.masuk, got, k.mau)
		}
	}
}
