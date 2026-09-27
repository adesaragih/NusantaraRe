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

// AC tiket 02: nomor berbentuk `<awalan>K<kode bisnis>.MM.YYYY.<5 digit>`.
//
// Contoh AC-nya sendiri dipakai sebagai kasus pertama: `RNML-KL1.08.2026.00936`.
func TestBentukNomorKlaim(t *testing.T) {
	kasus := []struct {
		awalan, kodeBisnis, mmYYYY string
		urut                       int
		mau                        string
	}{
		{"RNML-", "L1", "08.2026", 936, "RNML-KL1.08.2026.00936"},
		// Satu digit tetap menjadi lima.
		{"RNML-", "L1", "01.2026", 1, "RNML-KL1.01.2026.00001"},
		// ⛔ Urut yang sudah LEBIH panjang dari lima TIDAK dipotong.
		// Memotongnya menerbitkan nomor yang bertabrakan dengan nomor lain.
		{"RNML-", "L1", "12.2025", 123456, "RNML-KL1.12.2025.123456"},
		// Awalan datang dari basis data, jadi awalan lain harus ikut terbawa.
		{"UJI-", "L9", "03.2027", 7, "UJI-KL9.03.2027.00007"},
	}
	for _, k := range kasus {
		dapat := RakitNomorKlaim(k.awalan, k.kodeBisnis, k.mmYYYY, k.urut)
		if dapat != k.mau {
			t.Errorf("RakitNomorKlaim(%q,%q,%q,%d) = %q, mau %q",
				k.awalan, k.kodeBisnis, k.mmYYYY, k.urut, dapat, k.mau)
		}
	}
}

// AC tiket 02: nomor di-commit SEBELUM pendaftaran, bukan bersamanya.
//
// ⛔ Penjaga batas transaksi. Ronde sebelumnya menomori di dalam transaksi
// pendaftaran, sehingga kuncian `SELECT … FOR UPDATE` atas baris penghitung
// dipegang sampai seluruh pendaftaran selesai - termasuk selama pembacaan
// peserta dan penulisan pohon klaimnya. Satu pendaftar yang lambat menahan
// SEMUA pendaftar lain di seluruh instalasi.
//
// ⚠️ Penjaga STATIK, dan alasannya dinyatakan: membuktikannya lewat perilaku
// menuntut Oracle sungguhan dengan dua sambungan serentak. Yang dapat dijaga
// tanpa Oracle adalah bentuknya - dua transaksi, dan penomoran di yang
// pertama. Bila itu runtuh, kunciannya kembali menahan semua orang.
func TestPenomoranDiTransaksiSendiri(t *testing.T) {
	isi, err := os.ReadFile("pendaftaran.go")
	if err != nil {
		t.Fatal(err)
	}
	tubuh := ""
	for _, fn := range strings.Split(string(isi), "\nfunc ") {
		if strings.HasPrefix(fn, "(p *Pendaftaran) Daftar(") {
			tubuh = fn
			break
		}
	}
	if tubuh == "" {
		t.Fatal("fungsi Daftar tidak ditemukan; pembacanya yang rusak")
	}
	const mauTransaksi = 2
	if n := strings.Count(tubuh, "p.svc.DalamTransaksi("); n != mauTransaksi {
		t.Fatalf("Daftar membuka %d transaksi, mau %d - satu untuk nomornya, "+
			"satu untuk pohon klaimnya", n, mauTransaksi)
	}
	// ⚠️ Letak KEDUA dihitung dari awal teks yang sama, bukan dari potongan.
	// Ronde pertama penjaga ini membandingkan indeks ke dalam POTONGAN
	// dengan indeks ke dalam teks utuh - dan menuduh kode yang benar.
	const buka = "p.svc.DalamTransaksi("
	letakNomor := strings.Index(tubuh, "NomorBerikut(")
	letakPertama := strings.Index(tubuh, buka)
	letakKedua := letakPertama + 1 +
		strings.Index(tubuh[letakPertama+1:], buka)
	if letakNomor < 0 {
		t.Fatal("Daftar tidak lagi memanggil NomorBerikut")
	}
	if letakNomor > letakKedua {
		t.Error("NomorBerikut dipanggil di transaksi KEDUA; kuncian penghitung " +
			"akan menahan pendaftar lain selama seluruh pendaftaran")
	}
}

// Butir at + au: pendaftaran mengisi TAHAP, TGL_CREATE, dan CREATE_OP.
//
// ⛔ Penjaga penyambungan. Ketiganya kolom baru, dan kolom baru yang tidak
// pernah diisi TIDAK BERBUNYI sama sekali: barisnya tersimpan, layarnya
// tampil, dan kotak masuk hanya diam-diam kosong.
func TestPendaftaranMengisiTahapDanWaktuBuat(t *testing.T) {
	isi, err := os.ReadFile("pendaftaran.go")
	if err != nil {
		t.Fatal(err)
	}
	tubuh := ""
	for _, fn := range strings.Split(string(isi), "\nfunc ") {
		if strings.HasPrefix(fn, "(p *Pendaftaran) Daftar(") {
			tubuh = fn
			break
		}
	}
	if tubuh == "" {
		t.Fatal("fungsi Daftar tidak ditemukan; pembacanya yang rusak")
	}

	// ⚠️ Jarak spasinya TIDAK ikut diperiksa: gofmt menyejajarkan medan per
	// KELOMPOK, dan komentar di antaranya memutus kelompok. Penjaga yang
	// mengunci spasi akan jatuh karena pemformatan, bukan karena kolomnya
	// hilang.
	for _, wajib := range []string{
		// Kasus BARU berada di Outstanding Claim: Assignment2 (Input
		// Register) SELESAI saat borang register disimpan.
		"models.TahapOutstanding.String()",
		"models.PeranAdminLife",
		"Tahap:",
		"PyPosition:",
		"TglCreate:",
		"CreateOp:",
	} {
		if !strings.Contains(tubuh, wajib) {
			t.Errorf("Daftar tidak menulis %q; kolomnya akan tinggal kosong "+
				"dan kotak masuk diam-diam kosong", wajib)
		}
	}

	// ⛔ SATU jam untuk seluruh pendaftaran: dua `time.Now()` terpisah dapat
	// jatuh di sisi berlawanan tengah malam, dan periode nomor ditentukan
	// TANGGAL.
	//
	// ⚠️ Komentar dibuang lebih dulu. Ronde pertama mencacah teks mentah dan
	// menghitung 2 - yang kedua adalah komentar yang MENJELASKAN aturan ini.
	// Penjaga yang menuduh penjelasannya sendiri akan dihapus orang.
	var kode []string
	for _, b := range strings.Split(tubuh, "\n") {
		if !strings.HasPrefix(strings.TrimSpace(b), "//") {
			kode = append(kode, b)
		}
	}
	if n := strings.Count(strings.Join(kode, "\n"), "time.Now()"); n != 1 {
		t.Errorf("Daftar memanggil time.Now() %d kali, mau 1", n)
	}
}
