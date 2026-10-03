package services_test

// Lampiran + efek keluar stub (paket 8, tiket 08–09). Penyimpanan berkas di-FAKE
// di balik antarmuka; yang diperiksa EFEKNYA: rekam, objek, status, galat, ulang.

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"testing/iotest"
	"time"

	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/services"
)

// berkasPalsu - penyimpanan di memori; `Gagal` membuat Kirim gagal, `BukaGagal` membuat Buka gagal, `Segar`
// = objek baru yang Buka kembalikan (URL bertanda tangan diperbarui, `GetUrlGoogleStorage_Act`).
type berkasPalsu struct {
	mu        sync.Mutex
	antre     map[string][]byte
	simpan    map[string][]byte
	Gagal     error
	BukaGagal error
	Segar     *models.ObjekPenyimpanan
	kirim     int
	// ext, mime - yang diserahkan Kirim terakhir; dihapus - objek yang diserahkan Hapus (nil = belum terkirim);
	// nama - Namafile setiap percobaan Kirim; BacaGagal - isi yang dibuka gagal dibaca di tengah.
	ext, mime string
	dihapus   []*models.ObjekPenyimpanan
	nama      []string
	BacaGagal error
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

// Kirim meniru penyimpanan nyata: antrean TIDAK dibuang di sini - layanan membuangnya sesudah objek tercatat.
func (b *berkasPalsu) Kirim(_ context.Context, o models.ObjekPenyimpanan, ext, mime string) (models.ObjekPenyimpanan, error) {
	b.kirim++
	b.ext, b.mime = ext, mime
	b.nama = append(b.nama, o.AppFolder+o.FileName)
	if b.Gagal != nil {
		return models.ObjekPenyimpanan{}, b.Gagal
	}
	d, ada := b.antre[o.ImageID]
	if !ada {
		return models.ObjekPenyimpanan{}, services.ErrBerkasSumberHilang
	}
	b.simpan[o.ImageID] = d
	o.URLPublic, o.Exp = "UJI-URL-"+o.ImageID, "01/10/2026 13:32:03"
	return o, nil
}

func (b *berkasPalsu) Buka(_ context.Context, o models.ObjekPenyimpanan) (io.ReadCloser, *models.ObjekPenyimpanan, error) {
	if b.BukaGagal != nil {
		return nil, nil, b.BukaGagal
	}
	d, ada := b.simpan[o.ImageID]
	if !ada {
		return nil, nil, os.ErrNotExist
	}
	if b.BacaGagal != nil {
		return io.NopCloser(io.MultiReader(bytes.NewReader(d), iotest.ErrReader(b.BacaGagal))), b.Segar, nil
	}
	return io.NopCloser(bytes.NewReader(d)), b.Segar, nil
}

func (b *berkasPalsu) Hapus(_ context.Context, id string, o *models.ObjekPenyimpanan) error {
	b.dihapus = append(b.dihapus, o)
	delete(b.simpan, id)
	delete(b.antre, id)
	return nil
}

var jamLampiran = time.Date(2026, 10, 1, 13, 2, 3, 456_000_000, time.UTC)

func layananLampiran(t *testing.T) (*services.Layanan, *berkasPalsu, func() map[string]models.ObjekPenyimpanan) {
	t.Helper()
	l, g := layananMaster()
	g.IsiJSON("100007", `{"ID":"100007"}`, "")
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
	// InsertGoogleStorage_Act 8 b1339 (Asia/Jakarta, 12 jam): Folder + "/Doc/YYYY/MM/", "yyyyMMdd-hhmmss-S - nama" -
	// waktunya waktu REKAM (ID lampiran tiruan `20261001000000001`), sama di setiap percobaan; Durasi 1800
	// (`ProductNameSaveAttachment` 2.4 b904); yang DICATAT = objek jawaban penyimpanan (URLImage, exp).
	if !ada || o.AppFolder != "Contract/Doc/2026/10/" || o.FileName != "20261001-120000-1 - UJI Nota.PDF" ||
		o.AppName != "UJI-APP" || o.DurasiDetik != 1800 || o.URLPublic != "UJI-URL-"+a.StorageID || o.Exp != "01/10/2026 13:32:03" {
		t.Errorf("objek T_STORAGE_IMAGE: %+v %v", o, ada)
	}
	if string(b.simpan[a.StorageID]) != "isi pdf" || b.ext != "pdf" || b.mime != "application/pdf" {
		t.Errorf("berkas terkirim ke penyimpanan: ext %q mime %q", b.ext, b.mime)
	}
	if len(b.antre) != 0 {
		t.Error("antrean lokal dibuang sesudah objek tercatat")
	}
}

func TestUnggahGagalTercatatTerlihatLaluDiulangTanpaGanda(t *testing.T) {
	l, b, objek := layananLampiran(t)
	b.Gagal = errors.New("UJI penyimpanan tidak terjangkau")
	a, err := unggah(l, "UJI.pdf", "isi")
	if len(b.antre) != 1 {
		t.Error("antrean lokal DITAHAN selama belum terkirim - kirim ulang memakainya")
	}
	if err != nil {
		t.Fatalf("kegagalan efek keluar TIDAK membatalkan rekam lampiran: %v", err)
	}
	// Audit 02-10-2026: yang tersimpan dan tampil kalimat tetap; teks galat aslinya (bisa berisi ORA- atau jalur
	// folder stub) hanya di log.
	if a.Status != models.StatusGagal || a.Galat != services.PesanKirimGagal || strings.Contains(a.Galat, "tidak terjangkau") {
		t.Errorf("gagal tercatat dan terlihat dengan kalimat tetap: %+v", a)
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
	// `DeleteGoogleStorage_Act` 4 b639: objek dibaca dari T_STORAGE_IMAGE (GetLinkStorage_SQL) dan diserahkan utuh.
	if len(b.dihapus) != 1 || b.dihapus[0] == nil || b.dihapus[0].ImageID != a.StorageID || b.dihapus[0].URLPublic != "UJI-URL-"+a.StorageID {
		t.Errorf("objek terkirim diserahkan ke hapus: %+v", b.dihapus)
	}
	d, _ := l.DaftarLampiran(context.Background(), pelakuUji, "100007")
	if len(d) != 0 || len(objek()) != 0 {
		t.Errorf("rekam dan objek terhapus: %+v %v", d, objek())
	}
	if err := l.HapusLampiran(context.Background(), pelakuUji, "100007", a.ID); !errors.Is(err, services.ErrLampiranTidakAda) {
		t.Errorf("hapus dua kali: %v", err)
	}
}

func TestHapusLampiranBelumTerkirimTanpaObjek(t *testing.T) {
	l, b, _ := layananLampiran(t)
	b.Gagal = errors.New("UJI gagal")
	a, _ := unggah(l, "UJI.pdf", "x")
	if err := l.HapusLampiran(context.Background(), pelakuUji, "100007", a.ID); err != nil {
		t.Fatal(err)
	}
	if len(b.dihapus) != 1 || b.dihapus[0] != nil || len(b.antre) != 0 {
		t.Errorf("belum terkirim: hanya antrean yang dibuang: %+v %d", b.dihapus, len(b.antre))
	}
}

// `GetUrlGoogleStorage_Act` 12 b2375 (`Update_T_Storage_SQL`): URL bertanda tangan yang diperbarui DICATAT.
func TestUnduhMencatatURLYangDiperbarui(t *testing.T) {
	l, b, objek := layananLampiran(t)
	a, _ := unggah(l, "UJI.pdf", "isi")
	segar := objek()[a.StorageID]
	segar.URLPublic, segar.Exp, segar.TanggalUpload = "UJI-URL-BARU", "01/10/2026 14:00:00", "10/01/2026 13:30:00"
	b.Segar = &segar
	u, err := l.UnduhLampiran(context.Background(), pelakuUji, "100007", a.ID)
	if err != nil {
		t.Fatal(err)
	}
	_ = u.Isi.Close()
	if o := objek()[a.StorageID]; o.URLPublic != "UJI-URL-BARU" || o.Exp != "01/10/2026 14:00:00" || o.TanggalUpload != "10/01/2026 13:30:00" {
		t.Errorf("objek diperbarui: %+v", o)
	}
}

// Galat penyimpanan nyata diteruskan dengan jenisnya (502 / 503 / 409), tidak disamaratakan menjadi "tidak di stub";
// galat kirim bernama disimpan dengan kalimatnya sendiri.
func TestGalatPenyimpananNyataDiteruskan(t *testing.T) {
	l, b, _ := layananLampiran(t)
	b.Gagal = fmt.Errorf("%w: status 401", services.ErrStorageGagal)
	a, _ := unggah(l, "UJI.pdf", "isi")
	if a.Status != models.StatusGagal || !strings.Contains(a.Galat, "storage service") || !strings.Contains(a.Galat, "401") {
		t.Errorf("galat kirim bernama tersimpan berkalimat: %+v", a)
	}
	b.Gagal = nil
	a, _ = l.UlangiLampiran(context.Background(), pelakuUji, "100007", a.ID)
	for _, g := range []error{services.ErrStorageGagal, services.ErrStorageBelumSiap, services.ErrBerkasTidakDiStorage} {
		b.BukaGagal = fmt.Errorf("%w: UJI", g)
		if _, err := l.UnduhLampiran(context.Background(), pelakuUji, "100007", a.ID); !errors.Is(err, g) ||
			errors.Is(err, services.ErrBerkasTidakDiStub) {
			t.Errorf("%v diteruskan: %v", g, err)
		}
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
		t.Errorf("pdf: tautan hanya untuk xls/xlsx/doc/docx/ppt/pptx (b69291): %v", err)
	}
}

// Audit 02-10-2026: unduhan yang gagal tidak membawa jalur folder stub atau kunci objek ke layar, dan folder
// stub yang belum disetel tetap 503 (`ErrPenyimpananBelumDisetel`), bukan 409 "tidak di stub".
func TestUnduhGagalTanpaJalurServerDanBelumDisetelTetap503(t *testing.T) {
	l, b, _ := layananLampiran(t)
	a, err := unggah(l, "UJI.pdf", "isi")
	if err != nil || a.Status != models.StatusTerunggah {
		t.Fatalf("unggah: %+v %v", a, err)
	}
	b.BukaGagal = errors.New("open UJI-DIR/master-product-name-life/simpan/" + a.StorageID + ": not found")
	_, err = l.UnduhLampiran(context.Background(), pelakuUji, "100007", a.ID)
	if !errors.Is(err, services.ErrBerkasTidakDiStub) || strings.Contains(services.Pesan(err), "UJI-DIR") ||
		strings.Contains(services.Pesan(err), a.StorageID) || !strings.Contains(services.Pesan(err), "UJI.pdf") {
		t.Errorf("kalimat layar tanpa jalur server, menyebut nama berkas: %q %v", services.Pesan(err), err)
	}
	b.BukaGagal = services.ErrPenyimpananBelumDisetel
	if _, err = l.UnduhLampiran(context.Background(), pelakuUji, "100007", a.ID); !errors.Is(err, services.ErrPenyimpananBelumDisetel) ||
		errors.Is(err, services.ErrBerkasTidakDiStub) {
		t.Errorf("folder stub belum disetel = 503, bukan 409: %v", err)
	}
}

// Kirim ulang menulis ke objek yang SAMA (Folder + Namafile dari waktu rekam), bukan objek baru yang meninggalkan objek
// yatim di penyimpanan bila percobaan sebelumnya sudah sampai tetapi pencatatannya gagal.
func TestKirimUlangNamaObjekSama(t *testing.T) {
	l, b, _ := layananLampiran(t)
	b.Gagal = errors.New("UJI gagal")
	a, _ := unggah(l, "UJI.pdf", "x")
	b.Gagal = nil
	l = l.DenganJam(func() time.Time { return jamLampiran.Add(48 * time.Hour) })
	if _, err := l.UlangiLampiran(context.Background(), pelakuUji, "100007", a.ID); err != nil {
		t.Fatal(err)
	}
	if len(b.nama) != 2 || b.nama[0] != b.nama[1] {
		t.Errorf("Namafile setiap percobaan: %q", b.nama)
	}
}

// Isi yang putus di tengah arsip = galat penyimpanan (502), tanpa teks jaringan mentah (alamat) di layar.
func TestUnduhSemuaPutusDiTengah(t *testing.T) {
	l, b, _ := layananLampiran(t)
	_, _ = unggah(l, "UJI.pdf", "isi")
	b.BacaGagal = errors.New("read tcp UJI-ALAMAT: connection reset")
	var buf bytes.Buffer
	_, err := l.UnduhSemuaLampiran(context.Background(), pelakuUji, "100007", &buf)
	if !errors.Is(err, services.ErrStorageGagal) || strings.Contains(services.Pesan(err), "UJI-ALAMAT") {
		t.Errorf("putus: %v / %q", err, services.Pesan(err))
	}
}
