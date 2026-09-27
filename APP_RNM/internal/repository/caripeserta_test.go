package repository

import (
	"strings"
	"testing"
)

// Uji penyusun SQL pencarian peserta - LoadDataPesertaSpesifik_Act.
//
// Sumber, dibaca sebagai pohon langkah utuh:
//
//	`Activity/LoadDataPesertaSpesifik_Act.xml`
//	  b223  Page-Remove TempDetail1
//	  b405  SearchPolicyHolder.CARI3 = @toUpperCase(SearchPolicyHolder.CARI3)
//	  b485  RDB-List TempDetail1  -> `RDBList/GetPesertaClaim_sql1.xml:85`
//	          SELECT * FROM POOLDATA.M_LIFE_PREMIUM_DETAIL
//	          WHERE PL_NUMBER = {..PremiumListSummary.PL_NUMBER}
//	            AND CERTIFICATE_NO LIKE '%'||{SearchPolicyHolder.CARI2}||'%'
//	            AND UPPER(NAME_OF_INSURED) LIKE '%'||{SearchPolicyHolder.CARI3}||'%'

func TestSQLCariPesertaSelaluBerpagar(t *testing.T) {
	// ⛔ Dua pagar yang TIDAK BOLEH hilang dalam bentuk apa pun, sebab tabel
	// sumbernya 66,8 juta baris: penyaring ber-index dan batas hasil.
	kasus := []struct{ nama, sertifikat, cari string }{
		{"tanpa penyaring tambahan", "", ""},
		{"dengan sertifikat", "UJI-001", ""},
		{"dengan nama", "", "budi"},
		{"dengan keduanya", "UJI-001", "budi"},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			q, arg := sqlCariPeserta("SKEMAUJI.M_LIFE_PREMIUM_DETAIL",
				"UJI-PL-1", k.sertifikat, k.cari, 50)
			if !strings.Contains(q, "PL_NUMBER = :1") {
				t.Errorf("penyaring ber-index PL_NUMBER hilang:\n%s", q)
			}
			if !strings.Contains(q, "FETCH FIRST 50 ROWS ONLY") {
				t.Errorf("batas hasil hilang:\n%s", q)
			}
			if !strings.Contains(q, penyaringHidup) {
				t.Errorf("penyaring peserta hidup hilang:\n%s", q)
			}
			if arg[0] != "UJI-PL-1" {
				t.Errorf("bind pertama %v, mau nomor premium list", arg[0])
			}
			if err := PeriksaSQL(q); err != nil {
				t.Errorf("PeriksaSQL menolak: %v", err)
			}
		})
	}
}

func TestSQLCariPesertaPenyaringPilihan(t *testing.T) {
	t.Run("penyaring kosong TIDAK menambah LIKE", func(t *testing.T) {
		// ⛔ Ini penyimpangan SADAR dari Pega, dan sebabnya harus tertulis:
		// di Oracle `X LIKE '%'` bernilai FALSE ketika X NULL. Pega memasang
		// kedua LIKE tanpa syarat, sehingga kotak pencarian yang KOSONG pun
		// membuang setiap peserta yang NAME_OF_INSURED-nya NULL - baris yang
		// hilang tanpa seorang pun memintanya. Kita memasang LIKE hanya bila
		// kotaknya terisi. Dilaporkan sebagai OQ-E.
		q, arg := sqlCariPeserta("T", "PL", "", "", 10)
		if strings.Contains(q, "LIKE") {
			t.Errorf("penyaring kosong seharusnya tidak memasang LIKE:\n%s", q)
		}
		if len(arg) != 1 {
			t.Errorf("bind %d, mau 1 (hanya PL_NUMBER)", len(arg))
		}
	})

	t.Run("nama di-UPPERCASE, sertifikat TIDAK", func(t *testing.T) {
		// b405 hanya meng-uppercase CARI3 (nama). Sertifikat dibandingkan
		// apa adanya, dan SQL-nya pun tidak membungkus CERTIFICATE_NO dengan
		// UPPER - meng-uppercase keduanya akan membuat pencarian sertifikat
		// tidak lagi memakai index-nya.
		q, arg := sqlCariPeserta("T", "PL", "uji-abc", "budi santoso", 10)
		if !strings.Contains(q, "UPPER(NAME_OF_INSURED) LIKE") {
			t.Errorf("nama tidak dibandingkan ber-UPPER:\n%s", q)
		}
		if strings.Contains(q, "UPPER(CERTIFICATE_NO)") {
			t.Errorf("sertifikat TIDAK boleh dibungkus UPPER:\n%s", q)
		}
		if arg[1] != "uji-abc" {
			t.Errorf("bind sertifikat %q, mau apa adanya %q", arg[1], "uji-abc")
		}
		if arg[2] != "BUDI SANTOSO" {
			t.Errorf("bind nama %q, mau huruf besar", arg[2])
		}
	})

	t.Run("hanya nama: bind-nya bergeser ke :2", func(t *testing.T) {
		// Penomoran bind mengikuti penyaring yang benar-benar terpasang.
		// Nomor yang meleset membuat Oracle menjawab galat bind, atau lebih
		// buruk: membandingkan kolom dengan nilai milik kolom lain.
		q, arg := sqlCariPeserta("T", "PL", "", "ani", 10)
		if !strings.Contains(q, "UPPER(NAME_OF_INSURED) LIKE '%'||:2||'%'") {
			t.Errorf("nama tidak terpasang di :2:\n%s", q)
		}
		if len(arg) != 2 || arg[1] != "ANI" {
			t.Errorf("bind = %v, mau [PL ANI]", arg)
		}
	})

	t.Run("spasi di tepi dirapikan, isi tidak", func(t *testing.T) {
		_, arg := sqlCariPeserta("T", "PL", "  UJI-1  ", "  budi  ", 10)
		if arg[1] != "UJI-1" || arg[2] != "BUDI" {
			t.Errorf("bind = %v, mau [PL UJI-1 BUDI]", arg)
		}
	})
}
