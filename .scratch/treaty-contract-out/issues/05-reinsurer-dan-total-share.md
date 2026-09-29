# 05: Reinsurer + total share

**Status:** selesai (29-09-2026)

**Blocked by:** 04 (kombinasi tahun/grup/jenis dibuka oleh kontrak)

## Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin mencatat **para reinsurer** pada sebuah kombinasi tahun,
grup, dan jenis reasuransi, masing-masing dengan **share** dan **komisi reasuransi**-nya; dan
sebagai **underwriter**, saya ingin melihat **total share** seluruh reinsurer, supaya saya tahu
apakah penempatan sudah penuh. *(User story 11–12 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas reinsurer pada kombinasi |
| `internal/repository` | Upsert dikunci `ID`; baca daftar per kombinasi |
| `internal/services` | Penjumlahan total share |
| `internal/handlers` | Endpoint daftar + simpan reinsurer |
| `frontend/` | Grid reinsurer dengan baris total |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveTreatyReinsurerDetail1_Act` | `@BASECLASS` / `SAVETREATYREINSURERDETAIL1_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SaveTreatyReinsurerDetail1_Act.xml` | orkestrator simpan |
| `NewTreatyReinsurerDetail_Act` | `@BASECLASS` / `NEWTREATYREINSURERDETAIL_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/NewTreatyReinsurerDetail_Act.xml` | baris baru |
| `SetUbahTreatyReinsurerList_Act` | `@BASECLASS` / `SETUBAHTREATYREINSURERLIST_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetUbahTreatyReinsurerList_Act.xml` | muat untuk diubah |
| `SetErrorMessageReinsurer` | `@BASECLASS` / `SETERRORMESSAGEREINSURER` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SetErrorMessageReinsurer.xml` | ⚠️ normalisasi koma→titik |
| `SaveMasterTreatyReinsurer_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterTreatyReinsurer_SQL.xml` | → `POOLDATA.PEGA_TREATYREINSURER` |
| `GetMasterReinsurerList` | `ASM-FW-GISFW-INT-TREATYREINSURER` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/GetMasterReinsurerList.xml` | baca daftar per kombinasi |
| `BrowseDetailTreatyReisurer_RD` | `ASM-FW-GISFW-INT-TREATYREINSURER` / `RULE-OBJ-REPORT-DEFINITION` | `Treaty Contract Out/ReportDefinition/BrowseDetailTreatyReisurer_RD.xml` | daftar |

`[terverifikasi]` `GetMasterReinsurerList` menyaring **`TreatyYear` + `TreatyGroupID` +
`ReinsTypeID`** — kombinasi, bukan `ID` kontrak.
`[terverifikasi]` Parameter `PEGA_TREATYREINSURER` (19 in + 2 out): `ID`, `TreatyYear`,
`TreatyGroupID`, `TreatyGroupName`, `ReinsTypeID`, `ReinsTypeName`, `ReinsurerID`, `CLIENTID`,
`NAME`, `Ricomm`, `PctShare`, `IUDate`, `UserId`, `StartDate`, `EndDate`, `StatusOn`, `StdRating`,
`OperatorName`, `TglUpdate`.

`[data DBA]` `RICOMM` dan `PCTSHARE` sudah **`NUMBER`** di existing; sisanya `VARCHAR2`.
Upsert dikunci `ID`; identitas `'1' || lpad(M_TREATYREINSURER_SEQ.nextval, 6, '0')`; **tidak
`COMMIT` sendiri**; `StsSimpan` **1 = sukses / 0 = gagal**.

## ADR terkait

**ADR-0003** (uang & persen non-float), **ADR-0006** (identitas lewat sequence),
**ADR-0007** (jejak audit), **ADR-0015** (kegagalan eksplisit).

## Acceptance criteria

- [ ] Reinsurer dicatat pada kombinasi **(tahun, grup, jenis reasuransi)**, masing-masing dengan
      **share** dan **komisi reasuransi**. *(AC 14 spec; User story 11)*
- [ ] **Total share** seluruh reinsurer pada satu kombinasi **terlihat** bagi pengguna.
      *(AC 15 spec; User story 12)*
- [ ] ⚠️ Share dan komisi diperlakukan sebagai **desimal presisi arbitrer**; **tidak** melewati
      `float` dan **tidak** dibulatkan ke bilangan bulat. *(AC 16 spec; **ADR-0003**; penyimpangan
      sadar 6)*
- [ ] Identitas reinsurer **tidak pernah diketik pengguna**; berbentuk `'1'` + 6 digit.
      *(AC 5, 6 spec; **ADR-0006**)*
- [ ] Menyimpan reinsurer yang **sudah ada** memperbaruinya, **bukan** menambah baris baru.
      *(AC 8 spec)*
- [ ] Daftar reinsurer disaring **per kombinasi**, bukan per `ID` kontrak.
- [ ] Penyimpanan yang gagal menghasilkan kegagalan **terang-terangan**. *(AC 38, 39 spec;
      **ADR-0015**)*
- [ ] Setiap penyimpanan mencatat **jejak audit**. *(AC 41 spec; **ADR-0007**)*

## Blocker

**Tidak ada.**

⚠️ `[terbuka]` **Total share = 100% bukan gerbang di tiket ini.** Existing **tidak** menolak
kombinasi yang totalnya ≠ 100 pada tingkat reinsurer; yang ada hanyalah penegakan 100% pada baris
**anak klausul** (tiket 08). Apakah total share reinsurer wajib 100% adalah **fakta bisnis yang
belum ditetapkan** — AC di atas hanya mewajibkan totalnya **terlihat**, bukan ditegakkan. Bila
Product + UW menetapkan sebaliknya, itu AC tambahan, bukan perubahan tiket ini.

## Catatan

⚠️ **Angka hari ini disimpan sebagai teks ber-koma desimal.** `[terverifikasi]`
`Activity/SetErrorMessageReinsurer.xml` (`@BASECLASS!SETERRORMESSAGEREINSURER`) menjalankan
`@replaceAll(InputTreatyReinsurer.PctShare, ",", ".")` dan hal yang sama untuk `Ricomm` — bukti
bahwa nilai masuk dari layar dalam format berkoma. Di sistem baru **normalisasi dilakukan di batas
masukan**, sekali, bukan ditempel di tengah alur simpan.

⚠️ `[terverifikasi]` Kolom `STDRATING` adalah **field yang dipakai-ulang** — pola yang sama sudah
ditemukan di Master Contract Retro Life. Isi sebenarnya belum terverifikasi; bawa apa adanya dan
jangan menafsirkan.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — presisi desimal share dan komisi hanya
terbukti benar pada basis data sungguhan.

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Pembacaan ulang XML — 29-09-2026 (sesi modul, lanjutan 1)

Nomor baris = baris mentah berkas korpus; langkah aktivitas dibaca lengkap (prasyarat, `//` = dikomentari).

| Unsur | Bukti | Dibawa sebagai |
| --- | --- | --- |
| jalan masuk | `InputTreatyContractReinsType.xml` b11308 `Reinsurer List` → `BrowseTreatyReinsurerList_Act` b11325 (`CARI4 = Param.TreatyYear`, `CARI5 = TreatyGroupID`, `CARI6 = ReinsTypeID`) + `SetTreatyReinsurerList_Act`; panel `ViewDetailTreatyReinsurerGrid1` b13311 | tombol `Reinsurer List` membuka `PanelReinsurerKombinasi` |
| kunci kombinasi | `RDBList/GetMasterReinsurerList.xml` `TreatyYear={CARI4} and TreatyGroupID={CARI5} and ReinsTypeID={CARI6}`; `ReportDefinition/BrowseDetailTreatyReisurer_RD.xml` saringan b651/b668/b685, urut `.ID ASC` b727/b730 | kombinasi diturunkan server dari tahun (TREATYYEAR teks, TREATYGROUPID) + kontrak (REINSTYPEID) |
| total share | `BrowseTreatyReinsurerList_Act` langkah 3 (loop `pTotalShare + toDecimal(replaceAll(.PctShare,",","."))`), langkah 4 `TotalShare`; `SetUbahTreatyReinsurerList_Act` langkah 3 (`PctShare1 = .PctShare` baris yang diubah) | `models.TotalShareTCO` (apd) + `totalShare` di jawaban daftar |
| gerbang simpan | `SaveTreatyReinsurerDetail1_Act`: langkah 3 b694 `"Data tidak boleh kosong...!!!"` bila ReinsurerID kosong; langkah 7 b1357 `"Persentase tidak boleh lebih dari 100!"` kecuali `TotalShare + (PctShare - PctShare1) <= 100.000` (b1452); langkah 9 RDB `SaveMasterTreatyReinsurer_SQL` hanya bila kedua syarat itu lolos; langkah 4, 5, 8, 10–15 DIKOMENTARI | 422 dengan teks VERBATIM; total dihitung di dalam transaksi dengan kontrak dikunci |
| rentang dan koma | `SetErrorMessageReinsurer.xml` langkah 1 `@replaceAll(PctShare/Ricomm, ",", ".")`; langkah 2–3 `SetErrorMessageBetween` (b452/b631) bila `> 100 || < 0` | `models.UraiPersenMasukTCO` sekali di batas masukan: koma → titik, tolak koma+titik, ≤ 8 desimal, 0..100 |
| form | `ViewDetailTreatyReinsurerGrid1.xml` b7842 `ID`, b8042 `Reins.ID`, b8226 `Reinsurer` (pemilih `BrowseAgentReinsSOA_RD`: `.ID` → ReinsurerID, `.ClientName` → NAME, `.ClientID` → CLIENTID), b8522 `%Share`, b8800 `%Comm`, b9076 `Rating`, b11100 `Operator Name`, b11405 `Save`; delapan medan tanpa label (TreatyYear, TreatyGroupID, ReinsTypeID, UserId, IUDate, StartDate, EndDate, StatusOn, CLIENTID) tersembunyi permanen `pyCondition 1=2` | form panel dengan lima medan tampil; medan tersembunyi dipertahankan server |
| master reinsurer | `BrowseAgentReinsSOA_RD.xml` kelas `ASM-FW-GISFW-Int-AGENT`, `.ClientName Contains` tanpa beda huruf (b569–b572), `.StatusActive = 1` (b579–b591); kolom fisik `ID`/`CLIENTNAME`/`CLIENTID` terbukti di `Claim Fac In/RDBList/GetLeaderReport.xml` dan `GetAddressCeding.xml` | `GET /api/treaty-contract-out/reinsurer-master?cari=` (dibaca saja) |
| grid | b1730 `Add`, b1996 `Tambah` (aksi sama), kolom b2418 `ReinsID` · b2560 `Reinsurer` · b2702 `%Share` · b2844 `%Comm` · b2990 `Rating` · b3138 `Operator Name`; b4491 `Edit`; b4936 `Delete` → `DeleteTreatyReins_Act` (tiket 10); b5277 `Security Reinsurer` (tiket 06); b6186 `Total Share -->>` (b6330 nilai) | grid + kaki total; `Tambah` tidak dibawa sebagai tombol kedua |

### Ralat bertanggal 29-09-2026

1. **"Total share ≠ 100 bukan gerbang" hanya separuh benar.** Kurang dari 100 memang tidak ditolak. Lebih dari 100
   DITOLAK di Pega (`SaveTreatyReinsurerDetail1_Act` langkah 7 + 9, b1357/b1452) dengan pesan `Persentase tidak boleh
   lebih dari 100!` — dibawa sebagai 422. Pemeriksaan dilakukan di dalam transaksi dengan baris kontrak pembuka
   kombinasi dikunci (`FOR UPDATE`), supaya dua penulis serentak tidak sama-sama lolos.
2. **Share dan komisi wajib, masing-masing 0..100.** Rentang dari `SetErrorMessageReinsurer` (hidup di Pega, pada
   perubahan medan); kewajiban mengisi dari AC 14 ("masing-masing dengan share dan komisi") — Pega tidak
   menggerbanginya.
3. **Reinsurer dari master `AGENT`, nama dan client dari master.** Klien hanya mengirim `reinsurerId`; `NAME` dan
   `CLIENTID` diambil server. Kolom `STATUSACTIVE` diturunkan dari nama properti RD `.StatusActive` — **OQ-TCO-12**.
4. **`OPERATORNAME` menyimpan pengenal akun**, bukan `pyUserName` (nama tampilan orang) — nol nama orang di data baru.
   `USERID` = pembuat baris (dipertahankan saat diubah).
5. **Kombinasi dikunci teks `TREATYYEAR`, bukan `ID` tahun** (`BrowseTreatyReinsurerList_Act` langkah 1): dua tahun
   treaty berteks tahun dan grup sama berbagi reinsurer untuk jenis yang sama — dibawa apa adanya
   (`models.KombinasiTCO`).
6. **`STDRATING`** di tabel non-life adalah medan `Rating` biasa (b9076). Pemakaian ulangnya sebagai total share ada di
   varian Life (`SetTreatyReinsurerList_Act` langkah 4 `InputTreatyReinsurerLife.STDRATING = dPctShare`) dan tidak
   menyentuh tabel ini.

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| models | `tco_reinsurer.go`, `tco_kombinasi.go` (+uji) | `UraiPersenMasukTCO`, `TotalShareTCO`, `PeriksaTotalShareTCO`, `PeriksaReinsurerTCO`; uji uang dibuktikan merah terhadap float64 dulu |
| repository | `tco_reinsurer.go` (+uji), `tco_kontrak.go` (+`Kunci`) | daftar per kombinasi `ID ASC`, share lain `FOR UPDATE`, sisip/perbarui berbatas kombinasi; master `AGENT` dengan LIKE ber-ESCAPE |
| services | `tco_reinsurer.go` (+uji) | `ReinsurerTCO`: kombinasi dari tahun + kontrak; satu transaksi {kunci kontrak, total ≤ 100, tulis, jejak} |
| handlers | `tco_reinsurer.go` (+uji, +uji `db`) | 4 rute |
| frontend | `PanelReinsurerKombinasi.tsx` (+uji), `REINSURER_TCO` (21 baris diuji ke korpus), `api.ts` (+3), tombol `Reinsurer List` hidup | uang teks sepanjang jalan |

**Status:** selesai 29-09-2026 — commit `treaty-contract-out: tiket 05 — reinsurer + total share`.

## Keputusan work owner 29-09-2026

- **OQ-TCO-12 — ditutup.** Jawaban: *"benar"*. Kolom fisik `STATUSACTIVE` master `AGENT` (nama properti RD
  `BrowseAgentReinsSOA_RD`) dikonfirmasi `[keputusan work owner 29-09-2026]`.

## Ralat bertanggal 29-09-2026 — tco4 (nol tabel baru) `[keputusan work owner]`

- Tabel `TREATYREINSURER` warisan, sequence `M_TREATYREINSURER_SEQ`; `RICOMM`/`PCTSHARE` NUMBER, `STARTDATE`/
  `ENDDATE`/`TGLUPDATE` VARCHAR2 (stempel Pega). `USERID`/`TGLUPDATE` diisi layanan walau Pega mengosongkannya
  (`NewTreatyReinsurerDetail_Act` b917) — penyimpangan sadar kecil, **OQ-TCO-25**.
- **Jejak gugur**.
