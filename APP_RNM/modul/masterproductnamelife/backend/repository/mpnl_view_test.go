package repository

// Uji murni brief bab 4: daftar kunci yang dibaca view di kode SAMA dengan
// teks view DEV (`docs/dba-view-produk-life.md`), peka huruf besar-kecil.

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
)

func kunciDiDokumen(t *testing.T, view string) []string {
	t.Helper()
	isi, err := os.ReadFile("../../docs/dba-view-produk-life.md")
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	awal := strings.Index(s, "=== VIEW "+view+"\n")
	if awal < 0 {
		t.Fatalf("view %s tidak ada di dokumen", view)
	}
	blok := s[awal+len("=== VIEW "+view):]
	if akhir := strings.Index(blok, "=== VIEW "); akhir >= 0 {
		blok = blok[:akhir]
	}
	set := map[string]bool{}
	for _, m := range regexp.MustCompile(`a\.JSONDATA\.(\w+)`).FindAllStringSubmatch(blok, -1) {
		set[m[1]] = true
	}
	hasil := make([]string, 0, len(set))
	for k := range set {
		hasil = append(hasil, k)
	}
	sort.Strings(hasil)
	return hasil
}

func urut(d []string) []string {
	h := append([]string{}, d...)
	sort.Strings(h)
	return h
}

func TestKunciViewSamaDenganDokumenDBA(t *testing.T) {
	if got, mau := urut(KunciViewProduk), kunciDiDokumen(t, "PRODUCT_LIFE"); strings.Join(got, ",") != strings.Join(mau, ",") {
		t.Errorf("PRODUCT_LIFE:\nkode    %v\ndokumen %v", got, mau)
	}
	if got, mau := urut(KunciViewInward), kunciDiDokumen(t, "PRODUCTINWARD_LIFE"); strings.Join(got, ",") != strings.Join(mau, ",") {
		t.Errorf("PRODUCTINWARD_LIFE:\nkode    %v\ndokumen %v", got, mau)
	}
	if !strings.Contains(strings.Join(kunciDiDokumen(t, "PRODUCT_LIFE"), ","), "OutwardList") {
		t.Error("PRODUCT_LIFE membaca OutwardList[0].OVR_COMM (P3)")
	}
}
