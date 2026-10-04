package services

// Untuk apa berkas ini: PEMUAT DOKUMEN LAMA - tiket 22 (spec-penyimpanan
// ID-3, AC 55-59; KEPUTUSAN-RONDE-12 butir 5; K15, K17).
//
// Dijalankan MANUSIA lewat perintah `backend/alat/pemuatlama` (MODUL.md
// "Pemuat dokumen lama"), tidak pernah saat aplikasi menyala - modul.go tidak
// mendaftarkannya. Seluruh dokumen generasi NB di `POOLDATA.JSON_POLIS`
// (PRODKE 0) - setiap polis, tanpa penyaring status/tahun/lini usaha (butir 5).
// Generasi endorsemen (PRODKE > 0) milik pemuat EDM (edmtreatyin tiket 10).
//
// Per dokumen: pecah (models, murni) -> SATU transaksi -> tulis lewat
// antarmuka penyimpanan yang SAMA dengan jalur biasa (ID-3, AC 56):
//
//	SisipKasus      T_WORK_POLIS + generasi PRODKE 0 (pembuat tidak dikarang: NULL)
//	SimpanHalaman   T_GENERAL_POLIS + T_POLIS_QUOTATION/CEDING/INSTALMENT(_DETAIL)/
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

// Pemuat - pemuat dokumen lama. Penyimpanannya langsung `repository.Gudang`
// + transaksi `inti.Dasar` (`penyimpanOracle`): satu implementasi, tanpa
// tiruan - seam services tidak diuji (spec.md §6.2), jadi tidak perlu
// antarmuka tersendiri.
type Pemuat struct{ g penyimpanOracle }

// PemuatDariDasar menyusun pemuat di atas basis data bersama.
func PemuatDariDasar(d *inti.Dasar) (*Pemuat, error) {
	if d == nil || !d.PunyaDatabase() {
		return nil, ErrTanpaOracle
	}
	return &Pemuat{g: penyimpanOracle{Gudang: repository.Baru(d.DB()), dasar: d}}, nil
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
		case errors.Is(err, models.ErrBukanTreatyIn), errors.Is(err, models.ErrGenerasiEndorsemen):
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
		usulanDisalin := true // uji-kering: siap disalin
		if tulis {
			disalin, err := p.muat(ctx, b, h)
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
			usulanDisalin = disalin
		}
		if err := lap.Berhasil(h, usulanDisalin); err != nil {
			return err
		}
	}
	return nil
}

// muat menulis satu dokumen dalam SATU transaksi (P2, AC 45-46). Hasil
// pertama false = salinan SuggestList dilewati penjaga dobel IDPEGA (F3).
func (p *Pemuat) muat(ctx context.Context, b models.BarisJSONPolis, h models.HasilPecah) (bool, error) {
	usulanDisalin := false
	err := p.g.Transaksi(ctx, func(tx *db.Tx) error {
		ada, err := p.g.IDPegaKasus(ctx, tx, h.ID)
		switch {
		case err == nil && ada == b.IDPega:
			return errSudahDimuat
		case err == nil:
			return fmt.Errorf("%w: %s ber-IDPEGA %q", models.ErrIDKasusDipakai, h.ID, ada)
		case !errors.Is(err, repository.ErrKasusTidakAda):
			return err
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
		usulanDisalin, err = p.g.SalinUsulanLama(ctx, tx, b.IDPega, h.Usulan)
		return err
	})
	return usulanDisalin, err
}
