package models

// Uji fungsi murni penanda migrasi (tiket EDM 09): pemasangan menurut NOURUT (AC 40), pasangan bergeser tanpa
// mengubah satu angka (AC 41), rumus berlapis dari EDMNo generasi sebelumnya (AC 42). Fixture fiktif UJI-.

import (
	"reflect"
	"testing"
)

// halamanPenanda - generasi baru (data + selisih) dan generasi lama dengan kunci dagang yang diberikan.
func halamanPenanda(edmLama string) (baru, lama *Halaman) {
	baru, lama = HalamanBaru(), HalamanBaru()
	baru.SetelDaftar(DaftarSpreading, []Baris{
		{"TreatyType": "UJI-SP", "PremiumSpreaded": "72.00000011"},
		{"TreatyType": "UJI-QS", "PremiumSpreaded": "108.000000166"},
		{"TreatyType": "UJI-XL", "PremiumSpreaded": "1"}, // baris tambahan, NOURUT = maks + 1 (ID-17)
	})
	baru.SetelDaftar(DaftarSelisihSpreading, []Baris{
		{"TreatyType": "UJI-SP", "PremiumSpreaded": "-18.00000011"},
		{"TreatyType": "UJI-QS", "PremiumSpreaded": "48.000000166"},
		{"TreatyType": "UJI-XL", "PremiumSpreaded": "1"},
	})
	baru.SetelDaftar(DaftarAngsuran, []Baris{{"InstallmentNo": "01", "Premium": "5"}, {"InstallmentNo": "3", "Premium": "6"}})
	baru.SetelDaftar(DaftarSelisihAngsuran, []Baris{{"InstallmentNo": "01", "Premium": "1"}, {"InstallmentNo": "3", "Premium": "2"}})
	lama.SetelDaftar(DaftarSpreading, []Baris{{"TreatyType": "UJI-QS", "PremiumSpreaded": "60"}, {"TreatyType": "UJI-SP", "PremiumSpreaded": "90"}})
	lama.SetelDaftar(DaftarAngsuran, []Baris{{"InstallmentNo": "1", "Premium": "4"}, {"InstallmentNo": "2", "Premium": "4"}})
	lama.Setel(pmJalurEDMNo, edmLama)
	return baru, lama
}

func TestPenandaPasanganBergeserMenurutNourut(t *testing.T) { // AC 40, 41; ID-34, ID-35
	baru, lama := halamanPenanda("")
	sebelumBaru, sebelumLama := baru.Salin(), lama.Salin()
	p := HitungPenandaMigrasi(baru, lama)
	// Pasangan = posisi (NOURUT): baris 1 SP lawan QS, baris 2 QS lawan SP -> bergeser; baris 3 tanpa pasangan.
	harap := []bool{true, true, false}
	if len(p.Spreading) != 3 {
		t.Fatalf("penanda spreading %+v", p.Spreading)
	}
	for i, b := range p.Spreading {
		if b.PasanganBergeser != harap[i] || b.RumusBerlapis {
			t.Errorf("spreading NOURUT %d: %+v", i+1, b)
		}
	}
	// Angsuran: "01" = 1 (cacah), 3 lawan 2 -> bergeser.
	if len(p.Angsuran) != 2 || p.Angsuran[0].PasanganBergeser || !p.Angsuran[1].PasanganBergeser {
		t.Errorf("penanda angsuran %+v", p.Angsuran)
	}
	// ⛔ AC 41: tidak satu angka pun berubah - kedua halaman utuh.
	if !reflect.DeepEqual(baru, sebelumBaru) || !reflect.DeepEqual(lama, sebelumLama) {
		t.Error("penandaan mengubah halaman")
	}
	if s, a, x, b := p.Cacah(); s != 2 || a != 1 || x != 0 || b != 0 {
		t.Errorf("cacah %d %d %d %d", s, a, x, b)
	}
}

func TestPenandaRumusBerlapisSetiapBaris(t *testing.T) { // AC 42; ID-36
	baru, lama := halamanPenanda("UJI-QP.T1.10.2017.00001/E01")
	p := HitungPenandaMigrasi(baru, lama)
	for i, b := range append(append([]PenandaBaris{}, p.Spreading...), p.Angsuran...) {
		if !b.RumusBerlapis {
			t.Errorf("baris %d: rumus berlapis kosong padahal generasi sebelumnya ber-EDMNo", i)
		}
	}
	if _, _, _, b := p.Cacah(); b != 5 {
		t.Errorf("berlapis %d, harap 5", b)
	}
	// Generasi sebelumnya generasi NB (EDMNo kosong): varian langkah 1-3, bukan berlapis.
	baru, lama = halamanPenanda("  ")
	for _, b := range HitungPenandaMigrasi(baru, lama).Spreading {
		if b.RumusBerlapis {
			t.Error("EDMNo kosong ditandai berlapis")
		}
	}
}

func TestPenandaLapisanXOLMenurutSubskripMataUangDanLapisan(t *testing.T) { // AC 40; CalculateDifferenceEDM_act 1 / 1.2
	baru, lama := HalamanBaru(), HalamanBaru()
	lap := func(layer, part, cur string) Baris {
		return Baris{"Layer": layer, "LayerPart": part, "IDCurrency": cur}
	}
	baru.SetelDaftar(TabelXOL.Daftar, []Baris{{}, {}})
	baru.SetelDaftar(JalurAnak(TabelXOL.Daftar, 1, "ValueList"), []Baris{lap("1", "A", "UJI-1"), lap("2", "A", "UJI-1")})
	baru.SetelDaftar(JalurAnak(TabelXOL.Daftar, 2, "ValueList"), []Baris{lap("1", "A", "UJI-2")})
	baru.SetelDaftar(DaftarSelisihXOL, []Baris{{}, {}})
	baru.SetelDaftar(JalurAnak(DaftarSelisihXOL, 1, "ValueList"), []Baris{lap("1", "A", "UJI-1"), lap("2", "A", "UJI-1")})
	baru.SetelDaftar(JalurAnak(DaftarSelisihXOL, 2, "ValueList"), []Baris{lap("1", "A", "UJI-2")})
	lama.SetelDaftar(TabelXOL.Daftar, []Baris{{}, {}})
	lama.SetelDaftar(JalurAnak(TabelXOL.Daftar, 1, "ValueList"), []Baris{lap("1", "A", "UJI-1")})
	lama.SetelDaftar(JalurAnak(TabelXOL.Daftar, 2, "ValueList"), []Baris{lap("1", "B", "UJI-2")})
	lama.Setel(pmJalurEDMNo, "UJI/E01")
	p := HitungPenandaMigrasi(baru, lama)
	// (1,1) sama; (1,2) tanpa pasangan; (2,1) LayerPart B lawan A -> bergeser. Rumus berlapis tidak berlaku (363).
	harap := []bool{false, false, true}
	if len(p.Lapisan) != 3 {
		t.Fatalf("lapisan %+v", p.Lapisan)
	}
	for i, b := range p.Lapisan {
		if b.PasanganBergeser != harap[i] || b.RumusBerlapis {
			t.Errorf("lapisan NOURUT %d: %+v", i+1, b)
		}
	}
}

func TestTeksPenanda(t *testing.T) {
	if TeksPenanda(true) != "1" || TeksPenanda(false) != "0" {
		t.Error("teks penanda")
	}
}
