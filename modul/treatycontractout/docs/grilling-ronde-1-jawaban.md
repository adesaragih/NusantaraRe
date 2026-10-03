# Grilling — Treaty Contract Out — Ronde 1: JAWABAN FRONTIER (final)

Tanggal: 2026-09-15
Konteks: `treaty-contract-out` — konteks non-Life pertama. **303 file.**
Sumber jawaban: sesi grilling dengan work owner + body/DDL DBA.

> Konvensi: `[terverifikasi]` = terbukti korpus (class+nama+path); `[keputusan work owner]`;
> `[fakta bisnis — work owner]`; `[data DBA]`; `[terbuka]` = OQ.

**Frontier kosong — 12 pertanyaan terjawab. OQ-001 & OQ-002 ditutup.** Sumber kebenaran untuk `/to-spec`.

---

## Ringkasan: apa modul ini sebenarnya

`[terverifikasi]` **Editor master term/arrangement kontrak treaty NON-LIFE** — bukan treaty outward
(nama menyesatkan, OQ-022 terbukti bukan outward). Editor master murni: tanpa tangga persetujuan,
CRUD per klausul. Modul ini **satu-satunya penulis** master arrangement; Claim Prop/Fac In/Komite =
**read-only** `[keputusan work owner]`.

Entitas: `TREATYYEAR` → `TREATYCONTRACT` → { `TREATYREINSURER` → `MTREATYSECURITY`, `TREATYBUSINESS`,
`PROPORTIONALARRG` (klausul) }. Master dibaca (bukan dimiliki): `REINSURANCETYPE` (jenis reasuransi),
`TREATYDESC` (jenis klausul), currency.

---

## Temuan Ronde 1 diverifikasi Kiro ke korpus (LULUS)

1. ✅ **Tabel kembar JSON vs relasional** — kueri baca dari `PROPORTIONALARRG` (kolom datar) DAN
   `m_PROPORTIONALARRG` (JSON `a.JSONDATA.…`, JSON_TABLE). Terbukti di `GetMasterDescriptionLimitParentList`
   (relasional) vs `GetMasterDescriptionPLAParentList`/`GetMasterPortfolioListDetail`/dll (JSON).
2. ✅ **25 klausul → 2 procedure** (`PEGA_PROPORTIONALARRG` 35 param, `PEGA_M_PROPORTIONALARRG_CHILD`
   26 param); dibedakan `TreatyDescID`. **`[data DBA]` keduanya menulis ke SATU tabel `PROPORTIONALARRG`.**
3. ✅ **Kaskade hapus** `DeleteFromTREATYCONTRACT_SQL`: hapus `treatycontract`→`treatybusiness`→
   `MTREATYSECURITY`→`TREATYREINSURER` (4 tabel), **`PROPORTIONALARRG` TIDAK ikut**.
4. ✅ **Rule `_Old` justru terbaru & dominan** — `BrowseReinsuranceType_RD_Old_Ljt_id_isnotnull`
   (ruleset 01-01-91, 2026-02) dipakai **11 grid**; yang non-Old dipakai 1 layar. (RD mandiri, bukan
   memanggil rule "Old" lain.)

Koreksi D2/D3: `SaveTreatyArr*_Act` = 25 (bukan 27); 16 `CancelActivity*` bukan daftar klausul lengkap.

---

## FAKTA BISNIS + KEPUTUSAN (work owner)

### Q1 — Tabel kembar JSON vs relasional `[keputusan work owner]`
**(a) Skema relasional penuh; JSON dibuang.** `[keputusan work owner]` **`M_PROPORTIONALARRG` (JSON)
sudah lama tidak dipakai; kueri yang masih baca dari situ = lupa dihapus (dead code).** Sumber
kebenaran = **`PROPORTIONALARRG` (relasional)**. `[data DBA]` terbukti: 2 procedure penulis hanya
INSERT/UPDATE ke `PROPORTIONALARRG`, tak pernah ke JSON. Semua kueri `FROM m_PROPORTIONALARRG` →
migrasi baca dari `PROPORTIONALARRG`. Migrasi TIDAK ambil dari JSON.
> ⚠️ Penyimpangan sadar: buang dualitas JSON.

### Q2 — Satu tabel generik untuk semua klausul `[keputusan work owner]`
**(a) SATU tabel `proportionalarrg`.** `[data DBA]` terbukti: procedure "CHILD" menulis ke tabel yang
**sama** (`PROPORTIONALARRG`), bedanya hanya jumlah kolom (child 26, induk 35 — tanpa 9 kolom induk:
`ID_OCCUPATION, OCCUPATION, ID_CLAUSE, CLAUSE, TREATYLIMIT, COINS_MIN, COINS_MAX, MORERP, MOREUSD`).
Baris "child" = 9 kolom itu NULL. Jenis klausul dibedakan `TREATYDESCID`. Validasi per klausul di Go (Q9).

### Q3 — Arti/kepanjangan istilah klausul `[keputusan work owner]`
**Ikuti nama & tampilan Pega apa adanya** — jangan diterjemahkan/ditebak. Field (`Ydcf, PctMe,
LayerPartType, SpreadingOrder, Method`, dll) pakai nama asli; label klausul (EPI, PLA, Ricomm, MB,
LOL, BordereAux, ExGratia, CashLossLimit, ClaimCoorperation, ProfitCommision, pTerrLimit, Portfolio,
CoinsPanel, ExclutionTreaty) ambil dari teks tampilan Pega. Arti bisnis = **OQ terbuka Product+UW**,
tak memblokir; spec pakai istilah asli.

### Q4 — Daftar jenis klausul: data atau kode `[fakta bisnis + keputusan work owner]`
**(a)** Daftar jenis = **master data `TREATYDESC`** (master TERPISAH, di-input di layar sendiri;
modul ini baca saja). `[data DBA]` DDL `TREATYDESC(ID, DESCNAME, ISXOL, STATUSAKTIF)`. **Aturan
validasi per jenis = kode Go** (petakan kode → aturan wajib-isi).

### Q5 — Blacklist 12 ID jenis reasuransi `[keputusan work owner]`
**Ikuti apa adanya (masih benar).** Filter 11 grid: `.ID NotStartsWith
("10004","10011","10012","10021","10022","10025","10026","10028","10248","10249","10018","10217")
AND .Flag="active" AND .Type IN ("1","2","3")`. Tiru persis; jangan digeneralkan.

### Q6 — Rule `_Old` yang terbaru `[keputusan work owner]`
Migrasikan perilaku RD yang dipakai **11 grid** (`BrowseReinsuranceType_RD_Old_Ljt_id_isnotnull`,
ruleset 01-01-91). **Beri nama jujur** di sistem baru (buang "Old"). RD non-Old = perilaku layar
master jenis reasuransi (Q12i). OQ-066 (nama bohong).

### Q7 — Fitur copy tahun treaty `[keputusan work owner]`
**DIBUANG — tidak dimigrasikan.** Alasan work owner: menyalin membawa nilai lama (mis. QS 2025 1M) ke
tahun yang seharusnya beda → sumber kesalahan. `PROSESCOPY` + pemanggil (`BrowseCopyData`,
`SaveTreatyYearMultiple_Act`, bagian copy `BrowseDeleteRowTreatyInContract`) = **dead code**.
PROSESCOPY keluar dari daftar OQ-002.

### Q8 — Kaskade hapus & klausul yatim `[keputusan work owner]`
`[fakta bisnis]` `PROPORTIONALARRG` menggantung pada **(TreatyYear, TreatyGroupID, ReinsTypeID)** —
**bukan FK ke ID kontrak** (PK ke TreatyYear). Jadi klausul **milik level tahun/grup/jenis**, dipakai
bersama; tidak-ikut-terhapus = **desain, bukan bug**.
**Keputusan (i sesuai a):** hapus kontrak → **kaskade** `treatybusiness` + `MTREATYSECURITY` +
`TREATYREINSURER` + `treatycontract`, dengan **popup konfirmasi Ya/Batal** yang menyebut jumlah child
terdampak (pola Master Contract Retro Life). **`PROPORTIONALARRG` (klausul) TIDAK ikut — tetap hidup.**
> ⚠️ Penyimpangan sadar: kaskade + popup untuk 3 child; klausul dikecualikan sengaja.

### Q9 — Validasi wajib-isi per klausul `[keputusan work owner]`
**Pertahankan perbedaan** (tiap klausul minta field beda), validasi per jenis di Go. **LimitMB &
Portfolio** yang sekarang tanpa validasi → **OQ terbuka** field wajibnya (tak ditebak; bagian Q3),
bukan "tanpa validasi".

### Q10 — MTREATYSECURITY `[keputusan work owner + data DBA]`
`[data DBA]` DDL: `THN_TREATY VARCHAR2(4) DEFAULT '1' NOT NULL, TOP_ID VARCHAR2(9), TP_TREATY CHAR(2),
REAS_ID CHAR(7) NOT NULL, PCT_SHARE VARCHAR2(99), USER_ID CHAR(99), REAS_SECURITY CHAR(10) NOT NULL`.
Tanpa PK; INSERT posisional mengosongkan `TOP_ID, TP_TREATY, USER_ID`; kunci lama pakai `trim()`.
**Skema baru:** tabel anak `TREATYREINSURER` yang wajar — **PK surrogate**, kolom bernama,
`PCT_SHARE` → decimal, `REAS_SECURITY` atribut biasa (bukan bagian kunci).
> ⚠️ Penyimpangan sadar: bersihkan struktur kotor.

### Q11 — Kurs `[keputusan work owner]`
`[data DBA]` tabel kurs = **`TREATYEXCHANGEYEARLY`** (`TOIDR, TOUSD, IDCURRENCY, QUARTER`, semua
VARCHAR2). **`testingKurs` yang DIPAKAI** (meski nama "testing"); **`SetTreatyArrangementDesc_Act`
mati (di-remark)** — OQ-066 terbalik. **`IDCURRENCY='10001' = USD`** untuk konversi USD→IDR (memang
harus, bukan tambalan) → jadikan rujukan master. **`QUARTER='0'` ikuti apa adanya** (arti tak pasti →
OQ kecil). Dua kolom Rp & Usd dipertahankan (memang IDR & USD). Tanggal (teks format Pega) → DATE;
uang → decimal (ADR-0003).

### Q12 — Batas konteks `[keputusan work owner]`
- **(i)** Di daftar treaty year ada `REINSURANCETYPE` + `TREATYDESC` = **master terpisah, baca saja**.
- **(ii)** Attachment ke `M_ATTACHMENTTREATY_2` **BELUM selesai di-develop** di Pega → **BUAT BARU**:
  upload attachment **di level treaty year** (mirip Master Treaty In; Google Storage, opsional, efek
  keluar ADR-0015, endpoint config ADR-0013). Fitur dilengkapi, bukan sekadar migrasi.
- **(iii)** Modul ini **satu-satunya penulis** `TREATYREINSURER`/`TREATYBUSINESS`/`PROPORTIONALARRG`;
  Claim Prop/Komite/Fac In = **read-only**. Skema baru wajib sediakan bentuk relasional yang mereka
  baca sampai konteks itu ikut migrasi (OQ-042 terkonfirmasi).

---

## Tambahan validasi (pasca-spec) `[keputusan work owner]`

- **AC 9 dipertahankan** — masa berlaku (end date) yang **berakhir sebelum mulai (start date)
  DITOLAK**, untuk tahun treaty maupun kontrak. Gerbang wajar (bukan penyimpangan; tidak ada di Pega).
- **Anti-dobel periode+grup (BARU)** — tahun treaty dengan kombinasi **(StartDate, EndDate,
  TreatyGroupID) yang sudah ada DITOLAK** (mode **b = tolak**, bukan sekadar warning). Cegah duplikat
  periode per grup.
- **AC 23 (Claude, diterima)** — UPDATE `PEGA_TREATYBUSINESS` existing hanya set 5 kolom
  (ISACTIVE/BIZCODE/BIZNAME/USERID/TGLUPDATE); REINSTYPEID dll hanya saat INSERT. Sistem baru:
  **pembaruan memperbarui SELURUH field yang dikirim** (perbaikan sadar kejanggalan, bukan tiru bug).
- **AC 72 (Claude, diterima)** — penyimpangan 3 (fitur salin dibuang) diikat sebagai AC (test yang
  menemukan endpoint/fitur salin tahun → gagal).

## Penyimpangan sadar (rekap) — jadi AC ⚠️ di spec

| # | Penyimpangan |
| --- | --- |
| 1 | Buang dualitas JSON — satu sumber `PROPORTIONALARRG` relasional (JSON mati) |
| 2 | Satu tabel generik untuk semua klausul (child = subset kolom) |
| 3 | Fitur copy tahun treaty DIBUANG (dead code) |
| 4 | Kaskade hapus 3 child + popup konfirmasi; klausul (PROPORTIONALARRG) dikecualikan |
| 5 | MTREATYSECURITY: PK surrogate + tipe wajar (buang struktur posisional kotor) |
| 6 | Semua uang/persen → decimal (ADR-0003); semua tanggal → DATE (existing banyak teks) |
| 7 | `IDCURRENCY='10001'` (USD) jadi rujukan master, bukan hardcode; kurs USD→IDR |
| 8 | Nama jujur di sistem baru (buang "Old", "testing", "JSON_KLAIM" ErrMsg) |
| 9 | Fitur BARU: upload attachment di treaty year (belum ada di Pega) |

## Dead code (JANGAN dimigrasikan)
`M_PROPORTIONALARRG` (JSON) + semua kueri `FROM m_PROPORTIONALARRG` · `PROSESCOPY` + fitur copy
(`BrowseCopyData`, `SaveTreatyYearMultiple_Act`) · `SetTreatyArrangementDesc_Act` (di-remark) ·
RD non-Old kecuali layar master jenis reasuransi · kolom `PROPORTIONALLIST`/`OBJECT` (tak di-set
procedure — konfirmasi saat migrasi).

## Status OQ
| OQ | Status |
| --- | --- |
| OQ-002 | ✅ DITUTUP — 6 procedure penulis diterima (PEGA_PROPORTIONALARRG + CHILD + TREATYCONTRACT/YEAR/REINSURER/BUSINESS). Semua: upsert `ID`, ID `'1'+lpad(seq,6/7)`, TIDAK commit sendiri, StsSimpan 1/0. `dba-procedures.md` |
| OQ-001 | ✅ DITUTUP — DDL 8 tabel diterima. Uang campur NUMBER/VARCHAR2 → semua decimal; tanggal teks → DATE |
| OQ-020 | ✅ ditutup — ReinsTypeID master bersama `REINSURANCETYPE`, filter Flag=active+Type+blacklist 12 ID (Q5) |
| OQ-042 | ✅ terkonfirmasi — modul ini satu-satunya penulis; hilir read-only |
| OQ-047 | terbuka tak-memblokir — alamat Google Storage (ADR-0013) |
| OQ-022 | keputusan GATE (nama "Out") — tak memblokir |
| OQ kecil | terbuka tak-memblokir — arti istilah klausul (Q3), field wajib LimitMB/Portfolio (Q9), arti `QUARTER='0'`, kolom PROPORTIONALLIST/OBJECT |

## Kesiapan `/to-spec` — SEMUA PEMBLOKIR DITUTUP
- Frontier kosong; 12 pertanyaan terjawab; OQ-001 & OQ-002 ditutup.
- 9 penyimpangan sadar siap jadi AC ⚠️. Dead code teridentifikasi.
- Fitur baru: upload attachment treaty year.
- Ronde 2 (opsional) sempat disebut untuk detail jalur tulis/atomisitas — **sudah terjawab** oleh
  body procedure (tidak commit sendiri → Go bisa atomik). Tidak perlu Ronde 2.
