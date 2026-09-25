# Hasil sapuan S10 – S16, sapu ulang S1/S2, dan rekonsiliasi DDL

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas XML** pada empat lapisan (Activity 114, RDBList 50, Section 36, ReportDefinition 22, FlowAction 16, Harness 14, DataTransform 10, ConnectREST 7, When 7, DecisionTable 1, Flow 1, SystemSettings 1); ditambah `pengetahuan/ddl/` (49 berkas) dan kedua workbook DDL.
> Folder `Komite Claim Non Prop` **tidak dibuka**.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut.

Dijalankan 18 September 2026. Gerbang ADR-0028 berlaku utuh: tidak ada DDL, tipe kolom usulan, nama constraint, nama index, maupun singkatan. Tipe dan nama sistem lama yang dikutip adalah **bukti**, bukan rancangan.

## Pola literal yang dipakai, supaya dapat diperiksa ulang

Setiap sapuan perbandingan memakai pola yang mencakup **ketiga bentuk** sekaligus:

```
<properti>\s*(==?|!=|<>)\s*( &quot;…&quot; | "…" | '…' | -?[0-9]+(\.[0-9]+)? )
```

Pembacaan diambil dari 26 tag, termasuk `pyBrowseSQL`, `pyTargetProperty`, dan `pyFieldName` — bukan hanya `pyStepsPreCondParamsWhen`. Nama telanjang di dalam SQL dan Report Definition dihitung sebagai pembacaan; tanpa itu `STS_REJECT` akan salah tergolong yatim.

---

# 0. Empat perbaikan laporan

## 0.1 — 273 versus 279: cacat saya, bukan berkas yang hilang

**Ekspor memuat 279 XML. Seluruhnya terbaca.** Angka 273 di stempel laporan sebelumnya salah, dan sebabnya satu perintah:

```
find . -iname "*.xml" -not -path "*Komite*"
```

Saringan itu dimaksudkan mengecualikan **folder** `Komite Claim Non Prop`. Tetapi ia menyaring **jalur**, dan enam berkas di dalam folder yang terbuka memuat kata itu pada **namanya**:

| Berkas | Lapisan |
|---|---|
| `Activity\CreateChildKomiteCNP_Act.xml` | Activity |
| `Activity\CreateChildKomiteCloseNP_Act.xml` | Activity |
| `Activity\ProteksiSendKomiteCNP_Act.xml` | Activity |
| `Harness\KomiteCNP.xml` | Harness |
| `ReportDefinition\FilterEmailKomiteWithLimit.xml` | ReportDefinition |
| `Section\KomiteCLMNP.xml` | Section |

**Yang tersaring hanya pencacahannya, bukan sapuannya.** Seluruh sapuan S3–S9 dijalankan dengan `grep -r` tanpa saringan itu, dan tiga dari enam berkas tersebut dikutip langsung di laporan S3 — `CreateChildKomiteCNP_Act` untuk nilai `PaymentType` `7`, `ProteksiSendKomiteCNP_Act` untuk `''` dan `7`. Bukti bahwa mereka terbaca ada di dalam laporannya sendiri.

**Akibatnya**: setiap EVIDENCED-NIHIL di laporan S3–S9 **tidak bersyarat**, dan **ADR-0006 tidak bersyarat**. Yang salah angka stempelnya; sudah dikoreksi jadi 279.

> **KOREKSI 18 September 2026 — kalimat aslinya salah dan saya cabut.** Saya menulis bahwa `Struktur_Flow_TreatyIn.xlsx` “belum pernah disebut di dokumen mana pun”. **Ia disebut di lima tempat**: `BLUEPRINT.md` baris 14 dan 35 (yang mencacahnya dan menurunkan graf pemanggilan darinya), `FINDING-003` baris 35 (yang **mencabutnya** sebagai bukti karena ia turunan), dan `_selesai/OPEN-QUESTIONS.md` E3 dan E4 (yang ditutup dengan menyebutnya sebagai sumber).
>
> Berkas itu **sudah dibaca di sesi sebelumnya**. Yang luput bukan pembacaannya melainkan perambatan aturan `FINDING-003`: berkas itu dinyatakan turunan dan tidak boleh berdiri sebagai bukti, tetapi E3 dan E4 dibiarkan berdiri di atasnya. Keduanya sudah dibasiskan ulang ke XML.
>
> Berkas kunci `~$Struktur_Flow_TreatyIn.xlsx` (165 byte) memang tidak pernah disebut, dan memang bukan sumber — ia sisa sesi Excel yang pernah terbuka.

## 0.2 — tujuh baris, sembilan nama

Diperbaiki di `SAPUAN-S3-S9.md`. Tabelnya tujuh baris karena dua baris memuat dua nama: `AdjusterFee`/`AdjusterFeeValue` dan `Salvage`/`SalvageValue`. Sembilan nama seluruhnya.

## 0.3 — satuan

Stempel memakai **byte di cakram**, dan sekarang menyebutnya begitu: `DDL_Script_ClaimNonProp2.xls` 51 712 byte, `DDL_Script_ClaimNonProp.xls` 143 872 byte.

## 0.4 — `pzInsKey` bawaan Pega

Diterapkan ke `BLUEPRINT.md` **§8.4a** yang baru, bersebelahan dengan catatan `IndexObject`, dan ke kode D. Ringkasnya:

| Pengenal lama | Nasibnya | Alasan |
|---|---|---|
| `pzInsKey` | tabel korelasi saja | pegangan instance; **paling andal di sisi lama**, tidak bermakna di sisi baru |
| `pxCoverInsKey` | tabel korelasi saja | menyimpan `pzInsKey` induk; relasi induk-anak baru memakai kunci sendiri |
| `IndexObject` | tidak dipakai | posisi, bukan pengenal |
| `CASEID` | tabel korelasi saja | salinan `pzInsKey` di sisi basis data |

Bacaan `CLAIMREJECTED` dikoreksi separuh: preseden **tabel tersendiri yang membawa pelaku dan alasan** berdiri; preseden **berkunci `pzInsKey`** gugur.

---

# 1. Sapu ulang S1 dan S2 dengan pola lengkap

## S2 — `CNPStatusCase`

**Temuan baru: tidak ada.** Pola tiga-bentuk menemukan nilai yang sama persis.

**Temuan lama yang berubah label: satu, dan arahnya mengeras.**

Dua nilai yang didaftar `MEMORI_PEMAHAMAN.MD` §7.3 — `"CLAIM ACCEPTED"` dan `"CLAIM REJECTED"` — disapu sebagai literal telanjang di seluruh 279 berkas dan muncul **nol kali**. Bukan nol sebagai nilai `CNPStatusCase`; nol **di mana pun**. Sebelumnya D1 hanya menyatakan ada nilai keempat yang tidak tercatat di memori. Sekarang lebih tajam: dari empat nilai yang didaftar memori, **hanya dua yang ada**, dan salah satu yang ada tidak didaftar memori.

Yang benar-benar tertulis, seluruhnya:

| Nilai | Ditulis oleh |
|---|---|
| `"COMITEE ACCEPTANCE (DEPT. HEAD)"` | `CreateChildKomiteCNP_Act`, `CreateChildKomiteCloseNP_Act` |
| `"INPUT ACCEPTATION CLAIM"` | `InputOutStandingCTNP_PostAct` |

**Temuan baru kedua, dari sisi lain**: `CNPStatusCase` **tidak pernah dibaca**. Tiga penulisan, nol pembacaan pada empat lapisan. Ia masuk daftar penanda yatim S11.

**Temuan ejaan**: `COMITEE` (satu T) dipakai sebagai **nilai tersimpan**, 2 kali. `COMMITTEE` (dua T) muncul 10 kali di 7 berkas sebagai teks tampilan. Ejaan yang salah adalah yang masuk basis data.

**Tiket `06` tidak kembali ke papan aktif** — nilai yang ditemukannya tetap, dan yang bertambah adalah dua pernyataan nihil yang menguatkannya. Bila Anda menghendaki sebaliknya, ia dikembalikan.

## S1 — `InputParamOs.*`

**Temuan baru: tidak ada pada parameternya.** Tujuh belas nama, sama persis.

**Temuan baru pada cakupannya: lima berkas yang tidak disebut laporan S1.** Parameter yang sama juga diisi atau dibaca oleh `Activity\SaveToOS.xml`, `Activity\GetSelisihActual_Act.xml`, `Activity\CloseClaimTNonProp.xml`, `Harness\Hitung_Test.xml`, dan `Section\Hitung_Test.xml`. Dua yang terakhir adalah layar uji — nama `Hitung_Test` bukan rule produksi, dan ia memuat medan yang terikat langsung ke `InputParamOs.GrossValue` dan `InputParamOs.Value`.

**Temuan baru ketiga: sisi kasir belum pernah diinventarisasi.** `HitServiceToKasir_Act` memakai **24 slot `TempKasir.CARI1`…`CARI24` ditambah `CARIDATETIME`** — 25 medan bernama urut, bukan bernama isi. Laporan S1 memetakan sisi Arasapas (17 parameter bernama) tetapi tidak sisi kasir.

**Tiket `05` kembali ke papan aktif.** Alasannya bukan pola literal, melainkan cakupan: pernyataannya tentang "kolom yang sesungguhnya diterima **Arasapas dan kasir**", dan sisi kasir belum terbaca.

## Aturan yang berlaku ke depan

Setiap pola sapuan wajib mencakup ketiga bentuk literal, dan **polanya ditulis di badan laporan**. Sudah diterapkan di kepala dokumen ini.

---

# 2. Rekonsiliasi `DDL_Script_ClaimNonProp2.xls`

## Gabungan kedua berkas — dan koreksi atas laporan saya

| | Objek |
|---|---|
| `DDL_Script_ClaimNonProp.xls` | **48** |
| `DDL_Script_ClaimNonProp2.xls` | **39** |
| Tumpang-tindih | **39** |
| Hanya di berkas 2 | **0** |
| Hanya di berkas 1 | **9** |
| **Gabungan, objek berbeda** | **48** |

**Berkas kedua adalah himpunan bagian murni.** Ia tidak membawa satu objek pun yang belum ada. Kesembilan yang hanya ada di berkas pertama: `CLAIMXOL`, `GETCURRENCYSTANDARD`, `GET_TOKEN_STORAGE`, `PEGA_D_CAUSE_OF_LOSS`, `PEGA_JSON_KLAIM_PNC`, `PEGA_JSON_OS_AKSEP_KLAIMTNP`, `PEGA_M_CAUSE_OF_LOSS`, `PROC_GENERATE_SEQUENCE_NUMBER`, `V_MST_USER_TEKNIS` — yaitu view `CLAIMXOL` dan seluruh procedure serta function.

Laporan S3–S9 menyatakan berkas kedua "membawa objek yang belum pernah masuk daftar mana pun". **Itu salah dan sudah dicabut di tempatnya.** Keempat tabel yang saya soroti sudah punya berkasnya sendiri di `pengetahuan/ddl/` sejak awal. Yang baru bukan objeknya melainkan pembacaannya; D13 dan D14 tetap berlaku.

## REQ-017 — berapa dari 14 yang tertutup

**Nol.**

Sapuan atas seluruh 49 berkas `ddl/` mencari objek yang **dirujuk** oleh `FROM`, `JOIN`, `INTO`, `UPDATE` tetapi tidak punya berkas DDL, lalu menyingkirkan yang ternyata variabel PL/SQL:

| Objek dirujuk tanpa DDL | Dirujuk oleh |
|---|---|
| `M_CURRENCYSTANDARD` | `FUNCTION_GETCURRENCYSTANDARD`, `VIEW_CURRENCYSTANDARD` |
| `M_CURRENCY` | `VIEW_CURRENCY` |
| `M_NATION` | `VIEW_CURRENCY`, `VIEW_PROVINCE` |
| `M_PROVINCE` | `VIEW_PROVINCE` |
| `M_CAUSE_OF_LOSS` | `PROCEDURE_PEGA_M_CAUSE_OF_LOSS`, `VIEW_V_M_CAUSE_OF_LOSS` |
| `D_CAUSE_OF_LOSS` | `PROCEDURE_PEGA_D_CAUSE_OF_LOSS`, dua view |
| `M_SITE_DATABASE` | dua procedure cause-of-loss |
| `MST_USER_TEKNIS` | `VIEW_V_MST_USER_TEKNIS` |
| `OS_AKSEPTASI_SUBJECTIVITY` | `VIEW_CLAIMXOL` |
| `TANGGAL_CLOSING` | `PROCEDURE_PROC_GENERATE_SEQUENCE_NUMBER` |
| `GCP_IMAGE` | `PROCEDURE_GET_TOKEN_STORAGE` |
| `BRANCH`, `CITYINPUT`, `DISTRICTINPUT`, `RWINPUT` | `VIEW_CITY` |

Disingkirkan sebagai variabel PL/SQL, bukan objek: `V_DAY_CLOSING`, `V_PERIODE_DATE`, `V_SEQ`, `VHASIL`, `VAKSESTOKEN`, `ID_COUNT`, `ID_SITE`, `DULU` — kedelapannya dideklarasikan di badan procedure-nya (`v_day_closing NUMBER;` dan seterusnya di `PROC_GENERATE_SEQUENCE_NUMBER` baris 13–18).

### `CURRENCYSTANDARD` bukan `M_CURRENCYSTANDARD` — dan buktinya langsung

Peringatan Anda tepat. `pengetahuan/ddl/VIEW_CURRENCYSTANDARD.sql`:

```
CREATE OR REPLACE FORCE EDITIONABLE VIEW "POOLDATA"."CURRENCYSTANDARD" (...) AS
  SELECT a.ID, a.CurrencyValue, a.CurrencyDate, a.UserID, a.InputDate
     FROM m_CurrencyStandard a
```

View itu **membaca** `m_CurrencyStandard`. Keduanya objek berbeda: satu view, satu tabel basis. Memilikinya **tidak** menutup REQ-017, dan blocker-nya tetap berdiri. Ini contoh keenam dari kelas cacat nama nyaris kembar.

Hal yang sama berlaku untuk `V_MST_USER_TEKNIS` (dimiliki) versus `MST_USER_TEKNIS` (tidak dimiliki).

### Temuan yang menyentuh ADR-0016 — database link sudah dipakai hari ini

`VIEW_V_MST_USER_TEKNIS.sql` baris 10:

```
FROM hrdasm.v_hrd_mst@asmd.sinarmas.co.id b
```

**Ini database link ke instance Oracle lain, dan ia sudah berjalan di sistem lama.** ADR-0016 melarang database link di sistem baru; larangan itu tidak berubah. Yang berubah adalah kedudukannya: ia bukan lagi larangan atas sesuatu yang belum pernah ada, melainkan atas sesuatu yang **sudah dipakai produksi**. Satu-satunya pemakaian pada 49 berkas DDL; sapuan mencari pola `@<host>` di seluruh berkas.

Akibat praktisnya ada di REQ-021 dan ADR-0028: bila sistem baru membutuhkan data HRD yang sama, jalurnya harus ditetapkan tanpa link — dan itu belum dibahas di ADR mana pun. **Dicatat, tidak diputuskan.**

## `PULL-LIST.csv`

**25 dari 74 tetap belum punya DDL.** Cacah ulang tidak mengubah apa pun, karena berkas kedua tidak membawa objek baru. Ke-25 itu terdiri dari tiga jenis yang tidak setara:

| Jenis | Cacah | Contoh |
|---|---|---|
| Objek basis data nyata | 14 | `M_CURRENCYSTANDARD`, `MST_USER_TEKNIS`, `TANGGAL_CLOSING`, `OS_AKSEPTASI_SUBJECTIVITY`, `D_CAUSE_OF_LOSS`, `M_CAUSE_OF_LOSS`, `M_CURRENCY`, `M_NATION`, `M_PROVINCE`, `M_SITE_DATABASE`, `BRANCH`, `CITYINPUT`, `DISTRICTINPUT`, `RWINPUT` |
| Tarikan profil data, bukan DDL | 4 | `(CACAH BARIS TREATYTYPE=UR PER KLAIM)`, `(STATUS BUKA/TUTUP 8 KLAIM BERTAMBALAN)`, dan dua lainnya |
| Objek Pega / lintas instance | 7 | `PR4_*`, `DATA-ADMIN-DB-TABLE`, `TABEL WORK ASM-FW-GCNMFW-WORK`, `HRDASM.V_HRD_MST@ASMD.SINARMAS.CO.ID`, `F_GET_EMAIL`, `PLATNP_SEQ`, `TREATY_OUT`, `TRLOSS_DETAIL_T` |

Angka "25 belum punya DDL" karena itu **menyesatkan bila dibaca sebagai satu antrean**. Yang benar-benar menahan REQ-017 adalah 14 baris pertama.

## `SCHEMA-ACTUAL.csv`

**Tidak ada objek baru untuk diisikan** — berkas kedua tidak membawa satu pun. 698 baris tetap.

## REQ-033 — pertanyaannya menyempit

`PROPORTIONALARRG` dan `TREATYINDETAIL` sudah punya berkas DDL-nya (`TABLE_PROPORTIONALARRG.sql`, `TABLE_TREATYINDETAIL.sql`), jadi **strukturnya terbaca**. Yang tersisa murni **sebaran**: apakah kombinasi `(ID, LAYER, LAYERTYPE)` pernah melahirkan lebih dari satu `LAYERPART`. Dicatat di REQ-033.

---

# 3. S10 — rantai satu nilai uang, dari masuk sampai keluar

## 3.1 `.GrossValue`

### Temuan pertama: satu besaran, tiga nama

Uang yang sama berpindah nama dua kali sepanjang perjalanannya, dan tidak satu pun peralihan itu diberi keterangan:

| Tahap | Nama | Tempat |
|---|---|---|
| keluaran mesin alokasi | `SpreadingRisk.ClaimSpreaded` | `CountLossAllocation_act` 9474 |
| muatan ke Komite | `CNPCurrencyList.GrossValue` | `CreateChildKomiteCNP_Act` 4948, `SaveCNPLayerList_Act` 1049 |
| muatan keluar | `InputParamOs.GrossValue` | `SaveDataToOSAksep_Act` 7282, `SaveToOS` 2322 |

Peralihan pertama lurus: `GrossValue <- Primary.SpreadingRisk(Local.IdxAdjustment).ClaimSpreaded`, berpenjaga `Local.XOL==.XOL`.

### Temuan kedua: dua jalur menghitung "nilai yang sama" secara berbeda

| Jalur | Ekspresi | Berkas, baris |
|---|---|---|
| Komite | `GrossValue <- .ClaimSpreaded` | `CreateChildKomiteCNP_Act` 4948 |
| Arasapas | `InputParamOs.GrossValue <- .ClaimEstimation + .AdjClaimValue` | `SaveDataToOSAksep_Act` 7282; `SaveToOS` 2322 |

`.ClaimSpreaded` dan `.ClaimEstimation + .AdjClaimValue` **bukan besaran yang sama**. Mesin alokasi menulis keduanya secara terpisah (9258, 9474) dengan rumus berbeda. Dua hilir karena itu menerima angka yang berbeda untuk nama yang sama. **Ini temuan, bukan kesimpulan**: mana yang benar tidak terbaca dari XML.

### Temuan ketiga: pengurangan selisih terjadi ke dua arah

| Ekspresi | Penjaga | Berkas, baris |
|---|---|---|
| `InputParamOs.GrossValue - .TotalClaim` | `.TreatyName=="UR"` | `SaveToOS` 2860 |
| `InputParamOs.GrossValue - OutOSAcc.pxResults(1).GrossValue` | `.TreatyName=="UR"` | `SaveToOS` 3435; `SaveDataToOSAksep_Act` 7858 |
| `.TotalClaim - InputParamOs.GrossValue` | **tidak ada penjaga pada langkah itu** | `SaveToOS` 3999 |
| `.TotalClaim - InputParamOs.GrossValue` | `pyWorkPage.FlagActualPremium==true` | `GetSelisihActual_Act` 3538 |
| `InputParamOs.GrossValue - .TotalClaim` | tidak ada | `GetSelisihActual_Act` 2777 |

Dua bentuk saling terbalik tandanya, hidup berdampingan di **satu berkas yang sama** (`SaveToOS` 2860 dan 3999). Salah satunya berpenjaga `TreatyName=="UR"`, satunya tidak berpenjaga sama sekali. **Mana yang berlaku bergantung pada langkah mana yang tercapai lebih dulu**, dan itu tidak terbaca tanpa menjalankan.

### Temuan keempat: nol ditulis, bukan dikosongkan

`SaveAdjustmentToOSAksep_Act_Tes` menulis `GrossValue <- 0` di tiga tempat (2205, 2817, 3999) sebelum mengisinya. Itu inisialisasi memakai nol, bukan ketiadaan — bahan langsung bagi tiket `33` (*NULL bukan nol*).

### Konversi mata uang — di mana terjadi dan di mana tidak

| Titik | Konversi? | Bukti |
|---|---|---|
| mesin alokasi, `ClaimAmountIDR` | **ya, bersyarat** | `CountLossAllocation_act` 9585: `@if(Local.Currency=="IDR", ClaimEstimation, ClaimEstimation*Local.Kurs)` |
| mesin alokasi, `KursIDR` | disimpan | 9690: `KursIDR <- Local.KursIDR` |
| `SetActualPremium_ACT`, `ClaimAmountIDR` | **TIDAK** | 681: `ClaimAmountIDR <- SpreadingRisk(<LAST>).TotalClaim` |
| seluruh rantai `GrossValue` | **TIDAK, di mana pun** | tidak ada `*Kurs` maupun `KursIDR` pada rantai itu |
| `CountTotalInsterest_Act` | akumulasi ke `Local.GrossValueIDR` | 5605: `+ .ConvertGrossEstimasi` — konversinya terjadi di hulu, di properti lain |

**Temuan yang perlu ditandai**: `ClaimAmountIDR` diisi oleh **dua rule dengan perlakuan berbeda**. Mesin alokasi menguji mata uangnya lebih dulu; `SetActualPremium_ACT` menyalin `TotalClaim` apa adanya ke kolom berlabel IDR **tanpa menguji mata uang dan tanpa mengalikan kurs**. Bila `TotalClaim` bukan rupiah, kolom IDR memuat angka yang bukan rupiah.

Ini menguatkan ADR-0029 dari arah yang berbeda dari S4: bukan *kursnya tidak disimpan*, melainkan *kolom IDR dapat terisi tanpa kurs pernah dipakai*.

### Skala — di mana dinyatakan

| Tempat | Skala |
|---|---|
| seluruh rantai `.GrossValue` | **tidak dinyatakan di satu titik pun** |
| `SetActualPremium_ACT` 639 | `@divide(Local.ActualPremi, @divide(.ClaimPercentage,100,20), 20)` — **20**, dua kali |
| `CountLossAllocation_act` 9732 | `@toDecimal(...)` tanpa skala |

Rantai `GrossValue` melewati enam berkas dan **nol pernyataan skala**. Ia bergantung penuh pada perilaku desimal bawaan Pega — yang persis alasan ADR-0003 ada.

## 3.2 `.AdjusterFeeValue`

Rantainya **lebih pendek dan lebih rapi**, dan perbandingannya yang berguna.

| Tahap | Ekspresi | Berkas, baris |
|---|---|---|
| lahir | `AdjusterFee / @divide(ClaimPercentage,100,20)` | `CreateChildKomiteCNP_Act` 4706; `SaveCNPLayerList_Act` 823 |
| lahir, bentuk kedua | `@divide(.AdjusterFee, @divide(TreatyInMaster.RNMShare,100,20), 20)` | `CreateChildKomiteCNP_Act` 5473; `SaveCNPLayerList_Act` 1461 |
| selisih | `.AdjusterFeeValue - @divide(AlokasiXOLPaid(IdxLossOld).AdjusterFee, @divide(RNMShare,100,20), 20)` | `CreateChildKomiteCNP_Act` 7996; `SaveCNPLayerList_Act` 2683 |
| selisih, bentuk kedua | `.AdjusterFeeValue - PrevData.AdjusterFeeValue` | `SaveAdjustmentToOSAksep_Act_Tes` 4281 |
| dibaca keluar | `@toDecimal(TempKasir.CARI9) + .AdjusterFeeValue` | `HitServiceToKasir_Act` 7820, berpenjaga `PaymentType = 1 \|\| 2 \|\| 5` |

**Skala dinyatakan, dan konsisten: 20, di setiap `@divide`.** Berbeda tajam dari `GrossValue` yang tidak menyatakannya sama sekali. Ini bukti langsung bagi catatan `arithmetic-inventory.tsv`: skala mengikuti **modul**, bukan jenis nilai. Dua besaran uang yang berjalan berdampingan di berkas yang sama diperlakukan berbeda.

**Dua pembagi berbeda untuk besaran yang sama**: `ClaimPercentage` (4706) dan `TreatyInMaster.RNMShare` (5473). Keduanya melahirkan `AdjusterFeeValue` di rule yang sama. Apakah keduanya selalu bernilai sama tidak terbaca dari XML — **TIDAK DITEMUKAN**.

**Konversi mata uang: nol di seluruh rantai.** Konsisten dengan S4 — besaran ini memang tidak punya padanan IDR, dan rantai ini menunjukkan ia juga tidak pernah dikonversi di jalan. Ia masuk instruksi bayar dalam mata uang aslinya.

---

# 4. S11 — penanda yatim, tiga golongan

104 kandidat disapu pada empat lapisan: berawalan `Is`, `Has`, `Flag`, `Sts`, `Status`, ditambah setiap properti yang seluruh nilainya `0`/`1`/`true`/`false`/`Y`/`T`. Properti kerangka Pega dan HTML dikecualikan dengan daftar tertutup.

| Golongan | Cacah |
|---|---|
| 1 — ditulis **dan** dibaca | 40 |
| 2 — ditulis, **tidak pernah** dibaca | 14 |
| 3 — dibaca, **tidak pernah** ditulis | 42 |

## Golongan 2 — ditulis tanpa pernah dibaca (14)

| Penanda | Ditulis | Nilai | Tempat |
|---|---|---|---|
| `IsEditClaim` | 3x | `0`, `1` | `CountLossAllocation_act` 7322, 9795; `EditXOLAlokasi` |
| `CNPStatusCase` | 3x | dua teks status | `CreateChildKomiteCNP_Act`, `CreateChildKomiteCloseNP_Act`, `InputOutStandingCTNP_PostAct` |
| `FlagFullLimit` | 1x | `0` | `CountLossAllocation_act` 9073 |
| `FlagOutstanding` | 1x | `"1"` | `SaveDataToOSAksep_Act` 9087 |
| `IsFlagError` | 1x | `"true"` | `CreateChildKomiteCNP_Act` 14665 |
| `KomiteCount` | 2x | `1` | `CreateChildKomiteCloseNP_Act` 1121 |
| `IsUW` | 1x | `"true"` | `GetDetailPolis_act` 442 |
| `IsRealisation` | 1x | `@if(@SizeOfPropertyList(Polis.pxResults)>0,1,0)` | `SetValueClaimTNP_Act` 1613 |
| `AktifButton` | 2x | `"1"`, `0` | `AddAkseptasiCNP_Act` 3042 |
| `AlokasiIDR` | 1x | `1` | `CountLossAllocation_act` 6368 |
| `CNPFlagReinstate` | 1x | `.CNPFlagReinstate` | `CreateChildKomiteCNP_Act` 2088 |
| `DATASHOW` | 1x | `1` | `SetViewDataMasterTreaty_Act` 686 |
| `IndividualRiskPercentage` | 1x | `0` | `AddAkseptasiCNP_Act` 3218 |
| `HasilEnd` | 1x | ekspresi tanggal | `SetEndDate_Act` 751 |

`CNPStatusCase` adalah yang paling berat di daftar ini: ia **status kasus**, ditulis di tiga tempat, dan tidak satu pun rule pada empat lapisan yang mengujinya. Bila statusnya benar-benar mengatur sesuatu, pengaturnya berada di luar ekspor.

## Golongan 3 — dibaca tanpa pernah ditulis (42), dan tidak seluruhnya cacat

**Ini pembedaan yang harus dibuat, kalau tidak daftarnya menuduh terlalu banyak.** Golongan 3 memuat dua hal berbeda:

### 3a. Kolom sistem lain, dibaca lewat SQL atau Report Definition — bukan cacat

`STS_ADJ`, `STS_ADJUSTER`, `STS_REG`, `STS_SALVAGE`, `STS_SURVEY`, `STATUSWORK`, `STS_PKP`, `MOSTATUS`, `StatusActive`, `StatusRandomTF`, `StatusSyariah`, `StatusService`, `StatusServiceKasir`.

Ketiga belasnya dibaca dari **tabel atau layanan di luar modul ini**, lewat `pyBrowseSQL`, Report Definition, atau ConnectREST. Penulisnya memang di luar ekspor, dan itu wajar. Mereka **tidak** yatim; mereka milik orang lain.

### 3b. Penanda pada objek kerja sendiri, dibaca tanpa pernah diisi — **ini yang cacat**

| Penanda | Dibaca | Nilai yang diuji | Tempat utama |
|---|---|---|---|
| `IsAnyAcceptation` | 24x | `1` | `Harness\OutstandingClaim` — 7 `pyReadOnlyCondition`, 2 `pyDisabledWhen`, 2 `pyCondition` |
| `IsTreatyIn` | 20x | `0`, `1`, `"1"` | `CheckDateDOL_Act`, `CheckDateReceived_Act`, `CheckReportDate_Act`, `GetReportStatus_Act` |
| `PremiOgp` | 32x | `0`, kosong | `CountNetPremi_act` |
| `PremiOnp` | 25x | `0`, kosong | `CountNetPremi_act` |
| `StsKatastrofe` | 17x | `Catastrophe`, `Non-Catastrophe` | `SetDefNonCatastrope_Act` |
| `ReporterStatus` | 15x | `1`, `2`, `3` | `GetReportStatus_Act` |
| `Limit2` | 15x | `0` | `CountLossAllocation_act` 5931 |
| `DeductibleType` | 12x | `true`, `false` | `CountClaimTNP_Act` |
| `PrintFaceClaim` | 10x | `1`, kosong | `CopyOldataCurr_act` |
| `isDouble` | 10x | — | `CountResult1Onp_Act` |
| `FlagProrate` | 8x | `0`, `1`, `2`, `3` | `Section\ReinstatementPremiumDetails` |
| `FlagActualPremium` | 7x | `false` | `GetSelisihActual_Act` |
| `IsAnalisTransfer` | 7x | `1` | `SetTPLNote_Act` |
| `Deduction1` | 6x | `0`, kosong | `CountNetPremi_act` |
| `isEstimation` | 6x | `1` | `Harness\OutstandingClaim` |
| `IsSurveyReport` | 5x | `No`, kosong | `Section\ViewDetailDeptHeadTreatyIn_UW` |
| `IsEditRNMShare` | 4x | `true`, `false` | `Harness\OutstandingClaim` |
| `CNPFlagXOL` | 4x | — | `CountLossAllocation_act` 3714 |
| `ShareCeding` | 4x | `1` | `GenerateCACNP_Act` |
| `FlagPPH` | 2x | `true` | `SetPPNPPH` |
| `IsAdjVal` | 2x | — | `CountTotalInsterest_Act` |
| `IsOldData` | 2x | `Yes` | `CountTotalInsterest_Act` |
| `IsSubjectivity` | 2x | — | `CreateChildKomiteCNP_Act` |
| `AgingStatus` | 2x | — | `Harness\DetailPaymentStsCNP` |
| `StatusKonversi` | 1x | `1` | `HitServiceToKasir_Act` |
| `WorkStatus` | 1x | — | `ASMForceCaseClose` |

Dua puluh enam penanda. Sebagian punya penjelasan yang sah dan belum terbukti: `IsSubjectivity` ditulis balik oleh **modul Komite** menurut `BLUEPRINT.md` §8.3, dan modul itu tertutup. Tetapi `FlagProrate` dengan empat nilai (`0`–`3`), `ReporterStatus` dengan tiga, dan `StsKatastrofe` dengan dua nilai teks **bukan penanda biner** — mereka domain kecil yang menentukan cabang, dibaca belasan kali, dan tidak pernah diisi pada lapisan mana pun.

## Yang ditentukan daftar ini

Tiket `17` dan `31` membutuhkannya untuk memutuskan **dimigrasi atau ditinggalkan**. Garis yang terbaca dari bukti:

- Golongan 1 (40) — dimigrasi; mereka punya penulis dan pembaca.
- Golongan 2 (14) — **tidak ada kolomnya di sistem baru**, kecuali penulisnya ternyata bermakna di luar ekspor. `CNPStatusCase` perlu diputuskan tersendiri karena ia status kasus.
- Golongan 3a (13) — bukan milik modul ini; tidak dimigrasi, dibaca lewat batas sistem.
- Golongan 3b (26) — **tidak dapat diputuskan dari XML saja.** Penulisnya ada di modul Komite yang tertutup, atau tidak ada sama sekali. Membuat kolomnya berarti menyediakan tempat bagi nilai yang mungkin tidak pernah datang; tidak membuatnya berarti memutus belasan cabang layar.

**Saya tidak memutuskan golongan 3b.** Ia menunggu pembukaan folder Komite, dan sampai itu tiket `17` dan `31` punya penahan yang sekarang bernama.

---

# 5. S12 — ketidakkonsistenan tipe

**22 properti** dibandingkan sebagai angka di satu tempat dan sebagai teks di tempat lain. **11** di antaranya punya nilai yang **hanya** bermakna dalam bentuk teks.

| Properti | Bentuk angka | Bentuk teks | Nilai yang hanya ada sebagai teks |
|---|---|---|---|
| `PaymentType` | `1`–`5` — `CreateChildKomiteCNP_Act` 13 | `2`, `3`, `7`, kosong — `HitServiceToKasir_Act` 605 | **kosong** |
| `AcceptanceStatus` | `0`, `1`, `2` — `CountSpreadingXOL` 2603 | `0`, `1`, `2`, kosong — `CloseClaimTNonProp` 399 | **kosong** |
| `IsPLA` | `1` — `GeneratePlaCNP_Act` 616 | `1A` — `Section\OutstandingClaim(1)` 105513 | **`1A`** |
| `CARI1` | `1` — `AddAkseptasiCNP_Act` 1316 | `ClaimAmount`, `EDITCLAIMXOL`, `Estimation`, `Interest`, kosong — `CopyOldataCurr_act` 746 | kelimanya |
| `TypeDeductible` | `0`, `1`, `2` — `AddAkseptasiCNP_Act` 3292 | `1`, kosong — sama | **kosong** |
| `Payable` | `1`, `2`, `3` — `SetPayableTreatyNP_Act` 2569 | kosong — `ProteksiSendKomiteCNP_Act` 5484 | **kosong** |
| `Layer` | `1` — `DetailCalculation` 4969 | kosong — `DetailCalculation` 5091 | **kosong** |
| `IsCFS` | `1` — `Section\InputAcceptation` 111333 | `1`, kosong — `Harness\OutstandingClaim` 44860 | **kosong** |
| `CheckCurr` | `0` — `GenerateCFS_act` 3596 | kosong — `SetAccoutNo_Act` 3510 | **kosong** |
| `Deduction1` | `0` — `CountNetPremi_act` 1314 | kosong — sama | **kosong** |
| `PrintFaceClaim` | `1` — `CopyOldataCurr_act` 2448 | kosong — `Harness\OutstandingClaim` 38166 | **kosong** |
| `CARI2` | `1` — `AddAkseptasiCNP_Act` 1456 | kosong — `GetReportStatus_Act` 1203 | **kosong** |
| `IsTreatyIn` | `0`, `1` — `CheckDateDOL_Act` 3039 | `1` — `GetReportStatus_Act` 1800 | — |
| `FormType` | `1`, `2` | `1`, `2` — `CountClaimTNP_Act` 3017 | — |
| `TPLFormat` | `1`, `2` | `1`, `2` — `SetTPLNote_Act` 482 | — |
| `Type` | `3` — `Harness\KomiteCNP` 3577 | `3` — `HitServiceToKasir_Act` 605 | — |
| `IsKomite`, `IsOutstanding`, `CNPFlagOuts`, `Appendflag`, `CARI10`, `ViewState` | `1`/`0` | `1`/`0` | — |

## Yang ditentukan daftar ini

**Sepuluh properti tidak dapat diberi domain numerik tertutup**, karena nilai kosong terbukti diuji dan tidak punya padanan angka. Ini pola yang sama yang membuat saya menolak `CHECK` pada `JENIS_PEMBAYARAN` — dan penolakan itu tetap berlaku untuk bagian ini, meski alasan `4`,`5`,`6` sudah gugur.

`IsPLA` berdiri sendiri: nilai `1A` bukan angka sama sekali, dan tidak dapat disimpan di kolom numerik tanpa kehilangan. Nama "Is…" yang bernilai `1A` juga berarti ia **bukan penanda biner**, meski namanya mengaku begitu.

Kumpulan ini beku di balik gerbang; ia tinggal dipakai saat tiket domainnya dibuka.

---

# 6. S13 — `CountLossAllocation_act` dibaca utuh

**22 langkah tingkat atas.**

| # | Metode | Penjaga | Sasaran / keterangan |
|---|---|---|---|
| 1 | Page-Copy | — | — |
| 2 | Property-Set | — | — |
| 3 | Property-Set | `Param.Currency==.Currency` | `ClaimData.ListClaimAmount` |
| 4 | Property-Set | `Param.Currency==.Currency` | *pct* |
| 5 | Property-Set | `Param.Calculation=="pct"` | *amount* |
| 6 | Property-Set | `Param.Calculation=="amount"` | — |
| 7 | Property-Set | `Param.Calculation=="CountXOL"` | **"PERHITUNGAN XOL MULAI DARI SINI"** |
| 8 | Page-New | `Param.Calculation=="CountXOL"` | `TempTotalClaim` |
| 9 | Page-New | — | `TempLimit` |
| **10** | **Property-Remove** | — | **kosong — lihat S7** |
| 11 | Property-Set | `TreatyInMaster.TreatyGroup==.TreatyGroup` | `TreatyInMaster.Limits` |
| 12 | Property-Set | `Local.Check==0` | `TreatyInMaster.Limits`; deskripsi memuat `\\\.EDMState=="3"` |
| 13 | Property-Set | `Local.Check==0` | `ClaimData.CNPSpreadLoss` — *Get Total Claim per currency* |
| 14 | Property-Set | `Local.Check==0` | — |
| 15 | Property-Set | `Local.CurencyLayer==.Currency` | `TreatyInMaster.CurrencyList` |
| 16 | Property-Set | `Local.CurencyLayer==.Currency` | — |
| 17 | Property-Set | `ProtectTreatyName.CARI1==1` | `TempTotalClaim.pxResults` — **"Let the fun BEGAN (BUAT APPEND ALOKASI)"** |
| 18 | Property-Set | `@LengthOfPageList(TempTotalClaim.pxResults)==Local.LoopLimit` | — |
| 19 | Property-Set | `Local.Currency==.Currency` | `ClaimData.SpreadingRisk` — *UNTUK TOTAL ALOKASI XOL* |
| 20 | Property-Set | `Local.Currency==.Currency` | `ClaimData.SpreadingRisk` |
| 21 | Property-Set | — | `ClaimData.ListTotalEstimation` — *Spreading* |
| 22 | Call | — | `CountReinstatement_Act` |

## 6.1 Retensi Cedant di-`APPEND`, dan tidak ada penghapusan

Dua tempat menulis `SpreadingRisk(<APPEND>)`:

| Baris | Baris pertama blok | Isi blok |
|---|---|---|
| **7024** | `TreatyType <- "UR"`, `TreatyName <- "UR"` | baris Retensi Cedant — 14 penugasan, ditutup `IsEditClaim <- 0` di 7322 |
| **9191** | `TreatyType <- ListSpreading.pxResults(1).ID`, `TreatyName <- ...Note` | baris layer — 25 penugasan, ditutup `IsEditClaim <- 0` di 9795 |

**Mengapa tidak ada penghapusan daftar sebelum itu**: satu-satunya `Property-Remove` di seluruh berkas adalah **langkah 10 yang kosong** — tanpa sasaran, tanpa parameter, tanpa prasyarat. Tidak ada `Page-Remove` atas `SpreadingRisk` di mana pun.

**DERIVED**: langkah 10 adalah tempat penghapusan itu *seharusnya* berada, dan ia tidak pernah diisi. Daftar karena itu menumpuk setiap kali mesin dijalankan ulang — bukan karena rancangan memutuskan begitu, melainkan karena satu parameter tidak terisi. Ini menguatkan catatan BLUEPRINT dengan sebab yang sebelumnya tidak terbaca.

## 6.2 Kedua penulisan `IsEditClaim = 0`

Keduanya adalah **penugasan terakhir** dalam blok pembuatan baris — 7322 menutup blok Retensi Cedant, 9795 menutup blok layer. Yang mendahuluinya adalah seluruh nilai uang baris itu.

**DERIVED**: `IsEditClaim` bukan penjaga yang dilanggar; ia **medan inisialisasi**. Setiap baris baru lahir dengan `0`. `EditXOLAlokasi` kemudian menulis `1`, dan tidak ada yang pernah membacanya. Mekanisme perlindungan suntingan itu **tidak pernah lengkap**, bukan rusak belakangan.

## 6.3 `Local.URLimit` dan cabang `Deductible2`

`Local.URLimit` dihitung di **dua** tempat, 6038 dan 6065, keduanya dengan struktur bertingkat yang sama:

```
@if(.Currency=="IDR", (@divide(.Deductible, Local.Kurs, 10) * Local.ProrateClaim/100),
  @if(Local.Currency==.Currency, (.Deductible * Local.ProrateClaim/100),
    (@toDecimal(@replaceAll(.Deductible2, ",", ".")) * Local.Kurs * Local.ProrateClaim/100)))
```

Cabang ketiga memanggil `@replaceAll(.Deductible2, ",", ".")` — **perbaikan koma-jadi-titik atas medan teks**, lalu `@toDecimal`. Nilai uang disimpan sebagai teks berformat lokal, dan diperbaiki di titik pakai.

Catatan penulisnya sendiri ada di baris 12 berkas itu: `<pyMemo>Pastiin format data .Deductible2</pyMemo>`.

`Local.UR` dihitung di 6103 dan 6125 memakai `Local.URLimit` sebagai ambang. Di dalam 6125 terdapat:

```
(.Deductible * 0 * Local.ProrateClaim/100)
```

**Dikalikan nol.** Cabang itu selalu menghasilkan nol apa pun isinya. Skala yang dinyatakan di rantai ini: `@divide(..., 10)` — **10**, berbeda dari 20 yang dipakai rantai `AdjusterFeeValue`, di berkas yang sama.

---

# 7. S14 — apakah hanya `UR`

## Jawabannya: ya, dan karena itu ia satu kejadian, bukan pola

Dua belas pola perakitan teks ditemukan pada empat lapisan:

| Rakitan | Diuji utuh? |
|---|---|
| `"Previously Calculated "+.TreatyName` | **ya, 2 kali** |
| `" CLAIM HISTORY MASTER ID " + InputSpreading.CARI22` | tidak |
| `"( ROE 1 "+Local.Currency+" = IDR "+Local.Kurs+")"` | tidak |
| `"(ROE 1 "+Local.Currency+" = IDR "+Local.ROE+")"` | tidak |
| `"Admin Claim " + .PICSuggest` | tidak |
| `"CFS_"+@replaceAll(...InsuredName,",","")+".pdf"` | tidak |
| `"CTS-"+@pxReplaceAllViaRegex(@DateTime.CurrentDateTime(),"[GMT.]","")` | tidak |
| `"Ceding Share ("+Local.ShareCedant+"%)"` | tidak |
| `"Claim Analysis " + @replaceAll(PrintCA.CARI2,",","") + ".pdf"` | tidak |
| `"PLA " + TempDataPLA.CARI2 + ".pdf"` | tidak |
| `"Pengajuan "+local.SubSubject+" Klaim : "+…` | tidak |
| `"Pengajuan Akseptasi : "+…` | tidak |

Sebelas sisanya adalah teks tampilan: nama berkas PDF, subjek surel, label persentase. Tidak satu pun dibandingkan dengan `==`.

**Saya menurunkan klaim saya sendiri.** Di laporan S3–S9 saya menulis bahwa "nama dan keadaan bercampur di satu kolom". Itu tetap benar sebagai deskripsi kejadian ini, tetapi ia **satu kejadian**, bukan pola yang berulang. Naik jadi pola hanya bila ditemukan yang kedua; sampai sekarang tidak ada.

Sapuan tambahan: **tidak ada nilai `TreatyName` lain yang pernah dirakit.** Sapuan `<PropertiesName>…TreatyName</PropertiesName>` atas seluruh ekspor mengembalikan nol penugasan berbentuk perakitan selain yang ini.

## Kedua cabang penjaga: ternyata identik

```
pyStepsPreCondParamsWhen      : Local.TreatyName=="Previously Calculated UR"
pyStepsPreCondParamsWhenTrue  : 3
pyStepsPreCondParamsWhenFalse : 2
```

Keduanya — baris 7223 dan 7425 — **sama persis**, sampai ke kode transisinya. Jadi ia bukan dua cabang berbeda; ia **satu uji yang digandakan** di dua langkah.

**Arti kode transisi `2` dan `3` TIDAK DITEMUKAN.** Keduanya kode bawaan Pega, dan pemetaannya tidak ada di dalam ekspor. Perilaku yang dibawa nama itu karena itu **belum terbaca** — yang terbaca baru keberadaannya dan tempat ujinya. Menebak artinya akan menjadi ADR di atas dugaan.

---

# 8. S15 — sasaran pengganti `_HITUNG`/`_SUNTING`

**Pertanyaan dibalik menghasilkan jawaban. Bukan nol.**

Sapuan: setiap `Property-Set` yang sasarannya memuat `SpreadingRisk`, di seluruh 279 berkas.

| Rule | Penugasan | Bernilai uang |
|---|---|---|
| `CountLossAllocation_act` — **mesin alokasi** | 43 | 29 |
| `AdjClaimCNP_Act` | 18 | **18** |
| `AdjClaimAmount_Act` | 9 | **9** |
| `SetActualPremium_ACT` | 6 | **6** |
| `AddListClaimNP_Act` | 4 | 0 — hanya `px…` |
| `AddAkseptasiCNP_Act`, `CreateChildKomiteCloseNP_Act`, `GeneratePlaCNP_Act` | 1 masing-masing | salinan seluruh PageList |
| `GenerateCFS_act` | 4 | 1 |

**Tiga rule menulis nilai uang pada `SpreadingRisk` di luar mesin alokasi.** Kolom yang ditulis **oleh mesin dan oleh salah satu dari ketiganya** adalah kandidat pasangan hitung/sunting:

| Kolom | Ditulis mesin | Ditulis juga oleh |
|---|---|---|
| `ClaimEstimation` | ya | `AdjClaimCNP_Act`, `SetActualPremium_ACT` |
| `ClaimAmountAdjust` | ya | `AdjClaimCNP_Act`, `AdjClaimAmount_Act`, `SetActualPremium_ACT` |
| `ClaimSpreaded` | ya | ketiganya |
| `TotalClaim` | ya | `AdjClaimCNP_Act`, `AdjClaimAmount_Act`, `SetActualPremium_ACT` |
| `AdjClaimValue` | ya | `AdjClaimCNP_Act`, `AdjClaimAmount_Act` |
| `AdjusterFee` | ya | `AdjClaimCNP_Act` |
| `Salvage` | ya | `AdjClaimCNP_Act` |
| `CNPOthersFee` | ya | `AdjClaimCNP_Act`, `AdjClaimAmount_Act` |
| `ClaimAmountIDR` | ya | `SetActualPremium_ACT` |
| `TotalSpread` | ya | `SetActualPremium_ACT` |

**Sepuluh kolom ditulis dua kali oleh rule berbeda.** Itu daftar kandidat yang dicari tiket `17`, dan ia datang dari bukti, bukan dari perkiraan.

## Tetapi jawabannya tidak sepenuhnya melegakan

Ketiga rule itu **bukan penyuntingan manual petugas**. Namanya menunjukkan perhitungan: *AdjClaim*, *SetActualPremium*. Yang terbaca adalah **hitung ulang oleh rule lain**, bukan nilai yang diketik orang.

Jadi pertanyaan aslinya terbelah dua:

1. **Nilai mana yang dihitung lebih dari satu rule** — **terjawab**, sepuluh kolom di atas. ADR-0008 tetap berlaku untuk ini: tanpa pemisahan, rule yang jalan belakangan menimpa yang duluan, dan `IsEditClaim` terbukti tidak mencegahnya.
2. **Apakah ada jalur suntingan manual sama sekali** — **TIDAK DITEMUKAN** pada empat lapisan. `EditXOLAlokasi` hanya gerbang kata sandi; tidak ada Section maupun FlowAction yang mengikat medan masukan ke properti uang `SpreadingRisk`.

Bila pembacaan kedua bertahan setelah folder Komite dibuka, akibatnya berat dan harus dinyatakan apa adanya: **yang dilindungi ADR-0008 adalah sesuatu yang belum pernah ada jalurnya**, dan `IsEditClaim` adalah bekas rancangan yang tidak pernah selesai — konsisten dengan temuan bahwa ia tidak pernah dibaca.

Saya tidak menutup butir 2. Ia menunggu folder Komite.

## Temuan sampingan yang penting: nol ditulis ke tiga besaran sekaligus

`AdjClaimCNP_Act` baris 2578, 2625, 2647, 2669 menulis `ClaimSpreaded <- 0`, `AdjusterFee <- 0`, `Salvage <- 0`, `CNPOthersFee <- 0` berurutan. Bahan langsung bagi tiket `33`.

---

# 9. S16 — properti bawaan Pega versus properti bisnis

**47 nama** berawalan Pega di empat lapisan: `px` 16, `py` 28, `pz` 3.

## 9.1 Verifikasi konvensi — dan konvensi itu tidak bertahan utuh

Pemahaman kerja yang diminta diverifikasi: `px` disetel sistem dan hanya-baca; `py` standar, boleh diisi aplikasi; `pz` internal Pega, bukan untuk dipakai langsung.

**Bukti yang bertentangan ada, dan cukup untuk mengubah pemahaman itu.**

### `px` **ditulis aplikasi**, berulang kali dan dengan pola tetap

Setiap kali sebuah baris ditambahkan ke PageList tertanam, empat properti `px` diisi:

```
ClaimData.SpreadingRisk(Local.lengthpage+1).pxObjClass     <- "ASM-FW-GISFW-Data-SpreadingRisk"
ClaimData.SpreadingRisk(<LAST>).pxCreateDateTime           <- @CurrentDateTime()
ClaimData.SpreadingRisk(<LAST>).pxCreateOperator           <- OperatorID.pyUserIdentifier
ClaimData.SpreadingRisk(<LAST>).pxCreateOpName             <- OperatorID.pyUserName
```

Pola itu berulang untuk `ListClaimAmount`, `EstimationList`, `SpreadingClaim`, `ListClaimAcceptation`, dan `AdjustmentList` — seluruhnya di `AddListClaimNP_Act` dan `AddAkseptasiCNP_Act`. Cacah: `pxCreateOperator` ditulis 7×, `pxCreateOpName` 6×, `pxCreateDateTime` 5×, `pxObjClass` 8×.

**"Hanya-baca" tidak berlaku untuk baris tertanam.** Pega mengisi `px` secara otomatis hanya pada **work object**; pada PageList tertanam, aplikasilah yang mengisinya — dan di sini aplikasi memang mengisinya.

### `pz` — konvensinya bertahan, dengan satu keterangan

`pzInsKey` tercatat "ditulis" 2×, tetapi keduanya **penyalinan pegangan, bukan pembuatan**:

| Berkas, baris | Penugasan |
|---|---|
| `CreateChildKomiteCNP_Act` 1702 | `Local.pzInsKey <- pyWorkPage.pzInsKey` |
| `ASMForceCaseClose` 1315 | `WorkObjectsToUpdate.pxResults(<APPEND>).pzInsKey <- .pzInsKey` |

Tidak satu pun menulis `pzInsKey` milik objek kerja itu sendiri. **Konvensi `pz` bertahan**; yang perlu dikoreksi hanya "bukan untuk dipakai langsung" — ia dipakai langsung sebagai pegangan, 31 pembacaan dan 5 penjaga.

## 9.2 Tiga golongan

### Golongan 1 — bawaan Pega, tidak dimigrasi

`pxDeadlineTime`, `pxLinkedRefFrom`, `pxLinkedRefTo`, `pxListSubscript`, `pxLockHandle`, `pxRefObjectInsName`, `pxReqContextURI`, `pxSystemNodeID`, `pxTaskLabel`, `pxUserIdentifier`, `pxResults`, `pyDesktopType`, `pyDesktopSubType`, `pyGridPaginator`, `pyIsMobile`, `pyPosition`, `pyTemplateButton`, `pyTemplateInputBox`, `pyTemplateGeneric`, `pyTemplateDisplayText`, `pzPegaDefaultGridIcons`, `pzProductionLevel`, `pyPDFHeaderHTMLTemplate`, `pyPDFFooterHTMLTemplate`, `pyPDFPageSize`, `pyPDFPageOrientation`, `pyExpanded`, `pyPageName`.

Mereka hilang bersama Pega. **Tetapi dua di antaranya dipakai sebagai penjaga** dan karenanya membawa keputusan — lihat 9.3.

### Golongan 2 — jembatan migrasi

`pzInsKey`, `pxCoverInsKey`, `pxCoveredInsKeys`. Masuk `MIGRASI_KORELASI`, berumur terbatas, tidak pernah masuk skema kanonik. Dasarnya §8.4a.

### Golongan 3 — **nilai bisnis di dalam properti bawaan Pega**

Ini yang diminta dicari, dan ia lebih besar dari `pyID`.

| Properti | Nilai bisnis yang disimpan | Bukti |
|---|---|---|
| `pxCreateOperator` | **pelaku** tiap baris tertanam | `AddListClaimNP_Act` 651, 1088, 1777, 2182 dan seterusnya |
| `pxCreateOpName` | **nama pelaku** tiap baris tertanam | 671, 1108, 1797, 2202 |
| `pxCreateDateTime` | **waktu pembuatan** tiap baris tertanam | 623, 1068, 1757, 2162 |
| `pyID` | **nomor klaim** `CLMNP-…` | 100 pembacaan, 40 penjaga |
| `pyWorkIDPrefix` | **pembeda lini bisnis** `"CLM-"`, `"CLMP-"`, `"CLMNP-"` | 65 pembacaan, **17 penjaga** |
| `pyNote` | **uraian sebab kerugian** dan **status aktif** | `CNMInsertCauseOfLoss_act` 735 ← `OutputData.COL_DESC`; `CNMInsertDetailCauseOfLoss_act` 733 ← `OutputData.STS_AKTIF` |
| `pyCategory` | **kategori dokumen** | `GetBase64Attachment` 1534 ← `.KATEGORI_1` |
| `pyMemo` | **nama berkas lampiran** | `GetBase64Attachment` 1576 ← `.NAMAFILE` |
| `pyNotifyAccountName` | **pembeda konvensional dan syariah** | `SendEmailKlaim` 3021 ← `"NUSARE"`, 3281 ← `"NUSARESYARIAH"` |
| `pyLabel` | **nama aksi** `"Close"`, `"Update"` | empat rule cause-of-loss |
| `pyReportClass` | **kelas tabel hilir** | 13 penugasan, termasuk `"ASM-FW-GCNMFW-Int-CLAIMREJECTED"`, `"…-EMAILKOMITE"` |

**Yang paling berisiko adalah tiga baris pertama.** Kolom pelaku dan waktu untuk setiap baris tertanam — alokasi, estimasi, akseptasi, adjustment — **tersimpan di properti `px`**. Bila golongan 1 dibuang mentah karena berawalan `px`, **setiap baris kehilangan pencatat dan waktunya**. Itu tepat data yang dirancang kolom pelaku di spec, dan sumbernya ada di tempat yang paling mudah disangka sampah sistem.

`pyNotifyAccountName` menyentuh **G1**, pertanyaan syariah yang masih terbuka: pembedaan konvensional/syariah beredar sebagai nilai properti notifikasi Pega, bukan sebagai medan bisnis. Dicatat; tidak ditarik kesimpulan.

## 9.3 Properti Pega yang dipakai sebagai **keputusan**

Lima belas properti Pega muncul sebagai penjaga langkah, syarat tampil, atau syarat hanya-baca:

| Properti | Cacah penjaga | Keputusan apa |
|---|---|---|
| `pyID` | **40** | percabangan per nomor klaim — termasuk tambalan per-case |
| `pyWorkIDPrefix` | **17** | **pemilihan lini bisnis** — `CLMP-` prop, `CLM-` biasa, `CLMNP-` non-proporsional |
| `pxResults` | 9 | keberadaan hasil pencarian |
| `pyNote` | 9 | — |
| `pxCreateOperator` | 8 | percabangan per orang |
| `pyIsMobile` | 8 | jenis perangkat |
| `pxReqContextURI` | 6 | konteks permintaan |
| `pzInsKey` | 5 | pegangan instance |
| `pxDeadlineTime` | 4 | tenggat SLA |
| `pxListSubscript`, `pxLockHandle`, `pyHasAttachments`, `pyLabel`, `pyPosition`, `pyStatusMessage` | 1–2 | — |

**`pyWorkIDPrefix` yang paling berat.** Tujuh belas kali, lini bisnis ditentukan dengan memeriksa **awalan nomor kasus Pega**. Contohnya sudah terbaca di S3:

```
(pyWorkPage.pyWorkIDPrefix=="CLMP-" && .Type=="3") ||
(pyWorkPage.pyWorkIDPrefix=="CLM-"  && .PaymentType=="3") ||
(pyWorkPage.pyWorkIDPrefix=="CLMNP-"&& .PaymentType=="3")
```

Sistem baru tidak punya `pyWorkIDPrefix`. **Penggantinya harus disediakan secara eksplisit** — sebuah pembeda lini bisnis yang berdiri sendiri, bukan disimpulkan dari bentuk nomor. Hal yang sama berlaku untuk `pxCreateOperator` sebagai penjaga (8 kali) dan `pyID` (40 kali, sebagian besar tambalan per-case yang ADR-0012 sudah bekukan).

---

# 10. Yang diparkir, tidak dikejar

| Sasaran | Kedudukan |
|---|---|
| Nilai rupiah biaya penilaian dan salvage | ~~**`ASK-AKUNTANSI` butir 7**; dua pembacaan dicatat sebagai **cabang terbuka di badan usulan ADR-0029**; tidak ada rancangan yang mengandaikan salah satunya~~ — **DICABUT 19 September 2026.** Cabang A dan B dihapuskan; **ADR-0029 `accepted` dan berlaku seragam** untuk seluruh nilai uang, termasuk kesembilan nama itu. Butir 7 turun menjadi **verifikasi** dan tidak menahan apa pun. Bukti di bagian ini tidak berubah; yang berubah akibatnya |
| Sufiks `EDM` | tercatat, tidak dikejar |
| Penulis `AlokasiXOLPaid` | tercatat, tidak dikejar |
| `excludeXML/GetBase64Attachment.xml` | folder tertutup |
| Arti kode transisi Pega `2`/`3` | **TIDAK DITEMUKAN** pada empat lapisan |
| Apakah `ClaimPercentage` dan `RNMShare` selalu sama | **TIDAK DITEMUKAN** |
| Penulis 26 penanda golongan 3b | menunggu folder Komite |
| `Struktur_Flow_TreatyIn.xlsx` | sumber baru yang ditemukan, belum dibaca |

---

# 11. Kode D baru

| Kode | Temuan | Tiket yang disentuh |
|---|---|---|
| **D16** | Stempel 273 berkas salah; ekspor memuat **279** dan seluruhnya terbaca. Sebabnya saringan `-not -path "*Komite*"` yang ikut menyaring enam berkas ber-nama Komite di folder terbuka. EVIDENCED-NIHIL dan ADR-0006 **tidak bersyarat** | seluruh laporan sapuan |
| **D17** | `pzInsKey`, `pxCoverInsKey`, `IndexObject`, `CASEID` — keempatnya pengenal sistem lama; hanya `pzInsKey` dan `pxCoverInsKey` yang berguna, dan hanya di tabel korelasi | `15`, `23`, `34` |
| **D18** | `DDL_Script_ClaimNonProp2.xls` **himpunan bagian murni** dari berkas pertama; nol objek baru. Klaim "sumber baru" dicabut | — |
| **D19** | REQ-017 **tertutup nol dari 14**. `CURRENCYSTANDARD` adalah view yang membaca `m_CurrencyStandard`; keduanya objek berbeda | REQ-017 |
| **D20** | **Database link ke instance lain sudah dipakai produksi** — `hrdasm.v_hrd_mst@asmd.sinarmas.co.id` di `V_MST_USER_TEKNIS`. ADR-0016 kini melarang sesuatu yang sudah berjalan | ADR-0016, ADR-0028, REQ-021 |
| **D21** | Satu besaran, tiga nama: `ClaimSpreaded` → `GrossValue` → `InputParamOs.GrossValue`. Dua jalur menghitungnya berbeda — `.ClaimSpreaded` lawan `.ClaimEstimation + .AdjClaimValue` | `29`, `31` |
| **D22** | Pengurangan selisih terjadi **ke dua arah** dalam satu berkas yang sama, satu berpenjaga dan satu tanpa penjaga | `29`, `38` |
| **D23** | `ClaimAmountIDR` diisi dua rule: satu menguji mata uang dan mengalikan kurs, satu **menyalin apa adanya**. Kolom IDR dapat terisi tanpa kurs pernah dipakai | ADR-0029, `16`, `33` |
| **D24** | Skala mengikuti modul terbukti di satu berkas: rantai `AdjusterFeeValue` menyatakan **20** di setiap pembagian; rantai `GrossValue` **tidak menyatakan skala sama sekali**; `Local.URLimit` memakai **10** | ADR-0003, `16` |
| **D25** | Penanda yatim bukan dua melainkan **banyak**: 14 ditulis tanpa dibaca, 26 dibaca tanpa ditulis pada objek kerja sendiri, 13 lagi milik sistem lain | `17`, `31` |
| **D26** | **22 properti** dibandingkan sebagai angka **dan** teks; **10** punya nilai kosong tanpa padanan numerik; `IsPLA` punya nilai `1A` yang bukan angka | seluruh tiket berdomain |
| **D27** | Daftar alokasi menumpuk **karena langkah 10 `Property-Remove` tidak pernah diisi parameternya**, bukan karena rancangan | `17`, `19` |
| **D28** | `IsEditClaim` adalah **medan inisialisasi**, ditulis `0` sebagai penugasan terakhir tiap blok pembuatan baris. Perlindungan suntingan tidak pernah lengkap sejak awal | `17` |
| **D29** | `"Previously Calculated UR"` **satu-satunya** rakitan yang diuji utuh dari 12 rakitan. Klaim "pola" diturunkan jadi "satu kejadian". Kedua penjaganya identik; arti kode transisinya TIDAK DITEMUKAN | `17`, REQ-033 |
| **D30** | Sepuluh kolom `SpreadingRisk` ditulis mesin alokasi **dan** salah satu dari `AdjClaimCNP_Act`, `AdjClaimAmount_Act`, `SetActualPremium_ACT`. Tetapi jalur **suntingan manual** TIDAK DITEMUKAN | `17` |
| **D31** | Konvensi awalan Pega **tidak bertahan utuh**: `px` ditulis aplikasi pada setiap baris PageList tertanam. `pz` bertahan | `12`, `17`, `19` |
| **D32** | **Pelaku dan waktu setiap baris tertanam tersimpan di `pxCreateOperator`, `pxCreateOpName`, `pxCreateDateTime`.** Membuang `px` mentah berarti kehilangan seluruh kolom pelaku | `12`, `17`, `19`, `20` |
| **D33** | `pyWorkIDPrefix` memutuskan **lini bisnis** 17 kali. Sistem baru harus menyediakan pembeda yang berdiri sendiri, bukan disimpulkan dari bentuk nomor | `12`, `29` |
| **D34** | `pyNotifyAccountName` membawa pembedaan `"NUSARE"`/`"NUSARESYARIAH"` — pertanyaan syariah G1 beredar di properti notifikasi Pega | G1, `12` |
| **D35** | Dari empat nilai `CNPStatusCase` yang didaftar memori §7.3, **dua tidak ada di mana pun** di 279 berkas. `CNPStatusCase` juga tidak pernah dibaca | `06`, `31` |
| **D36** | Sisi kasir memakai **25 medan bernama urut** `TempKasir.CARI1`…`CARI24` dan `CARIDATETIME`, belum pernah diinventarisasi | `05`, `25`, `29` |
