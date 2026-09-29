# Paritas layar dan aksi — Treaty Contract Out

> Disusun 28-09-2026 (sesi modul, tiket 01 →). Satu baris per tombol/aksi korpus → rute/kontrol sistem
> baru. Nomor baris = `sed -e 's/></>\n</g'` atas berkas korpus `D:\XML\RNM_BRD\Treaty Contract Out\`.
> Keadaan: ✅ dibangun · ⚠️ belum · ➖ sengaja tidak dibawa (kode mati / penyimpangan sadar) · 🔜 tiket berikut.

## Menu — tiga harness portal (kelas `Data-Portal`)

| Korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `Harness/InboxTreatyContract.xml` b151 `<pyLabel>InboxTreatyContract</pyLabel>` → section `GridTreatyContract` (b931) → `InputTreatyContract` (`GridTreatyContract.xml` b2482) | sidebar **Treaty Contract Out → InboxTreatyContract** | 🔜 tiket 03 |
| `Harness/InboxTreatyContractReinsType.xml` b151 `<pyLabel>InboxTreatyContractReinsType</pyLabel>` → judul `ReinsType` b1670 → section `PanggilReinsType` (b1825) → `InputTreatyContractReinsType` (`PanggilReinsType.xml` b1300) | sidebar **Treaty Contract Out → InboxTreatyContractReinsType** | ✅ tiket 04 (pemilih tahun + editor) |
| `Harness/InboxTreatyContractDescription.xml` b359 `<pyLabel>InboxTreatyContractDescription</pyLabel>` → `NitipKurs` b3882, `SubViewDetailDescription` b11366, `ViewDetailDescriptionProp` b12583, `ViewDetailDescriptionNonProp` b13547 | sidebar **Treaty Contract Out → InboxTreatyContractDescription** | 🔜 tiket 08 |

⚠️ Kelompok **Treaty Contract Out** BELUM ada di sidebar (`labels.ts` `MODUL` memuat 17 kelompok; folder
korpus 20 — ralat §8 `PROMPT-EKSEKUSI-HULU-HILIR.md`). Ditambahkan ADITIF bersama layar pertama (tiket 03).

## Tiket 01 — skema + migrasi data (tanpa layar)

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `RDBList/SaveMasterTreatyYear_SQL.xml` → `POOLDATA.PEGA_TREATYYEAR` (10 param) | tabel `T_TREATYYEAR` (migrasi 300), sequence `SEQ_T_TREATYYEAR` | ✅ skema; penulisnya 🔜 tiket 03 |
| `RDBList/SaveMasterTreatyContract_SQL.xml` → `PEGA_TREATYCONTRACT` (8 param) | `T_TREATYCONTRACT` (301) | ✅ skema; 🔜 tiket 04 |
| `RDBList/SaveMasterTreatyReinsurer_SQL.xml` → `PEGA_TREATYREINSURER` (19 param) | `T_TREATYREINSURER` (302) | ✅ skema; 🔜 tiket 05 |
| `RDBList/InsertToMTreatySecurity.xml` INSERT posisional 7 nilai, `UpdateMTreatySecurity.xml` kunci `trim()` | `T_MTREATYSECURITY` (303) PK surrogate, kolom bernama, tanpa `trim()` | ✅ skema; 🔜 tiket 06 |
| `RDBList/SaveMasterTreatyBusiness_SQL.xml` → `PEGA_TREATYBUSINESS` (12 param) | `T_TREATYBUSINESS` (304) | ✅ skema; 🔜 tiket 07 |
| `RDBList/SaveMasterProportionalArrg.xml` (35) + `SaveMasterProportionalArrgChild.xml` (26) → satu tabel | `T_PROPORTIONALARRG` (305) — 35 kolom, anak NULL di 9 kolom induk | ✅ skema; 🔜 tiket 08 |
| — (tidak ada di korpus; ADR-0007) | `T_TREATYCO_JEJAK` (306) | ✅ skema; penulis 🔜 tiket 03+ |
| enam tabel warisan `POOLDATA.TREATYYEAR` … `PROPORTIONALARRG` | `-migrate-data-treaty-contract-out` → `repository.MigrasiTCO.Pindahkan` (baca teks → konversi → tulis → rekonsiliasi tepat → sequence) | ✅ kode + uji `db` (SKIP tanpa skema uji); ⛔ tidak dijalankan di DEV (tco2) |
| kueri hilir `Claim Prop/RDBList/GetLimitPLATreatyin.xml`, `GetListRetro_Sql.xml`, `GetTreatyGroupID.xml`; `Claim Fac In/RDBList/GetLimitPLADLA_Sql.xml`, `GetQuotaShare.xml`, `GetTreatyGroup_Sql.xml`, `GetDataTreatyLimit_Sql.xml`; `Komite Claim Prop/RDBList/GetListRetro_Sql.xml` | `repository.KontrakHilirTCO` — 6 pembaca read-only berkolom VERBATIM; uji tiga sisi (DDL, korpus, nol tulis) | ✅ tco3 |
| `M_PROPORTIONALARRG`, `M_TREATYCONTRACT`, `M_TREATYYEAR`, `M_TREATYBUSINESS` + seluruh kueri `FROM m_*` | — | ➖ MATI (penyimpangan sadar 1); penjaga `TestTCONolTabelDokumenWarisan` |
| `POOLDATA.PROSESCOPY`, `RDBList/SaveMasterCopyData_SQL.xml`, `Activity/BrowseCopyData.xml`, `SaveTreatyYearMultiple_Act.xml`, `NewInputTreatyYearMultiple_Act.xml`, bagian salin `BrowseDeleteRowTreatyInContract.xml`; form From/To `InputTreatyContract.xml` b2374/b3818, tombol `Proces` b5104 → `BrowseCopyData` b5128 | — | ➖ fitur salin DIBUANG (penyimpangan sadar 3, AC 72) |
| kolom warisan `PROPORTIONALLIST`, `OBJECT` | dicacah laporan migrasi (`KolomMatiBerisi`), tidak dibawa | ➖ AC 70 |

## Tiket 02 — jenis reasuransi: master dibaca + saringan non-life

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `ReportDefinition/BrowseReinsuranceType_RD_Old_Ljt_id_isnotnull.xml` — `.ID NotStartsWith` 12 nilai b573/b581, `.Flag="active"` b586, `.Type` 1/2/3 b603, sort `.Note` b667; dipakai 11 grid klausul | `GET /api/treaty-contract-out/jenis-reasuransi` → `repository.MasterJenisReasuransi.DaftarNonLife` (`FLAG=`, `TYPE IN`, `ID NOT LIKE` ×12, `ORDER BY NOTE`) | ✅ nama jujur (nol "Old") |
| master `POOLDATA.REINSURANCETYPE` (`ID`, `NOTE`, `TYPE`, `FLAG`) | dibaca saja; penjaga `TestTCOWarisanHanyaDibaca` (`masterDibacaSajaTCO`) | ✅ AC 12 |
| master kosong / tersaring habis | 503 `ErrMasterJenisReasuransiKosong` menyebut `REINSURANCETYPE` | ✅ ADR-0015 |
| pemilih `ReinsType` (`InputTreatyContractReinsType.xml` b2652) / `Reinsurance Type` (`InputTreatyContract.xml` b7925) | `components/treaty-contract-out/PilihJenisReasuransi.tsx` (`Pilih` ui/dasar; nilai `.ID` teks, label `.Note`) | ✅ komponen; dipasang di layar 🔜 tiket 03/04/08 |
| RD non-Old `BrowseReinsuranceType_RD.xml` (param `Flag="active"` b2768, pemilih form kontrak) | daftar tersaring yang sama (AC 5 tiket 02) | ⚠️ OQ-TCO-06 |
| `RDBList/GetMasterReinsTypeContract.xml` (`a.JSONDATA`) | — | ➖ MATI |

## Tiket 03 — layar tahun treaty (`InboxTreatyContract` → `GridTreatyContract` → `InputTreatyContract` + `InputDtlTreatyContact`)

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| menu `Harness/InboxTreatyContract.xml` b151 | sidebar **Treaty Contract Out → InboxTreatyContract** (`labels.ts` `MODUL.treatyContractOut` +1, `Shell.tsx` kelompok ke-18, `daftarMenu.ts` `tco-tahun`) | ✅ keterangan "belum dimigrasi" dicabut untuk kelompok ini |
| judul `TREATY CONTRACT OUT` (`GridTreatyContract.xml` b1057) | `TAHUN_TCO.judul` | ✅ |
| grid `BrowseTreatyYear_RD` b17355: `Underwriting Year` b17378 / `Transaction Year` b17531 / `StartDate` b17684 / `EndDate` b17837 / `Treaty Group` b17990 / `Reinsurance Type` b18136 (sel b18964/b19125/b19286/b19446/b19605/b19724); sort `.ID DESC` b672 | `GET /api/treaty-contract-out/tahun` → tabel enam kolom, `Halaman` 20/halaman | ✅ (OQ-TCO-05 label bersilang, dibawa apa adanya) |
| `Add` b16387 → `NewInputTreatyYear_Act` | tombol `Add` → `formKosong()` | ✅ |
| `Edit` b19939 → `SetTreatyYear_Act` b20030 | tombol `Edit` → `formDari(baris)` | ✅ |
| `ReinsType` b20778 → harness 2 | tombol membuka `PanelKontrakTahun` untuk baris itu | ✅ tiket 04 |
| `List Description` b22196 → harness 3 | tombol berdiri `disabled` "menunggu tiket 08" | 🔜 tiket 08 |
| `Copy` b20459 + `From`/`To` b2374/b3818 + `Proces` b5104 → `BrowseCopyData` | — | ➖ AC 72 |
| form `Input New Data` (`InputDtlTreatyContact.xml` b5437): `ID` b6379 · `Treaty Group` b6560 (`BrowseTreatyGroup_RD` b6624) · `Reinsurance Type` b6800 → `.Proportion` b6829 · `Start Date` b7532 · `End Date` b7816 · `Underwriting Year` b8004 → `.TreatyYear` · `Transaction Year` b8284 → `.UnderwritingYear` · `Modified Date` b9097 · `Username` b9282 | `Field`/`FieldTanggal`/`Pilih` (grup dari `GET /grup-treaty`) / `PilihJenisReasuransi`; ID, Modified Date, Username hanya dibaca | ✅ |
| `Save` b10332 → `SaveTreatyYear_Act` (b280 UserID, b327 TglUpdate, prasyarat b388/b411/b434, RDB `SaveMasterTreatyYear_SQL` → `PEGA_TREATYYEAR`) | `POST /tahun` (baru, ID dari `SEQ_T_TREATYYEAR`) / `PUT /tahun/{id}` (seluruh medan) → `T_TREATYYEAR` + jejak `T_TREATYCO_JEJAK`, satu transaksi | ✅ logika ditiru, procedure tidak dipanggil |
| `Cancel` b10622 → `CancelActivityTreatyContract` | tombol `Cancel` menutup form | ✅ |
| `CheckYear` b335 `isNumber(TreatyYear)` (`pyMessageLabel CheckYearly` — teks tidak diekspor) | `ErrTahunTreatyBukanAngka` → 422 | ✅ (teks pesan kosakata kami) |
| — (tidak ada di Pega) | gerbang periode terbalik (AC 9) → 422; anti-dobel (AC 73) → 409 menyebut ID baris lain | ✅ tambahan sadar |
| form kedua `Input New Data` `InputTreatyContract.xml` b6351 (tanpa Start/End/UW Year; `ReinsuranceType` b7954 bukan parameter procedure) | — | ⚠️ `[terbuka]` mana yang tampil di Pega; form lengkap yang dibangun |
| `GridTreatyArrangementAttachment` (`Attachment for` b11721) | `PanelLampiranTahun` di form tahun ber-ID | ✅ tiket 12 |

## Tiket 12 — lampiran tahun treaty (`GridTreatyArrangementAttachment`, FITUR BARU)

Nomor baris = `Section/GridTreatyArrangementAttachment.xml` kecuali disebut lain.

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| `InputTreatyContract.xml` b11721 `Attachment for`, b13074 include panel | `PanelLampiranTahun` di form tahun treaty yang sudah ber-ID; tahun baru menampilkan catatan "simpan dulu" | ✅ |
| b1785 `For Treaty Contract Out` | subjudul panel | ✅ |
| `Add attachment` b578 → `SetCategory_act` b596 → `TreatyOutAttachContent` b642 (`pyAttachmentScreen`) → `TreatyOutSaveAttachment` → `InsertAtatchment_Sql` → `PEGA_M_ATTACHMENT` | pemilih `Type` (master `CATEGORY_ATTACH_REAS`) + kotak berkas → `POST /api/treaty-contract-out/tahun/{id}/lampiran` (multipart) → rekam `T_TREATYYEAR_LAMPIRAN` + efek outbox `storage-unggah` + jejak, satu transaksi | ✅ prosedur tidak dipanggil |
| `Refresh` b1023 → `LoadAttachmentTreatyOut` | `GET /tahun/{id}/lampiran` | ✅ |
| `Download All` b2659 → `TreatyOutDownloadAll_Act` | `GET /tahun/{id}/lampiran/semua` (zip) lewat `fetch` berheader identitas | ✅ |
| `Download` b2391 → `DownloadAll_Act` | — | ➖ tombol kedua untuk aksi yang sama |
| sel `.pyFileName` b3428 → `TreatyOutDownloadOne` b3488 | tombol nama berkas → `GET /tahun/{id}/lampiran/{lid}/isi` lewat `fetch` berheader identitas | ✅ |
| kolom `File Name` b3032 · `Type` b3170 (`.pyCategory` b3705) | kolom tabel panel | ✅ |
| `Delete` b3897 → `DeleteAttachmentTreaty` → `DeleteAttachment2_Sql` | `DELETE /tahun/{id}/lampiran/{lid}` → rekam + efek `storage-hapus` + jejak; berkas yang sudah tidak ada tidak menggagalkan | ✅ |
| `TreatyOutSaveAttachment.xml` b376 `Tidak ada file yg diattach` | 400 dengan teks VERBATIM | ✅ |
| `GetAllAttachment2_Sql` / `GetAttachment2_Sql` (`M_ATTACHMENTTREATY_2 where treatyid = {TreatyIn.ID}`) | — | ➖ kunci treaty inward tidak dibawa (penjaga Go + JS) |
| — (tidak ada di Pega) | status terkirim / tertunda / gagal + galat terakhir; `Ulangi` (`POST …/{lid}/ulangi`); `Periksa keselarasan` (`GET …/selaras`) | ✅ tambahan AC 55, 58, 61 |
| `ConnectREST/ServiceGoogle.xml`, `LinkService`, `GetTokenStorage_SQL` | `PenyimpananJarakJauhTCO` (resolver `M_LINK_SERVICE` saat jalan, `CacheTokenTCO`), transport tidak diimplementasikan → `ErrPenyimpananBelumDisetujui`; yang dipasang di DEV adalah stub lokal `PenyimpananLokalTCO` | ⚠️ penyambungan nyata menuntut persetujuan manusia |

## Tiket 04 — kontrak treaty di dalam tahun (`InputTreatyContractReinsType`)

Nomor baris = `Section/InputTreatyContractReinsType.xml` kecuali disebut lain.

| Unsur korpus | Sistem baru | Keadaan |
| --- | --- | --- |
| harness `InboxTreatyContractReinsType` b151, judul `ReinsType` b1670 | butir menu kedua → `InboxTreatyContractReinsType.tsx` (pemilih tahun, `[tidak ada di korpus]`) → `PanelKontrakTahun` | ✅ |
| popup dari `InputTreatyContract.xml` b20778 (`BrowseReinsTypeYear` + `showHarness` + `KirimTahunGroupID`) | tombol `ReinsType` baris tahun → `PanelKontrakTahun` tahun itu | ✅ |
| kepala `Underwriting Year` b1145 · `ReinsType` b1358 (nilai `.TreatyGroupName` b1386) | dua medan baca-saja + catatan label bersilang | ✅ |
| `ReinsType` b2652 (RD non-Old, `Flag "active"` b2768) | `PilihJenisReasuransi` (daftar tersaring tiket 02) | ✅ OQ-TCO-06 tetap terbuka |
| `Start Date` b2905 → `SetTanggalTreatyContract` b3007 | isi tanggal mulai → `GET /tahun/{id}/kontrak/akhir-bawaan` mengisi tanggal akhir | ✅ OQ-TCO-10 |
| `End Date` b3244 (`IsEndDate==1`) | ubah tanggal akhir saja | ✅ |
| `Save` b3618 → `SaveTreatyContract_Act` → `SaveMasterTreatyContract_SQL` → `PEGA_TREATYCONTRACT` | `POST /tahun/{id}/kontrak` / `PUT /tahun/{id}/kontrak/{kid}` → `T_TREATYCONTRACT` + jejak, satu transaksi | ✅ prosedur tidak dipanggil |
| `Undo` b5343 → `UndoOperation` | kembalikan isian terakhir yang dimuat | ✅ |
| `Information` b6400 (`OutputData.HASIL1`) | baris status sesudah simpan | ✅ |
| `Modified Date` b4363 · `Username` b4547 · ID b2478 | baca-saja (Username = penulis terakhir rekam) | ✅ ralat 7 |
| `Add` b8528 → `NewInputTreatyContract_Act` | form kosong | ✅ |
| grid `BrowseTreatyContract_RD` (tidak diekspor): `Reins Type` b9164 · `Treaty Start` b9304 · `Treaty End` b9444 | `GET /tahun/{id}/kontrak`, `ID DESC` | ✅ |
| `Edit` b10519 → `SetUbahTreatyContract` | form dari baris | ✅ |
| `Business List` b10842 | tombol berdiri `disabled` | 🔜 tiket 07 |
| `Reinsurer List` b11308 | tombol berdiri `disabled` | 🔜 tiket 05 |
| `Delete` b11809 → `BrowseDeleteRowTreatyInContract` | tombol berdiri `disabled` | 🔜 tiket 10 |
| `ViewDetailTreatyReinsurerGrid1` b13311, `ViewDetailTreatyBusinessGrid` b14064, grid security `SelectSecurityReinsurer` b16069 | — | 🔜 tiket 05, 07, 06 |
| gerbang `SaveTreatyContract_Act` langkah 2–6, 12 (DIKOMENTARI) | jenis wajib dari daftar, tanggal wajib, periode (AC 9), dobel 409 `Data sudah pernah di Input` | ✅ ralat 1–2, OQ-TCO-11 |
| gerbang `SetTanggalTreatyContract` langkah 4 (HIDUP) | tahun mulai = tahun treaty → 422 | ✅ |

