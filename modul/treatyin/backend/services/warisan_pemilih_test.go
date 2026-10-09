package services_test

// Isi pemilih "Choose Ceding" / "Choose Source of Business".
//
// ⛔ Yang diuji di sini BUKAN pembacaannya melainkan PENANDAAN nama kembar —
// satu-satunya tafsir yang lapisan ini tambahkan, dan satu-satunya tempat
// kesalahan dapat menyelinap tanpa terlihat di layar.

import (
	"context"
	"errors"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

// NEGATIF — tanpa identitas, gudang tidak pernah disentuh.
func TestPemilihMenolakTanpaIdentitas(t *testing.T) {
	g := &gudangTiruan{cedant: []models.PilihanWarisan{{ID: "G0000002", Nama: "ARTHAGRAHA"}}}
	l := services.LayananDengan(g)

	if _, err := l.DaftarCedant(context.Background(), inti.Pelaku{}); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("cedant: mau ErrTanpaIdentitas, dapat %v", err)
	}
	if _, err := l.DaftarAsalBisnis(context.Background(), inti.Pelaku{}); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("asal bisnis: mau ErrTanpaIdentitas, dapat %v", err)
	}
}

// NEGATIF — galat gudang diteruskan apa adanya, tidak ditelan menjadi
// daftar kosong. Pemilih kosong yang sebenarnya galat membuat orang mengira
// memang tidak ada cedant.
func TestPemilihMeneruskanGalatGudang(t *testing.T) {
	rusak := errors.New("oracle mati")
	g := &gudangTiruan{galatPilihan: rusak}
	l := services.LayananDengan(g)

	if _, err := l.DaftarCedant(context.Background(), pelakuAda); !errors.Is(err, rusak) {
		t.Errorf("mau galat diteruskan, dapat %v", err)
	}
}

// POSITIF — nama yang dipakai LEBIH DARI SATU pengenal ditandai, yang
// tunggal tidak. Data contohnya nyata: `ASURANSI ADIRA DINAMIKA` memang
// tercatat dengan dua pengenal di POOLDATA pada 4 Oktober 2026, salah
// satunya ID kerja Pega yang menyelinap menjadi data.
func TestPemilihMenandaiNamaKembar(t *testing.T) {
	g := &gudangTiruan{cedant: []models.PilihanWarisan{
		{ID: "G0000002", Nama: "ARTHAGRAHA GENERAL INSURANCE"},
		{ID: "ASM-SFAGIS-WORK-ORG ORG-34", Nama: "ASURANSI ADIRA DINAMIKA"},
		{ID: "G0000006", Nama: "ASURANSI ADIRA DINAMIKA"},
	}}
	l := services.LayananDengan(g)

	hasil, err := l.DaftarCedant(context.Background(), pelakuAda)
	if err != nil {
		t.Fatalf("mau diterima, dapat %v", err)
	}
	if len(hasil) != 3 {
		t.Fatalf("mau 3 baris — kembar TIDAK disatukan, dapat %d", len(hasil))
	}
	mau := map[string]bool{
		"G0000002":                   false,
		"ASM-SFAGIS-WORK-ORG ORG-34": true,
		"G0000006":                   true,
	}
	for _, b := range hasil {
		if b.Kembar != mau[b.ID] {
			t.Errorf("%s (%s): mau kembar=%v, dapat %v", b.ID, b.Nama, mau[b.ID], b.Kembar)
		}
	}
}

// POSITIF — daftar KOSONG adalah jawaban yang sah, bukan galat.
func TestPemilihKosongBukanGalat(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{})

	hasil, err := l.DaftarAsalBisnis(context.Background(), pelakuAda)
	if err != nil {
		t.Fatalf("mau diterima, dapat %v", err)
	}
	if len(hasil) != 0 {
		t.Errorf("mau nol baris, dapat %d", len(hasil))
	}
}
