---
status: accepted
tanggal: 2026-09-14
sumber: grilling Ronde 1 Q2 (prinsip) + Ronde 2 Q10 (klasifikasi kolom) — `.scratch/claim-life/`, keputusan work owner
---

# Uang di Claim — Life tidak direpresentasikan sebagai `float`

Delapan kolom nilai pada tabel akseptasi klaim Life adalah **uang** dan **tidak** boleh
direpresentasikan sebagai `float` di sistem baru. Satu kolom adalah **persen**, bukan uang.
Representasi uang yang dipakai harus mempertahankan presisi desimal secara eksak
(mis. tipe desimal berskala tetap, atau bilangan bulat dalam satuan minor) — **bukan** biner
floating-point.

## Klasifikasi kolom

`[terverifikasi work owner 2026-09-14]` Kolom pada `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`:

| Kolom | Sifat |
| --- | --- |
| `SUM_INSURED` | **uang** |
| `CEDING_RETENTION` | **uang** |
| `SUM_REASURED` | **uang** |
| `SHARE_NUSANTARA_RE` | **uang** |
| `CLAIM_AMOUNT` | **uang** |
| `SHARE_RETRO` | **uang** |
| `CLAIM_RETRO` | **uang** |
| `RETROCEDED_SHARE` | **uang** |
| `EM_PERCENT` | **persen** |
| `CURRENCY` | kode mata uang |

`[terverifikasi]` Nama kolom **tidak dapat dipakai untuk menebak sifatnya**. Empat kolom bernama
`SHARE_*` / `*_SHARE` ternyata **uang**, bukan rasio — padahal namanya menyarankan sebaliknya.
Korpus ini sudah terbukti memakai nama yang menipu (lihat `STS_REJECT` di `CONTEXT.md`).
Klasifikasi di atas berasal dari work owner, **bukan dari penamaan**.

## Bukti sumber kolom

`[terverifikasi]` Tabel ditulis oleh
`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL`
(`Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml`) — blok PL/SQL `INSERT INTO … COMMIT`,
sehingga **55 nama kolomnya terbaca langsung**, bukan tersembunyi di stored procedure.

Perintah audit:
```
awk '/<pyBrowseSQL>/,/<\/pyBrowseSQL>/' "Komite Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml"
```

## Consequences

- Aritmetika atas delapan kolom uang wajib memakai tipe eksak di seluruh lapisan —
  `handlers` → `services` → `repository` — dan di kontrak API.
- `EM_PERCENT` diperlakukan terpisah dari uang; mencampurnya dalam tipe yang sama akan
  menyembunyikan perbedaan makna.
- Claim — Life **sadar mata uang**: ada kolom `CURRENCY` pada rekam akseptasi. Ini berbeda dari
  nilai ter-hardcode di konteks lain yang **tidak menyebut mata uang** sama sekali
  (`Claim Non Prop`: `LimitMax = 30000000.00`; `Komite Claim FacIn`: pita
  `> 30000000.00 && <= 57750000.00`; `NB FacIn`: pangsa `0.45`/`0.05`, ambang `"3000000000"`
  dibandingkan sebagai **string**) → OQ-046.
- Nilai uang menyeberang ke Komite lewat kontrak **ADR-0001** (`CLAIM_AMOUNT` + `CURRENCY`
  ditambahkan atas keputusan Ronde 2 Q15) — representasinya harus konsisten di kedua sisi batas.


## Penguatan dari DDL Oracle (2026-09-14) `[data DBA]`

OQ-001 ditutup untuk Claim — Life; DDL `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` kini diketahui.
Temuannya **mengunci** keputusan non-float menjadi syarat yang lebih tajam:

> **Kedelapan kolom uang bertipe Oracle `NUMBER` — tanpa presisi dan tanpa skala.**

`NUMBER` tanpa presisi berarti Oracle menyimpan desimal **presisi arbitrer** (hingga 38 digit
signifikan) apa adanya. Konsekuensinya mengikat:

- Go **wajib** memakai tipe desimal **presisi arbitrer**. **`float64` dilarang** — ia hanya punya
  ~15–17 digit signifikan, sehingga nilai yang sah di Oracle dapat **berubah diam-diam** saat
  dibaca.
- Larangan ini berlaku di **seluruh lapisan** dan di **kontrak API**, termasuk saat nilai menyeberang
  ke Komite (**ADR-0001**) dan saat data dipindahkan (**ADR-0009**).
- Karena kolomnya tanpa skala, **tidak ada pembulatan yang boleh diasumsikan**. Pembulatan apa pun
  harus keputusan eksplisit, bukan efek samping tipe.

Kolom penyerta: `CURRENCY VARCHAR2(100)`.

**Invariant mata uang** `[keputusan work owner]` — dari penutupan **OQ-060**: bentuk nilai uang
adalah `(amount, currency)` **per baris**, dengan syarat **semua baris satu klaim wajib bermata uang
sama**. `[terverifikasi]` Konsisten dengan `Claim Life/Activity/SetIndexAdjustmentList.xml`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `SETINDEXADJUSTMENTLIST` / `RULE-OBJ-ACTIVITY`), yang
menyalin `CURRENCY` dan `CURRENCYID` dari `AdjustmentList(1)` ke `AdjustmentList(<LAST>)`.

⚠️ `[data DBA]` Roster komite `POOLDATA.EMAILKOMITE` **tidak punya kolom mata uang**, sehingga pita
`LIMIT_BOTTOM`/`LIMIT_TOP` berlaku atas **satu mata uang implisit**. Selama invariant di atas
dipegang, pembandingan pita tetap sahih; bila kelak invariant itu dilonggarkan, pembandingan pita
menjadi tidak sahih.

## OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-001** (modul lain) | DDL di luar persistensi Claim — Life belum diserahkan; tipe kolom uang konteks lain belum diketahui |
| **OQ-060** | **TERTUTUP untuk Claim — Life** (2026-09-14) — invariant satu klaim satu mata uang. Cakupan konteks lain tetap terbuka |

`[terjawab 2026-09-14]` OQ-060 kini **tertutup untuk Claim — Life**: bentuk tipe uang adalah
`(amount, currency)` **per baris**, dengan invariant **semua baris satu klaim bermata uang sama**.
