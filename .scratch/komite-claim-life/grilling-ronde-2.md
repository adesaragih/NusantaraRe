# Grilling — Komite Claim Life — Ronde 2

Status: answered (work owner, 2026-09-15)
Konteks: `komite-claim-life`
Tanggal: 2026-09-15
Skill: `/mattpocock-skills:grill-with-docs` (grilling + domain-modeling)
Ronde sebelumnya: `grilling-ronde-1.md`

> ⚠️ **CATATAN AUDIT.** `Komite Claim Life/Activity/KomitePostAdjustment.xml` **diperbarui work
> owner 2026-09-15**. Seluruh bacaan di berkas ini dari **versi terbaru**; bacaan Ronde 1 yang
> bertentangan dibatalkan.
> `[terverifikasi]` Bukti pembaruan: `mtime` 15 Sep 2026 09:35, ukuran **515.675 byte**
> (Ronde 1: 513.661 — tumbuh ~2 KB), stempel simpan `20260915T023044.157 GMT` dan
> `20260915T023136.191 GMT`.

---

## Bagian 0 — Verifikasi status enable/disable dari versi terbaru

Metode `[terverifikasi]`: nomor langkah dari `<pyStepPageReference>`, status dari
`<pyStepsBlockName>` — berisi `//` = **REMARK**, kosong `<pyStepsBlockName/>` = **AKTIF**.

**Koreksi metode Ronde 1.** `<pyStepPageReference>` muncul **SETELAH** `<pyStepsActivityName>` milik
langkahnya, bukan sebelum. Penomoran Ronde 1 saya karena itu meleset satu. Nomor di bawah adalah
yang benar.

### Langkah tingkat atas

| Step | Baris ref | Langkah | Status | Deskripsi |
| ---: | ---: | --- | --- | --- |
| 1 | — | `Property-Set` | AKTIF | Open Claim |
| 2 | — | `Property-Set` | AKTIF | Set Index |
| 3 | — | — | AKTIF | set acc / reject adjustment in PNC dan komite |
| 4 | — | blok **AKSEP** (18 sub-step) | AKTIF | |
| 5 | — | blok **REJECT** (6 sub-step) | AKTIF | |
| 6 | — | `Obj-Save` | AKTIF | ← PERSIST |
| **8** | 8381 | `Call InsertJsonClaimLife_Act` | **AKTIF** | |
| **9** | 8470 | *(langkah gerbang, tanpa ActivityName)* | **AKTIF** | `EXIT JIKA RETROID "L0000141"` |
| **10** | 8589 | `Call serviceInsertArasapasClaimLife_act` | **AKTIF** | Arasapas |
| **11** | 8712 | `Call SendEmailKlaimLife` | **AKTIF** | |
| **12** | 8820 | `Call HitServiceToKasirKMTLife_Act` | **AKTIF** | Kasir |
| **13** | 9003 | `Property-Set` | **AKTIF** | kasir |
| **14** | 9137 | `Call SetInformationData` | **AKTIF** | |

`[terverifikasi]` Seluruh step 8–14 ber-`<pyStepsBlockName/>` kosong → **AKTIF**. Total di berkas:
**34 blockname kosong**, **5 berisi `//`**.

### Step 9 — gerbang EXIT

`[terverifikasi]` Step 9 adalah **langkah gerbang tersendiri** (tanpa `pyStepsActivityName`),
membawa `<pyStepsPreCondition>true` dan:

```
baris 8483 : pyStepsDescription = EXIT JIKA RETROID "L0000141"
baris 8525 : pyWorkCover.ClaimData.PolicyDataLife.RetroID=="L0000141"
             || pyWorkCover.ClaimData.PolicyDataLife.SecurityReinsurerID=="L0000134"
baris 8548 : pyWorkCover.ClaimData.PolicyDataLife.RetroID=="1000013"
```

⚠️ **Koreksi atas brief.** Brief menyebut transisi `WhenTrue=6` (Exit Activity). Yang terbaca di
berkas adalah nilai **`2`**, dan nilai `2` **juga muncul di step 8** yang bukan EXIT. Kode transisi
numerik Pega **tidak dapat dipetakan ke maknanya dari korpus** — tidak ada metadata runtime di
ekspor.

Jadi: `[terverifikasi]` deskripsi penulisnya berbunyi **EXIT**, langkahnya murni precondition, dan
tiga identitas retro ter-hardcode ada di sana. `[dugaan]` bahwa kode transisinya berarti "Exit
Activity" — masuk akal dan konsisten dengan label, tetapi **tidak terbukti dari ekspor**.

### Lima langkah ter-REMARK — **tiga lebih banyak** dari yang disebut brief

| Step | Baris `//` | Langkah | Deskripsi |
| --- | ---: | --- | --- |
| **4.4** | 1725 | `RDB-List` | Generate No Akseptasi (QP,QR) |
| **4.5** | 1943 | `RDB-List` | Generate No Akseptasi (TP,TR) |
| **4.6** | 2160 | `Property-Set` | Set Nilai Akseprtasi |
| **4.14** | 4947 | `Property-Set` | Tukar SecurityReinsurer dengan RetroName |
| **5.5** | 7673 | `Property-Set` | Tukar SecurityReinsurer dengan RetroName |

Brief menyebut 4.4 dan 4.5. **4.6, 4.14, dan 5.5 juga mati** — lihat koreksi di bawah.

### ⚠️ Koreksi penting: seluruh logika tukar-RetroName mati, termasuk **dua** precondition yang brief sebut "aktif"

Brief (Ronde 1 poin 4) menyatakan: *"Pengisian identitas retro cukup dua precondition aktif:
`Type=="TP"||"TR"` dan `SecurityReinsurerID!="" && SecurityReinsurer!=""`."*

`[terverifikasi]` **Keduanya tidak aktif.** Ketiga precondition berada **di dalam langkah yang
ter-remark**:

| Step | Rentang baris | Blockname `//` | Precondition di dalamnya |
| --- | --- | ---: | --- |
| **4.14** | 4937 – 5185 | 4947 | 5098 `Type=="TP"\|\|"TR"` · 5121 `SecurityReinsurerID!=""` · 5144 `ProdDateTime<"20250207…"` |
| **5.5** | 7663 – 7911 | 7673 | 7824 · 7847 · 7870 (ketiganya sama) |

Jadi yang mati bukan hanya cabang cutover, melainkan **seluruh langkah "Tukar SecurityReinsurer
dengan RetroName"** di kedua jalur (aksep dan reject). Kesimpulan Ronde 1 bahwa cutover 7 Feb 2025
tidak dipakai lagi **tetap benar dan kini terbukti korpus** — tetapi alasannya lebih luas daripada
yang dinyatakan brief.

`[pertanyaan terbuka]` Karena langkah itu mati, **`RetroName`/`RetroID` tidak lagi ditukar** sebelum
rekam akseptasi ditulis. Apa yang kini masuk ke kolom retro pada `OS_AKSEPTASI_KLAIM_LIFE`?
→ frontier Ronde 3.

### Daftar final efek keluar AKTIF

`[terverifikasi]` **Empat efek keluar**, seluruhnya **setelah `Obj-Save` (step 6)**, dalam urutan:

| Urutan | Step | Efek |
| ---: | ---: | --- |
| 1 | **8** | `Call InsertJsonClaimLife_Act` |
| 2 | **10** | `Call serviceInsertArasapasClaimLife_act` — endpoint lewat `M_LINK_SERVICE`, kunci `Klaim`/`insertClaimLife` (**ADR-0013**) |
| 3 | **11** | `Call SendEmailKlaimLife` |
| 4 | **12** | `Call HitServiceToKasirKMTLife_Act` — Kasir/pembayaran |

Step 9 (gerbang EXIT) berdiri **di antara** efek 1 dan 2. Step 14 `SetInformationData` **bukan**
efek keluar (lihat B3).

### Jalur penomoran akseptasi yang AKTIF

`[terverifikasi]` Dengan 4.4/4.5 mati, nomor akseptasi kini dibuat lewat rantai **aktif** di dalam
step 4:

| Sub-step | Rule | Deskripsi |
| --- | --- | --- |
| 4.1 | `RDB-List` | get tanggal produksi |
| **4.7** | `RDB-List` → `GetKodeProdLife_SQL` | **AMBIL KODE PROD** (prefix) |
| **4.9** | `RDB-List` → `GetSequenceNumber_SQL` | **generate MM.YYYY DAN SEQUENCE** |
| 4.11 | `Property-Set` | QR,QP |
| 4.12 | `Property-Set` | TR,TP |
| 4.13 | `Property-Set` | Set Param |
| **4.15** | `RDB-List` → `UpdateOsAkseptasiClaimLife_sql` | **Insert ke OS** |
| 4.17 / 4.18 | `Call PrintAkseptasiPDF` / `Call LoadDocumentLife_ACT` | |

**Ini jalur yang sama persis dengan Claim Life** (**ADR-0006**, OQ-002): `GetKodeProdLife_SQL` +
`GetSequenceNumber_SQL`. Percabangan per `Type` bertahan di 4.11/4.12.

**"Sekali di tingkat final" TETAP BENAR** `[terverifikasi]` — seluruh step 4 digerbangi
`pyWorkPage.AcceptStatus = 1 && pyWorkPage.KomiteCount == pyWorkPage.KomiteLoop` (baris **5695**);
step 5 oleh `AcceptStatus==2 && KomiteCount == KomiteLoop` (baris **8119**); step 10 dan step 12 oleh
gerbang yang sama (baris **8648**, **8887**). Mematikan 4.4/4.5 **tidak** mengubah kesimpulan itu —
ia hanya memindahkan pembuatan nomor dari rule `Generate_NoAccept_KMT_*` lama ke pasangan
`GetKodeProd` + `GetSequenceNumber`.

Ini **menguatkan** temuan Ronde 1 dengan sinyal kedua yang bebas: uji asimetri indeks menunjukkan
`Generate_NoAccept_KMT_Life` / `_LifeRetro` tidak terindeks aktif; blockname `//` kini
mengonfirmasinya secara harfiah.

---

## Bagian 1 — Lima jawaban work owner

### B1. Idempotensi & anti-dobel `[keputusan work owner]`

Keempat efek keluar **boleh di-hit berkali-kali** (retry). Syarat wajib sistem baru:

- Tiap kiriman membawa **ID idempoten unik** agar penerima dapat menolak duplikat.
- Khusus **Email** dan **Kasir**: sebelum kirim ulang, **WAJIB cek status "sudah terkirim sukses"** —
  jangan sampai email atau pembayaran dobel.

→ masuk syarat **ADR-0015**.

### B2. Gerbang EXIT retro dibuang `[keputusan work owner]`

Tiga identitas `1000013`, `L0000141`, `L0000134` beserta **gerbang EXIT step 9** dibuang.
**Konsekuensi yang diterima work owner:** klaim ber-retro tersebut yang tadinya di-EXIT (tidak
menjalankan efek keluar) kini menjalankan **seluruh** efek keluar termasuk **Kasir**. **Semua klaim
menjalankan keempat efek.** → **OQ-064 DITUTUP.**

### B3. `SetInformationData` bukan efek keluar `[keputusan work owner]`

Step 14 hanya **temporary** — menampung hasil submit (`"Approve"`/`"Reject"` + `AcceptedNo` +
`CLAIM_NO`). **Bukan** efek keluar, **bukan** penyimpanan permanen, **tidak** wajib-berhasil, tidak
perlu jaminan pengiriman.

`[terverifikasi]` Konsisten isi rule: `Komite Claim Life/Activity/SetInformationData.xml`
(`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `SETINFORMATIONDATA` / `RULE-OBJ-ACTIVITY`) mengisi
`ParamDataLife.CARI1…CARI6` dengan `AcceptedNo`, `CLAIM_NO`, `PolicyHolderName`, `Comment`, dan
literal `"Approve"` / `"Reject"`.

### B4. Pemantau kegagalan `[keputusan work owner]`

Kegagalan efek keluar yang perlu intervensi ditampilkan sebagai **status eksplisit di UI Komite** +
**laporan harian**. → **ADR-0015**.

### B5. Delegasi/eskalasi pemutus `[keputusan work owner]`

Inbox komite **hanya** menampilkan kasus sesuai posisi (roster tingkat itu). Bila anggota komite
absen, kasus **dipindahkan NAIK SATU TINGKAT** untuk diaksep tingkat di atasnya — **aksi manual
admin**.

Ini **melengkapi ADR-0014**: penegakan per `KomiteID` tetap berlaku, dengan **pengecualian sah**
berupa eskalasi manual ke tingkat lebih tinggi — bukan "siapa saja boleh".

---

## Status OQ setelah Ronde 2

| OQ | Verdict |
| --- | --- |
| **OQ-064** | **TERTUTUP** `[keputusan work owner]` — gerbang EXIT retro dibuang |
| **OQ-033**, **OQ-034** | tetap tertutup (Ronde 1) |
| **OQ-035** | terbuka, bukan pemblokir isi |
| **OQ-065** (baru) | Isi kolom retro pada rekam akseptasi setelah langkah tukar-RetroName mati |
