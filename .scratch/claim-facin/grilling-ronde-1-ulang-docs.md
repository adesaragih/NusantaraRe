# Grilling Ronde 1 ULANG — Claim Fac In diadu dengan DOKUMEN

**Tanggal:** 2026-09-19 · **Korpus:** `D:\XML\RNM_BRD\Claim Fac In\` (READ-ONLY, **482 berkas**)
**Dokumen:** 15 ADR · `ATURAN-BACA-KORPUS-PEGA.md` · `CONTEXT.md` · `docs/agents/domain.md`

> ⭐ **Ini ronde ULANG, bukan ronde 2.** Ronde 1 mengukur struktur. Ronde ini mengukur apakah
> struktur itu **bertabrakan dengan dokumen yang sudah disahkan**.
> ⛔ `grilling-ronde-1.md` **DIPERTAHANKAN apa adanya**. Seluruh ralat ditulis di sini dengan
> **mengutip angka lamanya** — nol angka lama dihapus.
>
> ⛔ Berkas lama NOL disunting · modul lain NOL · kode NOL · `CREATE TABLE` NOL · DDL NOL ·
> nilai rahasia NOL · nol nomor baris XML · nol butir `[terbuka]` ditutup ·
> ⛔ **nol usulan revisi ADR** · ⛔ **dua pertanyaan §G tidak dijawab**.
>
> ⚠️ **Tabrakan dilaporkan:** blok ini datang saat `/grill-with-docs` baru dipanggil. Blok
> GANTI KONTEKS **menang**; skill itu dibatalkan.

---

## §A — Tiga dugaan asisten: **ketiganya DIKUATKAN**

### A1 — `pxInsName` saja **bukan** identitas lengkap ⭐ **DIKUATKAN**

**Jendela (S1):** seluruh **482** berkas Claim Fac In + **329** Claim Prop; medan yang dibaca
`pxInsName` · `pyRuleSet` · `pyRuleSetVersion` · `pxUpdateDateTime`. Nol berkas dikecualikan.

| Cara hitung | Fac In unik | Prop unik | BERSAMA |
| --- | --- | --- | --- |
| **LAMA** — `pxInsName` saja *(ronde 1)* | 470 | 317 | **162** |
| ⭐ **BARU** — empat bagian | **478** | **324** | ⭐ **161** |

**Rincian baru:**

| | Ronde 1 *(lama)* | ⭐ Ronde ulang *(baru)* |
| --- | --- | --- |
| bersama | **162** | **161** |
| `pxUpdateDateTime` identik | **158** | ⭐ **159** |
| "beda versi" | **4** | ⭐ **2** — dan artinya berubah: **beda `pxUpdateDateTime` pada identitas empat bagian yang SAMA** |
| nama sama tetapi **ruleset/versi BEDA** | *(tidak terhitung)* | ⭐ **4 pasangan** |

**Keempat pasangan berbeda ruleset/versi:**

| Nama | Claim Fac In | Claim Prop |
| --- | --- | --- |
| `ASM-FW-GISFW-INT-CURRENCY!BROWSECURRENCY_RD` | **GISFW** 01-01-56 | **GCNMFW** 01-01-07 |
| `ASM-FW-GISFW-INT-CURRENCYSTANDARD!ASM!CURRENCYSTANDARD` | **GISFW** 01-01-64 | **GCNMFW** 01-01-06 |
| `DATA-PORTAL!MSTADJUSTERCONSULTANT` *(Harness)* | GISFW 01-01-56 | GISFW **01-01-83** |
| `DATA-PORTAL!MSTADJUSTERCONSULTANT` *(Section)* | GISFW 01-01-83 | GISFW **01-01-56** |

> ⭐ **Dua rule ber-`pxInsName` sama ternyata hidup di RULESET BERBEDA.** `BROWSECURRENCY_RD` dan
> `CURRENCYSTANDARD` bukan "rule yang sama versi berbeda" — mereka **dua rule berbeda yang kebetulan
> sama namanya**, satu di `GISFW`, satu di `GCNMFW`.
>
> ⚠️ ⭐ `MSTADJUSTERCONSULTANT` **belum pernah terlihat di ronde 1** — dan ia **bersilang**: versi
> Harness di Fac In cocok dengan versi Section di Prop, dan sebaliknya.

**Sebaran ruleset** *(penyebut 482 / 329)*:

| Ruleset | Fac In | Prop |
| --- | --- | --- |
| `GCNMFW` | 350 | 230 |
| `GISFW` | 113 | 81 |
| `GCNMFWInt` | 5 | 6 |
| `SFAGIS` | 5 | 5 |
| `ADESAMUEL@` | ⚠️ **5** | 2 |

⚠️ ⭐ **`ADESAMUEL@` adalah ruleset bernama akun orang** — 5 rule di Fac In, 2 di Prop. Dicatat
sebagai fakta; ⛔ nol kesimpulan.

### A2 — Korpus campuran ekspor dari beberapa sistem ⭐ **DIKUATKAN, dan LEBIH BESAR dari jendelanya**

**Jendela (S1):** seluruh **482** berkas, medan `pxUpdateSystemID`. **Terisi 482 · kosong 0.**

| Nilai | Fac In | % dari 482 | Claim Prop *(penyebut 329)* |
| --- | --- | --- | --- |
| `pega` | **320** | 66,4% | 222 |
| `pegadevnusare2` | **146** | 30,3% | 103 |
| ⭐ `pegaprdnusare` | **16** | 3,3% | 4 |

> ## ⭐ TEMUAN — dan ia menyentuh SELURUH modul yang sudah selesai
>
> **Bukan dua sistem, melainkan TIGA.** Dan polanya **sama persis di Claim Prop** — jadi ini
> **bukan ciri Claim Fac In**, melainkan **ciri seluruh korpus**.
>
> ⚠️ Artinya setiap berkas keluaran yang sudah jadi — spec Claim Prop, spec Komite Claim Prop,
> struktur tabel, relasi tabel, dan spec tujuh modul lain — disusun dari korpus yang
> **mencampur ekspor dari tiga sistem** tanpa pernah menyatakannya.
>
> ⛔ **Tidak saya tetapkan apa pun.** Ini fakta, dan akibatnya keputusan work owner.

### A3 — **Dua jalur halaman, bukan satu** ⭐ **DIKUATKAN — dan "38 dari 60" memang tidak bisa dipakai**

**Jendela (S1):** seluruh **482** berkas. *Menulis* dibaca dari `PropertiesName`; *membaca* dari
`PropertiesValue` · `pyConditionString` · `pyValue` · `pyCriteriaValue`.

| Jalur | Merujuk | ⭐ **MENULIS** | Membaca |
| --- | --- | --- | --- |
| **JALUR-1** `pyWorkPage.Quotation.*` | **59** dari 482 | ⭐ **2 berkas · 10 penugasan** | 34 |
| **JALUR-2** `pyWorkPage.OfferFacIn.QuotationData.*` | **28** dari 482 | ⭐ **4 berkas · 27 penugasan** | 19 |

**Penulis jalur-1:** `Activity/CopyNB_Act.xml` *(1)* · `Activity/SetDataSobCeding_Act.xml` *(9)*.
**Penulis jalur-2:** `Activity/CheckLimitSpreadingTreaty_Act.xml` *(14)* ·
`Activity/CheckLimit_Act1.xml` *(7)* · `Activity/GeneratePLATreaty_Act.xml` *(3)* ·
`Activity/SetDataSobCeding_Act.xml` *(3)*.

> ⭐ **KEDUA jalur punya penulis.** ⛔ Karena itu **tidak satu pun** rule `When` boleh disebut mati
> karena halamannya tidak ada. **Dugaan asisten dikuatkan sepenuhnya.**
>
> ⚠️ Angka ronde 1 *"38 dari 60 menguji `pyWorkPage.Quotation`"* **menggabung dua jalur**.
> Hitungan yang benar dari §D1: **34** menguji jalur-1, **4** menguji jalur-2.

---

## §B — Korpus diadu dengan 15 ADR

⚠️ **Fakta yang mengubah cara seluruh §B dibaca:** kelima belas ADR bersumber dari
**Claim Life** dan **Komite Claim Life** — tertulis di medan `sumber` masing-masing. **Nol ADR
menyebut Claim Fac In.** Karena itu "menentang" di bawah berarti: *korpus Fac In berperilaku
berbeda dari keputusan yang diambil untuk Life*, **bukan** *Fac In melanggar aturannya sendiri*.

⚠️ **ADR-0004 berstatus `superseded`**, digantikan ADR-0013.

### B1 · ADR-0003 — uang non-float

`[terverifikasi]` **1 208 titik penugasan bernuansa uang di 97 dari 179 Activity.**

⚠️ **Tambalan pemisah desimal ditemukan** — 5 titik disaring, **2 benar-benar menyentuh uang**:

| Rule | Isi |
| --- | --- |
| `Activity/CekPremiLunas_Act.xml` | ⭐ `@replaceAll(HasilPremi.pxResults(1).HASIL1, ",", ".")` — **nilai premi ditambal koma-ke-titik** |
| `RDBList/CariHistoryClaim_SQL.xml` | nilai diambil dari `DATA_JSON` CLOB dengan alias yang berbohong — lihat §E2 |

**Vonis: ⚠️ MENENTANG.** ADR-0003 menuntut uang tidak pernah lewat *float* dan tidak disimpan
sebagai teks; korpus ini **menambal pemisah desimal pada nilai premi**, yang hanya perlu bila
nilainya **teks**.

### B2 · ADR-0004 *(superseded)* · ADR-0013 — resolusi endpoint

`[terverifikasi]` Kedelapan `ConnectREST` diperiksa satu per satu. ⛔ **Nilai autentikasi tidak dicetak.**

| Rule | Sumber endpoint | Auth |
| --- | --- | --- |
| `HitDLAClaimFacin` | `SETTING` → `LinkService!LinkService` | false |
| `KonversiKlaimNonLife` | `SETTING` → `LinkService!LinkService` | false |
| `SendAcceptationToKasir` | `SETTING` → `LinkService!LinkService` | ⭐ **true** |
| `ServiceGoogle` | `SETTING` → `LinkService!LinkService` | false |
| `getPayAttachment` | `SETTING` → `LinkService!LinkService` | false |
| `getPaymentClaim` | `SETTING` → `LinkService!LinkService` | false |
| `getPremiumPaidOn` | `SETTING` → `LinkService!LinkService` | false |
| ⚠️ **`getPremiumPaidOnMarine`** | ⛔ **`URL` langsung** | false |

⭐ **7 lewat `M_LINK_SERVICE` · 1 URL langsung.** ⚠️ Ketujuh yang lewat setting **tetap menyimpan
`pyURLEntered` terisi** sebagai sisa — yang menentukan adalah `pyBaseURLSelectionType`.

**Vonis ADR-0013: ⚠️ MENENTANG sebagian** — satu dari delapan tidak lewat tabel tautan layanan.
**Vonis ADR-0004: tidak menyentuh** — ia sudah `superseded`.

### B3 · ADR-0005 — `IsPEGAPROD` menjadi flag lingkungan

`[terverifikasi]` `@BASECLASS!ISPEGAPROD` · ruleset **GISFW 01-01-91** · **10 perujuk**, di
antaranya `HitServiceToKasir_Act` · `InsertGoogleStorage_Act` · `KonversiKlaim_Act` ·
`SendCloseClaimToKomite` · `ChooseDla_Act` · `ProtectDownloadFaceClaim`.

**Kondisinya:** `pxProcess.pzProductionLevel = "5"`.

> ⭐ **Diadu dengan A2 — dan hasilnya TIDAK bertemu.** `IsPEGAPROD` menguji **tingkat produksi
> proses berjalan**, sedangkan `pxUpdateSystemID` mencatat **sistem tempat rule terakhir disimpan**.
> Keduanya **bukan pengukur yang sama**: yang satu perilaku saat berjalan, yang satu jejak
> penyuntingan. ⛔ **Campuran tiga sistem di A2 tidak dapat dijelaskan oleh `IsPEGAPROD`.**

**Vonis: ⚠️ MENENTANG.** ADR-0005 memutuskan rule itu **tidak ditiru** dan diganti flag lingkungan
eksplisit; di Fac In ia **hidup dan menggerbangi 10 rule**, sebagian besar efek keluar.

### B4 · ADR-0006 — penomoran via stored procedure

`[terverifikasi]` **9 rule `RDBList` membawa `COMMIT` tertanam.** Dua di antaranya generator nomor:

| Rule | Generator? | Objek Oracle |
| --- | --- | --- |
| ⭐ `GetSequenceNumber_SQL` | **YA** | `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` |
| ⭐ `InsertClaimPNC` | **YA** | `POOLDATA.PEGA_JSON_KLAIM_PNC` |
| `GetTokenStorage_SQL` | bukan | `pooldata.GET_TOKEN_STORAGE` |
| `InsertLOGDirectKasir_SQL` | bukan | `POOLDATA.DIRECTTOKASIR_LOG` |
| `InsertProgressClaim_SQL` | bukan | `POOLDATA.PEGA_PROGRESSCLAIM` |
| `InsertSUBProgressClaim_SQL` | bukan | `POOLDATA.PEGA_SUBPROGRESSCLAIM` |
| `Insert_T_Storage_SQL` | bukan | `t_storage_image` |
| `SaveOSClaim_SQL` | bukan | `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM` |
| `UpdateDCauseOfLoss` | bukan | blok `DECLARE … BEGIN` |

**Vonis: ✅ MENGUATKAN.** ADR-0006 menyebut `PROC_GENERATE_SEQUENCE_NUMBER`, dan **procedure yang
sama** dipakai di modul ini.

### B5 · ADR-0007 jejak audit · ADR-0011 unit status = baris `AdjustmentList`

`[terverifikasi]` Modul ini **punya** kelas `ASM-FW-GCNMFW-DATA-ADJUSTMENT` *(43 rule)* dan menulis
`ComiteeClaim(<APPEND>)` per baris — **bentuk yang sama** dengan Life dan Prop.

⚠️ **Tetapi unit strukturalnya berbeda:** modul ini punya **dua tingkat tambahan yang tidak ada di
Life** — `DATA-OBJECT` *(50 rule)* dan `DATA-OBJECTITEM` *(56 rule)*. `CreateKMTNo_Act` langkah
**6.9** membawa **empat** penunjuk posisional ke kasus komite: `IDObject` · `CoverageSubscript` ·
`IndexObjectItem` · `IndexObject`.

**Vonis ADR-0011: ⚠️ MENENTANG sebagian** — unit statusnya bukan hanya baris `AdjustmentList`;
ada tingkat objek dan item objek di bawahnya.
**Vonis ADR-0007: ✅ MENGUATKAN** — jejak kronologi ditulis (`SetchronologyKlaimFacIn`), tetapi
kelengkapannya belum diuji.

### B6 · ADR-0008 efek keluar asinkron · ADR-0015 outbox wajib berhasil

`[terverifikasi]` **Sensus efek keluar** *(jendela S1: 179 Activity + 63 RDBList + 8 ConnectREST)*:

| Jalur | Langkah | Berkas |
| --- | --- | --- |
| basis data — `RDB-List` | **118** | 69 |
| basis data — `Obj-Save` / `Commit` | **31** | *(termasuk di atas)* |
| layanan luar — `Connect-REST` | **11** | 9 |
| email | **7** | 7 |
| Google Storage | — | 5 |
| Kasir | — | 9 |
| dokumen *(PDF/DLA/PLA/FaceSheet)* | — | 6 |

⭐ **Dari 52 langkah `Obj-*`/`Commit` yang ronde 1 hitung, hanya 2 punya jalur kegagalan.**

**Vonis ADR-0015: ⚠️ MENENTANG** — nol jejak outbox, nol antre-ulang, dan **nol** penanganan
kegagalan pada hampir seluruh efek keluar.
**Vonis ADR-0008: ⚠️ MENENTANG sebagian** — kegagalan **tidak dicatat**, sedangkan ADR-0008
menuntut *dicatat, tidak memblokir*.

### B7 · ADR-0010 — penyimpanan berkas tetap di Google Storage

`[terverifikasi]` ⭐ **Keduanya beda LINGKUNGAN, bukan beda versi** — ruleset dan versi **identik**:

#### `GENERATEIMAGEID_SQL` — `GISFW 01-01-88` di kedua salinan

| | Claim Fac In | Claim Prop |
| --- | --- | --- |
| `pxUpdateSystemID` | ⭐ **`pegaprdnusare`** | `pegadevnusare2` |
| `pxUpdateDateTime` | 2025-09-10 14:53:38 | 2025-09-04 14:39:33 |
| **isi SQL** | `'ASMPP' \|\| TO_CHAR(SYSTIMESTAMP,'YYYYMMDDHH24MISSFF9') \|\| SYS_GUID()` | `'ASMPP' \|\| TO_CHAR(SYSTIMESTAMP,'DD/MM/YYYY HH24:MI:SS.FF3')` |

⭐ **Bedanya, langkah demi langkah:** salinan Fac In **(a)** menaikkan ketelitian cap waktu dari
**milidetik ke nanodetik**, dan **(b)** menambahkan **`SYS_GUID()`**. Keduanya menghilangkan
kemungkinan dua berkas yang diunggah pada saat yang sama mendapat pengenal yang sama.

#### `GETMIMETYPE` — `GISFW 01-01-88` di kedua salinan

| | Claim Fac In | Claim Prop |
| --- | --- | --- |
| `pxUpdateSystemID` | ⭐ **`pegaprdnusare`** | `pegadevnusare2` |
| `pxUpdateDateTime` | 2026-02-19 01:50:45 | 2025-09-04 15:16:37 |
| jenis berkas | **19** | **17** |
| ⭐ **hanya di Fac In** | **`avi` · `mp4`** | — |

⭐ **Salinan `pegaprdnusare` mengenali dua jenis berkas video yang salinan `pegadevnusare2` tidak.**

⛔ **Itu data, bukan keputusan** *(§G3)*. ⛔ Apakah `avi` ikut ditiru **tidak saya jawab** — §G2.

**Vonis ADR-0010: ✅ MENGUATKAN** — berkas memang disimpan di Google Storage
*(`ServiceGoogle` · `InsertGoogleStorage_Act` · `t_storage_image`)*.

### B8 · ADR-0012 wewenang kirim komite · ADR-0014 pemutus per `KomiteID`

`[terverifikasi]` Keempat rule dibaca:

| Rule (`pxInsName`) | Langkah | Yang dikerjakan |
| --- | --- | --- |
| `…DATA-ADJUSTMENT!SETKOMITELIST_ACT` | 8 | nol penugasan bernuansa komite di medan yang disisir |
| `…DATA-ADJUSTMENT!SETLISTKOMITE_ACT` | 14 | ⭐ **membangun roster**: `ComiteeClaim(<APPEND>).KomiteID := .OPERATOR_ID` · `.KomiteAproval := 0` · `.KomiteEmail := .EMAIL` · `.IDKomite := .JABATAN` · langkah 7 `.TotalKomite := @SizeOfPropertyList(...)` |
| `…DATA-ADJUSTMENT!VIEWKOMITE_ACT` | 4 | ⚠️ langkah 3 memakai `@substring(pyWorkPage.pxCoveredInsKeys…)` — **irisan posisi karakter atas kunci internal Pega**, pola yang sama dengan Claim Prop |
| `…DATA-ADJUSTMENT!CREATEKMTNO_ACT` | **42** | ⭐ pembuat kasus komite: `.IsKomite := 1` · `childPageKomite.TransferType := 2` · **empat penunjuk posisional** |

⭐ **Bentuk rosternya identik dengan Life dan Prop** — `KomiteID` = akun operator, `IDKomite` =
jabatan, `KomiteAproval` mulai `0`.

**Vonis ADR-0014: ✅ MENGUATKAN** — roster per `KomiteID` memang bentuknya.
**Vonis ADR-0012: ⚠️ tidak dapat diputuskan** — gerbang wewenang berdasarkan `Type` **tidak
ditemukan** di medan yang disisir pada keempat rule itu. ⛔ Bukan "tidak ada"; **belum disisir
sampai ke layar**. `[terbuka]`.

### B9 · ADR-0001 · ADR-0002 · ADR-0009

| ADR | Vonis |
| --- | --- |
| **0001** batas konteks Claim Life ↔ Komite Life | ⛔ **tidak menyentuh** — lingkupnya Life. ⚠️ Tetapi bentuknya **berulang di sini**: Fac In punya kasus anak komite sendiri lewat `CreateKMTNo_Act` |
| **0002** RBAC tiga peran Claim Life | ⛔ **tidak menyentuh** — ADR sendiri menyatakan model peran **dirancang, bukan dimigrasikan**, karena korpus Pega nol memuatnya |
| **0009** migrasi penuh data Claim Life | ⛔ **tidak menyentuh** — keputusan migrasi data Life |

### B10 — ⭐ REKAP 15 ADR

| ADR | Vonis | Bukti satu baris |
| --- | --- | --- |
| **0001** batas konteks | ⛔ tidak menyentuh | lingkupnya Life |
| **0002** RBAC | ⛔ tidak menyentuh | korpus nol memuat model peran |
| **0003** uang non-float | ⚠️ **MENENTANG** | `CekPremiLunas_Act` menambal koma-ke-titik pada nilai premi |
| **0004** endpoint env-var | ⛔ tidak menyentuh | berstatus `superseded` |
| **0005** `IsPEGAPROD` | ⚠️ **MENENTANG** | rule hidup, **10 perujuk**, menggerbangi efek keluar |
| **0006** penomoran via procedure | ✅ menguatkan | `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` dipakai |
| **0007** jejak audit tiap transisi | ✅ menguatkan | `SetchronologyKlaimFacIn` menulis kronologi |
| **0008** efek keluar asinkron | ⚠️ **MENENTANG sebagian** | kegagalan **tidak dicatat** — 2 dari 52 punya jalur gagal |
| **0009** migrasi penuh Life | ⛔ tidak menyentuh | lingkupnya Life |
| **0010** berkas di Google Storage | ✅ menguatkan | `ServiceGoogle` + `t_storage_image` |
| **0011** unit status = baris `AdjustmentList` | ⚠️ **MENENTANG sebagian** | ada dua tingkat lagi: `DATA-OBJECT` 50 · `DATA-OBJECTITEM` 56 |
| **0012** wewenang kirim komite per `Type` | ⚠️ tidak dapat diputuskan | belum disisir sampai layar — `[terbuka]` |
| **0013** endpoint via `M_LINK_SERVICE` | ⚠️ **MENENTANG sebagian** | `getPremiumPaidOnMarine` memakai URL langsung |
| **0014** pemutus per `KomiteID` | ✅ menguatkan | `SetListKomite_act` membangun roster per `KomiteID` |
| **0015** outbox wajib berhasil | ⚠️ **MENENTANG** | nol outbox, nol antre-ulang |

⭐ **Menguatkan 5 · MENENTANG 6 · tidak menyentuh 3 · tidak dapat diputuskan 1.**

⛔ **Nol revisi ADR diusulkan.** Keenam yang ditentang dinaikkan sebagai pertanyaan di §H4.

---

## §C — Korpus diadu dengan aturan baca

### C1 — Tujuh langkah bernilai asing: ⭐ **TERJAWAB DARI KORPUS**

Dibaca dengan **jangkar `Embed-ActivitySteps` dan `Embed-ActivityPreConditions`**, bukan jendela
karakter.

| Rule | Langkah | flag | Metode | Baris syarat | Isi syarat | Kode | Ulang |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `CLaimFaceSheet_Act` | 39 | `0` | `Property-Set-HTML` | **1** | ⭐ **KOSONG** | `2/6` | — |
| `DLAFacintoTreaty_Act` | 13.4 | `0` | `Property-Set-HTML` | **1** | ⭐ **KOSONG** | `2/6` | — |
| `DraftGenerateDLAFacin_Act` | 13 | `0` | `Property-Set-HTML` | **1** | ⭐ **KOSONG** | `2/6` | — |
| `GenerateDLAFacin_Act` | 16 | `0` | `Property-Set-HTML` | **1** | ⭐ **KOSONG** | `2/6` | — |
| `GeneratePLATreaty_Act` | 18.23 | `0` | `Property-Set-HTML` | **1** | ⭐ **KOSONG** | `2/6` | — |
| `PrintPDFAccep_MultiAksep` | 23 | `0` | `Property-Set-HTML` | **1** | ⭐ **KOSONG** | `2/6` | — |
| `SendEmail_ACT` | 19 | *(kosong)* | `Property-Set` | **1** | ⭐ **KOSONG** | `2/2` | `PROPERTYLIST`, obj kosong |

> ### ⭐ PERTANYAAN 2 RONDE 1 TERJAWAB DARI KORPUS
>
> **Ketujuhnya punya tepat satu baris syarat, dan isi syaratnya KOSONG.** Karena tidak ada yang
> diuji, **apa pun arti `0`, tidak ada gerbang yang berlaku** — hasilnya sama dengan flag kosong.
>
> ⛔ **Yang masih `[terbuka]`** hanyalah `PROPERTYLIST`: gerbangnya nihil, tetapi **halaman yang
> diulang tidak terbaca** karena `pyStepsObjectName` kosong. ⛔ Tidak ditebak.
>
> ⚠️ Keenam langkah `flag='0'` seluruhnya `Property-Set-HTML` — **metode yang nol di Claim Prop**.

### C2 — Selisih keluarga-2: ⭐ **KETEMU**

`[terverifikasi]` **`SetCauseOfLossValue_act` langkah 4.4** punya **nol** baris
`pyStepsTransParams`. Ia satu-satunya dari 2 509 langkah. **2 508 + 1 = 2 509.** ✅

⚠️ Aturan **D1** menyatakan *"satu baris per langkah"* — ⭐ **aturan itu punya pengecualian**, dan
ini yang pertama ditemukan di proyek ini. Dicatat; ⛔ aturan tidak disunting.

### C3 · C4 — Usulan tambahan aturan baca — ⭐ **4 butir, semuanya BERTANDA USULAN**

⛔ **`ATURAN-BACA-KORPUS-PEGA.md` tidak disunting.**

> #### USULAN 1 — Aturan K, `Rule-Declare-Pages` (DataPage)
> *(usulan ronde 1, masih menggantung — diulang di sini supaya tidak hilang)*
> K1 kelas rule `CODE-PEGA-LIST` · K2 bentuk di `pyStructure` · K3 umur di `pyScope` ·
> K4 **sumber di `pyDeclarePagesDataSource`** dengan nama di `pyLoadActivity` / medan transform /
> `pyConnectorList` · K5 parameter di `pyParametersParamName` · K6 `pyRefreshStrategy` ·
> K7 jumlah sumber di `pySourceCount`.
> **usulan, menunggu pengesahan work owner**

> #### USULAN 2 — ⭐ Identitas rule **EMPAT BAGIAN**
> **Identitas rule = `pxInsName` (kelas + nama) + `pyRuleSet` + `pyRuleSetVersion`.**
> `pxInsName` saja **tidak cukup**: terbukti dua rule bernama sama hidup di ruleset berbeda
> *(`BROWSECURRENCY_RD`, `CURRENCYSTANDARD`)*.
>
> ⚠️ ⭐ **AKIBATNYA TEGAS:** setiap pembandingan lintas-modul yang sudah ditulis memakai
> `pxInsName` saja — termasuk **`RELASI-TABEL-KOMITE-CLAIM-PROP.md` §D2**,
> **`STRUKTUR-TABEL-CLAIM-PROP.md`**, dan **`grilling-ronde-1.md` §B4** — berdiri di atas
> identitas yang belum lengkap. Ini menyentuh **20 modul yang belum disensus**.
> **usulan, menunggu pengesahan work owner**

> #### USULAN 3 — `pxUpdateSystemID` wajib dicatat di tiap sensus
> Korpus memuat **tiga** nilai (`pega` · `pegadevnusare2` · `pegaprdnusare`) di **kedua** modul
> yang diperiksa. Setiap klaim sensus sebaiknya menyebut sebarannya, karena dua salinan rule
> beridentitas **empat bagian sama** bisa **berbeda isinya** *(terbukti: `GENERATEIMAGEID_SQL`)*.
> **usulan, menunggu pengesahan work owner**

> #### USULAN 4 — Pengecualian aturan D1
> Aturan **D1** menyatakan `pyStepsTransParams` **satu baris per langkah**. Ditemukan **satu
> pengecualian**: `SetCauseOfLossValue_act` langkah 4.4 bernilai **nol**. Kalimat aturannya
> sebaiknya berbunyi *"satu baris per langkah, dengan pengecualian yang wajib dilaporkan"*.
> **usulan, menunggu pengesahan work owner**

---

## §D — 60 rule `When` · 130 kandidat `pyMemo`

### D1 · D2 — Keenam puluh dibaca satu per satu ✅

**Jendela (S1):** 60 berkas `When`; **tujuh medan** disisir untuk kondisi — `pyConditionString` ·
`pyConditionValue1String` · `pyConditionValue1` · `pySimpleCondition` · `pyWhenExpression` ·
`pyExpression` · `pyLabel`.

| | |
| --- | --- |
| Punya kondisi **terbaca** | ⭐ **60 dari 60** |
| Tidak ketemu di 7 medan | **0** |
| Perujuk **NOL** | **2** — `IsEDM` · `IsEdmAdjShareCedant` |

⭐ **Aturan I2 terbayar:** di Claim Prop `isMaintenance` tampak tanpa kondisi di medan *viewer*.
Di sini **tidak satu pun** kosong sesudah tujuh medan disisir.

**Halaman yang diuji** *(penyebut 60)*:

| Halaman | Rule |
| --- | --- |
| ⭐ `pyWorkPage.Quotation.*` | **34** |
| ⭐ `pyWorkPage.OfferFacIn.*` | **4** |
| `pyWorkPage.pyWorkIDPrefix` / `pyWorkCover.pyWorkIDPrefix` | beberapa — `IsCLM` · `IsCLMP` · `IsCLMNP` · `IsClaim` |
| `pxProcess.pzProductionLevel` | 1 — `IsPEGAPROD` |
| `pxProcess.pxSystemNodeID` | 1 |
| `pyWorkPage.IsSpreadingUW` | 1 |

**Penggolongan** *(penyebut 60)*:

| Golongan | Jumlah | Dasar |
| --- | --- | --- |
| **HIDUP** — halaman punya penulis **dan** ada perujuk | ⭐ **58** | §A3: kedua jalur punya penulis |
| **Perujuk NOL** *(bukan berarti mati)* | **2** | `IsEDM` · `IsEdmAdjShareCedant` |
| **MATI karena halaman tidak ada** | ⭐ **0** | §A3 membuktikan kedua jalur ditulis |
| **MATI karena syarat selalu salah** | ⛔ **belum diuji** — menuntut nilai data, bukan struktur |
| **SISA IMPOR** | ⛔ **belum dapat diputuskan** |

> ### ⭐ SILANG DENGAN A3 — DAN INI MEMBALIK CLAIM PROP
>
> Di **Claim Prop**, 52 dari 61 rule `When` menguji halaman yang **tidak ada**, semuanya bernilai
> salah, dan diputuskan **tidak dimigrasikan**.
>
> Di **Claim Fac In**, **kedua jalur halaman PUNYA PENULIS** — jadi **nol** rule `When` boleh
> disebut mati karena alasan itu. ⭐ **Pertanyaan 3 ronde 1 terjawab dari korpus.**
>
> ⛔ Itu **tidak** berarti keenam puluhnya hidup — hanya berarti **alasan "halaman tidak ada"
> tidak berlaku di sini**.

### D3 — Nama kembar dengan Claim Prop: ⭐ **60 dari 60, identitas empat bagian SAMA**

`[terverifikasi]` Keenam puluh nama `When` Fac In **ada juga** di Claim Prop, dan **keenam puluhnya
cocok pada keempat bagian** — `pxInsName` · `pyRuleSet` · `pyRuleSetVersion`. **Beda: 0.**

⭐ **Seluruh himpunan `When` dipakai bersama kedua modul.** Ruleset-nya beragam — `GISFW`,
`GCNMFW`, `SFAGIS` — dan versinya pun *(01-01-52 sampai 01-01-91)*.

⚠️ **Itu membuat perbedaan vonis antara kedua modul menjadi tajam:** rule yang **sama persis**
dinyatakan **mati** di Claim Prop dan **tidak dapat dinyatakan mati** di Claim Fac In —
semata karena halaman yang diujinya **ditulis di sini** dan **tidak di sana**.

### D4 — 130 kandidat diuji: ⭐ **6 dari 12 pindah golongan**

**Kedua belas kandidat (d) diuji dengan membuka rule-nya:**

| Rule | Catatan | Hasil |
| --- | --- | --- |
| `GetCurencyCoverage_Act` | *"Hapus jika ada spreadingClaim yg sama"* | ⚠️ **MASIH ADA** — 2 titik, keduanya hidup |
| `InsertJsonClaimNonMBU_act` | *"hapus yg hapus attachment list"* | ✅ sudah dipatuhi — **pindah ke (b)** |
| `SaveAdjustmenttoDB_ACT` | *"hapus yg hapus attachment list"* | ✅ sudah dipatuhi — **pindah ke (b)** |
| `SendEmailDLA_ACT` | *"hapus when step 6"* | ✅ sudah dipatuhi — **pindah ke (b)** |
| `SetCatastrope_act` | *"hapus save"* | ✅ sudah dipatuhi — **pindah ke (b)** |
| `SetDefNonCatastrope_Act` | *"hapus save"* | ✅ sudah dipatuhi — **pindah ke (b)** |
| `SetSalvageValue` | *"hapus yg set cronologu"* | ✅ sudah dipatuhi — **pindah ke (b)** |
| ⚠️ `SetValueAdjusterFee` | *"hapus yg set chronologi"* | ⚠️ **MASIH ADA** — 1 titik, hidup |
| `InputAdjustment` *(FlowAction)* | *"hapus pre"* | ⚠️ pola masih muncul — belum dapat diputuskan |
| ⭐ `GetLimitPLADLA_Sql` | *"ganti biar jgn ambil dari DLA tpai dari treatyLimit"* | ⚠️ **pola masih muncul 18 kali** — menyentuh **batas nilai uang** |
| `BrowseAgentNusaRe_RD` | *"hapus param yang ga dipake"* | ⚠️ belum dapat diputuskan |
| `isPA_PNC` *(When)* | *"tambah pyWorkPage dan hapus policy"* | ⚠️ belum dapat diputuskan |

⭐ **6 dari 12 PINDAH GOLONGAN (d → b).** Golongan (d) yang bertahan: **6**, dua di antaranya
terbukti masih hidup jejaknya.

⚠️ **Di Claim Prop cara golong-dari-kata-kunci meleset lima kali lipat.** Di sini arah lesetnya
**berlawanan**: ronde 1 **melebih-lebihkan** golongan (d), bukan mengecilkannya.

### D5 — ⛔ Batas kejujuran, ditulis tegas

> **136 catatan golongan (a) + 191 catatan golongan (b) = 327 catatan TIDAK DIUJI di ronde ulang
> ini.** Dan **118 kandidat golongan (c) juga TIDAK diuji satu per satu** — hanya golongan (d)
> yang dibuka rule-nya.
>
> **Jadi dari 459 catatan, yang benar-benar diuji dengan membuka rule: 12.**

---

## §E — Tiga kekurangan yang ronde 1 tandai sendiri

### E1 — `pyXMLSignature`: ⚠️ **gagal dibaca, dan itu dilaporkan**

`[terverifikasi]` **178 dari 179 Activity** punya `pyXMLSignature` **terisi**.

⛔ **Parameter tidak berhasil diekstraksi** — dua pola dicoba (`name="…"` dan `<ParamName>`), setelah
*unescape* dua tingkat; keduanya menghasilkan **nol**.

⚠️ Sesuai aturan **I2**: kalimat yang sah adalah ***"tidak ketemu dengan dua pola yang dicoba"***,
**bukan** *"nol parameter"*. ⭐ **Ronde 1 menandai ini sebagai kekurangan; ronde ulang ini
mempersempitnya tetapi belum menutupnya.** `[terbuka]`.

### E2 — Alias SQL: ⭐ **38 dari 50 pasangan TIDAK COCOK**

**Jendela (S1):** 63 `RDBList`, seluruhnya punya `pyBrowseSQL`. Pasangan `kolom → "alias"`
diperiksa: **50**. Tidak cocok: **38**.

⚠️ **Sebagian "tidak cocok" itu palsu** — penyaring saya ikut menangkap kata kunci `AS` sebagai
nama kolom *(3 kasus: `GenerateImageID_SQL`, `CariHistoryClaim_SQL`, `GetIDConsultanAdj_SQL`)*.
Dilaporkan apa adanya, **tidak dibersihkan supaya angkanya rapi**.

**Yang paling menyesatkan** — ⭐ **rule yang sama dengan AC 106 Claim Prop**:

| Rule | Kolom asli | Alias | Kebohongannya |
| --- | --- | --- | --- |
| `CariHistoryClaim_SQL` | `DATA_JSON.DateOfLoss` | `START_DATE` | tanggal kejadian dinamai *tanggal mulai* |
| `CariHistoryClaim_SQL` | `IDPEGA` | `BRANCH_CODE` | id Pega dinamai *kode cabang* |
| `CariHistoryClaim_SQL` | `DATA_JSON.ClaimNo` | `BRANCH_NAME` | nomor klaim dinamai *nama cabang* |
| `CariHistoryClaim_SQL` | `DATA_JSON.CauseOfLoss` | `BUSINESS_CODE` | penyebab kerugian dinamai *kode bisnis* |
| ⭐ `BrowseHistoryClaim` | `DATA_JSON.DateOfLoss` | `BRANCH_CODE` | **rule KEDUA dengan pola yang sama** |
| `BrowseHistoryClaim` | `DATA_JSON.ClaimNo` | `BUSINESS_CODE` | idem |
| `BrowseRW_SQL` | `NOTE` | `TERITORYNAME` | catatan dinamai *nama wilayah* |
| `BrowseRW_SQL` | `DISTRICTNAME` | `NAMEDISTRICT` | dibalik susunannya |

⭐ **`BrowseHistoryClaim` belum pernah tercatat di modul mana pun** — ia kembaran
`CariHistoryClaim_SQL` dengan kebohongan yang sama. ⛔ Nol DDL ditulis.

### E3 — Gerbang aksi tombol: **tidak ketemu di 11 medan yang disisir**

**Medan yang disisir:** `pyActionConditions` · `pyActions` · `pyActionSetName` · `pyActionName` ·
`pyWhenName` · `pyClientValidation` · `pyPreActivity` · `pyPostActivity` · `pyActionString` ·
`pyEventName` · `pyBehaviorName`.

| Medan yang **ADA** | Jumlah |
| --- | --- |
| `pyClientValidation` | **50** |
| `pyActionName` | **33** — tepat satu per FlowAction, jadi **nama flow action**, bukan gerbang |

⚠️ Kalimat yang sah: ***"`pyActionConditions` tidak ketemu di 11 medan yang disisir"*** — **bukan
"nol"**. Claim Prop punya 17 baris di medan itu; ⛔ apakah modul ini memakai medan kedua belas
**belum terbukti**. `[terbuka]`.

---

## §F — `GetAllData_Act` dan 23 lompatan

### F1 — `GetAllData_Act`: **50 langkah**

`[terverifikasi]` `ASM-FW-GCNMFW-DATA-OBJECT!GETALLDATA_ACT` — ⚠️ **berkelas `DATA-OBJECT`**, yaitu
kelas terbesar ketiga modul ini.

⭐ **Label yang ada di dalamnya, seluruhnya:**

```
//   MBU   MBU2   MBU3   PA   PA2   PA3   TRAVEL   TRAVEL2   TRAVEL3
```

**Sembilan label lini produk, tepat sembilan lompatan.** Satu label `//` adalah penanda remark,
bukan sasaran.

### F2 — Ke-23 lompatan se-modul

| Rule | Langkah | Kode | Sasaran |
| --- | --- | --- | --- |
| `CLaimFaceSheet_Act` | 4 · 5 | `1/_` · `_/1` | `NoEx` |
| `CloseClaim` | 6 | `1/_` | `END` |
| `DeleteLocation_Act` | 2 | `2/1` | `SAVE` |
| **`GetAllData_Act`** | **11** | `1/_` · `1/2` ×2 | `MBU` · `TRAVEL` · `PA` |
| **`GetAllData_Act`** | **16.1.1.1.1** | `1/2` ×3 | `TRAVEL2` · `MBU2` · `PA2` |
| **`GetAllData_Act`** | **16.1.2.1.1** | `1/2` ×3 | `TRAVEL3` · `MBU3` · `PA3` |
| `GettsiAneka_Act` | 1.1.1.1.2.1 · 1.1.1.1.2.2 · 4 | `1/2` | `MBD` ×2 · `MBD1` |
| `SetNilaiResikoSendiri` | 18 | `_/1` | `TO` |
| `SetProtectionEstimation` | 12 · 13 · 19.1 · 21 · **26** | `1/_` · `1/3` | `TO` · `B` · `TO` · `Er` · ⚠️ **`END` menggantung** |
| `TravelDocument_act` | 3 | `_/1` | `TO` |

> ### ⭐ BENTUK PERCABANGAN — ditulis sebagai FAKTA KORPUS
>
> **Claim Fac In bercabang per lini produk dengan LOMPATAN, bukan dengan gerbang.** Empat lini
> muncul sebagai label: **MBU** *(kendaraan bermotor)* · **PA** *(personal accident)* ·
> **TRAVEL** · **MBD**. Tiga di antaranya berulang tiga kali dengan akhiran `2` dan `3` —
> tiga jalur berbeda yang masing-masing bercabang ke tiga lini yang sama.
>
> ⭐ **Claim Prop punya NOL kode 1 dan NOL kode 4.** Bentuk ini **tidak ada padanannya** di sana.
>
> ⛔ **Nol rancangan untuk Go ditulis.**

### F3 — Satu lompatan kode 4 *(Exit Iteration)*

⚠️ Ronde 1 menghitung **kode 4 = 1**. ⛔ Rule dan langkahnya **belum saya isolasi** di ronde ulang
ini — penyaring §C1 menghitung kemunculan kode, bukan menyimpan asalnya. `[terbuka]`, dan ini
**kekurangan ronde ulang ini**.

### F4 — `.ObjectItemList` dan `ObjectList`

`[terverifikasi]` **Jendela (S1):** seluruh 482 berkas; *penulis* dari `PropertiesName`,
*pembaca/pengulang* dari `PropertiesValue` · `pyValue` · `pyStepsObjectName`.

| Halaman | Penulis | Pembaca / pengulang | Perulangan *(ronde 1)* |
| --- | --- | --- | --- |
| `.ObjectItemList` | **19 berkas** | **54 berkas** | 32 |
| `ClaimData.ObjectList` | **26 berkas** | **57 berkas** | 22 |

⭐ **Keduanya ditulis oleh puluhan rule**, bukan hanya dibaca — jadi keduanya **memegang data**,
bukan sekadar halaman kerja bentukan.

⚠️ **Bertahan antar-permintaan?** ⛔ **Tidak terbaca dari korpus.** Keduanya halaman di dalam
`pyWorkPage`, dan umur `pyWorkPage` ditentukan mesin Pega, bukan ekspor.

⛔ **Nol kesimpulan tabel-atau-bukan-tabel.** Itu keputusan work owner.

---

## §G — ⛔ Yang tidak dijawab

| # | Butir | Status |
| --- | --- | --- |
| **G1** | *Rule dari lingkungan mana yang ditiru aplikasi Go* | ⛔ **TIDAK dijawab, TIDAK ditebak, TIDAK dianggap tertutup** |
| **G2** | *Apakah tipe berkas `avi` yang hanya ada di salinan `pegaprdnusare` ikut ditiru* | ⛔ **TIDAK dijawab** |
| **G3** | ✅ Menuliskan **apa bedanya** keempat rule itu | ✅ **dikerjakan** — §B7, langkah demi langkah |

⭐ **Dikonfirmasi sendiri: kedua pertanyaan §G tidak dijawab dan tidak diulang di §H4.**

---

## §H — Apa lagi

### H1 — Sesudah ronde ulang ini

| # | Yang dikerjakan | Kenapa menahan |
| --- | --- | --- |
| **1** | ⭐ **Uji 118 kandidat `pyMemo` golongan (c)** dengan membuka rule-nya | hanya 12 dari 459 yang benar-benar diuji; golongan (c) belum disentuh sama sekali |
| **2** | ⭐ **Sisir ulang identitas empat bagian atas SELURUH berkas keluaran lama** | usulan 2 menyentuh 20 modul, dan tiga berkas yang sudah jadi memakai identitas lama |
| **3** | **Baca `GetAllData_Act` langkah demi langkah** *(50 langkah)* | tulang punggung percabangan; F1 baru memetakan labelnya |
| **4** | **Isolasi satu lompatan kode 4** *(F3)* | ronde ini gagal mengisolasinya |
| **5** | **Sisir `pyXMLSignature` dengan pola ketiga** *(E1)* | 178 dari 179 terisi, isinya masih gelap |
| **6** | **Sisir medan kedua belas untuk gerbang aksi tombol** *(E3)* | pernyataan bersandar pada 11 medan |
| **7** | **Uji ADR-0012** sampai ke layar | satu-satunya ADR yang vonisnya *tidak dapat diputuskan* |

### H2 — Kesimpulan yang PALING RAWAN salah

⛔ **Ditunjuk, tidak diperbaiki.**

> ⚠️ **Paling rawan: §D2 — penggolongan 58 rule `When` sebagai HIDUP.**

Ia bersandar pada **satu** mata rantai: §A3 membuktikan kedua jalur halaman **punya penulis**, lalu
saya menyimpulkan tidak satu pun rule `When` mati karena halamannya tidak ada. ⚠️ **Yang tidak
saya periksa:** apakah **properti tertentu** yang diuji tiap rule *(`BusinessType`, `BusinessCode`,
`BusinessName`, `StatusBusiness`)* benar-benar ditulis — saya hanya membuktikan **halamannya**
ditulis. Sebuah halaman yang ada dengan properti yang tidak pernah diisi **tetap menghasilkan rule
yang selalu salah**.

⭐ **Ini bentuk yang sama dengan pola yang sudah berkali-kali menggigit proyek ini** — menyimpulkan
dari satu lapis tanpa menyisir lapis berikutnya.

**Paling rawan kedua:** §E2 *"38 dari 50 alias tidak cocok"*. Saya sendiri menemukan **3 positif
palsu** di dalamnya dan **tidak membersihkannya**; angka sebenarnya lebih rendah, dan berapa tepatnya
**belum dihitung**.

**Paling rawan ketiga:** §B10 vonis *"MENENTANG 6"*. Keenamnya diukur terhadap ADR yang
**lingkupnya Claim Life** — dan sebuah modul tidak dapat "menentang" keputusan yang tidak pernah
mengikatnya.

### H3 — ⭐ RALAT terhadap ronde 1: **tujuh angka berubah**

⛔ **Angka lama dikutip semua. Nol dihapus.** ⛔ `grilling-ronde-1.md` **tidak disunting.**

| # | Butir | ⛔ Angka LAMA *(ronde 1)* | ⭐ Angka BARU *(ronde ulang)* | Sebabnya |
| --- | --- | --- | --- | --- |
| **1** | rule bersama dengan Claim Prop | **162** | **161** | identitas empat bagian |
| **2** | `pxUpdateDateTime` identik | **158** | **159** | idem |
| **3** | "beda versi" | **4** | **2** *(beda waktu)* **+ 4** *(beda ruleset/versi)* | dua hal berbeda yang ronde 1 satukan |
| **4** | `When` menguji `pyWorkPage.Quotation` | **38 dari 60** | **34** jalur-1 **+ 4** jalur-2 | dua jalur digabung |
| **5** | golongan (d) `pyMemo` | **12 kandidat** | **6 bertahan · 6 pindah ke (b)** | diuji dengan membuka rule |
| **6** | gerbang aksi tombol | *"tidak ketemu di 3 medan"* | *"tidak ketemu di **11** medan"* | sisiran diperluas |
| **7** | selisih keluarga-2 | *"⚠️ satu kurang, tidak dikejar"* | ⭐ **KETEMU** — `SetCauseOfLossValue_act` 4.4, nol baris | dicari dengan jangkar |

⭐ **Dua pertanyaan ronde 1 terjawab dari korpus** *(Pertanyaan 2 sebagian, Pertanyaan 3 penuh)* —
sesuai perintah *"coba jawab dulu dari korpus"*.

### H4 — Pertanyaan BARU untuk work owner — **2**

⛔ Dua pertanyaan §G **tidak diulang**. Pertanyaan 1 ronde 1 *(versi mana yang berlaku)* **tidak
diulang** — ia sudah menggantung.

#### Pertanyaan 1 — Enam ADR ditentang korpus Fac In, padahal lingkup ADR adalah Claim Life

**Apa yang ditanyakan.** Enam dari 15 ADR ditentang korpus Claim Fac In: **0003** *(uang ditambal
koma-ke-titik)* · **0005** *(`IsPEGAPROD` hidup, 10 perujuk)* · **0008** dan **0015** *(nol outbox,
2 dari 52 efek keluar punya jalur gagal)* · **0011** *(ada dua tingkat di bawah `AdjustmentList`)* ·
**0013** *(satu endpoint URL langsung)*. Tetapi **kelima belas ADR bersumber dari Claim Life**, dan
**nol** menyebut Fac In.

**Kenapa muncul.** Ronde ulang ini adalah kali pertama korpus non-Life diadu dengan ADR.

**Bedanya kalau A atau B.** **A — ADR berlaku lintas-modul:** keenam pertentangan itu adalah
**penyimpangan sadar** yang harus dicatat di spec Fac In kelak, dan sebagian menyentuh uang.
**B — ADR hanya mengikat Claim Life:** keenamnya **bukan pertentangan sama sekali**, dan Fac In
memerlukan **ADR-nya sendiri** — yang berarti keputusan yang sama harus diambil ulang untuk tiap
modul.

**Apa yang tertahan.** Bentuk seluruh spec modul non-Life, dan apakah 15 ADR perlu diberi
pernyataan lingkup yang eksplisit.

#### Pertanyaan 2 — 60 rule `When` identik dipakai dua modul, vonisnya berlawanan

**Apa yang ditanyakan.** Keenam puluh rule `When` Claim Fac In **identik pada keempat bagian**
dengan Claim Prop — kelas, nama, ruleset, versi. Tetapi di **Claim Prop** 52 dari 61 sudah
diputuskan **tidak dimigrasikan** karena halaman yang diujinya tidak ada; di **Claim Fac In**
kedua jalur halaman **punya penulis**, sehingga alasan itu tidak berlaku.

**Kenapa muncul.** Baru terlihat setelah identitas empat bagian dipakai dan kedua jalur disensus.

**Bedanya kalau A atau B.** **A — satu rule, satu nasib:** kalau rule yang sama dimigrasikan untuk
Fac In, ia otomatis hidup juga untuk Prop, dan keputusan Claim Prop perlu ditinjau. **B — nasibnya
per modul:** rule fisik yang sama diperlakukan berbeda menurut modul pemanggilnya, dan aplikasi Go
menyimpan **dua salinan logika klasifikasi**.

**Apa yang tertahan.** Klasifikasi lini bisnis di **kedua** modul, dan AC 72 Claim Prop yang
penggolongannya memang masih `[terbuka]`.

### H5 — ⭐ Yang seharusnya dikerjakan tetapi TIDAK diperintahkan blok ini

⛔ **Disebutkan, tidak dikerjakan.**

| # | Butir |
| --- | --- |
| **1** | **Menyisir `pxUpdateSystemID` di 20 modul korpus lain** — A2 menunjukkan pola yang sama di dua modul; apakah seluruh korpus begitu belum diketahui |
| **2** | **Memeriksa apakah properti yang diuji tiap `When` benar-benar ditulis** — bukan hanya halamannya *(H2 butir pertama)* |
| **3** | **Membersihkan 3 positif palsu pada hitungan alias** *(E2)* dan menghitung ulang angkanya |
| **4** | **Membaca `BrowseHistoryClaim`** — kembaran `CariHistoryClaim_SQL` yang belum pernah tercatat di modul mana pun |
| **5** | **Memeriksa ruleset `ADESAMUEL@`** — ruleset bernama akun orang, 5 rule di Fac In, 2 di Prop |
| **6** | **Mengadu Claim Prop dan Komite Claim Prop dengan 15 ADR** — keduanya sudah selesai tanpa pernah diadu dengan dokumen |
| **7** | **Memeriksa `CreateKMTNo_Act` (42 langkah)** — pembuat kasus komite dengan **empat** penunjuk posisional, dua lebih banyak dari Claim Prop |

---

## Lampiran — bukti berkas lain tidak disentuh

Sidik jari MD5 atas **18 berkas** diambil sebelum ronde ini dan dibandingkan sesudahnya:
`grilling-ronde-1.md` · **15 ADR** · `ATURAN-BACA-KORPUS-PEGA.md` · `CONTEXT.md`.

**Satu-satunya berkas baru: `grilling-ronde-1-ulang-docs.md`.**
