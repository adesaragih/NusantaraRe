# Hasil sapuan S3, S4, S5–S9

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas terbaca** pada empat lapisan (Activity, DataTransform, Section/Harness/FlowAction, RDBList/ConnectREST/ReportDefinition), ditambah `pengetahuan/DDL_Script_ClaimNonProp2.xls` versi 2026-09-18 10:57, 51 712 byte di cakram (**39 objek**; workbook-nya belum pernah dibuka, tetapi isinya sudah terekstrak ke `pengetahuan/ddl/` — lihat koreksi di bagian dua sumber).
> Folder `Komite Claim Non Prop` **tidak dibuka**. Satu sasaran sapuan berada di dalamnya dan dinyatakan tidak terjawab, bukan ditebak.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut.

Dijalankan 18 September 2026, atas tiket papan `07` `08` `09` `10`.

Label yang dipakai: **EVIDENCED** (terbaca di sumber), **EVIDENCED-NIHIL** (disapu, tidak ada, dengan lapisan dan pola disebut), **DERIVED** (kesimpulan atas bukti), **TIDAK DITEMUKAN** (tidak terjawab pada lapisan yang dibaca).

Tidak ada DDL, tipe kolom usulan, nama constraint, nama index, maupun singkatan di dokumen ini — gerbang ADR-0028 berlaku. Tipe kolom yang disebut adalah **tipe sistem lama yang terbaca**, yaitu bukti, bukan rancangan.

---

## S3 — perilaku per nilai `PaymentType`

**Temuan utama membatalkan penolakan saya sebelumnya.** Lihat bagian *Koreksi* di bawah.

Empat lapisan disapu dengan pola `PaymentType`. Enam berkas Activity, dua berkas Section. Nol di RDBList, ConnectREST, ReportDefinition, DataTransform.

### Nilai yang terbukti dipakai — EVIDENCED

| Nilai | Muncul sebagai | Berkas, baris |
|---|---|---|
| `''` kosong | penjaga langkah | `Activity\ProteksiSendKomiteCNP_Act.xml` 3448 |
| `1`, `2`, `5` | penjaga langkah, satu ekspresi | `Activity\HitServiceToKasir_Act.xml` 7666, 7754 |
| `2` | penjaga langkah | `Activity\SetInterimXOL_Act.xml` 391 |
| `'2'` | syarat tampil | `Section\AdjustmentDetailNP.xml` 1742 |
| `"3"` | penjaga langkah, tiga cabang prefiks | `Activity\HitServiceToKasir_Act.xml` 605, 617 |
| `3` | penjaga langkah | `Activity\HitServiceToKasir_Act.xml` 7950, 8038 |
| `4`, `6` | penjaga langkah, satu ekspresi | `Activity\HitServiceToKasir_Act.xml` 7808, 7896 |
| `7` | penjaga langkah | `Activity\ProteksiSendKomiteCNP_Act.xml` 1738, 1752 |
| `7` | pembanding negatif dalam hitungan | `Activity\CreateChildKomiteCNP_Act.xml` 8060, 8081 |
| `7` | penjaga langkah | `Activity\CreateChildKomiteCNP_Act.xml` 8749 |

**Ketujuh nilai ada, beserta nilai kosong.** Delapan kemungkinan, bukan tiga.

### Apa yang dilakukan tiap kelompok — EVIDENCED

`Activity\HitServiceToKasir_Act.xml` menyusun nilai instruksi bayar `TempKasir.CARI9` bertahap. Tiga langkah terakhir masing-masing **berpenjaga `PaymentType`**, dan tiap langkah menambahkan besaran yang **berbeda**:

| Penjaga | Baris deskripsi / penugasan / penjaga | Yang ditambahkan |
|---|---|---|
| `.PaymentType = 1 \|\| 2 \|\| 5` | 7666 / 7677 / 7754 | `.AdjustmentValue` |
| `.PaymentType = 4 \|\| 6` | 7808 / 7819 / 7896 | `.AdjusterFeeValue` |
| `.PaymentType == 3` | 7950 / 7961 / 8038 | `.SalvageValue` |

Urutan barisnya memastikan pasangannya: tiap langkah dibuka deskripsi, ditutup penjaga, dan penugasannya berada di antaranya.

**DERIVED**: nilai bayar **bukan** jumlah keempat besaran. Ia `(.TotalClaim − .PremiumSpreaded)` ditambah **tepat satu** besaran yang dipilih `PaymentType`. Ini mengubah pembacaan lama atas rangkaian `CARI9`, yang membacanya sebagai penjumlahan berantai tanpa syarat.

Penguat: `Activity\HitServiceToKasir_Act.xml` baris 560 memberi deskripsi langkah **"kalau salvage jangan lanjut"** pada penjaga baris 605 yang menguji `.PaymentType=="3"`. Jadi `3` berarti salvage, konsisten di dua tempat berbeda.

### Yang tetap tidak diketahui — TIDAK DITEMUKAN

**Arti tiap nilai** tidak terbaca. Definisi `Rule-Obj-Property PaymentType` (class `ASM-FW-GCNMFW-Data-Adjustment`, family `PAYMENTTYPE`) hanya muncul sebagai **rujukan indeks** di enam berkas; isinya ada di schema PegaRULES, yaitu **REQ-001**, yang belum kembali. Yang terbaca perilakunya, bukan namanya.

### Ketidakkonsistenan tipe — EVIDENCED

Properti yang sama dibandingkan sebagai angka di sebagian tempat dan sebagai teks di tempat lain: `.PaymentType = 1`, `.PaymentType==2`, `.PaymentType==3` di satu sisi; `.PaymentType=="3"`, `.PaymentType=='2'`, `Primary.PaymentType=="7"`, `.PaymentType==""` di sisi lain. Nilai kosong hanya bermakna pada perbandingan teks.

### Nama nyaris kembar — EVIDENCED

`Section\ViewDetailDeptHeadTreatyIn_UW.xml` memakai properti **`ClaimPaymentType`** (family `CLAIMPAYMENTTYPE`), berbeda rule dari `PaymentType`. Ini contoh keempat dari kelas cacat penamaan yang sudah didaftar `MEMORI_PEMAHAMAN.MD` §10.4 butir 4.

### Koreksi atas putusan saya sendiri

Saya sebelumnya **menolak** usulan §4.1 agar domain jenis pembayaran dibatasi tujuh nilai, dengan alasan `4`, `5`, dan `6` **TIDAK DITEMUKAN** dan `1` serta `2` hanya berasal dari `MEMORI_PEMAHAMAN.MD` §7.4 yang menyatakan dirinya *"disimpulkan dari pemakaian"*.

**Penolakan itu keliru, dan sebabnya cacat metode saya.** Saya menyapu literal **berkutip** (`"3"`, `"7"`) dan karenanya melewatkan perbandingan **tanpa kutip** — `.PaymentType = 1 || .PaymentType = 2 || .PaymentType = 5` dan `.PaymentType = 4 || .PaymentType = 6` — yang keduanya duduk di `HitServiceToKasir_Act.xml`, berkas yang sudah saya baca untuk keperluan lain.

Yang berubah: `1` sampai `7` kini **EVIDENCED**, bukan disimpulkan. Yang **tidak** berubah: nilai kosong `''` tetap terbukti diuji, sehingga domain tertutup tujuh nilai tetap akan menolak nilai yang sah. Bentuk domainnya ditetapkan saat tiketnya dikerjakan, dan sekarang ia beku di balik gerbang ADR-0028. Sebaran nilai yang sesungguhnya tetap diukur **REQ-027**.

---

## S4 — adakah padanan IDR untuk biaya penilaian, salvage, dan biaya lain

Pola: seluruh pengenal bersufiks `IDR` di empat lapisan, lalu dibandingkan terhadap daftar properti bernilai uang. Pola pembanding juga dijalankan untuk penamaan alternatif — `Rupiah`, `Idr`, `_IDR`, `Rp`, `InRp` — **nol hasil**.

### Yang punya padanan — EVIDENCED

| Properti uang | Padanan IDR | Kelas tempatnya |
|---|---|---|
| `GrossValue` | `GrossValueIDR` | `CountTotalInsterest_Act` |
| `ClaimAmount` | `ClaimAmountIDR` | `SpreadingRisk`, `SpreadingClaim`, `SpreadingBreakQS` |
| `Value` | `ValueIDR` | `CountTotalInsterest_Act` |
| `Amount` | `AmountIDR` | sembilan berkas, termasuk `SaveDataToOSAksep_Act` |
| `TotalSumInsured` | `TotalSumInsuredIDR` | `TreatyIn` |

### Yang **tidak** punya padanan — EVIDENCED-NIHIL

| Properti uang | Padanan IDR | Yang ada sebagai gantinya |
|---|---|---|
| `AdjusterFee`, `AdjusterFeeValue` | **tidak ada** | `AdjusterFeeRNM` |
| `Salvage`, `SalvageValue` | **tidak ada** | `SalvageRNM` |
| `CNPOthersFee` | **tidak ada** | `CNPOthersFeeRNM` |
| `OthersFee` | **tidak ada** | `OthersFeeRNM` |
| `AdjustmentValue` | **tidak ada** | — |
| `TotalClaim` | **tidak ada** | `TotalClaimKurs`, `TotalClaimRNM` |
| `PremiumSpreaded` | **tidak ada** | — |

Sufiks `RNM` bukan padanan IDR: ia porsi pihak, bukan mata uang. `TotalClaimKurs` adalah **kursnya**, bukan nilai terkonversinya.

### Kaitannya dengan ADR-0029 — dilaporkan, tidak diselesaikan

Inilah yang tiket `08` minta diperhatikan, dan jawabannya ya: **ada kelas nilai uang yang selama ini hidup tanpa pasangan IDR sama sekali.** **Tujuh baris, sembilan nama** — dua baris memuat dua nama masing-masing (`AdjusterFee`/`AdjusterFeeValue` dan `Salvage`/`SalvageValue`). Seluruhnya bernilai uang, dan lima di antaranya masuk hitungan instruksi bayar lewat `HitServiceToKasir_Act`.

Sementara itu `KursIDR` **ada** dan duduk di `SpreadingRisk` serta `CNPCurrencyList` — kursnya tersimpan, tetapi hanya sebagian nilai yang dikonversi memakainya. Ekspresi yang terbaca memakai `KursIDR` hanya untuk `.CNPReinstatement` dan `.GrossAdjustment`.

**Ini menyentuh ADR-0029 secara langsung dan tidak saya putuskan sendiri.** Dua pembacaan yang sama-sama berdiri di atas bukti yang sama:

1. Ketujuhnya memang **tidak pernah perlu** nilai IDR, karena tidak pernah masuk jurnal rupiah. Bila begini, kolom pasangannya tidak ditambahkan, dan ADR-0029 dinyatakan tidak berlaku untuk kelas ini.
2. Ketujuhnya **perlu** dan sistem lama tidak menyediakannya. Bila begini, kolomnya bertambah, dan pertanyaan "nilai rupiah dari biaya penilaian selama ini diambil dari mana" menjadi pertanyaan akuntansi, bukan pertanyaan rancangan.

Membedakan keduanya butuh jawaban orang, bukan sapuan lagi. Diajukan sebagai pertanyaan baru; lihat bagian penutup.

---

## S9 — apa persisnya yang ditulis `EditXOLAlokasi`

Empat lapisan. Empat berkas menyebut namanya: `Activity\`, `Harness\`, `Section\`, dan `Section\OutstandingClaim(1).xml`.

### Yang ditulisnya — EVIDENCED, dan daftarnya pendek

`Activity\EditXOLAlokasi.xml` melakukan **tiga** penugasan, seluruhnya:

| Sasaran | Nilai |
|---|---|
| `Local.ErrMsg` | `"Invalid Password"` |
| `.IsEditClaim` | `1` |
| `DataChronology.CARI1` | `EditCNP.CARI2` |

Ketiganya berpenjaga sama: `EditCNP.CARI1=="EDITCLAIMXOL"`.

`Section\EditXOLAlokasi.xml` (2 200 baris) memuat **dua** medan saja: `CARI1` dan `CARI2`. Tidak ada satu pun medan bernilai uang. `Harness\EditXOLAlokasi.xml` (3 572 baris) sama.

### Yang **tidak** disentuhnya — EVIDENCED-NIHIL

**Tidak satu pun properti bernilai uang.** Rule ini bukan penyunting alokasi; ia **gerbang kata sandi** yang menyalakan sebuah penanda dan menulis satu catatan kronologi. Namanya menjanjikan penyuntingan; isinya tidak melakukannya.

### Akibat bagi pertanyaan `_HITUNG`/`_SUNTING`

Tiket `17` menunggu S9 untuk tahu besaran mana yang benar-benar berpasangan hitung/sunting. **S9 tidak dapat menjawabnya**, karena rule yang dituju tidak menyebut satu pun besaran. Ini jawaban yang sah dan harus dicatat apa adanya: pertanyaannya tertuju pada rule yang keliru. Sasaran penggantinya belum ditetapkan.

### Temuan tambahan: penanda yang tidak pernah dibaca — EVIDENCED

`IsEditClaim` muncul **tiga kali** di seluruh empat lapisan:

| Berkas, baris | Peran |
|---|---|
| `Activity\EditXOLAlokasi.xml` | ditulis `1` |
| `Activity\CountLossAllocation_act.xml` 7322 | ditulis `0` pada `SpreadingRisk(<LAST>)` |
| `Activity\CountLossAllocation_act.xml` 9795 | ditulis `0` pada `SpreadingRisk(<LAST>)` |

**Nol pembacaan.** Tidak ada penjaga langkah, syarat tampil, ekspresi, maupun aturan `When` yang mengujinya. Penanda yang seharusnya melindungi suntingan petugas dari tertimpa hitung ulang **tidak pernah ditanya oleh siapa pun**. Hitung ulang menimpa apa pun, lalu menulis `0` ke penanda itu.

Ini kelas cacat yang sama dengan `CheckDateDOL_Act`: aturan yang ada bentuknya dan tidak ada akibatnya.

---

## S5 — siapa yang menulis `AlokasiXOLPaid`

Pola diperluas sesuai permintaan tiket: bukan hanya `Property-Set`, tetapi juga `Page-Copy`, `Page-New`, `RDB-List` ke page, dan hasil Report Definition. Enam berkas menyebutnya, seluruhnya Activity dan Section.

### EVIDENCED-NIHIL — tidak ada penulis pada empat lapisan

Setiap kemunculan berperan salah satu dari:

| Peran | Tempat |
|---|---|
| `pyStepsObjectName` — sasaran perulangan | 5 Activity; `GenerateCACNP_Act` memberi deskripsi **"Loop get alokasi XOL OLD"** |
| dibaca sebagai nilai selisih | 10 ekspresi di `CreateChildKomiteCNP_Act`, 9 di `SaveCNPLayerList_Act` |
| penjaga panjang `@LengthOfPageList(...)>0` | 4 tempat |
| `pyPageListProperty` di layar | `Section\AdjustmentDetailNP.xml` 21148 |

Nol `Property-Set` bersasaran `AlokasiXOLPaid`. Nol di RDBList, ConnectREST, ReportDefinition. Pengisinya berada **di luar ekspor yang dibaca**.

### Yang justru terbaca: daftar kolom muatan selisih — EVIDENCED

Ekspresinya seragam, berbentuk `<nilai sekarang> − Primary.AlokasiXOLPaid(Local.IdxLossOld).<nilai lama>`:

| Nilai sekarang | Dikurangi |
|---|---|
| `.AdjusterFeeRNM` | `.AdjusterFee` |
| `.CNPOthersFee`, `.CNPOthersFeeRNM` | `.CNPOthersFee` |
| `.CNPReinstatement`, `.CNPReinstatementRNM` | `.CNPReinstatement`, `.CNPReinstatementRNM` |
| `.GrossAdjustment` | `.ClaimAmountAdjust` |
| `.GrossValue` | `.ClaimSpreaded` |
| `.SalvageRNM`, `.SalvageValue` | `.Salvage` |

**DERIVED**: ini pola pengiriman selisih yang sama dengan temuan S1, kali ini pada jalur Komite alih-alih jalur kasir. Dua nama berbeda menunjuk besaran yang sama di dua sisi pengurangan — `.GrossAdjustment` lawan `.ClaimAmountAdjust`, `.GrossValue` lawan `.ClaimSpreaded`. Itu bukan salah tulis; sisi lama memakai nama lain.

Dua di antaranya berpenjaga `PaymentType`: `@if(.CNPReinstatement==0 && Local.PaymentType!=7, 0, ...)`. Jadi premi pemulihan **hanya** diselisihkan bila jenis pembayarannya `7`. Ini mengikat S3 ke S5.

---

## S6 — asal `"Previously Calculated UR"`

Pola sesuai permintaan tiket: potongan, bukan literal utuh.

### EVIDENCED — nilainya memang dirakit

`Activity\GenerateCACNP_Act.xml` baris **6643**:

```
Local.TreatyName  <-  "Previously Calculated " + .TreatyName
```

Lalu diuji utuh di baris **7223** dan **7425**:

```
Local.TreatyName == "Previously Calculated UR"
```

Kedua penjaga itu bercabang dua arah: benar → langkah 3, salah → langkah 2.

**DERIVED**: nilai itu lahir ketika `.TreatyName` bernilai `UR`. Artinya satu nama lapisan kasar punya **bentuk turunan kedua** yang beredar sebagai nilai `TreatyName` juga — bukan sebagai kolom penanda terpisah. Nama dan keadaan dicampur dalam satu kolom teks.

Penguat di lapisan layar: `Section\AdjustmentDetailNP.xml` baris 20626 memberi judul bagian **"Previously Calculated"**, dan baris 50190 memberi tombol **"Save Previously Paid"**. Jadi konsep ini punya tempatnya sendiri di layar.

### Akibat bagi tiket `17`

Tiket `17` menyatakan: bila S6 menemukan keadaan ketiga, kunci (klaim, mata uang) **belum cukup**. S6 menemukannya. Karena `"Previously Calculated UR"` beredar sebagai nilai `TreatyName`, dua baris dapat berbagi klaim dan mata uang yang sama sambil berbeda pada bentuk itu. Tiket `17` beku di balik gerbang; temuan ini menunggunya di sana, dan ia **menaikkan** pentingnya REQ-033.

---

## S7 — `Property-Remove` di `CountLossAllocation_act` langkah 10

Langkah dibaca utuh beserta tetangganya, baris 1990–2210.

### EVIDENCED — langkahnya memang kosong

| Bagian langkah | Isi |
|---|---|
| `pyStepsActivityName` | `Property-Remove` |
| `pyStepPageReference` | `RH_1.pySteps(10)` |
| `pyStepsObjectName` | **kosong** |
| `pyStepsDescription` | **kosong** |
| `pyStepsPreCondition` | **kosong** |
| `pyParamArray` | hanya cetakan parameter UI; **tidak ada satu pun pasangan nama-nilai** |

Parameter `Property` bagi metode ini bertanda `pyParametersParamReq: true` — **wajib** — dan tidak terisi.

**DERIVED**: langkah itu tidak membuang properti apa pun. Ia sisa yang tidak dihapus, bukan perilaku. Migrasi tidak perlu menirunya, dan tidak ada yang perlu ditanyakan kepada siapa pun tentangnya.

---

## S8 — empat penanda

Keempatnya disapu pada empat lapisan. Pernyataan "bersih" **tidak** berlaku untuk semuanya; hasilnya berbeda satu sama lain.

### `IsTreatyIn` — dibaca, EVIDENCED

Empat Activity mengujinya sebagai penjaga langkah: `CheckDateDOL_Act`, `CheckDateReceived_Act`, `CheckReportDate_Act`, `GetReportStatus_Act`. Ketiga yang pertama menguji **sebagai angka** (`==0`, `==1`); `GetReportStatus_Act` menguji **sebagai teks** (`=="1"`). Ketidakkonsistenan tipe yang sama dengan `PaymentType`.

### `IsReject` dan `IsCloseFile` — ditulis dan dibaca, EVIDENCED

Ditulis `Activity\CreateChildKomiteCloseNP_Act.xml`, dua kali masing-masing: sekali pada halaman sendiri (`.IsReject`, `.IsCloseFile`) dan sekali pada halaman anak (`ChildWorkPage.*`). Dibaca `Section\InputAcceptation.xml` sebagai syarat tampil wadah:

```
pyWorkPage.IsCloseFile==1 || pyWorkPage.IsReject==1
```

Keduanya berpasangan dalam satu syarat, tidak pernah sendiri.

### `IsAnyAcceptation` — dibaca sepuluh kali, **tidak pernah ditulis**, EVIDENCED-NIHIL

`Harness\OutstandingClaim.xml` mengujinya di sepuluh tempat: tujuh `pyReadOnlyCondition`, dua `pyDisabledWhen`, dan dua `pyCondition` — seluruhnya `=1` atau `!= 1`. Nol penugasan pada empat lapisan.

**DERIVED**: ini bayangan cermin `IsEditClaim`. Yang satu ditulis dan tidak pernah dibaca; yang lain dibaca dan tidak pernah ditulis. Keduanya kelas cacat yang sama: penanda tanpa pasangan. Sepuluh medan layar menjadi hanya-baca atau mati berdasar nilai yang tidak pernah diisi siapa pun pada lapisan yang terbaca.

---

## Dua sumber yang belum pernah dibaca

### `pengetahuan/DDL_Script_ClaimNonProp2.xls` — dibaca, EVIDENCED

**39 objek: 32 tabel dan 7 view**, seluruhnya milik `POOLDATA`.

> **KOREKSI 18 September 2026 — kalimat berikutnya di versi pertama salah, dan saya cabut.** Saya menulis bahwa berkas ini “membawa objek yang belum pernah masuk daftar mana pun”. **Tidak.** Rekonsiliasi terhadap berkas pertama menunjukkan berkas kedua **himpunan bagian murni**: ke-39 objeknya seluruhnya sudah ada di berkas pertama yang memuat 48, dan keempat tabel yang saya soroti di bawah **sudah punya berkasnya sendiri** di `pengetahuan/ddl/` sejak awal — `TABLE_JSON_KLAIM.sql`, `TABLE_DIRECTTOKASIR_LOG.sql`, `TABLE_CLAIMXOL2.sql`, `TABLE_CLAIMREJECTED.sql`.
>
> Yang benar-benar baru bukan objeknya melainkan **pembacaannya**: tidak ada dokumen sebelumnya yang menarik kesimpulan dari struktur keempat tabel itu. Temuan D13 dan D14 tetap berlaku apa adanya; klaim “sumber baru” yang tidak.

Empat di antaranya menyentuh tiket yang sudah ada.

#### `DIRECTTOKASIR_LOG` — arsip muatan keluar sudah ada di sistem lama

Kolomnya: `TGL_INPUT` (tanggal, bawaan `sysdate`), `DATA_JSON` (CLOB), `IDPEGA`, `NOAKSEPTASI`, `KET`.

**DERIVED**: sistem lama **sudah** mengarsipkan apa yang dikirim ke kasir — sebagai satu gumpalan JSON tanpa tipe. Ini menguatkan tiket `25` sekaligus menajamkan alasannya: yang ditambahkan sistem baru bukan arsipnya, melainkan **arsip yang berkolom**. Perhatikan juga apa yang **tidak ada** di sini: tidak ada penanda tolak. `STS_REJECT` yang terbaca pada jalur akseptasi tidak punya sepupu pada jalur kasir.

#### `JSON_KLAIM` — nomor klaim lama tidak dijamin unik

Kolomnya: `MNK_NO_KLAIM` (**NOT NULL**), `DATA_JSON` (CLOB, dengan `CHECK (DATA_JSON IS JSON)`), `IDPEGA`, `TGL_INPUT`, `TGL_KONVERSI`, `NOPOLIS` (**NOT NULL**), `IDPROD`, `STS_KONVERSI`.

Kunci primernya **`IDPEGA`**, bukan `MNK_NO_KLAIM`. Indeks penopangnya berkolom dua, `(IDPEGA, NOPOLIS)`, menopang kunci berkolom satu — sah sebagai awalan, dan patut dicatat.

**DERIVED, dan penting bagi REQ-018**: nomor klaim lama `MNK_NO_KLAIM` diwajibkan ada, tetapi **tidak diwajibkan unik oleh struktur mana pun**. Sampai sekarang kemungkinan nomor klaim ganda berstatus dugaan yang menunggu pencacahan. Sekarang ia **kemungkinan struktural**: tidak ada yang mencegahnya. REQ-018 berubah dari "mengukur apakah terjadi" menjadi "mengukur berapa banyak yang sudah terjadi".

#### `CLAIMXOL2` — kembaran, dan seluruh uangnya teks

Kolomnya: `CASEID`, `GrossAdjustment`, `CNPReinstatement`, `Currency`, `KursIDR`, `XOL` — **kelimanya `VARCHAR2(4000)`** — dan `TANGGAL`.

Dua hal. Pertama, ini nama nyaris kembar dengan `CLAIMXOL`, contoh kelima dari kelas cacat penamaan itu. Kedua, ia menguatkan `BLUEPRINT.md` §13.2 dan §13.5 dengan bukti kedua yang berdiri sendiri: nilai uang lama memang tersimpan sebagai teks. Dan `KursIDR` **berdampingan** dengan nilainya di satu baris — bukti langsung bagi usulan ADR-0029.

#### `CLAIMREJECTED` — daftar tolak, berkunci `pzInsKey`

Kolomnya: `INSKEY` (kunci primer), `ID`, `INSNAME`, `LABEL`, `STATUSWORK`, `CREATEOPNAME`, `CREATEOPERATOR`, `OBJCLASS`, `UPDATEDATETIME`, `UPDATEOPNAME`, `UPDATEOPERATOR`, `REMARK`.

**DERIVED — dikoreksi 18 September 2026, separuhnya gugur.**

Yang **berdiri**: penolakan dicatat sebagai **baris di tabel tersendiri** yang membawa pelaku dan alasan (`REMARK`), bukan sebagai kolom keadaan di tabel klaim. Itu preseden nyata bagi pola yang sama di sistem baru.

Yang **gugur**: “berkunci `pzInsKey`” bukan preseden yang ditiru. `pzInsKey` pegangan instance bawaan Pega, bukan pengenal bisnis, dan ADR-0021 menetapkan pengenal dinamai menurut isinya. Nasibnya sama dengan `CASEID`: tabel korelasi saja. Lihat `BLUEPRINT.md` §8.4a.

#### Objek lain yang baru terlihat

`ADJUSTERCONSULTANT`, `AGENT`, `BANKACCOUNT`, `BUSINESS`, `CATASTROPHE`, `EMAILKOMITE`, `JSON_POLIS`, `KODE_PRODUKSI`, `LST_BANK_GROUP`, `MARKETINGOFFICER`, `M_CLIENT`, `M_TREATY_IN`, `M_TREATY_IN_EDM`, `M_TREATY_OUT`, `M_TREATY_OUT_DETAIL`, `OS_AKSEPTASI_KLAIM`, `PROPORTIONALARRG`, `REINSURANCETYPE`, `RW`, `TREATYBUSINESS`, `TREATYCONTRACT`, `TREATYGROUP`, `TREATYINDETAIL`, `TREATYINDETAILEDM`, `TREATYINPRODUCTION`, `T_FOLDER_IMAGE`, `T_STORAGE_IMAGE`; view `CITY`, `CURRENCY`, `CURRENCYSTANDARD`, `PROVINCE`, `V_D_CAUSE_OF_LOSS`, `V_D_CAUSE_OF_LOSS_BUSINESS`, `V_M_CAUSE_OF_LOSS`, `V_POLIS`.

Dua di antaranya adalah persis dua tabel yang ditanyakan **REQ-033**: `PROPORTIONALARRG` dan `TREATYINDETAIL`. Strukturnya kini terbaca; sebarannya tetap butuh pencacahan.

Pasangan bersufiks `EDM` — `M_TREATY_IN` lawan `M_TREATY_IN_EDM`, `TREATYINDETAIL` lawan `TREATYINDETAILEDM` — belum diketahui maknanya. **TIDAK DITEMUKAN** pada lapisan yang dibaca.

### `excludeXML/GetBase64Attachment.xml` — **tidak dibaca**

Berkas itu berada di `D:\XML_NURE\Komite Claim Non Prop\excludeXML\`. **Folder itu tertutup dan saya tidak membukanya.** Sasaran ini dinyatakan tidak terjawab.

Yang **dapat** dibaca adalah saudaranya di folder yang terbuka, `Claim Non Prop\Activity\GetBase64Attachment.xml` (3 117 baris), dan itu saya baca. Ia memakai `Obj-Browse`, `Obj-Open-By-Handle`, `Page-Remove`, `Property-Remove`, `Property-Set`, `Java`, dan memanggil `ASM-FW-GISFW-Int-T_STORAGE_IMAGE.GetUrlGoogleStorage_Act`. Ia menulis ke `Primary.ClaimData.Attachment(...)`: `IMAGEID`, `AttachStream`, `PNOTE`, `URLPUBLIC`, `pyCategory`, `pyMemo`. Kategori yang diuji: `Premium`, `DLA`, `AcceptanceNote`, `DLARetro`.

**DERIVED**: lampiran klaim disimpan di luar basis data (Google Cloud Storage) dan dirujuk lewat `T_STORAGE_IMAGE`; yang tinggal di baris hanyalah pengenal, kategori, dan catatan. Ini bahan bagi tiket `20`. Apakah berkas di dalam folder tertutup berbeda dari yang ini **tidak diketahui**.

---

## Yang perlu masuk daftar

### Temuan yang mengubah putusan lama

| Kode | Temuan | Menyentuh |
|---|---|---|
| **D9** | Ketujuh nilai `PaymentType` EVIDENCED; penolakan saya atas §4.1 keliru karena pola sapuan hanya mencakup literal berkutip | tiket `07`, §4.1 |
| **D10** | `IsEditClaim` ditulis tanpa pernah dibaca; `IsAnyAcceptation` dibaca sepuluh kali tanpa pernah ditulis | tiket `17`, `31` |
| **D11** | `AlokasiXOLPaid` tanpa penulis pada empat lapisan | tiket `19` |
| **D12** | `"Previously Calculated UR"` dirakit, lalu diuji utuh — nama dan keadaan bercampur di satu kolom | tiket `17`, REQ-033 |
| **D13** | `JSON_KLAIM` mewajibkan `MNK_NO_KLAIM` ada tetapi tidak unik; kunci primernya `IDPEGA` | tiket `12`, REQ-018 |
| **D14** | `DIRECTTOKASIR_LOG` sudah menjadi arsip muatan keluar di sistem lama, tanpa penanda tolak | tiket `25` |
| **D15** | Nilai instruksi bayar menambahkan **tepat satu** besaran menurut `PaymentType`, bukan menjumlahkan seluruhnya | tiket `18`, `29` |

### Pertanyaan baru yang butuh jawaban orang

**Satu, dan hanya satu.** Tujuh properti bernilai uang hidup tanpa padanan IDR: `AdjusterFee`, `AdjusterFeeValue`, `Salvage`, `SalvageValue`, `CNPOthersFee`, `OthersFee`, `AdjustmentValue`, `TotalClaim`, `PremiumSpreaded`. Pertanyaannya bukan "kolomnya kurang atau tidak" — itu rancangan — melainkan **apakah nilai rupiah dari biaya penilaian dan salvage pernah dibutuhkan, dan bila ya, selama ini diambil dari mana**. Diajukan ke akuntansi, bukan ke DBA.

### Yang tetap tertutup

| Sasaran | Sebab |
|---|---|
| Arti tiap nilai `PaymentType` | **REQ-001**, definisi properti di schema PegaRULES |
| Penulis `AlokasiXOLPaid` | di luar ekspor yang dibaca |
| Isi `excludeXML/GetBase64Attachment.xml` | folder `Komite Claim Non Prop` tertutup |
| Makna sufiks `EDM` | tidak terbaca di lapisan mana pun |
| Sasaran pengganti bagi pertanyaan `_HITUNG`/`_SUNTING` | S9 membuktikan `EditXOLAlokasi` bukan sasarannya; penggantinya belum ditetapkan |
