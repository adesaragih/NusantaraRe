package rules

import (
	"sort"
	"strings"
	"testing"
)

// TestRegistryUtuh - table-driven atas SELURUH registry bangkitan (spec Modul 3):
// setiap predikat menyebut asal rule-nya, dan setiap predikat yang tidak panic
// punya logika yang terurai (dijamin `pohonLogika` saat paket dimuat), label
// yang semuanya berbaris kondisi, dan rujukan yang semuanya ada di registry.
func TestRegistryUtuh(t *testing.T) {
	// [terverifikasi] 210 berkas `NB FacIn\When`, 209 identitas `<pxInsName>`
	// unik, dihitung dua cara pada 01-10-2026 atas folder itu saja:
	//   cara 1: ls "NB FacIn/When" | wc -l                                  → 210
	//           grep -ho "<pxInsName>[^<]*</pxInsName>" "NB FacIn"/When/*.xml | sort -u | wc -l → 209
	//   cara 2: pengurai XML Python, kunci (class, nama) dari <pxInsName>    → 209, kembar 1
	if len(registry) != 210 {
		t.Errorf("registry %d predikat, mau 210 (209 NB + 1 pinjaman K-003)", len(registry))
	}
	for nama, pr := range registry {
		if nama != strings.ToUpper(nama) {
			t.Errorf("%s: kunci registry wajib huruf besar", nama)
		}
		if len(pr.asal) == 0 || !strings.Contains(pr.asal[0], ".xml (") || !strings.Contains(pr.asal[0], "!"+nama+")") {
			t.Errorf("%s: asal %v tidak menyebut berkas dan identitas", nama, pr.asal)
		}
		if pr.panik != "" {
			continue
		}
		label := pohonLogika[nama].Label()
		if len(label) != len(pr.kondisi) {
			t.Errorf("%s: %d label di logika, %d baris kondisi", nama, len(label), len(pr.kondisi))
		}
		for _, l := range label {
			kd, ada := pr.kondisi[l]
			if !ada {
				t.Errorf("%s: label %s tanpa baris kondisi", nama, l)
				continue
			}
			if kd.jenis == rujukWhen {
				if _, ada := registry[kd.rujukan]; !ada {
					t.Errorf("%s: merujuk %s yang tidak ada", nama, kd.rujukan)
				}
			}
		}
	}
}

// TestRegistryIdentitasKembar - ISFLAGOLDDATA tersimpan di dua berkas. Pembangkit
// menyatukannya hanya karena sidik isinya (logika + baris kondisi) sama; hash
// ternormalisasi seluruh berkas (23 tag volatil + isi pzIndexes dibuang) juga
// sama, diukur terpisah - `docs/KEPUTUSAN-30-09-2026.md` bab NB-08.
func TestRegistryIdentitasKembar(t *testing.T) {
	if asal := registry["ISFLAGOLDDATA"].asal; len(asal) != 2 {
		t.Errorf("asal ISFLAGOLDDATA %v, mau dua berkas", asal)
	}
}

// TestEvalRujukanMelingkarPanic - korpus NB tidak memuat rujukan melingkar, jadi
// dua predikat uji disisipkan sementara ke registry dan dibuang lagi. Tidak
// boleh `t.Parallel`: ia menyunting variabel paket.
func TestEvalRujukanMelingkarPanic(t *testing.T) {
	registry["UJI_LINGKAR_A"] = predikat{asal: []string{"uji"}, logika: "A",
		kondisi: map[string]kondisi{"A": {jenis: rujukWhen, rujukan: "UJI_LINGKAR_B"}}}
	registry["UJI_LINGKAR_B"] = predikat{asal: []string{"uji"}, logika: "A",
		kondisi: map[string]kondisi{"A": {jenis: rujukWhen, rujukan: "UJI_LINGKAR_A"}}}
	defer func() {
		delete(registry, "UJI_LINGKAR_A")
		delete(registry, "UJI_LINGKAR_B")
		if recover() == nil {
			t.Error("rujukan melingkar tidak panic")
		}
	}()
	_, _ = Eval("uji_lingkar_a", kasusUji{})
}

// medanIdentitasUji - huruf kecil; dicocokkan TIDAK peka huruf besar-kecil, sama dengan
// generator (`bangkit.cocokMedan`) - jebakan sensus no. 4.
var medanIdentitasUji = map[string]bool{"oldpolicyno": true, "policyno": true, "marketingname": true, "pytelephone": true,
	"pyuseridentifier": true, "pxcreateoperator": true, "pxupdateoperator": true}

// literalIdentitas - pelanggaran di `reg`: pembanding medan identitas dengan literal
// tidak-kosong, atau pesan panik yang menyebut medan identitas.
func literalIdentitas(reg map[string]predikat) []string {
	var langgar []string
	for nama, p := range reg {
		for m := range medanIdentitasUji {
			if strings.Contains(strings.ToLower(p.panik), "."+m) {
				langgar = append(langgar, nama+": pesan panik menyalin ekspresi atas "+m)
			}
		}
		for label, k := range p.kondisi {
			seg := strings.Split(k.kiri, ".")
			if k.jenis == banding && medanIdentitasUji[strings.ToLower(seg[len(seg)-1])] && k.kanan.teks != "" {
				langgar = append(langgar, nama+" baris "+label+": "+k.kiri+" dibandingkan dengan literal tidak-kosong")
			}
		}
	}
	sort.Strings(langgar)
	return langgar
}

// TestRegistryTanpaLiteralIdentitas - penjaga isi registry_gen.go TANPA korpus: tidak
// satu pembanding pun memuat literal tidak-kosong untuk medan identitas (nomor polis,
// nama marketing, telepon/login operator). CLAUDE.md §4 butir 10 dan §10; pasangan
// generatornya `medanLogin` + `medanIdentitas` di ./bangkit.
func TestRegistryTanpaLiteralIdentitas(t *testing.T) {
	for _, l := range literalIdentitas(registry) {
		t.Error(l)
	}
}

// TestPenjagaIdentitasMenggigit - uji instrumen dengan butir yang jawabannya sudah
// diketahui: huruf besar-kecil apa pun tertangkap, literal kosong dan medan lain lolos.
func TestPenjagaIdentitasMenggigit(t *testing.T) {
	reg := map[string]predikat{
		"UJI_KECIL":  {kondisi: map[string]kondisi{"A": {jenis: banding, kiri: "pyWorkPage.Quotation.oldpolicyno", kanan: operand{teksLiteral, "UJI-POLIS"}}}},
		"UJI_BESAR":  {kondisi: map[string]kondisi{"A": {jenis: banding, kiri: ".MARKETINGNAME", kanan: operand{teksLiteral, "UJI NAMA"}}}},
		"UJI_KOSONG": {kondisi: map[string]kondisi{"A": {jenis: banding, kiri: "pyWorkPage.PolicyTreatyIn.PolicyNo", kanan: operand{teksLiteral, ""}}}},
		"UJI_LAIN":   {kondisi: map[string]kondisi{"A": {jenis: banding, kiri: ".BusinessCode", kanan: operand{teksLiteral, "02"}}}},
		"UJI_PANIK":  {panik: "baris A belum diport: @String.equals(OperatorID.PYUSERIDENTIFIER,\"UJI\")"},
	}
	got := strings.Join(literalIdentitas(reg), "\n")
	for _, mau := range []string{"UJI_KECIL baris A", "UJI_BESAR baris A", "UJI_PANIK: pesan panik"} {
		if !strings.Contains(got, mau) {
			t.Errorf("tidak tertangkap: %s\n%s", mau, got)
		}
	}
	for _, lolos := range []string{"UJI_KOSONG", "UJI_LAIN"} {
		if strings.Contains(got, lolos) {
			t.Errorf("salah tuduh: %s\n%s", lolos, got)
		}
	}
}
