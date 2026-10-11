package repository

// Penanda nama kembar — pagar pencarian BALIK pengenal dari nama.
//
// ⛔ Tanpa penanda ini, layar yang mencari pengenal dari nama akan memilih
// salah satu dari dua baris bernama sama, dan layer menunjuk jenis treaty
// yang KELIRU tanpa ada yang tahu. Itu kesalahan yang baru ketahuan
// bertahun kemudian.

import (
	"testing"

	"nusantarare/modul/treatyin/backend/models"
)

func TestTandaiKembarMenandaiNamaYangDipakaiLebihDariSatu(t *testing.T) {
	d := []models.PilihanWarisan{
		{ID: "1", Nama: "QUOTA SHARE"},
		{ID: "2", Nama: "EXCESS OF LOSS"},
		{ID: "3", Nama: "EXCESS OF LOSS"},
		{ID: "4", Nama: "SURPLUS"},
	}
	tandaiKembar(d)
	mau := map[string]bool{"1": false, "2": true, "3": true, "4": false}
	for _, p := range d {
		if p.Kembar != mau[p.ID] {
			t.Errorf("%s (%q) kembar=%v, mau %v", p.ID, p.Nama, p.Kembar, mau[p.ID])
		}
	}
}

// ⛔ Daftar tanpa kembaran tidak boleh menandai apa pun — penanda yang selalu
// menyala membuat pencarian balik tidak pernah bekerja.
func TestTandaiKembarDiamBilaSeluruhNamaTunggal(t *testing.T) {
	d := []models.PilihanWarisan{
		{ID: "1", Nama: "QUOTA SHARE"},
		{ID: "2", Nama: "SURPLUS"},
	}
	tandaiKembar(d)
	for _, p := range d {
		if p.Kembar {
			t.Errorf("%q ditandai kembar padahal tunggal", p.Nama)
		}
	}
}

func TestTandaiKembarDaftarKosongTidakPanik(t *testing.T) {
	tandaiKembar(nil)
	tandaiKembar([]models.PilihanWarisan{})
}
