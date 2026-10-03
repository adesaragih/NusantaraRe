# 00: Skema tujuh tabel + `T_WORK_POLIS` + migrasi — **PREFACTOR**

**Status:** ready-for-agent

**Blocked by:** CL-01 (kerangka aplikasi + seam API — scaffolding lintas konteks)

⚠️ **PREFACTOR dan tiket PERTAMA.** Diberi nomor `00` supaya berada di depan tanpa menomori ulang
sembilan tiket yang sudah terbit. **Seluruh tiket 01–09 kini memblokir pada tiket ini.**

Alasannya: bentuk penyimpanan berubah total — dokumen JSON + tabel flat warisan → **tujuh tabel
relasional empat tingkat**. Tidak ada irisan lain yang dapat berdiri sebelum bentuk barunya ada.
*"Make the change easy, then make the easy change."*

## Hasil & nilai pengguna

Sebagai **tim migrasi**, saya ingin setiap atribut polis menjadi **kolom bernama** bertipe benar, dan
seluruh polis lama pindah **tanpa kehilangan satu nilai pun** — termasuk seluruh peserta, seluruh
rekap mata uang, seluruh spreading, dan seluruh riwayat penawaran. Sebagai **organisasi**, saya ingin
bentuk polis dapat diperiksa, dicari, dan divalidasi. *(Spec §12, US 24a–24f)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `migrations/` | DDL tujuh tabel + `T_WORK_POLIS` + sequence + **index pada setiap FK**; skrip migrasi & rekonsiliasi |
| `internal/models` | Agregat polis: header → {rekap mata uang, peserta → spreading → spreading retro, riwayat penawaran} |
| — | Skrip rekonsiliasi jumlah baris, nilai uang, dan tanggal |

## Bentuk baru — spec §12

```
T_PREMIUM_LIST (PK ID)
  ├─ T_PREMIUM_LIST_SUMMARY        1:N  FK PREMIUM_LIST_ID   ON DELETE CASCADE
  ├─ T_PREMIUM_LIST_DETAIL         1:N  FK PREMIUM_LIST_ID   ON DELETE CASCADE
  │     └─ T_PREMIUM_LIST_SPREADING       1:N  FK DETAIL_ID          ON DELETE CASCADE
  │            └─ T_PREMIUM_LIST_SPREADING_RETRO  1:N  FK SPREADING_ID  ON DELETE CASCADE
  └─ T_VIEW_SUGGEST                1:N  FK PREMIUM_LIST_ID   ON DELETE CASCADE
T_WORK_POLIS (PK ID)   ⬅ mandiri, LINTAS-LINI — bukan anak T_PREMIUM_LIST
```

Daftar kolom lengkap tiap tabel ada di **spec §12** dan
`.scratch/premiumlist-life/revisi-penyimpanan-premiumlist.md`.

## Sumber migrasi

| Sumber | Identitas | Peran |
| --- | --- | --- |
| `JSON_POLIS.DATA_JSON` (CLOB) | ditulis `PremiumList Life/RDBList/InsertJsonPolis.xml` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTJSONPOLIS`) → `POOLDATA.INSERTJSONPOLISLIFE` | pohon polis penuh |
| `JSON_OFFER_LIFE` (CLOB + 24 kolom flat) | `RDBList/SaveOfferJsonLife_SQL.xml` (`ASM-FW-GISFW-INT-OFFERJSON` / `ASM!SAVEOFFERJSONLIFE_SQL`); dibaca `RDBList/GetOfferLife_sql.xml` (`ASM-FW-GISFW-WORK-LIFE` / `RNM!GETOFFERLIFE_SQL`) | penawaran + riwayat |
| `M_LIFE_PREMIUM_DETAIL` | `RDBList/SaveMasterLPDet.xml` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET`) — **±81 kolom** | peserta |
| `M_LIFE_PREMIUM_SUMMARY` | `RDBList/InsertPLSummary.xml` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY`) → `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY` | rekap |
| `LIFEINPRODUCTION` | `RDBList/SaveLifeinProduction_SQL.xml` | produksi |

⚠️ `[terverifikasi]` `PEGA_M_LIFE_PREMIUM_SUMMARY` dipanggil **posisional** (31+ argumen); **nama
kolomnya tidak terbaca dari korpus**. Pemetaan kolom sumber → tujuan **wajib dicocokkan ke DDL
sebenarnya** saat tiket ini dikerjakan — jangan menebak dari urutan argumen.

⚠️ `[terverifikasi]` `GetOfferLife_sql` menyaring `WHERE STATUS = 'Bind' AND OLDID IS NULL`. Adanya
`OLDID` menandakan **versioning penawaran** — migrasi memutuskan apakah versi lama ikut dibawa;
temuannya **dilaporkan**, tidak ditebak.

## ADR terkait

**ADR-0003** (uang non-float), **ADR-0006** (identitas lewat sequence), **ADR-0009** (migrasi penuh).

## Acceptance criteria

- [ ] ⚠️ Skema polis **relasional penuh**: setiap atribut menjadi **kolom bernama**. Test yang
      menemukan kolom JSON menyimpan atribut polis **gagal**. *(AC 32 spec; penyimpangan sadar 1)*
- [ ] ⚠️ **Tujuh tabel** ada dengan PK sequence dan FK sesuai diagram; seluruh FK **`ON DELETE
      CASCADE`**. *(AC 46 spec)*
- [ ] ⚠️ FK spreading menunjuk **peserta**, dan FK spreading retro menunjuk **baris spreading** —
      bukan header. Test yang menemukan keduanya menggantung pada header **gagal**. *(AC 40 spec)*
- [ ] ⚠️ **`T_WORK_POLIS` ada sebagai tabel mandiri** — bukan anak `T_PREMIUM_LIST`; keadaan tangga
      **tidak** menjadi kolom header polis. Kolom minimum: identitas polis + lini, posisi/status
      tangga, audit. *(AC 45 spec; penyimpangan sadar 4)*
- [ ] ⚠️ **Setiap FK punya index.** Migrasi yang meninggalkan FK tanpa index **gagal**.
      *(AC 49 spec)*
- [ ] ⚠️ **`T_PREMIUM_LIST_DETAIL` memuat `PARENT_ID`** — FK self-reference ke `ID` peserta versi
      sebelumnya, **nullable**, ber-index. Ia dipakai konteks **Endorsement Life** untuk mencocokkan
      peserta lama↔baru. Ketiga tabel lain **tidak** memilikinya.
      *(`.scratch/endorsement-life/spec.md` §16, AC 60; `[keputusan work owner]`)*
- [ ] ⚠️ **`T_PREMIUM_LIST` memuat kolom EDM**, seluruhnya **nullable** dan kosong pada baris new
      business: `EDM_TYPE`, `OLD_POLICY_NO`, `EDM_DATE`, `EDM_NOTE`, `TYPE_CEDING`,
      `PREMI_PROPOSED`, `UANG_PERTANGGUNGAN`, `SUM_INSURED`, `JENIS_PRODUK`, `SISTEM_REASURANSI`,
      `STATUSS`, `STATUS_UPDATE`, `STATUS_SERVICE`, `START_DATE`, `END_DATE`, `NOENDORS`,
      `PL_NUMBER_EDM`, `PRODKE`, `EDMSTATUS`, `STATUSOLD`.
      *(`.scratch/endorsement-life/spec.md` §16, AC 57; `[keputusan work owner]`)*
- [ ] Versi berjalan sebuah polis dapat ditemukan sebagai baris ber-**`PRODKE` terbesar**; seluruh
      versi **hidup berdampingan**. *(Endorsement §16, AC 56)*
- [ ] ⚠️ **Properti bawaan Pega tidak menjadi kolom** — tidak ada `px*`, `py*`, `pz*`; single-page
      kosong (`Policy`, `Quotation`, `TempError`) juga tidak. *(AC 47 spec)*
- [ ] `T_PREMIUM_LIST_SUMMARY` memuat **rekap uang penuh per mata uang**, bukan hanya kode mata
      uang. *(AC 37 spec; penyimpangan sadar 2)*
- [ ] `T_PREMIUM_LIST_DETAIL` memuat **kedua jendela valuasi** — `GROSS_VALUATION_*` **dan**
      `RETROCESSION_VALUATION_*` — beserta `EFFECTIVE_DATE`, `LAPSE_DATE`, `PERIOD_MM`.
      *(AC 38 spec)*
- [ ] `FACTOR` bertipe **desimal**; nilai berdesimal tujuh angka pindah **tanpa berubah**.
      *(AC 39 spec; **ADR-0003**)*
- [ ] ⚠️ Seluruh uang dan share bertipe **desimal presisi arbitrer**; seluruh tanggal **`DATE`**;
      seluruh kolom **nullable**; identitas dari **sequence**. Test yang menemukan kolom uang
      bertipe teks atau melewati `float` **gagal**. *(AC 48 spec; penyimpangan sadar 5)*
- [ ] Seluruh polis lama pindah **tanpa kehilangan satu nilai pun**; jumlah baris per polis —
      peserta, rekap, spreading, riwayat — **sama** sebelum dan sesudah. *(AC 50 spec; **ADR-0009**)*
- [ ] Nilai uang pindah **tanpa berubah satu digit pun**; rekonsiliasi membandingkan **secara
      tepat**, bukan dengan toleransi. *(AC 50 spec; **ADR-0003**)*
- [ ] Tanggal yang berupa teks menjadi `DATE` **tanpa pergeseran zona waktu**; yang **tidak dapat
      diurai dilaporkan**, bukan didiamkan.
- [ ] Polis tanpa retrosesi pindah dengan **nol baris** spreading dan spreading retro — **bukan**
      kegagalan. *(AC 42 spec)*
- [ ] ⚠️ Migrasi **tidak mereplikasi** `@ASM.GetPageJSONString()` dan tidak menulis satu pun CLOB
      JSON. *(AC 34 spec; penyimpangan sadar 1)*
- [ ] Sequence polis pindah dengan **nilai berjalan yang benar**; penomoran **tidak melompat dan
      tidak mengulang** setelah migrasi. *(**ADR-0006**)*
- [ ] Migrasi dapat **dijalankan ulang dengan aman** dan punya **jalur mundur yang diuji**.

## Blocker

**Tidak ada pemblokir.** **OQ-001 ditutup 2026-09-16** `[keputusan work owner]` — ketujuh tabel
dirancang sendiri; presisi fisik dicocokkan DBA **di dalam tiket ini**, bukan sebagai prasyarat.

## Catatan

⚠️ **Skala: jutaan baris.** `[keputusan work owner]` `T_PREMIUM_LIST_DETAIL`, `_SPREADING`, dan
`_SPREADING_RETRO` tumbuh sangat besar — satu polis grup dapat memuat ribuan peserta, masing-masing
dengan beberapa baris spreading dan retro. Konsekuensi yang **harus** ditangani di tiket ini:
**index pada setiap FK**, **partisi per periode/tahun** dipertimbangkan, dan migrasi dijalankan
**bertahap**, bukan satu transaksi raksasa.

⚠️ **Hapus fisik polis besar itu berat** — kaskade dapat menyentuh jutaan baris.
`[keputusan work owner]` **arsip/soft-delete dipertimbangkan sebagai jalur normal**; hapus fisik
disediakan tetapi bukan jalur sehari-hari. Keputusan finalnya diambil di tiket ini dan **dicatat**.

⚠️ **Nama berbohong (OQ-066).** `[terverifikasi]` `InsertJsonPolis`, `InsertJsonPolisEDM`, dan
`SaveOfferJsonLife_SQL` benar-benar menulis **blob JSON** — ketiganya dibuang. Sebaliknya di konteks
lain nama "Json" justru menandai INSERT flat. **Baca kodenya, jangan namanya.**

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — kaskade empat tingkat, presisi desimal,
konversi tanggal, dan rekonsiliasi jumlah baris **hanya berperilaku benar pada basis data
sungguhan**; memalsukannya berarti tidak menguji apa pun yang penting.

```
go test ./internal/...
cd frontend && npm test
make check
```
