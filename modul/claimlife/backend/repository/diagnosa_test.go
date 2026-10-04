package repository

// Uji bentuk query diagnosa - butir bd.
//
// Nol Oracle: yang diperiksa TEKS query-nya. Kolom yang lupa di-bind dan
// kolom yang ikut tertulis padahal tidak boleh tidak terlihat dari daftar
// nama mana pun - hanya dari pernyataannya.

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

const (
	tabelUjiDiagnosa = "SKEMAUJI.T_CLAIMLF_DIAGNOSE"
	tabelUjiPesertaD = "SKEMAUJI.T_CLAIMLF_PREMIUMLIST_DETAIL"
)

func TestBacaDiagnosaDibatasiKlaimnya(t *testing.T) {
	q := sqlAmbilDiagnosa(tabelUjiDiagnosa, tabelUjiPesertaD)
	// ⛔ Dibatasi lewat PESERTA-nya, bukan lewat kolom klaim pada diagnosa.
	// Tabel diagnosa sengaja TIDAK punya CLAIM_ID: satu-satunya pemiliknya
	// peserta (`SetDisease.xml` b389 `Obj-Save pyWorkPage`), dan kolom kedua
	// yang menyebut klaim adalah kolom kedua yang dapat berbeda.
	if !strings.Contains(q, "p.CLAIM_ID = :1") {
		t.Errorf("pembacaan tidak dibatasi klaim:\n%s", q)
	}
	if strings.Contains(q, "d.CLAIM_ID") {
		t.Errorf("diagnosa dibaca lewat kolom klaimnya sendiri; "+
			"tabelnya tidak punya kolom itu:\n%s", q)
	}
	// ⛔ Urutan STABIL, dan dua kuncinya. `URUTAN` yang dilihat pemakai,
	// `ID` yang memutus seri - daftar tanpa pemutus seri dapat tampil
	// berbeda pada dua pembacaan yang sama.
	if !strings.Contains(q, "ORDER BY d.PREMIUM_LIST_DETAIL_ID, d.URUTAN, d.ID") {
		t.Errorf("urutan baca tidak stabil:\n%s", q)
	}
	for _, kolom := range []string{
		"d.ICD_CODE", "d.DISEASE", "d.GROUP_DIAGNOSE", "d.STS_REJECT",
	} {
		if !strings.Contains(q, kolom) {
			t.Errorf("kolom %s tidak dibaca:\n%s", kolom, q)
		}
	}
}

func TestSisipDiagnosaMenulisBarisKOSONG(t *testing.T) {
	q := sqlSisipDiagnosa(tabelUjiDiagnosa)
	// ⛔ `Add` b4700 `addRow` menambahkan anggota PageList yang propertinya
	// BELUM terisi; pemakai lalu menekan `Find Disease`. Penyisipan yang
	// menuntut ICD dan nama membalik urutan kerja orang.
	for _, kolom := range []string{"ICD_CODE", "DISEASE", "GROUP_DIAGNOSE"} {
		if strings.Contains(q, kolom) {
			t.Errorf("penyisipan mengisi %s; `addRow` menambah baris kosong:\n%s",
				kolom, q)
		}
	}
	// ⚠️ `STS_REJECT` justru IKUT sejak lahir - lihat komentar
	// sqlSisipDiagnosa. Tanpa itu baris yang lahir sesudah keputusan peserta
	// tidak akan pernah disentuh pencerminan, dan akan tampak belum diputus.
	if !strings.Contains(q, "STS_REJECT") {
		t.Errorf("penyisipan tidak mewarisi STS_REJECT peserta:\n%s", q)
	}
	if !strings.Contains(q, "URUTAN") {
		t.Errorf("penyisipan tanpa URUTAN:\n%s", q)
	}
}

func TestUbahDiagnosaTidakMenyentuhUrutanMaupunKeputusan(t *testing.T) {
	q := sqlPerbaruiDiagnosa(tabelUjiDiagnosa)
	// ⛔ INILAH penjaga yang paling berarti di berkas ini. `Choose` menulis
	// isi baris; ia TIDAK boleh menyentuh `STS_REJECT` - rute yang diam-diam
	// dapat menulisnya adalah rute yang membatalkan keputusan peserta tanpa
	// ada yang memintanya, dan tidak ada satu pun layar yang menampilkannya.
	if strings.Contains(q, "STS_REJECT") {
		t.Errorf("pembaruan isi menyentuh STS_REJECT:\n%s", q)
	}
	if strings.Contains(q, "URUTAN") {
		t.Errorf("pembaruan isi menyentuh URUTAN; itu milik Add/Delete:\n%s", q)
	}
	for _, kolom := range []string{"ICD_CODE", "DISEASE", "GROUP_DIAGNOSE"} {
		if !strings.Contains(q, kolom) {
			t.Errorf("kolom %s tidak ditulis:\n%s", kolom, q)
		}
	}
	if !strings.Contains(q, "WHERE ID = :4") {
		t.Errorf("pembaruan tanpa syarat baris:\n%s", q)
	}
}

func TestRapatkanUrutanHanyaMilikPesertaItu(t *testing.T) {
	q := sqlRapatkanUrutan(tabelUjiDiagnosa)
	// ⛔ DUA syarat, dan keduanya wajib. Tanpa `PREMIUM_LIST_DETAIL_ID`
	// perapatan menggeser urutan diagnosa SELURUH peserta di basis data -
	// satu penghapusan merusak setiap klaim yang pernah ada.
	if !strings.Contains(q, "PREMIUM_LIST_DETAIL_ID = :1") {
		t.Errorf("perapatan tidak dibatasi pesertanya:\n%s", q)
	}
	if !strings.Contains(q, "URUTAN > :2") {
		t.Errorf("perapatan tidak dibatasi baris sesudah yang dihapus:\n%s", q)
	}
	if !strings.Contains(q, "URUTAN = URUTAN - 1") {
		t.Errorf("perapatan bukan pengurangan satu:\n%s", q)
	}
}

func TestCerminkanStsRejectMenyapuSeluruhDaftar(t *testing.T) {
	q := sqlCerminkanStsReject(tabelUjiDiagnosa)
	// ⛔ SATU pernyataan untuk seluruh daftar peserta itu - padanan putaran
	// `EMBEDDED` b345. Syaratnya pesertanya, BUKAN satu baris: `SetSTS_Reject`
	// b241 memutari seluruh `.DiagnoseList`, dan separuh daftar yang
	// tercermin adalah keadaan yang tidak pernah ada di sistem lama.
	if !strings.Contains(q, "WHERE PREMIUM_LIST_DETAIL_ID = :2") {
		t.Errorf("pencerminan tidak dibatasi pesertanya:\n%s", q)
	}
	if strings.Contains(q, "AND ID") {
		t.Errorf("pencerminan dibatasi satu baris; rule memutari seluruh daftar:\n%s", q)
	}
	if !strings.Contains(q, "SET STS_REJECT = :1") {
		t.Errorf("pencerminan tidak menulis STS_REJECT:\n%s", q)
	}
}

func TestQueryDiagnosaMemakaiBindBukanTempelan(t *testing.T) {
	// Seluruh pernyataan berkas ini sekaligus: nol nilai yang ditempel ke
	// teks, nol pemisah pernyataan, nol komentar SQL.
	for nama, q := range map[string]string{
		"ambil":    sqlAmbilDiagnosa(tabelUjiDiagnosa, tabelUjiPesertaD),
		"sisip":    sqlSisipDiagnosa(tabelUjiDiagnosa),
		"perbarui": sqlPerbaruiDiagnosa(tabelUjiDiagnosa),
		"rapatkan": sqlRapatkanUrutan(tabelUjiDiagnosa),
		"cermin":   sqlCerminkanStsReject(tabelUjiDiagnosa),
	} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Errorf("%s: %v\n%s", nama, err, q)
		}
		if !strings.Contains(q, ":1") {
			t.Errorf("%s: nol bind:\n%s", nama, q)
		}
		if strings.Contains(q, ";") || strings.Contains(q, "--") {
			t.Errorf("%s: memuat pemisah pernyataan atau komentar:\n%s", nama, q)
		}
	}
}
