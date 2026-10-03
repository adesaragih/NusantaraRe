package rules

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// kasusUji - Kasus dari peta jalur → nilai. Jalur ditulis PERSIS seperti di
// korpus (`pyWorkPage.Quotation.BusinessCode`, `.Quotation.BusinessType`).
type kasusUji map[string]string

func (k kasusUji) Nilai(jalur string) (string, bool) {
	v, ada := k[jalur]
	return v, ada
}

// ujiEval - satu pertanyaan ke Eval dengan jawaban yang diharapkan.
type ujiEval struct {
	nama, predikat string
	kasus          kasusUji
	mau            bool
}

func periksaEval(t *testing.T, daftar []ujiEval) {
	t.Helper()
	for _, u := range daftar {
		got, err := Eval(u.predikat, u.kasus)
		if err != nil {
			t.Errorf("%s: %v", u.nama, err)
			continue
		}
		if got != u.mau {
			t.Errorf("%s: %s = %v, mau %v", u.nama, u.predikat, got, u.mau)
		}
	}
}

// TestEvalPredikatNyata - predikat dari registry bangkitan, dua kelompok tag
// (NB-08): kondisi yang tampak di `pyConditionString`, dan yang tampilannya
// placeholder `[Double click to add condition]` sehingga hanya terbaca dari
// ekspresi tersimpan.
func TestEvalPredikatNyata(t *testing.T) {
	periksaEval(t, []ujiEval{
		// IsAsuransiKredit - tampilan terbaca: A OR B OR C OR D atas
		// pyWorkPage.Quotation.BusinessCode = "10053"/"10243"/"10244"/"10245".
		{"tampilan terbaca, cocok baris B", "IsAsuransiKredit", kasusUji{"pyWorkPage.Quotation.BusinessCode": "10243"}, true},
		{"tampilan terbaca, tidak cocok", "IsAsuransiKredit", kasusUji{"pyWorkPage.Quotation.BusinessCode": "99999"}, false},
		{"tampilan terbaca, properti tidak ada", "IsAsuransiKredit", kasusUji{}, false},
		// IsBonding - tampilan PLACEHOLDER: A OR B OR C, B = evaluateWhen("IsBondingKBG").
		{"placeholder, baris A relatif", "IsBonding", kasusUji{".Quotation.BusinessType": "Bonding"}, true},
		{"placeholder, lewat rujukan IsBondingKBG", "IsBonding", kasusUji{"pyWorkPage.Quotation.BusinessType": "BondingKBG"}, true},
		{"placeholder, tidak cocok", "IsBonding", kasusUji{}, false},
		// IsCivilEngineeringCompletedRisks - A AND (B OR C).
		{"kurung: A dan C", "IsCivilEngineeringCompletedRisks", kasusUji{"pyWorkPage.Quotation.BusinessType": "Aneka", "pyWorkPage.Quotation.BusinessName": "CIVIL ENGINEERING COMPLETED RISKS"}, true},
		{"kurung: A saja", "IsCivilEngineeringCompletedRisks", kasusUji{"pyWorkPage.Quotation.BusinessType": "Aneka"}, false},
		{"kurung: B tanpa A", "IsCivilEngineeringCompletedRisks", kasusUji{"pyWorkPage.Quotation.BusinessCode": "10166"}, false},
		// IsChekCity - A AND B, keduanya `!= ""`.
		{"tidak sama: keduanya terisi", "IsChekCity", kasusUji{"InputCity.ID": "x", "InputCity.Note": "y"}, true},
		{"tidak sama: satu kosong", "IsChekCity", kasusUji{"InputCity.ID": "", "InputCity.Note": "y"}, false},
		// IsAdjustableFlag - .IsAdjustableFlag = true.
		{"literal true", "IsAdjustableFlag", kasusUji{".IsAdjustableFlag": "true"}, true},
		{"literal true, nilai false", "IsAdjustableFlag", kasusUji{".IsAdjustableFlag": "false"}, false},
		// IsFlagOldData - .FlagOldData = 1 (angka); identitas di dua berkas identik.
		{"angka sama", "IsFlagOldData", kasusUji{".FlagOldData": "1"}, true},
		{"angka beda", "IsFlagOldData", kasusUji{".FlagOldData": "2"}, false},
		// Nama predikat tidak peka huruf besar-kecil, seperti rujukan di korpus.
		{"nama huruf kecil", "isbonding", kasusUji{".Quotation.BusinessType": "Bonding"}, true},
	})
}

// TestEvalNegasiDanSimbol - `!` = NOT dan `&&` = AND [dugaan], butir 23.
func TestEvalNegasiDanSimbol(t *testing.T) {
	periksaEval(t, []ujiEval{
		// IsNotPAandNotMBU = !A AND !B, A = IsPA, B = IsMBU.
		{"bukan PA bukan MBU", "IsNotPAandNotMBU", kasusUji{"pyWorkPage.Quotation.BusinessType": "Fire"}, true},
		{"PA", "IsNotPAandNotMBU", kasusUji{"pyWorkPage.Quotation.BusinessType": "PA"}, false},
		{"MBU", "IsNotPAandNotMBU", kasusUji{"pyWorkPage.OfferFacIn.QuotationData.BusinessType": "MBUCar"}, false},
		// isTravelTime = A && B.
		{"travel berjalan", "isTravelTime", kasusUji{"pyWorkPage.Quotation.BusinessType": "Travel", "pyWorkPage.FlagOnGoingPolicy": "1"}, true},
		{"travel tidak berjalan", "isTravelTime", kasusUji{"pyWorkPage.Quotation.BusinessType": "Travel", "pyWorkPage.FlagOnGoingPolicy": "0"}, false},
	})
}

// TestEvalTafsirBerbedaDitolak - tipe properti tidak ada di korpus. Bila tafsir
// teks dan tafsir angka memberi hasil berbeda, atau nilai kosong dibandingkan
// dengan angka, Eval mengembalikan galat, bukan tebakan (butir 20).
func TestEvalTafsirBerbedaDitolak(t *testing.T) {
	for _, u := range []struct {
		nama, predikat string
		kasus          kasusUji
	}{
		{"1.0 vs angka 1", "IsFlagOldData", kasusUji{".FlagOldData": "1.0"}},
		{"kosong vs angka 1", "IsFlagOldData", kasusUji{".FlagOldData": ""}},
		{"tidak ada vs angka 1", "IsFlagOldData", kasusUji{}},
		{"10053.0 vs teks \"10053\"", "IsAsuransiKredit", kasusUji{"pyWorkPage.Quotation.BusinessCode": "10053.0"}},
	} {
		if _, err := Eval(u.predikat, u.kasus); !errors.Is(err, ErrTafsirBerbeda) {
			t.Errorf("%s: galat %v, mau ErrTafsirBerbeda", u.nama, err)
		}
	}
}

// TestEvalGagalKeras - kesalahan program → panic: nama tak dikenal (salah ketik
// tidak boleh menjadi gerbang yang selalu tertutup), IsPKSASM (spec Modul 3),
// bentuk kondisi yang belum diport, boolean berhuruf besar, dan login operator
// literal (nilainya tidak disalin, CLAUDE.md §4 butir 10).
func TestEvalGagalKeras(t *testing.T) {
	// alasan - potongan pesan panic: panic yang datang dari jalur lain (mis.
	// baris kondisi yang hilang) tidak boleh terbaca sebagai bukti.
	for _, u := range []struct{ nama, alasan string }{
		{"IsTidakAdaDiKorpus", "tidak ada di registry"},
		{"IsPKSASM", "kondisi belum terbukti dieksekusi"},
		{"pyIsMobile", "bentuk ekspresi lain"},                // pzIsMobile(tools)
		{"IsPASSG", "sisi kiri bukan jalur properti"},         // SizeOfPropertyList(...) > 0
		{"pyIsIpadOrDesktop", "merujuk predikat PYISIPAD"},    // pyIsIPad tidak ada di korpus
		{"IsVisible", "operand kanan boolean berhuruf besar"}, // = True
		{"IsGroup", "login operator"},                         // OperatorID.pyUserIdentifier
		{"IsGroupCreate", "login operator"},                   // pxCreateOperator
		{"IsSPVCreate", "login operator"},
	} {
		func() {
			defer func() {
				r := recover()
				if r == nil {
					t.Errorf("%s: tidak panic", u.nama)
				} else if s := fmt.Sprint(r); !strings.Contains(s, u.alasan) {
					t.Errorf("%s: panic %q, mau memuat %q", u.nama, s, u.alasan)
				}
			}()
			_, _ = Eval(u.nama, kasusUji{})
		}()
	}
}

// TestEvalKataTelanjangSebagaiTeks - operand kanan tanpa kutip dibaca sebagai teks
// (keputusan work owner 01-10-2026, butir 28). Operator tanpa workbasket tidak
// membuka IsUW: "" tidak sama dengan "ReasFacInDirector".
func TestEvalKataTelanjangSebagaiTeks(t *testing.T) {
	const wb = "pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName"
	for _, u := range []struct {
		nama  string
		kasus kasusUji
		mau   bool
	}{
		{"IsUW", kasusUji{wb: "ReasFacInDirector"}, true},
		{"IsUW", kasusUji{wb: "ReasFacInUnderwriting"}, true},
		{"IsUW", kasusUji{wb: ""}, false},
		{"IsUW", kasusUji{wb: "reasfacindirector"}, false}, // teks: peka huruf
		{"IsCedingConfirmOffer", kasusUji{".IsCedingConfirm": "accepted"}, false},
		{"IsCedingConfirmOffer", kasusUji{".IsCedingConfirm": "Offer"}, true},
	} {
		got, err := Eval(u.nama, u.kasus)
		if err != nil || got != u.mau {
			t.Errorf("%s %v: dapat %v (%v), mau %v", u.nama, u.kasus, got, err, u.mau)
		}
	}
}
