package services_test

// Uji tombol tulis Save / Submit / Actions / Decline offer — aturan ekspor
// dan keputusan pemilik proses 6–7 Oktober 2026, di atas gudang tiruan.

import (
	"context"
	"errors"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

func (g *gudangTiruan) BacaKepalaTreatyIn(_ context.Context, id string) (map[string]any, bool, error) {
	k, ada := g.kepalaTreatyIn[id]
	salin := map[string]any{}
	for a, b := range k {
		salin[a] = b
	}
	return salin, ada, nil
}

func (g *gudangTiruan) BacaDokumenPendaratan(_ context.Context, id string) (map[string]any, error) {
	d := map[string]any{}
	for a, b := range g.dokumenTersimpan[id] {
		d[a] = b
	}
	return d, nil
}

func (g *gudangTiruan) SimpanKontrak(_ context.Context, r models.RencanaSimpan) (string, error) {
	g.disimpan = append(g.disimpan, r)
	if r.ID == "" {
		return "1001857", nil
	}
	return r.ID, nil
}

func (g *gudangTiruan) KunciBelumTerpasang(_ context.Context, doc map[string]any) ([]string, error) {
	var out []string
	for _, k := range g.belumTerpasang {
		if _, ada := doc[k]; ada {
			out = append(out, k)
		}
	}
	return out, nil
}

func (g *gudangTiruan) PemegangPosisi(_ context.Context, wb string) ([]string, error) {
	return g.pemegangPosisi[wb], nil
}

var (
	admin    = inti.Pelaku{AkunID: "ADESAMUEL", Peran: []string{models.PosisiAdmin}}
	secHead  = inti.Pelaku{AkunID: "SEC1", Peran: []string{models.PosisiSecHead}}
	director = inti.Pelaku{AkunID: "DIR1", Peran: []string{models.PosisiDirector}}
)

func gudangSimpan() *gudangTiruan {
	return &gudangTiruan{
		kepalaTreatyIn: map[string]map[string]any{
			"1001001": {"TreatyContractName": "LAMA", "Ceding": "ASURANSI A", "LeadingReinsSource": "DIRECT", "ClassofBusiness": "FIRE"},
		},
		dokumenTersimpan: map[string]map[string]any{
			"1001001": {"Portfolio": []any{map[string]any{"Description": "tersimpan"}}, "StatusAkseptasi": "", "Position": ""},
		},
		pemegangPosisi: map[string][]string{
			models.PosisiSecHead:  {"SEC2", "SEC1"},
			models.PosisiDeptHead: {"DEPT1"},
		},
	}
}

func TestSaveMenimpaClipboardTanpaMenghapusTabYangTakDibuka(t *testing.T) {
	g := gudangSimpan()
	l := services.LayananDengan(g)
	h, err := l.SimpanKontrak(context.Background(), admin, services.MasukanSimpan{
		IDKontrak: "1001001",
		Dokumen: map[string]any{
			"TreatyContractName": "BARU",
			// ⛔ Milik server — diabaikan.
			"StatusAkseptasi": "Resolve Complete", "Position": "ReasTreatyInDirector",
			// Halaman sesi — bukan properti TreatyIn.
			"SearchData.CARI1": "Q1",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	d := g.disimpan[0].Dokumen
	if d["TreatyContractName"] != "BARU" || d["ClassofBusiness"] != "FIRE" {
		t.Errorf("kepala %v", d)
	}
	// Tab Portfolio tidak dikirim layar — tetap dari dokumen tersimpan.
	if p, _ := d["Portfolio"].([]any); len(p) != 1 {
		t.Errorf("Portfolio hilang: %v", d["Portfolio"])
	}
	// DT TreatyInAddNew: Position = ReasTreatyInAdmin; status tidak berubah.
	if d["Position"] != models.PosisiAdmin || d["StatusAkseptasi"] != "" {
		t.Errorf("Position %v Status %v", d["Position"], d["StatusAkseptasi"])
	}
	if _, ada := d["SearchData.CARI1"]; ada {
		t.Error("halaman sesi ikut tersimpan")
	}
	if h.Pesan != "Data Sudah Disimpan Dengan ID : 1001001" {
		t.Errorf("pesan %q", h.Pesan)
	}
}

func TestSaveKontrakBaruMendapatPengenalDariGudang(t *testing.T) {
	g := gudangSimpan()
	h, err := services.LayananDengan(g).SimpanKontrak(context.Background(), admin, services.MasukanSimpan{
		Dokumen: map[string]any{"TreatyContractName": "KONTRAK BARU", "RNMShareP": "10", "PropertiTanpaKolom": "x"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if g.disimpan[0].ID != "" || h.ID != "1001857" {
		t.Errorf("rencana ID %q, hasil %q", g.disimpan[0].ID, h.ID)
	}
	// ⛔ Properti tanpa kolom DILAPORKAN, tidak ditelan.
	// `RNMShareP` kini punya kolom (migrasi 448) — tidak lagi dilaporkan.
	if got := strings.Join(h.KunciTakTersimpan, ","); !strings.Contains(got, "PropertiTanpaKolom") || strings.Contains(got, "RNMShareP") {
		t.Errorf("kunci tak tersimpan %v", h.KunciTakTersimpan)
	}
}

func TestSaveDitolakSesudahResolveComplete(t *testing.T) {
	g := gudangSimpan()
	g.dokumenTersimpan["1001001"]["StatusAkseptasi"] = "Resolve Complete"
	_, err := services.LayananDengan(g).SimpanKontrak(context.Background(), admin, services.MasukanSimpan{IDKontrak: "1001001"})
	if !errors.Is(err, services.ErrTombolDitolak) {
		t.Fatalf("galat %v", err)
	}
	if len(g.disimpan) != 0 {
		t.Error("tetap tersimpan")
	}
}

func TestSubmitMemeriksaCedingLaluSoB(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		dok  map[string]any
		mau  string
		alas string
	}{
		{map[string]any{"Ceding": "", "LeadingReinsSource": ""}, "Please input Source of Business (SoB)", "keduanya kosong — langkah TERAKHIR menang"},
		{map[string]any{"Ceding": "", "LeadingReinsSource": "DIRECT"}, "Please input Ceding", "Ceding kosong"},
	} {
		g := gudangSimpan()
		_, err := services.LayananDengan(g).KirimKontrak(ctx, admin, services.MasukanKirim{
			MasukanSimpan: services.MasukanSimpan{IDKontrak: "1001001", Dokumen: c.dok}, Aksi: services.AksiSubmit,
		})
		if !errors.Is(err, services.ErrTombolDitolak) || err.Error() != c.mau {
			t.Errorf("%s: galat %v, mau %q", c.alas, err, c.mau)
		}
	}
}

func TestSubmitAdminNaikKeSecHead(t *testing.T) {
	g := gudangSimpan()
	h, err := services.LayananDengan(g).KirimKontrak(context.Background(), admin, services.MasukanKirim{
		MasukanSimpan: services.MasukanSimpan{IDKontrak: "1001001", Dokumen: map[string]any{"Comment": "mohon diperiksa"}},
		Aksi:          services.AksiSubmit,
	})
	if err != nil {
		t.Fatal(err)
	}
	d := g.disimpan[0].Dokumen
	if d["Position"] != models.PosisiSecHead || d["StatusAkseptasi"] != "Accept" || d["ChooseStatusAkseptasi"] != "Accept" {
		t.Errorf("langkah %v %v %v", d["Position"], d["StatusAkseptasi"], d["ChooseStatusAkseptasi"])
	}
	// ⛔ SATU nama, yang PERTAMA SECARA URUT — permintaan pemilik proses
	// 8 Oktober 2026. Gudang tiruan mengembalikan `{"SEC2", "SEC1"}` dengan
	// sengaja: urutan basis data BUKAN urutan hasilnya.
	if d["PositionUsername"] != "SEC1" || h.PemegangPosisi != "SEC1" {
		t.Errorf("PositionUsername %v", d["PositionUsername"])
	}
	// AddCommentList_Act — riwayat membaca status BARU.
	k, _ := d["CommentList"].([]any)
	el, _ := k[len(k)-1].(map[string]any)
	if el["OperatorName"] != "ADESAMUEL" || el["IsApproved"] != "Accept" || el["Suggest"] != "mohon diperiksa" || el["Date"] == "" {
		t.Errorf("komentar %v", el)
	}
}

func TestSubmitDitolakBagiYangBukanPemegangPosisi(t *testing.T) {
	g := gudangSimpan()
	_, err := services.LayananDengan(g).KirimKontrak(context.Background(), secHead, services.MasukanKirim{
		MasukanSimpan: services.MasukanSimpan{IDKontrak: "1001001"}, Aksi: services.AksiSubmit,
	})
	if !errors.Is(err, services.ErrBukanPemegangPosisi) || len(g.disimpan) != 0 {
		t.Errorf("galat %v, disimpan %d", err, len(g.disimpan))
	}
}

func TestActionsRejectKembaliKePengaju(t *testing.T) {
	g := gudangSimpan()
	g.dokumenTersimpan["1001001"]["Position"] = models.PosisiSecHead
	g.dokumenTersimpan["1001001"]["StatusAkseptasi"] = "Accept"
	g.dokumenTersimpan["1001001"]["CommentList"] = []any{
		map[string]any{"OperatorName": "ADESAMUEL", "IsApproved": "Accept"},
		map[string]any{"OperatorName": "LAIN", "IsApproved": "Accept"},
	}
	_, err := services.LayananDengan(g).KirimKontrak(context.Background(), secHead, services.MasukanKirim{
		MasukanSimpan: services.MasukanSimpan{IDKontrak: "1001001", Dokumen: map[string]any{"Comment": "lengkapi"}},
		Aksi:          services.AksiAkseptasi, Pilihan: "Reject",
	})
	if err != nil {
		t.Fatal(err)
	}
	d := g.disimpan[0].Dokumen
	// Reject: `CommentList(1).OperatorName` — penulis komentar PERTAMA.
	if d["Position"] != models.PosisiAdmin || d["StatusAkseptasi"] != "Reject" || d["PositionUsername"] != "ADESAMUEL" {
		t.Errorf("Reject %v %v %v", d["Position"], d["StatusAkseptasi"], d["PositionUsername"])
	}
}

func TestActionsDirectorAcceptMenuntaskan(t *testing.T) {
	g := gudangSimpan()
	g.dokumenTersimpan["1001001"]["Position"] = models.PosisiDirector
	g.dokumenTersimpan["1001001"]["StatusAkseptasi"] = "Accept"
	_, err := services.LayananDengan(g).KirimKontrak(context.Background(), director, services.MasukanKirim{
		MasukanSimpan: services.MasukanSimpan{IDKontrak: "1001001"}, Aksi: services.AksiAkseptasi, Pilihan: "Accept",
	})
	if err != nil {
		t.Fatal(err)
	}
	d := g.disimpan[0].Dokumen
	if d["Position"] != "" || d["StatusAkseptasi"] != models.StatusTuntas || d["PositionUsername"] != "" {
		t.Errorf("tuntas %v %v %v", d["Position"], d["StatusAkseptasi"], d["PositionUsername"])
	}
}

func TestActionsPilihanTakBerlakuDariAdmin(t *testing.T) {
	g := gudangSimpan()
	_, err := services.LayananDengan(g).KirimKontrak(context.Background(), admin, services.MasukanKirim{
		MasukanSimpan: services.MasukanSimpan{IDKontrak: "1001001"}, Aksi: services.AksiAkseptasi, Pilihan: "Reject",
	})
	if !errors.Is(err, services.ErrTombolDitolak) {
		t.Errorf("galat %v", err)
	}
}

func TestDeclineOfferMencatatDanMenandaiDecline(t *testing.T) {
	g := gudangSimpan()
	_, err := services.LayananDengan(g).KirimKontrak(context.Background(), admin, services.MasukanKirim{
		MasukanSimpan: services.MasukanSimpan{IDKontrak: "1001001", Dokumen: map[string]any{"Comment": "ditolak cedant"}},
		Aksi:          services.AksiDecline,
	})
	if err != nil {
		t.Fatal(err)
	}
	d := g.disimpan[0].Dokumen
	k, _ := d["CommentList"].([]any)
	el, _ := k[len(k)-1].(map[string]any)
	if d["StatusAkseptasi"] != "Decline" || el["IsApproved"] != "Decline" || el["Suggest"] != "ditolak cedant" {
		t.Errorf("decline %v %v", d["StatusAkseptasi"], el)
	}
	// Posisi tidak berubah — `TreatyInDeclineConfirmation_postact` tidak menyentuhnya.
	if d["Position"] != "" {
		t.Errorf("Position %v", d["Position"])
	}
}

// ⭐ 8 Oktober 2026 — laporan "TIDAK tersimpan" hanya atas properti BERISI:
// larik kosong yang pembaca pohon sisipkan bukan data yang hilang.
func TestLaporanTakTersimpanMengabaikanLarikKosong(t *testing.T) {
	g := gudangSimpan()
	h, err := services.LayananDengan(g).SimpanKontrak(context.Background(), admin, services.MasukanSimpan{
		IDKontrak: "1001001",
		Dokumen: map[string]any{
			"LarikAsingKosong": []any{},
			"TeksAsingKosong":  "",
			"LarikAsingBerisi": []any{map[string]any{"A": "1"}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	laporan := strings.Join(h.KunciTakTersimpan, ",")
	if strings.Contains(laporan, "LarikAsingKosong") || strings.Contains(laporan, "TeksAsingKosong") {
		t.Errorf("properti kosong dilaporkan: %v", h.KunciTakTersimpan)
	}
	if !strings.Contains(laporan, "LarikAsingBerisi") {
		t.Errorf("properti berisi TIDAK dilaporkan: %v", h.KunciTakTersimpan)
	}
	// Dokumen yang ditulis tetap utuh — penyaring hanya untuk laporan.
	if _, ada := g.disimpan[0].Dokumen["LarikAsingKosong"]; !ada {
		t.Error("dokumen yang ditulis ikut tersaring")
	}
}

// ⛔ TOMBOL SAVE TIDAK LAGI MENJATUHKAN BERKAS KEMBALI KE PENGAJU.
//
// ---------------------------------------------------------------------
// Laporan pemakai 8 Oktober 2026 atas kontrak uji `1002305`
// ---------------------------------------------------------------------
// *"di submit terus di actions sudah di accept terus tetapi belum pernah
// resolve complete"*.
//
// Terukur di basis data: `POSITION = ReasTreatyInDeptHead` dengan TIGA
// komentar akseptasi. Tanpa Save, tiga akseptasi berakhir di
// `ReasTreatyInDirector` — satu-satunya urutan yang menghasilkan Dept Head
// adalah Submit → SAVE → Accept → Accept.
//
// ⚠️ Uji ini memakai urutan ITU, bukan urutan yang disederhanakan: yang
// hendak dijaga adalah bahwa Save DI TENGAH tangga tidak lagi memundurkan.
func TestSaveTidakMemundurkanTanggaAkseptasi(t *testing.T) {
	g := gudangSimpan()
	g.pemegangPosisi[models.PosisiDeptHead] = []string{"DEPT1"}
	g.pemegangPosisi[models.PosisiDirector] = []string{"DIR1"}
	l := services.LayananDengan(g)
	ctx := context.Background()
	// Satu akun memegang seluruh anak tangga — seperti akun uji pemakai.
	semua := inti.Pelaku{AkunID: "DASTIN", Peran: []string{
		models.PosisiAdmin, models.PosisiSecHead, models.PosisiDeptHead, models.PosisiDirector,
	}}
	// Umpan balik ke "basis data" tiruan, seperti jalur nyata membacanya ulang.
	simpanBalik := func() map[string]any {
		d := g.disimpan[len(g.disimpan)-1].Dokumen
		g.dokumenTersimpan["1001001"] = d
		return d
	}
	kirim := func(aksi, pilihan string) map[string]any {
		t.Helper()
		if _, err := l.KirimKontrak(ctx, semua, services.MasukanKirim{
			MasukanSimpan: services.MasukanSimpan{IDKontrak: "1001001"}, Aksi: aksi, Pilihan: pilihan,
		}); err != nil {
			t.Fatalf("%s: %v", aksi, err)
		}
		return simpanBalik()
	}

	if d := kirim(services.AksiSubmit, ""); d["Position"] != models.PosisiSecHead {
		t.Fatalf("sesudah Submit posisi %v, mau SecHead", d["Position"])
	}
	// ⛔ INILAH LANGKAH YANG DULU MERUSAK.
	if _, err := l.SimpanKontrak(ctx, semua, services.MasukanSimpan{IDKontrak: "1001001"}); err != nil {
		t.Fatal(err)
	}
	if d := simpanBalik(); d["Position"] != models.PosisiSecHead {
		t.Fatalf("Save memundurkan posisi ke %v — seharusnya tetap SecHead", d["Position"])
	}
	if d := kirim(services.AksiAkseptasi, models.PilihAccept); d["Position"] != models.PosisiDeptHead {
		t.Fatalf("posisi %v, mau DeptHead", d["Position"])
	}
	if d := kirim(services.AksiAkseptasi, models.PilihAccept); d["Position"] != models.PosisiDirector {
		t.Fatalf("posisi %v, mau Director", d["Position"])
	}
	d := kirim(services.AksiAkseptasi, models.PilihAccept)
	if d["StatusAkseptasi"] != models.StatusTuntas {
		t.Fatalf("status %v, mau %q", d["StatusAkseptasi"], models.StatusTuntas)
	}
	if d["Position"] != models.PosisiKosong {
		t.Errorf("posisi %v, mau kosong saat tuntas", d["Position"])
	}
	// ⭐ Dan Save sesudah tuntas tetap DITOLAK — pagar lama tidak ikut longgar.
	if _, err := l.SimpanKontrak(ctx, semua, services.MasukanSimpan{IDKontrak: "1001001"}); err == nil {
		t.Error("Save sesudah Resolve Complete tidak ditolak")
	}
}

// ⭐ PASANGANNYA: di posisi Admin (atau belum berposisi) Save TETAP
// menyetel Admin — itu langkah `[1]` DT `TreatyInAddNew` yang sah, dan
// mencabutnya akan membuat kontrak baru nol berposisi.
func TestSaveTetapMenyetelAdminDiPangkalTangga(t *testing.T) {
	for _, awal := range []string{"", models.PosisiAdmin} {
		g := gudangSimpan()
		g.dokumenTersimpan["1001001"]["Position"] = awal
		if _, err := services.LayananDengan(g).SimpanKontrak(
			context.Background(), admin, services.MasukanSimpan{IDKontrak: "1001001"}); err != nil {
			t.Fatal(err)
		}
		if got := g.disimpan[0].Dokumen["Position"]; got != models.PosisiAdmin {
			t.Errorf("posisi awal %q → %v, mau Admin", awal, got)
		}
	}
}

// ⛔ SATU NAMA DI LAYAR, AKSES TETAP MILIK SEMUA PEMEGANG WORKBASKET.
//
// Permintaan pemilik proses 8 Oktober 2026: *"dibuat salah satu nya di
// tampilan tp bisa diakses semua yg dapat Workbasket itu"*.
//
// ⚠️ Kedua kalimat itu DUA tuntutan, dan yang kedua yang mudah terlanggar
// diam-diam: begitu `PositionUsername` menyusut jadi satu nama, siapa pun
// yang kelak menulis pagar akses tergoda membandingkan nama pemakai dengan
// kolom itu. Uji ini membuat godaan itu merah.
func TestSemuaPemegangTetapBolehBertindak(t *testing.T) {
	g := gudangSimpan()
	g.pemegangPosisi[models.PosisiSecHead] = []string{"SEC2", "SEC1"}
	g.pemegangPosisi[models.PosisiDeptHead] = []string{"DEPT1"}
	l := services.LayananDengan(g)
	ctx := context.Background()

	if _, err := l.KirimKontrak(ctx, admin, services.MasukanKirim{
		MasukanSimpan: services.MasukanSimpan{IDKontrak: "1001001"}, Aksi: services.AksiSubmit,
	}); err != nil {
		t.Fatal(err)
	}
	d := g.disimpan[0].Dokumen
	if d["PositionUsername"] != "SEC1" {
		t.Fatalf("PositionUsername %v, mau SEC1", d["PositionUsername"])
	}
	g.dokumenTersimpan["1001001"] = d

	// ⭐ SEC2 TIDAK disebut di `PositionUsername`, tetapi ia memegang
	// workbasket yang sama — dan ia HARUS tetap dapat menerima.
	sec2 := inti.Pelaku{AkunID: "SEC2", Peran: []string{models.PosisiSecHead}}
	if _, err := l.KirimKontrak(ctx, sec2, services.MasukanKirim{
		MasukanSimpan: services.MasukanSimpan{IDKontrak: "1001001"},
		Aksi:          services.AksiAkseptasi, Pilihan: models.PilihAccept,
	}); err != nil {
		t.Fatalf("pemegang yang tidak disebut ditolak: %v", err)
	}
	if got := g.disimpan[1].Dokumen["Position"]; got != models.PosisiDeptHead {
		t.Errorf("posisi %v, mau DeptHead", got)
	}

	// ⛔ Dan yang BUKAN pemegang tetap ditolak — pagarnya tidak ikut longgar.
	asing := inti.Pelaku{AkunID: "SEC1", Peran: []string{models.PosisiAdmin}}
	if _, err := l.KirimKontrak(ctx, asing, services.MasukanKirim{
		MasukanSimpan: services.MasukanSimpan{IDKontrak: "1001001"},
		Aksi:          services.AksiAkseptasi, Pilihan: models.PilihAccept,
	}); !errors.Is(err, services.ErrBukanPemegangPosisi) {
		t.Errorf("bukan pemegang tidak ditolak: %v", err)
	}
}
