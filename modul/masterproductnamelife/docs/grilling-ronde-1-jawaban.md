# Grilling — Master Product Name Life — Ronde 1: JAWABAN FRONTIER (final)

Tanggal: 2026-09-15
Konteks: `master-product-name-life` — **konteks/menu sendiri** `[keputusan work owner]`
Modul: **Master Product Name Life** (114 berkas)
Sumber jawaban: sesi grilling dengan work owner (Ronde 1).

> Konvensi: `[terverifikasi]` = terbukti korpus (class + nama + path); `[keputusan work owner]`;
> `[fakta bisnis — work owner]`; `[data DBA]`; `[terbuka]` = OQ.

**Frontier kosong — 10 pertanyaan terjawab.** Dokumen ini = sumber kebenaran untuk `/to-spec`.

---

## Ringkasan: modul ini JAUH lebih ramping dari dugaan awal

Empat sumber kompleksitas terbesar ternyata **dead code / salah ekspor** (dikonfirmasi work owner):

| Dead code | Bukti | OQ |
| --- | --- | --- |
| `SetTreatyIn_Act` + seluruh jalur treaty-in (`SaveTreatyIn`, `BrowseTreatyIn`, `revisionstate=1`, `PEGA_TREATY_IN`, `M_TREATY_IN`/`M_TREATY_IN_EDM`) | **salah ekspor XML** `[keputusan work owner]` | OQ-056 ditutup |
| `GetCountClaim` + `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` | **tidak dipakai** `[keputusan work owner]` | OQ-058 ditutup |
| `GetReinsTypeOR_Life` + `BrowseReinstypeOR_SQL` (jalur OR/10200) | **mati** — `[terverifikasi]` container di `InboxProductName.xml` bergerbang `pyContainerVisibleWhen>1==2` (selalu false) | OQ-057 ditutup |
| `SaveProductNameLIfeFlat` (UPDATE parsial RIRISKID/RIRISK) | **sudah tidak dipakai** `[keputusan work owner]` | — |

⚠️ Pelajaran OQ-066 berlaku berkali-kali di sini: **file ikut terekspor / rule ada ≠ dipakai.**
Semua dikonfirmasi work owner atau terbukti gerbang `1=2`.

**Sisa nyata: editor master produk life** — satu produk disimpan ke dua tabel, dengan lampiran ke
Google Storage.

---

## Temuan Ronde 1 yang DIVERIFIKASI Kiro ke korpus (LULUS)

1. ✅ **Simpan produk lewat JSON CLOB (di Pega)** — 2 procedure penulis terima 1 CLOB `DATAPEGA` +
   `IDPEGA`, keluaran **tiga** (`ERRMSG`, `IDPEGAOUT`, `STSSAVE`).
2. ✅ **`M_PRODUCT_LIFE` hibrida** — `SaveProductNameLIfeFlat`
   (`ASM-FW-GISFW-INT-PRODUCT_LIFE!ASM!SAVEPRODUCTNAMELIFEFLAT`, `RULE-CONNECT-SQL`):
   `UPDATE POOLDATA.M_PRODUCT_LIFE SET RIRISKID={CARI2}, RIRISK={CARI3} WHERE ID={CARI1}` —
   membuktikan tabel punya kolom relasional nyata, bukan murni JSON. (Rule ini kini dinyatakan mati.)
3. ✅ **OQ-057** — `BrowseReinstypeOR_SQL` memuat `REINSTYPEID = '10200'` **literal** + join
   `treatycontract_life × treatyyear_life`. OR = 10200 (enumerasi Retro Life). Pemakainya mati (`1=2`).
4. ✅ `SetTreatyIn_Act` punya 5 langkah `param.revisionstate==1` (memo `add RevisionState`,
   2025-11) — struktur nyata, tetapi **salah ekspor** (bukan fungsi modul ini).
5. ✅ Gerbang `1=2`/`1==2` tersebar di `InboxProductName.xml` termasuk `pyContainerVisibleWhen>1==2`
   yang membungkus pemanggil `GetReinsTypeOR_Life`.

Cacat laten (tak diminta): `PoductName` (typo, 20×, hanya di 2 activity simpan) vs `ProductName`
(714×); `IsORS` dibandingkan `=="true"`/`=="'true'"` tetapi **nol tempat men-set**.

---

## FAKTA BISNIS (work owner — tidak ditebak)

### Q1 — `M_PRODUCT_LIFE` vs `M_PRODUCTINWARD_LIFE` `[fakta bisnis — work owner]`
**SATU produk yang sama.** Datanya disimpan ke **dua tabel**, terhubung lewat **`ID` sama**. Bukan
dua produk berbeda. → satu operasi simpan menulis dua tabel.

### Q1b — GABUNG dua tabel induk `[keputusan work owner]`
Existing: `M_PRODUCT_LIFE` + `M_PRODUCTINWARD_LIFE` = **1:1 by `ID`**. **Keputusan: digabung jadi
SATU tabel `product_life`** di skema baru. Alasan: JSON dibuang (D2) + atomik (D1) → pemisahan tak
beralasan; gabung menghapus kelas bug "produk timpang antar dua induk" dan menyederhanakan atomik
(satu baris). Field yang muncul di kedua sisi (`ID`, `CEDING`, `POLICYHOLDER`, `POLICYHODERNAME`,
`BIRTHDAY`, `TREATYNUMBER`) → satu kolom. **Total tabel skema baru = 6** (1 induk + 5 anak).
⚠️ **Penyimpangan sadar (tambahan).**

### Q6 — 7 pemilih master `[fakta bisnis — work owner]`
| Pemilih | Arti |
| --- | --- |
| `ChooseSOB` | **SOB = Source of Business** |
| `ChooseRIRate` | master **RI Rate** (master tersendiri; kepanjangan `[terbuka]`) |
| `ChooseRIRisk` | master **RI Risk** — **berbeda** dari RI Rate (kepanjangan `[terbuka]`); mengisi `RIRISKID`/`RIRISK` |
| `ChooseCauseOfLoss` | penyebab klaim/rugi |
| `ChooseCeding` | perusahaan ceding |
| `ChooseCurrency` | mata uang |
| `ChoosePolicyHolder` | pemegang polis |

⚠️ RI Rate ≠ RI Risk — dua master berbeda, jangan digabung.

### Q8 — kontrak OR (10200) `[keputusan work owner + terverifikasi]`
**MATI.** `GetReinsTypeOR_Life`/`BrowseReinstypeOR_SQL` dimatikan lewat gerbang `1=2` di
`InboxProductName`. Tidak masuk cakupan modul ini. (OR = 10200 tetap benar sebagai enumerasi, tapi
tak dipakai di sini.)

### Q9 / OQ-056 — jalur `M_TREATY_IN` `[keputusan work owner]`
**SALAH EKSPOR, tidak dipakai.** Modul ini **tidak** menulis master treaty inward. Tidak ada efek
tulis lintas-domain.

---

## KEPUTUSAN DESAIN (rekomendasi Kiro, disetujui work owner)

### D1 — Simpan 1 produk → 2 tabel: transaksi atomik (opsi a) `[keputusan work owner]`
`[data DBA]` Procedure `PEGA_M_PRODUCT_LIFE` & `PEGA_M_PRODUCT_INWARD_LIFE` **COMMIT sendiri** →
tidak atomik di Pega. **Keputusan (a):** di Go **tulis INSERT/UPDATE sendiri dalam SATU transaksi**;
rollback → **keduanya rollback**. **Procedure JSON TIDAK dipanggil/dimigrasikan.**
> ⚠️ **Penyimpangan sadar besar:** beda dari pola "panggil procedure apa adanya" — karena procedure
> commit sendiri & menghalangi atomik yang diminta work owner.

### D2 — Simpan FLAT/relasional penuh (opsi a) `[keputusan work owner]`
`[data DBA]` DB existing menyimpan **JSON** (`M_PRODUCTINWARD_LIFE` = hanya `ID`+`JSONDATA`;
`M_PRODUCT_LIFE` = `JSONDATA` + 4 kolom flatten `RIRISKID`/`RIRISK`/`PRODUCTNAME`/`BEGIN_DATE`;
constraint `JSONDATA IS JSON`). **Keputusan (a):** skema baru **relasional penuh** — semua atribut
(MINAGE, MAXAGE, BROKERAGE, EXTRAPREMI, INSURED, MAXSUMREASURED, MAXEXPIREDCLAIM, MAXDATARECEIVE,
SUBJECTTO, RNMLIMITNUM, RIRISKID/RIRISK, PRODUCTNAME, BEGIN_DATE, dll) jadi **kolom bernama**;
**migrasi data JSON→kolom**.
> ⚠️ **Penyimpangan sadar besar:** menyimpang dari DB existing (JSON). D1+D2 sejalan — karena tak
> pakai procedure JSON, atomik jadi satu transaksi INSERT ke dua tabel relasional.
> ⚠️ **OQ BARU (OQ-JSON-PRODUCT):** **daftar field di dalam `JSONDATA`** belum terbaca dari korpus/DDL
> — dibutuhkan untuk mendesain kolom. **Butuh DBA/Product** (contoh JSON produk nyata) atau baca
> Section input `InboxProductName`/`SetProductNameInward` untuk menurunkan nama field. Memblokir
> desain skema final (tiket migrasi), **tidak** memblokir bentuk spec.

### D3 — `SaveProductNameLIfeFlat` `[keputusan work owner]`
**Tidak dipakai** — dead code, tidak dimigrasikan. Update RI Risk jadi bagian update produk biasa.

### D4 — Google Storage + lampiran + token `[keputusan work owner]`
Ikuti rekomendasi: endpoint di-resolve dari config/`M_LINK_SERVICE` (**ADR-0013**, bukan hardcode);
upload lampiran = **efek keluar**, gagal jangan diam-diam (**ADR-0015**). **Metadata produk & lampiran
DIPISAH** — produk tersimpan dulu, lampiran menyusul dengan status jelas. **Lampiran OPSIONAL** —
produk boleh tersimpan tanpa lampiran.

### D5 — Hasil simpan `[keputusan work owner]` + `[data DBA]`
`[data DBA]` Kontrak keluaran procedure lama: `StsSave` **100=sukses / 99=gagal**, `ErrMsg` = teks
(juga saat sukses), `IDPegaOut` = ID final. **Karena D1/D2 membuang procedure**, Go **tidak** memakai
3 keluaran itu — Go hasilkan **sukses/gagal + ID sendiri**, **gagal terang-terangan**. Format ID
**tetap ditiru:** `'1' + lpad(sequence, 5, '0')` (5 digit, mis. `100001`) via sequence
`M_PRODUCT_LIFE_SEQ` / `M_PRODUCT_INWARD_LIFE_SEQ`. Konsisten **ADR-0006**.

### D6 — Cacat laten `[keputusan work owner]`
**Jangan tiru typo** `PoductName` → pakai `ProductName` konsisten. **Buang `IsORS`** (cabang mati;
jalur OR sudah mati). Catat sebagai kejanggalan Pega, bukan fitur.

---

## Penyimpangan sadar (rekap) — jadi AC ⚠️ di spec

| # | Penyimpangan |
| --- | --- |
| 1 | Simpan 1 produk → 2 tabel dalam **satu transaksi atomik** (rollback keduanya); **tidak** pakai procedure lama (yang commit sendiri) |
| 2 | Simpan **flat/relasional penuh**, bukan JSON CLOB; migrasi JSON→kolom |
| 3 | Hasil simpan: gagal terang-terangan; ID `'1'+lpad(seq,5)` ditiru via sequence |
| 4 | Lampiran = efek keluar (ADR-0015), endpoint dari config (ADR-0013), **lampiran opsional** & terpisah dari metadata |
| 5 | Buang typo `PoductName` + cabang mati `IsORS` |

## Dead code (JANGAN dimigrasikan)

`SetTreatyIn_Act`+jalur treaty-in · `GetCountClaim`/`PC_ASM_FW_GCNMFW_WORK` ·
`GetReinsTypeOR_Life`/`BrowseReinstypeOR_SQL`/jalur OR (gerbang `1=2`) · `SaveProductNameLIfeFlat`.

## Status OQ

| OQ | Status |
| --- | --- |
| OQ-056 | ✅ DITUTUP `[keputusan work owner]` — SetTreatyIn salah ekspor |
| OQ-057 | ✅ DITUTUP — OR=10200; pemakai di modul ini mati (`1=2`) |
| OQ-058 | ✅ DITUTUP `[keputusan work owner]` — GetCountClaim tidak dipakai |
| OQ-002 | ✅ DITUTUP `[data DBA]` — body 3 procedure diterima (`dba-procedures-and-ddl.md`): keduanya upsert JSON, **COMMIT sendiri**, `StsSave` 100/99, ID `'1'+lpad(seq,5)`; `GET_TOKEN_STORAGE` (MD5, cache `GCP_IMAGE`, 1 menit) |
| OQ-001 (modul) | ✅ DITUTUP `[data DBA]` — DDL diterima: `M_PRODUCT_LIFE` (`ID`, `JSONDATA`+`JSON` constraint, kolom flatten `RIRISKID`/`RIRISK`/`PRODUCTNAME`/`BEGIN_DATE`), `M_PRODUCTINWARD_LIFE` (`ID`+`JSONDATA`). Tanpa PK (hanya index). Menegaskan D2 = redesign relasional |
| **OQ-JSON-PRODUCT** | ✅ DITUTUP `[data DBA]` — contoh JSON produk nyata (ID 100421) diterima. Field skalar + list bersarang (CommentList, DocumentClaim, PlanList, UnderwritingLimitList, dll) terurai. Lihat `dba-procedures-and-ddl.md`. Konsekuensi D2: relasional penuh = induk 2 tabel + ~7 tabel anak |
| OQ-047 | terbuka — alamat Google Storage dari `M_LINK_SERVICE`; ADR-0013 berlaku |
| OQ kecil | terbuka tak-memblokir — kepanjangan RI Rate/RI Risk |

## Kesiapan `/to-spec` — SEMUA PEMBLOKIR DITUTUP
- Frontier kosong; 4 fakta bisnis + 6 keputusan desain selesai; 4 jalur dead code dikonfirmasi.
- OQ-002, OQ-001, **OQ-JSON-PRODUCT** ✅ ditutup (body + DDL + contoh JSON diterima).
- **Struktur relasional (D2) terurai:** induk `product_life` + `productinward_life` (by ID) +
  **5 tabel anak** (comment, document_claim, plan, uw_limit, fin_uw). `LienClause` = field skalar
  (bukan tabel); `OutwardList` = dead code jalur OR. `[terverifikasi]`
- Spec + tiket (termasuk migrasi) dapat **ready-for-agent**.
- ⚠️ Penyimpangan sadar besar tercatat: D1 atomik (Go tulis sendiri, procedure commit-sendiri tak
  dipakai), D2 relasional penuh (migrasi JSON→kolom), typo `POLICYHODER`/`PoductName` dibuang.
- OQ tersisa tak-memblokir: OQ kecil (kepanjangan RI Rate/Risk), OQ-047 (alamat Google Storage).
