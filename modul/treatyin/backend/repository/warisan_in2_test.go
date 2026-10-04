package repository_test

// Penjaga pembacaan `M_TREATY_IN2` - NOL koneksi Oracle.
//
// ⛔ Yang dijaga di sini pemindaian 41 kolom BERPOSISI. Satu medan yang
// tergeser memindahkan SELURUH nilai sesudahnya ke kolom tetangganya, dan
// hasilnya tetap berupa grid yang terisi rapi - `LIMIT_100` sebuah layer
// tampil sebagai `ADJ_RATE`-nya, dan nol galat menyebutkannya. Uji basis
// data tidak menangkapnya: kueri dan pemindainya salah bersama-sama.

import (
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/repository"
)

// ⭐ Komentar di `pindaiLayer` adalah KONTRAK, bukan hiasan: tiap baris
// penugasan menyebut kolom yang indeksnya ia baca, dan uji ini mengadu
// sebutan itu dengan `kolomLayer` satu per satu.
func TestPindaiLayerSejajarDenganDaftarKolom(t *testing.T) {
	const berkas = "warisan_in2.go"
	f, err := os.Open(berkas)
	if err != nil {
		t.Fatalf("membaca %s: %v", berkas, err)
	}
	defer func() { _ = f.Close() }()
	isi, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	teks := string(isi)

	awal := strings.Index(teks, "func pindaiLayer(")
	if awal < 0 {
		t.Fatal("pindaiLayer tidak ditemukan; pembacanya yang rusak")
	}
	badan := teks[awal:]

	// `Medan: s[12].String, // NAMA_KOLOM`
	pola := regexp.MustCompile(`s\[(\d+)\]\.String,\s*//\s*([A-Z0-9_]+)`)
	cocok := pola.FindAllStringSubmatch(badan, -1)

	kolom := repository.KolomLayerWarisan()
	if len(kolom) != 41 {
		t.Fatalf("kolomLayer memuat %d kolom, mau 41", len(kolom))
	}
	if len(cocok) != len(kolom) {
		t.Fatalf("pindaiLayer menugaskan %d medan berkomentar kolom, kolomLayer %d; "+
			"tiap penugasan WAJIB menyebut kolomnya supaya pergeseran terlihat",
			len(cocok), len(kolom))
	}

	lihat := map[int]bool{}
	for _, m := range cocok {
		var idx int
		for _, c := range m[1] {
			idx = idx*10 + int(c-'0')
		}
		if idx < 0 || idx >= len(kolom) {
			t.Errorf("indeks s[%d] di luar jangkauan 0..%d", idx, len(kolom)-1)
			continue
		}
		if lihat[idx] {
			t.Errorf("s[%d] dibaca dua kali; satu kolom tidak pernah mengisi dua medan", idx)
		}
		lihat[idx] = true
		if kolom[idx] != m[2] {
			t.Errorf("s[%d] diberi komentar %s, tetapi kolomLayer[%d] adalah %s — "+
				"pemindainya bergeser, dan nilainya mendarat di medan tetangga",
				idx, m[2], idx, kolom[idx])
		}
	}
	for i, k := range kolom {
		if !lihat[i] {
			t.Errorf("kolom %s (indeks %d) dibaca dari Oracle tetapi tidak pernah ditugaskan; "+
				"nilainya dibuang diam-diam", k, i)
		}
	}
}

// Ke-41 nama kolom sama persis dengan yang Oracle punya — tidak kurang,
// tidak lebih, tidak ada yang kembar.
//
// ⚠️ Daftarnya dituliskan di sini SEKALI LAGI, dan itu disengaja. Uji yang
// membandingkan sebuah daftar dengan dirinya sendiri tidak pernah gagal;
// yang ini membandingkannya dengan daftar yang ditulis terpisah dari
// katalog Oracle, 3 Oktober 2026.
func TestKolomLayerSamaDenganKatalogOracle(t *testing.T) {
	mau := []string{
		"MASTERID", "PROPORTIONTYPE", "CEDINGID", "CEDING", "SOBID", "SOB",
		"TREATYGROUP", "TREATYCONTRACTNAME", "COMMENCEMENT", "TERMINATION",
		"TREATYTYPE", "CESSIONPCT", "CEDANT_RETENTION", "BASIS_COVER",
		"LAYERTYPE", "LAYER", "SPREADINGTYPE", "CURRENCY", "LIMIT_100",
		"ADJ_RATE", "EARN_PREMIUM", "MDP_RATIO", "MDP", "ROL",
		"CURRENCYRELATION", "CESSION_TO_RI", "EPI100", "RIOGR",
		"BROKERAGEPERCENTP", "EARTHQUAKE", "RNMSHARE", "CURRENCYLIMIT",
		"LIABILITY_RNM", "MDP_RNM_100", "QSOR", "QSRI", "LIABILITYQSRI",
		"LIABILITYQSOR", "EPIRNMQS100", "RNM_RETAINED_PREMI", "RNM_QS_PREMI",
	}
	dapat := repository.KolomLayerWarisan()
	if len(dapat) != len(mau) {
		t.Fatalf("%d kolom, mau %d", len(dapat), len(mau))
	}
	lihat := map[string]bool{}
	for i := range mau {
		if dapat[i] != mau[i] {
			t.Errorf("kolom ke-%d: %q, mau %q", i, dapat[i], mau[i])
		}
		if lihat[dapat[i]] {
			t.Errorf("kolom %s terdaftar dua kali", dapat[i])
		}
		lihat[dapat[i]] = true
	}
}

// ⛔ `ORDER BY` WAJIB, dan wajib berlapis. `LAYER` kolom TEKS: urutan teks
// menaruh layer 10 di antara 1 dan 2, dan grid limit yang layernya tertukar
// TERBACA BENAR.
func TestKueriLayerBerurutAngka(t *testing.T) {
	isi, err := os.ReadFile("warisan_in2.go")
	if err != nil {
		t.Fatal(err)
	}
	teks := string(isi)
	for _, wajib := range []string{
		"ORDER BY TO_NUMBER(LAYER DEFAULT 0 ON CONVERSION ERROR)",
		", LAYER, LAYERTYPE",
	} {
		if !strings.Contains(teks, wajib) {
			t.Errorf("kueri layer tanpa %q", wajib)
		}
	}
	// ⛔ Nol tulisan TIDAK diperiksa di sini. `TestWarisanHanyaDibaca` di
	// paket `backend` sudah menyapu SELURUH berkas Go modul ini untuk kelima
	// tabel warisan, dan ia membuang komentar lebih dulu — pemeriksaan
	// tekstual di sini justru merah oleh kalimat "nol INSERT" di kepala
	// berkasnya sendiri.
}
