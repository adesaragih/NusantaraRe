# 22: `loader.Flatten` — satu dokumen penawaran → baris 78 tabel flat (murni)

> ⚠️ **Disusun agent dari spec/XML atas perintah work owner — bukan hasil `/to-tickets`.** Perintah: arahan work owner
> 02-10-2026 (diteruskan sesi `nusantarare-0f`, sumbernya diperiksa agent) langkah 1–2 jalur pemuatan data lama; pilihan
> work owner **"Flatten dulu, migrasi ditahan"** (butir 66).

**What to build:** seam **`loader.Flatten`** spec 11 — satu dokumen penawaran (isi `POOLDATA.JSON_POLIS.DATA_JSON`
beserta `IDPEGA`-nya) menjadi himpunan baris untuk **78 tabel flat**, ditambah diagnostik berupa **hitungan saja**.
⛔ **Murni**: tanpa basis data, tanpa jam, tanpa berkas. ID dan `PARENT_ID` diberikan repository (tiket 24) dari
`Baris.Kunci` / `Baris.Induk`.

**Asal** (READ-ONLY, ditetapkan work owner 02-10-2026 — register bab *Ralat kelima*):
- `D:\migrasi\RNM\OUTPUT\04-spec\11-spec-pemuatan-data-lama.md` — bab *Implementation Decisions*, *Yang wajib berhenti
  keras*, *Testing Decisions*.
- `D:\migrasi\RNM\OUTPUT\08-flat\BAHAN-SPEC-PEMUATAN.md` §0–§8.
- `D:\migrasi\RNM\OUTPUT\08-flat\Tabel-Flat-Lintas-Siklus.xlsx` lembar `Kolom`, `Jalur Sumber`, `Audit Mata Uang`,
  `Daftar Relasi` — dibangkitkan ke `skema_gen.go` (`services/loader/bangkit`, dijalankan tangan).
- Teks keputusan V-xx: lembar `BACA-INI` `D:\migrasi\RNM\Claude outputs\Tabel-Flat-per-Grup-Bisnis.xlsx`.
- `_DAFTAR-ISSUE-TERBUKA.md` **F-3**: tiket yang dulu terhalang struktur tabel flat "kini tidak terhalang lagi".

**Sumber skema — dua workbook, diukur 02-10-2026 dua cara (pengurai xlsx vs pencacah baris DDL):**
`Tabel-Flat-Lintas-Siklus.xlsx` = `DDL-tabel-flat-draf.sql` **persis** (78 tabel, 1.329 kolom, selisih nol dua arah).
`Claude outputs\Tabel-Flat-per-Grup-Bisnis.xlsx` (75 / 1.290) adalah **himpunan bagiannya**: 39 kolom selisih = 3 tabel
wadah V-27 (20 kolom sistem) + 18 `CURRENCY_CODE` (K-063/K-069) + `OLD_POLIS_ID` (K-068); **nol** kolom ke arah
sebaliknya. Generator membaca Lintas-Siklus karena itulah yang ditetapkan work owner bersama DDL draf.

**Blocked by:** — (F-3)

**Status:** ready-for-human — diport 02-10-2026; keputusan work owner butir 68 diterapkan (kelima kasus nyata kini
**lolos mode ketat**); butir 69–72 diterapkan; uji mutasi 62/62; code review dua sumbu ditangani; menunggu tinjauan work owner

## Yang dibangun

- [x] `services/loader/bangkit` — generator `skema_gen.go` dari workbook (zip + XML stdlib, sel ditempatkan menurut
      huruf kolom); penjaga `TestSkemaSamaDenganKeluaranBangkit` (env `FLAT_RNM_XLSX`)
- [x] Pengurai JSON **berurutan** (urutan kunci dipertahankan; kunci ganda ditolak)
- [x] Mesin generik: jalur → tabel (Jalur Sumber), FIELD ASLI → kolom (lembar Kolom), kolom sistem `IDPEGA`,
      `COB_GROUP`, `SEQ_NO`, `ROW_UID`, `PARENT_TABLE`, `SRC_PATH`, `JENIS_WORK`/`NO_WORK`
- [x] Aturan khusus — hanya yang tertulis: V-16/V-16a (COB), V-17b, V-19, V-22/22a/22b, V-24b/V-26, V-25, V-27, V-28,
      V-28b, V-30, V-31, V-33, V-34, V-39, V-41, V-44, V-45, V-48, V-49, K-063 (c), K-069 (`aturan.go`)
- [x] Format angka BAHAN §7: bulat, titik desimal, koma desimal (pengecualian, terhitung); selainnya berhenti keras
      **tanpa menyebut nilainya**; tidak membulatkan (ADR-0005)
- [x] Diagnostik tanpa nilai data: baris per tabel, mata uang tak diketahui, koma desimal, angka > 38 digit, dibuang
      (per alasan), cabang kosong, tak terpetakan, penunjuk V-47 belum dikonversi; **setiap medan terisi masuk tepat
      satu ember** (`TestSetiapMedanTerhitung`)
- [x] Setiap kolom dari 1.329 punya asal tertulis: medan 813 · diisi Flatten 320 · diisi repository 164 · sengaja kosong
      32 (`TestSetiapKolomBerasal`, dihitung dua jalan)
- [ ] Bentuk masukan XML (spec 11 uji 11) — **tidak dibangun**: masukan produksi JSON (K-066); contoh XML hanya di korpus
      dan tidak boleh masuk repositori. Pengukuran agent memakai pengubah XML→JSON di luar repositori

## Uji (spec 11 *Yang diuji di seam `loader.Flatten`*)

| # | Spec | Uji |
| :-: | --- | --- |
| 1 | penawaran sederhana | `TestFlattenBentukDasar` |
| 2–4 | bersarang 8 tingkat, 12 induk ganda, wadah murni | `TestFlattenSetiapJalur` — **141 jalur** dibangkitkan dari Jalur Sumber (148 − akar − 2 ruas V-30 − 4 jalur tabel wadah, yang tercakup sebagai leluhur); kelima induk `T_COVERAGELIST`; `TestKedalamanRancangan` (8) |
| 5 | tiap kelompok lini, BusinessType ganda | `TestFlattenLiniBisnis` + cara kedua `TestKasusUkur` (BusinessType → lembar `Grup Bisnis`) |
| 6 | baris tanpa mata uang | `TestFlattenMataUang` |
| 7 | koma desimal | `TestFlattenFormatAngka` |
| 8–9 | cabang dikecualikan, tabel `Total*` | `TestFlattenCabangDibuang` |
| 10 | tiap berhenti keras | `TestFlattenBerhentiKeras` (termasuk: pesan galat tidak menyalin isi dokumen); kode enumerasi `BusinessCode`/`BusinessOldId` **tidak ditafsirkan** Flatten (A61); induk tak tetap dan kedalaman > 8 **tak tercapai** lewat dokumen sah — dijaga invarian `TestJalurIndukAdalahAwalan`, `TestKedalamanRancangan` |
| 11 | JSON = XML | ⛔ belum (lihat di atas) |
| — | kasus nyata | `TestKasusLolosKetat`, `TestKasusUkur` (angka emas; `T_COVERAGELIST` 45/5/4/1/104 = sensus README fixture; penampung = Diagnostik), `TestSetiapMedanTerhitung` (pencacah independen) |
| — | butir 68 | `TestFlattenTeksMenyimpang` (68.1), `TestFlattenKodeDariHalamanCurrency` (68.3), penampung di `TestFlattenTakTerpetakan` (68.5) |
| — | butir 69–70 | `TestFlattenIsCedingConfirmDiselamatkan` (69.4), `TestFlattenAmandemen70` (P1, P2–P3, P6), P5 di `TestFlattenLipatan`, `TestFlattenPenunjukV47` (P7) |

Uji mutasi `docs/alat/mutasi_loader.py` → **62/62** tertangkap (`mutasi_loader.log`); dua celah yang tertangkap
putaran pertama (unsur larik bernilai, faktor V-30 tanpa `ChechBoxN`) ditutup dengan kasus baru. Satu mutasi
dikeluarkan dengan alasan tertulis di skrip (cabang cadangan `sebabDokumen` tak tercapai lewat masukan apa pun).
Putaran butir 68: mutasi "perluas 68.3 ke kembaran FR" pertama LOLOS karena hanya menambah tabel ke peta — anak
`Currency` kembaran FR adalah `T_FR_CURRENCY`, yang tidak dibaca; diganti mutasi yang mengubah kedua syarat.

**Code review dua sumbu (02-10-2026)** — ditangani: nilai `BusinessType` dan potongan isi JSON tidak lagi masuk pesan
galat; V-16 langkah 2 melewati cabang dibuang ("LocationList bisnis"); `SEQ_NO` opsi V-30 = N (sejalan A50);
predikat `PAY_` dan pencatat subpohon disatukan; medan `jenis` yang tak terpakai dibuang dari generator; `wajib`
(NOT NULL) kini ditagih di `TestFlattenSetiapJalur`; kolom aturan yang tak ada = `panic` bug program. Dicatat sebagai
keputusan, bukan diubah: A61–A65 (register). Tidak diubah (penilaian): parameter berjalan bersama di pejalan, kait
`kumpulBukanAngka` khusus uji, nama variabel skrip mutasi.

**Code review dua sumbu putaran butir 68 (02-10-2026)** — ditangani: ⛔ **cacat** — entri penampung baris yang dibatalkan
V-27 tetap memegang `Kunci` baris batal, yang lalu dipakai ulang baris berikutnya (nilai menempel ke baris salah);
kini dipindah ke baris induk, diuji. Dua halaman `Currency` berkode berbeda di bawah satu baris kini berhenti keras
(dulu saling timpa). Uji teks menyimpang mencakup **8/8** kolom, termasuk koma desimal tanpa spasi. Urutan penampung
diuji atas JSON mentah. Berkas Go paket ini dikembalikan ke **LF** (README §6 butir 3; sempat CRLF). Kepala paket dan
komentar `Nilai` menyebut penampung dan kolom teks. Dicatat sebagai keputusan: A66 (register). Di luar kode: ADR-0023
akibat 2 ("penampung berisi = belum selesai") ditegakkan di tiket 24, bukan di seam murni ini.

## Pengukuran atas 115 contoh XML korpus — 02-10-2026

`D:\migrasi\RNM\DDL\CONTOH\*.xml` diubah ke JSON berbentuk `DATA_JSON` **di scratchpad agent** (bukan di repositori;
PageList → larik `rowdata`; nol PageGroup ditemukan), lalu dijalankan dalam **mode ukur** (konflik tipe dihitung, tidak
menghentikan). Yang dicetak hanya hitungan dan nama.

- **Mode ketat, SEBELUM butir 68.1:** ⛔ **106 dari 115** berhenti di konflik tipe pertamanya — `PPN_CHECK` 74,
  `SHARE_OF_CEDING` 27, `MASTER_RATE_COVERAGE` 4, `EDM_CHARGE_FEE` 1; **9** lolos.
- ⚠️ **Ralat jendela** (verifikasi independen sesi `nusantarare-0f`, dicek ulang agent): versi tiket sebelumnya menulis
  "dua kolom" di Status dan "delapan kolom" di butir terbuka. Keduanya benar, **jendelanya berbeda**: di **5 fixture**
  hanya **2** kolom memicu — `PPN_CHECK` 5/5 kasus, `SHARE_OF_CEDING` 4 kasus (di `edm-fire-1` medannya hanya terisi di
  bawah `OldData`, yang dibuang); di **115 contoh korpus** **8** kolom (tabel di bawah).
- **COB:** Aneka 81 · FIRE 12 · Life 5 · MBUCar 6 · MarineCargo 6 · PA 5 = 115. Lembar `Grup Bisnis` workbook asal
  (dihitung atas **114**): sama, kecuali FIRE 11 — selisih tepat satu berkas tambahan. `[terverifikasi]`
- **V-27 tereproduksi:** cabang kosong yang terdeteksi mesin — `CargoList` beserta `CoverageList` dan `SpreadingList`
  di **2** berkas, `QuotationData/QuotationList` **1** — persis yang disebut teks V-27 (dua berkas HEAVY EQUIPMENT grup
  Aneka; QuotationList).
- **Konflik tipe** — kolom `NUMBER` di DDL draf yang isinya bukan angka (dua cara: mesin vs pengurai Python, jendela
  sama — tanpa `OldData`/`FacRetro`; ⚠️ grep mentah memberi angka lebih besar karena ikut `OldData`):

  | Kolom | Nilai bukan angka | Bentuk |
  | --- | ---: | --- |
  | `T_COVERAGELIST.TYPE_OF_DISCOUNT` | 94 | `true` |
  | `T_QUOTATIONDATA.SHARE_OF_CEDING` | 87 | angka bersufiks `%` |
  | `T_GENERAL_POLIS.PPN_CHECK` | 74 | `true` / `false` |
  | `T_TABLEOFLIMIT.PCT_LIMIT` · `T_FR_TABLEOFLIMIT.PCT_LIMIT` | 56 · 1 | koma desimal **berakhiran spasi** |
  | `T_FR_PRINTRISLIP.WARR_PAYMENT` | 15 | teks bebas |
  | `T_PERSONLIST.MASTER_RATE_COVERAGE` | 5 | berhuruf |
  | `T_QUOTATIONDATA.EDM_CHARGE_FEE` | 1 | berhuruf |

  ⚠️ Cara kedua sempat memberi 0 untuk `EDM_CHARGE_FEE`: **keliru**, pengurai Python mencari `EDMChargeFee`, padahal
  medannya `EdmChargeFee` (jebakan sensus no. 4, peka huruf). Dengan ejaan yang benar: 1.
- **Angka > 38 digit BERMAKNA** — ⚠️ **ralat butir 69:** hanya `T_CURRENCYLIST.PAY_NET_PREMIUM`, **3** di korpus dan **1** di
  fixture. Angka di bawah ini (putaran pertama) menghitung **koefisien mentah termasuk nol di belakang koma** — cara yang
  keliru untuk batas presisi; dipertahankan sebagai riwayat. Korpus (cara lama): `T_COVERAGELIST.NET_RATE` 185
  (seluruhnya wadah `LocationList/Property/PropertyItemList/CoverageList`; **0** di fixture),
  `T_CURRENCYLIST.PAY_EDM_PREMI_MENJADI` / `PAY_NET_PREMIUM` / `SUM_TOTAL_PAYMENT` 8 masing-masing,
  `T_FR_FACOUTOBJECTLIST.OBJECT_PREMI` 3, `T_FR_CURRENCYLIST.SUM_TOTAL_PAYMENT` 2, tiga kolom `T_FR_*` lain 1
  masing-masing. Fixture: `edm-fire-1` `T_CURRENCYLIST.PAY_EDM_PREMI_MENJADI` / `PAY_NET_PREMIUM` / `SUM_TOTAL_PAYMENT`
  1 masing-masing. Bahan keputusan presisi: `docs/USULAN-PRESISI-TABEL-FLAT.md`.
- **Mata uang tak diketahui (UNKNOWN), SEBELUM butir 68.3:** `T_SPREADINGLIST` 2.595, `T_COVERAGELIST` 2.569, `T_DEDUCTIBLELIST` 1.719,
  `T_PROPERTY` 297 (seluruhnya), `T_ANEKALIST` 127, `T_ADDITIONALCOVERAGE` 125, `T_RETROLIST` 80, `T_FR_SPREADINGLIST`
  68, `T_FR_COVERAGELIST` 34, `T_FR_FACOUTOBJECTLIST` 26, `T_FR_ANEKALIST` 15, `T_FR_DEDUCTIBLELIST` 6,
  `T_COVERAGEDATALIST` 5.

## Penyimpangan sadar dari DDL draf — butir 68.1

⛔ Delapan kolom yang DDL draf beri `NUMBER` disimpan **teks apa adanya** (`kolomTeksMenyimpang`, `aturan.go`): `%` tidak
ditafsirkan, koma desimal tidak dinormalkan, spasi di ujung `PctLimit` tidak dipangkas. Panjang `VARCHAR2`-nya
keputusan tiket 23 (ditahan). Baris DDL yang disimpangi (`08-flat\DDL-tabel-flat-draf.sql`):

| Kolom | Baris DDL draf |
| --- | --- |
| `T_GENERAL_POLIS.PPN_CHECK` | L95 `PPN_CHECK NUMBER,` |
| `T_COVERAGELIST.TYPE_OF_DISCOUNT` | L206 `TYPE_OF_DISCOUNT NUMBER,` |
| `T_QUOTATIONDATA.EDM_CHARGE_FEE` | L233 `EDM_CHARGE_FEE NUMBER,` |
| `T_QUOTATIONDATA.SHARE_OF_CEDING` | L272 `SHARE_OF_CEDING NUMBER,` |
| `T_PERSONLIST.MASTER_RATE_COVERAGE` | L513 `MASTER_RATE_COVERAGE NUMBER,` |
| `T_TABLEOFLIMIT.PCT_LIMIT` | L914 `PCT_LIMIT NUMBER` |
| `T_FR_PRINTRISLIP.WARR_PAYMENT` | L1151 `WARR_PAYMENT NUMBER` |
| `T_FR_TABLEOFLIMIT.PCT_LIMIT` | L1519 `PCT_LIMIT NUMBER` |

## Keputusan work owner 02-10-2026 (butir 68, diteruskan sesi `nusantarare-0f`, "setuju") — diterapkan

| Butir terbuka lama | Keputusan | Diterapkan |
| --- | --- | --- |
| Tipe delapan kolom | teks apa adanya | bab di atas; kelima fixture lolos ketat |
| Presisi tabel flat | tetap keputusan tim inti; tulis usulan | `docs/USULAN-PRESISI-TABEL-FLAT.md`; tiket 23 tetap ditahan |
| Kode mata uang coverage/aneka | dari `Currency/Name` | `kodeDariCurrency` — **hanya** `T_COVERAGELIST`, `T_ANEKALIST`; kembaran `T_FR_*` tidak diperluas *(diubah butir 69: kembaran FR kini ikut)*. Fixture: marine 104 + kredit coverage/aneka/spreading kini berkode; coverage FIRE tanpa halaman `Currency` tetap UNKNOWN |
| `T_PROPERTY` | dari `Currency` "bila leluhurnya punya halaman Currency" | `[terverifikasi]` nol halaman `Currency` sebagai **anak langsung** `LocationList` maupun `LocationList/Property` (0/4 fixture, 0/297 korpus); `Currency` di akar ada (3/5, 41/115) tetapi `Name`-nya tidak pernah terisi → syarat tak terpenuhi, tetap UNKNOWN. ⚠️ Ralat butir 69: rumusan lama "nol halaman `Currency` di bawah `LocationList`" terlalu luas — di bawah `PropertyItemList`/`DeductibleList`/`AnekaList` ada halaman `Currency` |
| 29 FK V-47, `TSI_TOP_RISK` | biarkan kosong sampai aturan tertulis | tidak berubah |
| `CoverageInitial`, `AdditionalShip` (dan medan tak dikenal lain) | ke penampung ADR-0023 | `Hasil.Penampung` — nilai + kunci baris + jalur; juga penunjuk V-47 *(sejak butir 72 penunjuk berkolom)*. Marine: 520 (104 + 416). ⚠️ ADR-0023: penampung berisi = pekerjaan belum selesai |
| V-30 23 vs 24 | ikuti data (24) | sudah; ralat V-30 di register |
| A49–A65 | dikonfirmasi | register |

## Keputusan work owner 02-10-2026 (butir 69, diteruskan sesi `nusantarare-0f`, "mau") — diterapkan

| Butir | Keputusan | Diterapkan |
| --- | --- | --- |
| A66 | dikonfirmasi: berhenti keras; saat muat dilewati dan dicatat (68.8) | tidak berubah |
| Penampung di Oracle | **tidak** disimpan di Oracle; tiap medan butuh kolom atau keputusan "dibuang" eksplisit | usulan `docs/USULAN-KOLOM-PENAMPUNG.md` (P1–P7); kolom **tidak** ditambah agent |
| Kembaran FR | ikut 68.3 | `kodeDariCurrency` kini peta tabel → tabel anak Currency (`T_FR_*` → `T_FR_CURRENCY`); `TestFlattenKodeDariHalamanCurrency` |
| `IsCedingConfirm` | jalankan K-069 (7b) "kolom sendiri"; tabel/tipe tidak disebut → usulan | ⛔ tidak lagi terbuang bersama `ViewSuggest`: `medanDiselamatkan` → penampung (fixture 34 = pengurai Python independen); usulan P4 tiga pilihan, menunggu work owner |
| Presisi | tetap tim inti, sesudah tabel B dibetulkan | tabel B diralat; `Diagnostik.LebihDari38Digit` kini digit bermakna |
| J-16 | tetap terbuka | — |

## Amandemen rancangan — butir 70 (diteruskan sesi `nusantarare-0f`, "setuju")

⛔ Workbook dan DDL draf (`D:\migrasi\RNM\OUTPUT\08-flat\`) **tidak** berubah — READ-ONLY. Amandemen ditulis di
`loader/amandemen.go` dan digabung ke skema bangkitan saat paket dimuat; DDL/migrasinya menunggu tiket 23 (ditahan).

| P | Amandemen | Sumber medan | Uji |
| :-: | --- | --- | --- |
| P1 | `T_COVERAGELIST.COVERAGE_INITIAL VARCHAR2(500)` | `CoverageList.CoverageInitial` | `TestFlattenAmandemen70` |
| P2–P3 | tabel `T_ADDITIONALSHIP` (anak `T_SHIP`, berulang): `ID`, `IDPEGA`, `COB_GROUP`, `PARENT_ID`, `SEQ_NO`, `ROW_UID`, `DWT`/`GRT`/`NRT` `VARCHAR2(50)`, `ADDITIONAL_SHIP_REF_ID VARCHAR2(50)`; jalur `CargoList/Ship/AdditionalShip` | `Ship/AdditionalShip[n]` (`ID` pola V-49) | idem; `TestFlattenSetiapJalur` (142 jalur) |
| P5 | `T_CURRENCYLIST.CURRENCY_REF_ID VARCHAR2(50)` | `CurrencyList.ID` (pola V-49) | `TestFlattenLipatan` |
| P6 | `T_FR_CURRENCYLIST.POLICY_TSI` — tipe ditulis `NUMBER` tanpa presisi sebagai **penanda** "menunggu tim inti" | FR `CurrencyList/Policy.TSI` (lipatan Policy, medan khusus) | `TestFlattenAmandemen70` |
| P7 | ~~penunjuk "SUDAH DIHAPUS" (R1) dan R3 induk langsung dibuang~~ → **dicabut butir 71** (A67, A68 — data membantah premis V-47 dan "sudah diwakili kolom urutan baris") |
| butir 72 | **48 kolom penunjuk teks mentah** `VARCHAR2(50)`, nama SNAKE_CASE dari FIELD ASLI (`INDEX_CARGO`, `IDX_LOCATION`, …), di 16 tabel — satu aturan untuk semua penunjuk (R1 13 · induk langsung 16 · leluhur 16 · berselisih 2 · `IdxPerson` 1). Kolom FK V-47 tetap kosong | `amandemenPenunjuk` | `TestFlattenPenunjukKeKolom`, `TestKasusPenunjukBukanPosisi` | lembar **Kandidat Hapus** (`Claude outputs`), kini ikut dibangkitkan (`penunjukKandidat`, 50 kunci) | `TestFlattenPenunjukV47` |

Skema sesudah amandemen: **79 tabel, 1.390 kolom, 149 jalur** (bangkitan tetap 78 / 1.329 / 148; butir 70 +13, butir 72
+48). Asal kolom: medan 865 · Flatten 327 · repository 166 · sengaja kosong 32. `[terverifikasi]` lembar Kandidat Hapus **berselisih** untuk dua
kunci — `T_FR_ANEKALIST.IdxOccupation` dan `T_FR_DEDUCTIBLELIST.IndexProperty` (R2 "SUDAH DIHAPUS" dan R3 "SUDAH JADI FK")
— tidak dipilih, tetap di penampung. R3b "DIPERTAHANKAN" (tiga kunci) sudah berkolom kunci, jadi tidak ke penampung.
*(Riwayat — A68 "15 penunjuk induk langsung dibuang walau berkolom FK" dicabut butir 71; butir 72 menyimpan semua
penunjuk apa adanya di kolom teks mentah.)*
`[dugaan]` tabel sasaran penunjuk dibaca dari namanya (`IdxLocation` → `T_LOCATIONLIST`, …), melanjutkan contoh V-47.
P4 (`IsCedingConfirm`) tetap di penampung: kolomnya di tabel lama `POOLDATA.HISTORYAKSEPTASIPRODUCTION`, **menunggu DBA**.

### Sisa isi penampung sesudah butir 72 — 02-10-2026

| Isi | Fixture | Korpus | Menunggu |
| --- | ---: | ---: | --- |
| `ViewSuggest[n].IsCedingConfirm` (P4) | 34 | 805 | DBA |
| *(artefak spasi pengubah XML→JSON agent — bukan data)* | 0 | 214 | — |

Nilai penunjuk yang kini berkolom: fixture 600, korpus 21.023 (= 6.133 R1 + 6.138 induk langsung + 8.723 leluhur + 25
berselisih + 4 `IdxPerson`).

### Sisa isi penampung sesudah butir 71 — 02-10-2026 *(riwayat)*

| Isi | Fixture | Korpus | Menunggu |
| --- | ---: | ---: | --- |
| `ViewSuggest[n].IsCedingConfirm` (P4) | 34 | 805 | DBA |
| penunjuk R1 indeks-diri (13 kunci) | 227 | 6.133 | keputusan per kunci (A67: korpus 39 beda) |
| penunjuk R3 sasaran induk langsung (16 kunci) | 219 | 6.138 | keputusan V-47 (A68: fixture 103 beda, korpus 8) |
| penunjuk R3 leluhur | 154 | 8.723 | aturan isi FK V-47 |
| penunjuk dua kunci berselisih | 0 | 25 | work owner |
| penunjuk tak tercantum di Kandidat Hapus: `T_FR_PERSONLIST.IdxPerson` | 0 | 4 | work owner |
| `CoverageInitial`, `AdditionalShip` | 0 | — | — (berkolom; fixture 104 dan 416 medan) |
| `CurrencyList.ID`, FR `Policy.TSI` | — (medannya **tidak ada** di fixture) | 0 (5 dan 15 kemunculan, berkolom) | — |
| **jumlah** | **634** | **21.828** | |

Rekonsiliasi: korpus R3 6.138 + 8.723 = 14.861 = seluruh penunjuk R3 yang terukur; 8.723 + 25 + 4 = 8.752 = angka
"penunjuk leluhur" butir 70 (ralat label). Rincian pengukuran A67/A68: register butir 71.

### Sisa isi penampung sesudah butir 70 — 02-10-2026 *(riwayat; digantikan tabel di atas)*

Jendela: 5 fixture; 115 contoh korpus (diubah ke JSON di luar repositori, 214 entri artefak spasi dari pengubah XML tidak
dihitung, lalu salinannya dihapus).

| Isi | Fixture | Korpus | Menunggu |
| --- | ---: | ---: | --- |
| `ViewSuggest[n].IsCedingConfirm` (P4) | 34 | 805 | DBA |
| penunjuk leluhur V-47 — terbesar `…CoverageList.IndexProperty`, `…DeductibleList.IndexProperty`/`IndexPropertyItem` | 154 | 8.752 | aturan isi FK V-47 (butir 4 lama) |
| ↳ termasuk dua kunci berstatus berselisih (`T_FR_ANEKALIST.IdxOccupation` 15; `T_FR_DEDUCTIBLELIST.IndexProperty` 10) | 0 | 25 | work owner |
| `CoverageInitial`, `AdditionalShip`, `CurrencyList.ID`, `FR Policy.TSI` | **0** | **0** | — (berkolom) |
| **jumlah** | **188** | **9.557** | |

*(Riwayat butir 70.)* Penunjuk yang kini dibuang menurut P7: fixture 446 (R1 indeks-diri 227 — A67; V-47 induk langsung 219 — A68), korpus
12.271 (6.133 / 6.138).
⛔ ADR-0023: selama penampung berisi, pemuatan produksi/fase 1 **belum** boleh dinyatakan selesai.

## Butir terbuka — menunggu work owner

1. ~~Usulan kolom penampung P1–P7~~ → ✅ butir 70. Tersisa: **DBA** untuk kolom `IsCedingConfirm` di
   `POOLDATA.HISTORYAKSEPTASIPRODUCTION` (P4); **aturan isi FK V-47** untuk penunjuk leluhur; status berselisih
   `T_FR_ANEKALIST.IdxOccupation` dan `T_FR_DEDUCTIBLELIST.IndexProperty` di lembar Kandidat Hapus; ~~konfirmasi A67 dan
   A68~~ → diputuskan butir 71. ~~Tersisa: apakah V-47 diubah dan R1 per kunci~~ → ✅ **butir 72**: V-47
   diubah, satu aturan — semua penunjuk berkolom teks mentah. Tersisa: **aturan isi FK V-47** (kolom FK tetap kosong).
2. **V-22 "ID CurrencyList berisi kode mata uang" tidak cocok dengan data.** Jendela: **korpus** saja — 116 unsur
   `CurrencyList` akar di 115 contoh, `ID` terisi 5, seluruhnya angka lima digit; `Name` kode tiga huruf 116/116. Di
   **fixture** `CurrencyList` **tidak punya** medan `ID`. `CURRENCY_CODE` diisi dari `Name` (BAHAN §0, A60); `ID` ke
   penampung (usulan P5).
3. **Induk ganda 12 atau 13:** Daftar Relasi 12, Jalur Sumber 13 (+`T_FR_POLICY`); ketiga tabel wadah 24-09 tidak ada di
   Daftar Relasi. Flatten mengisi `PARENT_TABLE` di ke-23 tabel yang punya kolomnya.
4. **Bentuk baris per tabel di `models`** (A65) — tiket 24.
5. ⚠️ **Ralat agent:** daftar medan korpus "tak ada di rancangan" yang dulu ditulis di sini (`LocationList.TableOfLimit`
   186, `…OfferFacIn` 22, `LocationList/Property.Property` 3, `…SurveyAgent` 2) **bukan data** — `[terverifikasi]` elemen
   XML kosong berisi spasi saja, diubah pengubah XML→JSON agent menjadi medan bernilai spasi. Data korpus yang sungguh
   tak terpetakan: `CurrencyList.ID` (5) dan `FacOfferList/CurrencyList/Policy.TSI` (15) — ditambah `IsCedingConfirm`
   (805), yang sejak butir 69 diselamatkan dari `ViewSuggest` ke penampung (usulan P4).
6. **Penjaga modul lain merah karena berkas ini:** `modul/claimlife/backend/repository` `TestKolomTakDibawaHanyaAdaDiKatalog`
   menyatakan lingkup "MODUL CLAIM LIFE" tetapi memindai seluruh `APP_RNM` (`akarModul = "../../../.."`), sehingga
   menuduh `skema_gen.go` dan `aturan.go` yang menyebut `STS_KONVERSI`/`TGL_KONVERSI` (kolom `T_WORK_POLIS`, V-48).
   ⛔ Work owner 02-10-2026: **"FOKUS KE NB FACIN SAJA"** — claimlife tidak disentuh; butir untuk pemilik claimlife.

## Comments

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan. Tidak ada nilai kasus di tiket ini.*
