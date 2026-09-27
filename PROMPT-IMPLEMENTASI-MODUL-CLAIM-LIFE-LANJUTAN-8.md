# PROMPT — lanjutan 8: **F0.3 → F0.4 → F0.5 → A3 Register dalam SATU giliran**, Inbox dengan kolom `TAHAP` dan `TGL_CREATE` *(keputusan at, au)*, ralat tiket 08 dan 14

> Dibaca sesudah brief modul dan lanjutan 1–7; semuanya tetap berlaku kecuali yang diralat di sini.
> Perintah work owner 27 September 2026 tidak berubah: *"JANGAN HANYA 1 TIKET … JANGAN BANYAK YANG DI
> SKIP! BACA ULANG LAGI XML NYA! … JANGAN MEMBUAT MENU YANG TIDAK ADA DI SISTEM SEBELUMNYA (PEGA) …
> KALAU DI TIKET TIDAK BENAR … PERBAIKI JUGA DI TIKET NYA."*

---

## 0. KEADAAN SESUDAH `5b49954` — F0.2 DIVERIFIKASI ULANG

| Klaim laporan F0.2 | Diperiksa ulang | Hasil |
| --- | --- | --- |
| 12 berkas, sesi + klien fetch + Login stub | `git show --stat`: 12 berkas, +1.082/−353 | ✅ |
| `axios` dilepas; build 41 modul, 158,83 kB | `package.json` tanpa axios; build 41 modul, 158,83 kB | ✅ |
| 76 uji JS; Go 261 · 0 · 34 | vitest 6 berkas 76 lulus; `go test -tags=db` 261 PASS · 0 FAIL · 34 SKIP | ✅ |
| `hapusKlaim` memakai DELETE | `api.ts:535` `metode: 'DELETE'` | ✅ |
| ralat ADR-U-0004 → ADR-0013 | `api.ts:84–85` | ✅ |
| `\b` jadi 0x08 diperbaiki | `grep -P '\x08'` di `src`: nol | ✅ |
| peran asing dibuang saat dibaca ulang; gerbang `VITE_AUTH_STUB` gagal-tertutup | `sesi.ts` `saringPeran`; `bolehMasukStub` `=== 'true'` | ✅ |
| nol kebocoran | alamat/sandi/service/nama operator di `src`: nol | ✅ |

**Ralat brief 7 §1.1** *(catatan executor benar)*: Python **ada** lewat peluncur `py` *(3.14)*; yang
mati hanya alias `python` bawaan Store. Pakai `py`. `.scratch/alat/baca-xlsx.ps1` tetap sah.

---

## 1. MEKANISME GILIRAN — kenapa berhenti lagi, dan aturannya kini

Giliran F0.2 berakhir dengan pesan *"Lanjut F0.3: …"*. Pesan tanpa panggilan alat sesudahnya =
giliran selesai = manusia harus menempel lagi *(lanjutan 3 §1 sudah menjelaskan sebabnya dari teks
skill `implement`/`code-review`)*. Aturan yang berlaku giliran ini, tanpa tafsir:

1. Sesudah **setiap** commit paket: **nol teks ke manusia**. Tulis catatan kemajuan ke
   `.scratch/claim-life/LAPORAN-GILIRAN-F0.md` *(bab per paket: SHA, uji, XML dibaca, telemetri)*,
   lalu **panggilan alat berikutnya langsung** *(membaca XML paket berikutnya)*.
2. Laporan `/code-review` masuk berkas yang sama, bukan pesan.
3. Pesan ke manusia **hanya satu**, sesudah **A3 kelompok Register** ter-commit — atau sesudah berhenti
   sah *(lanjutan 2 §1–2: gerbang yang hanya manusia dapat buka, disebut namanya)*. Pertanyaan pohon
   XML **dijawab dengan pohon** *(PowerShell `[xml]`, lanjutan 4 §7)*, bukan dibawa ke manusia.
4. Konteks menipis bukan alasan berhenti: tulis catatan ke berkas di butir 1, lanjutkan.

---

## 2. KEPUTUSAN **at** DAN **au** — kolom `TAHAP` dan `TGL_CREATE` di `T_WORK_CLAIM`

### 2.1 Bukti XML yang memaksanya

| # | Bukti | Isi |
| ---: | --- | --- |
| 1 | `Flow\Register_Flow.xml` | empat assignment adalah **keadaan** kasus: `Input Register` 358, `Outstanding Claim` 343, `Medical Check` 268, `Claim Analis` 313; `Decision3` `IsSendtoAdmin` → **Assignment2 Input Register** *(sheet `Struktur_InboxClaimLife.xlsx` baris 172)* |
| 2 | `Section\InputOSClaimLife.xml` | tombol **`Send Back to Register`** 21404 → `<pyLocalAction>SendtoAdmin` 21433 *(flow action `SendtoAdmin` → `SendtoAdmin_Act`: `Property-Set pyWorkPage.SendtoAdmin="1"`, prasyarat `pyWorkPage.pyPosition=="ReasLifeMedicalAdvisor"` WhenTrue=2 WhenFalse=3)*; tombol **`Send to Medical Check`** 21349/21839 → `SendtoAdmin_Act1` 21863 *(`SendtoAdmin="0"`, prasyarat sama)* |
| 3 | `ReportDefinition\InboxPremiumList.xml` | kolom `Create Date/Time` = `pxCreateDateTime` 736 — Pega menyimpan **waktu buat** setiap work object, terpisah dari waktu ubah |
| 4 | `APP_RNM` sekarang | `models.TahapDariPeran` memetakan Admin **selalu** ke Outstanding *(`models/tahap.go:98–105`)*; `serahTerimaSah` tanpa pasangan Admin→Admin; `T_WORK_CLAIM` hanya punya `TGL_UPDATE` yang **ditimpa** tiap perpindahan *(`repository/klaimlife.go:502–503`)*; pendaftaran menulis `CreateOpName` saja, `CreateOp` dan `PyPosition` **kosong** *(`services/pendaftaran.go:195–213`)* |

Keadaan `Input Register` **ada** di XML dan dapat dituju kembali; waktu buat **ada** di Pega. Skema kita
tidak dapat menyatakan keduanya. Tiket 08 sudah mencatatnya `[terbuka — work owner]`; XML menutupnya.

### 2.2 Keputusan `[DIPUTUSKAN — 27 September 2026, diturunkan dari XML di 2.1; work owner dapat memveto sebelum migrasi dijalankan di Oracle mana pun]`

| Butir | Isi |
| --- | --- |
| **at** | `T_WORK_CLAIM.TAHAP VARCHAR2(32)` menyimpan nama assignment **VERBATIM `pyTaskName`** *(`Input Register` / `Outstanding Claim` / `Medical Check` / `Claim Analis`)*. `PY_POSITION` **tetap** *(peran pemegang, ADR-U-0002)*. Pendaftaran menulis `TAHAP='Outstanding Claim'` *(Assignment2 selesai saat form register disimpan)* dan `PY_POSITION=ReasLifeAdmin`. `models.Tahap` dibaca dari `TAHAP`; bila NULL *(baris lama)* jatuh ke `TahapDariPeran`. `serahTerimaSah` menjadi peta **tahap**, bukan peran: Register↔Outstanding *(Admin→Admin: maju saat register disimpan, balik lewat `Send Back to Register`)*, Outstanding→Medical, Medical→{Analis, Outstanding}, Analis→{Outstanding, Medical} *(ADR-U-0002 jalur balik tetap benar; yang berubah hanya butirannya)*. `Pindah` menerima tahap tujuan dan menulis keduanya |
| **au** | `T_WORK_CLAIM.TGL_CREATE DATE` = `pxCreateDateTime`; pendaftaran menulisnya sekali; `TGL_UPDATE` tetap waktu ubah. Pendaftaran juga menulis `CREATE_OP = pelaku.AkunID` *(padanan `pxCreateOperator`, yang dirutekan worklist)* di samping `CREATE_OP_NAME` |
| migrasi | **`016_kolom_tahap_dan_tgl_create.sql`** + `_down`: `ALTER TABLE {skema}.T_WORK_CLAIM ADD (TAHAP VARCHAR2(32), TGL_CREATE DATE)` *(pola 014; `KolomAlterTambah`; nullable, ADR-U-0027)*; pengisian baris lama *(bila ada)* di kode, bukan di DDL: `TAHAP` dari `PY_POSITION` lewat `TahapDariPeran`, `TGL_CREATE` dari `TGL_UPDATE` — dilakukan di `BongkarBarisLama`/migrasi data tiket 13 sebagai langkah bernama, dan dicatat |
| uji | murni: peta serah terima baru *(termasuk Admin→Admin dan penolakan lompatan)*; `db` *(SKIP tanpa skema uji)*: register → `TAHAP`, `TGL_CREATE`, `CREATE_OP` terisi; `Pindah` menulis `TAHAP` |
| tiket | **08**: blok `### Ralat menurut XML — 27 September 2026` menutup `[terbuka]` "Register vs Outstanding" dengan bukti 2.1 dan keputusan at; **14**: blok ralat kolom `TAHAP`, `TGL_CREATE` pada `T_WORK_CLAIM` *(AC 34/35 tetap)*; `STRUKTUR-TABEL-CLAIM-LIFE.md` diperbarui |

⚠️ Prasyarat `pyPosition=="ReasLifeMedicalAdvisor"` pada kedua activity di 2.1 butir 2 **belum
dijelaskan** *(dengan posisi Admin langkahnya dilewati — lalu bagaimana `Send Back to Register`
bekerja?)*. Itu pertanyaan **pohon** untuk A3 kelompok Outstanding *(baca `pyPropertiesValue`
`pyPosition` di `Register_Flow.xml` 583–772 dan urutan shape-nya)* — bukan alasan menunda at.

---

## 3. F0.3 → F0.4 → F0.5

### 3.1 F0.3 — Shell dan menu *(persis lanjutan 7 §1.2)*

Kelompok **Claim Life**: `Inbox Claim Life` `[tidak ada di korpus]` dan `Register` `[Register_Flow.xml:155]`.
Tidak ada yang lain. Cabang "sudah masuk" di `App.tsx` sekarang *(header `bilah-sesi` + dua halaman
ditumpuk tanpa gaya — yang terlihat work owner di `localhost:5173` sesudah masuk)* adalah keadaan
**sementara F0.2** dan **diganti Shell di F0.3**; sesudah F0.3 tidak boleh ada halaman yang dirender
di luar Shell. Uji: bukti `MENU.register`; sidebar tanpa nama modul lain, tanpa `BelumTersedia`;
menu per peran; palet; lipat/laci. Judul aplikasi "Nusantara Re" `[tidak ada di korpus]`.

### 3.2 F0.4 — Inbox Claim Life *(lanjutan 7 §1.3, dilengkapi)*

| Bagian | Ketentuan |
| --- | --- |
| Backend dulu | keputusan at + au *(migrasi 016, model, pendaftaran, `Pindah`, uji)* — commit `claim-life: at/au — kolom TAHAP dan TGL_CREATE T_WORK_CLAIM`; baru rute daftar |
| Rute | `GET /api/klaim-life?tahap=<1..4>&halaman=&ukuran=` — `tahap` = `models.Tahap`; sumber `T_WORK_CLAIM ⋈ T_GENERAL_CLAIM` **memakai `TAHAP`** *(fallback `PY_POSITION` bila NULL)*; `ukuran` bawaan 50, maksimum 500, urut `TGL_CREATE DESC` *(`InboxPremiumList.xml` 592 · 942 · 733)*; jawaban memuat `total`; `tahap` di luar 1–4 → 400; pelaku kosong → 403 |
| Gerbang *(services)* | Input Register dan Outstanding Claim = **worklist**: `CREATE_OP = pelaku.AkunID` **dan** `WajibPeran(PeranAdminLife)` *(Register_Flow 1508, 1351)*; Medical Check dan Claim Analis = **workbasket**: `WajibPemegangTahap` *(1072, 1198)* |
| Kolom dan labelnya | `Case ID` → `CASE_ID` `[InboxPremiumList.xml:721]`; `Create Date/Time` → `TGL_CREATE` `[736]`; `Create Operator Name` → `CREATE_OP_NAME` `[751]`; `Work Status` → status turunan klaim `[765]`; `Claim No` → `CLAIM_NO` `[InputOSClaimLife.xml:951]`; `Policy No` → `POLICY_NO` `[tidak ada sebagai label kolom; tombol "Choose Policy No" InputRegisterClaimLife.xml:3776]`; `BUSINESS_NAME` → label `Class of Business` `[InputOSClaimLife.xml:10036]` **hanya bila** kolom itu memang kelas bisnis *(periksa tiket 02/03: isinya nama bisnis dari produk)* — bila bukan, `[tidak ada di korpus]`; `CURRENCY` → `[tidak ada di korpus]`; tahap → `TAHAP` VERBATIM |
| Tab per peran, klik baris, tombol `Register`, lencana | persis lanjutan 7 §1.3; tab *Input Register* kini **berisi** kasus ber-`TAHAP='Input Register'` *(dikembalikan lewat `Send Back to Register`)* — tidak lagi `Kosong` bersyarat |
| Uji | murni services *(worklist milik sendiri, workbasket per peran, 400, 403, batas 500)*; handler; JS per peran *(tab tampil, lencana, klik baris → rute tahap, tombol Register)* |

### 3.3 F0.5 — halaman lama masuk shell

`RegisterKlaim` lewat menu dan tombol kepala Inbox; `KlaimLife` lewat klik baris; `App.tsx` lama
dibuang; `PANDUAN-MENJALANKAN.txt` bab 4 *(login stub: `VITE_AUTH_STUB=true` **dan** `AUTH_STUB=true`
backend; peran; Inbox; cacah uji)*. Bukti: `npm run dev` masuk → Inbox → Register → Detail, nol galat
konsol; build; tsc; vitest; Go.

---

## 4. A3 KELOMPOK 1 — REGISTER *(dikerjakan di giliran yang sama, sesudah F0.5)*

Aksi PARITAS baris 3, 10, 12, 13, 14, 15, 16 + flow action 12 *(`UploadCSV_ClaimLife`)*; commit
`claim-life: A3 — Register`. Bekal supaya tidak berhenti:

| Sumber | Isi yang sudah dipastikan |
| --- | --- |
| Label VERBATIM `Section\InputRegisterClaimLife.xml` | `Choose Policy No` 3776/3827 · `Find Insured` 7057/7111 · `Type` 9104 · `Marketing Officer` 9890 · `Ceding` 11541 · `Policy Holder` 11736 · `Class of Business` 12171 · `Date Received Email` 13384 · `Response Date` 13590 · `Confirmation Date` 13795 · `Status` 14002 · `Updated Status` 14406 · `Realization Date` 14601 · `Name of Insured` 16064 · `Select Insured` 20008/20060 |
| Popup `SearchPolicy_Harness` | dipanggil 3845–3847 dan 4006–4007 *(`pyTarget popup`)*; section `SearchPolicy_Section` memakai `InboxPremiumList_Claim` *(class `ASM-FW-GISFW-Work-LIFE`)* dan `InboxPremiumList` *(`Assign-Worklist`)*, aktivitas `SearchPolicyHolder_act`, `setDetailClaim_act` *(**menulis** — PARITAS Koreksi 3)*; tombol `Search`, `Choose` |
| Sumber dropdown *(sheet baris 17–19)* | `BrowseCedingCoLife_RD` class `Int-AGENT` *(field `ID`, `ClientName`, `StatusActive`, `Leader0`, `ChildCount`)*; `BrowseBusinessLife_RD` class `Int-BUSINESS` *(`ID`, `OLDID`, `GroupPanel`, `Note`, `NoteINA`, `BusinessGroupID`, `ContentNote`, `PolicyCost`, `MinPremi`, `MaxDisc`, …)*; `BrowseMarketingOfficer_RD` class `Int-marketingofficer` *(`ID`, `ClientID`, `ClientName`, `BranchDetailID`, `BranchDetailName`, `MOLeader`, `MOStatus`, `TeamGroup`, …)* |
| `[data DBA — katalog DEV, 27-09-2026]` | `POOLDATA.AGENT` **TABLE** 28 kolom 429 baris *(`ID`, `CLIENTID`, `CLIENTNAME`, `STATUSACTIVE`, `LEADER0`, `CHILDCOUNT`, …)*; `POOLDATA.BUSINESS` **TABLE** 18 kolom 177 baris *(`ID`, `OLDID`, `GROUPPANEL`, `NOTE`, `NOTEINA`, `BUSINESSGROUPID`, `BUSINESSGROUPNAME`, `CONTENTNOTE`, …)*; `POOLDATA.MARKETINGOFFICER` **TABLE** 14 kolom 74 baris *(`ID`, `CLIENTID`, `CLIENTNAME`, `BRANCHDETAILID`, `BRANCHDETAILNAME`, `MOLEADER`, `MOSTATUS`, `TEAMGROUP`, …)*; `POOLDATA.CURRENCY` **VIEW** 8 kolom 38 baris *(`ID`, `CURRENCY`, `ISOSYMBOL`, `CURRENCYSYMBOL`, …)*; `POOLDATA.DISEASE_LIFE` **TABLE** 3 kolom 97.586 baris *(`ID`, `ICD_CODE`, `DISEASE`)* — untuk popup Diagnose kelompok Medis. Nama kolom = nama field RD dalam huruf besar. ⛔ Hanya dibaca lewat query berbatas di repository; **nol** baris disalin ke dokumen/uji; fixture `UJI-*` |
| Yang sudah ada | `GET /api/peserta-life` *(`LoadDataPeserta_Act`)*, `POST /api/klaim-life`; `RegisterKlaim.tsx` |
| Cara | section dibaca sebagai **pohon**: tiap field, `<pyCondition>`/`pyVisible`, tombol → aktivitas → langkah *(`pyStepsActivityName`, `PropertiesName/Value`, prasyarat 2 = lanjut, 3 = lewati)*; yang tidak ditiru dinyatakan dengan bukti; ralat tiket 02 bila layar XML berbeda dari tiket |

Selesai Register → lanjut kelompok **Outstanding** *(lanjutan 6 §3, pertanyaan pohon 2.2 dijawab di
sana)* selama giliran masih hidup.

---

## 5. LANGKAH 0

`git add PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE-LANJUTAN-8.md` → commit `docs: brief lanjutan 8 — satu
giliran F0.3–A3 Register; keputusan at/au kolom TAHAP dan TGL_CREATE` → `git status --porcelain`
kosong → uji hijau *(261 · 34 SKIP · 76 JS · 41 modul)* → SHA = titik tetap **F0.3**.

## 6. LAPORAN AKHIR GILIRAN — satu-satunya pesan ke manusia

Persis lanjutan 7 §6, ditambah: tabel **keputusan at/au → migrasi → kolom → uji**, daftar ralat
tiket *(08, 14, 02 bila ada)* dengan barisnya, dan tabel **menu → bukti XML → halaman → rute**.

## 7. TELEMETRI EKSEKUSI

Per paket/kelompok, satu baris di `LAPORAN-GILIRAN-F0.md` dan diringkas di laporan akhir: SHA ·
berkas · uji *(Go PASS/FAIL/SKIP, JS, modul)* · XML dibaca ulang *(cacah + nama)* · taksiran token
*(±, sebut "taksiran")* · waktu. Tutup dengan total giliran.

---

*Disusun 27 September 2026 sesudah verifikasi `5b49954` (build, tsc, vitest, `go test -tags=db`
dijalankan ulang; `api.ts`, `sesi.ts`, `Masuk.tsx` dibaca), pembacaan ulang `Register_Flow.xml`,
`InputOSClaimLife.xml` (tombol 21349–22033), `SendtoAdmin_Act.xml`, `SendtoAdmin_Act1.xml`,
`InputRegisterClaimLife.xml` (label), empat `Browse*_RD.xml`, tiket 08 dan 14, ADR-U-0002, kode
`models/tahap.go`, `services/tahap.go`, `services/pendaftaran.go`, migrasi 001/002/012/014, dan
katalog DEV (objek, kolom, cacah baris — tanpa isi baris).*
