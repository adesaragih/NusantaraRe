package loader

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

// TestSkemaUkuran - angka rancangan yang sudah diukur DUA cara (02-10-2026): pengurai
// xlsx lembar Kolom vs pencacah baris `DDL-tabel-flat-draf.sql` - 78 tabel, 1.329
// kolom, selisih nol dua arah. 148 jalur = 151 baris lembar Jalur Sumber - baris
// catatan - baris kepala - T_WORK_POLIS ("objek kerja Pega", tanpa jalur dokumen).
func TestSkemaUkuran(t *testing.T) {
	n := 0
	for _, ks := range skemaTabel {
		n += len(ks)
	}
	// Sesudah amandemen butir 70 (amandemen.go): +1 tabel (T_ADDITIONALSHIP, 10 kolom),
	// +3 kolom di tabel lama, +1 jalur. Bagian bangkitan tetap 78 / 1.329 / 148.
	am := 0
	for _, ks := range amandemenKolom {
		am += len(ks)
	}
	for _, ks := range amandemenTabel {
		am += len(ks)
	}
	// Butir 72: +48 kolom penunjuk teks mentah (amandemenPenunjuk) -> 61 kolom amandemen.
	// Butir 76.2: +1 T_WORK_POLIS.LINI -> 62 (76.1/76.4 mengganti tipe/nama, tidak menambah).
	// Tiket 35 (A110): +3 kolom T_BUILDINGCONSTRUCTION -> 65. Tiket 38 (A130): +18 T_SURROUNDINGRISK -> 83.
	// Tiket 39 (A132): +6 T_PROPERTYITEMLIST -> 89. Tiket 41 (A142): +1 tabel T_FEALIST (15 kolom) + 1 jalur -> 104.
	// Tiket 42 (A145): +5 T_LISTCAUSEOFLOSS -> 109.
	if am != 109 || len(amandemenPenunjuk) != 48 || len(amandemenTabel) != 2 || len(amandemenJalur) != 2 {
		t.Fatalf("amandemen %d kolom (%d penunjuk) / %d tabel / %d jalur, mau 109 (48) / 2 / 2", am, len(amandemenPenunjuk), len(amandemenTabel), len(amandemenJalur))
	}
	if len(skemaTabel) != 78+2 || n != 1329+am || len(jalurSumber) != 148+2 || len(warisMataUang) != 8 || len(penunjukKandidat) != 50 {
		t.Fatalf("%d tabel, %d kolom, %d jalur, %d waris, %d penunjuk; mau 80/1438/150/8/50",
			len(skemaTabel), n, len(jalurSumber), len(warisMataUang), len(penunjukKandidat))
	}
	unknown := 0
	for _, ks := range skemaTabel {
		for _, k := range ks {
			if k.bawaanUnknown() {
				unknown++
			}
		}
	}
	// K-069 usulan 10: 24 pembawa mata uang NOT NULL DEFAULT 'UNKNOWN' (a=7 b=9 c=8).
	if unknown != 24 {
		t.Errorf("%d kolom ber-DEFAULT 'UNKNOWN', mau 24", unknown)
	}
}

// TestIndukGanda - spec 11 menyebut "dua belas tabel berinduk ganda", yang terparah
// T_COVERAGELIST dengan lima induk; semuanya wajib punya PARENT_TABLE.
//
// ⚠️ Dua sumber, dua angka - dicatat, tidak dipilih (02-10-2026): lembar Daftar
// Relasi memberi 12, lembar Jalur Sumber 13. Selisihnya T_FR_POLICY (induk
// T_FR_NETPERCURRENCY dan T_FR_OFFEREDPAYMENT): ketiga tabel wadah yang ditambahkan
// 24-09 (T_FR_FACOFFERLIST, T_FR_OBJECT, T_FR_POLICY) TIDAK ADA di Daftar Relasi
// sama sekali, padahal DDL draf memberi T_FR_POLICY PARENT_TABLE. Flatten tidak
// memilih: ia mengisi PARENT_TABLE di setiap tabel yang punya kolom itu (23).
func TestIndukGanda(t *testing.T) {
	induk := map[string]map[string]bool{}
	for _, j := range jalurSumber {
		if induk[j.tabel] == nil {
			induk[j.tabel] = map[string]bool{}
		}
		induk[j.tabel][j.induk] = true
	}
	var ganda []string
	for tb, s := range induk {
		if len(s) > 1 {
			ganda = append(ganda, tb)
			if !punyaKolom(tb, kolomTabelInduk) {
				t.Errorf("%s berinduk ganda tanpa PARENT_TABLE", tb)
			}
		}
	}
	sort.Strings(ganda)
	if len(ganda) != 13 || !contains(ganda, "T_FR_POLICY") {
		t.Errorf("%d tabel berinduk ganda menurut Jalur Sumber, mau 13 (12 Daftar Relasi + T_FR_POLICY): %v", len(ganda), ganda)
	}
	n := 0
	for tb := range skemaTabel {
		if punyaKolom(tb, kolomTabelInduk) {
			n++
		}
	}
	if n != 23 {
		t.Errorf("%d tabel ber-PARENT_TABLE, mau 23", n)
	}
	if len(induk["T_COVERAGELIST"]) != 5 {
		t.Errorf("T_COVERAGELIST %d induk, mau 5", len(induk["T_COVERAGELIST"]))
	}
}

// TestJalurIndukAdalahAwalan - invarian yang membuat ErrIndukTakTetap tak tercapai
// lewat dokumen: tabel induk tiap jalur SELALU tabel jalur awalannya (jalur minus
// ruas terakhir), jadi baris yang lahir di jalur itu pasti berinduk tabel yang
// Jalur Sumber tetapkan. Bila workbook berubah dan invarian ini patah, uji ini yang
// gagal lebih dulu.
func TestJalurIndukAdalahAwalan(t *testing.T) {
	tabelJalur := map[string]string{"": "T_GENERAL_POLIS"}
	for _, j := range jalurSumber {
		tabelJalur[j.jalur] = j.tabel
	}
	for _, j := range jalurSumber {
		if j.jalur == "" {
			continue
		}
		awal := ""
		if i := strings.LastIndexByte(j.jalur, '/'); i >= 0 {
			awal = j.jalur[:i]
		}
		if tabelJalur[awal] != j.induk {
			t.Errorf("%s: awalan %q bertabel %q, Jalur Sumber menetapkan induk %s", j.jalur, awal, tabelJalur[awal], j.induk)
		}
	}
}

// TestKedalamanRancangan - kedalaman larik terdalam yang dapat dicapai lewat jalur
// rancangan (ruas yang tabelnya ber-SEQ_NO). Hasilnya 8 - sama dengan kedalaman
// rowdata maksimum BAHAN §1 yang diukur di XML. Lebih dari 8 karena itu tak
// tercapai lewat jalur rancangan; ErrTerlaluDalam penjaga bila rancangan berubah.
func TestKedalamanRancangan(t *testing.T) {
	tabelJalur := map[string]string{}
	for _, j := range jalurSumber {
		tabelJalur[j.jalur] = j.tabel
	}
	maks := 0
	for _, j := range jalurSumber {
		ruas, d := strings.Split(j.jalur, "/"), 0
		for i := range ruas {
			if tb, ada := tabelJalur[strings.Join(ruas[:i+1], "/")]; ada && punyaKolom(tb, kolomSeq) {
				d++
			}
		}
		if d > maks {
			maks = d
		}
	}
	if maks != kedalamanMaks {
		t.Errorf("kedalaman larik rancangan %d, mau %d", maks, kedalamanMaks)
	}
}

// TestSetiapKolomBerasal - tidak satu pun kolom tanpa asal tertulis. Jumlah per asal
// dikunci supaya perubahan workbook atau aturan terlihat sebagai selisih angka.
func TestSetiapKolomBerasal(t *testing.T) {
	jumlah := map[string]int{}
	kosong := map[string]int{}
	for tb, ks := range skemaTabel {
		for _, k := range ks {
			asal, alasan := asalKolom(tb, k)
			if asal == "" {
				t.Errorf("%s.%s: %s", tb, k.nama, alasan)
				continue
			}
			jumlah[asal]++
			if asal == asalKosong {
				kosong[alasan]++
			}
		}
	}
	// Dihitung dua jalan (02-10-2026):
	//  per asal  - Flatten 320 = IDPEGA 78 + COB_GROUP 77 + SEQ_NO 49 + ROW_UID 49 +
	//              PARENT_TABLE 23 + SRC_PATH 23 + JENIS/NO_WORK 2 + V-30 5 + V-49 6 +
	//              CURRENCY_CODE dari Name 1 + warisan 7; repository 164 = ID 78 +
	//              PARENT_ID 76 + 10 bernama; kosong 32 = V-47 29 + V-37 2 + T_PROPERTY 1.
	//  per jenis - medan 813 = data 806 + PII 45 + kunci 6 (857 non-sistem) - 44
	//              non-sistem yang bukan dari medan (V-47 29, V-37 2, V-49 6, V-30 5,
	//              CURRENCY_CODE T_CURRENCYLIST 1, PROD_KE 1). 813+320+164+32 = 1.329.
	// Amandemen butir 70: medan +4 (COVERAGE_INITIAL, DWT, GRT, NRT), Flatten +7 (CURRENCY_REF_ID,
	// ADDITIONAL_SHIP_REF_ID, POLICY_TSI, IDPEGA/COB_GROUP/SEQ_NO/ROW_UID T_ADDITIONALSHIP),
	// repository +2 (ID, PARENT_ID). 817 + 327 + 166 + 32 = 1.342.
	// Butir 72: medan +48 (kolom penunjuk teks mentah). 865 + 327 + 166 + 32 = 1.390.
	// Butir 76: repository +1 (LINI); empat kolom gabungan tetap repository. 865 + 327 + 167 + 32 = 1.391.
	// Tiket 35 (A110): medan +3 (PartitionType, SupportWallType, OthersType). 868 + 327 + 167 + 32 = 1.394.
	// Tiket 38 (A130): medan +18 (empat sisi x 4, FloodArea, HousekeepingRemark). 886 + 327 + 167 + 32 = 1.412.
	// Tiket 39 (A132): medan +6 (PropertyYear, Unit, Condition, Year, NoOfTree, AreaHectar). 892 + 327 + 167 + 32 = 1.418.
	// Tiket 41 (A142): tabel T_FEALIST - medan +9, kolom sistem Flatten +4 (ID, PARENT_ID, SEQ_NO, ROW_UID), repository +2
	// (IDPEGA, COB_GROUP). 901 + 331 + 169 + 32 = 1.433.
	// Tiket 42 (A145): medan +5 (DateOfLoss, LossObject, Amount, PreventionOfLoss, CauseOfLoss). 906 + 331 + 169 + 32 = 1.438.
	mau := map[string]int{asalMedan: 906, asalFlatten: 331, asalRepository: 169, asalKosong: 32}
	for a, n := range mau {
		if jumlah[a] != n {
			t.Errorf("%s: %d kolom, mau %d", a, jumlah[a], n)
		}
	}
	if kosong[alasanFKV47] != 29 || kosong[kolomTanpaRumus["TSI_TOP_RISK"]] != 2 {
		t.Errorf("sengaja kosong %v, mau 29 V-47 + 2 V-37 + 1 T_PROPERTY", kosong)
	}
}

// TestMedanAkarTidakBertabrakan - V-48: medan akar dicari di T_GENERAL_POLIS lebih
// dulu, lalu T_WORK_POLIS. Keduanya tidak boleh berbagi FIELD ASLI.
func TestMedanAkarTidakBertabrakan(t *testing.T) {
	umum := map[string]bool{}
	for _, k := range skemaTabel["T_GENERAL_POLIS"] {
		umum[k.medan] = true
	}
	n := 0
	for _, k := range skemaTabel["T_WORK_POLIS"] {
		if k.medan == "" {
			continue
		}
		n++
		if umum[k.medan] {
			t.Errorf("FIELD ASLI %s ada di T_GENERAL_POLIS dan T_WORK_POLIS", k.medan)
		}
	}
	if n != 6 {
		t.Errorf("T_WORK_POLIS %d kolom ber-FIELD ASLI, mau 6 (V-48)", n)
	}
}

// TestRancanganTanpaMedanYangDibuang - kaidah buang (aturan.go) tidak boleh mengenai
// kolom yang rancangannya simpan: nol FIELD ASLI ber-awalan px/pz/py/EDMOld atau
// berakhiran Old. Pasangan TestAlasanDaun.
func TestRancanganTanpaMedanYangDibuang(t *testing.T) {
	for tb, ks := range skemaTabel {
		for _, k := range ks {
			if k.medan != "" && alasanDaun(k.medan) != "" {
				t.Errorf("%s.%s FIELD ASLI %s terkena kaidah %q", tb, k.nama, k.medan, alasanDaun(k.medan))
			}
		}
	}
}

// TestAlasanDaun - uji instrumen kaidah buang dengan jawaban yang sudah diketahui.
func TestAlasanDaun(t *testing.T) {
	for medan, mau := range map[string]string{
		"pxObjClass": alasanMeta, "pzInsKey": alasanMeta, "pyExpanded": alasanPy,
		"TSIOld": alasanSufiksOld, "EDMOldPremi": alasanEDMOld,
		"OldTSI": "", "OldPolicyNo": "", "TSI": "", "Older": "",
	} {
		if got := alasanDaun(medan); got != mau {
			t.Errorf("%s: %q, mau %q", medan, got, mau)
		}
	}
}

// TestAturanMenunjukKolomYangAda - setiap kolom yang disebut aturan.go ada di skema.
func TestAturanMenunjukKolomYangAda(t *testing.T) {
	for tb, kol := range medanID {
		if !punyaKolom(tb, kol) {
			t.Errorf("V-49 %s.%s tidak ada", tb, kol)
		}
	}
	for tb, m := range medanGanda {
		for _, kol := range m {
			if !punyaKolom(tb, kol) {
				t.Errorf("medanGanda %s.%s tidak ada", tb, kol)
			}
		}
	}
	for tk := range kolomV30 {
		tb, kol, _ := strings.Cut(tk, ".")
		if !punyaKolom(tb, kol) {
			t.Errorf("V-30 %s tidak ada", tk)
		}
	}
	for l := range lipat {
		if _, ada := skemaTabel[l.tabel]; !ada {
			t.Errorf("lipatan ke %s: tabel tidak ada", l.tabel)
		}
	}
	for _, tb := range []string{tabelSkoring, tabelFaktor, tabelOpsi} {
		if _, ada := skemaTabel[tb]; !ada {
			t.Errorf("tabel V-30 %s tidak ada", tb)
		}
	}
}

func contains(s []string, x string) bool {
	for _, v := range s {
		if v == x {
			return true
		}
	}
	return false
}

// TestAmandemenPenunjukKelompok - 48 kolom penunjuk butir 72 per kelompok, dicocokkan
// dengan lembar Kandidat Hapus (bangkitan): R1 13, R3 32 (induk langsung 16 + leluhur
// 16 - pemisahannya hanya label laporan), berselisih 2, tak tercantum 1
// (T_FR_PERSONLIST.IdxPerson); tiga kunci lembar lainnya (R3b) sudah berkolom kunci
// di skema bangkitan, jadi tidak diamandemen.
func TestAmandemenPenunjukKelompok(t *testing.T) {
	kel := map[string]int{}
	dipakai := map[string]bool{}
	for _, k := range amandemenPenunjuk {
		key := k.tabel + "." + k.medan
		dipakai[key] = true
		v, ada := penunjukKandidat[key]
		switch {
		case !ada:
			kel["tak tercantum "+key]++
		case strings.Contains(v.status, " | "):
			kel["berselisih"]++
		default:
			kel[v.kategori]++
		}
	}
	mau := map[string]int{"R1 indeks-diri": 13, "R3 penunjuk leluhur": 32, "berselisih": 2, "tak tercantum T_FR_PERSONLIST.IdxPerson": 1}
	if !reflect.DeepEqual(kel, mau) {
		t.Errorf("kelompok %v\nmau %v", kel, mau)
	}
	for k := range penunjukKandidat {
		if dipakai[k] {
			continue
		}
		tb, medan, _ := strings.Cut(k, ".")
		if medanKolom(tb, medan) == "" {
			t.Errorf("%s tidak diamandemen dan tidak berkolom di skema bangkitan", k)
		}
	}
}

// Status lembar Kandidat Hapus (catatan rancangan, bangkitan).
const (
	statusDihapus = "SUDAH DIHAPUS"
	statusJadiFK  = "SUDAH JADI FK (V-47)"
)

// TestPenunjukKandidat - isi bangkitan lembar Kandidat Hapus yang dirujuk register
// butir 70-71: 50 kunci, 13 R1 indeks-diri, dua kunci berstatus berselisih (disimpan
// keduanya, tidak dipilih).
func TestPenunjukKandidat(t *testing.T) {
	r1 := 0
	for _, v := range penunjukKandidat {
		if v.kategori == "R1 indeks-diri" {
			r1++
		}
	}
	if len(penunjukKandidat) != 50 || r1 != 13 {
		t.Errorf("%d kunci, %d R1; mau 50 / 13", len(penunjukKandidat), r1)
	}
	for _, k := range []string{"T_FR_ANEKALIST.IdxOccupation", "T_FR_DEDUCTIBLELIST.IndexProperty"} {
		if v := penunjukKandidat[k]; v.status != statusDihapus+" | "+statusJadiFK {
			t.Errorf("%s status %+v", k, v)
		}
	}
}

func punyaKolom(tb, kol string) bool {
	for _, k := range skemaTabel[tb] {
		if k.nama == kol {
			return true
		}
	}
	return false
}

// TestAmandemenWorkPolis - butir 76: T_WORK_POLIS rancangan selaras dengan tabel yang ADA
// (premiumlistlife 050/059/063, K-064). Dihitung dua jalan: daftar tabel anak dari
// jalurSumber (init) vs daftar tertulis di bawah - keduanya 10.
func TestAmandemenWorkPolis(t *testing.T) {
	tipe := func(tb, nama string) string {
		if i := indeksKolom(tb, nama); i >= 0 {
			return skemaTabel[tb][i].tipe
		}
		return "(tidak ada)"
	}
	for _, tb := range []string{"T_WORK_POLIS", "T_GENERAL_POLIS"} {
		if tipe(tb, "ID") != "VARCHAR2(32)" {
			t.Errorf("%s.ID %s, mau VARCHAR2(32) (butir 76.1)", tb, tipe(tb, "ID"))
		}
	}
	anak := []string{"T_CARGOLIST", "T_CEDINGCEDANTLIST", "T_CURRENCYLIST", "T_FACRETRODETAILS", "T_FACRETROLIST",
		"T_LOCATIONLIST", "T_PERSONLIST", "T_QUOTATIONDATA", "T_SCORINGRISK", "T_VEHICLELIST"}
	if strings.Join(anakIDKasus, ",") != strings.Join(anak, ",") {
		t.Errorf("anak ber-PARENT_ID kasus %v, mau %v", anakIDKasus, anak)
	}
	for _, tb := range anak {
		if tipe(tb, "PARENT_ID") != "VARCHAR2(32)" {
			t.Errorf("%s.PARENT_ID %s", tb, tipe(tb, "PARENT_ID"))
		}
	}
	// Tabel lain tetap surrogate NUMBER - penggantian tidak merembes.
	if tipe("T_COVERAGELIST", "ID") != "NUMBER" || tipe("T_SHIP", "PARENT_ID") != "NUMBER" {
		t.Error("ID/PARENT_ID tabel lain ikut berubah")
	}
	for lama, baru := range map[string]string{"POSISI": "POSITION VARCHAR2(255)", "STATUS_PROSES": "STATUS_WORK VARCHAR2(255)",
		"TGL_INPUT": "TGL_CREATE DATE", "USERNAME": "CREATE_OP VARCHAR2(64)", "": "LINI VARCHAR2(255)"} {
		n := strings.Fields(baru)
		if lama != "" && tipe("T_WORK_POLIS", lama) != "(tidak ada)" {
			t.Errorf("kolom rancangan %s masih ada (digabung butir 76.4)", lama)
		}
		if tipe("T_WORK_POLIS", n[0]) != n[1] {
			t.Errorf("T_WORK_POLIS.%s %s, mau %s", n[0], tipe("T_WORK_POLIS", n[0]), n[1])
		}
		if asal, _ := asalKolom("T_WORK_POLIS", skemaTabel["T_WORK_POLIS"][indeksKolom("T_WORK_POLIS", n[0])]); asal != asalRepository {
			t.Errorf("T_WORK_POLIS.%s asal %q, mau repository", n[0], asal)
		}
	}
	if _, alasan := asalKolom("T_WORK_POLIS", skemaTabel["T_WORK_POLIS"][indeksKolom("T_WORK_POLIS", "ID")]); !strings.Contains(alasan, "NO_WORK") {
		t.Errorf("asal T_WORK_POLIS.ID %q, mau aturan butir 76.1 (NO_WORK)", alasan)
	}
	if n := len(skemaTabel["T_WORK_POLIS"]); n != 19 {
		t.Errorf("T_WORK_POLIS %d kolom, mau 19 (18 rancangan + LINI)", n)
	}
}
