# Grilling — Komite Claim Life — Ronde 1

Status: answered (work owner, 2026-09-15) — **frontier Ronde 1 kosong**
Konteks: `komite-claim-life` (Komite Claim Life)
Tanggal: 2026-09-15
Skill: `/mattpocock-skills:grill-with-docs` (grilling + domain-modeling)
Modul: `Komite Claim Life` — 47 berkas, 11 tipe rule

> **Konvensi penandaan.** `[terverifikasi]` = terbukti korpus dengan path + rule;
> `[keputusan work owner]` = keputusan bisnis; `[data DBA]`; `[terbuka]` = OQ.

---

## Bagian 0 — KOREKSI atas klaim Ronde 1 saya

Dua klaim saya salah. Keduanya saya baca ulang ke korpus dan perbaiki di sini.

### K1. Gerbang EXIT ada di **step 9 (Arasapas)**, bukan di `InsertJsonClaimLife_Act`

**Klaim lama (SALAH):** "`InsertJsonClaimLife_Act` ← gerbang: `RetroID == "1000013"`".

`[terverifikasi]` Nomor langkah dibaca dari `<pyStepPageReference>RH_1.pySteps(n)` —
sumber paling tegas, bukan urutan baris. Peta langkah `KomitePostAdjustment.xml`
(`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`):

| Step | Langkah | Deskripsi |
| ---: | --- | --- |
| 1 | `Property-Set` | Open Claim |
| 2 | `Property-Set` | Set Index |
| 3 | — | set acc / reject adjustment in PNC and komite |
| **4** | blok besar (22 sub-step) | **jalur AKSEP** — generate no akseptasi, Insert ke OS, `Call PrintAkseptasiPDF`, `Call LoadDocumentLife_ACT` |
| **5** | blok besar (9 sub-step) | **jalur REJECT** — Set Reject Komite berjenjang, Set Nilai Akseptasi, `Call UpdateWorkObject` ("Insert ke OS") |
| 6 | `Obj-Save` | ← PERSIST |
| 7 | `Call InsertJsonClaimLife_Act` | |
| 8 | — | |
| **9** | **`Call serviceInsertArasapasClaimLife_act`** | **`EXIT JIKA RETROID "L0000141"`** |
| 10 | `Call SendEmailKlaimLife` | |
| 11 | `Call HitServiceToKasirKMTLife_Act` | |
| 12 | `Property-Set` | kasir |
| 13 | `Call SetInformationData` | |
| 14 | — | |

**Step 9** membawa `<pyStepsDescription>EXIT JIKA RETROID "L0000141"</pyStepsDescription>`
(baris 8474), `<pyStepsPreCondition>true`, dan **dua** baris precondition:

```
baris 8516 : pyWorkCover.ClaimData.PolicyDataLife.RetroID=="L0000141"
             || pyWorkCover.ClaimData.PolicyDataLife.SecurityReinsurerID=="L0000134"
baris 8539 : pyWorkCover.ClaimData.PolicyDataLife.RetroID=="1000013"
```

dengan `pyStepsTransParamsWhenTrue = 2` dan `…WhenFalse = 2`.

**Efek yang di-skip — inilah koreksi pentingnya.** Karena gerbang ini **EXIT** (menghentikan
activity), bukan sekadar memilih cabang, maka ketika identitas retro cocok, langkah **9 sampai 14
tidak berjalan**:

| Step | Efek | Nasib saat EXIT |
| ---: | --- | --- |
| 7 | `InsertJsonClaimLife_Act` | **tetap berjalan** (sebelum gerbang) |
| 9 | `serviceInsertArasapasClaimLife_act` | **DI-SKIP** |
| 10 | `SendEmailKlaimLife` | **DI-SKIP** |
| 11 | `HitServiceToKasirKMTLife_Act` (Kasir) | **DI-SKIP** |
| 13 | `SetInformationData` | **DI-SKIP** |

Jadi tiga identitas retro ter-hardcode itu membuat klaim tertentu **tidak pernah dikirim ke
Arasapas, tidak pernah dikirim email, dan tidak pernah masuk Kasir**.

`[dugaan]` Apakah kedua baris precondition itu ber-OR atau dievaluasi terpisah tidak dapat
dipastikan dari tag; yang pasti keduanya melekat pada step 9 dan keduanya bersifat EXIT.

### K2. Rekam akseptasi ditulis **sekali di tingkat final**, bukan tiap tingkat

**Klaim lama (SALAH, ditandai `[dugaan]`):** "rekam akseptasi tampaknya ditulis di setiap tingkat".

`[terverifikasi]` Keempat gerbang tingkat-final ada di baris **5686, 8110, 8639, 8878**:

| Baris | Gerbang | Melekat pada |
| ---: | --- | --- |
| 5686 | `AcceptStatus = 1 && KomiteCount == KomiteLoop` | step 4 — jalur aksep |
| 8110 | `AcceptStatus == 2 && KomiteCount == KomiteLoop` | step 5 — jalur reject |
| 8639 | `AcceptStatus = 1 && KomiteCount == KomiteLoop` | step 9 — Arasapas |
| 8878 | `AcceptStatus = 1 && KomiteCount == KomiteLoop` | step 11 — Kasir |

Step 4 mulai baris 1022, jadi gerbang 5686 berada **di dalam** blok aksep — yang tidak terlihat saat
saya hanya menyapu 4880–5250. **Kedua blok tulis digerbangi tingkat final.** Lebih jauh: **Arasapas
dan Kasir pun hanya berjalan pada tingkat final DAN saat aksep.**

---

## Bagian A — Temuan korpus Ronde 1 (tetap berlaku)

### A1. `KomiteRouter` punya **dua** precondition `[terverifikasi]`

`Komite Claim Life/Activity/KomiteRouter.xml`
(`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEROUTER` / `RULE-OBJ-ACTIVITY`, 26.387 byte):

```
Property-Set  param.AssignTo = .KomiteID          (baris ~294)
  ├─ .KomiteAproval == 0                          (baris ~382)
  └─ Primary.TransferType == '2'                  (baris ~442)
```

Artefak D2 hanya mencatat yang kedua. **`.KomiteAproval == 0` adalah mesin penggerak tangga.**

### A2. `AcceptStatus` tidak pernah ditulis rule mana pun `[terverifikasi]`

Nol `<PropertiesName>…AcceptStatus</PropertiesName>` di 47 berkas. Ia diisi lewat **dropdown wajib**
di `Komite Claim Life/Section/ShowTransfer.xml`
(`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `SHOWTRANSFER` / `RULE-HTML-SECTION`, 1.125.234 byte), baris
32607: `pyValue = .AcceptStatus`, `pyFormat = pxDropdown`, `pyRequired = true`,
`pyRequiredNew = always`.

### A3. Dua blok tulis kembar `[terverifikasi]`

`UpdateOsAkseptasiClaimLife_sql` dipanggil **2×** (`<RequestType>` = 2, terindeks 2 → keduanya
aktif). Diff struktural himpunan properti kedua blok: **nol beda**.

### A4. Dua rule penomoran Komite **dead** — pola sama ADR-0006 `[terverifikasi]`

| Rule | `<RequestType>` | terindeks | Verdict |
| --- | ---: | ---: | --- |
| `Generate_NoAccept_KMT_Life` | 1 | **0** | dead |
| `Generate_NoAccept_KMT_LifeRetro` | 1 | **0** | dead |
| `GetKodeProdLife_SQL` | 1 | 2 | aktif |
| `GetSequenceNumber_SQL` | 1 | 2 | aktif |

### A5. Lima efek keluar; Kasir tidak ada di Claim Life `[terverifikasi]`

`HitServiceToKasirKMTLife_Act` = `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `HITSERVICETOKASIRKMTLIFE_ACT`
/ `RULE-OBJ-ACTIVITY`. Ditambah `SetInformationData` (step 13) yang tidak tercatat di D2.

### A6. Hanya dua activity dipicu dari UI `[terverifikasi]`

`DownloadDocumentClaim` (4×), `SetRemarkKomiteLife` (2×). `KomitePostAdjustment` berjalan lewat
pemrosesan FlowAction `ViewTransferDtl`
(`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `VIEWTRANSFERDTL` / `RULE-OBJ-FLOWACTION`, 31.963 byte).

---

## Bagian B — Delapan keputusan work owner (2026-09-15)

### B1. `TransferType` = **dead code** → dibuang

`[terverifikasi]` Tidak pernah diisi — **nol `Property-Set`** di Claim Life maupun Komite Claim Life.
Hanya dibaca di dua tempat: precondition `KomiteRouter` (baris ~442) dan visible-when
`ShowTransfer.xml` `.TransferType==2` (baris ~29177). Karena tak pernah di-set,
`TransferType=='2'` **selalu FALSE** → cabang mati, sisa salinan dari modul Komite lain.

`[keputusan work owner]` **Buang `TransferType` dan UI `ShowTransfer` yang bergantung padanya.**
Routing tingkat = cari baris roster pertama ber-`KomiteAproval == 0`, sasaran `param.AssignTo =
.KomiteID`. → **OQ-033 DITUTUP sebagai dead code.**

### B2. Penegakan pemutus per tingkat

`[keputusan work owner]` Hanya pemilik `KomiteList(KomiteCount).KomiteID` pada tingkat berjalan yang
boleh memutuskan; pengguna lain **ditolak di lapisan layanan**. Pega hanya *menempatkan* tugas di
worklist — penempatan, bukan penegakan. Sistem baru memperketat secara sadar. → **ADR-0014**.

### B3. Dropdown keputusan = enum tertutup `{1, 2}`

`[keputusan work owner]` Hanya **Setuju** (`AcceptStatus = 1`) dan **Tolak** (`AcceptStatus = 2`).
Tidak ada opsi ketiga. Konsisten korpus: `When/IsKomiteLoop.xml` lanjut hanya pada `"1"`;
`KomitePostAdjustment` hanya menangani `1` dan `2`. Nilai lain **ditolak terang-terangan**
(fail-loud), bukan gagal diam-diam.

### B4. Cutover 7 Feb 2025 = tidak dipakai lagi

`[keputusan work owner]` Blok precondition
`…OfferFacIn.PolicyData.ProdDateTime < "20250207T000000.000 GMT"` (baris ~5135 dan ~7861) sudah
di-remark/tidak dipakai. **Jangan direplikasi.** Pengisian identitas retro cukup dua precondition
aktif: `Type=="TP"||Type=="TR"` dan `SecurityReinsurerID!="" && SecurityReinsurer!=""`.
→ **OQ-034 DITUTUP.**

### B5. Dua blok tulis kembar → **disatukan**

`[terverifikasi + keputusan work owner]` Blok #1 (aksep, `RequestType` baris ~5242, gerbang ~5686)
dan blok #2 (tolak, baris ~7968, gerbang ~8110) berisi identik, beda hanya nilai status.
**Sistem baru: satu jalur simpan berparameter status.** Penyimpangan sadar dari paritas struktural —
dua blok kembar adalah sarang bug. Unit keputusan tetap baris `AdjustmentList` (**ADR-0011**).

### B6. Nomor & rekam akseptasi dibuat **sekali, di tingkat final**

`[terverifikasi]` Kedua blok tulis digerbangi `KomiteCount == KomiteLoop`. Penomoran
`GetSequenceNumber_SQL` (baris ~2779) berada di jalur yang digerbangi
`AcceptStatus = 1 && KomiteCount == KomiteLoop` (baris ~8639, ~8878).
Di tingkat bukan-terakhir hanya `KomiteAproval = AcceptStatus` dicatat (baris ~792, ~908), lalu
`KomiteCount = KomiteCount + 1` (baris ~9020) — naik tingkat **tanpa** nomor akseptasi.
→ **OQ-006 DITUTUP.**

### B7. Semua efek keluar **wajib berhasil** — menyimpang dari ADR-0008

`[keputusan work owner]` Keempat efek keluar setelah `Obj-Save` — `InsertJsonClaimLife_Act` (step 7),
`serviceInsertArasapasClaimLife_act` (step 9, endpoint lewat `M_LINK_SERVICE` kunci
`Klaim`/`insertClaimLife`, **ADR-0013**), `SendEmailKlaimLife` (step 10),
`HitServiceToKasirKMTLife_Act` (step 11) — **wajib berhasil**, jaminan **at-least-once**, tidak boleh
gagal diam-diam.

⚠️ **Ini MENYIMPANG dari ADR-0008**, yang untuk Claim Life membolehkan efek asinkron non-blocking.
Pemicunya: **Kasir** — integrasi keuangan. Pola: **transactional outbox** — keputusan + daftar efek
disimpan dalam satu transaksi; worker mengirim tiap efek dengan retry sampai sukses; keputusan komite
**tidak dianggap tuntas** sampai semua efek berhasil atau ditandai perlu intervensi. → **ADR-0015**.

### B8. Tiga identitas retro ter-hardcode → **dibuang**

`[keputusan work owner]` `1000013`, `L0000141`, `L0000134` dibuang dari logika sistem baru, **tidak**
direplikasi sebagai konstanta. Bila kelak diperlukan, jadikan **data/konfigurasi** seperti roster
`EMAILKOMITE`.

⚠️ **Konsekuensi perilaku yang perlu disadari** (dari koreksi K1): ketiga identitas itu hari ini
memicu **EXIT di step 9**, sehingga klaim yang cocok **tidak** dikirim ke Arasapas, **tidak** dikirim
email, dan **tidak** masuk Kasir. Membuangnya berarti klaim-klaim itu **mulai menerima keempat efek
keluar**. Ini **perubahan perilaku**, bukan sekadar pembersihan kode — dan ia berinteraksi langsung
dengan B7 (semua efek wajib berhasil).

---

## Status OQ setelah Ronde 1

| OQ | Verdict |
| --- | --- |
| **OQ-033** | **TERTUTUP** untuk Komite Life — `TransferType` dead code |
| **OQ-034** | **TERTUTUP** untuk Komite Life — blok cutover di-remark, dibuang |
| **OQ-006** | **TERTUTUP** — nomor & rekam akseptasi sekali di tingkat final |
| **OQ-035** | tetap terbuka — `UpdateWorkObject` / `serviceInsertArasapas` lintas modul; **bukan pemblokir isi** |
