package models

import (
	"strings"
	"testing"
	"time"
)

// TestIDDokumenPolis - `yyyyMMddhhmmssSSS` Asia/Jakarta: `hh` jam 12-an (13:04 WIB -> "01"), milidetik 3 digit.
func TestIDDokumenPolis(t *testing.T) {
	for saat, mau := range map[time.Time]string{
		time.Date(2026, 10, 8, 6, 4, 5, 7e6, time.UTC):    "20261008010405007",
		time.Date(2026, 10, 8, 18, 4, 5, 7e6, time.UTC):   "20261009010405007", // 01:04 WIB esok hari
		time.Date(2026, 10, 7, 17, 0, 0, 999e6, time.UTC): "20261008120000999", // 00:00 WIB = jam 12
	} {
		if got := IDDokumenPolis(saat); got != mau {
			t.Errorf("%v: %q mau %q", saat, got, mau)
		}
	}
}

// TestLokasiUnggah - Folder "Policy/Doc/YYYY/MM/" (YYYY = tahun-minggu en_US), Namafile "yyyyMMdd-hhmmss-S - <nama>" (S tanpa
// nol di depan) dari ID dokumen; ID berbentuk lain ditolak.
func TestLokasiUnggah(t *testing.T) {
	for id, mau := range map[string]LokasiObjek{
		"20261008010405007": {ImageID: "K", Folder: "Policy/Doc/2026/10/", Namafile: "20261008-010405-7 - a b.pdf"},
		"20261008010405120": {ImageID: "K", Folder: "Policy/Doc/2026/10/", Namafile: "20261008-010405-120 - a b.pdf"},
		"20261227090000000": {ImageID: "K", Folder: "Policy/Doc/2027/12/", Namafile: "20261227-090000-0 - a b.pdf"}, // Minggu: minggu 1 Jan 2027
		"20261226090000000": {ImageID: "K", Folder: "Policy/Doc/2026/12/", Namafile: "20261226-090000-0 - a b.pdf"}, // Sabtu
		"20270101090000000": {ImageID: "K", Folder: "Policy/Doc/2027/01/", Namafile: "20270101-090000-0 - a b.pdf"},
	} {
		got, ok := LokasiUnggah(id, "K", "a b.pdf")
		if !ok || got != mau {
			t.Errorf("%s: %+v mau %+v", id, got, mau)
		}
	}
	for _, id := range []string{"P1", "2026100801040500", "2026100801040500x"} {
		if _, ok := LokasiUnggah(id, "K", "a.pdf"); ok {
			t.Errorf("%q diterima", id)
		}
	}
}

// TestLokasiDariLink - GetUrlGoogleStorage_Act 6.2 / DeleteGoogleStorage_Act 6 atas APPFOLDER tersimpan.
func TestLokasiDariLink(t *testing.T) {
	l := LinkStorage{ImageID: "K", App: "UJIAPP", Namafile: "20261008-010405-7 - a.pdf",
		AppFolder: AppFolderUji("UJIAPP", "Policy/Doc/2026/10/20261008-010405-7 - a.pdf")}
	if !strings.HasPrefix(l.AppFolder, "gs:") || !strings.Contains(l.AppFolder, "/UJIAPP/Policy/") {
		t.Fatalf("instrumen APPFOLDER %q", l.AppFolder)
	}
	if got := LokasiDariLink(l); got != (LokasiObjek{ImageID: "K", Folder: "Policy/Doc/2026/10/", Namafile: l.Namafile}) {
		t.Errorf("%+v", got)
	}
	if got := JalurHapusDariLink(l); got != "Policy/Doc/2026/10/20261008-010405-7 - a.pdf" {
		t.Errorf("hapus %q", got)
	}
	if got := JalurHapusDariLink(LinkStorage{AppFolder: "Policy/x.pdf"}); got != "Policy/x.pdf" {
		t.Errorf("App kosong %q", got)
	}
}

func TestStatusDokumen(t *testing.T) {
	for _, c := range []struct {
		d   DokumenPolis
		mau string
	}{
		{DokumenPolis{Tersimpan: true, StatusEfek: "GAGAL_PERMANEN"}, ""},
		{DokumenPolis{StatusEfek: "GAGAL_PERMANEN"}, StatusLampiranGagal},
		{DokumenPolis{}, StatusLampiranGagal},
		{DokumenPolis{StatusEfek: "ANTRE"}, StatusLampiranMenunggu},
		{DokumenPolis{StatusEfek: "SELESAI"}, StatusLampiranMenunggu},
	} {
		if got := StatusDokumen(c.d, "GAGAL_PERMANEN"); got != c.mau {
			t.Errorf("%+v: %q mau %q", c.d, got, c.mau)
		}
	}
}

func TestBentukStorage(t *testing.T) {
	if got := ExpStorage("2026-10-08T13:34:05.000Z"); got != "08/10/2026 13:34:05" {
		t.Errorf("exp %q", got)
	}
	if ExpStorage("besok") != "" || TanggalUploadStorage("2026-10-08") != "" || TanggalUploadStorage(" 10/08/2026 13:00:00 ") != "10/08/2026 13:00:00" {
		t.Error("bentuk tak terbaca")
	}
	kini := time.Date(2026, 10, 8, 6, 0, 0, 0, time.UTC)
	l := LinkStorage{URLPublic: "x", ExpUTC: "20261008060200"}
	if !LinkBerlaku(l, kini, time.Minute) || LinkBerlaku(l, kini, 3*time.Minute) {
		t.Error("berlaku")
	}
	l.URLPublic = ""
	if LinkBerlaku(l, kini, 0) || LinkBerlaku(LinkStorage{URLPublic: "x"}, kini, 0) {
		t.Error("tanpa URL / EXPDATE")
	}
	if NamaBerkasAntreLampiran("K1", `..\..\x.PDF`) != "K1.pdf" || NamaBerkasAntreLampiran("K1", "x.pdf;rm") != "K1" ||
		EkstensiBerkas("laporan.2026.XLSX") != "xlsx" || EkstensiBerkas("README") != "" || NamaDasarBerkas("..") != "" {
		t.Error("nama berkas")
	}
	if k, ok := KategoriLampiranSah([]string{"EXTEND WPC", "OFFER"}, " offer "); !ok || k != "OFFER" {
		t.Errorf("kategori %q", k)
	}
	if _, ok := KategoriLampiranSah([]string{"OFFER"}, ""); ok {
		t.Error("kategori kosong")
	}
}
