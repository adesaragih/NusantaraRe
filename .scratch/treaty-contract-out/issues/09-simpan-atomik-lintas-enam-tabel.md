# 09: Simpan atomik lintas enam tabel

**Status:** ready-for-agent

**Blocked by:** 05, 06, 07, 08 (seluruh penulis harus ada sebelum dapat dibungkus jadi satu)

## Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin seluruh perubahan pada satu kontrak — kontrak, reinsurer,
security, business, **dan seluruh klausul** — tersimpan **bersama atau tidak sama sekali**, sehingga
tidak pernah ada kontrak yang tersimpan separuh; dan sebagai **admin master**, saya ingin
**diberi tahu bila penyimpanan gagal**, supaya saya tidak mengira data tersimpan padahal tidak.
*(User story 24–25 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | **Satu transaksi** melintasi enam tabel; commit sekali |
| `internal/services` | Orkestrasi simpan menyeluruh; rollback total |
| `internal/handlers` | Endpoint simpan kontrak utuh |
| `frontend/` | Tombol simpan tunggal; tampilan galat |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveMasterTreatyContract_SQL` | `ASM-FW-GISFW-INT` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterTreatyContract_SQL.xml` | ⚠️ `COMMIT` di dalam SQL Pega |
| `SaveMasterTreatyYear_SQL` | `ASM-FW-GISFW-INT-TREATYYEAR` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterTreatyYear_SQL.xml` | ⚠️ idem |
| `SaveMasterTreatyReinsurer_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterTreatyReinsurer_SQL.xml` | ⚠️ idem |
| `SaveMasterTreatyBusiness_SQL` | `ASM-FW-GISFW-INT-TREATYBUSINESS` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterTreatyBusiness_SQL.xml` | ⚠️ idem |
| `SaveMasterProportionalArrg` / `…Child` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/` | ⚠️ idem |
| `RefreshErrorProportionalarrg` | `@BASECLASS` / `REFRESHERRORPROPORTIONALARRG` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/RefreshErrorProportionalarrg.xml` | tampilan galat |
| `SetErrorMessage` | `@BASECLASS` / `SETERRORMESSAGE` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetErrorMessage.xml` | tampilan galat |

⚠️ **Temuan yang membuat tiket ini mungkin** `[data DBA]`: keenam procedure penulis
(`PEGA_PROPORTIONALARRG`, `PEGA_M_PROPORTIONALARRG_CHILD`, `PEGA_TREATYCONTRACT`, `PEGA_TREATYYEAR`,
`PEGA_TREATYREINSURER`, `PEGA_TREATYBUSINESS`) **TIDAK `COMMIT` sendiri** — hanya `ROLLBACK` bila
galat; commit diserahkan kepada pemanggil.

`COMMIT` yang terlihat di rule Connect-SQL Pega berada **di luar** procedure, ditulis Pega sendiri.
Jadi **tidak ada titik potong transaksi di sisi basis data** — Go dapat dan harus membungkus
seluruhnya.

⚠️ Ini **berbeda dari lima konteks Life sebelumnya** (Claim Life, Komite, PremiumList, Endorsement,
Master Contract Retro Life), yang procedure-nya `COMMIT` sendiri dan karenanya memaksa titik potong.

## ADR terkait

**ADR-0006** (identitas lewat sequence), **ADR-0007** (jejak audit setiap transaksi),
**ADR-0015** (kegagalan ditangani eksplisit, tidak ditelan).

## Acceptance criteria

- [ ] ⚠️ Seluruh perubahan satu kontrak — kontrak, reinsurer, security, business, **dan seluruh
      klausul** — ditulis dalam **satu transaksi**; kegagalan di mana pun **membatalkan seluruhnya**.
      Dibuktikan dengan menyuntikkan kegagalan pada **satu baris klausul ke-N** lalu memastikan
      **tidak ada** perubahan tersimpan pada kelima tabel lain. *(AC 37 spec; User story 24;
      `[data DBA]` — procedure tidak commit sendiri)*
- [ ] Penyimpanan yang **gagal** menghasilkan kegagalan **terang-terangan** dengan pesan yang
      menyebut **apa** yang gagal — bukan diam-diam dianggap sukses. *(AC 38 spec; User story 6,
      25; **ADR-0015**)*
- [ ] Status simpan **1 = sukses / 0 = gagal** ditegakkan; nilai **selain `1`** — termasuk kosong
      dan NULL — **selalu** dibaca sebagai kegagalan. *(AC 39 spec; `[data DBA]`)*
- [ ] ⚠️ Pesan galat di sistem baru **tidak memuat teks `"JSON_KLAIM"`**. Test yang menemukannya
      **gagal**. *(AC 40 spec; penyimpangan sadar 8)*
- [ ] Setiap penyimpanan mencatat **jejak audit** — siapa dan kapan. *(AC 41 spec; **ADR-0007**)*
- [ ] Transaksi **di-commit sekali**, di akhir; **tidak ada** commit per tabel.
- [ ] Kegagalan mengembalikan basis data ke keadaan **persis sebelum** permintaan — termasuk
      **tidak menyisakan identitas terpakai** yang membuat nomor melompat tanpa alasan.

## Blocker

**Tidak ada.** **OQ-002 ditutup** — body keenam procedure sudah diterima; perilaku transaksinya
terbukti.

## Catatan

⚠️ **Inilah kebalikan dari Master Product Name Life.** Di sana `[data DBA]` membuktikan procedure
**`COMMIT` sendiri**, sehingga atomik hanya mungkin dengan **membuang** procedure-nya
(penyimpangan sadar di konteks itu). Di sini `[data DBA]` membuktikan **sebaliknya** — dan itulah
yang membuat simpan atomik lintas enam tabel menjadi mungkin **tanpa membuang satu pun procedure**.
Motif proyek: **"periksa apakah procedure commit sendiri."**

⚠️ **Nomor urut yang melompat bukan sekadar kosmetik.** Karena identitas dibuat basis data lewat
sequence (`ADR-0006`), rollback tidak mengembalikan `nextval` yang sudah diambil. AC terakhir
menuntut agar kegagalan **tidak mengambil identitas lebih dulu daripada perlu** — bukan agar
sequence dimundurkan.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — atomisitas lintas enam tabel adalah
**satu-satunya** hal yang diuji tiket ini, dan ia **hanya berperilaku benar pada basis data
sungguhan**. Memalsukan Oracle di sini berarti tidak menguji apa pun yang penting.

```
go test ./internal/...
cd frontend && npm test
make check
```
