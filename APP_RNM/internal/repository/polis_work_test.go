package repository

// Uji bentuk query kerja polis - tiket 01 PremiumList Life.
//
// ⛔ Inti berkas ini SATU ralat: kolom mana menyimpan apa. `POSITION` adalah
// posisi LAYAR (`Offer`/`Premium`, `ProtectAccept.xml` b1207/b2288);
// `STATUS` adalah `pyWorkStatus` - nama tahap selama berjalan, `Resolved-*`
// saat tertutup. Ronde pertama menukarnya, dan akibatnya kotak masuk yang
// kosong untuk pekerjaan yang benar-benar ada.

import (
	"strings"
	"testing"

	"nusantarare/internal/models"
)

const tabelUjiWorkPolis = "SKEMAUJI.T_WORK_POLIS"

func TestPindahTahapMenulisStatusBukanPosition(t *testing.T) {
	q := sqlPindahTahapPolis(tabelUjiWorkPolis)
	if !strings.Contains(q, "SET STATUS = :1") {
		t.Errorf("perpindahan tahap tidak menulis STATUS:\n%s", q)
	}
	// ⛔ POSITION TIDAK disentuh. Ia posisi layar, dan perpindahan tahap
	// bukan perpindahan layar - menimpanya akan membuat `ProtectAccept`
	// memilih pemeriksaan yang salah.
	if strings.Contains(q, "SET POSITION") || strings.Contains(q, "POSITION =") {
		t.Errorf("perpindahan tahap menyentuh POSITION:\n%s", q)
	}
	// Syarat optimistik atas nilai LAMA - baris dibaca di luar transaksi.
	if !strings.Contains(q, "STATUS = :3") {
		t.Errorf("perpindahan tanpa syarat status lama; dua permintaan "+
			"serentak akan sama-sama menang:\n%s", q)
	}
}

func TestTutupPolisMenolakPenutupanKedua(t *testing.T) {
	q := sqlTutupPolis(tabelUjiWorkPolis)
	// ⛔ POSITION dikosongkan: kotak masuk adalah worklist, dan kasus
	// tertutup tidak berdiri di antrean mana pun.
	if !strings.Contains(q, "POSITION = NULL") {
		t.Errorf("penutupan tidak mengosongkan POSITION:\n%s", q)
	}
	// ⛔ Syaratnya BUKAN `STATUS IS NULL`. Sejak ralat 28-09-2026 kolom itu
	// juga menyimpan nama tahap selama kasus berjalan, jadi `IS NULL` hanya
	// benar untuk kasus yang belum pernah bertahap - dan kasus bertahap
	// tidak akan pernah dapat ditutup.
	if !strings.Contains(q, "STATUS NOT IN (:3, :4)") {
		t.Errorf("syarat penutupan tidak menolak status akhir:\n%s", q)
	}
	if strings.Contains(q, "AND STATUS IS NULL\n") {
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
		if err := PeriksaSQL(q); err != nil {
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
	if !strings.Contains(q, "(STATUS = :5 OR (STATUS IS NULL AND :5 IS NULL))") {
		t.Errorf("penutupan polis tidak dikunci tahap asalnya:\n%s", q)
	}
	if !strings.Contains(q, "STATUS NOT IN (:3, :4)") {
		t.Errorf("penjaga kasus tertutup hilang:\n%s", q)
	}
}
