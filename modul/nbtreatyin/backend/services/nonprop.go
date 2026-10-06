package services

// Untuk apa berkas ini: JALUR NB NON-PROPORSIONAL (XOL) di layanan -
// `[keputusan work owner]` K8. Urutan langkah Pega yang membaca master dan tabel
// acuan; perhitungannya fungsi murni `models` (nonprop.go, nonprop_detail.go).
//
//	pilih bisnis  InputPolicyTreatyInDetail_preACT 16 (Quotation.ProportionalType
//	              NonProportional) -> InputPolicyTreatyInDetail_NonProp, lalu 18
//	pra-proses    InputPolicyTreatyInPre_Act 10 -> TreatyRealizationCheckXOLList
//	buka berkas   master dibaca ulang untuk subsection DetailPoliciesNonProportional
//
// Master dibaca lewat `PembacaMasterTreaty` (baca-saja); halaman `TreatyIn` tidak
// pernah disimpan.

import (
	"context"
	"errors"
	"fmt"

	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/repository"
)

// Galat pembacaan master XOL - handlers menjawabnya 422 (AC 37: kegagalan
// membaca data kontrak ditampilkan kepada pengguna).
var (
	ErrMasterXOLTidakAda = repository.ErrMasterXOLTidakAda
	ErrMasterXOLRusak    = repository.ErrMasterXOLRusak
)

// nonPropDipilih = syarat preACT 16: `pyWorkPage.Quotation.ProportionalType=="NonProportional"`.
func nonPropDipilih(h *models.Halaman) bool {
	return h.Ambil(models.HalamanQuotation+".ProportionalType") == models.JenisNonProporsional
}

// rusak membungkus galat hitung atas nilai master (bukan angka, bagi nol) -
// sumbernya dokumen master, bukan isian layar.
func rusak(err error) error {
	if err == nil || errors.Is(err, ErrMasterXOLRusak) {
		return err
	}
	if errors.Is(err, models.ErrBukanAngka) || errors.Is(err, models.ErrBagiNol) {
		return fmt.Errorf("%w: %w", ErrMasterXOLRusak, err)
	}
	return err
}

// bacaMaster membaca master kontrak polis. `ketat` = pilih bisnis: master yang
// tidak ada / rusak MENGHENTIKAN proses (AC 36-38). Tidak ketat = pra-proses dan
// pembukaan berkas: master kosong, persis halaman Pega yang `adoptJSONObject`-nya
// gagal (galat hanya di log) - TreatyRealizationCheckXOLList lalu memasang
// pesannya sendiri.
func (l *Layanan) bacaMaster(ctx context.Context, h *models.Halaman, ketat bool) (models.MasterXOL, error) {
	m, err := l.g.MasterXOL(ctx, h.Ambil(models.HalamanPolis+".NoOffer"))
	if err != nil {
		if !ketat && (errors.Is(err, ErrMasterXOLTidakAda) || errors.Is(err, ErrMasterXOLRusak)) {
			return models.MasterXOL{}, nil
		}
		return models.MasterXOL{}, err
	}
	return m, nil
}

// idMataUangMaster = RDB `GetDataCurrencyByName_SQL` per mata uang angsuran
// master (`select ID ... from CURRENCY where CURRENCY = {CARI8}` - pembacanya
// `IDMataUangDariNama`, tabel dan kolom yang sama dengan `GetCurrencyIDByName`).
func (l *Layanan) idMataUangMaster(ctx context.Context, h *models.Halaman) (models.IDMataUang, error) {
	ids := models.IDMataUang{}
	for _, n := range models.MataUangAngsuranMaster(h) {
		id, err := l.g.IDMataUangDariNama(ctx, n)
		if err != nil {
			return nil, err
		}
		ids[n] = id
	}
	return ids, nil
}

// pilihBisnisNonProp = preACT langkah 16 dan 18 untuk kontrak NonProportional.
//
//	16  InputPolicyTreatyInDetail_NonProp: 2-6 master (RDB BrowseTreatyInJoinEDM +
//	    adoptJSONObject) lalu `models.InputDetailNonProp`
//	18  FlagPPH: PPN/PPH per layer dan angsuran (`models.PPNPPHLapisanXOL`)
//
// `[tafsiran]` Langkah 18 tidak bersyarat jenis di XML, tetapi membaca
// `TreatyIn.LimitShareSummaryList` - ringkasan layer XOL yang hanya dimuat jalur
// ini (K8); untuk kontrak proporsional master XOL tidak dibaca dan langkahnya
// tidak berbuat apa pun.
func (l *Layanan) pilihBisnisNonProp(ctx context.Context, h *models.Halaman) error {
	if !nonPropDipilih(h) {
		return nil
	}
	m, err := l.bacaMaster(ctx, h, true)
	if err != nil {
		return err
	}
	models.TerapkanMasterXOL(h, m)
	ids, err := l.idMataUangMaster(ctx, h)
	if err != nil {
		return err
	}
	if err := models.InputDetailNonProp(h, ids); err != nil {
		return rusak(err)
	}
	return rusak(models.PPNPPHLapisanXOL(h))
}

// pajakNonProp - `[keputusan work owner 06-10-2026]` With Tax / Type Tax berubah sesudah Choose Business:
// pajak polis NonProp baru dihitung ulang dengan menjalankan lagi preACT 16 dan 18 (`pilihBisnisNonProp`),
// satu-satunya tempat XML menghitungnya (sel `.TypeTax` hanya postValue, `.FlagPPH` hanya RemoveTypeTax_ACT).
// Isian pengguna yang ikut ditulis ulang langkah 13 dan 23 - StartDate, EndDate, %Share spreading -
// dikembalikan (`models.SimpanIsianNonProp`).
func (l *Layanan) pajakNonProp(ctx context.Context, h *models.Halaman) error {
	if !models.PolisNonPropBaru(h) || !nonPropDipilih(h) {
		return nil
	}
	isian := models.SimpanIsianNonProp(h)
	if err := l.pilihBisnisNonProp(ctx, h); err != nil {
		return err
	}
	return isian.Kembalikan(h)
}

// siapkanNonProp - bagian NonProp pra-proses (dipanggil `siapkan` sesudah
// `muatMaster`):
//
//  1. kontrak NonProportional: master dibaca ulang untuk tampilan (halaman
//     `TreatyIn` tidak disimpan) - `models.TampilanMasterNonProp`;
//  2. InputPolicyTreatyInPre_Act 10 (`models.PerluCekDaftarXOL`) ->
//     TreatyRealizationCheckXOLList: daftar XOL kosong -> master (SetTreatyIn_Act
//     4-5, `Page-Copy` 5) -> `models.LengkapiDaftarXOL`.
//
// Master dibaca paling banyak sekali per pra-proses.
func (l *Layanan) siapkanNonProp(ctx context.Context, h *models.Halaman) error {
	var master *models.MasterXOL
	baca := func() (models.MasterXOL, error) {
		if master == nil {
			m, err := l.bacaMaster(ctx, h, false)
			if err != nil {
				return models.MasterXOL{}, err
			}
			master = &m
		}
		return *master, nil
	}
	if nonPropDipilih(h) {
		m, err := baca()
		if err != nil {
			return err
		}
		models.TerapkanMasterXOL(h, m)
		if err := models.TampilanMasterNonProp(h); err != nil {
			return rusak(err)
		}
	}
	if !models.PerluCekDaftarXOL(h) || !models.AwalCekDaftarXOL(h) {
		return nil
	}
	m, err := baca()
	if err != nil {
		return err
	}
	models.TerapkanMasterXOL(h, m) // Page-Copy TreatyIn -> pyWorkPage.TreatyIn (master mentah)
	ids, err := l.idMataUangMaster(ctx, h)
	if err != nil {
		return err
	}
	return rusak(models.LengkapiDaftarXOL(h, ids))
}
