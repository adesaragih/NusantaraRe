# 06: Validasi Date of Loss per `Type` + `ContentNote` dari `BusinessCode`

**Status:** ready-for-agent

**Blocked by:** 03 (baris `AdjustmentList` + Save ke Outstanding)

## Hasil & nilai pengguna

Sebagai **ReasLifeAdmin**, sistem menolak *Date of Loss* yang berada di luar jendela valuasi polis —
sehingga klaim yang tidak tertanggung tidak pernah masuk ke siklus — dan jenis klaim terisi otomatis
dari kode produk sehingga saya tidak salah memilih. *(User story 4 dan 5 di spec)*

## Area codebase

`internal/services` (aturan validasi + penurunan jenis klaim), `internal/models` (data acuan
`BusinessCode` → `ContentNote`), `internal/handlers` (pesan kesalahan yang dapat dibaca),
`frontend/` (penampilan pesan pada formulir).

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Activity/ValidasiDOL_Act.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `VALIDASIDOL_ACT` / `RULE-OBJ-ACTIVITY`, 59.747 byte | dua cabang jendela valuasi per `Type` |
| `Claim Life/Activity/SaveOutStandingLife_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY` | penurunan `ContentNote` dari `BusinessCode` |

`[terverifikasi]` Dua cabang validasi DOL:

| Cabang | Jendela | Argumen `@addCalendar` |
| --- | --- | --- |
| `Type=="QR" \|\| Type=="QP"` | `GROSS_VALUATION_BEGIN_DATE` / `_EXPIRED_DATE` | `(.DATE_OF_LOSS,0,0,0,0,0,0,0)` |
| `Type=="TR" \|\| Type=="TP"` | `RETROCESSION_VALUATION_BEGIN_DATE` / `_EXPIRED_DATE` | `(.DATE_OF_LOSS,0,0,0,1,0,0,0)` |

Gagal → `local.errmsg = "Invalid DOL"`, dengan `local.Begin==false || local.Expired==true`.

## ADR terkait

**ADR-0012** (`Type` menggerbangi dua hal: wewenang **dan** jendela validasi; `TP` = Payable,
`TR` = Receivable — `[keputusan work owner]`, tidak ada di korpus).

## Acceptance criteria

- [ ] Klaim ber-`Type` `QP`/`QR` dengan *Date of Loss* di luar `GROSS_VALUATION_BEGIN_DATE` …
      `_EXPIRED_DATE` ditolak dengan pesan yang setara `"Invalid DOL"`. *(AC 13 spec)*
- [ ] Klaim ber-`Type` `TP`/`TR` diuji terhadap `RETROCESSION_VALUATION_*`, dengan pergeseran
      tanggal yang sama seperti Pega. *(AC 14 spec)*
- [ ] `ContentNote` terisi sesuai tabel `BusinessCode` `L1`–`L21` di `CONTEXT.md`. *(AC 15 spec)*
- [ ] Pemetaan `BusinessCode` → `ContentNote` diwujudkan sebagai **data acuan**, bukan rangkaian
      `if` bercabang.
- [ ] Validasi dan penurunan jenis klaim membaca **satu** field `Type` yang sama — bukan dua salinan
      seperti di Pega.

## Catatan penutupan (2026-09-14)

**DOL dibawa sebagai paritas** `[keputusan work owner]` — validasi *Date of Loss* per `Type`
dipertahankan persis seperti perilaku Pega. Tercatat di `CONTEXT.md` (Lampiran, butir 2).

`[terbuka — catatan kecil, NON-PEMBLOKIR]` Satuan pergeseran tanggal pada cabang `TP`/`TR`
(`@addCalendar(.DATE_OF_LOSS,0,0,0,1,0,0,0)`) belum dikonfirmasi **Product+UW**; definisi
`@addCalendar` tidak ada di korpus. **Ini bukan blocker**: pergeseran direplikasi apa adanya dari
argumen yang terbaca. Bila kelak Product+UW menyatakan satuannya berbeda, itu perubahan satu baris
disertai testnya.

`[terbuka]` **OQ-020** — arti `QP`/`QR` belum dijawab. **Tidak memblokir**: perilaku kedua cabang
sudah terbaca penuh; hanya namanya yang belum.

`[data DBA]` Kolom terkait di `POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`: `TYPE VARCHAR2(10)`.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
