//go:build db

package repository_test

// Peserta versi terakhir terhadap Oracle NYATA - kasus (a)-(j) `versiterakhir_test.go` atas SQL sungguhan
// `Cari` dan `AmbilUntukKlaim`. Tanpa ORACLE_DSN MELEWATI. Fixture UJI-, nol nama orang, nol nomor polis nyata.
//
// ⚠️ Tiruan bersama `uji/skemauji` (di luar folder modul ini) belum memuat `PL_NUMBER_EDM` dan `TGL_INPUT`;
// kedua kolom itu ditambahkan ke TIRUAN skema uji oleh lengkapiKolomVersi - bukan DDL tabel warisan.

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"nusantarare/modul/claimlife/backend/repository"
	"nusantarare/uji/skemauji"
)

// lengkapiKolomVersi menambah dua kolom versi ke tiruan M_LIFE_PREMIUM_DETAIL skema uji bila belum ada
// (ORA-01430 = kolom sudah ada). Salinannya di `services/pendaftaran_db_test.go` (paket uji berbeda).
func lengkapiKolomVersi(t *testing.T, ctx context.Context, sqlDB *sql.DB, skema string) {
	t.Helper()
	for _, kolom := range []string{"PL_NUMBER_EDM VARCHAR2(100)", "TGL_INPUT DATE"} {
		if _, err := sqlDB.ExecContext(ctx, `ALTER TABLE `+skema+`.M_LIFE_PREMIUM_DETAIL ADD (`+kolom+`)`); err != nil &&
			!strings.Contains(err.Error(), "ORA-01430") {
			t.Fatalf("melengkapi tiruan peserta (%s): %v", kolom, err)
		}
	}
}

func TestDBPesertaVersiTerakhir(t *testing.T) {
	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
		if !skemauji.BolehDilewati(err) {
			t.Fatalf("skema uji menolak: %v", err)
		}
		t.Skipf("lewati: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()
	ctx := context.Background()
	if err := sqlDB.PingContext(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	if err := skemauji.Pasang(ctx, sqlDB, skema); err != nil {
		t.Fatalf("memasang skema uji: %v", err)
	}
	defer func() { _ = skemauji.Bongkar(ctx, sqlDB, skema) }()
	lengkapiKolomVersi(t, ctx, sqlDB, skema)

	// sisip - satu baris fixture; ID berawalan PL (unik antarkasus).
	sisip := func(kode, pl, sertifikat, id, edm, status, tgl, nama string, i int) {
		t.Helper()
		if _, err := sqlDB.ExecContext(ctx, `INSERT INTO `+skema+`.M_LIFE_PREMIUM_DETAIL
		    (ID, PL_NUMBER, POLICY_NO, CERTIFICATE_NO, CURRENCY, EDMSTATUS, PL_NUMBER_EDM, TGL_INPUT, NAME_OF_INSURED, SUM_INSURED)
		    VALUES (:1, :2, 'UJI-POL', :3, 'IDR', :4, :5, TO_DATE(:6, 'YYYY-MM-DD HH24:MI:SS'), :7, :8)`,
			pl+"-"+id, pl, sertifikat, nilaiAtauNil(status), nilaiAtauNil(edm), nilaiAtauNil(tgl), nilaiAtauNil(nama),
			1000*(i+1)); err != nil {
			t.Fatalf("(%s) fixture: %v", kode, err)
		}
	}
	kasus := repository.KasusVersiTerakhir()
	for _, k := range kasus {
		for i, b := range k.Baris {
			sisip(k.Kode, k.PL, k.Sertifikat, b.ID, b.PLNumberEDM, b.EDMStatus, b.TglInput, b.Nama, i)
		}
	}
	// Partisi: TIGA sertifikat di satu PL - satu dibawa Old, satu di-Delete, satu NB saja. Kasus (a)-(l)
	// masing-masing satu sertifikat per PL, sehingga kunci PARTITION BY tidak teruji tanpa ini.
	const plPartisi = "UJI-PLV-P"
	for i, b := range []struct{ sertifikat, id, edm, status, nama string }{
		{"UJI-C1", "101", "", "", "UJI-P1-LAMA"}, {"UJI-C1", "102", plPartisi + "/01", "Old", "UJI-P1-BARU"},
		{"UJI-C2", "103", "", "", "UJI-P2"}, {"UJI-C2", "104", plPartisi + "/01", "Delete", "UJI-P2"},
		{"UJI-C3", "105", "", "", "UJI-P3"},
	} {
		sisip("partisi", plPartisi, b.sertifikat, b.id, b.edm, b.status, "", b.nama, i)
	}
	repoDB, err := skemauji.BukaRepositori()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = repoDB.Close() }()
	r := repository.NewPesertaPolis(repoDB)

	for _, k := range kasus {
		t.Run(k.Kode, func(t *testing.T) {
			// ⚠️ ID fixture berawalan PL (unik antarkasus) - urutan ANGKA ID sesama kasus tetap terjaga karena
			// awalannya sama panjang.
			daftar, err := r.Cari(ctx, k.PL, "", k.CariNama, 10)
			if err != nil {
				t.Fatalf("Cari: %v", err)
			}
			tampil := false
			for _, c := range daftar {
				tampil = tampil || c.NomorSertifikat == k.Sertifikat
			}
			if tampil != k.Tampil || len(daftar) > 1 {
				t.Errorf("(%s) %s: Cari tampil %v (%d baris), mau %v", k.Kode, k.Nama, tampil, len(daftar), k.Tampil)
			}
			terpilih, hidup := repository.PilihAcuan(k)
			p, err := r.AmbilUntukKlaim(ctx, k.PL, []string{k.Sertifikat})
			switch {
			case hidup && (err != nil || len(p) != 1 || p[0].SumberID != k.PL+"-"+terpilih.ID):
				t.Errorf("(%s) AmbilUntukKlaim: %v %+v, mau baris %s", k.Kode, err, p, terpilih.ID)
			case !hidup && (err == nil || !strings.Contains(err.Error(), "tidak ada atau sudah batal/delete")):
				t.Errorf("(%s) AmbilUntukKlaim harus menolak dengan pesan tidak-ditemukan: %v %+v", k.Kode, err, p)
			}
		})
	}

	t.Run("partisi", func(t *testing.T) {
		daftar, err := r.Cari(ctx, plPartisi, "", "", 10)
		if err != nil {
			t.Fatalf("Cari: %v", err)
		}
		var dapat []string
		for _, c := range daftar {
			dapat = append(dapat, c.NomorSertifikat+"="+c.NamaTertanggung)
		}
		if got, mau := strings.Join(dapat, ","), "UJI-C1=UJI-P1-BARU,UJI-C3=UJI-P3"; got != mau {
			t.Errorf("Cari %s = %s, mau %s (C2 di-Delete /01, C1 dari /01)", plPartisi, got, mau)
		}
	})

	// Tiga penampung (:1 PL dan :2 sertifikat di dalam jendela, :3 nama di luar) atas godror - kasus (j)
	// hanya menuntut NOL baris, sehingga saringan nama yang tidak pernah cocok pun lolos tanpa ini.
	t.Run("nama-dan-sertifikat", func(t *testing.T) {
		for _, c := range []struct {
			sertifikat, nama string
			tampil           bool
		}{
			{"C1", "baru", true}, {"", "baru", true}, {"C1", "lama", false}, {"C9", "baru", false}, {"C1", "", true},
		} {
			daftar, err := r.Cari(ctx, "UJI-PLV-J", c.sertifikat, c.nama, 10)
			if err != nil {
				t.Fatalf("Cari(%q, %q): %v", c.sertifikat, c.nama, err)
			}
			if tampil := len(daftar) == 1 && daftar[0].NamaTertanggung == "UJI-NAMA-BARU"; tampil != c.tampil || len(daftar) > 1 {
				t.Errorf("Cari(%q, %q) = %+v, mau tampil %v", c.sertifikat, c.nama, daftar, c.tampil)
			}
		}
	})
}

// nilaiAtauNil - teks kosong menjadi NULL (Oracle), seperti baris new business.
func nilaiAtauNil(s string) any {
	if s == "" {
		return nil
	}
	return s
}
