package services

// Asal: salinan sebagian modul/nbtreatyin/backend/services/nonprop.go (06-10-2026): pra-proses
// `InputPolicyTreatyInPre_Act` langkah 10 (identik di korpus EDM). Jalur pilih bisnis NB (preACT 16-18) dibuang -
// EDM memakai `EDMChooseBusiness_Act` (tindakan.go PilihBisnis, models.PilihBisnisEDM).
//
// Untuk apa berkas ini: JALUR NON-PROPORSIONAL (XOL) di layanan -
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

	"nusantarare/modul/edmtreatyin/backend/models"
	"nusantarare/modul/edmtreatyin/backend/repository"
)

// Galat pembacaan master XOL - handlers menjawabnya 422 (AC 37: kegagalan
// membaca data kontrak ditampilkan kepada pengguna).
var (
	ErrMasterXOLTidakAda = repository.ErrMasterXOLTidakAda
	ErrMasterXOLRusak    = repository.ErrMasterXOLRusak
)

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
	// ⛔ NB menampilkan master XOL di subsection `DetailPoliciesNonProportional` (nonPropDipilih) - section itu
	// tidak ada di layar EDM; yang tersisa langkah 10 InputPolicyTreatyInPre_Act.
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
