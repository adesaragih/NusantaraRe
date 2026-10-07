package models

// Untuk apa berkas ini: AKTIVITAS INDUK jalur NonProp - `InputPolicyTreatyInDetail_NonProp`
// (preACT 16), langkah 18 `InputPolicyTreatyInDetail_preACT` (PPN/PPH per layer),
// pemulihan limit master (`TreatySetReinstatement` / `SetReinstatementPct`), dan
// pra-proses `TreatyRealizationCheckXOLList` (`InputPolicyTreatyInPre_Act` 10).
// Semantik yang ditiru: lihat kepala `nonprop.go`.

import (
	"strconv"

	"github.com/cockroachdb/apd/v3"
)

// PesanGagalXOL - VERBATIM `TreatyRealizationCheckXOLList` langkah 3 (`local.err`).
const PesanGagalXOL = "Error fetching XolList"

// jumlahMataUang menjumlah `.Value` baris bermata uang `mu` (perulangan
// bersyarat `.Currency==local.currency` per baris).
func jumlahMataUang(k *kalkulator, rows []Baris, mu string, awal *apd.Decimal) *apd.Decimal {
	for _, r := range rows {
		if r["Currency"] == mu {
			awal = k.Tambah(awal, angkaBaris(k, r, "Value"))
		}
	}
	return awal
}

// ---------------------------------------------------------------- pemulihan limit master

// SetReinstatementPct = `Activity/SetReinstatementPct` atas baris
// `TreatyIn.Limits(idx)`: `Reinstatement_List` dibangun ulang, satu baris per
// pemulihan (`.ReinstatementValue` kali).
//
//	3.1  ReinstatementPct "100", ReinstatementValue n, ReinstatementNote, Amount1/2 =
//	     Limit/Limit2, ID = pxListSubscript limit, AdditionalPct "100"
//	3.2  MDPList 1 baris: AdditionalAmount1 = MDP IDR, AdditionalAmount2 = MDP USD
//	3.3  MDPList >= 2:    AdditionalAmount1 = MDPList(1) bila IDR, Amount2 = MDPList(2) bila USD
//
// Nilai master ini TIDAK dibaca rule NB lain dan tidak disimpan (halaman master
// baca-saja; pemulihan limit milik modul `treatyin`, tabel PEMULIHAN_LIMIT).
func SetReinstatementPct(h *Halaman, idx int) error {
	limits := mDaftar(h, "Limits")
	if idx < 1 || idx > len(limits) {
		return nil
	}
	l := limits[idx-1]
	jalur := JalurAnak(jMaster+"Limits", idx, "Reinstatement_List")
	h.SetelDaftar(jalur, nil) // 1
	k := &kalkulator{}
	maks := angkaBaris(k, l, "ReinstatementValue") // 2: Local.maxidx (int)
	if k.err != nil {
		return k.err
	}
	// bilangan bulat; pecahan dibuang (`int`)
	c := k.ctx()
	c.Rounding = apd.RoundDown
	var bulat apd.Decimal
	if _, err := c.RoundToIntegralValue(&bulat, maks); err != nil {
		return err
	}
	n, err := bulat.Int64()
	if err != nil {
		return err
	}
	mdp := mAnak(h, "Limits", idx, "MDPList")
	nilaiMU := func(b Baris, mu string) string {
		if b["Currency"] == mu {
			return b["Value"]
		}
		return "0"
	}
	var baris []Baris
	for r := int64(0); r < n; r++ { // 3
		b := Baris{ // 3.1
			"ReinstatementPct": "100", "ReinstatementValue": strconv.FormatInt(r+1, 10),
			"ReinstatementNote": l["ReinstatementNote"], "ReinstatementAmount1": l["Limit"],
			"ReinstatementAmount2": l["Limit2"], "ID": strconv.Itoa(idx), "AdditionalPct": "100",
		}
		if len(mdp) == 1 { // 3.2
			b["AdditionalAmount1"], b["AdditionalAmount2"] = nilaiMU(mdp[0], "IDR"), nilaiMU(mdp[0], "USD")
		}
		if len(mdp) >= 2 { // 3.3
			b["AdditionalAmount1"], b["AdditionalAmount2"] = nilaiMU(mdp[0], "IDR"), nilaiMU(mdp[1], "USD")
		}
		baris = append(baris, b)
	}
	h.SetelDaftar(jalur, baris)
	return nil
}

// TreatySetReinstatement = `Activity/TreatySetReinstatement` (SetTreatyIn_Act 13):
// setiap `TreatyIn.Limits` yang `Reinstatement_List(1).ReinstatementValue` kosong
// -> `SetReinstatementPct`.
func TreatySetReinstatement(h *Halaman) error {
	for i := range mDaftar(h, "Limits") {
		rl := mAnak(h, "Limits", i+1, "Reinstatement_List")
		if len(rl) == 0 || rl[0]["ReinstatementValue"] == "" {
			if err := SetReinstatementPct(h, i+1); err != nil {
				return err
			}
		}
	}
	return nil
}

// ---------------------------------------------------------------- pra-proses

// PerluCekDaftarXOL = syarat `InputPolicyTreatyInPre_Act` langkah 10:
// `IsNewPolicyNonProp==1` dan `.PolicyTreatyIn.EDMType` bukan "3".
//
// `[penyimpangan sadar]` ditambah: jenis proporsi `QuotationData` bukan
// "Proportional". DASAR (diperiksa ulang audit silang P3 7.4, 04-10-2026):
// spec-penyimpanan AC 33 `[terverifikasi]` (ID-35, hasil grilling) - "Menyimpan
// baris XOL pada polis ber-`ProportionalType = 'Proportional'` ditolak"
// (`PeriksaBentukSimpan`). Penanda IsNewPolicyNonProp TIDAK PERNAH diturunkan
// pra-proses: `DataTransform/InputPolicyTreatyIn_preDT` langkah 4 `WHEN
// .Quotation.ProportionalType=="NonProportional"` -> "1", langkah 5 `WHEN
// .PolicyTreatyIn.IsNewPolicyNonProp != "1"` -> "0" (nilai "1" bertahan). Berkas
// NonProp yang dipilih ulang ke kontrak proporsional (preDT 14: QuotationData =
// Quotation "Proportional") akan dibuatkan baris XOL oleh langkah 10 apa adanya -
// dan SETIAP Save/Submit-nya lalu ditolak AC 33 (jalan buntu). Syarat ketiga
// mencegah baris yang penyimpanan tolak; perilaku Pega lain tidak berubah.
// Dicatat di tiket 19 (AC 33).
func PerluCekDaftarXOL(h *Halaman) bool {
	return samaDenganSatu(h.Ambil(pt+"IsNewPolicyNonProp")) &&
		h.Ambil(pt+"EDMType") != "3" &&
		h.Ambil(pt+"QuotationData.ProportionalType") != JenisProporsional
}

// AwalCekDaftarXOL = `TreatyRealizationCheckXOLList` langkah 1-2: pesan
// halaman dibersihkan; true = `TreatyXOLList` kosong, aktivitas berlanjut.
func AwalCekDaftarXOL(h *Halaman) bool {
	h.BersihkanPesan()
	return len(h.AmbilDaftar(DaftarXOL)) < 1
}

// LengkapiDaftarXOL = `TreatyRealizationCheckXOLList` langkah 4-8 sesudah master
// dimuat (`SetTreatyIn_Act` 4-5: RDB `BrowseTreatyIn` + `adoptJSONObject`; 5
// `Page-Copy` ke `pyWorkPage.TreatyIn` - `TerapkanMasterXOL`).
//
//	SetTreatyIn_Act 13  ProportionType "NonProportional" -> TreatySetReinstatement
//	6                   InsertToTreatyXOLList (BUKAN RetroShare)
//	8                   daftar tetap kosong -> pesan "Error fetching XolList"
//
// ⛔ Tidak diport: langkah 7 (`PolicyTreatyIn.OldData.TreatyXOLList`) - OldData
// adalah generasi sebelumnya, digantikan `OLD_POLIS_ID` (diagram F13-F14) yang
// di NB selalu kosong; SetTreatyIn_Act 1-3, 6-12, 14-15 (penanda tampilan dan
// komentar layar master `Data-Portal`, revisi master yang MENULIS JSON
// `SaveTreatyIn` - hanya bila `revisionstate==1`, tidak pernah dari sini;
// `CheckDuplicateOffer` langkah 1-4 `//`).
func LengkapiDaftarXOL(h *Halaman, idMU IDMataUang) error {
	if h.Ambil(jMaster+"ProportionType") == JenisNonProporsional {
		if err := TreatySetReinstatement(h); err != nil {
			return err
		}
	}
	if err := InsertToTreatyXOLList(h, idMU); err != nil {
		return err
	}
	if len(h.AmbilDaftar(DaftarXOL)) < 1 {
		h.TambahPesan("", PesanGagalXOL)
	}
	return nil
}
