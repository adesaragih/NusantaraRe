package services

// Sambungan penyimpanan dokumen klaim (11-10-2026). Transport Google belum dipasang (izin klien HTTP di penjaga inti); uji
// transport tersimpan bersama berkasnya untuk dipasang kemudian.

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/claimlife/backend/repository"
)

// Sambungan: unggah langsung lewat PelaksanaGoogle (T_STORAGE_IMAGE + T_STORAGE_ID); View/Delete memakai penyimpanan ini.
func TestUnggahanMemakaiPelaksanaGoogle(t *testing.T) {
	isi, err := os.ReadFile("unggahan.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	for _, mau := range []string{"NewPelaksanaGoogle(u.svc, u.penyimpanan).Laksanakan(ctx, tx, b)", "penyimpanan: s.penyimpananDokumen(),",
		"u.penyimpanan.Buka(ctx, o)", "u.penyimpanan.Hapus(ctx, objek)"} {
		if !strings.Contains(s, mau) {
			t.Errorf("unggahan.go tanpa %q", mau)
		}
	}
	g, err := os.ReadFile("penyimpanangoogle.go")
	if err != nil {
		t.Fatal(err)
	}
	pel := string(g)[strings.Index(string(g), "func (p *PelaksanaGoogle) Laksanakan("):]
	urut := []string{"p.baca.DokumenAda(", "p.baca.NamaAplikasiPenyimpanan(", "p.penyimpanan.Kirim(", "p.baca.SisipObjekPenyimpanan(",
		"p.baca.TandaiDokumenTerunggah("}
	at := -1
	for _, u := range urut {
		i := strings.Index(pel, u)
		if i < at {
			t.Errorf("urutan pelaksana: %s (%d < %d)", u, i, at)
		}
		at = i
	}
}

// Tanpa transport: penyimpanan nil, unggahan memakai pelaksana lokal; galat transport gagal tertutup dengan kalimat jelas.
func TestTanpaTransportGoogleGagalTertutup(t *testing.T) {
	if New(nil).penyimpananDokumen() != nil {
		t.Error("transport Google belum boleh dipasang")
	}
	if !errors.Is(ErrTransportBelumDisetujui, ErrStorageBelumSiap) || !strings.HasPrefix(PesanStorage(ErrTransportBelumDisetujui), "Storage is not configured") {
		t.Errorf("galat transport: %q", PesanStorage(ErrTransportBelumDisetujui))
	}
	saat := time.Date(2026, 10, 11, 9, 0, 0, 5_000_000, time.UTC) // 16:00:00.005 WIB
	if folder, file := namaObjekDokumen(saat, "UJI.pdf"); folder != "Claim/Doc/2026/10/" || file != "20261011-040000-5 - UJI.pdf" {
		t.Errorf("nama objek %q %q", folder, file)
	}
}

// penyimpananTiruan - tiruan antarmuka (tanpa HTTP) untuk menguji PelaksanaGoogle dan alur Unggahan kelak.
type penyimpananTiruan struct{ dikirim int }

func (p *penyimpananTiruan) Kirim(_ context.Context, o repository.ObjekPenyimpanan, _, _, _ string) (repository.ObjekPenyimpanan, error) {
	p.dikirim++
	o.URLPublic = "tiruan"
	return o, nil
}
func (p *penyimpananTiruan) Buka(context.Context, repository.ObjekPenyimpanan) (io.ReadCloser, *repository.ObjekPenyimpanan, error) {
	return io.NopCloser(strings.NewReader("ISI")), nil, nil
}
func (p *penyimpananTiruan) Hapus(context.Context, repository.ObjekPenyimpanan) error { return nil }

func TestPenyimpananTiruanMemenuhiAntarmuka(t *testing.T) {
	var p PenyimpananDokumen = &penyimpananTiruan{}
	u := New(nil).Dokumen().DenganPenyimpanan(p)
	if u.penyimpanan == nil {
		t.Error("DenganPenyimpanan tidak memasang penyimpanan")
	}
	_ = os.ErrNotExist
}
