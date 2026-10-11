package repository

import (
	"bufio"
	"context"
	"encoding/base64"
	"errors"
	"mime"
	"net"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestSusunSurel - header RFC 5322 (CRLF), Cc / Reply-To, Subject encoded-word UTF-8, badan HTML base64 ≤ 76 kolom;
// CR / LF di header ditolak.
func TestSusunSurel(t *testing.T) {
	s := models.Surel{To: []string{"uji.a@contoh.test", "uji.b@contoh.test"}, Cc: []string{"uji.p@contoh.test"},
		ReplyTo: "uji.p@contoh.test", Subject: "UJI Penawaran – Kebakaran", HTML: "<p>" + strings.Repeat("isi ", 40) + "</p>"}
	waktu := time.Date(2026, 10, 6, 3, 4, 0, 0, time.UTC)
	p, err := SusunSurel("uji.sistem@contoh.test", s, waktu)
	if err != nil {
		t.Fatal(err)
	}
	kepala, badan, _ := strings.Cut(p, "\r\n\r\n")
	for _, harus := range []string{"From: uji.sistem@contoh.test", "To: uji.a@contoh.test, uji.b@contoh.test", "Cc: uji.p@contoh.test",
		"Reply-To: uji.p@contoh.test", "Subject: =?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(s.Subject)) + "?=", "Date: Tue, 06 Oct 2026 03:04:00 +0000",
		"MIME-Version: 1.0", "Content-Type: text/html; charset=UTF-8", "Content-Transfer-Encoding: base64"} {
		if !strings.Contains(kepala+"\r\n", harus+"\r\n") {
			t.Errorf("kepala tanpa %q:\n%s", harus, kepala)
		}
	}
	var gabung string
	for _, b := range strings.Split(strings.TrimRight(badan, "\r\n"), "\r\n") {
		if len(b) > 76 {
			t.Errorf("baris badan %d kolom", len(b))
		}
		gabung += b
	}
	if isi, err := base64.StdEncoding.DecodeString(gabung); err != nil || string(isi) != s.HTML {
		t.Errorf("badan: %v %q", err, isi)
	}
	// Subject panjang: setiap baris ≤ 998 (RFC 5322) - ASCII 1000 bita dan multibita dilipat menjadi encoded-word; isinya utuh.
	for _, subjek := range []string{strings.Repeat("a", 1000), strings.Repeat("é", 400), "pendek biasa"} {
		p, err := SusunSurel("uji.sistem@contoh.test", models.Surel{To: []string{"uji.a@contoh.test"}, Subject: subjek}, waktu)
		if err != nil {
			t.Fatal(err)
		}
		crlf := string([]byte{13, 10})
		kepala, _, _ := strings.Cut(p, crlf+crlf)
		var nilai string
		for _, b := range strings.Split(kepala, crlf) {
			if len(b) > 998 {
				t.Errorf("baris kepala %d kolom", len(b))
			}
			if strings.HasPrefix(b, "Subject: ") {
				nilai = strings.TrimPrefix(b, "Subject: ")
			} else if strings.HasPrefix(b, " =?") && nilai != "" {
				nilai += b // baris lipatan encoded-word
			}
		}
		if dek, err := (&mime.WordDecoder{}).DecodeHeader(nilai); err != nil || dek != subjek {
			t.Errorf("subjek %d bita tidak utuh: %v", len(subjek), err)
		}
	}
	for _, rusak := range []models.Surel{{To: []string{"a@contoh.test\r\nBcc: b@contoh.test"}}, {Subject: "a\nBcc: x"}, {ReplyTo: "a\r"}} {
		if _, err := SusunSurel("uji.sistem@contoh.test", rusak, waktu); err == nil {
			t.Errorf("CRLF lolos: %+v", rusak)
		}
	}
}

// TestSQLKorespondensi - kolom tulis / baca ada di migrasi 207; jumlah bind = kolom; status draf / terkirim; saran akun
// aktif ber-EMAIL.
func TestSQLKorespondensi(t *testing.T) {
	kolom := map[string]map[string]bool{}
	for _, m := range regexp.MustCompile(`(?s)CREATE TABLE \{skema\}\.(\w+) \((.*?)\n\)`).FindAllStringSubmatch(bacaMigrasi(t, "207_korespondensi.sql"), -1) {
		kolom[m[1]] = map[string]bool{}
		for _, b := range strings.Split(m[2], "\n") {
			if f := strings.Fields(b); len(f) > 0 && f[0] != "CONSTRAINT" {
				kolom[m[1]][f[0]] = true
			}
		}
	}
	if len(kolom[TabelKorespondensi]) != 13 || len(kolom[TabelPenerimaKorespondensi]) != 7 {
		t.Fatalf("migrasi 207 %v", kolom)
	}
	tb := tabelKorespondensi{work: "W", general: "G", kor: TabelKorespondensi, pen: TabelPenerimaKorespondensi, login: "L"}
	for tabel, q := range map[string]string{TabelKorespondensi: sqlSisipKorespondensi(tb), TabelPenerimaKorespondensi: sqlSisipPenerima(tb.pen)} {
		m := regexp.MustCompile(`INSERT INTO \w+ \(([^)]*)\)`).FindStringSubmatch(q)
		for _, c := range strings.Split(m[1], ", ") {
			if !kolom[tabel][c] {
				t.Errorf("%s.%s ditulis, tidak ada di 207", tabel, c)
			}
		}
	}
	if n := strings.Count(sqlSisipKorespondensi(tb), ":"); n != 10 {
		t.Errorf("bind sisip %d, mau 10", n)
	}
	for _, m := range regexp.MustCompile(`\bk\.([A-Z_]+)`).FindAllStringSubmatch(sqlDaftarKorespondensi(tb), -1) {
		if !kolom[TabelKorespondensi][m[1]] {
			t.Errorf("%s dibaca, tidak ada", m[1])
		}
	}
	if !strings.Contains(sqlUbahKorespondensi(tb.kor), "IS_EDITABLE = '1'") || !strings.Contains(sqlTandaiTerkirim(tb.kor), "SET IS_EDITABLE = '0'") ||
		!strings.Contains(sqlSaranPenerima("L", true), "IS_ACTIVE = '1' AND EMAIL IS NOT NULL AND (UPPER(NAME) LIKE :1") {
		t.Error("status / saran")
	}
}

// pelayanSMTPUji - server SMTP minimal di loopback (TANPA STARTTLS) yang merekam DATA; tidak pernah ke luar mesin.
func pelayanSMTPUji(t *testing.T, terima chan<- string) (string, int) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	go func() {
		c, err := l.Accept()
		if err != nil {
			return
		}
		defer c.Close()
		r, w := bufio.NewReader(c), bufio.NewWriter(c)
		balas := func(s string) { w.WriteString(s + "\r\n"); w.Flush() }
		balas("220 uji")
		var data strings.Builder
		for {
			baris, err := r.ReadString('\n')
			if err != nil {
				return
			}
			perintah := strings.ToUpper(strings.TrimSpace(baris))
			switch {
			case strings.HasPrefix(perintah, "EHLO"):
				balas("250 uji")
			case strings.HasPrefix(perintah, "MAIL"), strings.HasPrefix(perintah, "RCPT"):
				data.WriteString(strings.TrimSpace(baris) + "\n")
				balas("250 ok")
			case perintah == "DATA":
				balas("354 lanjut")
				for {
					b, err := r.ReadString('\n')
					if err != nil || b == ".\r\n" {
						break
					}
					data.WriteString(b)
				}
				balas("250 diterima")
			case perintah == "QUIT":
				balas("221 sampai jumpa")
				terima <- data.String()
				return
			default:
				balas("502 tidak dikenal")
			}
		}
	}()
	host, port, _ := net.SplitHostPort(l.Addr().String())
	n, _ := strconv.Atoi(port)
	return host, n
}

// TestKirimSurelLoopback - server tiruan loopback TANPA STARTTLS: bawaan (TLS wajib) -> ErrSMTPTanpaTLS dan MAIL FROM tidak
// pernah terkirim; NBFACIN_SMTP_TANPA_TLS=true -> percakapan utuh (amplop = To + Cc, pesan tersusun); akun diisi -> tetap
// ErrSMTPTanpaTLS walau TanpaTLS (sandi tidak pernah lewat kanal terbuka).
func TestKirimSurelLoopback(t *testing.T) {
	terima := make(chan string, 1)
	host, port := pelayanSMTPUji(t, terima)
	s := models.Surel{To: []string{"uji.a@contoh.test"}, Cc: []string{"uji.p@contoh.test"}, Subject: "UJI", HTML: "<b>x</b>"}
	dari := "uji.sistem@contoh.test"
	tolak := make(chan string, 1)
	h0, p0 := pelayanSMTPUji(t, tolak)
	if err := NewPengirimSMTP(KonfigSMTP{Host: h0, Port: p0, Dari: dari}).KirimSurel(context.Background(), s); !errors.Is(err, ErrSMTPTanpaTLS) {
		t.Errorf("TLS wajib: %v", err)
	}
	select {
	case d := <-tolak:
		t.Errorf("server tanpa TLS menerima percakapan: %q", d)
	case <-time.After(200 * time.Millisecond):
	}
	if err := NewPengirimSMTP(KonfigSMTP{Host: host, Port: port, Dari: dari, TanpaTLS: true}).KirimSurel(context.Background(), s); err != nil {
		t.Fatal(err)
	}
	got := <-terima
	for _, harus := range []string{"MAIL FROM:<uji.sistem@contoh.test>", "RCPT TO:<uji.a@contoh.test>", "RCPT TO:<uji.p@contoh.test>",
		"To: uji.a@contoh.test", "Cc: uji.p@contoh.test"} {
		if !strings.Contains(got, harus) {
			t.Errorf("tanpa %q:\n%s", harus, got)
		}
	}
	host, port = pelayanSMTPUji(t, make(chan string, 1))
	k := KonfigSMTP{Host: host, Port: port, User: "uji-akun", Sandi: "uji-sandi", Dari: dari, TanpaTLS: true}
	if err := NewPengirimSMTP(k).KirimSurel(context.Background(), s); !errors.Is(err, ErrSMTPTanpaTLS) {
		t.Errorf("tanpa TLS: %v", err)
	}
}
