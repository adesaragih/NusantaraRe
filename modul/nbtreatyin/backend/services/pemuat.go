package services

// Untuk apa berkas ini: PEMUAT DOKUMEN LAMA - tiket 22 (spec-penyimpanan
// ID-3, AC 55-59; KEPUTUSAN-RONDE-12 butir 5; K15, K17).
//
// Dijalankan MANUSIA lewat perintah `backend/alat/pemuatlama` (MODUL.md
// "Pemuat dokumen lama") untuk seluruh dokumen, atau per dokumen dari layar
// lewat tombol Copy Old superadmin (`copyold.go`, perintah work owner
// 07-10-2026) - `muat` yang sama. Seluruh dokumen generasi NB di `POOLDATA.JSON_POLIS`
// (PRODKE 0) - setiap polis, tanpa penyaring status/tahun/lini usaha (butir 5).
// Generasi endorsemen (PRODKE > 0) milik pemuat EDM (edmtreatyin tiket 10).
//
// Per dokumen: pecah (models, murni) -> SATU transaksi -> tulis lewat
// antarmuka penyimpanan yang SAMA dengan jalur biasa (ID-3, AC 56):
//
//	SisipKasus      T_WORK_POLIS + generasi PRODKE 0 (pembuat tidak dikarang: NULL)
//	SimpanHalaman   T_GENERAL_POLIS_TREATY + T_POLIS_QUOTATION/CEDING/INSTALMENT(_DETAIL)/
//	                SPREADING/XOL/XOL_LAYER menurut katalog - penjaga yang sama
//	                (PeriksaBentukSimpan AC 31/33, konversi kolom ID-14..18)
//	SetelNomorPolis NOPOLIS = JSON_POLIS.NOPOLIS (indeks unik NOPOLIS, PRODKE)
//	SetelKolomDatarLama  IDPEGA, NOENDORS, TGL_INPUT, USERNAME apa adanya (ID-21)
//	TutupKasus      Resolved-Completed: dokumen JSON_POLIS hanya lahir di jalur
//	                Decision8 "Nopolis not empty" -> Utility1 -> Utility2 -> End3
//	                (Flow InputRealizationTreatyIn) - kasusnya sudah selesai.
//	SalinUsulanLama SuggestList dokumen -> POOLDATA.HISTORYAKSEPTASIPRODUCTION
//	                lewat `CatatUsulan` jalur biasa (F3, WO 04-10-2026), dengan
//	                penjaga dobel menurut IDPEGA (IDPEGA = JSON_POLIS.IDPEGA =
//	                pzInsKey kasus lama, sama dengan `InsertViewSuggest_SQL`).
//
// Baris hasil pemuat dikenali dari IDPEGA (`<kelas> <pyID>`, jalur biasa
// menulis ID kasus) dan status Resolved-Completed - tanpa kolom penanda
// SUMBER (F6, WO 04-10-2026: penanda `SUMBER='PEGA'` gugur).
//
// Gagal di tengah = dokumen itu tidak tersimpan sama sekali (P2, AC 45) dan
// sebabnya masuk berkas laporan galat (AC 58). Jalankan ulang aman: kasus yang
// barisnya sudah ada dengan IDPEGA sama dilewati, IDPEGA berbeda = galat.
// ⛔ Seam services TIDAK diuji (spec.md §6.2): bagian murninya di models
// (`PecahDokumenLama`, `LaporanPemuat`), penulisannya di repository (tag db).

import (
	"context"
	"errors"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/repository"
)

// GudangPemuat - kebutuhan penyimpanan pemuat: pembaca JSON_POLIS + antarmuka
// tulis jalur biasa + dua penulis tambahan (`repository/lama.go`). Implementasi
// Oracle: `penyimpanOracle`. Antarmuka ini ada sejak Copy Old (07-10-2026):
// seam HTTP-nya diuji di atas gudang tiruan (`handlers/copyold_test.go`).
type GudangPemuat interface {
	Transaksi(ctx context.Context, fn func(tx *db.Tx) error) error
	HitungJSONPolisLain(ctx context.Context) (int, error)
	KunciJSONPolis(ctx context.Context) ([]string, error)
	// KunciJSONPolisCopyOld - popup Copy Old: hanya dokumen yang kasus Pega-nya ada dan sudah berproduksi (WO 07-10-2026).
	KunciJSONPolisCopyOld(ctx context.Context) ([]string, error)
	BacaJSONPolis(ctx context.Context, kunci string) (models.BarisJSONPolis, error)
	AdaKasus(ctx context.Context, tx *db.Tx, id string) (bool, error)
	SisipKasus(ctx context.Context, tx *db.Tx, id, pembuat, namaPembuat string) error
	SimpanHalaman(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error
	SetelNomorPolis(ctx context.Context, tx *db.Tx, id, nopol string) error
	SetelKolomDatarLama(ctx context.Context, tx *db.Tx, id string, k models.KolomDatarLama) error
	TutupKasus(ctx context.Context, tx *db.Tx, id, statusLama, statusAkhir string) error
	SalinUsulanLama(ctx context.Context, tx *db.Tx, idPega string, baris []models.UsulanProduksi) (models.NasibUsulan, error)
}

var _ GudangPemuat = penyimpanOracle{}

// Pemuat - pemuat dokumen lama di atas `GudangPemuat`.
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

// errSudahDimuat - baris kasus sudah ada dengan IDPEGA yang sama (jalankan ulang).
var errSudahDimuat = errors.New("services: dokumen sudah dimuat")

// Jalankan memuat seluruh dokumen generasi NB. `tulis` false = uji-kering:
// hanya JSON_POLIS yang dibaca, nol pernyataan ke tabel mana pun - arsip
// medan dan laporan galat tetap ditulis lengkap; salinan SuggestList dihitung
// "siap disalin" (penjaga dobel IDPEGA hanya diperiksa saat menulis).
//
// Galat yang dikembalikan hanya galat yang menghentikan seluruh jalankan
// (JSON_POLIS tidak terbaca, berkas laporan tidak dapat ditulis, dibatalkan);
// galat per dokumen masuk laporan.
func (p *Pemuat) Jalankan(ctx context.Context, tulis bool, lap *models.LaporanPemuat) error {
	lain, err := p.g.HitungJSONPolisLain(ctx)
	if err != nil {
		return err
	}
	lap.ProdKeLain(lain)
	kunci, err := p.g.KunciJSONPolis(ctx)
	if err != nil {
		return err
	}
	terlihat := map[string]string{}
	for _, k := range kunci {
		if err := ctx.Err(); err != nil {
			return err
		}
		lap.Dibaca()
		b, err := p.g.BacaJSONPolis(ctx, k)
		if err != nil {
			if err := lap.Gagal(models.BarisJSONPolis{IDPega: "ROWID " + k}, []models.GalatDokumen{{Err: err}}); err != nil {
				return err
			}
			continue
		}
		h, err := models.PecahDokumenLama(b)
		switch {
		case errors.Is(err, models.ErrBukanTreatyIn), errors.Is(err, models.ErrGenerasiEndorsemen),
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
		terlihat[h.ID] = "ROWID " + k
		usulan := models.NasibUsulanUjiKering(h)
		if tulis {
			nasib, err := p.muat(ctx, b, h)
			switch {
			case errors.Is(err, errSudahDimuat):
				lap.SudahDimuat()
				continue
			case err != nil:
				if err := lap.Gagal(b, []models.GalatDokumen{{Err: err}}); err != nil {
					return err
				}
				continue
			}
			usulan = nasib
		}
		if err := lap.Berhasil(h, usulan); err != nil {
			return err
		}
	}
	return nil
}

// muat menulis satu dokumen dalam SATU transaksi (P2, AC 45-46) dan
// menjawab nasib salinan SuggestList-nya (F3, penjaga dobel IDPEGA).
func (p *Pemuat) muat(ctx context.Context, b models.BarisJSONPolis, h models.HasilPecah) (models.NasibUsulan, error) {
	var usulan models.NasibUsulan
	err := p.g.Transaksi(ctx, func(tx *db.Tx) error {
		// IDPEGA tidak lagi disimpan (06-10-2026): ID yang sudah ada = dokumen sudah dimuat - dilewati.
		ada, err := p.g.AdaKasus(ctx, tx, h.ID)
		if err != nil {
			return err
		}
		if ada {
			return errSudahDimuat
		}
		if err := p.g.SisipKasus(ctx, tx, h.ID, "", ""); err != nil {
			return err
		}
		if err := p.g.SimpanHalaman(ctx, tx, h.ID, h.Halaman); err != nil {
			return err
		}
		if err := p.g.SetelNomorPolis(ctx, tx, h.ID, h.NoPolis); err != nil {
			return err
		}
		if err := p.g.SetelKolomDatarLama(ctx, tx, h.ID, h.Datar); err != nil {
			return err
		}
		if err := p.g.TutupKasus(ctx, tx, h.ID, models.AssignmentAdmin, models.StatusSelesai); err != nil {
			return err
		}
		usulan, err = p.g.SalinUsulanLama(ctx, tx, b.IDPega, h.Usulan)
		return err
	})
	return usulan, err
}
