package repository

// Pengirim email SMTP korespondensi (tiket 56) - persetujuan work owner 06-10-2026 "iya kirim email sungguhan" (CLAUDE.md
// §10). Konfigurasi HANYA dari env lewat inti/backend/config (`NBFACIN_SMTP_*`, ADR-0004): nol literal host / akun /
// sandi di kode. Port 465 = TLS implisit; selain itu STARTTLS. TLS WAJIB (keputusan work owner 06-10-2026): server tanpa
// STARTTLS -> tidak ada yang dikirim, kecuali NBFACIN_SMTP_TANPA_TLS=true DAN akun kosong (sandi tidak pernah lewat kanal
// terbuka; smtp.PlainAuth juga menolaknya). Badan text/html UTF-8 base64; Subject encoded-word UTF-8 terlipat. CR / LF di
// alamat dan subject ditolak layanan SEBELUM sampai ke sini (dan diperiksa ulang di bawah).

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// portSMTPBawaan - submission (STARTTLS) bila NBFACIN_SMTP_PORT kosong.
	portSMTPBawaan = 587
	// portSMTPTLS - TLS implisit.
	portSMTPTLS = 465
	// waktuSambungSMTP / waktuKirimSMTP - batas waktu sambung dan seluruh percakapan SMTP.
	waktuSambungSMTP = 15 * time.Second
	waktuKirimSMTP   = 60 * time.Second
)

// ErrSMTPTanpaTLS - server tidak menawarkan STARTTLS (port bukan 465) dan pengiriman tanpa TLS tidak diizinkan (bawaan, atau
// akun diisi): sambungan diputus SEBELUM MAIL FROM.
var ErrSMTPTanpaTLS = errors.New("repository: server SMTP tidak mendukung TLS (STARTTLS); email tidak dikirim")

// PengirimSurel - satu email keluar.
type PengirimSurel interface {
	KirimSurel(ctx context.Context, s models.Surel) error
}

// KonfigSMTP - NBFACIN_SMTP_* (dibaca inti/backend/config, dirakit modul.go).
type KonfigSMTP struct {
	Host, User, Sandi, Dari string
	Port                    int // 0 = 587; 465 = TLS implisit
	TanpaTLS                bool
}

// PengirimSMTP - PengirimSurel atas net/smtp.
type PengirimSMTP struct{ k KonfigSMTP }

// NewPengirimSMTP merakit pengirim; port 0 = 587.
func NewPengirimSMTP(k KonfigSMTP) *PengirimSMTP {
	if k.Port == 0 {
		k.Port = portSMTPBawaan
	}
	return &PengirimSMTP{k: k}
}

// adaCRLF - CR / LF (injeksi header).
func adaCRLF(s string) bool { return strings.ContainsAny(s, "\r\n") }

// SusunSurel - pesan RFC 5322 lengkap (header CRLF, badan base64 76 kolom). waktu = header Date.
func SusunSurel(dari string, s models.Surel, waktu time.Time) (string, error) {
	semua := append(append([]string{dari, s.ReplyTo, s.Subject}, s.To...), s.Cc...)
	for _, v := range semua {
		if adaCRLF(v) {
			return "", errors.New("repository: CR / LF di header email ditolak")
		}
	}
	var b strings.Builder
	tulis := func(k, v string) { b.WriteString(k + ": " + v + "\r\n") }
	tulis("From", dari)
	tulis("To", strings.Join(s.To, ", "))
	if len(s.Cc) > 0 {
		tulis("Cc", strings.Join(s.Cc, ", "))
	}
	if s.ReplyTo != "" {
		tulis("Reply-To", s.ReplyTo)
	}
	tulis("Subject", subjekTerlipat(s.Subject))
	tulis("Date", waktu.Format(time.RFC1123Z))
	tulis("MIME-Version", "1.0")
	tulis("Content-Type", "text/html; charset=UTF-8")
	tulis("Content-Transfer-Encoding", "base64")
	b.WriteString("\r\n")
	isi := base64.StdEncoding.EncodeToString([]byte(s.HTML))
	for len(isi) > 76 {
		b.WriteString(isi[:76] + "\r\n")
		isi = isi[76:]
	}
	b.WriteString(isi + "\r\n")
	return b.String(), nil
}

// subjekTerlipat - Subject aman batas baris 998 (RFC 5322): ASCII cetak pendek (≤ 900 bita) apa adanya; selain itu
// encoded-word "B" UTF-8 berpotongan ≤ 45 bita (pada batas rune), dilipat CRLF + spasi - tidak ada baris yang panjang.
func subjekTerlipat(s string) string {
	polos := len(s) <= 900
	for i := 0; i < len(s) && polos; i++ {
		polos = s[i] >= 0x20 && s[i] < 0x7f && !strings.HasPrefix(s[i:], "=?")
	}
	if polos {
		return s
	}
	var potong []string
	for len(s) > 0 {
		n := len(s)
		if n > 45 {
			n = 45
			for n > 0 && !utf8.RuneStart(s[n]) {
				n--
			}
		}
		potong = append(potong, "=?UTF-8?B?"+base64.StdEncoding.EncodeToString([]byte(s[:n]))+"?=")
		s = s[n:]
	}
	return strings.Join(potong, "\r\n ")
}

// KirimSurel - lihat PengirimSurel. Penerima amplop = To + Cc.
func (p *PengirimSMTP) KirimSurel(ctx context.Context, s models.Surel) error {
	k := p.k
	pesan, err := SusunSurel(k.Dari, s, time.Now())
	if err != nil {
		return err
	}
	alamat := net.JoinHostPort(k.Host, strconv.Itoa(k.Port))
	tlsKonfig := &tls.Config{ServerName: k.Host, MinVersion: tls.VersionTLS12}
	penyambung := &net.Dialer{Timeout: waktuSambungSMTP}
	var conn net.Conn
	if k.Port == portSMTPTLS {
		conn, err = (&tls.Dialer{NetDialer: penyambung, Config: tlsKonfig}).DialContext(ctx, "tcp", alamat)
	} else {
		conn, err = penyambung.DialContext(ctx, "tcp", alamat)
	}
	if err != nil {
		return fmt.Errorf("repository: sambung SMTP: %w", err)
	}
	_ = conn.SetDeadline(time.Now().Add(waktuKirimSMTP))
	c, err := smtp.NewClient(conn, k.Host)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("repository: sapa SMTP: %w", err)
	}
	defer c.Close()
	if k.Port != portSMTPTLS {
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err := c.StartTLS(tlsKonfig); err != nil {
				return fmt.Errorf("repository: STARTTLS: %w", err)
			}
		} else if !k.TanpaTLS || k.User != "" {
			return ErrSMTPTanpaTLS
		}
	}
	if k.User != "" {
		if err := c.Auth(smtp.PlainAuth("", k.User, k.Sandi, k.Host)); err != nil {
			return fmt.Errorf("repository: autentikasi SMTP gagal: %w", err)
		}
	}
	if err := c.Mail(k.Dari); err != nil {
		return fmt.Errorf("repository: MAIL FROM: %w", err)
	}
	for _, r := range append(append([]string{}, s.To...), s.Cc...) {
		if err := c.Rcpt(r); err != nil {
			return fmt.Errorf("repository: RCPT TO ditolak: %w", err)
		}
	}
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("repository: DATA: %w", err)
	}
	if _, err := w.Write([]byte(pesan)); err != nil {
		return fmt.Errorf("repository: tulis pesan: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("repository: pesan ditolak server: %w", err)
	}
	return c.Quit()
}
