# SPEC — Model Data dan DDL, modul Claim Non Prop

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09, rule termutakhir `pxUpdateDateTime = 2026-08-30`); `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**, terurai ke `pengetahuan/ddl/`; `pengetahuan/ddl/TABLE_PC_ASM_FW_GCNMFW_WORK.sql`; dan `MEMORI_PEMAHAMAN.MD` sebagai otoritas tertinggi untuk maksud bisnis dan istilah.
> Disusun 18 September 2026. Tidak ada pengetahuan umum tentang reasuransi maupun tentang Pega yang dipakai sebagai sumber.

**Lingkup: model data dan DDL saja.** Golang, React, endpoint, kontrak API, layar, service, repository, ORM — tidak ada di dokumen ini, dan tidak boleh ditambahkan ke dalamnya.

---

## Problem Statement

Petugas klaim hari ini bekerja di atas sistem yang **tidak menyimpan angkanya**. Seluruh nilai uang klaim — estimasi, alokasi per layer, premi pemulihan, nilai adjustment — tersimpan sebagai **teks di dalam satu kolom BLOB** (`PZPVSTREAM`), dan ketika nilai itu disalin keluar ke tabel bisnis, ia tetap teks: view `CLAIMXOL` mengeluarkan `TotalXOLGross`, `CNPReinstatement`, dan `KursIDR` sebagai `varchar2`.

Akibatnya bagi orang yang memakainya:

- **Tidak ada yang bisa menjawab "berapa total klaim tahun ini per mata uang" tanpa membongkar BLOB.** Dari 659 properti Pega, 549 tidak punya kolom di mana pun.
- **Angka yang sama bisa berbeda tergantung urutan tombol yang ditekan.** `TotalClaim` ditimpa tiga kali di dalam satu rule, dan rumus premi pemulihan membacanya di antara dua penimpaan (FINDING-007).
- **Hasil suntingan manual petugas hilang tanpa jejak** setiap kali alokasi dihitung ulang, karena penanda `.IsEditClaim` ditulis tetapi tidak pernah dibaca (FINDING-004).
- **Nilai bermata uang asing dibandingkan terhadap ambang kewenangan tanpa dikonversi** (FINDING-001), dan kurs yang tidak ditemukan dikembalikan sebagai `1` tanpa peringatan (FINDING-006).
- **Baris ganda bukan kemungkinan teoretis**: `OS_AKSEPTASI_KLAIM` tidak punya primary key maupun unique constraint, dan prosedurnya selalu `INSERT` karena logika upsert-nya dikomentari seluruhnya.

Setiap satu dari lima hal itu berakar pada hal yang sama: **nilai tidak punya kolomnya sendiri**, sehingga tidak ada tipe, tidak ada presisi, tidak ada constraint, dan tidak ada yang bisa menolak angka yang salah.

## Solution

Model data relasional di mana **setiap nilai punya kolom, tipe, dan presisi sendiri**, dan setiap aturan yang bisa ditegakkan basis data ditegakkan di sana — bukan diserahkan pada ingatan penulis kode.

Bagi orang yang memakainya, yang berubah:

- Pertanyaan tentang angka dijawab dengan satu kueri, bukan dengan membongkar BLOB.
- Alokasi yang dihitung ulang menghasilkan keadaan yang identik (ADR-0011), sehingga angka tidak lagi bergantung urutan klik.
- Nilai hasil suntingan manual bertahan dan **terlihat sebagai suntingan**, lengkap dengan pelaku dan waktunya (ADR-0008).
- Nilai uang selalu berpasangan dengan mata uangnya, dan nilai IDR-nya disimpan terpisah beserta kurs, tanggal kurs, dan sumbernya (ADR-0007).
- Kurs yang tidak ada **tidak menghasilkan angka** — kolom IDR dibiarkan NULL dan barisnya ditandai menunggu kurs (ADR-0014, ADR-0019). NULL berarti belum dihitung; nol berarti sudah dihitung dan hasilnya nol.
- Satu akseptasi per klaim per layer per mata uang, ditegakkan `UNIQUE` (ADR-0024).

---

## User Stories

**Registrasi dan identitas klaim**

1. Sebagai petugas registrasi, saya ingin satu klaim menyimpan tepat satu Tanggal Kejadian yang tidak berubah sepanjang umurnya, agar seluruh perhitungan di bawahnya punya satu titik acuan waktu.
2. Sebagai petugas registrasi, saya ingin Nomor Polis dan Nomor Polis Cedant tersimpan sebagai dua kolom berdampingan, agar saya tidak perlu menebak nomor mana milik siapa.
3. Sebagai petugas registrasi, saya ingin Nomor Klaim dijamin unik oleh basis data, agar dua klaim tidak pernah berbagi nomor yang sama seperti di sistem lama.
4. Sebagai petugas registrasi, saya ingin masa berlaku treaty tersimpan sebagai tanggal dengan batas inklusif, agar klaim yang jatuh tepat di hari terakhir tidak tertolak karena perbandingan teks.
5. Sebagai petugas registrasi, saya ingin Penyebab Kerugian tersimpan sebagai rujukan ke daftar induk, bukan sebagai teks bebas, agar laporan per penyebab dapat dijumlahkan.
6. Sebagai kepala bagian klaim, saya ingin tahu klaim mana yang masuk lewat penjaga tanggal yang keliru di sistem lama, agar saya dapat meninjaunya tanpa menebak.

**Nilai uang dan mata uang**

7. Sebagai petugas akseptasi, saya ingin setiap nilai uang tersimpan bersama mata uangnya dalam satu baris, agar tidak ada angka yang satuannya harus ditebak.
8. Sebagai petugas akseptasi, saya ingin nilai IDR tersimpan terpisah dari nilai mata uang asli, agar keduanya dapat dibandingkan tanpa menghitung ulang.
9. Sebagai petugas akseptasi, saya ingin kurs yang dipakai tersimpan beserta tanggal dan sumbernya, agar angka konversi dapat ditelusuri kembali.
10. Sebagai petugas akseptasi, saya ingin baris yang kursnya belum ada ditandai dengan jelas dan kolom IDR-nya kosong, agar saya tahu angka itu belum jadi, bukan bernilai nol.
11. Sebagai akuntansi, saya ingin membedakan nilai yang belum dihitung dari nilai yang hasilnya nol, agar laporan tidak mencampur keduanya.
12. Sebagai akuntansi, saya ingin ambang kewenangan dibandingkan terhadap nilai IDR, agar nilai bermata uang asing tidak lolos hanya karena angkanya terlihat kecil.

**Alokasi kerugian**

13. Sebagai petugas klaim, saya ingin alokasi per Layer tersimpan di tabel tersendiri dari Retensi Cedant, agar setiap kueri menyatakan secara eksplisit apakah Retensi Cedant ikut dihitung.
14. Sebagai petugas klaim, saya ingin urutan pengisian layer tersimpan sebagai data, agar hasil alokasi dapat dibaca ulang tanpa mengetahui urutan tombol yang ditekan petugas.
15. Sebagai petugas klaim, saya ingin menjalankan ulang perhitungan alokasi dan mendapatkan keadaan yang identik, agar saya berani menghitung ulang tanpa takut merusak angka.
16. Sebagai petugas klaim, saya ingin nilai yang saya sunting manual bertahan terhadap hitung ulang, agar pekerjaan saya tidak hilang diam-diam.
17. Sebagai pemeriksa, saya ingin melihat nilai mana yang hasil hitung dan mana yang hasil suntingan, beserta siapa dan kapan, agar selisih dapat dijelaskan.
18. Sebagai petugas klaim, saya ingin Limit Layer dan Premi Deposit tersimpan pada baris alokasi, agar dasar perhitungan premi pemulihan tetap terbaca setelah treaty berubah.

**Akseptasi dan adjustment**

19. Sebagai petugas akseptasi, saya ingin basis data menolak akseptasi kedua untuk kombinasi klaim, layer, dan mata uang yang sama, agar baris ganda tidak lahir seperti di sistem lama.
20. Sebagai petugas akseptasi, saya ingin akseptasi berupa satu baris dengan keadaan, bukan dua tabel, agar riwayatnya tidak terpecah.
21. Sebagai petugas pembayaran, saya ingin satu klaim dapat memiliki banyak Adjustment, agar pembayaran bertahap tercatat sebagai transaksi terpisah.
22. Sebagai petugas pembayaran, saya ingin rekening penerima tersimpan sebagai baris tersendiri, agar satu Adjustment dapat menunjuk rekening tanpa menyalin datanya.
23. Sebagai akuntansi, saya ingin Nomor Akseptasi tersimpan pada Adjustment yang menerbitkannya, agar instruksi bayar dapat ditelusuri ke keputusannya.
24. Sebagai pemeriksa, saya ingin keputusan alur tersimpan sebagai field terstruktur, bukan dibaca dari teks komentar, agar satu salah ketik tidak mengubah jalur.

**Premi pemulihan**

25. Sebagai petugas klaim, saya ingin premi pemulihan tersimpan per layer beserta seluruh masukan rumusnya, agar angkanya dapat dihitung ulang dan diperiksa.
26. Sebagai akuntansi, saya ingin tahu nilai mana yang dipakai rumus premi pemulihan, agar selisih terhadap sistem lama dapat dijelaskan, bukan sekadar ditemukan.

**Dokumen**

27. Sebagai petugas klaim, saya ingin dokumen klaim dirujuk lewat metadata dan URL, bukan disalin ke basis data, agar tidak ada berkas fisik yang harus dipindahkan.
28. Sebagai petugas klaim, saya ingin masa berlaku URL dokumen tersimpan, agar sistem tahu kapan URL itu perlu disegarkan.

**Tutup buku dan tarif**

29. Sebagai akuntansi, saya ingin tanggal tutup buku tersimpan bertanggal berlaku, agar perubahan tanggal tidak mengubah tafsir atas transaksi lama.
30. Sebagai akuntansi, saya ingin tarif pajak dan brokerage tersimpan sebagai data bertanggal berlaku, bukan tertanam di kode, agar perubahan tarif tidak memerlukan perubahan program.
31. Sebagai kepala bagian, saya ingin nilai yang sistem lama patok per klaim punya tempatnya sendiri sebagai data, agar tambalan tidak lahir kembali.

**Migrasi dan paritas**

32. Sebagai penanggung jawab migrasi, saya ingin setiap baris baru dapat ditunjuk balik ke barisnya di sistem lama, agar hasil sistem baru dapat dibandingkan baris per baris.
33. Sebagai penanggung jawab migrasi, saya ingin nilai lama yang tidak lolos aturan penguraian punya tempat tercatat, agar tidak ada angka yang hilang tanpa jejak.
34. Sebagai penanggung jawab migrasi, saya ingin jembatan ke sistem lama berada di tabel terpisah, agar ia dapat dijatuhkan setelah paritas diterima tanpa menyentuh tabel yang hidup.
35. Sebagai penanggung jawab migrasi, saya ingin klaim tertutup ikut termigrasi apa adanya dan ditandai, agar riwayat tidak terpotong di tanggal cutover.

**Kepemilikan dan batas**

36. Sebagai pemilik modul, saya ingin hanya klaim, akseptasi, alokasi, dan adjustment yang dimiliki skema ini, agar batas modul terbaca dari skemanya sendiri.
37. Sebagai pemilik modul, saya ingin sistem hilir membaca lewat view, bukan menyalin tabel, agar hanya ada satu bentuk kanonik.
38. Sebagai pemilik modul, saya ingin Arasapas dan kasir tetap menerima bentuk lama lewat lapisan view, agar kontrak yang tidak berubah tidak memaksa perubahan di model kanonik.
39. Sebagai pemilik modul, saya ingin hanya akun aplikasi yang memegang hak tulis, agar satu pintu tulis tidak bergantung pada kesepakatan.

---

## Implementation Decisions

### 0. Definisi **Layer** — ditetapkan sekali, dipakai sama di seluruh spec

Sistem lama memakai **dua satuan berbeda** untuk hal yang sama-sama disebut "layer", dan keduanya tidak sepadan. Selama itu tidak dinyatakan, `ALOKASI_LAYER` dan `AKSEPTASI` akan diam-diam berbeda butir.

| Satuan | Dibentuk dari | Dipakai sebagai |
|---|---|---|
| **Baris batas treaty** — halus | `Layer`, `LayerType`, `LayerPart`, `LayerPartType`, disalin apa adanya dari `TreatyInMaster.Limits[]` | butir alokasi: `SpreadingRisk` menulis keempatnya; `ReinstatementList` membawa `Layer` + `LayerPart` |
| **Jenis reasuransi** — kasar | hasil lookup ke `REINSURANCETYPE`, dengan probe yang **hanya** memakai `Layer` dan `LayerType` | kunci akseptasi |

**Bukti probe** (`Activity\CountLossAllocation_act.xml`, `InputSpreading.CARI1`) — EVIDENCED:

```
"XL " + @substring(.Layer,0,1)
      + @if(@substring(.Layer,0,1)=="1","ST ",@if(...=="2","ND ",@if(...=="3","RD ","TH ")))
      + @if(@toUpperCase(.LayerType)=="SUBLAYER","SUB LAYER",@toUpperCase(.LayerType))
```

`LayerPart` **tidak ada di dalamnya.** Seluruh bagian dari satu layer memetakan ke satu hasil lookup yang sama.

**Bukti hasil lookup menjadi nama, bukan rumusnya** — EVIDENCED, dan ini mengoreksi pembacaan pertama saya atas `MEMORI_PEMAHAMAN.MD` §6.2:

| Penetapan | Nilai |
|---|---|
| `SpreadingRisk(<LAST>).TreatyName` | `ListSpreading.pxResults(1).Note` — atau literal `"UR"` |
| `SpreadingRisk(<APPEND>).TreatyType` | `ListSpreading.pxResults(1).ID` — atau literal `"UR"` |

Jadi string `"XL 1ST LAYER"` adalah **kunci pencarian**, bukan nama yang disimpan.

**Bukti satuan yang dipakai kunci akseptasi** (`Activity\SaveDataToOSAksep_Act.xml`) — EVIDENCED:

| Parameter | Diisi dari |
|---|---|
| `InputParamOs.TypeLoss` | `.TreatyName` |
| `InputParamOs.TypeLossID` | `.TreatyType` |

dan `RDBList\SaveDataToOsAkseptasiNP.xml` meneruskan `{InputParamOs.TypeLoss}` sebagai parameter `LayerT` ke `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP`, yang blok kuncinya berbunyi `a.data_json.TypeLoss = LayerT`.

#### Yang ditetapkan untuk spec ini

> **Layer** = **baris batas treaty**, yaitu satuan halus: `Layer`, `LayerType`, `LayerPart`, `LayerPartType`.

Alasannya glosarium: *Layer adalah satu lapisan proteksi yang dibatasi batas bawah dan batas atas sendiri.* Batas itu melekat pada baris `Limits[]`, bukan pada hasil lookup jenis reasuransi. Satuan kasar **tidak diberi nama "layer"** di dokumen ini; ia dirujuk sebagai **jenis reasuransi**, dan hidup sebagai atribut pada baris alokasi — bukan sebagai kunci.

Dengan definisi itu, keduanya memakai butir yang sama:

| Tabel | Kunci alami |
|---|---|
| `ALOKASI_LAYER` | (klaim, **layer**, mata uang) |
| `AKSEPTASI` | (klaim, **layer**, mata uang) — ADR-0024 |

#### Dikuatkan dari tiga arah, bukan satu — EVIDENCED

| Bukti | Isinya |
|---|---|
| Sasaran lookup punya kolomnya | `pengetahuan/ddl/TABLE_REINSURANCETYPE.sql` berkolom **`ID`** dan **`NOTE`** — persis pasangan yang dibaca `ListSpreading.pxResults(1).ID` menjadi `TreatyType` dan `.Note` menjadi `TreatyName` |
| Penanda jenis reasuransi muncul di `TreatyName` | `@contains(.TreatyName,"R/I")` muncul **4x** (FINDING-002 §6.1). `"R/I"` tidak pernah ada di rumus probe — bila `TreatyName` berisi `"XL 1ST LAYER"`, keempat kondisi itu **tidak akan pernah cocok** |
| Bentuk nilainya sendiri | `MEMORI_PEMAHAMAN.MD` §6.2 Langkah 10 memakai `"QS (OR)"` menjadi `"10028"` dan `"QS (R/I)"` menjadi `"10004"`. Itu bentuk `REINSURANCETYPE`, bukan bentuk `"XL <n>TH LAYER"` |

Koreksi atas `MEMORI_PEMAHAMAN.MD` §6.2 tercatat sebagai **D8** di `_selesai/OPEN-QUESTIONS.md`.

#### `"UR"` adalah artefak, bukan nilai domain — dan itu menghapus satu kolom

Baris Retensi Cedant menulis `TreatyName = "UR"` dan `TreatyType = "UR"` **secara literal, melewati lookup sepenuhnya** (FINDING-003 bagian 2). EVIDENCED.

Bila kolom itu identitas layer, nilai literal `"UR"` mustahil — ia bukan nomor layer. Bila ia jenis reasuransi, `"UR"` memang bukan salah satunya: ia penanda yang menumpang di kolom yang salah, persis penyatuan yang dibatalkan ADR-0010.

> **`RETENSI_CEDANT` tidak memerlukan kolom jenis reasuransi sama sekali.** DECIDED(ADR-0010). Setelah dipisahkan menjadi tabelnya sendiri, tidak ada yang perlu dibedakan — seluruh barisnya Retensi Cedant.

#### Bagaimana view kompatibilitas memetakan halus ke kasar

Satuan kasar adalah **identitas yang dilihat sistem hilir**, dan rantainya EVIDENCED sampai ke luar: `TreatyName` ke `InputParamOs.TypeLoss` ke parameter `LayerT` ke `a.data_json.TypeLoss`, lalu dibongkar `CLAIMXOL` sebagai `$.CNPLayerList[*].XOL`.

Arah risikonya perlu dinyatakan terang: **kunci yang lebih halus adalah constraint yang lebih lemah** — ia mengizinkan lebih banyak baris. Jadi risikonya bukan "akseptasi sah tertolak", melainkan "baris yang sistem lama anggap satu menjadi dua, dan hilir tidak tahu".

**Putusan — DECIDED-TEKNIS**: view kompatibilitas **beragregasi**, dan agregasinya dinyatakan, bukan diam-diam.

| Perlakuan | Kolom |
|---|---|
| **Dijumlahkan** atas seluruh baris dengan (klaim, jenis reasuransi, mata uang) sama | nilai kerugian teralokasi, biaya penilaian, salvage, biaya lain, premi pemulihan — seluruh besaran aditif |
| **Tidak boleh dijumlahkan**; harus sama di seluruh baris yang digabung, dan bila tidak sama view **menolak menghasilkan baris** | Limit Layer, Premi Deposit, persentase premi pemulihan, kurs, porsi reasuradur — besaran per-layer, bukan besaran aditif |
| **Tidak dibawa** ke view | `LayerPart`, `LayerPartType` — justru dimensi yang dilebur |

Alasan memilih agregasi ketimbang menegakkan kunci kasar sebagai `UNIQUE` kedua: `UNIQUE` kedua akan **menolak baris alokasi yang sah** di sisi kanonik hanya karena hilir tidak dapat membedakannya. Itu menjadikan keterbatasan bentuk lama sebagai aturan bisnis sistem baru. Agregasi menempatkan keterbatasan itu di tempatnya — di lapisan view, yang memang ada untuk melayani kontrak lama (ADR-0023).

**Batas klaim**: bahwa kolom di baris kedua *bisa* berbeda pada besaran per-layer adalah kemungkinan struktural; **TIDAK DITEMUKAN** bukti bahwa ia pernah terjadi. Penolakan view di atas adalah penjaga, bukan perkiraan frekuensi.

#### Akibatnya terhadap ADR-0024 — dinyatakan, tidak diputuskan sendiri

Kunci ini **lebih sempit daripada bukti yang melahirkan ADR-0024**. Blok yang dikomentari itu memakai `TypeLoss`, yaitu satuan kasar. Setiap baris ganda yang ditolak kunci lama juga ditolak kunci baru; tidak sebaliknya.

Keduanya berimpit **bila dan hanya bila** `LayerPart` tidak pernah berbeda di dalam satu pasangan `Layer` + `LayerType`.

**`REQ-018` tidak dapat menjawabnya.** Ia mencacah baris ganda `(CASEID, TypeLoss, Currency)` — kunci **kasar** — dan `LAYERPART` tidak pernah sampai ke `OS_AKSEPTASI_KLAIM`. Mencarinya dari sisi klaim berarti membongkar blob, karena `SpreadingRisk` IN-BLOB (`BLUEPRINT.md` §13.5).

Pertanyaannya tidak perlu lewat sana. `LayerPart` berasal dari struktur treaty, dan **kolomnya nyata** — EVIDENCED dari DDL:

| Tabel | Kolom |
|---|---|
| `POOLDATA.PROPORTIONALARRG` | `LAYER`, `LAYERPART`, `LAYERPARTTYPE`, `LAYERTYPE` — seluruhnya `VARCHAR2(1000)` |
| `POOLDATA.TREATYINDETAIL` | keempatnya, sama |

Satu `GROUP BY ... HAVING COUNT(DISTINCT LAYERPART) > 1` menjawabnya: read-only, tanpa data nasabah, tanpa membongkar blob. Diangkat sebagai **REQ-033**, bernomor sendiri, **BLOCKER, OPEN** — berhenti di situ.

**Yang DDL sudah jawab, dan yang tidak.** Terjawab: keempat kolom ada, jadi pertanyaannya dapat dijalankan sebagai satu kueri. **Tidak terjawab**: `PROPORTIONALARRG`, `TREATYINDETAIL`, dan `REINSURANCETYPE` **tidak punya satu pun constraint** — tanpa primary key, tanpa unique. Definisi tabelnya tidak melarang apa pun, jadi hanya data yang dapat menjawab.

Dan ADR-0024 sudah meramalkan bentuk jawabannya: kelompok ketiga, *"baris berbeda, tanpa pola waktu — berarti ada dimensi pembeda yang belum tertangkap"*. **`LayerPart` adalah kandidat dimensi itu**, dinamai sebelum datanya dilihat.

**Batas klaim**: bahwa `LayerPart` *dapat* membedakan dua baris alokasi — EVIDENCED, ia kolom tersendiri di `SpreadingRisk` dan di dua tabel struktur treaty. Bahwa ia *pernah* berbeda pada data nyata — **TIDAK DITEMUKAN** di XML maupun DDL.

---

### 1. Entitas yang dimiliki, beserta kunci alaminya

Batasnya ditetapkan ADR-0026: modul ini memiliki **klaim, akseptasi, alokasi, dan adjustment** — itu saja. Polis, treaty, cedant, mata uang, dokumen, dan daftar Komite dibaca dari modul lain.

| # | Entitas | Berdiri sendiri karena | Kunci alami |
|---|---|---|---|
| 1 | **KLAIM** | ADR-0001: aggregate root, satu klaim satu kejadian kerugian. Diperkuat sapuan 114 activity — **tidak satu pun `Property-Set` menulis `.DateOfLoss`** (`BLUEPRINT.md` §2.2) | `NOMOR_KLAIM` — lihat catatan (a) |
| 2 | **NILAI_KLAIM_MATA_UANG** | Satu klaim membawa nilai di lebih dari satu mata uang; sistem lama menyimpannya di PageList `ListClaimAmount` (37 rujukan) | (klaim, mata uang) |
| 3 | **ALOKASI_LAYER** | ADR-0010: Layer dan Retensi Cedant dua entitas terpisah | (klaim, **layer**, mata uang) — layer = komposit empat field, bagian 0 |
| 4 | **RETENSI_CEDANT** | ADR-0010. Sistem lama menaruhnya sebagai baris ber-`TreatyName="UR"` di dalam `SpreadingRisk`; bentuk itu **tidak dipindahkan**, dan **tanpa kolom jenis reasuransi** (bagian 0) | (klaim, mata uang) |
| 5 | **AKSEPTASI** | ADR-0015: satu entitas dengan keadaan, bukan dua tabel | **(klaim, layer, mata uang)** — ADR-0024, definisi layer yang sama persis, ditegakkan `UNIQUE` |
| 6 | **ADJUSTMENT** | `CONTEXT.md`: unit transaksi pembayaran klaim; satu klaim banyak Adjustment. `MEMORI_PEMAHAMAN.MD` §4.3 | (klaim, nomor urut) — lihat catatan (b) |
| 7 | **PREMI_PEMULIHAN** | `MEMORI_PEMAHAMAN.MD` §6.3: dihitung per layer, disimpan ke `ReinstatementList[]` dengan seluruh masukan rumusnya | (klaim, **layer**, mata uang) — **berubah** setelah putusan bagian 0; `ReinstatementList` membawa `Layer` dan `LayerPart`, keduanya bagian dari komposit yang sama |
| 8 | **REKENING_PENERIMA** | PageList `ReceiverClaim` (101 rujukan) dengan varian rekening kedua | (klaim, nomor urut) |
| 9 | **OBJEK_PERTANGGUNGAN** | PageList `InterestList`/`ObjectList` | (klaim, nomor urut) |
| 10 | **KRONOLOGI_KLAIM** | ADR-0009 + `CONTEXT.md` **Kronologi**: pelaku, peran, waktu, keterangan | (klaim, nomor urut) |
| 11 | **DOKUMEN_KLAIM** | ADR-0027: metadata dan URL berbatas waktu, berkasnya di Google Cloud Storage | (klaim, nomor urut) |

Tabel pendukung, bukan entitas bisnis:

| # | Tabel | Dasar |
|---|---|---|
| 12 | **TUTUP_BUKU** | ADR-0025, bertanggal berlaku, dengan kolom lingkup yang disediakan sekarang meski belum ada yang memerlukannya |
| 13 | **TARIF_BERLAKU** | ADR-0004: tarif 2,5% / 2% / 2,2% tertanam di `SetPPNPPH` tanpa tanggal berlaku. Nilainya **belum dikonfirmasi** (A8, ASK-AKUNTANSI no. 5); yang dibuat adalah tempatnya |
| 14 | **KOREKSI_NILAI** | ADR-0004: 29 langkah, 18 ekspresi, 8 rule menambal per-case. Model data harus **punya tempat** bagi Premi Pemulihan yang dipatok, pemetaan Biaya Penilaian per porsi, dan batas Retensi Cedant alternatif — bila tidak ada tempatnya, tambalan lahir lagi |
| 15 | **KLAIM_PENJAGA_TANGGAL** | ADR-0018: klaim terdampak dimigrasi apa adanya **dan didaftar**; daftar itu perlu tempat menetap |
| 16 | **MIGRASI_KORELASI** | Jembatan ke sistem lama — lihat bagian 5 |
| 17 | **MIGRASI_PENDARATAN** | Satu-satunya tempat bentuk lama boleh masuk utuh, termasuk JSON (ADR-0028) |
| 18 | **MIGRASI_NILAI_DITOLAK** | ADR-0014: nilai yang tidak lolos aturan penguraian ditolak dan **masuk daftar**, tidak dibulatkan |

**(a) Kunci alami KLAIM belum dapat dikunci.** `BLUEPRINT.md` §2.3: tidak ada rule yang mencegah duplikat `PolicyNo` + `DateOfLoss` + `IDMaster`, dan itu *disimpulkan dari ketiadaan rule penolakan, bukan dari adanya rule*. Sementara itu `PYID` di tabel work **tanpa unique constraint dan tanpa index tersendiri**. `NOMOR_KLAIM` diusulkan `UNIQUE`; apakah data lama melanggarnya diukur REQ-018 (BLOCKER, OPEN).

**(b) Nomor urut Adjustment menggantikan `IndexObject`.** Sistem lama memakai **posisi numerik di dalam `AdjustmentList`** sebagai penunjuk dari case Komite (`BLUEPRINT.md` §8.4) — bukan surrogate key. Posisi berubah bila daftar disusun ulang. Nomor urut baru bersifat tetap; pemetaan dari posisi lama ke nomor baru disimpan di `MIGRASI_KORELASI`.

### 2. Relasi dan kardinalitas

```
KLAIM 1 ──n NILAI_KLAIM_MATA_UANG
      1 ──n ALOKASI_LAYER
      1 ──0..n RETENSI_CEDANT        (satu per mata uang)
      1 ──n AKSEPTASI
      1 ──n ADJUSTMENT
      1 ──n PREMI_PEMULIHAN
      1 ──n REKENING_PENERIMA
      1 ──n OBJEK_PERTANGGUNGAN
      1 ──n KRONOLOGI_KLAIM
      1 ──n DOKUMEN_KLAIM

ADJUSTMENT n ──1 REKENING_PENERIMA    (rujukan, bukan salinan)
ADJUSTMENT 1 ──n AKSEPTASI            (lihat catatan)
```

**Arah sarang layer→mata uang ditetapkan oleh bentuk penyimpanan lama**, bukan oleh selera: `CLAIMXOL` membongkar `$.CNPLayerList[*]` yang **membungkus** `$.CNPCurrencyList[*]`. Layer di luar, mata uang di dalam. Bukti ini ada di badan ADR-0024.

Relasi `ADJUSTMENT → AKSEPTASI` **belum dapat dikunci arahnya**: di sistem lama satu Adjustment membawa `CNPLayerList[]` berisi `CNPCurrencyList[]`, sehingga satu Adjustment menyentuh banyak (layer, mata uang). Apakah itu berarti satu Adjustment menerbitkan banyak baris akseptasi, atau sebaliknya, bergantung pada REQ-018 dan REQ-023.

### 3. Presisi — tabel domain, seluruhnya **USULAN**

Tiga fakta bertemu dan spec ini tidak melangkahinya:

1. Seluruh nilai uang lama **IN-BLOB**, terkonfirmasi dari DDL tabel work: 19 kolom di luar bawaan Pega, **tidak satu pun nilai uang**. Tidak ada tipe kolom lama yang bisa disalin.
2. Bahkan ketika nilai dikeluarkan lewat view, tipenya `varchar2` (`CLAIMXOL`).
3. Dari 673 ekspresi aritmetika, **518 (77%) tanpa skala eksplisit**; yang menyatakan skala memakai 20 · 10 · 8 · 5 · 4 · 2 · 0, dan **skalanya mengikuti modul, bukan jenis nilai**.

Maka setiap angka di bawah ini **usulan dengan alasan**, bukan keputusan. Setiap kolom di kamus merujuk kelompoknya; menaikkan ADR-0003 nanti menjadi sapuan mekanis atas satu tabel ini.

| Kelompok | Skala sistem lama | Presisi diusulkan | Yang hilang bila diambil lebih rendah |
|---|---|---|---|
| **U1 — nilai uang antara** (alokasi per layer, nilai per mata uang, salvage, biaya penilaian) | `@divide(...,20)` di seluruh rule inti klaim | `NUMBER(23,6)` | Pada 2 desimal, selisih premi pemulihan terukur **Rp 1,26 juta per klaim**; pada 4 desimal masih **Rp 38 ribu**. Sumber: tabel presisi di `ASK-AKUNTANSI.md` no. 1 |
| **P1 — persentase dan porsi** (`ClaimPercentage`, `RNMShare`, `SharePercentage`, `PctProrateClaim`) | 20, dan 5 untuk `.PctProrateClaim` | `NUMBER(11,8)` | Porsi dipakai sebagai **pengali berantai** — alokasi dibagi porsi lalu dikali porsi lagi. Presisi rendah menggandakan galatnya di tiap tahap |
| **P2 — tarif pajak dan brokerage** | `@divide(...,8)` di `SetPPNPPH` | `NUMBER(11,8)` | Skala 8 dipakai sistem lama secara konsisten di keempat tarif; menurunkannya mengubah angka yang selama ini diterima akuntansi |
| **K1 — kurs** | tidak dinyatakan; `GETCURRENCYSTANDARD` mengembalikan `NUMBER` tanpa presisi | `NUMBER(19,8)` | Kurs mengalikan seluruh nilai IDR. Galat kurs menyebar ke setiap turunannya |
| **L1 — limit dan premi deposit** (`CNPLimit`, `CNPMDP`) | 20 | `NUMBER(23,6)` | Keduanya **penyebut** dalam rumus premi pemulihan; galat penyebut tidak dapat dikoreksi di hilir |
| **C1 — cacah dan nomor urut** | — | `NUMBER(9)` | — |

**Yang tidak diusulkan dan sengaja dibiarkan terbuka**: skala 10 yang dipakai `CountClaimTNP_Act` dan `CountLossAllocation_act` berdampingan dengan skala 20, dengan perimbangan 14:14 yang *terlalu rapi untuk kebetulan* (`_selesai/OPEN-QUESTIONS.md` E1). Sampai itu dijelaskan, tidak ada kelompok presisi yang dibangun di atasnya.

### 4. Keadaan di tingkat baris — ADR-0019, tiga lapis yang tidak dicampur

| Lapis | Bentuk | Aturan |
|---|---|---|
| 1 | `NULL` pada kolom nilai | **belum dihitung**. Nol berarti sudah dihitung dan hasilnya nol. Tidak pernah dipakai bergantian |
| 2 | Kolom keadaan **di tingkat baris** | satu kolom per baris, bukan satu kolom per nilai. Nilainya sekurangnya: lengkap · menunggu kurs · gagal urai |
| 3 | Penanda **per-nilai** | hanya di satu tempat: konversi mata uang — kurs, tanggal kurs, sumber kurs |

**Penularan keadaan "menunggu kurs" dimodelkan, bukan sekadar disebut.** ADR-0014 mengaturnya: bila kurs tidak ada, nilai IDR tidak dihitung, dan **setiap baris turunan yang membutuhkan nilai IDR itu ikut bertanda menunggu kurs**. Karena itu kolom keadaan hadir di seluruh tabel yang membawa nilai IDR turunan — `ALOKASI_LAYER`, `RETENSI_CEDANT`, `AKSEPTASI`, `ADJUSTMENT`, `PREMI_PEMULIHAN` — bukan hanya di tabel sumbernya.

Alasannya konkret: `GETCURRENCYSTANDARD` mengembalikan **`1`** ketika kurs tidak ditemukan (FINDING-006). Angka 1 tidak dapat dibedakan dari kurs yang sah. Di model baru, ketiadaan kurs menghasilkan **ketiadaan angka**, bukan angka yang menyamar.

### 5. `MIGRASI_KORELASI` — tabel terpisah, bukan kolom di tabel bisnis

Satu baris per baris kanonik yang lahir dari migrasi:

| Kolom | Isi |
|---|---|
| tabel tujuan | nama tabel kanonik |
| kunci tujuan | nilai kunci baris baru |
| pengenal lama, **komposit** | `pzInsKey`, `pyID`, `CASEID`, posisi `IndexObject`, urutan baris `SpreadingRisk` — apa adanya, tanpa penyatuan |
| sumber lama | tabel atau PageList asalnya |
| waktu migrasi | — |

Tiga alasan ia terpisah: ADR-0023 menuntut satu bentuk kanonik dan kolom korelasi bukan bagiannya; kolom bermasa pakai terbatas di dalam tabel permanen adalah pola yang melahirkan tambalan di sistem lama; dan menjatuhkan satu tabel jauh lebih murah daripada membuang kolom dari tabel yang sudah hidup.

**Bentuk kolom pengenal lamanya belum dapat dikunci.** `PYID` tidak dijamin unik dan `OS_AKSEPTASI_KLAIM` tanpa kunci sama sekali — jadi kolom itu harus menampung **komposit yang benar-benar dipakai**, dan komposit mana yang dipakai baru diketahui setelah REQ-018 (BLOCKER, OPEN). Tidak diselesaikan dengan tebakan.

### 6. Constraint yang menegakkan ADR

| Aturan | Bentuk | ADR |
|---|---|---|
| Satu akseptasi per klaim per layer per mata uang | `UNIQUE (klaim, layer, mata uang)` pada `AKSEPTASI` | 0024 |
| Setiap kolom uang punya penunjuk mata uang | `NOT NULL` pada kolom mata uang di setiap tabel bernilai uang, ditambah `CHECK` bahwa nilai IDR dan kurs sama-sama terisi atau sama-sama kosong | 0007 |
| Batas masa berlaku treaty inklusif | `CHECK (tanggal_mulai <= tanggal_akhir)`, tipe `DATE` | 0022 |
| Nomor Klaim unik | `UNIQUE (NOMOR_KLAIM)` — **menunggu REQ-018** | 0021 |
| Tutup buku bertanggal berlaku | `UNIQUE (lingkup, tanggal_berlaku)` | 0025 |
| Tidak ada JSON di data yang dimiliki | tidak ada kolom `CLOB IS JSON` di tabel kanonik; hanya di `MIGRASI_PENDARATAN` | 0028 |
| Satu pintu tulis | `GRANT`/`REVOKE` sebagai objek DDL: akun aplikasi memegang tulis atas tabel kanonik, seluruh akun lain termasuk `POOLDATA` hanya `SELECT` atas view | 0017, 0028 |

**Yang tidak dapat ditegakkan basis data, dan karena itu jatuh ke aplikasi:**

| Aturan | Sebab |
|---|---|
| Hitung ulang idempoten (ADR-0011) | Sifat prosedur, bukan sifat baris. Yang dapat ditegakkan hanya `UNIQUE` yang mencegah baris ganda lahir dari hitung ulang |
| Suntingan manual bertahan terhadap hitung ulang (ADR-0008) | Basis data dapat menyimpan penanda, pelaku, dan waktu, tetapi tidak dapat mencegah proses hitung menimpanya |
| Perhitungan sebagai turunan data (ADR-0013) | Tidak berbentuk constraint |
| Satu klaim satu kejadian (ADR-0001) | Dapat ditegakkan hanya bila kunci alami klaim diketahui — dan ia belum |
| NULL ≠ nol (ADR-0019) | `CHECK` dapat mencegah kombinasi yang mustahil, tetapi makna NULL adalah kesepakatan, bukan constraint |

### 7. View untuk sistem hilir — ADR-0023

Sistem hilir diberi **view**, bukan salinan. Karena skema baru duduk di instance yang sama dengan `POOLDATA` (ADR-0028, `proposed`), view lintas skema adalah view sungguhan — tanpa salinan, tanpa database link.

| View | Untuk | Bentuk |
|---|---|---|
| Kompatibilitas akseptasi | Arasapas dan kasir | Mempertahankan bentuk lama **termasuk `CASEID`**. ADR-0021 memensiunkan `CASEID` dari model internal dan percakapan, dengan pengecualian tepat di sini: kontrak integrasi tidak berubah. `CASEID` hidup **hanya di lapisan view**, bukan sebagai kolom kanonik |
| Rekap per klaim per mata uang | pelaporan | Agregasi `ALOKASI_LAYER` + `RETENSI_CEDANT`, dengan Retensi Cedant **sebagai kolom terpisah**, sehingga penjumlahan tidak pernah mencampurnya diam-diam |
| Paritas shadow-run | ADR-0005 | Lihat bagian 8 |

### 8. Kesiapan shadow-run — ADR-0005

Model ini **memungkinkan** perbandingan baris per baris, dengan satu syarat dan satu batas.

Syaratnya `MIGRASI_KORELASI`: tanpa jembatan itu, baris baru tidak dapat ditunjuk balik ke barisnya di sistem lama, karena `IndexObject` adalah posisi numerik yang hilang begitu baris berpindah ke kunci sendiri.

Yang dibutuhkan untuk perbandingan, dan sudah ada di model: nilai uang berpasangan (nilai asli, mata uang, nilai IDR, kurs) sehingga selisih dapat dipecah menurut sebabnya; kolom keadaan di tingkat baris sehingga baris yang belum lengkap dapat dikeluarkan dari perbandingan alih-alih dihitung sebagai selisih; dan `MIGRASI_NILAI_DITOLAK` sehingga nilai yang tidak lolos penguraian terlihat sebagai daftar, bukan sebagai selisih tak terjelaskan.

**Batasnya**: baseline pembandingnya sendiri belum tentu benar. `TotalUR` selalu nol pada cabang mata uang sama karena dikalikan nol (ADR-0010, tambahan 18 September), dan dua rumus premi pemulihan bekerja atas nilai yang berbeda (FINDING-007). Apakah perilaku itu dipertahankan atau diperbaiki menentukan apa yang dianggap "sama" — dan itu pertanyaan akuntansi (`ASK-AKUNTANSI.md` no. 2b), bukan pertanyaan teknis.

### 9. Field pendaratan dari modul Komite

Folder `Komite Claim Non Prop` tidak pernah dibuka. **Tidak ada tabel Komite di spec ini.** Yang ada di sisi Claim hanya field pendaratan yang EVIDENCED di `BLUEPRINT.md` §8.3 — dan tiap satunya ditandai bahwa **pemilik dan penulisnya belum dikonfirmasi**:

| Field pendaratan | Di tabel | Catatan |
|---|---|---|
| status akseptasi | `ADJUSTMENT` | Sistem lama menulisnya **satu kali saja**, `= 0`, di `CreateChildKomiteCNP_Act` langkah 31. Transisi ke 1 atau 2 **tidak ada di folder ini** — hipotesis I1/I2 belum ditutup |
| nomor akseptasi, tanggal akseptasi | `ADJUSTMENT` | penulis belum dikonfirmasi |
| status tampilan klaim | `KLAIM` | enum di `BLUEPRINT.md` §3.1 |
| penanda dan catatan Subjectivity | `KLAIM` dan `ADJUSTMENT` | `CONTEXT.md`: persetujuan bersyarat |

`ComiteeClaim[]` — daftar approver — **tidak dimodelkan**, karena ia keanggotaan komite. Lihat Out of Scope.

---

### 10. Putusan atas 18 PageList di bawah `.ClaimData`

Sumber daftar dan cacah rujukannya: `BLUEPRINT.md` §2.1 — EVIDENCED. Tiap baris diberi **tepat satu** putusan.

| PageList | Ref | Putusan | Dasar |
|---|---:|---|---|
| `SpreadingRisk` | 164 | **dua entitas**: `ALOKASI_LAYER` + `RETENSI_CEDANT` | DECIDED(ADR-0010) |
| `ReceiverClaim` | 101 | entitas `REKENING_PENERIMA` | lihat bagian 12 |
| `AdjustmentList` | 48 | entitas `ADJUSTMENT` | `CONTEXT.md`, `MEMORI_PEMAHAMAN.MD` §4.3 |
| `ListClaimAmount` | 37 | entitas `NILAI_KLAIM_MATA_UANG` — **masukan**, diisi pengguna | `MEMORI_PEMAHAMAN.MD` §4.2 |
| `CNPSpreadLoss` | 23 | entitas `PEMBAGIAN_KERUGIAN` — **masukan mesin alokasi, bukan keluarannya** | `MEMORI_PEMAHAMAN.MD` §6.2 Langkah 2; DECIDED(ADR-0011, ADR-0013) |
| `InterestList` | 19 | entitas `OBJEK_PERTANGGUNGAN` | `MEMORI_PEMAHAMAN.MD` §4.2 |
| `SpreadingClaim` | 19 | entitas `PENYEBARAN` | `CONTEXT.md` **Penyebaran** |
| `ListTotalEstimation` | 17 | **view turunan** — rekap total per mata uang | agregasi murni atas `NILAI_KLAIM_MATA_UANG` |
| `EstimationList` | 14 | entitas `ESTIMASI_AWAL` — **masukan** | membawa `TypeLoss`, `KursValue`, `TotalEstimasi` |
| `ReinstatementList` | 13 | entitas `PREMI_PEMULIHAN` | `MEMORI_PEMAHAMAN.MD` §6.3 |
| `SpreadingAdjustment` | 13 | **view turunan** — agregasi **per jenis reasuransi per mata uang** | kunci agregasinya `(Currency, TreatyName)`, `MEMORI_PEMAHAMAN.MD` §6.4 — satuan kasar, bagian 0 |
| `SpreadingAdjustmentQS` | 12 | **view turunan**, idem untuk quota share | idem |
| `SpreadingBreakQS` | 11 | entitas `PENYEBARAN`, baris ber-jenis quota share | `CountLossAllocation_act` Langkah 10 |
| `TotalInterestInsured` | 7 | **view turunan** — rekap nilai pertanggungan | agregasi atas `OBJEK_PERTANGGUNGAN` |
| `Attachment` | 7 | entitas `DOKUMEN_KLAIM` | DECIDED(ADR-0027); lihat bagian 12 |
| `ObjectList` | 7 | **lubang tercatat** — dua PageList untuk objek yang sama | lihat di bawah |
| `ListClaimAcceptation` | 4 | **view turunan** — rekap akseptasi | agregasi atas `AKSEPTASI` |
| `ClaimComitee` | 2 | **tidak dimiliki** — keanggotaan komite | tertahan, lihat *Out of Scope* |

Dua PageList di luar daftar 18 yang tetap harus tertangani:

| PageList | Letak | Putusan | Dasar |
|---|---|---|---|
| `SuggestList` | di bawah `.ClaimData` (`MEMORI_PEMAHAMAN.MD` §4.2) | entitas `KRONOLOGI_KLAIM` | `CONTEXT.md` **Kronologi**; DECIDED(ADR-0009) — isinya tidak pernah dibaca mesin |
| `AlokasiXOLPaid` | di bawah `.AdjustmentList[]` (`MEMORI_PEMAHAMAN.MD` §4.3) | **lubang tercatat** | lihat di bawah |

#### Mengapa `SpreadingAdjustment` dan `SpreadingAdjustmentQS` bukan entitas

Keduanya beragregasi dengan kunci `(Currency, TreatyName)` — dan `TreatyName` adalah **satuan kasar**, yaitu jenis reasuransi (bagian 0). Sebelum putusan (B), keduanya terbaca sebagai agregasi "per layer" yang tidak jelas butirnya. Dengan penamaan yang benar, keduanya terbaca terang: **agregasi per jenis reasuransi per mata uang**, seluruhnya dapat dihitung ulang dari `ALOKASI_LAYER` dan `PENYEBARAN`.

Menyimpannya sebagai tabel berarti menyimpan hasil yang dapat menyimpang dari sumbernya — persis yang dilarang ADR-0013, dan persis yang membuat `CountSpreadingXOL` harus **menghapus lalu membangun ulang** isinya setiap kali dijalankan.

#### Dua lubang, dicatat bukan dihilangkan

**`ObjectList` berdampingan dengan `InterestList`.** Keduanya membawa objek pertanggungan dengan properti yang bertindih (`ObjectName`, `TSIPerObject`, `Currency`). Apa yang membedakan keduanya **TIDAK DITEMUKAN DI XML**. Sampai itu jelas, hanya `InterestList` yang dimodelkan; `ObjectList` masuk daftar lubang. Menyatukannya sekarang berarti menebak bahwa keduanya benda yang sama.

**`AlokasiXOLPaid` dibaca, penulisnya tidak diketahui.** `CreateChildKomiteCNP_Act` membacanya untuk **mengurangi nilai layer dengan yang sudah dibayar** (`MEMORI_PEMAHAMAN.MD` §6.5) — jadi ia memengaruhi nilai yang masuk jenjang persetujuan. Penulisnya **TIDAK DITEMUKAN DI FOLDER INI** (`_selesai/OPEN-QUESTIONS.md` C5).

Itu alasan ia masuk daftar lubang, **bukan alasan ia menghilang dari spec**. Bila penulisnya ada di modul Komite, ia kopling lintas modul yang belum masuk Boundary Contract; bila tidak ada penulisnya, pengurangan itu selalu nol dan nilai yang masuk komite selalu nilai penuh. Keduanya mengubah model, dan keduanya menunggu C5.

---

### 11. Kolom berpasangan: nilai hitungan dan nilai suntingan

DECIDED(ADR-0008). Ini keputusan **model data**, bukan aplikasi — ADR-nya menyatakannya sendiri: *menambahkan ini setelah tabel terbentuk berarti membongkar setiap baris yang sudah ada.*

Bentuknya, untuk setiap besaran yang dapat disunting:

| Kolom | Isi | Arti NULL |
|---|---|---|
| `<besaran>_HITUNG` | hasil mesin alokasi, **selalu ditulis ulang** saat hitung ulang | belum pernah dihitung |
| `<besaran>_SUNTING` | nilai yang dimasukkan orang, **tidak pernah ditimpa** mesin | tidak pernah disunting — dan itu keadaan normal |
| `<besaran>_DIPAKAI` | kolom turunan: `COALESCE(<besaran>_SUNTING, <besaran>_HITUNG)` | kedua-duanya kosong |
| `DISUNTING_OLEH` | pelaku — bentuknya di bagian 14 | belum pernah disunting |
| `DISUNTING_PADA` | waktu | idem |
| `ALASAN_SUNTING` | teks bebas, **tidak pernah dibaca mesin** | idem — DECIDED(ADR-0009) |

Dengan bentuk ini, hitung ulang yang idempoten (ADR-0011) menulis hanya kolom `_HITUNG`, dan suntingan bertahan tanpa perlu bendera terpisah. Itu juga yang membuat FINDING-004 tidak dapat terulang: di sistem lama satu-satunya penjaga adalah `.IsEditClaim`, yang **ditulis tetapi tidak pernah dibaca** — EVIDENCED, `BLUEPRINT.md` §2.5.

**Berlaku di `ALOKASI_LAYER` dan juga di `RETENSI_CEDANT`.** Alasannya bukan simetri: `EditXOLAlokasi` menulis `.IsEditClaim` pada class `ASM-FW-GISFW-Data-SpreadingRisk` — dan di sistem lama class itu memuat **baris layer maupun baris Retensi Cedant** dalam satu PageList. Jadi jalur suntingan yang sama menyentuh keduanya. Setelah ADR-0010 memisahkannya menjadi dua tabel, kemampuan menyunting itu ikut ke keduanya.

**Batas klaim**: besaran mana saja yang *sesungguhnya* dapat disunting lewat `EditXOLAlokasi` belum disapu satu per satu; yang EVIDENCED adalah bahwa rule itu ada dan menulis penanda pada class tersebut.

---

### 12. Kepemilikan `DOKUMEN_KLAIM` dan `REKENING_PENERIMA`

**`DOKUMEN_KLAIM` — tabel penghubung, bukan salinan.** DECIDED(ADR-0027, ADR-0026). Berkasnya berada di Google Cloud Storage; yang tersimpan hanya metadata, URL, dan `EXPDATE` yang disegarkan (`MEMORI_PEMAHAMAN.MD` §5.7). Yang **dimiliki** modul ini adalah *fakta bahwa dokumen X melekat pada klaim Y* — itu pernyataan tentang klaim, dan klaim miliknya. Yang **tidak dimiliki** adalah dokumennya.

Karena URL berumur terbatas, ia **tidak boleh disimpan sebagai nilai tetap**: kolom URL berpasangan dengan kolom kedaluwarsa, dan URL yang lewat kedaluwarsa diperlakukan sebagai tidak ada, bukan sebagai nilai yang masih sah.

**`REKENING_PENERIMA` — salinan saat pembayaran, dan itu pengecualian sah terhadap ADR-0023.** Alasannya bukan kenyamanan: rekening penerima adalah **bagian dari instruksi bayar yang sudah diterbitkan**. Bila nomor rekening seseorang berubah tahun depan, instruksi bayar tahun lalu tidak boleh ikut berubah — ia catatan tentang ke mana uang benar-benar dikirim.

ADR-0023 melarang salinan karena salinan dapat menyimpang dari bentuk kanoniknya. Di sini menyimpang justru yang benar: yang disalin bukan *keadaan rekening sekarang*, melainkan *keadaan rekening pada saat pembayaran*. Keduanya benda berbeda, dan yang kedua tidak punya sumber lain.

#### Parameter perhitungan yang di-snapshot

Prinsipnya satu: **setiap nilai yang masuk rumus disimpan bersama hasilnya**, supaya hasil dapat dihitung ulang dan diperiksa tanpa membaca keadaan master hari ini. DECIDED(ADR-0011, ADR-0013).

| Parameter | Sumbernya di sistem lama | Kenapa harus di-snapshot |
|---|---|---|
| Limit Layer | `TreatyInMaster.Limits[].Limit` | penyebut rumus premi pemulihan |
| Premi Deposit | `Limits[].MDPList[]` | pengali rumus premi pemulihan |
| Persentase premi pemulihan | `Limits[].ReinstatementPct` | pengali |
| **Porsi cedant** (`ShareCeding`) | `.ClaimData.ShareCeding` | mengalikan Salvage, Biaya Penilaian, dan biaya lain di Mode 1 (`MEMORI_PEMAHAMAN.MD` §6.2) |
| **Prorata Klaim** (`PctProrateClaim`) | `.ClaimData` / `ListClaimAmount` | mengalikan batas layer dan nilai Retensi Cedant |
| **Nilai Retensi Cedant** | `Limits[].Deductible` | menentukan titik mulai layer pertama |
| **Kurs**, tanggal kurs, sumber kurs | `GETCURRENCYSTANDARD` | DECIDED(ADR-0007). Diperkuat FINDING-006: fungsi lama mengembalikan `1` bila tidak ditemukan, dan `1` tidak dapat dibedakan dari kurs sah |
| **Porsi reasuradur** (`RNMShare`) | `TreatyInMaster.RNMShare` | dipakai di hampir seluruh rumus; dan dinormalisasi jadi `100` pada cabang tertentu (`MEMORI_PEMAHAMAN.MD` §6.2 Langkah 0) — nilai yang dipakai harus terbaca, bukan disimpulkan |

---

### 13. Tabel koreksi bernilai tercatat — pengganti 29 langkah tambalan

DECIDED(ADR-0004). **`NILAI_DIPATOK_PER_KLAIM` yang disebut di bagian 1 adalah benda ini**, dan namanya diganti menjadi **`KOREKSI_NILAI`** supaya isinya tidak salah dibaca sebagai parameter biasa. Kelima kolom yang diwajibkan ADR-0004 ada seluruhnya:

| Kolom | Isi | Arti NULL |
|---|---|---|
| sasaran | tabel dan baris yang dikoreksi | — (`NOT NULL`) |
| besaran | kolom mana yang dikoreksi | — (`NOT NULL`) |
| `NILAI_SEBELUM` | nilai hasil hitung sebelum koreksi | tidak ada nilai sebelumnya — koreksi mengisi yang semula kosong |
| `NILAI_SESUDAH` | nilai yang berlaku | koreksi **membatalkan** nilai, bukan menggantinya |
| `ALASAN` | teks bebas, tidak pernah dibaca mesin | — (`NOT NULL`) |
| `PELAKU` | bentuknya di bagian 14 | — (`NOT NULL`) |
| `WAKTU` | — | — (`NOT NULL`) |

Bedanya dengan kolom suntingan di bagian 11: kolom `_SUNTING` adalah **cara kerja normal** petugas atas satu baris alokasi; `KOREKSI_NILAI` adalah **penyimpangan bernama** atas nilai apa pun, termasuk yang tidak punya jalur suntingan. Di sistem lama, ketiga hal yang dipatok — Premi Pemulihan yang dipatok per klaim, pemetaan Biaya Penilaian per porsi, dan batas Retensi Cedant alternatif — tidak punya jalur suntingan sama sekali; satu-satunya jalan adalah mengubah kode. Itu sebabnya tambalannya lahir.

---

### 14. Bentuk kolom pelaku selama tabel pengguna belum ada

Aturannya sudah ada di sumber, dan berasal dari kejadian nyata — FINDING-002 §9, penggantian `DARTO` menjadi `CHRISTINEANGELINA` di dalam kode. DECIDED(ADR-0006, FINDING-002).

> **Identitas pelaku tidak pernah ditimpa.** Perwakilan dicatat sebagai *"X bertindak untuk Y"*, dengan **kedua identitas tersimpan**.

Bentuk kolomnya, tanpa menunggu tabel pengguna:

| Kolom | Isi |
|---|---|
| `PELAKU` | pengenal orang yang benar-benar melakukan tindakan, disimpan sebagai teks, **tidak pernah diubah setelah ditulis** |
| `ATAS_NAMA` | pengenal orang yang diwakili. `NULL` berarti tidak ada perwakilan — bukan berarti tidak diketahui |

Keduanya kolom teks tanpa foreign key. Ini **DECIDED-TEKNIS dan bersifat sementara**: ketika tabel pengguna ada, keduanya menjadi foreign key tanpa mengubah arti maupun jumlah kolomnya. Yang tidak boleh terjadi adalah sebaliknya — menunggu tabel pengguna, lalu menemukan bahwa tindakan lama hanya menyimpan satu identitas dan perwakilan sudah hilang.

**Batas klaim**: tabel pengguna, jabatan, dan keanggotaan komite **tidak dirancang di sini**, dan sebabnya ada di *Out of Scope* — ADR tentang pemisahan itu tidak pernah ditulis.

#### Dua lubang yang belum terdaftar, sekarang terdaftar

| Lubang | Isi | Kenapa ia menyentuh model data | Register |
|---|---|---|---|
| Penulis `AlokasiXOLPaid` | dibaca `CreateChildKomiteCNP_Act` untuk mengurangi nilai layer dengan yang sudah dibayar | menentukan nilai yang masuk jenjang persetujuan | `_selesai/OPEN-QUESTIONS.md` **C5** |
| `TreatyName == "Previously Calculated UR"` | dibaca `GenerateCACNP_Act` baris **7223** dan **7425**; **tidak pernah ditulis** di folder Claim | bila ia **keadaan ketiga** Retensi Cedant — di samping "berlaku" dan "tidak ada" — maka kunci `(klaim, mata uang)` **belum cukup**, karena satu klaim dapat memuat Retensi Cedant berjalan dan Retensi Cedant terhitung-sebelumnya pada mata uang yang sama | `_selesai/OPEN-QUESTIONS.md` **C9** |

Keduanya **TIDAK DITEMUKAN** penulisnya di 279 berkas; disimpulkan dari ketiadaan penetapan, bukan dari adanya rule.

---

### 15. Pembulatan dua desimal — di tepi, bukan di kolom

DECIDED-TEKNIS (AK-1). Tidak ada kelompok presisi untuk "nilai uang final". Seluruh nilai uang disimpan pada kelompok **U1**; pembulatan dua desimal terjadi **hanya di tepi**.

Alasannya satu ekspresi, EVIDENCED — `Activity\HitServiceToKasir_Act.xml`:

```
TempKasir.CARI9 = .TotalClaim - .PremiumSpreaded
TempKasir.CARI9 = @toDecimal(TempKasir.CARI9) + .AdjustmentValue
TempKasir.CARI9 = @toDecimal(TempKasir.CARI9) + .AdjusterFeeValue
TempKasir.CARI9 = @toDecimal(TempKasir.CARI9) + .SalvageValue
```

`.TotalClaim` adalah **nilai antara** — hasil alokasi, dipakai rumus premi pemulihan — **sekaligus bahan instruksi bayar**. Satu kolom tidak dapat berskala 20 dan 2 sekaligus. Menyimpannya pada skala 2 berarti membulatkan sebelum rumus premi pemulihan membacanya.

Dua tempat pembulatan, keduanya di luar tabel bisnis:

| Tempat | Bentuk |
|---|---|
| Lapisan view kompatibilitas | pembulatan dilakukan di definisi view, sehingga hilir menerima bentuk yang sudah dikenalnya |
| Arsip muatan keluar | bila muatan yang dikirim perlu disimpan, arsipnya **tabel tersendiri** berisi muatan apa adanya — bukan kolom di tabel bisnis |

Aturan yang mengikat keduanya: **tidak ada kolom di tabel bisnis yang menyimpan nilai yang sudah dibulatkan.**

---

### 16. Penamaan — final, 19 September 2026

DECIDED-TEKNIS. Bahasa Indonesia, `UPPER_SNAKE_CASE`, mengikuti glosarium. Nama skema **`KLAIMNP`**, akun aplikasi **`KLAIMNP_APP`**.

#### 16.0 Batas 30 byte — aturan tetap, bukan pengamanan sementara

**Seluruh pengenal dijaga di bawah 30 byte, dan itu berlaku selamanya.**

Versi pertama bagian ini menjaga 30 byte sebagai **pengamanan** sampai nilai `COMPATIBLE` diketahui lewat REQ-032, dengan rencana melonggar ke 128 bila ternyata boleh. Rencana itu dibatalkan.

Alasannya bukan kehati-hatian: **30 byte sah di setiap versi Oracle; 128 hanya di sebagian.** Memilih yang berlaku di mana-mana membuat skema ini benar tanpa bergantung pada nilai yang tidak kita pegang — dan **pertanyaan yang jawabannya tidak mengubah apa pun bukan pertanyaan yang perlu ditunggu**. REQ-032 karena itu turun dari BLOCKER menjadi **VERIFIKASI**.

#### 16.1 Tabel singkatan tertutup — **DIBATALKAN, bukan ditunda**

Versi pertama memuat daftar 23 singkatan (`KLAIM`→`KLM`, `MATA_UANG`→`MTU`, dan seterusnya), tertutup dan sekali-tulis. **Daftar itu tidak dibuat, dan tidak perlu ada.**

Premisnya gugur. Tekanan panjang selama ini datang dari **satu kebiasaan, bukan dari bahasa**: mencantumkan daftar kolom di nama constraint. Contoh yang dulu dipakai sebagai pembenaran, `UQ_AKSEPTASI_KLAIM_LAYER_MATA_UANG` — 34 byte — panjang karena ia **mengeja isinya**, padahal isi itu sudah ada di katalog Oracle dan dapat dibaca dari sana kapan saja.

Begitu nama constraint berhenti mengeja kolom, tekanannya hilang. **Nama tabel dan kolom sendiri tidak pernah menabrak batas.**

Biaya membatalkannya nol; biaya membekukannya permanen — **singkatan tidak pernah dicabut sesudah dipakai.** Itu persis cara `UR` dan `MDP` lahir di sistem lama, dan persis yang glosarium tutup.

Tiket `02` dibatalkan dan pindah ke `_mati/`.

#### 16.2 Empat aturan

**1. Nama tabel dan kolom memakai kata utuh bahasa Indonesia**, `UPPER_SNAKE_CASE`, mengikuti glosarium. **Tidak ada singkatan buatan.**

Bila sebuah nama melewati 30 byte, yang diganti adalah **katanya** — dengan kata utuh yang lebih pendek — **bukan dipotong jadi singkatan.**

**2. Nama constraint dan index tidak mengeja kolom.** Bentuknya `<peran>_<tabel>`, ditambah **pembeda numerik** bila sebuah tabel punya lebih dari satu constraint sejenis. Isinya dibaca dari katalog, bukan dari nama.

```
PK_KLAIM                     CK_KLAIM_1   CK_KLAIM_2   CK_KLAIM_3
UQ_KLAIM_1                   IX_KLAIM_1
FK_ADJUSTMENT_1              FK_ADJUSTMENT_2
```

Kunci primer tidak berpembeda: sebuah tabel hanya punya satu.

**3. Peran dipakai sebagai awalan tetap**: `PK_` · `UQ_` · `FK_` · `CK_` · `IX_`, ditambah `V_` untuk view dan `SQ_` untuk sequence.

Kelimanya **bukan singkatan bahasa**; ia penanda jenis objek yang sudah baku di Oracle dan tidak menutupi arti apa pun. Aturan "tanpa singkatan" berlaku pada istilah bisnis, bukan pada penanda jenis.

**4. Tiap nama diperiksa panjangnya saat ditulis.** Nama yang melewati 30 byte **menghentikan pembangkitan**, bukan lolos diam-diam.

Daftar `_Avoid_` tetap berlaku penuh, di nama tabel maupun kolom.

#### 16.2a Pembangkit pengenal — sequence, bukan IDENTITY

Kedua puluh dua tabel berkunci primer `NUMBER(19)`, dan sampai 19 September 2026 **tidak satu pun menyatakan dari mana nilainya datang**: nol `IDENTITY`, nol `SEQUENCE`, nol `DEFAULT` di seluruh berkas. Hak `CREATE SEQUENCE` sudah diberikan di `00_SKEMA_DAN_AKUN.sql` — niatnya ada, objeknya tidak pernah dibuat.

**Dipilih sequence, satu per tabel, dinamai `SQ_<tabel>`.** Alasannya **persis alasan yang menetapkan batas 30 byte**: pilih bentuk yang sah di setiap versi. `IDENTITY` baru ada sejak 12.1; sequence sah jauh sebelumnya. REQ-032 diturunkan menjadi verifikasi dengan janji bahwa jawabannya **tidak mengubah apa pun** — memakai `IDENTITY` akan membatalkan janji itu.

Nilai diminta pemanggil lewat `NEXTVAL`, dan itu sejalan ADR-0017: satu pintu tulis, dan pintu itu yang meminta nomornya. `GRANT SELECT` atas tiap sequence diberikan ke `KLAIMNP_APP` saja.

#### 16.3 Satu nama tabel diganti — konsekuensi aturan 1 dan 2

`DAFTAR_KLAIM_PENJAGA_TANGGAL` **28 byte**, di bawah batas sebagai nama tabel. Tetapi constraint-nya `FK_DAFTAR_KLAIM_PENJAGA_TANGGAL_1` = **33 byte**, di atas batas.

Aturan 1 menentukan penyelesaiannya: **ganti katanya, jangan potong.** Kata `DAFTAR_` dibuang karena sebuah tabel memang sudah daftar — ia tidak menambah arti.

| | Lama | Baru | Byte |
|---|---|---|---:|
| Tabel | `DAFTAR_KLAIM_PENJAGA_TANGGAL` | **`KLAIM_PENJAGA_TANGGAL`** | 21 |
| Constraint terpanjangnya | `FK_DAFTAR_KLAIM_PENJAGA_TANGGAL_1` (33) | **`FK_KLAIM_PENJAGA_TANGGAL_1`** | 26 |

Berkasnya ikut: `ddl-usulan/18_KLAIM_PENJAGA_TANGGAL.sql`.

**Ini satu-satunya nama tabel yang berubah.** Kedua puluh satu lainnya sudah muat tanpa disentuh — diperiksa, bukan diperkirakan.

#### 16.4 Perhitungan panjang byte — dihitung atas berkas DDL yang ada

*Tabel ini disalin dari keluaran `alat/periksa-penamaan.py`, bukan ditulis ulang dengan tangan — dua angka sebelumnya keliru justru karena ditulis dari ingatan.*

| Jenis | Cacah | Pengenal terpanjang | Byte |
|---|---:|---|---:|
| tabel | 22 | `KLAIM_PENJAGA_TANGGAL` | 21 |
| **view** | 9 | **`V_TOTAL_NILAI_PERTANGGUNGAN`** | **27** |
| sequence | 22 | `SQ_MIGRASI_NILAI_DITOLAK` | 24 |
| constraint | 115 | `UQ_KLAIM_PENJAGA_TANGGAL_1` | 26 |
| index | 11 | `IX_MIGRASI_NILAI_DITOLAK_1` | 26 |
| akun, peran, tablespace | 8 | `KLAIMNP_PERAN_PEMASANGAN` | 24 |
| kebijakan pengawasan | 1 | `KLAIMNP_PENGAWASAN_TULIS` | 24 |
| kolom | 168 | `PREMI_PEMULIHAN_PORSI_IDR` | 25 |

**355 pengenal unik di seluruh `ddl-usulan/`. Di atas 30 byte: nol. Sisa singkatan lama: nol.**

Angka itu mencakup **seluruh** jenis pengenal, bukan hanya constraint dan index. Yang terpanjang ternyata **nama view**, bukan nama constraint — itu justru yang membuat pemeriksaan ini perlu mencacah semuanya, dan sebabnya pemeriksaan pertama (yang hanya melihat tabel dan constraint) melaporkan 26 byte dan keliru.

Dihitung oleh **`alat/periksa-penamaan.py`**, yang **hanya membaca** dan tidak pernah mengubah apa pun. Ia memeriksa tiga hal dan **keluar dengan kode 1** bila salah satunya dilanggar: nama di atas 30 byte, sisa singkatan dari daftar yang dibatalkan, dan nama constraint yang dipakai dua kali. Itu aturan 4 sebagai kode, bukan sebagai niat. Peta nama lama → baru tersimpan di `alat/peta-nama-2026-09-19.tsv`, 105 baris, supaya penggantian 19 September dapat ditelusuri balik.

Bandingkan bentuk lama dan baru pada contoh yang dulu menjadi pembenaran daftar singkatan:

| Bentuk | Nama | Byte | Terbaca? |
|---|---|---:|---|
| mengeja kolom | `UQ_AKSEPTASI_KLAIM_LAYER_MATA_UANG` | 34 | ya, tetapi **lewat batas** |
| bersingkatan (dibatalkan) | `UQ_AKS_KLM_LYR_MTU` | 18 | **tidak** — perlu daftar singkatan untuk membacanya |
| **berlaku** | `UQ_AKSEPTASI_1` | **14** | ya — kolomnya dibaca dari katalog |

Bentuk yang berlaku **lebih pendek daripada bentuk bersingkatan sekaligus lebih terbaca**. Itu yang membuat daftar singkatan bukan hanya tidak perlu, melainkan lebih buruk.

---

### 17. Kamus kolom dan DDL

| Keluaran | Berkas |
|---|---|
| Kamus kolom — nama, tipe, kelompok presisi, nullability, arti NULL, sumber | `KAMUS-KOLOM.md` |
| DDL per objek, dapat dijalankan | `ddl-usulan/01_KLAIM.sql` … `ddl-usulan/21_MIGRASI_PENDARATAN.sql` |

**21 tabel, 377 kolom, 9 view.** Keduanya **dibangkitkan dari satu definisi yang sama**, sehingga kamus dan DDL tidak dapat berbeda — kalau tipe berubah di satu tempat, ia berubah di keduanya atau tidak berubah sama sekali.

Setiap berkas DDL berkepala dua baris penanda:

```
-- USULAN. Presisi provisional: ADR-0003 (accepted 18 Sep) kelompok AK-1.
-- BELUM PERNAH DIJALANKAN. Nama <= 30 byte: berlaku di 12.1 maupun 12.2+.
```

Kolom jejak audit hadir di setiap tabel: `DIBUAT_OLEH`, `DIBUAT_ATAS_NAMA`, `DIBUAT_PADA`, `DIUBAH_OLEH`, `DIUBAH_ATAS_NAMA`, `DIUBAH_PADA` — bentuk dua-identitas dari bagian 14.

**Yang sengaja tidak ada di DDL**: tidak satu pun kolom bertipe JSON di tabel kanonik. `MUATAN CLOB CHECK (MUATAN IS JSON)` hanya ada di `MIGRASI_PENDARATAN` — DECIDED(ADR-0028).

---

### 18. Sapuan ulang XML — S1 dan S2

Dijalankan 18 September 2026 atas `D:\XML_NURE\Claim Non Prop`. **Koreksi metode §8.0 dipakai**: sapuan sebelumnya hanya membaca pasangan `<PropertiesName>`/`<PropertiesValue>` yang berdampingan di Activity. Sapuan ini membaca empat lapisan.

Lapisan yang ada di ekspor, dan seluruhnya disapu: `Activity`, `DataTransform`, `Section`, `Harness`, `FlowAction`, `Flow`, `RDBList`, `ConnectREST`, `ReportDefinition`, `DecisionTable`, `When`, `SystemSettings`.

#### S1 — kolom yang sesungguhnya diterima Arasapas dan kasir — EVIDENCED

Ini kontrak ke luar; ia **dibaca, bukan dirancang**.

**Akseptasi — `Activity\SaveDataToOSAksep_Act.xml`, 17 parameter `InputParamOs.*`:**

| Parameter | Diisi dari |
|---|---|
| `NoClaim` | `pyWorkPage.ClaimData.NoClaim` |
| `TypeLoss` · `TypeLossID` | `.TreatyName` · `.TreatyType` — **satuan kasar**, menutup rantai bagian 0 |
| `Currency` · `CurrencyID` · `KursValue` | `.Currency` · `.CurrencyID` · `.Kurs` |
| `GrossValue` | `.ClaimEstimation + .AdjClaimValue` |
| `Value` | `.ClaimSpreaded` |
| `Adjusterfee` · `Salvage` · `CNPOthersFee` | senama |
| `PersenRNM` | `pyWorkPage.TreatyInMaster.RNMShare` |
| `IDMasterTreaty` | `pyWorkPage.TreatyInMaster.ID` |
| `CauseOfLoss` · `CauseOfLossID` | `pyWorkPage.ClaimData.*` |
| `EstimationDate` | `""` — **selalu kosong** |
| `Type` | `0` — **selalu nol** |

**Temuan yang mengubah bentuk view**: keempat nilai uang dikirim sebagai **selisih**, bukan nilai penuh —

```
InputParamOs.Value       = .ClaimSpreaded
InputParamOs.Value       = InputParamOs.Value - OutOSAcc.pxResults(1).Value
InputParamOs.GrossValue  = InputParamOs.GrossValue - OutOSAcc.pxResults(1).GrossValue
```

`OutOSAcc` adalah isi baris yang **sudah ada** di tabel akseptasi. Jadi sistem lama mengirim *tambahan terhadap yang sudah tercatat*, bukan keadaan sekarang. Itu pasangan wajar dari prosedurnya yang **selalu `INSERT`** (ADR-0024): baris bertambah, dan totalnya adalah jumlah seluruh baris.

**Batas klaim**: bahwa pola selisih itu **dimaksudkan** tidak dapat disimpulkan dari XML; yang EVIDENCED hanya bahwa ekspresinya begitu. Apakah jumlah seluruh baris akseptasi lama sama dengan nilai sekarang adalah pertanyaan data — bertemu dengan REQ-018.

**Kasir — `Activity\HitServiceToKasir_Act.xml`, 25 parameter `TempKasir.CARI1..CARI24` + `CARIDATETIME`.** Yang membawa uang hanya `CARI9`, `CARI10`, `CARI11`; sisanya identitas, tanggal, dan konstanta — termasuk `CARI13 = "NUSARE"`, `CARI14 = "D0031"`, dan `CARI15` yang bercabang `"100081"` / `"100115"` (G1).

**Yang tidak ditemukan**: `InsertOSKlaimCNP` **TIDAK DITEMUKAN** sebagai nama berkas di ekspor ini. `KonversiKlaim_Act.xml` ada tetapi **nol** pasangan properti — muatannya dirakit di tempat lain; `ConnectREST\KonversiKlaimNonLife.xml` ada dan belum dibaca isinya. Dicatat, tidak ditebak.

#### S2 — enum `CNPStatusCase`: sapuan empat lapisan **menguatkan** BLUEPRINT §3.1, bukan mengubahnya

| Yang dicari | Hasil |
|---|---|
| Nilai yang **ditulis** | **2**: `"INPUT ACCEPTATION CLAIM"`, `"COMITEE ACCEPTANCE (DEPT. HEAD)"` |
| Berkas yang menulis | 3, seluruhnya di lapisan `Activity` |
| Lapisan lain yang menulis | **nihil** — DataTransform, UI, dan Integrasi tidak menyentuhnya |
| Rule yang **menguji** nilainya | **nihil**, di seluruh empat lapisan |
| `"CLAIM ACCEPTED"` / `"CLAIM REJECTED"` sebagai literal | **nihil di 279 berkas** |

`BLUEPRINT.md` §3.1 sudah menyatakan ketiga hal ini. Yang berubah bukan kesimpulannya melainkan **dasarnya**: semula berlaku untuk lapisan Activity saja, sekarang untuk keempat lapisan. Tidak ada butir D9; tidak ada yang ditulis ulang.

> **Domain `STATUS_KLAIM` tidak dapat ditutup dari folder ini.** Dua nilai ditulis di sini, dua lagi diketahui dari `MEMORI_PEMAHAMAN.MD` §7.3 — dan memori mencakup 341 berkas, yaitu Claim **dan** Komite. Kelengkapannya hanya dapat dibuktikan dengan membuka folder Komite.

Karena itu `STATUS_KLAIM` **tidak diberi `CHECK`**, dan itu keputusan, bukan kelalaian: domain tertutup yang isinya belum lengkap akan menolak nilai yang sah pada hari pertama modul Komite tersambung.

#### Catatan pencatatan — ADR-0012 dan PENGETAHUAN §14

Sapuan ini pelaksanaan ADR-0012, dan hasilnya dibandingkan terhadap `BLUEPRINT.md` §7.5: **tidak ada tambalan per-case baru** yang muncul di luar yang sudah terdaftar.

Pernyataan nihil di atas berbunyi **"tidak ditemukan di lapisan yang sudah disapu"**, bukan "tidak ada di sistem". Pola yang dipakai tercatat supaya sapuan berikutnya tidak mengulanginya: pasangan `<PropertiesName>`/`<PropertiesValue>`, `<pyPropertiesName>`/`<pyPropertiesValue>`, perbandingan `==`/`!=` terhadap literal, dan `@contains(...)`.

**S3–S9 belum dijalankan** dan tidak menahan; keduanya masuk tiket.

---

### 19. Hasil review yang ditutup di keluaran ini

| # | Temuan review | Penutupannya |
|---|---|---|
| 1 | Tidak ada satu pun view | **9 view** di `ddl-usulan/V01…V09`, ditambah `V00_HAK_AKSES.sql`. Daftar kolomnya dari S1 |
| 2 | Pasangan IDR–kurs hanya di 1 tabel | **23 CHECK**, dipasang **oleh pembangkit** atas setiap kolom `*_IDR` di tabel yang punya `KURS` — bukan didaftar tangan |
| 3 | Kolom keadaan tanpa domain | `KEADAAN_BARIS` ber-`CHECK` di **11 dari 11** tabel; `KEADAAN_AKSEPTASI` ber-`CHECK`. `STATUS_KLAIM` dan hasil keputusan Komite **sengaja tanpa** `CHECK` — lihat S2 dan I1/I2 |
| 4 | Nama nyaris kembar | `ADJUSTMENT.STATUS_AKSEPTASI` → **`KEPUTUSAN_KOMITE`**. Dinamai menurut isinya: ia hasil keputusan Komite atas Adjustment, bukan keadaan entitas akseptasi |
| 5 | `PENGENAL_LAMA` adalah penyatuan | diganti **lima kolom terpisah**: `PZINSKEY_LAMA`, `PYID_LAMA`, `CASEID_LAMA`, `INDEX_OBJECT_LAMA`, `URUTAN_BARIS_LAMA`. Seluruhnya EVIDENCED; yang menunggu REQ-018 hanyalah **`UNIQUE` mana** yang berlaku, dan itu belum dipasang |
| 6 | `FK_ADJUSTMENT_2` tanpa index | `IX_ADJUSTMENT_1`. Pembangkit kini memasang index penopang untuk **setiap** FK yang kolomnya belum jadi awalan index lain |
| 6b | `FK_DPT_KLM_FK` bersufiks ganjil | benar, ada cabang penanganan tabrakan nama. Diseragamkan jadi `FK_KLAIM_PENJAGA_TANGGAL_1`; tidak ada sufiks `_FK` tersisa |

**Butir 5 tidak ditahan, dan alasannya berbeda dari yang diperkirakan review.** Kelima pengenal lama itu sendiri EVIDENCED — `PZINSKEY` dan `PYID` dari DDL tabel work, `CASEID` dari `OS_AKSEPTASI_KLAIM` dan §13.2, `IndexObject` dari §8.4. Yang menunggu REQ-018 bukan *kolom mana yang ada*, melainkan *kombinasi mana yang unik*. Kolomnya dapat ditulis sekarang; constraint-nya yang menunggu — dan itu sudah dinyatakan di kamus.

**Butir 7 review dibawa ke tiket, bukan ke spec**: `NUMBER(38,20)` adalah plafon Oracle, menyisakan tepat 18 digit bulat, dan tidak dapat dipetakan ke `float64` tanpa kehilangan. Pustaka desimal di sisi Golang beserta perilakunya saat hasil antara melewati 38 digit signifikan harus diputuskan di tiket pertama yang menyentuh nilai uang. Itu lahir langsung dari AK-1.

---

### 20. Atribusi suntingan — mana yang berlaku bila berbeda

`ALOKASI_LAYER` membawa `DISUNTING_OLEH`, `DISUNTING_PADA`, dan `ALASAN_SUNTING` **di tingkat baris**; `KOREKSI_NILAI` menyimpan perinciannya **per kolom** lewat `KOLOM_SASARAN`.

> **Bila keduanya berbeda, `KOREKSI_NILAI` yang berlaku.** Kolom di tingkat baris menunjukkan suntingan **terakhir** atas baris itu; ledger menyimpan seluruhnya. Yang di baris adalah kenyamanan baca, bukan sumber kebenaran.

---

### 21. Aturan yang lahir dari penutupan register — 19 September 2026

*Delapan aturan. Seluruhnya lahir dari butir `_selesai/OPEN-QUESTIONS.md` yang ditutup pada tanggal itu, dan dicatat di sini karena di sinilah ia dipakai — register menyimpan keadaan, spec menyimpan bentuknya.*

#### 21.1 Nama tidak menentukan tipe (menutup E13)

**Tidak ada properti berawalan `Is` atau `Flag` yang diberi domain biner tanpa bukti domainnya.** Setiap satunya dimodelkan dengan domain yang **benar-benar terbaca di sumber**, bukan dengan domain yang namanya janjikan.

| Properti | Domain terbaca | Bukan |
|---|---|---|
| `IsPLA` | `1` sebagai angka **dan `"1A"` sebagai teks** — `1A` bukan angka sama sekali | biner |
| `FlagProrate` | **empat** nilai, `0`–`3` | biner |
| `ReporterStatus` | **tiga** nilai | biner |

Konsekuensinya melampaui ketiganya: **104 kandidat penanda** tersapu di S11, dan **22 properti** terbukti dibandingkan sebagai angka **dan** sebagai teks, **10** di antaranya punya nilai kosong tanpa padanan numerik (D26). Aturan ini berlaku untuk seluruhnya, bukan hanya untuk tiga yang sudah ketahuan.

**Kosong bukan sinonim nol** — ADR-0019, dan di sini ia punya contoh yang mahal: `AcceptanceStatus` punya **empat** keadaan termasuk kosong, dan justru yang kosong yang lolos penjaga `CloseClaimMD` (D2).

#### 21.2 Arti nilai `PaymentType` disimpan sebagai **data**, bukan sebagai cabang di kode (menutup A6b)

`PaymentType` bukan label yang menunggu diterjemahkan — ia **field kendali aritmetika di dua jalur berbeda**, kasir dan Komite.

**Nilai disimpan apa adanya, termasuk kosong. Domainnya tidak ditutup pada tujuh nilai** — nilai kosong terbukti diuji, sehingga domain tertutup akan menolak yang sah.

Yang menjadi **baris tabel**, bukan cabang di kode: **nilai mana menambahkan besaran mana**.

| Nilai | Perilaku terbaca | Bukti |
|---|---|---|
| `1`, `2`, `5` | menambahkan `.AdjustmentValue` | `HitServiceToKasir_Act` 7754 |
| `4`, `6` | menambahkan `.AdjusterFeeValue` | `HitServiceToKasir_Act` 7896 |
| `3` | menambahkan `.SalvageValue`; jalur salvage berhenti di sana | `HitServiceToKasir_Act` 560/605 |
| `7` | mengendalikan apakah premi pemulihan diselisihkan ke jalur komite | `CreateChildKomiteCNP_Act` 8060, 8081, 8749 |
| kosong | diuji, tidak ditolak | D9 |

Ini ADR-0004 diterapkan: yang dulu cabang di kode menjadi baris di tabel. **Labelnya dapat diisi kapan saja tanpa mengubah struktur**, jadi ia tidak lagi menahan apa pun. REQ-027 memverifikasi nilai mana yang benar-benar dipakai; ia tidak menahan.

#### 21.3 Penanda lini usaha adalah **atribut data eksplisit** (menutup G1 dan A15)

Penentu syariah di sistem lama ada **dua**, dan **tidak satu pun ikut pindah**: node yang mengeksekusi (`When\IsPEGASyariah` yang mengubah `TempKasir.CARI15` menjadi `"100115"`) dan nama akun notifikasi (`pyNotifyAccountName` = `"NUSARE"` / `"NUSARESYARIAH"`).

Karena keduanya lenyap bersama arsitektur lama, **kedua cabang G1a dan G1b menuntut hal yang sama**: sistem baru memerlukan penanda yang berdiri sendiri. Itu yang membuat butir ini dapat diputus tanpa jawabannya.

**Penanda lini usaha menjadi atribut data**, sekelas dengan pembeda yang menggantikan `pyWorkIDPrefix` (D33 — bentuk nomor memutuskan lini bisnis 17 kali; sistem baru tidak menyimpulkan lini bisnis dari bentuk nomor). Ia disediakan **sekarang meski belum ada yang mengisinya**, sebagaimana ADR-0025 menyediakan kolom lingkup sebelum ada yang memerlukannya.

Dicatat sebagai perilaku lama, bukan ditiru: **jalur pembayaran buta terhadap syariah** — `StsSyariah` selalu `0` dan `CompanyName` selalu `"NUSARE"`, keduanya tertanam di kode (D42). Keduanya masuk daftar nilai tertanam yang menjadi data (ADR-0004). REQ-028 memverifikasi apakah klaim non-prop syariah pernah ada; ia tidak menahan.

#### 21.4 Pemisahan tugas adalah **kemampuan**, bukan kebijakan yang ditanam (menutup A3)

Model **wajib mampu** menegakkan pemisahan tugas: **pelaku tiap tindakan tercatat terpisah per tindakan**, sehingga pertanyaan *"apakah Registrasi dan Akseptasi dilakukan orang yang sama?"* selalu dapat dijawab dari data.

Apakah hal itu **dilarang atau diizinkan** adalah **setelan, bukan struktur**. Sediakan setelannya; jangan tanam jawabannya. Menanamnya berarti mengubah kebijakan menjadi pekerjaan skema.

Bentuk kolom pelakunya sudah ditetapkan di bagian 14, dan ADR-0016 menegaskan bentuk itu **potret**: pengenal dan nama sebagaimana tercatat saat tindakan terjadi, tanpa rujukan hidup ke sistem kepegawaian.

#### 21.5 Tabel perbedaan perilaku bertambah tiga baris

Perilaku lama yang **dicatat dan tidak ditiru**:

| Perilaku lama | Yang berlaku di sistem baru | Asal |
|---|---|---|
| **Sepuluh kolom ditulis dua kali** oleh mesin alokasi **dan** salah satu dari `AdjClaimCNP_Act`, `AdjClaimAmount_Act`, `SetActualPremium_ACT`: `ClaimEstimation`, `ClaimAmountAdjust`, `ClaimSpreaded`, `TotalClaim`, `AdjClaimValue`, `AdjusterFee`, `Salvage`, `CNPOthersFee`, `ClaimAmountIDR`, `TotalSpread` | Nilai turunan dihitung **dari data, oleh satu jalur**, dan hasilnya tidak bergantung urutan (ADR-0011, ADR-0013). Dua penulis atas kolom yang sama **tidak dapat terjadi pada bentuk ini** — bukan dicegah, melainkan mustahil | E15, D30 |
| **Selisih dikurangkan ke dua arah di satu berkas**: `SaveToOS` 2860 `GrossValue − TotalClaim` (berpenjaga `TreatyName=="UR"`), 3999 `TotalClaim − GrossValue` (tanpa penjaga) | **Selisih selalu nilai baru dikurangi nilai lama, satu arah, di seluruh sistem.** Mana yang "benar" di sistem lama tidak dapat dijawab dan **tidak perlu** dijawab — ADR-0013 menghapus ketergantungan urutan yang membuat pertanyaannya ada | E16, D22 |
| **Tiga skala dalam satu modul**: `.AdjusterFeeValue` 20, `.GrossValue` nol pernyataan, `Local.URLimit` 10; `TREATYINPRODUCTION` `NUMBER(20,4)` | Skala kanonik **20** sepanjang rantai, pembulatan hanya di tepi (ADR-0003) | E1, E2, E12, D24 |

#### 21.6 Satu jalur penutupan, dan klaim lama ditandai (menutup B9 dan A7)

**Jalur `CloseClaimMD` tidak dimigrasi.** Statusnya berubah dari `CANDIDATE-NOT-MIGRATED` menjadi **`NOT-MIGRATED`**.

Sistem baru punya **satu** jalur penutupan klaim, dan setiap penutupan mencatat **pelaku, waktu, dan dasar wewenangnya sebagai field terstruktur** (ADR-0009 — keputusan terstruktur, bukan teks bebas).

Klaim lama yang tertutup lewat jalur itu **dimigrasi sebagai tertutup dan ditandai**, sehingga asal-usulnya terbaca tanpa jalurnya ikut hidup. Menandai bukan menghakimi: ia menyatakan lewat jalur mana baris itu sampai ke keadaannya.

Pertanyaan *"apakah `CloseClaimMD` memang boleh menutup tanpa Komite"* (A7) karena itu **tidak perlu dijawab** — jalur pintasnya tidak dimigrasi. REQ-006 memverifikasi berapa banyak yang lewat sana; ia tidak menahan.

#### 21.7 Index ditentukan pasca-migrasi, dan itu keputusan (menutup B7)

**Index ditentukan dari pola akses nyata setelah skema berdiri dan data masuk**, bukan dari cacah baris sebelum ada satu pun kueri.

Kardinalitas PageList lama tidak menjawab pertanyaan index sistem baru, karena bentuk penyimpanannya berbeda: yang lama PageList di dalam blob, yang baru baris di tabel. Menurunkan index dari cacah baris lama berarti mengoptimalkan untuk kueri yang tidak akan pernah ditulis.

REQ-010 turun dari penahan menjadi **PENTING**.

#### 21.8 Muatan JSON mendarat utuh, yang dikenali naik, sisanya tercatat (menutup B4)

`adoptJSONObject` di `GeneratePlaCNP_Act`, `GetDetailPolis_act`, dan `SetValueClaimTNP_Act` mengadopsi teks `HASIL1` **tanpa memeriksa bentuknya**, dan yang terakhir menulis hasilnya ke halaman objek kerja (D43). Bentuk muatan itu karena itu **bukan hanya tidak terdokumentasi — ia tidak dibatasi.**

**Bentuk yang tidak dibatasi tidak dapat dijawab dengan daftar medan.** Ia hanya dapat dijawab dengan wadah:

1. **Mendarat utuh** — tiket `13` menetapkan satu tabel pendaratan tempat bentuk lama masuk apa adanya.
2. **Yang dikenali naik ke kanonik** — tujuh medan sudah terbaca dari notasi titik di `RDBList\GetDataOS` dan `GetDataCNPOS`: `Adjusterfee`, `GrossValue`, `CNPOthersFee`, `Salvage`, `Value`, `TypeLoss`, `Currency`.
3. **Sisanya tercatat sebagai tak terurai beserta asalnya** — tiket `14`, tidak dibulatkan dan tidak dibuang.

REQ-009 tetap berguna untuk memperbanyak daftar medan yang dikenali; ia **tidak menahan**, karena butir 1 dan 3 menerima apa pun yang datang.

---

## Testing Decisions

**Seam-nya satu: skema itu sendiri.** Tidak ada seam di lapisan aplikasi, karena lapisan itu di luar lingkup. Karena seam-nya SQL, **DDL harus dapat dijalankan** — DDL yang tidak jalan tidak dapat diuji, dan satu-satunya seam yang ada menjadi mati. Presisi provisional tidak menghalangi itu; ia dicatat di tabel domain, bukan dibiarkan kosong.

Uji yang baik di sini menguji **perilaku yang terlihat dari luar skema** — baris apa yang diterima dan baris apa yang ditolak — bukan bentuk internalnya. Menguji "kolom X bertipe `NUMBER(23,6)`" adalah menguji implementasi; menguji "baris kedua dengan kombinasi klaim, layer, dan mata uang yang sama ditolak" adalah menguji perilaku.

| Yang diuji | Bentuk |
|---|---|
| Kunci alami akseptasi | `INSERT` kedua dengan (klaim, layer, mata uang) sama **ditolak** |
| Pasangan uang dan mata uang | `INSERT` nilai uang tanpa mata uang **ditolak**; nilai IDR terisi tanpa kurs **ditolak** |
| Batas tanggal inklusif | klaim ber-Tanggal Kejadian tepat di hari terakhir treaty **diterima** |
| NULL bukan nol | baris menunggu kurs punya nilai IDR `NULL` dan keadaan barisnya bukan "lengkap"; agregasi atasnya tidak memperlakukannya sebagai nol |
| Penularan keadaan | baris turunan dari nilai yang menunggu kurs ikut bertanda menunggu kurs |
| Satu pintu tulis | akun selain akun aplikasi gagal `INSERT` ke tabel kanonik, dan berhasil `SELECT` lewat view |
| Bentuk view kompatibilitas | view Arasapas/kasir menghasilkan kolom dan tipe yang sama dengan yang diterima hari ini, termasuk `CASEID` |
| Paritas | untuk satu klaim contoh, tiap baris baru punya tepat satu baris pasangan di `MIGRASI_KORELASI` |

**Prior art di repositori ini**: tidak ada uji otomatis apa pun — tidak ada kerangka uji, tidak ada berkas uji, dan `MEMORI_PEMAHAMAN.MD` §10.8 mencatat ketiadaan itu sebagai salah satu "aspek yang tidak ada" di sistem lama. Yang paling dekat dengan prior art adalah `pengetahuan/PREFLIGHT.sql`: SQL berblok, tiap blok berdiri sendiri, read-only, dapat dijalankan dari TOAD, dengan cabang cadangan ketika hak akses tidak ada. Bentuk itu yang diikuti.

**Kasus uji utama untuk paritas**: `CLMNP-975` (ADR-0005, lewat ADR-0004) — tambalan terbaru, 2026-07-16.

---

## Out of Scope

**Di luar lingkup dan tidak boleh ditambahkan ke spec ini**: Golang, React, endpoint, kontrak API, layar, komponen, service, repository, ORM, nama fungsi, urutan pemanggilan. Bila sesuatu hanya dapat dijelaskan lewat kode aplikasi, ia bukan bagian spec ini.

**Modul Komite.** Foldernya tidak pernah dibuka. Tidak ada tabel Komite di sini, dan sembilan butir menunggu di `_selesai/OPEN-QUESTIONS.md` bagian C.

**Pengguna, jabatan, dan keanggotaan komite.** Tertahan, dan sebabnya bukan lingkup melainkan temuan: `docs/adr/0016-tanpa-database-link-ke-sistem-lain.md:12` menyebut *"ADR-0024 tentang pemisahan pengguna, jabatan, dan keanggotaan komite"* — padahal ADR-0024 adalah kunci alami akseptasi, dan **ADR tentang pemisahan itu tidak pernah ditulis**. Keputusan model data itu hidup di dalam kurung milik ADR lain. Sampai ia ditulis, tidak ada tabel pengguna, jabatan, atau keanggotaan komite yang dirancang di sini. Diperkuat ADR-0006: seluruh RBAC dirancang dari nol, tidak ada yang diwarisi.

**Fakultatif.** ADR-0020: tidak dimiliki, tidak dimodelkan, tidak ditulis.

**Tabel work Pega `PC_ASM_FW_GCNMFW_WORK`.** Tidak dimigrasi, dan strukturnya tidak ditiru.

**Objek `-Int-`**: `V_POLIS`, `T_STORAGE_IMAGE`, `EMAILKOMITE`, `M_LINK_SERVICE`. Dibaca, tidak dimiliki (ADR-0026). Perlu dipahami, tidak dipindahkan.

**74 objek di `PULL-LIST.csv`.** Daftar tarikan itu untuk **memahami**, bukan untuk memindahkan. "Seluruhnya dimigrasi" pada Q5 berarti seluruh baris **entitas yang dimiliki**, bukan seluruh objek di daftar tarikan.

**298 kolom DDL tanpa pasangan properti Pega.** Sebagian besar milik modul lain. Angka 110 berpasangan / 298 yatim / 549 yatim adalah **batas atas** hasil pencocokan nama, bukan angka pasti — pemetaan yang sah hanya datang dari `Data-Admin-DB-Table` (REQ-012).

---

## Yang tidak boleh dikarang — dan karena itu belum ditulis

| Yang dibutuhkan | Tabel/kolom mana yang tertahan | Register asalnya |
|---|---|---|
| Angka versi instance tujuan dan nilai `COMPATIBLE` | **Seluruh nama constraint dan index.** Batas pengenal 30 byte bila `COMPATIBLE` < 12.2, 128 byte bila ≥ 12.2. `UQ_AKSEPTASI_KLAIM_LAYER_MATA_UANG` sudah 34 byte | REQ-032 |
| Tipe dan presisi sebenarnya di Oracle | Seluruh kolom nilai — setiap tipe yang diturunkan dari pencocokan nama adalah **hipotesis** | REQ-001, ASK-AKUNTANSI no. 1 |
| Arti tiap nilai `PaymentType` | Kolom jenis pembayaran di `ADJUSTMENT` — sebarannya terukur, maknanya tidak | A6b, REQ-027 |
| Apakah `.ValueAdjustment` sudah IDR di hulu | Kolom nilai adjustment: apakah ia butuh pasangan mata uang atau sudah IDR. **Tidak boleh ada kolom yang mengandaikan salah satunya** | H1/H2 |
| Transisi status akseptasi 0 → 1/2 | Kolom status akseptasi di `ADJUSTMENT`: apakah ia pernah berubah nilai sama sekali | I1/I2, REQ-019 |
| Duplikat `PYID` dan baris ganda `(CASEID, TypeLoss, Currency)` | `UNIQUE (NOMOR_KLAIM)`, dan bentuk kolom pengenal lama di `MIGRASI_KORELASI` | REQ-018 (BLOCKER) |
| Pemetaan `Data-Admin-DB-Table` untuk class `ASM-FW-%` | Setiap pemetaan kolom lama → kolom baru yang kini berdasar kecocokan nama | REQ-012 (BLOCKER) |
| Apakah lini usaha syariah nyata | **Catatan, bukan kolom.** Bila G1a benar, penanda syariah menjadi atribut pada kontrak treaty — dan itu menyentuh model data | G1a/G1b, REQ-028 |
| Ambang kewenangan 15 atau 30; tarif 2,5% / 2% / 2,2% | Isi `TARIF_BERLAKU` dan `NILAI_DIPATOK_PER_KLAIM`. **Tempatnya dibuat, nilainya tidak diisi** | A8, ASK-AKUNTANSI no. 5 dan 6.1 |
| Presisi lama yang sebenarnya diterima akuntansi | Naiknya ADR-0003 dari `proposed` ke `accepted` | ASK-AKUNTANSI no. 1 |
| Volume baris per entitas yang dimiliki | Ukuran `MIGRASI_KORELASI` dan rencana migrasi ADR-0018 | REQ-031 |
| Apakah `POOLDATA` benar-benar menulis ke tabel Pega | Bentuk penegakan grant ADR-0017 | REQ-021 (BLOCKER) |
| Model otorisasi | Tidak ada tabel peran yang dirancang. `pyPrivilegeName` kosong di seluruh 279 berkas | A2, ADR-0006 |

---

## Further Notes

**Enam dari tujuh keluaran selesai.** Yang belum, satu:

| Keluaran | Keadaan |
|---|---|
| Daftar entitas beserta kunci alaminya | selesai — bagian 1, disapu ulang setelah putusan bagian 0 |
| Kamus kolom | selesai — `KAMUS-KOLOM.md`, **396 kolom** di 22 tabel ditambah **56 kolom** di 9 view, tiap satunya berkolom sumber. *Angka 368, lalu 452, keduanya **dicabut**.* Perintah yang menghasilkan 396: `awk '/^## /{t=($0~/^## [0-9]+\. `/)} t && /^\| `/{n++} END{print n+0}' KAMUS-KOLOM.md` — dan 396 yang sama keluar dari DDL tanpa melewati kamus, lewat `python sampah/alat/buat-tabel-datar.py`. *Angka **452 kolom di 22 tabel** dicabut 19 September 2026 oleh rekonsiliasi 1 tabel datar: 452 = **396 kolom tabel + 56 kolom view**, dan perintah yang tertulis di sebelahnya menghasilkan **457**, bukan 452. Ketiga angka beserta perintahnya di `sampah/keluaran/RINGKASAN-TABEL-DATAR.md` bagian 2.1.* |
| DDL per berkas per objek | selesai — `ddl-usulan/`, 21 berkas, dapat dijalankan |
| Daftar constraint yang menegakkan ADR | selesai — bagian 6, beserta lima aturan yang jatuh ke aplikasi |
| View untuk sistem hilir | selesai — bagian 7, beserta pemetaan halus ke kasar di bagian 0 |
| Kesiapan shadow-run | selesai — bagian 8 |
| **Pemetaan migrasi kolom demi kolom**, termasuk kolom TIDAK DIMIGRASI | **DITAHAN sampai REQ-012** — satu-satunya keluaran yang belum selesai. Pemetaan yang sah hanya datang dari `Data-Admin-DB-Table`; menuliskannya sekarang berarti mengabadikan kecocokan nama sebagai kebenaran |

**Nama benda sistem lama** yang dikutip di dokumen ini — `SpreadingRisk`, `AcceptanceStatus`, `CNPLayerList`, `IndexObject`, `V_POLIS`, `CLAIMXOL`, `json_klaim`, `OS_AKSEPTASI_KLAIM`, `PC_ASM_FW_GCNMFW_WORK`, `GETCURRENCYSTANDARD`, `SetPPNPPH` — seluruhnya **nama lama**, dikutip untuk menunjuk, tidak dipakai sebagai nama di model baru.

**Swa-periksa**

| Pertanyaan | Jawaban |
|---|---|
| Setiap kolom punya sumber? | Kamus kolom belum ditulis; setiap **entitas** dan setiap **kelompok presisi** di sini punya sumbernya |
| Istilah `_Avoid_` bersih dari nama tabel dan kolom? | Ya — nama entitas memakai istilah glosarium. `UR`, `MDP`, `XOL`, `reinstatement`, `spreading`, `CASEID`, `pyID` tidak muncul sebagai nama apa pun. `CASEID` muncul hanya sebagai **nama lama** dan sebagai kolom di lapisan view kompatibilitas, sesuai pengecualian ADR-0021 |
| Setiap ADR relevan ditegakkan, dan yang tidak bisa dinyatakan? | Ya — tabel constraint bagian 6, beserta lima aturan yang jatuh ke aplikasi |
| Setiap presisi ditandai sebagai usulan? | Ya — seluruh tabel domain bagian 3 |
| Setiap lubang tercatat, bukan terisi tebakan? | Ya — 13 baris di bagian *Yang tidak boleh dikarang* |
| Tidak ada tabel Komite yang terbawa? | Ya. Yang ada hanya empat field pendaratan EVIDENCED, masing-masing ditandai bahwa penulisnya belum dikonfirmasi |
