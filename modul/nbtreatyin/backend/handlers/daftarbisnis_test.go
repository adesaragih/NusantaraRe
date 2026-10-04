package handlers_test

// Uji seam 1 isi popup pilih bisnis - audit silang putaran 3 W1, IKUTI XML:
//
//	tombol `Choose Business` (`Section/DetailPolicyTreatyIn`, wadah
//	  `.ClaimType != 'XOL Retro'`) -> showHarness `BusinessAndSOBList`
//	  (`pyWindowName` "Business And SOB List", `pySubmitData=Yes`)
//	`Section/BusinessAndSOBList`: grid lama RD `BrowseTreatyInDetail` (tabel
//	  TREATYINDETAIL) ber-`pyContainerVisibleWhen=1=2` - memo rule "Hidden the
//	  old one, now use treatyindetail join edm"; grid AKTIF `pyGridProps/pyRDName`
//	  `BrowseTreatyJoinEDM` (view TREATYINDETAILJOINEDM) dengan `pyRDParams`
//	  `PROPORTIONALTYPE = .QuotationData.ProportionalType`
//	RD `BrowseTreatyJoinEDM` filter H `.PROPORTIONTYPE = Param.PROPORTIONALTYPE`
//	  tanpa `pyUseNullIfEmpty` (= false): parameter kosong -> syarat diabaikan
//
// Helper (`baru`, `halamanLengkap`, pelaku) milik alur_test.go.

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

// kontrakPopup - empat baris view: dua Proportional, satu NonProportional, satu
// tanpa jenis proporsi.
func kontrakPopup(u *uji) {
	u.g.Kontrak["UJI-D-P1"] = models.BarisKontrak{"ID": "UJI-D-P1", "TREATYID": "UJI-T2", "PROPORTIONTYPE": "Proportional"}
	u.g.Kontrak["UJI-D-P2"] = models.BarisKontrak{"ID": "UJI-D-P2", "TREATYID": "UJI-T1", "PROPORTIONTYPE": "Proportional"}
	u.g.Kontrak["UJI-D-N1"] = models.BarisKontrak{"ID": "UJI-D-N1", "TREATYID": "UJI-T3", "PROPORTIONTYPE": "NonProportional"}
	u.g.Kontrak["UJI-D-K"] = models.BarisKontrak{"ID": "UJI-D-K", "TREATYID": "UJI-T4", "PROPORTIONTYPE": ""}
}

func (u *uji) daftarBisnis(id string, p pelakuUji, h *models.Halaman) (int, string) {
	u.t.Helper()
	return u.panggil("POST", "/kasus/"+id+"/bisnis", p, map[string]any{"halaman": h})
}

// idBisnis - ID baris jawaban popup, berurut (urutan jawaban = urutan RD).
func (u *uji) idBisnis(id string, h *models.Halaman) []string {
	u.t.Helper()
	kode, isi := u.daftarBisnis(id, admin, h)
	if kode != http.StatusOK {
		u.t.Fatalf("daftar bisnis: %d %s", kode, isi)
	}
	var b []models.BarisKontrak
	if err := json.Unmarshal([]byte(isi), &b); err != nil {
		u.t.Fatalf("%v: %s", err, isi)
	}
	out := []string{}
	for _, x := range b {
		out = append(out, x["ID"])
	}
	slices.Sort(out)
	return out
}

func halamanProporsi(jenis string) *models.Halaman {
	h := models.HalamanBaru()
	h.Setel("PolicyTreatyIn.QuotationData.ProportionalType", jenis)
	return h
}

// Filter H atas `PolicyTreatyIn.QuotationData.ProportionalType` halaman kasus
// saat tombol diklik (showHarness `pySubmitData=Yes`). Sel medan itu
// `pyReadOnly=true`: nilainya salinan `Quotation` (InputPolicyTreatyIn_preDT 14)
// atau hasil tombol Enable / Disable Input Type yang dipegang layar
// (`TreatyEnableDisableInput`, hanya bila TreatyType XOL - W3 audit silang P3).
func TestDaftarBisnisMenyaringJenisProporsiKasus(t *testing.T) {
	u := baru(t)
	kontrakPopup(u)
	id := u.buat()
	u.g.Halaman[id].Setel("Quotation.ProportionalType", "Proportional")
	if got := u.idBisnis(id, nil); !slices.Equal(got, []string{"UJI-D-P1", "UJI-D-P2"}) {
		t.Errorf("Proportional: %v", got)
	}
	u.g.Halaman[id].Setel("PolicyTreatyIn.TreatyType", "XOL")
	if got := u.idBisnis(id, halamanProporsi("NonProportional")); !slices.Equal(got, []string{"UJI-D-N1"}) {
		t.Errorf("NonProportional (hasil Enable / Disable Input Type): %v", got)
	}
}

// Kasus baru: `.QuotationData.ProportionalType` kosong -> filter H diabaikan
// (RD tanpa `pyUseNullIfEmpty`), seluruh baris view tampil - termasuk baris
// tanpa jenis proporsi.
func TestDaftarBisnisJenisProporsiKosongMengabaikanFilter(t *testing.T) {
	u := baru(t)
	kontrakPopup(u)
	id := u.buat()
	harap := []string{"UJI-D-K", "UJI-D-N1", "UJI-D-P1", "UJI-D-P2"}
	if got := u.idBisnis(id, nil); !slices.Equal(got, harap) {
		t.Errorf("tanpa kiriman layar: %v", got)
	}
	if got := u.idBisnis(id, halamanProporsi("")); !slices.Equal(got, harap) {
		t.Errorf("isian kosong: %v", got)
	}
}

// Nilai TERSIMPAN dipakai bila layar tidak mengirim medan itu: pilih bisnis
// menyimpan `Quotation.ProportionalType` (preACT langkah 11, satu kolom
// T_POLIS_QUOTATION), dan pra-proses `InputPolicyTreatyIn_preDT` langkah 14
// menyalinnya ke `PolicyTreatyIn.QuotationData` sebelum grid membacanya.
func TestDaftarBisnisMemakaiJenisProporsiTersimpan(t *testing.T) {
	u := baru(t)
	kontrakPopup(u)
	id := u.buat()
	u.g.Halaman[id].Setel("Quotation.ProportionalType", "NonProportional")
	if got := u.idBisnis(id, nil); !slices.Equal(got, []string{"UJI-D-N1"}) {
		t.Errorf("tersimpan NonProportional: %v", got)
	}
}

// Hak dan posisi sama dengan tindakan layar lain: pelaku wajib anggota antrean
// posisi kasus; tombol `Choose Business` hanya di layar admin
// (`Section/DetailPolicyTreatyIn`, flow action `InboxPolicyTreatyIn`) dan
// hanya bila wadahnya tampil (`.ClaimType != 'XOL Retro'`, isian layar).
func TestDaftarBisnisHakDanPosisi(t *testing.T) {
	u := baru(t)
	kontrakPopup(u)
	id := u.buat()
	if kode, _ := u.daftarBisnis(id, pelakuUji{}, nil); kode != http.StatusUnauthorized {
		t.Errorf("tanpa identitas: %d", kode)
	}
	if kode, _ := u.daftarBisnis(id, orang, nil); kode != http.StatusForbidden {
		t.Errorf("bukan anggota antrean admin: %d", kode)
	}
	if kode, _ := u.daftarBisnis("UJI-TIDAK-ADA", admin, nil); kode != http.StatusNotFound {
		t.Errorf("kasus tak ada: %d", kode)
	}
	retro := models.HalamanBaru()
	retro.Setel("PolicyTreatyIn.ClaimType", "XOL Retro")
	if kode, isi := u.daftarBisnis(id, admin, retro); kode != http.StatusConflict {
		t.Errorf("ClaimType XOL Retro (tombol tak tampil): %d %s", kode, isi)
	}
	if kode, isi := u.kirim(id, admin, halamanLengkap("1")); kode != http.StatusOK {
		t.Fatalf("kirim ke Sec Head: %d %s", kode, isi)
	}
	if kode, isi := u.daftarBisnis(id, secHead, nil); kode != http.StatusConflict {
		t.Errorf("layar Sec Head tanpa tombol Choose Business: %d %s", kode, isi)
	}
}

// showHarness tidak ber-Obj-Save: membuka popup tidak menyimpan apa pun,
// termasuk isian layar yang ikut dikirim (`pySubmitData=Yes` hanya ke
// clipboard).
func TestDaftarBisnisTidakMenyimpan(t *testing.T) {
	u := baru(t)
	kontrakPopup(u)
	id := u.buat()
	sebelum := len(u.g.Panggil)
	h := halamanProporsi("Proportional")
	h.Setel("PolicyTreatyIn.Remark", "UJI-belum-disimpan")
	if kode, isi := u.daftarBisnis(id, admin, h); kode != http.StatusOK {
		t.Fatalf("%d %s", kode, isi)
	}
	if slices.Contains(u.g.Panggil[sebelum:], "SimpanHalaman") {
		t.Fatal("membuka popup tidak boleh menyimpan halaman")
	}
	if got := u.g.Halaman[id].Ambil("PolicyTreatyIn.Remark"); got != "" {
		t.Fatalf("isian layar tersimpan oleh popup: %q", got)
	}
}

// Grid lama (RD `BrowseTreatyInDetail`, `GET /bisnis` tanpa kasus) tidak ada lagi.
func TestRuteDaftarBisnisLamaTidakAda(t *testing.T) {
	u := baru(t)
	if kode, _ := u.panggil("GET", "/bisnis", admin, nil); kode != http.StatusNotFound {
		t.Fatalf("GET /bisnis: %d", kode)
	}
}

// Tombol `Choose` hidup HANYA di dalam popup: bila tombol `Choose Business`
// tidak tampil (`.ClaimType = 'XOL Retro'`), pilih bisnis ditolak dan nol
// simpanan - sama dengan isi popup.
func TestPilihBisnisDitolakBilaTombolTakTampil(t *testing.T) {
	u := baru(t)
	kontrakPopup(u)
	id := u.buat()
	retro := models.HalamanBaru()
	retro.Setel("PolicyTreatyIn.ClaimType", "XOL Retro")
	kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D-P1", "halaman": retro})
	if kode != http.StatusConflict {
		t.Fatalf("ClaimType XOL Retro: %d %s", kode, isi)
	}
	if got := u.g.Halaman[id].Ambil("PolicyTreatyIn.NoOffer"); got != "" {
		t.Fatalf("pilih bisnis yang ditolak tidak boleh menyimpan: NoOffer %q", got)
	}
}

// Pola F4 (pengerasan tercatat, PERMINTAAN H2): tombol `Choose` hanya ada di
// baris grid AKTIF popup (`SetValue_Act(ID=.ID)`), jadi `idDetail` diterima hanya
// bila ada di RD `BrowseTreatyJoinEDM` yang dijalankan ulang dengan saringan
// kasus ini. UJI-D-N1 ada di view, tetapi tersaring keluar (filter H
// PROPORTIONTYPE = Proportional) - 422 dan nol simpanan; UJI-D-P1 di daftar - 200.
func TestPilihBisnisDiLuarDaftarPopupDitolak(t *testing.T) {
	u := baru(t)
	kontrakPopup(u)
	id := u.buat()
	u.g.Halaman[id].Setel("Quotation.ProportionalType", "Proportional")
	sebelum := len(u.g.Panggil)
	kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D-N1"})
	if kode != http.StatusUnprocessableEntity || !strings.Contains(isi, "tidak ada di daftar popup Choose Business") {
		t.Fatalf("ID di luar daftar popup tersaring: %d %s", kode, isi)
	}
	if slices.Contains(u.g.Panggil[sebelum:], "SimpanHalaman") || u.g.Halaman[id].Ambil("PolicyTreatyIn.NoOffer") != "" {
		t.Fatal("pilihan yang ditolak tidak boleh menyimpan")
	}
	if kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D-P1"}); kode != http.StatusOK {
		t.Fatalf("ID di daftar popup: %d %s", kode, isi)
	}
}
