# 00: Skema penyimpanan komite (`T_GENERAL_KOMITE` + `T_KOMITE_KOMITELIST`) + `COVER_KEY` — **PREFACTOR**

**Status:** ready-for-agent

**Blocked by:** **Claim Life tiket `14`** (skema klaim — `T_CLAIMLF_ADJUSTMENT` dan `T_WORK_CLAIM`
harus ada lebih dulu; tabel di sini merujuk keduanya)

⚠️ **PREFACTOR dan tiket PERTAMA konteks ini.** Diberi nomor `00` supaya berada di depan tanpa
menomori ulang sepuluh tiket yang sudah terbit. **Tiket 01, 02, 04b, dan 09 memblokir pada tiket
ini.**

## Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin **roster komite dan keputusan tiap anggota tersimpan di tabel
relasional** yang dapat ditelusuri per anggota per jenjang — supaya riwayat persetujuan dapat
dipertanggungjawabkan dan tidak hilang bersama halaman kerja. Dan sebagai **pengguna**, dari sebuah
baris adjustment saya ingin **langsung tahu kasus komite mana** yang memutuskannya.
*(User story 41–42 di spec)*

⚠️ **Lubang yang ditutup.** Sampai sebelum ini, roster dan keputusan komite hanya hidup sebagai
**page runtime Pega** (`ClaimData…AdjustmentList.KomiteList`) — tidak punya tabel di skema baru.

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `migrations/` | DDL `T_GENERAL_KOMITE` + `T_KOMITE_KOMITELIST` (PK sequence, FK, cascade, index); `ALTER T_WORK_CLAIM` (+`COVER_KEY`) |
| `internal/models` | Kasus komite: header → baris per anggota per jenjang |
| — | Skrip migrasi roster & keputusan lama (bila ada data hidup) + rekonsiliasi |

## Bentuk yang dibangun — spec §9

> ⚠️ **REVISI 2026-09-17** `[keputusan work owner]` — **`KMT_NO` DIBUANG, diganti `COVER_KEY`.**
> `T_WORK_CLAIM` adalah tabel **satu baris per work object**. Kasus komite mendapat **barisnya
> sendiri**: `ID` = identitas kasus komite, `COVER_KEY` = `ID` baris klaim. Kolom nomor terpisah
> tidak diperlukan lagi.

```
T_CLAIMLF_ADJUSTMENT (milik Claim Life) — KOMITE_ID = identitas kasus komite   ⬅ pintasan tampilan
     │   kirim komite → LAHIR BARIS BARU di T_WORK_CLAIM
     ▼
T_WORK_CLAIM (lintas-lini, satu baris per work object)
     │   ID = identitas kasus komite ; COVER_KEY = ID baris klaim ; posisi/status tangga
     │   (TANPA KMT_NO · TANPA ADJUSTMENT_ID)
     ▼   penghubung: ID YANG SAMA (shared primary key) — tanpa kolom penyambung
T_GENERAL_KOMITE (header kasus komite) — ID = T_WORK_CLAIM.ID baris komite
     └─ T_KOMITE_KOMITELIST   1:N   FK DATA_KOMITE_ID → T_GENERAL_KOMITE.ID   ON DELETE CASCADE
```

**`T_GENERAL_KOMITE`** — `ID` (**PK = `T_WORK_CLAIM.ID` baris komite, shared PK, teks berformat
`KMT-xxxxxx`** — **bukan** sequence, dan **tanpa** kolom `WORK_CLAIM_ID`), `ADJUSTMENT_ID` (→ `T_CLAIMLF_ADJUSTMENT.ID`),
`KOMITE_LOOP` (jumlah tingkat), `KOMITE_COUNT` (tingkat berjalan), `ACCEPT_STATUS` (`1` aksep /
`2` tolak), + audit (operator, tanggal).

**`T_KOMITE_KOMITELIST`** — `ID` (PK), `DATA_KOMITE_ID` (FK), `KOMITE_URUT` (jenjang), `KOMITE_ID`
(anggota pemutus), `ID_KOMITE` (jabatan), `KOMITE_EMAIL`, `KOMITE_APROVAL` (`0`/`1`/`2`),
`KOMITE_COMMENT`, `DATE_APPROVE` (**DATE**).

**`T_WORK_CLAIM`** — satu baris per work object: `ID` (identitas work object), **+`COVER_KEY`**
nullable (`ID` baris induk; `NULL` bila tidak punya induk).

✅ **Tipe `T_WORK_CLAIM.ID` — DIPUTUSKAN 2026-09-18, `[terbuka]` tipe DITUTUP.**
`[keputusan work owner]` Ia **teks berformat**: baris klaim `CLM-xxxxxx` (contoh `CLM-123456`),
baris komite `KMT-xxxxxx` (contoh `KMT-000789`). **Bukan angka sequence.** Kolom yang mengikuti
tipe ini: `COVER_KEY`, `T_GENERAL_CLAIM.ID` (shared PK), `T_GENERAL_KOMITE.ID` (shared PK), dan
`T_CLAIMLF_ADJUSTMENT.KOMITE_ID`.

⚠️ **Penyimpangan sadar dari ADR-0006 — dicatat, bukan dilanggar diam-diam.**
`[keputusan work owner]` **ADR-0006** menetapkan identitas dari **sequence**. `T_WORK_CLAIM` dan
kedua tabel ber-shared-PK memakai **nomor bisnis berformat**. ADR-0006 **tetap berlaku** untuk
`T_KOMITE_KOMITELIST.ID`, yang tetap dari sequence.

⛔ **Kolom `WORK_CLAIM_ID` DIBUANG 2026-09-18** `[keputusan work owner]` — bukan diganti nama,
**tidak ada**. Hubungan `T_WORK_CLAIM` ↔ `T_GENERAL_KOMITE` dijamin oleh **ID yang identik**
(shared PK), bukan oleh kolom penyambung. Teks lama menyebutnya *"nama usulan, ganti bila DBA
punya nama lain"* — **dicabut**; kolomnya tidak dipakai sama sekali.

⚠️ `[terbuka]` **Dua hal TETAP terbuka. Jangan tebak:** (a) **generator** nomor `CLM-`/`KMT-` —
siapa yang membuatnya, apakah ada sequence di belakang prefiks, apakah di-reset per tahun;
(b) apakah `COVER_KEY` dan `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` dipasangi
`REFERENCES T_WORK_CLAIM(ID)`. Pemilik keduanya **DBA / work owner**.

> ⚠️ **RALAT 2026-09-18 — `T_WORK_CLAIM` kini AKAR pohon lintas-modul, dan bertambah empat kolom.**
> `[keputusan work owner]` Perubahan ini **milik Claim Life tiket 14** (yang membuat tabelnya);
> dicatat di sini supaya kedua modul tidak berselisih. Tiket **00 Komite** tetap hanya
> **`ALTER`**-menambah penaut induk-anak — cakupannya **tidak berubah**.
>
> 1. **`T_WORK_CLAIM` adalah akar**, bukan tabel di samping pohon klaim. Baris **klaim** dan baris
>    **kasus komite** sama-sama hidup di dalamnya, dibedakan oleh penaut induk-anak. Pohon enam
>    tingkat yang berlaku ada di `.scratch/claim-life/spec.md` §2b RALAT D.
> 2. **+`LINI`** (kolom **ADA**; nilai Life = konstanta lini Life — `[terbuka — Non-Life]` hanya
>    **daftar enum lintas-lini**) dan **+`CREATE_OP`**, **`CREATE_OP_NAME`**, **`TGL_UPDATE`**,
>    **`CASEID`** — keempatnya **pindahan** dari header klaim. Dibuat oleh tiket 14.
> 3. ✅ **Ejaan penaut induk-anak — DIPUTUSKAN 2026-09-18, `[terbuka]` DITUTUP.**
>    `[keputusan work owner]` **Ejaan final `COVER_KEY` (snake_case);** `CoverKey` warisan Pega
>    **tidak dipakai** sebagai nama kolom. Alasan: seluruh kolom SQL baru proyek ini snake_case
>    (`DATA_KOMITE_ID`, `ADJUSTMENT_ID`, `KOMITE_URUT`, dst), dan Oracle melipat identifier
>    tanpa kutip menjadi huruf besar — `CoverKey` akan menjadi identifier **tanpa garis bawah**,
>    berbeda dari kolom-kolom sekitarnya. Konsistensi SQL menang. `[terverifikasi]` `CoverKey`
>    **nol kemunculan di korpus Pega**, jadi tidak ada fakta korpus yang dikorbankan.
> 4. ⚠️ **`T_CLAIM_POLICY` dan `T_CLAIM_MARKETING` DIHAPUS** di Claim Life (bukan diganti nama).
>    Tidak ada akibat bagi skema komite; dicatat agar tidak dicari.

## Rule Pega sumber

Field roster ber-class **`ASM-FW-GCNMFW-Data-Comitee`**.

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `CreateKMTLife_Act` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `CREATEKMTLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Claim Life/Activity/CreateKMTLife_Act.xml` | `[terverifikasi]` membentuk roster: `KomiteList(<APPEND>).KomiteID = .OPERATOR_ID`, `.KomiteAproval = 0`, `.KomiteEmail = .EMAIL` |
| `KomitePostAdjustment` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY` | `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `[terverifikasi]` saat putus: `KomiteList(Local.Komite).KomiteAproval = pyWorkPage.AcceptStatus`, `.KomiteComment = pyWorkPage.Comment`, `.DateApprove = @CurrentDateTime()` |

**Pemetaan kolom → field** `[terverifikasi]`: `KOMITE_ID`←`KomiteID` · `ID_KOMITE`←`IDKomite` ·
`KOMITE_EMAIL`←`KomiteEmail` · `KOMITE_APROVAL`←`KomiteAproval` · `KOMITE_COMMENT`←`KomiteComment` ·
`DATE_APPROVE`←`DateApprove`.

`[terverifikasi]` `AcceptStatus` ber-class `ASM-FW-GCNMFW-Work-KomiteLife`.

`[data DBA]` **`POOLDATA.EMAILKOMITE` tetap master yang DIBACA** — sumber `KomiteID`, email,
jabatan, dan pita nilai (`LIMIT_BOTTOM`/`LIMIT_TOP`, `STS_AKTIF`, `STS_KLAIM`). **Bukan** tabel
baru, **tidak** diganti.

⚠️ **Penyimpangan sadar 1 — roster & keputusan jadi tabel relasional.** `[keputusan work owner]`
Page/JSON Pega dibuang.

⚠️ **Penyimpangan sadar 2 — rujukan `ID` / `COVER_KEY`, bukan indeks posisi.** `[terverifikasi]`
`CreateKMTLife_Act` memakai `IndexAdjustment` / `IndexPremiumList` (`.pxListSubscript`) — rusak
begitu urutan baris bergeser.

⚠️ **Penyimpangan sadar 3 — tidak ada hapus fisik**; integritas lewat PK/FK + cascade header → anak.

⚠️ **Penyimpangan sadar 4 — `DATE_APPROVE` bertipe `DATE`**; kolom uang tetap desimal presisi
arbitrer (**ADR-0003**).

## ADR terkait

**ADR-0003** (uang non-float), **ADR-0006** (identitas lewat sequence), **ADR-0007** (jejak audit),
**ADR-0009** (migrasi penuh), **ADR-0014** (pemutus ditegakkan per `KomiteID` tingkat berjalan).

## Acceptance criteria

- [ ] ⚠️ **`T_GENERAL_KOMITE` ada** dengan **`ID` = `T_WORK_CLAIM.ID` baris komite (shared primary
      key, teks berformat `KMT-xxxxxx`)** beserta `ADJUSTMENT_ID`, `KOMITE_LOOP`, `KOMITE_COUNT`,
      `ACCEPT_STATUS`, dan audit. **REVISI 2026-09-18:** PK-nya **bukan** sequence.
      *(AC 30 spec; penyimpangan sadar 1)*
- [ ] ⚠️ ⛔ **Tidak ada kolom `WORK_CLAIM_ID` di mana pun.** Hubungan `T_WORK_CLAIM` ↔
      `T_GENERAL_KOMITE` dijamin oleh **`ID` yang identik** (shared PK), bukan oleh kolom
      penyambung. Test yang menemukan kolom `WORK_CLAIM_ID` **gagal**.
      *(REVISI 2026-09-18; `[keputusan work owner]`)*
- [ ] ⚠️ **`T_GENERAL_KOMITE.ADJUSTMENT_ID` TETAP ADA** — penutup lingkar ke baris adjustment; ia
      **bukan** bagian shared PK dan **tidak** ikut dibuang. *(REVISI 2026-09-18)*
- [ ] ⚠️ **`T_KOMITE_KOMITELIST` ada** dengan FK **`DATA_KOMITE_ID`** → `T_GENERAL_KOMITE.ID` dan
      **`ON DELETE CASCADE`**; satu baris **per anggota per jenjang**. *(AC 30 spec)*
- [ ] ⚠️ **`T_WORK_CLAIM` memuat `COVER_KEY`**, nullable dan ber-index — `ID` baris induk; `NULL`
      bila baris itu tidak punya induk. *(§9 spec; REVISI 2026-09-17)*
- [ ] ⚠️ **Kirim komite melahirkan BARIS BARU di `T_WORK_CLAIM`** — `ID` = identitas kasus komite,
      `COVER_KEY` = `ID` baris klaim. Kasus komite **bukan** kolom pada baris klaim.
      *(REVISI 2026-09-17; `[keputusan work owner]`)*
- [ ] ⚠️ **Tidak ada kolom `KMT_NO` di mana pun.** Identitas kasus komite hidup di
      `T_WORK_CLAIM.ID`; menyimpannya lagi sebagai kolom terpisah berarti duplikasi. Test yang
      menemukan `KMT_NO` **gagal**. *(REVISI 2026-09-17)*
- [ ] ⚠️ **`ADJUSTMENT_ID` berada di `T_GENERAL_KOMITE`, BUKAN di `T_WORK_CLAIM`.** Test yang
      menemukannya di tabel work **gagal**. *(§9 spec)*
- [ ] `DATE_APPROVE` bertipe **`DATE`**; `KOMITE_APROVAL` hanya menerima `0`, `1`, `2`.
      *(AC 31 spec; penyimpangan sadar 4)*
- [ ] `KOMITE_URUT` menyimpan jenjang tangga, sehingga riwayat dapat diurut **tanpa** bergantung
      urutan penyisipan baris. *(AC 33 spec)*
- [ ] ⚠️ Rujukan antar tabel memakai **`ID` / `COVER_KEY`**; **tidak ada** kolom yang menyimpan indeks
      posisi. Test yang menemukan padanan `IndexAdjustment` / `IndexPremiumList` **gagal**.
      *(AC 34 spec; penyimpangan sadar 2)*
- [ ] Setiap FK (`DATA_KOMITE_ID`, `ADJUSTMENT_ID`) dan `COVER_KEY` **ber-index**. **REVISI
      2026-09-18:** `WORK_CLAIM_ID` dihapus dari daftar — kolomnya tidak ada; `T_GENERAL_KOMITE.ID`
      adalah PK sehingga sudah ber-index dengan sendirinya.
- [ ] Integritas rujukan ditegakkan basis data: `ADJUSTMENT_ID` yang menunjuk baris adjustment
      **tidak ada** **ditolak**; `DATA_KOMITE_ID` yatim **ditolak**.
- [ ] ⚠️ **Tidak ada hapus fisik** baris keputusan; menghapus header komite mengkaskade ke anaknya.
      *(penyimpangan sadar 3)*
- [ ] ⚠️ **REVISI 2026-09-18 — identitas kedua tabel TIDAK seragam lagi.**
      `T_KOMITE_KOMITELIST.ID` berasal dari **sequence** (**ADR-0006**), tetapi
      `T_GENERAL_KOMITE.ID` adalah **shared PK** — ia **mengambil** `T_WORK_CLAIM.ID` baris komite
      (teks berformat `KMT-xxxxxx`), **bukan** sequence. Test yang menuntut sequence untuk
      `T_GENERAL_KOMITE.ID` **keliru** dan harus dibalik. *(penyimpangan sadar dari ADR-0006)*
- [ ] ✅ Tipe `T_WORK_CLAIM.ID` **SUDAH DITETAPKAN 2026-09-18** — **teks berformat**, baris klaim
      `CLM-xxxxxx` dan baris komite `KMT-xxxxxx`; **bukan** angka sequence. `COVER_KEY`,
      `T_GENERAL_CLAIM.ID`, `T_GENERAL_KOMITE.ID`, dan `T_CLAIMLF_ADJUSTMENT.KOMITE_ID`
      **mengikuti** tipe itu. *(`[keputusan work owner]`; ⚠️ penyimpangan sadar dari **ADR-0006**)*
- [ ] ⚠️ `[terbuka]` **Generator nomor `CLM-`/`KMT-` belum ditetapkan** — siapa yang membuatnya,
      apakah ada sequence di belakang prefiks, apakah di-reset per tahun. Pemilik **DBA / work
      owner**. **Jangan tebak.**
- [ ] ⚠️ `[terbuka]` **Apakah `COVER_KEY` dan `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` dipasangi
      `REFERENCES T_WORK_CLAIM(ID)`** atau dibiarkan tanpa constraint — **belum diputuskan**.
      Pemilik **DBA / work owner**. Tiket ini **tidak dinyatakan selesai** sebelum jawabannya ada.
- [ ] `POOLDATA.EMAILKOMITE` **tidak ditulis** — ia master yang dibaca. Test yang menemukan tulisan
      ke sana **gagal**. *(`[data DBA]`)*
- [ ] Bila ada roster/keputusan lama yang hidup, migrasi memindahkannya ke kedua tabel **tanpa
      kehilangan satu nilai pun**, dan rekonsiliasi membandingkan jumlah baris per kasus.
      *(**ADR-0009**)*
- [ ] Migrasi dapat **dijalankan ulang dengan aman** dan punya **jalur mundur yang diuji**.

## Blocker

**Tidak ada pemblokir.**

⚠️ `[terbuka — Non-Life]` **tidak memblokir:** kolom lintas-lini final `T_WORK_CLAIM` dan
`T_GENERAL_KOMITE` ditetapkan saat konteks **Non-Life** digarap. Untuk Komite Life yang mengikat hanya
kolom di atas.

## Catatan

⚠️ **Penunjuk dua arah harus konsisten.** `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` (pintasan tampilan) dan
`T_GENERAL_KOMITE.ADJUSTMENT_ID` (penutup lingkar) diisi dalam **satu transaksi** saat kirim komite —
AC penegakannya ada di tiket **01** dan **04b**. Tiket ini menyiapkan **kolom dan constraint**-nya.

⚠️ **`STS_REJECT` berbeda tipe antar tabel** `[data DBA]`: `VARCHAR2(15)` di `EMAILKOMITE` versus
`NUMBER(38)` di `OS_AKSEPTASI_KLAIM_LIFE`. Jangan menyamakan keduanya saat memetakan.

## Seam & perintah verifikasi

**Seam: API HTTP** — memakai ulang seam Claim — Life, **tidak menambah seam baru** (spec §Testing).
Diuji terhadap **skema uji Oracle nyata**: cascade, integritas rujukan, dan konsistensi penunjuk dua
arah **hanya berperilaku benar pada basis data sungguhan**.

```
go test ./internal/...
cd frontend && npm test
make check
```
