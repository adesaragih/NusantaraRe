package acceptance

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

// kasusTangga - kasus uji: jalur properti Pega → nilai. Jalur yang tidak ada
// dibaca kosong, seperti clipboard Pega.
type kasusTangga map[string]string

func (k kasusTangga) Nilai(jalur string) (string, bool) {
	v, ada := k[jalur]
	return v, ada
}

// dengan - salinan k dengan beberapa jalur ditimpa.
func (k kasusTangga) dengan(pasangan ...string) kasusTangga {
	b := kasusTangga{}
	for j, v := range k {
		b[j] = v
	}
	for i := 0; i < len(pasangan); i += 2 {
		b[pasangan[i]] = pasangan[i+1]
	}
	return b
}

// angkaLimit - angka CSV ekspor: titik pemisah ribuan (README testdata/limit).
var angkaLimit = regexp.MustCompile(`^\d{1,3}(\.\d{3})+$|^\d+$`)

// muatLimit - kelima tabel bentuk A dari fixture `testdata/limit`.
func muatLimit(t *testing.T) TabelLimit {
	t.Helper()
	tabel := TabelLimit{}
	for _, nama := range []NamaTabel{TabelProperty, TabelPropertyNonPreferred,
		TabelPropertyPreferredCommercial, TabelEngineering, TabelNonPropEng} {
		baris := bacaCSVLimit(t, filepath.Join(folderLimit, string(nama)+".csv"))
		kol := map[string]int{}
		for i, k := range baris[0] {
			kol[k] = i
		}
		tanpaTitik := func(s string) string {
			if !angkaLimit.MatchString(s) {
				t.Fatalf("%s: angka %q bukan bentuk ekspor yang dikenal", nama, s)
			}
			return strings.ReplaceAll(s, ".", "")
		}
		for _, b := range baris[1:] {
			tabel[nama] = append(tabel[nama], BarisLimit{Jabatan: Jabatan(b[kol["JABATAN"]]), TeamGroup: b[kol["TEAM_GROUP"]],
				LimitBottom: desimal(t, tanpaTitik(b[kol["LIMIT_BOTTOM"]])), LimitBottom2: desimal(t, tanpaTitik(b[kol["LIMIT_BOTTOM2"]]))})
		}
	}
	return tabel
}

// Jalur properti yang dibaca tangga (GetLimitAkseptasi_ActFlow, flow Decision23).
const (
	jTSI        = "pyWorkPage.OfferFacIn.TotalTSINusaRe"
	jTopRisk    = "pyWorkPage.OfferFacIn.TotalTSITopRisk"
	jAdaTop     = "pyWorkPage.IsAdaTopRisk"
	jStatus     = "pyWorkPage.OfferFacIn.QuotationData.StatusBusiness"
	jTSILama    = "pyWorkPage.OfferFacIn.OldData.TotalTSINusaRe"
	jTopLama    = "pyWorkPage.OfferFacIn.OldData.TotalTSITopRisk"
	jTeamGroup  = "pyWorkPage.OfferFacIn.QuotationData.TeamGroup"
	jBisnis     = "pyWorkPage.OfferFacIn.QuotationData.BusinessType"
	jBisnisLama = "pyWorkPage.OfferFacIn.QuotationData.BusinessOldId"
	jPreferred  = "pyWorkPage.OfferFacIn.IsPreferredRisk"
	jEngineer   = "pyWorkPage.Quotation.BusinessType"
	jBanding    = "pyWorkPage.OfferFacIn.IsBanding"
	jReject     = "pyWorkPage.OfferFacIn.IsFlagReject"
	jAntrean    = "pyWorkPage.PositionNote"
	jB2B        = "pyWorkPage.OfferFacIn.IsB2B"
)

// fireTG1 - kasus FIRE Preferred Risk team group 1, TSI 300 miliar. Tabel
// M_LIMIT_PROPERTYY tg 1 (LIMIT_BOTTOM; LIMIT_BOTTOM2): LEADER 0, JUW_B 0,
// UNDERWRITER 1, SENIORUW ×2 178.500.000.001, KADIVTEKNIK 255.000.000.001
// (289.000.000.001), MANAGERTEKNIK dan KADIVFACULTATIVE 255.000.000.001,
// DIREKTURMARKETING 331.500.000.001, DIREKTURTEKNIK 357.500.000.001.
var fireTG1 = kasusTangga{
	jTSI: "300000000000", jTopRisk: "0", jAdaTop: "false", jStatus: "1",
	jTeamGroup: "1", jBisnis: "Fire", jBisnisLama: "22", jPreferred: "Preferred Risk",
	jBanding: "false", jReject: "false", jAntrean: "ReasFacInSeniorUnderwriting",
}

func penggunaBiasa(j Jabatan) Pengguna { return Pengguna{Jabatan: j} }

func desimal(t *testing.T, s string) *apd.Decimal {
	t.Helper()
	v, err := utils.ParseDecimal(s)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

// TestNextSatuLangkah - satu panggilan = satu transisi, tidak ada rantai: hasil
// tiap baris adalah jabatan BERIKUTNYA saja, dan memanggil lagi dengan jabatan
// itu memberi langkah sesudahnya.
func TestNextSatuLangkah(t *testing.T) {
	tabel := muatLimit(t)
	tsi400 := fireTG1.dengan(jTSI, "400000000000")
	for _, u := range []struct {
		nama     string
		kasus    kasusTangga
		pengguna Jabatan
		mau      Transisi
	}{
		// Kandidat LIMIT_BOTTOM > 178,5 M dan < 300 M: KADIVTEKNIK, MANAGERTEKNIK,
		// KADIVFACULTATIVE (seri); langkah 18 membuang dua yang terakhir.
		{"SENIORUW 300 M", fireTG1, "SENIORUW",
			Transisi{JabatanTujuan: "KADIVTEKNIK", Antrean: "ReasFacInGroupLeader", PositionNoteDitulis: true}},
		{"KADIVTEKNIK 400 M", tsi400, "KADIVTEKNIK",
			Transisi{JabatanTujuan: "DIREKTURMARKETING", Antrean: "ReasFacInMarketingDirector", PositionNoteDitulis: true}},
		{"DIREKTURMARKETING 400 M", tsi400, "DIREKTURMARKETING",
			Transisi{JabatanTujuan: "DIREKTURTEKNIK", Antrean: "ReasFacInTechnicalDirector", PositionNoteDitulis: true}},
	} {
		got, err := Next(u.kasus, tabel, penggunaBiasa(u.pengguna))
		if err != nil || got != u.mau {
			t.Errorf("%s: dapat %+v (%v), mau %+v", u.nama, got, err, u.mau)
		}
	}
}

// TestNextTanggaSelesai - tidak ada jabatan tujuan yang cocok = penyelesaian
// normal, bukan galat.
func TestNextTanggaSelesai(t *testing.T) {
	tabel := muatLimit(t)
	for _, u := range []struct {
		nama     string
		kasus    kasusTangga
		pengguna Pengguna
		mau      Transisi
	}{
		{"DIREKTURTEKNIK, tidak ada di atasnya", fireTG1.dengan(jTSI, "400000000000"), penggunaBiasa("DIREKTURTEKNIK"),
			Transisi{Selesai: true}},
		{"SENIORUW, TSI di bawah limit berikutnya", fireTG1.dengan(jTSI, "200000000000"), penggunaBiasa("SENIORUW"),
			Transisi{Selesai: true}},
		// Kandidat pertama SENIORUW; Decision23 tidak punya konektor SENIORUW → Else.
		{"UNDERWRITER → SENIORUW tanpa konektor", fireTG1.dengan(jTSI, "200000000000"), penggunaBiasa("UNDERWRITER"),
			Transisi{Selesai: true, JabatanTujuan: "SENIORUW"}},
		// Langkah 3: IsGroup → Exit Activity; LetterNo tetap "" dari langkah 2.
		{"anggota grup", fireTG1, Pengguna{Jabatan: "SENIORUW", AnggotaGrup: true},
			Transisi{Selesai: true}},
		// Subquery LOGIN kosong di sistem lama → tidak ada kandidat (butir 44 c).
		{"jabatan pengguna tidak ada di tabel", fireTG1, penggunaBiasa("TIDAKADA"),
			Transisi{Selesai: true}},
		{"team group di luar tabel", fireTG1.dengan(jTeamGroup, "5"), penggunaBiasa("SENIORUW"),
			Transisi{Selesai: true}},
	} {
		got, err := Next(u.kasus, tabel, u.pengguna)
		if err != nil || got != u.mau {
			t.Errorf("%s: dapat %+v (%v), mau %+v", u.nama, got, err, u.mau)
		}
	}
}

// TestNextBasisTSI - langkah 2, 4, 5: TSI yang dibandingkan dengan limit.
func TestNextBasisTSI(t *testing.T) {
	tabel := muatLimit(t)
	keKadivTeknik := Transisi{JabatanTujuan: "KADIVTEKNIK", Antrean: "ReasFacInGroupLeader", PositionNoteDitulis: true}
	for _, u := range []struct {
		nama  string
		kasus kasusTangga
		mau   Transisi
	}{
		// Langkah 4: top risk menggantikan TSI.
		{"IsAdaTopRisk", fireTG1.dengan(jTSI, "100000000000", jAdaTop, "true", jTopRisk, "300000000000"), keKadivTeknik},
		{"TotalTSITopRisk > 0", fireTG1.dengan(jTSI, "100000000000", jTopRisk, "300000000000"), keKadivTeknik},
		// Langkah 5: StatusBusiness 3 (Endorsement) memakai selisih |500 − 200| = 300 M.
		{"endorsement selisih TSI", fireTG1.dengan(jTSI, "500000000000", jStatus, "3", jTSILama, "200000000000", jTopLama, "0"), keKadivTeknik},
		// Selisih top risk > 0 menang atas selisih TSI: |300 − 0| = 300 M.
		{"endorsement selisih top risk", fireTG1.dengan(jTSI, "500000000000", jStatus, "3", jTSILama, "490000000000",
			jTopRisk, "300000000000", jTopLama, "0"), keKadivTeknik},
		// Pembulatan @Math.divide(x,1,0): 255.000.000.001,4 → 255.000.000.001, bukan > limit → selesai.
		{"pembulatan ke bawah", fireTG1.dengan(jTSI, "255000000001.4"), Transisi{Selesai: true}},
	} {
		got, err := Next(u.kasus, tabel, penggunaBiasa("SENIORUW"))
		if err != nil || got != u.mau {
			t.Errorf("%s: dapat %+v (%v), mau %+v", u.nama, got, err, u.mau)
		}
	}
}

// TestNextPilihTabel - langkah 6–10 memilih satu tabel per lini.
func TestNextPilihTabel(t *testing.T) {
	tabel := muatLimit(t)
	// M_LIMIT_ENGINEERINGG tg 1: SENIORUW 85.000.000.001, lalu MANAGERTEKNIK /
	// KADIVFACULTATIVE / KADIVTEKNIK 144.500.000.001, DIREKTURMARKETING 238.000.000.001.
	eng := fireTG1.dengan(jBisnis, "Car", jEngineer, "Car", jPreferred, "", jTSI, "200000000000")
	got, err := Next(eng.dengan(jBisnis, ""), tabel, penggunaBiasa("SENIORUW"))
	if err != nil || got.JabatanTujuan != "KADIVTEKNIK" {
		t.Errorf("engineering: dapat %+v (%v)", got, err)
	}
	// M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL tg 3: SENIORUW 1, DEPHEADUNDERWRITER 255.000.000.001.
	tg3 := fireTG1.dengan(jTeamGroup, "3", jPreferred, "Preferred Risk Commercial", jTSI, "260000000000")
	got, err = Next(tg3, tabel, penggunaBiasa("SENIORUW"))
	mau := Transisi{JabatanTujuan: "DEPHEADUNDERWRITER", Antrean: "ReasFacInDepHeadUnderwriting", PositionNoteDitulis: true}
	if err != nil || got != mau {
		t.Errorf("preferred commercial: dapat %+v (%v), mau %+v", got, err, mau)
	}
	// ToDepHeadUW mensyaratkan IsB2B != "ASM" atau kosong; "ASM" → Else.
	got, err = Next(tg3.dengan(jB2B, "ASM"), tabel, penggunaBiasa("SENIORUW"))
	if err != nil || got != (Transisi{Selesai: true, JabatanTujuan: "DEPHEADUNDERWRITER"}) {
		t.Errorf("DEPHEAD B2B ASM: dapat %+v (%v)", got, err)
	}
}

// TestNextBanding - langkah 11: banding + reject dari antrean Kadiv Facultative
// memakai varian Banding (saring LIMIT_BOTTOM2, urut LIMIT_BOTTOM) dan langkah 18
// tidak membuang KADIVFACULTATIVE.
func TestNextBanding(t *testing.T) {
	tabel := muatLimit(t)
	banding := fireTG1.dengan(jBanding, "true", jReject, "true", jAntrean, "ReasFacInFacultativeDivHead")
	// LIMIT_BOTTOM2 > 178,5 M dan < 270 M: MANAGERTEKNIK, KADIVFACULTATIVE (KADIVTEKNIK 289 M tidak).
	got, err := Next(banding.dengan(jTSI, "270000000000"), tabel, penggunaBiasa("SENIORUW"))
	if mau := (Transisi{JabatanTujuan: "KADIVFACULTATIVE", Antrean: "ReasFacInFacultativeDivHead"}); err != nil || got != mau {
		t.Errorf("banding 270 M: dapat %+v (%v), mau %+v", got, err, mau)
	}
	// 300 M: KADIVFACULTATIVE dan KADIVTEKNIK seri di LIMIT_BOTTOM 255 M → urutan Oracle tidak pasti.
	if _, err := Next(banding, tabel, penggunaBiasa("SENIORUW")); !errors.Is(err, ErrUrutanTakPasti) {
		t.Errorf("banding 300 M: galat %v, mau ErrUrutanTakPasti", err)
	}
	// Banding tanpa reject → langkah 6–10 biasa.
	got, err = Next(banding.dengan(jReject, "false"), tabel, penggunaBiasa("SENIORUW"))
	if err != nil || got.JabatanTujuan != "KADIVTEKNIK" {
		t.Errorf("banding tanpa reject: dapat %+v (%v)", got, err)
	}
}

// TestNextDitolak - keadaan yang perilaku Pega-nya tidak pasti, atau yang belum
// diport, ditolak dengan galat bernama (keputusan work owner 01-10-2026, butir 44 d).
func TestNextDitolak(t *testing.T) {
	tabel := muatLimit(t)
	for _, u := range []struct {
		nama  string
		kasus kasusTangga
		mau   error
	}{
		{"IsFire dan IsEngineering", fireTG1.dengan(jEngineer, "Car"), ErrDaftarGanda},
		{"IsFire tanpa IsPreferredRisk dikenal", fireTG1.dengan(jPreferred, ""), ErrTabelTakTerpilih},
		{"bentuk B (Kredit Cash Loan)", fireTG1.dengan(jBisnisLama, "C2"), ErrBentukB},
		{"TSI kosong", fireTG1.dengan(jTSI, ""), utils.ErrBukanDesimal},
	} {
		if _, err := Next(u.kasus, tabel, penggunaBiasa("SENIORUW")); !errors.Is(err, u.mau) {
			t.Errorf("%s: galat %v, mau %v", u.nama, err, u.mau)
		}
	}
}

// TestNextDataTabelDijaga - kontrak bentuk data tabel limit (tabel sintetis).
func TestNextDataTabelDijaga(t *testing.T) {
	baris := func(j Jabatan, lb string) BarisLimit {
		return BarisLimit{Jabatan: j, TeamGroup: "1", LimitBottom: desimal(t, lb), LimitBottom2: desimal(t, lb)}
	}
	// Dua baris jabatan pengguna dengan limit berbeda: subquery Pega akan gagal.
	ganda := TabelLimit{TabelProperty: {baris("SENIORUW", "1"), baris("SENIORUW", "2"), baris("KADIVTEKNIK", "10")}}
	if _, err := Next(fireTG1, ganda, penggunaBiasa("SENIORUW")); !errors.Is(err, ErrBarisPenggunaGanda) {
		t.Errorf("baris pengguna ganda: galat %v", err)
	}
	// Dua baris identik (SENIORUW tg 1 di data nyata) diterima.
	kembar := TabelLimit{TabelProperty: {baris("SENIORUW", "1"), baris("SENIORUW", "1"), baris("KADIVTEKNIK", "10")}}
	if got, err := Next(fireTG1, kembar, penggunaBiasa("SENIORUW")); err != nil || got.JabatanTujuan != "KADIVTEKNIK" {
		t.Errorf("baris kembar: dapat %+v (%v)", got, err)
	}
	// Langkah 17 dibuang di sistem baru (butir 47): baris "DIREKTUR TEKNIK" berspasi tetap
	// kandidat, seperti di Pega saat gerbang nomor polis salah; tanpa konektor → Else.
	l17 := TabelLimit{TabelProperty: {baris("SENIORUW", "1"), baris("DIREKTUR TEKNIK", "10")}}
	if got, err := Next(fireTG1, l17, penggunaBiasa("SENIORUW")); err != nil || got != (Transisi{Selesai: true, JabatanTujuan: "DIREKTUR TEKNIK"}) {
		t.Errorf("langkah 17: dapat %+v (%v)", got, err)
	}
}

// TestNextCelahMutasi - kasus yang ditambahkan setelah uji mutasi menemukan
// aturan yang belum teruji.
func TestNextCelahMutasi(t *testing.T) {
	tabel := muatLimit(t)
	pengguna := penggunaBiasa("SENIORUW")
	bandingFac := fireTG1.dengan(jBanding, "true", jReject, "true", jAntrean, "ReasFacInFacultativeDivHead", jTSI, "270000000000")

	// Langkah 11 mensyaratkan antrean Kadiv Facultative. Dari antrean lain jalur
	// biasa (saring LIMIT_BOTTOM): KADIVTEKNIK dan KADIVFACULTATIVE (tidak dibuang,
	// banding + reject) seri di 255 M.
	if _, err := Next(bandingFac.dengan(jAntrean, "ReasFacInSeniorUnderwriting"), tabel, pengguna); !errors.Is(err, ErrUrutanTakPasti) {
		t.Errorf("banding dari antrean lain: galat %v, mau ErrUrutanTakPasti", err)
	}
	// Langkah 11.1: Preferred Risk Commercial memakai tabel PROPERTYY varian Banding.
	if got, err := Next(bandingFac.dengan(jPreferred, "Preferred Risk Commercial"), tabel, pengguna); err != nil || got.JabatanTujuan != "KADIVFACULTATIVE" {
		t.Errorf("banding Preferred Risk Commercial: dapat %+v (%v)", got, err)
	}
	// Langkah 5: TSI turun |200 − 500| = 300 M - nilai mutlak.
	turun := fireTG1.dengan(jTSI, "200000000000", jStatus, "3", jTSILama, "500000000000", jTopLama, "0")
	if got, err := Next(turun, tabel, pengguna); err != nil || got.JabatanTujuan != "KADIVTEKNIK" {
		t.Errorf("endorsement TSI turun: dapat %+v (%v)", got, err)
	}

	// Varian Banding menyaring LIMIT_BOTTOM2 tetapi mengurutkan LIMIT_BOTTOM:
	// A (LB 10, LB2 20) mendahului B (LB 15, LB2 12) walau LB2-nya lebih besar.
	d := func(s string) *apd.Decimal { return desimal(t, s) }
	sintetis := TabelLimit{TabelProperty: {
		{Jabatan: "SENIORUW", TeamGroup: "1", LimitBottom: d("5"), LimitBottom2: d("5")},
		{Jabatan: "B", TeamGroup: "1", LimitBottom: d("15"), LimitBottom2: d("12")},
		{Jabatan: "A", TeamGroup: "1", LimitBottom: d("10"), LimitBottom2: d("20")},
	}}
	if got, err := Next(bandingFac.dengan(jTSI, "100"), sintetis, pengguna); err != nil || got.JabatanTujuan != "A" {
		t.Errorf("urutan banding: dapat %+v (%v), mau A", got, err)
	}
}

// TestBulatNolSeriKeAtas - `@Math.divide(x,1,0)` tepat di tengah dibulatkan ke atas
// (A37). Diuji langsung: batas tabel limit semuanya berakhiran …001 (ganjil),
// sehingga lewat Next setengah-ke-atas dan setengah-genap tak terbedakan.
func TestBulatNolSeriKeAtas(t *testing.T) {
	for tsi, mau := range map[string]string{"2.5": "3", "3.5": "4", "2.4": "2"} {
		got, err := bulatNol(fireTG1.dengan(jTSI, tsi), jTSI)
		if err != nil || utils.FormatDecimal(got) != mau {
			t.Errorf("%s: dapat %v (%v), mau %s", tsi, got, err, mau)
		}
	}
}
