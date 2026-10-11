# PARITAS — tombol dan aksi XML Claim Fac In ↔ status

> 10-10-2026, tahap 1 dari 2 (prompt "IMPLEMENTASI CLAIM FAC IN, TAHAP 1 DARI 2"). Sumber inventaris: flow
> `Register_Flow`, Section `InputRegister`, `InputRegisterDetail`, `Catastrope_Sec`, `ViewHistoryClaim`,
> `InputInwardFacultativeDtl`, `ProgresClaim_Sec`, `SubProgresClaim_Sec`, `ViewPolis` (harness ChoosePolis), `ProtectDOL`,
> `ClaimComiteeReject` + `SureRejectClaim_section`, `InputEstimasiAdmin`, `InputEstimasiDetail`,
> `PropertyItemListGridEstimation*`, `ShowItemPA`, `Estimasi`, `EstimasiPA`, `EstimasiMarine`, `Pla_Dtl`, `ClaimSurvey`,
> `InputEstimasi`, `ShowObjectAdj`, `ItemListEstimation*`, `Adjusment_SC`, `InputAdjustment`, `ViewCedantPanel`,
> `ClaimComite`, `PrintDLA_dtl`, `PreventRejectClaim` di `D:\XML\RNM_BRD\Claim Fac In` (setiap pxButton / pxLink dan setiap
> medan ber-`on change` activity). Sumber status: definisi layar server `backend/models/layar_register.go`,
> `layar_estimasi.go`, `layar_adjustment.go`, `komite.go`, `tutup.go`, `dla.go` — satu-satunya tempat tombol lahir;
> frontend merender pohon tata apa adanya (`frontend/components/TataView.tsx`). Peta aksi server:
> `backend/services/aksi.go` (`penangan`, 82 aksi) dan `services/adjustment.go` (`penanganAdjustment`).
>
> **Status:**
> - **dibangun** — tampil menurut kondisi XML yang sama (visible-when / disabled-when / read-only-when), aksinya port
>   activity-nya.
> - **nonaktif-OQ** — tampil sesuai section, tetapi nonaktif; keterangan `OQ-CFI-nn` di atribut `title` (`OQ.md`).
> - **tidak tampil di XML** — kondisi tampilnya statis `1=2` / `NEVER` / tidak pernah benar: tidak dibangun.
> - **tahap 2** — milik modul `komiteclaimfacin` (prompt §2 "Di luar lingkup").
>
> Penanda: `[penyimpangan sadar]` = perilaku sengaja berbeda dari XML (alasan + perilaku Pega lama dicatat);
> `[inferensi]` = disimpulkan, tidak tertulis di XML.

## 1. Flow dan halaman awal

| XML | Aksi XML | Status |
| --- | --- | --- |
| `Register_Flow` Start1 → Decision4 B2B (`When/IsSPK`: `OfferFacIn.IsB2B = "SPK"`) | jalur B2B langsung ke Choose Surveyor | **dibangun sebagian** — kasus lahir lewat **Add Claim** tanpa polis, sehingga `IsSPK` salah dan kasus selalu masuk Input Register. Kasus B2B lahir dari harness `New` / `NewSample` yang tidak diekspor (OQ-CFI-29) |
| Assignment1 "Input Register" (operator saat ini, FlowAction `InputRegister`) | `TAHAP = InputRegister` | **dibangun** — tab *Process* (worklist pembuat) |
| Assignment7 "Input Estimasi" (WorkList Custom = pembuat, FlowAction `InputEstimasi`) | `TAHAP = InputEstimasi` | **dibangun** |
| Assignment3 "Choose Surveyor" (WorkBasket Custom, FlowAction `InputSurveyor` berlabel "Input Adjustment") | `TAHAP = InputSurveyor` | **dibangun** — workbasket `ReasKlaimTeknik` (bawaan b: `ReasPNCTeknik` tidak ada di DEV); switch **Teknik** di halaman awal bagi anggota workbasket (`GET /hak`) |
| Decision5 / Decision8 `IsBackStage` (`.pyNote = "Back"`) | Back Input Estimasi → Input Register; Back Choose Surveyor → Input Estimasi | **dibangun** (§3, §4) |
| END52 Resolved-Completed | `CloseClaim` (`ASMForceCaseClose`) | **dibangun** — tab *Resolve* |
| Pra-proses FlowAction `InputRegister` (`CallActivityInputRegister` + DT `InsertObjects_dt`) | — | **dibangun** (`models.PraRegister`, `objek.go`) |
| Pra-proses `InputEstimasi` (`InputEstimationPre` + DT `SetEstimation_DT`) | — | **dibangun** (`models.PraEstimasi`) |
| Pra-proses `InputSurveyor` (`SetTypePDFAdjustment` + DT `SetStartDateAdjustment`) | — | **dibangun** (`models.PraAdjustment`; StartDateAdjustment hanya bila kosong, §8) |
| `InsertProgressClaim` (setiap assignment dibuka) | `PROGRESSCLAIM` / `SUBPROGRESSCLAIM` | **dibangun** — idempoten per (kasus, posisi), isi `PEGA_PROGRESSCLAIM` / `PEGA_SUBPROGRESSCLAIM` ditulis ulang tanpa memanggil prosedur |
| Validate FlowAction `ValidateDate` | — | **tidak dibangun** — rule tidak diekspor (OQ-CFI-13); pemeriksaan tanggal `CheckDate*` tetap jalan |
| Kotak masuk Beranda | — | **tidak dibangun** — bukan bagian Pega (prompt §2 "Di luar lingkup"; `menu.ts` tanpa `antreanBeranda`) |
| Tabel komite di bawah inbox | — | **tahap 2** |

## 2. Input Register — Section `InputRegister` / `InputRegisterDetail`

`InputRegisterDetail` juga tampil hanya-baca sebagai tab **View Registration** (Input Estimasi, `IsTreatyIn != 1`) dan
**Registration** (Choose Surveyor); medan terkunci menurut kondisi section yang sama (`isEstimation`, `IsAnyAccept`).

**Rupa 11-10-2026 (screenshot Pega work owner, "UBAH LAYOUTNYA BEGINI"):** kontainer 650 (sel 652-673: Policy No +
Choose Polis | Insured Name, QQ Name | Ceding Co Name, Source Of Business | Date of Loss, Report Date | Received Date,
Reporter Name | Reporter Phone Number, Choose Cause of Loss + Cause of Loss | Catastrophe), 674 (676 PLA No Ceding /
Report Type / PLA No SOB | 681 Reporter Status + Specify... + grid ceding) dan 713 (Description Report | Reporter
Address) = dua kolom diisi per baris, label DI ATAS isian (letak `duaAtas`). Sebelumnya sel 650 ditumpuk di kolom kiri
dan 674 di kolom kanan, berlabel kiri. Blok sesudahnya (grid objek, Consultant / Adjuster, Location, Country ...,
Currency / Estimate) tidak ada di screenshot - rupanya tidak diubah.

| Label XML | Aksi XML | Status |
| --- | --- | --- |
| View Status Payment Premi | REST `getPremiumPaidOn` (M_LINK_SERVICE) | **nonaktif-OQ** OQ-CFI-14 |
| Reject Claim | local action `RejectSurveyClaim` (pra `SetRejectClaim_pre`) | **dibangun** — pop-up `ClaimComiteeReject` (`BukaRejectClaim`, modal `tolak`); tombol **Yes** = `SendRejectClaimToKomite2` **dibangun** tahap 2 (§5) |
| Tab Register / Policy Detail & Claims History / Progress Claim | — | **dibangun** — tab **Progress Claim disembunyikan** di ketiga layar kasus (work owner 11-10-2026 "HIDE AJA! GA KEPAKE JUGA", `TabProgresTampil`); PROGRESSCLAIM tetap ditulis |
| Save (VIS `PolicyNo` kosong) | click:save | **dibangun** (`Simpan`) — **11-10-2026**: SELALU tampil di samping Submit (work owner "SAVE NYA PINDAH KE SAMPING SUBMIT"); tombol Save dan strip tab "Lampiran" panel lampiran dibuang |
| Submit (VIS `IsError = 1`) | local action `ProtectDOL` | **dibangun** (`BukaProtectDOL`, modal `protectDOL`) |
| Submit (VIS `IsError = 0`) | `CheckListEstimasi_Act` + `ProteksiDataRegister_Act` + finishAssignment | **dibangun** (`Submit`) → Input Estimasi |
| Choose Polis (NA `isEstimation = 1`) | harness `ChoosePolis` (Section `ViewPolis`) | **dibangun** (`BukaPolis`, modal `pilihPolis`) |
| Date of Loss / Report Date / Received Date (change) | `CheckDate_Act` / `CheckDateReport_Act` / `CheckDateReceived_Act` | **dibangun** — langkah 9/13 `+1` atas teks tanggal = OQ-CFI-12 (dibangun menurut hasil XML) |
| Choose Cause of Loss | harness `CauseofLoss_Harness` → `GetNameCauseofLoss` | **dibangun** — pop-up daftar `pilihan/sebab` (`BukaSebab`), pilih = `GetNameCauseofLoss`. Master Cause of Loss milik modul lain, tidak diubah |
| Catastrophe (radio) / Non-Catastrophe type (change) | `SetDefNonCatastrope` | **dibangun** |
| Catastrophe List | harness `CatastrofeList` | **dibangun** — pop-up daftar `pilihan/katastrofe` (`BukaKatastrofe`); pilih `SetCatastrope`, input `InputCatastrope`, simpan `SaveCatasrtope` (tabel `CATASTROPHE`) |
| ikon Edit / Save Catastrophe | `SetEditCatastrope` Edit / Save | **dibangun** |
| Reporter Status (change) | `GetReportStatus_Act` | **dibangun** — label 1 Ceding Co Name / 2 Broker Name / 3 Others = prompt values property `InsuredRelationship` (screenshot work owner 11-10-2026) |
| grid Ceding Co (LS132) — Choose | `GetCeding_act` | **dibangun** — **ralat 11-10-2026** ("TOMBOL INI GA FUNGSI"): GetLeaderReport mengambil agen baris terakhir; DEV dua agen bernama klien sama, yang terakhir tanpa `CLIENT_ADDRESS` (alamat hanya di JSON M_CLIENT yang tidak dipakai) -> Reporter Address kosong. Agen yang kliennya beralamat didahulukan (`sqlAgenKlien`); DEV: alamat terisi. **11-10-2026** ("PERBAIKI CHOOSENYA INI!"): GetCeding_act 6 / GetReportStatus_Act 13 merangkai empat bagian alamat apa adanya - klien tanpa CLIENT_ADDRESS (DEV: EVEREST REINSURANCE COMPANY, 0 alamat) memberi Reporter Address ", , , "; bagian kosong kini dilewati `[penyimpangan sadar]` (`rangkaiAlamat`, `TestAlamatPelaporTanpaPemisahMenggantung`). Reporter Address tetap terkunci kecuali Reporter Status Others (XML RO); kosong tidak menahan Submit (medan wajib hanya-baca tidak ditagih) |
| grid pilih objek (LS137, per lini) — centang | `CheckListEstimasi_Act` | **dibangun** |
| Consultant ID / Adjuster ID (autocomplete) | `SetConsultant` / `SetAdjsuter` | **dibangun** |
| ikon + Consultant / Adjuster | harness `MstAdjusterConsultant` (Data-Portal) | **nonaktif-OQ** OQ-CFI-32 — Claim Prop menyimpan master barunya lewat rute pinjaman `POST /api/adjuster-consultant`, yang di `cmd/api/rakit.go` (`ruteDipinjam`) hanya dipinjamkan ke `claimprop`; memanggilnya dari modul ini juga gagal penjaga `TestPanggilanLintasModulTerdaftar`. Pemilihan adjuster / consultant yang sudah ada tetap jalan |
| Location of Loss (change) | `SetTempLocation` | **dibangun** |
| Country / Province / City / District / Region / Zip Code | `BrowseCountry_RD` … `InputKodePos1` | **dibangun** — saringan provinsi menurut negara diperbaiki (§8). **Logika diperbaiki 11-10-2026** ("PERBAIKI DEPAN DAN BELAKANG. PERBAIKI LOGIC NYA"): Section merujuk RD kelas Int-PROVINCE / Int-CITY (tidak diekspor; yang terekspor bernama sama kelas Int-RW) dan menyaring per ID induk, tetapi RW DEV ber-PROVINCEID / CITYID / DISTRICTID kosong 86.283/86.283 dan master CITY / DISTRICT bernama ganda -> District / Region kosong. Rantai = NAMA di RW (pola BrowseRW_SQL Zip Code): Province (negara) -> City (provinsi) -> District (provinsi + kota) -> Region (provinsi + kota + kecamatan, isi RWID + Zip Code); mengganti tingkat mengosongkan tingkat bawah; tingkat tanpa induk nonaktif berpetunjuk ("Choose Province first"); Country awal INDONESIA (sel 924 `pyDefaultValue`); Zip Code mengisi nama dari RW sendiri. DEV: 38 provinsi, DKI 6 kota, Jakarta Barat 8 kecamatan, kelurahan berkode pos. **Ralat 11-10-2026 sore** ("MASA AKU BISA ISI GINI? LOGIC NYA DIMANA?" - Province BALI + City JAKARTA TIMUR + District CENGKARENG): kiriman layar digabung SEBELUM aksi, sehingga "tingkat berubah?" selalu salah dan tingkat bawah tidak pernah kosong; kini dibandingkan dengan nilai tersimpan (`nilaiTersimpan`), dan nama wilayah TIDAK digabung dari kiriman layar (`BuangWilayahPilihan`) - hanya lewat aksi pemilihan / Zip Code. Uji meniru kiriman layar penuh. **11-10-2026** ("RAPIHKAN!" layar komite): `InputKodePos1_Act` 5 merangkai Location `lokasi sementara + ", " + RW + ...` apa adanya - lokasi kosong menghasilkan ", KEL, ..."; bagian kosong kini dilewati `[penyimpangan sadar]` (`TestInputKodePosLokasiTanpaPemisahMenggantung`). **Ralat 11-10-2026**: kelima kueri memilih kolom nama dua kali tanpa alias + `FETCH FIRST` = ORA-00960 di DEV (daftar gagal); salinan kedua kini beralias, dicoba di DEV (saringan nol baris) + uji penjaga `TestSQLTanpaKolomKembarTakBeralias` |
| Currency (change) | `InputCurrencyValueAct_Register` | **dibangun** |
| Claim Estimate (change) | `CheckDoubleClaim_Act` | **dibangun** — nomor polis tertulis mati dibuang (bawaan c) |
| Policy Detail — View Retro List | harness `RetroList_Harnness` | **dibangun** — pop-up hanya-baca `FacRetroList` dari halaman polis yang dibaca ulang (`BukaRetroList`, OQ-CFI-16) |
| Claim Status (kronologi) / Claim History | `T_VIEW_SUGGEST` / klaim lain polis ini | **dibangun** |
| Progress Claim — Input Progres Claim | FlowAction `InputSubProgressClaim` (`AddDelProgress`) | **nonaktif-OQ** OQ-CFI-15 (hanya mengubah clipboard; sumber Progres1/2 tidak diekspor) |
| Field `ClaimData.UserTeknis` (LS177) | — | **tidak tampil di XML** (VIS `1=2`) |

### Pop-up `ViewPolis` (Choose Polis) dan `ProtectDOL`

| Label XML | Aksi XML | Status |
| --- | --- | --- |
| Search Type / Search | `SearchPolis_act` + `GetPolisForClaim_SQL` | **dibangun** — teks cari **diikat** (§7 butir 1); hasil = `pilihan/polis` (`FACINPRODUCTION` berindeks, batas 500 baris) |
| tautan Policy No (grid hasil) | `CopyNB_Act` | **dibangun** (`CopyNB`, param `nopolis|prodke`) — salin `JSON_POLIS.DATA_JSON` → `OfferFacIn`. **10-10-2026 (perintah work owner)**: DATA_JSON tanpa isi polis (kosong / penunjuk kasus NB Fac In sistem baru) → dokumen berbentuk sama disusun dari tabel datar NB Fac In (`T_GENERAL_POLIS` dan pohon PARENT_ID / SRC_PATH-nya, `models/polis_datar.go`), juga setiap muat (`muatPolis`); §8 butir 25. Tanpa keduanya ditolak terang (OQ-CFI-11) |
| detail polis (`InputParam.CARI4 = 'true'`) | Section `InputInwardFacultativeDtl` | **dibangun** — **10-10-2026 (perintah work owner "IKUTI MODUL NB FACIN. TAMPILAN READ ONLY")**: blok "General" Pega diganti salinan hanya-baca layar Inward Facultative NB Fac In (unsur tata `polisNB`), juga di tab Policy Detail; grid hasil di atas detail ("RAPIHKAN INI"); §8 butir 26 |
| Cancel / Submit | — / `SetInputParam_Act` | **dibangun** — `CheckPeriodePolicy` tidak dibangun: langkah 1 `Exit-Activity` tanpa syarat (OQ-CFI-34) |
| ProtectDOL: grid kasus kembar, Cancel / Submit | `InsertObjectItemList_DT` + `CheckListEstimasi_Act` + finishAssignment | **dibangun** (`Submit`, konteks `protectDOL`) |

## 3. Input Estimasi — Section `InputEstimasiAdmin` / `InputEstimasiDetail`

| Label XML | Aksi XML | Status |
| --- | --- | --- |
| No Claim, View Status Payment Premi | — / REST | **dibangun** / **nonaktif-OQ** OQ-CFI-14 |
| Send to PIC Claim (NA `IsPicTransfer != 1 \|\| IsError > 1`) | `SetDisable_ACT` + finishAssignment (bila `ObjectList(1).IsFacretro != 1` atau `PlaStatus = 1`) | **dibangun** (`SetDisable`) → Choose Surveyor, `POSITION = ReasKlaimTeknik` |
| Tab Estimation / Policy Detail & Claims History / Progress Claim / View Registration | — | **dibangun** |
| Back (VIS `IsCFS != 1`) | DT `BackFromRegister` + `BackToRegister_act` + finishAssignment | **dibangun** — nilai uji langkah 8–11 tidak ditulis (§7 butir 3) |
| grid objek: Download Claim Face Sheet | `CLaimFaceSheet_Act` | **dibangun** — nomor klaim **K**, `OS_AKSEPTASI_KLAIM` STS 0 per estimasi, `JSON_KLAIM`, konversi (efek hanya produksi), `PrintFaceClaim` / `IsCFS`. Berkas PDF = OQ-CFI-20 **Diperiksa ulang 11-10-2026** terhadap `SaveCFS_ACT.xml` + isi procedure `PEGA_JSON_OS_AKSEP_KLAIM` (SYS.ALL_SOURCE DEV): kolom dan syarat sesuai; DATA_JSON Marine Cargo (TradingName / PackingName / Goodsname / ConveyanceName, langkah 12) dan MBU (BrandName / BrandID / ModelName / Model / TypeName / LicensePlate / TypeID, langkah 13) yang dulu tertinggal kini diisi (`TestBarisOSEstimasiIkutSaveCFS`) |
| grid objek: Print PLA (VIS `ClaimData.ExGratia = 0`, NA `IsFacretro != 1 \|\| PlaStatus = 1 \|\| IsError > 1`) | `SetIndexObj_Act` + harness `Pla_Dtl` | **dibangun** (`BukaPLA`, modal `pla:o`) — `PlaStatus` hanya sesudah Submit berhasil (§7 butir 2) |
| Pla_Dtl: Send Email / Message / Share Retro / Remarks / Submit | `GeneratePLA` → `GeneratePLATreaty_Act` / `ObjectItem.GeneratePLA` | **dibangun sebagian** — PLA treaty **G** + revisi `/n` dibangun tetapi gerbangnya `CekLimit.CARI1 = 1` tidak pernah tercapai (OQ-CFI-24); PLA fac retro (`ObjectItem.GeneratePLA`) tidak diekspor (OQ-CFI-21); penerima surel tidak dibangun (OQ-CFI-22). **11-10-2026**: Send Email / Message `change:refresh` kini aksi `PostValue` (Message tampil dan Submit NA dievaluasi ulang seketika; dahulu tanpa aksi) |
| grid objek: Outstanding Summary | `GetAllData_Act` | **dibangun** — pop-up ringkasan `pilihan/outstanding` (dua kolom mata uang diperbaiki, §8). **11-10-2026** ("PERBAIKI TAMPILAN INI"): TSI RNM / Estimation / Kurs / Convert bertumpuk di kolom TSI RNM - kelas angka (display block) terpasang di `<td>`; angka kini di `<span>` dalam sel (juga Retro List), judul kolom angka rata kanan, No Claim / Policy No / Accepted No tidak dilipat (`popupAngka.test.ts`) |
| panel objek: Add (item) | `addRow` | **dibangun** (`TambahItem`); nama item kembar diperiksa saat nama / coverage dipilih (`ProtectionObjectItem_Act`). **11-10-2026** ("INI DI INPUT DOUBLE. TADI ADA WARNING! TAPI SAYA LANJUT INPUT TETEP BISA"): XML hanya memasang pesan bila SATU item lain bernama sama (baru dibuang bila 3+), dan pesan itu hilang di aksi berikutnya - pilihan Object Name / Coverage Fire, Object Id / Coverage Aneka yang menghasilkan item kembar kini DIBATALKAN (422, tidak disimpan, pesan "Can not add Object Item with the same name"; `tolakKembar`), dan item kembar yang terlanjur tersimpan menonaktifkan Download Claim Face Sheet objeknya (`ObjekBerItemKembar` - tidak ada baris OS_AKSEPTASI_KLAIM; work owner "JIKA ADA KEMBAR JANGAN BISA KLIK PRINT CLAIM FACE SHEET") dan menolak Print PLA objeknya (`tahanItemKembar`; XML menonaktifkan tombol hanya bila IsError > 1, satu pasang kembar = IsError 1). Send to PIC Claim TIDAK ditahan ("SAAT DI KLIK SEND PIC DI BUANG AJA DEH"). **MBU / PA / Travel** (11-10-2026 "PERBAIKI JUGA UNTUK BISNIS LAIN!"): Coverage ID yang sudah dipakai item lain objek itu kini juga DIBATALKAN (`PilihCoverageObjek` lewat `tolakKembar`) - `ProtectCoverage_Act` membuang item TERAKHIR ber-coverage sama, sehingga mengganti coverage baris 1 ke coverage baris 2 menghapus item baris 2; item kembar tersimpan (kriteria DeleteTest: Coverage ID sama) menahan Face Sheet / PLA seperti Fire / Aneka (`TestCoverageKembarPADitolak`). Marine Cargo tanpa Add item |
| panel item: rupa | — | **11-10-2026** ("tolong rapihin. buat lebih user frendly"): total estimasi dan total spreading = kaki grid-nya (bilah ringkasan, nilai kosong "—"), Deductible Info satu baris berlabel, grid kosong "No items", panel baris dan kaki grid selebar area terlihat (tabel lebar menggulir sendiri). Medan / label / kondisi tetap section Estimasi. **Semua lini** (11-10-2026 "PERBAIKI JUGA UNTUK BISNIS LAIN!"): EstimasiMarine kini berkaki empat total yang sama; footer XML "Total" tanpa label (TotalClaimSpreaded, TotalEstimationValueinIDR) mengulang dua di antaranya - dilebur `[penyimpangan sadar]`. Uji `TestPanelItemEstimasiTotalDiKakiGrid` per lini (Fire, Aneka, Golf, MBU, PA, Travel, Marine) |
| panel objek: Delete item | `HapusItem` | **dibangun** |
| panel objek: Object Name / Occupation / Coverage / Object Id / Coverage ID (change) | `SetObjectItem` / `GetOccupation` / `GetCoverage*` / `GetCoverageAneka_Act` / `GetCurencyCoverage_Act` | **dibangun** (`SetObjectItem`, `PilihOkupasi`, `PilihCoverageFire`, `PilihAneka`, `PilihCoverageAneka`, `PilihCoverageObjek`). **Ralat 11-10-2026** ("KENAPA MASIH BISA EDIT SETELAH KLIK PRINT CLAIM FACESHEET?!"): Coverage Name Fire terlewat RO `.PrintFaceClaim = 1` dan NA `IsAdjustment = 1 \|\| .IsFacretro = 1 \|\| IsError > 1` (port, bukan XML); Coverage Name Aneka / Golf terlewat NA `IsError > 1` (`InputParam.CARI39` = Send to PIC, sesudahnya IsAdjustment 1). Server menolak aksi sel terkunci (`AksiTerbuka`) - `TestCoverageItemTerkunciSesudahFaceSheet`, `TestCoverageFireDitolakSesudahFaceSheet`. Currency Value(IDR) / TSI RNM in IDR grid item Fire (XML Text Input) kini hanya-baca seperti lini lain `[penyimpangan sadar]` (work owner "INI READ ONLY. BUAT FIRE ANEKA DAN BISNIS LAIN NYA!", `TestKursDanTSIRNMIDRGridItemHanyaBaca`) |
| panel item MBU (EditAction `PropertyItemListGridEstimationMBU_FA`) | merender grid item itu sendiri | **dibangun** — `[penyimpangan sadar]` baris item MBU membuka section `Estimasi` (Pega tidak dapat mengisi estimasi MBU) |
| panel item: Add estimasi | `ValidateInputEstimate_act` | **dibangun** |
| Estimation Date (change) | `ProtectionDate_Act` (+ `Set7Hours`) | **dibangun** — `Set7Hours` tanpa padanan (tanggal sudah WIB, §8) |
| Currency (change) | `SetConvertValueKurs_Estimation` | **dibangun** |
| Estimation Gross / Estimasi > TSI (change) | `CheckEstimateValue` | **dibangun** — pzInsKey tertulis mati langkah 25 dibuang (bawaan c) |
| Delete estimasi | `DeleteValueEstimation` + `CountSpreadingClaim_ACT` | **dibangun** |
| Deductible (% / Type / Amount, change) | `CountTSI_Act` | **dibangun** |
| Deductible Info (centang) / Format (change) | `change:postValue` | **ralat 11-10-2026** ("CETANGAN Deductible Info TIDAK BERFUNGSI!"): keduanya tanpa aksi - nilai hanya berubah di layar, Format (VIS `.DeductibleType=='true'`) dan baris Format 1 / 2 tidak pernah muncul. Kini aksi `PostValue` (isian digabung, layar dievaluasi ulang; `TestDeductibleInfoMemunculkanFormat`) |
| Spreading Claim: View Retro (TreatyType FAC retro) | local action `ShowRetro` | **dibangun** — pop-up hanya-baca dari halaman (`BukaRetro`) |
| Spreading Policy / Spreading Claim yang dapat disunting, **Save Spreading** (`LS42`) | `SetValueSaveSpreading` → `SaveSpreadingSP_act` (kata sandi) | **tidak tampil di XML** — bergerbang `pyWorkPage.ClaimData.ExGratia = 1`; penulis satu-satunya `InsertObjects_dt` 11 menulis `0`. Kata sandi tidak ditulis di mana pun (OQ-CFI-17) |
| Spreading: Share % (change) | `CheckTotalSpreadingPct_Act` | **tidak tampil di XML** — medan RO=ALWAYS, event tidak pernah terpicu (§7 butir 10) |
| Tambah / Delete baris spreading | `addRow` / `deleteRow` | **tidak tampil di XML** (VIS `… && 1=2`) |
| baris Total `InputParam.CARI41 / CARI42` | — | **tidak dibangun** — halaman requestor tanpa penulis |

## 4. Choose Surveyor / Input Adjustment — Section `ClaimSurvey` / `InputAdjustment`

| Label XML | Aksi XML | Status |
| --- | --- | --- |
| Claim No / View Status Payment Premi / Remark from Inputor | — / REST / — | **dibangun** / **nonaktif-OQ** OQ-CFI-14 / **dibangun** |
| Tab Claim Details (Registration / Estimation / Adjustment & Acceptation) / Policy Details and Claim History / Progress Claim | Section `InputEstimasi` / `ShowObjectAdj` | **dibangun** |
| Back | `BackToEstimasi` + finishAssignment | **dibangun** → Input Estimasi. **11-10-2026** ("TOMBOL BACK KALO UDAH ADA AKSEPTASI, HIDE AJA") `[penyimpangan sadar]`: tersembunyi bila ada adjustment ber-AcceptanceStatus 1 (`AdaAkseptasi`, dibaca dari baris adjustment - penanda IsAnyAcceptation tidak dihitung ulang sesudah komite); server menolak aksinya |
| Close Claim | local action `PreventRejectClaim` | **dibangun** (modal `tutup`, §5) |
| grid objek: Print DLA (VIS `IsFacretro = 1`, NA `DLAStatus` bukan 0/kosong) | `SetIndexObject` + harness `PrintDLA_dtl` | **dibangun** (`BukaDLA`, modal `dla:o`) |
| PrintDLA_dtl: Send Email / Message / Share Retro / RemarksDLA / Submit | `ChooseDla_Act` → `GenerateDLAFacin_Act` (P) / `DLAFacintoTreaty_Act` (S) | **dibangun** — DLA fac **P**; DLA treaty **S** bergerbang `CekLimit.CARI1` (OQ-CFI-24); tanpa lompatan mundur, `HitDLAClaimFacin` hanya produksi (§7 butir 4–5); `OS_AKSEPTASI_KLAIM` `UPDATE DLA_NO, DLA_DATE`; surel lewat outbox produksi; Draft DLA saat Send to Committe |
| grid item PA / Travel: Send / Sent Claim to committee | harness `ClaimCommittee` | **nonaktif-OQ** OQ-CFI-25 |
| grid Adjustment / Spreading Adjustment / BreakDown QS / Committee Accept Status: kolom nomor | `.pxListSubscript` (judul kosong) | **ubah 11-10-2026** ("KODE INI PERBAIKI"): nomor baris tanpa judul terbaca sebagai nilai Status di sebelahnya - berjudul "No" (`JudulNomorBaris`, semua lini, `TestKolomNomorBarisBerjudulNo`). Status baris disetujui (`IsApproved = 1`) = label "Approve" / "Reject" |
| grid item: View Status Payment Claim / View Payment Attachment | REST `getPaymentClaim` / `getPayAttachment` | **nonaktif-OQ** OQ-CFI-14 |
| grid Adjustment: + (VIS item `IsKomite` kosong/0) | `SetIndexObject_ACT` + `CountTotalEstimasi_Act` | **dibangun** (`CountTotalEstimasi`) — `[inferensi]` baris baru (OQ-CFI-23) |
| grid Adjustment: ikon hapus (VIS `IsKomite != 1`, NA `IsApproved = 1`) | `DisableSendComite` | **dibangun** — baris yang punya `KOMITE_ID` tidak terhapus (server) |
| grid `.KomiteList` ItemListEstimation | — | **tidak tampil di XML** (VIS `1=2`); daftar komite tampil di panel adjustment |
| footer `CountSpread.CARI24–27` | — | **tidak dibangun** — halaman requestor tanpa penulis |
| View Inc Fac Cedant Panel | harness `ViewCedantpanels` (`ViewCedantPanel`) | **dibangun** (`BukaCedant`, modal `cedant:…`); Choose / Choose All = `SetCedant_act` 1 / 2 |
| DLA No Ceding / DLA No SOB (change) | `SetDLACedingSOB` | **dibangun** |
| Payment Type (change) | `SetAdjTypePayment` | **dibangun** — kode 1–7 tampil apa adanya (label OQ-CFI-18). **Label terjawab 11-10-2026** (screenshot property rule `ASM-FW-GCNMFW-Data-Adjustment.PaymentType`): 1 Final, 2 Claim, 3 Salvage, 4 Adjuster Fee, 5 Adjustment, 6 Consultant Fee, 7 Cancellation |
| Choose Currency (change) | `CheckCurrency` | **dibangun** |
| Gross Adjustment (change) | `SetGrossAdjustment` | **dibangun** — langkah 6 (pre=false) selalu menolkan deductible % (literal XML) |
| Deductible Type / % / Value (change) | `SetNilaiResikoSendiri` | **dibangun** |
| VAT (change, PT 4 / 6) | `SetValueAdjusterFee` | **dibangun** |
| Payable To / Specify (change) | `SetPayable_Act` | **dibangun** |
| grid Ceding Co: Choose (Payable) | `SetPayableTo` | **dibangun** — langkah 1–4 ber-REMARK di Pega (pembaruan 08-01-2025): Choose tidak berbuat apa-apa. **11-10-2026** ("CHOOSENYA INI GA BERFUNGSI!") `[penyimpangan sadar]`: langkah 1–3 dihidupkan - Specify adjustment + klaim = ceding baris itu, rekening pertama ceding itu (GetDataBankAccount_sql: CedingCo + mata uang adjustment), pilihan Name of Bank sesudahnya = rekening ceding terpilih (`PilihPayableCeding`, `RencanaPayableTerpilih`, `TestChoosePayableToCedingLain`). Ubah Payable To / Specify / Gross sesudahnya menjalankan SetPayable_Act lagi (kembali ke ceding bawaan, XML). DEV: EVEREST REINSURANCE COMPANY tanpa baris BANKACCOUNT - rekening tetap kosong |
| Name of Bank (autocomplete) | `Result.pxResults` `SetPayable_Act` | **dibangun** (`PilihRekening`) |
| View KMT NO | `ViewKomite_act` | **nonaktif** — `KomiteNo` sudah tampil sebagai medan (dibaca ulang saja) |
| Ex Gratia (change) | `CekExGratia` | **dibangun**. **11-10-2026**: baris adjustment baru ber-`ExGratia` 0 (LS41 `pyDefaultValue` 0, yang dikirim Pega pada kiriman pertama panel); label 0 **NO** dari layout lama, 1 **YES** `[inferensi]` (OQ-CFI-37) |
| Spreading Adjustment: Add (VIS adjustment `ExGratia = 1`) / Share % (change) / Delete | `addRow` / `CountSpreadingAdjustment` / `deleteRow` | **dibangun** (`TambahSpreadAdj`, `CountSpreadingAdjustment`, `HapusSpreadAdj`) |
| Spreading Adjustment: View Retro | local action `ShowRetro` kelas SpreadingRisk | **nonaktif-OQ** OQ-CFI-25 |
| BreakDown Spreading Quota Share | — | **dibangun** (hanya-baca) |
| Save | click:save | **dibangun** (`SimpanAdjustment`) |
| Send to Committe (NA `IsKomite = 1` atau spreading kosong) | `SendPICProtect_Act` → harness `Comittee` | **dibangun** (`SendPICProtect`, modal `komite:…`, §5). **11-10-2026**: Gross / deductible / VAT yang ditolak proteksi dan masih di isian layar tidak lagi lolos - aksi medannya dijalankan dulu (§9 medan tertunda; `TestSendKomiteDitahanProteksiNilaiAdjustment`) |
| panel adjustment: layout group Tab (LS1) | — | **11-10-2026** ("INI BENTUK TAB!"): tab "Detail Chronology to Comitee" (LS2, VIS IsKomite = 1, tiga kolom berlabel atas) dan "Adjusment Detail" (LS5); spreading (LS39), Committee Accept Status (LS52), dan tombol (LS57) di dalam tab Adjusment Detail, seperti XML. Kolom **Status** grid Adjustment berlabel ("STATUS NYA BENERIN JANGAN ANGKA"): 0 Transfer to committee, 1 Approve, 2 Reject - label yang sama dengan History Adjustment layar Komite |
| panel adjustment: rupa | — | **11-10-2026** ("PERBAIKI TAMPILAN ADJUSTMENT! INI SAYA LAMPIRKAN LAYOUT LAMA"): LS12 dan LS15 = Inline grid triple berlabel di atas (`tiga`, paling banyak tiga sel per baris, panel sempit turun ke dua / satu); setiap sel = kontainer XML - medan sejajar berlabel atas (`kolom`: LS16 Payment Type | Detail, LS17 tiga medan estimasi, LS20-LS22, LS27), Payable To / Specify dan bank Stacked with labels left (`kiri`, LS29 / LS30 / LS37), Komite Number (LS38; label XML `Text Input` bawaan, nama dari layout lama); kontainer tersembunyi tidak memakan sel. Kotak centang berkapsi kanan (Net For Collection, Transfer Direct to Kasir - kapsi layout lama); Status Kasir tanpa label. Total spreading adjustment / QS (LS46 / LS51) = kaki grid-nya. Medan / kondisi / aksi tetap InputAdjustment (`TestPanelAdjustmentTigaKolomIkutLayoutLama`) |
| Acceptation (VIS `AcceptanceStatus = 1`) | `SaveAcceptation` + `HitServiceToKasir_Act` | **dibangun** — aktif sesudah komite tahap 2 menulis `AcceptanceStatus`; kasir lewat outbox produksi, sekali per klik (§8). RALAT 10-10-2026: `HitServiceToKasir_Act` 3 `getStatusKonversi_Act` bertransisi pasca `.StatusKonversi=="1"` (COUNT `reinsurance.trloss_detail_t`, hanya produksi) - selainnya keluar; kini ditegakkan (`Acuan.StatusKonversi`). Langkah 12 `PrintPDFAccep_MultiAksep` **dibangun** (10-10-2026, stream `AcceptanceNotePDF` diberikan work owner): salinan halaman klaim diambil sebelum langkah 13; PDF disusun, digambar, diunggah sesudah aksi tersimpan (InsertGoogleStorage_Act Folder "Claim" Durasi 1800, T_STORAGE_IMAGE + DOCUMENT_CLAIM KATEGORI_1 "AcceptanceNote" MIME "pdf"); lini MBU / Travel / PA → `info` OQ-CFI-20; §8 butir 24 |

## 5. Komite, penutupan, penolakan

| Label XML | Aksi XML | Status |
| --- | --- | --- |
| ClaimComite: Date / Initial / Chronology / Extent / Legal Liability / Salvage / Adjuster Fee / Remarks | harness `Comittee` | **dibangun**. **11-10-2026** ("RAPIHKAN DAN TOMBOL SEND NYA GA ADA!"): isian berlabel di atas, dua kolom (Chronology / Extent of Loss, Legal Liability / Remarks; Salvage / Adjuster Fee sebelum Remarks menurut Payment Type). Tombol **Send Claim to Committee** tampil selama Protect.CARI1 / CARI2 lolos: VIS LS7 (isian wajib terisi) dievaluasi server dan mengetik di pop-up tidak memanggil server, jadi tombol tak pernah muncul `[penyimpangan sadar]`; isian wajibnya ditegakkan server saat diklik (422 "Remarks: Value cannot be blank") |
| Send Claim to Committee (VIS `Protect.CARI1 = 1 && CARI2 = 1` + isian wajib) | `DraftGenerateDLAFacin_Act` + `CreateKMTNo_Act` + `SetListKomite_act` | **dibangun** (`KirimKomite`) — kasus `KMT-` TT2 lahir (`T_WORK_CLAIM` LINI FACIN, `TAHAP Komite_Flow`, `POSITION` = `OPERATOR_ID` tingkat 1, `T_GENERAL_KOMITE`, `T_KOMITE_KOMITELIST`), `KomiteNo`, progres "Waiting Committee", kronologi, email outbox produksi. Bundel PT 2 = OQ-CFI-28. Tanggal / circumstances tidak ditimpa (§7 butir 6–7) |
| `SendPICProtect_Act` proteksi lampiran LOD / DLA / SPGR / ADU / Invoice / Salvage, rekening, premi | — | **dibangun** — cacah `DOCUMENT_CLAIM` × master `T_KATEGORI_DOC_KLAIM` FAC; premi lunas pola Claim Prop (saldo premi non-produksi dilewati) |
| PreventRejectClaim: Close Without Payment (centang) | postValue | **dibangun** (`SetelAlokasiSalvage`) |
| PreventRejectClaim: Yes (CWP) | `SendCloseClaimToKomite` (TT4) | **dibangun** tahap 2 (`SendCloseClaimToKomite`, KCF-03) — 2 `Remark` / `Remark_Close` := Remarks; 4–6 adjustment belum diputus → pesan VERBATIM, kasus tidak lahir; 7.1–7.4 kasus `KMT-` tanpa adjustment, `TRANSFER_TYPE` 4, tangga satu tingkat `ReasClaimDeptHead` (prompt §5 butir 7: akun + email tertulis mati diganti workbasket); 7.5–7.7 `Komite.CARI1`, kronologi "Request close claim without payment " + KMT; 7.9 email outbox produksi. Penjaga permintaan ganda (`PesanTutupKomiteGanda`) `[penyimpangan sadar]`. Sesudah lahir: LS23 (isian + Yes) tersembunyi, LS28 label sukses, checkbox LS21 tetap tampil (tanpa syarat). Chronology / Extent / Policy Liability pop-up disimpan di kepala kasus komite (migrasi komiteclaimfacin 643, OQ-KCFI-03) |
| PreventRejectClaim: Yes (close) | `CloseClaim` | **dibangun** — `ValidationAdjustmentKomite`, setiap pesan menghentikan (§7 butir 8), `OS_AKSEPTASI_KLAIM` STS 4, konversi STS 4 (produksi), kronologi, `JSON_KLAIM`, Resolved-Completed; 12 `ASMForceCaseClose` `CloseAllSubCases=true`: kasus komite KMT- klaim ini yang masih menunggu ikut ditutup (10-10-2026, `TutupKomiteAnak`) |
| `UpdateTotalJob_sql` (`mst_user_teknis.TOTAL_JOB`) | CloseClaim 7 | **tidak dibangun** — `ClaimData.UserTeknis` tanpa penulis (LS177 VIS `1=2`); langkah tidak pernah berjalan (OQ-CFI-33) |
| ClaimComiteeReject: Yes (naJika `IsReject = 1`, tampil selama `Komite.CARI1` kosong) | `SendRejectClaimToKomite2` (TT3) | **dibangun** tahap 2 (`SendRejectClaimToKomite2`, KCF-03) — 2 `Remark` / `Remark_Close` := Remarks; 4–6 sudah ada akseptasi / estimasi belum Face Claim → pesan VERBATIM (pesan terakhir menang); 7.1–7.4 kasus `KMT-` tanpa adjustment, `TRANSFER_TYPE` 3, tangga `ReasClaimDeptHead`; 7.5–7.7 kronologi "Request Reject claim " + KMT; 7.9 email outbox produksi. Sesudah lahir: isian LS4 tetap tampil berisi ketikan (clipboard Pega; `models.IsianPopUpTutup` dibawa `BawaSementara`), Yes diganti label sukses. Chronology / Extent / Policy Liability pop-up TT3 / TT4 disimpan di kepala kasus komite (`repository.TeksKomite`, migrasi komiteclaimfacin 643, keputusan OQ-KCFI-03) dan tampil di layar komite. Keputusan komite (CLAIMREJECTED, OS STS 2, Resolved-Rejected) = modul `komiteclaimfacin` |
| Perluasan tangga `ApprovalKomite_Act`, keputusan komite, `SaveAccept_ACT` (OS STS 1) | — | **tahap 2 dibangun** — modul `komiteclaimfacin` (`modul/komiteclaimfacin/docs/PARITAS.md`), tulis balik lewat kontrak `kontrak.KlaimFacInKomite` (`services/kontrak_komite.go`) |

## 6. Tidak tampil di XML (tidak dibangun)

- Grid / tombol bergerbang `pyWorkPage.ClaimData.ExGratia = 1` (Spreading Policy / Claim yang dapat disunting, Save
  Spreading): `ExGratia` tingkat klaim selalu `0`.
- Tombol Add / Delete spreading estimasi (`1=2`), medan `ClaimData.UserTeknis` (`1=2`), grid `.KomiteList`
  ItemListEstimation (`1=2`).
- Grid estimasi Travel di section `Estimasi` (Travel memakai `EstimasiPA`).
- ~~Bagian "Success Create Request to Committee" PreventRejectClaim (`Komite.CARI1` hanya terisi oleh TT4, OQ-CFI-27).~~
  Tahap 2 (10-10-2026): `Komite.CARI1` kini terisi TT3 / TT4 — label tampil di pop-up Reject / Close sesudah kasus
  komite lahir, tombol Yes tersembunyi.

## 7. Sebelas arah perbaikan prompt §5 — `[penyimpangan sadar]`

| # | Rule | Perilaku Pega lama | Dibangun | Uji |
| ---: | --- | --- | --- | --- |
| 1 | `SearchPolis_act` / `GetPolisForClaim_SQL` | teks cari disisipkan ke SQL (`{ASIS:InputData.CARI1}`) | bind `:1` | `repository/sql_test.go` `TestCariPolisTeksTerikat` |
| 2 | `SetIndexObj_Act` 1 | `PlaStatus = 1` saat pop-up Print PLA dibuka | `PlaStatus = 1` hanya sesudah `GeneratePLA` berhasil | `models/perbaikan_test.go` `TestPlaStatusHanyaSesudahGeneratePLA` |
| 3 | `BackFromRegister` 8–11 | medan kosong diisi "Test", "0812", "IDR", tanggal sekarang | tidak mengisi apa pun | `TestBackFromRegisterTanpaNilaiUji` |
| 4 | `ChooseDla_Act` | `RetroList.TotalClaim = "123"`; transisi kosong #9 lompat mundur ke `JMP1` | tanpa lompatan: log `AKSEPATSI` dan `DLA` masing-masing sekali; `DLAList` / `TotalClaim` tidak dibangun — nol pembaca di korpus (OQ-CFI-25) | `handlers/perbaikan_test.go` `TestDLATanpaLompatanMundurDanHitDLAHanyaProduksi` |
| 5 | `ChooseDla_Act` 10 | REST `HitDLAClaimFacin` tanpa penjaga produksi | efek `dla-klaim` di outbox hanya bila produksi | uji yang sama (non-produksi: nol efek) |
| 6 | `CreateKMTNo_Act` L4 | memotong `pxCreateDateTime`, `DateOfLoss`, tanggal polis ke 8 karakter lalu menyimpannya | tanggal tersimpan tidak diubah | `models/komite_test.go` `TestTandaiKirimKomiteTanpaMengubahTanggalDanKronologi` |
| 7 | `CreateKMTNo_Act` L11 | `TempCommiteClaim.CircumtansesCouseOfLoss` ditimpa pzInsKey | circumstances tidak ditimpa | uji yang sama |
| 8 | `CloseClaim` 6 | `Param.ErrorMsg` dibaca tanpa halaman parameter: pesan tampil, klaim tetap tertutup | setiap pesan `ValidationAdjustmentKomite` / proteksi DLA menghentikan penutupan | `handlers/perbaikan_test.go` `TestCloseClaimBerhentiPadaProteksiDLA`, `models/komite_test.go` `TestCekDLATutup` |
| 9 | `InsertProgressClaim` 2.1 | `ProdKe := REGISTER` tanpa kutip → posisi kosong; dipanggil setiap assignment dibuka | posisi `REGISTER`; sisip hanya bila (kasus, posisi) belum ada, posisi lain `Done` (isi `PEGA_PROGRESSCLAIM`, upsert) | `TestProgresPosisiRegister` |
| 10 | `CheckTotalSpreadingPct_Act` | nilai tanpa kutip di #1 | **tidak berlaku** — pemicunya `change` medan Share % yang RO=ALWAYS di tiga section estimasi; activity tak terjangkau, tidak dibangun | — (bukti: `Section/Estimasi.xml`, `EstimasiPA.xml`, `EstimasiMarine.xml`) |
| 11 | `GetSequenceNumber_SQL` | `ParamSeq.CARI3` (tanggal procedure) tidak pernah diisi | periode dari tanggal tutup buku lewat `inti/backend/penomor` (`HitungPeriodeNomor`), pola Claim Prop / Claim Non Prop | uji milik `inti/backend/penomor`; `models/perbaikan_test.go` `TestRakitNomor` (bentuk nomor) |

## 8. Penyimpangan sadar lain dan `[inferensi]`

**Penyimpangan sadar** (perilaku Pega lama → di sini):

1. **Write-through** (pola Claim Prop): Pega menyimpan clipboard hanya pada Save / Obj-Save; di sini setiap aksi
   menyimpan halaman. Aksi yang keluar karena pesan validasi dibatalkan seluruhnya.
2. **Pra-proses setiap aksi**: Pega menjalankan pra-proses FlowAction saat assignment dibuka; di sini setiap muat. Akibat
   yang terlihat: `SetStartDateAdjustment` 1 menimpa `StartDateAdjustment` setiap kali dibuka → di sini hanya bila kosong;
   item `IsKomite` di-reset `0` (bila `AktifButton` kosong) sehingga tombol + adjustment muncul lagi sesudah aksi
   berikutnya — Pega sesudah assignment dibuka ulang.
3. `SetKomiteList_ACT` 2 menghitung ulang `ValueAdjustment` SEMUA baris adjustment item (kondisi
   `PaymentType != "4" || PaymentType != "3"` selalu benar) → baris fee / salvage menjadi 0. Tidak ditiru.
4. `CLaimFaceSheet_Act` 23.4 `TotalEstimasiConvertIDR` dari akumulator yang tidak dinolkan antar item → di sini
   `TotalEstimasiValue` adjustment = Σ ConvertValue item itu sendiri.
5. `GenerateDLAFacin_Act` dipanggil sekali untuk SETIAP baris spreading 10015 setiap item (objek ber-dua item retro
   menerbitkan dua nomor DLA) → di sini setiap jenis nomor (P / S) sekali per Submit.
6. `SaveAcceptation` menjalankan `HitServiceToKasir_Act` dua kali pada klik pertama → sekali per klik, lewat outbox.
7. `SetValueAdjusterFee` 13: kondisi `IsError < 2 || …` dibangun menurut teks XML.
8. `InsertObjects_dt` 8.6.25 menulis `OccupationList` hanya ke calon terakhir → setiap calon Aneka menerima
   `OccupationList` lokasinya sendiri.
9. `ProtectDownloadFaceClaim` 2 membaca klaim tersimpan hanya di produksi (pola Claim Prop OQ-CP-11): di luar produksi
   CFS objek kedua tidak tertahan.
10. `GetAllData_Act` (Outstanding Summary): kurs Travel / PA dicocokkan pada `.CARI15` yang tidak dipilih SQL, dan kolom
    Currency Marine menampilkan `.CARI10` yang tidak dipilih → di sini mata uang baris.
11. `GeneratePLATreaty_Act` 14–15 menolkan `TotalEstimasi` / `TotalGrossEstimasi` item (indeks salah) → tidak dibawa.
12. `BrowseProvince_RD`: section mengirim `ID2 = CountryID` dan membaca `.Note` / `.ID` yang tidak ada → disaring nama
    negara, `ProvinceID` ikut terisi.
13. `GetAddressCeding` membaca `M_CLIENT.JSONDATA` → tabel datar `CLIENT_ADDRESS` (keputusan work owner Claim Prop
    08-10-2026).
14. `SetchronologyKlaimFacIn` 1.1 memetakan tiga nama orang tertulis mati ke jabatan untuk setiap baris kronologi pelaku
    itu → jabatan pelaku dari roster `EMAILKOMITE` FACIN untuk setiap baris, selainnya "Claim Admin" (bawaan c, pola
    Claim Prop). Jabatan roster yang tidak dipetakan XML (mis. Claim Supervisor) kini tampil sebagai jabatannya.
15. `Set7Hours` (Estimation Date): pergeseran GMT → WIB tanpa padanan; tanggal disimpan dalam zona Jakarta.
16. Satu perlakuan presisi (pola Claim Prop): `@divide` berskala 4 / 8 di beberapa rumus → presisi penuh, pembulatan
    hanya di penyimpanan `NUMBER(38,10)` dan tampilan.
17. `IsBondingAndCustomBonds` C merujuk `IsCustomBonds` yang tidak ada (yang ada `IsCustomBond`) → bernilai salah.
18. Medan tanpa pembaca di korpus tidak disimpan: `ShareNusare`, `PUCLStatus`, `ReceiverClaim`, `TFAllObj`, `OurRef`,
    `StartPeriod` / `EndPeriod`, `DLAList`, `PLAList`.
19. `CheckDate_Act` 9 / 13 `+1` atas teks tanggal (penggabungan teks) → pesan "Start Date Policy" tidak pernah terpasang.
    **Keputusan work owner 11-10-2026** (OQ-CFI-12, "PROTEKSI INI KO GA ADA?"): DOL SEBELUM tanggal mulai polis ditolak
    dengan pesan itu; DOL tepat di tanggal mulai boleh. Uji `TestProteksiTanggalRegister` (semua pesan CheckDate* lewat
    jalur HTTP).
20. `GetSpreadingMarine_Act` (pra-proses `EstimasiMarine_FA`, Pega: setiap pane item Marine dibuka) → pra-proses Input
    Estimasi hanya bila Spreading Policy item masih kosong (`models.LengkapiMarine`): pra-proses di sini berjalan setiap
    aksi, dan menjalankannya ulang menimpa Spreading Claim hasil estimasi.
21. Pesan validasi membatalkan seluruh aksi (write-through); Pega membiarkan perubahan clipboard sebelum pesan (mis.
    baris spreading yang dihapus tetap hilang walau total < 100%). Layar validasi = halaman sebelum penangan + isian
    layar + pesan; contoh: 50 / 30 / 20 → nolkan baris yang akan dihapus dan naikkan baris lain dalam satu kiriman, lalu
    hapus.
22. Kolom milik komite (`ACCEPTED_NO`, `ACCEPTED_DATE`, `ACCEPTANCE_STATUS`, `IS_APPROVED` di `T_CLAIM_ADJUSTMENT`) tidak
    ikut UPDATE simpan halaman klaim (`Kolom.MilikKomite`): klaim tidak pernah menulisnya, dan simpan yang membaca
    halaman sebelum komite memutuskan tidak boleh menimpa keputusan itu.
23. Submit pop-up (Print PLA `GeneratePLA`, Comittee `KirimKomite`) memeriksa ulang tombol pembukanya di server (Pega:
    pop-up hanya dapat dibuka lewat tombol yang aktif); hapus adjustment bertaut komite = 409.
24. Nota akseptasi `PrintPDFAccep_MultiAksep` (SaveAcceptation 12; `models/akseptasi_pdf.go`, `models/pdf_akseptasi.go`,
    `services/akseptasi_pdf.go`, salinan Komite Claim Fac In - bukan impor):
    - urutan (prompt tahap 2 §6 butir 7, pola Komite Prop): salinan halaman diambil di langkah 12 (saringan S9.3.11
      `IsPrintAccept == "" && AcceptanceStatus == 1 && AcceptedNo == Param.NoAkseptasi` VERBATIM); PDF disusun,
      digambar, diunggah, dan dicatat SESUDAH aksi tersimpan (tidak menahan kunci kasus; pemutusan klien tidak
      membatalkannya); gagal (termasuk panik penggambar) → akseptasi tetap, `info` `PesanDokumenGagal`;
    - angka `DecimalFormat("#,###.####")` id_ID atas nilai DESIMAL tersimpan, bukan `double` Java (nol float untuk
      uang): beda hanya pada seri tepat ...5 di desimal ke-5 dan nilai di atas 15 digit bermakna;
    - `HTMLToPDF` → go-pdf/fpdf yang menggambar susunan stream (label VERBATIM, ukuran CSS stream);
    - perbaikan (maksud pasti): `(6 || 4) && U` di S9.3.6 / S9.3.9; `TemporaryUang` segar per cetak, bertambah selama
      iterasi; Property-Remove di dalam For Each = saringan biasa; "Swift Code" hanya bila terisi;
    - VERBATIM (maksud tidak pasti): "Location of Loss" = CauseOfLoss (S13); nominal / Payable To / spreading hanya item
      terakhir yang lolos; tabel QS tanpa penulis tidak tercetak; "Premium Paid On" kosong (OQ-CFI-14).
25. Polis dari tabel datar NB Fac In (OQ-CFI-11, perintah work owner 10-10-2026): peta kolom → medan / jalur DISALIN
    dari skema loader nbfacin sesudah init (`polis_datar_peta.go`: skema_gen + amandemen + penunjuk Idx*/Index*; medan
    `T_GENERAL_POLIS` berlokasi `PolicyData` menurut JSON_DATAGUIDE dokumen FAC DEV; CURRENCY_CODE CoverageList /
    AnekaList = `Currency.Name`). Peta itu RANCANGAN: repository hanya memilih tabel / kolom yang ada di katalog
    (ALL_TAB_COLUMNS, sekali per proses) - DEV 10-10-2026: 34 dari 76 tabel, 45 kueri builder lolos di DEV (COUNT).
    Kasus = penunjuk DATA_JSON, cadangan IDPEGA JSON_POLIS / FACINPRODUCTION. Medan tanpa kolom datar diisi dari sumber
    layar / rule NB Fac In: BusinessCode / BusinessName / BusinessOldId / GroupPanel = master BUSINESS menurut
    `FACINPRODUCTION.BUSINESSCODE` kasus itu (cadangan NOTE = Class Of Business + BUSINESSGROUPID); BusinessType = port
    decision table `BusinessType_DeT` (SetBusinessType_Act 2, bawaan "UNKNOWN"); Insured = `T_NB_OPPORTUNITY` +
    `T_M_ACCOUNT`; StatusBusiness "New Business" = kode 1 (Section Periode); Marketing / TeamGroup / Branch =
    `MARKETINGOFFICER` menurut MOID; PolicyNo = nomor JSON_POLIS; indeks = port `SetIndexCoverage_DT` (Fire) dan
    `SetIndexObject_DT` (Aneka); CurrencyID item / `Currency.ID` = `CURRENCY.ID` menurut kode; sentinel `UNKNOWN`
    kolom mata uang = kosong. Galat penyusunan saat muat = kasus tetap terbuka tanpa halaman polis (dicatat log); Choose
    Polis menolak. `[penyimpangan sadar]` tabel datar dibaca hidup setiap muat (bukan cuplikan seperti JSON_POLIS Pega).
    Kinerja (pengukuran DEV 10-10-2026, izin work owner baca-saja): tabel datar nbfacin berisi 2-4 baris tetapi
    bersegmen 19-240 MB dengan statistik basi -> setiap kueri anak menyapu segmen (100-190 ms), satu susun ~0,7-0,8 s
    dan setiap aksi memuat dua kali. Perbaikan: hanya tabel yang dibaca klaim (`tabelDipakaiKlaim`), hint indeks induk
    dari katalog (`sqlIndeksDatar`), tembolok dokumen 5 menit per polis (`tembolokDatar`) -> susun pertama ~0,24 s, muat
    berikutnya ~15 ms. Sisa: `T_PROPERTYITEMLIST` tanpa indeks PARENT_ID (~0,12 s) - urusan DBA / nbfacin (shrink /
    statistik / indeks), bukan modul ini. **11-10-2026 ("SAAT BUKA CLM NYA LAMBAT!")**: DEV kini ~1 juta baris tabel datar;
    `T_PROPERTYITEMLIST` hanya ber-PK sehingga setiap susun menyapu seluruh tabel (0,1 s hangat, buka kasus pertama
    terukur 5,9 s saat cache Oracle dingin; buka berikutnya 21 ms). Umur tembolok 5 menit -> 8 jam (polis produksi satu
    PRODKE tidak berubah). Usul DBA / nbfacin: indeks `T_PROPERTYITEMLIST (PARENT_ID)`.
26. Tampilan polis NB Fac In hanya-baca (perintah work owner 10-10-2026: "SAAT CHOOSE, TAMPILAN POLIS NYA MENGIKUTI
    TAMPILAN POLIS JUGA. IKUTI MODUL NB FACIN. TAMPILAN READ ONLY" + "DI TAB POLICY DETAIL JUGA BEGITU"; pilihan "salin
    tampilan di Claim Fac In"). `detailPolis()` (detail Choose Polis + tab Policy Detail) = unsur `polisNB`: blok
    General Pega (Offer Slip Number, Policy No, Period, Prorate / Short period, ...) DIGANTI rupa layar Inward
    Facultative `modul/nbfacin/frontend` yang disalin (pola, bukan impor; nbfacin tidak disentuh) ke
    `frontend/polisNB.ts` + `components/PolisNB*.tsx`: General dua kolom (label VERBATIM `nbfacin/labels.ts`) lalu
    sepuluh tab LANGSUNG (pemisah kosong dan centang Show Detail nbfacin dibuang - work owner 11-10-2026 "BUAT LEBIH
    SIMPEL"). Seperti nbfacin, isi tab hanya kasus FIRE (When `IsFire` dievaluasi server, `kendali`
    "fire") dan hanya Object (tujuh sub-tab) / Coverage / Spreading / Inw Fac Cedant Panels; tab lain = pesan inti
    `BelumTersedia`. Clauses = `BelumTersedia`: `ClauseList` milik `pyWorkPage` kasus NB, bukan halaman polis
    `OfferFacIn`. Total yang nbfacin HITUNG (TotalTSIList, TotalTSIPremiGrossList, TotalTSIPremiSpreadRNM, ringkasan
    TempTotal* / TotalSpreading*) dihitung ulang dengan rumus nbfacin (eksak BigInt; spreading = port `totalLokasi` /
    `totalMataUang` / `ringkasan` tanpa pemeriksaan batas treaty + kurs, nama treaty dari TreatyName baris spreading);
    nilai tersimpan polis (TotalTSINusaReSpreading, TotalPremiNusaRe, ClaimEstimation baris) apa adanya. Tanggal
    dd/mm/yyyy, stempel GMT digeser ke WIB (nbfacin `services/kasus.go`). Tombol (Add / Delete / Choose / Save / CSV /
    Upload) tidak dibawa; View Retro List tetap di bawahnya (baris Inline - tidak terbawa ke kaki modal). Tabel datar
    sub-tab (BuildingConstruction, SurroundingRisk, TableOfLimit, FEA, ListCauseOfLoss + CoinsData, Deductible) ikut
    dimuat HANYA bila katalog menemukan indeks induknya (`tabelTampilanNB`): DEV 10-10-2026 `T_LISTCAUSEOFLOSS` ~1 juta
    baris tanpa indeks PARENT_ID = ~10 detik per kueri -> Loss Record polis dari tabel datar kosong sampai indeksnya
    dibuat (OQ-CFI-35); lainnya 1-3 ms. Kode Policy Type Pega tidak diketahui (OQ-CFI-36).
    **11-10-2026** ("tampilan view polis yang belum tergenerate tolong diperbaiki. cek tampilan nya di nb facin!"):
    tab dan isinya per lini seperti nbfacin (`InwardFacultative.tsx` `tabTampil` + dispatch): Clauses / Scoring Risk hanya
    FIRE, Loss Record / Correspondence hanya akun IT (keputusan work owner: ikut NB); COB di luar FIRE / Aneka / Marine
    Cargo = pesan nbfacin `cobBelum`. **Keluarga Aneka** (Aneka not bonding, HE, Golf, Bonding / Custom Bond; lini dan
    varian kontainer C1-C9 = When nbfacin `lini.go` + registry, `polisNBAneka.ts` - When klaim tidak dipakai karena
    berbeda: NB memasukkan CustomBond ke IsAneka, IsElectronicEquipment NB berspasi tidak pernah benar): tab Object
    (lokasi -> Risk Address + Top Risk, grid Golf, okupasi -> objek K2/K3/K4 -> form objek berkartu Object / Vessel /
    Vehicle / Aircraft / Sum Insured / Identity / Period / Loan Details / Note / Akseptasi Memorandum), Coverage (okupasi /
    Golf -> objek -> coverage -> form + Deductible; Grand Total per mata uang `totalaneka.go`) dan Spreading (detail
    empat angka + baris spreading; ringkasan `spreadinganeka.go` totalAneka + ringkasan FIRE). Data: kolom T_ANEKALIST /
    T_MARINEHULL / T_VEHICLEHE amandemen nbfacin 214 masuk peta, T_MARINEHULL / T_AVIATIONHULL / T_MEMORANDUMACCEPT
    dimuat sebagai tabel tampilan (DEV 11-10-2026: semuanya berindeks PARENT_ID, 0-29 baris; kueri dicoba read-only).
    **Marine Cargo** (`polisNBCargo.ts`): tab Object (grid No. / Top Risk / SailDate / Goods / Ship Name / Conveyance ->
    form cargo: kepala Trading / Conveyance / Goods / Packing, Survey Agent, LC, Ship Data, data pelayaran, Additional
    Ship), Coverage (coverage + form selalu terbuka seperti nbfacin, Warranty Extension, Deductible ber-Deskripsi, Grand
    Total `totalGrossCargo`), Spreading (detail lima angka + baris spreading; ringkasan `totalCargo` + ringkasan FIRE). Data
    pelayaran dibaca dari dua bentuk: polis Pega `CargoList(n).PolicyData.*` dan polis tabel datar (loader V-24b melipat
    `.PolicyData` ke baris cargo). Data: kolom T_CARGOLIST / T_LC / T_SHIP / T_ADDITIONALSHIP migrasi nbfacin 256 masuk peta,
    T_SURVEYAGENT baru di peta, T_SHIP / T_LC / T_SURVEYAGENT / T_ADDITIONALSHIP / T_ADDITIONALCOVERAGE dimuat sebagai tabel
    tampilan (DEV 11-10-2026: berindeks PARENT_ID, 0-2 baris; kueri dicoba read-only). Penjaga SQL baca-saja repository kini
    memakai batas kata (kolom FLAG_DELETE bukan pernyataan DELETE).
    **Premium Deduction / Brokerage Fee, Payment, Inw Fac Cedant Panels** (FIRE, keluarga Aneka, Marine Cargo -
    `polisNBBayar.ts`; keputusan work owner 11-10-2026 "Nilai tersimpan"): Deduction / Payment = nilai tersimpan tabel NB
    (T_GENERAL_POLIS PAY_*, T_CURRENCYLIST PAY_* / SUM_TOTAL_PAYMENT, T_LISTINSTALLMENT), bukan rumus NB; kepala PPH / PPN
    tanpa tarif (tarif master TAX terbaru dapat berbeda dari tarif saat disimpan). Rincian Cedant per mata uang tidak
    disimpan NB -> dihitung seperti nbfacin `ringkasanRNM` (dasar x ShareCeding / 100). **Clauses** (FIRE): T_CLAUSELIST
    (satu baris per argumen) dimuat sebagai tabel tampilan, baris ber-SEQ_NO sama dilipat menjadi satu kartu (nbfacin
    `kumpulkanKlausa`); kartu Clause Code / bahasa / Total Argument / Title / Description tanpa tombol Choose Clause / See
    Clause and Argument (isi klausa dan argumen tidak dimuat) - menggantikan catatan `BelumTersedia` di atas.
    **Scoring Risk** (FIRE, `polisNBScoring.ts`; "Nilai tersimpan"): T_SCORINGRISK -> T_SCORINGRESULT /
    T_DATASCORINGRISKLIST -> T_SCORING_FACTOR -> T_SCORING_OPTION (nbfacin 205) dimuat sebagai tabel tampilan + kolom
    penentu INDEX_LOCATION / FACTOR_GROUP / FACTOR_NAME / PERCENT_LOSS_RATIO / CHECK_BOX / OPTION_NO; rupa TabScoring
    nbfacin baca-saja (kepala, tabel 8 faktor / 24 sub-faktor berlabel VERBATIM, centang terkunci + skor tersimpan,
    subtotal faktor = jumlah eksak skor sub, Total / Scoring Result, kartu Classification / Recomendation / Required
    approval); opsi tercentang ganda -> nomor terkecil seperti nbfacin. Blok General "Risk Scoring" = FinalScore + lencana
    NoteFinalScore. Polis Pega (JSON_POLIS) dibaca juga dari bentuk halaman `.<Grup>.<Nama>.ChechBoxN`.
    **Loss Record / Correspondence** (keputusan work owner "Ikut NB: hanya akun IT"): data TIDAK ikut halaman polis -
    rute sendiri `GET /api/claim-fac-in/kasus/{id}/polis-nb-it` menegakkan M_LOGIN_GO.JOB_POSITION "IT" di layanan
    (401 tanpa identitas, 403 selain IT; nbfacin `wajibJabatanIT`); layar memanggilnya dan menampilkan kedua tab hanya bila
    dijawab 200 (lunak seperti nbfacin `AksesTabIT`). Kasus NB = penunjuk JSON_POLIS, selain itu dicari seperti tabel datar
    polis. Loss Record (isi hanya FIRE; lini lain pesan `cobBelum` seperti nbfacin): premi, grid kerugian, Total Loss per
    Class Of Business + Currency (jumlah eksak baris), Loss Ratio 1 / 3 / 5 Year - nilai tersimpan Save NB terakhir (Uw
    Year, Loss Amount ( IDR ), rasio tidak dihitung ulang). Correspondence (setiap COB): grid Date / Sender / Subject /
    Recipient / Status, baris dibuka = pengirim, penerima, subjek, pesan; pesan HTML hanya ditampilkan sesudah pembersih
    nbfacin `bersihkanHtml` (tag polos tanpa atribut, script / style dibuang). Tanpa Add / Remove / Save / Send Email /
    unggah CSV. Semua kueri baru dicoba read-only di DEV 11-10-2026 (berindeks PARENT_ID; hitungan baris saja).

**`[inferensi]`:**

1. Tombol + grid Adjustment tidak ber-`addRow`; `CountTotalEstimasi_Act` menulis `.Adjustment(<LAST>)` → baris baru
   ditambahkan lalu diisi; pemeriksaan 20–21 atas baris terakhir yang sudah ada (OQ-CFI-23).
2. Langkah Java "hapus spreading / treaty / value yang sama" tidak diekspor → baris berkunci sama (TreatyType; atau
   CurrencyID + TreatyType) dibuang kecuali yang pertama, menurut agregasi sesudahnya.
3. `.CurencyAdjustment` (pilihan Choose Currency) disusun ulang dari estimasi item ber-`PrintFaceClaim = 1`, bukan kolom.
4. Pilihan Cedant Panel disimpan sebagai `CEDANT_CHOICE` (`*` = Choose All); daftar disusun ulang dari polis.
5. `KunciCedingKasir` `@substring(pyWorkPage.Quotation.CedingCo, 0, 8)`: `pyWorkPage.Quotation` tanpa penulis →
   `OfferFacIn.QuotationData.CedingCo`.
6. Pembagi nol (`@divide` / `/`) bernilai 0 di rumus yang membagi properti yang dapat kosong (`hasilAkhirDLA`, dll.).
7. `IsError` kosong sama dengan 0 (tombol Submit Input Register).
8. Label ReportType / IndividualRiskType / Payable dari screenshot Claim Prop (kelas properti sama).
9. Perbandingan `CheckLimit_Act1` 9 numerik (`FindData.HASILD2` desimal); `@divide(x,1,2)` = ROUND_HALF_UP.

### RALAT 10-10-2026 — transisi pasca-langkah yang terlewat alat dump tahap 1

Alat dump activity tahap 1 hanya mencetak prasyarat langkah (`pyStepsPreCondParams`), tidak transisi pasca-langkah
(`pyStepsTransition` + `pyStepsTransParams`). Sensus ulang seluruh `Claim Fac In/Activity` (tahap 2, review kode):

| Rule / langkah | Transisi pasca | Status |
| --- | --- | --- |
| `HitServiceToKasir_Act` 3 | `.StatusKonversi=="1"` T=2 F=6 | **diperbaiki** — gerbang `Acuan.StatusKonversi` (hanya produksi) sebelum IDOfBank / muatan |
| `ValidateInputEstimate_act` 6 (pre: tanpa spreading) | `true` T=6 | **diperbaiki** — `TambahEstimasi` keluar sesudah baris dibuang |
| `SetNilaiResikoSendiri` 15 (pesan "can not be filled by Zero") | `true` T=6 | **diperbaiki** — keluar; pesan langkah 17 tidak ikut tampil |
| `CLaimFaceSheet_Act` 1 (`ProtectDownloadFaceClaim`) | `Protect.HASIL1==1` T=6 | sudah ditegakkan (`validasi`) |
| `CLaimFaceSheet_Act` 6 | `1==1` T=6 | jalur mati (`ExGratia` selalu 0, langkah 4–6 tidak dibangun) |
| `CloseClaim` 5.2 | `.IsFacretro==1 && .RemarksDLA==""` T=6 | sudah ditegakkan (§7 butir 8) |
| `SetProtectionEstimation` 23 | F=1 tanpa label | tidak dibangun (pemanggil = tombol OQ-CFI-25) |

## 9. Frontend

`frontend/` merender pohon tata server apa adanya (pola Claim Non Prop, disalin, tidak diimpor); bahasa layar Inggris,
kulit token `--cfi-*` salinan `.kelola-user`, kelas `claimfacin__*`.

- **Halaman awal** (`pages/ClaimFacIn.tsx`): tab *Process* (worklist pembuat) / *Resolve*, switch **Teknik** (workbasket
  `ReasKlaimTeknik`, nonaktif bagi akun bukan anggota, `GET /hak`), **Add Claim** (`POST /kasus`). Tanpa tabel komite
  (tahap 2). **11-10-2026** ("KOLOM INI DI HAPUS AJA DARI TAMPILAN!"): kolom **Assignment** (nama tahap / statusWork)
  dibuang dari tabel kasus; kolom lain tetap (Case ID, Claim No, Policy No, Insured Name, Date of Loss, Create Operator,
  Create Date/Time).
  > **RALAT 10-10-2026** (tahap 2). Kalimat lamanya dikutip utuh, tidak dihapus: *"Tanpa tabel komite (tahap 2)."* → tabel
  > **Committee** (`components/TabelKomite.tsx`) di bawah inbox, tanpa switch, tidak ikut tab, tampil bila `GET /hak` →
  > `komite` (pemegang `ReasClaimSPVB` atau akun / workbasket roster FACIN aktif); isinya `GET
  > /api/komite-claim-fac-in/kasus` (rute pinjaman); klik baris = layar komite di tempat (`onBukaModul`,
  > `MODUL_DIPINJAM`), Back / Submit kembali ke inbox. "View more details" di layar komite membuka berkas ini hanya-baca
  > (`?lihat=1`).
- **Panel bersarang** (`components/rincian.ts`, `TataView.tsx`): grid ber-`rincian` membuka panel baris
  `layar.panel["<prefiks>:<jalur>(n)"]`; panel dirender ulang lewat `TataView`, sehingga grid di dalam panel membuka
  panelnya sendiri (objek → item → estimasi / adjustment → detail adjustment). Satu baris terbuka per grid; status buka
  dikunci per tab dan per grid (`kunciGrid`) — temuan cek layar: tanpa itu baris terbuka bocor antar-tab.
- **Alamat aksi**: setiap unsur membawa `konteks` (`''` layar utama, kunci panel, atau kunci modal) dan `indeks` (baris
  grid di konteks itu; 0 = bukan sel grid / tombol Add kepala grid). Permintaan memuat `masukan` (layar utama + panel +
  modal), `mode` (`layar.mode` ditimpa isian lokal 14 jalur mode), dan `tahap`.
- **422**: pesan + pesan per medan tampil, layar yang dikirim server dirender ulang tanpa membuang isian lokal (juga di
  dalam modal / pop-up).
- **Pop-up sisi klien** sesudah aksi server tanpa efek: Cause of Loss (`pilihan/sebab` → `GetNameCauseofLoss`),
  Catastrophe (`pilihan/katastrofe` → `SetCatastrope` / `InputCatastrope` / `SaveCatasrtope`), Outstanding Summary
  (`pilihan/outstanding`), View Retro / View Retro List (hanya-baca dari halaman), hasil Choose Polis (`pilihan/polis`,
  tautan nomor polis → `CopyNB` `nopolis|prodke`).
- **Lampiran** (`components/PanelLampiran.tsx`, pola Claim Prop): kategori dari server (`T_KATEGORI_DOC_KLAIM` FAC).
- Isian angka: hanya angka, ribuan otomatis, rata kanan; tanggal `dd-mm-yyyy`. **11-10-2026** ("UNTUK CLAIM FACIN DAN KOMITE FACIN, UNTUK ANGKA PROTEKSI INPUTAN HANYA DECIMAL DAN TAMPILAN 4 DIGIT BELAKANG KOMA"): tampilan SEMUA angka (medan, grid, kaki, pop-up, tampilan polis NB, layar Komite Claim Fac In) selalu 4 desimal (`20.000.000,0000`; Loss Ratio polis NB tidak lagi 2 / 3 desimal, Share % spreading estimasi kini angka); server menolak kiriman medan angka terbuka yang bukan bilangan desimal mentah berdesimal paling banyak 4 (422 "Value must be a decimal number (max 4 digits after the comma)" di medannya, `models.AngkaTidakSah`). Minus di depan tetap diterima (pxNumber).
- **Medan berkalkulasi tertunda (11-10-2026)** ("INI KENA PROTEK, TAPI KLIK SEND KOMITE MASA BISA!. LOGIC PROTEKSINYA PERBAIKI"): nilai yang ditolak aksi medannya (422) tetap di isian layar dan ikut terkirim bersama aksi berikutnya; dahulu server menggabungnya tanpa menjalankan aksi medan itu (Gross Adjustment baru tanpa hitung ulang / proteksi -> Send to Committe membuka pop-up komite, Save menyimpannya). Kini medan ketik (angka / teks / area / telepon) beraksi yang nilainya berubah tanpa aksinya sendiri dijalankan dulu aksinya, seperti kursor meninggalkan medan (`services.medanTertunda` / `jalankanTertunda`); galat proteksinya membatalkan seluruh permintaan berpesan sama. Medan pilihan / tanggal / centang tetap digabung seperti semula.
- **Tab (11-10-2026)** ("UBAH TAMPILAN TAB NYA INI JADI LEBIH BAGUS DAN MODERN! DAN BUAT PERPINDAHAN TAB DAN HALAMAN NYA ITU SMOOTH" + "INTINYA SEMUA TAB DI MENU KLAIM FACIN"): SEMUA strip tab modul (Process / Resolve, layout group Tab layar kasus, tab dan sub-tab objek polis NB) = `components/TabGeser.tsx` - kontrol bersegmen, kapsul penanda meluncur ke tab terpilih (diukur dari tombol aktif, left / top / width tanpa transform), tab bersarang bergaya tenang, panah kiri / kanan / Home / End berpindah tab. Isi tab baru, layar kasus yang dibuka / ditutup, dan layar tahap berikutnya muncul memudar (opacity saja, tanpa sisa sesudah selesai); `prefers-reduced-motion` mematikan keduanya. Strip tab inti `StripTab` tidak dipakai lagi di modul ini.

## 10. Isolasi kotak masuk dan pembaca per-ID (prompt §6 butir 1)

Kasus Claim Fac In: `T_WORK_CLAIM.LINI = 'FACIN'`, `TAHAP` selalu terisi (`InputRegister` / `InputEstimasi` /
`InputSurveyor` / `Komite_Flow`; penutupan hanya mengisi `STATUS_WORK`, `repository/gudang.go` `sqlTutupKasus`),
awalan `CLM-` / `KMT-` (OQ-CFI-02), nomor `SEQ_WORK_CLAIM` (unik lintas lini). Dibaca dari SQL modul-modul itu, tidak
disunting.

| Modul | Saringan kotak masuk | Pembaca / penulis per-ID | Baris FACIN |
| --- | --- | --- | --- |
| Claim Life | `NVL(w.TAHAP, :tahapCadangan) = :tahap`, tahap "Input Register" / "Outstanding Claim" / "Medical Check" / "Claim Analis" (`claimlife/backend/repository/inbox.go` `sqlInboxWhere`, `models/tahap.go`) | **`GET /api/klaim-life/{id}`** → `KlaimLife.AmbilHeader` `SELECT … FROM T_GENERAL_CLAIM WHERE ID = :1` **tanpa LINI**. Penulis per-ID memeriksa tahap: `TahapBerlaku` → `TahapDariNama("InputRegister")` tidak dikenal → `TahapDariPeran(POSITION)` (pembuat / `ReasKlaimTeknik` ≠ `ReasLife*`) tidak dikenal → ditolak; UPDATE tahap / tutup bersyarat `NVL(TAHAP, …) = <label Life>`; `DELETE /api/klaim-life/{id}` gagal terang (ADR-U-0031) | **kotak masuk: tidak cocok** (TAHAP tanpa spasi, tidak pernah NULL). **Detail per-ID: kepala klaim FACIN terbaca** bila ID-nya diketik di URL Claim Life (hanya-baca; keadaan yang sama sudah berlaku untuk `CLMP-` / `CLMNP-`) → **OQ-CFI-31** (perbaikan = saringan LINI di claimlife, tidak disunting) |
| Komite Claim Life | `AND (w.LINI = :lini OR w.LINI IS NULL)` (`komiteclaimlife/backend/repository/komite_inbox.go`) | `sqlKasusKomite` `WHERE g.ID = :2 AND (w.LINI = :3 OR w.LINI IS NULL)` — membuka, memutuskan, riwayat semuanya lewat `Kasus` | tersaring (`FACIN` ≠ `LIFE`, tidak NULL) |
| Claim Prop | `w.LINI = 'PROP' AND w.ID LIKE 'CLMP-%'` (`claimprop/backend/repository/gudang.go`) | `Keadaan` `WHERE w.ID = :1 AND w.LINI = :2`; lampiran lewat `Keadaan` | tersaring |
| Komite Claim Prop | `w.LINI = 'PROP'` ketat + awalan `TKMT-` (`komiteclaimprop/backend/repository/gudang.go`, `tangga.go`) | `WHERE g.ID = :1 AND w.LINI = :2 AND w.ID LIKE :3` | tersaring |
| Claim Non Prop | `w.LINI = 'NONPROP' AND w.ID LIKE 'CLMNP-%'` (`claimnonprop/backend/repository/gudang.go`) | `WHERE w.ID = :1 AND w.LINI = :2` | tersaring |
| Komite Claim Non Prop | `w.LINI = 'NONPROP'` ketat + awalan `KMTNP-` | `WHERE g.ID = :1 AND w.LINI = :2 AND w.ID LIKE :3` | tersaring |

Arah sebaliknya: setiap kueri `T_WORK_CLAIM` modul ini menyaring `w.LINI = 'FACIN'` (kotak masuk juga `w.TAHAP <> 'Komite_Flow'`,
`repository/gudang.go`), sehingga kasus Life `CLM-` tidak terbaca di Claim Fac In.
