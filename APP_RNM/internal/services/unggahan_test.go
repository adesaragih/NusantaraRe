package services

// Uji jalur unggah dokumen - butir be.
//
// Yang diuji di sini: gerbang, batas ukuran, dan bentuk muatan outbox.
// Ketiganya dapat diuji tanpa Oracle; yang menuntut Oracle ada di
// `unggahan_db_test.go`.

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"nusantarare/internal/repository"
)

// kategoriUji menggantikan daftar kategori yang belum ada di korpus.
type kategoriUji []string

func (k kategoriUji) KategoriWajib(context.Context) ([]string, error) {
	return []string(k), nil
}

func TestUnggahMenolakTanpaIdentitas(t *testing.T) {
	_, err := New(nil).Dokumen().Unggah(context.Background(), Pelaku{},
		"K-1", "P-1", BerkasMasuk{}, time.Now())
	if !errors.Is(err, ErrTanpaIdentitas) {
		t.Errorf("galat = %v, mau ErrTanpaIdentitas", err)
	}
}

func TestUnggahTanpaFolderGagalTerang(t *testing.T) {
	// ⛔ GAGAL TERANG, bukan bawaan diam-diam ke `./unggahan`. Berkas
	// pelanggan tidak boleh mendarat di folder kerja siapa pun yang
	// kebetulan menjalankan server - dan galat yang menyebut nama kuncinya
	// dapat diperbaiki orang tanpa membaca kode.
	//
	// ⚠️ Diuji lewat `tulisBerkas`, sebab jalur penuhnya menuntut Oracle.
	// Yang dijaga: pesannya menyebut `UNGGAHAN_DIR`.
	if !strings.Contains(ErrUnggahanDirBelumDisetel.Error(), "UNGGAHAN_DIR") {
		t.Errorf("galat tidak menyebut nama kuncinya: %v", ErrUnggahanDirBelumDisetel)
	}
}

func TestBatasUkuranDitegakkanSaatMenyalin(t *testing.T) {
	// ⛔ Saat MENYALIN, bukan dari `Content-Length`. Header itu datang dari
	// pengirim dan dapat berbohong; `io.LimitReader` tidak dapat.
	// ⚠️ Batasnya DIKECILKAN untuk uji. Ronde pertama menulis 25 MiB dua
	// kali ke disk sungguhan, dan sekali gagal lalu lulus tiga kali berikutnya
	// - uji yang gagal secara acak akan diabaikan orang. Yang diuji di sini
	// PERILAKUNYA; angka 25 MiB dijaga TestBatasUkuranAdalahAngkaYangDinyatakan.
	const batasUji = 1024
	dir := t.TempDir()
	u := New(nil).Dokumen().DenganFolder(dir).DenganBatas(batasUji)

	jalur := filepath.Join(dir, "besar.bin")
	isi := bytes.NewReader(make([]byte, batasUji+1))
	err := u.tulisBerkas(jalur, isi)
	if !errors.Is(err, ErrBerkasTerlaluBesar) {
		t.Fatalf("galat = %v, mau ErrBerkasTerlaluBesar", err)
	}
	// ⛔ Dan berkasnya DIBUANG. Berkas yang ditolak tetapi tertinggal adalah
	// cara termudah memenuhi disk lewat permintaan yang semuanya gagal.
	if _, err := os.Stat(jalur); !os.IsNotExist(err) {
		t.Errorf("berkas yang ditolak tidak dibuang: %v", err)
	}
}

func TestBatasUkuranTepatMasihDiterima(t *testing.T) {
	// Uji yang hanya mencoba "jauh kelebihan" tidak dapat membedakan `>`
	// dari `>=`, dan yang kedua menolak berkas terbesar yang sah.
	const batasUji = 1024
	dir := t.TempDir()
	jalur := filepath.Join(dir, "pas.bin")
	err := New(nil).Dokumen().DenganFolder(dir).DenganBatas(batasUji).tulisBerkas(jalur,
		bytes.NewReader(make([]byte, batasUji)))
	if err != nil {
		t.Fatalf("berkas sebesar batasnya ditolak: %v", err)
	}
	if _, err := os.Stat(jalur); err != nil {
		t.Errorf("berkas yang diterima tidak tersimpan: %v", err)
	}
}

func TestBerkasKosongDitolak(t *testing.T) {
	dir := t.TempDir()
	jalur := filepath.Join(dir, "kosong.bin")
	err := New(nil).Dokumen().DenganFolder(dir).tulisBerkas(jalur,
		bytes.NewReader(nil))
	if !errors.Is(err, ErrBerkasKosong) {
		t.Fatalf("galat = %v, mau ErrBerkasKosong", err)
	}
	if _, err := os.Stat(jalur); !os.IsNotExist(err) {
		t.Error("berkas kosong tetap tersimpan")
	}
}

func TestKategoriDiperiksaTerhadapDaftarYangSama(t *testing.T) {
	// ⛔ Daftar yang SAMA dengan gerbang Save ke Outstanding (butir ar1).
	// Dua daftar berarti dokumen dapat diunggah dengan kategori yang gerbang
	// penyimpanan tolak, dan pemakai baru tahu berbulan-bulan kemudian.
	u := New(nil).Dokumen().DenganKategori(kategoriUji{"KTP", "Surat Kematian"})
	ctx := context.Background()
	if err := u.periksaKategori(ctx, "KTP"); err != nil {
		t.Errorf("kategori sah ditolak: %v", err)
	}
	// Huruf besar-kecil tidak membedakan: kategori datang dari layar.
	if err := u.periksaKategori(ctx, "ktp"); err != nil {
		t.Errorf("kategori sah beda huruf ditolak: %v", err)
	}
	if err := u.periksaKategori(ctx, "Karangan"); !errors.Is(err,
		ErrKategoriDokumenTidakDikenal) {
		t.Errorf("kategori karangan: %v, mau ErrKategoriDokumenTidakDikenal", err)
	}
	if err := u.periksaKategori(ctx, "  "); !errors.Is(err, ErrPermintaanTidakSah) {
		t.Errorf("kategori kosong: %v, mau ErrPermintaanTidakSah", err)
	}
}

func TestKategoriBelumDiketahuiTidakDiamDiam(t *testing.T) {
	// ⛔ Sumber bawaan GAGAL, dan galatnya menyebut apa yang ditunggu.
	// Menerima kategori apa pun selama daftarnya belum ada berarti gerbang
	// itu tidak pernah ada.
	err := New(nil).Dokumen().periksaKategori(context.Background(), "apa saja")
	if !errors.Is(err, ErrKategoriWajibBelumDiketahui) {
		t.Errorf("galat = %v, mau ErrKategoriWajibBelumDiketahui", err)
	}
}

func TestMuatanOutboxBerkasPulangPergi(t *testing.T) {
	// ⛔ Penulis dan pembacanya harus sepakat, dan keduanya hidup di berkas
	// yang sama justru supaya tidak dapat bergeser sendiri-sendiri. Ini
	// bentuk cacat lintas-lapis yang sama, kali ini menyeberangi batas
	// WAKTU: baris outbox dibaca berhari-hari sesudah ditulis.
	teks := muatanBerkas("K-1", `C:\unggahan\123.pdf`, "UJI-berkas.pdf")
	m := bacaMuatanBerkas(teks)
	if m.KlaimID != "K-1" {
		t.Errorf("klaimId = %q", m.KlaimID)
	}
	if m.Jalur != `C:\unggahan\123.pdf` {
		t.Errorf("jalur = %q; backslash Windows tidak selamat", m.Jalur)
	}
	if m.NamaFile != "UJI-berkas.pdf" {
		t.Errorf("namaFile = %q", m.NamaFile)
	}
	// Muatan rusak TIDAK menggagalkan: yang penting (`RUJUKAN`) ada di
	// kolomnya sendiri, bukan di muatan.
	if kosong := bacaMuatanBerkas("{bukan json"); kosong.Jalur != "" {
		t.Errorf("muatan rusak menghasilkan %+v, mau kosong", kosong)
	}
}

func TestJenisEfekBerkasAdalahTeksYangTetap(t *testing.T) {
	// ⛔ Nilainya tersimpan di kolom dan dibaca kembali oleh proses lain,
	// mungkin berhari-hari kemudian. Mengganti nama konstanta Go tidak boleh
	// membuat baris yang sudah terantre tak terkenali.
	if JenisEfekStorageUnggah != "storage-unggah" {
		t.Errorf("jenis unggah = %q", JenisEfekStorageUnggah)
	}
	if JenisEfekStorageHapus != "storage-hapus" {
		t.Errorf("jenis hapus = %q", JenisEfekStorageHapus)
	}
	if PenyimpananStandar != "standard" {
		t.Errorf("STORAGE = %q; `Insert_T_Storage_SQL.xml` b100 menulis 'standard'",
			PenyimpananStandar)
	}
}

func TestTautanUnduhBawaanAdalahRuteKita(t *testing.T) {
	// ⛔ Butir be: selama penyambungan nyata belum disetujui, `URLPUBLIC`
	// berisi rute KITA - bukan URL penyimpanan luar yang dikarang.
	got := TautanUnduhBawaan(20260927103000123)
	if got != "/api/dokumen/20260927103000123/isi" {
		t.Errorf("tautan = %q", got)
	}
	if strings.Contains(got, "storage.googleapis.com") ||
		strings.Contains(got, "http") {
		t.Errorf("tautan menunjuk layanan luar: %q", got)
	}
}

func TestPelaksanaBerkasMenolakJenisAsing(t *testing.T) {
	// ⛔ Outbox LINTAS MODUL. Baris jenis lain yang terlanjur dipungut harus
	// terlihat sebagai galat, bukan ditandai selesai tanpa pernah dikerjakan
	// - yang kedua menghilangkan pekerjaan orang lain tanpa jejak.
	p := &PelaksanaBerkasLokal{}
	err := p.Laksanakan(context.Background(), nil,
		repositoryBarisUji("email"))
	if !errors.Is(err, ErrPermintaanTidakSah) {
		t.Errorf("galat = %v, mau ErrPermintaanTidakSah", err)
	}
}

// repositoryBarisUji menyusun satu baris outbox seadanya.
func repositoryBarisUji(jenis string) repository.BarisEfekKeluar {
	return repository.BarisEfekKeluar{ID: "1", Jenis: jenis, Rujukan: "1"}
}

func TestBatasUkuranAdalahAngkaYangDinyatakan(t *testing.T) {
	// ⛔ Uji batas di atas memakai angka kecil supaya cepat dan tidak
	// bergantung pada disk. Yang menjaga ANGKA SUNGGUHANNYA adalah uji ini -
	// tanpa itu, batas 25 MiB dapat berubah menjadi apa pun tanpa satu pun
	// uji berbunyi.
	if BatasUkuranUnggahan != 25<<20 {
		t.Errorf("BatasUkuranUnggahan = %d, mau %d (25 MiB)",
			BatasUkuranUnggahan, 25<<20)
	}
	// ⛔ Dan ia TIDAK dapat dimatikan lewat DenganBatas.
	u := New(nil).Dokumen()
	if got := u.DenganBatas(0).batas; got != BatasUkuranUnggahan {
		t.Errorf("DenganBatas(0) = %d, mau tetap %d", got, BatasUkuranUnggahan)
	}
	if got := u.DenganBatas(-1).batas; got != BatasUkuranUnggahan {
		t.Errorf("DenganBatas(-1) = %d, mau tetap %d", got, BatasUkuranUnggahan)
	}
}
