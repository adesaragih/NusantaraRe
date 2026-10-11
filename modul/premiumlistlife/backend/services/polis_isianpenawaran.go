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
// dan satu baris riwayat dalam SATU transaksi. Langkah 3-4 ("Pega to
// json_offer_life", `SaveOfferJsonLife_SQL`) DITIRU sejak 06-10-2026 (keputusan
// work owner) dengan SQL langsung ke `JSON_OFFER_LIFE` di transaksi yang sama -
// repository/polis_jsonoffer.go; prosedur `INSERTJSONOFFERLIFE` tidak dipanggil,
// JSONDATA/OLDID NULL. Langkah 5-6 (`GetIdOffer_SQL` -> NoOffer) TIDAK ditiru:
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
	"nusantarare/inti/backend/jejak"
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
type FormPenawaran struct {
	svc *Service
	// jejak - dipakai bila Save Offer Bind/Closed menutup kasus (08-10-2026).
	jejak jejak.Jejak
}

// FormPenawaran menyusun layanannya.
// ErrPenawaranMelebihiKolom - isian tidak muat di kolom JSON_OFFER_LIFE (400, berpesan).
var ErrPenawaranMelebihiKolom = errors.New("services: isian penawaran melebihi lebar kolom")

func (s *Service) FormPenawaran() *FormPenawaran {
	return &FormPenawaran{svc: s, jejak: jejak.JejakBelumDiputuskan{}}
}

// DenganJejak mengganti perekam jejak penutupan Bind/Closed.
func (f *FormPenawaran) DenganJejak(j jejak.Jejak) *FormPenawaran {
	salin := *f
	salin.jejak = j
	return &salin
}

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
	// StatusTutup - status kerja bila Save Offer ini MENUTUP kasus (Bind/Closed,
	// keputusan work owner 08-10-2026); kosong selain itu.
	StatusTutup string `json:"statusTutup"`
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
// ⭐ Keputusan work owner 08-10-2026: kasus bendera "0" dengan Status `Bind` /
// `Closed` sesudah tersimpan juga DITUTUP (Resolved-Completed / Resolved-
// Rejected) lewat `Penawaran.terapkanDalam` - jalur tutup yang SAMA dengan
// Confirm/Decline, di transaksi yang SAMA: gagal menutup = Save Offer batal.
// Bind menjalankan gerbang Confirm tahap penawaran (isian wajib lengkap) atas
// isian yang SEDANG disimpan. Akibatnya dikembalikan; kosong = tetap terbuka.
//
// ⛔ Gerbang tahap: layar ini milik `Assignment2` saja. Kasus yang sudah
// maju atau tertutup tidak dapat ditulisi penawarannya - isian yang berubah
// sesudah Premium List Detail berjalan akan membuat rekap tidak lagi cocok
// dengan penawaran yang melahirkannya.
func (f *FormPenawaran) Simpan(ctx context.Context, pelaku inti.Pelaku, polisID string,
	isi models.IsianPenawaran, saat time.Time) (models.AkibatKeputusan, error) {

	if err := f.pagari(pelaku, polisID); err != nil {
		return models.AkibatKeputusan{}, err
	}
	kerja := repository.NewWorkPolis(f.svc.DB())
	keadaan, err := kerja.Keadaan(ctx, polisID)
	if errors.Is(err, repository.ErrWorkPolisTidakAda) {
		return models.AkibatKeputusan{}, fmt.Errorf("%w: %s", ErrPolisTakDitemukan, polisID)
	}
	if err != nil {
		return models.AkibatKeputusan{}, err
	}
	if models.KasusPolisTertutup(keadaan.Status) {
		return models.AkibatKeputusan{}, fmt.Errorf("%w: polis %q berstatus %q", ErrKasusPolisTertutup, polisID, keadaan.Status)
	}
	if keadaan.Status != models.TahapPolisPenawaran {
		return models.AkibatKeputusan{}, fmt.Errorf("%w: polis %q di tahap %q", ErrPenawaranBukanTahapnya, polisID, keadaan.Status)
	}
	flag, err := kerja.Bendera(ctx, polisID)
	if err != nil {
		return models.AkibatKeputusan{}, err
	}
	siap, err := models.SusunPenawaran(flag, isi)
	if err != nil {
		return models.AkibatKeputusan{}, fmt.Errorf("%w: %w", galat.ErrPermintaanTidakSah, err)
	}
	// Bind/Closed (bendera "0") menutup kasus - keputusan work owner 08-10-2026.
	akibat, keputusan, sebab, tutup, err := models.AkibatStatusOffer(flag, siap.Status)
	if err != nil {
		return models.AkibatKeputusan{}, fmt.Errorf("polis %q: %w", polisID, err)
	}
	if tutup && akibat.SimpanPolis {
		// Jalur Offer tidak pernah melewati Utility1 - nomor PL, LIFEINPRODUCTION,
		// JSON_POLIS, M_LIFE_PREMIUM_* TIDAK boleh terbit dari Save Offer.
		return models.AkibatKeputusan{}, fmt.Errorf("%w: Save Offer %q tidak boleh menyimpan polis",
			galat.ErrPermintaanTidakSah, siap.Status)
	}
	// Gerbang Confirm tahap penawaran (periksaPenawaranLengkap) atas isian yang SEDANG
	// disimpan - nilai tersimpan baru terlihat di dalam transaksi ini. Decline tidak digerbangi.
	if tutup && keputusan == models.KeputusanConfirm {
		if kurang := models.KekuranganPenawaran(models.WajibPenawaran{
			CedingCoName: siap.CedingCoName, PolicyHolderName: siap.PolicyHolderName,
			TypeCeding: siap.TypeCeding, BusinessCode: siap.BusinessCode, Description: siap.Description,
			DateReceived: siap.DateReceived,
		}); len(kurang) > 0 {
			return models.AkibatKeputusan{}, fmt.Errorf("%w: %s", ErrPenawaranBelumLengkap, models.GabungPesanPenawaran(kurang))
		}
	}
	repo := repository.NewPenawaran(f.svc.DB())
	err = f.svc.DalamTransaksi(ctx, func(tx *db.Tx) error {
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
		// JSON_OFFER_LIFE dari isian yang BARU tersimpan (keputusan work owner
		// 06-10-2026): ID lama dipertahankan; gagal = seluruh Save Offer batal.
		if _, err := repository.NewJSONOfferLife(f.svc.DB()).Ganti(ctx, tx, polisID); err != nil {
			if errors.Is(err, repository.ErrNilaiJSONOfferTerlaluPanjang) {
				return fmt.Errorf("%w: %s", ErrPenawaranMelebihiKolom,
					strings.TrimPrefix(err.Error(), repository.ErrNilaiJSONOfferTerlaluPanjang.Error()+": "))
			}
			return err
		}
		if _, err := repo.SisipSuggest(ctx, tx, polisID,
			models.SuggestBaru(flag, siap.Status, siap.Description, pelaku.AkunID, saat)); err != nil {
			return err
		}
		if !tutup {
			return nil
		}
		// Penutupan Bind/Closed: jalur Confirm/Decline yang sama (TutupKasus + jejak),
		// SESUDAH penawaran tersimpan, di transaksi ini.
		_, err = f.svc.Penawaran().DenganJejak(f.jejak).terapkanDalam(ctx, tx, pelaku, keadaan, akibat, sebab, saat)
		return err
	})
	if err != nil {
		return models.AkibatKeputusan{}, err
	}
	if !tutup {
		return models.AkibatKeputusan{}, nil
	}
	return akibat, nil
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
