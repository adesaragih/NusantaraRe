package handlers

// Pintu HTTP data polis layar Input Premium Detail - tiket 03 bagian 2.
//
//	GET /api/polis-life/{id}/data-polis          isian + pilihan
//	PUT /api/polis-life/{id}/data-polis          Save Data: medan wajib + WPC (05-10-2026)
//	GET /api/polis-life/{id}/cari-produk?cari=   popup Choose Product Name
//	GET /api/polis-life/cari-marketing?cari=     autocomplete Marketing Officer
//	GET /api/polis-life/cari-rislip?cari=        autocomplete R/I SLIP RNM No.
//	GET /api/polis-life/rincian-produk?id=       isi Product Name (tombol View, 05-10-2026)
//	GET /api/polis-life/rate-produk?id=          isi R/I Rate baris PLAN LIST (05-10-2026)
//	GET /api/polis-life/risk-produk?id=          isi R/I Risk Name produk (05-10-2026)
//
// Billing Name dan Retrocessionaire memakai `GET /api/polis-life/cari-ceding`
// - report definition yang sama (BrowseCedingCoLife_RD).
//
// Nol aturan dagang di sini; seluruhnya di services/polis_datapolis.go.

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/premiumlistlife/backend/models"
	"nusantarare/modul/premiumlistlife/backend/repository"
	"nusantarare/modul/premiumlistlife/backend/services"
)

// isiDataPolis adalah badan PUT - angka desimal sebagai TEKS (ADR-U-0003).
type isiDataPolis struct {
	Type                string `json:"type"`
	ProductNameID       string `json:"productNameId"`
	ProductName         string `json:"productName"`
	SourceOfBusiness    string `json:"sourceOfBusiness"`
	SobName             string `json:"sobName"`
	CedingCo            string `json:"cedingCo"`
	CedingCoName        string `json:"cedingCoName"`
	PolicyHolder        string `json:"policyHolder"`
	PolicyHolderName    string `json:"policyHolderName"`
	RISlipRNM           string `json:"riSlipRnm"`
	ProRateType         string `json:"proRateType"`
	MoID                string `json:"moId"`
	MarketingCode       string `json:"marketingCode"`
	MarketingName       string `json:"marketingName"`
	AnnuityInterest     string `json:"annuityInterest"`
	PremiumRefundFactor string `json:"premiumRefundFactor"`
	RetroID             string `json:"retroId"`
	RetroName           string `json:"retroName"`
	SecurityReinsurerID string `json:"securityReinsurerId"`
	SecurityReinsurer   string `json:"securityReinsurer"`
	// DateReceived - `YYYY-MM-DD` (masukan `type="date"`); kosong = tidak diisi.
	DateReceived string `json:"dateReceived"`
}

func (i isiDataPolis) keIsian() (models.IsianDataPolis, error) {
	hasil := models.IsianDataPolis{
		Type: i.Type, ProductNameID: i.ProductNameID, ProductName: i.ProductName,
		SourceOfBusiness: i.SourceOfBusiness, SobName: i.SobName,
		CedingCo: i.CedingCo, CedingCoName: i.CedingCoName,
		PolicyHolder: i.PolicyHolder, PolicyHolderName: i.PolicyHolderName,
		RISlipRNM: i.RISlipRNM, ProRateType: i.ProRateType, MoID: i.MoID,
		MarketingCode: i.MarketingCode, MarketingName: i.MarketingName,
		RetroID: i.RetroID, RetroName: i.RetroName,
		SecurityReinsurerID: i.SecurityReinsurerID, SecurityReinsurer: i.SecurityReinsurer,
	}
	if v := strings.TrimSpace(i.DateReceived); v != "" {
		d, err := time.Parse("2006-01-02", v)
		if err != nil {
			return hasil, errors.New("dateReceived harus berbentuk YYYY-MM-DD")
		}
		hasil.DateReceived = &d
	}
	for _, d := range []struct {
		medan, nilai string
		tujuan       **apd.Decimal
	}{
		{"annuityInterest", i.AnnuityInterest, &hasil.AnnuityInterest},
		{"premiumRefundFactor", i.PremiumRefundFactor, &hasil.PremiumRefundFactor},
	} {
		if v := strings.TrimSpace(d.nilai); v != "" {
			x, err := utils.ParseDecimal(v)
			if err != nil {
				return hasil, errors.New(d.medan + " harus angka desimal (titik sebagai pemisah)")
			}
			*d.tujuan = x
		}
	}
	return hasil, nil
}

func bacaDataPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		hasil, err := svc.FormDataPolis().Baca(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("id"))
		if jawabGalatDataPolis(w, err) {
			return
		}
		galat.TulisJSON(w, hasil)
	}
}

func simpanDataPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		var badan isiDataPolis
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
		form := svc.FormDataPolis()
		rekapDihapus, err := form.Simpan(r.Context(), pelaku, r.PathValue("id"), isi)
		if jawabGalatDataPolis(w, err) {
			return
		}
		hasil, err := form.Baca(r.Context(), pelaku, r.PathValue("id"))
		if jawabGalatDataPolis(w, err) {
			return
		}
		// Type / Product Name berganti → rekap summary dihapus; layar meminta Calculate CSV
		// ulang sebelum Confirm (keputusan work owner 05-10-2026).
		hasil.RekapDihapus = rekapDihapus
		galat.TulisJSON(w, hasil)
	}
}

func cariProdukPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		hasil, err := svc.FormDataPolis().CariProduk(r.Context(), inti.PelakuDari(r, stubPelaku),
			r.PathValue("id"), r.URL.Query().Get("cari"))
		if jawabGalatDataPolis(w, err) {
			return
		}
		galat.TulisJSON(w, hasil)
	}
}

func cariMarketingPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		hasil, err := svc.FormDataPolis().CariMarketing(r.Context(), inti.PelakuDari(r, stubPelaku),
			r.URL.Query().Get("cari"))
		if jawabGalatDataPolis(w, err) {
			return
		}
		galat.TulisJSON(w, hasil)
	}
}

// rincianProdukPolis - tombol View di samping Product Name (05-10-2026).
func rincianProdukPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		hasil, err := svc.FormDataPolis().RincianProduk(r.Context(), inti.PelakuDari(r, stubPelaku),
			r.URL.Query().Get("id"))
		if errors.Is(err, repository.ErrRincianProdukTidakAda) {
			galat.Tulis(w, http.StatusNotFound, "Product not found in Master Product Name Life.")
			return
		}
		if jawabGalatDataPolis(w, err) {
			return
		}
		galat.TulisJSON(w, hasil)
	}
}

// rateProdukPolis - View Rate di popup Product Name (05-10-2026).
func rateProdukPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		hasil, err := svc.FormDataPolis().RateProduk(r.Context(), inti.PelakuDari(r, stubPelaku),
			r.URL.Query().Get("id"))
		if jawabGalatDataPolis(w, err) {
			return
		}
		galat.TulisJSON(w, hasil)
	}
}

// riskProdukPolis - View R/I Risk di popup Product Name (05-10-2026).
func riskProdukPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		hasil, err := svc.FormDataPolis().RiskProduk(r.Context(), inti.PelakuDari(r, stubPelaku),
			r.URL.Query().Get("id"))
		if jawabGalatDataPolis(w, err) {
			return
		}
		galat.TulisJSON(w, hasil)
	}
}

func cariRISlipPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		hasil, err := svc.FormDataPolis().CariRISlip(r.Context(), inti.PelakuDari(r, stubPelaku),
			r.URL.Query().Get("cari"))
		if jawabGalatDataPolis(w, err) {
			return
		}
		galat.TulisJSON(w, hasil)
	}
}

// jawabGalatDataPolis - galat khas data polis, lalu sisanya ke penawaran/polis.
func jawabGalatDataPolis(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, models.ErrDataPolisBelumLengkap):
		galat.Tulis(w, http.StatusBadRequest, sesudahSebab(err, models.ErrDataPolisBelumLengkap))
	case errors.Is(err, services.ErrDataPolisBukanTahapnya):
		galat.Tulis(w, http.StatusConflict, err.Error())
	default:
		return jawabGalatPenawaran(w, err)
	}
	return true
}
