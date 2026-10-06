package menu

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
)

// barisUji - potongan menu datar (901), SENGAJA diacak urutannya: menu tidak
// boleh bergantung pada ORDER BY pembacanya.
func barisUji() []Baris {
	return []Baris{
		{ID: 5, Kode: "treatycontractout", Label: "Treaty Contract Out", Golongan: "MASTER", Modul: "treatycontractout", Urutan: 3, Dimigrasi: true},
		{ID: 1, Kode: "claimfacin", Label: "Claim Fac In", Golongan: "KLAIM", Modul: "claimfacin", Urutan: 1},
		{ID: 2, Kode: "claimlife", Label: "Claim Life", Golongan: "KLAIM", Modul: "claimlife", Urutan: 2, Dimigrasi: true},
		{ID: 3, Kode: "nbfacin", Label: "NB FacIn", Golongan: "FACULTATIVE", Modul: "nbfacin", Urutan: 1},
		{ID: 4, Kode: "premiumlistlife", Label: "PremiumList Life", Golongan: "TREATY", Modul: "premiumlistlife", Urutan: 5, Dimigrasi: true},
		{ID: 7, Kode: "komiteclaimlife", Label: "Komite Claim Life", Golongan: "KLAIM", Modul: "komiteclaimlife", Urutan: 6, Dimigrasi: true},
		{ID: 6, Kode: "nbtreatyin", Label: "NB Treaty In", Golongan: "TREATY", Modul: "nbtreatyin", Urutan: 1},
	}
}

var semuaAktif = []string{"claimlife", "premiumlistlife", "komiteclaimlife", "treatycontractout"}

// ringkas menulis menu sebagai teks satu baris per simpul, supaya selisihnya
// terbaca di pesan uji.
func ringkas(m Menu) []string {
	var out []string
	for _, g := range m.Golongan {
		out = append(out, g.Kode)
		for _, k := range g.Modul {
			out = append(out, "  "+k.Kode)
		}
	}
	return out
}

// SATU tingkat di bawah golongan: golongan menurut `Golongan`, modul menurut
// URUTAN lalu ID. Menggantikan `TestSusunGolonganKelompokButirBerurutan`
// (pohon golongan -> kelompok -> butir).
func TestSusunSatuTingkatDiBawahGolongan(t *testing.T) {
	dapat := ringkas(Susun(barisUji(), semuaAktif))
	mau := []string{
		"TREATY", "  nbtreatyin", "  premiumlistlife",
		"FACULTATIVE", "  nbfacin",
		"KLAIM", "  claimfacin", "  claimlife", "  komiteclaimlife",
		"MASTER", "  treatycontractout",
	}
	if !reflect.DeepEqual(dapat, mau) {
		t.Errorf("menu:\n%s\nmau:\n%s", strings.Join(dapat, "\n"), strings.Join(mau, "\n"))
	}
	for _, g := range Susun(barisUji(), semuaAktif).Golongan {
		for _, m := range g.Modul {
			if m.Kode == "claimfacin" && (m.Dimigrasi || m.Label != "Claim Fac In" || m.Urutan != 1) {
				t.Errorf("modul claimfacin: %+v", m)
			}
		}
	}
}

// Modul DIMIGRASI di luar MODUL_AKTIF tidak dikirim (tidak ada halamannya di
// proses ini); modul yang BELUM dimigrasi tetap dikirim - frontend menampilkan
// tombol nonaktif "belum dimigrasi". Menggantikan
// `TestSusunButirModulNonaktifTidakDikirim`.
func TestSusunModulDimigrasiDiLuarModulAktifTidakDikirim(t *testing.T) {
	var kode []string
	for _, g := range Susun(barisUji(), []string{"premiumlistlife"}).Golongan {
		for _, m := range g.Modul {
			kode = append(kode, m.Kode)
		}
	}
	sort.Strings(kode)
	if mau := []string{"claimfacin", "nbfacin", "nbtreatyin", "premiumlistlife"}; !reflect.DeepEqual(kode, mau) {
		t.Errorf("modul dikirim %v, mau %v", kode, mau)
	}
	// Golongan yang seluruh modulnya tersaring tidak dikirim (MASTER).
	for _, g := range Susun(barisUji(), []string{"premiumlistlife"}).Golongan {
		if g.Kode == "MASTER" {
			t.Errorf("golongan MASTER tanpa modul tampil ikut dikirim: %+v", g)
		}
	}
}

// Baris yang BUKAN baris modul (KODE <> MODUL) - lima butir anak 900 selama 901
// belum dijalankan, atau jalur mundurnya - tidak dikirim, juga bila pembaca
// meloloskannya. Menggantikan `TestSusunButirYatimDanTingkatKetigaDibuang`.
func TestSusunBarisBukanModulDibuang(t *testing.T) {
	baris := append(barisUji(),
		Baris{ID: 21, Kode: "inbox", Label: "Inbox Claim Life", Golongan: "KLAIM", Modul: "claimlife", Urutan: 1, Dimigrasi: true},
		Baris{ID: 51, Kode: "tco-tahun", Label: "Treaty Contract Out", Golongan: "MASTER", Modul: "treatycontractout", Urutan: 1, Dimigrasi: true},
	)
	for _, s := range ringkas(Susun(baris, semuaAktif)) {
		if strings.Contains(s, "inbox") || strings.Contains(s, "tco-tahun") {
			t.Errorf("baris bukan modul %q ikut dikirim", strings.TrimSpace(s))
		}
	}
}

// JSON: daftar kosong adalah `[]`, bukan `null` - frontend membaca panjangnya.
func TestSusunDaftarKosongBukanNull(t *testing.T) {
	isi, err := json.Marshal(Susun(barisUji(), nil))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(isi), "null") {
		t.Errorf("JSON memuat null: %s", isi)
	}
	isi, _ = json.Marshal(Susun(nil, nil))
	if string(isi) != `{"golongan":[]}` {
		t.Errorf("menu kosong = %s", isi)
	}
}

// Saringan per akun (Kelola User, keputusan work owner 01-10-2026): hanya
// modul yang KODE-nya dimiliki akun; golongan yang kosong hilang; menu
// aplikasi (Kelola User) tampil di golongan ADMIN, PALING BAWAH, hanya bila
// dimiliki. Modul yang belum dimigrasi ikut disaring - kelak ia tampil hanya
// bagi yang diberi.
func TestSaringMenuUntukAkun(t *testing.T) {
	m := Susun(barisUji(), semuaAktif)
	dapat := ringkas(SaringMenuUntukAkun(m, []string{"claimlife", "nbfacin", KodeKelolaUser, "tidakada"}))
	mau := []string{"FACULTATIVE", "  nbfacin", "KLAIM", "  claimlife", GolonganAdmin, "  " + KodeKelolaUser}
	if !reflect.DeepEqual(dapat, mau) {
		t.Errorf("dapat %v\nmau   %v", dapat, mau)
	}
	if dapat := SaringMenuUntukAkun(m, nil); len(dapat.Golongan) != 0 || dapat.Golongan == nil {
		t.Errorf("akun tanpa menu: %+v, mau golongan kosong (bukan null)", dapat)
	}
	// Menu asal tidak berubah - saringan menyalin.
	if !reflect.DeepEqual(ringkas(m), ringkas(Susun(barisUji(), semuaAktif))) {
		t.Error("SaringMenuUntukAkun mengubah menu asalnya")
	}
}

// Menu aplikasi BUKAN baris tabel menu: golongannya di luar CHECK GROUPMENU,
// dan kodenya tidak bertabrakan dengan nama modul mana pun.
func TestMenuAplikasi(t *testing.T) {
	mau := []Modul{
		{Kode: KodeKelolaUser, Label: "Kelola User", Modul: KodeKelolaUser, Urutan: 1, Dimigrasi: true},
		{Kode: KodeTemplateManager, Label: "Template Manager", Modul: KodeTemplateManager, Urutan: 2, Dimigrasi: true},
	}
	if !reflect.DeepEqual(MenuAplikasi, mau) {
		t.Errorf("MenuAplikasi = %+v", MenuAplikasi)
	}
	for _, g := range Golongan {
		if g == GolonganAdmin {
			t.Errorf("golongan %s ada di Golongan (CHECK GROUPMENU) - menu aplikasi bukan baris tabel", g)
		}
	}
	for _, b := range barisUji() {
		if b.Kode == KodeKelolaUser {
			t.Errorf("kode %s bertabrakan dengan modul", b.Kode)
		}
	}
}

// SQL tanpa PARENT_ID, disaring KODE = MODUL: benar SEBELUM 901 (lima butir
// lama tidak lolos - `inbox` bukan `claimlife`) dan SESUDAHNYA (kolom itu tidak
// ada lagi). Menggantikan `TestSQLMenuHanyaBarisAktifBerurutan`.
func TestSQLMenuBarisModulAktifTanpaParentID(t *testing.T) {
	q := sqlMenu("SKEMAUJI.M_NAV_MENU")
	if benderaYa != "1" {
		t.Errorf("bendera aktif %q, mau \"1\" (konvensi data warisan)", benderaYa)
	}
	for _, mau := range []string{"FROM SKEMAUJI.M_NAV_MENU", "WHERE STATUS_AKTIF = :1 AND KODE = MODUL", "ORDER BY GROUPMENU, URUTAN, ID"} {
		if !strings.Contains(q, mau) {
			t.Errorf("SQL tidak memuat %q: %s", mau, q)
		}
	}
	if strings.Contains(q, "PARENT_ID") {
		t.Errorf("SQL masih menyebut PARENT_ID - mati di ORA-00904 sesudah 901: %s", q)
	}
	if err := db.PeriksaSQL(q); err != nil {
		t.Error(err)
	}
}

// pembacaUji memenuhi PembacaMenu tanpa Oracle.
type pembacaUji struct {
	baris []Baris
	err   error
}

func (p pembacaUji) Baca(context.Context) ([]Baris, error) { return p.baris, p.err }

func minta(t *testing.T, h http.HandlerFunc) (*httptest.ResponseRecorder, map[string]any) {
	t.Helper()
	w := httptest.NewRecorder()
	h(w, httptest.NewRequest(http.MethodGet, "/api/menu", nil))
	var badan map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &badan); err != nil {
		t.Fatalf("badan bukan JSON: %q", w.Body.String())
	}
	return w, badan
}

// Bentuk jawaban PERSIS: {"golongan":[{"kode","modul":[{kode,label,modul,
// urutan,dimigrasi}]}]} - satu tingkat, tanpa "kelompok" maupun "butir".
// Menggantikan `TestRuteMenjawabPohon`.
func TestRuteMenjawabSatuTingkatDiBawahGolongan(t *testing.T) {
	w, badan := minta(t, Rute(pembacaUji{baris: barisUji()}, semuaAktif, true))
	if w.Code != http.StatusOK {
		t.Fatalf("kode %d: %v", w.Code, badan)
	}
	golongan, _ := badan["golongan"].([]any)
	if len(golongan) != 4 {
		t.Fatalf("golongan %d, mau 4: %v", len(golongan), badan)
	}
	for _, g := range golongan {
		bagian := g.(map[string]any)
		if kunci := kunciUrut(bagian); !reflect.DeepEqual(kunci, []string{"kode", "modul"}) {
			t.Errorf("kunci golongan %v, mau [kode modul]", kunci)
		}
		for _, m := range bagian["modul"].([]any) {
			if kunci := kunciUrut(m.(map[string]any)); !reflect.DeepEqual(kunci, []string{"dimigrasi", "kode", "label", "modul", "urutan"}) {
				t.Errorf("kunci modul %v, mau [dimigrasi kode label modul urutan]", kunci)
			}
		}
	}
}

func kunciUrut(m map[string]any) []string {
	var k []string
	for n := range m {
		k = append(k, n)
	}
	sort.Strings(k)
	return k
}

// Tanpa Oracle: 503 berbadan galat - bukan menu kosong diam-diam.
func TestRuteTanpaDatabase503(t *testing.T) {
	w, badan := minta(t, Rute(nil, semuaAktif, false))
	if w.Code != http.StatusServiceUnavailable || badan["galat"] == nil {
		t.Errorf("kode %d, badan %v - mau 503 dengan galat", w.Code, badan)
	}
}

// Tabel belum ada (migrasi 900 belum dijalankan): 503 yang MENYEBUT sebabnya.
func TestRuteTabelBelumAda503MenyebutMigrasi(t *testing.T) {
	err := errors.New("ORA-00942: table or view does not exist")
	w, badan := minta(t, Rute(pembacaUji{err: err}, semuaAktif, true))
	pesan, _ := badan["galat"].(string)
	if w.Code != http.StatusServiceUnavailable || !strings.Contains(pesan, "900") || !strings.Contains(pesan, "-migrate") {
		t.Errorf("kode %d, galat %q - mau 503 yang menyebut migrasi 900 dan -migrate", w.Code, pesan)
	}
}

func TestRuteGagalBaca500(t *testing.T) {
	w, badan := minta(t, Rute(pembacaUji{err: errors.New("ORA-12541: no listener")}, semuaAktif, true))
	pesan, _ := badan["galat"].(string)
	if w.Code != http.StatusInternalServerError || !strings.Contains(pesan, "M_NAV_MENU") {
		t.Errorf("kode %d, galat %q - mau 500 yang menyebut M_NAV_MENU", w.Code, pesan)
	}
	// Galat driver tidak ke badan jawaban - ia tinggal di log server.
	if strings.Contains(pesan, "ORA-12541") || strings.Contains(pesan, "listener") {
		t.Errorf("galat driver bocor ke badan jawaban: %q", pesan)
	}
}

// mintaDengan - seperti `minta`, dengan context permintaan yang disiapkan.
func mintaDengan(t *testing.T, h http.HandlerFunc, siap func(*http.Request) *http.Request) (*httptest.ResponseRecorder, Menu) {
	t.Helper()
	w := httptest.NewRecorder()
	h(w, siap(httptest.NewRequest(http.MethodGet, "/api/menu", nil)))
	var m Menu
	if w.Code == http.StatusOK {
		if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
			t.Fatalf("badan bukan menu: %q", w.Body.String())
		}
	}
	return w, m
}

// Siapa melihat apa:
//   - sesi login: HANYA menu akunnya, ditambah golongan ADMIN bila ia memegang
//     Kelola User - stub atau bukan
//   - tanpa sesi, AUTH_STUB=true: seluruh menu tabel, TANPA menu aplikasi
//     (Kelola User menuntut login sungguhan)
//   - tanpa sesi, AUTH_STUB=false: 401
func TestRuteMenyaringMenurutSesi(t *testing.T) {
	p := pembacaUji{baris: barisUji()}
	dariSesi := func(kode ...string) func(*http.Request) *http.Request {
		return func(r *http.Request) *http.Request {
			return r.WithContext(inti.DenganAksesMenu(r.Context(), kode))
		}
	}
	tanpaSesi := func(r *http.Request) *http.Request { return r }
	for _, stub := range []bool{false, true} {
		w, m := mintaDengan(t, Rute(p, semuaAktif, stub), dariSesi("premiumlistlife", KodeKelolaUser))
		if mau := []string{"TREATY", "  premiumlistlife", GolonganAdmin, "  " + KodeKelolaUser}; w.Code != http.StatusOK ||
			!reflect.DeepEqual(ringkas(m), mau) {
			t.Errorf("stub=%v sesi: kode %d menu %v, mau %v", stub, w.Code, ringkas(m), mau)
		}
	}
	w, m := mintaDengan(t, Rute(p, semuaAktif, true), tanpaSesi)
	if w.Code != http.StatusOK || !reflect.DeepEqual(ringkas(m), ringkas(Susun(barisUji(), semuaAktif))) {
		t.Errorf("stub tanpa sesi: kode %d menu %v, mau seluruh menu tabel", w.Code, ringkas(m))
	}
	w, _ = mintaDengan(t, Rute(p, semuaAktif, false), tanpaSesi)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("tanpa sesi dan tanpa stub: kode %d, mau 401", w.Code)
	}
}

// Hak menu LIHAT (migrasi 914): hanya menu ber-hak LIHAT yang menolak tulis - juga bagi superadmin (05-10-2026, B).
func TestBolehUbah(t *testing.T) {
	ctx := inti.DenganMenuLihat(inti.DenganAksesMenu(context.Background(), []string{"accounts", "aggregate"}), []string{"accounts"})
	if BolehUbah(ctx, "accounts") || !BolehUbah(ctx, "aggregate") {
		t.Error("accounts View only, aggregate penuh")
	}
	super := inti.DenganMenuLihat(inti.DenganAksesMenu(context.Background(), []string{"accounts", KodeKelolaUser}), []string{"accounts"})
	if BolehUbah(super, "accounts") {
		t.Error("superadmin ber-hak View only tetap tidak boleh menulis")
	}
	if !BolehUbah(context.Background(), "accounts") {
		t.Error("tanpa sesi (AUTH_STUB) = penuh; gerbang menu yang menjaga")
	}
}
