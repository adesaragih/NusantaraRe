package services_test

// Lampiran + efek keluar stub (paket 8, tiket 08–09). Penyimpanan berkas di-FAKE
// di balik antarmuka; yang diperiksa EFEKNYA: rekam, objek, status, galat, ulang.

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/services"
)

// berkasPalsu - penyimpanan di memori; `Gagal` membuat Kirim gagal.
type berkasPalsu struct {
	mu     sync.Mutex
	antre  map[string][]byte
	simpan map[string][]byte
	Gagal  error
	kirim  int
}

func baruBerkasPalsu() *berkasPalsu {
	return &berkasPalsu{antre: map[string][]byte{}, simpan: map[string][]byte{}}
}

func (b *berkasPalsu) SimpanAntrean(_ context.Context, id string, isi io.Reader) error {
	d, err := io.ReadAll(isi)
	if err != nil {
		return err
	}
	b.antre[id] = d
	return nil
}

func (b *berkasPalsu) BuangAntrean(_ context.Context, id string) { delete(b.antre, id) }

func (b *berkasPalsu) Kirim(_ context.Context, id string) error {
	b.kirim++
	if b.Gagal != nil {
		return b.Gagal
	}
	if _, ada := b.simpan[id]; ada {
		return nil
	}
	d, ada := b.antre[id]
	if !ada {
		return services.ErrBerkasSumberHilang
	}
	b.simpan[id] = d
	delete(b.antre, id)
	return nil
}

func (b *berkasPalsu) Buka(_ context.Context, id string) (io.ReadCloser, error) {
	d, ada := b.simpan[id]
	if !ada {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(bytes.NewReader(d)), nil
}

func (b *berkasPalsu) Hapus(_ context.Context, id string) error {
	delete(b.simpan, id)
	delete(b.antre, id)
	return nil
}

var jamLampiran = time.Date(2026, 10, 1, 13, 2, 3, 456_000_000, time.UTC)

func layananLampiran(t *testing.T) (*services.Layanan, *berkasPalsu, func() map[string]models.ObjekPenyimpanan) {
	t.Helper()
	l, g := layananMaster()
	g.Umum["100007"] = `{"ID":"100007"}`
	g.AppName = "UJI-APP"
	b := baruBerkasPalsu()
	return l.DenganJam(func() time.Time { return jamLampiran }).DenganPenyimpanan(b), b,
		func() map[string]models.ObjekPenyimpanan { return g.Objek }
}

func unggah(l *services.Layanan, nama, isi string) (models.Lampiran, error) {
	return l.UnggahLampiran(context.Background(), pelakuUji, "100007", nama, strings.NewReader(isi))
}

func TestUnggahMerekamLaluMengirimLewatStub(t *testing.T) {
	l, b, objek := layananLampiran(t)
	a, err := unggah(l, "UJI Nota.PDF", "isi pdf")
	if err != nil {
		t.Fatal(err)
	}
	if a.Category != "File" || a.FileMimeType != "pdf" || a.UserName != "UJI-PELAKU" || a.ProdukID != "100007" ||
		a.Status != models.StatusTerunggah || len(a.StorageID) != 32 {
		t.Errorf("rekam InsertAttachProdName_Sql + status: %+v", a)
	}
	o, ada := objek()[a.StorageID]
	// InsertGoogleStorage_Act 8 b1337 (Asia/Jakarta, 12 jam): Folder + "/Doc/YYYY/MM/", "yyyyMMdd-hhmmss-S - nama".
	if !ada || o.AppFolder != "Contract/Doc/2026/10/" || o.FileName != "20261001-080203-456 - UJI Nota.PDF" ||
		o.AppName != "UJI-APP" || o.DurasiDetik != 1800 {
		t.Errorf("objek T_STORAGE_IMAGE: %+v %v", o, ada)
	}
	if string(b.simpan[a.StorageID]) != "isi pdf" {
		t.Error("berkas terkirim ke penyimpanan stub")
	}
}

func TestUnggahGagalTercatatTerlihatLaluDiulangTanpaGanda(t *testing.T) {
	l, b, objek := layananLampiran(t)
	b.Gagal = errors.New("UJI penyimpanan tidak terjangkau")
	a, err := unggah(l, "UJI.pdf", "isi")
	if err != nil {
		t.Fatalf("kegagalan efek keluar TIDAK membatalkan rekam lampiran: %v", err)
	}
	if a.Status != models.StatusGagal || !strings.Contains(a.Galat, "tidak terjangkau") {
		t.Errorf("gagal tercatat dan terlihat: %+v", a)
	}
	if _, err := l.UnduhLampiran(context.Background(), pelakuUji, "100007", a.ID); !errors.Is(err, services.ErrLampiranBelumTerkirim) {
		t.Errorf("unduh sebelum terkirim: %v", err)
	}
	b.Gagal = nil
	a2, err := l.UlangiLampiran(context.Background(), pelakuUji, "100007", a.ID)
	if err != nil || a2.Status != models.StatusTerunggah {
		t.Fatalf("ulang: %+v %v", a2, err)
	}
	if _, err := l.UlangiLampiran(context.Background(), pelakuUji, "100007", a.ID); !errors.Is(err, services.ErrLampiranSudahTerkirim) {
		t.Errorf("ulang atas lampiran terkirim: %v", err)
	}
	if len(objek()) != 1 {
		t.Errorf("pengulangan tidak menggandakan objek: %d", len(objek()))
	}
}

func TestUnggahDitolakTerang(t *testing.T) {
	l, _, _ := layananLampiran(t)
	for _, k := range []struct {
		nama, isi string
		galat     error
		pesan     string
	}{
		{"", "x", services.ErrMasukanTidakSah, "Tidak ada file yg diattach"},
		{"UJI.xyzzy", "x", services.ErrMasukanTidakSah, `File type "xyzzy" is not supported`},
		{"UJI tanpa titik", "x", services.ErrMasukanTidakSah, "is not supported"},
	} {
		_, err := unggah(l, k.nama, k.isi)
		if !errors.Is(err, k.galat) || !strings.Contains(services.Pesan(err), k.pesan) {
			t.Errorf("%q: %v", k.nama, err)
		}
	}
	if _, err := unggah(l, "UJI.pdf", "a"); err != nil {
		t.Fatal(err)
	}
	if _, err := unggah(l, "uji.PDF", "b"); !errors.Is(err, services.ErrNamaLampiranSudahAda) {
		t.Errorf("nama sama tidak menimpa diam-diam: %v", err)
	}
	if _, err := l.UnggahLampiran(context.Background(), pelakuUji, "100999", "UJI.pdf", strings.NewReader("x")); !errors.Is(err, services.ErrProdukTidakAda) {
		t.Errorf("produk tidak ada: %v", err)
	}
}

func TestUnduhSatuDanSemua(t *testing.T) {
	l, b, _ := layananLampiran(t)
	a, _ := unggah(l, "UJI satu.pdf", "SATU")
	b.Gagal = errors.New("UJI gagal")
	_, _ = unggah(l, "UJI dua.pdf", "DUA")
	b.Gagal = nil
	u, err := l.UnduhLampiran(context.Background(), pelakuUji, "100007", a.ID)
	if err != nil {
		t.Fatal(err)
	}
	isi, _ := io.ReadAll(u.Isi)
	_ = u.Isi.Close()
	if string(isi) != "SATU" || u.Nama != "UJI satu.pdf" || u.Mime != "application/pdf" {
		t.Errorf("unduh satu: %q %+v", isi, u)
	}
	var buf bytes.Buffer
	n, err := l.UnduhSemuaLampiran(context.Background(), pelakuUji, "100007", &buf)
	if err != nil || n != 1 {
		t.Fatalf("unduh semua: %d %v", n, err)
	}
	z, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil || len(z.File) != 1 || z.File[0].Name != "UJI satu.pdf" {
		t.Errorf("arsip hanya lampiran terkirim: %v %v", z, err)
	}
}

func TestHapusLampiranRekamObjekDanBerkas(t *testing.T) {
	l, b, objek := layananLampiran(t)
	a, _ := unggah(l, "UJI.pdf", "x")
	delete(b.simpan, a.StorageID) // berkasnya sudah tidak ada di penyimpanan
	if err := l.HapusLampiran(context.Background(), pelakuUji, "100007", a.ID); err != nil {
		t.Fatalf("berkas yang sudah hilang tidak menggagalkan hapus rekam: %v", err)
	}
	d, _ := l.DaftarLampiran(context.Background(), pelakuUji, "100007")
	if len(d) != 0 || len(objek()) != 0 {
		t.Errorf("rekam dan objek terhapus: %+v %v", d, objek())
	}
	if err := l.HapusLampiran(context.Background(), pelakuUji, "100007", a.ID); !errors.Is(err, services.ErrLampiranTidakAda) {
		t.Errorf("hapus dua kali: %v", err)
	}
}

func TestViewOfficeOnlineStub(t *testing.T) {
	l, _, _ := layananLampiran(t)
	x, _ := unggah(l, "UJI.xlsx", "x")
	p, _ := unggah(l, "UJI.pdf", "x")
	if err := l.LihatOffice(context.Background(), pelakuUji, "100007", x.ID); !errors.Is(err, services.ErrOfficeStub) ||
		!strings.Contains(services.Pesan(err), "OQ-MPNL-11") {
		t.Errorf("xlsx: penampil luar tidak dipanggil, 503 berkalimat: %v", err)
	}
	if err := l.LihatOffice(context.Background(), pelakuUji, "100007", p.ID); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("pdf: tautan hanya untuk xls/xlsx/doc/docx/ppt/pptx (b69247): %v", err)
	}
}
