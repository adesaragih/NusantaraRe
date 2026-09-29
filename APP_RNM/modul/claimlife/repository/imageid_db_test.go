//go:build db

package repository_test

// Rumus `IMAGEID` diadu dengan Oracle sendiri - butir be, 28-09-2026.
//
// ⛔ SEBAB UJI INI ADA. `models.ImageIDDari` menghitung MD5 di Go dan
// menuliskannya heksa huruf besar, dengan klaim bahwa itu sama dengan
// `STANDARD_HASH(x,'MD5')` di Oracle. Klaim itu tidak dapat dibuktikan oleh
// uji murni mana pun - yang murni hanya dapat membuktikan Go konsisten dengan
// dirinya sendiri. Di sinilah ia diadu dengan pihak yang sebenarnya.
//
// ⚠️ Masukannya TETAP, bukan `SYSTIMESTAMP || SYS_GUID()`. Yang diuji
// RUMUSNYA; memakai masukan yang berubah membuat kedua sisi tidak pernah
// menghitung hal yang sama, dan uji yang selalu hijau karena tidak pernah
// membandingkan apa pun.
//
// ⛔ Melewati bila Oracle belum dikonfigurasi. Melewati bukan lulus.

import (
	"context"
	"testing"

	"nusantarare/inti/unggah"
	"nusantarare/uji/skemauji"
)

func TestImageIDGoSamaDenganStandardHashOracle(t *testing.T) {
	sqlDB, _, err := skemauji.Buka()
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

	// Beberapa masukan, termasuk yang kosong dan yang panjang: hash yang
	// cocok untuk satu nilai dapat cocok secara kebetulan.
	for _, masukan := range []string{
		"ASMPP202609281345071234567890123456789ABCDEF0123456789ABCDEF",
		"ASMPP" + "20260928134507000000000" + "00000000000000000000000000000000",
		"ASMPP",
		"",
	} {
		var dariOracle string
		// ⛔ `STANDARD_HASH` dipanggil atas BIND, bukan atas literal yang
		// ditempel: teks yang ditempel akan diurai Oracle sebagai bagian
		// pernyataan, dan tanda kutip di dalamnya mengubah artinya.
		err := sqlDB.QueryRowContext(ctx,
			`SELECT STANDARD_HASH(:1, 'MD5') FROM DUAL`, masukan).Scan(&dariOracle)
		if err != nil {
			t.Fatalf("masukan %q: %v", masukan, err)
		}
		if got := unggah.ImageIDDari(masukan); got != dariOracle {
			t.Errorf("masukan %q:\n  Go     = %s\n  Oracle = %s",
				masukan, got, dariOracle)
		}
	}
}
