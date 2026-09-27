package services

// Pendaftaran klaim - TANPA Oracle.
//
// Pemilik: tiket 02.
//
// Dibaca sesudah: pendaftaran.go.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"nusantarare/internal/repository"
)

// permintaanUji adalah permintaan yang sah, untuk diubah satu-satu.
func permintaanUji() PermintaanDaftar {
	return PermintaanDaftar{
		NomorPremiList: "UJI-PL-1",
		NomorPolis:     "UJI-POL-0001",
		Type:           "QP",
		KodeBisnis:     "L1",
		MataUang:       "IDR",
		Sertifikat:     []string{"006"},
	}
}

// Jalur yang DITOLAK diuji satu per satu, bukan hanya jalur yang berhasil.
func TestPermintaanDaftarMenolakYangTidakLengkap(t *testing.T) {
	kasus := map[string]func(*PermintaanDaftar){
		"tanpa nomor premium list": func(p *PermintaanDaftar) { p.NomorPremiList = "" },
		"nomor premium list spasi": func(p *PermintaanDaftar) { p.NomorPremiList = "   " },
		"tanpa Type":               func(p *PermintaanDaftar) { p.Type = "" },
		"tanpa kode bisnis":        func(p *PermintaanDaftar) { p.KodeBisnis = "" },
		"nol peserta":              func(p *PermintaanDaftar) { p.Sertifikat = nil },
	}
	for nama, ubah := range kasus {
		t.Run(nama, func(t *testing.T) {
			p := permintaanUji()
			ubah(&p)
			err := p.Periksa()
			if err == nil {
				t.Fatal("permintaan tidak lengkap diterima")
			}
			if !errors.Is(err, ErrPermintaanTidakSah) {
				t.Errorf("galatnya bukan ErrPermintaanTidakSah: %v", err)
			}
		})
	}
	if err := permintaanUji().Periksa(); err != nil {
		t.Errorf("permintaan yang sah ditolak: %v", err)
	}
}

// ⛔ Penomoran yang belum diputuskan GAGAL TERANG, bukan mengarang nomor.
//
// Keputusan work owner o melarang memanggil stored procedure, sedangkan AC 2,
// 3, dan 7-11 tiket 02 masih menuntutnya. Sampai teks AC itu ditulis ulang,
// nomor karangan yang tampak benar jauh lebih berbahaya daripada galat:
// nomor klaim dibaca manusia dan dipakai di luar sistem ini.
func TestPenomorBelumDiputuskanGagalTerang(t *testing.T) {
	_, err := PenomorBelumDiputuskan{}.NomorBerikut(context.Background(), nil, "L1", time.Now())
	if err == nil {
		t.Fatal("penomoran yang belum diputuskan mengembalikan nomor")
	}
	if !errors.Is(err, ErrPenomorBelumDiputuskan) {
		t.Errorf("galatnya bukan ErrPenomorBelumDiputuskan: %v", err)
	}
	// Pesannya harus menyebut APA yang ditunggu, bukan sekadar "belum ada".
	for _, mau := range []string{"butir o", "AC 2"} {
		if !strings.Contains(err.Error(), mau) {
			t.Errorf("pesan tidak menyebut %q: %v", mau, err)
		}
	}
}

// Validasi terjadi SEBELUM nomor diambil.
//
// Nomor yang sudah terbentuk tidak dapat dikembalikan ke urutannya. Mengambil
// nomor lebih dulu berarti membuang satu nomor setiap kali borang salah isi,
// dan lubang nomor itu terlihat oleh orang di luar sistem.
func TestNomorTidakDiambilBilaPermintaanDitolak(t *testing.T) {
	dipanggil := 0
	p := &Pendaftaran{penomor: penomorPencatat{n: &dipanggil}}
	minta := permintaanUji()
	minta.Type = ""
	if _, err := p.Daftar(context.Background(), pelakuUji(), minta); err == nil {
		t.Fatal("permintaan tidak sah diterima")
	}
	if dipanggil != 0 {
		t.Errorf("penomor dipanggil %d kali padahal permintaan ditolak", dipanggil)
	}
}

// penomorPencatat menghitung berapa kali nomor diminta.
type penomorPencatat struct{ n *int }

func (p penomorPencatat) NomorBerikut(context.Context, *repository.Tx, string, time.Time) (string, error) {
	*p.n++
	return "UJI-NOMOR", nil
}

// pelakuUji adalah pelaku yang membawa identitas, seperlunya saja.
// pelakuUji berperan Input Register sejak tiket 07: pendaftaran menulis
// status Outstanding, jadi ia salah satu jalur pengubah status.
func pelakuUji() Pelaku {
	return Pelaku{AkunID: "UJI-OPERATOR", Peran: []string{PeranInputRegister}}
}

// Pendaftaran tanpa Oracle gagal terang, bukan diam.
func TestDaftarTanpaOracleGagal(t *testing.T) {
	p := &Pendaftaran{svc: New(nil), penomor: PenomorBelumDiputuskan{}}
	_, err := p.Daftar(context.Background(), pelakuUji(), permintaanUji())
	if !errors.Is(err, repository.ErrTanpaOracle) {
		t.Errorf("galatnya bukan ErrTanpaOracle: %v", err)
	}
}

// ⛔ Pendaftaran TANPA identitas pelaku ditolak - fail-closed, bukan anonim.
//
// Tanpa pagar ini, permintaan tanpa pelaku - keadaan BAWAAN saat stub mati,
// yaitu keadaan produksi - tetap menulis klaim dengan CREATE_OP_NAME kosong.
// Jejaknya hilang, dan tidak ada yang tahu siapa mendaftarkannya.
func TestDaftarTanpaPelakuDitolak(t *testing.T) {
	dipanggil := 0
	p := &Pendaftaran{svc: New(nil), penomor: penomorPencatat{n: &dipanggil}}
	for _, pelaku := range []Pelaku{{}, {AkunID: "   "}} {
		_, err := p.Daftar(context.Background(), pelaku, permintaanUji())
		if !errors.Is(err, ErrTanpaIdentitas) {
			t.Errorf("pelaku %+v diterima; galatnya %v", pelaku, err)
		}
	}
	if dipanggil != 0 {
		t.Errorf("penomor dipanggil %d kali untuk pelaku anonim", dipanggil)
	}
}

// AC tiket 02: awalan nomor klaim DI-LOOKUP, bukan konstanta.
//
// ⛔ Penjaga arah-balik. Ronde sebelumnya menanam `"RNML-"` sebagai konstanta
// Go - benar untuk lingkungan yang kebetulan dipakai saat kode ditulis, dan
// salah di mana pun `KODE_PRODUKSI` berisi awalan lain. Salahnya tidak
// terlihat: nomornya tetap terbentuk, tetap tersimpan, dan baru ketahuan
// ketika seseorang mencarinya dan tidak menemukannya.
//
// ⚠️ Yang dijaga NILAINYA, bukan nama konstantanya - mengganti nama konstanta
// tidak memperbaiki apa pun.
func TestNolAwalanNomorKlaimSebagaiLiteral(t *testing.T) {
	// Awalan produksi yang `[data DBA]` sebut untuk lini Life.
	const awalanDBA = "RNML-"
	berkas, err := filepath.Glob(filepath.Join(".", "*.go"))
	if err != nil {
		t.Fatal(err)
	}
	diperiksa := 0
	for _, nama := range berkas {
		if strings.HasSuffix(nama, "_test.go") {
			continue
		}
		isi, err := os.ReadFile(nama)
		if err != nil {
			t.Fatal(err)
		}
		diperiksa++
		for _, baris := range strings.Split(string(isi), "\n") {
			potong := strings.TrimSpace(baris)
			// Komentar boleh menyebutnya - justru di sanalah alasannya
			// dicatat. Yang dilarang literalnya di dalam kode.
			if strings.HasPrefix(potong, "//") {
				continue
			}
			if !strings.Contains(baris, `"`+awalanDBA) {
				continue
			}
			// Awalan AKSEPTASI (`RNML-A`, `RNML-AR`) dikecualikan dan
			// alasannya dinyatakan: keduanya TIDAK ada di `KODE_PRODUKSI`.
			if strings.Contains(baris, `"`+awalanDBA+"A") {
				continue
			}
			t.Errorf("%s memuat awalan nomor %q sebagai literal; ia milik "+
				"KODE_PRODUKSI dan harus di-lookup:\n\t%s", nama, awalanDBA, potong)
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol berkas diperiksa; pembacanya yang rusak")
	}
}
