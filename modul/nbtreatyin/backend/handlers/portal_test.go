package handlers_test

// Uji seam 1 - GERBANG DAFTAR PORTAL (putaran 2, paket P8; AC 11, 14, 92).
//
// XML `Section/SFAPortal_OpportunitiesList.xml`: satu-satunya grid portal
// (badan REPEATING `pgRepPgSubSectionSFAPortal_OpportunitiesListBBBBBBBB.
// pxResults`, `pyGridProps/pyRDName = GetListOpportunity`) bersarang di dua
// wadah bersyarat:
//
//	luar  pyContainerVisibleWhen = OperatorID.pyWorkGroup!='ReasLife' &&
//	      OperatorID.pyWorkBasketList(2).pyWorkBasketName=='ReasTreatyInAdmin'
//	dalam pyContainerVisibleWhen = !IsOperatorLife
//
// Di luar kedua wadah hanya baris saringan (`.FilterTermForOpportunity`,
// tombol Filter). Sec Head dan Dept Head TIDAK melihat grid itu; tugas
// mereka dirutekan `Flow/InputRealizationTreatyIn` ke workbasket
// (Assignment4/6 `ReasTreatyInSecHead`, Assignment3 `ReasTreatyInDeptHead`,
// `ToWorkBasket`) - padanannya: antrean posisi yang dipegang pelaku.
// Keanggotaan menurut NAMA workbasket, bukan nomor urut (AC 14).

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func (u *uji) daftar(p pelakuUji, kueri string) (int, []string, string) {
	u.t.Helper()
	kode, isi := u.panggil("GET", "/kasus"+kueri, p, nil)
	if kode != http.StatusOK {
		return kode, nil, isi
	}
	var baris []models.RingkasanKasus
	if err := json.Unmarshal([]byte(isi), &baris); err != nil {
		u.t.Fatalf("daftar bukan JSON: %v %s", err, isi)
	}
	id := []string{}
	for _, b := range baris {
		id = append(id, b.ID)
	}
	sort.Strings(id)
	return kode, id, isi
}

func TestGerbangDaftarPortal(t *testing.T) { // AC 11, 14, 92 - wadah grid SFAPortal_OpportunitiesList
	u := baru(t)
	diAdmin := u.buat()
	diSec := u.buat()
	if kode, isi := u.kirim(diSec, admin, halamanLengkap("1")); kode != http.StatusOK {
		t.Fatalf("admin menyetujui: %d %s", kode, isi)
	}
	diDept := u.buat()
	u.kirim(diDept, admin, halamanLengkap("1"))
	if kode, isi := u.kirim(diDept, secHead, putusan("1")); kode != http.StatusOK {
		t.Fatalf("Sec Head menyetujui: %d %s", kode, isi)
	}
	ditolak := u.buat() // admin menolak -> Resolved-Rejected (filter F GetListOpportunity)
	if kode, isi := u.kirim(ditolak, admin, halamanLengkap("0")); kode != http.StatusOK {
		t.Fatalf("admin menolak: %d %s", kode, isi)
	}
	for id, pos := range map[string]string{diAdmin: models.PosisiAdmin, diSec: models.PosisiSecHead, diDept: models.PosisiDeptHead} {
		if u.g.Kasus[id].PositionNote != pos {
			t.Fatalf("persiapan: %s di %q, harap %q", id, u.g.Kasus[id].PositionNote, pos)
		}
	}
	urut := func(id ...string) []string { sort.Strings(id); return id }
	semua := urut(diAdmin, diSec, diDept)

	for _, c := range []struct {
		nama  string
		p     pelakuUji
		kueri string
		kode  int
		harap []string
	}{
		// wadah luar: ReasTreatyInAdmin di urutan KEDUA daftar workbasket
		{"admin di urutan 2 melihat grid GetListOpportunity", admin, "", http.StatusOK, semua},
		// AC 14: menurut NAMA - urutan pertama pun anggota
		// RALAT 06-10-2026 (filter A GetListOpportunity, keputusan work owner): admin hanya melihat buatannya
		{"admin di urutan 1 tanpa berkas buatannya", pelakuUji{"UJI-ADMIN2", models.PosisiAdmin}, "", http.StatusOK, []string{}},
		{"admin menyaring satu posisi", admin, "?posisi=" + models.PosisiSecHead, http.StatusOK, urut(diSec)},
		// RALAT 06-10-2026 (keputusan work owner "hanya filter berdasarkan create operator aja"): antrean atasan
		// TIDAK lagi tampil di portal - setiap akun hanya melihat buatannya (berkas atasan dibuka dari Beranda)
		{"Sec Head tanpa berkas buatannya", secHead, "", http.StatusOK, []string{}},
		{"Dept Head tanpa berkas buatannya", deptHead, "", http.StatusOK, []string{}},
		{"Sec Head + Dept Head", pelakuUji{"UJI-SHDH", models.PosisiDeptHead + "," + models.PosisiSecHead}, "", http.StatusOK, []string{}},
		{"Sec Head mencari kasus admin", secHead, "?cari=" + diAdmin, http.StatusOK, []string{}},
		{"Sec Head meminta posisi admin", secHead, "?posisi=" + models.PosisiAdmin, http.StatusOK, []string{}},
		{"bukan anggota antrean tangga", orang, "", http.StatusOK, []string{}},
		{"GroupLeader (posisi buangan)", pelakuUji{"UJI-GL", "ReasTreatyInGroupLeader"}, "", http.StatusOK, []string{}},
		{"tanpa peran", pelakuUji{"UJI-KOSONG", ""}, "", http.StatusOK, []string{}},
		// switch Proses / Resolved (keputusan work owner 06-10-2026, bawaan Proses): yang ditolak hanya di Resolved
		{"admin: Resolved", admin, "?status=selesai", http.StatusOK, urut(ditolak)},
		// Resolved = SEMUA berkas selesai, siapa pun pembuatnya (WO 06-10-2026)
		{"Sec Head: Resolved semua", secHead, "?status=selesai", http.StatusOK, urut(ditolak)},
		{"admin lain: Resolved semua", pelakuUji{"UJI-ADMIN2", models.PosisiAdmin}, "?status=selesai", http.StatusOK, urut(ditolak)},
		{"tanpa peran: Resolved semua", pelakuUji{"UJI-KOSONG", ""}, "?status=selesai", http.StatusOK, urut(ditolak)},
		{"admin: Resolved disaring posisi", admin, "?status=selesai&posisi=" + models.PosisiSecHead, http.StatusOK, []string{}},
		{"admin: status tak dikenal = Proses", admin, "?status=UJI", http.StatusOK, semua},
	} {
		kode, id, isi := u.daftar(c.p, c.kueri)
		if kode != c.kode {
			t.Fatalf("%s: kode %d, harap %d (%s)", c.nama, kode, c.kode, isi)
		}
		if kode == http.StatusForbidden {
			if !strings.Contains(isi, "bukan anggota antrean") {
				t.Fatalf("%s: pesan 403 %s", c.nama, isi)
			}
			continue
		}
		if strings.Join(id, " ") != strings.Join(c.harap, " ") {
			t.Fatalf("%s: dapat %v, harap %v", c.nama, id, c.harap)
		}
	}
}

// RD `GetListOpportunity` (grid satu-satunya portal), `pyMaxRecords` = 500. Pencarian filter G
// `.TextNoQuotation Contains` DIPERLUAS (perintah work owner 07-10-2026: "pencarian nya pada nb dan edm treaty buat
// bisa mencari nomor nb/edm. insured name dll"): setiap kata wajib termuat di salah satu medan berkas
// (`models.CocokCariPortal`), tanpa beda huruf besar/kecil.
func TestDaftarPortalSesuaiGetListOpportunity(t *testing.T) { // P8, temuan tinjauan P9
	u := baru(t)
	id := u.buat()
	u.g.Halaman[id].Setel(models.HalamanQuotation+".BusinessName", "UJI-BISNIS-CARI")
	u.g.Halaman[id].Setel(models.HalamanQuotation+".InsuredName", "UJI-TERTANGGUNG-CARI")
	for _, c := range []struct {
		kueri string
		harap []string
	}{
		{"?cari=" + strings.ToLower(id), []string{id}},  // tanpa beda huruf besar/kecil
		{"?cari=uji-bisnis-cari", []string{id}},         // group business ikut dicari
		{"?cari=UJI-TERTANGGUNG-CARI", []string{id}},    // insured name ikut dicari
		{"?cari=tertanggung+" + id, []string{id}},       // dua kata, keduanya cocok (kolom berbeda)
		{"?cari=tertanggung+UJI-TIDAK-ADA", []string{}}, // satu kata tidak cocok = tidak tampil
	} {
		kode, ids, isi := u.daftar(admin, c.kueri)
		if kode != http.StatusOK || strings.Join(ids, " ") != strings.Join(c.harap, " ") {
			t.Fatalf("%s: %d %v (%s), harap %v", c.kueri, kode, ids, isi, c.harap)
		}
	}
	for i := 1; i <= models.BatasDaftarPortal; i++ { // 501 kasus terbuka
		u.buat()
	}
	if kode, ids, _ := u.daftar(admin, ""); kode != http.StatusOK || len(ids) != 500 {
		t.Fatalf("pyMaxRecords 500: %d baris (kode %d)", len(ids), kode)
	}
}
