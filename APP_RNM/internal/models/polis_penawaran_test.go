package models_test

// Uji tangga kerja polis - tiket 01 PremiumList Life.
//
// ⛔ Yang dijaga di sini bukan "kode berjalan" melainkan bahwa TABELNYA SAMA
// dengan peta konektor `InputPolicyHolder.xml`. Tabel transisi yang meleset
// satu sel memindahkan kasus ke tempat yang salah, dan tidak ada satu pun
// galat yang akan berbunyi - kasusnya hanya hilang dari antrean yang benar.

import (
	"errors"
	"os"
	"strings"
	"testing"

	"nusantarare/internal/models"
)

func TestTabelTransisiSamaDenganPetaKonektor(t *testing.T) {
	for _, u := range []struct {
		nama      string
		tahap     string
		keputusan string
		mauTahap  string
		mauStatus string
		mauGolong bool
		konektor  string
	}{
		{"penawaran + Confirm -> penggolong", models.TahapPolisPenawaran,
			models.KeputusanConfirm, "", "", true, "Transition4 b1574"},
		{"penawaran + Decline -> tutup ditolak", models.TahapPolisPenawaran,
			models.KeputusanDecline, "", models.StatusPolisDitolak, false, "Transition5 b2235"},
		{"detail + Confirm -> tutup selesai", models.TahapPolisDetail,
			models.KeputusanConfirm, "", models.StatusPolisSelesai, false, "Transition7 b2090"},
		{"detail + Decline -> tutup ditolak", models.TahapPolisDetail,
			models.KeputusanDecline, "", models.StatusPolisDitolak, false, "Transition6 b2162"},
		{"detail + Reject -> kembali ke penawaran", models.TahapPolisDetail,
			models.KeputusanReject, models.TahapPolisPenawaran, "", false, "Transition9 b2306"},
	} {
		got, err := models.TransisiPenawaran(u.tahap, u.keputusan)
		if err != nil {
			t.Errorf("%s (%s): %v", u.nama, u.konektor, err)
			continue
		}
		if got.TahapTujuan != u.mauTahap {
			t.Errorf("%s: tahap = %q, mau %q", u.nama, got.TahapTujuan, u.mauTahap)
		}
		if got.StatusWork != u.mauStatus {
			t.Errorf("%s: status = %q, mau %q", u.nama, got.StatusWork, u.mauStatus)
		}
		if got.MenungguPenggolong != u.mauGolong {
			t.Errorf("%s: menunggu penggolong = %v, mau %v",
				u.nama, got.MenungguPenggolong, u.mauGolong)
		}
	}
}

func TestRejectDiTahapPenawaranTidakPunyaJalur(t *testing.T) {
	// ⛔ RALAT ATAS TIKET 01, dan uji inilah yang menguncinya.
	//
	// Tiket menulis "Reject pada tahap MANA PUN mengembalikan case ke Input
	// Offer". Konektornya membantah: `Reject` b2306 hanya ada pada
	// `Decision2` - penggolong SESUDAH Input Premium Detail. `Decision1`,
	// penggolong sesudah tahap penawaran, hanya punya `Confirm` b1574 dan
	// `Decline` b2235.
	//
	// Menyediakan Reject di tahap penawaran berarti membuat jalur yang tidak
	// pernah ada di sistem lama.
	_, err := models.TransisiPenawaran(models.TahapPolisPenawaran, models.KeputusanReject)
	if !errors.Is(err, models.ErrKeputusanTidakAdaDiTahapIni) {
		t.Fatalf("galat = %v, mau ErrKeputusanTidakAdaDiTahapIni", err)
	}
	// Dan galatnya menyebut KEDUANYA - keputusan dan tahapnya - supaya yang
	// membacanya tahu bukan keputusannya yang salah, melainkan tempatnya.
	if !strings.Contains(err.Error(), models.KeputusanReject) ||
		!strings.Contains(err.Error(), models.TahapPolisPenawaran) {
		t.Errorf("pesan tidak menyebut keputusan dan tahapnya: %v", err)
	}
}

func TestTahapSummaryNolKonektorKeputusan(t *testing.T) {
	// `Assignment1` hanya punya `Transition2` b1866 yang KELUAR, dan itu
	// penyelesaian layar summary - bukan keputusan.
	for _, k := range []string{
		models.KeputusanConfirm, models.KeputusanReject, models.KeputusanDecline,
	} {
		if _, err := models.TransisiPenawaran(models.TahapPolisSummary, k); !errors.Is(
			err, models.ErrKeputusanTidakAdaDiTahapIni) {
			t.Errorf("%s di summary: %v, mau ErrKeputusanTidakAdaDiTahapIni", k, err)
		}
	}
}

func TestKeputusanDanTahapAsingDitolak(t *testing.T) {
	if _, err := models.TransisiPenawaran(models.TahapPolisDetail, "Approve"); !errors.Is(
		err, models.ErrKeputusanTidakDikenal) {
		t.Errorf("keputusan karangan: %v", err)
	}
	if _, err := models.TransisiPenawaran("Input Something", models.KeputusanConfirm); !errors.Is(
		err, models.ErrTahapPolisTidakDikenal) {
		t.Errorf("tahap karangan: %v", err)
	}
}

func TestPenggolongOfferMenutupPremiumMelanjutkan(t *testing.T) {
	// ⛔ `Offer` MENUTUP (Transition11 b1807 -> END52), tidak menunggu.
	offer, err := models.LanjutanPenggolong(models.LanjutOffer)
	if err != nil {
		t.Fatal(err)
	}
	if offer.StatusWork != models.StatusPolisSelesai || offer.TahapTujuan != "" {
		t.Errorf("Offer = %+v, mau tutup Resolved-Completed", offer)
	}
	premium, err := models.LanjutanPenggolong(models.LanjutPremium)
	if err != nil {
		t.Fatal(err)
	}
	if premium.TahapTujuan != models.TahapPolisDetail || premium.Ditutup() {
		t.Errorf("Premium = %+v, mau pindah ke %q", premium, models.TahapPolisDetail)
	}
}

func TestKasusPolisTertutupTepatBukanAwalan(t *testing.T) {
	// ⛔ Perbandingan TEPAT. Pega punya banyak status berawalan `Resolved-`;
	// menerima awalan berarti menerima status yang alur ini tidak pernah
	// hasilkan - dan kasus yang berstatus begitu akan tampak tertutup
	// padahal ia tersesat.
	if !models.KasusPolisTertutup(models.StatusPolisDitolak) ||
		!models.KasusPolisTertutup(models.StatusPolisSelesai) {
		t.Error("status akhir yang sah tidak terbaca tertutup")
	}
	for _, asing := range []string{
		"", "Resolved-Discarded", "Resolved", "Open", "resolved-completed",
	} {
		if models.KasusPolisTertutup(asing) {
			t.Errorf("status %q terbaca tertutup", asing)
		}
	}
}

// TestNamaTahapDanStatusVERBATIMDariKorpus membaca flow-nya LANGSUNG.
//
// ⛔ Penjaga yang membandingkan salinan dengan salinan tidak menjaga apa pun.
// Nilai-nilai ini masuk `T_WORK_POLIS.POSITION` dan `.STATUS`; yang meleset
// membuat kotak masuk kosong untuk baris yang sebenarnya ada.
func TestNamaTahapDanStatusVERBATIMDariKorpus(t *testing.T) {
	const letak = `D:\XML\RNM_BRD\PremiumList Life\InputPolicyHolder.xml`
	isi, err := os.ReadFile(letak)
	if err != nil {
		t.Skipf("korpus tidak terjangkau di mesin ini (%v); nilai tidak terperiksa", err)
	}
	teks := string(isi)
	for _, nilai := range []string{
		models.StatusPolisDitolak, models.StatusPolisSelesai,
		models.TahapPolisPenawaran, models.TahapPolisDetail, models.TahapPolisSummary,
	} {
		if !strings.Contains(teks, "<pyWorkStatus>"+nilai+"</pyWorkStatus>") {
			t.Errorf("%q bukan pyWorkStatus mana pun di flow", nilai)
		}
	}
	for _, k := range []string{
		models.KeputusanConfirm, models.KeputusanReject, models.KeputusanDecline,
		models.LanjutOffer, models.LanjutPremium,
	} {
		if !strings.Contains(teks, "<pyExpression>"+k+"</pyExpression>") {
			t.Errorf("%q bukan pyExpression konektor mana pun", k)
		}
	}
	// ⛔ Dan `Reject` muncul TEPAT SEKALI - kalau suatu hari dua, ralat
	// tiket di atas harus dibaca ulang.
	if n := strings.Count(teks, "<pyExpression>Reject</pyExpression>"); n != 1 {
		t.Errorf("pyExpression Reject muncul %d kali, mau 1 (hanya Decision2 b2306); "+
			"bila bertambah, ralat tiket 01 harus dibaca ulang", n)
	}
}
