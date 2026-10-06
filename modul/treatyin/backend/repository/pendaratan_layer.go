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
//	MDP, MDP kedua           <- Limits[].MDPList[]            (T_TREATY_LIMIT_MEASURE, belum ada)
//	EARN_PREMIUM, keduanya   <- Limits[].PremiumEarnedList[]  (T_TREATY_LIMIT_MEASURE, belum ada)
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
	return RangkaiLayer(limits, detail, cob), nil
}

// RangkaiLayer merangkai ketiga tabel menjadi baris layar.
//
// ⛔ DIPISAH dari pembacaannya supaya ia dapat diuji TANPA Oracle — pola
// yang sama dengan `LayerDariDokumen` yang ia gantikan. Perangkaian adalah
// tempat kesalahan bentuk hidup; kueri hanya mengambil baris.
func RangkaiLayer(limits, detail, cob []barisPendaratan) []models.BarisLayerWarisan {
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

// Kunci pohon tab Limits proporsional, berpasangan kunci tampil <- kolom.
//
// ⛔ Nama KUNCI-nya tetap ejaan Pega (`TreatyGroup`, `QSPct`), bukan nama
// kolom Oracle: layar `TabLimitsProp` mengikat medannya dengan ejaan itu,
// dan menukar sumber data tidak boleh menukar nama medan di layar.
var (
	pohonKunciLimit = [][2]string{
		{"TreatyType", "TREATYTYPE"}, {"TreatyTypeID", "TREATYTYPEID"},
	}
	pohonKunciDetail = [][2]string{
		{"TreatyGroup", "TREATYGROUP"}, {"TreatyGroupID", "TREATYGROUPID"},
		{"TreatyType", "TREATYTYPE"}, {"QSPct", "QSPCT"}, {"Surplus", "SURPLUS"},
		{"RetentionPct", "RETENTIONPCT"}, {"CessionPct", "CESSIONPCT"},
		{"CurrencyRSMD", "CURRENCYRSMD"}, {"RSMDLimit", "RSMDLIMIT"},
		{"CurrencyEarthquake", "CURRENCYEARTHQUAKE"}, {"Earthquake", "EARTHQUAKE"},
		{"CurrencyFloodJab", "CURRENCYFLOODJAB"}, {"FloodJab", "FLOODJAB"},
		{"CurrencyFloodNat", "CURRENCYFLOODNAT"}, {"FloodNation", "FLOODNATION"},
		{"RIOGR", "RIOGR"}, {"RIONR", "RIONR"},
		{"PremiumReservePct", "PREMIUMRESERVEPCT"}, {"ProfitCommision", "PROFITCOMMISION"},
		{"ProfitME", "PROFITME"}, {"ProfitYDCF", "PROFITYDCF"},
		{"LowerBand", "LOWERBAND"}, {"UpperBand", "UPPERBAND"},
		{"ReisuredParticipant", "REISUREDPARTICIPANT"}, {"Periode", "PERIODE"},
	}
	pohonKunciCOB = [][2]string{{"ClassOfBusiness", "CLASSOFBUSINESS"}}
	pohonKunciAch = [][2]string{
		{"Quarter", "QUARTER"}, {"QUARTERYEAR", "QUARTERYEAR"}, {"Currency", "CURRENCY"},
		{"PREMIUM", "PREMIUM"}, {"RICOMM", "RICOMM"}, {"BROKERAGE", "BROKERAGE"},
		{"NETPREMIUM", "NETPREMIUM"}, {"PaidClaim", "PAIDCLAIM"}, {"CASHCALL", "CASHCALL"},
		{"OutstandingClaim", "OUTSTANDINGCLAIM"}, {"IncuredClaim", "INCUREDCLAIM"},
		{"Total", "TOTAL"}, {"LossRatio", "LOSSRATIO"},
	}
)

func kolomDari(pasang [][2]string) []string {
	out := make([]string, 0, len(pasang))
	for _, p := range pasang {
		out = append(out, p[1])
	}
	return out
}

func simpulDari(n map[string]string, pasang [][2]string) map[string]any {
	s := make(map[string]any, len(pasang)+1)
	for _, p := range pasang {
		s[p[0]] = n[p[1]]
	}
	return s
}

// BacaPohonLimitsPendaratan menyusun POHON tab Limits proporsional dari
// tabel pendaratan: `T_TREATY_LIMITS -> T_TREATY_LIMIT_DETAIL -> COBList`,
// beserta `AchievementLists`.
//
// ⚠️ Sepuluh grid lain di dalam `Detail` (`IOOLimitList`, `RetentionList`,
// `CessionList`, `DeductionList`, `DeductionTotalList`, `ReserveList`,
// `PLAList`, `CashLossList`, `ClaimCoopList`, `EPIList`) TIDAK punya tabel
// di antara keempat belas, sehingga simpulnya tidak memuatnya. Layar
// merendernya sebagai `No items`, dan itu keadaan yang jujur: sumbernya
// memang belum ada, bukan datanya yang kosong.
func (g *Gudang) BacaPohonLimitsPendaratan(ctx context.Context, masterID string) ([]map[string]any, error) {
	limits, err := g.bacaPendaratan(ctx, "T_TREATY_LIMITS", false, kolomDari(pohonKunciLimit), masterID)
	if err != nil {
		return nil, err
	}
	detail, err := g.bacaPendaratan(ctx, "T_TREATY_LIMIT_DETAIL", true, kolomDari(pohonKunciDetail), masterID)
	if err != nil {
		return nil, err
	}
	cob, err := g.bacaPendaratan(ctx, "T_TREATY_LIMIT_COB", true, kolomDari(pohonKunciCOB), masterID)
	if err != nil {
		return nil, err
	}
	ach, err := g.bacaPendaratan(ctx, "T_TREATY_LIMIT_ACHIEVEMENT", true, kolomDari(pohonKunciAch), masterID)
	if err != nil {
		return nil, err
	}

	return RangkaiPohonLimits(limits, detail, cob, ach), nil
}

// RangkaiPohonLimits merangkai keempat tabel menjadi pohon tiga tingkat.
// Dipisah dengan alasan yang sama seperti `RangkaiLayer`.
func RangkaiPohonLimits(limits, detail, cob, ach []barisPendaratan) []map[string]any {
	anakCOB := map[int64][]map[string]any{}
	for _, c := range cob {
		anakCOB[c.Induk] = append(anakCOB[c.Induk], simpulDari(c.Nilai, pohonKunciCOB))
	}
	anakAch := map[int64][]map[string]any{}
	for _, a := range ach {
		anakAch[a.Induk] = append(anakAch[a.Induk], simpulDari(a.Nilai, pohonKunciAch))
	}
	anakDetail := map[int64][]barisPendaratan{}
	for _, d := range detail {
		anakDetail[d.Induk] = append(anakDetail[d.Induk], d)
	}

	out := []map[string]any{}
	for _, l := range limits {
		s := simpulDari(l.Nilai, pohonKunciLimit)
		rinci := []map[string]any{}
		for _, d := range anakDetail[l.ID] {
			sd := simpulDari(d.Nilai, pohonKunciDetail)
			// ⛔ Larik KOSONG, bukan nihil: `TabLimitsProp` membaca
			// `.length` di atasnya.
			sd["COBList"] = isiAtauKosong(anakCOB[d.ID])
			sd["AchievementLists"] = isiAtauKosong(anakAch[d.ID])
			rinci = append(rinci, sd)
		}
		s["Detail"] = rinci
		out = append(out, s)
	}
	return out
}

func isiAtauKosong(v []map[string]any) []map[string]any {
	if v == nil {
		return []map[string]any{}
	}
	return v
}
