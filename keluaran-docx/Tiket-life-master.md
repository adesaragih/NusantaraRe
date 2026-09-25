# Tiket - Data Induk Jiwa

> Dokumen ini memuat **badan tiket lengkap**, disusun per modul lalu per nomor.
> Disusun 25 September 2026 dari berkas tiket proyek migrasi Nusantara Re.

## Matriks status

| Modul | Tiket | Siap | Tertahan | needs-info | wontfix | Lain |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| Master Product Name Life | **9** | 9 | 0 | 0 | 0 | 0 |
| Master Contract Retro Life | **13** | 13 | 0 | 0 | 0 | 0 |
| **Jumlah** | **22** | **22** | **0** | **0** | **0** | **0** |

---

# Master Product Name Life

Jumlah tiket: **9**

## Master Product Name Life - 01 - Skema relasional penuh + migrasi JSON → kolom — **PREFACTOR**

**Status:** ready-for-agent

**Blocked by:** **CL-01** (kerangka aplikasi + seam API — scaffolding lintas konteks, tidak dibuat
di sini)

⚠️ **Ini PREFACTOR, dan ia tiket PERTAMA — bukan tiket migrasi yang biasanya terakhir.** Bentuk
skema berubah total (JSON → relasional), sehingga **tidak ada slice lain yang dapat berdiri**
sebelum bentuk barunya ada. *"Make the change easy, then make the easy change."*

#### Hasil & nilai pengguna

Sebagai **tim migrasi**, saya ingin setiap atribut produk menjadi **kolom bernama** dan setiap daftar
bersarang menjadi **tabel anak**, dengan data lama pindah tanpa kehilangan satu nilai pun — sehingga
bentuk produk dapat diperiksa, dicari, dan divalidasi, bukan tersembunyi di dalam satu dokumen.
*(User story 29–32 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `migrations/` | DDL **satu tabel induk** `product_life` + **lima** tabel anak + sequence; skrip migrasi JSON → relasional (gabung dua tabel existing jadi satu) |
| `internal/models` | Bentuk produk (sisi umum + sisi inward) dan **lima** jenis baris anak |
| — | Skrip rekonsiliasi |

#### Keadaan lama `[data DBA]`

| Tabel | Bentuk sekarang |
| --- | --- |
| `POOLDATA.M_PRODUCTINWARD_LIFE` | **hanya** `ID` + `JSONDATA` |
| `POOLDATA.M_PRODUCT_LIFE` | `JSONDATA` (constraint `IS JSON`) + empat kolom hasil flatten: `RIRISKID`, `RIRISK`, `PRODUCTNAME`, `BEGIN_DATE` |
| Kunci | **tanpa PK** — hanya index |

⚠️ **Penyimpangan sadar 2 — skema relasional penuh.** `[keputusan work owner]` Seluruh atribut
menjadi kolom bernama; seluruh daftar bersarang menjadi tabel anak; data JSON dimigrasikan.

#### Bentuk baru

```
product_life   (ID + ~20 field umum + 40 field inward, field duplikat → satu kolom;
                TERMASUK LIENCLAUSE skalar)
  ├─ product_life_comment          ← CommentList
  ├─ product_life_document_claim   ← DocumentClaim
  ├─ product_life_plan             ← PlanList
  ├─ product_life_uw_limit         ← UnderwritingLimitList
  └─ product_life_fin_uw           ← FinancialUnderwritingList
```

⚠️ **Penyimpangan sadar (baru) — dua tabel induk existing DIGABUNG jadi satu `product_life`.**
`[keputusan work owner]` Pega memisah `M_PRODUCT_LIFE` + `M_PRODUCTINWARD_LIFE` (dua JSON blob, 1:1
by `ID`). Karena JSON dibuang & simpan atomik, pemisahan tak beralasan lagi; gabung menghapus kelas
bug "produk timpang". Field yang muncul di kedua sisi (`ID`, `CEDING`, `POLICYHOLDER`, `POLICYHODERNAME`,
`BIRTHDAY`, `TREATYNUMBER`) → **satu kolom**.

`[terverifikasi]` **Lima tabel anak.** Dua yang semula disangka tabel dicoret:
- **`LienClause` → BUKAN tabel** — `LIENCLAUSE` adalah **field skalar** (kolom di `product_life`).
- **`OutwardList` → BUKAN tabel** — hanya diisi `GetReinsTypeOR_Life` (jalur OR **mati** `1==2`).

##### Kolom sisi inward — **40 field** `[terverifikasi]`

Sumber kebenaran: field yang **di-SET** di `SetProductNameInward`
(`ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SETPRODUCTNAMEINWARD` / `RULE-OBJ-ACTIVITY`):

```
ID, PRODUCTID, BEGIN, MATURE, STNC, BIRTHDAY, MONTHS, CURRENCY, PAYMENT, MAXCONTRACT,
MINAGE, MAXAGE, MINSUMINSURED, MAXSUMINSURED, MAXSUMREASURED,
CEDING, CEDINGLIMIT, CEDINGLIMITXPN, CEDINGRETENTIONNUM, CEDINGRETENTIONPCT,
RNMLIMITNUM, RNMLIMITPCT, RNMSHARE, BROKERAGE, EXTRAPREMI, EXTRAMORTALITY,
MAXDATARECEIVE, MAXEXPIREDCLAIM, INSURED, SUBJECTTO,
ADDENDUMNO, ADDENDUMWORD, AMANDEMENTNO, AMANDEMENTSCHD,
INWARDTREATYNM, PROPORTIONALTABLE, LIENCLAUSE, TREATYNUMBER,
POLICYHODER (sic), POLICYHODERNAME
```

⚠️ **Keempat puluh dipakai** `[keputusan work owner]` — termasuk **`ADDENDUMNO`**,
**`AMANDEMENTNO`**, dan **`CEDING`** yang hilang dari rekap catatan sumber.

##### Kolom sisi umum — **20 field** `[terverifikasi]`

Dari `SetProductName` (`ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SETPRODUCTNAME` / `RULE-OBJ-ACTIVITY`):

```
ID, PRODUCTNAME, PRODUCTCODE, PRODUCTTYPE, TYPE, GRUP, INWARDNAME,
CEDING, CAUSE, RIRISK, RIRATE, RICOMM, BENEFIT, BIRTHDAY, PAYMENTTYPE,
ISFACULTATIVE, IsView, OUTWARDNAME, OUTWARDRATE, OUTWARDCOMM
```

Ditambah lima field yang di-set `SaveProductName_Act` (`ASM-FW-GISFW-INT-PRODUCT_LIFE` /
`SAVEPRODUCTNAME_ACT`): `ID`, `CREATEOP`, `UPDATEOP`, `POLICYHODER` (sic), `POLICYHODERNAME`.

##### Kolom tabel anak `[data DBA]`

| Tabel | Kolom baris |
| --- | --- |
| `product_life_comment` | tanggal, nama operator, isi saran |
| `product_life_document_claim` | nama dokumen (+ audit) |
| `product_life_plan` | plan, id plan, nama, benefit, **RI Rate + id-nya** |
| `product_life_uw_limit` | deskripsi, status medis (`FCL`/`NM`/`MEDIS`), usia min/maks, uang pertanggungan min/maks |
| `product_life_fin_uw` | `[terverifikasi]` dari `CopyFinancialWriting`: `MinInsured`, `MaxInsured`, `Employee`, `Non_Employee` |

##### Tipe kolom `[data DBA]`

| Kelompok | Tipe |
| --- | --- |
| Uang: `*LIMIT*`, `*SUMINSURED*`, `*SUMREASURED*`, batas uang pada baris underwriting | **desimal presisi arbitrer** (**ADR-U-0003**) |
| Persen: `*PCT`, `RNMSHARE`, `BROKERAGE`, `RICOMM` | desimal |
| Tanggal: `BEGIN`, `MATURE`, `STNC`, `BIRTHDAY` | **`DATE`** |
| Usia, jumlah hari, jumlah kontrak | bilangan bulat |
| Teks panjang: `INSURED`, `SUBJECTTO`, `INWARDNAME` | teks besar |

##### Aturan pemilihan kolom `[keputusan work owner]`

1. **Kolom = field yang di-SET di Activity simpan.** Itu data yang benar-benar diisi.
2. Field yang **hanya** muncul di Section dengan visibilitas mati (`1=2` / `1==2` / `never`) **dan
   tidak di-set Activity** — **jangan diambil**. Tampilan mati bukan data.

#### ADR terkait

**ADR-U-0003** (uang non-float), **ADR-U-0006** (identitas lewat sequence basis data),
**ADR-U-0009** (migrasi penuh; koeksistensi ditolak).

#### Acceptance criteria

- [ ] ⚠️ Skema target **relasional penuh**: setiap atribut menjadi **kolom bernama**, setiap daftar
      bersarang menjadi **tabel anak**. Test yang menemukan kolom JSON sebagai penyimpan atribut
      produk **gagal**. *(AC 40 spec; `[keputusan work owner]` — penyimpangan sadar 2)*
- [ ] Daftar kolom sisi inward memuat **keempat puluh field** yang di-set Activity — termasuk
      `ADDENDUMNO`, `AMANDEMENTNO`, dan `CEDING`. *(AC 41 spec)*
- [ ] Field yang hanya tampil dengan **visibilitas mati** dan tidak di-set Activity **tidak** menjadi
      kolom. *(AC 42 spec)*
- [ ] ⚠️ **Seluruh kolom hasil migrasi NULLABLE.** Test yang menemukan `NOT NULL` pada kolom hasil
      migrasi **gagal**. Wajib-isi ditegakkan **di Go**. *(AC 43, 17 spec; `[keputusan work owner]`)*
- [ ] ⚠️ Field yang **absen** di dokumen JSON lama menjadi **kolom kosong**, bukan kegagalan migrasi.
      *(AC 44 spec; `[keputusan work owner]` — field bernilai kosong tidak muncul di JSON)*
- [ ] Migrasi mengurai **seluruh daftar bersarang** menjadi baris tabel anak; **tidak ada** daftar
      yang tertinggal di dalam dokumen. *(AC 45 spec)*
- [ ] Daftar yang **kosong** menghasilkan **nol baris anak**, bukan kegagalan. *(AC 27 spec)*
- [ ] Nilai uang pindah **tanpa berubah satu digit pun**; rekonsiliasi membandingkan **secara tepat**,
      bukan dengan toleransi. *(AC 46 spec; **ADR-U-0003**)*
- [ ] Nilai tanggal pindah **tanpa pergeseran zona waktu**.
- [ ] Sequence `M_PRODUCT_LIFE_SEQ` dan `M_PRODUCT_INWARD_LIFE_SEQ` pindah dengan **nilai berjalan
      yang benar**, sehingga identitas baru **tidak bertabrakan** dengan yang lama. *(AC 47 spec)*
- [ ] ⚠️ Kolom baru bernama **`POLICYHOLDER`**, bukan `POLICYHODER`; migrasi **memetakan** ejaan lama
      ke ejaan benar. *(AC 51 spec; `[keputusan work owner]` — penyimpangan sadar 5)*
- [ ] ⚠️ **Tidak ada kolom `IsORS`** di skema baru. *(AC 52 spec)*
- [ ] Migrasi dapat **dijalankan ulang dengan aman** dan punya **jalur mundur yang diuji**.
      *(AC 48 spec)*
- [ ] Skema uji yang dipakai seluruh tiket lain dibangun **dari DDL yang sama** dengan produksi —
      bukan dari tiruan yang ditulis terpisah. *(AC 49 spec)*

#### Blocker

**Tidak ada.** ✅ **OQ-001, OQ-002, dan OQ-JSON-PRODUCT ditutup** `[data DBA]` — DDL, body procedure,
dan isi `JSONDATA` seluruhnya diterima.

#### Catatan

⚠️ **Procedure lama tidak ikut pindah.** `[data DBA]` `POOLDATA.PEGA_M_PRODUCT_LIFE` dan
`POOLDATA.PEGA_M_PRODUCT_INWARD_LIFE` adalah **upsert JSON yang `COMMIT` sendiri**. Keduanya
**tidak** dipanggil dan **tidak** dimigrasikan — alasannya di tiket **03**.

⚠️ **Rekap catatan sumber keliru.** `dba-procedures-and-ddl.md` menyebut "28 field" lalu mendaftar
37; sensus korpus atas `SetProductNameInward` menghasilkan **40**. Aturan pemilihan kolom menetapkan
**Activity sebagai sumber kebenaran**, jadi yang berlaku **40**. `[keputusan work owner]`

✅ **Tabel anak terselesaikan dari korpus (bukan ditunda).** `FinancialUnderwritingList` = tabel
nyata, kolom `[terverifikasi]` dari `CopyFinancialWriting` (`MinInsured`, `MaxInsured`, `Employee`,
`Non_Employee`). `LienClause` = field skalar (bukan tabel). `OutwardList` = dead code jalur OR.
Semua bentuk pasti — tidak ada yang ditebak dari nama.

#### Seam & perintah verifikasi

**Seam: API HTTP** (dipakai ulang dari CL-01) terhadap **skema uji Oracle nyata**. Migrasi diuji
terhadap basis data sungguhan; procedure **tidak** di-mock karena memang tidak dipakai.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Master Product Name Life - 02 - Produk sisi umum — CRUD, identitas dari sequence, gagal terang-terangan

**Status:** ready-for-agent

**Blocked by:** 01 (skema relasional — bentuk barunya harus ada lebih dulu)

#### Hasil & nilai pengguna

Sebagai **admin master produk life**, saya dapat membuat, mengubah, dan membaca **definisi produk**
tanpa menentukan nomor identitasnya sendiri — dan bila penyimpanan gagal, saya **diberi tahu**,
bukan dibiarkan mengira produk sudah tersimpan. *(User story 1–4, 10 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas produk sisi umum |
| `internal/repository` | Tulis/baca `product_life`; identitas dari sequence |
| `internal/services` | Orkestrasi buat/ubah; pemetaan kegagalan menjadi galat domain |
| `internal/handlers` | Endpoint CRUD produk |
| `frontend/` | Layar produk — bagian umum |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `InwardProductName` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `INWARDPRODUCTNAME` / `RULE-HTML-HARNESS` | `Master Product Name Life/Harness/InwardProductName.xml` (540.636 byte) | **titik masuk, satu-satunya Harness** |
| `SetProductName` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SETPRODUCTNAME` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/SetProductName.xml` | **20 field sisi umum** |
| `SaveProductName_Act` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SAVEPRODUCTNAME_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/SaveProductName_Act.xml` (165.682 byte, 16 langkah) | orkestrator simpan; set `CREATEOP`/`UPDATEOP` |
| `NewProductLife` | `ASM-FW-GISFW-…` / `NEWPRODUCTLIFE` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/NewProductLife.xml` | produk baru |

`[terverifikasi]` **Bukan proses berjenjang** — nol rujukan `StatusAkseptasi`, tanpa
`Akseptasi_DT`, tanpa tombol Submit/Decline, **nol rule `Flow`**. Ia editor master murni.

`[data DBA]` **Format identitas ditiru**: **`'1' + lpad(sequence, 5, '0')`** — lima digit, misalnya
`100001`. Sumbernya `M_PRODUCT_LIFE_SEQ`.

`[data DBA]` Kontrak keluaran procedure lama — `StsSave` **100 = sukses / 99 = gagal**, `ErrMsg`
teks yang **terisi juga saat sukses**, `IDPegaOut` identitas final — **tidak dipakai**, karena
procedure-nya dibuang (tiket 03). Dicatat sebagai **rujukan**, bukan perilaku.

#### ADR terkait

**ADR-U-0006** (identitas lewat sequence basis data — aplikasi tidak menyusunnya), **ADR-U-0007**
(jejak audit), **ADR-U-0003** (uang non-float pada atribut produk).

#### Acceptance criteria

- [ ] Produk dapat dibuat dengan nama, tipe, grup, dan atribut sisi umum lainnya. *(AC 1 spec)*
- [ ] **Aplikasi tidak menetapkan identitas produk** — identitas dibuat basis data lewat sequence.
      Test yang menemukan pembentukan identitas di sisi aplikasi **gagal**. *(AC 2 spec; **ADR-U-0006**)*
- [ ] ⚠️ Identitas berbentuk **`'1'` + lima digit** (misalnya `100001`) — format lama **ditiru**.
      *(AC 3 spec; `[keputusan work owner]` — penyimpangan sadar 3)*
- [ ] Menyimpan produk yang **sudah ada** memperbarui baris itu; **tidak** membuat duplikat.
      *(AC 4 spec)*
- [ ] Produk dapat dibaca kembali utuh lewat API. *(AC 5 spec)*
- [ ] ⚠️ **Gagal terang-terangan**: penyimpanan yang gagal **ditampilkan** dan **tidak pernah** tampak
      berhasil. *(AC 10 spec; `[keputusan work owner]` — penyimpangan sadar 3)*
- [ ] Penyimpanan mencatat **siapa** pembuat dan **siapa** pengubah terakhir. *(**ADR-U-0007**)*
- [ ] Nilai uang pada atribut produk diperlakukan sebagai **desimal presisi arbitrer**; **tidak**
      melewati `float` di lapisan mana pun maupun di JSON API. *(AC 28 spec; **ADR-U-0003**)*
- [ ] Nilai uang yang ditulis dan dibaca kembali **identik** — tidak ada pembulatan diam.
      *(AC 29 spec)*
- [ ] ⚠️ Tidak ada **`PoductName`** di kode baru — hanya `ProductName`. *(AC 50 spec;
      `[keputusan work owner]` — penyimpangan sadar 5)*

#### Blocker

**Tidak ada.**

#### Catatan

⚠️ **Salah ketik yang mengenai langkah simpan.** `[terverifikasi]` Ejaan **`PoductName`** (tanpa `r`)
muncul **20 kali di dua berkas saja** — keduanya activity simpan — dan **tidak pernah dideklarasikan
sebagai halaman**, sementara `ProductName` muncul **714 kali**. Yang bersalah ketik justru menjaga
**langkah yang memanggil penyimpanan**.

⚠️ Membuangnya **mengubah perilaku**: guard yang selama ini mungkin tak pernah menyala akan **mulai
menolak** penyimpanan yang dahulu lolos. Itu disadari, bukan diselundupkan. Validasi lengkapnya ada
di tiket **05**.

⚠️ `ErrMsg` procedure lama **terisi juga saat sukses** `[data DBA]` — jangan jadikan keberadaan teks
sebagai penanda kegagalan. Di sistem baru, kegagalan ditentukan hasil transaksi Go sendiri.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — pembuatan identitas lewat sequence tidak
dapat difake dengan jujur.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Master Product Name Life - 03 - Sisi inward (kolom pada satu tabel) + simpan atomik produk + anak

**Status:** ready-for-agent

**Blocked by:** 02 (sisi umum harus ada — keduanya berbagi identitas yang sama)

#### Hasil & nilai pengguna

Sebagai **underwriter**, saya menetapkan **syarat penerimaan** produk — rentang usia, batas uang
pertanggungan, batas retensi ceding, share Nusantara Re, brokerage, masa berlaku; dan sebagai
**organisasi**, saya ingin kedua sisi produk tersimpan **bersama atau tidak sama sekali**, sehingga
tidak pernah ada produk yang tersimpan separuh. *(User story 5–9 di spec)*

⚠️ Inilah tiket dengan penyimpangan paling dalam di konteks ini.

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas produk sisi inward — 40 field |
| `internal/repository` | **Satu transaksi** menulis baris `product_life` beserta tabel anak |
| `internal/services` | Orkestrasi simpan atomik; rollback menyeluruh |
| `internal/handlers` | Endpoint simpan produk utuh |
| `frontend/` | Layar produk — bagian syarat underwriting |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SetProductNameInward` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SETPRODUCTNAMEINWARD` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/SetProductNameInward.xml` (125.938 byte) | **40 field sisi inward** |
| `SaveInwardProductName_Act` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SAVEINWARDPRODUCTNAME_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/SaveInwardProductName_Act.xml` (64.796 byte, 6 langkah) | orkestrator simpan inward |
| `SaveProductNameLIfe` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!SAVEPRODUCTNAMELIFE` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/SaveProductNameLIfe.xml` | ⚠️ **rujukan — tidak dimigrasikan** |
| `SaveProductNameInwardLIfe` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!SAVEPRODUCTNAMEINWARDLIFE` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/SaveProductNameInwardLIfe.xml` | ⚠️ **rujukan — tidak dimigrasikan** |

`[fakta bisnis — work owner]` **SATU produk.** Existing: dua tabel (1:1 by `ID`). Skema baru:
**digabung jadi satu tabel `product_life`** (`[keputusan work owner]`).

⚠️ **Penyimpangan sadar 1 — satu transaksi atomik, procedure lama dibuang.**
`[keputusan work owner]`

`[data DBA]` Kedua procedure penulis (`POOLDATA.PEGA_M_PRODUCT_LIFE`,
`POOLDATA.PEGA_M_PRODUCT_INWARD_LIFE`) adalah **upsert JSON** yang **`COMMIT` sendiri**. Selama
keduanya dipakai, "tersimpan bersama atau tidak sama sekali" **mustahil**.

⚠️ **Penggabungan tabel induk `[keputusan work owner]`.** Sisi umum + sisi inward (dulu 2 tabel
existing, 1:1 by `ID`) **digabung jadi SATU tabel `product_life`** (tiket 01). Maka "sisi inward" =
sekumpulan kolom pada baris produk yang sama, bukan tabel terpisah.

**Keputusan:** sistem baru **tidak memanggil dan tidak memigrasikan** kedua procedure. Go menulis
**satu baris `product_life`** (memuat sisi umum + inward) **beserta seluruh baris tabel anak** dalam
**satu transaksi**; kegagalan di mana pun **membatalkan seluruhnya**.

⚠️ Ini **berbeda** dari pola "panggil procedure apa adanya" yang berlaku di lima konteks Life
sebelumnya (**ADR-U-0006**). Alasannya tunggal: **procedure itu menghalangi atomik yang diminta.**
Efek samping baik penggabungan: atomik induk jadi **trivial** (satu baris), bukan lintas dua tabel.

#### ADR terkait

**ADR-U-0003** (uang non-float — batas pertanggungan, retensi, limit), **ADR-U-0006** (identitas dari
sequence; ⚠️ pengecualian sadar pada pemanggilan procedure), **ADR-U-0007** (jejak audit),
**ADR-U-0015** (kegagalan ditangani eksplisit, tidak ditelan).

#### Acceptance criteria

- [ ] Sisi inward dapat diisi dengan **keempat puluh field**, termasuk rentang usia, batas uang
      pertanggungan, batas retensi ceding, share Nusantara Re, brokerage, dan masa berlaku.
      *(User story 5–6; AC 41 spec)*
- [ ] ⚠️ Baris `product_life` (sisi umum + inward) **beserta seluruh baris tabel anak** ditulis
      dalam **satu transaksi**: kegagalan pada tabel anak mana pun **membatalkan seluruhnya**.
      Dibuktikan dengan menyuntikkan kegagalan pada satu tabel anak lalu memastikan **tidak ada baris**
      `product_life`. *(AC 6, 7 spec; `[keputusan work owner]` — penyimpangan sadar 1)*
- [ ] ⚠️ Sisi umum & sisi inward = **satu baris** `product_life` (dulu dua tabel existing digabung);
      field yang muncul di kedua sisi jadi **satu kolom**. *(AC 8 spec; `[keputusan work owner]`)*
- [ ] ⚠️ **Procedure JSON lama tidak dipanggil.** Test yang menemukan pemanggilan
      `PEGA_M_PRODUCT_LIFE` atau `PEGA_M_PRODUCT_INWARD_LIFE` **gagal**. *(AC 9 spec;
      `[keputusan work owner]`)*
- [ ] Nilai uang (`*LIMIT*`, `*SUMINSURED*`, `*SUMREASURED*`) diperlakukan sebagai **desimal presisi
      arbitrer**; **tidak** melewati `float`. *(AC 28 spec; **ADR-U-0003**)*
- [ ] Persentase (share, brokerage, retensi) juga desimal — **tidak** dibulatkan ke bilangan bulat.
      *(AC 30 spec)*
- [ ] Tanggal (`BEGIN`, `MATURE`, `STNC`, `BIRTHDAY`) tersimpan sebagai **tanggal**, bukan teks.
      *(AC 31 spec)*
- [ ] Batas bawah yang **lebih besar** dari batas atas — usia maupun uang pertanggungan — **ditolak**,
      dengan pesan yang menyebut field-nya. *(AC 32 spec)*
- [ ] Menyimpan produk yang sudah ada memperbarui **kedua** sisi, bukan menambah baris baru di salah
      satunya. *(AC 4 spec)*
- [ ] ⚠️ Kolom pemegang polis bernama **`POLICYHOLDER`**, bukan `POLICYHODER`. *(AC 51 spec)*

#### Blocker

**Tidak ada.**

#### Catatan

⚠️ **Mengapa konteks ini menyimpang dari lima konteks Life sebelumnya.** Di Claim Life, Komite,
PremiumList Life, Endorsement Life, dan Master Contract Retro Life keputusannya selalu **"panggil
procedure apa adanya"**. Di sini tidak — dan D1 (atomik) menarik D2 (relasional) sebagai akibat:
bila Go menulis sendiri, menulis **JSON** tidak lagi masuk akal. **Satu keputusan, dua akibat.**

⚠️ **Rujukan yang tidak dipakai** `[data DBA]`: `StsSave` 100/99, `ErrMsg` (terisi juga saat
sukses), `IDPegaOut`. Ketiganya milik procedure yang dibuang; dicatat sebagai jejak, bukan perilaku.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — atomik lintas satu tabel induk
`product_life` dan lima tabel anak **hanya berperilaku benar pada basis data sungguhan**;
memalsukannya berarti tidak menguji apa pun yang penting.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Master Product Name Life - 04 - Tujuh pemilih master — dan RI Rate yang bukan RI Risk

**Status:** ready-for-agent

**Blocked by:** 02 (pemilih mengisi field pada produk)

#### Hasil & nilai pengguna

Sebagai **admin master**, saya memilih ceding, pemegang polis, mata uang, sumber bisnis, dan penyebab
klaim dari **master masing-masing**; dan sebagai **underwriter** saya memilih **RI Risk** dan
**RI Rate** dari dua master yang **berbeda** — sehingga tidak ada nama yang diketik bebas dan tidak
ada dua master yang tertukar. *(User story 11–17 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Tujuh jenis nilai master beserta identitasnya |
| `internal/repository` | Pembacaan tiap master |
| `internal/services` | Penolakan nilai di luar master |
| `internal/handlers` | Tujuh endpoint daftar master |
| `frontend/` | Tujuh dialog pemilih |

#### Rule Pega sumber

`[terverifikasi]` Sebelas FlowAction modul ini seluruhnya bercorak **pemilih dan konfirmasi**, bukan
langkah proses. Tujuh di antaranya pemilih master:

| Pemilih | Class / Nama / Tipe | Arti `[fakta bisnis — work owner]` |
| --- | --- | --- |
| `ChooseCeding` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSECEDING` / `RULE-OBJ-FLOW-ACTION` | perusahaan ceding |
| `ChoosePolicyHolder` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSEPOLICYHOLDER` / `RULE-OBJ-FLOW-ACTION` | pemegang polis |
| `ChooseCurrency` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSECURRENCY` / `RULE-OBJ-FLOW-ACTION` | mata uang |
| `ChooseSOB` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSESOB` / `RULE-OBJ-FLOW-ACTION` | **SOB = Source of Business** |
| `ChooseCauseOfLoss` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSECAUSEOFLOSS` / `RULE-OBJ-FLOW-ACTION` | penyebab klaim/rugi |
| `ChooseRIRisk` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSERIRISK` / `RULE-OBJ-FLOW-ACTION` | master **RI Risk** — mengisi `RIRISKID`/`RIRISK` |
| **`ChooseRIRate`** | ⚠️ **`ASM-FW-GISFW-DATA-PLAN`** / `CHOOSERIRATE` / `RULE-OBJ-FLOW-ACTION` | master **RI Rate** — **class berbeda** |

⚠️ `[fakta bisnis — work owner]` **RI Rate ≠ RI Risk.** Dua master berbeda; **jangan digabung**.
**RI Risk** melekat pada **produk induk**; **RI Rate** melekat pada **baris plan** (tabel anak,
tiket 06).

`[terverifikasi]` Pendukung: `SearchPolicyHolder_act` (`ASM-FW-GISFW-…` / `SEARCHPOLICYHOLDER_ACT` /
`RULE-OBJ-ACTIVITY`), `SetRIRate` (`…` / `SETRIRATE`), `SetParamRate` (`…` / `SETPARAMRATE`).

#### ADR terkait

**ADR-U-0001** (master pendukung dikelola konteks lain — modul ini **memilih dari** mereka),
**ADR-U-0003** (nilai bermata uang mengikuti mata uang terpilih).

#### Acceptance criteria

- [ ] Ketujuh nilai dipilih dari **masternya masing-masing** — **tidak** diketik bebas.
      *(AC 18 spec)*
- [ ] ⚠️ **RI Rate dan RI Risk adalah dua master berbeda** dan **tidak pernah tertukar**. RI Risk
      mengisi produk induk; RI Rate mengisi **baris plan**. Test yang menemukan satu sumber untuk
      keduanya **gagal**. *(AC 19 spec; `[fakta bisnis — work owner]`)*
- [ ] Nilai yang dipilih tersimpan **beserta identitas masternya**, bukan hanya namanya —
      `CEDING`+`CEDINGID`, `SOBNAME`+`SOBID`, `CAUSE`+`CAUSEID`, `RIRISK`+`RIRISKID`. *(AC 20 spec)*
- [ ] Pilihan yang **tidak ada** di master **ditolak**. *(AC 21 spec)*
- [ ] Nama yang ditampilkan dapat **dibangun ulang** dari identitasnya — nama tersimpan adalah
      salinan tampilan, bukan sumber kebenaran.
- [ ] Pemilih pemegang polis dapat **dicari**, tidak hanya digulir — daftar master dapat besar.
- [ ] Mata uang terpilih berlaku pada **seluruh nilai uang** produk itu. *(**ADR-U-0003**)*
- [ ] Ketujuh daftar dibaca dari basis data, **tidak ditanam** sebagai konstanta di kode.

#### Blocker

**Tidak ada pemblokir.**

⚠️ `[terbuka]` **OQ kecil — tidak memblokir:** kepanjangan **"RI Rate"** dan **"RI Risk"** belum
diketahui. Pemilik: **Product+UW**. Tiket ini dapat selesai tanpa jawaban itu — yang mengikat adalah
**keduanya master terpisah**, bukan namanya.

#### Catatan

⚠️ **Satu pemilih ber-class berbeda.** `[terverifikasi]` Enam pemilih ber-class
`ASM-FW-GISFW-INT-PRODUCT_LIFE`; **`ChooseRIRate` ber-class `ASM-FW-GISFW-DATA-PLAN`**. Ini sejalan
dengan perannya: RI Rate dipakai pada **plan**, bukan pada produk induk.

⚠️ **Master pendukung di luar cakupan.** Modul ini **memilih dari** ketujuh master; ia **tidak
mengelolanya**. Lihat §Out of Scope spec.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Master Product Name Life - 05 - Validasi wajib-isi — lima pemeriksaan, seragam di kedua sisi

**Status:** ready-for-agent

**Blocked by:** 03 (sisi inward — pemegang polis diperiksa dari sana), 04 (sumber bisnis dan ceding
berasal dari pemilih master)

#### Hasil & nilai pengguna

Sebagai **admin master**, saya ingin produk **ditolak** bila field penentu belum terisi — dengan
pesan yang menyebut apa yang kurang — dan saya ingin aturannya **sama** dari jalur mana pun saya
menyimpan, sehingga tidak ada pintu yang lebih longgar. *(User story 1, 10 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/services` | **Satu** aturan wajib-isi, dipanggil kedua jalur simpan |
| `internal/handlers` | Pesan penolakan menyebut field yang kurang |
| `frontend/` | Galat tampil di dekat field yang bersangkutan |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path |
| --- | --- | --- |
| `SaveProductName_Act` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SAVEPRODUCTNAME_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/SaveProductName_Act.xml` |
| `SaveInwardProductName_Act` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SAVEINWARDPRODUCTNAME_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/SaveInwardProductName_Act.xml` |

`[terverifikasi]` **Lima pemeriksaan di `SaveProductName_Act`:**

| Step | Pesan | Precondition korpus |
| ---: | --- | --- |
| 1 | — | `ProductName.TYPE=="" \|\| ProductName.GRUP==""` |
| 2 | "error Product Name" | `ProductName.PRODUCTNAME==""` |
| 3 | "error Ceding" | `ProductName.CEDING==""` |
| 4 | "error Policy Holder" | `ProductNameInward.POLICYHODERNAME==""` |
| 5 | "error SOB" | `ProductName.SOBNAME==""` |

⚠️ `[terverifikasi]` Di **`SaveInwardProductName_Act`**, dua pemeriksaan padanannya **ter-remark**:
step 2 "error Type" (`<pyStepsBlockName>//` baris 414) dan step 3 "error Grup" (baris 601).
**Jalur inward lebih longgar dari jalur utama.**

⚠️ `[data DBA]` Basis data **tidak menegakkan apa pun** — seluruh kolom hasil migrasi **nullable**
(tiket 01). **Wajib-isi seluruhnya ditegakkan di Go.**

#### ADR terkait

**ADR-U-0007** (jejak audit penolakan), **ADR-U-0009** (migrasi penuh).

#### Acceptance criteria

- [ ] Produk **ditolak** bila **tipe** atau **grup** kosong. *(AC 11 spec)*
- [ ] Produk **ditolak** bila **nama produk** kosong. *(AC 12 spec)*
- [ ] Produk **ditolak** bila **ceding** kosong. *(AC 13 spec)*
- [ ] Produk **ditolak** bila **pemegang polis** kosong. *(AC 14 spec)*
- [ ] Produk **ditolak** bila **sumber bisnis** kosong. *(AC 15 spec)*
- [ ] ⚠️ Kelima pemeriksaan berlaku **SAMA pada kedua sisi produk** — **tidak ada** sisi yang lebih
      longgar. Test yang menemukan jalur simpan dengan pemeriksaan lebih sedikit **gagal**.
      *(AC 16 spec)*
- [ ] Wajib-isi ditegakkan **di Go**, bukan oleh basis data — seluruh kolom hasil migrasi nullable.
      Test yang mengandalkan basis data untuk menolak nilai kosong **gagal**. *(AC 17 spec;
      `[data DBA]`)*
- [ ] Pesan penolakan **menyebut field** yang kurang, bukan galat umum.
- [ ] Bila **beberapa** field kurang sekaligus, **semuanya** dilaporkan — bukan hanya yang pertama.
- [ ] Aturan wajib-isi berada di **satu tempat**, dipanggil kedua jalur — bukan disalin dua kali.

#### Blocker

**Tidak ada.**

#### Catatan

⚠️ **Pega lebih longgar di satu sisi, dan itu tidak ditiru.** `[terverifikasi]` Dua pemeriksaan pada
jalur inward ter-remark. Sistem baru **menyeragamkan** — dua jalur atas satu produk tidak boleh
berbeda ketatnya.

⚠️ **OQ-066 diterapkan:** remark itu **tidak** saya simpulkan sendiri sebagai kelalaian; penyeragaman
adalah **keputusan** yang tercatat di spec §8, bukan tafsir dari penanda korpus.

⚠️ **Perubahan perilaku yang disadari.** Bersama pembuangan salah ketik `PoductName` (tiket 02),
penegakan ini akan **mulai menolak** penyimpanan yang dahulu lolos. Itu disengaja.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Master Product Name Life - 06 - Daftar plan dan batas underwriting

**Status:** ready-for-agent

**Blocked by:** 03 (baris anak ikut dalam transaksi atomik produk)

#### Hasil & nilai pengguna

Sebagai **underwriter**, saya mencatat **beberapa plan** dalam satu produk — masing-masing dengan
benefit dan RI Rate-nya — dan **beberapa batas underwriting** dengan status medis serta rentang usia
dan uang pertanggungan, sehingga satu produk dapat menawarkan beberapa paket dengan syarat
bertingkat. *(User story 18–19, 23 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Baris plan; baris batas underwriting |
| `internal/repository` | Tulis/baca kedua tabel anak **di dalam transaksi produk** |
| `internal/services` | Aturan per baris; keterkaitan RI Rate pada plan |
| `internal/handlers` | Endpoint kedua daftar |
| `frontend/` | Dua grid di layar produk |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `ProteksiPlanListLife` | `ASM-FW-GISFW-…` / `PROTEKSIPLANLISTLIFE` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/ProteksiPlanListLife.xml` | daftar plan |
| `CopyUnderWritingLimit` | `ASM-FW-GISFW-…` / `COPYUNDERWRITINGLIMIT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/CopyUnderWritingLimit.xml` | salin batas underwriting |
| `BrowseUnderwritingList` | `ASM-FW-GISFW-…` / `BROWSEUNDERWRITINGLIST` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/BrowseUnderwritingList.xml` | baca daftar |
| `ChooseRIRate` | `ASM-FW-GISFW-DATA-PLAN` / `CHOOSERIRATE` / `RULE-OBJ-FLOW-ACTION` | `Master Product Name Life/FlowAction/ChooseRIRate.xml` | pemilih RI Rate (tiket 04) |
| `CountMaxReasured_Act`, `CountMaxSumReasured_Act` | `ASM-FW-GISFW-…` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/` | hitungan terkait batas |

`[data DBA]` Isi baris, dari contoh data nyata:

| Tabel anak | Kolom baris |
| --- | --- |
| `product_life_plan` | plan, id plan, nama, **benefit**, **RI Rate + id-nya** |
| `product_life_uw_limit` | deskripsi, **status medis** (`FCL` / `NM` / `MEDIS`), usia min/maks, uang pertanggungan min/maks |

⚠️ `[fakta bisnis — work owner]` **RI Rate melekat pada baris plan**, bukan pada produk induk —
berbeda dari **RI Risk** yang melekat pada induk.

#### ADR terkait

**ADR-U-0003** (uang non-float — batas uang pertanggungan pada baris underwriting),
**ADR-U-0007** (jejak audit), **ADR-U-0015** (kegagalan ditangani eksplisit).

#### Acceptance criteria

- [ ] Produk dapat memuat **beberapa plan**, masing-masing dengan benefit dan RI Rate-nya.
      *(AC 22 spec)*
- [ ] Produk dapat memuat **beberapa batas underwriting** dengan deskripsi, status medis, rentang
      usia, dan rentang uang pertanggungan. *(AC 23 spec)*
- [ ] Status medis hanya menerima **`FCL`**, **`NM`**, atau **`MEDIS`**; nilai lain **ditolak**.
      *(`[data DBA]` dari contoh nyata)*
- [ ] Baris dapat **ditambah dan dihapus** tanpa menyentuh produk induk — kecuali ketika penyimpanan
      dilakukan sebagai satu kesatuan. *(AC 26 spec)*
- [ ] Daftar yang **kosong** tersimpan sebagai daftar kosong, **bukan** kegagalan. *(AC 27 spec)*
- [ ] ⚠️ Kegagalan pada baris plan atau batas underwriting **membatalkan seluruh simpan** (baris
      `product_life` + seluruh anak) — baris anak ikut dalam transaksi produk.
      *(AC 7 spec; `[keputusan work owner]`)*
- [ ] Batas uang pertanggungan pada baris underwriting diperlakukan sebagai **desimal presisi
      arbitrer**; **tidak** melewati `float`. *(AC 28 spec; **ADR-U-0003**)*
- [ ] Batas bawah yang **lebih besar** dari batas atas — usia maupun uang pertanggungan — **ditolak**.
      *(AC 32 spec)*
- [ ] ⚠️ **RI Rate pada baris plan** berasal dari master RI Rate, **bukan** dari master RI Risk.
      *(AC 19 spec; `[fakta bisnis — work owner]`)*
- [ ] Urutan baris dalam daftar **dipertahankan** saat dibaca kembali.

#### Blocker

**Tidak ada.**

#### Catatan

⚠️ **Dua tabel anak paling berisi.** `[data DBA]` Dari ketujuh daftar bersarang, hanya `PlanList` dan
`UnderwritingLimitList` yang **terisi** pada contoh produk nyata — karena itu bentuk kolomnya
terbaca. Lima sisanya ditangani tiket **07**.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — keikutsertaan baris anak dalam transaksi
produk hanya berperilaku benar pada basis data sungguhan.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Master Product Name Life - 07 - Komentar, dokumen klaim, dan underwriting finansial

**Status:** ready-for-agent

**Blocked by:** 03 (baris anak ikut dalam transaksi atomik produk)

#### Hasil & nilai pengguna

Sebagai **admin master**, saya mencatat **dokumen klaim** yang disyaratkan produk dan meninggalkan
**komentar** beserta tanggal dan nama penulisnya; dan sebagai **underwriter** saya mencatat
**underwriting finansial** bila produk memerlukannya — sehingga riwayat pertimbangan dan syarat
pelengkap tidak hilang. *(User story 20–23 di spec)*

> ✅ **Blocker "tiga daftar belum berbentuk" DISELESAIKAN dari korpus** (bukan ditunda). Sensus
> Activity menunjukkan hanya **satu** yang benar-benar tabel anak nyata:
> - `[terverifikasi]` **`LIENCLAUSE` = field skalar**, bukan list — ia `ProductNameInward.LIENCLAUSE`
>   (satu dari 40 field inward, tiket 03). **Bukan tabel anak.** ❌ `product_life_lien` dicoret.
> - `[terverifikasi]` **`OutwardList` = dead code** — hanya diisi `GetReinsTypeOR_Life` (jalur OR,
>   mati oleh gerbang `1==2`). **Bukan tabel.** ❌ `product_life_outward` dicoret.
> - `[terverifikasi]` **`FinancialUnderwritingList` = tabel anak nyata**, kolom terbaca dari
>   `CopyFinancialWriting`: `MinInsured`, `MaxInsured`, `Employee`, `Non_Employee`.
>   ✅ `product_life_fin_uw` — kolom pasti, tidak ditebak.

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Lima jenis baris anak |
| `internal/repository` | Tulis/baca kelima tabel anak **di dalam transaksi produk** |
| `internal/services` | Aturan per jenis baris |
| `internal/handlers` | Endpoint kelima daftar |
| `frontend/` | Lima grid di layar produk |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `AddCommentList_Act` | `ASM-FW-GISFW-…` / `ADDCOMMENTLIST_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/AddCommentList_Act.xml` | tambah komentar — dipanggil `SaveProductName_Act` step 7 |
| `ConvertHistoryDate` | `ASM-FW-GISFW-…` / `CONVERTHISTORYDATE` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/ConvertHistoryDate.xml` | format tanggal riwayat |
| `CopyFinancialWriting` | `ASM-FW-GISFW-…` / `COPYFINANCIALWRITING` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/CopyFinancialWriting.xml` | salin underwriting finansial |

`[terverifikasi]` Tabel anak nyata modul ini (kolom dari Activity, bukan tebakan):

| Tabel anak | Kolom baris | Sumber |
| --- | --- | --- |
| `product_life_comment` | `Date`, `OperatorName`, `Suggest` (isi saran) | `AddCommentList_Act` / contoh JSON |
| `product_life_document_claim` | `Document` (+ audit `pxCreate*`) | contoh JSON |
| `product_life_plan` | (di tiket 06) | — |
| `product_life_uw_limit` | (di tiket 06) | — |
| `product_life_fin_uw` | `MinInsured`, `MaxInsured`, `Employee`, `Non_Employee` | `[terverifikasi]` `CopyFinancialWriting` |

**Bukan tabel (dicoret):** `product_life_lien` (LIENCLAUSE = field skalar di `productinward_life`),
`product_life_outward` (OutwardList = dead code jalur OR).

#### ADR terkait

**ADR-U-0007** (jejak audit — komentar adalah riwayat pertimbangan), **ADR-U-0003** (bila daftar yang
belum berbentuk ternyata memuat nilai uang).

#### Acceptance criteria

- [ ] Produk dapat memuat **dokumen klaim** yang disyaratkan. *(AC 24 spec)*
- [ ] Produk dapat memuat **komentar**, tersimpan beserta **tanggal dan nama penulisnya**.
      *(AC 24–25 spec)*
- [ ] Produk dapat memuat **underwriting finansial** (`product_life_fin_uw`) dengan kolom
      `MinInsured`, `MaxInsured`, `Employee`, `Non_Employee`. *(AC 24 spec; `[terverifikasi]`)*
- [ ] `LIENCLAUSE` diperlakukan sebagai **field skalar** di sisi inward (tiket 03), **bukan** tabel.
- [ ] **Tidak ada** tabel `product_life_outward` — `OutwardList` milik jalur OR yang mati.
- [ ] Baris dapat **ditambah dan dihapus** tanpa menyentuh produk induk — kecuali ketika penyimpanan
      dilakukan sebagai satu kesatuan. *(AC 26 spec)*
- [ ] Daftar yang **kosong** tersimpan sebagai daftar kosong, **bukan** kegagalan. *(AC 27 spec)*
- [ ] ⚠️ Kegagalan pada baris mana pun **membatalkan seluruh simpan** (baris `product_life` +
      seluruh anak) — baris anak ikut dalam transaksi produk. *(AC 7 spec; `[keputusan work owner]`)*
- [ ] ⚠️ Komentar **tidak dapat diubah** setelah tersimpan — ia riwayat, bukan field.
      `[keputusan work owner]` (dikonfirmasi: immutable).
- [ ] Urutan baris dalam tiap daftar **dipertahankan** saat dibaca kembali.

#### Blocker

**Tidak ada.** ✅ Bentuk seluruh tabel anak **terselesaikan dari korpus** — tidak ada yang ditunda
ke implementasi.

#### Catatan

⚠️ **Aturan pemilihan kolom tetap berlaku** `[keputusan work owner]`: kolom = field yang **di-SET di
Activity**. Field yang hanya muncul di Section dengan **visibilitas mati** (`1=2` / `1==2` /
`never`) dan tidak di-set Activity **jangan diambil**. Modul ini `[terverifikasi]` memang memuat
banyak gerbang `1==2` di Harness-nya.

⚠️ **Komentar melekat pada jalur simpan.** `[terverifikasi]` `SaveProductName_Act` **step 7**
memanggil `AddCommentList_Act` — jadi komentar terbentuk saat menyimpan produk, bukan lewat jalur
terpisah. Pertahankan keterkaitan itu.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Master Product Name Life - 08 - Lampiran — opsional, terpisah dari metadata, dengan status yang terlihat

**Status:** ready-for-agent

**Blocked by:** 02 (lampiran melekat pada produk yang sudah tersimpan)

#### Hasil & nilai pengguna

Sebagai **admin master**, saya melampirkan berkas pada produk, mengunduhnya kembali, dan
menghapusnya — dan bila unggahan gagal, **produk saya tetap tersimpan** dengan status lampiran yang
jelas, sehingga gangguan penyimpanan berkas tidak menghentikan pekerjaan saya.
*(User story 24–27 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Metadata lampiran + statusnya |
| `internal/repository` | Rekam lampiran; **klien penyimpanan berkas terpisah** |
| `internal/services` | **Pemisahan metadata dan lampiran**; status per lampiran |
| `internal/handlers` | Endpoint unggah, unduh, hapus, daftar lampiran |
| `frontend/` | Panel lampiran dengan status per berkas |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `ProductNameSaveAttachment` | `ASM-FW-GISFW-…` / `PRODUCTNAMESAVEATTACHMENT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/ProductNameSaveAttachment.xml` | simpan lampiran |
| `LoadAttachmentProdName` | `ASM-FW-GISFW-…` / `LOADATTACHMENTPRODNAME` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/LoadAttachmentProdName.xml` | muat daftar |
| `DownloadAttProdName_Act` | `ASM-FW-GISFW-…` / `DOWNLOADATTPRODNAME_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/DownloadAttProdName_Act.xml` | unduh satu |
| `DownloadAll_Act` | `ASM-FW-GISFW-…` / `DOWNLOADALL_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/DownloadAll_Act.xml` | unduh seluruhnya |
| `DeleteAttacProdName_act` | `ASM-FW-GISFW-…` / `DELETEATTACPRODNAME_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/DeleteAttacProdName_act.xml` | hapus |
| `InsertAttachProdName_Sql` | `ASM-FW-GISFW-…` / `ASM!INSERTATTACHPRODNAME_SQL` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/InsertAttachProdName_Sql.xml` | rekam lampiran |
| `GetAttachmentProdName_Sql` | `ASM-FW-GISFW-…` / `ASM!GETATTACHMENTPRODNAME_SQL` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/GetAttachmentProdName_Sql.xml` | baca lampiran |
| `DeleteAttachProdName_Sql` | `ASM-FW-GISFW-…` / `ASM!DELETEATTACHPRODNAME_SQL` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/DeleteAttachProdName_Sql.xml` | hapus rekam |
| `GetMimeType` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `GETMIMETYPE` / `RULE-OBJ-DECISIONTABLE` | `Master Product Name Life/DecisionTable/GetMimeType.xml` | jenis berkas |
| `SetCategory_act` | `ASM-FW-GISFW-…` / `SETCATEGORY_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/SetCategory_act.xml` | kategori lampiran |

⚠️ **Penyimpangan sadar 4 — lampiran opsional dan terpisah.** `[keputusan work owner]` Lampiran
adalah **pelengkap** master produk, **bukan syarat sahnya**. Memblokir penyimpanan produk karena
penyimpanan berkas sedang tak dapat dihubungi akan menghentikan pekerjaan tanpa perlu.

#### ADR terkait

**ADR-U-0010** (penyimpanan berkas tetap di Google Storage), **ADR-U-0015** (efek keluar — kegagalan
tidak boleh diam-diam), **ADR-U-0007** (jejak audit).

#### Acceptance criteria

- [ ] ⚠️ **Lampiran OPSIONAL**: produk tersimpan **tanpa** lampiran tanpa keluhan. *(AC 33 spec;
      `[keputusan work owner]` — penyimpangan sadar 4)*
- [ ] ⚠️ **Lampiran TERPISAH dari metadata**: kegagalan unggah **tidak** membatalkan produk yang
      sudah tersimpan. Produk tersimpan **lebih dulu**; lampiran menyusul. *(AC 34 spec)*
- [ ] ⚠️ Kegagalan unggah **tercatat dan terlihat** — bukan senyap — dan **dapat diulang**.
      *(AC 35 spec; **ADR-U-0015**)*
- [ ] **Status tiap lampiran** terbaca lewat API: **terunggah**, **gagal**, atau **belum**.
      *(AC 36 spec)*
- [ ] Lampiran dapat **diunduh** satu per satu **dan** seluruhnya sekaligus. *(AC 38 spec)*
- [ ] Lampiran dapat **dihapus**; penghapusan mencakup rekam **dan** berkasnya.
- [ ] Jenis berkas ditentukan dari isinya, dan berkas yang jenisnya tidak didukung **ditolak** dengan
      pesan yang menyebut jenisnya.
- [ ] Menghapus **produk** tidak meninggalkan lampiran yatim.
- [ ] Mengunggah berkas bernama sama **tidak** menimpa diam-diam — pengguna diberi tahu.

#### Blocker

**Tidak ada.** Resolusi alamat penyimpanan, token, dan pengulangan kegagalan ditangani tiket **09**.

#### Catatan

⚠️ **Urutannya mengikat: produk dulu, lampiran menyusul.** `[keputusan work owner]` Bukan sekadar
preferensi — ia yang membuat AC 34 dapat diuji: menyuntikkan kegagalan unggah lalu memastikan produk
**tetap ada**.

⚠️ **Berbeda dari Master Contract Retro Life.** Konteks itu **tanpa integrasi luar sama sekali**;
konteks ini punya efek keluar nyata, sehingga **ADR-U-0010** dan **ADR-U-0015** berlaku di sini.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**; **penyimpanan berkas di-fake** di balik
interface — yang diperiksa adalah **efeknya**: unggah terjadi atau tidak, status lampiran, kegagalan
tercatat dan dapat diulang.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Master Product Name Life - 09 - Efek keluar penyimpanan berkas — resolusi alamat, token, dan pengulangan

**Status:** ready-for-agent

**Blocked by:** 08 (lampiran harus ada lebih dulu — tiket ini mengeraskan jalur keluarnya)

#### Hasil & nilai pengguna

Sebagai **tim operasi**, saya ingin alamat penyimpanan berkas **dibaca dari konfigurasi saat
dijalankan** — bukan ditanam di kode — dan saya ingin unggahan yang gagal **tercatat dan dapat
diulang**, sehingga perpindahan lingkungan tidak menuntut penempelan ulang dan tidak ada berkas yang
hilang diam-diam. *(User story 28 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Resolusi alamat penyimpanan; pengambilan dan **cache token** |
| `internal/clients` | Klien penyimpanan berkas **di balik interface** |
| `internal/services` | Orkestrasi efek keluar; penandaan kegagalan; pengulangan |
| `internal/handlers` | Status efek keluar terbaca API |
| `frontend/` | Penanda "terkirim" / "tertunda" pada lampiran |

#### Rule Pega sumber

`[terverifikasi]` Seluruh rantai berkas ber-class **`ASM-FW-GISFW-INT-T_STORAGE_IMAGE`**:

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `ServiceGoogle` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `SERVICEGOOGLE` / `RULE-CONNECT-REST` | `Master Product Name Life/ConnectREST/ServiceGoogle.xml` | **satu-satunya ConnectREST** modul ini |
| `GetTokenStorage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RNM!GETTOKENSTORAGE_SQL` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/GetTokenStorage_SQL.xml` | ambil token |
| `GetLinkStorage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RNM!GETLINKSTORAGE_SQL` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/GetLinkStorage_SQL.xml` | ambil tautan |
| `Insert_T_Storage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RNM!INSERT_T_STORAGE_SQL` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/Insert_T_Storage_SQL.xml` | rekam berkas |
| `Update_T_Storage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RNM!UPDATE_T_STORAGE_SQL` / `RULE-CONNECT-SQL` | `Master Product Name Life/RDBList/Update_T_Storage_SQL.xml` | perbarui rekam |
| `GetLinkService` | `ASM-FW-GISFW-…` / `GETLINKSERVICE` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/GetLinkService.xml` | **resolusi alamat** |
| `InsertGoogleStorage_Act` | `ASM-FW-GISFW-…` / `INSERTGOOGLESTORAGE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/InsertGoogleStorage_Act.xml` | unggah |
| `GetUrlGoogleStorage_Act` | `ASM-FW-GISFW-…` / `GETURLGOOGLESTORAGE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/GetUrlGoogleStorage_Act.xml` | ambil URL |
| `DeleteGoogleStorage_Act` | `ASM-FW-GISFW-…` / `DELETEGOOGLESTORAGE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Product Name Life/Activity/DeleteGoogleStorage_Act.xml` | hapus berkas |

`[terverifikasi]` `ServiceGoogle` **tanpa URL literal** — alamat datang dari `M_LINK_SERVICE`.

`[data DBA]` Token mengikuti pola `GET_TOKEN_STORAGE`: **MD5**, **di-cache**, **berlaku 1 menit**.

#### ADR terkait

**ADR-U-0013** (alamat di-resolve **runtime**; **dilarang** sebagai literal, konstanta, **maupun env
var** — menggantikan ADR-U-0004), **ADR-U-0015** (efek keluar — kegagalan tidak boleh diam-diam),
**ADR-U-0010** (penyimpanan berkas tetap Google Storage), **ADR-U-0005** (flag lingkungan bila
diperlukan).

#### Acceptance criteria

- [ ] Alamat penyimpanan di-resolve **runtime** dari konfigurasi. Test yang memindai kode untuk URL
      sebagai **literal, konstanta, atau pembacaan env var** **gagal** bila menemukannya.
      *(AC 37 spec; **ADR-U-0013**)*
- [ ] Token **di-cache** dan **diperbarui sebelum kedaluwarsa**; kedaluwarsanya **tidak** menggagalkan
      permintaan pengguna secara langsung. *(AC 39 spec; `[data DBA]` berlaku 1 menit)*
- [ ] Token yang **gagal diambil** menghasilkan kegagalan yang **terlihat**, bukan unggahan yang diam
      saja tidak terjadi.
- [ ] ⚠️ Kegagalan unggah **tercatat dan dapat diulang**; pengulangan **tidak** menggandakan berkas.
      *(AC 35 spec; **ADR-U-0015**)*
- [ ] Klien penyimpanan berkas berada **di balik interface** dan **di-fake** di test; yang diperiksa
      adalah **efeknya**.
- [ ] ⚠️ Kegagalan efek keluar **tidak** membatalkan produk maupun rekam lampirannya — hanya status
      lampiran yang berubah. *(AC 34 spec)*
- [ ] Rekam berkas di basis data dan berkas di penyimpanan **tetap sejalan**: rekam tanpa berkas
      terdeteksi dan dapat diperbaiki.
- [ ] Penghapusan berkas yang **sudah tidak ada** di penyimpanan **tidak** menggagalkan penghapusan
      rekamnya.

#### Blocker

**Tidak ada pemblokir.**

⚠️ `[terbuka]` **OQ-047 — tidak memblokir:** alamat fisik Google Storage tersimpan di
`M_LINK_SERVICE`, dan isinya belum dibaca. **ADR-U-0013** sudah mengatur **cara** meresolusinya, jadi
tiket ini dapat selesai tanpa mengetahui alamatnya — yang mengikat adalah **alamat tidak ditanam**.

#### Catatan

⚠️ **Berbeda dari Master Contract Retro Life.** Konteks itu **tanpa `ConnectREST` sama sekali**,
sehingga ADR-U-0013 dan ADR-U-0015 tidak berlaku di sana. Di sini **berlaku keduanya** — inilah
satu-satunya efek keluar konteks ini.

⚠️ **Jangan samakan dengan ADR-U-0015 versi Komite.** Di Komite Claim Life, efek keluar **wajib
berhasil** (transactional outbox). Di sini `[keputusan work owner]` lampiran **opsional** — yang
wajib adalah **kegagalannya terlihat dan dapat diulang**, bukan berhasil. Taruhannya berbeda.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**; **penyimpanan berkas di-fake** di balik
interface.

```
go test ./internal/...
cd frontend && npm test
make check
```

# Master Contract Retro Life

Jumlah tiket: **13**

## Master Contract Retro Life - 01 - Tahun treaty — CRUD dan aturan abadi

**Status:** ready-for-agent

**Blocked by:** **CL-01** (kerangka aplikasi + seam API — scaffolding lintas konteks, tidak dibuat
di sini)

#### Hasil & nilai pengguna

Sebagai **admin master retro life**, saya dapat membuat dan mengubah **tahun treaty** — akar seluruh
hierarki retrosesi life — dan saya tidak dapat menghapusnya, sehingga riwayat kontrak tahun-tahun
lampau tidak pernah hilang. *(User story 1–4 di spec)*

⚠️ Tiket ini juga **menetapkan pola jalur simpan** yang diwarisi keempat entitas berikutnya: upsert
dikunci identitas, identitas dibuat basis data, dan hasil procedure diperiksa.

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas tahun treaty |
| `internal/repository` | Pemanggilan procedure penulis; pemeriksaan `o_message` |
| `internal/services` | Wajib-isi; larangan hapus |
| `internal/handlers` | Endpoint buat/ubah/daftar tahun treaty |
| `frontend/` | Grid tahun treaty |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveMasterTreatyYear_Life_SQL` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` / `ASM!SAVEMASTERTREATYYEAR_LIFE_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/SaveMasterTreatyYear_Life_SQL.xml` | `POOLDATA.INSERTTREATYYEAR_LIFE` |
| `SaveTreatyYearLife_Act` | `ASM-FW-GISFW-…` / `SAVETREATYYEARLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SaveTreatyYearLife_Act.xml` | orkestrator simpan |
| `NewInputTreatyYear_Life_Act` | `ASM-FW-GISFW-…` / `NEWINPUTTREATYYEAR_LIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/NewInputTreatyYear_Life_Act.xml` | baris baru |
| `SetTreatyYearLife_Act` | `ASM-FW-GISFW-…` / `SETTREATYYEARLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SetTreatyYearLife_Act.xml` | isi form |
| `BrowseTreatyYear_Life_RD` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` / `BROWSETREATYYEAR_LIFE_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Master Contract Retro Life/ReportDefinition/BrowseTreatyYear_Life_RD.xml` | daftar |

`[terverifikasi]` Wajib-isi saat simpan (`SaveTreatyYearLife_Act`): `UNDERWRITINGYEAR`,
`TREATYYEAR`, `STARTDATE`, `ENDDATE`.

`[terverifikasi]` **Tidak ada penghapus** untuk `TREATYYEAR_LIFE` — lima penulis, hanya **empat**
penghapus di seluruh modul.

`[data DBA]` Perilaku `INSERTTREATYYEAR_LIFE`: **upsert dikunci `ID`**
(`SELECT COUNT(1) … WHERE ID = p_ID` → ada `UPDATE`, tidak ada `INSERT`); identitas baris baru
`'1' || lpad(TREATYYEAR_LIFE_SEQ.nextval, 6, '0')`; `p_ID` **diabaikan** saat INSERT;
`TGLUPDATE` = **`SYSDATE`** (parameter `p_TGLUPDATE` diabaikan); `USERID` dari aplikasi;
**`COMMIT` di dalam procedure**; **nol validasi bisnis**.

`[data DBA]` Kolom (`ddl-tables-from-dba.md`): `ID VARCHAR2(100)` (**PK**), `TREATYYEAR VARCHAR2(100)`,
`UNDERWRITINGYEAR VARCHAR2(100)`, `USERID VARCHAR2(100)`, `TGLUPDATE DATE`, `STARTDATE DATE`,
`ENDDATE DATE`. **Semua nullable** — wajib-isi ditegakkan **di Go**.
Sequence `TREATYYEAR_LIFE_SEQ` `START WITH 44` → identitas berikutnya `1000044`.

#### ADR terkait

**ADR-U-0006** (identitas dibuat lewat procedure/basis data — aplikasi tidak menyusunnya),
**ADR-U-0007** (jejak audit), **ADR-U-0009** (migrasi penuh).

#### Acceptance criteria

- [ ] Tahun treaty dapat dibuat dengan tahun underwriting, tahun treaty, tanggal mulai, dan tanggal
      akhir; keempatnya **wajib** — ditegakkan **di Go**, karena basis data **nol `NOT NULL`**.
      *(AC 1 spec)*
- [ ] ⚠️ **Tahun treaty TIDAK DAPAT DIHAPUS** lewat jalur mana pun. Test yang menemukan endpoint
      hapus tahun **gagal**. *(AC 2 spec; `[fakta bisnis — work owner]` — tahun treaty **abadi**,
      disengaja)*
- [ ] Tahun treaty dapat diubah, dan perubahannya mencatat **siapa** pelakunya. *(AC 3 spec)*
- [ ] ⚠️ **Aplikasi tidak menetapkan identitas baris baru** — ia mengirim identitas kosong dan basis
      data yang membuatnya lewat sequence. Test yang menemukan pembentukan identitas di sisi aplikasi
      **gagal**. *(AC 4 spec; **ADR-U-0006**)*
- [ ] Menyimpan tahun yang **sudah ada** memperbarui baris itu — **upsert dikunci identitas**, bukan
      baris kedua. Dibuktikan dengan menyimpan dua kali lalu **menghitung baris**. *(AC 45 spec)*
- [ ] Setiap penyimpanan mengirim **identitas pengguna**; **cap waktu tidak dikirim aplikasi** —
      basis data yang menetapkannya. *(AC 48 spec; `[data DBA]`)*
- [ ] `o_message` **diperiksa** setelah pemanggilan procedure: **kosong/NULL = sukses**, berisi teks
      = gagal. Kegagalan **ditampilkan**, tidak ditelan. *(AC 46 spec — pembersihan HTML dan
      konformansi lintas jalur ada di **tiket 10**)*
- [ ] Daftar tahun treaty dapat dibaca lewat API dan tampil di layar.

#### Blocker

**Tidak ada.** Menunggu scaffolding **CL-01** yang sudah `ready-for-agent` di konteks Claim Life.

#### Catatan

⚠️ **Pola yang ditetapkan di sini diwarisi tiket 02, 05, 06, dan 07.** Keempat jalur simpan lain
berperilaku identik: upsert dikunci `ID`, identitas dari sequence, `TGLUPDATE` = `SYSDATE`, `COMMIT`
internal, nol validasi di basis data.

⚠️ **Setiap simpan adalah transaksi mandiri** — `COMMIT` berada di dalam procedure. Tidak ada
transaksi lintas-baris di sisi basis data, dan itu **tidak boleh disembunyikan** dari pengguna.
*(AC 50 spec)*

#### Seam & perintah verifikasi

**Seam: API HTTP** (dipakai ulang dari CL-01) terhadap **skema uji Oracle nyata** — procedure
**tidak di-mock**: perilaku upsert dan pembuatan identitas lewat sequence tidak dapat difake dengan
jujur.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Master Contract Retro Life - 02 - Kontrak treaty — batas proteksi per mata uang dan lebar layer yang dihitung

**Status:** ready-for-agent

**Blocked by:** 01 (kontrak lahir di bawah tahun treaty; pola jalur simpan ditetapkan di sana)

#### Hasil & nilai pengguna

Sebagai **underwriter**, saya menetapkan **batas bawah dan batas atas** proteksi per mata uang pada
sebuah kontrak treaty, dan **lebar layer dihitung sistem** dari kedua batas itu — sehingga angkanya
tidak pernah menyimpang dari batasnya sendiri. *(User story 5–11 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas kontrak; nilai batas sebagai desimal presisi arbitrer |
| `internal/repository` | Pemanggilan procedure penulis kontrak; pemeriksaan `o_message` |
| `internal/services` | Hitung lebar layer; wajib-isi; pemeriksaan batas bawah ≤ batas atas |
| `internal/handlers` | Endpoint CRUD kontrak |
| `frontend/` | Grid kontrak + batas proteksi; lebar layer **hanya tampil**, tidak dapat diketik |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveMasterTreatyContract_Life_SQL` | `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` / `ASM!SAVEMASTERTREATYCONTRACT_LIFE_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/SaveMasterTreatyContract_Life_SQL.xml` | `POOLDATA.INSERTTREATYCONTRACT_LIFE` |
| `SaveTreatyLimit_Act` | `ASM-FW-GISFW-…` / `SAVETREATYLIMIT_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SaveTreatyLimit_Act.xml` | orkestrator simpan + wajib-isi |
| `NewInputTreatyLimit_Life` | `ASM-FW-GISFW-…` / `NEWINPUTTREATYLIMIT_LIFE` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/NewInputTreatyLimit_Life.xml` | baris baru |
| `SetValueRetroLimit_TreatyYearLife` | `@BASECLASS` / `SETVALUERETROLIMIT_TREATYYEARLIFE` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SetValueRetroLimit_TreatyYearLife.xml` | ⚠️ **penyalur konteks**, bukan rumus |
| `BrowseTreatyContract_Life_RD` | `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` / `BROWSETREATYCONTRACT_LIFE_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Master Contract Retro Life/ReportDefinition/BrowseTreatyContract_Life_RD.xml` | daftar |
| `InboxRetroLimitReinsurers` | `DATA-PORTAL` / `INBOXRETROLIMITREINSURERS` / `RULE-HTML-HARNESS` | `Master Contract Retro Life/Harness/InboxRetroLimitReinsurers.xml` (575.179 byte) | layar grid |

`[fakta bisnis — work owner]` **Arti keenam kolom limit:**

| Kolom | Arti | Tipe `[data DBA]` |
| --- | --- | --- |
| `B_IDR` / `B_USD` | **batas bawah** (prefix `B_` = bawah) | `NUMBER` |
| `IDR` / `USD` | **batas atas** | `NUMBER` |
| `IDR_SELISIH` / `USD_SELISIH` | **atas − bawah** = **lebar layer** | `NUMBER` |

`[fakta bisnis — work owner]` **Layer tersusun menaik** — batas atas satu layer menjadi batas bawah
layer berikutnya: `QS → 2nd QS → Surplus → 2nd Surplus`.

`[terverifikasi]` Wajib-isi (`SaveTreatyLimit_Act`): `REINSTYPEID`, `TREATYSTARTDATE`,
`TREATYENDDATE`, `B_IDR`, `IDR`, `B_USD`. ⚠️ **`USD` TIDAK wajib** — `[data DBA]` kolomnya nullable,
jadi ini **aturan sah**, bukan kelalaian.

⚠️ `[data DBA]` `INSERTTREATYCONTRACT_LIFE` **menulis `IDR_SELISIH`/`USD_SELISIH` apa adanya dari
parameter** — basis data **tidak menghitungnya**. Tidak ada yang menjaganya sinkron dengan kedua
batas.

#### ADR terkait

**ADR-U-0003** (uang non-float — `[data DBA]` kolom bertipe `NUMBER` tanpa presisi, jadi desimal
presisi arbitrer aman), **ADR-U-0006** (identitas dari basis data), **ADR-U-0007** (jejak audit).

#### Acceptance criteria

- [ ] Kontrak lahir **di bawah** satu tahun treaty; kontrak tanpa tahun induk **ditolak**.
      *(AC 5 spec)*
- [ ] Kontrak wajib memuat `REINSTYPEID`, tanggal mulai, tanggal akhir, **batas bawah IDR**, **batas
      atas IDR**, dan **batas bawah USD** — ditegakkan **di Go**. *(AC 6 spec)*
- [ ] **`USD` boleh kosong** — kontrak ber-IDR saja tersimpan tanpa keluhan. *(AC 7 spec;
      `[data DBA]` aturan sah)*
- [ ] ⚠️ **Lebar layer DIHITUNG** dari `batas atas − batas bawah`, per mata uang. **Tidak ada** jalur
      yang membiarkan pengguna mengetiknya, dan **tidak ada** jalur yang mengirimkannya mentah dari
      masukan. *(AC 8 spec; `[keputusan work owner]` — penyimpangan sadar 2)*
- [ ] Lebar layer yang tersimpan **selalu** sama dengan selisih kedua batas, **termasuk setelah salah
      satu batas diubah**. *(AC 9 spec)*
- [ ] Batas bawah yang **lebih besar** dari batas atas **ditolak**, dengan pesan yang menyebut mata
      uangnya. *(AC 10 spec)*
- [ ] Susunan layer menaik dapat direkam: batas atas satu layer boleh menjadi batas bawah layer
      berikutnya **tanpa** dianggap tumpang tindih. *(AC 11 spec)*
- [ ] Seluruh nilai batas diperlakukan sebagai **desimal presisi arbitrer**; **tidak ada** yang
      melewati `float` di lapisan mana pun maupun di JSON API. *(AC 12 spec; **ADR-U-0003**)*
- [ ] Nilai batas yang ditulis dan dibaca kembali **identik** — tidak ada pembulatan diam, termasuk
      pada nilai berpecahan panjang. *(AC 13 spec)*
- [ ] Menyimpan kontrak yang sudah ada = **upsert**, bukan baris kedua; identitas baru dibuat basis
      data. *(AC 45, 4 spec)*
- [ ] `o_message` diperiksa; kegagalan ditampilkan. *(AC 46 spec — HTML dan konformansi di
      **tiket 10**)*

#### Blocker

**Tidak ada.**

#### Catatan

⚠️ **`REINSTYPEID` diterima apa adanya di tiket ini.** Enumerasi yang sah (`Flag = 1`, lima nilai),
penolakan nilai di luar itu, dan **normalisasi** (`REINSTYPEID` hanya di kontrak) ditambahkan
**tiket 04**. Sampai tiket 04 selesai, kontrak menyimpan nilai apa pun yang dikirim.

⚠️ **Nama menyesatkan.** `SetValueRetroLimit_TreatyYearLife` terbaca seperti rumus limit; `[terverifikasi]`
ia **hanya menyalin konteks tahun treaty** ke tiga halaman input dan mengatur tampilan
(`OutputParam.DATASHOW` / `STSSAVE`). **Tidak ada perhitungan.** Jangan tiru namanya.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — procedure tidak di-mock.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Master Contract Retro Life - 03 - Gerbang konsistensi tahun — dibandingkan dari nilai tanggal, bukan potongan teks

**Status:** ready-for-agent

**Blocked by:** 02 (gerbang ini berlaku saat menyimpan kontrak)

#### Hasil & nilai pengguna

Sebagai **admin master**, saya ingin sistem menolak kontrak yang **tanggal mulainya jatuh di tahun
berbeda** dari tahun treaty induknya, sehingga kontrak tidak pernah nyasar tahun dan laporan per
tahun treaty tidak bocor. *(User story 9 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/services` | Aturan gerbang tahun — perbandingan tahun dari nilai tanggal |
| `internal/handlers` | Pesan penolakan yang menyebut kedua tahun |
| `frontend/` | Galat tampil di form kontrak, di dekat field tanggal mulai |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path |
| --- | --- | --- |
| `SaveBusinessLife_Act` | `ASM-FW-GISFW-…` / `SAVEBUSINESSLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SaveBusinessLife_Act.xml` |
| `SaveSecurityLife_Act` | `ASM-FW-GISFW-…` / `SAVESECURITYLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SaveSecurityLife_Act.xml` |
| `SaveSecurityReinsurerLife_Act` | `ASM-FW-GISFW-…` / `SAVESECURITYREINSURERLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SaveSecurityReinsurerLife_Act.xml` |

`[terverifikasi]` **Ketiganya berbagi gerbang yang sama:**

```
@substring(InputRetrocessionLife.TREATYSTARTDATE,6,10) <> InputRetrocessionLifeTreatyType.TREATYYEAR_LIFE
```

yakni **memotong teks tanggal pada posisi 6–10** lalu membandingkannya dengan tahun treaty.

⚠️ `[data DBA]` **DDL membuktikan cara itu keliru sejak awal**: `TREATYSTARTDATE`, `TREATYENDDATE`,
`STARTDATE`, `ENDDATE`, dan `TGLUPDATE` seluruhnya bertipe **`DATE`** — bukan teks. Pemotongan posisi
tetap atas nilai `DATE` bergantung pada format tampilan, yang dapat berubah.

#### ADR terkait

**ADR-U-0007** (jejak audit — penolakan tercatat), **ADR-U-0009** (migrasi penuh).

#### Acceptance criteria

- [ ] ⚠️ **Gerbang tahun**: tahun pada **tanggal mulai kontrak** harus **sama** dengan tahun treaty
      induknya; yang berbeda **ditolak**. *(AC 14 spec; `[keputusan work owner]` — penyimpangan
      sadar 4)*
- [ ] ⚠️ Perbandingan memakai **tahun dari nilai bertipe tanggal**. Test yang menemukan **pemotongan
      teks posisi tetap** di jalur ini **gagal**. *(AC 14 spec)*
- [ ] Pesan penolakan menyebut **kedua tahun** — tahun pada tanggal mulai dan tahun treaty induk —
      sehingga pengguna tahu mana yang harus diperbaiki.
- [ ] Gerbang berlaku pada **penyimpanan kontrak**, dan tidak dapat dilewati lewat jalur API mana
      pun.
- [ ] Kontrak yang **tahunnya cocok** tersimpan tanpa keluhan — gerbang tidak menghasilkan positif
      palsu pada tanggal awal maupun akhir tahun (1 Januari dan 31 Desember diuji).
- [ ] Perubahan tanggal mulai pada kontrak yang sudah tersimpan **diuji ulang** oleh gerbang yang
      sama.

#### Blocker

**Tidak ada.**

#### Catatan

⚠️ **Aturannya dipertahankan, caranya diperbaiki.** `[keputusan work owner]` Aturan "tahun tanggal
mulai = tahun treaty induk" **sah dan berguna** — ia mencegah kontrak nyasar tahun. Yang dibuang
hanyalah **cara Pega membandingkannya**.

`[terverifikasi]` Gerbang yang sama dipasang pada tiga activity simpan berbeda di Pega
(business, security, security reinsurer). Di sistem baru ia **satu aturan di satu tempat**, dipanggil
dari jalur simpan kontrak — bukan disalin tiga kali.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**, dengan **jam yang dapat dikendalikan** agar
kasus pergantian tahun dapat diuji.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Master Contract Retro Life - 04 - Jenis reasuransi — enumerasi `Flag = 1` dan normalisasi `REINSTYPEID`

**Status:** ready-for-agent

**Blocked by:** 02 (jenis reasuransi melekat pada kontrak)

#### Hasil & nilai pengguna

Sebagai **admin master**, saya memilih jenis reasuransi dari **daftar yang sah untuk lini life** dan
melihat **namanya**, bukan kodenya — sehingga tidak ada jenis dari lini lain yang menyelinap masuk
dan layar terbaca manusia. *(User story 30–31 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Jenis reasuransi; `REINSTYPENAME` ditandai **turunan** |
| `internal/repository` | Pembacaan enumerasi dengan penyaring `Flag = 1` |
| `internal/services` | Penolakan nilai di luar enumerasi; **normalisasi** — anak mewarisi dari induk |
| `internal/handlers` | Endpoint daftar jenis reasuransi |
| `frontend/` | Dropdown jenis reasuransi menampilkan nama |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `BrowseReinsuranceTypeLimit_RD` | `ASM-FW-GISFW-INT-REINSURANCETYPE` / `BROWSEREINSURANCETYPELIMIT_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Master Contract Retro Life/ReportDefinition/BrowseReinsuranceTypeLimit_RD.xml` | **sumber enumerasi**, penyaring `.Flag = 1` |
| `TreatyLimit_TypeProtect` | `@BASECLASS` / `TREATYLIMIT_TYPEPROTECT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/TreatyLimit_TypeProtect.xml` | ⚠️ **lookup**, bukan rumus |

`[terverifikasi]` `TreatyLimit_TypeProtect` mencocokkan `REINSTYPEID` dengan `.ID`, lalu menyalin
**`.Note`** ke `REINSTYPENAME`. Kolom yang tersedia di class: `.ID`, `.Code`, `.Note`, `.Type`,
`.Flag`, `.SOANote`, `.UserID`. Penyaring RD: **`.Flag = 1`** (baris 581/585 dan salinan indeks
887/893).

`[fakta bisnis — work owner]` **Lima jenis untuk lini life:**

| ID | Nama (`.Note`) |
| --- | --- |
| 10196 | QS |
| 10197 | 2ND QS |
| 10198 | SURPLUS |
| 10199 | 2ND SURPLUS |
| 10200 | OR |

`[fakta bisnis — work owner]` **`Flag = 1` = "for life"** — penanda lini, **bukan** aktif/nonaktif.
Sejalan dengan memo korpus `[terverifikasi]` `pyUsage: "Parameter Flag, 1 for life"`.
Kolom `.Code` dan `.Type` **tidak dipakai** sistem.

⚠️ `"OR"` dicatat **apa adanya**; artinya **tidak ditafsirkan**.

`[terverifikasi]` **Denormalisasi di Pega**: `REINSTYPEID`/`REINSTYPENAME` tersimpan **tiga kali**
(kontrak, reinsurer, business) dan `TREATYYEAR` **dua kali** (tahun, business), **tanpa penjaga
konsistensi** — dan `[data DBA]` procedure **tetap menulisnya** ke tabel anak.

#### ADR terkait

**ADR-U-0001** (batas konteks — enumerasi ini dipakai lintas lini), **ADR-U-0009**.

#### Acceptance criteria

- [ ] Daftar jenis reasuransi yang ditawarkan **hanya** yang **`Flag = 1`** — lima jenis, ID
      **10196**–**10200**. *(AC 38 spec)*
- [ ] Layar menampilkan **nama** dari kolom `.Note`, **bukan** kode. *(AC 39 spec)*
- [ ] `REINSTYPEID` di luar kelima nilai itu **ditolak** saat menyimpan kontrak. *(AC 40 spec)*
- [ ] ⚠️ **Normalisasi**: `REINSTYPEID` tersimpan **hanya pada kontrak**; reinsurer dan business
      mewarisinya lewat `TREATYCONTRACTID`. Test yang menemukan `REINSTYPEID` yang **dapat ditulis**
      pada tabel anak **gagal**. *(AC 41 spec; `[keputusan work owner]` — penyimpangan sadar 1)*
- [ ] ⚠️ `TREATYYEAR` tersimpan **hanya pada tahun treaty**; business mewarisinya.
      *(AC 42 spec; `[keputusan work owner]`)*
- [ ] `REINSTYPENAME` yang disimpan ditandai **turunan/cache** — **bukan sumber kebenaran** — dan
      **selalu dapat dibangun ulang** dari `REINSTYPEID`. *(AC 43 spec)*
- [ ] Mengubah jenis reasuransi pada kontrak **tidak** menuntut pembaruan manual di tabel anak —
      anak membacanya lewat induk.
- [ ] Enumerasi dibaca dari basis data, **tidak ditanam** sebagai konstanta di kode. Test yang
      menemukan kelima ID sebagai literal di lapisan services **gagal**.

#### Blocker

**Tidak ada pemblokir.**

⚠️ `[terbuka]` **OQ kecil** — **tidak memblokir**: nama **tabel fisik Oracle** untuk class
`ASM-FW-GISFW-INT-REINSURANCETYPE` belum diketahui (korpus hanya menyebut *class*), dan isi kolom
`.Code` serta `.Type` belum terbaca. `[fakta bisnis — work owner]` keduanya **tidak dipakai sistem**,
jadi tiket ini dapat selesai tanpa jawaban itu. Pemilik: **DBA**.

#### Catatan

`[terverifikasi]` RD serupa (`BrowseReinsuranceType_RD`, class yang sama) ada di **sembilan modul
lain** — enumerasi ini **dipakai lintas lini**, bukan khas life. Versi modul lain difilter lewat
`Param.ID` / `Param.Name` / `Param.Note` / `Param.Type`, **bukan** `.Flag`.

⚠️ **Nama menyesatkan.** `TreatyLimit_TypeProtect` terbaca seperti logika "jenis proteksi limit";
`[terverifikasi]` ia **hanya lookup nama**. Jangan tiru namanya.

⚠️ **Normalisasi adalah keputusan skema baru, bukan tiruan.** `[data DBA]` Procedure **tetap
menerima dan menulis** `REINSTYPEID`/`REINSTYPENAME` ke tabel anak — jadi **lapisan repository yang
memutuskan apa yang dikirim**.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Master Contract Retro Life - 05 - Reinsurer — share, komisi, dan total yang ditampilkan tanpa memblokir

**Status:** ready-for-agent

**Blocked by:** 02 (reinsurer lahir di bawah kontrak), 04 (jenis reasuransi diwarisi dari kontrak)

#### Hasil & nilai pengguna

Sebagai **admin master**, saya menambahkan **reinsurer** ke sebuah kontrak beserta persentase share,
komisi, dan overriding commission-nya; dan sebagai **underwriter** saya melihat **total share** yang
sudah teralokasi — tanpa dihalangi menyimpan ketika kontrak masih saya susun bertahap.
*(User story 12–17 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas reinsurer; share/komisi sebagai desimal |
| `internal/repository` | Pemanggilan procedure penulis reinsurer; kueri total share per kontrak |
| `internal/services` | Validasi 0–100 per baris; penjumlahan total share |
| `internal/handlers` | Endpoint CRUD reinsurer; total share pada respons kontrak |
| `frontend/` | Grid reinsurer; **total share mencolok** + tanda "belum 100%" |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveMasterTreatyReinsurer_Life_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!SAVEMASTERTREATYREINSURER_LIFE_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/SaveMasterTreatyReinsurer_Life_SQL.xml` | `POOLDATA.INSERTREINSURER_LIFE` |
| `SaveSecurityLife_Act` | `ASM-FW-GISFW-…` / `SAVESECURITYLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SaveSecurityLife_Act.xml` | orkestrator simpan |
| `NewInputSecurityLife_Act` | `ASM-FW-GISFW-…` / `NEWINPUTSECURITYLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/NewInputSecurityLife_Act.xml` | baris baru |
| `SetSecurityLife_Act` | `ASM-FW-GISFW-…` / `SETSECURITYLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SetSecurityLife_Act.xml` | isi form |
| `SetErrorMessageReinsurer` | `@BASECLASS` / `SETERRORMESSAGEREINSURER` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SetErrorMessageReinsurer.xml` | **validasi 0–100** |
| `CountingPercentShare_Act` | `@BASECLASS` / `COUNTINGPERCENTSHARE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/CountingPercentShare_Act.xml` | ⚠️ **penjumlah**, bukan pembagi |
| `GetMasterReinsurerLifeList_SQl` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!GETMASTERREINSURERLIFELIST_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/GetMasterReinsurerLifeList_SQl.xml` | sumber total share |
| `BrowseDetailTreatyReisurerLife_RD` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `BROWSEDETAILTREATYREISURERLIFE_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Master Contract Retro Life/ReportDefinition/BrowseDetailTreatyReisurerLife_RD.xml` | daftar |

`[terverifikasi]` **Validasi per baris** (`SetErrorMessageReinsurer`):

```
@toDecimal(PctShare) > 100 || @toDecimal(PctShare) < 0   → galat
@toDecimal(Ricomm)   > 100 || @toDecimal(Ricomm)   < 0   → galat
```

`[terverifikasi]` **Total share** dihitung `CountingPercentShare_Act`:
`Local.TotalShare += @toDecimal(.PCTSHARE)` atas hasil
`SELECT PCTSHARE FROM POOLDATA.TREATYREINSURER_LIFE WHERE TREATYYEARID = … AND TREATYCONTRACTID = …`,
lalu hasilnya ditaruh di field **`STDRATING`** untuk **ditampilkan**.
⚠️ **Tidak ada perbandingan terhadap 100** di seluruh modul, dan `[data DBA]` **tidak ada pula di
procedure**.

`[data DBA]` Kolom: `PCTSHARE`, `COMMISION`, `OVR_COMM` bertipe **`NUMBER`** — desimal, bukan teks.

#### ADR terkait

**ADR-U-0003** (uang & persentase non-float), **ADR-U-0006** (identitas dari basis data), **ADR-U-0007**.

#### Acceptance criteria

- [ ] Reinsurer lahir **di bawah** satu kontrak; tanpa kontrak induk **ditolak**. *(AC 15 spec)*
- [ ] `PCTSHARE` dan komisi di luar rentang **0–100 ditolak**, dengan pesan yang **menyebut
      kolomnya**. *(AC 16 spec)*
- [ ] Total share seluruh reinsurer pada satu kontrak **dihitung dan ditampilkan**. *(AC 17 spec)*
- [ ] ⚠️ Total share **TIDAK memblokir penyimpanan** — kontrak dengan total ≠ 100% **tetap
      tersimpan**. *(AC 18 spec; `[keputusan work owner]`)*
- [ ] Kontrak dengan total ≠ 100% **ditandai mencolok** di layar. *(AC 19 spec)*
- [ ] Share dan komisi diperlakukan sebagai **desimal presisi arbitrer**; **tidak** melewati `float`.
      *(AC 49 spec; **ADR-U-0003**)*
- [ ] `REINSTYPEID` **tidak ditulis** pada baris reinsurer — diwarisi dari kontrak. *(AC 41 spec;
      tiket 04)*
- [ ] Menyimpan reinsurer yang sudah ada = **upsert**, bukan baris kedua; identitas baru dibuat basis
      data. *(AC 45, 4 spec)*
- [ ] `o_message` diperiksa; kegagalan ditampilkan. *(AC 46 spec — HTML dan konformansi di
      **tiket 10**)*

#### Blocker

**Tidak ada pemblokir.**

⚠️ **Keputusan tiket yang belum diambil — anti-dobel logis.** `[data DBA]` Upsert dikunci `ID`
mencegah dua baris ber-`ID` sama, **tidak** mencegah **dua reinsurer yang sama pada satu kontrak**.
Basis data tidak menjaganya (**nol `UNIQUE`**), dan Pega pun tidak. Bila anti-dobel dikehendaki, itu
aturan Go dan idealnya `UNIQUE` di basis data. **Putuskan sebelum tiket ini ditutup.**

#### Catatan

⚠️ **Nama menyesatkan — dua di tiket ini.** `[terverifikasi]`

| Yang tertulis | Yang sebenarnya |
| --- | --- |
| `CountingPercentShare_Act` "rumus pembagian share retro" | hanya **menjumlahkan** `PCTSHARE` |
| field `STDRATING` "standard rating" | menyimpan **total share** |

Di sistem baru, **beri nama yang jujur** — jangan bawa `STDRATING` sebagai tempat total share.

⚠️ `COMMISION` dieja demikian di basis data (satu `S`). Bawa nama kolomnya apa adanya di lapisan
repository; gunakan ejaan yang benar di lapisan domain.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Master Contract Retro Life - 06 - Security reinsurer — retrosesi atas retrosesi dan eksposur berjenjang

**Status:** ready-for-agent

**Blocked by:** 05 (security reinsurer lahir di bawah reinsurer)

#### Hasil & nilai pengguna

Sebagai **admin master**, saya menambahkan **security reinsurer** di bawah seorang reinsurer; dan
sebagai **underwriter** saya melihat **eksposur efektifnya terhadap treaty** — bukan angka mentah
yang menyesatkan. *(User story 18–20 di spec)*

⚠️ Inilah tiket dengan risiko salah-baca terbesar di modul ini.

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas security reinsurer; rujukan ke reinsurer induk |
| `internal/repository` | Pemanggilan procedure penulis security reinsurer |
| `internal/services` | **Perhitungan eksposur berjenjang** — share anak × share induk |
| `internal/handlers` | Endpoint CRUD; eksposur efektif pada respons |
| `frontend/` | Grid security reinsurer; **eksposur efektif ditampilkan berdampingan** dengan share mentah |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveMasterTreatySecurityReinsurer_Life_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!SAVEMASTERTREATYSECURITYREINSURER_LIFE_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/SaveMasterTreatySecurityReinsurer_Life_SQL.xml` | `POOLDATA.INSERTSECURITYREINSURER_LIFE` |
| `SaveSecurityReinsurerLife_Act` | `ASM-FW-GISFW-…` / `SAVESECURITYREINSURERLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SaveSecurityReinsurerLife_Act.xml` | orkestrator simpan |
| `SetSecurityReinsurerLife_Act` | `ASM-FW-GISFW-…` / `SETSECURITYREINSURERLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SetSecurityReinsurerLife_Act.xml` | isi form |
| `BrowseSecurityReinsurer_Life_RD` | `ASM-FW-GISFW-INT-TREATYSECURITYREINSURER_LIFE` / `BROWSESECURITYREINSURER_LIFE_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Master Contract Retro Life/ReportDefinition/BrowseSecurityReinsurer_Life_RD.xml` | daftar |
| `InboxSecurityReinsurerLife` | `DATA-PORTAL` / `INBOXSECURITYREINSURERLIFE` / `RULE-HTML-HARNESS` | `Master Contract Retro Life/Harness/InboxSecurityReinsurerLife.xml` (444.874 byte) | layar grid |

`[terverifikasi]` `INSERTSECURITYREINSURER_LIFE` menerima **`p_TREATYREINSURERID`** — penunjuk baris
**reinsurer induk**. `[data DBA]` Kini ditegakkan **FK** di basis data
(`TREATYSECURITYREINSURER_LIFE.TREATYREINSURERID` → `TREATYREINSURER_LIFE.ID`).

`[fakta bisnis — work owner]` **`PCTSHARE` security reinsurer adalah persentase DARI SHARE REINSURER
INDUKNYA**, **bukan** dari keseluruhan treaty. Ini **retrosesi atas retrosesi**.

> Reinsurer A memperoleh **40%** treaty. Security Reinsurer X di bawah A memperoleh **10%**.
> **Eksposur X terhadap treaty = 10% × 40% = 4%.**

`[data DBA]` `PCTSHARE` bertipe **`NUMBER`** — desimal, bukan teks.

#### ADR terkait

**ADR-U-0003** (persentase non-float), **ADR-U-0006**, **ADR-U-0007**,
**ADR-U-0001** (angka eksposur ini dikonsumsi konteks hilir).

#### Acceptance criteria

- [ ] Security reinsurer lahir **di bawah** satu reinsurer; tanpa reinsurer induk **ditolak**.
      *(AC 21 spec)*
- [ ] ⚠️ `PCTSHARE` security reinsurer dibaca sebagai **persentase dari share induknya**, **bukan**
      dari treaty. *(AC 22 spec; `[fakta bisnis — work owner]`)*
- [ ] **Eksposur efektif** terhadap treaty dihitung **share anak × share induk**, dan **itulah** angka
      yang ditampilkan sebagai eksposur. Test memuat kasus **`10% × 40% = 4%`**. *(AC 23 spec)*
- [ ] Mengubah share **induk** mengubah eksposur efektif **seluruh anaknya** — tanpa menyentuh baris
      anak. *(AC 24 spec)*
- [ ] Layar menampilkan **share mentah dan eksposur efektif berdampingan**, dengan label yang
      membedakan keduanya — sehingga `10%` tidak pernah terbaca sebagai eksposur terhadap treaty.
- [ ] `PCTSHARE` di luar rentang **0–100 ditolak**. *(sejalan AC 16 spec)*
- [ ] Share diperlakukan sebagai **desimal presisi arbitrer**; **tidak** melewati `float`.
      *(AC 49 spec; **ADR-U-0003**)*
- [ ] Menyimpan security reinsurer yang sudah ada = **upsert**, bukan baris kedua; identitas baru
      dibuat basis data. *(AC 45, 4 spec)*
- [ ] `o_message` diperiksa; kegagalan ditampilkan. *(AC 46 spec — HTML dan konformansi di
      **tiket 10**)*

#### Blocker

**Tidak ada.**

#### Catatan

⚠️ **Jebakan terbesar modul ini.** Menampilkan `10%` tanpa konteks induknya **salah besar** — ia
tampak empat kali lebih besar dari eksposur sebenarnya. Angka ini merambat ke perhitungan klaim dan
eksposur di **Claim Life** dan **Komite Claim Life**.

⚠️ **Anti-dobel logis belum diputuskan** — sama seperti tiket 05. `[data DBA]` Basis data tidak
mencegah dua security reinsurer yang sama di bawah satu induk.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — kaskade FK dan perilaku upsert hanya
berperilaku benar pada basis data sungguhan.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Master Contract Retro Life - 07 - Business — cakupan treaty, rate reasuransi, dan tampilan tabel rate

**Status:** ready-for-agent

**Blocked by:** 02 (business lahir di bawah kontrak), 04 (jenis reasuransi diwarisi dari kontrak)

#### Hasil & nilai pengguna

Sebagai **admin master**, saya menambahkan **jenis business** yang tercakup sebuah kontrak beserta
rate reasuransinya, dan saya dapat **melihat tabel rate** lebih dulu supaya memilih yang benar —
dengan rate tersimpan **persis seperti saya menuliskannya**. *(User story 21–23 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas business; **`RIRATE` sebagai teks** |
| `internal/repository` | Pemanggilan procedure penulis business; pembacaan tabel rate |
| `internal/services` | Wajib-isi; pewarisan jenis reasuransi dari kontrak |
| `internal/handlers` | Endpoint CRUD business; endpoint lihat tabel rate |
| `frontend/` | Grid business; dua layar tampilan rate |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveMasterTreatyBusiness_Life_SQL` | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE` / `ASM!SAVEMASTERTREATYBUSINESS_LIFE_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/SaveMasterTreatyBusiness_Life_SQL.xml` | `POOLDATA.INSERTBUSINESS_LIFE` |
| `SaveBusinessLife_Act` | `ASM-FW-GISFW-…` / `SAVEBUSINESSLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SaveBusinessLife_Act.xml` | orkestrator simpan |
| `NewInputBusinessLife_Act` | `ASM-FW-GISFW-…` / `NEWINPUTBUSINESSLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/NewInputBusinessLife_Act.xml` | baris baru |
| `SetBusinessListLife_Act` | `ASM-FW-GISFW-…` / `SETBUSINESSLISTLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SetBusinessListLife_Act.xml` | isi daftar |
| `SetParamRate` / `SetParamRateTable` | `ASM-FW-GISFW-…` / `SETPARAMRATE`, `SETPARAMRATETABLE` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/` | menyiapkan tampilan rate |
| `ViewRate` / `ViewRateTable` | `@BASECLASS` / `VIEWRATE`, `VIEWRATETABLE` / `RULE-OBJ-FLOW-ACTION` | `Master Contract Retro Life/FlowAction/` | **tampilan saja** |
| `BrowseRateLife_RD` | `ASM-FW-GISFW-INT-M_RATE_LIFE` / `BROWSERATELIFE_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Master Contract Retro Life/ReportDefinition/BrowseRateLife_RD.xml` | tabel rate |
| `BrowseRateLifeSummary` | `ASM-FW-GISFW-INT-RATE_LIFE_SUMMARY` / `BROWSERATELIFESUMMARY` / `RULE-OBJ-REPORT-DEFINITION` | `Master Contract Retro Life/ReportDefinition/BrowseRateLifeSummary.xml` | ringkasan rate |
| `BrowseBusinessLife_RD` | `ASM-FW-GISFW-INT-BUSINESS` / `BROWSEBUSINESSLIFE_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Master Contract Retro Life/ReportDefinition/BrowseBusinessLife_RD.xml` | daftar business |
| `InboxBusinessLifeReinsurers` | `DATA-PORTAL` / `INBOXBUSINESSLIFEREINSURERS` / `RULE-HTML-HARNESS` | `Master Contract Retro Life/Harness/InboxBusinessLifeReinsurers.xml` (486.911 byte) | layar grid |

`[terverifikasi]` `SetParamRate` menyalin `InputBusinessLife.RIRATEID` dan `RIRATE` ke `ParamID.*`;
`SetParamRateTable` menyalin dari `Param.*`. Keduanya menyiapkan **popup tampilan**, tidak mengubah
apa pun.

⚠️ `[data DBA]` **`RIRATE` bertipe `VARCHAR2(1000)` — TEKS**, bukan `NUMBER`. Rate reasuransi
disimpan sebagai string; kemungkinan memuat format bertingkat. **Jangan paksa menjadi angka.**

`[data DBA]` `TREATYBUSINESS_LIFE` adalah **satu-satunya tabel yang sudah punya PK sejak semula**
(`TREATYBUSINESS_LIFE_PK`, unique index, `ENABLE VALIDATE`).

#### ADR terkait

**ADR-U-0003** (uang non-float — **tidak berlaku pada `RIRATE`** yang memang teks),
**ADR-U-0006**, **ADR-U-0007**.

#### Acceptance criteria

- [ ] Business lahir **di bawah** satu kontrak, dengan kode business, nama, dan rate reasuransi.
      *(AC 25 spec)*
- [ ] ⚠️ **`RIRATE` tersimpan sebagai teks apa adanya** — **tidak** dikonversi, **tidak** diformat
      ulang, **tidak** dipaksa menjadi angka. Nilai yang ditulis dan dibaca kembali **identik
      karakter demi karakter**. *(AC 26 spec; `[data DBA]` kolom `VARCHAR2(1000)`)*
- [ ] Tabel rate dapat **dilihat** sebelum memilih, **tanpa mengubah apa pun** — popup murni baca.
      *(AC 27 spec)*
- [ ] Ringkasan rate juga dapat dilihat, dari sumber yang terpisah dari tabel rate.
- [ ] `REINSTYPEID` dan `TREATYYEAR` **tidak ditulis** pada baris business — diwarisi dari kontrak
      dan tahun. *(AC 41–42 spec; tiket 04)*
- [ ] Menyimpan business yang sudah ada = **upsert**, bukan baris kedua; identitas baru dibuat basis
      data. *(AC 45, 4 spec)*
- [ ] `o_message` diperiksa; kegagalan ditampilkan. *(AC 46 spec — HTML dan konformansi di
      **tiket 10**)*
- [ ] Daftar business yang dapat dipilih dibaca dari master business, bukan ditanam sebagai
      konstanta.

#### Blocker

**Tidak ada.**

#### Catatan

⚠️ **Jangan ubah tipe `RIRATE`.** `[data DBA]` Ia `VARCHAR2(1000)`. Mengubahnya menjadi angka
menuntut **konfirmasi format terlebih dahulu** dari Product+UW — di luar cakupan tiket ini, dan
tercatat di §Out of Scope spec.

⚠️ **OQ-066 — jangan buang perilaku karena memo.** `[data DBA]` `SetOutputParam_DT` (`@BASECLASS` /
`SETOUTPUTPARAM_DT`) bermemo **`not used`** **tetapi masih dirujuk tiga Harness dan dua Section**.
Bila tiket ini menyentuh Harness business, **telusuri rujukannya** sebelum menyimpulkan ia mati.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Master Contract Retro Life - 08 - Terapkan ke semua — pratinjau, konfirmasi, jejak audit, dan laporan sebagian-gagal

**Status:** ready-for-agent

**Blocked by:** 07 (fitur ini bekerja atas baris business)

#### Hasil & nilai pengguna

Sebagai **admin master**, saya menerapkan satu business ke **seluruh kontrak berjenis reasuransi
sama** sekaligus — tetapi hanya setelah sistem memberi tahu **berapa baris akan terpengaruh** dan
saya menyetujuinya, dan hasilnya tercatat sehingga dapat ditelusuri.
*(User story 24–26 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Kueri baris terdampak; penerapan per baris |
| `internal/services` | Pratinjau; orkestrasi penerapan; **penghitungan berhasil/gagal** |
| `internal/handlers` | Endpoint pratinjau; endpoint eksekusi; respons memuat rekap |
| `frontend/` | Dialog konfirmasi dengan jumlah baris; ringkasan hasil |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveBusinessToAllLife_Act` | `ASM-FW-GISFW-…` / `SAVEBUSINESSTOALLLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SaveBusinessToAllLife_Act.xml` | orkestrator penerapan massal |
| `SaveTreatyBusinessAll_Life_SQL` | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE` / `ASM!SAVETREATYBUSINESSALL_LIFE_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/SaveTreatyBusinessAll_Life_SQL.xml` | penulisan massal |

`[terverifikasi]` Gerbang pemilihan baris: **`.REINSTYPEID == Param.REINSTYPEID`** — penerapan
menyasar seluruh baris berjenis reasuransi sama.

`[terverifikasi]` Di Pega **tidak ada konfirmasi dan tidak ada jejak audit** — satu klik dapat
mengubah puluhan baris tanpa peringatan dan tanpa jejak.

⚠️ `[data DBA]` **`COMMIT` berada di dalam tiap procedure** — penerapan massal karena itu **tidak
atomik**. Bila gagal di tengah, sebagian baris **sudah tersimpan** dan tidak dapat di-rollback.

#### ADR terkait

**ADR-U-0007** (jejak audit — inti tiket ini), **ADR-U-0003**, **ADR-U-0015** (batas transaksi dipegang
aplikasi; di sini **tidak ada** transaksi menyeluruh yang mungkin).

#### Acceptance criteria

- [ ] ⚠️ **Terapkan ke semua** menampilkan **jumlah baris yang akan terpengaruh** dan **menunggu
      konfirmasi** sebelum dijalankan. *(AC 28 spec; `[keputusan work owner]` — penyimpangan sadar 5)*
- [ ] Pratinjau menyebut **jenis reasuransi** yang menjadi dasar pemilihan, sehingga pengguna tahu
      cakupannya.
- [ ] ⚠️ Hasil penerapan massal **tercatat di jejak audit**, termasuk **berapa baris berubah** dan
      **siapa** pelakunya. *(AC 29 spec; `[keputusan work owner]`)*
- [ ] Penerapan massal yang **dibatalkan** pada dialog konfirmasi **tidak mengubah apa pun**.
      *(AC 30 spec)*
- [ ] ⚠️ Bila penerapan gagal di tengah, sistem melaporkan **berapa berhasil dan berapa gagal** —
      **tidak** menampilkan sukses tunggal yang menyesatkan. *(AC 31 spec; `[data DBA]` tiap
      procedure commit sendiri)*
- [ ] Baris yang gagal **disebutkan** — pengguna tahu mana yang perlu diulang.
- [ ] `o_message` diperiksa **per baris**, bukan hanya pada baris terakhir. *(AC 46 spec)*
- [ ] Pratinjau yang menunjukkan **nol baris** memberi tahu pengguna, dan tidak menampilkan dialog
      konfirmasi yang sia-sia.

#### Blocker

**Tidak ada.**

#### Catatan

⚠️ **Sifat tidak-atomik harus terlihat, bukan disembunyikan.** `[data DBA]` Karena `COMMIT` ada di
dalam tiap procedure, tidak ada cara membungkus penerapan massal dalam satu transaksi. Menyajikannya
seolah "semua atau tidak sama sekali" akan **berbohong kepada pengguna**. Yang benar: laporkan
hasilnya apa adanya.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — sifat commit per baris tidak dapat difake
dengan jujur, dan justru itulah yang diuji.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Master Contract Retro Life - 09 - Kaskade hapus induk → anak, dengan popup konfirmasi sebelum apa pun terhapus

**Status:** ready-for-agent

**Blocked by:** 02 (kontrak), 05 (reinsurer), 06 (security reinsurer), 07 (business) — keempat
entitas harus ada agar kaskade dapat diuji utuh

#### Hasil & nilai pengguna

Sebagai **admin master**, saya menghapus sebuah induk dan **seluruh anaknya ikut terhapus** —
sehingga tidak ada baris yatim yang tertinggal — tetapi **hanya setelah** sistem memberi tahu apa dan
berapa yang akan ikut hilang, dan saya menyetujuinya. *(User story 27–29 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Penghitungan anak per induk; penghapusan |
| `internal/services` | **Aturan kaskade + pratinjau dampak**; larangan hapus tahun treaty |
| `internal/handlers` | Endpoint pratinjau hapus; endpoint eksekusi hapus |
| `frontend/` | **Dialog konfirmasi Ya/Batal** yang menyebut apa & berapa |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Sasaran |
| --- | --- | --- | --- |
| `DeleteTreatyLimit_SQL` | `ASM-FW-GISFW-INT-RETROCESSIONLIFE` / `ASM!DELETETREATYLIMIT_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/DeleteTreatyLimit_SQL.xml` | `treatycontract_life` |
| `DeleteSecurityReinsurer_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!DELETESECURITYREINSURER_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/DeleteSecurityReinsurer_SQL.xml` | `TREATYREINSURER_LIFE` |
| `DeleteSecurityReinsurerLife_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!DELETESECURITYREINSURERLIFE_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/DeleteSecurityReinsurerLife_SQL.xml` | `TREATYSECURITYREINSURER_LIFE` |
| `DeleteRowBusinessList` | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE` / `ASM!DELETEROWBUSINESSLIST` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/DeleteRowBusinessList.xml` | `treatybusiness_life` |
| `DeleteTreatyLimit_Act`, `DeleteSecurityLife_Act`, `DeleteSecurityReinsurerLife_Act`, `DeleteRowBusiness` | `ASM-FW-GISFW-…` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/` | orkestrator hapus |

`[terverifikasi]` **Keempatnya `DELETE … WHERE ID = …` datar, tanpa kaskade** — menghapus induk
meninggalkan anak yatim. `[terverifikasi]` **Tidak ada penghapus untuk `TREATYYEAR_LIFE`**.

`[data DBA]` **Keempat FK sudah terpasang** dengan mode **`ON DELETE CASCADE`**:

| FK | Kolom anak | → induk |
| --- | --- | --- |
| 1 | `TREATYCONTRACT_LIFE.IDTREATYYEAR` | `TREATYYEAR_LIFE.ID` |
| 2 | `TREATYREINSURER_LIFE.TREATYCONTRACTID` | `TREATYCONTRACT_LIFE.ID` |
| 3 | `TREATYSECURITYREINSURER_LIFE.TREATYREINSURERID` | `TREATYREINSURER_LIFE.ID` |
| 4 | `TREATYBUSINESS_LIFE.TREATYCONTRACTID` | `TREATYCONTRACT_LIFE.ID` |

#### ADR terkait

**ADR-U-0007** (jejak audit penghapusan), **ADR-U-0001** (master ini dirujuk Claim Life, Komite Claim
Life, dan Master Product Name Life — penghapusan berdampak lintas konteks).

#### Acceptance criteria

- [ ] ⚠️ Menghapus **kontrak** memunculkan **konfirmasi lebih dulu** yang menyebut **apa dan berapa**
      yang akan ikut terhapus — reinsurer, security reinsurer, dan business di bawahnya. *(AC 32
      spec; `[keputusan work owner]` — penyimpangan sadar 3)*
- [ ] ⚠️ Menghapus **reinsurer** memunculkan konfirmasi yang menyebut berapa **security reinsurer**
      akan ikut terhapus. *(AC 33 spec)*
- [ ] ⚠️ Menekan **Ya** menghapus induk **beserta seluruh sub-pohonnya**; **tidak ada baris yatim**
      yang tertinggal — dibuktikan dengan menghitung baris anak sesudahnya. *(AC 34 spec)*
- [ ] ⚠️ Menekan **Batal** membuat **tidak ada satu pun** baris terhapus — induk maupun anak.
      Dibuktikan dengan menghitung baris sebelum dan sesudah. *(AC 35 spec)*
- [ ] Menghapus baris **tanpa anak** tetap memerlukan konfirmasi, dengan pesan yang menyatakan
      **tidak ada anak** yang terpengaruh. *(AC 36 spec)*
- [ ] Menghapus **security reinsurer** atau **business** (baris daun) menghapus **hanya baris itu**.
      *(AC 37 spec)*
- [ ] ⚠️ **Tahun treaty tetap tidak dapat dihapus** — kaskade **tidak** membuka jalur hapus untuknya.
      *(AC 2 spec; `[fakta bisnis — work owner]`)*
- [ ] Pratinjau dampak dihitung **sesaat sebelum** dialog tampil, bukan dari data lama yang mungkin
      sudah berubah.
- [ ] Penghapusan tercatat di **jejak audit**: siapa, kapan, dan berapa baris ikut terhapus.
- [ ] Menghapus baris yang **sudah tidak ada** memberi pesan yang jelas, bukan galat mentah.

#### Blocker

**Tidak ada.**

#### Catatan

⚠️ **Basis data adalah lapis kedua, bukan pengganti.** `[data DBA]` FK `ON DELETE CASCADE` akan ikut
menghapus anak — tetapi **FK tidak dapat menjelaskan apa pun kepada pengguna**. Pesan yang menyebut
apa dan berapa anak **tetap tanggung jawab Go**, dan pratinjau dihitung sebelum penghapusan dimulai.

⚠️ **Perubahan arah keputusan.** Versi awal grilling mencatat "tolak hapus bila punya anak".
`[keputusan work owner]` **Direvisi menjadi kaskade + konfirmasi** — hapus induk memang dimaksudkan
menghapus seluruh sub-pohonnya. Tiket ini mengikuti keputusan **final**.

⚠️ **Kejanggalan penamaan** `[terverifikasi]`: `DeleteTreatyLimit_SQL` ber-class
`ASM-FW-GISFW-INT-RETROCESSIONLIFE` tetapi menghapus `treatycontract_life` — class tidak sejalan
dengan tabel sasarannya. Jangan tiru penamaannya.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — kaskade FK hanya berperilaku benar pada
basis data sungguhan; memalsukannya berarti tidak menguji apa pun yang penting.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Master Contract Retro Life - 10 - Penegakan `HASIL1` di kelima jalur simpan — gagal terang-terangan, HTML dibersihkan

**Status:** ready-for-agent

**Blocked by:** 01 (tahun), 02 (kontrak), 05 (reinsurer), 06 (security reinsurer), 07 (business) —
kelima jalur simpan harus ada agar konformansi dapat diuji

#### Hasil & nilai pengguna

Sebagai **admin master**, saya **diberi tahu bila penyimpanan ditolak** basis data — dengan pesan
yang **terbaca manusia**, bukan potongan markup — sehingga saya tidak pernah mengira data tersimpan
padahal tidak. *(User story 32–33 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Pemeriksaan hasil di **kelima** jalur simpan; **pembersihan HTML** di satu tempat |
| `internal/services` | Pemetaan hasil basis data menjadi galat domain |
| `internal/handlers` | Galat sampai ke pemanggil sebagai kegagalan, bukan sukses |
| `frontend/` | Pesan galat tampil bersih di dekat form yang bersangkutan |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Procedure |
| --- | --- | --- |
| `SaveMasterTreatyYear_Life_SQL` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` / `ASM!SAVEMASTERTREATYYEAR_LIFE_SQL` / `RULE-CONNECT-SQL` | `INSERTTREATYYEAR_LIFE` |
| `SaveMasterTreatyContract_Life_SQL` | `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` / `ASM!SAVEMASTERTREATYCONTRACT_LIFE_SQL` / `RULE-CONNECT-SQL` | `INSERTTREATYCONTRACT_LIFE` |
| `SaveMasterTreatyReinsurer_Life_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!SAVEMASTERTREATYREINSURER_LIFE_SQL` / `RULE-CONNECT-SQL` | `INSERTREINSURER_LIFE` |
| `SaveMasterTreatySecurityReinsurer_Life_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!SAVEMASTERTREATYSECURITYREINSURER_LIFE_SQL` / `RULE-CONNECT-SQL` | `INSERTSECURITYREINSURER_LIFE` |
| `SaveMasterTreatyBusiness_Life_SQL` | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE` / `ASM!SAVEMASTERTREATYBUSINESS_LIFE_SQL` / `RULE-CONNECT-SQL` | `INSERTBUSINESS_LIFE` |

`[terverifikasi]` Kelimanya mengembalikan `{OutputData.HASIL1 out}`. Sensus **27 Activity** modul
ini: **hanya `DeleteRowBusiness.xml`** yang menyebut `HASIL1`; **tidak satu pun activity `Save*`
membacanya**. **Galat procedure jatuh diam-diam.**

`[data DBA]` **Semantiknya:**

| Nilai `o_message` | Arti |
| --- | --- |
| **kosong / `NULL`** | **sukses** |
| **berisi teks** | **gagal** |

⚠️ `[data DBA]` Isinya **pesan bergaya UI Pega yang mengandung HTML** —
`<span style="color:red">…</span>` — ditambah `SQLERRM`.

#### ADR terkait

**ADR-U-0007** (jejak audit kegagalan), **ADR-U-0015** (kegagalan ditangani eksplisit, tidak ditelan).

#### Acceptance criteria

- [ ] ⚠️ **`o_message` DIPERIKSA** setelah **setiap** pemanggilan procedure di **kelima** jalur
      simpan: **kosong/NULL = sukses**, **berisi teks = gagal**. *(AC 46 spec;
      `[keputusan work owner]` — penyimpangan sadar 6)*
- [ ] ⚠️ Penyimpanan yang ditolak basis data **tidak pernah tampak berhasil** — API mengembalikan
      kegagalan, layar menampilkan galat. *(AC 46 spec)*
- [ ] ⚠️ **HTML dari pesan galat basis data tidak bocor** ke API maupun layar; pesan yang sampai ke
      pengguna **bersih dan terbaca**. Test yang menemukan tag markup pada respons API **gagal**.
      *(AC 47 spec; `[keputusan work owner]`)*
- [ ] **Uji konformansi lintas jalur**: test yang menemukan **jalur simpan mana pun** tanpa
      pemeriksaan hasil **gagal**. Kelima jalur diperiksa, bukan sebagian.
- [ ] Pembersihan HTML berada di **satu tempat**, bukan diulang lima kali.
- [ ] Pesan asli dari basis data **tetap tercatat** di log/jejak audit untuk penelusuran — yang
      dibersihkan hanya yang **sampai ke pengguna**.
- [ ] Kegagalan pada **penerapan massal** (tiket 08) diperiksa **per baris**, memakai jalur
      pemeriksaan yang sama. *(AC 31 spec)*
- [ ] Galat basis data yang **tidak dikenali** tetap dilaporkan sebagai kegagalan — tidak ada jalur
      yang menganggap galat tak dikenal sebagai sukses.

#### Blocker

**Tidak ada.**

#### Catatan

⚠️ **Tiket ini menutup celah, bukan membangun dari nol.** Tiket 01, 02, 05, 06, dan 07 masing-masing
sudah mengikat "`o_message` diperiksa, kegagalan ditampilkan" sebagai AC — karena jalur simpan yang
menelan galat bukanlah slice yang lengkap. Yang ditambahkan **di sini**: **pembersihan HTML yang
benar**, **pencatatan pesan asli**, dan **uji konformansi** yang menjamin tidak ada jalur tertinggal.

⚠️ **Nama menyesatkan** `[data DBA]`: parameter bernama `HASIL1` terbaca seperti kode hasil berangka;
ia sebenarnya **pesan galat teks**. Di lapisan domain, beri nama yang jujur.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — galat procedure hanya dapat dipicu dengan
jujur pada basis data sungguhan.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Master Contract Retro Life - 11 - Laporan kontrak dengan total share ≠ 100%

**Status:** ready-for-agent

**Blocked by:** 05 (total share dihitung di sana)

#### Hasil & nilai pengguna

Sebagai **manajemen**, saya ingin melihat **seluruh kontrak yang total sharenya belum 100%** dalam
satu tempat — sehingga celah alokasi risiko tidak tersembunyi di antara ratusan kontrak, meski
sistem sengaja tidak memblokir penyimpanannya. *(User story 17 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Kueri agregat total share per kontrak |
| `internal/services` | Penyaringan kontrak yang totalnya ≠ 100% |
| `internal/handlers` | Endpoint laporan |
| `frontend/` | Layar laporan; tiap baris dapat dibuka ke kontraknya |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `CountingPercentShare_Act` | `@BASECLASS` / `COUNTINGPERCENTSHARE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/CountingPercentShare_Act.xml` | penjumlah share per kontrak |
| `GetMasterReinsurerLifeList_SQl` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!GETMASTERREINSURERLIFELIST_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/GetMasterReinsurerLifeList_SQl.xml` | `SELECT PCTSHARE … WHERE TREATYYEARID = … AND TREATYCONTRACTID = …` |
| `BrowseTreatyContract_Life_RD` | `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` / `BROWSETREATYCONTRACT_LIFE_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Master Contract Retro Life/ReportDefinition/BrowseTreatyContract_Life_RD.xml` | daftar kontrak |

⚠️ `[terverifikasi]` **Laporan ini tidak ada di Pega.** Korpus hanya **menjumlahkan dan menampilkan**
total per kontrak yang sedang dibuka; tidak ada satu pun rule yang mencari kontrak bercelah secara
menyeluruh. Ini **kemampuan baru**, konsekuensi langsung dari keputusan **tidak memblokir**
penyimpanan (Q4).

#### ADR terkait

**ADR-U-0003** (persentase non-float — perbandingan terhadap 100 harus tepat, tanpa galat pembulatan),
**ADR-U-0001** (celah alokasi berdampak pada konteks hilir).

#### Acceptance criteria

- [ ] Tersedia **laporan kontrak dengan total share ≠ 100%**. *(AC 20 spec)*
- [ ] Laporan menampilkan, per kontrak: tahun treaty, jenis reasuransi, **total share terhitung**,
      dan **selisihnya terhadap 100%**.
- [ ] Kontrak yang totalnya **kurang dari 100%** dan yang **lebih dari 100%** keduanya muncul —
      laporan tidak hanya mencari yang kurang.
- [ ] Kontrak **tanpa reinsurer sama sekali** (total `0`) muncul di laporan, tidak terlewat karena
      dianggap tidak punya data.
- [ ] Perbandingan terhadap 100 memakai **desimal presisi arbitrer**; kontrak yang totalnya tepat
      `100` **tidak** muncul karena galat pembulatan. *(**ADR-U-0003**)*
- [ ] Tiap baris laporan dapat **dibuka ke kontraknya**, sehingga celah dapat langsung diperbaiki.
- [ ] Laporan dapat disaring setidaknya per **tahun treaty**.
- [ ] Laporan **tidak mengubah apa pun** — murni baca.

#### Blocker

**Tidak ada.**

#### Catatan

⚠️ **Mengapa laporan ini perlu ada.** `[keputusan work owner]` Total share **sengaja tidak
memblokir** penyimpanan, supaya kontrak dapat disusun bertahap. Konsekuensinya: celah alokasi
**pasti akan ada** pada suatu waktu. Tanpa laporan ini, celah itu hanya terlihat oleh orang yang
kebetulan membuka kontraknya — dan itu membuat keputusan "tidak memblokir" menjadi berbahaya.
Laporan inilah yang membuatnya aman.

⚠️ **Nama menyesatkan** `[terverifikasi]`: di Pega, total share ditaruh di field **`STDRATING`**
("standard rating"). Jangan bawa nama itu ke laporan — beri nama yang jujur.

#### Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Master Contract Retro Life - 12 - Migrasi skema — DDL lima tabel, PK, FK, dan sequence

**Status:** ready-for-agent — ✅ **OQ-001 ditutup** `[data DBA]`

**Blocked by:** None (can start immediately)

#### Hasil & nilai pengguna

Sebagai **tim migrasi**, saya ingin skema Oracle sistem baru menyimpan setiap nilai retro life
**tanpa berubah satu digit pun**, dengan **kunci primer dan kunci asing** yang menjamin keunikan dan
keterhubungan — sehingga rekonsiliasi tidak menemukan selisih dan tidak ada baris yatim yang mungkin.
*(User story 36–38 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `migrations/` | DDL kelima tabel + PK + FK + sequence |
| `internal/repository` | Pemetaan tipe kolom ↔ desimal presisi arbitrer; `RIRATE` sebagai teks |
| — | Skrip rekonsiliasi |

#### Tabel dan tipe `[data DBA]`

Sumber: `.scratch/master-contract-retro-life/ddl-tables-from-dba.md`. Tablespace `TBS_POOLDATA`.

| Tabel | Kolom kunci |
| --- | --- |
| `POOLDATA.TREATYYEAR_LIFE` | `ID`, `TREATYYEAR`, `UNDERWRITINGYEAR`, `USERID`, `TGLUPDATE`, `STARTDATE`, `ENDDATE` |
| `POOLDATA.TREATYCONTRACT_LIFE` | `ID`, `IDTREATYYEAR`, `REINSTYPEID`, `REINSTYPENAME`, `USERID`, `TGLUPDATE`, `IDR`, `USD`, `B_IDR`, `B_USD`, `IDR_SELISIH`, `USD_SELISIH`, `TREATYSTARTDATE`, `TREATYENDDATE` |
| `POOLDATA.TREATYREINSURER_LIFE` | `ID`, `TREATYYEARID`, `TREATYCONTRACTID`, `REINSTYPEID`, `REINSTYPENAME`, `REINSURERID`, `REINSURERNAME`, `PCTSHARE`, `COMMISION`, `OVR_COMM`, `USERID`, `TGLUPDATE` |
| `POOLDATA.TREATYSECURITYREINSURER_LIFE` | `ID`, `TREATYYEARID`, `TREATYCONTRACTID`, `TREATYREINSURERID`, `REINSURERID`, `REINSURERNAME`, `PCTSHARE`, `USERID`, `TGLUPDATE` |
| `POOLDATA.TREATYBUSINESS_LIFE` | `ID`, `TREATYYEARID`, `TREATYYEAR`, `TREATYCONTRACTID`, `REINSTYPEID`, `REINSTYPENAME`, `BIZCODE`, `BIZNAME`, `RIRATEID`, `RIRATE`, `USERID`, `TGLUPDATE` |

| Kelompok | Tipe `[data DBA]` | Catatan |
| --- | --- | --- |
| Uang: `IDR`, `USD`, `B_IDR`, `B_USD`, `IDR_SELISIH`, `USD_SELISIH` | **`NUMBER`** (tanpa presisi) | ✅ **ADR-U-0003 aman** — angka presisi penuh Oracle, bukan teks |
| Share/komisi: `PCTSHARE`, `COMMISION`, `OVR_COMM` | **`NUMBER`** | desimal |
| ⚠️ `RIRATE` | **`VARCHAR2(1000)`** | **teks** — bawa apa adanya |
| Tanggal: `TGLUPDATE`, `STARTDATE`, `ENDDATE`, `TREATYSTARTDATE`, `TREATYENDDATE` | **`DATE`** | menguatkan gerbang tahun (tiket 03) |
| Identitas: `ID`, `*ID` | `VARCHAR2(100)` | |
| Nama: `REINSTYPENAME`, `REINSURERNAME`, `BIZNAME` | `VARCHAR2(1000)` | lebar untuk cache nama |
| `USERID` | `VARCHAR2(100)` di year/contract/business; `VARCHAR2(1000)` di reinsurer/security | **tidak seragam** |
| Nullability | **semua kolom nullable, nol `NOT NULL`** | wajib-isi ditegakkan **di Go** |

#### Kunci dan sequence `[data DBA]`

⚠️ **Penyimpangan sadar 7 — `ID` menjadi PRIMARY KEY di kelima tabel.** Semula hanya
`TREATYBUSINESS_LIFE` yang punya (`TREATYBUSINESS_LIFE_PK`, unique index, `ENABLE VALIDATE`); empat
lainnya tidak — padahal upsert berkunci `ID`.

⚠️ **Penyimpangan sadar 8 — FK antar tabel, mode `ON DELETE CASCADE`** (selaras kaskade tiket 09):

| FK | Kolom anak | → induk |
| --- | --- | --- |
| 1 | `TREATYCONTRACT_LIFE.IDTREATYYEAR` | `TREATYYEAR_LIFE.ID` |
| 2 | `TREATYREINSURER_LIFE.TREATYCONTRACTID` | `TREATYCONTRACT_LIFE.ID` |
| 3 | `TREATYSECURITYREINSURER_LIFE.TREATYREINSURERID` | `TREATYREINSURER_LIFE.ID` |
| 4 | `TREATYBUSINESS_LIFE.TREATYCONTRACTID` | `TREATYCONTRACT_LIFE.ID` |

`[data DBA]` **Sequence**: `TREATYYEAR_LIFE_SEQ` `START WITH 44`, `MINVALUE 1`, `NOCACHE`,
`NOCYCLE`, `NOORDER` → identitas tahun berikutnya `1000044`. Empat sequence lain dirujuk procedure:
`TREATYCONTRACT_LIFE_seq`, `TREATYREINSURER_LIFE_SEQ`, `TREATYSECURITYREINSURER_LIFE_SEQ`,
`TREATYBUSINESS_LIFE_SEQ`.

#### ADR terkait

**ADR-U-0003** (uang non-float — `NUMBER` tanpa presisi → desimal presisi arbitrer di aplikasi),
**ADR-U-0006** (identitas dari sequence basis data), **ADR-U-0009** (migrasi penuh; koeksistensi
ditolak).

#### Acceptance criteria

- [ ] DDL kelima tabel target menyalin **tipe dan panjang kolom** dari sumber apa adanya; **tidak
      ada** kolom uang yang menjadi `FLOAT`/`BINARY_DOUBLE`. *(**ADR-U-0003**)*
- [ ] ⚠️ **`ID` adalah PRIMARY KEY di kelima tabel.** Skrip migrasi **memverifikasi keberadaannya**
      dan **menambahkannya bila belum ada** — tidak berasumsi. *(AC 51 spec;
      `[data DBA]` + `[keputusan work owner]`)*
- [ ] ⚠️ **Keempat FK terpasang** dengan mode **`ON DELETE CASCADE`**. Skrip **memverifikasi mode
      `ON DELETE`-nya**, bukan hanya keberadaan FK-nya — mode yang keliru (`RESTRICT`/`NO ACTION`)
      **bertentangan** dengan kaskade tiket 09. *(AC 52 spec)*
- [ ] ⚠️ `RIRATE` **tetap `VARCHAR2`** — **tidak** diubah menjadi numerik. *(AC 26 spec;
      `[data DBA]`)*
- [ ] Kelima **sequence** pindah dengan **nilai berjalan yang benar**, sehingga identitas baru
      **tidak pernah bertabrakan** dengan yang lama. Diuji dengan membuat satu baris di tiap tabel
      sesudah migrasi.
- [ ] Format identitas `'1' || lpad(seq, 6, '0')` **tetap berlaku** sesudah migrasi.
- [ ] Nilai uang dan share pindah **tanpa berubah satu digit pun**; rekonsiliasi membandingkan nilai
      lama dan baru **secara tepat**, bukan dengan toleransi.
- [ ] Nilai `DATE` pindah **tanpa pergeseran zona waktu**.
- [ ] Wajib-isi **tidak** ditambahkan sebagai `NOT NULL` tanpa keputusan terpisah — `[data DBA]`
      basis data lama **nol `NOT NULL`**, dan data lama mungkin memuat kolom kosong yang akan
      menolak migrasi. *(AC 53 spec — penegakan tetap di Go)*
- [ ] Index pendukung ada untuk kolom yang dipakai kueri hilir — khususnya kolom FK, dan
      `TREATYCONTRACTID` yang dipakai penghitungan total share.
- [ ] Migrasi dapat **dijalankan ulang dengan aman** dan punya **jalur mundur yang diuji**.
- [ ] Skema uji yang dipakai seluruh tiket lain dibangun **dari DDL yang sama** dengan produksi —
      bukan dari tiruan yang ditulis terpisah. *(AC 54 spec)*

#### Blocker

**Tidak ada.** ✅ **OQ-001 ditutup** `[data DBA]` 2026-09-15 — tipe, nullability, PK, FK, dan
sequence seluruhnya diketahui.

⚠️ **Verifikasi keadaan nyata, jangan berasumsi** `[keputusan work owner]`. Keadaan yang dinyatakan:
**kelima tabel sudah punya PK `ID`** (business dari DBA; empat lainnya ditambahkan work owner) dan
**keempat FK sudah terpasang dengan mode `ON DELETE CASCADE`**. Meski demikian, skrip migrasi
**memeriksa keadaan sebenarnya lebih dulu** lalu menambahkan yang belum ada — termasuk
**memverifikasi mode `ON DELETE`**, karena mode `RESTRICT`/`NO ACTION` akan **bertentangan langsung**
dengan kaskade tiket 09. Ini sudah tercermin di AC di atas.

⚠️ **Catatan sumber yang sudah usang.** `.scratch/master-contract-retro-life/ddl-tables-from-dba.md`
pada bagian akhirnya masih memuat dua pernyataan yang **tidak lagi berlaku**: (1) FK "masih nol", dan
(2) Q5 sebagai "tolak hapus induk berpunya-anak" dengan `ON DELETE CASCADE` disebut *bertentangan*.
Keduanya **sudah direvisi** — lihat `grilling-ronde-1-jawaban.md` §Q5 (kaskade + popup konfirmasi)
dan §Status OQ. Bila membaca berkas DDL itu, **pakai keadaan final di sini**.

#### Catatan

⚠️ **`USERID` tidak seragam** `[data DBA]`: `VARCHAR2(100)` di tiga tabel, `VARCHAR2(1000)` di dua
lainnya. Diseragamkan atau dibawa apa adanya — **putuskan di tiket ini**, jangan diam-diam
mengubahnya.

⚠️ **Anti-dobel logis belum diputuskan.** `[data DBA]` Basis data **nol `UNIQUE`** selain PK `ID`.
Bila anti-dobel logis dikehendaki (mis. satu reinsurer hanya sekali per kontrak), constraint-nya
ditambahkan **di sini** — tetapi keputusannya milik tiket 05/06.

#### Seam & perintah verifikasi

Migrasi diuji terhadap **skema uji Oracle nyata** — bukan mock.

```
go test ./internal/...
```

## Master Contract Retro Life - TAM - TAMBAHAN TIKET — hasil Grilling Ronde 2

**Tanggal:** 2026-09-19 · **Modul:** Master Contract Retro Life
**Menunjuk balik ke:** `issues/01-…` sampai `issues/12-…` *(dua belas tiket, tidak disunting)*

> ⛔ **NOL berkas tiket lama disunting.** ⛔ `CREATE TABLE` NOL · DDL NOL · kode NOL.
> ⛔ Nol butir `[terbuka]` / OQ dinyatakan tertutup.
> Bentuk tiket baru mengikuti `issues/05-reinsurer-share-dan-total.md` — tidak dikarang baru.

---

### BAGIAN (a) — TAMBAHAN untuk tiket 01–12 yang sudah ada

#### Tiket **03** — Gerbang tahun lewat nilai tanggal

| | |
| --- | --- |
| **Yang ditambahkan** | ⭐ **Gerbang konsistensi tahun di Pega TIDAK PERNAH BERJALAN.** Ketiga baris syaratnya di `Activity/SaveSecurityLife_Act.xml` langkah **4 · 5 · 6** dan tiga kembarannya di `Activity/SaveSecurityReinsurerLife_Act.xml` langkah **4 · 5 · 6** berflag prakondisi **`false`** — tersimpan tetapi dimatikan. Langkahnya berjalan **tanpa saringan**. |
| **Bukti** | kedua berkas di `Master Contract Retro Life/Activity/`, langkah 4 · 5 · 6, medan `pyStepsPreCondition` |
| **Menambah / meralat** | ⭐ **MERALAT dasar faktualnya** — dari *"Pega memeriksanya dengan cara keliru"* menjadi *"Pega tidak memeriksanya sama sekali"*. ⛔ Keputusan tiketnya **tidak berubah**: gerbangnya tetap **dibangun** memakai nilai tanggal |
| **Akibat untuk pengujian** | Test tidak boleh berbunyi *"meniru perilaku Pega"* — tidak ada perilaku untuk ditiru. Test menguji **aturan baru** |

#### Tiket **05** — Reinsurer, share, dan total

| | |
| --- | --- |
| **Yang ditambahkan** | ⭐ **Nol penjaga total share di kedua sisi.** Pega menjumlah lalu menampilkan *(`CountingPercentShare_Act` langkah **3.1** dan **4**)* — nol perbandingan terhadap 100. `[data DBA]` procedure juga **nol validasi bisnis**. ⚠️ Hasil penjumlahan ditulis ke medan bernama **`STDRATING`** |
| **Bukti** | `Activity/CountingPercentShare_Act.xml` langkah 3.1 · 4 · `procedure-bodies-from-dba.md` |
| **Menambah / meralat** | **MENAMBAH** — menguatkan keputusan Q4 *(tampilkan, jangan blokir)* dengan bukti dua sisi |
| **Peringatan migrasi** | Bila kolom bernama *rating* ditemui saat migrasi, **periksa isinya** — bisa jadi ia total share |

#### Tiket **07** — Business dan tampilan rate

| | |
| --- | --- |
| **Yang ditambahkan** | ⭐ ⚠️ **`RIRATE` disimpan sebagai TEKS sepanjang 1000 karakter**, bukan angka. `RIRATEID` berdampingan dengannya. Keduanya punya **4 penulis** *(`NewInputBusinessLife_Act` · `SetBusinessListLife_Act` · `SetParamRate`)* dan **8–9 pembaca**, jadi keduanya **dipakai**, bukan sisa |
| **Bukti** | `ddl-tables-from-dba.md` temuan 3 · sensus 66 berkas di `grilling-ronde-2.md` §B3 |
| **Menambah / meralat** | ⭐ **MERALAT** — ronde 1 menulis `RIRATE` *"belum saya telusur"*; kini terbaca tipenya, tetapi **isinya belum** |
| **⛔ MEMBLOKIR** | ⭐ **Tiket ini MACET** sampai **Pertanyaan A** dijawab: bentuk kolom, validasi, dan migrasi nilainya ketiganya bergantung apakah `RIRATE` satu angka atau teks berstruktur |

#### Tiket **08** — Terapkan ke semua dengan pratinjau

| | |
| --- | --- |
| **Yang ditambahkan** | ⚠️ **Gerbangnya terbelah.** `Activity/SaveBusinessToAllLife_Act.xml` langkah **3** berflag **`false`** — **mati**; tetapi langkah **3.1** dan **3.2** berflag **`true`** — **berlaku**. Syaratnya sama persis: `.REINSTYPEID == Param.REINSTYPEID` |
| **Bukti** | `Activity/SaveBusinessToAllLife_Act.xml` langkah 3 · 3.1 · 3.2 |
| **Menambah / meralat** | **MENAMBAH** — ronde 1 mengutip gerbang ini tanpa membaca arahnya |
| **Akibat** | Penyaringan *"hanya baris berjenis sama"* **berjalan di tingkat anak**, bukan di tingkat induk. Pratinjau harus menunjukkan baris mana yang benar-benar tersentuh |

#### Tiket **09** — Kaskade hapus dengan popup konfirmasi

| | |
| --- | --- |
| **Yang ditambahkan** | ⭐ `[data DBA]` **Basis data kini mengaskade sendiri** — **empat** kunci tamu bermode `ON DELETE CASCADE`: kontrak→tahun · reinsurer→kontrak · security→reinsurer · business→kontrak. Kelima tabel juga sudah punya kunci utama |
| **Bukti** | `ddl-tables-from-dba.md`, bagian *KEADAAN FINAL* |
| **Menambah / meralat** | **MENAMBAH** — *"nol kaskade"* **tetap benar untuk Pega**; yang berubah adalah **siapa** yang mengaskade |
| **Akibat untuk pengujian** | ⚠️ Test kaskade **tidak boleh** mengandaikan aplikasi yang menghapus anak. Test harus membuktikan **hasil akhirnya**, dan memastikan popup konfirmasi muncul **sebelum** basis data bertindak |

#### Tiket **10** — Penegakan `HASIL1` di lima jalur simpan

| | |
| --- | --- |
| **Yang ditambahkan** | ⭐ **Tujuh berkas dari 66 menyebut `HASIL1`**, bukan satu — kelima rule SQL simpan **mendeklarasikannya**. ⭐ Tetapi **nol activity membacanya**. Satu-satunya activity yang menyentuh penampung serupa memakai **`HASIL12`**, nama yang **berbeda**. `[data DBA]` Saat gagal, procedure menaruh pesan galat berikut teks galat basis data ke penampung itu, **menjalankan pembatalan**, lalu selesai — sehingga **seluruh perubahan dibatalkan tanpa ada yang melihatnya** |
| **Bukti** | lima berkas di `RDBList/` · `Activity/DeleteRowBusiness.xml` · `procedure-bodies-from-dba.md` |
| **Menambah / meralat** | ⭐ **MERALAT** *(jumlah berkas: 1 → 7)* **dan MENAMBAH** *(apa yang hilang)* |
| **Tambahan bentuk** | Pesan galat berisi **HTML**; sistem baru **tidak boleh meneruskannya mentah** |
| **`[terbuka]` yang menyentuhnya** | **Pertanyaan D** — apakah `HASIL12` salah ketik atau penampung lain |

#### Tiket **11** — Laporan kontrak dengan total share tidak 100%

| | |
| --- | --- |
| **Yang ditambahkan** | ⚠️ **Sebelas dari dua belas laporan modul ini belum pernah dibaca**, termasuk `BrowseRateLife_RD` dan `BrowseRateLifeSummary`. ⛔ Kelas sumber, kolom, dan saringannya **tidak berhasil dibaca** di ronde ini — disisir tiga medan, ketiganya nol |
| **Bukti** | sensus `ReportDefinition/` di `grilling-ronde-2.md` §B5 |
| **Menambah / meralat** | **MENAMBAH** — butir `[terbuka]`, ⛔ bukan pembatalan |
| **Akibat** | Bila salah satu laporan itu **sudah** menghitung total share, tiket ini berubah dari *"bangun laporan"* menjadi *"tiru laporan"* |

#### Tiket **12** — Migrasi skema lima tabel

| | |
| --- | --- |
| **Yang ditambahkan** | ⭐ `[data DBA]` **Kunci utama kini ada di kelima tabel** · **empat kunci tamu bermode kaskade** sudah dipasang · **wajib-isi tetap nol** di basis data. ⚠️ Kolom angka **tidak dibatasi digitnya** oleh basis data — batas ketelitian sepenuhnya keputusan aplikasi. ⭐ Satu kolom, `RIRATE`, bertipe **teks** |
| **Bukti** | `ddl-tables-from-dba.md` temuan 1 · 2 · 3 · 4 · 6, dan bagian *KEADAAN FINAL* |
| **Menambah / meralat** | **MENAMBAH** |
| **Peringatan** | ⚠️ Kunci tamu kaskade **mungkin dipasang setelah data lama dibuat** — data yatim lama bisa masih ada. Migrasi wajib **memverifikasi keadaan nyata**, bukan berasumsi |

---

### BAGIAN (b) — TIKET BARU

⭐ **Dua tiket baru**, diberi nomor mulai **13** supaya tidak bertabrakan.

---

### 13: Validasi simpan yang di Pega tertulis tetapi tidak berjalan

**Status:** ⚠️ **menunggu keputusan work owner** *(Pertanyaan B)* — bukan `ready-for-agent`

**Blocked by:** 02 (kontrak lahir), 05 (reinsurer lahir), 06 (security reinsurer lahir)

#### Hasil & nilai pengguna

Sebagai **admin master**, saya diberi tahu **saat menyimpan** bila nama reinsurer atau tanggal mulai
belum saya isi — bukan menemukan barisnya kosong berhari-hari kemudian saat laporan dibuat.

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/services` | Validasi wajib-isi sebelum memanggil jalur simpan |
| `internal/handlers` | Pesan galat per medan pada respons simpan |
| `frontend/` | Penanda medan wajib; galat ditampilkan di medannya |

#### Perilaku Pega yang ditiru, berikut buktinya

⭐ **Yang ditiru adalah ketiadaannya** — dan itu disengaja untuk dicatat, bukan untuk diteruskan.

| Rule | Langkah | Syarat yang tertulis | Flag |
| --- | --- | --- | --- |
| `Activity/SaveSecurityLife_Act.xml` | **4 · 5 · 6** | nama reinsurer kosong · tanggal mulai kosong · konsistensi tahun | ⛔ **`false` — mati** |
| `Activity/SaveSecurityReinsurerLife_Act.xml` | **4 · 5 · 6** | tiga syarat yang **sama persis** | ⛔ **`false` — mati** |
| `Activity/SaveBusinessLife_Act.xml` | **6** | tiga syarat yang **sama persis** | ⛔ **`false` — mati** |
| `Activity/SaveBusinessToAllLife_Act.xml` | **3** | jenis reasuransi sama | ⛔ **`false` — mati** |

⭐ **21 dari 42 baris syarat di jalur simpan modul ini bergerbang mati.**
`[data DBA]` Dan kelima procedure penulis **nol validasi bisnis**. ⭐ **Tidak ada penjaga di kedua
sisi.**

#### Keputusan work owner yang mengikat

⚠️ **BELUM ADA.** Tiket ini **menunggu Pertanyaan B** di `grilling-ronde-2.md` §G:

> *Apakah sistem baru **membangun** pemeriksaan itu, atau **meniru keadaan sekarang** yang tanpa
> pemeriksaan?*

⭐ **Rekomendasi asisten — bukan keputusan:** **bangun**, dengan dua syarat: *(i)* dinyatakan
sebagai **penyimpangan sadar** di spec, dan *(ii)* migrasi data lama **tidak menolak** baris yang
sudah terlanjur kosong. Alasannya: pemeriksaan itu **ditulis sendiri oleh pengembangnya** lalu
dimatikan — bentuk yang menyerupai keputusan sementara, bukan keputusan rancangan.

#### Apa yang harus diuji

- [ ] simpan reinsurer **tanpa nama** → **ditolak**, dan pesannya menyebut medan mana
- [ ] simpan **tanpa tanggal mulai** → **ditolak**
- [ ] simpan dengan tahun kontrak **tidak sesuai tahun treaty** → **ditolak**
- [ ] baris **lama** yang sudah terlanjur kosong → **masih dapat dibuka dan diperbaiki**, tidak
      dikunci oleh validasi baru
- [ ] bila keputusan jatuh ke **meniru**: seluruh test di atas **dibalik** menjadi *"diterima"*

#### Butir `[terbuka]` yang menyentuhnya

- ⭐ **Pertanyaan B** — dibangun atau ditiru *(memblokir tiket ini)*
- ⚠️ **H2 `grilling-ronde-2.md`** — arti flag `false` **belum diuji di layar Pega** untuk modul ini.
  Bila artinya ternyata lain, **21 baris berbalik menjadi hidup** dan tiket ini **gugur seluruhnya**

---

### 14: Empat layar dan sebelas laporan yang belum pernah dibaca

**Status:** ⚠️ **belum dapat dikerjakan** — pekerjaan pembacaan korpus, bukan pembangunan

**Blocked by:** tidak ada

#### Hasil & nilai pengguna

Sebagai **pembangun**, saya tidak menemukan di tengah pekerjaan bahwa layar yang saya bangun
ternyata **sudah ada bentuknya** di sistem lama dan saya membangunnya berbeda tanpa alasan.

#### Area codebase

⛔ **Belum dapat ditentukan** — tiket ini menghasilkan **bacaan**, bukan kode.

#### Perilaku Pega yang ditiru, berikut buktinya

⭐ **Belum terbaca.** Yang terbaca hanyalah **keberadaan dan ukurannya**:

| Section belum dibaca | Byte | Disambung oleh |
| --- | --- | --- |
| `InputRetroLimitReinsurers` | **541 845** | `Harness/InboxRetroLimitReinsurers` |
| `InputSecurityLifeReinsurers` | **496 437** | `Harness/InboxRetroLifeReinsurersList` |
| `InputBusinessLifeReinsurers` | **449 805** | `Harness/InboxBusinessLifeReinsurers` |
| `InputDtlRetrocessionLife` | **203 911** | ⛔ belum terbaca |

⭐ **Total 1,69 MB — lebih besar dari keempat layar yang SUDAH dibaca (1,10 MB).**

**Sebelas laporan belum dibaca:** `BrowseBusinessLife_RD` · `BrowseCedingCoLife_RD` ·
`BrowseDetailTreatyReisurerLife_RD` · **`BrowseRateLifeSummary`** · **`BrowseRateLife_RD`** ·
`BrowseRetrocessionLife_RD` · `BrowseSecurityReinsurer_Life_RD` · `BrowseTreatyBusiness_Life_RD` ·
`BrowseTreatyContract_Life_RD` · `BrowseTreatyYear_Life_RD` · `BrowseTreatyYear_RD`.

⚠️ **Gerbang layar yang sudah terbaca dari keempat pembungkusnya:** **40 hidup · 15 mati**.
Penanda mati di sini berarti **niat menyembunyikan**, bukan gerbang rusak.

#### Keputusan work owner yang mengikat

⛔ **Belum ada** — dan **belum diperlukan**. Tiket ini membaca, bukan memutuskan.

#### Apa yang harus diuji

⛔ **Belum ada test** — hasilnya berupa bacaan yang menambal **tiket 04 · 05 · 06 · 07 · 11**.

**Yang harus dihasilkan tiket ini:**

- [ ] untuk tiap **4 Section**: kolom yang ditampilkan, gerbang layarnya, dan medan yang dapat disunting
- [ ] untuk tiap **11 laporan**: kelas sumbernya, saringannya, kolom yang ditampilkan
- [ ] jawaban apakah **`BrowseRateLife_RD` / `BrowseRateLifeSummary`** adalah sumber tampilan rate
- [ ] jawaban apakah salah satu laporan **sudah menghitung total share** — bila ya, **tiket 11**
      berubah bentuk

#### Butir `[terbuka]` yang menyentuhnya

- ⚠️ medan kolom Harness **tidak ketemu** di `pyPropertyName` — medan yang benar belum diketahui
- ⚠️ kelas/kolom/saringan laporan **tidak ketemu** di tiga medan yang disisir
- ⭐ **Pertanyaan A** — bila `BrowseRateLife_RD` menerangkan bentuk `RIRATE`, Pertanyaan A bisa
  terjawab dari korpus tanpa menunggu work owner

---

#### Urutan ketergantungan

```
14 (baca 4 layar + 11 laporan)   ──┐
   tidak diblokir siapa pun        │  hasilnya menambal 04 · 05 · 06 · 07 · 11
                                   │  dan BISA menjawab Pertanyaan A
                                   ▼
Pertanyaan A  (RIRATE teks?)  ───► 07  (business dan tampilan rate)   ⛔ MACET tanpa ini
                                   ▼
Pertanyaan B  (validasi?)     ───► 13  (validasi simpan)              ⛔ MACET tanpa ini
                                   │
                                   └─► menyentuh 05 · 06 · 07 · 10
```

⭐ **Tiket 14 didahulukan** — ia satu-satunya yang **tidak diblokir apa pun**, dan hasilnya
**berpeluang menjawab Pertanyaan A dari korpus** sehingga tiket 07 tidak perlu menunggu work owner.

⛔ **Tiket 13 tidak boleh dimulai** sebelum Pertanyaan B dijawab — dan ⚠️ **bisa gugur seluruhnya**
bila arti flag `false` ternyata berbeda *(lihat `grilling-ronde-2.md` §H2)*.

---

#### Rekap

| | Jumlah |
| --- | --- |
| Tiket lama **ditambal** | ⭐ **8** — 03 · 05 · 07 · 08 · 09 · 10 · 11 · 12 |
| di antaranya **MERALAT** | **3** — 03 · 07 · 10 |
| di antaranya **MENAMBAH** | **5** — 05 · 08 · 09 · 11 · 12 |
| Tiket lama **tidak tersentuh** | **4** — 01 · 02 · 04 · 06 |
| ⭐ **Tiket BARU** | **2** — **13** *(validasi simpan yang tidak berjalan)* · **14** *(empat layar dan sebelas laporan yang belum dibaca)* |

⛔ **NOL berkas tiket lama disunting.**
