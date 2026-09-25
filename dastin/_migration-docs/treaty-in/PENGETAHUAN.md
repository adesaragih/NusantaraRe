# Modul Treaty In (Pega) — Pengetahuan AS-IS untuk Migrasi ke React + Golang + Oracle

Status: analisis AS-IS, per 22 September 2026 (diperbarui sore hari yang sama setelah user
menyerahkan sumber prosedur + DDL dan menjawab pertanyaan terbuka). Belum ada kode baru yang dibuat.
Sumber: 329 berkas XML export Pega di `D:\XML_NURE\Treaty In` (±80 MB), 5 stored procedure dan
7 DDL Oracle di `_migration-docs/treaty-in/PROCEDURE` dan `/Table`, + indeks pohon aturan
`Struktur_InputTreatyInOffer.xlsx` (1.138 baris pemanggilan, seluruhnya berstatus "Sudah diupload",
artinya rantai aturan modul ini lengkap — tidak ada XML yang masih kurang).

---

## 1. Inventaris berkas

| Folder | Jumlah | Ukuran | Isi |
|---|---|---|---|
| Activity | 149 | 21 MB | seluruh logika prosedural (perhitungan, simpan, konversi) |
| Section | 50 | 48 MB | definisi layar; `InputTreatyInOffer.xml` sendirian 8 MB |
| FlowAction | 32 | 1,1 MB | pembungkus tipis; hampir semuanya hanya memanggil satu Section senama |
| RDBList | 40 | 336 KB | seluruh akses Oracle (SQL mentah + panggilan stored procedure) |
| DataTransform | 35 | 700 KB | aturan deklaratif kecil, termasuk mesin approval `Akseptasi_DT` |
| ReportDefinition | 16 | 932 KB | pencarian data master (currency, client, treaty group, dsb.) |
| Harness | 3 | 8,9 MB | `InputTreatyInOffer` (entry point), `ShowAttachmentTreaty`, `TreatyInFacultativeShareCalculation` |
| When / DecisionTable / ConnectREST / SystemSettings | 1 tiap | kecil | `TreatyMasterInEDM`, `GetMimeType`, `ServiceGoogle`, `LinkService` |

Ada folder sejajar `D:\XML_NURE\Treaty In Adjustment` yang belum dibedah pada sesi ini.

---

## 2. Arsitektur aplikasi lama

Rantainya: **Harness `InputTreatyInOffer` (kelas `Data-Portal`) → Section → Flow Action → Activity →
RDB List (Oracle) / Report Definition (tabel master)**.

Seluruh modul adalah **satu layar tunggal** yang sangat besar. Tidak ada Pega Flow / case type —
tidak ada satu pun berkas Flow di export. Status dan approval dikelola manual lewat properti
`StatusAkseptasi` + `Position` yang diubah oleh Data Transform, bukan oleh workflow engine Pega.
Konsekuensi migrasi: **tidak perlu BPM engine**; state machine sederhana di Golang sudah cukup.

Struktur Section utama (dari `InputTreatyInOffer`):

```
InputTreatyInOffer (harness, Data-Portal)
├── TreatyInActionButtons        ← tombol Save/Submit/Copy/Revision/View/Edit
├── TreatyInTabsProportional     ← cabang UI untuk treaty proporsional
├── TreatyInNONProportional      ← cabang UI non-proporsional (7 MB, terbesar kedua)
│   ├── TreatyInTabsNonProportional
│   │   ├── TreatyInFacultativeShareCalculation
│   │   └── TreatyInTabsNonProportionalValueDifference
│   │       ├── TreatyInTabsNPValueDifferenceProRate
│   │       └── TreatyInTabsNPValueDifference_NoProRate
│   ├── TreatyInTabsNonProportionalAdjustPremi
│   ├── TreatyInActualLimits / TreatyInActualShare / TreatyInActualSumary
│   └── TreatyInFacultativeRetro / TreatyInShareProp
├── TreatyInTabsAchievement      ← realisasi premi & klaim per kuartal
├── TreatyInfoSubmit             ← komentar, informasi, tombol submit
└── WorkAttachments              ← lampiran dokumen
```

Flow Action yang dipakai sebagai modal/grid: `Layers`, `LayersEDM`, `Share`, `ShareOldData`,
`ShareRetro`, `DetailLimits`, `DetailShare`, `DetailEGNPI`, `Installments`, `CoBList`,
`MaxRetention`, `LimitProportional`, `SpreadingTPDtl`, `SpreadingTXOLDtl`, `TreatyRetroList`.

---

## 3. Model data clipboard (calon model domain)

Akar halaman adalah `TreatyIn`, kelas **`ASM-FW-GISFW-Int-TREATY_IN`**. Kelas anaknya:

| Kelas Pega | Peran | Properti kunci |
|---|---|---|
| `Int-TREATY_IN` | header kontrak treaty | ID, OLDID, ProportionType, TreatyContractName, TeritorialScope, Commencement, Termination, TreatyYear, Ceding/CedingID, LeadingReinsSource/ID, LeadingReinsName, RNMShare, RNMShareP, BrokeragePercent, BrokeragePercentP, StatusAkseptasi, ChooseStatusAkseptasi, Position, PositionUsername, ViewState, RevisionState, IsEditData |
| `Data-TreatyInLimits` | satu **layer** (non-prop) / satu kelompok limit (prop) | LayerType, Layer, LayerPartType, LayerPart, Currency, Currency2, Limit, Limit2, Deductible, Deductible2, AgregateLimit(2), AdjRate, MDPPct, MDPMinPct, ROLPct, Cover, CurrencyRelation, IsCombineMDP, NoRIPCalculation, Reinstatement_List, MDPList, MDPMinList, PremiumEarnedList, EgnpiTotalList, TreatyGroupList, Detail |
| `Data-TreatyInLimitsDetail` | rincian per Treaty Group / CoB di dalam layer | TreatyType, TreatyGroup(ID), ClassOfBusiness(ID), QSPct, Surplus, RetentionPct, CessionPct, RIOGR, RIONR, IOOLimitList, RetentionList, CessionList, EPIList, PLAList, CashLossList, ClaimCoopList, ReserveList, DeductionList, RNMShare, RNMShareList, SpreadingList, SpreadingTypeID, SpreadingTotalPct, RNMSpreadedList(RI), COBList, AchievementLists, LossRatio, ProfitCommision, UpperBand/LowerBand, PremiumReservePct |
| `Data-TreatyInShare` | share NuRe per layer (non-prop) | RnmLimitList, GrossPremiumList, GrossPremiumMinList, DeductionList, DeductionTotalList, NetPremiumList, SpreadingListXOL, SpreadingTypeXOL(ID), SpreadingTotalPctXOL, dan 10 daftar `RNMSpreadedList…XOL` (OR/RI × Limit/Gross/GrossMin/Deduct/Net), RnmLimitListDisplay, RnmGrossPremiDisplay |
| `Data-TreatyInLimitsSpreading` | satu baris penyebaran (QS OR / QS RI / ORS) | ReinsTypeID, ReinsTypeName, ParentReinsTypeID, Pct, Rp, Usd, BreakDownSprdList(XOL) |
| `Data-TreatyInShareReins` | retrosesioner / broker pada share fakultatif | ReinsID, ReinsName, BrokerID, BrokerName, SharePct, Amount, Amount2, FacultativeLimits |
| `Data-TreatyInRetroShare` | paket retro | RetroCoName, RetroType(ID), RetroPct, RetroBrokerage, RetroSpreadingList, RetroTotalPct, Share, ShareSumary, Total* |
| `Data-TreatyInEGNPI` | estimasi GNPI per treaty group | TreatyGroup, Currency, Amount, AmountIDR, AsDate, Proportion |
| `Data-TreatyInRetention` | retensi maksimum cedant | TreatyGroup, Currency, Amount, Note |
| `Data-TreatyInInstallment` | termin pembayaran premi | Installment, InstallmentPct, Amount, DueDate, PaymentDate, WPC, PctTotal, AmountTotal |
| `Data-TreatyInAccountReport` | jadwal pelaporan akun | Period, InitialDate, SubmissionDue, ConfirmationDue, SettlementDue, AutoCalculate |
| `Data-TreatyInAccumulation` | jadwal akumulasi | Period, ReportDate, SubDays, SubDueDate |
| `Data-TreatyInCurrencyList` | kurs berlaku per periode | CurrencyID, Currency, Conversion, PeriodStart, PeriodEnd |
| `Data-TreatyInCoInScale` | skala ko-asuransi | PctLimit, CoInShare |
| `Data-TreatyInPortfolio` | portofolio masuk/keluar | Type, TypePortfolio, Description |
| `Data-TreatyInTotal` | baris total generik multi-mata-uang | Currency, CurrencyID, Value, IDR, USD, Note |
| `Data-LimitSummaryList` | baris ringkasan per layer | Note, Limit(2), Deductible(2), MDP(2), NetPremi(2), AggregateLimit(2) |
| `Data-SuggestList` | riwayat komentar/approval | Date, OperatorName, IsApproved, Suggest |

**Pola multi-mata-uang.** Hampir setiap nilai uang disimpan sebagai *page list* `(Currency, Value)`,
bukan kolom skalar. Total dihitung dengan pola berulang: loop daftar → cocokkan mata uang → kalau
sudah ada tambahkan, kalau belum `<APPEND>` baris baru. Di Golang ini menjadi satu tipe
`Money{CurrencyID, Currency, Value}` dan satu fungsi agregasi, menggantikan ±40 salinan logika
yang sama di Pega. Sebagian layar lama masih memakai pasangan kolom kembar `Limit`/`Limit2`
dan `Deductible`/`Deductible2` yang secara de facto berarti IDR/USD — dua representasi hidup
berdampingan di modul yang sama.

---

## 4. Persistensi Oracle

Seluruh akses DB lewat RDB List. Dua gaya berdampingan:

**a. Simpan lewat stored procedure di skema `POOLDATA`** (jalur utama):

| Prosedur | Dipanggil dari | Isi |
|---|---|---|
| `POOLDATA.PEGA_TREATY_IN(DATAPEGA, ID, ProportionType, …, TreatyYear, ERRMSG out, IDPEGAOUT out, STSSAVE out)` | `SaveTreatyIn.xml` | simpan header + JSON penuh |
| `POOLDATA.PEGA_M_TREATY_IN_EDM(… , OLDID, STSCARI1, …, EDMState, EDMMaterialType, out×3)` | `SaveTreatyInEDM.xml` | simpan versi addendum |
| `POOLDATA.PEGA_M_TREATY_IN_DETAIL(…67 parameter…, out×3)` | `SaveTreatyInDetail.xml` | flat table detail |
| `POOLDATA.PEGA_M_TREATY_IN_DETAIL_EDM(…65 parameter…, out×3)` | `SaveTreatyInDetailEdm.xml` | flat table detail versi EDM |
| `POOLDATA.GET_TOKEN_STORAGE(App, UserID, Kodestring out, ResponseMsg out)` | `GetTokenStorage_SQL.xml` | token akses Google Cloud Storage |

Parameter pertama `DATAPEGA` selalu berisi **seluruh halaman clipboard yang di-serialisasi jadi JSON**
(`@ASM.GetPageJSONString()`), disimpan ke kolom `JSONDATA` bertipe JSON Oracle. Jadi kebenaran data
ada di dokumen JSON; kolom relasional adalah proyeksi denormalisasi untuk pelaporan.
Sumber prosedur dan DDL-nya sudah diterima dari pihak user pada 22 September 2026 dan dibedah di
**bagian 4A** di bawah.

**b. SQL langsung** untuk baca dan lampiran:

- `select JSONDATA from pooldata.M_TREATY_IN where ID = :id` **union all** `pooldata.M_TREATY_IN_EDM`
  (`BrowseTreatyIn.xml`) — data live dan addendum dibaca sebagai satu aliran.
- `TREATYINDETAIL`, `TREATYINDETAILEDM`, `M_TREATY_IN_DETAIL`, `M_TREATY_IN_DETAIL_EDM` — dihapus
  massal per `TREATYID` sebelum setiap simpan (pola delete-then-insert, `RemoveTreatyInDetail.xml`).
- `POOLDATA.TREATYINOFFER` — insert baris ringkas untuk integrasi hilir (`InsertToTable.xml`).
- `POOLDATA.TREATYINPRODUCTION`, `POOLDATA.ACHIEVEMENT`, `POOLDATA.LOG_ACHIEVEMENT` — realisasi.
- `POOLDATA.OS_AKSEPTASI_KLAIM` — klaim; dibaca lewat kolom JSON `a.data_json.Value`, dibedakan
  klaim final vs estimasi oleh `a.DATA_JSON.AcceptedNo IS NULL`.
- `M_ATTACHMENTTREATY_2`, `T_STORAGE_IMAGE`, `T_FOLDER_IMAGE`, `CATEGORY_ATTACH_REAS`,
  `M_KATEGORIMASTERTREATY` — lampiran.
- `TREATYEXCHANGEYEARLY` — kurs (`select TOIDR … order by startdate desc … rownum = 1`).
- `PROPORTIONALARRG` — master susunan treaty NuRe sendiri (sumber persentase penyebaran).
- `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` — cek apakah master treaty sudah punya klaim.

Tabel master lain dibaca lewat Report Definition (kelas `Int-…`): `AGENT`, `CLIENT`, `CURRENCY`,
`BUSINESS`, `BUSINESSGROUP`, `TREATYGROUP`, `TREATYBUSINESS`, `REINSURANCETYPE`,
`TREATYREINSURER`, `TREATYEXCHANGEYEARLY`, `PROPORTIONALARRG`.

Catatan penyaring yang bersifat aturan bisnis, bukan sekadar UI:
`BrowseAgentNusaRe_RD` menolak `AgentType2 = "LIFE INSURANCE"` dan hanya `StatusActive = 1`;
`BrowseCurrency_RD` membuang mata uang `ITL`; `BrowseBusinessGroup_RD` membuang grup berakhiran
`SYARIAH`; `BrowseTreatyTypeParent` hanya mengambil `ParentReinsTypeID = "00"` yang namanya
mengandung `TRT` atau `XOL` dan periodenya melingkupi tanggal mulai treaty.

---

## 4A. Sumber Oracle yang sudah diterima (22 September 2026)

Berkas ada di `_migration-docs/treaty-in/PROCEDURE/` (5 prosedur) dan `_migration-docs/treaty-in/Table/`
(7 DDL). Ini menutup dua gap terbesar dari sesi sebelumnya.

### 4A.1 Objek Oracle yang sebelumnya tidak terlihat dari XML

Prosedur menyingkap tabel dan sequence yang tidak pernah disebut satu pun RDB List:

| Objek | Peran |
|---|---|
| `POOLDATA.TREATY_IN` | **proyeksi relasional header** (20 kolom); pasangan dari `M_TREATY_IN` yang berisi JSON |
| `POOLDATA.TREATY_IN_EDM` | proyeksi relasional header addendum (+ `OLDID`, `EDMDATE`, `EDMSTATE`, `EDMMATERIALTYPE`) |
| `POOLDATA.M_TREATY_IN_EDM` | JSON header addendum (`ID`, `OLDID`, `EDMDATE`, `JSONDATA`) |
| `POOLDATA.M_TREATY_IN_DETAIL` / `_EDM` | JSON detail (`ID`, `JSONDATA`) |
| `POOLDATA.TREATYINDETAILEDM` | proyeksi relasional detail addendum |
| `POOLDATA.M_SITE_DATABASE` | menyimpan kode situs aktif (`CURRENT_SITE = '1'`) — jadi prefiks ID |
| `M_TREATY_IN_SEQ`, `POOLDATA.M_TREATY_IN_DETAIL_SEQ` | sequence penomoran |
| `POOLDATA.GCP_IMAGE` | cache token Google Cloud Storage |

Jadi polanya **selalu berpasangan**: satu tabel `M_…` menyimpan dokumen JSON utuh, satu tabel
relasional menyimpan kolom pilihan untuk pelaporan. Keduanya ditulis dalam satu prosedur.

### 4A.2 Cara ID dibentuk

```sql
SELECT ID INTO id_site FROM M_SITE_DATABASE WHERE CURRENT_SITE = '1';
id_ins := id_site || lpad(to_char(M_TREATY_IN_SEQ.nextval), 6, '0');
```

`M_TREATY_IN.ID` bertipe `VARCHAR2(10)`, jadi kode situs maksimal 4 karakter. Pega mengirim
`ID = 'UnknownId'` untuk menandai record baru, lalu prosedur melakukan
`replace(DataPega, 'UnknownId', id_ins)` **atas seluruh teks JSON** sebelum disimpan. Di sistem baru,
pakai ID yang dihasilkan aplikasi (atau sequence) dan jangan menyunting dokumen dengan string replace.

### 4A.3 Perilaku tiap prosedur

- **`PEGA_TREATY_IN`** — bila `IDPega = 'UnknownId'`: terbitkan ID, `INSERT` ke `M_TREATY_IN` **dan**
  `TREATY_IN`, `COMMIT`. Selain itu: `UPDATE` keduanya berdasarkan ID. Mengembalikan
  `ErrMsg = 'Data Sudah Disimpan Dengan ID : …'`, `StsSave = 1`, `IDPegaOut`. Gagal → `StsSave = 0`
  dan `ROLLBACK`.
- **`PEGA_M_TREATY_IN_EDM`** — dikendalikan `STSINPUT` (yang di Pega diisi dari hasil hitung
  `EDMCheckExistingData`): `'0'` berarti belum ada → `INSERT` ke `M_TREATY_IN_EDM` + `TREATY_IN_EDM`
  dengan `EDMDATE = SYSDATE`; selain itu `UPDATE`. ID **tidak** diterbitkan di sini — dipakai apa adanya
  dari Pega.
- **`PEGA_M_TREATY_IN_DETAIL`** dan **`PEGA_M_TREATY_IN_DETAIL_EDM`** — **hanya menangani cabang
  `IDPega = 'UnknownId'`, tidak ada cabang `ELSE`.** Ini konsisten dengan pola hapus-lalu-sisip:
  Pega selalu menghapus seluruh detail per `TREATYID` lebih dulu, lalu menyisipkan ulang semuanya
  dengan ID baru. Konsekuensinya **ID baris detail tidak stabil antar penyimpanan** — tidak bisa
  dijadikan acuan oleh sistem lain.
  Kedua prosedur mengambil `COMMENCEMENT`/`TERMINATION` bukan dari parameter melainkan dengan query
  ke header (`TREATY_IN` / `TREATY_IN_EDM`) memakai `to_date(…,'YYYYMMDD')` — artinya **tanggal di
  tabel header disimpan sebagai teks `YYYYMMDD`**, bukan `DATE`.
- **`GET_TOKEN_STORAGE`** — mencari token di `GCP_IMAGE` yang `INPUTDATE > SYSDATE`; bila tidak ada,
  menerbitkan token baru `RAWTOHEX(STANDARD_HASH('ASMAPP'||timestamp,'MD5'))` dengan masa berlaku
  `SYSDATE + 1 menit`. Kolomnya bernama `INPUTDATE` tetapi **semantiknya tanggal kedaluwarsa** —
  jangan tertukar saat migrasi.

### 4A.4 Temuan baru dari pembacaan DDL/prosedur

- **TD-14 — tanda tangan prosedur tertinggal dari aturan Pega.** `SaveTreatyInDetail.xml` mengirim dua
  parameter terakhir `SPREAD_RNM_SHARE_PCT` dan `SPREAD_RNM_SHARE_VALUE` (memo aturannya memang
  berbunyi "ADD …"), tetapi `PEGA_M_TREATY_IN_DETAIL` yang diterima berhenti di `P_MDP` dan DDL
  `TREATYINDETAIL` tidak punya kolom `SPREAD_RNM_SHARE_*`. Salinan yang diterima lebih tua dari aturan
  Pega yang berjalan. Sesuai kebijakan proyek, **XML diperlakukan sebagai sumber kebenaran**: model
  baru tetap menyediakan kedua field itu, dan salinan prosedur ini dianggap tidak mutakhir.
- **TD-15 — urutan parameter spreading berbeda antara dua prosedur.** Versi non-EDM menerima
  `P_SPREADINGTYPE` lalu `P_SPREADINGTYPEID`; versi EDM kebalikannya, `P_SPREADINGTYPEID` lalu
  `P_SPREADINGTYPE`. Masing-masing cocok dengan pemanggilnya, jadi tidak ada bug hari ini, tetapi ini
  jebakan klasik saat orang menyalin kode dari satu jalur ke jalur lain.
- **TD-16 — kegagalan ditelan diam-diam.** Pada kedua prosedur detail, blok penangan galat untuk
  `M_SITE_DATABASE` dan sequence hanya `ROLLBACK; RETURN;` tanpa mengisi `ErrMsg` **maupun**
  `StsSave`, sehingga Pega menerima nilai kosong dan tidak bisa membedakan gagal dari sukses. Bila
  `IDPega != 'UnknownId'`, prosedur juga langsung jatuh ke `StsSave := 1` padahal tidak ada yang
  ditulis. Sistem baru harus mengembalikan galat yang eksplisit.
- **TD-17 — `COMMIT` di dalam prosedur per baris.** Tiap baris detail disimpan dan di-`COMMIT`
  sendiri-sendiri, jadi penyimpanan satu treaty bukan satu transaksi. Kalau gagal di baris ke-N,
  data tertinggal separuh (sebagian detail lama sudah dihapus, sebagian baru sudah masuk). Di sistem
  baru satu penyimpanan treaty harus satu transaksi.
- **TD-18 — kunci lampiran rawan tabrakan.** `M_ATTACHMENTTREATY_2.ID` adalah primary key
  `VARCHAR2(100)` yang diisi `TO_CHAR(SYSTIMESTAMP,'YYYYMMDDHH24MISSFF3')`; dua unggahan dalam
  milidetik yang sama akan bentrok. Sama halnya `T_STORAGE_IMAGE.IMAGEID` yang memakai hash MD5
  dari timestamp.
- **TD-19 — tipe data terlalu longgar.** `PROPORTIONALARRG` menyimpan **semua** kolom sebagai
  `VARCHAR2(1000)`, termasuk `PCT`, `LINE`, `RP`, `USD` — itulah sebabnya XML harus memanggil
  `to_number(pct)` dan memo aturannya berbunyi "added to_number function in pct, it works in
  production server". `TREATYINDETAIL` pun mencampur: sebagian nilai numerik bertipe `NUMBER`,
  sebagian (`LIMIT_RSMD`, `SPL_LINE`, `QS_PCT`, dan seterusnya) diterima prosedur sebagai `VARCHAR2`.
  Skema baru harus bertipe tegas sejak awal.
- **TD-20 — tidak ada indeks pada kolom penyaring utama.** DDL `TREATYINDETAIL`, `ACHIEVEMENT`,
  `TREATYINPRODUCTION`, dan `PROPORTIONALARRG` tidak mendefinisikan primary key maupun indeks,
  padahal `TREATYINDETAIL` dihapus massal per `TREATYID` pada setiap simpan dan `ACHIEVEMENT`
  dicari dengan `SUBSTR(NOOFFER,1,7)` (predikat yang memang tidak bisa memakai indeks biasa).
- Catatan administratif: berkas `Table/PEGA_M_TREATY_IN_DETAIL.txt` sebenarnya berisi DDL tabel
  **`M_TREATY_IN`** (`ID VARCHAR2(10)`, `JSONDATA CLOB` dengan `CHECK (JSONDATA IS JSON)`), bukan
  tabel detail. Namanya saja yang tidak cocok dengan isinya.

### 4A.5 DDL yang masih belum ada

`TREATY_IN`, `TREATY_IN_EDM`, `M_TREATY_IN_EDM`, `M_TREATY_IN_DETAIL(_EDM)`, `TREATYINDETAILEDM`,
`M_SITE_DATABASE`, `GCP_IMAGE`, `LOG_ACHIEVEMENT`, `TREATYEXCHANGEYEARLY`, `M_KATEGORIMASTERTREATY`,
`CATEGORY_ATTACH_REAS`, `M_LINK_SERVICE`, `T_FOLDER_IMAGE`, `TREATYINOFFER`, dan tabel master
(`AGENT`, `CLIENT`, `CURRENCY`, `BUSINESS`, `TREATYGROUP`, `REINSURANCETYPE`, `TREATYREINSURER`).
Struktur kolomnya sudah bisa disimpulkan dari SQL dan parameter prosedur, jadi ini tidak menghambat;
cukup dikonfirmasi saat verifikasi skema.

---

## 5. Alur persetujuan (`Akseptasi_DT`)

Mesin approval seluruhnya ada di satu Data Transform. Jalur normal (`RevisionState != 1`):

```
(kosong)/ReasTreatyInAdmin → ReasTreatyInSecHead → ReasTreatyInDeptHead → ReasTreatyInDirector → Resolve Complete
```

`ReasTreatyInGroupLeader` adalah jalur alternatif yang langsung ke Director. Di tiap simpul,
`ChooseStatusAkseptasi` menentukan: `Accept` → naik satu tingkat; `Reject` → kembali ke
`ReasTreatyInAdmin` dengan `PositionUsername` diisi nama pembuat dari `CommentList(1)`;
`Decline` → `Position` dikosongkan dan status jadi `Decline` (final).

Jalur revisi (`RevisionState == 1`) dipendekkan: Admin → SecHead → langsung `Resolve Complete`,
sekaligus mengosongkan `RevisionState` dan `ViewState`.

**Temuan penting — TD-01.** Nama penerima tugas **di-hardcode sebagai string nama orang** di dalam
aturan: `IRVANDY`, `AGUNGPUTRAANDALAS`, `YOHANESKRISTIAWAN`, `NANDINA`. Pemilihan antara IRVANDY
dan AGUNGPUTRAANDALAS ditentukan oleh nilai `OperatorID.pyTelephone` yang diisi kode
`SPVTREATY1`/`SPVTREATY2`/`TREATY1`/`TREATY2` — kolom telepon dipakai sebagai kode peran. Di sistem
baru ini harus menjadi tabel peran-dan-penugasan, bukan konstanta.

Pemicu dan efek samping:
- `TreatyInSubmit` → validasi `TreatyInCheckID` + `TreatyInCheckError` → `Akseptasi_DT` →
  `AddCommentList_Act` → `SaveTreatyIn_Act` → `TreatyInInputVis`.
- `TreatyInSubmitEDM` sama, tetapi menyisipkan `TreatyEDMCalculateDifference` bila
  `EDMState` bernilai `"1"` atau `"2"`, lalu menyimpan lewat `SaveTreatyIn_EDM_Act`.
- `TreatyInCopy` menyalin kontrak: `CommentList` dikosongkan, `OLDID` diisi ID lama,
  `ID` di-set `"UnknownId"` supaya prosedur Oracle menerbitkan ID baru.
- `TreatyInDeclineConfirmation_postactEDM` menghapus baris di `M_TREATY_IN_EDM` dan `TREATY_IN_EDM`
  setelah decline.

### 5.x Validasi — dipisah antara yang TERTULIS dan yang BERJALAN

> **KOREKSI BESAR 23 Sep 2026.** Uraian lama di bagian ini mendaftar seluruh aturan validasi apa
> adanya, seolah semuanya berlaku. Setelah tiap langkah diperiksa hidup atau mati, ternyata
> **hampir seluruh daftar itu tidak berjalan.** Kesalahan jenis ini yang paling mahal: ia membuat
> pembaca percaya ada penjaga yang tidak ada, lalu perkiraan pekerjaan menganggap pemulihannya
> gratis. Mulai sekarang seluruh dokumen ini memisahkan keduanya secara tegas.

**Yang BENAR-BENAR BERJALAN di jalur pengajuan hari ini — seluruhnya:**

| Pemeriksaan | Di mana | Akibat bila gagal |
|---|---|---|
| `TreatyIn.Ceding` tidak boleh kosong | `TreatyInCheckError` langkah 3 | `OutputParam.ERRMSG` diisi |
| `TreatyIn.LeadingReinsSource` (SoB) tidak boleh kosong | `TreatyInCheckError` langkah 4 | `OutputParam.ERRMSG` diisi |
| Master sudah punya klaim → revisi ditolak | `CheckDuplicateOffer` langkah hidup terakhir | `Page-Set-Messages`, revisi dibatalkan |

Pemblokirannya nyata: di `TreatyInSubmit`, langkah sesudah pemanggilan berkondisi
`OutputParam.ERRMSG == ""` dengan **WhenFalse = 6 (keluar dari aktivitas)**, sehingga pengajuan
benar-benar berhenti. Tetapi yang diperiksa hanya dua field itu.

**Yang TERTULIS tetapi TIDAK BERJALAN:**

| Aturan yang tertulis | Sebab tidak berjalan |
|---|---|
| `CedingID` + `LeadingReinsSourceID` wajib; `ProportionType` wajib; `CurrencyID` wajib; `ClassOfBusiness` + `ClassOfBusinessID` wajib **di kedua cabang** (bukan hanya proporsional, seperti tertulis di versi lama dokumen ini); untuk proporsional `TreatyType` wajib; untuk non-proporsional `Layer`, `LayerPart`, `LayerType`, `LayerPartType` wajib dan `.Value` tidak boleh kosong atau nol | seluruhnya ada di `TreatyInCheckID`, dan pemanggilan `Call TreatyInCheckID` di `TreatyInSubmit` **ber-blok `//` alias mati** |
| Pesan galat ditempelkan ke field yang salah | kedua langkah `Property-Set-Messages` di `TreatyInCheckError` **mati**; yang hidup hanya mengisi `OutputParam.ERRMSG` |
| Peringatan duplikat Ceding + SoB + periode | **seluruh** logikanya di `CheckDuplicateOffer` mati: penyiapan parameter, pemanggilan `BrowseTREATY_IN`, penyusunan daftar ID serupa, sampai pesan *"Have similarities in SoB, Ceding, and Period, please check for duplicates"* |
| Ceding dalam *Agent Negative List* ditolak | `TreatyInCheckCedingBlacklist` **tidak ada di jalur submit** — ia kontrol layar saja |

Sisa pesan galat *"ID error: Found one or more empty field/ ID value"* masih tertulis di
`TreatyInSubmit` dan menjanjikan pemeriksaan yang jauh lebih luas daripada yang benar-benar
dijalankan.

**Akibat perencanaan.** Deteksi duplikat, daftar wajib-isi, dan penolakan daftar hitam adalah
**kemampuan baru** di sistem hasil migrasi, bukan pelestarian perilaku lama. Jangan dihitung
gratis. Sebaliknya, daftar di `TreatyInCheckID` yang mati berguna sebagai **titik awal** — ia
pernah dianggap benar oleh seseorang di NuRe — dan **Uji S** mengukur apa yang lolos selama ia
mati.

---

## 6. Dua cabang perhitungan

### 6.1 Proporsional (Quota Share / Surplus)

`LimitCalculation` adalah inti:
- **Quota Share**: `RetentionPct = 100 − QSPct`, lalu `CessionPct = 100 − RetentionPct`
  (jadi `CessionPct == QSPct`; perhitungan bolak-balik ini redundan tetapi dipertahankan).
  `RetentionList` dan `CessionList` dihitung per mata uang dari `IOOLimitList`.
- **Surplus**: menelusuri daftar layer dari atas untuk menemukan baris `QUOTA SHARE`,
  memakai limit retensinya sebagai basis, lalu `IOOPct = CessionPct = Surplus × 100` dan
  `IOOLimitList = basis × Surplus` (`Surplus` disimpan sebagai jumlah *line*, bukan persen).
- Ada parameter `autocalculate`; bila mati, aktivitas langsung keluar dan nilai diisi manual.

`CalculateShareList` menghitung porsi NuRe: `RNMShareList(i).Value = Value × RNMShare/100`.
Bila `RNMShareAcrossTheBoard` menyala dipakai `TreatyIn.RNMShareP`, kalau tidak dipakai
`RNMShare` per baris detail.

> **Koreksi 23 September 2026.** Kalimat ini semula menyebut `RNMShareP` sebagai bagian **"global"**.
> Itu keliru: akhiran `P` berarti **Proportional**, bukan global — dibuktikan oleh pasangan kembar
> `TreatyInMappingDataconvert` (menulis `RNMShare`) dan `TreatyInMappingDataconvertProp` (menulis
> `RNMShareP`), pola yang sama dengan tiga pasangan bersufiks `P` lainnya. `RNMShareP` adalah
> **bagian NuRe tingkat kontrak pada cabang proporsional**. Mekanisme yang diuraikan kalimat ini
> tetap benar; yang diperbaiki hanya arti namanya. Rinci di `SPEC-MODEL-DATA.md` §13. `ShareNote` dirakit jadi teks `" of {CessionPct}% of 100%"` untuk
Quota Share, atau `" of 100%"` untuk Surplus.

`TreatyInPropshare` / `TreatyInPropshareDetail` mengorkestrasi: hitung share → penyebaran
(`FetchQSfromMaster` untuk skema lama, `SetSpreadName` untuk skema baru) → total.

### 6.2 Non-proporsional (XOL)

`DetailCalculation` dengan parameter `type`:
- `type=copy` — menyalin nilai EGNPI yang cocok CoB-nya ke dalam layer, lalu mengagregasi per
  mata uang (`TotalEgnpi`).
- `type=adj` — `PremiumEarnedList = EGNPI × AdjRate/100`. Bila parameter `TNOP == "tnop"`
  dikalikan lagi prorata bulanan: `(bulan(Commencement→Termination)+1)/12`. Selisih tanggal
  dihitung dengan menempelkan literal `"T150000.000 GMT"` ke tanggal.
- `type=mdp` — `MDPList = PremiumEarned × MDPPct/100`, `MDPMinList = PremiumEarned × MDPMinPct/100`,
  lalu memanggil `DetailCalculationROL`.

`DetailCalculationROL` menghitung *Rate on Line*: semua premi dan semua limit dikonversi ke IDR
memakai `CurrencyList(i).Conversion`, lalu `ROLPct = TotalPremi / TotalLimit × 100`.
**Temuan TD-02:** bila `TotalLimit == 0`, hasilnya di-set ke angka ajaib `9989998` alih-alih nol
atau error, dan angka itu ikut tersimpan.

`SetReinstatementPct` membuat `Reinstatement_List` sebanyak `ReinstatementValue`, tiap baris
default `ReinstatementPct = 100` dan `AdditionalPct = 100`.
`CalculateReinstatement`: `ReinstatementAmount1/2 = AdditionalPct/100 × Limit/Limit2`.
`CalculateReinstatementPct`: arah sebaliknya, `ReinstatementPct = Amount1 / Limit × 100`.
`ReCalculateReinstatement`: `AdditionalAmount1/2 = ReinstatementPct/100 × MDPList(n).Value`.
**Temuan TD-03:** ketiganya mengasumsikan `MDPList(1)` adalah IDR dan `MDPList(2)` adalah USD —
mata uang lain diam-diam menghasilkan nol.

> **KOREKSI 23 Sep 2026 — dua klaim di sini dicabut.** Kalimat lama menyatakan bahwa
> `SetReinstatementPct` bercabang pada `==1` dan `>2` sehingga **kasus tepat 2 mata uang tidak
> tertangani**, dan bahwa kondisinya tidak konsisten dengan `>1` di `ReCalculateReinstatement`.
> Keduanya salah. Yang tertulis `>2` hanyalah **label langkahnya**; kondisi yang benar-benar
> dijalankan adalah `>=2`, dan memo aturannya sendiri berbunyi *"perbaiki when 3.3"*. Pada ukuran
> daftar, `>=2` dan `>1` identik — tidak ada inkonsistensi dan tidak ada kasus yang bolong.
> Yang tersisa dari TD-03 hanyalah asumsi slot, dan cacatnya berupa **nol diam-diam**, bukan
> aritmetika yang keliru.

**Temuan TD-22 (baru, 23 Sep 2026):** `CalculateReinstatementPct` menyetel
`ReinstatementPct = ReinstatementAmount1 / Limit × 100` — yaitu mengambil rasio **jumlah limit**
yang dipulihkan lalu **menimpa tarif premi** reinstatement dengannya. Untuk layer berketentuan
"satu reinstatement at 50%", angka 50 adalah syarat kontrak yang dinegosiasikan, dan ia lenyap
begitu seseorang mengetik jumlah reinstatement di grid. Cacat ini tak pernah terlihat karena kedua
persen disemai `100`, dan pada nilai itu menimpa yang satu dengan yang lain tidak mengubah apa pun.
Diuji oleh **Uji O**.

---

## 7. Penyebaran (spreading) ke susunan treaty NuRe sendiri

Setiap share NuRe dipecah lagi ke susunan retro internal NuRe, diambil dari tabel
`PROPORTIONALARRG` lewat dua Report Definition berurutan:
`BrowseTreatyArrangement_ParentReinsMasterTrt` (mencari induk, `ParentReinsTypeID = "00"`,
`TreatyDescID = "10001"`, nama mengandung `TRT` atau `ORS`, periode mencakup `Commencement`)
lalu `BrowseTreatyArrangement_Limit_MstTrt_RD` / `_Limit_RD` (mengambil anak-anaknya beserta `Pct`).

Dua ID jenis reasuransi dipakai sebagai konstanta di banyak tempat:
**`10004` = QS (R/I)** dan **`10028` = QS (OR)**; `ORS` (own retention share) diperlakukan sebagai
`QSPCT = 100`. Pembagian selalu `OR = QSPCT/100 × nilai` dan `RI = (100−QSPCT)/100 × nilai`,
diterapkan paralel ke lima besaran: Limit, Gross Premium, Gross Premium Min, Deduction, Net Premium
— itulah asal 10 daftar `RNMSpreadedList…XOL`.

Ada dua generasi implementasi yang hidup bersamaan dan dipilih lewat kondisi runtime:
`FetchQSfromMaster` / `FetchQSfromMasterXOL` (**"Spreading Lama"**) dan
`SetSpreadName` / `SetSpreadingXOL` (**"Spreading Baru"**, mendukung `BreakDownSprdList` bertingkat).
Migrasi harus memutuskan satu di antaranya; mempertahankan keduanya berarti mempertahankan bug.

**Temuan TD-04 (serius).** `FetchQSfromMasterXOL` mengandung cabang `TreatyIn.ID == "1000951"`
yang mengalikan hasil dengan `85/100` dan menyisipkan baris `ORS` 15% dengan `ReinsTypeID = "10007"`
secara literal. Ini tambalan untuk satu kontrak tertentu yang tertanam di kode. Harus dikonfirmasi
ke pengguna bisnis apakah perlakuan khusus ini masih berlaku sebelum data dimigrasikan.

`CalculateSharePctToRetro` + `CalculateORRItoRetro` + `CalculateRetroSumary` + `CalculateRetroTotal`
menangani retro ke pihak luar: `Share(i).RnmLimitList = Limit × RetroPct/100`, komisi
`"Comm to NuRe" = Gross × RetroBrokerage/100`, lalu `Net = Gross − Deduction`.

---

## 8. Addendum / EDM (perubahan kontrak berjalan)

Tiga salinan data hidup berdampingan di satu halaman:

| Cabang | Arti |
|---|---|
| `TreatyIn.*` | nilai estimasi/kontrak berjalan |
| `TreatyIn.ActualValue.*` | nilai aktual hasil pembaruan |
| `TreatyIn.OLDDATA.*` | snapshot sebelum addendum |
| `TreatyIn.ValueDifference.*` | selisih, yang inilah yang dibukukan |
| `TreatyIn.ValueBeforeProrate.*` | salinan selisih sebelum prorata diterapkan |

`TreatyIn.EDMState` bernilai `"1"`, `"2"`, atau `"3"` (when rule `TreatyMasterInEDM`);
`EDMMaterialType` membedakan addendum material dan non-material — beberapa Section
(`DetailEGNPI`, `Installments`) dinonaktifkan untuk tipe non-material.

Dua orkestrator selisih yang nyaris kembar:
- `TreatyInActualUpdateValue` — jalur non-EDM, memanggil 15 aktivitas `TreatyInDifference*`.
- `TreatyEDMCalculateDifference` — jalur EDM, memanggil `TreatyEDMDifference*` (Premium, Limits,
  Share, Deduction) lalu `TreatyEDMProRateCalculation`.

**Perbedaan semantik yang mudah terlewat:** rumpun `TreatyInDifference*` menjepit hasil negatif
ke nol (`@if(.Value < lama, 0, .Value − lama)`, dikomentari "If premi actual < estimate then set to 0"),
sedangkan rumpun `TreatyEDMDifference*` melakukan pengurangan langsung sehingga **boleh negatif**.
Dua modul memakai definisi "selisih" yang berbeda — ini harus ditetapkan eksplisit di sistem baru.
Bahkan di dalam satu aktivitas pun tidak konsisten: pada `TreatyInDifferenceLimits`, `Limit`,
`Deductible`, `AdjRate`, `ROLPct`, `AggregateLimit` dikurangi langsung, tetapi `MDP` dan
`EgnpiTotalList` dijepit ke nol.

**Prorata** (`TreatyCalculateProratePct`):
`ProRateDays = EDMEffective → Termination`, `ProRateTotalDays = Commencement → Termination`,
`ProRatePercent = ProRateDays / ProRateTotalDays × 100`. Bila `IsProRate` menyala,
`TreatyEDMProRateCalculation` mengalikan **seluruh** nilai di `ValueDifference` dengan persen itu —
sekitar 20 blok berulang yang isinya identik (`.Value × ProRatePercent/100`).

**Temuan TD-05.** `TreatyInDifferenceFacShare` dan `TreatyEDMDifferenceShare` menulis hasil ke
`TreatyIn.ActualValue.LimitShareSummaryList`, padahal semua saudaranya menulis ke
`TreatyIn.ValueDifference.*`. Berdasarkan pola di sekitarnya ini tampak sebagai salah ketik yang
menimpa nilai aktual dengan nilai selisih. Perlu dikonfirmasi ke pengguna sebelum ditiru.

---

## 9. Achievement (realisasi premi & klaim)

`GetAchievement` menarik dari `POOLDATA.ACHIEVEMENT` (dicocokkan lewat `SUBSTR(NOOFFER,1,7)`),
lalu menggabungkan klaim dari `GetHistoryClaim_act` (klaim disetujui) dan `GetEstimasiClaim_act`
(estimasi/outstanding). Rumusnya:

```
IncuredClaim = PaidClaim + CASHCALL + OutstandingClaim + OutstandingCASHCALL
Total        = NETPREMIUM − IncuredClaim
LossRatio    = IncuredClaim / NETPREMIUM × 100          (0 bila NETPREMIUM = 0)
AchievementPct = (TotalAchPremium / (RNMShareP/100)) / EPI × 100
```

Nilai klaim dalam IDR dibagi `Conversion` untuk kembali ke mata uang polis.

> **KOREKSI 23 Sep 2026 — TD-06 dibaca ulang, ARAH DAMPAKNYA BERUBAH.**
> Uraian lama di sini berbunyi *"bila kurs tidak ditemukan, `Conversion` dipaksa `1`"*. **Itu
> salah.** Kondisi langkahnya adalah `.Currency == "IDR"` — penyetelan `Conversion = 1` hanya
> berjalan ketika mata uangnya memang rupiah, dan itu benar. Untuk mata uang asing yang kursnya
> tidak ditemukan **tidak ada penggantian sama sekali**: `.Conversion` tinggal kosong, lalu
> `.PREMIUMtoIDR = .PREMIUM * .Conversion` menghasilkan **0** dan `@divide(ClaimIDR, .Conversion)`
> membagi dengan kosong.
>
> Perkiraan dampaknya karena itu berubah arah. Yang lama memperkirakan nilai **menggelembung**
> (valas dihitung 1:1 terhadap rupiah); yang benar adalah nilai **menghilang** — premi dan klaim
> valas lenyap dari agregat pencapaian. Kesalahan yang menggelembungkan mengundang pertanyaan;
> kesalahan yang menghilangkan tidak. Karena itu TD-06 dinaikkan prioritasnya: ia jauh lebih
> mungkin sudah berlangsung lama tanpa ada yang menyadarinya.
Hasilnya diarsipkan ke `POOLDATA.LOG_ACHIEVEMENT` dan bisa diekspor ke CSV/Excel
(`GenerateCSVTreaty` → `MSOGenerateExcelFile`).

---

## 10. Lampiran dokumen

Berkas **tidak** disimpan di database, melainkan di Google Cloud Storage lewat REST:

1. `GetLinkService` membaca URL layanan dari `M_LINK_SERVICE` berdasarkan `KATEGORI_1`/`KATEGORI_2`;
   pemilihan lingkungan (Sandbox/Dev/QA/Staging/Production) ada di System Setting `LinkService`.
2. `GetTokenStorage_SQL` (prosedur `pooldata.GET_TOKEN_STORAGE`) mengambil token.
3. `InsertGoogleStorage_Act` menentukan MIME type lewat Decision Table `GetMimeType`
   (memetakan ekstensi ke ±37 MIME type), mengunggah lewat Connect-REST `ServiceGoogle`,
   lalu mencatat metadata di `T_STORAGE_IMAGE` dengan `IMAGEID = STANDARD_HASH('ASMPP'||timestamp,'MD5')`.
   Folder mengikuti pola `{app}/Doc/{YYYY}/{MM}/`, nama berkas diberi awalan
   `yyyyMMdd-hhmmss-S - `.
4. `GetUrlGoogleStorage_Act` memberi URL bertanda tangan; bila `EXPDATE` sudah lewat, URL
   diperbarui dan `T_STORAGE_IMAGE` di-update. Dokumen Office bisa dibuka lewat
   `https://view.officeapps.live.com/op/view.aspx?src=…`.
5. Relasi dokumen ke treaty ada di `M_ATTACHMENTTREATY_2` (TREATYID, CATEGORY_ID, FILENAME,
   FILEMIMETYPE, T_STORAGE_ID). Kategori diambil dari `M_KATEGORIMASTERTREATY` dan
   `CATEGORY_ATTACH_REAS` (OFFER, QUOTATION/PLACING SLIP, EMAIL, SOA, R/I SLIP, CLAUSES,
   EXTENDWPC, NOTA/INVOICE, OBJECTLIST, OTHERS, PHOTO, SURVEYREPORT).

---

## 11. Penjadwalan: reporting, akumulasi, installment

`TreatyInSetReport` membangun `ReportingPeriodList` dari `ReportingStart`→`ReportingEnd`:
interval 3 bulan untuk `quarter` (label `"Q n"`), 6 untuk `half` (`"H n"`), 1 untuk `month` (`"M n"`),
atau nilai bebas `ReportingInterval` (`"T n"`). Tiap periode menghitung tiga jatuh tempo bertingkat:
`SubmissionDue = akhir periode + ReportingSubmission hari`,
`ConfirmationDue = akhir periode + (Submission + Confirmation) hari`,
`SettlementDue = akhir periode + (Submission + Confirmation + Settlement) hari`.
`TreatyInSetAccountReport` memakai pola sama untuk `AccumulationList`;
`TreatyInAccumulationSetSubDue` menghitung `SubDueDate = ReportDate + SubDays`.

`TreatyInSetValueInstallment` membagi `TotalShareNetNP` menjadi `InstallmentNo` termin,
default `100 / InstallmentNo` persen per termin; `SetTotalInstallment` menjumlahkan
`AmountTotal` dan `PctTotal` untuk validasi bahwa totalnya 100%.

---

## 12. Ringkasan utang teknis

| Kode | Temuan | Dampak migrasi |
|---|---|---|
| TD-01 | Nama approver di-hardcode; `pyTelephone` dipakai sebagai kode peran | wajib diganti tabel peran/penugasan |
| TD-02 | `ROLPct` diisi `9989998` bila limit nol; pengali `100` berada **di luar** percabangan sehingga yang **tersimpan** adalah `998999800`. Ditambah: limit dijumlahkan sekali per mata uang EGNPI (penyebut berlipat), dan di bawah `CurrencyRelation="AND"` limit ditambahkan dua kali dalam satu iterasi | ditutup ADR "kegagalan bukan nilai"; penyebut dihitung sekali. Diukur Uji P. **Uji K dianggap bersih 23 Sep 2026: penanda itu tidak pernah sampai ke data, jadi tidak ada pemetaan data yang perlu dikerjakan saat migrasi — tetapi mekanismenya tetap ada dan tetap dibuang** |
| TD-03 | **(direvisi)** Reinstatement menganggap `MDPList(1)=IDR`, `(2)=USD`; urutan terbalik atau mata uang ketiga menghasilkan **nol diam-diam**. Klaim lama soal ambang `>2` vs `>1` **dicabut** — itu label langkah, kondisi sebenarnya `>=2`, dan tidak ada kasus yang bolong | tulis ulang berbasis lookup mata uang; ditutup ADR "kegagalan bukan nilai". Diukur Uji J |
| TD-04 | `TreatyIn.ID=="1000951"` mendapat faktor 85% + baris ORS 15% hardcoded | **Sudah diputuskan user 22 Sep 2026: JANGAN diberlakukan di aplikasi hasil migrasi.** |
| TD-05 | `TreatyInDifferenceFacShare` menulis ke `ActualValue.*` alih-alih `ValueDifference.*` | **Sudah diputuskan user 22 Sep 2026: JANGAN diberlakukan di aplikasi hasil migrasi**; selisih fac share ditulis ke `ValueDifference.*` seperti saudara-saudaranya |
| TD-06 | **(direvisi)** Kurs valas tak ditemukan → `Conversion` kosong → nilai IDR jadi **0**, bukan 1:1. Dampaknya menghilangkan, bukan menggelembungkan | ditutup ADR "kegagalan bukan nilai": tidak ada nilai pengganti, ketidakadaan membawa sebabnya |
| TD-07 | Selisih kadang dijepit ke nol, kadang tidak — beda antar aktivitas dan bahkan di dalam satu aktivitas | tetapkan satu definisi selisih |
| TD-08 | Dua generasi spreading ("lama" & "baru") hidup bersamaan | pilih satu |
| TD-09 | `Limit`/`Limit2` (IDR/USD implisit) berdampingan dengan page list multi-mata-uang | satukan ke satu representasi |
| TD-10 | Pola hapus-lalu-sisip seluruh detail pada tiap simpan, tanpa transaksi lintas tabel yang terlihat | pakai transaksi dan upsert berbasis kunci |
| TD-11 | 15+ aktivitas berisi logika agregasi multi-mata-uang yang identik; `TreatyEDMProRateCalculation` mengulang blok yang sama ±20 kali | ganti satu fungsi generik |
| TD-12 | Beberapa aktivitas memuat catatan penulisnya sendiri: "status: not completed", "step 7-13 is problematic", "ROL BELUM SELESAI" | perlakukan bagian ini sebagai belum tervalidasi; konfirmasi perilaku yang diinginkan |
| TD-13 | Tombol pengembang masih ada di layar produksi ("1 Add from JsonOffer(dev)", "save new value (dev)", "testagent(dev)") | jangan dibawa |
| TD-14 | Tanda tangan `PEGA_M_TREATY_IN_DETAIL` yang diterima tertinggal 2 parameter (`SPREAD_RNM_SHARE_PCT/VALUE`) dari aturan Pega | ikuti XML; model baru tetap punya kedua field |
| TD-15 | Urutan parameter spreading terbalik antara prosedur EDM dan non-EDM | jangan salin-tempel antar jalur |
| TD-16 | Prosedur detail menelan galat tanpa mengisi `ErrMsg`/`StsSave`; `StsSave := 1` walau tidak ada yang ditulis | galat harus eksplisit |
| TD-17 | `COMMIT` per baris detail — satu penyimpanan treaty bukan satu transaksi | satukan jadi satu transaksi |
| TD-18 | Kunci lampiran dari timestamp milidetik (`M_ATTACHMENTTREATY_2.ID`, `T_STORAGE_IMAGE.IMAGEID`) rawan tabrakan | pakai UUID/sequence |
| TD-19 | `PROPORTIONALARRG` semua kolom `VARCHAR2(1000)` termasuk `PCT`; `TREATYINDETAIL` campur `NUMBER`/`VARCHAR2` | skema baru bertipe tegas |
| TD-20 | Tidak ada PK/indeks pada `TREATYINDETAIL`, `ACHIEVEMENT`, `TREATYINPRODUCTION`, `PROPORTIONALARRG` | rancang indeks sejak awal |
| TD-21 | `replace(DataPega,'UnknownId',id)` menyunting seluruh teks JSON secara buta | jangan tiru; pakai ID dari aplikasi |

---

## 13. Arahan kasar ke arsitektur baru

**Oracle.** Pertahankan pemisahan header/detail, tetapi ganti kolom `JSONDATA` sebagai sumber
kebenaran dengan skema relasional yang tegas; JSON boleh disimpan sebagai arsip muatan masuk saja.
Kandidat tabel: `TREATY`, `TREATY_LAYER`, `TREATY_LAYER_DETAIL`, `TREATY_SHARE`,
`TREATY_SPREADING`, `TREATY_RETRO`, `TREATY_EGNPI`, `TREATY_RETENSI`, `TREATY_INSTALLMENT`,
`TREATY_PERIODE_LAPORAN`, `TREATY_AKUMULASI`, `TREATY_KOMENTAR`, `TREATY_LAMPIRAN`,
`TREATY_NILAI` (tabel nilai uang generik: pemilik, jenis, mata uang, nilai) dan
`TREATY_ADDENDUM` + `TREATY_ADDENDUM_NILAI` untuk cabang OLDDATA/ActualValue/ValueDifference.
Jangan pakai empat cabang salinan di satu baris; pakai versi bernomor.

**Golang.** Satu paket `treaty` dengan: tipe `Money` + agregator multi-mata-uang; kalkulator
proporsional (QS/Surplus); kalkulator non-proporsional (adjustment, MDP, ROL, reinstatement);
mesin spreading (satu generasi saja); mesin selisih addendum dengan kebijakan penjepitan yang
eksplisit sebagai parameter; state machine approval berbasis tabel peran. Kalkulator harus murni
dan teruji terpisah dari HTTP dan DB — inilah bagian yang paling rawan salah.

**React.** Satu halaman besar Pega dipecah jadi rute: Daftar Treaty → Detail Treaty dengan tab
(Informasi, Limits/Layers, Share, Spreading, Retro, EGNPI, Installment, Reporting, Achievement,
Lampiran, Riwayat & Approval). Grid yang bisa diedit inline adalah pola dominan; siapkan satu
komponen tabel editable multi-mata-uang yang dipakai ulang. Hak akses dan kondisi `ViewState` /
`RevisionState` / `IsEditData` menentukan mode baca-saja — petakan itu ke satu hook izin.

---

## 14. Keputusan user (22 September 2026) dan status gap

Seluruh pertanyaan terbuka dari sesi pertama sudah dijawab. Ini yang berlaku sekarang:

| Gap sesi sebelumnya | Status |
|---|---|
| Isi stored procedure | **Tertutup.** Kelima prosedur sudah diterima; lihat bagian 4A. |
| DDL tabel `POOLDATA` | **Tertutup sebagian.** 7 DDL diterima (bagian 4A.5 mencatat sisanya, yang strukturnya sudah bisa disimpulkan dan tidak menghambat). |
| Arti `ReinsTypeID` dan `TreatyDescID` | **Terjawab.** `ReinsTypeID` (`10004`, `10028`, `10007`) adalah **ID tipe reasuransi** dan `TreatyDescID` (`10001`) adalah **ID deskripsi** — keduanya sekadar kunci asing ke tabel master, bukan enum dengan aturan bisnis tersembunyi. Di skema baru cukup jadi referensi ke tabel master tipe reasuransi dan tabel deskripsi treaty; jangan diubah menjadi konstanta di kode. |
| TD-04 (perlakuan khusus kontrak `1000951`) | **Diputuskan: JANGAN diberlakukan** di aplikasi hasil migrasi. |
| TD-05 (penulisan ke `ActualValue.*`) | **Diputuskan: JANGAN diberlakukan** di aplikasi hasil migrasi. |
| Folder `Treaty In Adjustment` | **Ditahan.** User meminta folder ini **tidak dibaca sampai ada instruksi khusus**. |
| Sumber Java yang tidak ter-export | **Ditutup permanen.** User meminta hal ini diabaikan; **acuan tunggal adalah apa yang ada di XML.** Jangan mengangkatnya lagi sebagai pertanyaan atau blocker. |

Pegangan yang berlaku untuk pekerjaan lanjutan: perilaku yang terlihat di XML adalah sumber
kebenaran. Bila prosedur Oracle dan aturan Pega berbeda (lihat TD-14), **XML yang menang** dan
salinan prosedur dianggap tidak mutakhir. Asumsi ditandai, pekerjaan jalan terus.
