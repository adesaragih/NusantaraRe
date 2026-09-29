package services

// Riwayat tangga Komite - tiket 09. TANPA Oracle.

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"
	"time"

	"nusantarare/inti"
	"nusantarare/modul/komite/repository"
)

// TestRiwayatMembedakanDilewatiDariMenunggu - AC tiket 09.
func TestRiwayatMembedakanDilewatiDariMenunggu(t *testing.T) {
	saat := time.Date(2026, 9, 28, 9, 0, 0, 0, time.UTC)
	k := repository.KasusKomite{AdjID: "UJI-ADJ",
		Baris: repository.BarisInboxKomite{KasusID: "KMTLF-UJI"},
		Tangga: []repository.AnggotaKasus{
			{Urut: 1, OperatorID: "UJI-A", Jabatan: "UJI-J1", Approval: "1", Komentar: "ok",
				TglAprove: sql.NullTime{Time: saat, Valid: true}},
			{Urut: 2, OperatorID: "UJI-B", Jabatan: "UJI-J2", Approval: ""},
			{Urut: 3, OperatorID: "UJI-C", Jabatan: "UJI-J3", Approval: "0"},
		}}
	jejak := []repository.JejakKomite{
		{Dari: "Komite tingkat 2", Ke: "Eskalasi ke tingkat 3 (KMTLF-UJI)", AkunID: "UJI-ADM", Waktu: saat},
		{Dari: "Komite tingkat 1", Ke: "Setuju (KMTLF-UJI)", AkunID: "UJI-A", Waktu: saat},
	}
	r := susunRiwayat(k, jejak)
	mau := []string{"Setuju", KataTingkatDilewati, KataTingkatMenunggu}
	for i, b := range r.Tangga {
		if b.Status != mau[i] {
			t.Errorf("tingkat %d: %q, mau %q", b.Urut, b.Status, mau[i])
		}
	}
	if r.Tangga[0].Committee != "UJI-J1" || r.Tangga[0].TanggalPutus != "2026-09-28 09:00:00" {
		t.Errorf("baris 1 %+v", r.Tangga[0])
	}
	if len(r.Eskalasi) != 1 || r.Eskalasi[0].DariTingkat != 2 || r.Eskalasi[0].KeTingkat != 3 ||
		r.Eskalasi[0].Oleh != "UJI-ADM" {
		t.Errorf("eskalasi %+v - hanya jejak eskalasi yang masuk", r.Eskalasi)
	}
}

// TestRiwayatUntukSiapaPun - melihat bukan memutuskan.
func TestRiwayatUntukSiapaPun(t *testing.T) {
	isi, err := os.ReadFile("komite_riwayat.go")
	if err != nil {
		t.Fatal(err)
	}
	badan := badanFungsi(t, string(isi), "func (i *InboxKomite) Riwayat(")
	if strings.Contains(badan, "periksaGiliran(") || strings.Contains(badan, "PunyaPeran(") ||
		strings.Contains(badan, "susunKasusTampil(") {
		t.Error("riwayat bergerbang wewenang memutuskan/keanggotaan")
	}
	if !strings.Contains(badan, "WajibIdentitas(pelaku)") {
		t.Error("riwayat tanpa identitas pelaku")
	}
	if _, err := New(nil).InboxKomite().Riwayat(context.Background(), inti.Pelaku{}, "K"); err == nil {
		t.Error("riwayat tanpa identitas diterima")
	}
}

// TestPenulisDanPembacaJejakSatuBentuk - teks jejak dari konstanta yang sama.
func TestPenulisDanPembacaJejakSatuBentuk(t *testing.T) {
	isi, err := os.ReadFile("komite_keputusan.go")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(isi), `"Komite tingkat "`) || strings.Contains(string(isi), `"Eskalasi ke tingkat "`) {
		t.Error("penulis jejak mengetik ulang teksnya; pembaca riwayat akan menyimpang")
	}
}

// OQ-K-05 (GILIRAN-17): tingkat yang tertimpa 5.1 menampilkan keputusan dan
// komentar ASLI-nya dari jejak "ditimpa"; eskalasi tetap terbaca.
func TestRiwayatMenampilkanKeputusanAsliTingkatTertimpa(t *testing.T) {
	k := repository.KasusKomite{AdjID: "UJI-ADJ", Baris: repository.BarisInboxKomite{KasusID: "KMTLF-UJI"},
		Tangga: []repository.AnggotaKasus{{Urut: 1, Approval: "2"}, {Urut: 2, Approval: "2"}}}
	jejak := []repository.JejakKomite{
		{Dari: awalanJejakTingkat + "1", Ke: awalanJejakTimpa + "Setuju (KMTLF-UJI)", Komentar: "UJI ok", AkunID: "UJI-C"},
		{Dari: awalanJejakTingkat + "1", Ke: "Setuju (KMTLF-UJI)", AkunID: "UJI-A"},
	}
	r := susunRiwayat(k, jejak)
	if r.Tangga[0].Asli == nil || r.Tangga[0].Asli.Status != "Setuju" || r.Tangga[0].Asli.Comment != "UJI ok" {
		t.Errorf("tingkat 1 tanpa keputusan asli: %+v", r.Tangga[0])
	}
	if r.Tangga[1].Asli != nil {
		t.Errorf("tingkat 2 tidak tertimpa, tetapi berketerangan asli: %+v", r.Tangga[1].Asli)
	}
	if len(r.Eskalasi) != 0 {
		t.Errorf("jejak keputusan terbaca sebagai eskalasi: %+v", r.Eskalasi)
	}
}
