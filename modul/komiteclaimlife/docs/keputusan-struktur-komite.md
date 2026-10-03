# Keputusan Struktur Penyimpanan Komite Claim Life — JSON dibuang, relasional

Tanggal: 2026-09-16
Sumber: keputusan work owner + sensus korpus (KomiteList di Claim Life + Komite Claim Life).
Status: `[keputusan work owner]` — menutup lingkar penyimpanan komite yang tertunda dari revisi Claim Life.

> Konteks: spec + tiket Komite Claim Life sudah terbit (Task 1). Revisi ini menetapkan **tabel
> penyimpanan roster + keputusan komite** (belum ada di skema baru) + tutup lingkar dengan
> `T_CLAIMLF_ADJUSTMENT.KOMITE_ID`.

---

## Alur (dari work owner) `[keputusan work owner]`

> ⚠️ **REVISI 2026-09-17** `[keputusan work owner]` — **`KMT_NO` DIBUANG.** `T_WORK_CLAIM` adalah
> tabel **satu baris per work object**, dan kasus komite mendapat **barisnya sendiri**: identitasnya
> di kolom **`ID`**, induknya (baris klaim) di kolom **`COVER_KEY`**. Jadi tidak perlu kolom nomor
> terpisah. Versi sebelumnya menaruh `KMT_NO` sebagai kolom di baris klaim — itu **dibatalkan**.
> Bentuk `COVER_KEY` ini juga lebih setia ke Pega, yang menautkan lewat mekanisme *cover*
> (`pxAddChildWork`, `pzInsKey`) dan bukan lewat kolom nomor.

```
1. Di layar adjustment (Claim Life) → klik kirim komite → dibuat kasus komite
2. Lahir BARIS BARU di T_WORK_CLAIM: ID = identitas kasus komite, COVER_KEY = ID baris klaim
3. Dari baris itu → dibuat header T_GENERAL_KOMITE (WORK_CLAIM_ID = T_WORK_CLAIM.ID)
   ⚠️ **RALAT 2026-09-18:** *"(WORK_CLAIM_ID = T_WORK_CLAIM.ID)"* **DICABUT** — kolom itu
   dibuang; header memakai **`ID` yang sama persis** (**shared primary key**).
4. T_GENERAL_KOMITE punya anak T_KOMITE_KOMITELIST (list komite per jenjang + keputusan)
```

Isi `T_WORK_CLAIM` setelah kirim komite:

| `ID` | `COVER_KEY` | keterangan |
| --- | --- | --- |
| *(identitas klaim)* | `NULL` | baris klaim — tidak punya induk |
| *(identitas komite)* | *(identitas klaim)* | baris komite — menunjuk induknya |

Rantai relasi:
```
T_CLAIMLF_ADJUSTMENT  (baris klaim)   —  KOMITE_ID = identitas kasus komite (DIPERTAHANKAN: penunjuk
      │                                  langsung, supaya user langsung tahu merujuk komite mana)
      │  kirim komite → lahir baris komite di T_WORK_CLAIM
      ▼
T_WORK_CLAIM  (lintas-lini, satu baris per work object)
      │        ID = identitas kasus komite ; COVER_KEY = ID baris klaim ; posisi/status tangga
      │        (TANPA KMT_NO · TANPA ADJUSTMENT_ID)
      │  (penghubung: ID)
      ▼
T_GENERAL_KOMITE  (header data komite)   —  WORK_CLAIM_ID (→ T_WORK_CLAIM.ID),
      !!! RALAT 2026-09-18: WORK_CLAIM_ID DICABUT. Header memakai ID yang
      !!! sama persis dengan T_WORK_CLAIM.ID baris komite (shared PK).
      │                                  ADJUSTMENT_ID (→ T_CLAIMLF_ADJUSTMENT.ID),
      │                                  KOMITE_LOOP, KOMITE_COUNT, ACCEPT_STATUS
      └─ T_KOMITE_KOMITELIST  1:N  FK DATA_KOMITE_ID   ← list komite sesuai jenjang + keputusan per anggota
```

`[terbuka]` **Tipe `T_WORK_CLAIM.ID`** — teks (mis. `"KMT-0001"`) atau angka sequence, **belum
ditetapkan**. Tipe `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` mengikuti tipe itu, apa pun hasilnya. **Jangan
tebak.** Nama kolom `WORK_CLAIM_ID` di `T_GENERAL_KOMITE` adalah **nama usulan**, disetujui work owner
2026-09-17; ganti bila DBA sudah punya nama lain.

> ✅ **RALAT 2026-09-18 — seluruh paragraf di atas DICABUT.** `[keputusan work owner]`
> Tipe `T_WORK_CLAIM.ID` = **teks berformat** `CLM-xxxxxx` / `KMT-xxxxxx`, **bukan** sequence.
> Dan kolom `WORK_CLAIM_ID` **dibuang** — bukan "nama usulan yang bisa diganti DBA", melainkan
> **tidak ada**; hubungannya **shared primary key**. Rincian di blok RALAT pada
> §`T_GENERAL_KOMITE` di bawah.

`[keputusan work owner]`:
- **`T_WORK_CLAIM.ID`** = penghubung antara baris work komite dan `T_GENERAL_KOMITE`.
- **`T_WORK_CLAIM.COVER_KEY`** = `ID` baris induk. Inilah yang menyatakan kasus komite ini milik
  klaim mana. Menggantikan mekanisme *cover* Pega.
- **`ADJUSTMENT_ID` ada di `T_GENERAL_KOMITE`** (BUKAN di `T_WORK_CLAIM`) — inilah yang tahu kasus komite
  ini milik baris adjustment mana (tutup lingkar dari sisi komite).
- **`T_CLAIMLF_ADJUSTMENT.KOMITE_ID` DIPERTAHANKAN** = identitas kasus komite (penunjuk langsung/shortcut). Alasan
  work owner: **saat user melihat baris adjustment, ia langsung tahu merujuk ke komite yang mana**
  tanpa harus query terbalik. Jadi ada penunjuk DUA ARAH:
  - `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` = identitas kasus komite (= `T_WORK_CLAIM.ID` baris komite) →
    shortcut tampilan (adjustment → komite mana).
  - `T_GENERAL_KOMITE.ADJUSTMENT_ID` → tutup lingkar (komite ini milik adjustment mana).
- Kedua tabel komite (`T_GENERAL_KOMITE`, `T_KOMITE_KOMITELIST`) **lintas-lini** (bisa dipakai Non-Life).

---

## Tabel

### T_GENERAL_KOMITE (header kasus komite — PK ID, LINTAS-LINI)

> ⚠️ **RALAT 2026-09-18 — SHARED PRIMARY KEY; kolom `WORK_CLAIM_ID` DIBUANG.**
> `[keputusan work owner]` Teks di bawah **tidak dihapus** sebagai jejak, tetapi yang **mengikat**
> adalah ralat ini.
>
> 1. ⛔ **Kolom `WORK_CLAIM_ID` tidak ada** — bukan diganti nama, **dibuang**. Hubungan
>    `T_WORK_CLAIM` ↔ `T_GENERAL_KOMITE` dijamin oleh **`ID` yang identik**: `T_GENERAL_KOMITE.ID`
>    **sama persis** dengan `T_WORK_CLAIM.ID` baris komite (**shared PK**, 1:1). Begitu pula di sisi
>    klaim: `T_GENERAL_CLAIM.ID` = `T_WORK_CLAIM.ID` baris klaim.
>    Setiap penyebutan *"`WORK_CLAIM_ID` penghubung"* atau *"nama usulan"* di bawah **dicabut**.
> 2. **`T_GENERAL_KOMITE.ADJUSTMENT_ID` TETAP ADA** — penutup lingkar ke baris adjustment. Ia
>    **bukan** bagian shared PK dan **tidak** ikut dibuang.
> 3. ✅ **Tipe `T_WORK_CLAIM.ID` DIPUTUSKAN** — **teks berformat**: baris klaim `CLM-xxxxxx`
>    (contoh `CLM-123456`), baris komite `KMT-xxxxxx` (contoh `KMT-000789`). **Bukan angka
>    sequence.** `COVER_KEY`, `T_GENERAL_CLAIM.ID`, `T_GENERAL_KOMITE.ID`, dan
>    `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` **mengikuti** tipe itu. Setiap `[terbuka]` *"teks atau angka
>    sequence — belum ditetapkan"* di bawah **dicabut**.
> 4. ⚠️ **Penyimpangan sadar dari ADR-0006**, dicatat bukan dilanggar diam-diam: identitas
>    `T_WORK_CLAIM` dan kedua tabel ber-shared-PK adalah **nomor bisnis berformat**, bukan sequence.
>    ADR-0006 **tetap berlaku** untuk `T_KOMITE_KOMITELIST.ID`.
> 5. ⚠️ `[terbuka]` **TETAP terbuka, jangan tebak:** (a) **generator** nomor `CLM-`/`KMT-` — siapa
>    yang membuatnya, sequence di belakang prefiks atau tidak, reset per tahun atau tidak;
>    (b) apakah `COVER_KEY` dan `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` dipasangi
>    `REFERENCES T_WORK_CLAIM(ID)`. Pemilik keduanya **DBA / work owner**.

`ID` (PK, sequence), `WORK_CLAIM_ID` (penghubung ke `T_WORK_CLAIM.ID` — baris work komite),
  ⚠️ *baris ini **DICABUT** 2026-09-18 — `ID` **bukan** sequence melainkan **shared PK** =
  `T_WORK_CLAIM.ID` baris komite, dan `WORK_CLAIM_ID` **tidak ada**,*
`ADJUSTMENT_ID` (→ `T_CLAIMLF_ADJUSTMENT.ID`),
`KOMITE_LOOP` (jumlah tingkat = COUNT roster aktif), `KOMITE_COUNT` (tingkat berjalan),
`ACCEPT_STATUS` (hasil final: 1 aksep / 2 tolak), + audit (op/tanggal).

> `WORK_CLAIM_ID` menunjuk `T_WORK_CLAIM.ID`. `KOMITE_LOOP` = COUNT roster `EMAILKOMITE` aktif
> ber-`LIMIT_BOTTOM <= |CLAIM_AMOUNT|` (terverifikasi, spec Komite). `KOMITE_COUNT` naik per tingkat.

### T_KOMITE_KOMITELIST (list komite per jenjang + keputusan — 1:N, INDUK = T_GENERAL_KOMITE)
`ID` (PK), `DATA_KOMITE_ID` (FK → T_GENERAL_KOMITE.ID, ON DELETE CASCADE),
`KOMITE_URUT` (jenjang/tingkat tangga; urutan roster),
`KOMITE_ID` (← `.OPERATOR_ID` di `CreateKMTLife_Act` — anggota pemutus),
`ID_KOMITE` (← `JABATAN`), `KOMITE_EMAIL` (← `EMAIL`),
`KOMITE_APROVAL` (0 belum / 1 setuju / 2 tolak — ← `AcceptStatus`),
`KOMITE_COMMENT` (← `Comment`), `DATE_APPROVE` DATE (← `@CurrentDateTime()`).

> `[terverifikasi sensus]` field roster ber-class `ASM-FW-GCNMFW-Data-Comitee`:
> - `Claim Life/Activity/CreateKMTLife_Act.xml` (`CREATEKMTLIFE_ACT`) meng-`Property-Set`
>   `KomiteList(<APPEND>).KomiteID = .OPERATOR_ID`, `.KomiteAproval = 0` (awal), `.KomiteEmail = .EMAIL`.
> - `Komite Claim Life/Activity/KomitePostAdjustment.xml` (`KOMITEPOSTADJUSTMENT`) saat putus meng-set
>   `KomiteList(...).KomiteAproval = pyWorkPage.AcceptStatus`, `.KomiteComment = pyWorkPage.Comment`,
>   `.DateApprove = @CurrentDateTime()`. `AcceptStatus` ber-class `ASM-FW-GCNMFW-Work-KomiteLife`.
> - `IDKomite`/`KomiteAproval`/`DateApprove`/`KomiteComment` juga muncul sebagai kolom tampilan di
>   `ClaimComite.xml` / `AdjustmentDetail_Section.xml` / `Committe_Life.xml`.

### T_WORK_CLAIM (lintas-lini — sudah ada di Claim Life; DITAMBAH untuk komite)
**Satu baris per work object.** `ID` (identitas work object), **`+COVER_KEY`** (`ID` baris induk;
`NULL` bila tidak punya induk), + posisi/status tangga. Kolom final lintas-lini menyusul saat
Non-Life.

⚠️ **REVISI 2026-09-17** `[keputusan work owner]` — **`KMT_NO` DIBUANG, diganti `COVER_KEY`.**
Saat kirim komite, kasus komite **tidak** ditulis sebagai kolom pada baris klaim; ia mendapat
**baris sendiri** — `ID` = identitas kasus komite, `COVER_KEY` = `ID` baris klaim. Menambahkan
`KMT_NO` sekarang berarti menyimpan identitas yang sama dua kali. Test yang menemukan kolom
`KMT_NO` di mana pun **gagal**.

⚠️ **KOREKSI 2026-09-17** `[keputusan work owner]` — `T_WORK_CLAIM` **TANPA `ADJUSTMENT_ID`**.
Penutup lingkar ke baris adjustment ada di **`T_GENERAL_KOMITE.ADJUSTMENT_ID`**, bukan di sini.
Versi sebelumnya menulis "Sudah punya `ADJUSTMENT_ID`" — itu **salah tulis**, bertentangan dengan
diagram rantai di atas ("TANPA ADJUSTMENT_ID") dan dengan tiga sumber lain yang sepakat:
`.scratch/komite-claim-life/issues/00-skema-penyimpanan-komite-dan-migrasi.md` (diagram + AC
"`ADJUSTMENT_ID` berada di `T_GENERAL_KOMITE`, BUKAN di `T_WORK_CLAIM`"), `.scratch/claim-life/spec.md`
§2b (kolom `T_WORK_CLAIM` = `pyPosition`, `AcceptStatus`, `SendtoAdmin`, `SendtoMedical`, `Type`),
dan `.scratch/claim-life/issues/14-skema-relasional-klaim-dan-migrasi.md` (daftar kolom wajib).

---

## Tutup lingkar dengan Claim Life
- `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` = identitas kasus komite **DIPERTAHANKAN** (penunjuk langsung,
  supaya user tahu baris adjustment ini merujuk komite mana). BUKAN redundan — ini shortcut
  tampilan yang disengaja.
- Rantai penuh untuk lihat keputusan komite dari sebuah adjustment:
  ```
  T_CLAIMLF_ADJUSTMENT.KOMITE_ID
    → T_WORK_CLAIM   (ID = KOMITE_ID ; COVER_KEY = ID baris klaim)
    → T_GENERAL_KOMITE  (WORK_CLAIM_ID = T_WORK_CLAIM.ID, ADJUSTMENT_ID = adjustment.ID)
      !!! RALAT 2026-09-18: ID = T_WORK_CLAIM.ID (shared PK); WORK_CLAIM_ID dicabut
    → T_KOMITE_KOMITELIST (list per jenjang + keputusan)
  ```
- Penunjuk dua arah dijaga konsisten: isi `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` dan
  `T_GENERAL_KOMITE.ADJUSTMENT_ID` di transaksi yang sama saat kirim komite.

## Penyimpangan sadar (⚠️)
1. Roster + keputusan komite → tabel relasional (`T_GENERAL_KOMITE` + `T_KOMITE_KOMITELIST`); JSON/page Pega
   (`ClaimData.AdjustmentList.KomiteList`) dibuang.
2. Referensi pakai `ID` / `COVER_KEY` stabil, bukan index posisi Pega (`IndexPremiumList`/`IndexAdjustment`).
3. Lintas-lini (Life + Non-Life).

## OQ tersisa
- `EMAILKOMITE` = master roster (sumber KomiteID/Email/Jabatan) — dibaca, bukan tabel baru. DDL sudah
  `[data DBA]` (Task 1).
- Kolom lintas-lini final `T_WORK_CLAIM`/`T_GENERAL_KOMITE` — saat Non-Life.

## Dampak artefak
- Revisi spec + tiket **Komite Claim Life**: tambah `T_GENERAL_KOMITE` + `T_KOMITE_KOMITELIST` +
  `T_WORK_CLAIM.COVER_KEY` (**bukan** `KMT_NO` — lihat REVISI 2026-09-17).
- Selaraskan **Claim Life** tiket 10/14: `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` = identitas kasus komite,
  DIPERTAHANKAN (penunjuk langsung); rantai penuh keputusan lewat `T_WORK_CLAIM` (`ID`/`COVER_KEY`)
  → `T_GENERAL_KOMITE` → `T_KOMITE_KOMITELIST`.
