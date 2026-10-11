package repository

import (
	"strings"
	"testing"
)

// TestDuaBentukIDPega - keputusan work owner 11-10-2026: tulisan = case id polos; pembaca mencocokkan (polos, berawalan) dari case
// id ATAU IDPEGA tersimpan bentuk mana pun (baris lama Pega / sistem baru berawalan tetap terbaca, CekFacin / cek JSON_OFFER tidak
// menyisip ganda).
func TestDuaBentukIDPega(t *testing.T) {
	for _, masuk := range []string{"NB-151033", "ASM-FW-GISFW-WORK NB-151033"} {
		p, a := DuaBentukIDPega(masuk)
		if p != "NB-151033" || a != "ASM-FW-GISFW-WORK NB-151033" || IDPegaProduksi(masuk) != "NB-151033" {
			t.Errorf("%q: %q %q %q", masuk, p, a, IDPegaProduksi(masuk))
		}
	}
	// Ekspresi case id SQL: berawalan -> SUBSTR(19), polos -> apa adanya (View Policy / Following).
	e := sqlCaseIDDariIDPega("p.IDPEGA")
	if e != "CASE WHEN SUBSTR(p.IDPEGA, 1, 18) = 'ASM-FW-GISFW-WORK ' THEN SUBSTR(p.IDPEGA, 19) ELSE p.IDPEGA END" {
		t.Error(e)
	}
	// Setiap pembaca IDPEGA nbfacin ber-IN dua bentuk (sensus kebalikan: tidak ada lagi `IDPEGA = :n` / `ID_PEGA = :n` tunggal di
	// pembaca tabel-tabel itu).
	for nama, q := range map[string]string{"CekFacin": sqlAdaFacInProduction("P"), "JSON_POLIS": sqlCariJSONPolis("J"),
		"hapus ceding": sqlHapusCedingProduksi("C"), "ceding PRODKE": sqlSisipCedingProduksi("C", "J"), "JSON_OFFER": sqlCariJSONOffer("JO"),
		"ubah JSON_OFFER": sqlUbahJSONOffer("JO"), "TGL_BIND": sqlTglBind("JO"), "hapus FACINOFFER": sqlHapusFacInOffer("FO"),
		"riwayat": sqlRiwayatAkseptasi("P"), "NOURUT": sqlSisipRiwayatProduksi("P"), "Banding": sqlOpsiBanding("L", "HP"),
		"Ask": sqlCatatanPegaNonMarketing("HP"), "lampiran jumlah": sqlJumlahLampiran("D"), "lampiran daftar": sqlDaftarLampiran("D", "O", "L", false),
		"lampiran hapus": sqlHapusLampiran("D"), "kasus ganda": sqlKasusGanda("P"), "akumulasi": sqlHitungAkumulasi("P", "FO"),
		"following ada": sqlAdaPolisFollowing("P"), "polis baris": sqlBarisPolis("P"), "polis ceding": sqlCedingPolis("C")} {
		u := strings.ReplaceAll(q, "\n", " ")
		if !strings.Contains(u, "IDPEGA IN (") && !strings.Contains(u, "ID_PEGA IN (") && !strings.Contains(u, "IDPEGA NOT IN (") {
			t.Errorf("%s tanpa dua bentuk: %s", nama, q)
		}
	}
}
