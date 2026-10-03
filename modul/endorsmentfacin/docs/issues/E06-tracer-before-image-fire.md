# E06: Tracer — before-image tiga lapis untuk satu lini

**What to build:** Ketika kasus endorsement dibuka, "keadaan sebelum" tersedia dalam **tiga lapis
terpisah**, dan satu berkas kasus kebakaran nyata **rekonsiliasi eksak** terhadap Pega.

Ini irisan **tertipis yang lengkap** untuk before-image. Begitu hijau, seluruh premisnya terbukti:
tiga mekanisme berbeda, siklus hidup berbeda, konsumen berbeda.

`[terverifikasi]` **Lapis A** memuat dokumen polis versi terakhir lalu menjalankan penyalinan
sehingga data kerja **identik** dengan polis lama. **Lapis B** mengisi nilai lama per baris.
**Lapis C** menandai baris warisan.

⛔ **Delta dihitung terhadap polis lama, bukan terhadap kosong.** Implementasi **tidak boleh** mulai
dari struktur kosong lalu menambahkan perubahan.

⚠️ Istilah "before-image" adalah **nama konsep**; tidak ada properti bernama itu di korpus.

**Asal (Pega).** `Activity\SetValueToEDMWork` langkah 14 (13 sub-langkah) dan langkah 15 ·
`RDBList\GetEDMOldData_SQL` · `Activity\SetOldData` blok 4.1–4.2 ·
`Activity\SetOLDValueToEDMWork_FIRE` · `Activity\InputAddendumFacIn_PreAct` langkah 28

**Keputusan.** K-047 (Seam 4) · K-010/K-012 · K-018 · K-027 · K-048

**Blocked by:** `..\01-tracer-money-ratio-premi-pa.md` · E01
⛔ Tiket NB itu **belum dikerjakan** — lihat §0 indeks.

**Status:** sebagian — 01-10-2026, `backend/services/{beforeimage,lapisc,nilailama,porsiperiode}.go`.
Dikerjakan TANPA E01: predikat lini dan `IsEDM` diterima sebagai masukan, bukan registry lokal
(pilihan pemegang modul pada sesi 01-10-2026 — belum keputusan work owner). Sisa: pemuatan lapis A
dari Oracle (E17) dan konfirmasi satuan selisih DateTime untuk porsi periode. Temuan korpus yang menyimpang dari tiket: `../LAPORAN-IMPLEMENTASI-BEFORE-IMAGE.md`.

- [x] **Seam 4 hidup** — `services.PrepareBeforeImage` (lahir kasus) + `services.IsiNilaiLama` (lapis B, tiap layar dibuka). ⚠️ Dua fungsi, bukan satu seperti spec §Testing — lihat laporan §3
- [x] **Tiga lapis tetap terpisah** di model domain — menyatukannya mengubah angka
- [ ] Lapis A memuat **versi terakhir** polis — query `GetEDMOldData_SQL` milik repository (E17, ter-block); seam menerima dokumen yang sudah dimuat. Penyalinan berpola ✅
- [x] Data kerja setelah penyalinan **identik** dengan polis lama — sebatas properti yang dimodelkan
- [x] Lapis B mengisi properti nilai lama lini kebakaran; nilai kosong → **nol**, bukan nilai tak-ada
- [x] Lapis C menandai baris warisan sesuai kedalaman daftar
- [x] Porsi periode: presisi 20 desimal + **guard pembagian nol**; bertipe **`Ratio`**; pembulatan **HALF_UP** menjauhi nol (A37, dikonfirmasi work owner 01-10-2026), dihitung eksak dengan bilangan bulat. ⚠️ Selisih yang bukan hari bulat tetap **ditolak** (`GalatPorsiPeriode`) — satuan selisih DateTime (lokal `int`) masih `[dugaan]` hari. Lapis A/B/C tidak ikut gagal
- [x] ⚠️ Porsi periode **bersifat sementara** (K-048) — tercatat di kode; modul selisih belum ada
- [x] Nilai uang bertipe `Money`; `Money + Ratio` **gagal saat kompilasi** — dijamin `inti/backend/uang` (KEPUTUSAN-30-09-2026 butir 3)
- [x] **Satu fixture kebakaran rekonsiliasi eksak** — `TestRekonsiliasiEksakFireEDM` atas fixture NB-15 `edm-fire-1.json` (dibaca di tempat, A01): **78 nilai lapis B/C cocok persis** dengan Pega, termasuk `PremiumOld`/`RateOld` = 0 hasil guard K-046; mutasi yang meluruskan guard memerahkannya. Porsi periode TIDAK direkonsiliasi: nilai tersimpannya ditulis juga oleh `CountPaymentEdm_Act`/`CountPremiEDM_DT` (K-048)
- [x] Tiap penyalinan menyebut rule Pega asalnya (§4.6)
