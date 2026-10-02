package main

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// TestPecahArgumen - koma di dalam kutip dan kurung bukan pemisah.
func TestPecahArgumen(t *testing.T) {
	for _, u := range []struct {
		masuk string
		mau   []string
	}{
		{`pyWorkPage.Quotation.BusinessCode, "=", "10053"`, []string{"pyWorkPage.Quotation.BusinessCode", `"="`, `"10053"`}},
		{`.X, "=", "a, b"`, []string{".X", `"="`, `"a, b"`}},
		{`@Utilities.SizeOfPropertyList(pyWorkPage.A, B), ">", 0`, []string{"@Utilities.SizeOfPropertyList(pyWorkPage.A, B)", `">"`, "0"}},
	} {
		if got := pecahArgumen(u.masuk); !reflect.DeepEqual(got, u.mau) {
			t.Errorf("%q → %q, mau %q", u.masuk, got, u.mau)
		}
	}
}

// TestKenaliKondisi - bentuk yang diport menjadi literal; selain itu alasan.
// Contoh diambil dari baris korpus `NB FacIn\When`.
func TestKenaliKondisi(t *testing.T) {
	const awal = "@(Pega-RULES:ExpressionEvaluators)."
	for _, u := range []struct {
		masuk, literal, alasan string
	}{
		{awal + `compareTwoValues(pyWorkPage.Quotation.BusinessCode, "=", "10053")`,
			`{jenis: banding, kiri: "pyWorkPage.Quotation.BusinessCode", tidakSama: false, kanan: operand{teksLiteral, "10053"}}`, ""},
		{awal + `compareTwoValues(.FlagOldData, "=", 1)`,
			`{jenis: banding, kiri: ".FlagOldData", tidakSama: false, kanan: operand{angkaLiteral, "1"}}`, ""},
		{awal + `compareTwoValues(InputCity.ID, "!=", "")`,
			`{jenis: banding, kiri: "InputCity.ID", tidakSama: true, kanan: operand{teksLiteral, ""}}`, ""},
		{awal + `compareTwoValues(.IsAdjustableFlag, "=", true)`,
			`{jenis: banding, kiri: ".IsAdjustableFlag", tidakSama: false, kanan: operand{booleanLiteral, "true"}}`, ""},
		{awal + `evaluateWhen("isGrowingTrees")`, `{jenis: rujukWhen, rujukan: "ISGROWINGTREES"}`, ""},
		// Kata telanjang dibaca sebagai teks (keputusan work owner 01-10-2026, butir 28).
		{awal + `compareTwoValues(pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName, "=", ReasFacInDirector)`,
			`{jenis: banding, kiri: "pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName", tidakSama: false, kanan: operand{teksLiteral, "ReasFacInDirector"}}`, ""},
		{awal + `compareTwoValues(pyWorkPage.OfferFacIn.QuotationData.BusinessOldId, "=", C2)`,
			`{jenis: banding, kiri: "pyWorkPage.OfferFacIn.QuotationData.BusinessOldId", tidakSama: false, kanan: operand{teksLiteral, "C2"}}`, ""},
		// Boolean berhuruf besar bukan kata biasa: tetap ditolak.
		{awal + `compareTwoValues(.IsVisible, "=", True)`, "", "operand kanan boolean berhuruf besar"},
		{awal + `compareTwoValues(.X, "=", .Y)`, "", "operand kanan bukan literal"},
		// Login operator tidak disalin (CLAUDE.md §4 butir 10); nilai uji rekaan.
		{awal + `compareTwoValues(OperatorID.pyUserIdentifier, "=", "UJI-LOGIN")`, "", alasanLogin},
		{awal + `compareTwoValues(pyWorkPage.pxCreateOperator, "!=", "UJI-LOGIN")`, "", alasanLogin},
		// Nomor polis, nama marketing, telepon operator: literal tidak disalin (CLAUDE.md
		// §4 butir 10, §10); nilai uji rekaan. Literal KOSONG tetap boleh (cek isian).
		{awal + `compareTwoValues(pyWorkPage.Quotation.OldPolicyNo, "=", "UJI-POLIS-1")`, "", alasanIdentitas},
		{awal + `compareTwoValues(.OfferFacIn.QuotationData.MarketingName, "!=", "UJI NAMA")`, "", alasanIdentitas},
		{awal + `compareTwoValues(OperatorID.pyTelephone, "=", "UJI-TELP")`, "", alasanIdentitas},
		{awal + `compareTwoValues(pyWorkPage.PolicyTreatyIn.PolicyNo, "=", 12345)`, "", alasanIdentitas},
		{awal + `compareTwoValues(pyWorkPage.PolicyTreatyIn.PolicyNo, "=", "")`,
			`{jenis: banding, kiri: "pyWorkPage.PolicyTreatyIn.PolicyNo", tidakSama: false, kanan: operand{teksLiteral, ""}}`, ""},
		{awal + `compareTwoValues(@Utilities.SizeOfPropertyList(pyWorkPage.PersonListPA), ">", 0)`, "", "sisi kiri bukan jalur properti"},
		{awal + `compareTwoValues(.X, ">", 0)`, "", "operator"},
		{"pyPortal.IsEdit", "", "bentuk ekspresi lain"},
	} {
		kd, alasan := kenaliKondisi(u.masuk)
		if u.alasan != "" {
			if !strings.HasPrefix(alasan, u.alasan) {
				t.Errorf("%s: alasan %q, mau berawal %q", u.masuk, alasan, u.alasan)
			}
			continue
		}
		if alasan != "" || kd.literal != u.literal {
			t.Errorf("%s:\n  dapat %q (%s)\n  mau   %q", u.masuk, kd.literal, alasan, u.literal)
		}
	}
}

// TestLoginTidakDisalin - predikat yang membandingkan login operator ditolak,
// dan nilai login-nya tidak muncul di pesan panik (CLAUDE.md §4 butir 10).
func TestLoginTidakDisalin(t *testing.T) {
	const xmlWhen = `<pagedata><pxInsName>ASM-FW-GISFW-WORK!ISREKAAN</pxInsName><pyLogic>A OR B</pyLogic>
<pyCondition><rowdata><pyConditionLabel>A</pyConditionLabel><pyConditionValue1>@(Pega-RULES:ExpressionEvaluators).compareTwoValues(.Status, "=", "1")</pyConditionValue1></rowdata>
<rowdata><pyConditionLabel>B</pyConditionLabel><pyConditionValue1>@(Pega-RULES:ExpressionEvaluators).compareTwoValues(OperatorID.pyUserIdentifier, "=", "UJI-LOGIN")</pyConditionValue1></rowdata></pyCondition></pagedata>`
	jalur := filepath.Join(t.TempDir(), "IsRekaan.xml")
	if err := os.WriteFile(jalur, []byte(xmlWhen), 0o600); err != nil {
		t.Fatal(err)
	}
	e, err := baca(jalur, "uji", "", sikapKhusus)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(e.panik, alasanLogin) || strings.Contains(e.panik, "UJI-LOGIN") {
		t.Fatalf("panik %q: mau berisi alasan login, tanpa nilainya", e.panik)
	}
}

// TestRegistrySamaDenganKeluaranBangkit - penjaga `registry_gen.go`: hasil
// bangkit ulang dari korpus harus sama persis dengan yang di-commit. Korpus tidak
// ada di repositori, jadi uji ini butuh env KORPUS_RNM_AKAR (akar korpus, berisi `NB FacIn\When`)
// dan dilewati tanpanya.
func TestRegistrySamaDenganKeluaranBangkit(t *testing.T) {
	akar := os.Getenv("KORPUS_RNM_AKAR")
	if akar == "" {
		t.Skip("KORPUS_RNM_AKAR kosong: penjaga registry_gen.go dilewati (korpus tidak ada di repositori)")
	}
	nb, err := pilihProfil(folderUtama, "rules")
	if err != nil {
		t.Fatal(err)
	}
	baru, _, err := bangkitkan(akar, nb)
	if err != nil {
		t.Fatal(err)
	}
	lama, err := os.ReadFile("../registry_gen.go")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(bytes.ReplaceAll(lama, []byte("\r\n"), []byte("\n")), baru) {
		t.Error("registry_gen.go berbeda dari keluaran bangkit - jalankan ulang pembangkit (lihat main.go)")
	}
}

// TestPilihProfil - folder NB memakai pinjaman K-003 dan catatan/sikap NB; folder
// lain (mis. siklus EDM, opsi (b) registry predikat EDM) TIDAK mewarisi satu pun:
// keputusan itu milik predikat NB (CLAUDE.md §4.6).
func TestPilihProfil(t *testing.T) {
	nb, err := pilihProfil(folderUtama, "rules")
	if err != nil || len(nb.pinjaman) != len(pinjaman) || len(nb.catatan) != len(catatanSikap) || len(nb.sikap) != len(sikapKhusus) {
		t.Fatalf("profil NB %+v (%v)", nb, err)
	}
	edm, err := pilihProfil(`Endorsment Fac In\When`, "edm")
	if err != nil || edm.folder != `Endorsment Fac In\When` || edm.paket != "edm" ||
		len(edm.pinjaman) != 0 || len(edm.catatan) != 0 || len(edm.sikap) != 0 {
		t.Fatalf("profil folder lain %+v (%v)", edm, err)
	}
	for _, u := range [][2]string{{"", "edm"}, {`Endorsment Fac In\When`, ""}, {`Endorsment Fac In\When`, "Edm-1"}} {
		if _, err := pilihProfil(u[0], u[1]); err == nil {
			t.Errorf("%q/%q diterima", u[0], u[1])
		}
	}
	// NB ke paket selain rules: registry NB hanya untuk paket rules nbfacin.
	if _, err := pilihProfil(folderUtama, "edm"); err == nil {
		t.Error("folder NB ke paket lain diterima")
	}
}

// TestBangkitFolderLain - folder When siklus lain dibangkitkan apa adanya ke paket
// yang diminta, tanpa pinjaman. Butuh KORPUS_RNM_AKAR (berisi `Endorsment Fac In\When`).
func TestBangkitFolderLain(t *testing.T) {
	akar := os.Getenv("KORPUS_RNM_AKAR")
	if akar == "" {
		t.Skip("KORPUS_RNM_AKAR kosong: dilewati (korpus tidak ada di repositori)")
	}
	p, err := pilihProfil(`Endorsment Fac In\When`, "edm")
	if err != nil {
		t.Fatal(err)
	}
	src, ringkas, err := bangkitkan(akar, p)
	if err != nil {
		t.Fatal(err)
	}
	s := string(src)
	for _, mau := range []string{"package edm\n", "// Sumber: Endorsment Fac In\\When - ", " + 0 pinjaman,"} {
		if !strings.Contains(s, mau) {
			t.Errorf("keluaran tanpa %q", mau)
		}
	}
	if strings.Contains(s, "K-002") || strings.Contains(s, "K-019") || strings.Contains(s, "dipinjam") {
		t.Error("keluaran folder lain mewarisi catatan/pinjaman NB")
	}
	t.Log(ringkas)
}

// TestIdentitasTidakDisalinDalamBentukLain - ekspresi yang menyebut medan identitas
// atau login dalam bentuk yang BELUM diport (mis. @String.equals) tetap tidak disalin
// ke pesan panik; pencocokan nama medan tidak peka huruf besar-kecil.
func TestIdentitasTidakDisalinDalamBentukLain(t *testing.T) {
	for nama, ekspresi := range map[string]string{
		"bentuk lain":    `@String.equals(pyWorkPage.Quotation.OldPolicyNo,"UJI-POLIS-1")`,
		"huruf kecil":    `@(Pega-RULES:ExpressionEvaluators).compareTwoValues(pyWorkPage.Quotation.oldpolicyno, "=", "UJI-POLIS-1")`,
		"login lain":     `@String.equals(OperatorID.pyUserIdentifier,"UJI-LOGIN")`,
		"login kecil":    `@(Pega-RULES:ExpressionEvaluators).compareTwoValues(OperatorID.pyuseridentifier, "=", "UJI-LOGIN")`,
		"marketing lain": `@String.contains(.OfferFacIn.QuotationData.MarketingName,"UJI NAMA")`,
	} {
		xmlWhen := `<pagedata><pxInsName>ASM-FW-GISFW-WORK!ISREKAAN</pxInsName><pyLogic>A</pyLogic>
<pyCondition><rowdata><pyConditionLabel>A</pyConditionLabel><pyConditionValue1>` + ekspresi + `</pyConditionValue1></rowdata></pyCondition></pagedata>`
		jalur := filepath.Join(t.TempDir(), "IsRekaan.xml")
		if err := os.WriteFile(jalur, []byte(xmlWhen), 0o600); err != nil {
			t.Fatal(err)
		}
		e, err := baca(jalur, "uji", "", sikapKhusus)
		if err != nil {
			t.Fatal(err)
		}
		if e.panik == "" || strings.Contains(e.panik, "UJI-") || strings.Contains(e.panik, "UJI NAMA") {
			t.Errorf("%s: panik %q", nama, e.panik)
		}
	}
	// Garis miring biasa = folder NB yang sama.
	if p, err := pilihProfil("NB FacIn/When", "rules"); err != nil || len(p.pinjaman) != len(pinjaman) {
		t.Errorf("NB FacIn/When: %+v (%v)", p, err)
	}
}
