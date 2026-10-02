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

**Status:** sebagian — 01-10-2026, paket `backend/services/medis/` (`ScoreMedical`). Aturan **dibangkitkan**
dari korpus (`docs/alat/bangkit_medis.py` → `aturan_gen.go`) dan dijalankan penafsir subset ekspresi
Pega (A04) — ambang tidak disalin tangan. Tiga activity: `CalculateScorLife_Act` (210 langkah, 338
penugasan) → `ScoreMedical`; `CalculatePhysicalExam` (3/7) → `PeriksaFisik`; `SetParamLab_Act` (9/438)
→ `NilaiRujukanLab`. Cara kedua: `grep -c '<pxObjClass>Embed-ActivitySteps</pxObjClass>'` dan
`grep -c '<PropertiesName'` per berkas, sepakat; dijaga `TestTabelBangkitanUtuh`. Sisa: `SaveMedical`
(hanya `Obj-Save` — repository E17); keputusan akhir gabungan TIDAK ada di activity mana pun ini.

- [x] **Seam 6 hidup** — `medis.ScoreMedical` (paket `backend/services/medis`; tata letak modul, bukan `services/underwriting`)
- [x] Keluaran berupa **keputusan**, bukan angka uang — teks skor per pemeriksaan apa adanya. ⚠️ Korpus memuat juga "Invalid Age", "150%", "175", dan tiga ejaan "Postpone…"; artinya `[pertanyaan terbuka]`, tidak dipetakan
- [x] Skor yang bergantung **usia** bercabang pada rentang yang benar — ekspresi korpus dijalankan apa adanya (`TestSeluruhAturanBerjalan`). Batas usia tidak diuji per angka karena angkanya tidak boleh masuk test
- [x] ⛔ Ambang klinis **diport apa adanya** — dibangkitkan dari `<PropertiesValue>`, bukan disalin tangan
- [x] ⛔ **Nilai ambang, hasil lab, dan data medis TIDAK disalin** ke tiket, test, fixture, log, maupun dokumen — ambang hanya hidup di kode bangkitan `aturan_gen.go` (port itu sendiri); paket tidak mencatat log
- [x] Test memakai **nilai sintetis** — semantik batas diuji pada ekspresi sintetis (`ekspresi_test.go`), bukan pada ambang korpus
- [x] Keluaran "ditolak" dapat muncul — `TestSerologiReaktifDitolak`
- [x] **K-046** ambang dipertahankan — `TestTabelBangkitanUtuh` (tabel = korpus, setiap ekspresi terurai)
- [x] Menyebut rule Pega asalnya — nomor langkah korpus di tiap baris tabel bangkitan dan di `Galat.Langkah`

⚠️ **`[pertanyaan terbuka]`** — apa yang terjadi pada kasus berskor "ditolak" bila jiwa **melewati
tangga akseptasi** (E18). Milik work owner dan Underwriting.
