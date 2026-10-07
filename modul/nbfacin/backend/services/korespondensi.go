package services

// Tab Correspondence kasus Inward Facultative (tiket 56; kasus apa pun - Pega S95 selalu tampil). Bukti `[terverifikasi]`
// `D:\migrasi\RNM\NB FacIn\`: `DataTransform\CorrespondenceAdd.xml` (draf: IsEditable 1, Sender / SenderMail = operator
// login, Date = sekarang), `Section\CorrespondenceContent.xml` (Save, Send Email; lampiran tersembunyi `1=2`),
// `Activity\SendCorrespondence.xml` (kirim lalu IsEditable 0).
// Keputusan work owner 06-10-2026 "iya kirim email sungguhan. saya setuju dengan rekomendasi kamu" (persetujuan manusia
// CLAUDE.md §10 untuk SMTP nyata): To = penerima (Pega To = pengirim sendiri, CC = penerima - DIBETULKAN), Cc + Reply-To =
// pengirim; Send menyimpan dulu (Pega tanpa Obj-Save - dibetulkan); lampiran tidak dibuat; penerima = akun M_LOGIN_GO.
// Keputusan work owner 06-10-2026 "ini visibility menu nya juga cuma it ya": seluruh tab (baca, tulis, kirim, saran
// penerima) hanya JOB_POSITION "IT" - aturan sama dengan Loss Record (WajibAksesKorespondensi).

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

var (
	// ErrMasukanKorespondensi - isian tidak sah / penerima email tidak sah. 400.
	ErrMasukanKorespondensi = errors.New("services: isian korespondensi tidak sah")
	// ErrKorespondensiTidakAda - nomor bukan milik kasus. 404.
	ErrKorespondensiTidakAda = errors.New("services: korespondensi tidak ada di kasus ini")
	// ErrKorespondensiTerkirim - sudah terkirim (terkunci). 409.
	ErrKorespondensiTerkirim = errors.New("services: korespondensi sudah terkirim dan tidak dapat diubah")
	// ErrKorespondensiTanpaDatabase - pembaca korespondensi / akun tidak terpasang. 503.
	ErrKorespondensiTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, korespondensi tidak terbaca")
	// ErrSMTPBelumDikonfigurasi - NBFACIN_SMTP_HOST / _FROM kosong: tidak ada yang dikirim. 503.
	ErrSMTPBelumDikonfigurasi = errors.New("services: SMTP belum dikonfigurasi")
	// ErrKirimSurelGagal - server SMTP menolak / tak terjangkau; isi TETAP tersimpan, status tetap draf. 502.
	ErrKirimSurelGagal = errors.New("services: email gagal dikirim")
	// ErrSMTPTanpaTLS - server SMTP tidak mendukung TLS dan pengiriman tanpa TLS tidak diizinkan: TIDAK ada yang dikirim;
	// isi TETAP tersimpan, status tetap draf. 502.
	ErrSMTPTanpaTLS = errors.New("services: server SMTP tidak mendukung TLS; email tidak dikirim")
	// ErrKorespondensiStatusGagal - email SUDAH terkirim tetapi status terkirim gagal ditulis: JANGAN kirim ulang. 500.
	ErrKorespondensiStatusGagal = errors.New("services: email terkirim tetapi status gagal disimpan")
)

const (
	// lebar* - kolom migrasi 207: SUBJECT 1000, RECIPIENT 150 (= M_LOGIN_GO.NAME), RECIPIENT_MAIL 254 (= M_LOGIN_GO.EMAIL).
	lebarSubjekKorespondensi = 1000
	lebarNamaPenerima        = 150
	lebarSurelPenerima       = 254
	// batasPesanKorespondensi - HTML bersih paling banyak 4000 bita (MESSAGE VARCHAR2(4000), migrasi 207; CLOB dilarang
	// penjaga inti); batasPenerimaKorespondensi - 100 penerima.
	batasPesanKorespondensi    = 4000
	batasPenerimaKorespondensi = 100
	// bentukTanggalKorespondensi - kabel `.Date`: DD-MM-YYYY HH:MM WIB.
	bentukTanggalKorespondensi = "02-01-2006 15:04"
)

// DenganKorespondensi memasang penyimpan korespondensi dan pembaca akun login (tiket 56).
func (s *Service) DenganKorespondensi(p repository.PenyimpanKorespondensi, a repository.PembacaAkunLogin) *Service {
	s.korespondensi, s.akunLogin = p, a
	return s
}

// DenganJabatan memasang pembaca JOB_POSITION akun (akses tab IT); DenganLossRecord juga memasangnya.
func (s *Service) DenganJabatan(j repository.PembacaJabatan) *Service { s.jabatan = j; return s }

// WajibAksesKorespondensi - tanpa identitas 401, bukan "IT" 403 (lihat wajibJabatanIT). Dipanggil handler SEBELUM badan
// diurai dan oleh setiap metode tab SEBELUM status baris dibaca.
func (s *Service) WajibAksesKorespondensi(ctx context.Context, pelaku inti.Pelaku) error {
	return s.wajibJabatanIT(ctx, pelaku, "Correspondence")
}

// DenganPengirimSurel memasang pengirim SMTP; nil = SMTP belum dikonfigurasi (Send -> 503).
func (s *Service) DenganPengirimSurel(p repository.PengirimSurel) *Service { s.surel = p; return s }

// DenganSMTP - konfigurasi NBFACIN_SMTP_* (inti/backend/config): host atau from kosong = tidak dipasang (Send 503).
func (s *Service) DenganSMTP(k repository.KonfigSMTP) *Service {
	if k.Host == "" || k.Dari == "" {
		s.surel = nil
		return s
	}
	return s.DenganPengirimSurel(repository.NewPengirimSMTP(k))
}

// TampilanKorespondensi - GET daftar.
type TampilanKorespondensi struct {
	Daftar []models.Korespondensi
	Pesan  []string
}

// waktuKorespondensiKeKabel - teks Pega DateTime -> DD-MM-YYYY HH:MM WIB; teks lain apa adanya.
func waktuKorespondensiKeKabel(s string) string {
	t, err := time.Parse(BentukWaktuPega, s)
	if err != nil {
		return s
	}
	return t.In(WIB).Format(bentukTanggalKorespondensi)
}

// keKabelKorespondensi - tanggal ke kabel; Message dibersihkan ULANG saat dibaca (baris yang ditulis di luar layanan ini -
// mis. loader / SQL langsung - tidak sampai ke layar mentah). Pembersih idempoten atas keluarannya sendiri.
func keKabelKorespondensi(k models.Korespondensi) models.Korespondensi {
	k.Message = BersihkanHTMLKorespondensi(k.Message)
	k.Tanggal, k.DikirimPada = waktuKorespondensiKeKabel(k.Tanggal), waktuKorespondensiKeKabel(k.DikirimPada)
	return k
}

// noKorespondensiSah - nomor = ID baris (angka ≤ 19 digit).
func noKorespondensiSah(no string) bool {
	if no == "" || len(no) > 19 {
		return false
	}
	for _, c := range no {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func (s *Service) siapKorespondensi(id string) error {
	if s.korespondensi == nil || s.akunLogin == nil {
		return ErrKorespondensiTanpaDatabase
	}
	if !idKasusSah(id) {
		return ErrKasusTidakAda
	}
	return nil
}

// galatKorespondensi - galat repository -> galat layanan (404 / 409 / 404 kasus).
func galatKorespondensi(err error) error {
	switch {
	case errors.Is(err, repository.ErrKorespondensiTidakAda):
		return ErrKorespondensiTidakAda
	case errors.Is(err, repository.ErrKorespondensiTerkirim):
		return ErrKorespondensiTerkirim
	}
	return galatKasus(err)
}

// DaftarKorespondensi - GET …/korespondensi: urut Add.
func (s *Service) DaftarKorespondensi(ctx context.Context, pelaku inti.Pelaku, id string) (TampilanKorespondensi, error) {
	if err := s.WajibAksesKorespondensi(ctx, pelaku); err != nil {
		return TampilanKorespondensi{}, err
	}
	return s.daftarKorespondensi(ctx, id)
}

// daftarKorespondensi - daftar tanpa pemeriksaan akses (pemanggil sudah memeriksa).
func (s *Service) daftarKorespondensi(ctx context.Context, id string) (TampilanKorespondensi, error) {
	if err := s.siapKorespondensi(id); err != nil {
		return TampilanKorespondensi{}, err
	}
	d, err := s.korespondensi.DaftarKorespondensi(ctx, id)
	if err != nil {
		return TampilanKorespondensi{}, galatKorespondensi(err)
	}
	t := TampilanKorespondensi{Daftar: make([]models.Korespondensi, 0, len(d)), Pesan: []string{}}
	for _, k := range d {
		t.Daftar = append(t.Daftar, keKabelKorespondensi(k))
	}
	return t, nil
}

// satuKorespondensi - satu baris daftar (sesudah tulis).
func (s *Service) satuKorespondensi(ctx context.Context, id, no string) (models.Korespondensi, error) {
	t, err := s.daftarKorespondensi(ctx, id)
	if err != nil {
		return models.Korespondensi{}, err
	}
	for _, k := range t.Daftar {
		if k.No == no {
			return k, nil
		}
	}
	return models.Korespondensi{}, ErrKorespondensiTidakAda
}

// TambahKorespondensi - POST …/korespondensi (CorrespondenceAdd): draf kosong, tanggal = sekarang, pengirim = akun login
// (NAME / EMAIL M_LOGIN_GO; akun tak ada -> kosong).
func (s *Service) TambahKorespondensi(ctx context.Context, pelaku inti.Pelaku, id string) (models.Korespondensi, error) {
	if err := s.WajibAksesKorespondensi(ctx, pelaku); err != nil {
		return models.Korespondensi{}, err
	}
	if err := s.wajibTahapUbah(ctx, id); err != nil {
		return models.Korespondensi{}, err
	}
	if err := s.siapKorespondensi(id); err != nil {
		return models.Korespondensi{}, err
	}
	if s.transaksi == nil {
		return models.Korespondensi{}, ErrKorespondensiTanpaDatabase
	}
	nama, surel, _, err := s.akunLogin.AkunLogin(ctx, pelaku.AkunID)
	if err != nil {
		return models.Korespondensi{}, err
	}
	k := models.Korespondensi{Tanggal: s.sekarang().UTC().Format(BentukWaktuPega), Sender: nama, SenderMail: surel}
	var no string
	err = s.transaksi(ctx, func(tx *db.Tx) error {
		var err error
		no, err = s.korespondensi.TambahKorespondensi(ctx, tx, id, k, pelaku.AkunID)
		return err
	})
	if err != nil {
		return models.Korespondensi{}, galatKorespondensi(err)
	}
	return s.satuKorespondensi(ctx, id, no)
}

// rapikanIsiKorespondensi - subject dipangkas (tanpa CR / LF), message dibersihkan, penerima dipangkas (baris kosong
// dibuang); batas lebar / jumlah -> 400.
func rapikanIsiKorespondensi(isi models.IsiKorespondensi) (models.IsiKorespondensi, error) {
	var masalah []string
	out := models.IsiKorespondensi{Subject: strings.TrimSpace(isi.Subject), Message: BersihkanHTMLKorespondensi(isi.Message),
		Penerima: []models.PenerimaKorespondensi{}}
	if adaCRLFKorespondensi(out.Subject) {
		masalah = append(masalah, "subject tidak boleh memuat baris baru")
	}
	if len(out.Subject) > lebarSubjekKorespondensi {
		masalah = append(masalah, fmt.Sprintf("subject paling banyak %d bita", lebarSubjekKorespondensi))
	}
	if len(out.Message) > batasPesanKorespondensi {
		masalah = append(masalah, fmt.Sprintf("message paling banyak %d bita", batasPesanKorespondensi))
	}
	for i, p := range isi.Penerima {
		p.Recipient, p.RecipientMail = strings.TrimSpace(p.Recipient), strings.TrimSpace(p.RecipientMail)
		if p.Recipient == "" && p.RecipientMail == "" {
			continue
		}
		if len(p.Recipient) > lebarNamaPenerima || len(p.RecipientMail) > lebarSurelPenerima {
			masalah = append(masalah, fmt.Sprintf("penerima[%d] paling banyak %d / %d bita", i, lebarNamaPenerima, lebarSurelPenerima))
		}
		if adaCRLFKorespondensi(p.Recipient) || adaCRLFKorespondensi(p.RecipientMail) {
			masalah = append(masalah, fmt.Sprintf("penerima[%d] tidak boleh memuat baris baru", i))
		}
		out.Penerima = append(out.Penerima, p)
	}
	if len(out.Penerima) > batasPenerimaKorespondensi {
		masalah = append(masalah, fmt.Sprintf("penerima paling banyak %d", batasPenerimaKorespondensi))
	}
	if !utf8.ValidString(out.Subject) {
		masalah = append(masalah, "subject bukan UTF-8 sah")
	}
	if len(masalah) > 0 {
		return out, fmt.Errorf("%w: %s", ErrMasukanKorespondensi, strings.Join(masalah, "; "))
	}
	return out, nil
}

func adaCRLFKorespondensi(s string) bool { return strings.ContainsAny(s, "\r\n") }

// simpanIsiKorespondensi - identitas, rapikan, tulis draf (409 terkirim / 404 bukan milik kasus). Isian tidak sah atas
// korespondensi yang sudah terkirim / bukan milik kasus -> 409 / 404 lebih dulu dari 400 (statusnya yang menentukan).
func (s *Service) simpanIsiKorespondensi(ctx context.Context, pelaku inti.Pelaku, id, no string, isi models.IsiKorespondensi) (models.IsiKorespondensi, error) {
	if err := s.WajibAksesKorespondensi(ctx, pelaku); err != nil {
		return isi, err
	}
	if err := s.wajibTahapUbah(ctx, id); err != nil {
		return isi, err
	}
	if err := s.siapKorespondensi(id); err != nil {
		return isi, err
	}
	if !noKorespondensiSah(no) {
		return isi, ErrKorespondensiTidakAda
	}
	bersih, err := rapikanIsiKorespondensi(isi)
	if err != nil {
		if k, errBaca := s.satuKorespondensi(ctx, id, no); errBaca != nil {
			return isi, errBaca
		} else if k.Terkirim {
			return isi, ErrKorespondensiTerkirim
		}
		return isi, err
	}
	if s.transaksi == nil {
		return isi, ErrKorespondensiTanpaDatabase
	}
	if err := s.transaksi(ctx, func(tx *db.Tx) error { return s.korespondensi.UbahKorespondensi(ctx, tx, id, no, bersih) }); err != nil {
		return isi, galatKorespondensi(err)
	}
	return bersih, nil
}

// SimpanKorespondensi - PUT …/korespondensi/{no}.
func (s *Service) SimpanKorespondensi(ctx context.Context, pelaku inti.Pelaku, id, no string, isi models.IsiKorespondensi) (models.Korespondensi, error) {
	if _, err := s.simpanIsiKorespondensi(ctx, pelaku, id, no, isi); err != nil {
		return models.Korespondensi{}, err
	}
	return s.satuKorespondensi(ctx, id, no)
}

// HapusKorespondensi - DELETE …/korespondensi/{no}: draf saja.
func (s *Service) HapusKorespondensi(ctx context.Context, pelaku inti.Pelaku, id, no string) error {
	if err := s.WajibAksesKorespondensi(ctx, pelaku); err != nil {
		return err
	}
	if err := s.wajibTahapUbah(ctx, id); err != nil {
		return err
	}
	if err := s.siapKorespondensi(id); err != nil {
		return err
	}
	if !noKorespondensiSah(no) {
		return ErrKorespondensiTidakAda
	}
	if s.transaksi == nil {
		return ErrKorespondensiTanpaDatabase
	}
	return galatKorespondensi(s.transaksi(ctx, func(tx *db.Tx) error { return s.korespondensi.HapusKorespondensi(ctx, tx, id, no) }))
}

// alamatSurelSah - satu alamat polos (tanpa nama tampilan), net/mail mengurainya kembali sama persis.
func alamatSurelSah(a string) bool {
	if a == "" || adaCRLFKorespondensi(a) {
		return false
	}
	p, err := mail.ParseAddress(a)
	return err == nil && p.Name == "" && p.Address == a
}

// susunSurelKorespondensi - To = recipientMail unik (tanpa beda huruf), Cc / Reply-To = senderMail sah bila bukan
// salah satu To; alamat tidak sah / tanpa penerima ber-email -> 400 per alamat.
func susunSurelKorespondensi(k models.Korespondensi, bersih models.IsiKorespondensi) (models.Surel, error) {
	sur := models.Surel{Subject: bersih.Subject, HTML: bersih.Message}
	var masalah []string
	lihat := map[string]bool{}
	for i, p := range bersih.Penerima {
		switch {
		case p.RecipientMail == "":
			continue
		case !alamatSurelSah(p.RecipientMail):
			masalah = append(masalah, fmt.Sprintf("penerima[%d] alamat email %q tidak sah", i, p.RecipientMail))
		case !lihat[strings.ToLower(p.RecipientMail)]:
			lihat[strings.ToLower(p.RecipientMail)] = true
			sur.To = append(sur.To, p.RecipientMail)
		}
	}
	if len(masalah) == 0 && len(sur.To) == 0 {
		masalah = append(masalah, "tidak ada penerima ber-email")
	}
	if len(masalah) > 0 {
		return sur, fmt.Errorf("%w: %s", ErrMasukanKorespondensi, strings.Join(masalah, "; "))
	}
	if pengirim := strings.TrimSpace(k.SenderMail); alamatSurelSah(pengirim) {
		sur.ReplyTo = pengirim
		if !lihat[strings.ToLower(pengirim)] {
			sur.Cc = []string{pengirim}
		}
	}
	return sur, nil
}

// KirimKorespondensi - POST …/korespondensi/{no}/kirim: simpan dulu (seperti PUT), lalu kirim. 400 penerima, 409 terkirim,
// 503 SMTP belum dikonfigurasi, 502 SMTP gagal - pada 400 / 503 / 502 isi TETAP tersimpan dan status tetap draf.
// Jejak: satu baris log berhasil / gagal (kasus, nomor, jumlah penerima, galat) - tanpa isi pesan, tanpa kredensial.
func (s *Service) KirimKorespondensi(ctx context.Context, pelaku inti.Pelaku, id, no string, isi models.IsiKorespondensi) (models.Korespondensi, error) {
	bersih, err := s.simpanIsiKorespondensi(ctx, pelaku, id, no, isi)
	if err != nil {
		return models.Korespondensi{}, err
	}
	k, err := s.satuKorespondensi(ctx, id, no)
	if err != nil {
		return models.Korespondensi{}, err
	}
	sur, err := susunSurelKorespondensi(k, bersih)
	if err != nil {
		return models.Korespondensi{}, err
	}
	if s.surel == nil {
		return models.Korespondensi{}, ErrSMTPBelumDikonfigurasi
	}
	if err := s.surel.KirimSurel(ctx, sur); err != nil {
		log.Printf("nbfacin: korespondensi %s kasus %s GAGAL dikirim ke %d alamat oleh %s: %v", no, id, len(sur.To)+len(sur.Cc), pelaku.AkunID, err)
		if errors.Is(err, repository.ErrSMTPTanpaTLS) {
			return models.Korespondensi{}, ErrSMTPTanpaTLS
		}
		return models.Korespondensi{}, fmt.Errorf("%w: %v", ErrKirimSurelGagal, err)
	}
	log.Printf("nbfacin: korespondensi %s kasus %s terkirim ke %d alamat oleh %s", no, id, len(sur.To)+len(sur.Cc), pelaku.AkunID)
	waktu := s.sekarang().UTC().Format(BentukWaktuPega)
	// Email SUDAH terkirim: status ditulis walau permintaan dibatalkan (klien menutup koneksi) supaya tidak tampil draf.
	tetap := context.WithoutCancel(ctx)
	if err := s.transaksi(tetap, func(tx *db.Tx) error { return s.korespondensi.TandaiTerkirim(tetap, tx, id, no, waktu) }); err != nil {
		log.Printf("nbfacin: korespondensi %s kasus %s terkirim tetapi status gagal disimpan: %v", no, id, err)
		return models.Korespondensi{}, fmt.Errorf("%w: %v", ErrKorespondensiStatusGagal, err)
	}
	return s.satuKorespondensi(ctx, id, no)
}

// SaranPenerimaKorespondensi - GET /api/nbfacin/korespondensi/penerima/saran?q=.
func (s *Service) SaranPenerimaKorespondensi(ctx context.Context, pelaku inti.Pelaku, kata string) ([]models.SaranAkumulasi, error) {
	if err := s.WajibAksesKorespondensi(ctx, pelaku); err != nil {
		return nil, err
	}
	if err := periksaParamKueri(ErrMasukanKorespondensi, paramAkumulasi{"q", &kata}); err != nil {
		return nil, err
	}
	if s.akunLogin == nil {
		return nil, ErrKorespondensiTanpaDatabase
	}
	return s.akunLogin.SaranPenerima(ctx, kata)
}
