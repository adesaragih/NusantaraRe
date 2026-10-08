package models

// Satu baris `POOLDATA.M_TREATY_IN2` — satu LAYER dari satu kontrak.
//
// ⛔ Keempat tab Limits · Share · Event Limits · RNM Share dibaca dari SATU
// baris ini. Mereka bukan empat sumber; mereka empat PROYEKSI atas satu
// baris, dan memecahnya menjadi empat struct akan membuat empat kueri untuk
// satu baca.
//
// ⚠️ Seluruh medan bertipe `string`, termasuk yang di Oracle bertipe
// `NUMBER`. Alasannya sama dengan tabel pendaratan: lapisan baca tidak
// menafsirkan. `NUMBER(22)` Oracle membawa presisi yang `float64` Go tidak
// sanggup bawa, dan yang memformat untuk layar adalah `formatNumber` /
// `formatPersen` di frontend.
//
// Nama medan Go mengikuti nama KOLOM, bukan nama layar. Nama layar hidup di
// `labels.ts`; menerjemahkannya di sini akan membuat dua kamus.
//
// ⛔ Nama kolom Oracle BUKAN nama properti Pega. Dari 41 kolom, hanya 8 yang
// cocok secara harfiah dengan `pyValue` di ekspor Section. Pemetaannya
// diturunkan dari NILAI — lihat `docs/PEMETAAN-M-TREATY-IN2.md`.
type BarisLayerWarisan struct {
	// --- kepala kontrak, berulang pada setiap baris layer ---
	MasterID       string `json:"masterID"`
	SifatProporsi  string `json:"sifatProporsi"`
	IDCedant       string `json:"idCedant"`
	Cedant         string `json:"cedant"`
	IDAsalBisnis   string `json:"idAsalBisnis"`
	AsalBisnis     string `json:"asalBisnis"`
	KelompokTreaty string `json:"kelompokTreaty"`
	// ⭐ TINGKAT KETIGA pohon tab Limits — `Detail[].COBList[]`.
	//
	// Gambar `05` (prop) dan `31` (non-prop) memperlihatkan grid
	// `Class of Business` DI DALAM baris treaty group yang dibuka, dan pada
	// kontrak contohnya ia berbunyi `No items`. Terukur: 15.742 elemen di
	// seluruh dokumen, jadi ia berisi pada kontrak lain.
	//
	// ⚠️ Nama kelas bisnisnya saja; `ClassOfBusinessID` tidak dibawa — nol
	// layar menampilkannya, dan pengenal yang dibawa tanpa pemakai adalah
	// medan yang kelak dikira berarti.
	KelasBisnis     []string `json:"kelasBisnis"`
	NamaKontrak     string   `json:"namaKontrak"`
	TanggalMulai    string   `json:"tanggalMulai"`
	TanggalBerakhir string   `json:"tanggalBerakhir"`

	// --- layer ---
	JenisTreaty     string `json:"jenisTreaty"`
	PersenCession   string `json:"persenCession"`
	RetensiCedant   string `json:"retensiCedant"`
	DasarCover      string `json:"dasarCover"`
	JenisLayer      string `json:"jenisLayer"`
	Layer           string `json:"layer"`
	JenisPenyebaran string `json:"jenisPenyebaran"`
	MataUang        string `json:"mataUang"`
	Limit100        string `json:"limit100"`
	AdjRate         string `json:"adjRate"`
	PremiEarned     string `json:"premiEarned"`
	RasioMDP        string `json:"rasioMDP"`
	MDP             string `json:"mdp"`

	// ⛔⛔ MATA UANG KEDUA — cacat yang ditemukan 5 Oktober 2026 lewat
	// pengukuran, dan ditutup 6 Oktober.
	//
	// Keempat medan di atas dahulu hanya membawa nilai PERTAMA. Terukur atas
	// 1.854 dokumen, 4.210 elemen `Limits[]`:
	//
	//	MDPList            berpanjang 2 pada 106 layer
	//	PremiumEarnedList  berpanjang 2 pada 100 layer
	//	⇒ dan pada SETIAP satunya kedua elemen bermata uang BERBEDA
	//
	// Jadi yang hilang bukan nilai kembar melainkan nilai yang berdiri
	// sendiri — 106 angka MDP dan 100 angka premi yang tidak pernah sampai
	// ke layar.
	//
	// ⭐ Bentuk DUA KOLOM dibaca dari GAMBAR, bukan diputuskan dari kode:
	// gambar `30` memperlihatkan grid Layers berkolom `100% Limits ( IDR )`
	// · `Deductible ( IDR )` · `100% Limits ( USD )` · `Deductible ( USD )`,
	// dan panel `Summary of Limit` berkolom `MDP (IDR)` dan `MDP (USD)`.
	// Pega menampilkan KEDUANYA berdampingan, bukan memilih salah satu.
	//
	// ⚠️ Pasangannya sejajar dan itu TERBUKTI: pada kontrak `1000003` tiap
	// layer ber-`Currency = IDR` dan `Currency2 = USD`, sementara
	// `MDPList[0]` bermata uang IDR dan `MDPList[1]` USD. Elemen ke-0 milik
	// `MataUang`, elemen ke-1 milik `MataUangLimit`.
	Limit100Kedua      string `json:"limit100Kedua"`
	RetensiCedantKedua string `json:"retensiCedantKedua"`
	MDPKedua           string `json:"mdpKedua"`
	PremiEarnedKedua   string `json:"premiEarnedKedua"`

	ROL             string `json:"rol"`
	RelasiMataUang  string `json:"relasiMataUang"`
	CessionKeRI     string `json:"cessionKeRI"`
	EPI100          string `json:"epi100"`
	RIOGR           string `json:"riogr"`
	PersenBrokerage string `json:"persenBrokerage"`
	Gempa           string `json:"gempa"`
	// ⭐ TIGA BATAS EVENT LIMITS yang tabel lama TIDAK punya.
	//
	// `M_TREATY_IN2` hanya berkolom `EARTHQUAKE`; ketiga saudaranya tidak
	// ada di sana, dan layar lama menampilkan keempatnya (gambar 28 dokumen
	// desain: RSMD Limit · Earthquake Limit · Flood Limit (Jabodetabek) ·
	// Flood Limit (Nationwide), masing-masing dengan mata uangnya).
	//
	// Dokumen punya keempatnya — terukur di `Limits[].Detail[]`:
	// `RSMDLimit` 493 · `Earthquake` 569 · `FloodJab` 218 ·
	// `FloodNation` 550, beserta `CurrencyRSMD`, `CurrencyEarthquake`,
	// `CurrencyFloodJab`, `CurrencyFloodNat`.
	//
	// ⛔ Jadi pencabutan `M_TREATY_IN2` bukan hanya memindahkan sumber: ia
	// mengembalikan tiga kolom yang tabel itu tidak pernah bisa berikan.
	BatasRSMD         string `json:"batasRSMD"`
	BatasBanjirJab    string `json:"batasBanjirJab"`
	BatasBanjirNas    string `json:"batasBanjirNas"`
	MataUangRSMD      string `json:"mataUangRSMD"`
	MataUangGempa     string `json:"mataUangGempa"`
	MataUangBanjirJab string `json:"mataUangBanjirJab"`
	MataUangBanjirNas string `json:"mataUangBanjirNas"`
	RNMShare          string `json:"rnmShare"`
	MataUangLimit     string `json:"mataUangLimit"`
	LiabilityRNM      string `json:"liabilityRNM"`
	MDPRNM100         string `json:"mdpRNM100"`
	QSOR              string `json:"qsor"`
	QSRI              string `json:"qsri"`
	LiabilityQSRI     string `json:"liabilityQSRI"`
	LiabilityQSOR     string `json:"liabilityQSOR"`
	EPIRNMQS100       string `json:"epiRNMQS100"`
	RNMRetainedPremi  string `json:"rnmRetainedPremi"`
	RNMQSPremi        string `json:"rnmQSPremi"`
}
