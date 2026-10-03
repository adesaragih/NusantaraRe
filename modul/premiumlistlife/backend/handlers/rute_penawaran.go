package handlers

// Pintu HTTP form penawaran - tiket 01 bagian 3 (layar Input Offer).
//
//	GET /api/polis-life/{id}/penawaran            isian + riwayat + pilihan
//	PUT /api/polis-life/{id}/penawaran            simpan (Save Offer)
//	GET /api/polis-life/cari-ceding?cari=         popup Choose Ceding Name
//	GET /api/polis-life/cari-pemegang-polis?cari= popup Policy Holder
//
// ⛔ Rute pencarian SATU segmen literal di bawah `/api/polis-life`, bukan
// `/api/polis-life/rujukan/{jenis}`: pola dua segmen itu bertabrakan dengan
// `/api/polis-life/{id}/peserta` (keduanya cocok dengan
// `/api/polis-life/rujukan/peserta`), dan ServeMux panik saat didaftarkan.
// Segmen literal selalu lebih spesifik dari `{id}` - pola `ringkas`.
//
// Nol aturan dagang di sini; seluruhnya di services/polis_isianpenawaran.go.

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/premiumlistlife/backend/models"
	"nusantarare/modul/premiumlistlife/backend/services"
)

// isiPenawaran adalah badan PUT - kode pilihan, bukan namanya.
//
// ⚠️ `dateReceived` bertanggal `YYYY-MM-DD` (masukan `type="date"`); kosong
// berarti tidak diisi - selnya tidak ber-pyRequired.
type isiPenawaran struct {
	CedingCo         string `json:"cedingCo"`
	CedingCoName     string `json:"cedingCoName"`
	PolicyHolder     string `json:"policyHolder"`
	PolicyHolderName string `json:"policyHolderName"`
	TypeCeding       string `json:"typeCeding"`
	BusinessCode     string `json:"businessCode"`
	DateReceived     string `json:"dateReceived"`
	Description      string `json:"description"`
	Status           string `json:"status"`
	// ⛔ Angka dikirim sebagai TEKS - `sumInsured` uang (ADR-U-0003), dan
	// teks kosong berbeda dari nol (ADR-U-0027).
	BatasUsiaPeserta       string `json:"batasUsiaPeserta"`
	PeriodePertanggungan   string `json:"periodePertanggungan"`
	SumInsured             string `json:"sumInsured"`
	TanggalPenawaran       string `json:"tanggalPenawaran"`
	TanggalRespon          string `json:"tanggalRespon"`
	TanggalKonfirmasi      string `json:"tanggalKonfirmasi"`
	TBC                    string `json:"tbc"`
	StatusUpdate           string `json:"statusUpdate"`
	KeteranganMarketing    string `json:"keteranganMarketing"`
	QQName                 string `json:"qqName"`
	JenisUsaha             string `json:"jenisUsaha"`
	KetentuanUnderwriting  string `json:"ketentuanUnderwriting"`
	TanggalKonfirmasiBalik string `json:"tanggalKonfirmasiBalik"`
	TanggalRealisasi       string `json:"tanggalRealisasi"`
	TanggalBind            string `json:"tanggalBind"`
	StatusFinal            string `json:"statusFinal"`
}

// keIsian mengurai badan menjadi isian model.
func (i isiPenawaran) keIsian() (models.IsianPenawaran, error) {
	hasil := models.IsianPenawaran{
		CedingCo: i.CedingCo, CedingCoName: i.CedingCoName,
		PolicyHolder: i.PolicyHolder, PolicyHolderName: i.PolicyHolderName,
		TypeCeding: i.TypeCeding, BusinessCode: i.BusinessCode,
		Description: i.Description, Status: i.Status,
		PeriodePertanggungan: i.PeriodePertanggungan,
		StatusUpdate:         i.StatusUpdate, KeteranganMarketing: i.KeteranganMarketing,
		QQName: i.QQName, JenisUsaha: i.JenisUsaha, KetentuanUnderwriting: i.KetentuanUnderwriting,
		StatusFinal: i.StatusFinal,
	}
	for _, t := range []struct {
		medan, nilai string
		tujuan       **time.Time
	}{
		{"dateReceived", i.DateReceived, &hasil.DateReceived},
		{"tanggalPenawaran", i.TanggalPenawaran, &hasil.TanggalPenawaran},
		{"tanggalRespon", i.TanggalRespon, &hasil.TanggalRespon},
		{"tanggalKonfirmasi", i.TanggalKonfirmasi, &hasil.TanggalKonfirmasi},
		{"tanggalKonfirmasiBalik", i.TanggalKonfirmasiBalik, &hasil.TanggalKonfirmasiBalik},
		{"tanggalRealisasi", i.TanggalRealisasi, &hasil.TanggalRealisasi},
		{"tanggalBind", i.TanggalBind, &hasil.TanggalBind},
	} {
		if v := strings.TrimSpace(t.nilai); v != "" {
			d, err := time.Parse("2006-01-02", v)
			if err != nil {
				return hasil, fmt.Errorf("%s harus berbentuk YYYY-MM-DD", t.medan)
			}
			*t.tujuan = &d
		}
	}
	for _, b := range []struct {
		medan, nilai string
		tujuan       **int
	}{
		{"batasUsiaPeserta", i.BatasUsiaPeserta, &hasil.BatasUsiaPeserta},
		{"tbc", i.TBC, &hasil.TBC},
	} {
		if v := strings.TrimSpace(b.nilai); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				return hasil, fmt.Errorf("%s harus bilangan bulat", b.medan)
			}
			*b.tujuan = &n
		}
	}
	if v := strings.ReplaceAll(strings.TrimSpace(i.SumInsured), ",", ""); v != "" {
		d, err := utils.ParseDecimal(v)
		if err != nil {
			return hasil, errors.New("sumInsured harus angka desimal")
		}
		hasil.SumInsured = d
	}
	return hasil, nil
}

// bacaPenawaranPolis melayani GET /api/polis-life/{id}/penawaran.
func bacaPenawaranPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		hasil, err := svc.FormPenawaran().Baca(r.Context(),
			inti.PelakuDari(r, stubPelaku), r.PathValue("id"))
		if jawabGalatPenawaran(w, err) {
			return
		}
		galat.TulisJSON(w, hasil)
	}
}

// simpanPenawaranPolis melayani PUT /api/polis-life/{id}/penawaran.
//
// ⛔ Menjawab isi layar TERBARU sesudah menyimpan - nama turunan dan baris
// riwayat baru dihitung server, dan layar yang menebaknya sendiri akan
// menampilkan yang tidak tersimpan.
func simpanPenawaranPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		var badan isiPenawaran
		if err := json.NewDecoder(r.Body).Decode(&badan); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON yang sah")
			return
		}
		isi, err := badan.keIsian()
		if err != nil {
			galat.Tulis(w, http.StatusBadRequest, err.Error())
			return
		}
		pelaku := inti.PelakuDari(r, stubPelaku)
		form := svc.FormPenawaran()
		if err := form.Simpan(r.Context(), pelaku, r.PathValue("id"), isi, time.Now()); jawabGalatPenawaran(w, err) {
			return
		}
		hasil, err := form.Baca(r.Context(), pelaku, r.PathValue("id"))
		if jawabGalatPenawaran(w, err) {
			return
		}
		galat.TulisJSON(w, hasil)
	}
}

// cariCedingPolis melayani GET /api/polis-life/cari-ceding.
func cariCedingPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		hasil, err := svc.FormPenawaran().CariCeding(r.Context(),
			inti.PelakuDari(r, stubPelaku), r.URL.Query().Get("cari"))
		if jawabGalatPenawaran(w, err) {
			return
		}
		galat.TulisJSON(w, hasil)
	}
}

// cariPemegangPolis melayani GET /api/polis-life/cari-pemegang-polis.
func cariPemegangPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		hasil, err := svc.FormPenawaran().CariPemegangPolis(r.Context(),
			inti.PelakuDari(r, stubPelaku), r.URL.Query().Get("cari"))
		if jawabGalatPenawaran(w, err) {
			return
		}
		galat.TulisJSON(w, hasil)
	}
}

// jawabGalatPenawaran menerjemahkan galat khas form penawaran, lalu
// menyerahkan sisanya ke `jawabGalatPolis`.
//
// ⛔ Isian belum lengkap dijawab dengan KALIMAT PEGA SAJA (satu per baris),
// tanpa awalan paket: layar menampilkannya apa adanya, sebagaimana
// `Page-Set-Messages` menampilkan `Local.Err`.
func jawabGalatPenawaran(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, models.ErrIsianPenawaranBelumLengkap):
		galat.Tulis(w, http.StatusBadRequest, sesudahSebab(err, models.ErrIsianPenawaranBelumLengkap))
	case errors.Is(err, services.ErrPenawaranBelumLengkap):
		// 409: keputusannya sah, isian polisnya yang belum siap.
		galat.Tulis(w, http.StatusConflict, sesudahSebab(err, services.ErrPenawaranBelumLengkap))
	case errors.Is(err, services.ErrPenawaranBukanTahapnya):
		galat.Tulis(w, http.StatusConflict, err.Error())
	default:
		return jawabGalatPolis(w, err)
	}
	return true
}

// sesudahSebab mengambil teks SESUDAH kalimat galat penanda di rantai galat.
func sesudahSebab(err, sebab error) string {
	teks := err.Error()
	if i := strings.Index(teks, sebab.Error()+": "); i >= 0 {
		return teks[i+len(sebab.Error())+2:]
	}
	return teks
}
