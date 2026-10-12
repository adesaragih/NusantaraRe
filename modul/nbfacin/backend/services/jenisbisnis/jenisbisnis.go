// Package jenisbisnis - BusinessType_DeT (`ASM-FW-GISFW-WORK!BUSINESSTYPE_DET`, `DDL\\BusinessType_DeT.xml`): first match atas
// (GroupPanel, BusinessOldId) baris BUSINESS; tanpa baris cocok -> "UNKNOWN" (pyDefaultResult). Tabel = urai XML, 36 baris (sensus dua
// arah sesi c3, tiket 73 P0). Dipindah dari services/lini.go supaya dipakai bersama services (lini kasus) dan services/kontrakfacin
// (kontrak.JenisBisnisFacIn untuk Endorsement Fac In, keputusan work owner 12-10-2026 opsi "1a") tanpa siklus impor.
package jenisbisnis

// Bawaan - pyDefaultResult BusinessType_DeT.
const Bawaan = "UNKNOWN"

// baris - satu baris BusinessType_DeT; oldID kosong = cocok apa pun.
type baris struct {
	groupPanel string
	oldID      []string
	hasil      string
}

// tabel - BusinessType_DeT baris 1-36 urut (first match).
var tabel = []baris{
	{"001", nil, "Medicare"},
	{"002", []string{"03"}, "PA"},
	{"003", []string{"36", "37", "38", "39", "40", "41", "42", "43", "44", "45", "46", "47", "48", "49", "50", "51", "52", "C2", "C3",
		"C4", "C5", "C6", "C7", "C8", "C9", "D1", "D2", "D3", "D4", "D5", "D6", "D7", "D8", "D9", "F1", "F2", "F3", "F4", "F5"}, "Bonding"},
	{"003", []string{"95", "96", "97", "98", "99", "A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "B4", "B5", "B6", "B8", "B9",
		"C1"}, "BondingKBG"},
	{"003", []string{"20"}, "HE"},
	{"003", []string{"05", "55"}, "MarineHull"},
	{"003", []string{"21"}, "Glass"},
	{"003", []string{"15", "24", "89", "90", "91", "E7", "E9"}, "Liability"},
	{"003", []string{"16"}, "AllRisk"},
	{"003", []string{"06"}, "AviationHull"},
	{"003", []string{"14"}, "Burglary"},
	{"003", []string{"07"}, "Car"},
	{"003", []string{"08"}, "Ear"},
	{"003", []string{"09"}, "ElectronicEquipment"},
	{"003", []string{"19"}, "Fidelity"},
	{"003", []string{"27"}, "GolfInsurance"},
	{"003", []string{"17"}, "CIT"},
	{"003", []string{"18"}, "CIS"},
	{"003", []string{"11"}, "Boiler"},
	{"003", []string{"93"}, "Workmen"},
	{"003", []string{"10"}, "MBD"},
	{"003", []string{"12"}, "ContractorsPM"},
	{"003", []string{"79", "80", "81", "82", "83", "84", "85", "86", "87", "88"}, "CustomBond"},
	{"003", []string{"28"}, "FireStyle1"},
	{"003", []string{"B3"}, "LandRig"},
	{"003", nil, "Aneka"},
	{"004", nil, "MarineCargo"},
	{"005", nil, "Travel"},
	{"006", []string{"78"}, "OilGas"},
	{"006", []string{"22", "26", "29", "57"}, "FireStyle1"},
	{"006", []string{"01", "34", "59"}, "FireStyle2"},
	{"006", nil, "Fire"},
	{"007", []string{"02"}, "MBUCar"},
	{"007", []string{"60"}, "MBUMotorCycle"},
	{"009", []string{"L1", "L2", "L3", "L4", "L5", "L6", "L7", "L8", "L9", "L10", "L11", "L12", "L13", "L14", "L15", "L16", "L17", "L18",
		"L19", "L20", "L21"}, "Life"},
	{"008", nil, "Medicare"},
}

// JumlahBaris - banyak baris tabel (uji sensus).
func JumlahBaris() int { return len(tabel) }

// Dari - BusinessType_DeT first match atas (GroupPanel, BusinessOldId); tanpa baris cocok -> Bawaan.
func Dari(groupPanel, oldID string) string {
	for _, b := range tabel {
		if b.groupPanel != groupPanel {
			continue
		}
		if len(b.oldID) == 0 {
			return b.hasil
		}
		for _, o := range b.oldID {
			if o == oldID {
				return b.hasil
			}
		}
	}
	return Bawaan
}
