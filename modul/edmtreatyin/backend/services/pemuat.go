package services

// Untuk apa berkas ini: PEMUAT DOKUMEN LAMA ENDORSEMEN - tiket EDM 10 (spec-penyimpanan ID-3, AC 44) dan 09 (ID-33..
// ID-37, AC 39-43); KEPUTUSAN 23-09-2026 butir 4: setiap polis, SETIAP generasi, tanpa penyaring. Asal: pola
// `modul/nbtreatyin/backend/services/pemuat.go` (tiket NB 22; NB memuat PRODKE 0 dan hanya MENGHITUNG PRODKE > 0).
//
// Dijalankan MANUSIA lewat `backend/alat/pemuatlama` (seluruh dokumen), atau per dokumen dari layar lewat tombol Copy
// Old superadmin (`copyold.go`, perintah work owner 07-10-2026) - `muat` yang sama. Urutan: generasi NB lebih dulu
// oleh pemuat NB, lalu pemuat ini.
//
// ⚠️ URUTAN PEMUATAN (tiket 10): dokumen dimuat menurut NOPOLIS lalu PRODKE BILANGAN naik
// (`models.UrutKunciGenerasi`) - OLD_POLIS_ID menuntut generasi sebelumnya sudah ada. OLD_POLIS_ID = ID generasi
// bernomor polis sama ber-PRODKE - 1 (`GenerasiSebelumnya`); tidak ada = galat dokumen, bukan tebakan.
//
// Per dokumen: pecah (models, murni) -> SATU transaksi (ID-44) -> tulis lewat antarmuka penyimpanan yang SAMA dengan
// jalur biasa `simpanGenerasi` (tindakan.go) - nol pintu belakang, penjaga tiket 02/03/04 ikut berjalan:
//
//	KunciGenerasiLama       ID sudah ada dengan kunci sama = sudah dimuat (dilewati); kunci beda = galat
//	GenerasiSebelumnya      OLD_POLIS_ID (NOPOLIS, PRODKE - 1) + BacaGenerasi - pembanding (ID-8)
//	BarisSpreadingHilang    PENJAGA KEUTUHAN jalur biasa (ID-15, AC 8; validasiKirim)
//	HitungPenandaMigrasi    penanda tiket 09 (murni; angka tidak disentuh)
//	SisipKasus              T_WORK_POLIS + generasi PRODKE / NOENDORS / OLD_POLIS_ID / EDM_TYPE - PENJAGA PERCABANGAN
//	                        UNIQUE OLD_POLIS_ID di basis data (ID-10, AC 3)
//	SimpanHalaman           T_GENERAL_POLIS_TREATY + T_POLIS_* menurut katalog (RapikanBentukSimpan sudah di pemecah)
//	SetelNomorPolisSelesai  NOPOLIS = nomor polis induk (Utility1 SaveJsonPolisTreatyInEDM_Act; UQ NOPOLIS, PRODKE)
//	SetelKolomDatarLamaEDM  TGL_INPUT, USERNAME json_polis apa adanya
//	SimpanSelisih           proyeksi selisih SUMBER 'PEGA' dari TreatyDifference / TreatyXOLDifferenceList dokumen
//	                        APA ADANYA (ID-33, AC 39) - beku (ID-38)
//	SetelPenandaMigrasi     PASANGAN_BERGESER / RUMUS_BERLAPIS, SUMBER 'PEGA' saja (AC 43)
//	SalinUsulanLamaEDM      SuggestList -> HISTORYAKSEPTASIPRODUCTION lewat `CatatUsulan` (penjaga dobel IDPEGA)
//	TutupKasus              Resolved-Completed: json_polis endorsemen hanya lahir di Utility1 sesudah disetujui
//
// Uji-kering (`tulis` false): hanya JSON_POLIS yang dibaca, nol pernyataan ke tabel mana pun. Pembanding generasi
// sebelumnya = dokumen json_polis-nya (generasi NB lewat `BacaJSONPolisNB`, endorsemen dari jalankan ini), sehingga
// keutuhan, percabangan (dua dokumen ber-NOPOLIS + PRODKE sama), dan penanda ikut dilaporkan sebelum menulis.
//
// Gagal di tengah = dokumen itu tidak tersimpan sama sekali dan sebabnya masuk berkas laporan galat. Jalankan ulang
// aman: generasi yang sudah ada dengan kunci sama dilewati.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/edmtreatyin/backend/models"
	"nusantarare/modul/edmtreatyin/backend/repository"
)

// GudangPemuat - kebutuhan penyimpanan pemuat: pembaca JSON_POLIS + antarmuka tulis jalur biasa + dua penulis
// tambahan (`repository/lama_edm.go`). Implementasi Oracle: `penyimpanOracle` (gudang.go); uji memakai tiruan.
type GudangPemuat interface {
	Transaksi(ctx context.Context, fn func(tx *db.Tx) error) error

	// baca saja (repository/lama_edm.go)
	KunciJSONPolisEDM(ctx context.Context) ([]models.KunciJSONPolis, error)
	// KunciJSONPolisEDMCopyOld - popup Copy Old: hanya dokumen yang kasus Pega-nya ada dan sudah berproduksi (WO 07-10-2026).
	KunciJSONPolisEDMCopyOld(ctx context.Context) ([]models.KunciJSONPolis, error)
	BacaJSONPolisEDM(ctx context.Context, kunci string) (models.BarisJSONPolis, error)
	BacaJSONPolisNB(ctx context.Context, nopolis string) ([]models.BarisJSONPolis, error)
	KunciGenerasiLama(ctx context.Context, tx *db.Tx, id string) (models.KunciGenerasi, bool, error)

	// antarmuka jalur biasa (repository/kasus.go, generasi.go, polis.go, selisih.go, usulan.go)
	GenerasiSebelumnya(ctx context.Context, tx *db.Tx, k models.Kasus, nopolis string) (string, error)
	BacaGenerasi(ctx context.Context, tx *db.Tx, id string) (*models.Halaman, error)
	SisipKasus(ctx context.Context, tx *db.Tx, id, pembuat, namaPembuat string, gb repository.GenerasiBaru) error
	SimpanHalaman(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error
	SetelNomorPolisSelesai(ctx context.Context, tx *db.Tx, id, nopolis string) error
	SimpanSelisih(ctx context.Context, tx *db.Tx, polisID string, h *models.Halaman, k repository.KunciSelisih) error
	TutupKasus(ctx context.Context, tx *db.Tx, id, statusLama, statusAkhir string) error

	// penulis tambahan pemuat (repository/lama_edm.go)
	SetelKolomDatarLamaEDM(ctx context.Context, tx *db.Tx, id string, k models.KolomDatarLama) error
	SetelPenandaMigrasi(ctx context.Context, tx *db.Tx, polisID string, p models.PenandaMigrasi) error
	SalinUsulanLamaEDM(ctx context.Context, tx *db.Tx, idPega string, baris []models.UsulanProduksi) (models.NasibUsulan, error)
}

// Pemuat - pemuat dokumen lama endorsemen.
type Pemuat struct{ g GudangPemuat }

// PemuatBaru menyusun pemuat di atas gudang `g`.
func PemuatBaru(g GudangPemuat) *Pemuat { return &Pemuat{g: g} }

// PemuatDariDasar menyusun pemuat di atas basis data bersama.
func PemuatDariDasar(d *inti.Dasar) (*Pemuat, error) {
	if d == nil || !d.PunyaDatabase() {
		return nil, ErrTanpaOracle
	}
	return PemuatBaru(penyimpanOracle{Gudang: repository.Baru(d.DB()), dasar: d}), nil
}

// pmErrSudahDimuat - generasi sudah ada dengan kunci yang sama (jalankan ulang).
var pmErrSudahDimuat = errors.New("services: dokumen sudah dimuat")

// pmRantai - generasi satu nomor polis yang sudah dipecah di jalankan ini (pembanding uji-kering).
type pmRantai struct {
	nopolis string
	gen     map[int]*models.Halaman
}

// Jalankan memuat seluruh dokumen generasi endorsemen. Galat yang dikembalikan hanya yang menghentikan seluruh
// jalankan (JSON_POLIS tidak terbaca, berkas laporan tidak dapat ditulis, dibatalkan); galat per dokumen masuk
// laporan.
func (pm *Pemuat) Jalankan(ctx context.Context, tulis bool, lap *models.LaporanPemuat) error {
	kunci, err := pm.g.KunciJSONPolisEDM(ctx)
	if err != nil {
		return err
	}
	models.UrutKunciGenerasi(kunci)
	terlihat := map[string]string{}
	rantai := pmRantai{}
	for _, k := range kunci {
		if err := ctx.Err(); err != nil {
			return err
		}
		lap.Dibaca()
		b, err := pm.g.BacaJSONPolisEDM(ctx, k.Kunci)
		if err != nil {
			b = models.BarisJSONPolis{IDPega: "ROWID " + k.Kunci, NoPolis: k.NoPolis, ProdKe: k.ProdKe}
			if err := lap.Gagal(b, []models.GalatDokumen{{Err: err}}); err != nil {
				return err
			}
			continue
		}
		if rantai.nopolis != strings.TrimSpace(b.NoPolis) {
			rantai = pmRantai{nopolis: strings.TrimSpace(b.NoPolis), gen: map[int]*models.Halaman{}}
		}
		h, err := models.PecahDokumenEDM(b)
		switch {
		case errors.Is(err, models.ErrBukanTreatyIn), errors.Is(err, models.ErrBukanGenerasiEndorsemen),
			errors.Is(err, models.ErrBarisAplikasiBaru):
			lap.Lewat(err)
			continue
		case err != nil:
			if err := lap.Gagal(b, []models.GalatDokumen{{Err: err}}); err != nil {
				return err
			}
			continue
		case len(h.Galat) > 0:
			if err := lap.Gagal(b, h.Galat); err != nil {
				return err
			}
			continue
		}
		if pertama, ganda := terlihat[h.ID]; ganda {
			g := models.GalatDokumen{Jalur: "IDPEGA", Nilai: b.IDPega,
				Err: fmt.Errorf("%w: %s sudah diambil dari baris %s", models.ErrDokumenGanda, h.ID, pertama)}
			if err := lap.Gagal(b, []models.GalatDokumen{g}); err != nil {
				return err
			}
			continue
		}
		terlihat[h.ID] = "ROWID " + k.Kunci
		var p models.PenandaMigrasi
		usulan := models.NasibUsulanUjiKering(h)
		if tulis {
			p, usulan, err = pm.muat(ctx, h)
		} else {
			p, err = pm.ujiKering(ctx, h, rantai)
		}
		switch {
		case errors.Is(err, pmErrSudahDimuat):
			lap.SudahDimuat()
			rantai.gen[h.ProdKe] = h.Halaman
			continue
		case err != nil:
			if err := lap.Gagal(b, []models.GalatDokumen{{Err: err}}); err != nil {
				return err
			}
			continue
		}
		rantai.gen[h.ProdKe] = h.Halaman
		if err := lap.Berhasil(h, usulan, p); err != nil {
			return err
		}
	}
	return nil
}

// pmPeriksaRantai - pemeriksaan generasi baru `h` terhadap generasi sebelumnya `lama` yang dijalankan SEBELUM
// menulis, sama di uji-kering dan `-jalankan`: penjaga keutuhan jalur biasa (ID-15), pemilih varian rumus Pega
// (OldData.EDMNo dokumen lawan EDMNo generasi OLD_POLIS_ID - ID-36), lalu kedua penanda (halaman tidak diubah).
func pmPeriksaRantai(h models.HasilPecahEDM, lama *models.Halaman) (models.PenandaMigrasi, error) {
	cek := h.Halaman.Salin()
	models.PasangOldData(cek, lama)
	if pesan := models.BarisSpreadingHilang(cek); len(pesan) > 0 {
		return models.PenandaMigrasi{}, fmt.Errorf("%w: %s", models.ErrKeutuhan, strings.Join(pesan, "; "))
	}
	lamaEDMNo := strings.TrimSpace(lama.Ambil(models.HalamanPolis + ".EDMNo"))
	if (strings.TrimSpace(h.OldDataEDMNo) == "") != (lamaEDMNo == "") {
		return models.PenandaMigrasi{}, fmt.Errorf("%w: OldData.EDMNo %q, generasi sebelumnya %q",
			models.ErrOldDataTakSesuai, h.OldDataEDMNo, lamaEDMNo)
	}
	return models.HitungPenandaMigrasi(h.Halaman, lama), nil
}

// ujiKering - pemeriksaan rantai tanpa basis data: pembanding = dokumen json_polis generasi sebelumnya.
func (pm *Pemuat) ujiKering(ctx context.Context, h models.HasilPecahEDM, r pmRantai) (models.PenandaMigrasi, error) {
	if _, ada := r.gen[h.ProdKe]; ada {
		return models.PenandaMigrasi{}, fmt.Errorf("%w: PRODKE %d nomor polis ini sudah diambil dokumen lain", models.ErrPercabangan, h.ProdKe)
	}
	lama := r.gen[h.ProdKe-1]
	if lama == nil && h.ProdKe == 1 {
		nb, err := pm.g.BacaJSONPolisNB(ctx, h.NoPolis)
		if err != nil {
			return models.PenandaMigrasi{}, err
		}
		if len(nb) > 1 {
			return models.PenandaMigrasi{}, fmt.Errorf("%w: %d dokumen generasi NB ber-NOPOLIS sama", models.ErrDokumenGanda, len(nb))
		}
		if len(nb) == 1 {
			if lama, err = models.HalamanPembanding(nb[0]); err != nil {
				return models.PenandaMigrasi{}, fmt.Errorf("%w: dokumen generasi NB: %w", models.ErrGenerasiSebelumnyaTidakAda, err)
			}
		}
	}
	if lama == nil {
		return models.PenandaMigrasi{}, fmt.Errorf("%w: NOPOLIS %s PRODKE %d", models.ErrGenerasiSebelumnyaTidakAda, h.NoPolis, h.ProdKe-1)
	}
	return pmPeriksaRantai(h, lama)
}

// muat menulis satu dokumen dalam SATU transaksi (ID-44) lewat antarmuka jalur biasa.
func (pm *Pemuat) muat(ctx context.Context, h models.HasilPecahEDM) (models.PenandaMigrasi, models.NasibUsulan, error) {
	var p models.PenandaMigrasi
	usulan := models.UsulanTanpaBaris
	err := pm.g.Transaksi(ctx, func(tx *db.Tx) error {
		kg, ada, err := pm.g.KunciGenerasiLama(ctx, tx, h.ID)
		if err != nil {
			return err
		}
		if ada {
			if kg.NoPolis == h.NoPolis && kg.ProdKe == h.ProdKe && kg.NoEndors == h.EDMNo {
				return pmErrSudahDimuat
			}
			return fmt.Errorf("%w: %s = NOPOLIS %q PRODKE %d NOENDORS %q", models.ErrIDKasusBentrok, h.ID, kg.NoPolis, kg.ProdKe, kg.NoEndors)
		}
		lamaID, err := pm.g.GenerasiSebelumnya(ctx, tx, models.Kasus{ProdKe: h.ProdKe}, h.NoPolis)
		if err != nil {
			return err
		}
		if lamaID == "" {
			return fmt.Errorf("%w: NOPOLIS %s PRODKE %d", models.ErrGenerasiSebelumnyaTidakAda, h.NoPolis, h.ProdKe-1)
		}
		lama, err := pm.g.BacaGenerasi(ctx, tx, lamaID)
		if err != nil {
			return err
		}
		if p, err = pmPeriksaRantai(h, lama); err != nil {
			return err
		}
		gb := repository.GenerasiBaru{ProdKe: h.ProdKe, EDMNo: h.EDMNo, OldPolisID: lamaID, EDMType: h.EDMType}
		if err := pm.g.SisipKasus(ctx, tx, h.ID, "", "", gb); err != nil {
			if errors.Is(err, repository.ErrGenerasiSudahDiendorse) {
				return fmt.Errorf("%w: OLD_POLIS_ID %s: %w", models.ErrPercabangan, lamaID, err)
			}
			return err
		}
		if err := pm.g.SimpanHalaman(ctx, tx, h.ID, h.Halaman); err != nil {
			return err
		}
		if err := pm.g.SetelNomorPolisSelesai(ctx, tx, h.ID, h.NoPolis); err != nil {
			return err
		}
		if err := pm.g.SetelKolomDatarLamaEDM(ctx, tx, h.ID, h.Datar); err != nil {
			return err
		}
		if err := pm.g.SimpanSelisih(ctx, tx, h.ID, h.Halaman, repository.KunciSelisih{
			NoPolis: h.NoPolis, ProdKe: h.ProdKe, EDMNo: h.EDMNo, IDPega: models.KunciInstans(h.ID), Sumber: models.SumberPega,
		}); err != nil {
			return err
		}
		if err := pm.g.SetelPenandaMigrasi(ctx, tx, h.ID, p); err != nil {
			return err
		}
		if usulan, err = pm.g.SalinUsulanLamaEDM(ctx, tx, models.KunciInstans(h.ID), h.Usulan); err != nil {
			return err
		}
		return pm.g.TutupKasus(ctx, tx, h.ID, models.AssignmentAdmin, models.StatusSelesai)
	})
	return p, usulan, err
}
