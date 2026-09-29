# 04: Kontrak treaty di dalam tahun treaty

**Status:** selesai (29-09-2026)

**Blocked by:** 02 (pemilih jenis reasuransi), 03 (tahun treaty sebagai induk)

## Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin membuat **kontrak treaty** di dalam sebuah tahun
treaty — dengan **jenis reasuransinya** dan **masa berlakunya sendiri** — dan mengubahnya kemudian
tanpa membuat duplikat. *(User story 4–5, 7 di spec)*

Kontrak inilah yang membuka kombinasi **(tahun, grup, jenis reasuransi)** yang dipakai reinsurer,
business, dan seluruh klausul di tiket-tiket berikutnya.

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas kontrak treaty |
| `internal/repository` | Upsert dikunci `ID`; rujukan ke tahun treaty |
| `internal/services` | Gerbang masa berlaku |
| `internal/handlers` | Endpoint kontrak |
| `frontend/` | Grid kontrak di dalam layar tahun treaty |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveTreatyContract_Act` | `@BASECLASS` / `SAVETREATYCONTRACT_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SaveTreatyContract_Act.xml` | orkestrator simpan |
| `NewInputTreatyContract_Act` | `@BASECLASS` / `NEWINPUTTREATYCONTRACT_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/NewInputTreatyContract_Act.xml` | baris baru |
| `SetUbahTreatyContract` | `@BASECLASS` / `SETUBAHTREATYCONTRACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetUbahTreatyContract.xml` | muat untuk diubah |
| `SetTanggalTreatyContract` | `@BASECLASS` / `SETTANGGALTREATYCONTRACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetTanggalTreatyContract.xml` | isi tanggal |
| `SaveMasterTreatyContract_SQL` | `ASM-FW-GISFW-INT` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterTreatyContract_SQL.xml` | → `POOLDATA.PEGA_TREATYCONTRACT` |
| `InboxTreatyContract` | `DATA-PORTAL` / `INBOXTREATYCONTRACT` / `RULE-HTML-HARNESS` | `Treaty Contract Out/Harness/InboxTreatyContract.xml` | titik masuk modul |

`[terverifikasi]` Parameter `PEGA_TREATYCONTRACT`: `ID`, `IDTreatyYear`, `ReinsTypeID`,
`ReinsTypeName`, `TreatyStartDate`, `TreatyEndDate`, pengguna, `TglUpdate`, + 2 keluaran.

`[data DBA]` Upsert dikunci `ID`; identitas `'1' || lpad(treatycontract_seq.nextval, 6, '0')`;
`TREATYSTARTDATE`/`TREATYENDDATE` di-`to_date(…,'DD/MM/YYYY')`; **tidak `COMMIT` sendiri**;
`StsSimpan` **1 = sukses / 0 = gagal**.

## ADR terkait

**ADR-0006** (identitas lewat sequence), **ADR-0007** (jejak audit), **ADR-0015** (kegagalan
eksplisit).

## Acceptance criteria

- [ ] Kontrak treaty dibuat **di dalam** sebuah tahun treaty dan menyimpan **rujukan ke tahun itu**.
      *(AC 7 spec; User story 4)*
- [ ] Menyimpan kontrak yang **sudah ada** memperbaruinya, **bukan** menambah baris baru.
      *(AC 8 spec — sisi kontrak)*
- [ ] Masa berlaku kontrak yang **berakhir sebelum dimulai** **ditolak**, dengan pesan yang
      menyebut field-nya. *(AC 9 spec — sisi kontrak)*
- [ ] Jenis reasuransi kontrak dipilih dari daftar tersaring tiket **02** — tidak diketik bebas.
      *(AC 10, 11 spec)*
- [ ] Identitas kontrak **tidak pernah diketik pengguna**; berbentuk `'1'` + 6 digit.
      *(AC 5, 6 spec; **ADR-0006**)*
- [ ] Tanggal mulai dan akhir kontrak tersimpan sebagai **tanggal**, bukan teks. *(AC 53 spec)*
- [ ] Kontrak **tidak dapat dibuat** tanpa tahun treaty induk yang ada.
- [ ] Penyimpanan yang gagal menghasilkan kegagalan **terang-terangan**. *(AC 38, 39 spec;
      **ADR-0015**)*
- [ ] Setiap penyimpanan mencatat **jejak audit**. *(AC 41 spec; **ADR-0007**)*

## Blocker

**Tidak ada.**

## Catatan

⚠️ **Kombinasi, bukan foreign key.** `[fakta bisnis — work owner]` Reinsurer, business, dan klausul
**tidak** menyimpan `ID` kontrak. Mereka menggantung pada **(TreatyYear, TreatyGroupID,
ReinsTypeID)**. Kontrak "membuka" kombinasi itu, tetapi tidak memilikinya. Ini **berbeda** dari
Master Contract Retro Life, di mana anak-anak menggantung pada `ID` kontrak — jangan menyalin pola
dari sana.

⚠️ **Editor master, bukan proses berjenjang.** `[terverifikasi]` Nol rule `Flow`, nol rule `When`,
nol `StatusAkseptasi` di seluruh 303 berkas modul. Tidak ada Submit/Decline, tidak ada assignment,
tidak ada SLA.

`[terverifikasi]` **RBAC tidak dapat direkonstruksi** dari korpus — modul ini tidak memuat ekspresi
visibilitas ber-workbasket (`discovery/flows/_METHOD-noflow.md` §3.5). Siapa yang boleh mengedit
kontrak adalah keputusan terpisah, di luar tiket ini.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata**.

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Pembacaan ulang XML — 29-09-2026 (sesi modul, lanjutan 1)

Nomor baris = baris mentah berkas korpus (satu tag per baris), diverifikasi dengan `awk 'NR==n'`. Langkah aktivitas
dibaca lengkap dengan prasyarat dan penandanya; langkah ber-`pyStepsBlockName = //` adalah langkah YANG DIKOMENTARI.

| Unsur | Bukti | Dibawa sebagai |
| --- | --- | --- |
| jalan masuk | `Section/InputTreatyContract.xml` b20778 `ReinsType` → `BrowseReinsTypeYear` (b20903, `IDTreatyYear = .ID`) → `showHarness` popup `InboxTreatyContractReinsType` (b20947) → `KirimTahunGroupID` (b20953: TreatyGroupName, ID, TreatyGroupID, TreatyYear, StartDate, EndDate, UnderwritingYear) | tombol `ReinsType` membuka `PanelKontrakTahun` untuk baris itu |
| harness menu | `Harness/InboxTreatyContractReinsType.xml` b151, judul `ReinsType` b1670, `PanggilReinsType` b1799 → `InputTreatyContractReinsType` (`PanggilReinsType.xml` b1247) | butir menu kedua; pemilih tahun lebih dulu (harness tanpa konteks) |
| kepala | `InputTreatyContractReinsType.xml` b1145 `Underwriting Year` (`.UnderwritingYear` b1176); b1358 `ReinsType` pada `.TreatyGroupName` b1386 | dua medan baca-saja; label bersilang dibawa VERBATIM + catatan |
| form | b2478 ID (`Formatted Text`), b2652 `ReinsType` (RD non-Old, `Flag "active"` b2768, tampil `.Note` b2724), b2905 `Start Date` → `SetTanggalTreatyContract` b3007, b3244 `End Date`, b3618 `Save` → `SaveTreatyContract_Act` b3642, b4363 `Modified Date`, b4547 `Username`, b5343 `Undo` → `UndoOperation` b5366, b6400 `Information` | form panel; pemilih tiket 02 |
| grid | `BrowseTreatyContract_RD` (b9145; rule-nya TIDAK diekspor korpus) disaring `InputData.HASIL13 = Param.IDTreatyYear` (`BrowseReinsTypeYear.xml` b273); `Add` b8528 → `NewInputTreatyContract_Act` b8552; kolom `Reins Type` b9164 / `Treaty Start` b9304 / `Treaty End` b9444 | `GET /tahun/{id}/kontrak`, `ORDER BY ID DESC` `[keputusan kami]` |
| tombol baris | `Edit` b10519 → `SetUbahTreatyContract` b10543 (ID, ReinsTypeID, ReinsTypeName, TreatyStartDate, TreatyEndDate); `Business List` b10842 (tiket 07); `Reinsurer List` b11308 (tiket 05); `Delete` b11809 → `BrowseDeleteRowTreatyInContract` (tiket 10) | `Edit` hidup; tiga lainnya berdiri menyebut tiketnya |
| `SaveTreatyContract_Act` | langkah 7 (HIDUP, tanpa prasyarat): `UserID ← OperatorID.pyUserName`, `TglUpdate ← @getCurrentTimeStamp()`, tanggal `dd/MM/yyyy`; langkah 8 (HIDUP): RDB `SaveMasterTreatyContract_SQL`, prasyarat `ID != ""` → jalankan, lalu `ReinsTypeID=="" && ReinsTypeName=="" && TreatyStartDate==""` → lewati; langkah 2–6 dan 12 DIKOMENTARI | `services.KontrakTreatyTCO.Simpan` |
| `SetTanggalTreatyContract` | langkah 2 `IsEndDate==1` → hanya akhir, keluar; langkah 3 `StartDate = CARIDATETIME + 8 jam`, `JumlahHari = 365`; langkah 4 (HIDUP) `ASMMessageStartDate` bila `@substring(StartDate,0,4) <> TreatyYear` (b847); langkah 5 `JumlahHari = 366` bila enam prasyarat lolos (b1039–b1191); langkah 6 `EndDate = StartDate + JumlahHari`; langkah 7–9 dikomentari | `models.AkhirKontrakBawaanTCO` + `GET /tahun/{id}/kontrak/akhir-bawaan`; gerbang tahun mulai |
| `PEGA_TREATYCONTRACT` `[data DBA]` | upsert dikunci ID; `'1'‖lpad(seq,6)`; tidak COMMIT | `repository.MasterKontrakTCO` (ID dari `SEQ_T_TREATYCONTRACT`), prosedur tidak dipanggil |

### Ralat bertanggal 29-09-2026

1. **Gerbang simpan di Pega sebagian besar mati.** `SaveTreatyContract_Act` langkah 2 (`ASMMessageReinstype`), 3
   (`ASMMessageStartDateKosong`), 5–6 (`"Data Reins Masih Kosong!!!"`), dan 12 (`"Data sudah pernah di Input"` bila
   `STSSAVE == 3`) DIKOMENTARI. Yang berjalan di Pega: kontrak baru disimpan kecuali ketiga medannya kosong sekaligus.
   Sistem baru menegakkan AC tiket: jenis reasuransi wajib dan harus ada di daftar tersaring tiket 02 (nama diambil dari
   master, bukan dari klien); kedua tanggal wajib; periode terbalik ditolak (AC 9).
2. **Satu jenis reasuransi satu kontrak per tahun** `[keputusan kami, berdasarkan fakta bisnis kombinasi]` — pesan
   VERBATIM langkah 12 yang dikomentari (`Data sudah pernah di Input`), 409 menyebut kontrak lain. Sebabnya: anak-anak
   menggantung pada kombinasi (tahun, grup, jenis); dua kontrak berjenis sama di satu tahun berbagi anak yang sama, dan
   kaskade hapus tiket 10 akan menghapus anak keduanya. **OQ-TCO-11** meminta konfirmasi work owner.
3. **Tahun tanggal mulai = tahun treaty** adalah gerbang yang HIDUP di Pega (`SetTanggalTreatyContract` langkah 4
   b847, pesan `ASMMessageStartDate` tidak diekspor) dan dibawa sebagai 422 saat simpan.
4. **Tanggal akhir bawaan ditiru APA ADANYA, termasuk anomalinya (OQ-TCO-10).** Prasyarat langkah 5 memotong cap waktu
   `yyyyMMdd…` seolah `dd/MM/yyyy` (contoh uji penulisnya `01/02/2018`). Akibatnya 366 hari hanya untuk mulai
   Oktober–Desember tahun kabisat; maksud penulisnya tampak Januari–Februari. Tanggal akhir tetap dapat diubah pemakai.
5. **`BrowseTreatyContract_RD` tidak ada di korpus** — urutan grid `ID DESC` adalah keputusan kami.
6. **Harness menu tanpa konteks** — Pega membukanya sebagai popup berkonteks; butir menu mendapat pemilih tahun.
7. **`Username` b4547** menampilkan operator yang sedang masuk (`OperatorID.pyUserName`); layar baru menampilkan
   `USERID` rekam (penulis terakhir), sejalan dengan layar tahun treaty.
8. **Placeholder berulang** di kueri anti-dobel tahun treaty (tiket 03) diganti placeholder berbeda (`:4`, `:5`) dengan
   nilai yang diikat dua kali — benar untuk pengikatan per nama maupun per kemunculan; penjaga
   `TestTCOPlaceholderTidakBerulang` menguncinya.

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| models | `tco_kontrak.go` (+uji) | `PeriksaKontrakTreaty`, `AkhirKontrakBawaanTCO` (anomali ditiru, 8 kasus uji) |
| repository | `tco_kontrak.go` (+uji) | daftar/ambil berbatas tahun, sisip (sequence), perbarui seluruh medan berbatas tahun, `CariDobel` |
| services | `tco_kontrak.go` (+uji) | `KontrakTreatyTCO`: tahun induk wajib ada, jenis dari master, satu transaksi {dobel, tulis, jejak} |
| handlers | `tco_kontrak.go` (+uji, +uji `db`) | 4 rute; POST menolak `id`, PUT menolak `id` berbeda |
| frontend | `PanelKontrakTahun.tsx` (+uji), `InboxTreatyContractReinsType.tsx` (+uji), `KONTRAK_TCO` (19 baris diuji ke korpus), `api.ts` (+3), tombol `ReinsType` hidup, butir menu kedua | |

**Status:** selesai 29-09-2026 — commit `treaty-contract-out: tiket 04 — kontrak treaty di dalam tahun`.

## Keputusan work owner 29-09-2026

- **OQ-TCO-11 — ditutup.** Jawaban: *"benar"*. Satu jenis reasuransi satu kontrak per tahun treaty (anti-dobel
  `CariDobel`, 409 `Data sudah pernah di Input`) dikonfirmasi; labelnya kini `[keputusan work owner 29-09-2026]`.
- **OQ-TCO-10 — ditutup.** Jawaban: *"Ganti jadi mulai + 1 tahun kalender"*. **Penyimpangan sadar, ralat bertanggal
  29-09-2026:** `SetTanggalTreatyContract` (dipanggil `Section/InputTreatyContractReinsType.xml` b3007/b3155) — langkah 3
  `JumlahHari = 365` (b620–b621), langkah 5 `JumlahHari = 366` bila prasyaratnya lolos (b918–b1191; prasyaratnya membaca
  cap waktu `yyyyMMdd…` seolah `dd/MM/yyyy`, sehingga yang berjalan "mulai Oktober–Desember tahun kabisat = 366 hari"),
  langkah 6 `TreatyEndDate = TreatyStartDate + JumlahHari` (b1273–b1274); langkah 7–9 ter-remark (`//` b1444, b1603,
  b1752). Kini `models.AkhirKontrakBawaanTCO(mulai)` = mulai + 1 tahun kalender dengan semantik `ADD_MONTHS(mulai, 12)`:
  29 Februari → 28 Februari, hari terakhir bulan tetap hari terakhir bulan (28 Februari tahun biasa → 29 Februari tahun
  kabisat). Uji: tahun biasa, kabisat Januari (dulu 2024-12-31), kabisat Oktober (dulu 366 hari), melewati 29 Februari,
  29 Februari, akhir bulan — dibuktikan merah terhadap rumus lama (4 kasus).

### Ralat bertanggal 29-09-2026 — tco5 menu satu butir

**Keputusan work owner tco5**: butir menu kedua `InboxTreatyContractReinsType` **dibuang**. XML membenarkannya: harness
itu bukan menu portal, melainkan popup dari form kontrak — `Section/InputTreatyContract.xml` tombol `ReinsType` b20778 →
harness b20947/b21627. Editor kontrak (`PanelKontrakTahun`) tetap, dibuka tombol `ReinsType` baris tahun treaty.
Halaman pembungkus `InboxTreatyContractReinsType.tsx` (pemilih tahun, penyimpangan 6 di atas) tidak lagi dirujuk menu;
rutenya masih ada di `App.tsx` yang memuat suntingan work owner belum di-commit — dibuang saat berkas itu di-commit.
