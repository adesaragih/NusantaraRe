package services

// `On Retention` → `OutwardList` (paket 9, RALAT R9/P3, OQ-MPNL-09).
//
// Di Pega checkbox b47312 memanggil `GetReinsTypeOR_Life` b47488 SAAT DIUBAH -
// keempat prakondisi langkahnya PRE=false, jadi daftar diganti setiap kali,
// dicentang maupun tidak. Di sini penggantian terjadi SAAT SIMPAN bila klien
// menandai checkbox diubah (`hitungOutward`), memakai `BEGIN`/`MATURE` yang
// disimpan. Produk baru (bukan salinan) yang `IsORS` = true juga dihitung -
// satu-satunya jalan `IsORS` menjadi true di layar adalah mengubah checkbox.

import (
	"context"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/repository"
)

// GudangOutward - pembaca kontrak On Retention (bagian Gudang).
type GudangOutward interface {
	DaftarReinstypeOR(ctx context.Context, tx *db.Tx, begin, mature string) ([]models.BarisOutward, error)
}

// perluHitungOutward - kapan `OutwardList` diganti hasil `BrowseReinstypeOR_SQL`.
func perluHitungOutward(m *models.Produk, baru bool) bool {
	return m.HitungOutward || (baru && m.SalinanDari == "" && m.Umum.IsORS)
}

// hitungOutward mengganti `OutwardList` (langkah 1 hapus, 3 baca, 4.1 tambah).
func (l *Layanan) hitungOutward(ctx context.Context, tx *db.Tx, m *models.Produk) error {
	baris, err := l.gudang.DaftarReinstypeOR(ctx, tx,
		repository.TanggalKePega(m.Inward.Begin), repository.TanggalKePega(m.Inward.Mature))
	if err != nil {
		return err
	}
	m.OutwardList = baris
	return nil
}
