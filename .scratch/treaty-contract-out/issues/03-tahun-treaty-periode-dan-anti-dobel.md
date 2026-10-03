# 03: Tahun treaty — CRUD, gerbang periode, dan anti-dobel

**Status:** ready-for-agent

**Blocked by:** 01 (skema + sequence harus ada)

## Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin membuat dan mengubah **tahun treaty** beserta grup
treaty, underwriting year, proporsi, dan masa berlakunya — **tanpa pernah mengetik nomor
identitas**, **tanpa** dapat membuat periode yang berakhir sebelum dimulai, dan **tanpa** dapat
membuat tahun yang sama dua kali. *(User story 1–3, 5–6 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas tahun treaty |
| `internal/repository` | Upsert dikunci `ID`; identitas dari sequence |
| `internal/services` | Gerbang periode; gerbang anti-dobel |
| `internal/handlers` | Endpoint tahun treaty |
| `frontend/` | Layar daftar + form tahun treaty |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveTreatyYear_Act` | `@BASECLASS` / `SAVETREATYYEAR_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SaveTreatyYear_Act.xml` | orkestrator simpan |
| `NewInputTreatyYear_Act` | `@BASECLASS` / `NEWINPUTTREATYYEAR_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/NewInputTreatyYear_Act.xml` | baris baru |
| `SetTreatyYear_Act` | `@BASECLASS` / `SETTREATYYEAR_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetTreatyYear_Act.xml` | isi form dari baris terpilih |
| `CheckYear` | `@BASECLASS` / `CHECKYEAR` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/CheckYear.xml` | ⚠️ **satu-satunya** validasi: `isNumber(TreatyYear)` |
| `SaveMasterTreatyYear_SQL` | `ASM-FW-GISFW-INT-TREATYYEAR` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterTreatyYear_SQL.xml` | → `POOLDATA.PEGA_TREATYYEAR` |
| `BrowseTreatyYear_RD` | `ASM-FW-GISFW-INT-TREATYYEAR` / `BROWSETREATYYEAR_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Treaty Contract Out/ReportDefinition/BrowseTreatyYear_RD.xml` | daftar tahun |

`[data DBA]` `POOLDATA.PEGA_TREATYYEAR` — **upsert dikunci `ID`**; identitas baru dibuat di basis
data sebagai `'1' || lpad(TreatyYear_seq.nextval, 6, '0')`; `TGLUPDATE` diisi `SYSDATE`;
**tidak `COMMIT` sendiri**; keluaran `StsSimpan` **1 = sukses / 0 = gagal**.

## ADR terkait

**ADR-0006** (identitas lewat sequence basis data), **ADR-0007** (jejak audit),
**ADR-0015** (kegagalan ditangani eksplisit, tidak ditelan).

## Acceptance criteria

- [ ] Tahun treaty dapat dibuat dengan **grup treaty, underwriting year, proporsi, dan masa
      berlakunya**. *(AC 4 spec; User story 1)*
- [ ] Identitas tahun treaty **tidak pernah diketik pengguna** — dibuat dari sequence.
      *(AC 5 spec; **ADR-0006**)*
- [ ] Identitas berbentuk `'1'` diikuti nomor urut ber-*padding* nol **6 digit**.
      *(AC 6 spec; `[data DBA]`)*
- [ ] Menyimpan tahun treaty yang **sudah ada** memperbaruinya, **bukan** menambah baris baru.
      *(AC 8 spec — sisi tahun; `[data DBA]` upsert dikunci `ID`)*
- [ ] Masa berlaku yang **berakhir sebelum dimulai** **ditolak**, dengan pesan yang menyebut
      field-nya. *(AC 9 spec — sisi tahun)*
- [ ] **Tahun treaty yang sama tidak dapat dibuat dua kali**: kombinasi **(`STARTDATE`, `ENDDATE`,
      `TREATYGROUPID`)** yang **sudah ada** ditolak, dengan pesan yang menyebut tahun treaty mana
      yang sudah memakainya. *(AC 73 spec; `[keputusan work owner]` — **aturan baru**)*
- [ ] Gerbang anti-dobel **tidak** menghalangi pembaruan baris itu sendiri — memperbarui tahun
      treaty yang sudah ada dengan periode & grup yang tidak berubah tetap **berhasil**.
- [ ] Tanggal mulai dan tanggal akhir tersimpan sebagai **tanggal**, bukan teks. *(AC 53 spec)*
- [ ] Penyimpanan yang gagal menghasilkan kegagalan **terang-terangan**; nilai status selain `1`
      **selalu** dibaca sebagai kegagalan. *(AC 38, 39 spec; **ADR-0015**)*
- [ ] Setiap penyimpanan mencatat **jejak audit** — siapa dan kapan. *(AC 41 spec; **ADR-0007**)*
- [ ] ⚠️ **Fitur salin tahun treaty tidak dibangun.** Tidak ada jalur — layar, endpoint, maupun
      pekerjaan latar — yang menyalin isi satu tahun treaty ke tahun lain. Test yang menemukan
      jalur semacam itu **gagal**. *(AC 72 spec; `[fakta bisnis — work owner]` — penyimpangan
      sadar 3)*

## Blocker

**Tidak ada.**

## Catatan

⚠️ **Mengapa fitur salin dibuang.** `[fakta bisnis — work owner]` `POOLDATA.PROSESCOPY` menyalin
isi satu tahun treaty ke tahun lain, dan dalam praktiknya membawa **nilai tahun lalu** — misalnya
batas QS tahun 2025 — ke tahun yang semestinya berbeda. Itu sumber kesalahan, bukan penghemat
waktu. Yang **tidak** dimigrasikan: `RDBList/SaveMasterCopyData_SQL.xml`,
`Activity/BrowseCopyData.xml`, `Activity/SaveTreatyYearMultiple_Act.xml`,
`Activity/NewInputTreatyYearMultiple_Act.xml`, dan bagian salin
`Activity/BrowseDeleteRowTreatyInContract.xml`.

⚠️ **Validasi existing sangat tipis.** `[terverifikasi]` `Activity/CheckYear.xml` hanya memeriksa
`@Default.isNumber(InputTreatyYear.TreatyYear)`. Tidak ada gerbang periode, tidak ada anti-dobel.
Kedua gerbang di tiket ini adalah **tambahan sadar** `[keputusan work owner]`, bukan tiruan.

⚠️ `[terverifikasi]` `RDBList/SaveMasterCopyData_SQL.xml` menulis
`dbms_output.put_line(errmsg)` dengan variabel lokal `errmsg` yang **dideklarasikan tetapi tidak
pernah diisi** — selalu mencetak NULL. Dicatat sebagai jejak; tidak dimigrasikan.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — gerbang anti-dobel adalah pertanyaan
keunikan di basis data; memalsukannya berarti tidak mengujinya.

```
go test ./internal/...
cd frontend && npm test
make check
```
