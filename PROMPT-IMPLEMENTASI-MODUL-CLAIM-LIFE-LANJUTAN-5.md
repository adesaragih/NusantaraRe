# PROMPT — lanjutan 5: sisa A2 → A3 paritas layar → A4 migrasi data → **Bagian B Komite Claim Life**, satu giliran

> **Skill:** ketik `/mattpocock-skills:implement` sebagai manusia, lalu tempel berkas ini **utuh**.
> Brief modul **`PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE.md`**, brief lanjutan 1–3, dan **brief lanjutan 4
> bab 0–8** berlaku **seluruhnya** *(bab 7 = jawaban XML soal Save Adjustment dan cara membaca pohon
> XML; bab 8 = ralat `@addCalendar` dan dua catatan kecil)*. Berkas ini memuat keadaan sesudah
> `f8bbf7a`, **daftar sisa yang persis**, dan dua fakta katalog baru yang membuka dua gerbang yang
> selama ini "menunggu DBA" *(§2 **ar**, §2 resolver)*.
>
> **GILIRAN INI:** Langkah 0 → sisa **A2** *(§2)* → **A3** *(§3)* → **A4** *(§4)* → **Bagian B** *(§5)*.
> Rantai tanpa mengakhiri giliran *(lanjutan 3 §1)*; pertanyaan yang XML jawab **tidak** dibawa ke
> manusia *(lanjutan 4 bab 7)*; pesan ke manusia **satu**, di akhir.

---

## 0. KEADAAN AWAL — 27 September 2026, sesudah `f8bbf7a`

| | Keadaan |
| --- | --- |
| `HEAD` | `f8bbf7a` *penutupan A2 — penomoran di aplikasi (o1, ADR-U-0043)*; sebelumnya `5011ee9` *(jejak audit sungguhan)*, `f336f7a` *(A1 migrasi paket)*, `08c682f`, `a67b24e`, `9aca1a7` *(A0 kelompok 1–3)*. Working tree: **`PROMPT-…-LANJUTAN-4.md` berubah** *(bab 7–8)* + berkas ini → Langkah 0 |
| Uji | vet · vet db · gofmt nol · build · **239 PASS · 0 FAIL · 34 SKIP** · **16** JS · **88** modul |
| Skema | migrasi **14** *(`011` bank · `012` `T_CLAIMLF_JEJAK` · `013` `T_GENERAL_KOMITE` + `T_KOMITE_KOMITELIST` + `SEQ_KOMITE_KOMITELIST` · `014` `BUSINESS_CODE`)*; `COVER_KEY` sudah di `001` |
| Permukaan | rute HTTP **11**; halaman React **2** *(`KlaimLife`, `RegisterKlaim`)*; `api.ts` 17 fungsi; 1 berkas test JS; `PARITAS-LAYAR-DAN-AKSI.md` **belum ada** |
| Stub yang masih hidup di kode | `RosterBelumDiputuskan`, `KasusKomiteBelumDiputuskan`, `ResolverBelumDiputuskan`, `AntreanBelumDiputuskan`, `KategoriWajibBelumDiketahui` *(+ `PenomorBelumDiputuskan`, `JejakBelumDiputuskan` yang sudah tergantikan — dibuang atau dipindah ke `_test.go`)* |
| AC | 123/217 pada teks tiket *(ralat A0 belum menggeser centang; A2–A3 yang menggesernya)*; Komite: 11 tiket `ready-for-agent`, 0 AC |

---

## 1. LANGKAH 0

`git add PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE-LANJUTAN-4.md PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE-LANJUTAN-5.md`
→ commit `docs: brief lanjutan 4 bab 7–8 (jawaban XML, ralat addCalendar) + brief lanjutan 5` →
`git status --porcelain` kosong → uji hijau *(239 · 34 SKIP · 16 JS · 88 modul)* → SHA = titik tetap
**A2 sisa**.

---

## 2. SISA A2 — satu commit per stub yang diganti, `claim-life: penutupan A2 — <stub>`

| Stub | Diganti dengan | Sumber XML / katalog `[data DBA — agregat]` |
| --- | --- | --- |
| `RosterBelumDiputuskan` *(af)* | pembaca `EMAILKOMITE` *(`ID`, `NAME`, `EMAIL`, `LIMIT_BOTTOM`, `LIMIT_TOP`, `DEGREE`, `STS_AKTIF`, `OPERATOR_ID`, `STS_REJECT`, `STS_ADJ`, `TYPE_BUSINESS`, `TYPE_KOMITE`, `STS_REG`, `STS_SURVEY`, `STS_SALVAGE`, `STS_ADJUSTER`, `STS_KLAIM`, `JABATAN`; 19 baris)* dengan **filter persis** `ReportDefinition\FilterEmailKomiteWithLimit.xml` *(baca utuh: kolom, filter, parameter `LIMIT_BOTTOM`, `STS_KLAIM`, urutan)*; tiruan skema uji berisi baris `UJI-*` saja — **nol** nama/email nyata di fixture, tiket, log | RD + katalog |
| `KasusKomiteBelumDiputuskan` *(af)* | penulis kasus Komite mengikuti `Activity\CreateKMTLife_Act.xml` **sepuluh langkah utuh**: baris `T_WORK_CLAIM` `KMT-` *(`RakitPengenalWork(AwalanKomite, …)`, `COVER_KEY` = ID klaim, `LINI`, `TYPE`, `PY_POSITION`)*, `T_GENERAL_KOMITE` *(shared PK; `ADJUSTMENT_ID`; `KOMITE_LOOP`, `KOMITE_COUNT`, `ACCEPT_STATUS`)*, `T_KOMITE_KOMITELIST` per anggota roster *(`KOMITE_URUT`, `KOMITE_ID`, `ID_KOMITE`, `KOMITE_EMAIL`)*, `T_CLAIMLF_ADJUSTMENT.KOMITE_ID`, sisi induk `.IsKomite`, `.KomiteNo`, `.TotalKomite` — muatan `childPageKomite.*` dipetakan kolom demi kolom dengan path + baris; satu transaksi; jejak direkam | `CreateKMTLife_Act` 121.652 byte |
| `ResolverBelumDiputuskan` | pembaca `M_LINK_SERVICE` *(`URL`, `KATEGORI_1`, `KATEGORI_2`, `USERNAME`)* berkunci **`(KATEGORI_1, KATEGORI_2)`** — 18 baris DEV: `Google/upload`, `Google/geturl`, `Google/delete`, `Google/changeStorageclass`, `Google/getAI`; `Klaim/insertClaimLife`, `Klaim/insertClaimAccept`, `Klaim/insertClaimDLA`; `SendEmail/OpenProteksiKlaim_RNM`; `Kasir/insertAllPaymentKasir`, `Kasir/getPaymentSlip`, `Kasir/InjectDataAccountChange`; `OpenProteksi_1/2/OpenProteksiKlaimConfirm_RNM`; `Production/convertJsonNusareToProduction` *(dibuang, spec §2b)*, `Production/getPayment*`. Kunci yang **dipakai tiap efek** dibaca dari activity-nya *(`InsertGoogleStorage_Act`, `GetUrlGoogleStorage_Act`, `DeleteGoogleStorage_Act`, `SendEmailKlaimLF`, `serviceInsertArasapasClaimLife_act`)* — jangan menebak dari nama; `URL`/`USERNAME` tidak pernah masuk log/artefak; tiruan skema uji berisi kunci + URL palsu `http://uji.invalid/...` | katalog + 5 activity |
| `AntreanBelumDiputuskan` *(tiket 12)* | **aq** `[DIPUTUSKAN — turunan tiket 12 AC "antre-ulang, kegagalan tercatat" + ADR-0008; bentuk milik executor]`: tabel **outbox lintas modul** `T_EFEK_KELUAR` *(migrasi bernomor baru: `ID` sequence, `LINI`, `MODUL`, `JENIS_EFEK`, `RUJUKAN` = ID work, `MUATAN` CLOB JSON, `STATUS` tertutup, `PERCOBAAN`, `JADWAL_BERIKUT`, `GALAT_TERAKHIR`, `DIBUAT`, `DIPERBARUI`)* ditulis **dalam transaksi yang sama** dengan aksi bisnisnya; worker mengambil `JADWAL_BERIKUT <= now` dengan `FOR UPDATE SKIP LOCKED`, backoff, anti-dobel lewat `STATUS`; kegagalan permanen → jejak audit *(bukan berputar selamanya — temuan review 13)*. **Komite tiket 06–07** memakai tabel yang sama | tiket 12, ADR-0008 |
| **an** `TokenStorage` | persis lanjutan 4 §1: `GCP_IMAGE` *(`APPNAME`, `KODEAKSES`, `USERINPUT`, `INPUTDATE`)*, token baru = MD5-hex dari `'ASMAPP' ‖ garam ‖ stempel waktu`, kedaluwarsa 1 menit; garam dari env `STORAGE_TOKEN_SALT`, kosong → gagal terang | procedure `[data DBA]` |
| `KategoriWajibBelumDiketahui` *(gerbang "dokumen lengkap")* | **ar** — sumber daftar kategori wajib. `[data DBA]` tabel **`POOLDATA.CATEGORY_ATTACH_CLAIMLIFE`** *(`NOTE VARCHAR2(30)`, `BISNIS VARCHAR2(100)`, `POSITION VARCHAR2(100)`, `NOU VARCHAR2(2)`)* berisi **6 baris**, seluruhnya `BISNIS = 'ALL'`, `POSITION` dan `NOU` NULL: *Notice of Claim · Internal Claim Requitition · Treaty/Nota Penutupan · Email · Photo · Document Claim*. Kandidat kedua: view `DOCUMENTCLAIM_LIFE` = `M_PRODUCT_LIFE.JSONDATA` `DocumentClaim[*].Document` *(per produk; jalur JSON produk yang AC 38 larang)*. Rule `GetCategoryLife_SQL` **tidak ada di korpus**, jadi mana yang Pega baca **tidak terverifikasi**. **`ar1` `[DIPUTUSKAN — 27 September 2026, ditafsirkan asisten dari "jangan banyak yang di-skip"; cabut dengan mengganti kata ini]`:** daftar wajib dibaca saat jalan dari `CATEGORY_ATTACH_CLAIMLIFE` *(filter `BISNIS IN ('ALL', <kode bisnis>)`; `POSITION` bila terisi)*, dibandingkan **cacah** seperti gerbang XML *(tiket 03)*; tiruan skema uji memuat keenam label itu *(label konfigurasi, bukan data orang)*; tiket mencatat kandidat kedua dan bahwa keputusan ini `[USULAN yang disahkan]`, bukan `[terverifikasi]` | katalog |

**Perapian yang ikut di A2** *(lanjutan 4 bab 8)*: `hapus.go` → **405/409** untuk `ErrHapusFisikDilarang`;
cabang 501 mati untuk `ErrPenomorBelumDiputuskan`/`ErrJejakBelumDiputuskan` dibuang beserta tipe stub
yang tak terpakai *(atau dipindah ke `_test.go` sebagai test double bernama)*; setiap penggantian stub
menggeser centang AC tiket asalnya *(02, 09, 10, 11, 12, 03)* **dengan sebabnya**.

---

## 3. A3 — PARITAS LAYAR DAN AKSI DARI XML

Persis lanjutan 4 §3 A3 *(sensus 31 aksi + 16 flow action + 3 harness + 20 section; aturan paritas;
tujuh kelompok; satu commit per kelompok)*, ditambah yang sudah diketahui sejak itu:

| Kelompok | Tambahan fakta `[terverifikasi]` |
| --- | --- |
| Register | precondition `SavePesertaClaim` 1812 *(`IsCheck && L1 \|\| … \|\| L11` — cacat kurung, tiru bacaan yang dimaksud)*; jalur `ContentNote == DEATH` / `!= DEATH` di 1997/2413; `Protect.pxResults` *(`ProtectCloseClaim_act`)* di 2555/2704 |
| Outstanding · Medis · Akseptasi | detail klaim `ClaimLifeDetailGCNM` diluncurkan lewat `pyEditAction` → `ViewClaimDetailLifeGCNM` di **ketiga** layar tahap, bergerbang `pyWorkPage.Save = 1`; tombol **Save Adjustment** tanpa gerbang peran → di React: panel detail yang sama dipakai tiga halaman, tombol tampil bila gerbang **activity** terpenuhi *(bab 7)* dan pelaku memegang tahapnya |
| Medis | **al** diputuskan dari `Section\Diagnose_Section.xml` *(29 rujukan kelas `Data-DiagnoseLife`)*, `Harness\Diagnose_Harness.xml`, `SearchDiagnose_act`, `SetDisease`, dan grid `.DiagnoseList` di `ClaimLifeDetailGCNM` — bukan dari `MedicalCheckClaimLife` yang memang nol; petunjuk katalog: kolom peserta `DISEASE`, `ICD_CODE`, `NOTES`, `KETERANGAN` di tabel warisan → kemungkinan **satu diagnosa per peserta** *(kolom `003`)*, daftar hanya hasil pencarian master *(`M_DISEASE_LIFE_SEQ` ada di katalog → master penyakit `[data DBA — belum dibaca]`)* |
| Dokumen | kategori wajib **ar** *(§2)*; penyimpanan warisan `ATTACHDOCUMENTCLAIM_LIFE` *(`ID`, `CASEID`, `DATA_JSON` CLOB; 727 baris)* dan `DOCUMENT_CLAIM` *(14 kolom)* — layar `AttachDocScreenLife`/`DocumentLife` menulis ke mana **dibaca dari `InsertDocument_Act` + `InsertGoogleStorage_Act`**, bukan ditebak dari dua tabel ini |
| Komite | `GetListKomiteLife` = roster A2; `CreateKMTLife_Act` = kasus A2 |

Keluaran A3 wajib: `PARITAS-LAYAR-DAN-AKSI.md` lengkap *(31 + 16 baris berverdict)*, halaman React per
flow action, komponen per section, `api.ts` per rute, test JS per gerbang; laporan akhir menyebut
**cacah aksi yang kini punya rute dan kontrol** dari 31.

---

## 4. A4 — MIGRASI DATA 13

Persis lanjutan 4 §3 A4 *(pengurai nomor klaim o1; posisi tahap **ao** `[USULAN]` di balik saklar; skema
uji saja)*, ditambah: dokumen warisan bersumber `ATTACHDOCUMENTCLAIM_LIFE.DATA_JSON` per `CASEID`
*(JSON lampiran — bentuknya dibaca dari agregat/`InsertDocument_Act`, bukan dari baris nyata)* →
`T_CLAIMLF_DOCUMENT`; rekonsiliasi cacah dokumen per klaim ikut laporan. ⛔ `POOLDATA` tidak disentuh.

---

## 5. BAGIAN B — MODUL KOMITE CLAIM LIFE, seluruhnya

Persis lanjutan 4 §4 *(urutan 00 → 01 → 02 → 03 → 04a → 04b → 05 → 06 → 07 → 08 → 09; aturan sama;
`komite-claim-life: tiket NN — <judul>`)*, dengan bekal yang kini ada:

| Bekal | Isi |
| --- | --- |
| Tabel | `013` sudah membuat `T_GENERAL_KOMITE`, `T_KOMITE_KOMITELIST`, `SEQ_KOMITE_KOMITELIST`; `COVER_KEY` di `001` → tiket **00** = verifikasi + penjaga + yang kurang |
| Roster, kasus, outbox, resolver, jejak | dari A2 — dipakai ulang, **tidak** ditulis kembar; Kasir lewat `M_LINK_SERVICE` `Kasir/*` untuk tiket **07** *(`HitServiceToKasirKMTLife_Act`)* |
| Penulis keputusan | `KomitePostAdjustment` tersensus *(`1` ×3, `2` ×1 pada baris **dan** peserta; `ACCEPTATION_DATE` ×3; gerbang `KomiteCount == KomiteLoop`)* → tiket **05** |
| ⚠️ Nomor akseptasi Komite | `SUMBER-PENOMORAN-DBA.md`: `(ASM-FW-GCNMFW-Work-KomiteLife, RNML-A)` adalah pasangan **counter procedure**, sedangkan Claim Life `SaveAdjustment_Act` memakai **sequence** `ACCEPTATIONNOLIFE_SEQ` dengan format `RNML-A<BusinessCode>.<MM>.<YY>.<NNNNN>` *(lanjutan 4 bab 7)*. **Dua mekanisme untuk awalan yang sama.** Tiket **04a** membaca `Generate_NoAccept_KMT_Life` *(ada di korpus, tiket 14 sensus)* utuh dan mengikuti XML-nya; bila keduanya memang berbeda, keduanya ditiru masing-masing dan selisihnya dicatat `[terverifikasi]` untuk work owner |
| Decision table | `DecisionTable\` 1 berkas → **tabel data** di `models`, nol `if` bercabang; testnya membandingkan seluruh baris tabel |
| Aksi layar | `SetRemarkKomiteLife`, `DownloadDocumentClaim`; 3 section, 3 flow action, 1 flow → `PARITAS-LAYAR-DAN-AKSI.md` milik modul ini |

---

## 6. ATURAN YANG TETAP — pengingat satu paragraf

XML menang dan tiket diralat dengan bukti *(modul §1.2)*; penyarangan dibaca sebagai pohon *(lanjutan 4
bab 7)*; verifikasi penuh sebelum commit **dan** sesudah perbaikan review; satu review per tiket/paket;
prosa tiket dibatasi *(lanjutan 2 §1-4)*; nol nama orang, nomor polis, kredensial, URL, atau username di
artefak; `POOLDATA` tidak pernah menjadi sasaran `-migrate`/test db; migrasi baru hanya dari keputusan
yang tercatat *(aj, am, af, al bersyarat, aq, ar1, an)*; berhenti hanya sesudah commit hijau dan
hanya karena gerbang di luar executor.

---

## 7. LAPORAN AKHIR GILIRAN

Persis lanjutan 4 §5 *(tabel A0, Bagian A, Bagian B; keputusan yang dipakai; terbuka menurut pemilik;
telemetri)*, ditambah satu baris: **stub yang tersisa di kode sesudah giliran ini** *(nama dan sebabnya)*.

---

*Disusun 27 September 2026 sesudah verifikasi `f8bbf7a` (lanjutan 4 bab 8): stub yang masih hidup
dibaca dari kode, permukaan HTTP/React dihitung, katalog `CATEGORY_ATTACH_CLAIMLIFE`, `M_ATTACHMENT_CATEGORY`,
`ATTACHDOCUMENTCLAIM_LIFE`, view `DOCUMENTCLAIM_LIFE`, dan sebaran kunci `M_LINK_SERVICE` dibaca
(definisi dan agregat; label konfigurasi saja, nol baris data orang).*
