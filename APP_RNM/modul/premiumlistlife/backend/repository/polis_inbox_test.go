package repository

// Uji bentuk query kotak masuk PremiumList - tiket 01 bagian 2.

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

const (
	tabelUjiWorkPolisInbox = "SKEMAUJI.T_WORK_POLIS"
	tabelUjiPolis          = "SKEMAUJI.T_PREMIUM_LIST"
	tabelUjiPolisDetail    = "SKEMAUJI.T_PREMIUM_LIST_DETAIL"
)

func TestInboxPolisMemakaiLeftJoinKeDetail(t *testing.T) {
	q := sqlInboxPolis(tabelUjiWorkPolisInbox, tabelUjiPolis, tabelUjiPolisDetail)
	// ⛔ LEFT JOIN, bukan JOIN: polis yang CSV-nya belum diunggah belum punya
	// baris detail sama sekali, dan ia justru yang paling perlu tampil di
	// kotak masuk — pekerjaannya belum selesai.
	if !strings.Contains(q, "LEFT JOIN") {
		t.Errorf("polis tanpa detail akan hilang dari kotak masuk:\n%s", q)
	}
	// ⛔ PL_NUMBER dari tabel DETAIL. Jalur propertinya di Pega menyebut
	// `PremiumListSummary`, tetapi di skema kami kolomnya hidup di 052.
	if !strings.Contains(q, "MAX(d.PL_NUMBER)") {
		t.Errorf("PL_NUMBER tidak dibaca dari tabel detail:\n%s", q)
	}
	// Urutan STABIL dengan pemutus seri: daftar tanpa pemutus seri
	// menampilkan baris yang berpindah sendiri di antara dua halaman, dan
	// baris yang berpindah dapat TERLEWAT.
	if !strings.Contains(q, "ORDER BY p.TGL_INPUT DESC, w.ID") {
		t.Errorf("urutan tidak stabil:\n%s", q)
	}
	if !strings.Contains(q, "OFFSET :3 ROWS FETCH NEXT :4 ROWS ONLY") {
		t.Errorf("halaman tidak dijepit:\n%s", q)
	}
}

func TestCacahInboxPolisQueryTersendiri(t *testing.T) {
	// ⛔ Query TERSENDIRI, bukan `COUNT(*) OVER ()` yang ditempel. Jendela
	// itu dihitung per baris halaman; pada halaman KOSONG — yang terjadi
	// setiap kali seseorang membuka halaman terakhir lalu satu baris hilang —
	// ia tidak mengembalikan apa pun, dan totalnya menjadi nol untuk antrean
	// yang tidak kosong.
	q := sqlCacahInboxPolis(tabelUjiWorkPolisInbox)
	if !strings.Contains(q, "SELECT COUNT(*)") {
		t.Errorf("pencacah bukan COUNT:\n%s", q)
	}
	if strings.Contains(q, "OVER ()") {
		t.Errorf("pencacah memakai fungsi jendela:\n%s", q)
	}
	if strings.Contains(q, "OFFSET") || strings.Contains(q, "FETCH") {
		t.Errorf("pencacah ikut terjepit halaman:\n%s", q)
	}
}

func TestSaringanPosisiOpsional(t *testing.T) {
	// ⚠️ Posisi KOSONG berarti seluruh posisi, bukan "tidak ada". Tab
	// "semua" memakainya, dan `POSITION = ''` akan mengembalikan nol baris.
	// ⛔ Kueri HALAMAN memakai penampung unik (`:1`, `:2`) - penampung
	// berulang bersama `OFFSET … FETCH` dijawab ORA-01008 (uji asap DEV,
	// GILIRAN-12). Kueri cacah tanpa pembatas baris terbukti aman berulang.
	for q, mau := range map[string]string{
		sqlInboxPolis(tabelUjiWorkPolisInbox, tabelUjiPolis, tabelUjiPolisDetail): ":1 IS NULL OR w.POSITION = :2",
		sqlCacahInboxPolis(tabelUjiWorkPolisInbox):                                ":1 IS NULL OR w.POSITION = :1",
	} {
		if !strings.Contains(q, mau) {
			t.Errorf("saringan posisi tidak opsional (mau %q):\n%s", mau, q)
		}
	}
}

func TestBatasUkuranHalamanPolisDijepit(t *testing.T) {
	// ⛔ `?ukuran=1000000` adalah permintaan yang memuat seluruh tabel polis
	// ke memori satu proses.
	if got := BatasUkuranHalamanPolis(0); got != UkuranHalamanPolisBawaan {
		t.Errorf("ukuran 0 -> %d, mau %d", got, UkuranHalamanPolisBawaan)
	}
	if got := BatasUkuranHalamanPolis(-5); got != UkuranHalamanPolisBawaan {
		t.Errorf("ukuran negatif -> %d", got)
	}
	if got := BatasUkuranHalamanPolis(1000000); got != 200 {
		t.Errorf("ukuran sejuta -> %d, mau 200", got)
	}
	if got := BatasUkuranHalamanPolis(50); got != 50 {
		t.Errorf("ukuran wajar -> %d, mau 50", got)
	}
}

func TestQueryInboxPolisBersih(t *testing.T) {
	for nama, q := range map[string]string{
		"halaman": sqlInboxPolis(tabelUjiWorkPolisInbox, tabelUjiPolis, tabelUjiPolisDetail),
		"cacah":   sqlCacahInboxPolis(tabelUjiWorkPolisInbox),
	} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Errorf("%s: %v\n%s", nama, err, q)
		}
		if strings.Contains(q, ";") || strings.Contains(q, "--") {
			t.Errorf("%s: memuat pemisah pernyataan atau komentar:\n%s", nama, q)
		}
	}
}

// TestKetentuanUnderwritingTidakDikarang menjaga satu celah tetap TERBUKA.
//
// ⛔ `InboxPremiumList.xml` b891 mendaftarkan kolom
// `A.KetentuanUnderwriting`, dan migrasi 050-056 TIDAK punya kolomnya.
// Kolom itu karena itu tidak ditampilkan — bukan diisi teks kosong, yang di
// layar terbaca "memang kosong" alih-alih "kami tidak punya datanya".
//
// Uji ini menyala bila seseorang menambahkannya tanpa menambah kolomnya.
func TestKetentuanUnderwritingTidakDikarang(t *testing.T) {
	q := sqlInboxPolis(tabelUjiWorkPolisInbox, tabelUjiPolis, tabelUjiPolisDetail)
	if strings.Contains(strings.ToUpper(q), "KETENTUAN") {
		t.Error("KetentuanUnderwriting dibaca padahal kolomnya tidak ada di " +
			"migrasi 050-056; bila kolomnya sudah dibuat, ralat uji ini " +
			"beserta tiket 01")
	}
}
