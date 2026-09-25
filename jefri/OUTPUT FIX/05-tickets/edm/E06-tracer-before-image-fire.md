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

**Status:** blocked

- [ ] **Seam 4 `services/endorsement.PrepareBeforeImage` hidup**
- [ ] **Tiga lapis tetap terpisah** di model domain — menyatukannya mengubah angka
- [ ] Lapis A memuat **versi terakhir** polis dan menjalankan penyalinan berpola
- [ ] Data kerja setelah penyalinan **identik** dengan polis lama
- [ ] Lapis B mengisi properti nilai lama lini kebakaran; nilai kosong → **nol**, bukan nilai tak-ada
- [ ] Lapis C menandai baris warisan sesuai kedalaman daftar
- [ ] Porsi periode: presisi desimal tinggi + **guard pembagian nol**; bertipe **`Ratio`**, bukan `Money`
- [ ] ⚠️ Porsi periode **bersifat sementara** — ditimpa modul selisih; **urutan eksekusi before-image → selisih wajib dijaga** (K-048)
- [ ] Nilai uang bertipe `Money`; `Money + Ratio` **gagal saat kompilasi**
- [ ] **Satu fixture kebakaran rekonsiliasi eksak**
- [ ] Tiap penyalinan menyebut rule Pega asalnya (§4.6)
