package repository

// Uji bentuk query kerja polis - tiket 01 PremiumList Life.
//
// ⛔ Inti berkas ini SATU ralat: kolom mana menyimpan apa. `POSITION` adalah
// posisi LAYAR (`Offer`/`Premium`, `ProtectAccept.xml` b1207/b2288);
// `STATUS_WORK` (bernama `STATUS` sampai migrasi 059) adalah `pyWorkStatus` -
// nama tahap selama berjalan, `Resolved-*` saat tertutup. Ronde pertama menukarnya, dan akibatnya kotak masuk yang
// kosong untuk pekerjaan yang benar-benar ada.

import (
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/premiumlistlife/backend/models"
)

const tabelUjiWorkPolis = "SKEMAUJI.T_WORK_POLIS"

func TestPindahTahapMenulisStatusBukanPosition(t *testing.T) {
	q := sqlPindahTahapPolis(tabelUjiWorkPolis)
	if !strings.Contains(q, "SET STATUS_WORK = :1, TGL_UPDATE = SYSDATE") {
		t.Errorf("perpindahan tahap tidak menulis STATUS_WORK dan TGL_UPDATE:\n%s", q)
	}
	// ⛔ POSITION TIDAK disentuh. Ia posisi layar, dan perpindahan tahap
	// bukan perpindahan layar - menimpanya akan membuat `ProtectAccept`
	// memilih pemeriksaan yang salah.
	if strings.Contains(q, "SET POSITION") || strings.Contains(q, "POSITION =") {
		t.Errorf("perpindahan tahap menyentuh POSITION:\n%s", q)
	}
	// Syarat optimistik atas nilai LAMA - baris dibaca di luar transaksi.
	if !strings.Contains(q, "STATUS_WORK = :3") {
		t.Errorf("perpindahan tanpa syarat status lama; dua permintaan "+
			"serentak akan sama-sama menang:\n%s", q)
	}
}

func TestTutupPolisMenolakPenutupanKedua(t *testing.T) {
	q := sqlTutupPolis(tabelUjiWorkPolis)
	// ⛔ POSITION dikosongkan: kotak masuk adalah worklist, dan kasus
	// tertutup tidak berdiri di antrean mana pun.
	if !strings.Contains(q, "POSITION = NULL, TGL_UPDATE = SYSDATE") {
		t.Errorf("penutupan tidak mengosongkan POSITION dan menulis TGL_UPDATE:\n%s", q)
	}
	// ⛔ Syaratnya BUKAN `STATUS IS NULL`. Sejak ralat 28-09-2026 kolom itu
	// juga menyimpan nama tahap selama kasus berjalan, jadi `IS NULL` hanya
	// benar untuk kasus yang belum pernah bertahap - dan kasus bertahap
	// tidak akan pernah dapat ditutup.
	if !strings.Contains(q, "STATUS_WORK NOT IN (:3, :4)") {
		t.Errorf("syarat penutupan tidak menolak status akhir:\n%s", q)
	}
	if strings.Contains(q, "AND STATUS_WORK IS NULL\n") {
		t.Errorf("syarat masih `STATUS IS NULL`; kasus bertahap tidak akan "+
			"pernah dapat ditutup:\n%s", q)
	}
}

func TestQueryKerjaPolisMemakaiBind(t *testing.T) {
	for nama, q := range map[string]string{
		"keadaan": sqlKeadaanPolis(tabelUjiWorkPolis),
		"pindah":  sqlPindahTahapPolis(tabelUjiWorkPolis),
		"tutup":   sqlTutupPolis(tabelUjiWorkPolis),
	} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Errorf("%s: %v\n%s", nama, err, q)
		}
		if !strings.Contains(q, ":1") {
			t.Errorf("%s: nol bind:\n%s", nama, q)
		}
		// Nol nilai status tertempel ke teks.
		for _, nilai := range []string{
			models.StatusPolisDitolak, models.StatusPolisSelesai,
			models.PosisiOffer, models.PosisiPremium,
		} {
			if strings.Contains(q, nilai) {
				t.Errorf("%s: nilai %q tertempel ke teks query:\n%s", nama, nilai, q)
			}
		}
	}
}

// TestTutupPolisDikunciTahapYangDibaca - temuan /code-review giliran 10.
//
// ⛔ Keadaan dibaca di luar transaksi; penutupan yang tidak menyebut tahap
// asalnya menutup kasus yang sudah berpindah tahap di antara baca dan tulis.
func TestTutupPolisDikunciTahapYangDibaca(t *testing.T) {
	q := sqlTutupPolis("S.W")
	if !strings.Contains(q, "(STATUS_WORK = :5 OR (STATUS_WORK IS NULL AND :5 IS NULL))") {
		t.Errorf("penutupan polis tidak dikunci tahap asalnya:\n%s", q)
	}
	if !strings.Contains(q, "STATUS_WORK NOT IN (:3, :4)") {
		t.Errorf("penjaga kasus tertutup hilang:\n%s", q)
	}
}

// TestBenderaPolisDibacaTerpisah - GILIRAN-14 butir bq.
//
// ⛔ Keadaan kasus TIDAK membaca kolom 057; hanya pembaca bendera yang
// membacanya, supaya `Reject`/`Decline` tidak ikut gagal selama 057 belum
// berjalan di sebuah skema.
func TestBenderaPolisDibacaTerpisah(t *testing.T) {
	if q := sqlKeadaanPolis(tabelUjiWorkPolis); strings.Contains(q, "FLAG_ONGOING_POLICY") {
		t.Errorf("keadaan polis membaca kolom 057:\n%s", q)
	}
	q := sqlBenderaPolis(tabelUjiWorkPolis)
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(q, "FLAG_ONGOING_POLICY") || !strings.Contains(q, ":1") {
		t.Errorf("pembaca bendera:\n%s", q)
	}
}

// TestSQLKerjaPolisSeragamDenganWorkClaim - brief seragam kolom 01-10-2026 §4.
//
// ⛔ SETIAP SQL `T_WORK_POLIS` memakai `STATUS_WORK`; nama `STATUS` sudah tidak
// ada sejak 059 dan pernyataan yang masih menyebutnya mati di ORA-00904.
// Setiap ubah baris kasus menulis `TGL_UPDATE = SYSDATE` di pernyataan yang
// sama. JSON API dan teks layar TIDAK berubah - hanya nama kolomnya.
func TestSQLKerjaPolisSeragamDenganWorkClaim(t *testing.T) {
	statusLama := regexp.MustCompile(`\bSTATUS\b`)
	for nama, q := range map[string]string{
		"keadaan": sqlKeadaanPolis(tabelUjiWorkPolis),
		"pindah":  sqlPindahTahapPolis(tabelUjiWorkPolis),
		"tutup":   sqlTutupPolis(tabelUjiWorkPolis),
		"sisip":   sqlSisipKasusPolis(tabelUjiWorkPolis),
		"inbox":   sqlInboxPolis(tabelUjiWorkPolis, "S.T_PREMIUM_LIST", "S.T_PREMIUM_LIST_DETAIL"),
	} {
		if statusLama.MatchString(q) {
			t.Errorf("%s masih menyebut kolom STATUS:\n%s", nama, q)
		}
		if !strings.Contains(q, "STATUS_WORK") {
			t.Errorf("%s tidak menyebut STATUS_WORK:\n%s", nama, q)
		}
	}
	for nama, q := range map[string]string{
		"pindah": sqlPindahTahapPolis(tabelUjiWorkPolis),
		"tutup":  sqlTutupPolis(tabelUjiWorkPolis),
	} {
		if !strings.Contains(q, "TGL_UPDATE = SYSDATE") {
			t.Errorf("%s mengubah baris kasus tanpa TGL_UPDATE:\n%s", nama, q)
		}
	}
	// Keadaan membaca pembuat dan waktu - empat medan + COVER_KEY (brief §4).
	if q := sqlKeadaanPolis(tabelUjiWorkPolis); !strings.Contains(q,
		"SELECT ID, LINI, POSITION, STATUS_WORK, COVER_KEY, CREATE_OP, CREATE_OP_NAME, TGL_CREATE, TGL_UPDATE FROM") {
		t.Errorf("keadaan polis:\n%s", q)
	}
}
