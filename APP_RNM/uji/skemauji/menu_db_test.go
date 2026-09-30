//go:build db

package skemauji_test

// M_NAV_MENU diadu dengan Oracle - migrasi 900 dan pembacanya (brief menu
// 30-09-2026).
//
// Yang hanya Oracle dapat buktikan: INSERT isi awal (NEXTVAL di dalam
// INSERT ... SELECT, butir membaca induknya) benar-benar berjalan, mengulangnya
// tidak menggandakan baris, dan CHECK GROUPMENU menolak golongan lain.
//
// ⛔ Melewati bila Oracle belum dikonfigurasi. Melewati bukan lulus.

import (
	"context"
	"strings"
	"testing"

	"nusantarare/inti"
	"nusantarare/inti/menu"
	"nusantarare/inti/migrasi"
	"nusantarare/uji/skemauji"
)

func TestMenuIsiAwalDariOracleIdempotenDanBerCheck(t *testing.T) {
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
		t.Fatal(err)
	}
	defer func() { _ = skemauji.Bongkar(ctx, sqlDB, skema) }()
	repo, err := skemauji.BukaRepositori()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = repo.Close() }()

	hitung := func() (kelompok, butir int) {
		t.Helper()
		baris, err := menu.NewPembaca(repo).Baca(ctx)
		if err != nil {
			t.Fatal(err)
		}
		for _, b := range baris {
			if b.IndukID == 0 {
				kelompok++
			} else {
				butir++
			}
		}
		return kelompok, butir
	}
	if k, b := hitung(); k != 20 || b != 5 {
		t.Fatalf("sesudah migrasi: %d kelompok, %d butir - mau 20 dan 5", k, b)
	}

	// Mengulang INSERT isi awal - keadaan langkah yang gagal separuh jalan lalu
	// diulang pelari - tidak menggandakan satu baris pun.
	pernyataan, err := migrasi.PernyataanLangkah(inti.SumberMigrasi(), "900_m_nav_menu.sql")
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range pernyataan {
		if strings.HasPrefix(p, "INSERT") {
			if _, err := sqlDB.ExecContext(ctx, strings.ReplaceAll(p, "{skema}", skema)); err != nil {
				t.Fatalf("INSERT ulang: %v", err)
			}
		}
	}
	if k, b := hitung(); k != 20 || b != 5 {
		t.Errorf("sesudah INSERT ulang: %d kelompok, %d butir - mau tetap 20 dan 5", k, b)
	}

	// Pohon lengkap: empat golongan, butir mewarisi golongan induknya.
	baris, err := menu.NewPembaca(repo).Baca(ctx)
	if err != nil {
		t.Fatal(err)
	}
	m := menu.Susun(baris, []string{"claimlife", "premiumlistlife", "komiteclaimlife", "treatycontractout"})
	if len(m.Golongan) != 4 {
		t.Errorf("golongan %d, mau 4", len(m.Golongan))
	}
	for _, b := range baris {
		if b.Kode == "tco-tahun" && b.Golongan != "MASTER" {
			t.Errorf("butir tco-tahun ber-GROUPMENU %q, mau warisan induknya MASTER", b.Golongan)
		}
	}

	// CHECK GROUPMENU menolak golongan di luar keempatnya.
	_, err = sqlDB.ExecContext(ctx, `INSERT INTO `+skema+`.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN) `+
		`VALUES (`+skema+`.SEQ_M_NAV_MENU.NEXTVAL, 'uji-lain', 'Uji', 'LAIN', 'uji', 1)`)
	if err == nil || !strings.Contains(err.Error(), "ORA-02290") {
		t.Errorf("GROUPMENU 'LAIN' diterima atau galatnya bukan ORA-02290: %v", err)
	}
}
