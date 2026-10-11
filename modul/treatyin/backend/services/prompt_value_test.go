package services_test

import (
	"sort"
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/services"
)

// ⭐ Label yang TERBUKTI — dari berkas rule atau tangkapan layarnya.
func TestLabelDariEkspor(t *testing.T) {
	for _, u := range []struct{ properti, nilai, mau string }{
		{"OptionLimit", "1", "Of Cession to R/I"},
		{"OptionLimit", "2", "Of 100% Limit"},
		{"EDMState", "1", "Internal"},
		{"EDMState", "2", "External"},
		{"EDMMaterialType", "1", "Material"},
		{"EDMMaterialType", "2", "Non Material"},
		// ⭐ Ketiganya lewat penerjemah yang SUDAH ada, bukan salinan di
		// peta — dua sumber untuk satu label adalah cara termudah keduanya
		// berbeda diam-diam.
		{"ReportingPeriod", "quarter", "Quarter Year"},
		{"ReportingPeriod", "half", "Half Year"},
		{"AccountingMode", "underwriting", "Underwriting Year"},
		{"AccountingModeNonProp", "loss", "Loss Occuring"},
		{"Bordeaux", "reporting", "Reporting"},
	} {
		if got := services.LabelPrompt(u.properti, u.nilai); got != u.mau {
			t.Errorf("%s.%s = %q, mau %q", u.properti, u.nilai, got, u.mau)
		}
	}
}

// ⭐ LABEL YANG DINYATAKAN PEMILIK PROSES 7 Oktober 2026:
// *"riskcat itu seharusnya yang diambil risk & cat"*.
//
// ⛔ Satu contoh itu mengunci `Cover`, BUKAN yang lain.
func TestLabelCoverDariPemilikProses(t *testing.T) {
	for _, u := range []struct{ nilai, mau string }{
		{"riskcat", "risk & cat"},
		{"risk", "risk"},
		{"cat", "cat"},
	} {
		if got := services.LabelPrompt("Cover", u.nilai); got != u.mau {
			t.Errorf("Cover.%s = %q, mau %q", u.nilai, got, u.mau)
		}
	}
	if asal := services.AsalLabel("Cover", "riskcat"); asal != services.AsalPemilik {
		t.Fatalf("asal Cover.riskcat = %q, mau %q", asal, services.AsalPemilik)
	}
}

// ⛔ LABEL YANG DISIMPULKAN — dan tiap satunya WAJIB mengaku.
//
// Perintah pemilik proses menyuruh seluruh dropdown diperbaiki, sementara
// contohnya hanya satu. Sisanya disimpulkan, dan kesimpulan yang tidak
// menyatakan dirinya kesimpulan akan dibaca sebagai fakta oleh ronde
// berikutnya.
func TestLabelYangDisimpulkanMengakuiDirinya(t *testing.T) {
	for _, u := range []struct{ properti, nilai, mau string }{
		{"LayerType", "sublayer", "sub layer"},
		{"AccountingModeNonProp", "risk", "Risk Attaching"},
		{"Bordeaux", "nonreporting", "Non Reporting"},
	} {
		if got := services.LabelPrompt(u.properti, u.nilai); got != u.mau {
			t.Errorf("%s.%s = %q, mau %q", u.properti, u.nilai, got, u.mau)
		}
		if asal := services.AsalLabel(u.properti, u.nilai); asal != services.AsalKesimpulan {
			t.Errorf("%s.%s asal %q — kesimpulan HARUS mengaku", u.properti, u.nilai, asal)
		}
	}
}

// ⛔ YANG TERBUKTI SELALU MENANG atas yang disimpulkan.
func TestYangTerbuktiMenangAtasKesimpulan(t *testing.T) {
	for _, u := range []struct{ properti, nilai, mau string }{
		{"Bordeaux", "reporting", "Reporting"},
		{"AccountingModeNonProp", "loss", "Loss Occuring"},
		{"OptionLimit", "2", "Of 100% Limit"},
		// ⛔ `ReinstatementNote.xml` (8 Oktober 2026) — dan kedatangannya
		// MEMBATALKAN kesimpulan yang sebelumnya ada di sini. Dugaannya
		// `as amount` / `as time`, dasarnya "kode gandeng dipecah jadi
		// kata"; yang sebenarnya kalimat penuh yang nol hubungannya dengan
		// ejaan kodenya.
		{"ReinstatementNote", "asamount", "Additional Premium as to amount"},
		{"ReinstatementNote", "astime", "Additional Premium as to time"},
	} {
		if got := services.LabelPrompt(u.properti, u.nilai); got != u.mau {
			t.Errorf("%s.%s = %q, mau %q", u.properti, u.nilai, got, u.mau)
		}
		if asal := services.AsalLabel(u.properti, u.nilai); asal != services.AsalEkspor {
			t.Errorf("%s.%s asal %q, mau %q", u.properti, u.nilai, asal, services.AsalEkspor)
		}
	}
}

// ⭐ Properti tak dikenal jatuh ke nilainya sendiri, dan nol asal.
func TestPropertiTakDikenalJatuhKeNilainya(t *testing.T) {
	if got := services.LabelPrompt("EntahApa", "xyz"); got != "xyz" {
		t.Fatalf("= %q, mau xyz", got)
	}
	if asal := services.AsalLabel("EntahApa", "xyz"); asal != "" {
		t.Fatalf("asal = %q, mau kosong", asal)
	}
}

// ⭐ CERMIN: berapa label yang masih DISIMPULKAN, bukan terbukti.
//
// ⛔ Angkanya DIPAKU. Ralat dari pemilik proses memindahkan satu baris dari
// `promptDisimpulkan` ke `promptValue`, dan uji ini merah sampai angkanya
// diperbarui — jadi kesenjangan ini menyusut TERUKUR.
//
// ⚠️ Yang dihitung BUKAN "belum punya label": seluruh dropdown kini punya.
// Yang dihitung "label yang belum terbukti".
func TestBerapaLabelMasihDisimpulkan(t *testing.T) {
	var kesimpulan []string
	for _, d := range services.DomainProperty {
		for _, v := range d.Nilai {
			if services.AsalLabel(d.Properti, v) == services.AsalKesimpulan {
				kesimpulan = append(kesimpulan, d.Properti+"."+v)
			}
		}
	}
	sort.Strings(kesimpulan)
	t.Logf("label yang masih DISIMPULKAN (%d):\n  %s",
		len(kesimpulan), strings.Join(kesimpulan, "\n  "))

	// Terukur 8 Oktober 2026 — TIGA:
	//   LayerType.sublayer · AccountingModeNonProp.risk ·
	//   Bordeaux.nonreporting
	//
	// ⭐ Turun dari LIMA sejak `ReinstatementNote.xml` tiba, dan cara
	// turunnya layak dicatat: kedua label itu TIDAK terbukti benar, mereka
	// terbukti SALAH. Dua dari dua kesimpulan yang akhirnya diadu dengan
	// ekspornya meleset (`OptionLimit` lebih dahulu), jadi ketiga sisa di
	// bawah ini lebih baik dianggap salah sampai berkasnya datang.
	//
	// ⚠️ `LayerType.layer` dan `CurrencyRelation.*` TIDAK terhitung:
	// labelnya sama dengan nilainya, jadi nol yang disimpulkan di sana.
	if len(kesimpulan) != 3 {
		t.Fatalf("masih disimpulkan = %d, terakhir diukur 3:\n  %s",
			len(kesimpulan), strings.Join(kesimpulan, "\n  "))
	}
}

// ⛔ Tiap domain WAJIB punya nama rule Property-nya. Domain tanpa nama
// tidak dapat dimintakan ekspornya, dan akan diam selamanya.
func TestSetiapDomainPunyaNamaProperty(t *testing.T) {
	for _, d := range services.DomainProperty {
		if strings.TrimSpace(d.Properti) == "" {
			t.Fatalf("domain tanpa nama properti: %v", d.Nilai)
		}
		if len(d.Nilai) == 0 {
			t.Fatalf("properti %q nol nilai", d.Properti)
		}
	}
}
