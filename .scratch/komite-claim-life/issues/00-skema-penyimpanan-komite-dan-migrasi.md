# 00: Skema penyimpanan komite (`T_GENERAL_KOMITE` + `T_KOMITE_KOMITELIST`) + `COVER_KEY` — **PREFACTOR**

**Status:** sebagian — DDL `013`/`030` + penjaga statik ada; belum: kolom audit `T_GENERAL_KOMITE`, `CHECK` approval, FK `ADJUSTMENT_ID`, migrasi data lama, dan `-migrate` naik/turun terhadap Oracle

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
      *(AC 30 spec; penyimpangan sadar 1)* — belum: kolom audit (operator, tanggal) tidak ada di `013_tabel_komite.sql` — namanya masih `[terbuka]` di STRUKTUR; sisanya ada (shared PK berawalan `KMTLF-`, `ADJUSTMENT_ID`, `KOMITE_LOOP`, `KOMITE_COUNT`, `ACCEPT_STATUS`)
- [x] ⚠️ ⛔ **Tidak ada kolom `WORK_CLAIM_ID` di mana pun.** Hubungan `T_WORK_CLAIM` ↔
      `T_GENERAL_KOMITE` dijamin oleh **`ID` yang identik** (shared PK), bukan oleh kolom
      penyambung. Test yang menemukan kolom `WORK_CLAIM_ID` **gagal**.
      *(REVISI 2026-09-18; `[keputusan work owner]`)* — bukti: uji `TestNamaYangDibuangTidakAda` (seluruh migrasi); di luar komentar nol kemunculan di kode Go/TS
- [x] ⚠️ **`T_GENERAL_KOMITE.ADJUSTMENT_ID` TETAP ADA** — penutup lingkar ke baris adjustment; ia
      **bukan** bagian shared PK dan **tidak** ikut dibuang. *(REVISI 2026-09-18)* — bukti: `repository/migrations/013_tabel_komite.sql:T_GENERAL_KOMITE` (`ADJUSTMENT_ID NOT NULL` + `UX_GENERAL_KOMITE_ADJ`)
- [x] ⚠️ **`T_KOMITE_KOMITELIST` ada** dengan FK **`DATA_KOMITE_ID`** → `T_GENERAL_KOMITE.ID` dan
      **`ON DELETE CASCADE`**; satu baris **per anggota per jenjang**. *(AC 30 spec)* — bukti: `repository/migrations/030_komite_kaskade_dan_lebar_id.sql:FK_KOMITELIST_KOMITE` (`ON DELETE CASCADE`), uji `TestKaskadeHanyaPadaRelasiTerdaftar`; `-migrate` belum pernah dijalankan
- [x] ⚠️ **`T_WORK_CLAIM` memuat `COVER_KEY`**, nullable dan ber-index — `ID` baris induk; `NULL`
      bila baris itu tidak punya induk. *(§9 spec; REVISI 2026-09-17)* — bukti: `repository/migrations/001_t_work_claim.sql:IX_WORK_CLAIM_COVER_KEY` (kolom nullable), uji `TestKunciTamuBerIndex`
- [x] ⚠️ **Kirim komite melahirkan BARIS BARU di `T_WORK_CLAIM`** — `ID` = identitas kasus komite,
      `COVER_KEY` = `ID` baris klaim. Kasus komite **bukan** kolom pada baris klaim.
      *(REVISI 2026-09-17; `[keputusan work owner]`)* — bukti: `repository/kasuskomite.go:PohonKlaim.BuatKasusKomite` (INSERT `T_WORK_CLAIM`: `ID` = `KMTLF-…`, `COVER_KEY` = klaim)
- [x] ⚠️ **Tidak ada kolom `KMT_NO` di mana pun.** Identitas kasus komite hidup di
      `T_WORK_CLAIM.ID`; menyimpannya lagi sebagai kolom terpisah berarti duplikasi. Test yang
      menemukan `KMT_NO` **gagal**. *(REVISI 2026-09-17)* — bukti: uji `TestNamaYangDibuangTidakAda`
- [x] ⚠️ **`ADJUSTMENT_ID` berada di `T_GENERAL_KOMITE`, BUKAN di `T_WORK_CLAIM`.** Test yang
      menemukannya di tabel work **gagal**. *(§9 spec)* — bukti: uji `TestKolomDDLCocokDenganStruktur` (kolom DDL di luar STRUKTUR gagal); `001_t_work_claim.sql` tanpa `ADJUSTMENT_ID`
- [ ] `DATE_APPROVE` bertipe **`DATE`**; `KOMITE_APROVAL` hanya menerima `0`, `1`, `2`.
      *(AC 31 spec; penyimpangan sadar 4)* — belum: `DATE_APPROVE DATE` ada, tetapi `KOMITE_APPROVAL VARCHAR2(8)` tanpa `CHECK`, dan eskalasi menulis `NULL` (`repository/komite_keputusan.go:sqlLewatiAnakTangga`)
- [x] `KOMITE_URUT` menyimpan jenjang tangga, sehingga riwayat dapat diurut **tanpa** bergantung
      urutan penyisipan baris. *(AC 33 spec)* — bukti: `repository/komite_inbox.go:sqlTanggaKasus` (`ORDER BY KOMITE_URUT, ID`)
- [x] ⚠️ Rujukan antar tabel memakai **`ID` / `COVER_KEY`**; **tidak ada** kolom yang menyimpan indeks
      posisi. Test yang menemukan padanan `IndexAdjustment` / `IndexPremiumList` **gagal**.
      *(AC 34 spec; penyimpangan sadar 2)* — bukti: uji `TestAC45Dan48RujukanKomiteBukanIndeksPosisi`, uji `TestNolPenyimpanKeputusanKomiteDiKonteksIni` (pola indeks posisi di seluruh kode Go)
- [x] Setiap FK (`DATA_KOMITE_ID`, `ADJUSTMENT_ID`) dan `COVER_KEY` **ber-index**. **REVISI
      2026-09-18:** `WORK_CLAIM_ID` dihapus dari daftar — kolomnya tidak ada; `T_GENERAL_KOMITE.ID`
      adalah PK sehingga sudah ber-index dengan sendirinya. — bukti: `013_tabel_komite.sql:IX_KOMITELIST_KOMITE` + `UX_GENERAL_KOMITE_ADJ`, `001_t_work_claim.sql:IX_WORK_CLAIM_COVER_KEY`
- [ ] Integritas rujukan ditegakkan basis data: `ADJUSTMENT_ID` yang menunjuk baris adjustment
      **tidak ada** **ditolak**; `DATA_KOMITE_ID` yatim **ditolak**. — belum: `DATA_KOMITE_ID` ber-FK, tetapi `ADJUSTMENT_ID` sengaja tanpa `REFERENCES` (dua tabel tujuan menurut `LINI`); keutuhannya hanya dijaga kode Go
- [x] ⚠️ **Tidak ada hapus fisik** baris keputusan; menghapus header komite mengkaskade ke anaknya.
      *(penyimpangan sadar 3)* — bukti: `030_komite_kaskade_dan_lebar_id.sql:FK_KOMITELIST_KOMITE` (kaskade); nol `DELETE` atas `T_GENERAL_KOMITE`/`T_KOMITE_KOMITELIST` di kode — eskalasi mengosongkan, bukan menghapus
- [x] ⚠️ **REVISI 2026-09-18 — identitas kedua tabel TIDAK seragam lagi.**
      `T_KOMITE_KOMITELIST.ID` berasal dari **sequence** (**ADR-0006**), tetapi
      `T_GENERAL_KOMITE.ID` adalah **shared PK** — ia **mengambil** `T_WORK_CLAIM.ID` baris komite
      (teks berformat `KMT-xxxxxx`), **bukan** sequence. Test yang menuntut sequence untuk
      `T_GENERAL_KOMITE.ID` **keliru** dan harus dibalik. *(penyimpangan sadar dari ADR-0006)* — bukti: `repository/kasuskomite.go:PohonKlaim.BuatKasusKomite` (anak dari `SEQ_KOMITE_KOMITELIST`; kepala = `PengenalWorkBerikut`, bukan sequence tabelnya)
- [x] ✅ Tipe `T_WORK_CLAIM.ID` **SUDAH DITETAPKAN 2026-09-18** — **teks berformat**, baris klaim
      `CLM-xxxxxx` dan baris komite `KMT-xxxxxx`; **bukan** angka sequence. `COVER_KEY`,
      `T_GENERAL_CLAIM.ID`, `T_GENERAL_KOMITE.ID`, dan `T_CLAIMLF_ADJUSTMENT.KOMITE_ID`
      **mengikuti** tipe itu. *(`[keputusan work owner]`; ⚠️ penyimpangan sadar dari **ADR-0006**)* — bukti: uji `TestAC33IdentitasBertipeSama`, uji `TestLebarKolomPenunjukSamaDenganIndukNya`, uji `TestIdentitasWorkClaimBerupaTeks` (awalan baris komite `KMTLF-` menurut STRUKTUR)
- [ ] ⚠️ `[terbuka]` **Generator nomor `CLM-`/`KMT-` belum ditetapkan** — siapa yang membuatnya,
      apakah ada sequence di belakang prefiks, apakah di-reset per tahun. Pemilik **DBA / work
      owner**. **Jangan tebak.** — belum: sebagian diputuskan (butir aa 26-09-2026: `SEQ_WORK_CLAIM` + awalan, `repository/pengenalwork.go:PohonKlaim.PengenalWorkBerikut`); reset per tahun tidak diputuskan dan STRUKTUR masih `[terbuka]`
- [x] ⚠️ `[terbuka]` **Apakah `COVER_KEY` dan `T_CLAIMLF_ADJUSTMENT.KOMITE_ID` dipasangi
      `REFERENCES T_WORK_CLAIM(ID)`** atau dibiarkan tanpa constraint — **belum diputuskan**.
      Pemilik **DBA / work owner**. Tiket ini **tidak dinyatakan selesai** sebelum jawabannya ada. — bukti: diputuskan butir d 26-09-2026 — `FK_WORK_COVER_KEY` (001) dan `FK_ADJ_KOMITE` (004), uji `TestAC50PenunjukKeAtasBerReferences`
- [ ] `POOLDATA.EMAILKOMITE` **tidak ditulis** — ia master yang dibaca. Test yang menemukan tulisan
      ke sana **gagal**. *(`[data DBA]`)* — belum: kode hanya membaca (`repository/roster.go:PohonKlaim.AmbilRosterKomite`, satu-satunya rujukan), tetapi tidak ada uji yang gagal bila ada tulisan ke `EMAILKOMITE`
- [ ] Bila ada roster/keputusan lama yang hidup, migrasi memindahkannya ke kedua tabel **tanpa
      kehilangan satu nilai pun**, dan rekonsiliasi membandingkan jumlah baris per kasus.
      *(**ADR-0009**)* — belum: tidak ada skrip migrasi roster/keputusan lama maupun rekonsiliasi jumlah barisnya
- [ ] Migrasi dapat **dijalankan ulang dengan aman** dan punya **jalur mundur yang diuji**. — belum: `TestMigrasiIdempoten` dan `TestJalurMundurDiuji` (uji db) ada, tetapi belum pernah dijalankan — `-migrate` belum pernah jalan terhadap Oracle

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

## Implementasi — 27 September 2026 (verifikasi `013` + migrasi `030`)

⚠️ **Butir km1: tiket ini VERIFIKASI, bukan pembuatan ulang.** `013_tabel_komite.sql` sudah membuat
kedua tabel di `main`. Yang dikerjakan: mencocokkan bentuknya dengan
`STRUKTUR-TABEL-KOMITE-CLAIM-LIFE.md`, menambal yang **kurang**, dan memasang penjaga.

### ⛔ Ralat tiket ini sendiri — tiga nama kolom

Bab *"Bentuk yang dibangun"* di atas menyebut `KOMITE_ID` *(anggota pemutus)*, `ID_KOMITE`
*(jabatan)*, dan `KOMITE_APROVAL`. **Ketiganya sudah diganti** `[keputusan work owner]`
2026-09-18 sore, dan STRUKTUR mencatat penggantiannya di bab *"Tiga nama kolom diganti"*:

| Nama di tiket | Nama sebenarnya | Sebab |
| --- | --- | --- |
| `KOMITE_ID` | **`KOMITE_OPERATORID`** | isinya akun **operator** (`KomiteID := .OPERATOR_ID`) |
| `ID_KOMITE` | **`KOMITE_JABATAN`** | isinya **jabatan** (`IDKomite := .JABATAN`) |
| `KOMITE_APROVAL` | **`KOMITE_APPROVAL`** | ejaan dibetulkan *(korpus satu P)* |

DDL `013` sudah memakai nama yang **benar**; yang tertinggal adalah teks tiket ini.

### Dua celah bentuk yang ditemukan, dan ditambal `030`

**Celah 1 — kaskade yang dijanjikan tetapi tidak ada.** STRUKTUR menyebut relasi
`T_GENERAL_KOMITE → T_KOMITE_KOMITELIST` sebagai `ON DELETE CASCADE` di **tiga** tempat
*(baris 159, 238, 253)*, dan tiket ini mengulanginya. `013` membuat `FK_KOMITELIST_KOMITE`
**tanpa `ON DELETE`**.

⚠️ Akibatnya nyata: menghapus satu kasus komite akan **ditolak Oracle** (ORA-02292) selama masih ada
baris roster yang menunjuknya — dan itu jalur yang **tiket 05** perlukan.

⛔ **Dan penjaga kita sendiri menegakkan cacat itu.** `TestKaskadeHanyaPadaEmpatRelasi` ditulis
ketika hanya migrasi Claim Life ada; daftarnya menuntut `013` **tidak** berkaskade. Jadi ia bukan
sekadar melewatkan cacat — **ia menahan perbaikannya**. Daftarnya diralat.

**Celah 2 — lebar kolom penunjuk berbeda dari induknya.**

| Kolom | `013` | Induknya |
| --- | --- | --- |
| `T_GENERAL_KOMITE.ID` | `VARCHAR2(40)` | `T_WORK_CLAIM.ID` `VARCHAR2(32)` |
| `T_GENERAL_KOMITE.ADJUSTMENT_ID` | `VARCHAR2(40)` | `T_CLAIMLF_ADJUSTMENT.ID` `VARCHAR2(32)` |
| `T_KOMITE_KOMITELIST.ID` / `DATA_KOMITE_ID` | `VARCHAR2(40)` | `T_GENERAL_KOMITE.ID` |

⚠️ Oracle **menerima** kunci tamu antarlebar berbeda, jadi ini tidak pernah gagal — **ia hanya
berbohong**. Kolom 40 karakter yang menunjuk kolom 32 karakter menjanjikan ruang yang tidak dapat
dipakai: nilai ke-33 sampai ke-40 tidak akan pernah punya induk. **Shared PK yang lebarnya berbeda
dari induknya bukan shared PK, melainkan kebetulan yang sedang cocok.**

### Yang DICOCOKKAN dan ternyata benar

| Hal | Hasil |
| --- | --- |
| `ADJUSTMENT_ID` **tanpa** `REFERENCES` *(dua tabel tujuan menurut `LINI`)* | ✅ sesuai |
| Index **UNIK** pada `ADJUSTMENT_ID` | ✅ `UX_GENERAL_KOMITE_ADJ` |
| `ACCEPT_STATUS` / `KOMITE_APPROVAL` **teks**, bukan angka *(ADR-U-0022)* | ✅ |
| Nol kolom `WORK_CLAIM_ID`, nol `KMT_NO` | ✅ |
| `T_WORK_CLAIM.COVER_KEY` + `LINI` ada dan nullable | ✅ migrasi `001` |
| `SEQ_KOMITE_KOMITELIST` ada *(ADR-0006 untuk `T_KOMITE_KOMITELIST.ID`)* | ✅ |
| Kesembilan nama kolom `T_KOMITE_KOMITELIST` | ✅ cocok STRUKTUR |

### Cacat di penjaga bersama yang ikut ketahuan

Dokumen STRUKTUR kedua tidak dapat ditambahkan sebelum **dua** cacat pemecahnya diperbaiki:

1. **Pemecah bab hanya mengenali `## `.** Dokumen Komite punya dua bab ber-`### ` yang tabelnya
   berisi nama kolom *(daftar ganti-nama dan daftar penutup)*, sehingga `T_KOMITE_KOMITELIST` tampak
   punya **26** kolom padahal sembilan. Komentar penjaga itu sendiri mencatat cacat serupa pernah
   terjadi — tambalannya hanya menutup separuh.
2. **Pembaca TIPE tidak pernah punya penjaga bab sama sekali.** Bab `### Daftar penutup` bertabel
   **dua** kolom *(nama + keterangan)*, sehingga **keterangan terbaca sebagai tipe** — *"penyetuju ke
   berapa"* menjadi golongan tipe, dan penjaga lalu menuduh DDL yang benar.

Dan `tipeMenurutDDL` hanya membaca `CREATE TABLE`, sehingga ia membaca bentuk saat tabel **lahir** —
bukan bentuknya sesudah seluruh migrasi. Tanpa `ALTER … MODIFY`, penjaga lebar akan menuduh skema
yang justru sudah diperbaiki `030`.

### Penjaga BARU, tiap satunya dibuat gagal lebih dulu

| Penjaga | Dibuat gagal dengan | Berbunyi |
| --- | --- | --- |
| `TestKaskadeHanyaPadaEmpatRelasi` *(diralat)* | cabut `ON DELETE CASCADE` dari `030` | `ada=false, mau=true` |
| `TestLebarKolomPenunjukSamaDenganIndukNya` | kembalikan `ADJUSTMENT_ID` ke 40 | menyebut kedua tipe |
| `TestDokumenSTRUKTURSepakatAtasTabelBersama` | sisipkan kolom karangan di satu dokumen | `digambarkan BERBEDA` |

### ⛔ Yang BELUM dijalankan

**`-migrate` belum pernah dijalankan.** `030` berisi `ALTER … MODIFY` yang **menyempitkan** kolom —
aman hanya bila tidak ada nilai lebih panjang; bila kelak ada, Oracle sendiri menolak (ORA-01441),
gagal terang. Menjalankannya menuntut skema uji dan **persetujuan manusia**.

**AC yang tercentang:** bentuk kedua tabel + shared PK + FK + index unik + sequence; `COVER_KEY`
nullable; nol `WORK_CLAIM_ID`/`KMT_NO`. **Belum:** yang menuntut basis data sungguhan.
