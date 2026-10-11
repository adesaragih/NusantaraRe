package tiruan

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sort"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/repository"
)

// KategoriLampiran - lihat `repository.Gudang.KategoriLampiran` (master FAC x cacah dokumen klaim berkategori itu).
func (g *Gudang) KategoriLampiran(_ context.Context, id string) ([]models.KategoriLampiran, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := []models.KategoriLampiran{} // kosong = [] (bukan null) di JSON
	for k, label := range g.KategoriDok {
		n := 0
		for _, d := range g.Dokumen {
			if d.Kategori1 == k && (d.IDPega == models.KunciInstans(id) || d.IDPega == models.KunciPegaLama(id)) {
				n++
			}
		}
		out = append(out, models.KategoriLampiran{ID: k, Label: label, CountAttach: n})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

// DaftarLampiran - lihat `repository.Gudang.DaftarLampiran`.
func (g *Gudang) DaftarLampiran(_ context.Context, id string) ([]models.Lampiran, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := []models.Lampiran{} // kosong = [] (bukan null) di JSON
	for i := len(g.Dokumen) - 1; i >= 0; i-- {
		d := g.Dokumen[i]
		if d.IDPega == models.KunciInstans(id) || d.IDPega == models.KunciPegaLama(id) {
			out = append(out, models.Lampiran{ID: d.ID, NamaFile: d.NamaFile, Kategori: d.Kategori1, MIME: d.MIME,
				Tanggal: d.Tanggal.In(models.Jakarta).Format("02/01/2006 15:04"), Operator: d.Operator,
				StorageID: d.StorageID, AdaObjek: d.StorageID != ""})
		}
	}
	return out, nil
}

// SisipDokumenKlaim - dokumen klaim tiruan; ID kembar = `repository.ErrIDDokumenTerpakai` (PK).
func (g *Gudang) SisipDokumenKlaim(_ context.Context, _ *db.Tx, d models.BarisDokumenKlaim) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, x := range g.Dokumen {
		if x.ID == d.ID {
			return repository.ErrIDDokumenTerpakai
		}
	}
	g.Dokumen = append(g.Dokumen, d)
	return nil
}

// HapusDokumenKlaim - lihat `repository.Gudang.HapusDokumenKlaim`.
func (g *Gudang) HapusDokumenKlaim(_ context.Context, _ *db.Tx, id, lid string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.GagalHapusDok != nil {
		return g.GagalHapusDok
	}
	for i, d := range g.Dokumen {
		if d.ID == lid && (d.IDPega == models.KunciInstans(id) || d.IDPega == models.KunciPegaLama(id)) {
			g.Dokumen = append(g.Dokumen[:i:i], g.Dokumen[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("tiruan: dokumen klaim %s tidak ada", lid)
}

// PindahKategoriDokumen - lihat `repository.Gudang.PindahKategoriDokumen`.
func (g *Gudang) PindahKategoriDokumen(_ context.Context, _ *db.Tx, id, lid, kategori string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	for i, d := range g.Dokumen {
		if d.ID == lid && (d.IDPega == models.KunciInstans(id) || d.IDPega == models.KunciPegaLama(id)) {
			g.Dokumen[i].Kategori1 = kategori
			return nil
		}
	}
	return fmt.Errorf("tiruan: dokumen klaim %s tidak ada", lid)
}

// Berkas - penyimpanan berkas tiruan (`services.PenyimpananBerkas`); catatan objeknya ikut transaksi gudang.
type Berkas struct {
	g *Gudang
	// Gagal - galat Unggah (layanan penyimpanan tak terjangkau / berkas ditolak).
	Gagal    error
	Unggahan []penyimpanan.MasukUnggah
	// Dihapus - IMAGEID objek yang dihapus dari penyimpanan.
	Dihapus []string
	// GagalHapus - galat HapusObjek (layanan penyimpanan menolak penghapusan).
	GagalHapus error
}

// BerkasBaru membuat penyimpanan tiruan di atas gudang ini.
func (g *Gudang) BerkasBaru() *Berkas { return &Berkas{g: g} }

// Unggah - InsertGoogleStorage_Act tiruan: IMAGEID "UJI-IMG-n".
func (b *Berkas) Unggah(_ context.Context, m penyimpanan.MasukUnggah) (penyimpanan.Objek, error) {
	if b.Gagal != nil {
		return penyimpanan.Objek{}, b.Gagal
	}
	b.Unggahan = append(b.Unggahan, m)
	return penyimpanan.Objek{ImageID: fmt.Sprintf("UJI-IMG-%d", len(b.Unggahan)), FileName: m.NamaFile}, nil
}

// Buka - GetUrlGoogleStorage_Act tiruan: isi berkas yang diunggah dengan IMAGEID itu.
func (b *Berkas) Buka(_ context.Context, imageID string, _ int, _ string) (io.ReadCloser, error) {
	for i, m := range b.Unggahan {
		if fmt.Sprintf("UJI-IMG-%d", i+1) == imageID {
			return io.NopCloser(bytes.NewReader(m.Isi)), nil
		}
	}
	return nil, penyimpanan.ErrObjekTidakAda
}

// Tautan - GetUrlGoogleStorage_Act tiruan: URL bertanda tangan palsu.
func (b *Berkas) Tautan(_ context.Context, imageID string, _ int, _ string) (string, error) {
	return "UJI-URL-" + imageID, nil
}

// HapusObjek - DeleteGoogleStorage_Act tiruan.
func (b *Berkas) HapusObjek(_ context.Context, imageID, _ string) error {
	if b.GagalHapus != nil {
		return b.GagalHapus
	}
	b.Dihapus = append(b.Dihapus, imageID)
	return nil
}

// HapusCatatan - DeleteStorage_SQL tiruan.
func (b *Berkas) HapusCatatan(_ context.Context, _ *db.Tx, imageID string) error {
	b.g.mu.Lock()
	defer b.g.mu.Unlock()
	for i, o := range b.g.Storage {
		if o.ImageID == imageID {
			b.g.Storage = append(b.g.Storage[:i:i], b.g.Storage[i+1:]...)
			return nil
		}
	}
	return nil
}

// Catat - Insert_T_Storage_SQL tiruan.
func (b *Berkas) Catat(_ context.Context, _ *db.Tx, o penyimpanan.Objek) error {
	b.g.mu.Lock()
	defer b.g.mu.Unlock()
	b.g.Storage = append(b.g.Storage, o)
	return nil
}
