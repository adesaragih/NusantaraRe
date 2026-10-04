# Bahan spec — Modul 3: Alur Masuk Endorsement (EDM)

> **Ini BAHAN untuk `/to-spec`, bukan spec.** `/to-spec` ber-`disable-model-invocation: true` dan
> **belum dijalankan**. Work owner yang menjalankannya manual.
>
> **Sumber:** `07-edm\03-e2-verifikasi-alur.md` · arsip `01-flow\03-alur-endorsement.md` §1, §1.3,
> §1.5. Prior art bentuk: `04-spec\05-spec-edm-before-image.md` (modul 1) dan `04-spec\04-spec-rnw.md`.
>
> **Keputusan mengikat:** K-005 · K-010/K-012 · K-027 · K-029 · K-044 (Life ikut; pilih-tertanggung
> dibuang) · K-046 (kode usang diport apa adanya) · **K-047 (total 4 seam)** · K-048.
>
> ✅ **Seam 5 disetujui work owner 19 September 2026 (K-049)** — `services/endorsement.OpenCase`.
> Total seam kini **5**. Modul selisih dan modul lain **tidak** mendapat seam baru. Lihat §8.
>
> ✅ **§9 dikoreksi (K-049):** rule pengecekan `BrowseOpenProteksiEdm_RD` **MASUK lingkup** dan
> diport apa adanya; hanya **isi tabel** `OPENPROTEKSI_EDM` yang dikelola di luar.

---

## 0. Posisi modul ini

Modul 3 adalah **pintu**: ia memutuskan apakah sebuah kasus endorsement **boleh lahir**, dan bila
boleh, menempatkannya di fase yang benar. Ia berjalan **sebelum** modul 1 (before-image) dan modul 2
(selisih).

⛔ Modul ini **orkestrasi dan gerbang**, bukan perhitungan. Tidak ada aritmetika uang di dalamnya.

---

## 1. Model case — dua kelas, satu penautan

### 1.1 Dua kelas case

`[terverifikasi]` Rule alur masuk hidup pada **`ASM-SFAGIS-Work-Endorsement`** (kelas portal),
sedangkan kasus yang **dilahirkan** berkelas **`ASM-FW-GISFW-Work-Endorsement`** berprefiks `EDM-`:

| Rule | `pxObjClass` | `pyClassName` |
| --- | --- | --- |
| `Activity\SetValueToEDMWork` | `Rule-Obj-Activity` | `ASM-SFAGIS-Work-Endorsement` |
| `DataTransform\DataToEDM` | **`Rule-Obj-Model`** | `ASM-SFAGIS-Work-Endorsement` |
| `Activity\SetErrorBatalEndorsement_Act` | `Rule-Obj-Activity` | `ASM-SFAGIS-Work-Endorsement` |
| `Activity\CheckEDMPolisDate` | `Rule-Obj-Activity` | `ASM-SFAGIS-Work-Endorsement` |
| `Activity\SetEdmType` | `Rule-Obj-Activity` | `ASM-SFAGIS-Work-Endorsement` |

⛔ **Dua kelas, bukan satu.** Portal tempat pengguna memilih polis berbeda kelas dari kasus
endorsement yang dihasilkan. Model domain harus memisahkannya.

### 1.2 Pembuatan kasus — `SetValueToEDMWork` langkah 7–11

`[terverifikasi]`

```
[7]  Property-Set  « Prepare the values for "createWorkPage" and "SpinoffNewWork" »
       SET param.classname      = "ASM-FW-GISFW-Work-Endorsement"
       SET param.IDPrefix       = "EDM-"
       SET param.modelname      = "pyDefault"
       SET param.workPage       = "curWorkPage"
       SET param.FlowType       = "pyStartCase"
       SET curWorkPage.PolicyNumber = .PolicyNo
       SET TempCase.PolicyNumber    = .PolicyNo
[8]  Call svcAddWorkObject                      ← kasus EDM lahir di sini
[9]  Property-Set  « untuk disable button OK »
[10] Obj-Open-By-Handle  « Re-open the page that refers to the new work created »
[11] Apply-DataTransform                        ← penautan, lihat §1.3
```

Lalu:

```
[12] SET newWorkPage.OfferFacIn.QuotationData.OldPolicyNo = .PolicyNo
[13] SET newWorkPage.OfferFacIn.EndorsmentReason          = .Note
[14] « Copy policy data to old data »           ← MODUL 1 (before-image) mulai di sini
[15] « Count StartProRate and EndProRate »      ← porsi periode (lihat K-048 §8.1)
```

📌 **Batas modul jelas:** modul 3 berakhir di langkah 13; langkah 14 adalah modul 1.

### 1.3 Penautan dua arah — `DataToEDM`

`[terverifikasi]` `DataTransform\DataToEDM` (`Rule-Obj-Model`), dipanggil langkah 11:

```
SET newWorkPage.pyLabel       = newWorkPage.pxInsName
SET newWorkPage.EndorsementID = .pzInsKey            ← handle kasus PORTAL disimpan di kasus EDM
UPDATE_PAGE newWorkPage
SET .EDMHandle                = newWorkPage.pzInsKey ← handle kasus EDM disimpan di kasus portal
```

⛔ **Tautan dua arah.** Kasus EDM menyimpan handle portal; kasus portal menyimpan handle EDM.
Implementasi harus menulis keduanya, bukan salah satu.

`[terverifikasi]` Langkah **22–25** `SetValueToEDMWork` menelusuri `Assign-Worklist` dan
`Assign-WorkBasket` dengan `.pxRefObjectKey == Primary.EDMHandle`, lalu menyimpan
`Primary.EDMHandle2 = .pzInsKey` — **handle assignment**, tautan ketiga.

---

## 2. Gerbang siklus — `IsEDM` dan `IsNotEDM`

`[terverifikasi]`

| Rule | `pxObjClass` | Kondisi (`pyConditionValue1String`) |
| --- | --- | --- |
| `When\IsEDM` | `Rule-Obj-When` | `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3` |
| `When\IsNotEDM` | `Rule-Obj-When` | `pyWorkPage.Quotation.StatusBusiness != 3` |

⛔ **`IsNotEDM` BUKAN negasi `IsEDM` (E-Q13).** Jalur propertinya berbeda:

| | Membaca |
| --- | --- |
| `IsEDM` | `OfferFacIn.QuotationData.StatusBusiness` — **agregat tersimpan** |
| `IsNotEDM` | `Quotation.StatusBusiness` — **halaman aktif** |

Karena sumbernya berbeda, **keduanya dapat bernilai benar bersamaan, atau salah bersamaan** bila
kedua halaman tidak sinkron.

⛔ **Mengikat implementasi:** resolver predikat **per-rule**, bukan satu predikat dengan negasi.
`IsNotEDM` wajib punya implementasi sendiri yang membaca halaman aktif. Kasus uji: kedua halaman
sengaja dibuat tidak sinkron, membuktikan keduanya bisa bernilai sama.

📌 `StatusBusiness = 3` sebagai penanda endorsement sudah terkunci (K-029 keluarga diskriminator
siklus: 1 = NB, 2 = RNW, 3 = EDM).

---

## 3. Enam gerbang penolakan

`[terverifikasi]` `Activity\SetErrorBatalEndorsement_Act` — **19 langkah tingkat atas**; seluruh
gerbang ada di **langkah 5** yang punya **24 sub-langkah**.

### 3.1 Pola seragam per gerbang

Setiap gerbang mengikuti rangkaian yang sama:

```
1. RDB-List / ReportDefinition   → ambil data pemeriksa
2. Property-Set                  → SET Param.Type = <1..4>          (siapkan klep)
3. Call pxRetrieveReportData     → isi page OpenProteksiEDM<X>      (klep pembatal)
4. Property-Set-Messages         → pasang galat, KECUALI klep terisi
```

### 3.2 Tabel keenam gerbang

`[terverifikasi]`

| # | Sub-langkah | Sumber data | Kondisi memblokir | Arti galat |
| ---: | --- | --- | --- | --- |
| **1** | 5.2 → 5.5 | `RDBList\GetEDMStatus_SQL` → `OutputData1.pxResults(1).CARI20` | `CARI20==1 \|\| CARI20==2` | polis sudah dibatalkan lewat endorsement |
| **2** | 5.7 → 5.9 | ReportDefinition `GetListEdm` → `ListEdm.pxResults` | `@LengthOfPageList(ListEdm.pxResults)>0` | ada EDM lain atas nopolis yang belum selesai |
| **3** | 5.11 → 5.14 | `RDBList\GetDataClaim_SQL` → `ListClaim.pxResults` | `@LengthOfPageList(ListClaim.pxResults)>0 && ListClaim.pxResults(1).CARI1!="2"` | sudah ada klaim atas nopolis itu |
| **4** | 5.15 → 5.18 | `RDBList\SearcStatusBayarArasaps_SQL` → `ListPembayaran.pxResults` | `@LengthOfPageList(ListPembayaran.pxResults)>0 && (Quotation.EdmType==1 \|\| ==2)` | sudah ada pembayaran **dan** endorsement bersifat pembatalan |
| **5** | 5.19 → 5.22 | `RDBList\GetListRNWbyNopolis_SQL` → `ListRNW.pxResults` | `@LengthOfPageList(ListRNW.pxResults)>0` | nopolis sudah punya renewal |
| **6** | 5.23 | `RDBList\GetFacoutList_SQL` | `.EndorsementInternalRetro==1` **dan** tidak ada baris fac out | polis tidak ter-spreading fac out |

`[terverifikasi]` Sub-langkah **5.24** merangkum: `IF @hasMessages(myStepPage)` → tandai gagal.

⛔ **Gerbang 4 hanya berlaku untuk endorsement pembatalan** (`EdmType` 1 atau 2 per K-029). Untuk
jenis endorsement lain, adanya pembayaran **tidak** memblokir.

### 3.3 Bypass "EDM RI Slip"

`[terverifikasi]` Sub-langkah **5.1** dan **5.10** melewati blok gerbang bila:

```
InputData.CARI40=="4" && InputData.CARI41=="4"     « lewatin utk edm ri slip »
```

⛔ Endorsement RI Slip **melewati sebagian gerbang penolakan**. Diport apa adanya.

### 3.4 Klep pembatal — empat, bukan enam

`[terverifikasi]` Klep memakai activity **platform Pega** `pxRetrieveReportData`, dengan nama dan
kelas laporan diberikan sebagai **parameter string runtime**:

```
pyReportName  = "BrowseOpenProteksiEdm_RD"
pyReportClass = "ASM-FW-GISFW-Int-OPENPROTEKSI_EDM"
```

Empat page hasil, satu per `Param.Type`:

| `Param.Type` | Sub-langkah | Page hasil | Gerbang yang diklep |
| ---: | --- | --- | --- |
| **1** | 5.3 → 5.4 | `OpenProteksiEDMBatal` | gerbang 1 (sudah dibatalkan) |
| **2** | 5.12 → 5.13 | `OpenProteksiEDMKlaim` | gerbang 3 (klaim) |
| **3** | 5.16 → 5.17 | `OpenProteksiEDMPembayaran` | gerbang 4 (pembayaran) |
| **4** | 5.20 → 5.21 | `OpenProteksiEDMRenewal` | gerbang 5 (renewal) |

Kondisi klep: `@LengthOfPageList(OpenProteksiEDM<X>.pxResults)>0` → **galat tidak dipasang**.

⛔ **Hanya 4 dari 6 gerbang punya klep.** Gerbang **2** (EDM belum selesai) memakai
`pxRetrieveReportData` untuk **mengambil** `ListEdm`, bukan sebagai klep. Gerbang **6** (fac out)
**tidak punya klep sama sekali** — tidak dapat dibatalkan.

### 3.5 Status `OPENPROTEKSI_EDM` — rule di dalam lingkup, isi tabel di luar

> ✅ **Dikoreksi atas K-049.** Versi pertama bahan ini menyatakan seluruh klep "di luar lingkup".
> **Itu salah.** Yang dibagi dengan menu lain adalah **tabelnya**, bukan rule pengecekannya.

⛔ **Pemisahan yang benar:**

| Bagian | Lingkup | Sikap |
| --- | --- | --- |
| **Rule pengecekan** `BrowseOpenProteksiEdm_RD` — dipanggil saat membuat EDM untuk memeriksa polis | ✅ **MASUK lingkup EDM** | **diport apa adanya**; klep untuk 4 gerbang (`Param.Type` 1=Batal · 2=Klaim · 3=Pembayaran · 4=Renewal) |
| **Isi tabel** `OPENPROTEKSI_EDM` — baris yang membatalkan penolakan | ⛔ **di luar korpus**, dikelola menu lain / DBA | `[pertanyaan terbuka]` bentuk dan sumbernya — **jangan ditebak**, dan **jangan `panic`**: tabel ini punya pemilik |
| **Klep gerbang 2 dan 6** | — | **tidak ada klep**; diport apa adanya (§3.4) |

⚠️ `[terverifikasi]` Berkas `BrowseOpenProteksiEdm_RD` **tidak ada** di `D:\migrasi\RNM\` — ia
dirujuk sebagai **string runtime** lewat activity platform `pxRetrieveReportData`. Token
`OPENPROTEKSI` muncul di **satu berkas saja** (`SetErrorBatalEndorsement_Act`, 55×), seluruhnya
sebagai nilai parameter.

⛔ **Mengikat implementasi:** rule pengecekan **wajib diimplementasikan** — ia bagian sah alur
penolakan EDM. Yang tidak dapat divalidasi saat kompilasi hanyalah **keberadaan laporan itu**,
karena rujukannya runtime. Lihat §9 untuk konsekuensinya.

---

## 4. Validasi tanggal — `CheckEDMPolisDate`

`[terverifikasi]` `Activity\CheckEDMPolisDate` — **6 langkah**:

```
[1] Property-Set
[2] « bypass » IF Primary.PolicyNo == <NOPOLIS-LITERAL>
[3] RDB-List  « select begin,end dari facinproduction »   RequestType = GetStartDate
[4] Property-Set  — bandingkan 8 karakter pertama tanggal mulai vs tanggal EDM
      [4.1] @substring(Local.StartDate,0,8)==@substring(Local.EdmDate,0,8)
      [4.2] « temporary utk monitoring »
[5] IF InputSearch.CARI20=="true"
      [5.1] set flag error   IF InputSearch.CARI1 > Local.StartDate
      [5.2] Property-Set-Messages  « tampilin error msg »
[6]
```

Aturan yang diuji `[terverifikasi]`:

```
@CompareDates(@toDate(Local.StartDate), @toDate(Local.EdmDate)) == "true"
  || @CompareDates(@toDate(Local.EdmDate), @toDate(Local.EndDate)) == "true"
```

⛔ **Tanggal endorsement wajib berada di dalam periode polis** (`facinproduction.begindate` …
`enddate`). Di luar itu → galat.

⚠️ `[terverifikasi]` Langkah 4.1 membandingkan **8 karakter pertama** string tanggal
(`@substring(…,0,8)`), bukan objek tanggal. Perbandingan tanggal sebagai string adalah bahaya yang
sudah terverifikasi di proyek ini (`CLAUDE.md` §4.1). **Diport apa adanya**, dicatat sebagai
kandidat perbaikan.

### 4.1 Nomor polis literal — dihitung, tidak disalin

`[terverifikasi]` **`CheckEDMPolisDate` memuat 2 nomor polis literal** yang membypass validasi
tanggal (langkah 2).

Sebaran di seluruh `Endorsment Fac In\Activity\` — **nilainya tidak disalin** (`CLAUDE.md` §3.5):

| Berkas | Kemunculan |
| --- | ---: |
| `SetErrorBatalEndorsement_Act` | **502** |
| `SumTSIPremiSpreadedRNM_Act` | 96 |
| `ProtectFIREMBUPA_Act` | 12 |
| `GetLimitAkseptasi_Act` | 8 |
| `SetValueToEDMWork` | 6 |
| `CheckEDMPolisDate` | **2** |
| `InputAddendumFacIn_PreAct` | 2 |
| `GetPaymentList_Act` | 1 |
| **TOTAL** | **629 kemunculan · 8 berkas · 176 nomor unik** |

⛔ **Diport apa adanya** (K-046) — kandidat perbaikan milik bisnis. ⚠️ Nomor polis produksi
**tidak boleh** menjadi literal di kode target; ia harus menjadi **data konfigurasi** yang dapat
diaudit. Itu keputusan rancangan, **bukan** perubahan perilaku — daftar pengecualiannya tetap sama.

📌 `SetValueToEDMWork` langkah 5 juga memuat bypass serupa pada gerbang "sudah dibatalkan"
(2 nomor literal di kondisinya).

---

## 5. Jenis endorsement — `SetEdmType`

`[terverifikasi]` `Activity\SetEdmType` · `pxRuleClassName` = `ASM-FW-GISFW-Data-Quotation` ·
4 langkah. Arti nilai sudah terkunci **K-029**:

| `EdmType` | Arti |
| ---: | --- |
| 1 | Batal Sejak Semula |
| 2 | Batal Prorata |
| 4 | Penambahan / Pengurangan / Perubahan |
| 3 | **usang** — cabang tetap diport apa adanya |

⚠️ `[pertanyaan terbuka]` yang tetap dari arsip §10 #12: nilai `Type = 3` (Adj Rate) dan `Type = 5`
(Add Object) punya rule `When` tetapi tidak muncul di dropdown `SetEdmType`; nilai `10` tidak dipakai.
Bagaimana pengguna memilihnya — belum terjawab. **Tidak memblokir modul ini.**

---

## 6. Gerbang masuk flow — EDM langsung berfase Policy

`[terverifikasi]` EDM punya **tepat 2 Flow**: `Flow\InputAddendumFacultativeIn.xml`
(`pxObjClass` = `Rule-Obj-Flow`, **69 shape / 138 konektor**) dan `Flow\OfferFacRetro.xml`.
**Tidak ada sub-proses siklus polis terpisah** — berbeda dari NB (6 flow) dan RNW (4 flow).

Konektor **`Start2 → Assignment7`** `[terverifikasi]`:

```
SET .FlagOnGoingPolicy = 1
SET .IsCedingConfirm   = Policy
SET .Position          = 1
SET .PositionNote      = "ReasFacInMarketing"
SET .NBStatus          = "NEW EDM"
SET .NBStatusNew       = "NEW EDM"
```

Shape `Assignment7` `[terverifikasi]`:

| Tag | Nilai |
| --- | --- |
| `pxObjClass` | `Data-MO-Activity-Assignment` |
| `pyMOName` | **`MARKETING`** |
| `pyImplementation` | `WorkBasket` |
| `pyRouteTo` | `Custom` · `pyIsCustomRouter` = `false` |
| `pyCallParams` / `Workbasket` | **`ReasFacInMarketing`** |
| `pyTicketShapes` | tiket **`AdminPolicy`** |

### 6.1 Kontras dengan NB dan RNW

⛔ **EDM tidak punya fase penawaran maupun binding sama sekali.** Ia **lahir** di fase polis dan
tidak pernah keluar darinya:

| Siklus | Urutan fase |
| --- | --- |
| NB / RNW | penawaran (`IsCedingConfirm="Offer"`, `FlagOnGoingPolicy=0`) → binding (`="Binding"`, `=2`) → polis (`="Policy"`, `=1`) |
| **EDM** | **langsung polis** (`="Policy"`, `=1`) |

`[terverifikasi]` Tidak ada satu pun konektor di flow EDM yang menetapkan `"Offer"`, `"Binding"`,
`0`, atau `2` untuk kedua properti itu.

⛔ **Mengikat model domain:** mesin fase endorsement **bukan** mesin fase NB dengan dua fase
dilewati. Ia mesin yang **hanya punya fase polis**. Memodelkannya sebagai NB-dengan-skip akan
memaksakan state yang tidak pernah ada.

---

## 7. Kontrak modul

### 7.1 Masukan

| Masukan | Sumber |
| --- | --- |
| Nomor polis yang di-endors | pilihan pengguna di portal (`ASM-SFAGIS-Work-Endorsement`) |
| Jenis endorsement (`EdmType`) | `SetEdmType` |
| Tanggal endorsement (`QuotationData.EdmDate`) | masukan pengguna |
| Alasan endorsement (`.Note`) | masukan pengguna |
| Penanda retro internal (`.EndorsementInternalRetro`) | data polis |
| Data pemeriksa 6 gerbang | 5 rule SQL + 1 ReportDefinition (§3.2) |
| Klep pembatal | tabel `OPENPROTEKSI_EDM` — **di luar lingkup** (§3.5) |

### 7.2 Keluaran

**Jalur berhasil:**

- kasus baru berkelas `ASM-FW-GISFW-Work-Endorsement`, prefiks `EDM-`;
- tautan **tiga arah**: `EndorsementID` (portal ← EDM), `EDMHandle` (EDM ← portal),
  `EDMHandle2` (handle assignment);
- `OldPolicyNo` dan `EndorsmentReason` terisi;
- kasus **berfase Policy** di workbasket `ReasFacInMarketing`, tiket `AdminPolicy`;
- siap diserahkan ke **modul 1** (before-image) mulai langkah 14.

**Jalur ditolak:** kasus **tidak lahir**; pesan galat terpasang pada field nomor polis, dengan
**alasan spesifik** per gerbang (§3.2).

⛔ **Penolakan bukan pengecualian teknis** — ia hasil bisnis yang sah dan harus dapat dibedakan dari
kegagalan sistem.

### 7.3 Tipe

⛔ **Tidak ada nilai uang di modul ini.** Yang mengalir: nomor polis (string), tanggal, kode jenis
endorsement (enumerasi terkunci K-029), dan flag. `Money` tidak dipakai.

⚠️ Tanggal wajib dibandingkan sebagai **tanggal**, bukan string — kecuali pada titik yang memang
membandingkan substring, yang **diport apa adanya** (§4).

---

## 8. Seam — ✅ Seam 5 DISETUJUI (K-049)

```
Seam 5 — services/endorsement.OpenCase(nomorPolis, jenisEndorsement, tanggalEndorsement)
              → Case  |  Penolakan
```

**Disetujui work owner 19 September 2026.** Satu pintu masuk yang mencakup seluruh modul alur masuk,
sejajar bentuk Seam 4.

| | Seam | Modul |
| ---: | --- | --- |
| 1 | `internal/rules.Eval` | registry predikat `When` |
| 2 | `services/acceptance.Next` | tangga akseptasi |
| 3 | `services/premium.Calculate` | perhitungan premi **dan selisih endorsement** (K-048) |
| 4 | `services/endorsement.PrepareBeforeImage` | before-image tiga lapis (K-047) |
| **5** | **`services/endorsement.OpenCase`** | **alur masuk EDM (K-049)** |

⛔ **Total 5 seam. Disiplin dijaga:** modul selisih dan modul berikutnya **tidak** mendapat seam
baru. Penambahan seam keenam butuh keputusan tersendiri.

### 8.1 Mengapa Seam 1 tidak cukup — bukti

`[terverifikasi]` `Activity\SetErrorBatalEndorsement_Act` memuat **37 gerbang berkondisi**
(`<pyStepsPreCondParamsWhen>` terisi). Pemilahannya:

| Bentuk kondisi | Cacah | Dapat diuji Seam 1? |
| --- | ---: | --- |
| **Ekspresi inline** (`@LengthOfPageList(...)>0`, `==`, `!=`, `>`) | **30** | ⛔ **tidak** |
| Nama rule polos — `IsFire` · `IsAneka` · `IsGolfInsurance` · `IsLife` · `IsMarineCargo` · `IsMBU` · `IsPA` | 7 | ✅ ya |

⚠️ **Gambarannya bernuansa, dan itu penting:** berkas ini **memang** memakai rule `When` — tetapi
hanya untuk **predikat lini bisnis**, bukan untuk gerbang penolakan. **Keenam gerbang penolakan
seluruhnya ekspresi inline.**

`[terverifikasi]` Token `@LengthOfPageList` muncul **49×** di berkas ini, **16** di antaranya di
dalam `<pyStepsPreCondParamsWhen>` (sisanya di ekspresi lain). `@hasMessages` muncul 4×.

📌 Catatan angka: penugasan menyebut 52×; penghitungan langsung memberi **49×**. Selisihnya tidak
mengubah kesimpulan — pemilahan 30 inline vs 7 rule `When` tetap sama.

```powershell
$t=Get-Content "D:\migrasi\RNM\Endorsment Fac In\Activity\SetErrorBatalEndorsement_Act.xml" -Raw
([regex]::Matches($t,'@LengthOfPageList')).Count          # => 49
$w=[regex]::Matches($t,'<pyStepsPreCondParamsWhen>([^<]*)</pyStepsPreCondParamsWhen>') |
   % { $_.Groups[1].Value.Trim() } | ? { $_ }
$w.Count                                                   # => 37
(@($w | ? { $_ -match '^[A-Za-z][A-Za-z0-9_]*$' })).Count  # => 7  (rule When)
(@($w | ? { $_ -match '@|==|!=|>|<' })).Count              # => 30 (inline)
```

Ketiga seam lain juga tidak cocok: `acceptance.Next` mengurus tangga persetujuan,
`premium.Calculate` mengurus angka, `PrepareBeforeImage` mengurus nilai lama — semuanya berjalan
**setelah** kasus lahir.

### 8.2 Yang diuji di Seam 5

Table-driven, satu kasus per baris:

- **Keenam gerbang** menolak pada kondisi yang benar, dan **hanya** pada kondisi itu — termasuk
  gerbang 4 yang **hanya** berlaku bila `EdmType` 1 atau 2.
- **Keempat klep** (`Param.Type` 1–4) membatalkan penolakan yang benar; **gerbang 2 dan 6 tetap
  menolak** meski tabel klep terisi — mereka tidak punya klep.
- **Bypass "EDM RI Slip"** (`CARI40=="4" && CARI41=="4"`) melewati gerbang yang benar.
- Kasus yang lahir berkelas `ASM-FW-GISFW-Work-Endorsement`, berprefiks `EDM-`, dengan **tautan tiga
  arah** lengkap.
- Kasus lahir **berfase Policy** — kasus uji negatif: tidak pernah melewati Offer atau Binding.
- **`IsNotEDM` bukan negasi `IsEDM`**: kedua halaman sengaja tidak sinkron, membuktikan keduanya
  dapat bernilai sama.
- Validasi tanggal menolak tanggal di luar periode polis, dan **bypass nomor polis literal bekerja**
  (memakai daftar pengecualian, bukan nilai tertanam — §4.1).
- **Penolakan terbedakan dari kegagalan sistem** — ia hasil bisnis yang sah dengan alasan spesifik.

---

## 9. Klep pembatal — pembagian lingkup (dikoreksi, K-049)

> ✅ **Dikoreksi atas K-049.** Versi pertama menempatkan seluruh klep "di luar lingkup". **Itu
> salah** dan sudah diperbaiki di sini serta di §3.5.

### 9.1 Rule pengecekan MASUK lingkup

`[terverifikasi]` `BrowseOpenProteksiEdm_RD` **dipakai untuk memeriksa polis saat membuat EDM** — ia
membaca tabel `OPENPROTEKSI_EDM` dan menentukan apakah penolakan dibatalkan. Itu **perilaku alur
masuk EDM**, bukan perilaku menu lain.

⛔ **Diport apa adanya**, sebagai klep untuk **4 gerbang**:

| `Param.Type` | Page hasil | Gerbang yang diklep |
| ---: | --- | --- |
| **1** | `OpenProteksiEDMBatal` | gerbang 1 — polis sudah dibatalkan |
| **2** | `OpenProteksiEDMKlaim` | gerbang 3 — ada klaim |
| **3** | `OpenProteksiEDMPembayaran` | gerbang 4 — ada pembayaran |
| **4** | `OpenProteksiEDMRenewal` | gerbang 5 — ada renewal |

⛔ **Gerbang 2** (`GetListEdm` — EDM lain belum selesai) dan **gerbang 6** (fac out) **tidak punya
klep**. Penolakannya **tidak dapat dibatalkan**. Diport apa adanya.

### 9.2 Isi tabel di luar korpus — punya pemilik

| | Status |
| --- | --- |
| **Bentuk dan sumber baris** `OPENPROTEKSI_EDM` | `[pertanyaan terbuka]` — dikelola **menu lain / DBA** |
| Sikap | ⛔ **jangan ditebak** · ⛔ **jangan `panic`** — tabel ini punya pemilik, ketiadaannya bukan cacat |
| Blocker? | **bukan** — modul ini dapat dispec dan dibangun tanpa menunggu isinya |

⚠️ Bedakan dari kasus `panic` di `CLAUDE.md` §4.5: di sana rule **dirujuk tetapi tidak ada pemiliknya**.
Di sini tabelnya **ada pemiliknya**, hanya isinya tidak terlihat dari korpus.

### 9.3 Konsekuensi teknis: rujukan runtime

`[terverifikasi]` Klep memanggil activity platform `pxRetrieveReportData` dengan nama dan kelas
laporan sebagai **literal string**, bukan rujukan rule:

```
pyReportName  = "BrowseOpenProteksiEdm_RD"
pyReportClass = "ASM-FW-GISFW-Int-OPENPROTEKSI_EDM"
```

⛔ Sistem baru **tidak dapat memvalidasi keberadaan laporan itu saat kompilasi**. Perilaku yang
mengikuti dari struktur ini `[terverifikasi]`: bila hasil klep kosong
(`@LengthOfPageList(OpenProteksiEDM<X>.pxResults) == 0`), **galat tetap dipasang** — yaitu
**penolakan berlaku penuh**.

📌 **Itu perilaku sistem lama, dan diport apa adanya.** Tabel kosong atau tidak terjangkau
menghasilkan penolakan paling ketat — arah yang aman. Kasus uji di Seam 5 wajib mencakup keadaan ini.

---

## 10. Out of Scope

1. **Fitur pilih-tertanggung** — **K-044: DIBUANG**. Keempat berkasnya ada di korpus tetapi tidak
   diport.
2. **Isi tabel `OPENPROTEKSI_EDM`** dan menu yang mengisinya — dikelola di luar (K-049).
   ⚠️ **Rule pengecekannya `BrowseOpenProteksiEdm_RD` TIDAK di luar lingkup** — ia diport apa adanya
   (§9.1). Yang di luar hanyalah **isi tabelnya**.
3. **Tangga akseptasi endorsement** — termasuk klaim arsip §4.3 *"nilai dasar akseptasi = selisih
   TSI"*, yang tetap **`[belum diuji]`**.
4. **Modul 1 (before-image)** mulai `SetValueToEDMWork` langkah 14 — modul tersendiri, sudah dispec.
5. **`Flow\OfferFacRetro`** — flow kedua EDM, jalur retro/fac out; modul tersendiri.
6. **Alur endorsement Life** setelah kasus lahir — **K-044 Life IKUT**, tetapi jalurnya terpisah dan
   dispec di modulnya sendiri. Alur **masuk**-nya sama dengan lini lain.
7. **Lapisan tampilan portal** — pemilihan polis, dropdown jenis endorsement.

---

## 11. Kesiapan

✅ **Bahan lengkap untuk `/to-spec`:**

| Butir | Status |
| --- | --- |
| Dua kelas case + pembuatan (langkah 7–11) + penautan tiga arah | §1 |
| Gerbang siklus `IsEDM`/`IsNotEDM` + penegasan E-Q13 | §2 |
| **Enam gerbang penolakan** — tabel sumber → kondisi → arti | §3.2 |
| Bypass "EDM RI Slip" | §3.3 |
| **Klep pembatal — 4 dari 6**, dengan `Param.Type` 1–4 dan page hasilnya | §3.4 |
| Status `OPENPROTEKSI_EDM` (K-049) — **rule pengecekan MASUK lingkup**, isi tabel di luar, **bukan blocker** | §3.5, §9 |
| Validasi tanggal + perbandingan substring yang diport apa adanya | §4 |
| Nomor polis literal — **629 / 8 berkas / 176 unik**, nilai tidak disalin | §4.1 |
| Jenis endorsement K-029 | §5 |
| Gerbang masuk flow + kontras NB/RNW (tanpa fase penawaran/binding) | §6 |
| Kontrak masukan/keluaran, jalur berhasil dan ditolak | §7 |
| ✅ **Seam 5 disetujui (K-049)** + bukti 30 inline vs 7 rule `When` + daftar uji | §8 |
| ✅ **Klep dikoreksi (K-049)** — rule masuk lingkup, isi tabel di luar, rujukan runtime | §9 |

✅ **Tidak ada pertanyaan terbuka yang memblokir.** Seam sudah diputuskan; pembagian lingkup klep
sudah dikoreksi. Bahan **siap untuk `/to-spec`**.

⛔ **`/to-spec` TIDAK dijalankan** — menunggu perintah work owner.

---

*Tanpa nama orang, tanpa alamat email, tanpa nomor polis, tanpa data pelanggan.*
