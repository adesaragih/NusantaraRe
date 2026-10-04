package services

// Form penawaran polis - tiket 01 bagian 3 (layar Input Offer).
//
// Untuk apa berkas ini: membaca dan menyimpan isian layar
// `Section/InputOfferLife.xml`, riwayatnya (`AddHistorySuggest`), dan kedua
// popup pilihan master (Ceding, Policy Holder). Aturannya murni di
// models/polis_isianpenawaran.go; yang di sini gerbang dan transaksi.
//
// ⛔ Tombol `Save Offer` Pega menjalankan `AddHistorySuggest` lalu
// `InputOfferLife_ACT` (langkah 2-8 hidup) - SATU simpanan di sini: header
// dan satu baris riwayat dalam SATU transaksi. Langkah 4-6
// (`INSERTJSONOFFERLIFE`, `GetIdOffer_SQL` -> NoOffer) TIDAK ditiru (spec §12):
// `NO_OFFER` tetap kosong untuk kasus baru, dan sel `Confirmation Number`
// memang hanya tampil bila tidak kosong (`pyVisible NOTBLANK`).
//
// Dibaca sesudah: models/polis_isianpenawaran.go, repository/polis_penawaran.go.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/premiumlistlife/backend/models"
	"nusantarare/modul/premiumlistlife/backend/repository"
)

var (
	// ErrPenawaranBukanTahapnya - layar Input Offer hanya hidup di tahap
	// penawaran (`Assignment2`, flow action `InputDataOfferLife`).
	ErrPenawaranBukanTahapnya = errors.New("services: isian penawaran hanya dapat disimpan di tahap Input Offer Life")
	// ErrPenawaranBelumLengkap - `Confirm` atas penawaran yang isian wajibnya
	// belum tersimpan.
	ErrPenawaranBelumLengkap = errors.New("services: isian penawaran belum lengkap")
)

// FormPenawaran melayani layar Input Offer.
type FormPenawaran struct{ svc *Service }

// FormPenawaran menyusun layanannya.
func (s *Service) FormPenawaran() *FormPenawaran { return &FormPenawaran{svc: s} }

// PilihanPenawaran - pilihan tertutup layar, dikirim BERSAMA isiannya supaya
// layar tidak menyimpan salinan kedua daftar itu.
type PilihanPenawaran struct {
	TypeCeding      []models.Pilihan `json:"typeCeding"`
	ClassOfBusiness []models.Pilihan `json:"classOfBusiness"`
	Status          []models.Pilihan `json:"status"`
}

// JawabanPenawaran adalah isi layar Input Offer satu polis.
type JawabanPenawaran struct {
	CaseID string `json:"caseId"`
	// Flag - `FlagOnGoingPolicy`; menentukan radio Status mana yang tampil
	// dan `Initial` riwayat (Offer/Bind).
	Flag             string     `json:"flag"`
	Tahap            string     `json:"tahap"`
	BolehDisimpan    bool       `json:"bolehDisimpan"`
	NoOffer          string     `json:"noOffer"`
	CedingCo         string     `json:"cedingCo"`
	CedingCoName     string     `json:"cedingCoName"`
	PolicyHolder     string     `json:"policyHolder"`
	PolicyHolderName string     `json:"policyHolderName"`
	TypeCeding       string     `json:"typeCeding"`
	TypeCedingName   string     `json:"typeCedingName"`
	JenisAsuransi    string     `json:"jenisAsuransi"`
	BusinessCode     string     `json:"businessCode"`
	BusinessName     string     `json:"businessName"`
	DateReceived     *time.Time `json:"dateReceived"`
	Description      string     `json:"description"`
	// Sel penawaran migrasi 059 (+ SUM_INSURED, STATUS_UPDATE dari 051).
	BatasUsiaPeserta     *int   `json:"batasUsiaPeserta"`
	PeriodePertanggungan string `json:"periodePertanggungan"`
	// SumInsured - uang sebagai TEKS desimal (ADR-U-0003); "" bila kosong.
	SumInsured             string     `json:"sumInsured"`
	TanggalPenawaran       *time.Time `json:"tanggalPenawaran"`
	TanggalRespon          *time.Time `json:"tanggalRespon"`
	TanggalKonfirmasi      *time.Time `json:"tanggalKonfirmasi"`
	TBC                    *int       `json:"tbc"`
	TanggalTBC             *time.Time `json:"tanggalTbc"`
	StatusUpdate           string     `json:"statusUpdate"`
	KeteranganMarketing    string     `json:"keteranganMarketing"`
	QQName                 string     `json:"qqName"`
	JenisUsaha             string     `json:"jenisUsaha"`
	KetentuanUnderwriting  string     `json:"ketentuanUnderwriting"`
	TanggalKonfirmasiBalik *time.Time `json:"tanggalKonfirmasiBalik"`
	TanggalRealisasi       *time.Time `json:"tanggalRealisasi"`
	TanggalBind            *time.Time `json:"tanggalBind"`
	StatusFinal            string     `json:"statusFinal"`
	// Status - status TERAKHIR disimpan (`STATUS_PENAWARAN` baris utama, 062);
	// layar memilih radio-nya dari sini.
	Status  string                `json:"status"`
	Riwayat []models.BarisSuggest `json:"riwayat"`
	Pilihan PilihanPenawaran      `json:"pilihan"`
}

// pagari - identitas, bentuk permintaan, database, lalu keberadaan kasus.
func (f *FormPenawaran) pagari(pelaku inti.Pelaku, polisID string) error {
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

// Baca membaca isi layar Input Offer satu polis.
func (f *FormPenawaran) Baca(ctx context.Context, pelaku inti.Pelaku, polisID string) (JawabanPenawaran, error) {
	if err := f.pagari(pelaku, polisID); err != nil {
		return JawabanPenawaran{}, err
	}
	kerja := repository.NewWorkPolis(f.svc.DB())
	keadaan, err := kerja.Keadaan(ctx, polisID)
	if errors.Is(err, repository.ErrWorkPolisTidakAda) {
		return JawabanPenawaran{}, fmt.Errorf("%w: %s", ErrPolisTakDitemukan, polisID)
	}
	if err != nil {
		return JawabanPenawaran{}, err
	}
	flag, err := kerja.Bendera(ctx, polisID)
	if err != nil {
		return JawabanPenawaran{}, err
	}
	repo := repository.NewPenawaran(f.svc.DB())
	isi, err := repo.Baca(ctx, polisID)
	if errors.Is(err, repository.ErrHeaderPolisTidakAda) {
		return JawabanPenawaran{}, fmt.Errorf("%w: %s", ErrPolisTakDitemukan, polisID)
	}
	if err != nil {
		return JawabanPenawaran{}, err
	}
	riwayat, err := repo.Riwayat(ctx, polisID)
	if err != nil {
		return JawabanPenawaran{}, err
	}
	return JawabanPenawaran{
		CaseID: keadaan.ID, Flag: flag, Tahap: keadaan.Status,
		BolehDisimpan: keadaan.Status == models.TahapPolisPenawaran,
		NoOffer:       isi.NoOffer,
		CedingCo:      isi.CedingCo, CedingCoName: isi.CedingCoName,
		PolicyHolder: isi.PolicyHolder, PolicyHolderName: isi.PolicyHolderName,
		TypeCeding: isi.TypeCeding, TypeCedingName: isi.TypeCedingName,
		JenisAsuransi: jenisAsuransiTersimpan(isi),
		BusinessCode:  isi.BusinessCode, BusinessName: isi.BusinessName,
		DateReceived: isi.DateReceived, Description: isi.Description,
		BatasUsiaPeserta: isi.BatasUsiaPeserta, PeriodePertanggungan: isi.PeriodePertanggungan,
		SumInsured:       teksUang(isi.SumInsured),
		TanggalPenawaran: isi.TanggalPenawaran, TanggalRespon: isi.TanggalRespon,
		TanggalKonfirmasi: isi.TanggalKonfirmasi, TBC: isi.TBC, TanggalTBC: isi.TanggalTBC,
		StatusUpdate: isi.StatusUpdate, KeteranganMarketing: isi.KeteranganMarketing,
		QQName: isi.QQName, JenisUsaha: isi.JenisUsaha, KetentuanUnderwriting: isi.KetentuanUnderwriting,
		TanggalKonfirmasiBalik: isi.TanggalKonfirmasiBalik, TanggalRealisasi: isi.TanggalRealisasi,
		TanggalBind: isi.TanggalBind, StatusFinal: isi.StatusFinal,
		Status:  isi.Status,
		Riwayat: riwayat,
		Pilihan: PilihanPenawaran{
			TypeCeding:      models.PilihanTypeCeding,
			ClassOfBusiness: models.PilihanClassOfBusiness,
			Status:          models.PilihanStatusMenurutBendera(flag),
		},
	}, nil
}

// Simpan menyimpan isian penawaran dan menambah satu baris riwayat.
//
// ⛔ Gerbang tahap: layar ini milik `Assignment2` saja. Kasus yang sudah
// maju atau tertutup tidak dapat ditulisi penawarannya - isian yang berubah
// sesudah Premium List Detail berjalan akan membuat rekap tidak lagi cocok
// dengan penawaran yang melahirkannya.
func (f *FormPenawaran) Simpan(ctx context.Context, pelaku inti.Pelaku, polisID string,
	isi models.IsianPenawaran, saat time.Time) error {

	if err := f.pagari(pelaku, polisID); err != nil {
		return err
	}
	kerja := repository.NewWorkPolis(f.svc.DB())
	keadaan, err := kerja.Keadaan(ctx, polisID)
	if errors.Is(err, repository.ErrWorkPolisTidakAda) {
		return fmt.Errorf("%w: %s", ErrPolisTakDitemukan, polisID)
	}
	if err != nil {
		return err
	}
	if models.KasusPolisTertutup(keadaan.Status) {
		return fmt.Errorf("%w: polis %q berstatus %q", ErrKasusPolisTertutup, polisID, keadaan.Status)
	}
	if keadaan.Status != models.TahapPolisPenawaran {
		return fmt.Errorf("%w: polis %q di tahap %q", ErrPenawaranBukanTahapnya, polisID, keadaan.Status)
	}
	flag, err := kerja.Bendera(ctx, polisID)
	if err != nil {
		return err
	}
	siap, err := models.SusunPenawaran(flag, isi)
	if err != nil {
		return fmt.Errorf("%w: %w", galat.ErrPermintaanTidakSah, err)
	}
	repo := repository.NewPenawaran(f.svc.DB())
	return f.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
		// Migrasi 062 - satu baris per status (repository/polis_barisstatus.go):
		// kunci baris utama, pindahkan isian status lama ke barisnya sendiri
		// bila statusnya berganti, BARU timpa baris utama dengan status baru.
		statusLama, err := repo.StatusTerkunci(ctx, tx, polisID)
		if err != nil {
			return err
		}
		if _, err := repo.PindahkanBarisStatus(ctx, tx, polisID, statusLama, siap.Status); err != nil {
			return err
		}
		if err := repo.Simpan(ctx, tx, polisID, siap); err != nil {
			return err
		}
		_, err = repo.SisipSuggest(ctx, tx, polisID,
			models.SuggestBaru(flag, siap.Status, siap.Description, pelaku.AkunID, saat))
		return err
	})
}

// CariCeding - popup `Choose Ceding Name`.
func (f *FormPenawaran) CariCeding(ctx context.Context, pelaku inti.Pelaku, teks string) (
	[]repository.BarisRujukan, error) {

	if err := f.pagariRujukan(pelaku); err != nil {
		return nil, err
	}
	return repository.NewRujukan(f.svc.DB()).CariCeding(ctx, teks)
}

// CariPemegangPolis - popup `Policy Holder`.
func (f *FormPenawaran) CariPemegangPolis(ctx context.Context, pelaku inti.Pelaku, teks string) (
	[]repository.BarisRujukan, error) {

	if err := f.pagariRujukan(pelaku); err != nil {
		return nil, err
	}
	return repository.NewRujukan(f.svc.DB()).CariPemegangPolis(ctx, teks)
}

func (f *FormPenawaran) pagariRujukan(pelaku inti.Pelaku) error {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return err
	}
	if f == nil || f.svc == nil || !f.svc.PunyaDatabase() {
		return db.ErrTanpaOracle
	}
	return nil
}

// periksaPenawaranLengkap - gerbang `Confirm` di tahap penawaran.
//
// ⛔ Di Pega `Confirm` di tahap ini adalah `finishAssignment` layar yang SAMA,
// dan finishAssignment menegakkan `pyRequired` selnya: System Reinsurance,
// Class of Business, Comment - ditambah Ceding Name dan Policy Holder
// (keputusan work owner 01-10-2026). Status tidak diperiksa: kolomnya tidak
// ada di header. Yang diperiksa di sini nilai yang TERSIMPAN -
// penawaran yang belum pernah disimpan tidak boleh maju.
func (p *Penawaran) periksaPenawaranLengkap(ctx context.Context, polisID string) error {
	isi, err := repository.NewPenawaran(p.svc.DB()).Baca(ctx, polisID)
	if errors.Is(err, repository.ErrHeaderPolisTidakAda) {
		return fmt.Errorf("%w: %s", ErrPolisTakDitemukan, polisID)
	}
	if err != nil {
		return err
	}
	if kurang := models.KekuranganPenawaran(models.WajibPenawaran{
		CedingCoName: isi.CedingCoName, PolicyHolderName: isi.PolicyHolderName,
		TypeCeding: isi.TypeCeding, BusinessCode: isi.BusinessCode, Description: isi.Description,
		DateReceived: isi.DateReceived,
	}); len(kurang) > 0 {
		return fmt.Errorf("%w: %s", ErrPenawaranBelumLengkap, models.GabungPesanPenawaran(kurang))
	}
	return nil
}

// teksUang - desimal ke teks layar; "" bila kosong (ADR-U-0027: kosong bukan nol).
func teksUang(d *apd.Decimal) string {
	if d == nil {
		return ""
	}
	return utils.FormatDecimal(d)
}

// jenisAsuransiTersimpan - "Reinsurance Type" yang tersimpan, atau turunannya
// dari System Reinsurance bila belum pernah tersimpan (kasus sebelum 061).
func jenisAsuransiTersimpan(isi repository.IsianTersimpan) string {
	if strings.TrimSpace(isi.JenisAsuransi) != "" {
		return isi.JenisAsuransi
	}
	return models.JenisAsuransi(isi.TypeCeding)
}
