package repository

import (
	"regexp"
	"strconv"
	"testing"
)

// Oracle mengikat bind SQL menurut urutan kemunculan, bukan nomor `:n`: nomor bind SQL lampiran wajib naik 1, 2, 3, ...
// sesuai urutan teks (09-10-2026: `:2, :3 ... :1` membuat nol kategori terbaca di DEV).
func TestBindLampiranUrutKemunculan(t *testing.T) {
	pola := regexp.MustCompile(`:(\d+)`)
	for nama, q := range map[string]string{
		"kategori": sqlKategoriLampiran("S.K", "S.D"),
		"daftar":   sqlDaftarLampiran("S.D"),
		"sisip":    sqlSisipDokumenKlaim("S.D"),
		"hapus":    sqlHapusDokumenKlaim("S.D"),
		"pindah":   sqlPindahKategoriDokumen("S.D"),
		"anak":     sqlAnakSpreading("S.P"),
	} {
		for i, m := range pola.FindAllStringSubmatch(q, -1) {
			if n, _ := strconv.Atoi(m[1]); n != i+1 {
				t.Errorf("%s: bind ke-%d bernomor :%d - nomor harus mengikuti urutan kemunculan", nama, i+1, n)
			}
		}
	}
}
