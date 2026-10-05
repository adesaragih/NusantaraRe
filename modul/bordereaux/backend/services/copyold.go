package services

// Copy Old Data (keputusan work owner 04-10-2026: "buat copy old data", tombolnya hanya untuk superadmin). Popup berisi
// berkas Pega yang isinya masih tertinggal di JSON `M_BORDEREAUX.DATA_JSON`; yang dicentang disalin lewat `Process
// Copy`, SATU TRANSAKSI PER BERKAS - satu yang gagal tidak membatalkan yang lain. JSON hanya DIBACA.
//
// Yang disalin per berkas:
//   - header BORDEREAUX - hanya bila barisnya belum ada (TANGGAL = pxCreateDateTime, USER_INPUT = pxCreateOperator);
//   - detail - page list kombinasi berkas (`models.KombinasiBdx.DaftarJSON`), hanya bila tabel detailnya masih nol
//     baris untuk berkas itu (DEV 04-10-2026: detail tidak pernah setengah tersalin - nol atau sama jumlahnya), jadi
//     isi yang sudah disimpan aplikasi baru tidak pernah ditimpa;
//   - riwayat - `CommentList` ke BORDEREAUX_HISTORY, hanya bila riwayat berkas itu masih kosong.
//
// Kolom yang namanya berganti sejak JSON ditulis dibaca dari nama lamanya (aliasLama, terbukti nilai-sama di DEV).
//
// Bentuk nilai JSON Pega (DEV 04-10-2026, pola saja): tanggal DateTime `yyyyMMddTHHmmss.SSS GMT` - tanggal Pega
// tersimpan sebagai tengah malam WIB (17:00 GMT hari sebelumnya), maka dibaca di WIB; satu header ber-`yyyyMMdd`.
// Angka kadang berformat ribuan Indonesia (`1.234.567`), kadang desimal titik (`1234.5`).

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"nusantarare/modul/bordereaux/backend/models"
	"nusantarare/modul/bordereaux/backend/repository"
)

// MaksSalinLama - ID paling banyak dalam satu permintaan Process Copy (layar mengirim per UKURAN_SALIN_LAMA).
const MaksSalinLama = 200

// Batas kolom tujuan (byte - semantik VARCHAR2 bawaan).
const (
	batasPIC       = 64
	batasUserInput = 50
)

// picKosong - PIC riwayat lama tanpa OperatorName (kolom PIC NOT NULL).
const picKosong = "-"

// wajibSuperadmin - Copy Old Data: superadmin dengan menu Bordereaux ber-hak PENUH (View only juga berlaku bagi
// superadmin, keputusan work owner 05-10-2026).
func wajibSuperadmin(a Aktor) error {
	if a.AkunID == "" || !a.Superadmin {
		return dilarang("Copy Old Data is only for super admin")
	}
	if !a.Penuh {
		return dilarang("Your access to the Bordereaux menu is View only")
	}
	return nil
}

// berkasJSON - isi JSON Pega satu berkas yang dipakai Copy Old Data.
type berkasJSON struct {
	header   models.Header
	dibuat   string // pxCreateDateTime, `DD-MM-YYYY HH24:MI:SS` WIB
	daftar   map[string][]map[string]any
	komentar []komentarJSON
}

type komentarJSON struct {
	Tanggal    string
	PIC        string
	IsApproved bool
	Komentar   string
}

// teksJSON - nilai skalar JSON menjadi teks (angka apa adanya, null = "").
func teksJSON(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case json.Number:
		return x.String()
	case bool:
		if x {
			return "true"
		}
		return "false"
	}
	return ""
}

// bacaJSON mengurai DATA_JSON; page list dibaca sebagai baris peta properti.
func bacaJSON(isi string) (berkasJSON, error) {
	d := json.NewDecoder(strings.NewReader(isi))
	d.UseNumber()
	var akar map[string]any
	if err := d.Decode(&akar); err != nil {
		return berkasJSON{}, fmt.Errorf("the old JSON cannot be read: %v", err)
	}
	s := func(k string) string { return strings.TrimSpace(teksJSON(akar[k])) }
	b := berkasJSON{daftar: map[string][]map[string]any{}}
	b.header = models.Header{BdxID: s("BDX_ID"), UserInput: s("pxCreateOperator"), Type: s("TYPE"), TypeBusiness: s("TYPE_BUSINESS"),
		MasterID: s("MASTERID"), CedingID: s("CEDINGID"), CedingName: s("CEDINGNAME"), SobID: s("SOBID"), SobName: s("SOBNAME"),
		TreatyName: s("TREATYNAME"), ReffNoSOA: s("REFFNO_OF_SOA"), ReffNoBDX: s("REFFNO_OF_BDX"), Position: s("POSITION"),
		StatusAksep: s("STATUSAKSEP")}
	b.header.ReportStart, _ = TanggalLama(s("BDXREPORT_START"))
	b.header.ReportEnd, _ = TanggalLama(s("BDXREPORT_END"))
	b.dibuat = WaktuLama(s("pxCreateDateTime"))
	for k, v := range akar {
		larik, ok := v.([]any)
		if !ok {
			continue
		}
		var baris []map[string]any
		for _, x := range larik {
			if m, ok := x.(map[string]any); ok {
				baris = append(baris, m)
			}
		}
		b.daftar[k] = baris
	}
	for _, m := range b.daftar["CommentList"] {
		tgl := WaktuLama(strings.TrimSpace(teksJSON(m["Date"])))
		b.komentar = append(b.komentar, komentarJSON{Tanggal: tgl, PIC: strings.TrimSpace(teksJSON(m["OperatorName"])),
			IsApproved: strings.TrimSpace(teksJSON(m["IsApproved"])) == "1", Komentar: strings.TrimSpace(teksJSON(m["Suggest"]))})
	}
	return b, nil
}

// waktuPega - DateTime Pega `yyyyMMddTHHmmss.SSS GMT` (detik pecahan boleh tiada) dalam WIB.
func waktuPega(s string) (time.Time, bool) {
	s = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(s), "GMT"))
	for _, pola := range []string{"20060102T150405.000", "20060102T150405"} {
		if t, err := time.ParseInLocation(pola, s, time.UTC); err == nil {
			return t.In(lokasiJakarta), true
		}
	}
	return time.Time{}, false
}

// tahunKosong - tanggal "kosong" Pega (epoch) dibaca sebagai kosong.
func tahunKosong(t time.Time) bool { return t.Year() <= 1970 }

// TanggalLama - nilai tanggal JSON Pega menjadi `DD-MM-YYYY`: DateTime GMT dibaca di WIB, `yyyyMMdd` apa adanya,
// `dd/MM/yyyy` seperti CSV; kosong atau epoch = "". ok false = bukan tanggal.
func TanggalLama(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", true
	}
	if t, ok := waktuPega(s); ok {
		if tahunKosong(t) {
			return "", true
		}
		return t.Format("02-01-2006"), true
	}
	if t, err := time.Parse("20060102", s); err == nil {
		if tahunKosong(t) {
			return "", true
		}
		return t.Format("02-01-2006"), true
	}
	return TanggalCSV(s)
}

// WaktuLama - DateTime Pega menjadi `DD-MM-YYYY HH24:MI:SS` WIB; `yyyyMMdd` = tengah malam; selain itu "".
func WaktuLama(s string) string {
	s = strings.TrimSpace(s)
	if t, ok := waktuPega(s); ok && !tahunKosong(t) {
		return t.Format("02-01-2006 15:04:05")
	}
	if t, err := time.Parse("20060102", s); err == nil && !tahunKosong(t) {
		return t.Format("02-01-2006 15:04:05")
	}
	return ""
}

// AngkaLama - nilai angka JSON Pega menjadi teks desimal kanonik. Berkoma atau bertitik lebih dari satu = format
// Indonesia (`1.234.567,89`); selain itu desimal titik (`1234.5`). `%` dan spasi dibuang; kosong atau `-` = "".
func AngkaLama(s string) (string, bool) {
	s = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s, "%", ""), " ", ""))
	if s == "" || s == "-" {
		return "", true
	}
	if strings.Contains(s, ",") || strings.Count(s, ".") > 1 {
		return AngkaIndonesia(s)
	}
	if !polaAngka.MatchString(s) {
		return "", false
	}
	return kanonikAngka(s), true
}

func normalKunci(s string) string { return strings.ToUpper(strings.ReplaceAll(s, "_", "")) }

// aliasLama - nama properti Pega LAMA di JSON untuk kolom yang namanya sudah berganti. Terbukti DEV 04-10-2026 (agregat
// saja): pada berkas yang detailnya ada di JSON dan di tabel, nilai kunci ini SAMA PERSIS dengan kolom tabelnya di
// seluruh baris (PremiumEngineering 21/21, ClaimFire 22/22, ClaimEngineering 21/21). Hanya dipakai bila kunci kolomnya
// sendiri tidak ada di baris JSON; urutan = urutan coba.
var aliasLama = map[string][]string{
	"PREMIUM_100":                  {"PREMIUM"},
	"PREMIUM_REINSURER":            {"PREMIUM_CEDED100"},
	"LOL_PML_EML":                  {"AMOUNT"},
	"LOL_PML_EML_PCT":              {"LOSESTIMATION"},
	"EFFECTIFDATE_EDM":             {"EDO_ENDORSEMENT"},
	"PAID_CLAIMS_REINSURER":        {"CLAIM_CEDED100"},
	"OUTSTANDING_CLAIMS_REINSURER": {"CEDED_OUTSTANDING_CLAIMS100"},
	"PAID_CLAIMS_RNM":              {"PAID_CLAIMS_RNMSHARE", "PAID_CLAIMS"},
	"OUTSTANDING_CLAIMS_RNM":       {"OUTSTANDING_CLAIMS"},
	"ZIP_CODE":                     {"ZIPCODEE"},
}

// nilaiKolom - nilai satu kolom tabel dari baris JSON: nama persis, lalu nama tanpa garis bawah (properti Engineering
// Pega beda garis bawah dari kolomnya), lalu nama lama (aliasLama).
func nilaiKolom(baris map[string]any, normal map[string]string, kolom string) string {
	if v, ok := baris[kolom]; ok {
		return teksJSON(v)
	}
	if v, ok := normal[normalKunci(kolom)]; ok {
		return v
	}
	for _, a := range aliasLama[kolom] {
		if v, ok := baris[a]; ok {
			return teksJSON(v)
		}
	}
	return ""
}

// BarisLama - page list JSON satu kombinasi menjadi baris kanonik tabel detail; seluruh kesalahan dikumpulkan.
func BarisLama(k models.KombinasiBdx, daftar []map[string]any) ([]models.Baris, error) {
	var p pengumpul
	hasil := make([]models.Baris, 0, len(daftar))
	for i, m := range daftar {
		normal := map[string]string{}
		for kk, v := range m {
			normal[normalKunci(kk)] = teksJSON(v)
		}
		b := models.Baris{}
		for _, c := range k.Kolom {
			mentah := strings.TrimSpace(nilaiKolom(m, normal, c.Kolom))
			v, ok := mentah, true
			switch c.Jenis {
			case models.Angka:
				v, ok = AngkaLama(mentah)
			case models.Tanggal:
				v, ok = TanggalLama(mentah)
			}
			masalah := ""
			if !ok {
				masalah = map[models.Jenis]string{models.Angka: "is not a number", models.Tanggal: "is not a date"}[c.Jenis]
			} else {
				masalah = periksaNilai(c, v)
			}
			if masalah != "" {
				p.tambah(fmt.Sprintf("%s row %d, column %s: %q %s", k.DaftarJSON, i+1, c.Judul, mentah, masalah))
				continue
			}
			b[c.Kolom] = v
		}
		hasil = append(hasil, b)
	}
	if err := p.galat(); err != nil {
		return nil, err
	}
	return hasil, nil
}

// potongByte memotong teks ke n byte tanpa memotong satu karakter.
func potongByte(s string, n int) string {
	if len(s) <= n {
		return s
	}
	s = s[:n]
	for !utf8.ValidString(s) {
		s = s[:len(s)-1]
	}
	return s
}

// kombinasiLama - kombinasi berkas: dari header BORDEREAUX bila ada, selain itu dari JSON.
func kombinasiLama(j models.JSONLama, b berkasJSON) (models.KombinasiBdx, bool) {
	tipe, bisnis := b.header.Type, b.header.TypeBusiness
	if j.AdaHeader {
		tipe, bisnis = j.Type, j.Business
	}
	k, err := models.CariKombinasi(strings.ToUpper(tipe), strings.ToUpper(bisnis))
	return k, err == nil
}

// DaftarLama - isi popup Copy Old Data: berkas JSON yang header, detail, atau riwayatnya belum ada di tabel. JSON
// yang tidak terbaca tetap tampil (Process Copy melaporkan alasannya).
func (l *Layanan) DaftarLama(ctx context.Context, a Aktor) ([]models.BerkasLama, error) {
	if err := wajibSuperadmin(a); err != nil {
		return nil, err
	}
	js, err := l.gudang.JSONLama(ctx, nil, "")
	if err != nil {
		return nil, err
	}
	riwayat, err := l.gudang.CacahRiwayat(ctx)
	if err != nil {
		return nil, err
	}
	cacahTabel := map[string]map[string]int{}
	hasil := []models.BerkasLama{}
	for _, j := range js {
		bl := models.BerkasLama{ID: j.ID, Riwayat: riwayat[j.ID], TanpaHeader: !j.AdaHeader, Type: j.Type, Business: j.Business}
		b, err := bacaJSON(j.Isi)
		if err == nil {
			h := b.header
			if !j.AdaHeader {
				bl.Type, bl.Business = h.Type, h.TypeBusiness
			}
			bl.Ceding, bl.Treaty, bl.Status, bl.Komentar = h.CedingName, h.TreatyName, h.StatusAksep, len(b.komentar)
			if k, ok := kombinasiLama(j, b); ok {
				bl.BarisJSON = len(b.daftar[k.DaftarJSON])
				if _, ada := cacahTabel[k.Tabel]; !ada {
					if cacahTabel[k.Tabel], err = l.gudang.CacahDetail(ctx, k); err != nil {
						return nil, err
					}
				}
				bl.BarisTabel = cacahTabel[k.Tabel][j.ID]
			}
		}
		perlu := err != nil || bl.TanpaHeader || (bl.BarisJSON > 0 && bl.BarisTabel == 0) || (bl.Komentar > 0 && bl.Riwayat == 0)
		if perlu {
			hasil = append(hasil, bl)
		}
	}
	return hasil, nil
}

// galatTolak - alasan berkas ditolak (dilaporkan, transaksi berkas itu dibatalkan).
type galatTolak struct{ pesan []string }

func (g galatTolak) Error() string { return strings.Join(g.pesan, "; ") }

// SalinLama - Process Copy: ID yang dicentang disalin satu per satu. Galat basis data satu berkas dicatat log dan
// dilaporkan `gagal` untuk berkas itu saja.
func (l *Layanan) SalinLama(ctx context.Context, a Aktor, ids []string) (models.JawabanSalinLama, error) {
	j := models.JawabanSalinLama{Hasil: []models.HasilSalinLama{}}
	if err := wajibSuperadmin(a); err != nil {
		return j, err
	}
	var unik []string
	dilihat := map[string]bool{}
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" || dilihat[id] {
			continue
		}
		dilihat[id] = true
		unik = append(unik, id)
	}
	switch {
	case len(unik) == 0:
		return j, tolak("Select at least one old bordereaux to copy")
	case len(unik) > MaksSalinLama:
		return j, tolak("At most %d old bordereaux can be copied at once", MaksSalinLama)
	}
	for _, id := range unik {
		h := models.HasilSalinLama{ID: id, Pesan: []string{}}
		var (
			ditulis bool
			catatan []string
		)
		err := l.tx(ctx, func(tx *dbTx) error {
			var err error
			ditulis, catatan, err = l.salinSatu(ctx, tx, id)
			return err
		})
		var tolakan galatTolak
		switch {
		case errors.As(err, &tolakan):
			h.Status, h.Pesan = models.SalinDitolak, tolakan.pesan
		case errors.Is(err, repository.ErrBelumDimigrasi):
			return j, err
		case err != nil:
			log.Printf("bordereaux: copy old %s: %v", id, err)
			h.Status, h.Pesan = models.SalinGagal, []string{"database error; details in the server log"}
		case ditulis:
			h.Status, h.Pesan = models.SalinDisalin, append(h.Pesan, catatan...)
			j.Disalin++
		default:
			h.Status, h.Pesan = models.SalinSudahAda, []string{"header, detail, and history are already in the tables"}
		}
		j.Hasil = append(j.Hasil, h)
	}
	return j, nil
}

// salinSatu menyalin satu berkas di dalam transaksinya: seluruh isi dibaca dan dinilai DULU, baru ditulis - berkas
// yang ditolak tidak meninggalkan apa pun. ditulis false = tidak ada yang perlu disalin.
func (l *Layanan) salinSatu(ctx context.Context, tx *dbTx, id string) (ditulis bool, catatan []string, err error) {
	js, err := l.gudang.JSONLama(ctx, tx, id)
	if err != nil {
		return false, nil, err
	}
	if len(js) == 0 {
		return false, nil, galatTolak{[]string{"not an old bordereaux document"}}
	}
	j := js[0]
	b, err := bacaJSON(j.Isi)
	if err != nil {
		return false, nil, galatTolak{[]string{err.Error()}}
	}
	var (
		k     models.KombinasiBdx
		baris []models.Baris
	)
	if kk, ok := kombinasiLama(j, b); ok && len(b.daftar[kk.DaftarJSON]) > 0 {
		ada, err := l.gudang.Detail(ctx, tx, kk, id)
		if err != nil {
			return false, nil, err
		}
		if len(ada) == 0 {
			k = kk
			if baris, err = BarisLama(k, b.daftar[k.DaftarJSON]); err != nil {
				var g GalatCSV
				if errors.As(err, &g) {
					return false, nil, galatTolak{g.Pesan}
				}
				return false, nil, err
			}
		}
	}
	var komentar []komentarJSON
	if len(b.komentar) > 0 {
		riw, err := l.gudang.Riwayat(ctx, id)
		if err != nil {
			return false, nil, err
		}
		if len(riw) == 0 {
			komentar = b.komentar
		}
	}
	if !j.AdaHeader {
		h := b.header
		h.BdxID, h.Tanggal = id, b.dibuat
		h.UserInput = potongByte(h.UserInput, batasUserInput)
		if err := l.gudang.SisipHeader(ctx, tx, h); err != nil {
			return false, nil, err
		}
		ditulis = true
		catatan = append(catatan, "header created from the old JSON")
	}
	if len(baris) > 0 {
		if err := l.sisipBaris(ctx, tx, k, id, baris); err != nil {
			return false, nil, err
		}
		ditulis = true
		catatan = append(catatan, fmt.Sprintf("%d detail rows", len(baris)))
	}
	for _, c := range komentar {
		pic := potongByte(c.PIC, batasPIC)
		if pic == "" {
			pic = picKosong
		}
		if err := l.gudang.SisipRiwayatLama(ctx, tx, id, c.Tanggal, pic, c.IsApproved, potongByte(c.Komentar, BatasKomentar)); err != nil {
			return false, nil, err
		}
	}
	if len(komentar) > 0 {
		ditulis = true
		catatan = append(catatan, fmt.Sprintf("%d history rows", len(komentar)))
	}
	return ditulis, catatan, nil
}

// sisipBaris menyisipkan baris detail ber-ID `BDX_DTL-...`; ID yang bentrok digeser satu milidetik (seperti Save).
func (l *Layanan) sisipBaris(ctx context.Context, tx *dbTx, k models.KombinasiBdx, id string, baris []models.Baris) error {
	jam, geser := l.sekarang(), 0
	for _, b := range baris {
		for {
			if geser >= cobaIDMaksimum+len(baris) {
				return fmt.Errorf("services: tidak menemukan ID detail kosong")
			}
			idDetail := IDDetail(jam.Add(time.Duration(geser) * time.Millisecond))
			geser++
			err := l.gudang.SisipDetail(ctx, tx, k, id, idDetail, b)
			if errors.Is(err, repository.ErrIDTerpakai) {
				continue
			}
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}
