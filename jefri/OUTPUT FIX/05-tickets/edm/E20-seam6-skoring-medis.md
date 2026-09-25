# E20: Skoring medis — Seam 6

**What to build:** Hasil pemeriksaan laboratorium dan usia tertanggung diubah menjadi **keputusan
underwriting**: diterima standar, ditunda, atau ditolak.

⛔ **Ini BUKAN perhitungan premi.** `[terverifikasi]` Berkas sumbernya **tidak menyentuh** rate,
nilai pertanggungan, premi, maupun pembagian — nol kemunculan keempatnya. Keluarannya **keputusan**,
bukan uang. Karena itu ia tidak masuk seam perhitungan mana pun.

⚠️ **Ambang klinisnya adalah aturan medis dan bisnis**, bukan konstanta teknis. **Diport apa
adanya** — jangan dibulatkan, jangan diseragamkan, jangan ditebak artinya.

⛔ **Domain sensitif.** Nilai ambang, hasil laboratorium, dan data medis **tidak pernah** disalin ke
tiket, test, fixture, log, maupun dokumen. Test memakai **nilai sintetis**.

⚠️ Skoring **risiko** yang bersama New Business adalah sistem **berbeda** dan **bukan** lingkup tiket
ini.

**Asal (Pega).** `Activity\CalculateScorLife_Act` · `SetParamLab_Act` · `CalculatePhysicalExam` ·
`SaveMedical` · `Harness\Medical_Harnes` · `Section\Medical_Sec`

**Keputusan.** **K-052** (Seam 6) · K-046 · `CLAUDE.md` §3.4, §3.5, §4.6

**Blocked by:** E18

**Status:** blocked

- [ ] **Seam 6 `services/underwriting.ScoreMedical` hidup**
- [ ] Keluaran berupa **keputusan** — diterima standar, ditunda, atau ditolak — **bukan angka uang**
- [ ] Skor yang bergantung **usia** bercabang pada rentang yang benar
- [ ] ⛔ Ambang klinis **diport apa adanya**; nilainya dibaca dari korpus, tidak ditebak
- [ ] ⛔ **Nilai ambang, hasil lab, dan data medis TIDAK disalin** ke artefak mana pun
- [ ] Test memakai **nilai sintetis** yang menguji batas, bukan data pasien
- [ ] Keluaran "ditolak" dapat muncul — perilakunya teruji sebagai **keputusan**, bukan galat
- [ ] **K-046** `K046_SkorMedis_AmbangKlinisDipertahankan`
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)

⚠️ **`[pertanyaan terbuka]`** — apa yang terjadi pada kasus berskor "ditolak" bila jiwa **melewati
tangga akseptasi** (E18). Milik work owner dan Underwriting.
