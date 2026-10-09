package tiruan

import (
	"bytes"
	"context"
	"io"
	"strconv"
	"sync"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/penyimpanan"
)

// Penyimpanan tiruan - objek di memori (services.PenyimpananBerkas); `Gagal` diisi = setiap panggilan jarak jauh
// menjawab galat itu.
type Penyimpanan struct {
	mu    sync.Mutex
	Objek map[string]penyimpanan.Objek
	Isi   map[string][]byte
	Hapus []string
	Gagal error
	nomor int
	Masuk []penyimpanan.MasukUnggah
}

// BaruPenyimpanan - penyimpanan tiruan kosong.
func BaruPenyimpanan() *Penyimpanan {
	return &Penyimpanan{Objek: map[string]penyimpanan.Objek{}, Isi: map[string][]byte{}}
}

// Unggah - objek baru UJI-IMG-n; isinya disimpan di memori sampai Catat.
func (p *Penyimpanan) Unggah(_ context.Context, m penyimpanan.MasukUnggah) (penyimpanan.Objek, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Gagal != nil {
		return penyimpanan.Objek{}, p.Gagal
	}
	p.nomor++
	id := "UJI-IMG-" + strconv.Itoa(p.nomor)
	p.Masuk = append(p.Masuk, m)
	p.Isi[id] = append([]byte(nil), m.Isi...)
	return penyimpanan.Objek{ImageID: id, FileName: m.NamaFile, URLPublic: "URL-" + id}, nil
}

// Catat - Insert_T_Storage_SQL tiruan.
func (p *Penyimpanan) Catat(_ context.Context, _ *db.Tx, o penyimpanan.Objek) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Objek[o.ImageID] = o
	return nil
}

// Buka - isi objek tercatat.
func (p *Penyimpanan) Buka(_ context.Context, id string, _ int, _ string) (io.ReadCloser, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Gagal != nil {
		return nil, p.Gagal
	}
	if _, ada := p.Objek[id]; !ada {
		return nil, penyimpanan.ErrObjekTidakAda
	}
	return io.NopCloser(bytes.NewReader(p.Isi[id])), nil
}

// Tautan - URL tiruan objek tercatat.
func (p *Penyimpanan) Tautan(_ context.Context, id string, _ int, _ string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Gagal != nil {
		return "", p.Gagal
	}
	o, ada := p.Objek[id]
	if !ada {
		return "", penyimpanan.ErrObjekTidakAda
	}
	return o.URLPublic, nil
}

// HapusObjek - DeleteGoogleStorage_Act tiruan.
func (p *Penyimpanan) HapusObjek(_ context.Context, id, _ string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Gagal != nil {
		return p.Gagal
	}
	p.Hapus = append(p.Hapus, id)
	return nil
}

// HapusCatatan - DeleteStorage_SQL tiruan.
func (p *Penyimpanan) HapusCatatan(_ context.Context, _ *db.Tx, id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.Objek, id)
	return nil
}
