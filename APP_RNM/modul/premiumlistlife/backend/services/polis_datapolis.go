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
	// Peringatan - pesan langkah 9 SavePremiumList_Act sesudah Save Data; kosong
	// saat membaca.
	Peringatan []string `json:"peringatan"`
	Pilihan    struct {
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
// Simpan menjalankan `Save Data` = SavePremiumList_Act (tanpa Calculate1_Act,
// keputusan work owner 01-10-2026): simpan data polis, lalu periksa umur dan
// sum insured peserta terhadap batas produk (langkah 6-8).
//
// ⛔ Mengembalikan PERINGATAN, bukan galat: Pega menyimpan dengan
// `Obj-Save WithErrors=true` (langkah 15) dan menampilkan pesannya lewat
// `Page-Set-Messages` (langkah 9) - data polis tetap tersimpan.
func (f *FormDataPolis) Simpan(ctx context.Context, pelaku inti.Pelaku, polisID string,
	isi models.IsianDataPolis) ([]string, error) {

	if err := f.pagari(pelaku, polisID); err != nil {
		return nil, err
	}
	keadaan, err := repository.NewWorkPolis(f.svc.DB()).Keadaan(ctx, polisID)
	if errors.Is(err, repository.ErrWorkPolisTidakAda) {
		return nil, fmt.Errorf("%w: %s", ErrPolisTakDitemukan, polisID)
	}
	if err != nil {
		return nil, err
	}
	if models.KasusPolisTertutup(keadaan.Status) {
		return nil, fmt.Errorf("%w: polis %q berstatus %q", ErrKasusPolisTertutup, polisID, keadaan.Status)
	}
	if keadaan.Status != models.TahapPolisDetail {
		return nil, fmt.Errorf("%w: polis %q di tahap %q", ErrDataPolisBukanTahapnya, polisID, keadaan.Status)
	}
	siap, err := models.SusunDataPolis(isi)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", galat.ErrPermintaanTidakSah, err)
	}
	repo := repository.NewPenawaran(f.svc.DB())
	if err := f.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		return repo.SimpanDataPolis(ctx, tx, polisID, siap)
	}); err != nil {
		return nil, err
	}
	// Langkah 6-7: batas produk. Produk tanpa baris PRODUCTINWARD_LIFE tidak
	// diperiksa (lihat models.PeriksaBatasProduk).
	batas, ada, err := repo.BatasProduk(ctx, siap.ProductNameID)
	if err != nil || !ada {
		return nil, err
	}
	peserta, err := repo.PesertaBatas(ctx, polisID)
	if err != nil {
		return nil, err
	}
	return models.PeriksaBatasProduk(siap.Type, siap.RISlipRNM, batas, peserta), nil
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
