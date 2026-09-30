# Keputusan Struktur Endorsement Life (EDM) — berbagi tabel PremiumList Life

Tanggal: 2026-09-16
Sumber: keputusan work owner + peta discovery Endorsement Life + verifikasi korpus (MappingEDMLife).
Status: `[keputusan work owner]` — konteks Endorsement Life memakai tabel PremiumList Life yang sama.

> EDM Life = **endorsement/perubahan polis life**. Berbagi seluruh struktur 7 tabel PremiumList Life
> (`revisi-penyimpanan-premiumlist.md`). Bukan tabel baru. Yang beda: versi + pencocokan old↔new.

---

## Keputusan inti `[keputusan work owner]`

1. **OldData TIDAK disimpan** — dibaca dari **polis versi sebelumnya** (PRODKE lebih kecil / PARENT_ID).
   `Addendum.OldData` di Pega = page kerja sementara (tampil before/after), tak di-persist.
2. **Struktur = 7 tabel PremiumList Life yang SAMA** (opsi A). Endorsement = **versi/baris baru** di
   `T_PREMIUM_LIST`, bukan tabel/entitas tersendiri. Semua versi hidup berdampingan.
3. **Pembeda versi/endorsement = kolom di `T_PREMIUM_LIST`** `[keputusan work owner: opsi A-revisi]`.
   EDM TIDAK punya tabel header sendiri — tambah kolom EDM (nullable, kosong untuk NB, terisi saat
   endorsement) ke `T_PREMIUM_LIST`. Class EDM `ASM-FW-GISFW-Work-EndorsementLife` (beda dari
   `Work-LIFE`) tetap dipetakan ke tabel yang sama.

   **Kolom EDM `[terverifikasi dari pyFields DATA_JSON EDM nyata]`:**
   | Kolom | Sumber property | Arti |
   | --- | --- | --- |
   | `EDM_TYPE` | `EdmType` | **1=Perubahan Data, 3=Batal** — dropdown korpus (bukan tebakan) |
   | `OLD_POLICY_NO` | `OldPolicyNo` | no polis lama yang di-endorse |
   | `EDM_DATE` DATE | `EdmDate` | tanggal endorsement |
   | `EDM_NOTE` | `EDMNote` | catatan endorsement |
   | `TYPE_CEDING` | `TypeCeding` | **1=QS, 2=SURPLUS, 3=QS+SURPLUS, 4=XOL** — dropdown korpus |
   | `PREMI_PROPOSED` NUMBER | `PremiProposed` | premi diusulkan |
   | `UANG_PERTANGGUNGAN` NUMBER | `UangPertanggungan` | |
   | `SUM_INSURED` (header) | `SumInsured` | |
   | `JENIS_PRODUK` | `JenisProduk` | |
   | `SISTEM_REASURANSI` | `SistemReasuransi` | |
   | `STATUSS`, `STATUS_UPDATE`, `STATUS_SERVICE` | `Statuss`/`StatusUpdate`/`StatusService` | status proses EDM |
   | `START_DATE`, `END_DATE` DATE | `StartDate`/`EndDate` | |
   Ditambah pembeda versi umum:
   - Work ID **`EDMLF-<n>`** `[fakta bisnis — work owner]` (paritas `NBLF-`; nol match korpus)
   - `NOENDORS` `[terverifikasi]` (`Generate_NoEndorsmentLife`)
   - `PL_NUMBER_EDM`, `PRODKE` (versi berjalan = PRODKE terbesar `[terverifikasi]`)
   - `EDMSTATUS` (Old/New/Delete/Batal `[terverifikasi]` set di `MappingEDMLife`), `STATUSOLD`
   > ⚠️ Kolom EDM nullable — kosong untuk baris new business (NB), terisi untuk endorsement.

## Pencocokan peserta old↔new — MEKANISME COPY + PARENT_ID `[keputusan work owner]`

**Alur (dari work owner):**
```
1. User pilih NOPOLIS → Submit
2. Sistem COPY premium list versi berjalan → versi baru:
   - header + PESERTA (detail) + spreading + retro di-COPY, tiap baris ID BARU + PARENT_ID = ID asal
   - SUMMARY (currency) TIDAK di-copy — dihitung ULANG dari peserta versi baru
3. User ubah data peserta di versi baru
4. Selisih new − old = nilai(ID) − nilai(PARENT_ID) per baris peserta
5. Summary currency versi baru = rekap dari peserta versi baru (AppendCurrencySummary_DT)
```

**Apa yang di-copy vs tidak `[keputusan work owner]`:**
| Tabel | Copy saat endorsement? |
| --- | --- |
| `T_PREMIUM_LIST` (header) | ✅ copy — versi baru + kolom EDM |
| `T_PREMIUM_LIST_DETAIL` (peserta) | ✅ copy + PARENT_ID |
| `T_PREMIUM_LIST_SPREADING` / `_SPREADING_RETRO` | ✅ copy (ikut peserta) — **TANPA PARENT_ID**, cukup FK ke peserta/spreading versi baru |
| `T_PREMIUM_LIST_SUMMARY` (currency/rekap) | ❌ **TIDAK** — rekap turunan, dihitung ulang dari peserta versi baru |
| `T_VIEW_SUGGEST` | ❌ **TIDAK** `[keputusan work owner]` — riwayat penawaran milik proses awal; endorsement punya jejaknya sendiri (bila ada) |

**Kenapa PARENT_ID (bukan kunci bisnis / indeks):**
- Pointer eksplisit → **tidak pernah salah pasang peserta** (kekhawatiran work owner).
- Tak bergantung CERTIFICATE_NO (bisa berubah) maupun urutan baris (bergeser saat tambah/hapus).
- ⚠️ **Penyimpangan sadar dari Pega:** sistem lama pakai loop indeks + `EDMStatus` di `MappingEDMLife`
  (`PremiumListDetail(Local.idxLocationEDM)` vs `(Param.EDMList=.pxListSubscript)`) — rawan salah
  pasang saat urutan bergeser, dan **kunci pencocokannya tak terbaca dari korpus**. Diganti pointer
  `PARENT_ID` eksplisit — lebih benar & auditable.

**Status peserta = turunan dari PARENT_ID + aksi:**
| Kondisi | EDMSTATUS | Pengurang |
| --- | --- | --- |
| hasil copy, nilai diubah | Change | selisih = new − old |
| peserta baru (setelah copy) | New | PARENT_ID NULL → tak ada pengurang |
| baris copy ditandai keluar | Delete | pengurang penuh (old) |

**Kolom BARU `PARENT_ID`** `[keputusan work owner]` — HANYA di `T_PREMIUM_LIST_DETAIL` (peserta),
FK self-reference ke `ID` peserta versi sebelumnya (NULL untuk NB / peserta baru). Inilah satu-satunya
tingkat yang di-selisih new−old, jadi hanya ia yang butuh pointer.
- `T_PREMIUM_LIST_SPREADING` / `_SPREADING_RETRO`: **TANPA PARENT_ID** — ikut ter-copy di bawah
  peserta versi baru, cukup FK `DETAIL_ID`/`SPREADING_ID` menunjuk induk versi baru.
- `T_PREMIUM_LIST_SUMMARY`: tak di-copy (dihitung ulang), tanpa PARENT_ID.
- `T_PREMIUM_LIST` (header): versi dibedakan `PRODKE` + kolom EDM; `PARENT_ID` header opsional untuk
  melacak versi induk (keputusan minor, boleh ada untuk jejak versi).

Query selisih (bersih):
```sql
SELECT n.SUM_INSURED - NVL(o.SUM_INSURED,0) AS selisih, ...
FROM   T_PREMIUM_LIST_DETAIL n
LEFT JOIN T_PREMIUM_LIST_DETAIL o ON o.ID = n.PARENT_ID
WHERE  n.PREMIUM_LIST_ID = :versiEndorsement
```

---

## Frontier EDM yang sudah tertutup
- Struktur penyimpanan: tabel PremiumList bersama, versi via baris + PARENT_ID ✅
- OldData: baca versi sebelumnya, tak disimpan ✅
- Pencocokan: copy + PARENT_ID (anti salah-pasang) ✅

## ⚠️ KOREKSI 2026-09-16 — label dropdown BUKAN bukti korpus
Klaim sebelumnya "arti kode dari dropdown pyFields korpus" **KELIRU**. Nilai `pyLocalizedValue` di
`DATA_JSON`/pyFields adalah **isi payload runtime satu instance**, BUKAN rule korpus. Sensus:
`Receivable`/`Payable` = **NOL berkas** di `Endorsement Life/`. Status yang benar:
- **`EdmType` 1/3** (Perubahan Data / Batal): label field ada di `InputEDMLife.xml`/
  `EndorsmentLife_Section.xml`, tapi pemetaan ke nilai = `[keputusan work owner]`, bukan dropdown korpus.
- **`Type` QR/QP/TR/TP**: arti **TIDAK terbukti korpus** → **OQ-020 TETAP TERBUKA**. Yang terbukti
  cuma perilaku (`Q*`→gross, `T*`→retrosesi). JANGAN klaim Receivable/Payable sebagai fakta korpus.
- **`TypeCeding` / `ProRateType`**: label dari instance JSON — perlakukan sebagai `[dugaan]` sampai
  ada bukti rule, bukan `[terverifikasi]`.
- ⚠️ **JANGAN selaraskan penutupan OQ-020 ke Claim Life** — sudah dibatalkan di open-questions.md &
  CONTEXT.md.

## Alur EdmType 1 vs 3 `[keputusan work owner]`
Keduanya pakai jalur simpan SAMA (copy + PARENT_ID). Beda di nilai & flag yang disimpan:

**EdmType 1 (Perubahan Data):**
- Copy normal (header+detail+spreading+retro, PARENT_ID di detail).
- Peserta diubah bila perlu → selisih = new − old (via PARENT_ID).
- Peserta yang **dihapus**: BUKAN delete fisik → beri **flag `EDMSTATUS='Delete'`** + **nilai
  dijadikan MINUS** (pengurang). Baris tetap tersimpan (auditable).
- Peserta baru: `EDMSTATUS='New'`, PARENT_ID NULL.
- Peserta tak berubah: tetap (selisih 0).

**EdmType 3 (Batal):**
- **SEMUA peserta di-minus-kan** — seluruh polis dibatalkan, tiap peserta jadi pengurang penuh (−old).

⚠️ **Prinsip: hapus = flag + nilai minus, TIDAK PERNAH delete fisik.** Auditable (sama semangat
Claim Life). Ini penyimpangan/penegasan sadar dari sekadar "hapus baris".

## Frontier EDM — praktis KOSONG untuk struktur/penyimpanan
- ~~Ambang 50.000 baris~~ **DIBUKA/mengikuti PremiumList Life** (tiket 04 "Unggah CSV"). Bukan OQ baru.
- Buang JSON: `InsertJsonPolisEDM`/`convertJsonNusareToProduction` → hilir baca flat (sama PremiumList).
- Penomoran `NOENDORS` (`Generate_NoEndorsmentLife`) — format & kontrak (DBA, non-pemblokir).
- `EDMLF-` format & generator — `[fakta bisnis — work owner]` (non-pemblokir).
- `EDMSTATUS` Old/New/Delete/Batal — nilai penanda per baris (Old=baris versi lama, New=peserta baru,
  Delete=dikeluarkan, Batal=pembatalan). Dicatat apa adanya; efek ke selisih lewat PARENT_ID.

## Dampak artefak
- Tambah kolom `PARENT_ID` ke tiket skema PremiumList (tiket 00) ATAU tiket skema EDM sendiri.
- Kolom EDM (`PL_NUMBER_EDM`/`EDMSTATUS`/`STATUSOLD`/`NOENDORS`/`PRODKE`/`EDM_TYPE`/`EDMLF-`) yang tadi
  ditandai "menyusul saat Endorsement" di PremiumList → sekarang dipastikan masuk.
