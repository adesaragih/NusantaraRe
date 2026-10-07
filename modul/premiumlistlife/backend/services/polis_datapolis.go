package services

// Data polis layar Input Premium Detail - tiket 03 bagian 2 (bagian 1-3 layar).
//
// Untuk apa berkas ini: membaca dan menyimpan isian data polis tahap Input
// Premium Detail, dan ketiga pencarian master layar itu. Aturannya murni di
// models/polis_datapolis.go.
//
// ⛔ `Save Data` = `SavePremiumList_Act` SAJA; `Calculate1_Act` (hitung premi,
// retensi, spreading) TIDAK dijalankan - keputusan work owner 01-10-2026.
//
// Dibaca sesudah: models/polis_datapolis.go, repository/polis_datapolis.go.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/penomor"
	"nusantarare/modul/premiumlistlife/backend/models"
	"nusantarare/modul/premiumlistlife/backend/repository"
)

// ErrDataPolisBukanTahapnya - data polis hanya dapat disimpan di tahap
// Input Premium Detail (`ASSIGNMENT63`, flow action `ShowLifePremiumDetail`).
var ErrDataPolisBukanTahapnya = errors.New(
	"services: data polis hanya dapat disimpan di tahap Input Premium Detail")

// FormDataPolis melayani isian data polis layar Input Premium Detail.
type FormDataPolis struct{ svc *Service }

// FormDataPolis menyusun layanannya.
func (s *Service) FormDataPolis() *FormDataPolis { return &FormDataPolis{svc: s} }

// JawabanDataPolis adalah isi bagian data polis layar Input Premium Detail.
type JawabanDataPolis struct {
	CaseID              string     `json:"caseId"`
	Tahap               string     `json:"tahap"`
	BolehDisimpan       bool       `json:"bolehDisimpan"`
	Type                string     `json:"type"`
	ProductNameID       string     `json:"productNameId"`
	ProductName         string     `json:"productName"`
	SourceOfBusiness    string     `json:"sourceOfBusiness"`
	SobName             string     `json:"sobName"`
	CedingCo            string     `json:"cedingCo"`
	CedingCoName        string     `json:"cedingCoName"`
	PolicyHolder        string     `json:"policyHolder"`
	PolicyHolderName    string     `json:"policyHolderName"`
	RISlipRNM           string     `json:"riSlipRnm"`
	ProRateType         string     `json:"proRateType"`
	MoID                string     `json:"moId"`
	MarketingCode       string     `json:"marketingCode"`
	MarketingName       string     `json:"marketingName"`
	AnnuityInterest     string     `json:"annuityInterest"`
	PremiumRefundFactor string     `json:"premiumRefundFactor"`
	RetroID             string     `json:"retroId"`
	RetroName           string     `json:"retroName"`
	SecurityReinsurerID string     `json:"securityReinsurerId"`
	SecurityReinsurer   string     `json:"securityReinsurer"`
	WPC                 *time.Time `json:"wpc"`
	DateReceived        *time.Time `json:"dateReceived"`
	// RekapDihapus - Save Data mengganti Type atau Product Name sehingga rekap summary dihapus
	// (keputusan work owner 05-10-2026); Calculate CSV wajib dijalankan ulang
	// sebelum Confirm. Selalu false saat membaca.
	RekapDihapus bool `json:"rekapDihapus"`
	Pilihan      struct {
		Type        []models.Pilihan `json:"type"`
		ProRateType []models.Pilihan `json:"proRateType"`
	} `json:"pilihan"`
}

func (f *FormDataPolis) pagari(pelaku inti.Pelaku, polisID string) error {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return err
	}
	if strings.TrimSpace(polisID) == "" {
		return fmt.Errorf("%w: pengenal polis wajib diisi", galat.ErrPermintaanTidakSah)
	}
	if f == nil || f.svc == nil || !f.svc.PunyaDatabase() {
		return db.ErrTanpaOracle
	}
	return nil
}

// Baca membaca data polis satu kasus.
func (f *FormDataPolis) Baca(ctx context.Context, pelaku inti.Pelaku, polisID string) (JawabanDataPolis, error) {
	if err := f.pagari(pelaku, polisID); err != nil {
		return JawabanDataPolis{}, err
	}
	keadaan, err := repository.NewWorkPolis(f.svc.DB()).Keadaan(ctx, polisID)
	if errors.Is(err, repository.ErrWorkPolisTidakAda) {
		return JawabanDataPolis{}, fmt.Errorf("%w: %s", ErrPolisTakDitemukan, polisID)
	}
	if err != nil {
		return JawabanDataPolis{}, err
	}
	d, err := repository.NewPenawaran(f.svc.DB()).BacaDataPolis(ctx, polisID)
	if errors.Is(err, repository.ErrHeaderPolisTidakAda) {
		return JawabanDataPolis{}, fmt.Errorf("%w: %s", ErrPolisTakDitemukan, polisID)
	}
	if err != nil {
		return JawabanDataPolis{}, err
	}
	j := JawabanDataPolis{
		CaseID: keadaan.ID, Tahap: keadaan.Status,
		BolehDisimpan: keadaan.Status == models.TahapPolisDetail,
		Type:          d.Type, ProductNameID: d.ProductNameID, ProductName: d.ProductName,
		SourceOfBusiness: d.SourceOfBusiness, SobName: d.SobName,
		CedingCo: d.CedingCo, CedingCoName: d.CedingCoName,
		PolicyHolder: d.PolicyHolder, PolicyHolderName: d.PolicyHolderName,
		RISlipRNM: d.RISlipRNM, ProRateType: d.ProRateType, MoID: d.MoID,
		MarketingCode: d.MarketingCode, MarketingName: d.MarketingName,
		AnnuityInterest: teksUang(d.AnnuityInterest), PremiumRefundFactor: teksUang(d.PremiumRefundFactor),
		RetroID: d.RetroID, RetroName: d.RetroName,
		SecurityReinsurerID: d.SecurityReinsurerID, SecurityReinsurer: d.SecurityReinsurer,
		WPC:          d.WPC,
		DateReceived: d.DateReceived,
	}
	j.Pilihan.Type = models.PilihanTypePolis
	j.Pilihan.ProRateType = models.PilihanProRateType
	return j, nil
}

// Simpan menulis data polis - hanya di tahap Input Premium Detail.
//
// Keputusan work owner 05-10-2026: `Save Data` hanya memeriksa medan wajib,
// menghitung WPC, lalu menyimpan data polis + WPC dalam satu transaksi.
// Rekap summary TIDAK dihitung di sini (itu tugas Calculate CSV); bila Type
// atau Product Name berganti, rekap lama yang bergantung padanya DIHAPUS supaya Confirm tidak
// memakai rekap basi - Confirm menolak polis tanpa rekap.
//
// Mengembalikan true bila rekap dihapus.
func (f *FormDataPolis) Simpan(ctx context.Context, pelaku inti.Pelaku, polisID string,
	isi models.IsianDataPolis) (bool, error) {

	if err := f.pagari(pelaku, polisID); err != nil {
		return false, err
	}
	keadaan, err := repository.NewWorkPolis(f.svc.DB()).Keadaan(ctx, polisID)
	if errors.Is(err, repository.ErrWorkPolisTidakAda) {
		return false, fmt.Errorf("%w: %s", ErrPolisTakDitemukan, polisID)
	}
	if err != nil {
		return false, err
	}
	if models.KasusPolisTertutup(keadaan.Status) {
		return false, fmt.Errorf("%w: polis %q berstatus %q", ErrKasusPolisTertutup, polisID, keadaan.Status)
	}
	if keadaan.Status != models.TahapPolisDetail {
		return false, fmt.Errorf("%w: polis %q di tahap %q", ErrDataPolisBukanTahapnya, polisID, keadaan.Status)
	}
	siap, err := models.SusunDataPolis(isi)
	if err != nil {
		return false, fmt.Errorf("%w: %w", galat.ErrPermintaanTidakSah, err)
	}
	// ⛔ Save Data HANYA menyimpan dan memeriksa medan wajib (SusunDataPolis di
	// atas) - keputusan work owner 05-10-2026. Hitung summary pindah ke
	// Calculate CSV; pemeriksaan batas produk SavePremiumList_Act menjadi
	// penolakan Validate CSV (UnggahPremiumList.periksaDanHitung).
	// WPC dihitung dan ditulis saat Save Data (keputusan work owner 05-10-2026):
	// Type yang baru disimpan + periode produksi `TANGGAL_CLOSING` saat ini -
	// aturan periode yang SAMA dengan PL Number (`penomor.HitungPeriodeNomor`).
	// Confirm tetap menulis ulang WPC dari periode nomor yang terbit.
	hariClosing, err := repository.NewTutupBuku(f.svc.DB()).Tanggal(ctx)
	if err != nil {
		return false, err
	}
	periode, err := penomor.HitungPeriodeNomor(time.Now(), hariClosing)
	if err != nil {
		return false, err
	}
	wpc, err := models.WPCPolis(siap.Type, periode.MMYYYY)
	if err != nil {
		return false, fmt.Errorf("%w: %w", galat.ErrPermintaanTidakSah, err)
	}
	repo := repository.NewPenawaran(f.svc.DB())
	// Type dan Product Name TERSIMPAN dibaca sebelum UPDATE: rekap
	// T_PREMIUM_LIST_SUMMARY bergantung padanya (keputusan work owner 05-10-2026;
	// Product Name ikut sejak permintaan work owner berikutnya di hari yang sama).
	lama, err := repo.BacaDataPolis(ctx, polisID)
	if errors.Is(err, repository.ErrHeaderPolisTidakAda) {
		return false, fmt.Errorf("%w: %s", ErrPolisTakDitemukan, polisID)
	}
	if err != nil {
		return false, err
	}
	rekapBasi := strings.TrimSpace(lama.Type) != strings.TrimSpace(siap.Type) ||
		strings.TrimSpace(lama.ProductNameID) != strings.TrimSpace(siap.ProductNameID)
	if err := f.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		if err := repo.SimpanDataPolis(ctx, tx, polisID, siap); err != nil {
			return err
		}
		if err := repo.TulisWPC(ctx, tx, polisID, wpc); err != nil {
			return err
		}
		// Type / Product Name berganti: rekap lama basi - dihapus, BUKAN dihitung ulang
		// (hitung ulang hanya di Calculate CSV). Keduanya sama: rekap tidak disentuh.
		if rekapBasi {
			if _, err := repository.NewSummaryPolis(f.svc.DB()).HapusRekap(ctx, tx, polisID); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return false, err
	}
	return rekapBasi, nil
}

// CariProduk - popup Choose Product Name, disaring Ceding KASUS
// (`Ceding = pyWorkPage.CedingCo`).
func (f *FormDataPolis) CariProduk(ctx context.Context, pelaku inti.Pelaku, polisID, teks string) (
	[]repository.BarisProduk, error) {

	if err := f.pagari(pelaku, polisID); err != nil {
		return nil, err
	}
	d, err := repository.NewPenawaran(f.svc.DB()).BacaDataPolis(ctx, polisID)
	if errors.Is(err, repository.ErrHeaderPolisTidakAda) {
		return nil, fmt.Errorf("%w: %s", ErrPolisTakDitemukan, polisID)
	}
	if err != nil {
		return nil, err
	}
	return repository.NewRujukan(f.svc.DB()).CariProduk(ctx, d.CedingCo, teks)
}

// RincianProduk - isi satu Product Name untuk tombol "View" layar Input
// Premium Detail, dari tabel flat Master Product Name Life (keputusan work
// owner 05-10-2026). Baca-saja; ID diambil dari isian layar sehingga produk
// yang baru dipilih (belum disimpan) pun dapat dilihat.
func (f *FormDataPolis) RincianProduk(ctx context.Context, pelaku inti.Pelaku, produkID string) (
	repository.RincianProduk, error) {

	if err := f.svc.FormPenawaran().pagariRujukan(pelaku); err != nil {
		return repository.RincianProduk{}, err
	}
	id := strings.TrimSpace(produkID)
	if id == "" {
		return repository.RincianProduk{}, fmt.Errorf("%w: ID produk kosong", galat.ErrPermintaanTidakSah)
	}
	return repository.NewRujukan(f.svc.DB()).RincianProduk(ctx, id)
}

// RateProduk - isi satu R/I Rate baris PLAN LIST di popup Product Name
// (permintaan work owner 05-10-2026). Baca-saja.
func (f *FormDataPolis) RateProduk(ctx context.Context, pelaku inti.Pelaku, riRateID string) (
	repository.RateProduk, error) {

	if err := f.svc.FormPenawaran().pagariRujukan(pelaku); err != nil {
		return repository.RateProduk{}, err
	}
	id := strings.TrimSpace(riRateID)
	if id == "" {
		return repository.RateProduk{}, fmt.Errorf("%w: ID R/I Rate kosong", galat.ErrPermintaanTidakSah)
	}
	return repository.NewRujukan(f.svc.DB()).RateProduk(ctx, id)
}

// RiskProduk - isi R/I Risk Name produk di popup Product Name (permintaan
// work owner 05-10-2026). Baca-saja.
func (f *FormDataPolis) RiskProduk(ctx context.Context, pelaku inti.Pelaku, riRiskID string) (
	repository.RiskProduk, error) {

	if err := f.svc.FormPenawaran().pagariRujukan(pelaku); err != nil {
		return repository.RiskProduk{}, err
	}
	id := strings.TrimSpace(riRiskID)
	if id == "" {
		return repository.RiskProduk{}, fmt.Errorf("%w: ID R/I Risk kosong", galat.ErrPermintaanTidakSah)
	}
	return repository.NewRujukan(f.svc.DB()).RiskProduk(ctx, id)
}

// CariMarketing - autocomplete Marketing Officer.
func (f *FormDataPolis) CariMarketing(ctx context.Context, pelaku inti.Pelaku, teks string) (
	[]repository.BarisMarketing, error) {

	if err := f.svc.FormPenawaran().pagariRujukan(pelaku); err != nil {
		return nil, err
	}
	return repository.NewRujukan(f.svc.DB()).CariMarketing(ctx, teks)
}

// CariRISlip - autocomplete R/I SLIP RNM No.
func (f *FormDataPolis) CariRISlip(ctx context.Context, pelaku inti.Pelaku, teks string) (
	[]repository.BarisRujukan, error) {

	if err := f.svc.FormPenawaran().pagariRujukan(pelaku); err != nil {
		return nil, err
	}
	return repository.NewRujukan(f.svc.DB()).CariRISlip(ctx, teks)
}
