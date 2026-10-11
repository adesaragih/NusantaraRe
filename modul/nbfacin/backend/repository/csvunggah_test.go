package repository

import (
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestSQLUnggahCsv - tiket 75 + 74 P6: bind berurutan tanpa celah; kolom penampung (fire DAN cargo) / pekerjaan / pesan ADA di
// migrasi 257 + 258 dan setiap kolom migrasi dipakai (dua arah); T_LOCATIONLIST.SEQ_NO diperlebar; nilai bind baris penampung
// (kosong -> NULL, kurang -> NULL).
func TestSQLUnggahCsv(t *testing.T) {
	for nama, c := range map[string]struct {
		q   string
		mau int
	}{
		"buat": {sqlBuatUnggah("S.U"), 9}, "ubah": {sqlUbahUnggah("S.U"), len(kolomUnggahUbah) + 1}, "baca": {sqlBacaUnggah("S.U"), 1},
		"terakhir": {sqlUnggahTerakhir("S.U"), 1}, "sisipBaris": {sqlSisipBarisUnggah("S.B", models.JenisCsvFire), 2 + len(models.KolomCsvObjek)},
		"sisipBarisCargo": {sqlSisipBarisUnggah("S.B", models.JenisCsvCargo), 2 + len(models.KolomCsvCargo)},
		"sisipBarisAneka": {sqlSisipBarisUnggah("S.B", models.JenisCsvAneka), 2 + len(models.KolomCsvAneka)}, // tiket 73 P6
		"bacaBarisAneka":  {sqlBacaBarisUnggah("S.B", models.JenisCsvAneka), 3},
		"bacaBaris":       {sqlBacaBarisUnggah("S.B", models.JenisCsvFire), 3}, "bacaBarisCargo": {sqlBacaBarisUnggah("S.B", models.JenisCsvCargo), 3}, "sisipPesan": {sqlSisipPesanUnggah("S.P"), 5}, "bacaPesan": {sqlBacaPesanUnggah("S.P"), 3},
	} {
		if n := bindTerbesar(t, c.q); n != c.mau {
			t.Errorf("%s: %d bind, mau %d: %s", nama, n, c.mau, c.q)
		}
	}
	mig := kolomSemuaMigrasi(t)
	pakai := map[string][]string{
		TabelBarisUnggahCsv: append(append([]string{"UNGGAH_ID", "NO_BARIS"}, kolomBarisUnggah(models.JenisCsvFire)...),
			kolomBarisUnggah(models.JenisCsvCargo)...),
		// tiket 73 P6: penampung Aneka tabel sendiri (migrasi 250, K73-52 opsi 1).
		TabelBarisUnggahAneka: append([]string{"UNGGAH_ID", "NO_BARIS"}, kolomBarisUnggah(models.JenisCsvAneka)...),
		TabelUnggahCsv: append([]string{"ID", "CASE_ID", "NAMA_BERKAS", "DIBUAT_OLEH", "DIBUAT", "DIUBAH", "JENIS", "KOLOM"},
			kolomUnggahUbah...),
		TabelPesanUnggahCsv: {"UNGGAH_ID", "URUT", "JENIS", "NO_BARIS", "TEKS"},
	}
	for tabel, kolom := range pakai {
		dipakai := map[string]bool{}
		for _, k := range kolom {
			dipakai[k] = true
			if !mig[tabel][k] {
				t.Errorf("%s.%s dipakai kode, tidak ada di migrasi", tabel, k)
			}
		}
		for k := range mig[tabel] {
			if !dipakai[k] {
				t.Errorf("%s.%s di migrasi 257 / 258 / 250 tidak dipakai kode", tabel, k)
			}
		}
		if len(mig[tabel]) == 0 {
			t.Errorf("%s tidak ada di migrasi", tabel)
		}
	}
	if !strings.Contains(bacaMigrasi(t, "257_csv_objek_unggah.sql"), "ALTER TABLE {skema}.T_LOCATIONLIST MODIFY (SEQ_NO NUMBER(10))") {
		t.Error("SEQ_NO lokasi tidak diperlebar")
	}
	arg := argBarisUnggah("9", models.JenisCsvFire, []models.BarisUnggahCsv{{No: 3, Isi: []string{"1", "", "X"}}})
	if len(arg[0]) != 63 || arg[0][0] != "9" || arg[0][1] != 3 || arg[0][2] != "1" || arg[0][3] != nil || arg[0][4] != "X" || arg[0][62] != nil {
		t.Errorf("arg: %v", arg[0][:6])
	}
	if err := periksaBindBaris(sqlSisipBarisUnggah("S.B", models.JenisCsvFire), TabelBarisUnggahCsv, arg, 0); err != nil {
		t.Errorf("periksaBindBaris: %v", err)
	}
	argC := argBarisUnggah("9", models.JenisCsvCargo, []models.BarisUnggahCsv{{No: 1, Isi: []string{"DOMESTIK"}}})
	if len(argC[0]) != 20 || argC[0][2] != "DOMESTIK" || argC[0][19] != nil {
		t.Errorf("arg cargo: %v", argC[0])
	}
	if q := sqlSisipBarisUnggah("S.B", models.JenisCsvCargo); !strings.Contains(q, "C_TRADING, C_DEPARTURE") || !strings.Contains(q, "C_TOP_RISK)") {
		t.Error(q)
	}
	if jenisUnggah("") != models.JenisCsvFire || jenisUnggah(models.JenisCsvCargo) != models.JenisCsvCargo {
		t.Error("jenisUnggah")
	}
	if potong4000(strings.Repeat("é", 2500)) != strings.Repeat("é", 2000) {
		t.Error("potong4000 memotong di tengah rune")
	}
}
