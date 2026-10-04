# E01: Registry predikat EDM — sumber data per-rule

**What to build:** Registry predikat yang sudah ada diperluas sehingga **tiap predikat mendeklarasikan
sumber datanya sendiri**. Sebagian besar predikat endorsement membaca properti kasus, persis seperti
New Business; **enam predikat lini bisnis** membaca kolom hasil query atas **polis lama**.

`[terverifikasi]` Dari **202** rule `When` endorsement, **196 membaca properti kasus** dan **hanya 6
membaca hasil query**. Registry tetap **satu** — sumber data adalah atribut **per-rule**, bukan
per-siklus.

`[terverifikasi]` Kolom yang dibaca keenamnya berisi **jenis bisnis polis lama**, diisi satu query
yang berjalan di alur masuk sebelum kasus lahir. Karena itu **urutan eksekusi mengikat**: query
dijalankan lebih dulu, baru predikat dapat dievaluasi.

**Asal (Pega).** `When\IsAneka` · `IsFire` · `IsPA` · `IsMBU` · `IsMarineCargo` · `IsGolfInsurance` ·
`When\IsEDM` · `When\IsNotEDM` · pengisi `RDBList\GetBusinessType_Sql` lewat
`Activity\SetErrorBatalEndorsement_Act` langkah 7.

**Keputusan.** K-050 · K-018 · K-029 · `CLAUDE.md` §4.5, §4.6

**Blocked by:** `..\08-registry-rules-eval.md`
⛔ Tiket itu **belum dikerjakan** — lihat §0 indeks.

**Status:** sebagian — 01-10-2026. Registry varian EDM ada: paket `backend/services/predikat/`
(`Eval`, `Sumber`, `KasusEDM`), `registry_gen.go` dibangkitkan generator nbfacin opsi (b) atas
`Endorsment Fac In\When` (202 berkas, 201 predikat; sensus dua cara cocok) lalu disaring
`docs/alat/saring_registry.py` (lapis kedua: medan identitas — nomor polis, nama marketing, peran
operator — dengan literal → panic tanpa nilai; penjaga `TestRegistryTanpaLiteralIdentitas`). Mesin evaluator SALINAN
nbfacin `rules` (asal + sha256 di `predikat.go`; dasar: tafsiran jawaban work owner "EDM DAN NB MENU
YANG TERPISAH"). Terpenuhi: registry satu + sumber per-rule, 37 cabang, urutan query mengikat (panic),
ekspresi tersimpan, nama tak dikenal panic, asal per predikat, ketiga K-046. Belum: resolver skala
K-018 (milik mesin premi).

- [ ] Registry **satu**; sumber data dideklarasikan **per-rule**
- [ ] **37 cabang** lini bisnis terimplementasi (5 · 25 · 1 · 4 · 1 · 1) — **bukan 203**
- [ ] Kolom sumber = jenis bisnis **polis lama**, versi terakhir (`COUNT−1`; **aman**, tidak ada penghapusan baris — K-050)
- [ ] Urutan mengikat: query pengisi berjalan **sebelum** predikat lini bisnis dievaluasi; bila belum, keenamnya bernilai salah dan percabangan runtuh **diam-diam**
- [ ] **Kedua** tag kondisi dibaca; yang mengikat **ekspresi tersimpan**, bukan label tampilan
- [ ] Nama rule tak dikenal → `panic`
- [ ] Predikat lini bisnis menyalakan resolver skala (K-018) — salah lini berarti salah satuan rate
- [ ] Tiap predikat menyebut rule Pega asalnya dalam komentar (§4.6)
- [ ] **K-046** `K046_IsMBU_CabangGanda_HasilTidakBerubah` — dua nilai muncul dua kali; redundan, hasil tidak berubah
- [ ] **K-046** `K046_IsAneka_LabelKodeBisnis_Diabaikan_IkutiValue1` — label tampilan tidak berkaitan dengan ekspresi tersimpan
- [ ] **K-046** `K046_COB_EDM_TanpaDelegasiSubRule` — endorsement membandingkan literal, tidak memanggil sub-rule
