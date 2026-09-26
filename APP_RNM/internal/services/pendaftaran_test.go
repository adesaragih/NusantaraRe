package services

// Pendaftaran klaim - TANPA Oracle.
//
// Pemilik: tiket 02.
//
// Dibaca sesudah: pendaftaran.go.

import (
	"context"
	"errors"
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
func pelakuUji() Pelaku { return Pelaku{AkunID: "UJI-OPERATOR"} }

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
		if !errors.Is(err, ErrTanpaWewenang) {
			t.Errorf("pelaku %+v diterima; galatnya %v", pelaku, err)
		}
	}
	if dipanggil != 0 {
		t.Errorf("penomor dipanggil %d kali untuk pelaku anonim", dipanggil)
	}
}
