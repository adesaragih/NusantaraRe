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
		{"admin di urutan 1", pelakuUji{"UJI-ADMIN2", models.PosisiAdmin}, "", http.StatusOK, semua},
		{"admin menyaring satu posisi", admin, "?posisi=" + models.PosisiSecHead, http.StatusOK, urut(diSec)},
		// di luar wadah: antrean workbasket posisi yang dipegang (Flow ToWorkBasket)
		{"Sec Head hanya antreannya", secHead, "", http.StatusOK, urut(diSec)},
		{"Dept Head hanya antreannya", deptHead, "", http.StatusOK, urut(diDept)},
		{"Sec Head + Dept Head", pelakuUji{"UJI-SHDH", models.PosisiDeptHead + "," + models.PosisiSecHead}, "", http.StatusOK, urut(diSec, diDept)},
		{"Sec Head mencari kasus admin", secHead, "?cari=" + diAdmin, http.StatusOK, []string{}},
		{"Sec Head meminta posisi admin", secHead, "?posisi=" + models.PosisiAdmin, http.StatusForbidden, nil},
		{"bukan anggota antrean tangga", orang, "", http.StatusForbidden, nil},
		// AC 10: posisi buangan tidak memberi akses
		{"GroupLeader (posisi buangan)", pelakuUji{"UJI-GL", "ReasTreatyInGroupLeader"}, "", http.StatusForbidden, nil},
		{"tanpa peran", pelakuUji{"UJI-KOSONG", ""}, "", http.StatusForbidden, nil},
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

// RD `GetListOpportunity` (grid satu-satunya portal): pencarian hanya filter G
// `.TextNoQuotation Contains Param.Search` (`pyCaseInsensitive=true`) = pengenal
// kasus; filter C `.Name` (kelas CRM, ditulis nol rule) tidak dibangun - nama
// bisnis/tertanggung BUKAN medan pencarian. `pyMaxRecords` = 500.
func TestDaftarPortalSesuaiGetListOpportunity(t *testing.T) { // P8, temuan tinjauan P9
	u := baru(t)
	id := u.buat()
	u.g.Halaman[id].Setel(models.HalamanQuotation+".BusinessName", "UJI-BISNIS-CARI")
	u.g.Halaman[id].Setel(models.HalamanQuotation+".InsuredName", "UJI-TERTANGGUNG-CARI")
	for _, c := range []struct {
		kueri string
		harap []string
	}{
		{"?cari=" + strings.ToLower(id), []string{id}}, // tanpa beda huruf besar/kecil
		{"?cari=UJI-BISNIS-CARI", []string{}},          // .Name / nama bisnis bukan saringan RD
		{"?cari=UJI-TERTANGGUNG-CARI", []string{}},
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
