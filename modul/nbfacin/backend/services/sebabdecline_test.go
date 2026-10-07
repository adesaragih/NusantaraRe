package services

// Berkas BANGKITAN (tiket 57): jawaban dihasilkan frontend/components/sebabDecline.ts sendiri (Node 24, type stripping,
// validasiSebab + komentarDecline) atas 20 masukan; port Go wajib sama persis. Jangan disunting tangan. Membangkitkan
// ulang: salin frontend/components/sebabDecline.ts (impor '../labels' -> './labels.ts') dan frontend/labels.ts ke satu
// folder, lalu `node` sebuah modul yang memanggil validasiSebab / komentarDecline untuk setiap masukan di bawah dan
// mencetak JSON {pesan, komentar}.

import (
	"reflect"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestSebabDeclineParitas - validasiSebab / komentarDecline = sebabDecline.ts.
func TestSebabDeclineParitas(t *testing.T) {
	kasus := []struct {
		daftar   []models.PilihanSebabTc
		nonTc    string
		pesan    []string
		komentar string
	}{
		{[]models.PilihanSebabTc{{TermCondition: "Rate", Prioritas: "I", Description1: "", Description2: "", PctRNM: "1.5", PctSOB: "2"}}, "", nil, "Rate   1.5 % RNM 2 % SOB \n"},
		{[]models.PilihanSebabTc{{TermCondition: "Rate", Prioritas: "I", Description1: "", Description2: "", PctRNM: "1.5", PctSOB: ""}}, "", []string{" PctRNM or PctSOB Can't be NULL"}, "Rate   1.5 % RNM \n"},
		{[]models.PilihanSebabTc{{TermCondition: "R/I Comm, Deduction", Prioritas: "I", Description1: "", Description2: "", PctRNM: "", PctSOB: "3"}}, "", []string{" PctRNM or PctSOB Can't be NULL"}, "R/I Comm, Deduction   3 % SOB \n"},
		{[]models.PilihanSebabTc{{TermCondition: "OJK Discount", Prioritas: "II", Description1: "", Description2: "", PctRNM: "", PctSOB: ""}}, "", []string{" PctRNM or PctSOB Can't be NULL", "Pilih minimal satu Cause of Decline."}, "\nOJK Discount   "},
		{[]models.PilihanSebabTc{{TermCondition: "Loss Ratio (net)", Prioritas: "I", Description1: "", Description2: "", PctRNM: "", PctSOB: ""}}, "", []string{" PctRNM Can't be NULL"}, "Loss Ratio (net)   \n"},
		{[]models.PilihanSebabTc{{TermCondition: "Loss Ratio (net)", Prioritas: "I", Description1: "", Description2: "", PctRNM: "65", PctSOB: ""}}, "", nil, "Loss Ratio (net)   65 % RNM \n"},
		{[]models.PilihanSebabTc{{TermCondition: "Occupancy", Prioritas: "I", Description1: "", Description2: "", PctRNM: "", PctSOB: ""}}, "", []string{" Remaks Can't be NULL"}, "Occupancy   \n"},
		{[]models.PilihanSebabTc{{TermCondition: "Occupancy", Prioritas: "I", Description1: "", Description2: "pabrik kimia", PctRNM: "", PctSOB: ""}}, "", nil, "Occupancy  pabrik kimia \n"},
		{[]models.PilihanSebabTc{{TermCondition: "TSI (in IDR)", Prioritas: "II", Description1: "", Description2: "", PctRNM: "", PctSOB: ""}}, "Capacity", []string{" Remaks Can't be NULL"}, "Capacity\nTSI (in IDR)   "},
		{[]models.PilihanSebabTc{{TermCondition: "Clauses", Prioritas: "I", Description1: "", Description2: "", PctRNM: "", PctSOB: ""}}, "", nil, "Clauses   \n"},
		{[]models.PilihanSebabTc{{TermCondition: "Period Of Insurance (POI)", Prioritas: "I", Description1: "", Description2: "", PctRNM: "", PctSOB: ""}}, "", []string{"Option Can't be NULL"}, "Period Of Insurance (POI)   \n"},
		{[]models.PilihanSebabTc{{TermCondition: "Extension of coverage", Prioritas: "I", Description1: "Others", Description2: "", PctRNM: "", PctSOB: ""}}, "", []string{"Remarks Can't be NULL"}, "Extension of coverage Others  \n"},
		{[]models.PilihanSebabTc{{TermCondition: "Extension of coverage", Prioritas: "I", Description1: "Others", Description2: "banjir", PctRNM: "", PctSOB: ""}}, "", nil, "Extension of coverage Others banjir \n"},
		{[]models.PilihanSebabTc{{TermCondition: "Risk Exposure", Prioritas: "II", Description1: "Others", Description2: "", PctRNM: "", PctSOB: ""}}, "Capacity", []string{"Remarks Can't be NULL"}, "Capacity\nRisk Exposure Others  "},
		{[]models.PilihanSebabTc{{TermCondition: "Risk Exposure", Prioritas: "II", Description1: "Gempa", Description2: "", PctRNM: "", PctSOB: ""}}, "", []string{"Pilih minimal satu Cause of Decline."}, "\nRisk Exposure Gempa  "},
		{[]models.PilihanSebabTc{}, "", []string{"Pilih minimal satu Cause of Decline."}, "\n"},
		{[]models.PilihanSebabTc{}, "Capacity", nil, "Capacity\n"},
		{[]models.PilihanSebabTc{{TermCondition: "Rate", Prioritas: "I", Description1: "", Description2: "", PctRNM: "1.5", PctSOB: "2"}, {TermCondition: "Occupancy", Prioritas: "II", Description1: "", Description2: "gudang", PctRNM: "", PctSOB: ""}}, "", nil, "Rate   1.5 % RNM 2 % SOB \nOccupancy  gudang "},
		{[]models.PilihanSebabTc{{TermCondition: "Rate", Prioritas: "I", Description1: "", Description2: "", PctRNM: "1.5", PctSOB: "2"}, {TermCondition: "Occupancy", Prioritas: "II", Description1: "", Description2: "gudang", PctRNM: "", PctSOB: ""}}, "Capacity", nil, "Capacity\nOccupancy  gudang "},
		{[]models.PilihanSebabTc{{TermCondition: "Others T/C", Prioritas: "I", Description1: "", Description2: "alasan lain", PctRNM: "", PctSOB: ""}, {TermCondition: "Lack of data for acceptance", Prioritas: "II", Description1: "", Description2: "", PctRNM: "", PctSOB: ""}}, "", []string{" Remaks Can't be NULL"}, "Others T/C  alasan lain \nLack of data for acceptance   "},
	}
	for i, k := range kasus {
		if p := validasiSebab(k.daftar, k.nonTc); !reflect.DeepEqual(p, k.pesan) {
			t.Errorf("%d pesan %q, mau %q", i, p, k.pesan)
		}
		if c := komentarDecline(k.daftar, k.nonTc); c != k.komentar {
			t.Errorf("%d komentar %q, mau %q", i, c, k.komentar)
		}
	}
}
