package repository

// Tab Limits, Share, Event Limits, dan RNM Share dibaca dari TABEL
// PENDARATAN — bukan lagi dari `M_TREATY_IN.JSONDATA`.
//
// ---------------------------------------------------------------------
// ⛔ KEPUTUSAN PEMILIK PROSES, 6 Oktober 2026
// ---------------------------------------------------------------------
//
//	"nilai yang ditarik dari JSON data itu dilarang keras, gunakan table
//	 baru yang pernah saya berikan xlsx nya"
//
// Lalu dipersempit ke daftar tabel yang BOLEH dipakai — keempat belas yang
// migrasi `437` dan `430` sudah dirikan:
//
//	T_TREATY_LIMITS        T_TREATY_LIMIT_DETAIL     T_TREATY_LIMIT_COB
//	T_TREATY_SHARE         T_TREATY_LIMIT_GROUP      T_TREATY_LIMIT_GROUP_COB
//	T_TREATY_RETRO_SHARE   T_TREATY_SHARE_SPREADING  T_TREATY_LIMIT_ACHIEVEMENT
//	T_TREATY_FAC_SHARE     T_TREATY_SHARE_DEDUCTION  T_TREATY_INSTALLMENT
//	T_TREATY_FAC_REINSURER T_TREATY_FAC_SHARE_DEDUCTION
//
// Berkas ini menggantikan `warisan_layer_dokumen.go` dan
// `warisan_limits_pohon.go` sebagai sumber keempat tab itu. Keduanya TIDAK
// dihapus di langkah ini: keduanya memegang pengetahuan bentuk dokumen yang
// pemuat pendaratan masih perlukan, dan membuangnya bersamaan dengan
// menukar sumber akan menyatukan dua perubahan yang harus dapat dibalik
// sendiri-sendiri.
//
// ---------------------------------------------------------------------
// ⚠️ YANG TIDAK PUNYA RUMAH DI ANTARA KEEMPAT BELAS, DAN KARENA ITU KOSONG
// ---------------------------------------------------------------------
// Delapan nilai grid lahir dari larik DI DALAM elemen `Limits[]` atau
// `Detail[]` yang keempat belas tabel itu tidak memuat:
//
//	MDP, MDP kedua           <- Limits[].MDPList[]            (T_TREATY_LIMIT_MEASURE)
//	EARN_PREMIUM, keduanya   <- Limits[].PremiumEarnedList[]  (T_TREATY_LIMIT_MEASURE)
//
// ⭐ Kedua baris di atas pernah berakhiran "belum ada", dan itu benar
// pada hari ditulis: `T_TREATY_LIMIT_MEASURE` lahir migrasi 439. Ia ADA
// sejak 6 Oktober 2026, berisi 15.232 baris, dan keempat medannya kini
// dibaca di bawah.
//
// ⛔ YANG MEMBUKTIKAN LUBANGNYA NYATA: `TestCacahLayerMataUangDuaTerukur`
// merah dengan `MDPKedua terisi pada 0 baris, terukur >= 106`. Jalur
// dokumen memberi 106 angka MDP dan 100 angka premi bermata uang kedua;
// jalur pendaratan memberi NOL. Pengurai dokumen yang Tugas D hendak
// buang adalah yang memperlihatkannya.
//	CESSION_TO_RI            <- Detail[].CessionList[]        (T_TREATY_LIMIT_AMOUNT, belum ada)
//	LIABILITY_RNM            <- Detail[].RNMShareList[]       (belum ada)
//	LIABILITYQSOR/QSRI       <- Detail[].RNMSpreadedList*[]   (belum ada)
//	QSOR, QSRI               <- Detail[].SpreadingList[]      (belum ada)
//
// ⛔ Kolomnya DIKOSONGKAN, bukan diisi dari dokumen diam-diam. Mengambil
// sebagian dari tabel dan sebagian dari JSON akan membuat larangan di atas
// terlihat dipatuhi padahal tidak — dan yang membacanya tidak akan tahu
// mana yang mana. Daftar di atas adalah tagihannya, dan ia terbaca di satu
// tempat.
//
// ⭐ `EPI100` adalah PENGECUALIAN yang dinyatakan: ia dahulu dibaca dari
// `Detail[].EPIList[0].Value`, dan `T_TREATY_LIMIT_DETAIL` memuat kolom
// skalar `EPI` berdampingan dengan `CURRENCYEPI`. Yang dipakai kolom itu.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/models"
)

// barisPendaratan adalah satu baris tabel pendaratan: pengenalnya, pengenal
// induknya, dan muatannya berkunci NAMA KOLOM.
type barisPendaratan struct {
	ID    int64
	Induk int64
	Nilai map[string]string
}

// bacaPendaratan membaca satu tabel pendaratan milik satu kontrak.
//
// ⛔ `ORDER BY URUTAN`, dan itu wajib: urutan larik di dokumen Pega adalah
// urutan yang layar lama tampilkan, dan tanpa klausa ini Oracle bebas
// mengembalikan baris dalam urutan apa pun.
//
// ⚠️ `MASTERID` ada di tabel ANAK juga, jadi tiap tabel dibaca SENDIRI lalu
// dirangkai di Go lewat `IDINDUK`. Itu disengaja: satu kueri ber-join empat
// tingkat mengembalikan hasil kali kartesian antar grid bersaudara, dan
// yang membacanya harus membongkarnya kembali.
func (g *Gudang) bacaPendaratan(ctx context.Context, tabel string, punyaInduk bool,
	kolom []string, masterID string) ([]barisPendaratan, error) {
	nama, err := g.db.Qualify(tabel)
	if err != nil {
		return nil, err
	}
	pilih := "ID"
	if punyaInduk {
		pilih += ", IDINDUK"
	} else {
		pilih += ", 0"
	}
	if len(kolom) > 0 {
		pilih += ", " + strings.Join(kolom, ", ")
	}
	q := fmt.Sprintf("SELECT %s FROM %s WHERE MASTERID = :1 ORDER BY URUTAN", pilih, nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, masterID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s kontrak %s: %w", tabel, masterID, err)
	}
	defer func() { _ = rows.Close() }()

	var out []barisPendaratan
	for rows.Next() {
		var id, induk int64
		sel := make([]sql.NullString, len(kolom))
		tuju := make([]any, 0, len(kolom)+2)
		tuju = append(tuju, &id, &induk)
		for i := range sel {
			tuju = append(tuju, &sel[i])
		}
		if err := rows.Scan(tuju...); err != nil {
			return nil, fmt.Errorf("repository: membaca baris %s: %w", tabel, err)
		}
		b := barisPendaratan{ID: id, Induk: induk, Nilai: make(map[string]string, len(kolom))}
		for i, k := range kolom {
			b.Nilai[k] = sel[i].String
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// Kolom yang keempat tab baca. Disebut di sini, bukan `SELECT *`: kolom yang
// bertambah di DDL tidak boleh diam-diam masuk ke jawaban API.
var (
	kolomLimits = []string{
		"LAYER", "LAYERTYPE", "COVER", "TREATYTYPE", "TREATYTYPEID", "CURRENCY", "CURRENCY2",
		"LIMIT", "LIMIT2", "DEDUCTIBLE", "DEDUCTIBLE2", "ADJRATE", "ROLPCT", "MDPPCT",
		"CURRENCYRELATION",
	}
	kolomLimitDetail = []string{
		"TREATYGROUP", "TREATYGROUPID", "TREATYTYPE", "CESSIONPCT", "BROKERAGE", "RNMSHARE",
		"RIOGR", "SPREADINGTYPE", "EPI",
		"RSMDLIMIT", "CURRENCYRSMD", "EARTHQUAKE", "CURRENCYEARTHQUAKE",
		"FLOODJAB", "CURRENCYFLOODJAB", "FLOODNATION", "CURRENCYFLOODNAT",
	}
	kolomLimitCOB = []string{"CLASSOFBUSINESS"}
	// ⭐ `T_TREATY_LIMIT_MEASURE` — migrasi 439, 15.232 baris terukur.
	// `JENIS` membedakan ketiga larik yang mendarat ke tabel ini.
	kolomLimitMeasure = []string{"JENIS", "VALUE"}
)

// BacaLayerPendaratan menyusun baris LAYER dari tabel pendaratan.
//
// Bentuk keluarannya SAMA PERSIS dengan `LayerDariDokumen` dahulu — satu
// baris per `Detail[]`, atau satu baris per `Limits[]` bila ia tidak punya
// detail — supaya keempat tab tidak perlu berubah bersama penukaran sumber
// ini. Yang berubah dari mana nilainya datang, bukan bentuknya.
func (g *Gudang) BacaLayerPendaratan(ctx context.Context, masterID string) ([]models.BarisLayerWarisan, error) {
	limits, err := g.bacaPendaratan(ctx, "T_TREATY_LIMITS", false, kolomLimits, masterID)
	if err != nil {
		return nil, err
	}
	detail, err := g.bacaPendaratan(ctx, "T_TREATY_LIMIT_DETAIL", true, kolomLimitDetail, masterID)
	if err != nil {
		return nil, err
	}
	cob, err := g.bacaPendaratan(ctx, "T_TREATY_LIMIT_COB", true, kolomLimitCOB, masterID)
	if err != nil {
		return nil, err
	}
	ukur, err := g.bacaPendaratan(ctx, "T_TREATY_LIMIT_MEASURE", true, kolomLimitMeasure, masterID)
	if err != nil {
		return nil, err
	}
	return RangkaiLayer(limits, detail, cob, ukur), nil
}

// RangkaiLayer merangkai ketiga tabel menjadi baris layar.
//
// ⛔ DIPISAH dari pembacaannya supaya ia dapat diuji TANPA Oracle — pola
// yang sama dengan `LayerDariDokumen` yang ia gantikan. Perangkaian adalah
// tempat kesalahan bentuk hidup; kueri hanya mengambil baris.
func RangkaiLayer(limits, detail, cob, ukur []barisPendaratan) []models.BarisLayerWarisan {
	// ⭐ Besaran per LAYER, dipisah menurut `JENIS`-nya lalu URUT `URUTAN`.
	//
	// ⛔ Urutannya BAGIAN DARI ARTINYA, bukan kerapian: elemen ke-0 milik
	// mata uang pertama dan ke-1 milik yang kedua (`Currency`/`Currency2`
	// baris layer yang sama). Membacanya tanpa urutan membuat angka USD
	// tampil di kolom IDR — salah yang terlihat benar.
	//
	// `bacaPendaratan` sudah `ORDER BY URUTAN`, dan `URUTAN` dihitung PER
	// INDUK, jadi pengelompokan di bawah mempertahankannya.
	besaran := map[int64]map[string][]string{}
	for _, u := range ukur {
		j := u.Nilai["JENIS"]
		if j == "" {
			continue
		}
		if besaran[u.Induk] == nil {
			besaran[u.Induk] = map[string][]string{}
		}
		besaran[u.Induk][j] = append(besaran[u.Induk][j], u.Nilai["VALUE"])
	}
	// keN mengambil elemen ke-i satu larik besaran, kosong bila tidak ada.
	keN := func(induk int64, jenis string, i int) string {
		l := besaran[induk][jenis]
		if i >= len(l) {
			return ""
		}
		return l[i]
	}
	// Kelas bisnis per detail, urut `URUTAN` sebagaimana terbaca.
	kelas := map[int64][]string{}
	for _, c := range cob {
		if v := c.Nilai["CLASSOFBUSINESS"]; v != "" {
			kelas[c.Induk] = append(kelas[c.Induk], v)
		}
	}
	anak := map[int64][]barisPendaratan{}
	for _, d := range detail {
		anak[d.Induk] = append(anak[d.Induk], d)
	}

	// ⛔ Irisan KOSONG, bukan nil — layar membedakan "nol baris" dari
	// "belum dibaca", dan `null` di JSON jawaban merender berbeda.
	hasil := []models.BarisLayerWarisan{}
	for _, l := range limits {
		n := l.Nilai
		dasar := models.BarisLayerWarisan{
			Layer:              n["LAYER"],
			JenisLayer:         n["LAYERTYPE"],
			DasarCover:         n["COVER"],
			JenisTreaty:        n["TREATYTYPE"],
			MataUang:           n["CURRENCY"],
			MataUangLimit:      n["CURRENCY2"],
			Limit100:           n["LIMIT"],
			RetensiCedant:      n["DEDUCTIBLE"],
			AdjRate:            n["ADJRATE"],
			ROL:                n["ROLPCT"],
			RasioMDP:           n["MDPPCT"],
			RelasiMataUang:     n["CURRENCYRELATION"],
			Limit100Kedua:      n["LIMIT2"],
			RetensiCedantKedua: n["DEDUCTIBLE2"],
			MDP:                keN(l.ID, "MDPList", 0),
			MDPKedua:           keN(l.ID, "MDPList", 1),
			PremiEarned:        keN(l.ID, "PremiumEarnedList", 0),
			PremiEarnedKedua:   keN(l.ID, "PremiumEarnedList", 1),
			KelasBisnis:        []string{},
		}
		rinci := anak[l.ID]
		if len(rinci) == 0 {
			hasil = append(hasil, dasar)
			continue
		}
		for _, d := range rinci {
			m := d.Nilai
			b := dasar
			b.KelompokTreaty = m["TREATYGROUP"]
			// `TreatyType` ada di DUA tingkat; yang di detail lebih khusus,
			// jadi ia menang bila terisi — aturan yang sama seperti dahulu.
			if v := m["TREATYTYPE"]; v != "" {
				b.JenisTreaty = v
			}
			b.PersenCession = m["CESSIONPCT"]
			b.PersenBrokerage = m["BROKERAGE"]
			b.RNMShare = m["RNMSHARE"]
			b.RIOGR = m["RIOGR"]
			b.JenisPenyebaran = m["SPREADINGTYPE"]
			b.EPI100 = m["EPI"]
			b.Gempa = m["EARTHQUAKE"]
			b.BatasRSMD = m["RSMDLIMIT"]
			b.BatasBanjirJab = m["FLOODJAB"]
			b.BatasBanjirNas = m["FLOODNATION"]
			b.MataUangRSMD = m["CURRENCYRSMD"]
			b.MataUangGempa = m["CURRENCYEARTHQUAKE"]
			b.MataUangBanjirJab = m["CURRENCYFLOODJAB"]
			b.MataUangBanjirNas = m["CURRENCYFLOODNAT"]
			// ⛔ Irisan KOSONG per baris, bukan nil — grid `Class of
			// Business` yang kosong berbunyi `No items`.
			b.KelasBisnis = []string{}
			b.KelasBisnis = append(b.KelasBisnis, kelas[d.ID]...)
			hasil = append(hasil, b)
		}
	}
	return hasil
}

// ⭐ POHON tab Limits — Prop DAN Non-Prop — dari tabel pendaratan, 6 Oktober
// 2026.
//
// ⛔ Kunci dan kolomnya TIDAK ditulis ulang di sini. Sumbernya satu:
// `PetaPendaratan` — peta yang sama yang MEMUAT tabel-tabel ini. Daftar
// kedua di berkas ini akan menyimpang dari pemuatnya tanpa ada yang tahu,
// dan medan yang dimuat tetapi tidak dibaca akan tampil kosong di layar.
//
// ⛔ Kuncinya ejaan PEGA (`TreatyGroup`, `QSPct`, `Limit2`), bukan nama kolom
// Oracle: layar mengikat medannya dengan ejaan itu.
//
//	Limits (T_TREATY_LIMITS)
//	├ Detail (T_TREATY_LIMIT_DETAIL)                       — Prop
//	│ ├ COBList (T_TREATY_LIMIT_COB)
//	│ ├ AchievementLists (T_TREATY_LIMIT_ACHIEVEMENT)
//	│ ├ SpreadingList (T_TREATY_LIMIT_SPREADING)              — 449
//	│ ├ DeductionList (T_TREATY_LIMIT_DEDUCTION)              — 450
//	│ ├ CurrencyList (T_TREATY_LIMIT_ACH_PARAM)               — 450
//	│ └ IOOLimitList · RetentionList · CessionList · EPIList ·
//	│   RNMShareList · RNMSpreadedList · RNMSpreadedListRI ·
//	│   ReserveList · PLAList · CashLossList · ClaimCoopList ·
//	│   DeductionTotalList (T_TREATY_LIMIT_AMOUNT, JENIS)
//	├ TreatyGroupList (T_TREATY_LIMIT_GROUP)               — Non-Prop
//	│ └ ClassOfBusinessList (T_TREATY_LIMIT_GROUP_COB)
//	├ Reinstatement_List (T_TREATY_LIMIT_REINSTATEMENT)      — 450, Non-Prop
//	└ MDPList · PremiumEarnedList · EgnpiTotalList · MDPMinList (T_TREATY_LIMIT_MEASURE, JENIS)
//
// ⚠️ Larik yang TIDAK punya tabel pendaratan tetap hadir sebagai larik
// KOSONG — layar membaca `.length` di atasnya. Kosongnya jujur: sumbernya
// belum ada, bukan datanya yang nol.
//
// ⭐ 450 (8 Oktober 2026): kesembilan larik yang dahulu di sini kini punya
// tabel (lihat pohon di atas), jadi daftarnya kosong. `Reinstatement_List`
// yang belum pernah tersimpan tetap dibangun services saat kontrak Non-Prop
// dibuka (`SiapkanReinstatementPohon` — yang tersimpan tidak ditimpa).
var (
	larikTanpaTabelDetail = []string{}
	larikTanpaTabelLimit  = []string{}
)

// entriPeta mencari entri `PetaPendaratan` untuk satu tabel.
func entriPeta(tabel string) (Pendaratan, bool) {
	for _, p := range PetaPendaratan {
		if p.Tabel == tabel {
			return p, true
		}
	}
	return Pendaratan{}, false
}

// kolomJenis — kolom pembeda larik gabungan; tidak ada di dokumen.
const kolomJenis = "JENIS"

// bacaEntri membaca satu tabel menurut entri petanya, beserta `JENIS` bila
// tabelnya menampung beberapa larik.
func (g *Gudang) bacaEntri(ctx context.Context, tabel, masterID string) ([]barisPendaratan, error) {
	asli, ada := entriPeta(tabel)
	if !ada {
		return nil, fmt.Errorf("repository: %s tidak ada di PetaPendaratan", tabel)
	}
	// ⭐ Hanya kolom yang SUDAH terpasang — peta boleh mendahului migrasinya
	// (`449`). Tabel yang belum ada = nol baris; kolom yang belum ada
	// terbaca kosong, sama dengan `NULL`.
	p, terpasang, err := g.petaTerpasang(ctx, asli)
	if err != nil {
		return nil, err
	}
	if !terpasang {
		return nil, nil
	}
	kolom := append([]string{}, p.Kolom...)
	if len(p.LarikGabung) > 0 {
		kolom = append(kolom, kolomJenis)
	}
	return g.bacaPendaratan(ctx, p.Tabel, p.Induk != "", kolom, masterID)
}

// simpulPeta mengubah satu baris menjadi simpul berkunci ejaan dokumen.
func simpulPeta(b barisPendaratan, p Pendaratan) map[string]any {
	s := make(map[string]any, len(p.Kunci)+4)
	for i, k := range p.Kunci {
		if i < len(p.Kolom) {
			s[k] = b.Nilai[p.Kolom[i]]
		}
	}
	return s
}

// tabelPohonLimits — urutan baca; kuncinya nama tabel.
var tabelPohonLimits = []string{
	"T_TREATY_LIMITS", "T_TREATY_LIMIT_DETAIL", "T_TREATY_LIMIT_COB",
	"T_TREATY_LIMIT_ACHIEVEMENT", "T_TREATY_LIMIT_SPREADING", "T_TREATY_LIMIT_DEDUCTION",
	"T_TREATY_LIMIT_ACH_PARAM", "T_TREATY_LIMIT_REINSTATEMENT", "T_TREATY_LIMIT_AMOUNT",
	"T_TREATY_LIMIT_GROUP", "T_TREATY_LIMIT_GROUP_COB", "T_TREATY_LIMIT_MEASURE",
}

// BacaPohonLimitsPendaratan menyusun POHON tab Limits (Prop dan Non-Prop).
func (g *Gudang) BacaPohonLimitsPendaratan(ctx context.Context, masterID string) ([]map[string]any, error) {
	baris := make(map[string][]barisPendaratan, len(tabelPohonLimits))
	for _, t := range tabelPohonLimits {
		b, err := g.bacaEntri(ctx, t, masterID)
		if err != nil {
			return nil, err
		}
		baris[t] = b
	}
	return RangkaiPohonLimitsPeta(baris), nil
}

// RangkaiPohonLimits — bentuk empat tabel lama, dipertahankan untuk uji.
func RangkaiPohonLimits(limits, detail, cob, ach []barisPendaratan) []map[string]any {
	return RangkaiPohonLimitsPeta(map[string][]barisPendaratan{
		"T_TREATY_LIMITS":            limits,
		"T_TREATY_LIMIT_DETAIL":      detail,
		"T_TREATY_LIMIT_COB":         cob,
		"T_TREATY_LIMIT_ACHIEVEMENT": ach,
	})
}

// anakMenurutInduk mengelompokkan baris anak menurut `IDINDUK`, dan — bila
// tabelnya gabungan — menurut `JENIS` pula.
type anakTabel struct {
	p     Pendaratan
	biasa map[int64][]map[string]any
	jenis map[int64]map[string][]map[string]any
}

func kelompokkan(tabel string, baris []barisPendaratan) anakTabel {
	p, _ := entriPeta(tabel)
	a := anakTabel{p: p, biasa: map[int64][]map[string]any{}, jenis: map[int64]map[string][]map[string]any{}}
	for _, b := range baris {
		s := simpulPeta(b, p)
		if len(p.LarikGabung) == 0 {
			a.biasa[b.Induk] = append(a.biasa[b.Induk], s)
			continue
		}
		j := b.Nilai[kolomJenis]
		if a.jenis[b.Induk] == nil {
			a.jenis[b.Induk] = map[string][]map[string]any{}
		}
		a.jenis[b.Induk][j] = append(a.jenis[b.Induk][j], s)
	}
	return a
}

// pasangLarikGabung memasang SETIAP larik gabungan — yang tak berbaris
// sebagai larik kosong.
func pasangLarikGabung(s map[string]any, a anakTabel, induk int64) {
	for _, nama := range a.p.LarikGabung {
		s[nama] = isiAtauKosong(a.jenis[induk][nama])
	}
}

// RangkaiPohonLimitsPeta merangkai pohon dari baris per tabel.
func RangkaiPohonLimitsPeta(baris map[string][]barisPendaratan) []map[string]any {
	pLim, _ := entriPeta("T_TREATY_LIMITS")
	pDet, _ := entriPeta("T_TREATY_LIMIT_DETAIL")
	pGrp, _ := entriPeta("T_TREATY_LIMIT_GROUP")
	cob := kelompokkan("T_TREATY_LIMIT_COB", baris["T_TREATY_LIMIT_COB"])
	ach := kelompokkan("T_TREATY_LIMIT_ACHIEVEMENT", baris["T_TREATY_LIMIT_ACHIEVEMENT"])
	sebar := kelompokkan("T_TREATY_LIMIT_SPREADING", baris["T_TREATY_LIMIT_SPREADING"])
	deduksi := kelompokkan("T_TREATY_LIMIT_DEDUCTION", baris["T_TREATY_LIMIT_DEDUCTION"])
	param := kelompokkan("T_TREATY_LIMIT_ACH_PARAM", baris["T_TREATY_LIMIT_ACH_PARAM"])
	reinst := kelompokkan("T_TREATY_LIMIT_REINSTATEMENT", baris["T_TREATY_LIMIT_REINSTATEMENT"])
	amt := kelompokkan("T_TREATY_LIMIT_AMOUNT", baris["T_TREATY_LIMIT_AMOUNT"])
	gcob := kelompokkan("T_TREATY_LIMIT_GROUP_COB", baris["T_TREATY_LIMIT_GROUP_COB"])
	ukur := kelompokkan("T_TREATY_LIMIT_MEASURE", baris["T_TREATY_LIMIT_MEASURE"])

	anakDetail := map[int64][]barisPendaratan{}
	for _, d := range baris["T_TREATY_LIMIT_DETAIL"] {
		anakDetail[d.Induk] = append(anakDetail[d.Induk], d)
	}
	anakGrup := map[int64][]barisPendaratan{}
	for _, gr := range baris["T_TREATY_LIMIT_GROUP"] {
		anakGrup[gr.Induk] = append(anakGrup[gr.Induk], gr)
	}

	out := []map[string]any{}
	for _, l := range baris["T_TREATY_LIMITS"] {
		s := simpulPeta(l, pLim)
		rinci := []map[string]any{}
		for _, d := range anakDetail[l.ID] {
			sd := simpulPeta(d, pDet)
			// ⛔ Larik KOSONG, bukan nihil: `TabLimitsProp` membaca
			// `.length` di atasnya.
			sd["COBList"] = isiAtauKosong(cob.biasa[d.ID])
			sd["AchievementLists"] = isiAtauKosong(ach.biasa[d.ID])
			sd["SpreadingList"] = isiAtauKosong(sebar.biasa[d.ID])
			sd["DeductionList"] = isiAtauKosong(deduksi.biasa[d.ID])
			sd["CurrencyList"] = isiAtauKosong(param.biasa[d.ID])
			pasangLarikGabung(sd, amt, d.ID)
			for _, nama := range larikTanpaTabelDetail {
				sd[nama] = []map[string]any{}
			}
			rinci = append(rinci, sd)
		}
		s["Detail"] = rinci

		grup := []map[string]any{}
		for _, gr := range anakGrup[l.ID] {
			sg := simpulPeta(gr, pGrp)
			sg["ClassOfBusinessList"] = isiAtauKosong(gcob.biasa[gr.ID])
			grup = append(grup, sg)
		}
		s["TreatyGroupList"] = grup
		pasangLarikGabung(s, ukur, l.ID)
		s["Reinstatement_List"] = isiAtauKosong(reinst.biasa[l.ID])
		for _, nama := range larikTanpaTabelLimit {
			s[nama] = []map[string]any{}
		}
		out = append(out, s)
	}
	return out
}

// Larik total AKAR yang tab Limits Non-Prop tampilkan (`Total All Layers`).
var larikTotalLimitNP = []string{
	"TotalLimitIOONP", "TotalLimitDeductblNP", "TotalLimitPremiEarnNP", "TotalLimitMDPNP",
}

// BacaLimitsAkarPendaratan membaca larik AKAR tab Limits Non-Prop:
// `LimitSummaryList` dan keempat `TotalLimit…NP`. `TotalLimitsROL` skalar
// akar — services mengisinya dari `T_TREATY_REVISION`.
func (g *Gudang) BacaLimitsAkarPendaratan(ctx context.Context, masterID string) (models.LimitsAkar, error) {
	akar := models.LimitsAkar{
		LimitSummaryList: []map[string]any{},
		Total:            map[string][]map[string]any{},
	}
	ring, err := g.bacaEntri(ctx, "T_TREATY_LIMIT_SUMMARY", masterID)
	if err != nil {
		return models.LimitsAkar{}, err
	}
	pRing, _ := entriPeta("T_TREATY_LIMIT_SUMMARY")
	for _, b := range ring {
		akar.LimitSummaryList = append(akar.LimitSummaryList, simpulPeta(b, pRing))
	}
	tot, err := g.bacaEntri(ctx, "T_TREATY_TOTAL", masterID)
	if err != nil {
		return models.LimitsAkar{}, err
	}
	pTot, _ := entriPeta("T_TREATY_TOTAL")
	for _, nama := range larikTotalLimitNP {
		akar.Total[nama] = []map[string]any{}
	}
	for _, b := range tot {
		j := b.Nilai[kolomJenis]
		if _, dipakai := akar.Total[j]; dipakai {
			akar.Total[j] = append(akar.Total[j], simpulPeta(b, pTot))
		}
	}
	return akar, nil
}

func isiAtauKosong(v []map[string]any) []map[string]any {
	if v == nil {
		return []map[string]any{}
	}
	return v
}
