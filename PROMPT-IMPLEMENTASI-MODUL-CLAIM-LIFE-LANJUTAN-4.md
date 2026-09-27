# PROMPT — lanjutan 4: **penutupan Claim Life** (paket disahkan, paritas layar dan aksi dari XML) lalu **modul Komite Claim Life** — satu giliran

> **Skill:** ketik `/mattpocock-skills:implement` sebagai manusia, lalu tempel berkas ini **utuh**.
> Brief modul **`PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE.md`** §1–§11 dan brief lanjutan 1–3 berlaku
> **seluruhnya** — terutama **XML menang** *(modul §1.2)*, cara membaca XML *(modul §2)*, dan
> **mekanisme rantai tanpa mengakhiri giliran** *(lanjutan 3 §1)*. Berkas ini menambah: paket
> keputusan yang kini **disahkan** *(§1)*, **paritas layar dan aksi** yang work owner tuntut *(§3)*,
> dan modul berikutnya *(§4)*.
>
> **GILIRAN INI:** Langkah 0 → **Bagian A** *(**A0 audit ulang XML atas yang sudah dibangun** → A1 →
> A2 → A3 → A4)* → **Bagian B** *(Komite Claim Life 00 → 01 → 02 → 03 → 04a → 04b → 05 → 06 → 07 →
> 08 → 09)*. Pesan ke manusia **hanya satu**, di akhir.
>
> **Kata work owner, 27 September 2026, `[DIPUTUSKAN]`:** *"Termasuk yang sudah dibuat tadi: cek
> kembali dari XML-nya; kalau ada yang kurang tepat, perbaiki."* — itulah A0.

---

## 0. KEADAAN AWAL — 27 September 2026, sesudah `a2102f0`

| | Keadaan |
| --- | --- |
| `HEAD` | `a2102f0` *(tiket 13)*; rantai satu giliran: `f610178` 05 → `c7fcd27` 15 → `9b48e32` 07 → `7c17bbf` 08 → `2598778` 09 → `237055b` 10 → `c6f6691` 12 → `b0afeb4` 11 → `a2102f0` 13. Working tree bersih. Nol remote, nol push |
| Uji | vet · vet db · gofmt nol · build · **228 PASS · 0 FAIL · 34 SKIP** *(seluruhnya "ORACLE_DSN belum dikonfigurasi")* · `tsc` · **13** JS · **88** modul |
| Verifikasi independen | 65 berkas +7.864/−127 sejak `65e937e`; angka tereproduksi; nol kebocoran; nol procedure dipanggil; nol `JSONDATA` dibaca; klaim review *(`PeriksaSatuMataUang` dua kolom, antrean, `IS_CHECK`)* terbaca di kode |
| AC modul | **123 / 217** tertutup *(01 4/14 · 02 7/26 · 03 11/28 · 04 8/8 · 05 8/8 · 06 4/5 · 07 4/6 · 08 3/7 · 09 1/6 · 10 12/18 · 11 6/8 · 12 8/12 · 13 2/6 · 14 40/53 · 15 5/12)*; **94 terbuka**, hampir seluruhnya oleh gerbang §1, Oracle *(G1)*, dan layar yang belum ada *(§3)* |
| Yang membuat work owner tidak puas | rute HTTP **9**, halaman React **2**, kontrol **5** — sedangkan Pega Claim Life punya **3 harness, 20 section, 16 flow action, 31 aksi layar** *(§3)*; lima pintu tulis menjawab **501** karena paket §1 dulu `[USULAN]` |

---

## 1. PAKET KEPUTUSAN — DISAHKAN

**PAKET MODUL: `[DIPUTUSKAN — 27 September 2026]`** — ditafsirkan asisten dari instruksi work owner
*"jangan banyak yang di-skip; ambil logic yang benar dari XML, termasuk button atau method yang
dipakai"*, diberikan **sesudah** tiga kali dijelaskan bahwa tanpa paket ini tiket 09–11 berhenti di
antarmuka. Work owner **mencabut** dengan mengganti kata ini kembali menjadi `[USULAN]` sebelum menempel.

| | Keputusan | Dasar | Yang dikerjakan |
| ---: | --- | --- | --- |
| **aj** | migrasi kolom bank `BRANCH_OF_BANK`, `SWIFT_CODE`, `PAYABLE_TO` di `T_CLAIMLF_ADJUSTMENT` | XML `AdjustmentDetail_Section` enam field bank | A1 |
| **am** | tabel jejak audit `T_CLAIMLF_JEJAK` | tiket 09 AC *(spec, ADR-0007 — keputusan work owner yang sudah ada)*; bentuknya milik executor | A1 → mengganti `JejakBelumDiputuskan` di 04, 05, 08, 09, 10, 11 |
| **af** | `T_GENERAL_KOMITE`, `T_KOMITE_KOMITELIST`, `ALTER T_WORK_CLAIM ADD COVER_KEY` dibuat **di sini** memakai bentuk **tiket 00 Komite Claim Life** §"Bentuk yang dibangun" *(shared PK teks `KMT-xxxxxx`; `ADJUSTMENT_ID` di `T_GENERAL_KOMITE`; `KOMITE_LOOP`, `KOMITE_COUNT`, `ACCEPT_STATUS`; daftar: `KOMITE_URUT`, `KOMITE_ID`, `ID_KOMITE`, `KOMITE_EMAIL`, `KOMITE_APROVAL` 0/1/2, `KOMITE_COMMENT`, `DATE_APPROVE DATE`; FK dan `COVER_KEY` ber-index; `T_KOMITE_KOMITELIST.ID` dari sequence)* | tiket 14 relasi 8·9·10 + tiket 00 Komite | A1 → mengganti `KasusKomiteBelumDiputuskan`, `RosterBelumDiputuskan` *(10, 11)*; tiket 00 Komite lalu **memverifikasi**, bukan membuat ulang |
| **al** | tabel diagnosa per peserta `T_CLAIMLF_DIAGNOSE` — **hanya bila** bab pembacaan XML tiket 08 membuktikan `MedicalCheckClaimLife`/`Diagnose_Section` mengisi `.DiagnoseList` *(kolom dari properti yang diisi, path + baris)* | XML | A1 bersyarat; lalu pencerminan `SetSTS_Reject` di 04 |
| **o1** | penomoran klaim: logika `PROC_GENERATE_SEQUENCE_NUMBER` **ditulis ulang di Go** *(`PenomorCounter`)* dari `SUMBER-PENOMORAN-DBA.md` — termasuk `FOR UPDATE` atas `(CLASS, JENIS, TAHUN)` dalam transaksi pendaftaran dan perakitan format `<prefix>K<kode>.MM.YYYY.<5 digit>` yang di Pega dirakit pemanggil | kata work owner sebelumnya: *"jangan ada lagi pemanggilan procedure, segala procedure hardcode dalam skrip"* | A2 → `PenomorBelumDiputuskan` diganti; AC 2, 3, 7–11 tiket 02; pengurai nomor untuk migrasi 13 |
| **o2, o3** | AC tiket 02 yang menyebut procedure diralat *(modul §1.2-b)*; **ADR-U-0043** ditulis: *penomoran di aplikasi, meng-supersede ADR-U-0006*, mengutip kata work owner di atas | sama | A2 |
| **an** *(baru)* | token storage: logika `GET_TOKEN_STORAGE` **ditulis ulang di Go** — `[data DBA]` 56 baris, nol `COMMIT`: bila `APPNAME` kosong → galat; ambil `KODEAKSES` terbaru dari `POOLDATA.GCP_IMAGE` *(`APPNAME VARCHAR2(20)`, `KODEAKSES VARCHAR2(100)`, `USERINPUT VARCHAR2(50)`, `INPUTDATE DATE`; 2 baris di DEV)* dengan `INPUTDATE > SYSDATE`; bila tidak ada → token baru = `RAWTOHEX(STANDARD_HASH('ASMAPP' ‖ <garam rahasia 27 karakter> ‖ TO_CHAR(SYSTIMESTAMP,'YYYYMMDDHH24MISSFF3'), 'MD5'))`, disimpan dengan `INPUTDATE = SYSDATE + 1 menit`, `USERINPUT = NVL(masukan,'Job')` | prinsip **o** | A2: `TokenStorage` di `repository` menulis `GCP_IMAGE` *(tulisan warisan yang **disengaja**, dicatat)*; ⛔ **garam TIDAK disalin** dari procedure ke mana pun — env `STORAGE_TOKEN_SALT`, diminta ke DBA lewat jalur rahasia; kosong → gagal terang; test memakai garam palsu |
| — | baca `M_LINK_SERVICE` *(`URL`, `KATEGORI_1`, `KATEGORI_2`, `USERNAME`; 18 baris)* saat jalan | **sudah diputuskan ADR-U-0013**; bukan persetujuan baru | A2: `ResolverBelumDiputuskan` diganti pembaca; URL/USERNAME **tidak pernah** masuk log/artefak |
| — | roster Komite dari `EMAILKOMITE` *(18 kolom: `NAME`, `EMAIL`, `LIMIT_BOTTOM`, `LIMIT_TOP`, `STS_KLAIM`, `TYPE_BUSINESS`, `TYPE_KOMITE`, `JABATAN`, …; 19 baris — **data orang**)* | `FilterEmailKomiteWithLimit` | A2: dibaca saat jalan saja; fixture `UJI-*`; nol baris disalin |
| **ao** *(baru)* | posisi tahap saat cutover *(tiket 13)*: sumber warisan tidak memuatnya → aturan turunan: klaim yang masih punya baris `STS_REJECT = 0` → `PY_POSITION = Outstanding Claim`; yang seluruh barisnya final → `NULL` *(selesai)*; Medical Check / Claim Analis **tidak dapat direkonstruksi** dan dicatat per klaim di laporan rekonsiliasi | `[USULAN]` — rekomendasi setujui | A4 |
| **ag**, **ah**, **ak** | tetap *(modul berikutnya Master Product Name Life; Product+UW)* | | — |

Tiap tabel/kolom baru = langkah migrasi bernomor baru berurutan *(`011`…)* + `_down` + STRUKTUR bertanggal
+ satu baris di tiket 14 + penjaga cacah `CREATE`/`ALTER` **diperbarui dengan angkanya**, bukan
dilonggarkan. Migrasi tetap **belum pernah dijalankan** di Oracle mana pun *(G1)*.

---

## 2. LANGKAH 0

`git add PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE-LANJUTAN-4.md` → commit `docs: brief lanjutan 4 — paket
disahkan, paritas layar dan aksi, modul Komite` → `git status --porcelain` kosong → uji hijau *(228 ·
34 SKIP · 13 JS · 88 modul)* → SHA = titik tetap **A1**.

---

## 3. BAGIAN A — PENUTUPAN CLAIM LIFE, lima paket kerja (A0–A4), satu commit per paket atau per kelompok

Aturan rantai lanjutan 3 §1 berlaku per paket: `claimed` di tiket yang tersentuh → XML → test → kode →
verifikasi → bab Implementasi *(≤ 50 baris, menyebut AC mana yang berpindah dari terbuka ke tertutup
dan **sebabnya**)* → `/code-review` → perbaiki → verifikasi → commit `claim-life: penutupan A<n> — <isi>`.

### A0 — Audit ulang XML atas SELURUH yang sudah dibangun (tiket 01–15), perbaiki yang kurang tepat

**Yang diaudit:** setiap fungsi di `internal/services`, `internal/repository`, `internal/handlers`,
dan setiap halaman/komponen di `frontend/src` yang meniru rule Pega *(komentar kepalanya menyebut rule
dan langkah)*, ditambah setiap AC yang tercentang `[x]` di tiket 01–15. Rule yang **dirujuk** kode
tetapi belum pernah dicatat di bab pembacaan tiket mana pun **dibaca utuh sekarang**.

**Cara audit, per rule** *(modul §2; berkas pecahan dan nomor barisnya dicatat)*: baca **blok langkah
utuh**, bukan hasil `grep` — lalu bandingkan dengan kode pada sembilan sisi:

| # | Sisi | Yang dibandingkan |
| ---: | --- | --- |
| 1 | rumus dan pembulatan | tiap `@divide(...,n)`, `@toDecimal`, `@round`, urutan operasi |
| 2 | precondition **beserta aksinya** | `WhenTrue`/`WhenFalse` `2` = lanjut, `3` = lewati; cabang mana yang berjalan untuk nilai mana |
| 3 | sumber nilai | `pyStepsObjectName` / halaman langkah: peserta, baris adjustment, header, `pyWorkPage` — kolom kita yang mana |
| 4 | penulis kolom | sensus `<PropertiesName>.*KOLOM` *(pola dengan `.*`, bukan `[^<]*`)* di seluruh `Activity\` kedua modul untuk setiap kolom yang kode kita tulis |
| 5 | gerbang layar | setiap `<pyCondition>` section/harness → gerbang `services` **dan** kontrol React |
| 6 | pesan | kalimat persis XML; data tambahan di medan terpisah |
| 7 | himpunan field | properti yang section tampilkan/isi vs medan model dan layar kita |
| 8 | urutan langkah dan transaksi | apa yang ditulis sebelum apa; apa yang dibatalkan bila gagal |
| 9 | jalur warisan | Connect-SQL yang activity panggil sesudahnya *(`UpdateOsAkseptasiClaimLife_sql`, `InsertJsonKlaimLife_sql`, dll.)* — apakah jalur tulis warisan kita menulis kolom yang sama |

**Keluaran:** `.scratch/claim-life/AUDIT-XML-ULANG.md` — satu baris per rule: rule *(path)* · kode yang
menirunya · tiket/AC · verdict **tepat** / **kurang tepat** *(apa persisnya, bukti baris)* /
**belum ditiru** · commit perbaikan. Perbaikan = kode + test + ralat tiket *(modul §1.2-b)*, satu
commit per kelompok `claim-life: audit XML — <kelompok>`. AC yang ternyata tercentang di atas bacaan
yang salah **dicabut** centangnya sampai kodenya benar. Verdict "tepat" pun wajib membawa bukti baris
— tanpa bukti ia bukan audit.

**Titik yang sudah dicurigai asisten selama verifikasi — wajib masuk audit, dijawab satu per satu:**

| # | Titik | Rule | Yang dicek |
| ---: | --- | --- | --- |
| 1 | pemeriksaan urutan kaskade memakai kapasitas mata uang klaim, sumber mengurutkan `TO_NUMBER(IDR)` untuk kedua mata uang | `GetJsonProductLife` | `HitungSpreading` — periksa pada `IDR`, kasus USD ditest |
| 2 | jalur akseptasi **tanpa Komite**: `.AdjustmentList(<LAST>).STS_REJECT = 1` + `ACCEPTEDNO` + `ACCEPTATION_DATE` di Claim Life sendiri | `SaveAdjustment_Act`, flow action `AkseptasiClaimLife`, `InputAkseptasiClaimLife` | gerbangnya *(peran? `Type`? limit?)*; tiket 04/08 diralat; `Status.Ubah` mendapat jalur Aksep dari Claim Life |
| 3 | tulisan datar warisan sesudah reject | `RejectOSClaimLife_Act` langkah 3–7 → `UpdateOsAkseptasiClaimLife_sql` | kolom yang ditulis vs `BarisLamaDari`; nama orang **tidak** disalin ke tabel baru |
| 4 | dua activity pengembalian | `SendtoAdmin_Act` **dan** `SendtoAdmin_Act1`, `SendtoMedical_Act` | apa beda keduanya; mana yang dipanggil dari mana; `tahap.go` |
| 5 | tanggal kejadian: siapa **penulis**-nya | `UpdateDateClaimLife_Act` vs `ValidasiDOL_Act` *(validasi saja)* | rute `PUT …/tanggal-kejadian` meniru penulis yang benar, termasuk kolom lain yang ikut ditulis |
| 6 | field pendaftaran | `InputRegisterClaimLife` *(section + flow action)*, `Register_Flow` | himpunan field vs `PermintaanDaftar` *(hari ini: `PL_NUMBER`, `Type`, `KodeBisnis`, `MataUang`, `Sertifikat`)* — yang kurang ditambah, yang dikarang dibuang |
| 7 | pemilihan peserta dan tulisnya | `SavePesertaClaim`, `SaveInsuredClaim_Act`, `SelectAllClaimLife_act`, `DeletePesertaClaimLife`, `LoadDataPesertaSpesifik_Act` | `Daftar` dan `IS_CHECK`; hapus peserta *(bukan hapus klaim)* — tiket 15 dan 02 |
| 8 | peran per tahap | `Register_Flow` *(`pyPosition` tiap `Assignment`)*, tiga `<pyCondition>` section | `wewenang.go`, `tahap.go`, `PeranSimpanOutstanding` |
| 9 | `ContentNote` di luar `DEATH` | `SaveOutStandingLife_Act` jalur `!= "DEATH"` *(baris 4054, 4199)* | apa yang berbeda untuk `HEALTH`/`CI`/`TPD`/`TI`; sudah ditiru? |
| 10 | dokumen: 7 kolom diadopsi dari 12 properti | `InsertDocument_Act`, `DocumentLife`, `AttachDocScreenLife` | kolom yang dibuang benar-benar tidak dipakai layar/gerbang |
| 11 | muatan penyerahan Komite dan roster | `CreateKMTLife_Act` sepuluh langkah, `FilterEmailKomiteWithLimit` | medan `PenyerahanKomite`, filter roster *(`LIMIT_BOTTOM`, `STS_KLAIM`, …)* |
| 12 | efek keluar | `SendEmailKlaimLF`, `serviceInsertArasapasClaimLife_act`, `InsertGoogleStorage_Act` | penerima/isi email, muatan Arasapas, langkah storage — vs antarmuka `EfekKeluar` |
| 13 | hapus klaim | seluruh korpus: `Delete*`, `Obj-Delete`, `pxDelete` | kesimpulan tiket 15 *("penolakan hapus fisik")* berdiri di atas bukti apa; `DeletePesertaClaimLife` ditiru? |
| 14 | `IsCheck` | `SetIndexAdjustmentList` *(`true`)*, `RejectOSClaimLife_Act` *(`false`)*, `KomitePostAdjustment` | pulih/tidaknya `IS_CHECK` di seluruh jalur *(temuan review tiket 11)* |
| 15 | dua asumsi pustaka Pega | `@CompareDates`, `@addCalendar` | tetap `[dugaan]` — audit hanya memastikan keduanya konstanta bernama dan tabel batasnya lengkap |

A0 **mendahului** A2–A4: logika yang salah tidak boleh dibangun di atasnya. Bila audit menemukan hal
yang membutuhkan keputusan di luar §1, tandai `[terbuka — pemilik]` dan lanjut.

### A1 — Migrasi paket §1 *(aj, am, af, al bersyarat)*

DDL dari sumber yang disebut §1; `Bongkar`/skema uji ikut; `TestSeluruhCreateDapatDibacaNamanya`,
`TestKolomDDLCocokDenganStruktur`, `KolomAlterTambah` diperbarui; tiket 14 satu baris per langkah.
Nol logika di paket ini.

### A2 — Mengganti setiap `…BelumDiputuskan` yang paketnya sudah sah

`Jejak` *(am)* → `PerekamJejak` di `repository`, dipanggil **di dalam** transaksi oleh setiap transisi
dan jalur balik *(04, 05, 08, 11)*; `KasusKomite` + `Roster` *(af)* → penulis `T_WORK_CLAIM` `KMT-`,
`T_GENERAL_KOMITE`, `T_KOMITE_KOMITELIST`, dan pembaca roster `EMAILKOMITE` mengikuti filter
`FilterEmailKomiteWithLimit` *(baca ReportDefinition-nya utuh)*; `Resolver` → pembaca `M_LINK_SERVICE`;
`Penomor` *(o1)* → `PenomorCounter`; `TokenStorage` *(an)*. Kelima pintu yang menjawab **501**
*(`hapus.go`, `komite.go`, `putaran.go`, `register.go`, `tolak.go`)* menjawab **sungguhan**. AC yang
terbuka karena stub-stub itu ditutup di tiket masing-masing **dengan sebabnya**; yang tetap terbuka
hanya karena Oracle diberi tanda `[menunggu G1]`.

### A3 — Paritas layar dan aksi Claim Life dari XML — inilah yang work owner tuntut

**Sensus korpus `[terverifikasi — grep `<pyActivity>` atas `Section\` + `Harness\`]`:** 3 harness
*(`Committe_Life`, `Diagnose_Harness`, `SearchPolicy_Harness`)*, 20 section, 16 flow action, **31
aksi layar** berbeda:

```
NextPrev · CheckTotalAdjustmentClaim · SearchDiagnose_act · setDetailClaim_act · ValidasiDOL_Act ·
SetDisease · SendtoAdmin_Act1 · SearchPolicyHolder_act · DownloadDocumentClaim · CreateKMTLife_Act ·
CountClaimAmountLife_Act · ValidasiClaimReceived_Act · setVisibility_Act · ValidasiSTNC_Act ·
UpdateDateClaimLife_Act · SetIndexAdjustmentList · SetClaimXOL_Act · SendtoMedical_Act · SendtoAdmin_Act ·
SelectAllClaimLife_act · SavePesertaClaim · SaveOutstandingLife_Act · SaveOutStandingLife_Act ·
SaveInsuredClaim_Act · SaveAdjustment_Act · RejectOSClaimLife_Act · ProtectCloseClaim_act ·
LoadDataPeserta_Act · LoadDataPesertaSpesifik_Act · GetListKomiteLife · DeletePesertaClaimLife
```

Flow action *(= layar per tahap)*: `InputRegisterClaimLife` · `OSClaimLife` · `Adjustment_Detail` ·
`RetroClaimLife` · `AttachDocumentLife` · `UploadCSV_ClaimLife` · `ConfirmDeleteAttachment` ·
`MedicalCheck` · `SendtoMedical` · `SendtoAdmin` · `AkseptasiClaimLife` · `RejectOSClaimLife` ·
`CloseClaim` · `ShowEditClaimLife` · `ViewClaimDetailLifeGCNM` · `PL_DetailAction_ViewPolis`.

**Aturan paritas `[DIPUTUSKAN — work owner]`:** setiap **flow action** = satu layar React; setiap
**section** yang dipakainya = komponen dengan **himpunan field persis** dari XML *(label boleh
Indonesia, field tidak dikarang, tidak dihilangkan)*; setiap **`<pyActivity>`** yang dipicu tombol =
satu metode `services` + satu rute HTTP + satu kontrol; setiap **`<pyCondition>`** = gerbang di
`services` **dan** kontrol tersembunyi; setiap **precondition/when** di activity = cabang yang ditest.
Yang boleh "tidak ditiru" hanya mekanik UI Pega *(`NextPrev`, `setVisibility_Act`, `setDetailClaim_act`
bila hanya menata halaman)* — dinyatakan per aksi dengan bukti bahwa ia tidak menulis data.

**Langkahnya:** *(1)* tulis **`.scratch/claim-life/PARITAS-LAYAR-DAN-AKSI.md`** — tabel 31 aksi + 16
flow action: XML *(path)* · yang dilakukannya *(dari langkah activity)* · `services` · rute · kontrol
React · keadaan *(ada / dibuat di A3 / tidak ditiru + bukti)*; *(2)* kerjakan yang **belum ada**,
berkelompok per tahap dan **satu commit per kelompok**:

| Kelompok | Layar / aksi | Catatan XML yang sudah diketahui |
| --- | --- | --- |
| Register | `InputRegisterClaimLife`, `SearchPolicy_Harness/Section`, `SearchPolicyHolder_act`, `LoadDataPeserta_Act`, `LoadDataPesertaSpesifik_Act`, `SavePesertaClaim`, `SaveInsuredClaim_Act`, `SelectAllClaimLife_act`, `DeletePesertaClaimLife`, `UpdateDateClaimLife_Act`, `ValidasiDOL_Act`/`ValidasiSTNC_Act`, `CountClaimAmountLife_Act`, `CheckTotalAdjustmentClaim`, `ValidasiClaimReceived_Act`, `UploadCSV_ClaimLife` | `LoadDataPeserta_Act` menyalin `Type` *(tiket 07)*; `SavePesertaClaim` bukti induk adjustment = peserta *(tiket 14)*; nama orang tampil di layar dari tabel polis **hidup**, tidak disalin ke tabel klaim |
| Outstanding | `OSClaimLife`, `InputOSClaimLife`, `Adjustment_Detail`, `AdjustmentDetail_Section`, `SaveOutStandingLife_Act`, `SetIndexAdjustmentList`, `RejectOSClaimLife`, `RetroClaimLife`, `RetroDetailClaimLife`, `SetClaimXOL_Act` | gerbang peran tiga *(tiket 07)*; enam field bank *(aj)*; spreading *(tiket 03; sumber rate tetap **ag**)* |
| Dokumen | `AttachDocumentLife`, `AttachDocScreenLife`, `DocumentLife`, `DownloadDocumentClaim`, `ConfirmDeleteAttachment` | token *(an)*, storage lewat antarmuka efek keluar *(tiket 12)*; kategori wajib `[terbuka — DBA]` |
| Medis | `MedicalCheck`, `MedicalCheckClaimLife`, `Diagnose_Harness/Section`, `SearchDiagnose_act`, `SetDisease`, `SendtoMedical`, `SendtoAdmin`, `SendtoAdmin_Act`/`_Act1` | bukti **al**; `IsSendtoAdmin`/`IsSendtoMedical` *(tiket 08)* |
| Akseptasi | `AkseptasiClaimLife`, `InputAkseptasiClaimLife`, **`SaveAdjustment_Act`** | ⚠️ `SaveAdjustment_Act` menulis `.AdjustmentList(<LAST>).STS_REJECT = 1` + `ACCEPTEDNO` + `ACCEPTATION_DATE` **di Claim Life sendiri** *(sensus tiket 04)* — jalur akseptasi **tanpa Komite** yang tiket 08 *("tahap 2–3 tidak mengubah status")* tidak kenal. Baca flow action + section + activity utuh, temukan **gerbangnya** *(limit? `Type`? peran?)*, lalu **XML menang**: ralat tiket 08/04 dengan bukti |
| Komite | `ClaimComite`, `Committe_Life`, `CreateKMTLife_Act`, `GetListKomiteLife` | tiket 10 + **af** |
| Tutup & lihat | `CloseClaim`, `CloseClaim_Section`, `ProtectCloseClaim_act`, `ViewClaimDetailLifeGCNM`, `ClaimLifeDetailGCNM`, `ShowEditClaimLife`, `EditDateClaimLife_Section`, `PL_DetailAction_ViewPolis`, `DetailPolisLife`, `PL_DetailViewPolis_Sec` | `ProtectCloseClaim_act` hanya **membaca** `STS_REJECT` *(sensus)*; tutup klaim = tahap akhir yang tiket 04 turunkan |

Frontend: satu halaman per flow action di `frontend/src/pages/`, komponen per section, `api.ts`
satu fungsi per rute, uang dan kode teks, test JS per gerbang tampil/sembunyi. Backend: `handlers`
tipis, aturan di `services`, SQL di `repository`, `Qualify` dan `PeriksaSQL` tetap.

### A4 — Migrasi data 13 dengan **o1** dan **ao**

Pengurai nomor klaim warisan dari format o1; posisi tahap menurut **ao**; rekonsiliasi per klaim;
seluruhnya di skema uji *(SKIP tanpa Oracle)*. ⛔ `POOLDATA` tidak disentuh.

---

## 4. BAGIAN B — MODUL KOMITE CLAIM LIFE, seluruhnya

**Korpus:** `D:\XML\RNM_BRD\Komite Claim Life\` — 18 activity, 14 RDBList, 3 section, 3 flow action,
1 flow, 1 decision table, 3 report definition, 2 when, 1 ConnectREST; aksi layar: `SetRemarkKomiteLife`,
`DownloadDocumentClaim`. **Tiket:** `.scratch/komite-claim-life/issues/` — 11 tiket, seluruhnya
`ready-for-agent`, 0 AC tertutup; spec dan struktur: `spec.md`, `STRUKTUR-TABEL-KOMITE-CLAIM-LIFE.md`,
`keputusan-struktur-komite.md` *(baca utuh sebelum tiket pertama)*.

**Urutan dari kolom Blocked by:** **00** *(skema — memverifikasi tabel yang A1 buat, menambah yang
kurang: sequence `T_KOMITE_KOMITELIST`, index, penjaga)* → **01** *(terima kasus + inbox per posisi)* →
**02** *(mesin tangga)* → **03** *(wewenang + eskalasi)* → **04a** *(nomor akseptasi sekali di tingkat
final — **o1** berlaku: `JENIS` Komite dari `SUMBER-PENOMORAN-DBA.md`)* → **04b** *(rekam akseptasi satu
jalur simpan)* → **05** *(jalur balik `STS_REJECT` dua tingkat — penulis yang tiket 11 Claim Life
baca)* → **06** *(transactional outbox)* → **07** *(worker retry anti-dobel)* → **08** *(status perlu
intervensi + laporan harian)* → **09** *(riwayat tangga di UI)*.

**Aturan yang sama persis** dengan Claim Life: per tiket *(claimed → baca ulang XML dengan bab
`## Pembacaan ulang XML` → test di seam → kode → verifikasi → Implementasi → review → commit
`komite-claim-life: tiket NN — <judul>`)*, XML menang dan tiket diralat dengan bukti, paritas layar
dan aksi *(`PARITAS-LAYAR-DAN-AKSI.md` milik modul ini)*, rantai tanpa mengakhiri giliran, decision
table dibaca sebagai **tabel data** *(nol `if` bercabang)*, `KomitePostAdjustment` sudah tersensus
*(`1` ×3, `2` ×1 pada baris **dan** peserta; `ACCEPTATION_DATE` ×3)*. Kode Go modul ini hidup di
`APP_RNM` yang sama *(`internal/services/komite*`, rute `/api/komite-life/...`)*; tabel milik A1.
Pagar keamanan sama: `EMAILKOMITE` berisi nama dan email orang — dibaca saat jalan, **tidak pernah**
disalin ke fixture, tiket, atau log.

---

## 5. LAPORAN AKHIR GILIRAN

Persis lanjutan 2 §6, dengan tiga tabel: **A0** *(rule diaudit · tepat · kurang tepat dan diperbaiki
· belum ditiru · AC dicabut/dipulihkan)*, **Bagian A** *(per paket: AC yang berpindah tertutup per
tiket, aksi layar dibuat / tidak ditiru dengan bukti, rute dan halaman sesudahnya)* dan **Bagian B**
*(per tiket Komite)*; keputusan §1 yang dipakai; terbuka menurut pemilik *(DBA: G1, garam token, kategori
dokumen wajib, `treatyyear_life` ber-`IDR`/`USD`; Product+UW: ah, ak, satuan DOL; work owner: ao bila
belum, ag)*; telemetri modul §10 + cacah aksi layar *(31 → berapa yang punya rute dan kontrol)*.

---

## 6. PERSETUJUAN MANUSIA

Modul §9. Yang **sudah** disahkan §1 tidak ditanyakan lagi. Yang tetap butuh manusia: `git push`,
migrasi ke skema mana pun selain user kosong DBA, garam token *(jalur rahasia DBA — bukan lewat repo
atau chat)*, endpoint storage/email/Arasapas **nyata**, menulis ke tabel warisan selain `GCP_IMAGE`
*(an)* dan baris datar yang sudah disepakati.

---

*Disusun 27 September 2026 sesudah verifikasi independen `a2102f0` (uji dijalankan ulang, 65 berkas
dibaca statistiknya, klaim review dicek di kode), sensus aksi layar dan flow action dari korpus,
pembacaan struktur `GET_TOKEN_STORAGE` (literal rahasianya disamarkan, tidak disalin), katalog
`M_LINK_SERVICE`, `EMAILKOMITE`, `GCP_IMAGE` (agregat), dan tiket 00 modul Komite.*
