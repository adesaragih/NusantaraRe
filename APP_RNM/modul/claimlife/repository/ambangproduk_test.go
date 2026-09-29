package repository

// Ambang produk - butir ba/bh. TANPA Oracle.

import (
	"strings"
	"testing"

	"nusantarare/inti/db"
)

// TestAmbangProdukHanyaDuaKolom - batas izin butir bh.
//
// ⛔ Izinnya SEMPIT, dan kesempitannya yang membuatnya aman: view ini 40
// kolom `[data DBA]`, dan `SELECT *` akan membawa kolom yang izin ini TIDAK
// cakup.
func TestAmbangProdukHanyaDuaKolom(t *testing.T) {
	q := sqlAmbangProduk("SKEMAUJI." + NamaViewProdukLife)
	if strings.Contains(q, "SELECT *") {
		t.Errorf("query membawa seluruh kolom:\n%s", q)
	}
	for _, k := range KolomAmbangProduk {
		if !strings.Contains(q, k) {
			t.Errorf("kolom %q tidak dibaca:\n%s", k, q)
		}
	}
	if !strings.Contains(q, "WHERE v.ID = :1") {
		t.Errorf("query tidak berkunci ID:\n%s", q)
	}
	if err := db.PeriksaSQL(q); err != nil {
		t.Errorf("%v\n%s", err, q)
	}
}

// TestAmbangKosongBukanNol - gagal terang, bukan ambang bawaan.
func TestAmbangKosongBukanNol(t *testing.T) {
	if (AmbangProduk{MaxExpiredClaim: "30", MaxDataReceive: "60"}).Lengkap() != true {
		t.Error("ambang lengkap dianggap tidak lengkap")
	}
	for _, a := range []AmbangProduk{
		{},
		{MaxExpiredClaim: "30"},
		{MaxDataReceive: "60"},
		{MaxExpiredClaim: "  ", MaxDataReceive: "60"},
	} {
		if a.Lengkap() {
			t.Errorf("%+v dianggap lengkap", a)
		}
	}
}

// TestPesanAmbangTidakMenyebutNamaObjekWarisan.
//
// ⛔ Pesan galat mendarat di log, dan nama objek warisan tidak pernah masuk
// log atau artefak. Keterangannya ada di komentar kepala berkas, tempat yang
// memang untuk itu.
func TestPesanAmbangTidakMenyebutNamaObjekWarisan(t *testing.T) {
	pesan := ErrAmbangProdukTakDitemukan.Error()
	if strings.Contains(strings.ToUpper(pesan), NamaViewProdukLife) {
		t.Errorf("pesan %q menyebut nama view warisan", pesan)
	}
	if !strings.Contains(pesan, "TIDAK ada nilai pengganti") {
		t.Errorf("pesan %q tidak menyatakan bahwa nol nilai pengganti dipakai", pesan)
	}
}
