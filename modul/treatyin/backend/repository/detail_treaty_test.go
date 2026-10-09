package repository_test

// Bentuk SQL tulisan detail (`SaveTreatyInDetail(Edm)`), tanpa basis data.
// Ekspresi VALUES-nya sudah dijalankan BACA-SAJA di DEV sebagai
// `SELECT … FROM DUAL` (9 Oktober 2026); uji ini menjaga bentuknya.

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/repository"
)

func TestSQLDetailTreatyBentukProsedur(t *testing.T) {
	for _, c := range []struct {
		nama, kepala string
		kolom        []models.KolomDetail
		jumlahKolom  int
	}{
		{"S.TREATYINDETAIL", "S.TREATY_IN", models.KolomDetailTreatyIn, 64},
		{"S.TREATYINDETAILEDM", "S.TREATY_IN_EDM", models.KolomDetailTreatyInEDM, 55},
	} {
		hapus, tertinggi, sisip, err := repository.SQLDetailTreaty(c.nama, c.kepala, "S.TREATYINDETAIL", "S.TREATYINDETAILEDM", c.kolom)
		if err != nil {
			t.Fatal(err)
		}
		// RemoveTreatyInDetail(Edm): satu tabel, satu kontrak.
		if hapus != "DELETE FROM "+c.nama+" WHERE TREATYID = :1" {
			t.Errorf("hapus %q", hapus)
		}
		// Sequence Pega dibagi KEDUA tabel — angka tertinggi membaca keduanya.
		if !strings.Contains(tertinggi, "S.TREATYINDETAIL WHERE") || !strings.Contains(tertinggi, "S.TREATYINDETAILEDM WHERE") {
			t.Errorf("tertinggi tidak membaca kedua tabel: %s", tertinggi)
		}
		kolom := sisip[strings.Index(sisip, "(")+1 : strings.Index(sisip, ") VALUES")]
		if n := len(strings.Split(kolom, ", ")); n != c.jumlahKolom {
			t.Errorf("%s: %d kolom, mau %d (kolom tabel di DEV)", c.nama, n, c.jumlahKolom)
		}
		// Bind :1..:n berurutan tanpa celah dan tanpa dipakai ulang — Oracle
		// mengikat menurut urutan kemunculan.
		bind := regexp.MustCompile(`:(\d+)`).FindAllStringSubmatch(sisip, -1)
		angka := 0
		for _, k := range c.kolom {
			if k.Angka {
				angka++
			}
		}
		if mau := 1 + len(c.kolom) + angka + 2; len(bind) != mau {
			t.Errorf("%s: %d bind, mau %d", c.nama, len(bind), mau)
		}
		for i, b := range bind {
			if b[1] != strconv.Itoa(i+1) {
				t.Fatalf("%s: bind ke-%d bernomor :%s", c.nama, i+1, b[1])
			}
		}
		// Tanggal dari kepala yang benar, NULL bila tidak tepat satu baris.
		if !strings.Contains(sisip, "FROM "+c.kepala+" WHERE ID = ") || !strings.Contains(sisip, "COUNT(*) = 1") {
			t.Errorf("%s: tanggal tidak dari %s", c.nama, c.kepala)
		}
	}
}
