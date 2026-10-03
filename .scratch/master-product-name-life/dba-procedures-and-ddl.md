# Body Procedure + DDL — dari DBA (menutup OQ-002 & OQ-001, Master Product Name Life)

Tanggal: 2026-09-15
Sumber: **DBA / work owner** (dikirim langsung; bukan korpus Pega).
Status: `[data DBA]` — menutup **OQ-002** & **OQ-001** untuk modul ini.

> ⚠️ **Dua temuan memaksa peninjauan ulang keputusan D1 & D2.** Lihat bagian "Konflik dengan
> keputusan sebelumnya".

---

## 1. `PEGA_M_PRODUCT_LIFE(DataPega CLOB, IDPega, ErrMsg OUT, IDPegaOut OUT, StsSave OUT)`

- **Mode:** `IDPega IS NULL` → INSERT (ID baru); selain itu → UPDATE `WHERE ID=IDPega`.
- **ID baru:** `concat('1', lpad(M_PRODUCT_LIFE_SEQ.nextval, 5, '0'))` → mis. `100001` (**5 digit**).
- **Tulis:** hanya `INSERT INTO M_PRODUCT_LIFE(ID, JSONDATA)` / `UPDATE … SET JSONDATA=…`.
  **Hanya kolom `JSONDATA`** — seluruh atribut produk ada di dalam JSON.
- **Trik insert:** `replace(DataPega, 'UnknownId', id_baru)` — placeholder `UnknownId` di JSON
  diganti ID asli setelah sequence keluar.
- ⚠️ **`COMMIT` DI DALAM procedure** (setelah INSERT/UPDATE). Rollback hanya pada exception.
- **Keluaran:** `StsSave = 100` sukses / `99` gagal. `ErrMsg` = teks (juga saat sukses:
  "Data Already Saved With ID : <id>"). `IDPegaOut` = ID final.

## 2. `PEGA_M_PRODUCT_INWARD_LIFE(...)` — pola IDENTIK

- Sama persis, tabel `m_productinward_life`, sequence `M_PRODUCT_INWARD_LIFE_SEQ`.
- **Tulis hanya `(ID, JSONDATA)`.** `COMMIT` di dalam. `StsSave` 100/99.

## 3. `GET_TOKEN_STORAGE(VAPPNAME, VUSERINPUT, VAKSESTOKEN OUT, VERRMSG OUT)`

- Validasi `VAPPNAME` tidak boleh null.
- Ambil token belum-expired terbaru dari `GCP_IMAGE` (`INPUTDATE > SYSDATE`, order desc, 1 row).
- Bila `NO_DATA_FOUND` → buat token baru: `RAWTOHEX(STANDARD_HASH('ASMAPP'||timestamp,'MD5'))`,
  masa berlaku **1 menit** (`SYSDATE + 1 MINUTE`), insert ke `GCP_IMAGE`.
- Pola sama dengan token storage Claim Life (`GCP_IMAGE`).

---

## DDL

### `M_PRODUCT_LIFE`
| Kolom | Tipe |
| --- | --- |
| `ID` | `VARCHAR2(6)` |
| `JSONDATA` | `CLOB` (SECUREFILE) — **constraint `JSONDATA IS JSON` (ENABLE VALIDATE)** |
| `RIRISKID` | `VARCHAR2(10)` |
| `RIRISK` | `VARCHAR2(100)` |
| `PRODUCTNAME` | `VARCHAR2(1000)` |
| `BEGIN_DATE` | `DATE` |

Index: `INDEX2 (ID)`, `M_PRODUCT_LIFE_INDEX1 (ID, PRODUCTNAME)`. **Tanpa PK** (hanya index non-unique).

### `M_PRODUCTINWARD_LIFE`
| Kolom | Tipe |
| --- | --- |
| `ID` | `VARCHAR2(6)` |
| `JSONDATA` | `CLOB` — constraint `JSONDATA IS JSON` (ENABLE VALIDATE) |

**Hanya 2 kolom.** Tanpa PK, tanpa index tambahan yang dikirim.

---

## Temuan mengikat

1. ⚠️ **Kedua tabel PENYIMPAN UTAMA = JSON.** `M_PRODUCTINWARD_LIFE` **hanya `ID`+`JSONDATA`** —
   semua atribut inward (MINAGE/MAXAGE/BROKERAGE/EXTRAPREMI/…) ada **di dalam JSON**, bukan kolom.
   `M_PRODUCT_LIFE` = JSON + **4 kolom flatten** (`RIRISKID`, `RIRISK`, `PRODUCTNAME`, `BEGIN_DATE`)
   yang diisi terpisah (mis. `SaveProductNameLIfeFlat`) untuk index/query — bukan oleh procedure JSON.
2. ⚠️ **Procedure COMMIT sendiri** → dua tabel disimpan dua procedure yang masing-masing commit.
   **TIDAK atomik di sisi DB.**
3. ✅ ID = `'1'+lpad(seq,5)` (5 digit). Sequence: `M_PRODUCT_LIFE_SEQ`, `M_PRODUCT_INWARD_LIFE_SEQ`.
4. ✅ `StsSave` 100=sukses/99=gagal; `ErrMsg` teks; `IDPegaOut` = ID final. (3 keluaran — D5.)
5. ⚠️ Trik `replace('UnknownId', id)` — JSON dikirim dengan placeholder, DB menukar dengan ID nyata.
6. ✅ `JSONDATA IS JSON` constraint — DB menolak JSON tak valid.

## Konflik dengan keputusan sebelumnya (PERLU KEPUTUSAN ULANG WORK OWNER)

- **D1 (atomik, rollback keduanya)** vs realita **procedure COMMIT sendiri** → di Pega TIDAK atomik.
  Pertanyaan ke work owner: di sistem baru, apakah dua tabel digabung dalam satu transaksi Go
  (menyimpang dari Pega, butuh tulis ulang persist bukan panggil procedure apa adanya), ATAU tiru
  Pega (dua commit terpisah, risiko timpang)?
- **D2 (simpan flat, bukan JSON)** vs realita **penyimpan utama = JSON** (`JSONDATA IS JSON`,
  `M_PRODUCTINWARD_LIFE` cuma ID+JSON). Pertanyaan: sistem baru benar-benar pindah ke kolom
  relasional penuh (menyimpang jauh dari DB existing + migrasi JSON→kolom), ATAU pertahankan JSON?


---

## Isi `JSONDATA` — dari contoh data nyata `[data DBA]` (menutup OQ-JSON-PRODUCT)

Sumber: contoh produk `ID=100421` (CAR — Asuransi Jiwa Kumpulan GTL/PA-AB/CI-ADD 2026-2027).
Menutup **OQ-JSON-PRODUCT**. Karena D2 = relasional penuh, list bersarang → **tabel anak**.

### `M_PRODUCTINWARD_LIFE.JSONDATA` — field skalar (→ kolom `m_productinward_life`)

| Field | Contoh | Dugaan tipe (untuk skema baru) |
| --- | --- | --- |
| `ID` / `PRODUCTID` | `100421` | VARCHAR2(6) — sama, PK |
| `BEGIN` | `01/03/2026` | DATE (mulai) |
| `MATURE` | `28/02/2027` | DATE (akhir) |
| `STNC` | `26/03/2026` | DATE |
| `CURRENCY` | `IDR` | VARCHAR2 |
| `CEDINGLIMIT` | `150000000` | NUMBER (uang) |
| `CEDINGRETENTIONNUM` | `50` | NUMBER (persen/retensi) |
| `MINAGE`/`MAXAGE` | `22`/`70` | NUMBER |
| `MINSUMINSURED`/`MAXSUMINSURED` | `0`/`1175000000` | NUMBER (uang) |
| `MAXSUMREASURED` | `1025000000` | NUMBER (uang) |
| `RNMLIMITNUM` | `1500000000` | NUMBER (uang) |
| `RNMSHARE` | `100` | NUMBER (persen) |
| `MAXCONTRACT` | `1` | NUMBER |
| `MAXDATARECEIVE` | `90` | NUMBER (hari) |
| `MAXEXPIREDCLAIM` | `180` | NUMBER (hari) |
| `PAYMENT` | `1` | NUMBER/kode |
| `INSURED` | teks panjang | VARCHAR2 besar/CLOB |
| `SUBJECTTO` | teks multiline | CLOB |
| `POLICYHODER`(sic)/`POLICYHODERNAME` | `ASM-SFAGIS-WORK-ORG ORG-66` / nama | VARCHAR2 — ⚠️ typo `POLICYHODER` (tanpa L kedua) |

### `M_PRODUCT_LIFE.JSONDATA` — field skalar (→ kolom `m_product_life`)

| Field | Contoh | Catatan |
| --- | --- | --- |
| `ID` | `100421` | sama dengan inward (D1: ID sama) |
| `PRODUCTNAME` | `ASURANSI JIWA KUMPULAN` | sudah jadi kolom flatten di DDL |
| `INWARDNAME` | teks panjang (judul nota) | VARCHAR2 besar |
| `TREATYNUMBER` | `007/RNML-CAR/III/2026` | VARCHAR2 |
| `CEDING`/`CEDINGID` | nama / `L0000014` | dari ChooseCeding |
| `SOBID`/`SOBNAME` | `L0000014` / nama | Source of Business |
| `POLICYHODER`/`POLICYHODERNAME` | org / nama | dari ChoosePolicyHolder (typo sama) |
| `CAUSE`/`CAUSEID` | `ANY CAUSE` / `100004` | dari ChooseCauseOfLoss |
| `RIRISK`/`RIRISKID` | teks / `1000117` | dari ChooseRIRisk (kolom flatten) |
| `RICOMM` | `0` | NUMBER (persen komisi) |
| `IsORS` | `false` | ⚠️ jalur OR mati — abaikan (D6) |
| `IsView` | `false` | flag tampilan |
| `CREATEOP`/`UPDATEOP` | operator | audit |
| `Comment` | teks | ringkasan |

### List bersarang di `M_PRODUCT_LIFE.JSONDATA` → **tabel anak** (D2 relasional)

| List | Isi baris | → tabel anak (usulan) |
| --- | --- | --- |
| `CommentList` | `Date`, `OperatorName`, `Suggest` | `product_life_comment` |
| `DocumentClaim` | `Document` (+ audit) | `product_life_document_claim` |
| `PlanList` | `Plan`, `PlanID`, `Name`, `Benefit`, `RIRATE`, `RIRATEID` | `product_life_plan` — dari ChooseRIRate |
| `UnderwritingLimitList` | `Description`, `Medical`(FCL/NM/MEDIS), `MinAge`, `MaxAge`, `MinInsured`, `MaxInsured` | `product_life_uw_limit` |
| `FinancialUnderwritingList` | (kosong di contoh) | `product_life_fin_uw` |
| `LienClause` | (kosong) | `product_life_lien` |
| `OutwardList` | (kosong) | `product_life_outward` |

⚠️ **Konsekuensi D2:** "relasional penuh" = **produk induk (2 tabel: product_life + productinward_life
by ID) + 5 tabel anak** (comment, document_claim, plan, uw_limit, fin_uw). `[terverifikasi]`
`LienClause` = field skalar (bukan tabel), `OutwardList` = dead code jalur OR. Migrasi JSON→relasional
mengurai 5 list ini.
Uang (`*LIMIT*`, `*SUMINSURED*`, `*SUMREASURED*`, `MaxInsured`) = decimal (**ADR-0003**).

⚠️ **Typo di data:** `POLICYHODER` (kurang huruf L) — nama field existing; di skema baru pakai
`POLICYHOLDER` yang benar (penyimpangan sadar D6, sekeluarga dengan `PoductName`).

## OQ-JSON-PRODUCT → ✅ DITUTUP `[data DBA]` (contoh data diterima 2026-09-15).
Semua OQ modul ini tertutup kecuali OQ kecil (kepanjangan RI Rate/Risk) & OQ-047 (alamat storage).


---

## ⚠️ KRITIS: field null TIDAK muncul di JSON → daftar dari contoh ≠ lengkap `[keputusan work owner]`

Work owner: *"jika value null, property tidak muncul di JSON"*. Jadi contoh produk CAR (ID 100421)
hanya memperlihatkan field yang **terisi**. Daftar kolom lengkap diturunkan dari **Section input Pega
`SetProductNameInward`** (semua field yang bisa diisi, terlepas terisi/tidak) — `[terverifikasi]`.

### Field `ProductNameInward` (→ tabel `productinward_life`) — dari `SetProductNameInward` (**40 field**, sensus korpus `[terverifikasi]`)

`ADDENDUMNO`, `ADDENDUMWORD`, `AMANDEMENTNO`, `AMANDEMENTSCHD`, `BEGIN`, `BIRTHDAY`, `BROKERAGE`,
`CEDING`, `CEDINGLIMIT`, `CEDINGLIMITXPN`, `CEDINGRETENTIONNUM`, `CEDINGRETENTIONPCT`, `CURRENCY`,
`EXTRAMORTALITY`, `EXTRAPREMI`, `ID`, `INSURED`, `INWARDTREATYNM`, `LIENCLAUSE`, `MATURE`, `MAXAGE`,
`MAXCONTRACT`, `MAXDATARECEIVE`, `MAXEXPIREDCLAIM`, `MAXSUMINSURED`, `MAXSUMREASURED`, `MINAGE`,
`MINSUMINSURED`, `MONTHS`, `PAYMENT`, `POLICYHODER`(sic), `POLICYHODERNAME`, `PRODUCTID`,
`PROPORTIONALTABLE`, `RNMLIMITNUM`, `RNMLIMITPCT`, `RNMSHARE`, `STNC`, `SUBJECTTO`, `TREATYNUMBER`.

> **Sensus korpus `[terverifikasi]` (2026-09-15): TOTAL = 40.** Koreksi catatan awal yang keliru
> menulis "28"/37. Tiga yang sempat terlewat: `ADDENDUMNO`, `AMANDEMENTNO`, `CEDING` — **ketiganya
> di-set di `SetProductNameInward`, jadi AMBIL** (sesuai aturan pemilihan kolom). Pasangan
> NO/WORD & NO/SCHD: `ADDENDUMNO`/`ADDENDUMWORD`, `AMANDEMENTNO`/`AMANDEMENTSCHD` = field terpisah.

> Field yang **absen** di contoh CAR tapi ada di Section (contoh: `ADDENDUMWORD`, `AMANDEMENTSCHD`,
> `BIRTHDAY`, `CEDINGLIMITXPN`, `CEDINGRETENTIONPCT`, `EXTRAMORTALITY`, `EXTRAPREMI`, `INWARDTREATYNM`,
> `MONTHS`, `PROPORTIONALTABLE`, `RNMLIMITPCT`, `BROKERAGE`) — persis kasus "null tak muncul di JSON".

Pasangan NUM/PCT (`CEDINGRETENTIONNUM`/`PCT`, `RNMLIMITNUM`/`PCT`) & XPN (`CEDINGLIMITXPN`) =
varian nilai; uang vs persen. Tipe pasti (uang→NUMBER ADR-0003, tanggal→DATE) tetap perlu konfirmasi
per-field, tapi **daftar nama kolom kini lengkap dari Section**.

### Field `ProductName` (→ tabel `product_life`) + list bersarang
Skalar (dari contoh + Section): `ID`, `PRODUCTNAME`, `INWARDNAME`, `TREATYNUMBER`, `CEDING`/`CEDINGID`,
`SOBID`/`SOBNAME`, `POLICYHODER`/`POLICYHODERNAME`, `CAUSE`/`CAUSEID`, `RIRISK`/`RIRISKID`, `RICOMM`,
`IsORS`(buang), `IsView`, `CREATEOP`/`UPDATEOP`, `Comment`.
List bersarang → tabel anak: `CommentList`, `DocumentClaim`, `PlanList`, `UnderwritingLimitList`,
`FinancialUnderwritingList`, `LienClause`, `OutwardList`.

### Aturan PEMILIHAN kolom `[keputusan work owner]`
1. **Field yang di-SET di Activity simpan (`SetProductNameInward` / setara `ProductName`) = AMBIL
   SEMUA jadi kolom** — ia benar-benar diisi saat simpan = data nyata. **Ini sumber kebenaran daftar
   kolom** (bukan visibilitas Section).
2. **Field yang HANYA muncul di Section dengan visibilitas `1=2` / `1==2` / `never` dan TIDAK di-set
   di Activity = JANGAN AMBIL** — tampilan mati, bukan data.
3. 28 field inward di atas **semuanya di-set di `SetProductNameInward`** → **semua diambil**.
   Visibilitas mati Section tidak mengurangi daftar ini.

### Aturan migrasi mengikat `[keputusan work owner]`
1. **Field absen di JSON = NULL** (karena null tak muncul di JSON). Migrasi
   `JSON_VALUE(jsondata,'$.FIELD')` → NULL bila absen.
2. **Semua kolom hasil migrasi NULLABLE** — tak boleh `NOT NULL` berdasar asumsi "selalu ada".
   Wajib-isi (jika ada) ditegakkan di Go, bukan constraint DB.
3. **Daftar kolom = field yang di-set di Activity** (bukan superset visibilitas Section).

## OQ-JSON-PRODUCT → ✅ DITUTUP `[data DBA + terverifikasi Section]` (2026-09-15).
