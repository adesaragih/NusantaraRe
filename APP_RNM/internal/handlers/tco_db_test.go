//go:build db

// Seam HTTP Treaty Contract Out terhadap skema uji Oracle NYATA.
//
// Jalankan: make test-db (perlu ORACLE_DSN, ORACLE_SCHEMA, ORACLE_SKEMA_UJI).
// Tanpa Oracle seluruhnya MELEWATI dengan pesan.
//
// Tiket 02: saringan jenis reasuransi hanya terbukti benar terhadap master
// yang benar-benar memuat ID yang dikecualikan. Tiket 03: identitas dari
// sequence, upsert ber-ID, anti-dobel sebagai pertanyaan keunikan di basis
// data, jejak di transaksi yang sama.
package handlers_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/internal/config"
	"nusantarare/internal/handlers"
	"nusantarare/internal/repository"
	"nusantarare/internal/repository/skemauji"
	"nusantarare/internal/services"
)

type ujiTCO struct {
	srv      *httptest.Server
	sqlDB    interface{ Close() error }
	db       *repository.DB
	skema    string
	ctx      context.Context
	isiJenis func([]skemauji.JenisReasuransiUji)
	isiGrup  func([]skemauji.GrupTreatyUji)
	// Tiket 12: master kategori lampiran dan folder unggahan uji.
	isiKategori func([]string)
	unggahan    string
	// Tiket 05: master reinsurer AGENT.
	isiAgent func([]skemauji.AgentUji)
	// Tiket 07: master bisnis BUSINESS.
	isiBusiness func([]skemauji.BusinessUji)
	// Tiket 08: sambungan mentah skema uji - untuk memeriksa pola NULL kolom.
	mentah *sql.DB
}

func (u *ujiTCO) sqlDBMentah() *sql.DB { return u.mentah }

// serverTCO memasang skema uji dan Router BER-STUB identitas: rute modul ini
// bergerbang identitas (401 tanpa X-Pelaku), jadi header harus terbaca.
// serverTCO - server uji dengan master grup bawaan 10001 "UJI GRUP" (tahun
// treaty memeriksa grupnya ke master sejak temuan /code-review).
func serverTCO(t *testing.T) (*ujiTCO, func()) {
	u, bersihkan := serverTCOAwal(t)
	u.isiGrup([]skemauji.GrupTreatyUji{{ID: "10001", TreatyGroupName: "UJI GRUP"}})
	return u, bersihkan
}

// serverTCOAwal - server uji dengan seluruh master kosong.
func serverTCOAwal(t *testing.T) (*ujiTCO, func()) {
	t.Helper()
	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
		if !skemauji.BolehDilewati(err) {
			t.Fatalf("skema uji menolak: %v", err)
		}
		t.Skipf("lewati: %v", err)
	}
	ctx := context.Background()
	if err := sqlDB.PingContext(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	if err := skemauji.Pasang(ctx, sqlDB, skema); err != nil {
		t.Fatalf("memasang skema uji: %v", err)
	}
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	db, err := repository.Open(cfg)
	if err != nil {
		t.Fatal(err)
	}
	// Tiket 12: UNGGAHAN_DIR uji - penyimpanan lampiran stub lokal di bawahnya.
	unggahan := t.TempDir()
	u := &ujiTCO{
		srv:   httptest.NewServer(handlers.Router(services.New(db).DenganUnggahanDir(unggahan), true)),
		sqlDB: sqlDB, db: db, skema: skema, ctx: ctx, unggahan: unggahan, mentah: sqlDB,
	}
	u.isiBusiness = func(baris []skemauji.BusinessUji) {
		if err := skemauji.IsiBusinessTCO(ctx, sqlDB, skema, baris); err != nil {
			t.Fatalf("mengisi master bisnis: %v", err)
		}
	}
	u.isiAgent = func(baris []skemauji.AgentUji) {
		if err := skemauji.IsiAgentTCO(ctx, sqlDB, skema, baris); err != nil {
			t.Fatalf("mengisi master reinsurer: %v", err)
		}
	}
	u.isiKategori = func(note []string) {
		if err := skemauji.IsiKategoriLampiranTCO(ctx, sqlDB, skema, note); err != nil {
			t.Fatalf("mengisi master kategori lampiran: %v", err)
		}
	}
	u.isiJenis = func(baris []skemauji.JenisReasuransiUji) {
		if err := skemauji.IsiJenisReasuransiTCO(ctx, sqlDB, skema, baris); err != nil {
			t.Fatalf("mengisi master jenis reasuransi: %v", err)
		}
	}
	u.isiGrup = func(baris []skemauji.GrupTreatyUji) {
		if err := skemauji.IsiGrupTreatyTCO(ctx, sqlDB, skema, baris); err != nil {
			t.Fatalf("mengisi master grup treaty: %v", err)
		}
	}
	return u, func() {
		u.srv.Close()
		_ = skemauji.Bongkar(ctx, sqlDB, skema)
		_ = db.Close()
		_ = sqlDB.Close()
	}
}

func (u *ujiTCO) minta(t *testing.T, metode, jalur string, badan any, beridentitas bool) (int, string) {
	t.Helper()
	var isi io.Reader
	if badan != nil {
		b, err := json.Marshal(badan)
		if err != nil {
			t.Fatal(err)
		}
		isi = bytes.NewReader(b)
	}
	r, err := http.NewRequest(metode, u.srv.URL+jalur, isi)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("Content-Type", "application/json")
	if beridentitas {
		r.Header.Set("X-Pelaku", "UJI-ADMIN")
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	return res.StatusCode, string(b)
}

// Tiket 02 AC 11: dua belas awalan blacklist, Flag active, Type 1/2/3 -
// ditegakkan terhadap master yang memuat yang harus tersingkir.
func TestJenisReasuransiDisaringPersisRD(t *testing.T) {
	u, bersihkan := serverTCO(t)
	defer bersihkan()
	u.isiJenis([]skemauji.JenisReasuransiUji{
		{ID: "10003", Note: "UJI QUOTA SHARE", Tipe: "1", Flag: "active"},
		{ID: "10005", Note: "UJI SURPLUS", Tipe: "2", Flag: "active"},
		{ID: "10006", Note: "UJI XOL", Tipe: "3", Flag: "active"},
		{ID: "10004", Note: "UJI BLACKLIST", Tipe: "1", Flag: "active"},
		{ID: "100041", Note: "UJI AWALAN", Tipe: "1", Flag: "active"},
		{ID: "10217", Note: "UJI BLACKLIST 12", Tipe: "2", Flag: "active"},
		{ID: "10007", Note: "UJI TIPE 4", Tipe: "4", Flag: "active"},
		{ID: "10008", Note: "UJI NONAKTIF", Tipe: "1", Flag: "inactive"},
		{ID: "10009", Note: "UJI FLAG LIFE", Tipe: "1", Flag: "1"},
	})

	kode, badan := u.minta(t, http.MethodGet, "/api/treaty-contract-out/jenis-reasuransi", nil, true)
	if kode != http.StatusOK {
		t.Fatalf("status = %d, badan = %s", kode, badan)
	}
	var jawab struct {
		Daftar []struct{ ID, Note, Tipe string } `json:"daftar"`
		Total  int                               `json:"total"`
	}
	if err := json.Unmarshal([]byte(badan), &jawab); err != nil {
		t.Fatalf("badan bukan JSON yang diharapkan: %v\n%s", err, badan)
	}
	if jawab.Total != 3 || len(jawab.Daftar) != 3 {
		t.Fatalf("total %d / %d baris, mau 3:\n%s", jawab.Total, len(jawab.Daftar), badan)
	}
	for i, id := range []string{"10003", "10005", "10006"} {
		if jawab.Daftar[i].ID != id {
			t.Errorf("baris %d = %s, mau %s (urutan NOTE ASC)", i, jawab.Daftar[i].ID, id)
		}
	}
	for _, r := range jawab.Daftar {
		if !repository.LolosSaringanNonLifeTCO(r.ID, "active", r.Tipe) {
			t.Errorf("SQL meloloskan %s tetapi tabel kebenaran Go menolaknya", r.ID)
		}
	}
}

func TestJenisReasuransiMasterKosongMenjawab503(t *testing.T) {
	u, bersihkan := serverTCO(t)
	defer bersihkan()
	if kode, badan := u.minta(t, http.MethodGet, "/api/treaty-contract-out/jenis-reasuransi", nil, true); kode != http.StatusServiceUnavailable {
		t.Fatalf("master kosong: status = %d, mau 503 (ADR-0015); badan = %s", kode, badan)
	}
}

func TestTreatyContractOutTanpaIdentitasMenjawab401(t *testing.T) {
	u, bersihkan := serverTCO(t)
	defer bersihkan()
	for _, jalur := range []string{"/api/treaty-contract-out/jenis-reasuransi",
		"/api/treaty-contract-out/grup-treaty", "/api/treaty-contract-out/tahun"} {
		if kode, _ := u.minta(t, http.MethodGet, jalur, nil, false); kode != http.StatusUnauthorized {
			t.Errorf("%s tanpa identitas: status = %d, mau 401", jalur, kode)
		}
	}
}

type tahunJSON struct {
	ID, TreatyYear, UnderwritingYear, TreatyGroupID, TreatyGroupName, Proportion string
	StartDate, EndDate, UserID, TglUpdate                                        string
}

func badanTahun(mulai, akhir string) map[string]string {
	return map[string]string{"treatyYear": "2026", "underwritingYear": "2026", "treatyGroupId": "10001",
		"treatyGroupName": "UJI GRUP", "proportion": "10003", "startDate": mulai, "endDate": akhir}
}

// Tiket 03: baru -> ID '1'+6 digit dari sequence; daftar ID DESC; perbarui
// menimpa tanpa baris baru; anti-dobel 409 menyebut baris lain; periode
// terbalik 422; jejak dua catatan.
func TestTahunTreatyLingkaranPenuh(t *testing.T) {
	u, bersihkan := serverTCO(t)
	defer bersihkan()

	kode, badan := u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun", badanTahun("2026-01-01", "2026-12-31"), true)
	if kode != http.StatusOK {
		t.Fatalf("POST: status = %d, badan = %s", kode, badan)
	}
	var baru tahunJSON
	if err := json.Unmarshal([]byte(badan), &baru); err != nil {
		t.Fatal(err)
	}
	if len(baru.ID) != 7 || !strings.HasPrefix(baru.ID, "1") {
		t.Errorf("ID baru %q, mau '1' + 6 digit (AC 6)", baru.ID)
	}
	if baru.UserID != "UJI-ADMIN" || baru.TglUpdate == "" || baru.StartDate != "2026-01-01" {
		t.Errorf("baru: %+v", baru)
	}

	// Baris kedua, periode lain - untuk urutan ID DESC.
	kode, badan = u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun", badanTahun("2027-01-01", "2027-12-31"), true)
	if kode != http.StatusOK {
		t.Fatalf("POST 2: status = %d, badan = %s", kode, badan)
	}
	var kedua tahunJSON
	_ = json.Unmarshal([]byte(badan), &kedua)

	kode, badan = u.minta(t, http.MethodGet, "/api/treaty-contract-out/tahun?halaman=1&ukuran=20", nil, true)
	if kode != http.StatusOK {
		t.Fatalf("GET daftar: %d %s", kode, badan)
	}
	var hal struct {
		Baris []tahunJSON `json:"baris"`
		Total int         `json:"total"`
	}
	_ = json.Unmarshal([]byte(badan), &hal)
	if hal.Total != 2 || len(hal.Baris) != 2 || hal.Baris[0].ID != kedua.ID {
		t.Errorf("daftar harus 2 baris, ID DESC (terbaru dahulu): %+v", hal)
	}

	// Perbarui: seluruh medan tertimpa, cacah tetap 2 (AC 8).
	ubah := badanTahun("2026-01-01", "2026-12-31")
	ubah["treatyGroupName"] = "UJI GRUP DIUBAH"
	ubah["underwritingYear"] = "2025"
	kode, badan = u.minta(t, http.MethodPut, "/api/treaty-contract-out/tahun/"+baru.ID, ubah, true)
	if kode != http.StatusOK {
		t.Fatalf("PUT: status = %d, badan = %s", kode, badan)
	}
	kode, badan = u.minta(t, http.MethodGet, "/api/treaty-contract-out/tahun/"+baru.ID, nil, true)
	var dibaca tahunJSON
	_ = json.Unmarshal([]byte(badan), &dibaca)
	if kode != http.StatusOK || dibaca.TreatyGroupName != "UJI GRUP DIUBAH" || dibaca.UnderwritingYear != "2025" {
		t.Errorf("GET sesudah PUT: %d %+v", kode, dibaca)
	}
	_, badan = u.minta(t, http.MethodGet, "/api/treaty-contract-out/tahun", nil, true)
	_ = json.Unmarshal([]byte(badan), &hal)
	if hal.Total != 2 {
		t.Errorf("pembaruan menambah baris: total %d", hal.Total)
	}

	// Anti-dobel (AC 73): periode + grup baris pertama dipakai baris baru -> 409
	// yang menyebut ID baris pertama; memperbarui baris pertama sendiri lolos.
	kode, badan = u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun", badanTahun("2026-01-01", "2026-12-31"), true)
	if kode != http.StatusConflict || !strings.Contains(badan, baru.ID) {
		t.Errorf("dobel: status = %d, badan = %s (mau 409 menyebut %s)", kode, badan, baru.ID)
	}
	if kode, badan := u.minta(t, http.MethodPut, "/api/treaty-contract-out/tahun/"+baru.ID, badanTahun("2026-01-01", "2026-12-31"), true); kode != http.StatusOK {
		t.Errorf("memperbarui diri sendiri dengan periode yang sama harus lolos: %d %s", kode, badan)
	}

	// Periode terbalik (AC 9) -> 422 menyebut medannya.
	kode, badan = u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun", badanTahun("2028-12-31", "2028-01-01"), true)
	if kode != http.StatusUnprocessableEntity || !strings.Contains(badan, "EndDate") {
		t.Errorf("periode terbalik: status = %d, badan = %s", kode, badan)
	}

	// Jejak (AC 41): dua catatan untuk baris pertama (baru + diperbarui ×2 = 3).
	jejak, err := u.db.JejakTCO(u.ctx, repository.TabelTahunTCO, baru.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(jejak) != 3 || jejak[0].AkunID != "UJI-ADMIN" || jejak[0].Aksi != repository.AksiJejakSimpan {
		t.Errorf("jejak: %+v", jejak)
	}

	// Tidak ada -> 404.
	if kode, _ := u.minta(t, http.MethodGet, "/api/treaty-contract-out/tahun/1999999", nil, true); kode != http.StatusNotFound {
		t.Errorf("tidak ada: status %d, mau 404", kode)
	}
}

func TestGrupTreatyDibacaSajaUrutIDDesc(t *testing.T) {
	u, bersihkan := serverTCOAwal(t)
	defer bersihkan()
	if kode, _ := u.minta(t, http.MethodGet, "/api/treaty-contract-out/grup-treaty", nil, true); kode != http.StatusServiceUnavailable {
		t.Errorf("master kosong: status %d, mau 503", kode)
	}
	u.isiGrup([]skemauji.GrupTreatyUji{{ID: "10001", TreatyGroupName: "UJI GRUP A"}, {ID: "10002", TreatyGroupName: "UJI GRUP B"}})
	kode, badan := u.minta(t, http.MethodGet, "/api/treaty-contract-out/grup-treaty", nil, true)
	if kode != http.StatusOK {
		t.Fatalf("status %d %s", kode, badan)
	}
	var jawab struct {
		Daftar []struct{ ID, TreatyGroupName string } `json:"daftar"`
	}
	_ = json.Unmarshal([]byte(badan), &jawab)
	if len(jawab.Daftar) != 2 || jawab.Daftar[0].ID != "10002" {
		t.Errorf("urutan ID DESC (b587): %+v", jawab.Daftar)
	}
}
