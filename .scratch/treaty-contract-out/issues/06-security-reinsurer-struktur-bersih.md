# 06: Security reinsurer — struktur bersih

**Status:** selesai (29-09-2026)

**Blocked by:** 05 (security menggantung pada reinsurer)

## Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin mencatat **security** di bawah seorang reinsurer beserta
porsinya, supaya eksposur berjenjang terlihat; dan sebagai **organisasi**, saya ingin baris security
punya **identitas sendiri**, supaya mengubah nama security tidak memutus rujukannya.
*(User story 13–14 di spec)*

⚠️ Inilah tiket dengan pembersihan struktur paling dalam di konteks ini.

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas security reinsurer |
| `internal/repository` | Tulis/baca dengan **kolom bernama** dan **PK surrogate** |
| `internal/services` | Aturan per baris |
| `internal/handlers` | Endpoint security |
| `frontend/` | Grid security di bawah baris reinsurer |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveSecurityReinsurer_Act` | `@BASECLASS` / `SAVESECURITYREINSURER_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SaveSecurityReinsurer_Act.xml` | orkestrator simpan |
| `ShowEditSecurityReinsurer` | `@BASECLASS` / `SHOWEDITSECURITYREINSURER` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/ShowEditSecurityReinsurer.xml` | muat untuk diubah |
| `InsertToMTreatySecurity` | `ASM-FW-GISFW-INT-MTREATYSECURITY` / `ASM!INSERTTOMTREATYSECURITY` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/InsertToMTreatySecurity.xml` | ⚠️ INSERT **posisional** |
| `UpdateMTreatySecurity` | `ASM-FW-GISFW-INT-MTREATYSECURITY` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/UpdateMTreatySecurity.xml` | ⚠️ kunci `trim()` |
| `DeleteSecurityReinsurer` | `ASM-FW-GISFW-INT-MTREATYSECURITY` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/DeleteSecurityReinsurer.xml` | ⚠️ kunci `trim()` |
| `SelectSecurityReinsurer` | `ASM-FW-GISFW-INT-MTREATYSECURITY` / `RULE-OBJ-REPORT-DEFINITION` | `Treaty Contract Out/ReportDefinition/SelectSecurityReinsurer.xml` | daftar |

`[terverifikasi]` Keadaan existing — **satu-satunya master di modul ini yang ditulis SQL mentah**:

```
insert into mtreatysecurity values ({…}, '', '', {…}, {…PCT_SHARE}, '', {…REAS_SECURITY})
update mtreatysecurity set THN_TREATY=…, PCT_SHARE=…, REAS_SECURITY=…
 where REAS_ID = … and trim(REAS_SECURITY) = trim(…)
delete from mtreatysecurity where REAS_ID = … and trim(REAS_SECURITY) = trim(…)
```

⚠️ INSERT **tanpa daftar kolom** — tujuh nilai berposisi, **tiga di antaranya string kosong**.
⚠️ Ketiganya **tanpa `COMMIT`**.
⚠️ Kunci pembaruan dan penghapusan memakai **`trim(REAS_SECURITY)`** — nama dipakai sebagai bagian
kunci, dan `trim()` mengakui datanya bertabur spasi.

`[data DBA]` DDL existing: `THN_TREATY VARCHAR2(4) DEFAULT '1' NOT NULL`, `TOP_ID VARCHAR2(9)`,
`TP_TREATY CHAR(2)`, `REAS_ID CHAR(7) NOT NULL`, `PCT_SHARE VARCHAR2(99)`, `USER_ID CHAR(99)`,
`REAS_SECURITY CHAR(10) NOT NULL` — **tanpa primary key**. INSERT posisional mengosongkan
`TOP_ID`, `TP_TREATY`, `USER_ID`.

⚠️ **Penyimpangan sadar 5 — struktur dibersihkan.** `[keputusan work owner]` Tabel anak
`TREATYREINSURER` yang wajar: **PK surrogate**, seluruh kolom **bernama dan diisi eksplisit**,
`PCT_SHARE` **desimal**, `REAS_SECURITY` **atribut biasa** — bukan bagian kunci.

## ADR terkait

**ADR-0003** (persen non-float), **ADR-0007** (jejak audit), **ADR-0015** (kegagalan eksplisit).

## Acceptance criteria

- [ ] Security dicatat **di bawah seorang reinsurer** beserta porsinya. *(AC 17 spec; User story 13)*
- [ ] ⚠️ Baris security punya **primary key surrogate** sendiri; `REAS_SECURITY` adalah **atribut
      biasa**, **bukan** bagian kunci. Test yang menemukan kunci berbasis nama security **gagal**.
      *(AC 18 spec; User story 14; penyimpangan sadar 5)*
- [ ] ⚠️ Seluruh kolom baris security **diisi eksplisit dan bernama**; **tidak ada** penulisan
      berposisi. Test yang menemukan penulisan tanpa daftar kolom **gagal**. *(AC 19 spec;
      penyimpangan sadar 5)*
- [ ] ⚠️ Pencocokan baris security **tidak** memakai `trim()` atas nilai kunci. *(AC 20 spec;
      penyimpangan sadar 5)*
- [ ] ⚠️ `PCT_SHARE` bertipe **desimal**, bukan teks. *(AC 51 spec; **ADR-0003**; penyimpangan
      sadar 6)*
- [ ] Mengubah **nama security** **tidak** memutus rujukan barisnya ke reinsurer. *(User story 14)*
- [ ] Baris security dapat **ditambah dan dihapus** tanpa menyentuh baris reinsurer induknya —
      kecuali ketika penyimpanan dilakukan sebagai satu kesatuan (tiket 09).
- [ ] Menghapus seorang reinsurer **menghapus juga** baris security di bawahnya. *(lihat tiket 10
      untuk kaskade dari kontrak)*
- [ ] Penyimpanan yang gagal menghasilkan kegagalan **terang-terangan**. *(AC 38, 39 spec;
      **ADR-0015**)*
- [ ] Setiap penyimpanan mencatat **jejak audit**. *(AC 41 spec; **ADR-0007**)*

## Blocker

**Tidak ada.**

⚠️ `[terbuka]` **Tiga kolom yang dikosongkan INSERT posisional** — `TOP_ID`, `TP_TREATY`,
`USER_ID`. `[data DBA]` DDL-nya sudah diketahui, tetapi **apakah ketiganya masih punya arti** tidak
terbaca dari korpus. Tiket ini **membawanya sebagai kolom bernama**; bila migrasi (tiket 01)
membuktikan seluruh barisnya kosong, pembuangannya adalah keputusan terpisah — **jangan dibuang di
sini atas inisiatif sendiri**.

## Catatan

⚠️ **Mengapa `trim()` tidak dibawa.** Pemakaian `trim()` di kunci adalah **pengakuan bahwa datanya
kotor**, bukan aturan bisnis. Membawanya ke sistem baru berarti membawa kekotorannya, dan menutup
kemungkinan dua security yang namanya hanya berbeda spasi tetap dianggap berbeda. Migrasi (tiket 01)
memberi setiap baris identitas sendiri; sesudah itu nama tidak perlu lagi memikul beban kunci.

⚠️ `RDBList/DeleteFromTreatyReinsurer_Act.xml` menghapus `MTREATYSECURITY` lalu `TREATYREINSURER`
**tanpa `COMMIT`** — konsisten dengan temuan bahwa modul ini menyerahkan batas transaksi kepada
pemanggil (tiket 09).

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — keutuhan rujukan setelah pemberian PK dan
hilangnya kunci berbasis nama **hanya berperilaku benar pada basis data sungguhan**.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Pembacaan ulang XML — 29-09-2026 (sesi modul, lanjutan 1)

Nomor baris = baris mentah berkas korpus `Treaty Contract Out/`; langkah aktivitas dibaca lengkap.

| Unsur | Bukti | Dibawa sebagai |
| --- | --- | --- |
| jalan masuk | `Section/ViewDetailTreatyReinsurerGrid1.xml` b5277 `Security Reinsurer` → `DataTransform/SetSecurityReinsurer.xml` (b520 `CARI8 = Param.THN_TREATY` = `.TreatyYear`, b550 `CARI9 = Param.REAS_ID` = `.ID`, b278 `HASILD21 = 1`); bagian security `Section/InputTreatyContractReinsType.xml` b14885 `HASILD21 == 1` | tombol baris reinsurer membuka `PanelSecurityReinsurer` |
| daftar | `ReportDefinition/SelectSecurityReinsurer.xml` saringan `A AND D` (b550): `.THN_TREATY = Param.THN_TREATY`, `.REAS_ID = Param.REAS_ID`; pyMaxRecords 500 (b758); tanpa urutan | `GET .../reinsurer/{rid}/security`, saringan sama, `ORDER BY ID ASC` |
| grid | b16088 `Reas Security` (`.REAS_SECURITY` b16724), b16228 `Security Name` (`.CLIENTNAME` b16878), b16368 `Percent Share` (`.PCT_SHARE` b17013); `Edit` b17252 → `ShowEditSecurityReinsurer` (b360/b381 `HASILD2 = HASILD3 = 1`); `Delete` b17559 → `DeleteSecurityReinsurer` (aksi `refresh`, tanpa konfirmasi) | grid + tombol baris |
| form | `Add` b15459 → `InputNewSecurityReinsurer` (b172 `HASILD3 = 0`); b19468 `Security ID` (`InputTreatySecurity.REAS_SECURITY` b19499, nonaktif b19515); b19648 `Security Name` wajib (b19642/b19693), pemilih `BrowseAgentReinsSOA_RD` b19711 (`.ID` → `REAS_SECURITY` b19739, `.ClientName` tampil); b19888 `%Share` (`PCT_SHARE` b19919, onchange `SetErrorMessageReinsurer` b19952); `Save` b20246 | `Security ID` baca-saja; `Security Name` = pemilih master AGENT tiket 05 |
| simpan | `Activity/SaveSecurityReinsurer_Act.xml`: langkah 1 `CARI10 = InputTreatySecurity.REAS_SECURITY` (b420); langkah 3 `Local.IsUpdate = true` bila nama sama (b760/b837); langkah 4 `THN_TREATYID = IDTreatyYear` (b999); INSERT bila `HASILD3==0` (b1225), UPDATE bila `HASILD3==1` (b1409) | satu transaksi {kunci kontrak, dobel 409, sisip/perbarui berkunci ID, jejak} |
| SQL warisan | `RDBList/InsertToMTreatySecurity.xml` b60 (7 nilai posisional, tiga `''`); `UpdateMTreatySecurity.xml` b85–b89 (`where REAS_ID = CARI9 and trim(REAS_SECURITY) = trim(CARI10)`); `DeleteSecurityReinsurer.xml` b85 (`... trim(REAS_SECURITY) = trim(CARI12)`) | kolom bernama, `ID` + `REAS_ID`, nol pemangkas spasi |
| hapus reinsurer | `Activity/DeleteTreatyReins_Act.xml` b408 → `RDBList/DeleteFromTreatyReinsurer_Act.xml` b60–b64 (security lalu reinsurer) | FK `ON DELETE CASCADE` (migrasi 303) — tombol hapus reinsurer milik tiket 10 |

### Ralat bertanggal 29-09-2026

1. **Ganti nama security memutus rujukan di Pega — dan diam.** `CARI10` diisi nama BARU (b420, b953) lalu dipakai sebagai
   kunci UPDATE (b89): baris lama tidak pernah cocok. Bila nama baru kebetulan dipegang baris lain reinsurer yang sama,
   baris LAIN itulah yang tertimpa. Sistem baru mencocokkan `ID` (User story 14).
2. **DELETE berkunci nama menghapus SEMUA baris bernama sama** di bawah reinsurer itu. Sistem baru menghapus satu `ID`.
3. **`Local.IsUpdate` mati**: dihitung (langkah 3) tetapi sisip/perbarui diputus `HASILD3`, sehingga security yang sama
   dapat tersisip berulang. Sistem baru menolak dobel per reinsurer dengan 409 `[keputusan kami]` (**OQ-TCO-17**).
4. **`%Share` tidak diperiksa di Pega**: onchange `SetErrorMessageReinsurer` memeriksa `InputTreatyReinsurer.PctShare`
   (`Activity/SetErrorMessageReinsurer.xml` b518) — halaman REINSURER, residu salin-tempel. AC 17 ("beserta porsinya")
   menjadikannya wajib; rentang 0..100 dari maksud gerbang yang sama (`SetErrorMessageBetween`) (**OQ-TCO-17**).
5. **`CLIENTNAME` bukan kolom DDL** (tujuh kolom `[data DBA]`): nama tampil dibaca dari master `AGENT` (LEFT JOIN, agen
   nonaktif tetap bernama), tidak disimpan.
6. **`THN_TREATYID` (langkah 4) tidak ditulis SQL mana pun** — residu; tidak dibawa. `THN_TREATY` = `TreatyYear`
   reinsurer, diturunkan server.
7. `TOP_ID`, `TP_TREATY`, `USER_ID` ditulis **NULL eksplisit** saat sisip dan tidak disentuh saat diperbarui — setara
   warisan, sesuai blocker tiket ini.

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| models | `tco_security.go` (+uji) | `PeriksaSecurityTCO`, `UraiShareSecurityTCO`, penjaga nol float |
| repository | `tco_security.go` (+uji) | `MasterSecurityTCO`: kolom bernama, kunci `ID` + `REAS_ID`, nol `TRIM`; penjaga lintas-modul `TestTCONolInsertPosisional` |
| services | `tco_security.go` (+uji) | `SecurityTCO`: reinsurer induk diturunkan dari jalur, dobel 409, ganti nama = baris yang sama, hapus satu baris |
| handlers | `tco_security.go` (+uji, +uji `db` termasuk kaskade FK) | 4 rute |
| frontend | `PanelSecurityReinsurer.tsx` (+uji), `SECURITY_TCO` (19 baris diuji ke korpus), `api.ts` (+3), tombol `Security Reinsurer` hidup | |

**Status:** selesai 29-09-2026 — commit `treaty-contract-out: tiket 06 — security reinsurer, struktur bersih`.

## Keputusan work owner 29-09-2026

- **OQ-TCO-17 — ditutup.** Jawaban: *"setuju"*. Security dobel per reinsurer ditolak 409 dan `%Share` wajib 0..100 —
  keduanya **penyimpangan sadar dari Pega** (Pega tidak menegakkan keduanya) `[keputusan work owner 29-09-2026]`.

## Ralat bertanggal 29-09-2026 — tco4 (nol tabel baru) `[keputusan work owner]`

- **RALAT AC 18 dan AC 20** — `MTREATYSECURITY` warisan TANPA identitas `[data DBA]`: kunci baris
  **(`REAS_ID`, `TRIM(REAS_SECURITY)`)**, PERSIS `UpdateMTreatySecurity` b89 dan `DeleteSecurityReinsurer` b85.
  PK surrogate dan larangan `TRIM` gugur. `ID` security di API = `REAS_SECURITY` terpangkas; mengganti nama security
  menimpa baris yang sama (UPDATE berkunci nama LAMA) dan kuncinya menjadi nama baru.
- AC 19 (kolom bernama) tetap: sisipan menulis daftar kolom dengan nilai yang sama dengan sisipan posisional
  `InsertToMTreatySecurity` b60 (`TOP_ID`, `TP_TREATY`, `USER_ID` = NULL). `PCT_SHARE` VARCHAR2(99): desimal teks.
- Kaskade reinsurer → security: tanpa FK; dijalankan layanan (`DeleteFromTreatyReinsurer_Act`), bukan Oracle.
- **Jejak gugur**.
