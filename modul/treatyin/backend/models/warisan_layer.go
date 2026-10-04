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
	MasterID        string `json:"masterID"`
	SifatProporsi   string `json:"sifatProporsi"`
	IDCedant        string `json:"idCedant"`
	Cedant          string `json:"cedant"`
	IDAsalBisnis    string `json:"idAsalBisnis"`
	AsalBisnis      string `json:"asalBisnis"`
	KelompokTreaty  string `json:"kelompokTreaty"`
	NamaKontrak     string `json:"namaKontrak"`
	TanggalMulai    string `json:"tanggalMulai"`
	TanggalBerakhir string `json:"tanggalBerakhir"`

	// --- layer ---
	JenisTreaty      string `json:"jenisTreaty"`
	PersenCession    string `json:"persenCession"`
	RetensiCedant    string `json:"retensiCedant"`
	DasarCover       string `json:"dasarCover"`
	JenisLayer       string `json:"jenisLayer"`
	Layer            string `json:"layer"`
	JenisPenyebaran  string `json:"jenisPenyebaran"`
	MataUang         string `json:"mataUang"`
	Limit100         string `json:"limit100"`
	AdjRate          string `json:"adjRate"`
	PremiEarned      string `json:"premiEarned"`
	RasioMDP         string `json:"rasioMDP"`
	MDP              string `json:"mdp"`
	ROL              string `json:"rol"`
	RelasiMataUang   string `json:"relasiMataUang"`
	CessionKeRI      string `json:"cessionKeRI"`
	EPI100           string `json:"epi100"`
	RIOGR            string `json:"riogr"`
	PersenBrokerage  string `json:"persenBrokerage"`
	Gempa            string `json:"gempa"`
	RNMShare         string `json:"rnmShare"`
	MataUangLimit    string `json:"mataUangLimit"`
	LiabilityRNM     string `json:"liabilityRNM"`
	MDPRNM100        string `json:"mdpRNM100"`
	QSOR             string `json:"qsor"`
	QSRI             string `json:"qsri"`
	LiabilityQSRI    string `json:"liabilityQSRI"`
	LiabilityQSOR    string `json:"liabilityQSOR"`
	EPIRNMQS100      string `json:"epiRNMQS100"`
	RNMRetainedPremi string `json:"rnmRetainedPremi"`
	RNMQSPremi       string `json:"rnmQSPremi"`
}
