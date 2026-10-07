package repository

// Untuk apa berkas ini: HALAMAN KASUS ENDORSEMEN - penyusun `pyWorkPage` satu kasus EDM dari tabel relasional:
//
//	PolicyTreatyIn.*                 generasi ini (T_GENERAL_POLIS_TREATY + T_POLIS_*), BacaGenerasi
//	PolicyTreatyIn.EDMNo / ProdKe    NOENDORS / PRODKE generasi ini
//	PolicyTreatyIn.PolicyNo          NOPOLIS (terisi sesudah selesai) - selama berjalan = OldData.PolicyNo
//	                                 (`SetEDMTNoPolis` langkah 3: PolicyNo = OldData.PolicyNo)
//	PolicyTreatyIn.OldData.*         generasi TEPAT sebelumnya (OLD_POLIS_ID) - `CreateEDMT` langkah 7-10 menyalin
//	                                 json_polis generasi terakhir ke OldData; di sini dibaca dari baris generasi itu,
//	                                 yang beku (generasi tertutup tidak dapat disunting - ID-13, P58)
//	PolicyTreatyIn.OldData.EDMNo     NOENDORS generasi sebelumnya (pemilih tab `PropOldData` / `PropOldData2`)
//	PolicyTreatyIn.OldData.TreatyDifference.*  selisih generasi sebelumnya (tab Old Data `PropOldData2`)
//	PolicyTreatyIn.TreatyDifference.* / TreatyXOLDifferenceList  proyeksi selisih generasi ini (360-363)
//	PolicyTreatyIn.SuggestList       HISTORYAKSEPTASIPRODUCTION (ketetapan NB K4, models/usulan.go)
//
// ⛔ `OldData` di dalam `OldData` (Pega `CreateEDMT` langkah 10 menyalin seluruh halaman, nol pembaca - spec ID-8,
// AC 9) tidak dibangun: hanya SATU tingkat OldData.

import (
	"context"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/edmtreatyin/backend/models"
)

// BacaHalaman memuat halaman kerja satu kasus endorsemen.
func (g *Gudang) BacaHalaman(ctx context.Context, tx *db.Tx, id string) (*models.Halaman, error) {
	k, err := g.Keadaan(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	h, err := g.BacaGenerasi(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	nopolis := h.Ambil(models.HalamanPolis + ".PolicyNo")
	if nopolis == "" {
		nopolis = models.NilaiQuotation(h, "OldPolicyNo")
	}
	lama, err := g.GenerasiSebelumnya(ctx, tx, k, nopolis)
	if err != nil {
		return nil, err
	}
	if lama != "" {
		hl, err := g.BacaGenerasi(ctx, tx, lama)
		if err != nil {
			return nil, err
		}
		if err := g.BacaSelisih(ctx, tx, lama, hl, ""); err != nil {
			return nil, err
		}
		models.PasangOldData(h, hl)
	}
	if h.Ambil(models.HalamanPolis+".PolicyNo") == "" {
		h.Setel(models.HalamanPolis+".PolicyNo", h.Ambil(models.HalamanPolis+".OldData.PolicyNo"))
	}
	if err := g.BacaSelisih(ctx, tx, id, h, ""); err != nil {
		return nil, err
	}
	if err := g.bacaUsulan(ctx, tx, id, h); err != nil {
		return nil, err
	}
	h.Setel("pyID", id)
	return h, nil
}
