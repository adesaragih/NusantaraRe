# 09: Simpan atomik lintas enam tabel

**Status:** selesai (29-09-2026)

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

## Pembacaan ulang XML — 29-09-2026 (sesi modul, lanjutan 1)

Nomor baris = baris mentah berkas korpus `Treaty Contract Out/`.

| Unsur | Bukti | Dibawa sebagai |
| --- | --- | --- |
| COMMIT Pega | `RDBList/SaveMasterTreatyContract_SQL.xml` b95, `SaveMasterTreatyYear_SQL.xml` b86, `SaveMasterTreatyReinsurer_SQL.xml` b102, `SaveMasterTreatyBusiness_SQL.xml` b100, `SaveMasterProportionalArrg.xml` b123, `SaveMasterProportionalArrgChild.xml` b114 — `COMMIT;` di teks SQL rule, SESUDAH panggilan prosedur; `InsertToMTreatySecurity.xml` tanpa COMMIT | nol COMMIT di teks SQL (`PeriksaSQL`); commit SEKALI oleh `DalamTransaksi` |
| satu tombol simpan | tidak ada di korpus: setiap panel punya `Save`-nya sendiri (kontrak b3642, reinsurer `ViewDetailTreatyReinsurerGrid1.xml` b11405, security b20246, business b6966, 25 form klausul) dan setiap `Save` meng-COMMIT sendiri | `Save` per panel dipertahankan (paritas); **simpan utuh** = pintu API baru satu transaksi |
| status prosedur | `SaveMasterProportionalArrg.xml` b120–b121 keluaran `HASIL1`, `HASIL2`; `[data DBA]` `StsSimpan` 1 = sukses / 0 = gagal | prosedur tidak dipanggil (keputusan o); aturan "hanya 1 sukses" berlaku pada `status` jawaban simpan utuh |
| `JSON_KLAIM` | nol kemunculan di korpus modul ini; `[data DBA]` teks `ErrMsg` sebagian prosedur (salin-tempel) | penjaga teks pesan di seluruh sumber produksi modul |
| galat tampil | `Activity/RefreshErrorProportionalarrg.xml`, `Activity/SetErrorMessage.xml` | galat menyebut baris yang gagal (`GalatSimpanUtuhTCO`) |

### Ralat bertanggal 29-09-2026

1. **Tidak ada "simpan kontrak utuh" di Pega** — setiap panel menyimpan dan meng-COMMIT sendiri, sehingga kontrak separuh
   tersimpan memang mungkin di Pega. Sistem baru menambah `POST /tahun/{id}/kontrak-utuh` dan
   `PUT /tahun/{id}/kontrak/{kid}/utuh`; `Save` per panel tetap ada dan masing-masing atomik untuk barisnya + jejaknya.
2. **Tombol simpan tunggal di layar BELUM dibangun**: panel-panel memegang isiannya sendiri (AC 29 tiket 08) dan layar
   Pega tidak punya tombol itu. Klien API `simpanKontrakUtuh` tersedia; letak dan bentuk tombolnya menunggu keputusan
   (**OQ-TCO-19**).
3. **Identitas tidak terpakai saat gagal** dicapai dengan identitas SEMENTARA selama transaksi (`S#########T`) dan
   pengambilan nomor sequence SESUDAH seluruh baris lolos (`TetapkanIdentitasTCO`) — bukan dengan memundurkan sequence.
   Reinsurer baru yang punya security diganti identitasnya dengan salin → alihkan security → buang (FK tanpa ON UPDATE).
4. **Pembaca modul melihat tulisan permintaan yang sama** (`bacaTCO`): security membaca reinsurer barunya, anak klausul
   membaca induk barunya — tanpa itu penulis per baris tidak dapat dipakai ulang di satu transaksi.
5. **Urutan permintaan dipakai apa adanya**: kontrak → reinsurer (+ security-nya) → business → klausul (induk sebelum anak).

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| models | `tco_status_simpan.go` (+uji) | `StatusSimpanSuksesTCO` — hanya "1"; penjaga teks `JSON_KLAIM` di sumber produksi |
| repository | `tco_transaksi_utuh.go` (+uji), `tco_identitas.go`, 11 berkas `tco_*.go` | mode transaksi utuh, identitas sementara, penetapan identitas + jejak; 29 pembaca dialihkan ke `bacaTCO` |
| services | `tco_simpan_utuh.go` (+uji) | `SimpanUtuhTCO` — kelima penulis per baris di SATU transaksi, `GalatSimpanUtuhTCO` menyebut baris |
| handlers | `tco_simpan_utuh.go` (+uji, +uji `db`) | 2 rute; kode HTTP dari galat baris, pesan menyebut baris, galat server tidak bocor |
| frontend | `api.ts` (+`simpanKontrakUtuh`, `statusSimpanSukses`), `simpanUtuh.test.ts` | tanpa tombol (OQ-TCO-19) |

**Status:** selesai 29-09-2026 — commit `treaty-contract-out: tiket 09 — simpan atomik lintas enam tabel`.

## Keputusan work owner 29-09-2026

- **OQ-TCO-19 — ditutup.** Jawaban: *"tidak perlu"*. **Ralat bertanggal 29-09-2026 — sebagian `wontfix` (keputusan work
  owner):** rute `POST /tahun/{id}/kontrak-utuh` dan `PUT /tahun/{id}/kontrak/{kid}/utuh` DIBUANG beserta handler,
  orkestrator `SimpanUtuhTCO`, cache master per permintaan, identitas sementara (`S…T`, `TetapkanIdentitasTCO`), aturan
  status `"1"` di jawaban, dan klien frontend `simpanKontrakUtuh` — seluruhnya tanpa pemanggil. Penyimpanan tetap **per
  panel seperti Pega**.
- **Yang tetap berlaku** (dipakai panel): setiap `Save` panel satu transaksi, commit sekali oleh `DalamTransaksi`, nol
  COMMIT di teks SQL, jejak setiap simpan, galat terang, penjaga teks `JSON_KLAIM` (AC 40,
  `models/tco_teks_galat_test.go`), dan pembaca-lewat-transaksi `DenganBacaTxTCO` (`repository/tco_baca_tx.go`).
- **`wontfix`**: AC 37 (satu transaksi untuk seluruh isi satu kontrak), AC 39 (status `"1"` jawaban simpan utuh), dan
  "identitas tidak terpakai saat gagal lintas baris" — ketiganya hanya bermakna untuk simpan utuh.

## Ralat bertanggal 29-09-2026 — tco4 (nol tabel baru) `[keputusan work owner]`

- Satu transaksi per panel tetap, kini atas tabel warisan. **Jejak gugur** (bagian "jejak setiap simpan" di keputusan
  OQ-TCO-19 tidak berlaku lagi).
