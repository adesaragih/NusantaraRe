# 03: Baris `AdjustmentList` + Save ke Outstanding

**Status:** sebagian — Save to RNM ADA sejak GILIRAN-11 paket 1 (`services/simpanrnm.go`); spreading tanpa pemanggil produksi; bendera `Save` = keadaan turunan, tanpa kolom (**OQ-N1 ditutup** GILIRAN-17); cermin mengisi nama/DOB/CEDINGCO (**OQ-N2 ditutup**, sisa OQ-N13); tukar retro dua syarat (**OQ-N5 ditutup**); cabut peserta = penanda (**OQ-M6 ditutup**, migrasi 022); baris pertama **lahir saat Submit Register** sejak GILIRAN-14 (butir bp; `Add` = putaran saja — bo diralat); `Delete` **tidak berlaku, tidak dirender** (OQ-N7 ditutup 29-09-2026); sunting sel: **tidak ada — ikut XML** (butir br dan OQ-N8 ditutup work owner 29-09-2026); `CLAIM_GROSS` = `CLAIM_AMOUNT` **sementara** (OQ-N12 (a) ditutup 29-09-2026) sampai pemilik ekspor menjawab (OQ-N11)

**Blocked by:** 02 (register klaim + penomoran), **14 (skema relasional klaim — PREFACTOR)**

> ⚠️ **Diselaraskan 2026-09-16 — revisi penyimpanan.** Dua perubahan mengikat:
> **(1)** baris adjustment melekat pada **PESERTA** (`T_CLAIMLF_PREMIUMLIST_DETAIL`), bukan pada
> klaim; **(2)** dokumen **per peserta** menjadi **gerbang simpan** ke Outstanding.
> Lihat blok AC "Penyimpanan relasional" dan "Dokumen per peserta" di bawah.

## Hasil & nilai pengguna

Sebagai **ReasLifeSPV**, saya dapat menginput baris **AdjustmentList** pertama pada sebuah klaim
dan menyimpannya ke Outstanding — sehingga klaim itu tercatat sebagai sedang berjalan dan menjadi
antrean kerja yang terlihat. *(User story 3 di spec)*

Ini tiket yang **memperkenalkan unit keputusan** sistem ini: barisnya, bukan klaimnya.

## Area codebase

`internal/models` (baris adjustment sebagai entitas dengan statusnya sendiri), `internal/repository`
(penulisan baris), `internal/services` (aturan penyimpanan ke Outstanding), `internal/handlers`
(endpoint input baris), `frontend/` (formulir dan daftar baris adjustment).

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Activity/SaveOutStandingLife_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` **satu-satunya penulis nilai `STS_REJECT = 0`** di seluruh korpus Claim Life |
| `Claim Life/Section/AdjustmentDetail_Section.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `ADJUSTMENTDETAIL_SECTION` / `RULE-OBJ-HTML-SECTION` | `[terverifikasi]` tombol "Save to Outstanding", dan bentuk baris adjustment |
| `Claim Life/Activity/SetIndexAdjustmentList.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `SETINDEXADJUSTMENTLIST` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` baris baru mewarisi **8 kolom** dari `AdjustmentList(1)`: `SHARE_NUSANTARA_RE`, `CEDING_RETENTION`, `SUM_REASURED`, `SUM_INSURED`, `SHARE_RETRO`, `RETROCEDED_SHARE`, `CURRENCYID`, `CURRENCY` — **tanpa** `STS_REJECT` |

## ADR terkait

**ADR-0011** (unit status = baris `AdjustmentList`), **ADR-0003** (8 kolom uang + `EM_PERCENT`
persen; nama kolom tidak dapat dipakai menebak sifatnya).

## Acceptance criteria

- [x] Baris `AdjustmentList` yang baru disimpan ke Outstanding berstatus **Outstanding** (`0`).
      *(AC 1 spec)* — bukti: `APP_RNM/internal/services/simpanrnm.go:SimpanRNM.Simpan` (langkah 22.1.3.2: `STS_REJECT=0` pada baris tanpa status, lewat `PerbaruiStatusBaris` + jejak); uji `TestStatusBarisKodeLamaKosongMenjadiISNULL` (penulisnya kini dapat menulis baris NULL)
- [x] Sebuah klaim dapat memuat **beberapa** baris adjustment sekaligus, masing-masing dengan
      statusnya sendiri. *(AC 25 spec)* — bukti: `APP_RNM/internal/repository/klaimlife.go:KlaimLife.AmbilBaris` (daftar baris per peserta, `KodeStatus` per baris); uji `TestBarisMelekatPadaPesertanya`, `TestStatusKlaimTurunan`
- [x] Baris yang ditambahkan setelah baris pertama **mewarisi delapan kolom** di atas dari baris
      pertama, dan **tidak** mewarisi status. *(AC 5 spec)* — bukti: `APP_RNM/internal/services/adjustment.go:WarisiKolom` (dipakai `BarisLanjutan`); uji `TestBarisKeduaMewarisiDelapanKolomTanpaStatus`, `TestBarisBaruTidakMewarisiStatus`
- [x] Kedelapan kolom uang diperlakukan sebagai uang; `EM_PERCENT` **tidak** diperlakukan sebagai
      uang. *(AC 23 spec)* — bukti: `APP_RNM/internal/repository/kolompeserta.go:rakitPeserta` (delapan kolom `Money`, `EM_PERCENT` sebagai `Ratio`); uji `TestUraiRasioMengembalikanGalatDanBukanUang`
- [x] Tidak ada nilai uang sebagai *binary floating point* di lapisan mana pun maupun di kontrak API.
      *(AC 22 spec)* — bukti: `APP_RNM/internal/models/money.go:Money.MarshalJSON`; uji `TestUangDiJSONAdalahTeks`, `TestUraiUangTanpaBerubahSatuDigit`

### Penyimpanan relasional ⚠️ BARU 2026-09-16 — spec §2b

- [x] ⚠️ Setiap baris adjustment menunjuk **satu peserta** lewat `PREMIUM_LIST_DETAIL_ID`. Test yang
      menemukan FK adjustment menunjuk **header klaim** **gagal**. *(AC 33 spec; penyimpangan
      sadar 2 — perbaikan relasi, bukan peniruan)* — bukti: `APP_RNM/internal/repository/pohonklaim.go:PohonKlaim.sisipBarisAdjustment` (`PREMIUM_LIST_DETAIL_ID`); uji `TestAdjustmentMenggantungPadaPeserta`
- [ ] ⚠️ Dua peserta dengan masing-masing dua putaran adjustment menghasilkan **empat baris yang
      seluruhnya dapat ditelusuri ke peserta yang benar**. *(AC 34 spec)* — belum: belum ada uji skenario dua peserta × dua putaran (yang terdekat, `TestBarisMenunjukPesertaYangBenar`, memakai tiga baris migrasi), dan belum ada jalur yang melahirkan baris pertama
- [x] ⚠️ Peserta menyimpan **penanda dipilih-untuk-diklaim** (`IS_CHECK`) — inilah penyimpan aturan
      "hanya peserta yang diklaim". *(AC 39 spec; penyimpangan sadar 3)* — bukti: `APP_RNM/internal/repository/kolompeserta.go:insertPeserta` (kolom `IS_CHECK`), `APP_RNM/internal/repository/klaimlife.go:KlaimLife.PasangPenandaDipilih`; uji `TestTambahBarisMenandaiPesertaDipilih`, `TestPenandaDipilihSatuNilaiSaja`
- [x] ⚠️ Tanggal **diterima**, **konfirmasi**, dan **penyelesaian** tersimpan **per peserta**, bukan
      di header. *(AC 40 spec)* — bukti: `APP_RNM/internal/repository/klaimlife.go:KlaimLife.PerbaruiTanggalKlaim` (kolom peserta `CLAIM_RECEIVED_DATE`, `CONFIRMATION_DATE`, `COMPLETE_DATE`; nol kolom itu di header); uji `TestKolomTanggalKlaimAdaDiMigrasi003`, `TestSQLTanggalKlaimBerurutPosisi`
- [ ] Peserta menyimpan `STATUS`, `RECOMMENDATION`, `SOURCE_ID`, dan `CEDING_RETENTION`.
      *(AC 41 spec)* — belum: `SOURCE_ID` dan `CEDING_RETENTION` ditulis (`kolompeserta.go:insertPeserta`), tetapi `STATUS` dan `RECOMMENDATION` hanya ada di DDL 003 — nol penulis maupun pembaca
- [ ] ⚠️ `STS_REJECT` baris adjustment diisi **nilai sebenarnya menurut aksi** — Admin insert
      Outstanding → `0`; SPV tambah Outstanding → `0`. Test yang menemukan nilai di-hardcode
      **gagal**. *(AC 42 spec; penyimpangan sadar 4)* — belum: jalur SPV (`BarisLanjutan` di `Putaran.Tambah`) menulis `0` menurut aksi, tetapi jalur insert Outstanding (Save to RNM) belum ada — `TandaiOutstanding` nol pemanggil produksi
- [x] ⚠️ `ACCEPTATION_DATE` **tidak** distempel saat insert; ia diisi **tanggal akseptasi
      sebenarnya** saat baris benar-benar diaksep. *(AC 43 spec; penyimpangan sadar 4)* — bukti: `APP_RNM/internal/services/statusbaris.go:Transisi` (stempel hanya pada Aksep); uji `TestSimpanKeOutstandingTidakMenstempelTanggalAkseptasi`, `TestAksepMenstempelTanggalAkseptasi`
- [x] Peserta beserta seluruh baris adjustment-nya tersimpan dalam **satu transaksi**. *(AC 49 spec)* — bukti: `APP_RNM/internal/repository/pohonklaim.go:PohonKlaim.Simpan` (satu `tx` dari `Service.DalamTransaksi`); uji `TestSimpanPohonMenulisDuaTempatDalamSatuTransaksi` (uji db)
- [ ] Baris adjustment menyimpan **nama bank**, **id bank**, dan **nomor rekening**, dan ketiganya
      dapat diisi dari layar rincian adjustment. *(AC 56 spec; `[terverifikasi]` class
      `ASM-FW-GISFW-Data-AdjustmentLife`, tampil di `Claim Life/Section/AdjustmentDetail_Section.xml`)* — belum: kolom ada dan ditulis `sisipBarisAdjustment`, tetapi nol layar maupun rute yang mengisinya — tidak ada isian bank di `frontend/src`
- [x] Ketiga field bank **boleh kosong saat Save ke Outstanding** — ia baru menjadi gerbang pada
      **penyerahan ke Komite** (tiket 10). *(AC 57 spec)* — bukti: `APP_RNM/internal/services/simpanrnm.go:PeriksaSimpanRNM` (nol gerbang bank di 29 langkah); uji `TestSimpanRNMLolosSeluruhGerbang` (peserta tanpa field bank lolos)

### Dokumen per peserta — **gerbang simpan** ⚠️ BARU 2026-09-16

`[terverifikasi]` `Claim Life/Activity/SaveOutStandingLife_Act.xml`
(`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY`) beriterasi atas
`.DocumentList` **per peserta** dan menolak simpan dengan
`"The document hasn't been uploaded person number "+<nomor>` serta
`"Documents are incomplete, please complete the documents"`.

- [x] ⚠️ Dokumen tersimpan **per peserta** di `DOCUMENT_CLAIM` dan dapat dibaca dengan `SELECT`
      biasa — **bukan** lewat mekanisme lampiran bawaan. *(AC 44 spec; penyimpangan sadar 5)* — bukti: `APP_RNM/internal/repository/klaimlife.go:KlaimLife.AmbilDokumen` (`SELECT` biasa atas `T_CLAIMLF_DOCUMENT`), `APP_RNM/internal/services/unggahan.go:Unggahan.Unggah`; uji `TestAC05DokumenMenunjukPeserta`
- [x] ⚠️ Menyimpan ke Outstanding **ditolak** bila ada peserta yang dokumennya belum **diunggah**, dengan
      pesan yang **menyebut peserta mana**. *(AC 45 spec; penyimpangan sadar 5)* — bukti: gerbang 1 `APP_RNM/internal/services/simpanrnm.go:PeriksaSimpanRNM` (langkah 3–4, "The document hasn’t been uploaded person number N"); uji `TestSimpanRNMDokumenBelumDiunggahMenyebutSetiapNomor`, `TestSimpanRNMDokumenTidakLengkapLolos`. ⛔ *Teks AC disunting di tempat 28-09-2026 (butir **bl**, OQ-N6 ditutup):* aslinya "belum **lengkap**"; kelengkapan per kategori adalah langkah 12 yang ter-remark (b6178) dan tidak pernah berlaku — `PeriksaDokumenLengkap` dibuang
- [x] Kolom isian `DOCUMENT_CLAIM` **diturunkan dari sensus `.DocumentList`** pada activity di atas,
      dan **keputusannya dicatat** — **jangan tebak dari nama tabel**. *(tiket 14 §Catatan)* — bukti: `APP_RNM/internal/repository/migrations/010_kolom_t_claimlf_document.sql` (tujuh kolom dari sensus `InsertDocument_Act`); uji `TestKolomDDLCocokDenganStruktur`
- [x] Halaman React menampilkan daftar baris adjustment dengan status masing-masing sebagai kata,
      bukan angka. *(AC 26 spec)* — bukti: `APP_RNM/frontend/src/pages/claimlife/KlaimLife.tsx:KlaimLife` (`b.status`, kata dari `BarisAdjustment.MarshalJSON`); uji `TestStatusBarisDariKode`

### Spreading adjustment ⚠️ BARU 2026-09-16

`[terverifikasi]` `Claim Life/Activity/SpreadingClaimLife_Act.xml`
(`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `SPREADINGCLAIMLIFE_ACT` / `RULE-OBJ-ACTIVITY`) menghitung dan
mengisi `.SpreadingList` (per treaty-year, class `ASM-FW-GISFW-Int-TREATYYEAR_LIFE`) beserta
`.RetroLifeList` (per reinsurer, class `ASM-FW-GISFW-Int-RETROCESSIONLIFE`) pada **baris
adjustment**.

- [ ] ⚠️ Menyimpan baris adjustment **menghitung dan menyimpan** spreading-nya: satu baris per
      treaty-year, dan di bawahnya satu baris per reinsurer. *(AC 59 spec; tiket 14; penyimpangan sadar —
      spreading dibekukan)* — belum: `HitungSpreading` nol pemanggil produksi dan masukan rate-nya tanpa pembaca (OQ-M7) — tidak ada spreading yang tersimpan
- [ ] ⚠️ Setiap baris spreading menunjuk **satu baris adjustment**; setiap baris spreading retro
      menunjuk **satu baris spreading**. Test yang menemukan keduanya menggantung pada peserta atau
      pada header klaim **gagal**. *(AC 58 spec; tiket 14)* — belum: DDL benar (`FK_SPR_ADJ`, `FK_SPR_RETRO_SPR` di migrasi 005/006), tetapi nol uji yang gagal bila FK itu menunjuk peserta/header — `TestKaskadeHanyaPadaRelasiTerdaftar` hanya memeriksa ada-tidaknya `ON DELETE CASCADE`
- [ ] Nilai turunan tersimpan sesuai perhitungan yang terbukti, **dengan pembulatan dan pembagian
      seribu**: `RATE_tersimpan = bulat(rate_mentah, 4)`; `rate_pakai = bulat(rate_mentah ÷ 1000, 10)`;
      `AMOUNT = RetrocadedShare × bulat(PERCENT_SHARE ÷ 100, 4)`;
      `PREMIUM_SPREADED_GROSS = rate_pakai × (1 + EM_PERCENT) × AMOUNT`, dengan `EM_PERCENT`
      dipakai sebagai **pecahan langsung** (tanpa ÷ 100).
      *(`[terverifikasi]` `SpreadingClaimLife_Act` langkah 8.2.1.9.2 baris 4396-4509)* — belum: rumus teruji (`TestRateDibagiSeribuSekaliSaja`, `TestUangDibulatkanDelapanDesimal`), tetapi `HitungSpreading` nol pemanggil — nilainya tidak pernah tersimpan
- [x] ⭐ `[ditutup oleh XML — 26-09-2026]` **`PREMIUM_SPREADED_NET` dipilih oleh tahun polis.**
      Kedua cabang memang ada, dan **aksi precondition-nya** yang memutuskan: langkah 8.2.1.9.3
      (`pyStepsDescription` = *"Tahun pertama"*, `local.Year==1` **WhenFalse=3 lewati**) memakai
      `NET = GROSS − Discount − Comm` dengan `Discount = GROSS × bulat(COMMISION ÷ 100, 5)` dan
      `Comm = (GROSS − Discount) × bulat(OVR_COMM ÷ 100, 5)`; langkah 8.2.1.9.2
      (*"Bukan tahun pertama"*, `local.Year==1` **WhenTrue=3 lewati**) memakai `NET = GROSS − Comm`
      dengan `Comm = GROSS × bulat(OVR_COMM ÷ 100, 5)`.
      `local.Year = tahun(GROSS_VALUATION_BEGIN_DATE) − tahun(BEGIN_DATE) + 1`.
      *(baris 4372-4599 dan 4640-4907; meralat `[terbuka]` tiket 14 §Blocker)* — bukti: `APP_RNM/internal/services/spreading.go:HitungSpreading`, `TahunPolis`; uji `TestTahunPertamaMemakaiDiscount`, `TestBukanTahunPertamaTanpaDiscount`, `TestTahunPolisDihitungDariTahunSaja`
- [ ] ⚠️ Nilai spreading **dibekukan**: perubahan master treaty sesudahnya **tidak mengubah** angka
      yang sudah tersimpan pada adjustment itu. *(AC 59 spec)* — belum: tidak ada spreading yang tersimpan (`HitungSpreading` nol pemanggil); skemanya memang tabel, bukan view (`TestAC42SpreadingDibekukanBukanTurunan`)
- [x] ⭐ **BARU menurut XML 26-09-2026 — kaskade kapasitas per treaty-year.** Sisa klaim mengalir
      dari treaty-year satu ke berikutnya: bila `CURRENCY=="IDR"` dan sisa `≤ IDR` treaty-year itu,
      `RetrocadedShare = sisa` dan sisa menjadi `0`; bila sisa `> IDR`, `RetrocadedShare = IDR` dan
      sisa berkurang `IDR`. Sama persis untuk `"USD"` dengan kolom `USD`. Mata uang selain keduanya
      tidak punya cabang di XML. *(`[terverifikasi]` langkah 8.2.1.4-7, baris 3030-3729)* — bukti: `APP_RNM/internal/services/spreading.go:HitungSpreading`; uji `TestKaskadeKapasitasMengalirKeTreatyBerikut`, `TestMataUangTanpaCabangDitolak`
- [x] ⭐ **BARU — pemilihan rate.** Baris `RATE_LIFE` dipilih per baris retro: bila `GENDER`
      **memuat `U`** cocokkan `AGE` saja; bila tidak, cocokkan `AGE` **dan** `SEX`. `CONTRACT` **tidak**
      ikut dalam precondition yang aktif. *(`[terverifikasi]` langkah 8.2.1.9.1, baris 4084-4284)* — bukti: `APP_RNM/internal/services/spreading.go:PilihRate`; uji `TestRateUnisexMengabaikanJenisKelamin`, `TestRateBerjenisKelaminMenuntutKecocokan`
- [ ] Adjustment tanpa retrosesi tersimpan dengan **nol baris** spreading — **bukan** kegagalan. — belum: aturan murni menghasilkan nol baris tanpa galat (`TestTanpaTreatyNolBaris`), tetapi jalur simpan adjustment + spreading belum ada
- [x] Baris adjustment beserta seluruh spreading dan spreading retro-nya tersimpan dalam **satu
      transaksi**. *(AC 49 spec)* — bukti: `APP_RNM/internal/repository/pohonklaim.go:PohonKlaim.Simpan` (adjustment, spreading, dan retro dalam satu `tx`)

## Catatan

`[terverifikasi]` Empat kolom bernama `SHARE_*` / `*_SHARE` ternyata **uang**, bukan rasio — nama
kolom di korpus ini terbukti menipu (bandingkan `STS_REJECT`, yang nilai `1`-nya berarti *diaksep*).
Klasifikasi mengikat ada di **ADR-0003**; jangan menyimpulkan dari nama.

`[terbuka]` **OQ-060** (pemilik **Product+UW**) — apakah seluruh baris satu klaim wajib bermata uang
sama. `[terverifikasi]` dalam praktiknya seragam karena `CURRENCY` disalin dari baris 1, tetapi
strukturnya membolehkan campur. **Tidak memblokir tiket ini**; memblokir bentuk akhir tipe uang.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Implementasi — 26 September 2026 malam (prasyarat, sesi batch 02–06)

**Status: `claimed`** — **0 dari 26 AC tertutup.** Yang dikerjakan commit ini adalah **dua prasyarat
§10 brief**, bukan AC tiket ini. Dinyatakan terus terang supaya tidak terbaca sebagai kemajuan yang
tidak ada.

Test Go **111 → 112**, nol FAIL; test bertag `db` **23 → 24**, seluruhnya SKIP; 88 modul.

### ⛔ Cacat tersembunyi yang ditutup

`services/pendaftaran_db_test.go` (tiket 02) memanggil `AmbilUntukKlaim`, yang membaca
`M_LIFE_PREMIUM_DETAIL` — tabel yang **skema uji tidak pernah buat**. Hari ini ia SKIP bersama yang
lain, jadi cacatnya tak terlihat; **dengan Oracle ia akan GAGAL di pembacaan peserta**, bukan
menguji pendaftaran. Tiga test db tiket 02 karena itu selama ini kosong isinya.

Sekarang skema uji membuat tiruannya: **24 kolom `kolomSalin` + `EDMSTATUS`**, tipe `[data DBA]`
dari katalog — bukan diturunkan dari nama. Fixture dua peserta: satu ber-`EDMSTATUS` **NULL**
(new business, harus muncul) dan satu **`Batal`** (harus disaring keluar), sehingga penyaring hidup
akhirnya teruji terhadap Oracle dan bukan hanya terhadap dirinya sendiri.

⛔ Kolom `KTP` **tidak ikut ditiru**: ia tidak pernah dibaca, dan tabel uji yang menyediakan tempat
untuk nomor identitas adalah undangan.

### ⛔ Satu cacat saya sendiri, ditemukan saat memverifikasi tipe kolom

`kolomSalin` membaca **`STNC` apa adanya**, padahal katalog menyebutnya **`DATE`** (kolom 29).
Membaca kolom tanggal tanpa `TO_CHAR` membuat bentuknya bergantung `NLS_DATE_FORMAT` sesi — jebakan
yang sama yang diperangi sepanjang tiket 14, dan saya sendiri yang memasangnya di tiket 02 lanjutan.
Sudah dibungkus `TO_CHAR`.

⚠️ `[terbuka]` Kolom tujuannya di `003` bernama `STNC_TREATY` dan bertipe `VARCHAR2(64)`; sumbernya
`DATE`. Selisih tipe itu belum diputuskan siapa pun, dan executor tidak mengubah DDL tanpa keputusan.

### Penjaga posisi `kolomSalin`

`salinKePeserta` membaca hasil `SELECT` lewat **indeks tetap 0–23**. Satu kolom yang disisipkan di
tengah menggeser seluruh sisanya **tanpa satu pun galat** — nilai hanya mendarat di medan yang
salah, dan itu baru terlihat jauh di hilir kalau pernah terlihat. `TestUrutanKolomSalinDikunci`
mengunci cacahnya (24) dan nama kolom pada tiap posisi.

### Tiruan tabel treaty

`RETROCESSIONLIFE` ditiru **dengan tipe aslinya** — di instance pengembangan ia **VIEW** ber-13
kolom yang seluruhnya `VARCHAR2(4000)`, termasuk `PERCENTSHARE`, `RATE`, `COMMISION`, `OVR_COMM`.
⭐ Tiruannya sengaja **tidak** dibuat `NUMBER` yang "lebih benar": justru jalur *teks → `ParseDecimal`
→ laporkan yang gagal* itulah yang perlu diuji, dan tiruan bertipe `NUMBER` membuat Oracle mengurai
angkanya lebih dulu sehingga pembacanya tidak pernah menemui teks. `TREATYYEAR_LIFE` ikut, 7 kolom.

### ⛔ `RATE_LIFE` tidak ditiru — dan itu memblokir satu masukan rumus spreading

Katalog baru memuat **enam kolom pertamanya** (`ID`, `IDUSEDBY`, `USEDBY`, `TYPE`, `GENDER`,
`CONTRACT`), dan **tidak satu pun di antaranya kolom rate**. Menirunya berarti mengarang bentuk, dan
membaca `RATE` dari tabel yang bentuknya dikarang berarti mengarang angkanya.

Akibatnya pada AC 22: rumus `PREMIUM_SPREADED_GROSS = RATE × (1 + EM_PERCENT) × AMOUNT` kehilangan
sumber `RATE`-nya. ⚠️ `[data DBA]` — daftar kolom `RATE_LIFE` beserta tipenya diperlukan sebelum
pembacanya ditulis.

### Penjaga AC 29 dipertajam, bukan dilonggarkan

Tiruan tabel peserta **mendeklarasikan** kolom `EDMSTATUS`, dan penjaga lama menuduhnya sebagai
penyaring kedua. Aturannya kini membedakan **menyaring** dari **menyebut**: yang dilarang adalah
bentuk `EDMSTATUS IS`/`NOT IN`/`IN`/`=` di luar pembacanya. ⚠️ Percobaan pertama saya memakai pola
`"EDMSTATUS)"`, yang cocok dengan **daftar kolom INSERT** — penjaga yang menuduh hal yang bukan
aturan. Diperbaiki, lalu dibuktikan masih menangkap penyaring kedua yang sungguhan.

### Yang BELUM dikerjakan — seluruh AC tiket ini

26 AC masih terbuka. Yang belum ada sama sekali: `Adjustment.Tambah` dan pewarisan delapan kolom,
`Spreading.Hitung`, pembaca treaty, gerbang dokumen, langkah migrasi `010` (butir **ad**), kedua
pintu HTTP, dan seluruh frontend-nya. Sensus `.DocumentList` untuk **ad** baru sampai pada daftar
berkas korpusnya, belum pada kolomnya.

---

## Pembacaan ulang XML — 26 September 2026 malam (aturan brief modul §1.3)

Perintah pecahan yang dipakai, supaya nomor baris di bawah dapat direproduksi orang lain: `sed`
memecah tiap `><` menjadi dua baris, lalu `sed -E` mengembalikan `&lt; &gt; &quot; &amp;` ke bentuk
aslinya; hasilnya ditulis ke scratchpad sebagai `<Rule>-split.xml`. Seluruh nomor baris di bawah
adalah nomor baris di berkas hasil pecahan itu.

**Cara membaca precondition Pega — ditetapkan dari korpus, bukan dari dokumentasi.** Tiap baris
`pyStepsPreCondParamsWhen` berpasangan dengan `…WhenTrue` dan `…WhenFalse`; nilainya **`2` = lanjut**
(evaluasi baris berikutnya, lalu jalankan langkah) dan **`3` = lewati langkah**. Buktinya sepasang
langkah yang saling melengkapi: 8.2.1.9.1 berjudul *"set Rate jika ada U"* memberi
`@contains(.CARI3,"U")` True=2 / False=3, sedangkan 8.2.1.9.2 berjudul *"set Rate"* memberi kondisi
yang **sama** dengan True=3 / False=2. Hanya satu pembacaan yang membuat keduanya masuk akal
sekaligus. Dengan kunci ini seluruh cabang di bawah terbaca — termasuk dua cabang
`PREMIUM_SPREADED_NET` yang selama ini `[terbuka]`.

| Rule | Byte | Tag yang dibaca | Yang diambil |
| --- | ---: | --- | --- |
| `Activity/SpreadingClaimLife_Act.xml` | 266.985 | 9 langkah beserta anaknya, `pyStepsObjectName`, precondition **beserta aksinya** | seluruh rumus, kaskade kapasitas, pemilihan rate, cabang tahun polis |
| `Activity/SaveOutStandingLife_Act.xml` | 642.787 | langkah 3, 4, 12, 22 | dua gerbang dokumen yang **berbeda**, penulis `STS_REJECT=0` |
| `Activity/SetIndexAdjustmentList.xml` | 47.748 | tiga langkah utuh | pewarisan 8 kolom **dan** `.IsCheck=true` |
| `Activity/InsertDocument_Act.xml` | 89.600 | `Obj-Save` + seluruh Property-Set | sensus kolom dokumen (butir **ad**) |
| `Section/AdjustmentDetail_Section.xml` | 557.563 | `pyCondition`, `pyActivity`, nama properti | gerbang peran, himpunan field, **enam** field bank |
| `Section/RetroDetailClaimLife.xml` | 121.355 | `pyPageListProperty`, `pyValue` | kolom tabel retro di layar |
| `Section/DocumentLife.xml` | 201.204 | `pyPageListProperty`, `pyCondition` | daftar dokumen memakai `.DocumentClaimList` |
| `RDBList/GetRateRetro.xml` | 6.243 | `pyBrowseSQL` utuh | SQL rate |
| `RDBList/GetRetroLife_SQL.xml` | 5.962 | `pyBrowseSQL` utuh | SQL retrosesi |
| `RDBList/GetProductLife.xml` | 6.348 | `pyBrowseSQL` utuh | SQL produk |
| `RDBList/GetJsonProductLife.xml` | 6.615 | `pyBrowseSQL` utuh | SQL treaty-year |

### ⭐ Yang berubah dari yang brief dan tiket kira

**1. `GetJsonProductLife` bukan `select * from treatyyear_life` belaka.** SQL utuhnya
`[terverifikasi]` `GetJsonProductLife.xml` baris 84–94 memilih `* from treatyyear_life` dengan
`where id in (...)`, dan subkueri itu membaca `m_product_life a` lewat
`json_table(A.JSONDATA, '$' COLUMNS(NESTED PATH '$.OutwardList[*]' COLUMNS(OUTWARDNAME, OUTWARDNAMEID)))`
menyaring `a.id = {ParamData.CARI1}`, lalu ditutup `order by TO_NUMBER(IDR) asc`.

Tiga akibat. **(a)** Objek `treatyyear_life` yang dilihat koneksi Pega **punya kolom `IDR`** — ia
dipakai di `ORDER BY`. Katalog `POOLDATA.TREATYYEAR_LIFE` yang tujuh kolom tanpa `IDR`/`USD` karena
itu **bukan objek yang sama**; selisih ini `[data DBA]`, bukan tebakan. **(b)** Urutan kaskade
kapasitas **ditentukan XML**: `TO_NUMBER(IDR)` menaik — kapasitas terkecil lebih dulu. **(c)** Daftar
treaty-year yang berlaku bagi sebuah klaim datang dari `OutwardList[*].OUTWARDNAMEID` **di dalam
JSON produk** — sumber yang tiket 14 AC 38 larang dibaca. Jadi bukan hanya `RATE` yang menggantung
pada modul Master Product Name Life (butir **ag**): **daftar treaty-year pun**.

Ini sekaligus meralat ralat katalog tadi malam. `GetProductLife` memang yang mengembalikan
`M_PRODUCT_LIFE.JSONDATA AS CARI1` (join ke view `PRODUCT_LIFE` untuk `RICOMM`), tetapi
**`GetJsonProductLife` membaca JSON yang sama** sebagai penyaring. Larangan AC 38 mengenai
**keduanya**.

**2. Peran yang boleh Save to Outstanding adalah `ReasLifeSPV`, bukan `ReasLifeAdmin`.**
`[terverifikasi]` `AdjustmentDetail_Section.xml`: dua tombol `SaveOutstandingLife_Act` (baris 16249
dan 16396) diikuti `pyCondition` `pyWorkPage.pyPosition =='ReasLifeSPV'` (baris 16468). Pola
pasangan yang sama berlaku pada jalur Komite (`GetListKomiteLife` 15564/15873, kondisi 16091
`pyPosition=='ReasLifeSPV' || Type='TP' || Type='TR'`) dan pada Reject (kondisi 15399
`pyPosition=='ReasLifeAdmin' && CLAIM_NO!='' && STS_REJECT==0`). Bab "Hasil & nilai pengguna" tiket
ini **diralat** dari *ReasLifeAdmin* menjadi *ReasLifeSPV*; penegakannya milik **tiket 07**.

**3. Ada DUA gerbang dokumen, bukan satu, dan keduanya memakai daftar yang berbeda.**

> ⛔ **Ralat 28-09-2026 (temuan /code-review GILIRAN-11).** Langkah 12 ber-`pyStepsBlockName = //`
> (b6178): gerbang **kedua** di bawah TER-REMARK dan tidak pernah jalan. Yang hidup hanya gerbang
> pertama. Tabelnya dibiarkan sebagai catatan pembacaan; lihat bab GILIRAN-11 di akhir tiket.

| Gerbang | Langkah | Daftar | Syarat | Pesan |
| --- | --- | --- | --- | --- |
| dokumen **ada** | 3, 3.1, 4 | `.DocumentClaimList` | jalan hanya bila `Type` **bukan** `TP`/`TR`, peserta `.IsAccept=="true"`, dan `@SizeOfPropertyList(.DocumentClaimList)<1` | `local.Errmsg10` bertambah *"The document hasn't been uploaded person number "* + **nomor urut peserta** (`local.idx2 = .pxListSubscript`) |
| dokumen **lengkap** | 12, 12.1–12.4 | `.DocumentList` (`.pyCategory`) | `Category2` diisi tiap `.pyCategory`, di-**dedup** oleh langkah Java atas `.CARI1`, lalu dibandingkan; pesan muncul bila `@LengthOfPageList(Category2.pxResults)` **tidak sama dengan** `@LengthOfPageList(Category.pxResults)` | `Local.Err4` = *"Documents are incomplete, please complete the documents"* |

⭐ Aturan **"dokumen lengkap"** karena itu **terbaca**: *cacah kategori berbeda yang terunggah sama
dengan cacah baris kategori wajib*. Yang **tidak** terbaca: isi daftar wajibnya —
`GetCategoryLife_SQL` **tidak ada di korpus** (29 berkas `RDBList/` sudah dicacah seluruhnya).
`[terbuka — data DBA/work owner]`.

⚠️ "person number" adalah **nomor urut peserta dalam daftar**, bukan pengenal peserta.

**4. `STS_REJECT = 0` ditulis per baris adjustment dan bergerbang `PrintFaceClaim`.**
`[terverifikasi]` langkah 22.1.3 mengulang `.AdjustmentList`; 22.1.3.1 menjalankan insert ke OS
akseptasi dengan precondition `.PrintFaceClaim==""` (True=2 lanjut, False=3 lewati); 22.1.3.2 lalu
menyetel `.PrintFaceClaim=1` **dan** `.STS_REJECT=0` (baris 10380–10383). Jadi `0` memang **nilai
menurut aksi**, dan hanya untuk baris yang belum pernah disimpan.

**5. `SetIndexAdjustmentList` juga menyetel `.IsCheck = true`** (baris 328–329) — penanda
dipilih-untuk-diklaim (AC 39) lahir di sini, bukan di layar. Pewarisan delapan kolom
(`.AdjustmentList(<LAST>).X = .AdjustmentList(1).X`, baris 570–744) cocok **persis** dengan teks
tiket, termasuk **tanpa** `STS_REJECT`. Langkah 2 menyetel `.IndexPremiumList` — subscript peserta,
yang di skema relasional digantikan `PREMIUM_LIST_DETAIL_ID` (penyimpangan sadar 2, tetap).

**6. Field bank pada layar adjustment ada ENAM, tiket menulis tiga.** `[terverifikasi]` nama
properti yang dirujuk `AdjustmentDetail_Section`: `NAMEOFBANK`, `IDOFBANK`, `ACCOUNTNO`,
`BRANCHOFBANK`, `SWIFTCODE` (huruf besar — kelas adjustment), di samping pasangan CamelCase
(`NameOfBank`, `IDOfBank`, `NoAccount`, `BranchOfBank`, `SwiftCode`, `PayableTo`) yang tampak
halaman master bank; salah satu gerbangnya `.SwiftCode != ''`. `T_CLAIMLF_ADJUSTMENT` hari ini punya
tiga: `NAME_OF_BANK`, `ID_BANK`, `ACCOUNT_NO`. ⛔ Menambah `BRANCH_OF_BANK`, `SWIFT_CODE`,
`PAYABLE_TO` menuntut **langkah migrasi `011`**, yang **tidak** termasuk langkah yang brief modul §4
izinkan sesi ini. Karena itu **tidak dikerjakan**; `[terbuka — work owner §9: izin langkah migrasi 011]`.

**7. Sumber nilai per lapisan** `[terverifikasi]` dari `pyStepsObjectName` tiap langkah:
`EM_PERCENT`, `AGE`, `PERIOD_MM`, `BEGIN_DATE`, `GROSS_VALUATION_BEGIN_DATE` dibaca dari **peserta**
(`Int-LIFE_PREMIUM_DETAIL`, langkah 8.1); `CURRENCY` dan `CLAIM_GROSS` dibaca dari **baris
adjustment** (langkah 8.2, objek `.AdjustmentList`). Kolom kita untuk `CLAIM_GROSS` bernama
`CLAIM_AMOUNT` (tiket 14) — satu nilai, dua nama; dicatat supaya tidak lahir kolom ketiga.
⚠️ `[sementara — menunggu OQ-N11 pemilik ekspor]` — dipertahankan work owner 29-09-2026 (OQ-N12 (a), bab bertanggal
GILIRAN-16 di akhir tiket).

**8. Kolom tabel dokumen — sensus butir `ad`, selesai.** `InsertDocument_Act` membuat halaman
`NewDocument`, mengisinya, memanggil `InsertGoogleStorage_Act`, lalu `Obj-Save`. Properti yang
diisinya `[terverifikasi]` baris 647–940: `.ID` (`@CurrentDate("yyyyMMddhhmmssSSS","Asia/Jakarta")`),
`.TANGGAL` (`@CurrentDateTime()`), `.IDPEGA`, `.NAMAFILE`, `.MIME` (`@toLowerCase`), `.KATEGORI_1`,
`.KATEGORI_2`, `.NOAKSEP`, `.NOPREKAS`, `.PAYMENTDATE`, `.pxCreateOperator`; `.T_STORAGE_ID` diisi
activity storage dan menjadi precondition `Obj-Save`. `Section/DocumentLife.xml` menampilkan
`.NAMAFILE` dan `.KATEGORI_2` atas daftar `.DocumentClaimList`, dengan pratinjau bergerbang `.MIME`
dan `.T_STORAGE_ID`.

⛔ Sensus `.DocumentList` yang tiket minta memang **tidak mungkin** menghasilkan kolom: `.DocumentList`
adalah daftar **lampiran bawaan Pega** (propertinya `.pyCategory`), dan itulah sebabnya tiket 14
menandainya. Sumber kolom yang benar adalah activity di atas. Teks AC-nya dibiarkan apa adanya
karena ia menyebut "activity di atas" dan mensyaratkan keputusannya **dicatat** — inilah catatannya.

### Cacat nyata di rule Pega yang TIDAK ditiru

`local.Rate` **tidak pernah dikosongkan** sebelum putaran pemilihan rate (penulisnya hanya langkah
8.2.1.9.1 dan 8.2.1.9.2), sedangkan langkah 8.2.1.9.2/9.3 menimpanya dengan nilai yang **sudah
dibagi 1000** (baris 4448 dan 4716). Bila untuk sebuah baris retro **tidak ada** baris `RATE_LIFE`
yang cocok, Pega memakai rate baris retro sebelumnya — yang sudah dibagi 1000 — lalu membaginya 1000
**sekali lagi**. Itu tidak mungkin disengaja. Go **menggagalkan** perhitungan secara terang bila
tidak ada rate yang cocok. `[penyimpangan sadar — dilaporkan ke work owner]`

Putaran pemilihan rate juga **tidak berhenti** pada baris pertama yang cocok (nol transisi keluar):
baris cocok **terakhir** yang menang, sedangkan `GetRateRetro` tidak memakai `ORDER BY`. Dengan 2.693
kombinasi ganda di `RATE_LIFE` `[data DBA]`, "terakhir" tidak tertentu. Go mengurutkan dan
**melaporkan** ambiguitas, tidak memilih diam-diam (butir **ah**).

### Pertanyaan yang XML tidak jawab

| # | Pertanyaan | Pemilik |
| ---: | --- | --- |
| 1 | Isi daftar kategori wajib (`GetCategoryLife_SQL` tidak ada di korpus) | DBA / work owner |
| 2 | Kolom mana pada tabel dokumen relasional yang memegang kategori pembanding — `KATEGORI_1` atau `KATEGORI_2`. `DocumentLife` menampilkan `.KATEGORI_2`; gerbang kelengkapan memakai `.pyCategory` dari daftar lampiran Pega, bukan dari kolom ini | work owner |
| 3 | Sumber `OUTWARDRATEID` **dan** daftar treaty-year tanpa membaca JSON produk | butir **ag**, Master Product Name Life |
| 4 | Kolom `IDR`/`USD` pada objek `treatyyear_life` yang dilihat Pega | `[data DBA]` |
| 5 | `TempInputData.SEX` tidak pernah disetel di activity ini; pemanggilnya yang mengisi | work owner |
| 6 | Mata uang selain `IDR`/`USD`: XML tidak punya cabang sama sekali | Product + UW |

---

## Implementasi — 26 September 2026 malam (tiket 03, sesi modul)

**Status: `claimed`** — **11 dari 28 AC tertutup** *(14 sebelum `/code-review`; tiga dibatalkan centangnya karena teksnya menuntut PENYIMPANAN, dan tidak ada jalur yang menyimpan)*. Jumlah AC bertambah dari 26 menjadi 28: dua AC
**baru** lahir dari pembacaan ulang XML (kaskade kapasitas dan pemilihan rate), dan tiga AC lama
**diralat teksnya** di tempat. Titik tetap `4c10059`.

Verifikasi penuh sesudah perbaikan review: `go vet` bersih · `go vet -tags=db` bersih ·
`gofmt -l` nol · `go build` · **144 PASS · 0 FAIL** *(dari 112)* · **24 SKIP** bertag `db` ·
`tsc --noEmit` · **5** test JS · `vite build` **88 modul**.

### Keputusan §5 yang dipakai, dan yang ditunggu

Dikerjakan: **ad** *(kolom dokumen — sensusnya selesai; keputusannya sendiri masih `[USULAN]` di
brief §5)*. Dipakai: **XML menang** *(aturan baru; dipakai lima kali di bawah)*. ⛔ **ab** TIDAK
dipakai — nol baris commit ini menyentuh `AUTH_STUB`. Ditunggu: **ag** *(sumber
`OUTWARDRATEID` **dan** daftar treaty-year)* — inilah yang menahan setengah tiket ini; **ah**
*(aturan bila rate ganda)*; **o1–o3** *(penomoran)*; izin langkah migrasi `011` *(field bank)*.

### ⭐ Lima ralat menurut XML — teks AC disunting di tempat (§1.2-b)

| # | Teks lama | Teks baru | Bukti |
| ---: | --- | --- | --- |
| 1 | *"Sebagai **ReasLifeAdmin**, saya dapat … menyimpannya ke Outstanding"* | *"Sebagai **ReasLifeSPV** …"* | `AdjustmentDetail_Section.xml` 16249, 16396, gerbang 16468 |
| 2 | `AMOUNT = RetrocadedShare × PERCENT_SHARE ÷ 100`; `GROSS = RATE × (1+EM_PERCENT) × AMOUNT` | rumus yang sama **ditambah** `RATE` tersimpan dibulatkan 4, rate pakai `÷1000` dibulatkan 10, `PERCENT_SHARE ÷ 100` dibulatkan 4, dan `EM_PERCENT` sebagai pecahan langsung | `SpreadingClaimLife_Act` 4396–4509 |
| 3 | `[terbuka]` *"mana yang berlaku tidak terbaca. **Jangan tebak.**"* | `[ditutup oleh XML]` — **tahun polis** yang memilih: tahun ke-1 memakai `GROSS − Discount − Comm`, selain itu `GROSS − Comm` | 4372–4599 *("Bukan tahun pertama", `local.Year==1` WhenTrue=3)* dan 4640–4907 *("Tahun pertama", WhenFalse=3)* |
| 4 | *(tidak ada)* | **AC baru**: kaskade kapasitas `IDR`/`USD` per treaty-year, sisa mengalir | 3030–3729 |
| 5 | *(tidak ada)* | **AC baru**: pemilihan rate — `GENDER` memuat `U` cocokkan umur saja, selain itu umur **dan** jenis kelamin; `CONTRACT` tidak ikut | 4084–4284 |

Ralat 3 menutup satu butir `[terbuka]` yang sudah berdiri sejak tiket 14 dan yang ronde-ronde
sebelumnya tandai *"tidak terbaca"*. Ia **terbaca** — yang belum terbaca dulu adalah **cara membaca
aksi precondition Pega**. Kuncinya ditemukan dari sepasang langkah yang saling melengkapi *(bab
Pembacaan ulang XML)*, bukan dari dokumentasi: `2` = lanjut, `3` = lewati.

### Yang dibangun

| Berkas | Isi |
| --- | --- |
| `services/spreading.go` *(baru)* | `HitungSpreading` murni: kaskade kapasitas, pemilihan rate, kedua cabang NET, seluruh pembulatan XML sebagai konstanta bernama. `TahunPolis`, `PilihRate` |
| `services/dokumen.go` *(baru)* | `PeriksaDokumenAda` *(gerbang 1)*, `PeriksaDokumenLengkap` *(gerbang 2)*, `KategoriBerbeda`, `SumberKategoriWajib` + `KategoriWajibBelumDiketahui` yang gagal terang |
| `services/adjustment.go` | `TambahBaris` *(pewarisan 8 kolom + `IsCheck`)*, `TandaiOutstanding` *(`STS_REJECT=0` bergerbang `PrintFaceClaim`)*, `PeranSimpanOutstanding` |
| `models/pohonklaim.go`, `models/klaimlife.go` | tiga ralat tipe *(di bawah)*, `Dokumen` bertujuh kolom isi, `Peserta.Dokumen`, konstanta `KodeOutstanding`/`KodeAksep`/`KodeDitolak` |
| `migrations/010_kolom_t_claimlf_document.sql` *(+`_down`)* | tujuh kolom isi `T_CLAIMLF_DOCUMENT` *(butir **ad**)* |
| `repository/migrasi.go`, `strukturkolom_test.go` | `KolomAlterTambah` — kolom yang lahir lewat `ALTER` kini terlihat penjaga bentuk |

### ⛔ Tiga ralat tipe di `models` — XML membetulkan tiket 14

1. **`Spreading.RetrocadedShare` bertipe `Ratio`, seharusnya `Money`.** XML mengisinya dengan sisa
   `CLAIM_GROSS` atau dengan kapasitas `IDR`/`USD` treaty-year *(langkah 8.2.1.4–7)*; keduanya uang.
   Ini persis jebakan yang bab Catatan tiket ini sendiri peringatkan: *nama berakhiran `_SHARE` di
   korpus ini menipu*.
2. **`SpreadingRetro.Commision` dan `.OvrComm` bertipe `Money`, seharusnya `Ratio`.** XML membagi
   keduanya seratus sebelum memakainya. Menyimpan persen sebagai `Money` membuat persen dapat
   dijumlahkan dengan uang — justru yang ADR-F-0004 larang.
3. Keduanya lolos selama ini karena **tidak ada satu pun test yang menghitung dengan nilainya**;
   yang ada hanya test pulang-pergi penyimpanan, dan pulang-pergi tidak peduli artinya apa.

### Penjaga baru, masing-masing DIBUKTIKAN dapat gagal

Tiap penjaga dirusak sengaja, kegagalannya dilihat, lalu dipulihkan:

| Cacat yang dipasang | Yang menangkap |
| --- | --- |
| pembagian per mil dihapus | 3 test rumus |
| cabang tahun polis ditukar | 3 test rumus |
| sisa klaim tidak mengalir ke treaty berikutnya | `TestKaskade…`, `TestTreatyKedua…` |
| pesan gerbang tidak menyebut peserta mana | `TestPesertaTanpaDokumenDitolakDenganNomorUrut` |
| dedup kategori dimatikan | `TestKelengkapanMembandingkanCacahKategoriBerbeda` *(uji dan fungsinya dibuang 28-09-2026 — butir bl)* |
| status ikut diwarisi | `TestBarisBaru…`, `TestBarisKedua…` |
| gerbang `PrintFaceClaim` dilepas | `TestTandaiOutstandingHanyaMenyentuhBarisTanpaStatus` |

⚠️ Percobaan **pertama** memasang cacat ketujuh **tidak cocok polanya** dan karena itu tidak
mengubah apa pun — dan test tetap hijau. Hijau yang berarti *"cacatnya tidak pernah terpasang"*,
bukan *"penjaganya lemah"*. Percobaan kedua memakai pola tanpa backslash dan penjaga itu menggigit.
Dicatat karena inilah cara sebuah pembuktian penjaga dapat berbohong kepada penulisnya sendiri.

### Dua cacat rule Pega yang TIDAK ditiru — dilaporkan, bukan didiamkan

1. **`local.Rate` tidak pernah dikosongkan.** Bila nol baris `RATE_LIFE` cocok untuk sebuah baris
   retro, Pega memakai rate baris sebelumnya — yang **sudah** dibagi seribu — lalu membaginya
   seribu **lagi**. Go menggagalkan perhitungan secara terang.
2. **Putaran rate tidak berhenti pada yang pertama cocok** dan `GetRateRetro` tanpa `ORDER BY`:
   yang menang tidak tertentu. Go melaporkan ambiguitas beserta `ID` barisnya *(butir **ah**)*.

Keduanya `[penyimpangan sadar — menunggu work owner]`.

### ⭐ Butir `ad` selesai — dua cacah dari sumber yang berbeda

Aturan sensus CLAUDE.md §4a dipenuhi: cacah **pertama** dari katalog `POOLDATA.DOCUMENT_CLAIM`
*(14 kolom)*, cacah **kedua** dari seluruh `Property-Set` `InsertDocument_Act` sebelum `Obj-Save`
*(12 properti, baris 647–940)*. Keduanya cocok. Tujuh kolom diadopsi; tujuh kolom warisan sengaja
tidak ikut, dengan sebabnya masing-masing, di `STRUKTUR-TABEL-CLAIM-LIFE.md`.

Sensus `.DocumentList` yang brief minta memang tidak dapat menghasilkan kolom: kelasnya lampiran
bawaan Pega. Itu dilaporkan, bukan dibulatkan menjadi "nol kolom".

### ⭐ Aturan "dokumen lengkap" — mekanismenya kini terbaca

Gerbang kedua membandingkan **cacah kategori berbeda yang terunggah** dengan **cacah baris
kategori wajib**. Perbandingan cacah itu ditiru **apa adanya**, termasuk kelemahannya *(dua dokumen
berkategori salah dengan cacah yang kebetulan pas akan lolos)*: memperbaikinya diam-diam berarti
sistem baru menolak klaim yang sistem lama terima, tanpa seorang pun memutuskannya.

⛔ Yang tetap terbuka: **isi** daftar wajibnya — `GetCategoryLife_SQL` **tidak ada di korpus**,
seluruh 29 berkas `Claim Life/RDBList/` sudah dicacah. `KategoriWajibBelumDiketahui` menggagalkannya
secara terang, seperti `PenomorBelumDiputuskan` pada tiket 02.

### ⛔ Yang BELUM dikerjakan — 14 AC, beserta sebab masing-masing

| AC | Sebab |
| --- | --- |
| spreading benar-benar **disimpan** saat baris disimpan; nilai dibekukan; satu transaksi | menunggu **ag**: kunci rate `OUTWARDRATEID` **dan** daftar treaty-year keduanya berasal dari JSON produk yang AC 38 larang dibaca. Tanpa keduanya, jalur ujung-ke-ujung tidak dapat ditulis tanpa mengarang masukannya |
| pembaca Oracle `TREATYYEAR_LIFE` / `RETROCESSIONLIFE` / `RATE_LIFE`, tiruan `RATE_LIFE` di skema uji | sama; ditambah kolom `IDR`/`USD` yang `[data DBA]` |
| `AGE`, `SEX`, `PERIOD_MM` peserta sebagai masukan spreading | menuntut `kolomSalin` tiket 02 bertambah, beserta `TestUrutanKolomSalinDikunci` dan tiruan skema uji — perubahan lintas tiket yang tidak dikerjakan diam-diam |
| gerbang dokumen **terpasang** pada jalur simpan; dua pintu HTTP; layar React | aturannya sudah ada dan teruji sebagai fungsi murni; yang belum ada perkabelannya |
| dua peserta × dua putaran = empat baris tertelusur *(test db)*; `SELECT` dokumen biasa | menunggu Oracle *(G1)* |
| tiga field bank dapat diisi dari layar; boleh kosong saat simpan | layar belum ada. ⚠️ XML menunjukkan **enam** field bank, bukan tiga — menambah tiganya menuntut langkah migrasi `011` yang brief §4 tidak izinkan sesi ini |
| peserta menyimpan `STATUS` dan `RECOMMENDATION` | keduanya belum ada di daftar kolom peserta; `SOURCE_ID` dan `CEDING_RETENTION` sudah |

### Hasil `/code-review` atas titik tetap `4c10059` — dan apa yang berubah karenanya

Dua sub-agen berjalan paralel, satu sumbu Standards dan satu sumbu Spec. **Delapan temuan
diterima dan diperbaiki, satu ditolak dengan bukti.** Jumlah AC tercentang **turun dari 14 menjadi
11** karena review ini — turun, bukan naik, dan itu memang gunanya.

| # | Sumbu | Temuan | Tindakan |
| ---: | --- | --- | --- |
| 1 | Standards | **ADR-U-0016 Akibat 3 dilanggar**: *"lapisan `services` membulatkan pada desimal kedelapan"*. `PREMIUM_SPREADED_GROSS` keluar dengan **sembilan** desimal, dan `NUMBER(38,8)` membuat **Oracle** yang membulatkannya diam-diam — angka yang dihitung berbeda dari angka yang tersimpan, tanpa satu pun galat | ✅ `desimalUang = 8` dipasang pada `share`, `AMOUNT`, `GROSS`, `NET`. GROSS dibulatkan **sebelum** NET diturunkan darinya, supaya `NET = GROSS − Comm` berlaku pada angka yang benar-benar masuk basis data. Test baru dengan rate `2,4681` yang **sengaja** menghasilkan sembilan desimal |
| 2 | Spec | **`WarisiKolom` tidak mewarisi apa pun** selain mata uang, padahal AC 5 menjanjikan delapan kolom — dan testnya tidak mungkin gagal pada tujuh di antaranya | ✅ Cacat **lebih dalam** dari yang dilaporkan: `models.BarisAdjustment` **tidak punya** ketujuh medan itu sama sekali, sedangkan DDL `004` menyediakan kolomnya. INSERT pun hanya menulis `CURRENCY`. Jadi ketujuh kolom itu **selalu NULL**. Medan ditambahkan, `WarisiKolom` menyalin kedelapannya, INSERT dan SELECT diperluas, dan test kini memberi tiap kolom nilai **berbeda** — dibuktikan gagal saat satu kolom mana pun dilepas |
| 3 | Spec | **Empat AC tercentang berbunyi "disimpan"/"tersimpan"** padahal tidak ada satu pun pemanggil non-test | ✅ Tiga dibatalkan centangnya. Yang tetap tercentang adalah yang benar-benar tentang **asal nilai**, bukan tentang penyimpanan |
| 4 | Spec | **Pesan gerbang kedua dikarang**: XML hanya menulis *"Documents are incomplete, please complete the documents"*, tanpa nomor peserta; Go menambahkan `person number` | ✅ Pesan dikembalikan **persis** seperti XML. Nomor peserta dipindah ke `PesertaDokumenTidakLengkap` yang terpisah, supaya layar tetap dapat menunjukkannya tanpa kami menaruh kalimat yang tidak pernah ada di mulut sistem lama |
| 5 | Spec | **Urutan kaskade tidak ditegakkan**: tiket menyatakan urutannya ditentukan `order by TO_NUMBER(IDR) asc`, kode hanya menitipkannya pada pemanggil yang belum ada | ✅ `HitungSpreading` kini **memeriksa** urutan menaik dan menolak yang tidak terurut. Diperiksa, bukan diurutkan diam-diam: mengurutkan sendiri akan menyembunyikan pembaca yang lupa `ORDER BY` |
| 6 | Standards | **`Dokumen.Tanggal`/`PaymentDate` bertipe `time.Time` biasa** untuk kolom nullable — ADR-U-0022 Akibat 2: *"kolom angka atau tanggal menjadi KOSONG, bukan nol dan bukan tanggal nol"* | ✅ Keduanya menjadi `*time.Time`. ⚠️ Medan tanggal lain di model punya masalah sama dan **tidak** diubah: itu keputusan sekali untuk seluruh model, milik tiket 14, dan ditandai `[terbuka]` di tempatnya |
| 7 | Standards | **Komentar DDL `005` dan `006` kini berbohong**: `005` masih menyebut `RETROCADED_SHARE` "bukan uang", `006` masih menandai `PREMIUM_SPREADED_NET` `[terbuka]` | ✅ Keduanya diralat di tempat, dengan catatan bahwa **nol langkah migrasi baru** diperlukan — tipe kolomnya sudah benar |
| 8 | Spec | **Bukti sensus saling bertentangan**: header `010` menulis *"tidak ada daftar kolom eksplisit untuk disensus"*, sedangkan tiket menulis sensus 12 properti | ✅ Header `010` ditulis ulang. Keduanya sekarang berkata sama: tidak ada daftar kolom pada SQL, tetapi **ada** daftar Property-Set, dan itulah cacah keduanya |
| 9 | Spec | *"Gerbang kedua melewatkan penjaga `@LengthOfPageList(Category2.pxResults)>0`, sehingga peserta nol dokumen di-skip Pega tetapi ditolak Go"* | ⛔ **DITOLAK.** Precondition itu ada di baris **6782**, sedangkan langkah 12.2.3 mulai di **6704** dan langkah 12.2.4 di **6831**. Jadi ia milik langkah **dedup Java**, bukan langkah pesan — langkah pesan punya precondition sendiri di baris **6939** (`Category2 = Category`). Peserta nol dokumen tetap dibandingkan `0 ≠ N` dan tetap ditolak. Perilaku Go sudah cocok |

⚠️ Satu koreksi pada bab ini sendiri: klaim *"keputusan **ab** dipakai"* di atas **salah** — nol
baris di commit ini menyentuh `AUTH_STUB`. Dan **ad** masih `[USULAN]` di brief §5; yang benar
adalah *tiket ini mengerjakan pekerjaan yang brief tugaskan padanya*, bukan bahwa keputusannya
sudah disahkan. `PeranSimpanOutstanding` juga belum punya pemanggil non-test — ia sengaja mendahului
tiket 07, yang akan memakainya; tanpa itu fakta XML-nya hilang.

## Implementasi — 27 September 2026 (A2, penutupan stub ar1)

**Satu centang bergeser, dan sebabnya:**

| AC | Sebab bergeser |
| --- | --- |
| baris adjustment menyimpan nama bank, id bank, nomor rekening | migrasi 011 menambah `BRANCH_OF_BANK`, `SWIFT_CODE`, `PAYABLE_TO` di samping ketiganya — kolomnya kini ada, dan penulisnya mengisinya |

**Yang TIDAK bergeser, dan sebabnya:** sisa AC tiket ini *(satu peserta per baris, empat baris untuk
dua peserta × dua putaran, tanggal per peserta, satu transaksi, gerbang bank pada Save Adjustment)*
menuntut **jalan Oracle sungguhan**; bentuknya sudah ada, perilakunya belum dijalankan.

**Gerbang "dokumen lengkap"** *(butir **ar1**)*: daftar kategori wajib dibaca saat jalan dengan
filter `BISNIS IN ('ALL', <kode bisnis>)`, dibandingkan **cacah** seperti gerbang XML.

⚠️ Keputusannya `[USULAN yang disahkan]`, **bukan** `[terverifikasi]`: rule `GetCategoryLife_SQL`
**tidak ada di korpus**, jadi mana yang Pega baca tidak terverifikasi. Kandidat kedua — view
`DOCUMENTCLAIM_LIFE` dari `M_PRODUCT_LIFE.JSONDATA` — dicatat dan **tidak** dipakai, sebab jalur
JSON produk dilarang AC 38.

## Ralat menurut XML — 27 September 2026 (`DeletePesertaClaimLife` dibaca sebagai pohon)

**Yang diralat:** brief lanjutan 10 §1 menyebut `DeletePesertaClaimLife` sebagai penghapus
peserta, lengkap dengan rute `DELETE /api/klaim-life/{id}/peserta/{pesertaId}` bergerbang tahap
dan pemegang. **KELIRU — activity itu tidak menghapus apa pun.**

### Apa yang benar-benar dilakukannya

`Activity/DeletePesertaClaimLife.xml`, dibaca sebagai pohon langkah utuh:

| Baris | Langkah | Isi |
| ---: | --- | --- |
| 225 | `Property-Set` pada `pyWorkPage.ClaimData.PremiumListSummary.PremiumListDetail` | `Local.IndexPremium = .pxListSubscript` |
| 583 | — | loop `pyStepsRepeatDefHasRepeat = EMBEDDED` atas daftar peserta |
| 314 | `Property-Set` pada `.AdjustmentList` | `.IndexPremiumList = Local.IndexPremium` |
| 417 | — | loop `EMBEDDED` kedua, atas baris adjustment peserta itu |
| 441 | `Obj-Save` pada `pyWorkPage` | menyimpan seluruh kasus |

Ia **mengindeks ulang penunjuk balik** setiap baris adjustment ke posisi peserta pemiliknya,
lalu menyimpan. Namanya menyesatkan.

⚠️ `pyStepsRepeatDefHasRepeat` bernilai **`EMBEDDED`**, bukan `true`. Pembacaan pertama mencari
`true` dan karena itu **kedua loopnya sempat tidak terlihat** — tanpa loop itu, activity ini
terbaca seolah hanya menyentuh satu baris.

### Siapa yang menghapus, dan kapan tombolnya ada

Penghapusan barisnya dikerjakan **klien**. Tombolnya di `Section/InputOSClaimLife.xml` — layar
**Outstanding**, bukan Register:

| Butir | Nilai | Baris |
| --- | --- | ---: |
| Label | `DELETE` | 17865 |
| Aksi 1 | `pyAction = deleteRow` *(hapus baris grid di klien)* | 18017 |
| Konfirmasi | `pyNextGenGridDeleteConfirm = false` — **tanpa** popup | 18021 |
| Aksi 2 | `pyAction = refresh` → `pyActivity = DeletePesertaClaimLife` | 18032, 18039 |
| Syarat tampil | `pyCondition = pyWorkPage.ClaimData.PremiumListSummary.CLAIM_NO == ''` | 18082 |

⛔ **Tombolnya hanya ada selama klaim belum bernomor.** Sesudah `CLAIM_NO` terisi, kolom DELETE
itu tidak tampil sama sekali — jadi tidak ada "hapus peserta pada klaim yang sudah terdaftar"
untuk ditiru, dan rute `DELETE` bergerbang tahap+pemegang itu **tidak dibangun**.

⚠️ `pyCondition` berdiri di `<pyUserData>` milik **sel**, sedangkan `pyModes` yang memuat
tombolnya berakhir di b17938. Jendela baca **maju** dari label `DELETE` berhenti sebelum b18082
dan melewatkan syarat tampilnya — pengulangan persis kekeliruan butir av. Ia ditemukan dengan
**menaiki** pohon dari `pyActivity` ke blok pembungkusnya.

### `SelectAllClaimLife_act` — pilih/lepas SEMUA, juga di layar Outstanding

`Activity/SelectAllClaimLife_act.xml`: b247 `Select.CARI1 = @if(Select.CARI1=="","true",
@if(Select.CARI1=="true","false","true"))` — **penjungkit tiga keadaan**: kosong dianggap belum
pernah dipakai dan menjadi `"true"`. Lalu loop `EMBEDDED` b550 menyetel `.IsAccept = Select.CARI1`
*(b429)* pada setiap baris. Pemanggilnya `Section/InputOSClaimLife.xml` b16633, b16710.

**AC:** tidak ada AC tiket ini yang berubah centangnya. Yang diralat adalah **siapa** yang
menghapus *(klien, bukan activity)*, **di layar mana** *(Outstanding, bukan Register)*, dan
**kapan** *(hanya selama `CLAIM_NO` kosong)*.

## Pembacaan ulang XML — 27 September 2026 (`SaveOutStandingLife_Act` langkah 23: enam total peserta)

`SaveOutStandingLife_Act` adalah rule milik tiket ini *(tombol `Save to RNM`, `InputOSClaimLife.xml`
b21102 → b21126)*. Langkah **23**-nya tidak pernah dibaca sampai selesai di ronde mana pun, dan di
situ ada perilaku yang tiket ini belum sebut.

### Yang langkah 23 lakukan

| Langkah | Baris | Isi |
| --- | ---: | --- |
| 23 | b10638 | `Property-Set` pada `pyWorkPage.ClaimData.PremiumListSummary.PremiumListDetail` *(= **peserta**)*, **ULANG** `EMBEDDED`. Reset enam penampung `local.Total*` ke literal **`0`** b10663–b10795 |
| 23.1 | b10841 | **ULANG** `EMBEDDED` b11046 atas **`.AdjustmentList`**, **prasyarat kosong**. `local.TotalCedingRetention = .CEDING_RETENTION + local.TotalCedingRetention` b10860; lalu `.SHARE_NUSANTARA_RE` b10906, `.SUM_INSURED` b10926, `.SUM_REASURED` b10946, `.SHARE_RETRO` b10966, `.CLAIM_AMOUNT` b10986 |
| 23.2 | b11067 | Tanpa ulang. Keenam total ditulis ke halaman **peserta**: `.TotalCedingRetention` b11091, `.TotalShareRNM` b11137, `.TotalSumInsured` b11157, `.TotalSumReasured` b11177, `.TotalShareRetro` b11197, `.TotalClaimAmount` b11217 |

Rumus yang sama persis ada di `SavePesertaClaim.xml` langkah 8 b4002 / 8.1 b4221 / 8.2 b4592 —
jalur `Submit` pada Register. **Dua tombol, satu rumus.**

### Tiga hal yang ini tegaskan, dan satu yang ia ralat

1. ⛔ **Tidak ada penyaringan.** Langkah 23.1 b10841 dan 8.1 b4221 sama-sama **berprasyarat kosong**
   *(`<pyStepsPreCondParamsWhen/>`, b4562 pada yang kedua)*. Baris ber-`STS_REJECT = 2` **ikut
   dijumlah**, dan baris yang `IsCheck`-nya dicabut pun ikut. Dugaan wajar *("tentu yang ditolak
   tidak ikut")* justru yang salah.
2. ⛔ **Penampungnya mulai dari literal `0`**, bukan kosong. Peserta **tanpa** baris adjustment
   karena itu bertotal **`0`** di Pega, bukan kosong — dan itu beda yang terlihat di layar.
3. ⚠️ **Totalnya ENAM.** `Total Ceding Retention` ikut dihitung *(b10860, b11091)* dan tampil di
   `ClaimLifeDetailGCNM.xml` b20629. Ia sempat luput dari dokumen kami karena pencacahannya memakai
   rujukan `CheckTotalAdjustmentClaim`, dan hanya total itu yang **tidak** punya aksi refresh.
4. ⛔ **RALAT atas OQ-H** *(dokumen `OQ-untuk-tim.md`, blok "Ralat kedua")*: kesimpulan bahwa keenam
   angka itu "tidak dapat ditiru" **dicabut**. Yang hilang dari ekspor hanya pemanggil **refresh**.

### Yang berubah di kode

Keenam total **dihitung saat Detail dibaca** *(`models.HitungTotalPeserta`, dirakit
`services.KlaimLife.Ambil`)* dan **tidak disimpan** — nol kolom `TOTAL_*` ditambahkan, dan ada
penjaga statik yang menolak migrasi yang menambahkannya.

⚠️ **Penyimpangan sadar:** Pega **menyimpan** hasilnya ke halaman peserta saat `Save to RNM`,
sehingga angkanya dapat **basi** sesudah putaran atau akseptasi sampai tombol itu ditekan lagi. Di
sini ia tidak pernah basi. Ditanyakan ke pemilik ekspor *(OQ-H versi sempit)* kalau-kalau ada
laporan yang justru mengandalkan angka tersimpan.

**AC:** tidak ada AC tiket ini yang berubah centangnya. Yang bertambah adalah **perilaku langkah 23**
yang sebelumnya tidak tercatat di tiket mana pun.

## Implementasi — 28 September 2026 (GILIRAN-11 paket 1: `Save to RNM`)

`Activity/SaveOutStandingLife_Act.xml` dibaca **utuh sebagai pohon**: 29 langkah teratas, 58 langkah
berikut anaknya, ±14.800 baris pecahan. Pohonnya dirakit dengan pengurai XML yang menyimpan nomor
baris (bukan grep), supaya medan milik langkah INDUK yang tertulis SESUDAH anak-anaknya tidak
terbaca sebagai milik anak terakhir. Petanya — langkah → padanan — ada di kepala
`APP_RNM/internal/services/simpanrnm.go`.

Yang dibangun: `POST /api/klaim-life/{id}/outstanding` (`handlers/simpanrnm.go` →
`SimpanRNM.Simpan` → `PeriksaSimpanRNM`, murni), bergerbang tahap **Outstanding Claim** +
pemegangnya; seluruh gerbang XML diperiksa SEBELUM satu tulisan pun; satu transaksi (nomor bila
`CLAIM_NO` kosong, `STS_REJECT=0` bagi baris tanpa status, jejak); tabel warisan **hanya dibaca**
(klaim ganda); Arasapas sesudah commit. Tombol `Save to RNM` (b21102) di `OutstandingClaimLife.tsx`.

### ⛔ Cara membaca yang menentukan — dari korpus

Baris `WHEN` hanya berlaku bila langkahnya ber-`pyStepsPreCondition=true`, baris `TRANS` hanya bila
`pyStepsTransition=true`. Langkah 11.1 (`BusinessCode` L1–L11), 22.1.1 (`.IsCheck=="'true'"`), dan
22.1.3 (`.PrintFaceClaim==1`) ber-WHEN **tanpa** bendera — WHEN-nya mati, dan langkahnya selalu
berjalan. Langkah 15 ber-TRANS tanpa bendera pula. Kode aksi `6` = keluar activity
(`claim-prop/grilling-ronde-2.md`).

### ⛔ Ralat bertanggal atas bab-bab di atas — 28-09-2026

1. **"`STS_REJECT = 0` ditulis per baris adjustment dan bergerbang `PrintFaceClaim`"** (pembacaan
   26-09 butir 4) — **keliru separuh.** 22.1.3.1 (insert warisan) memang bergerbang
   `.PrintFaceClaim==""`, tetapi **22.1.3.2 tidak punya precondition sama sekali**: ia menyetel
   `.PrintFaceClaim=1`, `.STS_REJECT=0`, dan `.ADJUSTMENT_DATE` pada **setiap** baris, setiap kali.
   Di Pega itu aman hanya karena bendera `pyWorkPage.Save=1` mematikan tombolnya sesudah simpan
   pertama (b21095). Aplikasi ini tanpa kolom bendera (**OQ-N1**), jadi hurufnya **tidak** ditiru:
   hanya baris **tanpa status** yang ditulis `0` — menulis semuanya berarti membatalkan penolakan
   Admin diam-diam bila tombolnya ditekan lagi.
2. **"Peran yang boleh Save to Outstanding adalah `ReasLifeSPV`"** — **benar untuk satu dari dua
   pemanggil.** `SaveOutStandingLife_Act` dipanggil DUA tombol: dua di `AdjustmentDetail_Section`
   (b16249/b16396, gerbang `pyPosition=='ReasLifeSPV'` b16468 — jalur SPV menambah putaran, tiket
   11) **dan** `Save to RNM` di layar Outstanding `InputOSClaimLife` b21102 → b21126, yang dipegang
   **Admin**. Rute ini meniru yang kedua (brief GILIRAN-11 §2 butir 1): gerbangnya pemegang tahap
   Outstanding, `WajibPemegangTahap`. Bab "Hasil & nilai pengguna" karena itu hanya separuh.
3. ~~**Gerbang dokumen LENGKAP (langkah 12) menyaring SELURUH peserta**~~ — **DICABUT** (ralat
   bertanggal di bawah): langkah 12 ter-remark, jadi tidak ada gerbang dokumen lengkap sama sekali.
   Gerbang dokumen ADA (langkah 3) memang menyaring `.IsAccept=="true"` (b1181; `IS_CHECK` padanan
   terdekat).
4. **DOL di Save BERBEDA dengan `ValidasiDOL_Act`**: langkah 11.8 memakai jendela retro **tanpa**
   pergeseran satu hari (b4464 `@addCalendar(.DATE_OF_LOSS,0,0,0,0,0,0,0)`), sedangkan
   `ValidasiDOL_Act` menggeser (b698). Keduanya ditiru apa adanya; uji
   `TestSimpanRNMDOLJendelaTanpaGeserRetro` menjaga bedanya.
5. **Langkah 5 residu**: `@contains(.Protect,"1")` → pesan "Claim gross tidak boleh lebih besar dari
   Share Nusantara Re" — `.Protect` **nol penulis** di seluruh korpus Claim Life, jadi pesan itu tidak
   pernah muncul (OQ-N4).
6. **Klaim ganda (11.2–11.6) membaca `OS_AKSEPTASI_KLAIM_LIFE`** dengan nama dan tanggal lahir
   tertanggung — keduanya tidak disalin ke tabel klaim, jadi dicocokkan di SQL dengan baris sumber
   (`repository/gandawarisan.go`); baris warisan milik aplikasi ini tidak pernah cocok (OQ-N2).
7. **Langkah 16–20 (nomor)** bergerbang `CLAIM_NO==""` di setiap langkahnya (13–15, jalur
   `Generate_NoKlaim_Life*`, ter-remark); aplikasi menomori saat
   pendaftaran, jadi cabang ini hanya berjalan bagi klaim tanpa nomor, memakai penomor yang sama.
8. **Layar Outstanding** menyalin tombol perpindahannya sendiri tanpa konfirmasi
   `Send Back to Admin?`; kini memakai `PanelPindahTahap` (satu daftar, satu dialog).

### Pertanyaan terbuka yang lahir

OQ-N1 (bendera `Save`), OQ-N2 (klaim ganda antarklaim baru, cacat SQL health), OQ-N3 (gerbang retro
langkah 27 lawan OQ-064 Komite), OQ-N4 (residu `.Protect`, `ADJUSTMENT_DATE`/`PrintFaceClaim`,
`.IsAccept`) — `OQ-untuk-tim.md`.

### ⛔ Ralat bertanggal — 28-09-2026 (temuan /code-review GILIRAN-11)

Sebabnya satu: `pyStepsBlockName = //` berarti langkahnya **ter-remark** dan tidak pernah jalan
(`claim-prop/grilling-ronde-2.md` Aturan 2), dan pembaca pohon yang dipakai paket 1 tidak mencetak
medan itu. Delapan langkah activity ini ber-remark: 11.3 b3495, 11.9 b4632, 11.11 b5009, 12 b6178,
13 b7074, 14 b7293, 15 b7512, 23 b10649.

1. **Gerbang STNC (11.9/11.11) dan dokumen lengkap (12) DIBUANG.** Keduanya menolak simpan yang
   tidak pernah ditolak sistem lama. `MasukanRNM` tidak lagi membawa ambang `MAXDATARECEIVE`,
   `DateReceived`, atau daftar kategori; `Simpan` tidak lagi membaca ambang produk dan daftar
   kategori. Gerbangnya kini empat: dokumen ada, klaim ganda, DOL, medan kosong. Uji
   `TestSimpanRNMTreatyTanpaDokumenLolos` dan `TestSimpanRNMDokumenTidakLengkapLolos` menggantikan
   uji yang menuntut pesan langkah 12; uji STNC dibuang. Kelengkapan per kategori: **OQ-N6** —
   ditutup GILIRAN-12 (butir **bl**, di bawah).
2. **Langkah 23 (enam total) ter-remark** — tidak ada perubahan perilaku: totalnya dihitung saat
   baca dari `SavePesertaClaim` langkah 8, yang hidup (`models.HitungTotalPeserta`).
3. **Langkah 27 membaca salinan yang MUNGKIN DITUKAR.** `pyWorkPage.ClaimData.PolicyDataLife` diisi
   `InsertJsonClaimLife_Act` (langkah 25); langkah 2-nya (tidak ter-remark) menukar `RetroID` dan
   `SecurityReinsurerID` bila TP/TR, security reinsurer terisi, dan `ProdDateTime` < 7 Feb 2025.
   `ArasapasDilewatiRetro(PolisRetro)` menilai kedua kemungkinan; bila hasilnya berbeda, Arasapas
   **ditahan** dengan alasan terbaca (**OQ-N5**). Keluar di langkah 27 juga melewati Obj-Save di
   sistem lama (tambahan **OQ-N3**).
4. **Langkah 22.1.3.2 menulis baris adjustment SAJA.** Paket 1 memakai `PerbaruiStatusBaris`, yang
   ikut menulis `STS_REJECT` peserta serta mengosongkan `ACCEPTED_NO`/`ACCEPTATION_DATE` —
   bertentangan dengan sensus penulis `STS_REJECT` di `repository/klaimlife.go`. Kini
   `KlaimLife.TandaiBarisOutstanding` (`… SET STS_REJECT = :1 WHERE ID = :2 AND STS_REJECT IS NULL`,
   nol baris = galat); namanya ditambahkan ke keempat penjaga statik (wewenang, jejak, dua penjaga
   Komite), dan kedua penjaga pertama terbukti merah terhadap penulis liar. ⚠️ Cermin header
   (`CerminkanHeader`) tetap di transaksi simpan — selisih sadar dengan XML, yang menulis header
   hanya lewat Arasapas langkah 28 (`serviceInsertArasapasClaimLife_act` b371, b417); tambahan OQ-N3.
5. **Klaim yang tidak ada** dijawab 404 `klaim tidak ada`, bukan 500.
6. **Bukti `IsCheck == "true"` (AC 8)** dipindah dari b1812/b1997/b2413 — WHEN langkah 7.1
   (precondition false) serta 7.2/7.4 (`//`) `SavePesertaClaim` — ke b3631 dan b3919 (7.7, 7.8,
   hidup). Nilainya tidak berubah.

### ⛔ Butir bl — 28-09-2026 (GILIRAN-12 paket 1): gerbang dokumen lengkap DIBUANG

`[DIPUTUSKAN; veto work owner]` — OQ-N6 **ditutup**. Kelengkapan dokumen per kategori ter-remark di
XML (langkah 12, b6178), jadi di sistem lama ia tidak pernah berlaku. Mengikuti XML dan aturan "kode
mati dibuang": `services.PeriksaDokumenLengkap`, `PesertaDokumenTidakLengkap`, `KategoriBerbeda`, dan
`ErrDokumenTidakLengkap` dibuang beserta empat ujinya (`services/dokumen_test.go`). Bila bisnis
menghendaki gerbang itu, ia keputusan **baru**, bukan replikasi.

Yang tetap: gerbang "belum diunggah" (langkah 3–4) dan daftar kategori butir **ar1**, yang kini
hanya dipakai validasi unggahan. Bab 26-09 butir 3 (tabel dua gerbang) dan ralat #3 GILIRAN-11
tetap sebagai catatan pembacaan. Teks AC 45 disunting di tempat (lihat AC-nya).

## ⛔ Ralat bertanggal — 29 September 2026 (GILIRAN-13 paket 2, butir **bo**: `Add` membuat baris pertama)

**Celah yang ditutup.** Tiket ini dan tiket 11 hanya punya jalur yang **menuntut** baris pertama sudah ada
(`…/putaran` mewarisi dari `.AdjustmentList(1)`; akseptasi, tolak, dan Komite menyunting baris yang ada).
Pendaftaran melahirkan peserta **tanpa** baris (`Baris: []`), jadi grid adjustment layar Detail kosong untuk
selamanya — dan `…/putaran` atas peserta itu menjawab **400 "bukan milik klaim"**, sebab `AmbilBaris`
menggabung ke baris adjustment dan peserta tanpa baris tidak muncul di sana.

**Pembacaan XML — sebagai pohon, `pyStepsBlockName` dicetak** (seluruh activity Claim Life yang menyebut
`AdjustmentList`; nol langkah di bawah bertanda `//`):

| Activity | Langkah | Apa yang dilakukannya pada `AdjustmentList` |
| --- | --- | --- |
| `SetIndexAdjustmentList` | 1 b328, 2, 3 b570–744 | `.IsCheck = true` peserta; `.IndexPremiumList`; baris `(<LAST>)` mewarisi delapan kolom dari `(1)` — pada grid kosong keduanya baris yang SAMA |
| `SaveInsuredClaim_Act` | 2.2 b1341 | menulis `TempDetail.pxResults(<LAST>).AdjustmentList(<LAST>)` — halaman **sementara** unggahan |
| `SavePesertaClaim` | 7.8 b3671 *(hidup, WHEN b3919 `.IsCheck=="true"`)* | menulis delapan kolom baris `(<LAST>)` peserta `(<LAST>)` — **baris pertama lahir saat pendaftaran** |
| `SaveOutStandingLife_Act` | 22.1 b10548 | `@LengthOfPageList(.AdjustmentList)=0` → **lewati** peserta tanpa baris |
| `SpreadingClaimLife_Act` | 7.1 b1962, 8.2 b2394 | mengisi `SpreadingList`/`RetroLifeList` baris yang ada |
| `SaveAdjustment_Act` | 1.5 b1808, 1.6 b1979 | akseptasi atas baris `(<LAST>)` yang ada |
| `DeletePesertaClaimLife` | 1.1 b316 | menomori ulang `.IndexPremiumList` baris yang ada |
| `serviceInsertArasapasClaimLife_act` | 1.1 b311 | membaca baris yang ada |

Yang **melahirkan** baris hanya dua: tombol `Add` (`Section/ClaimLifeDetailGCNM.xml` b17937 → `addRow`
b17947 + `SetIndexAdjustmentList` b17991; tampil bila `pyWorkPage.pyPosition =='ReasLifeSPV'` b18160) dan
pendaftaran 7.8.

**Yang dibangun** (`[DIPUTUSKAN; veto work owner]` butir bo):

- **Satu rute, bukan dua** — `POST /api/klaim-life/{id}/peserta/{pesertaId}/putaran`. Tombol XML-nya pun
  satu; yang membedakan kedua jalur hanya keadaan grid, dan itu dibaca layanan. Peserta tanpa baris →
  `Putaran.tambahPertama`; peserta berbaris → jalur putaran tiket 11, **tidak diubah**.
- **Baris pertama KOSONG** — `services.BarisPertama`, nilai nol seluruhnya, termasuk mata uang
  (`TambahBaris` mengisi mata uang dari peserta untuk baris pertama; itu bukan bentuk `Add`). Penanda
  dipilih (`IS_CHECK = 'true'`, `models.PenandaDipilih`) dipasang di transaksi yang sama; jejak kelahiran
  baris direkam (DARI/KE kosong — barisnya belum berstatus); header **tidak** dicerminkan.
- **Gerbang**, di atas gerbang `Tambah` (identitas, `PeranSimpanOutstanding`, kasus terbuka): (1) tahap
  **Claim Analis** (b18160) → `ErrTahapTanpaAddAdjustment` 409; (2) pemegangnya — sama dengan
  `PeranSimpanOutstanding`, dijaga `TestPemegangClaimAnalisSamaDenganPeranPutaran`; (3) **tujuh gerbang
  `.STS_REJECT=='1' || '2'`** peserta (b2628, b4682, b5059, b5870, b6152, b7335, b15234) →
  `ErrPesertaSudahDiputus` 409. `Add` b17937 sendiri **tidak** membawa `pyDisabledWhen` — gerbang (3)
  diberlakukan karena butir bo memintanya. Pada jalur putaran gerbang (1) dan (3) **tidak** dipasang:
  (3) bertentangan dengannya (peserta yang ditolak berkode `"2"`, justru yang dibuka putarannya), dan (1)
  akan mengubah jalur yang sudah berjalan tanpa bukti di tahap mana kasus berdiri sesudah Komite menolak.
- **Layar** — `KlaimLife.tsx`: tombol `Add` (`DETAIL.tambahAdjustment`, `bolehAddAdjustment`: Claim
  Analis, grid kosong, peserta belum diputus, kasus terbuka) memanggil rute yang sama dengan "Putaran
  berikutnya".
- ⛔ **`Delete` b19120 BERDIRI TETAPI MATI** (`DETAIL.hapusAdjustment`, tombol `disabled`). `deleteRow`
  b19130 menghapus baris; **ADR-U-0031** menetapkan nol hapus fisik di jalur pengguna — penghapusan adalah
  PENANDA + nilai balik — dan `T_CLAIMLF_ADJUSTMENT` tidak punya kolom penanda. Membuatnya keputusan skema
  (migrasi) yang tidak ada di brief → **OQ-N7**.

**Bukti:** `APP_RNM/internal/services/adjustment.go:BarisPertama`,
`APP_RNM/internal/services/hasilkomite.go:Putaran.tambahPertama`, `APP_RNM/internal/handlers/putaran.go:jawabGalatPutaran`;
uji `TestBarisPertamaLahirKosong`, `TestBarisPertamaHanyaUntukGridKosong`,
`TestBarisPertamaMenghormatiGerbangSTSReject`, `TestPemegangClaimAnalisSamaDenganPeranPutaran`,
`TestTombolAddAdjustmentVERBATIMDariKorpus`, `TestGalatPutaranDipetakanKeKodeYangBenar`,
`TestPutaranTanpaOracle503`; `db`: `TestAddMelahirkanBarisPertamaKosong`,
`TestAddBarisPertamaMenjagaGerbangnya` (**melewati** tanpa `ORACLE_DSN` — melewati bukan lulus); layar:
`src/services/adjustmentpertama.test.ts`, `src/assets/labels.test.ts` (b17937, b19120).

**Pertanyaan terbuka yang lahir:**

- **OQ-N7** — `Delete` baris adjustment: kolom penanda hapus di `T_CLAIMLF_ADJUSTMENT` (dan penyaring di
  setiap pembaca hilir, ADR-U-0031 Akibat 2), atau `Delete` dinyatakan tidak berlaku? XML menampilkannya
  hanya untuk baris yang belum disimpan ke Outstanding (`.PrintFaceClaim == ''` b19399).
- **OQ-N8** — baris pertama lahir **kosong**, dan **nol rute** menyunting sel baris adjustment
  (`CLAIM_AMOUNT`, `CURRENCY`, kolom uang lain). `Save to RNM` langkah 11.17 menolak baris tanpa
  `CLAIM_GROSS` (b6043), jadi kasus yang barisnya lahir lewat `Add` belum dapat disimpan ke Outstanding
  sampai rute sunting sel ada.
- **OQ-N9** — di Pega baris pertama lahir **saat pendaftaran** (`SavePesertaClaim` 7.8); di sini pendaftaran
  melahirkan peserta tanpa baris (AC 32 tiket 02). `Add` menutup jalan layarnya, tetapi urutan kerjanya
  berbeda: Pega mengisi baris saat Register di tahap Admin, sedangkan `Add` hanya tampil di Claim Analis.
  Apakah pendaftaran semestinya melahirkan baris pertama (dan `Add` hanya untuk baris tambahan)?

### Tambahan tinjauan — 29 September 2026 (GILIRAN-13, `/code-review`)

1. **`pyStepsBlockName` kedua activity besar, dicetak.** Tabel pembacaan di atas memuat langkah yang
   menyentuh `AdjustmentList` saja, dan tidak satu pun bertanda `//`. Langkah LAIN yang ter-remark:
   `SavePesertaClaim` **tujuh** — 3 b777, 4 b922, 7.2 b1863, 7.3 b2039, 7.4 b2279, 7.6 b2606, 9 b4875;
   `SaveOutStandingLife_Act` **delapan** — 11.3 b3495, 11.9 b4632, 11.11 b5009, 12 b6178, 13 b7074,
   14 b7293, 15 b7512, **23 b10649**.
2. ⚠️ **Langkah 23 `SaveOutStandingLife_Act` ter-remark**, beserta anaknya 23.1 b10841 / 23.2 b11066.
   Bab 27-09 ("enam total peserta") dan komentar `KlaimLife.tsx` (panel total) mengutip 23.1/23.2
   sebagai penghitung total — kutipan itu menunjuk langkah MATI. Nilainya tetap benar lewat
   `SavePesertaClaim` 8.1 b4220 / 8.2 b4591 (langkah 8 hidup); yang salah hanya kutipannya. Di luar
   diff GILIRAN-13; dicatat, tidak disunting.
3. **`Delete` di layar** mengabaikan syarat tampil `.PrintFaceClaim == ''` (b19399) — kolom itu tidak ada
   di kontrak API, dan tombolnya mati di baris mana pun.

## ⛔ Ralat bertanggal — 29 September 2026 (GILIRAN-14 paket 1: butir **bp** meralat **bo**)

Bab GILIRAN-13 di atas menulis bahwa baris pertama dibuat tombol `Add` pada grid kosong, dan bahwa
"yang melahirkan baris hanya dua: `Add` dan pendaftaran 7.8" — lalu **tetap** tidak membangun yang kedua
(**OQ-N9**). Keliru jalan lahirnya: `SavePesertaClaim` 7.8 b3671 (hidup, WHEN `.IsCheck=="true"` b3919)
melahirkan baris pertama setiap peserta terpilih **saat Submit Register**, jadi grid peserta terpilih tidak
pernah kosong. Tiket 02 memuat buktinya dan yang dibangun (`services.LahirkanBarisPendaftaran`).

**Ralat bo (dibuang, bukan dibiarkan mati):** `services.BarisPertama`, `Putaran.tambahPertama`,
`ErrTahapTanpaAddAdjustment`, `ErrPesertaSudahDiputus`, dan kedua uji `db`-nya. `Add` kembali = jalur
**putaran** tiket 11. Peserta tanpa baris — kini hanya ada pada klaim LAMA (A4 migrasi data; DEV: nol klaim
terlihat) — dijawab jujur **409** "tidak ada baris" (`BarisLanjutan`), bukan 400 "bukan milik klaim":
kepemilikannya diperiksa lewat `AmbilPeserta` (`pesertaMilikKlaim`); uji `db`
`TestPutaranPesertaTanpaBarisDijawabJujur`.

**Layar:** tombol "Putaran berikutnya" dan `Add` untuk grid kosong digabung menjadi SATU tombol **`Add`**
(`DETAIL.tambahAdjustment`), tampil di tahap Claim Analis (b18160) bila baris terakhir ditolak
(`bolehAddAdjustment`). XML hanya punya satu tombol; dua label untuk satu rute adalah temuan tinjauan
GILIRAN-13. Uji `src/services/tomboladd.test.ts`.

**OQ-N9 ditutup** oleh bp. **OQ-N8** (sunting sel) tetap terbuka sampai butir br (paket 3).

## ⛔ Ralat bertanggal — 29 September 2026 (GILIRAN-14 paket 3, butir **br**: sunting sel menurut grid XML)

`[DIPUTUSKAN — XML; veto work owner]` butir **br**: "hanya kolom yang `ClaimLifeDetailGCNM.xml` grid b17126
tampilkan dapat disunting". Dibaca per kolom — dan **vonisnya nol**.

**Cara membaca, dikalibrasi lebih dulu.** `pyEditOptions` sebuah sel ditulis SEBELUM `pyValue`-nya (rentang sesudah
`pyValue` sebelumnya). Kosakata di seluruh `Section/` Claim Life: `Auto` 450, `Editable` 15, `Read-only` 315 —
`[terverifikasi]` dua metode yang sepakat: (1) `grep -rhoE "<pyEditOptions>[^<]*</pyEditOptions>" Section | sort | uniq -c`
dari folder `Claim Life`; (2) `py` menyisir setiap baris setiap `Section/*.xml` dengan regex yang sama dan mencacahnya
per nilai. Jawaban
yang sudah diketahui: `.CLAIM_RECEIVED_DATE` bertanda **`Editable`** di dialog `EditDateClaimLife_Section.xml` (b1034 —
dialog yang memang menyuntingnya, rute `PUT …/tanggal-klaim`) dan **`Read-only`** di layar Detail. `pyReadOnlyCondition`
yang tampak di rentang sebuah sel tidak mengubah opsi selnya (`.DOB`, `.CLAIM_RECEIVED_DATE` di layar Detail tetap
`Read-only`).

**Grid `.AdjustmentList` b17126** — empat kolom data, keempatnya `pyReadOnly true` / `pyEditOptions Read-only`:
`.STS_REJECT` (b18203/b18213), `.ADJUSTMENT_DATE` (b18386/b18395), `.ACCEPTEDNO` (b18627/b18635, tampil bila
`.PrintFaceClaim=1` b18597), `.ACCEPTATION_DATE` (b18831/b18839).

**Panel yang grid buka** — `pyEditingMode expandPane` b19566, `pyEditAction Adjustment_Detail` b19583 →
`Section/AdjustmentDetail_Section.xml`: SETIAP medan data `Read-only`, nol `Editable` — antara lain `.CURRENCY`
b1448, `.PCTClaim` b1614, `.CLAIM_PAID` b1937, `.SHARE_NUSANTARA_RE` b2644, `.CEDING_RETENTION` b2807,
**`.CLAIM_GROSS` b2970**, `.CLAIM_RETRO` b3232, `.SUM_REASURED` b4097, `.SUM_INSURED` b4260, `.SHARE_RETRO` b4423,
`.RETROCEDED_SHARE` b4614, rekening (`.PayableTo` b5884, `.NameOfBank` b6080, `.SwiftCode` b6620, `.BranchOfBank`
b6825, `.NoAccount` b7004). Hanya tombol yang `Auto`.

**Yang dibangun: TIDAK ada rute sunting.** Rute `PUT …/adjustment/{adjId}` yang menjawab 422 untuk setiap kolom adalah
kode mati. Vonisnya dikunci `TestGridAdjustmentNolSelDapatDisunting`, dengan tafsir yang diuji dulu atas jawaban yang
diketahui (`TestKalibrasiOpsiSuntingDariDialogEditDate`); bila korpus kelak berkata lain, uji itu yang merah lebih dulu.

**OQ-N8 TETAP TERBUKA** *(diralat tinjauan GILIRAN-14 — draf pertama bab ini menulis "ditutup oleh bp"; menutup
OQ bukan wewenang executor, CLAUDE.md §6)*. Yang OQ-N8 cemaskan adalah baris yang tidak dapat disimpan ke Outstanding:
`Save to RNM` 11.17 b6043 menolak `.CLAIM_GROSS` kosong. bp memang melahirkan baris dengan `CLAIM_AMOUNT` (7.8 b3828) —
tetapi 7.8 **tidak** menulis `.CLAIM_GROSS`, dan gerbangnya terpenuhi hanya lewat pemetaan aplikasi `.CLAIM_GROSS` →
`CLAIM_AMOUNT` (catatan 7), pemetaan yang justru OQ-N11 pertanyakan. Keduanya karena itu diserahkan bersama ke work
owner.

⚠️ **Vonis br bertentangan dengan premis keputusannya** ("tombol grid dapat disunting", `PUT …/adjustment/{adjId}`):
XML menjawab nol sel. Executor mengikuti klausa pengaturnya ("HANYA kolom yang … dapat disunting") dan tidak membangun
rute — **work owner yang memutuskan** apakah vonis itu diterima atau `CLAIM_GROSS` dibuka sebagai penyimpangan sadar.

**Pertanyaan terbuka baru:**
- **OQ-N11** — `.CLAIM_GROSS` di panel `Read-only` (b2961/b2970), tetapi **wajib** (`pyRequired true` b2977/b3026,
  `pyRequiredNew always` b3021) dan punya aksi `change` → `CountClaimAmountLife_Act` (b3036/b3049) — bentuk medan
  masukan yang dikunci. Nol activity/data transform di korpus Claim Life yang MENULIS `.CLAIM_GROSS`
  (`CountClaimAmountLife_Act`, `RejectOSClaimLife_Act`, `SaveOutStandingLife_Act`, `SpreadingClaimLife_Act` hanya
  membacanya) — `[terverifikasi]` dua metode: (1) `grep -rlE "CLAIM_GROSS" .` di folder `Claim Life` → 5 berkas
  (keempat activity itu + `AdjustmentDetail_Section.xml`); (2) pohon langkah keempat activity (`pohon.py`,
  `pyStepsBlockName` dicetak) disisir untuk `PropertiesName = …CLAIM_GROSS` → nol. Siapa yang mengisinya di Pega — pengguna (dan tanda `Read-only` itu keliru dibaca), atau aturan di luar
  ekspor ini? Sampai dijawab, aplikasi memetakan `.CLAIM_GROSS` ke `CLAIM_AMOUNT` yang lahir di pendaftaran.

### Tambahan tinjauan — 29 September 2026 (GILIRAN-14, `/code-review`)

- **Gerbang tahap `Add` kini di layanan juga.** "Add tetap = putaran bergerbang ReasLifeSPV b18160" — gerbang itu
  semula hanya di layar (`bolehAddAdjustment`). `Putaran.Tambah` kini menolak di luar Claim Analis
  (`ErrTahapTanpaAddAdjustment`, 409); uji `db` `TestPutaranHanyaDiClaimAnalis`. Fixture `pohonUjiKomite` kini di
  Claim Analis (penyerahan ke Komite tidak bergerbang tahap).
- **Syarat "baris terakhir ditolak" di layar adalah penyimpangan sadar**: b18160 menampilkan `Add` tanpa syarat baris;
  syaratnya milik layanan (`BarisLanjutan`).
- **Jejak kelahiran baris saat Register tidak direkam** — ADR-U-0007 menjejak TRANSISI STATUS, dan baris 7.8 lahir
  tanpa status; transisi pertamanya (`0`, Save to RNM) direkam jalur itu. Cabang `tambahPertama` yang dibuang dahulu
  merekamnya (DARI/KE kosong); dicatat, tidak ditiru di pendaftaran.

## ⛔ Keputusan work owner bertanggal — 29 September 2026 (GILIRAN-15 paket 1: br, N7, N11 — "ikuti rekomendasi")

**br dan OQ-N8 — DITUTUP.** *Ikut XML: tidak ada sunting sel adjustment.* Vonis bab GILIRAN-14 diterima: keempat kolom
grid b17126 dan seluruh medan panel `Adjustment_Detail` `Read-only`. Tidak ada rute `PUT …/adjustment/{adjId}`;
`TestGridAdjustmentNolSelDapatDisunting` tetap sebagai penjaganya. Keberatan GILIRAN-14 ("bp hanya memenuhi 11.17 lewat
pemetaan") dijawab dengan keputusan ini sendiri: nilai baris lahir dari polisnya saat Register (7.8), dan sunting sel
bukan bagian sistem lama.

**OQ-N7 — DITUTUP.** *`Delete` baris adjustment TIDAK BERLAKU.* ADR-U-0031 melarang hapus fisik di jalur pengguna, dan
`T_CLAIMLF_ADJUSTMENT` tidak punya kolom penanda. Tombol `Delete` b19120 kini **tidak dirender** sama sekali — bukan
tombol mati: tombol yang berdiri tetapi tidak dapat ditekan menjanjikan aksi yang tidak akan pernah ada.
`DETAIL.hapusAdjustment` dibuang; uji `src/services/tomboladd.test.ts` menuntut ketiadaannya. Alasannya juga di
`PARITAS-LAYAR-DAN-AKSI.md`, bab "29-09-2026 (GILIRAN-15) — `Delete` baris adjustment TIDAK BERLAKU".

**OQ-N11 — dipindah ke daftar pemilik ekspor** (`OQ-untuk-tim.md`). *Sementara kosong:* aplikasi TIDAK menulis
`CLAIM_GROSS` dan tidak mengarang penulisnya — buktinya 4 pembaca, 0 penulis (dua metode, bab br di atas).
⚠️ Hubungannya dengan **catatan 7** (tiket 14: "kolom kita untuk `CLAIM_GROSS` bernama `CLAIM_AMOUNT` — satu nilai, dua
nama"): pemetaan itu keputusan lama dan TIDAK diubah; gerbang `Save to RNM` 11.17 tetap membaca `CLAIM_AMOUNT`. Bila
pemilik ekspor menjawab bahwa `CLAIM_GROSS` nilai TERSENDIRI (bukan `CLAIM_AMOUNT`), catatan 7 dibuka ulang — dan
gerbang 11.17 akan menolak setiap baris sampai penulisnya ada.

### Tambahan tinjauan — 29 September 2026 (GILIRAN-15, `/code-review`)

- Perintah audit metode kedua OQ-N11: `py pohon.py <activity>.split <keluaran>` (pemecah `><` → `>\n<`, `pyStepsBlockName`
  dicetak) untuk keempat activity, lalu `grep -n "PropertiesName = .*CLAIM_GROSS"` atas keempat keluaran → nol baris.
- **OQ-N12** *(untuk work owner — lahir dari tinjauan GILIRAN-15)* — dua jawaban 29-09-2026 bertemu dengan catatan 7
tiket 03. (1) N11: "`CLAIM_GROSS` sementara **kosong**". (2) "Kosong = nol di 7.8": `CLAIM_AMOUNT` baris yang lahir saat
Submit tidak pernah kosong (paling kecil `0`). (3) Catatan 7 (tiket 14): `CLAIM_GROSS` = `CLAIM_AMOUNT`, sehingga gerbang
`Save to RNM` 11.17.1 (`simpanrnm.go`, b6043 `.CLAIM_GROSS==""`) membaca `CLAIM_AMOUNT`. Akibatnya gerbang 11.17.1 **tidak
pernah menolak** baris yang lahir saat Register (dan putaran yang mewarisinya), sedangkan di Pega — `CLAIM_GROSS` tanpa
penulis — gerbang itu menolak. Perilaku TIDAK diubah executor. Pilih: (a) pertahankan catatan 7 (gerbang praktis mati untuk
baris 7.8); (b) pisahkan `CLAIM_GROSS` dari `CLAIM_AMOUNT` (kolom baru — migrasi — dan gerbang 11.17.1 akan menolak setiap
baris sampai penulisnya ada); (c) tunggu jawaban OQ-N11 dari pemilik ekspor.

## ⛔ Keputusan work owner 29-09-2026 — OQ-N12 (a) (GILIRAN-16)

*Kutipan: "A".* **Catatan 7 dipertahankan:** `CLAIM_GROSS` dibaca sama dengan `CLAIM_AMOUNT` baris adjustment,
`[sementara — menunggu OQ-N11 pemilik ekspor]`. Perilaku tidak berubah: gerbang `Save to RNM` 11.17.1 (`services/simpanrnm.go`) tetap
membaca `CLAIM_AMOUNT` dan, karena 7.8 membaca sumber kosong sebagai 0, praktis tidak menolak baris yang lahir saat
Register (dan putaran yang mewarisinya). Dasar — dibaca asisten, diperiksa ulang executor atas
`Activity/SaveOutStandingLife_Act.xml` (pecahan `><` → `>\n<`) dan pohon langkahnya (`pyStepsBlockName` dicetak):

| Fakta | Bukti |
| --- | --- |
| gerbang Save to RNM menolak bila `CLAIM_GROSS` kosong | langkah **11.17.1** `Page-Set-Messages` (nama activity b5935, `pyStepPageReference` `RH_1.pySteps(11).pySteps(17).pySteps(1)` b5936), prasyarat `true` b5954, syarat `.CLAIM_GROSS==""` b6043 |
| langkah itu **hidup** | `pyStepsBlockName` di dalam 11.17 kosong (b5904 untuk 11.17, b5946 untuk 11.17.1); `//` terdekat b6178 milik langkah **12** (`RH_1.pySteps(12)` b6167) |
| `CLAIM_GROSS` tanpa penulis di korpus | pembaca: `CountClaimAmountLife_Act`, `RejectOSClaimLife_Act`, `SaveOutStandingLife_Act`, `SpreadingClaimLife_Act` (+ `Section/AdjustmentDetail_Section.xml`); tag `<PropertiesName>` yang memuat `CLAIM_GROSS` di keempat activity: **0** |
| kesimpulan | bila `CLAIM_GROSS` benar-benar kosong di produksi, **tidak satu klaim pun** dapat melewati Save to RNM di Pega — mustahil untuk sistem yang dipakai; nilainya pasti diisi rule yang tidak diekspor `[dugaan: Declare Expression]` |

Dibuka ulang bila OQ-N11 terjawab dan rule itu menghitung `CLAIM_GROSS` dari sesuatu selain `CLAIM_AMOUNT`: catatan 7
dicabut, `CLAIM_GROSS` mendapat kolom sendiri (migrasi) beserta penulis yang meniru rule tersebut. Komentar gerbang
11.17.1 di `simpanrnm.go` merujuk OQ-N12/N11; OQ-N12 ditutup di `OQ-untuk-tim.md`.

## Keputusan bertanggal — 29 September 2026 (GILIRAN-17 paket 1: OQ-N1, OQ-N3 ditutup) `[keputusan work owner 29-09-2026 — lembar keputusan, "rekomendasi"]`

| OQ | Keputusan | Bukti | Kode |
| --- | --- | --- | --- |
| **N1** | bendera simpan `pyWorkPage.Save` = **keadaan turunan**, tanpa kolom | langkah 21/24 (b11448), `pyDisabledWhen` b21095, kontainer *Participant Details* b15490 | tidak berubah — tombol tetap hidup, tulisan idempoten; 22.1.3.2 tetap tidak ditiru hurufnya |
| **N3** | gerbang retro langkah 27 **dipertahankan** | `RetroID=="L0000141" \|\| SecurityReinsurerID=="L0000134"` b11794, `RetroID=="1000013"` b11817 — XML Claim Life hidup; OQ-064 milik Komite | tidak berubah (`ArasapasDilewatiRetro`, dibaca sesudah tukar — lihat OQ-N5) |
| **N3 tambahan** | cermin header `T_GENERAL_CLAIM.STS_REJECT`/`ACCEPTED_NO` **tetap di transaksi simpan** | sistem lama: `serviceInsertArasapasClaimLife_act` 1.1.1 b371/b417, lewat langkah 28 | tidak berubah |

## Keputusan bertanggal — 29 September 2026 (GILIRAN-17 paket 2: OQ-M6 ditutup) `[keputusan work owner 29-09-2026 — lembar keputusan, "rekomendasi"]`

Mencabut peserta = **penanda**; layar menyembunyikannya.

| Hal | Isi |
| --- | --- |
| XML | tombol `DELETE` `InputOSClaimLife` b17865 (judul kolom b16839) → `deleteRow` b17874 tanpa konfirmasi (b18021) → `DeletePesertaClaimLife` b17909 (langkah 1 b227, 1.1, 1.2 `Obj-Save` b443; `pyStepsBlockName` kosong b238/b326/b453; nol SQL); tampil bila `CLAIM_NO == ''` b18082 |
| kolom | `T_CLAIMLF_PREMIUMLIST_DETAIL.STS_HAPUS` VARCHAR2(1) — migrasi `022_kolom_sts_hapus_peserta.sql` (+ `_down`); NULL aktif, `'1'` dicabut (ADR-U-0031, pola `STS_*`); **dijalankan work owner** |
| rute | `POST /api/klaim-life/{id}/peserta/{pesertaId}/cabut` — Admin, kasus terbuka, `models.BolehCabutPeserta` (Outstanding **dan** belum Save to RNM — padanan b18082, OQ-M1), peserta milik klaim; penanda + jejak (`peserta <id>` → `dicabut`) satu transaksi |
| pembaca | setiap pembaca dan penulis tabel peserta menyaring `STS_HAPUS IS NULL`: peserta, baris adjustment, dokumen (daftar, satu, pemilik), diagnosa, spreading, penanda Save to RNM, tanggal, status, `IS_CHECK`. Pengecualian beralasan: `Simpan` (sisip), `HapusFisik` (uji), `Dampak` (hapus klaim utuh) |
| penjaga | `TestSetiapPenyentuhTabelPesertaMenyaringPenandaCabut` — fungsi baru yang menyentuh tabel peserta gagal sampai menyaring atau dikecualikan; dibuktikan merah lewat mutasi |
| cermin warisan | ⛔ *Ralat /code-review GILIRAN-17:* baris cermin `OS_AKSEPTASI_KLAIM_LIFE` peserta itu (status NULL, sejak OQ-N2 bernama) **dibuang** di transaksi yang sama (`HapusCerminBelumDisimpan`, `CASEID` + `STS_REJECT IS NULL`). Di Pega baris cermin baru lahir saat Save Outstanding, jadi peserta yang dilepas sebelumnya tidak pernah terlihat hilir |
| uji | `TestBolehCabutPesertaSebelumSaveRNM`, `TestSQLCabutPesertaMenandaiBukanMenghapus`, `TestMigrasi022PenandaCabutPeserta`, `TestCabutPesertaMenjagaPagarnya`, `TestRuteCabutPesertaDanGalatnya`, `cabutpeserta.test.ts`, `labels.test.ts`, `db` `TestCabutPesertaMenandaiDanMenyembunyikan` |

## Keputusan bertanggal — 29 September 2026 (GILIRAN-17 paket 3: OQ-N2, OQ-N5, OQ-M7 ditutup) `[keputusan work owner 29-09-2026 — lembar keputusan, "rekomendasi"]`

| OQ | Keputusan | Bukti XML | Kode |
| --- | --- | --- | --- |
| **N2** | cermin mengisi `NAME_OF_INSURED`, `DOB`, `CEDINGCO` seperti Pega, **di dalam SQL** | `SaveOutStandingLife_Act` 22.1.1 b8747 (CARI6 b8906, CARI8 b8946, CARI27 b9226) → `InsertJsonKlaimLife_sql` b93/b95/b114 | `PohonKlaim.Simpan` → `KlaimLife.IsiTertanggungCermin` (`UPDATE … SELECT` dari `M_LIFE_PREMIUM_DETAIL`, `TRUNC(m.DOB)`, `EXISTS` = sumber hilang jadi galat) + `sqlIsiCedingCermin` (`T_PREMIUM_LIST.CEDING_CO`, PROD_KE terakhir); pemeriksa ganda mengecualikan `CASEID` sendiri dan (kematian) baris tanpa status — baris aplikasi yang baru terdaftar tidak menutupi baris era Pega |
| **N5** | tukar retro **dua** WHEN; `ProdDateTime` tidak dipakai | `InsertJsonClaimLife_Act` langkah 2 b1016 (`pyStepsBlockName` kosong b1028): b1268, b1291; b1314 dibuang | `ArasapasDilewatiRetro` murni `bool`; `ErrGerbangRetroTakTerputuskan` dan jawaban `ditahan` dibuang |
| **M7** | izin baca `RATE_LIFE` sempit seperti butir bh | `GetRateRetro` b84 (`SpreadingClaimLife_Act` langkah 6 b1692/b1750, `CARI3 = .OUTWARDRATEID` b1527–b1528) | `repository/ratelife.go` `RateLife.Baca`, lima kolom; penjaga `TestMasterViewTidakDisentuh` = peta izin per berkas (`izinViewRate`); `OUTWARDRATEID` `[terbuka — DBA]`, Spreading tidak dipanggil |

Yang tetap terbuka: **OQ-N13** (baru). Status cermin `'0'` saat Save to RNM (b176) belum ditulis, karena langkah itu
memperlakukan tabel warisan baca-saja. Karena itu klaim yang menunggu Komite belum tertangkap sebagai ganda.

Uji: `TestSQLIsiTertanggungCerminDariSumber`, `TestSQLGandaMengecualikanKlaimSendiriDanBarisTanpaStatus`,
`TestSQLGandaBerurutPosisi`, `TestSimpanMengisiTertanggungCermin`, `TestArasapasRetroTukarDuaSyarat`,
`TestRateLifeKolomTetapBerkunciIDUSEDBY`, `TestRateLifePengenalKosongGagalTerang`, `TestMasterViewTidakDisentuh` (mutasi `USEDBY`
menjadi merah), dan `db` `TestCerminMengisiTertanggungDariSumber`. Uji `db` itu menghitung kecocokan nama/DOB di Oracle; nilainya
tidak dibaca ke Go.

## Perbaikan /code-review GILIRAN-17 — 29 September 2026

| Temuan | Tindakan |
| --- | --- |
| peserta tercabut meninggalkan baris cermin bernama di `OS_AKSEPTASI_KLAIM_LIFE` | dibuang dalam transaksi pencabutan (baris milik `CASEID` itu yang belum berstatus) |
| dokumen dapat diunggah ke peserta tercabut atau milik klaim lain | `Unggahan.Unggah` memeriksa peserta aktif milik klaim (`pesertaMilikKlaim`) sebelum berkas ditulis |
| balapan Save to RNM dengan cabut | `sqlCabutPeserta` memeriksa ulang "belum Save to RNM" di dalam pernyataan (`NOT EXISTS`) |
| cermin tertanggung/ceding dikunci `ID` saja | kini `ID` **dan** `CASEID` |
| `CASEID` kosong menyaring seluruh baris warisan di pemeriksa ganda | jatuh ke pengenal work (`CASEID` klaim aplikasi, butir ae1) |
| galat muat-ulang sesudah tolak/cabut hilang di dialog yang tertutup | dilaporkan di layar dengan kalimat yang benar ("sudah ditolak/dicabut, tetapi gagal dibaca ulang") |
| penjaga saringan `STS_HAPUS` lolos lewat komentar dan nama fungsi senama | diperiksa atas kode tanpa komentar, pengecualian berkunci `berkas:fungsi` |
| tombol `DELETE` bernama sama untuk tiap peserta; Remarks > 4000 byte tanpa pesan | `aria-label` bernomor sertifikat; pesan batas byte di dialog |

Sisa yang dicatat, **tidak** diubah: peserta **terakhir** dapat dicabut, dan Save to RNM atas klaim tanpa peserta aktif
tidak menandai apa pun (Pega `deleteRow` pun tidak membatasinya). Bila bisnis menghendaki batas itu, ia keputusan baru.
