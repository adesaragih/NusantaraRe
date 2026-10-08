# Modul `claimprop` — Claim Prop

Klaim treaty inward proporsional - kasus `ASM-FW-GCNMFW-Work-ClaimTreaty` (`Flow/Flow_TreatyIn.xml`: Outstanding
Claim -> Input Acceptation -> Resolved-Completed). Dimigrasi 07-10-2026 (prompt
`_brief/PROMPT-IMPLEMENTASI-MODUL-CLAIM-PROP.md`); paritas tombol XML di `docs/PARITAS.md`, OQ di `docs/OQ.md`.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `claimprop` |
| Folder korpus | `Claim Prop` |
| GROUPMENU | `KLAIM` |
| Pemilik | `@PEMILIK-CLAIMPROP` |
| Status | dimigrasi |
| Rentang migrasi | `520-559` |
| Slot menu | `980-981` |
| Prefix rute API | `/api/claim-prop` |
| Kontrak disediakan | — |
| Kontrak dipakai | — (nol kontrak `inti/backend/kontrak`; pola modul lain disalin, tidak diimpor) |

`Pemilik` adalah penanda pemegang modul. Wilayah berkas yang boleh disentuh cabang
`module/<nama>` dijaga `.github/workflows/penjaga-wilayah-cabang.yml` - CODEOWNERS
dipensiunkan 1 Oktober 2026.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | spec, tiket (`issues/`), grilling, catatan — dulu `.scratch/claim-prop/` (dipindah dengan `git mv`, isi byte-identik) |
| `backend/` | `modul.go` (`Pendaftaran()`), `models/` (port activity, katalog, tata layar), `repository/` (Oracle), `services/`, `handlers/`, `tiruan/` (uji), `alat/pemuatlama/` (pemuat data lama, uji-kering) |
| `frontend/` | `menu.ts`, `rute.tsx`, renderer tata (`components/TataView.tsx`), halaman awal (`pages/ClaimProp.tsx`) |

## Migrasi

Rentang `520-559` (tabel R2, urut hulu ke hilir: migrasi modul hilir yang merujuk tabel modul hulu
selalu berjalan sesudahnya). Slot menu `980-981` hanya menyalakan `DIMIGRASI` baris modul ini (satu `UPDATE`,
nol `INSERT` — menu datar 30-09-2026) saat modul mendapat layar pertamanya, di folder
`backend/migrations/` modul ini sendiri — bentuk SQL-nya di `APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`
bab 6. Nomor selalu tiga digit.

## Pernyataan untuk penjaga

⛔ **Dibaca penjaga** `inti/backend/penjaga` — satu jenis pernyataan per judul `###`, satu baris per butir (pola
MODUL.md Claim Life). Sebelum 07-10-2026 kedua judul di bawah berada di bab Migrasi sehingga TIDAK terbaca penjaga
(mutasi 532 tanpa CASCADE tetap hijau); dipindah ke bab ini.

### Kaskade ON DELETE CASCADE

Kaskade HANYA pada berkas migrasi modul ini yang berawalan di bawah; berkas lain modul ini tanpa
`ON DELETE CASCADE` (`TestKaskadeHanyaPadaRelasiTerdaftar`).

| Awalan berkas | Relasi |
| --- | --- |
| `521_` | T_CLAIM_ESTIMATION.CLAIM_ID -> T_GENERAL_CLAIM - `.ClaimData.EstimationList` hidup di halaman klaim |
| `522_` | T_CLAIM_INTEREST.CLAIM_ID -> T_GENERAL_CLAIM - `.ClaimData.InterestList` |
| `523_` | T_CLAIM_CLAIM_AMOUNT.CLAIM_ID -> T_GENERAL_CLAIM - `.ClaimData.ListClaimAmount` |
| `524_` | T_CLAIM_LOSS_ALLOCATION.CLAIM_ID -> T_GENERAL_CLAIM - `.ClaimData.SpreadingRisk` |
| `525_` | T_CLAIM_SPREADING.CLAIM_ID -> T_GENERAL_CLAIM - `.ClaimData.SpreadingClaim` |
| `526_` | T_CLAIM_BREAK_QS.CLAIM_ID -> T_GENERAL_CLAIM - `.ClaimData.SpreadingBreakQS` |
| `527_` | T_CLAIM_FAC_RETRO.CLAIM_ID -> T_GENERAL_CLAIM - `.ClaimData.FacRetroList` |
| `528_` | T_CLAIM_ADJUSTMENT.CLAIM_ID -> T_GENERAL_CLAIM - `.ClaimData.AdjustmentList` (baris berkomite tidak pernah dihapus aplikasi) |
| `529_` | T_CLAIM_ADJ_SPREADING.ADJUSTMENT_ID -> T_CLAIM_ADJUSTMENT - `.AdjustmentList(n).SpreadingAdjustment` |
| `530_` | T_CLAIM_ADJ_QUOTA_SHARE.ADJUSTMENT_ID -> T_CLAIM_ADJUSTMENT - `.AdjustmentList(n).SpreadingQuotaShare` |
| `531_` | T_CLAIM_ADJ_LOSS_ALLOCATION.ADJUSTMENT_ID -> T_CLAIM_ADJUSTMENT - `.AdjustmentList(n).LossAllocation` |
| `532_` | T_VIEW_SUGGEST.CLAIM_ID -> T_GENERAL_CLAIM - riwayat `.ClaimData.SuggestList` (induk kedua, CHECK tepat satu) |

### Nama terlarang di migrasi

Nama yang tidak boleh muncul di migrasi modul MANA PUN (`TestNamaYangDibuangTidakAda`).

| Nama | Sebab |
| --- | --- |
| `FLAG_ON_GOING_COMMITTEE` | `AddKomiteTreatyChild_ACT` langkah 3 `.FlagOnGoingCommitte` dibuang (keputusan work owner 19-09-2026) - tidak ada kolom, tidak ada properti tersimpan |

## Keputusan work owner 07-10-2026 sesudah laporan pertama

- **Migrasi `532` (`T_VIEW_SUGGEST.CLAIM_ID` + CHECK satu induk) dipasang** — "1 tabel aja gabung life dan non life": Claim History satu tabel
  dengan riwayat penawaran PremiumList Life. Baris kolomnya di STRUKTUR PremiumList Life; uji modul itu
  (`strukturtipe_polis_test.go`) menghitung 244 kolom dan melewati kolom buatan modul lain (`kolomModulLain`).
- **Baris kolom 520 di STRUKTUR Claim Life** (`## T_GENERAL_CLAIM`, dokumen saja) = keputusan "Tabel bersama".
- **Penyerahan ke komite aktif** (jawaban "b"): kode batas komite dikumpulkan di `backend/models/komite.go` dan
  `backend/repository/komite.go`; keduanya masuk daftar pengecualian penjaga batas Claim Life `komite_statik_test.go`
  (satu-satunya suntingan di modul itu). Gerbang lampiran tetap menolak selama OQ-CP-12 terbuka.
- **Kotak masuk Komite Claim Life disaring LINI** ("tambahkan!"): `komiteclaimlife/backend/repository/komite_inbox.go`
  `sqlSaringInboxKomite` + `(w.LINI = :lini OR w.LINI IS NULL)`, argumen `inti.LiniLife`. Kasus komite PROP tidak
  muncul di inbox Komite Life. Pembaca satu kasus `sqlKasusKomite` disaring sama (membuka, memutuskan, riwayat).
- **Spreading klaim** (08-10-2026, "Add pada spreading kok ga bisa?" lalu "tabel anak tidak muncul" -> pilihan work
  owner; dikerjakan sesi ASIS CLAIM PROP, ditinjau dan di-commit sesi CLAIM PROP). SpreadingClaim terisi otomatis
  saat Policy No dipilih (`models.IsiSpreadingPolis` sesudah CheckNoPolicy) dari TREATYINPRODUCTION polis (JN_REAS,
  PCT_SHARE_PREMI, CURR_ID -> CURRENCY.ID). SpreadingBreakQS = `TreatyInMaster.Limits(1).Detail(1).SpreadingList`
  (SetTreatyNameSpreading_Act langkah 11-15, XML); percobaan PROPORTIONALARRG (24a6e315) dibatalkan sesudah uji data
  DEV - master cocok dengan seluruh SpreadingBreakQS klaim CLMP lama, PROPORTIONALARRG tidak. Add / Delete AKTIF
  (membatalkan "non aktifkan" 07-10-2026; AddSpreading_Act / DeleteSpreading_Act tidak diekspor, perilaku ditetapkan
  work owner): Add LANGSUNG mengisi Treaty Type (hanya-baca karena terisi; "begitu add langsung set spreading type nya
  dan readonly") - satu klik = satu baris per mata uang yang belum ada untuk treaty spreading polis pertama yang masih
  kurang ("jika add langsung kedetek 2 currency, langsung add 2 mengikuti currency"), Share diisi manual; spreading
  sama (Treaty Type + Currency) ditolak dengan pesan, juga saat memilih di baris kosong ("tidak ada spreading sama ...
  kecuali currency beda"); tanpa spreading polis baris kosong dan Treaty Type dipilih; layar menggulir ke kotak pesan.
  Delete nonaktif bila `IsOldData='Yes'` (XML); tabel bawah dan turunan disusun ulang sesudah setiap perubahan. Dropdown
  Treaty Type (kedua layar) = spreading polis + treaty baris yang sudah ada: koreksi atas `[dugaan]` SpreadingList
  master yang dibantah data DEV (TreatyType klaim lama = induk, SpreadingList = anak).
- **`STS_REJECT = 1` ditunda** sampai modul Komite Claim Prop ("itu nanti kan dari komite") — pemuat mencatatnya
  "ditunda", bukan gagal.
- **Halaman awal dua tab, rupa Kelola User** ("cuman ada 2 tab process dan resolve"): tab Process bawaan = worklist
  pembuat (Assignment2 Outstanding Claim, daftar `saya`, tanpa cek workbasket - XML `ToCurrentOperator` apa adanya;
  `ReasKlaimAdmin` dibuang - migrasi 534). Dropdown Admin / Teknik diganti **switch Teknik** (08-10-2026: "admin nya
  buang ... seperti toggle ... kalo wb nya ada ReasKlaimTeknik baru switch nya di aktifkan"): nyala = Assignment1
  Input Acceptation (daftar `workbasket`); dapat dinyalakan hanya anggota `ReasKlaimTeknik` - halaman membacanya dari
  `GET /api/claim-prop/hak` karena sesi frontend inti hanya meloloskan peran Life. Tab Resolve = daftar `selesai`. Add
  Claim hanya di Process saat switch mati (`pages/inbox.ts`). Tampil "Inbox ( ) Technical" (label Inggris); switch dan
  kapsul tab aktif bergerak halus (`transition` 200-220 ms, tanpa `transform`; mati bila prefers-reduced-motion).
  Medan, tombol, tabel, toolbar memakai kelas inti (`field__input`, `btn`, `inbox__tabel`, `toolbar`);
  `claimprop.css` hanya menata letak pohon tata, warna lewat token inti.
- **Workbasket `ReasKlaimAdmin` dibuang** (08-10-2026, "hapus dari master wb ReasKlaimAdmin dan akun yang pake ini"):
  migrasi `534` menghapus barisnya di `M_LOGIN_GO_WORKBASKET` (DEV: 3 pemegang) lalu di `M_WORKBASKET`; pola
  Bordereaux 893. Pernyataannya diurai di DEV (`DBMS_SQL.PARSE`, tanpa eksekusi); berjalan saat work owner `-migrate`.
- **Popup Choose Master** (08-10-2026, "kenapa ga bisa filter dan cuman ada 1 master" -> "lanjut"): saringan XML
  `PROPORTIONTYPE = 'Proportional'` (parameter section MasterTreatyInList ke RD BrowseCLAIM_MASTER_TREATY) dipasang;
  filter per kolom digabung AND di server (`repository.sqlDaftarMaster`), urut Treaty Year lalu Treaty ID menurun,
  batas 500 (`models.BatasMaster`), 50 per halaman di layar. Sebab lama: 20 baris pertama urut TREATYID = satu
  treaty (rata-rata 39 baris per treaty di view CLAIM_MASTER_TREATY). DEV: 0,37 s tanpa filter, 0,05 s berfilter.
- **Tombol View polis** (08-10-2026, "buat view untuk polis itu sesuai dengan modul pada nb/edm treaty"; diralat
  hari yang sama: "jangan tab baru", "biarkan di layar utama", "hanya tampilan polisnya aja, ga usah sampe menu
  menunya ikut kebuka"): harness DetailPolisRealization tidak diekspor, tetapi `GetDetailPolis_act` mengurai dokumen
  polis ke kelas berkas polis Treaty In. View mencari berkasnya - `GET /api/claim-prop/berkas-polis`
  (T_GENERAL_POLIS_TREATY.NOPOLIS + T_WORK_POLIS, PRODKE terbesar; 0 = nbtreatyin, >= 1 = edmtreatyin) - lalu
  `PropsRute.onLihatBerkas` (inti `modul.ts`): `frontend/App.tsx` memasang rute modul NB / EDM sekali lagi di modal
  selebar layar (`Modal penuh`) di atas layar Claim Prop, berkas dibuka lewat jalur kotak masuk Beranda (`bukaKasus`);
  tanpa sidebar / topbar, tanpa tab baru, layar klaim tetap. Back berkas atau X menutupnya. Tab baru dibatalkan: di
  server dev vite tab baru memuat ulang seluruh aplikasi (22-30 detik layar putih), dan peramban tidak mengizinkan
  halaman membuka tab di latar. Tanpa berkas (404) / modul tidak dipasang bagi akun = pesan dalam modal. Modul NB /
  EDM tidak diubah. DEV: hanya berkas aplikasi baru / Copy Old punya berkas (15 dari 41.945 baris TREATYINPRODUCTION).
- **Panel AdjustmentDetail hanya di grid Adjustment** (08-10-2026, "kenapa ada itu?" - screenshot grid Insured
  Interests / Total in Original Currency): panel rinci baris dulu menempel di semua grid, jadi klik baris mana pun
  membuka "AdjustmentDetail" kosong. Kini hanya grid `ClaimData.AdjustmentList` (expand pane XML) yang barisnya dapat
  dibuka (`frontend/components/rincian.ts`, diuji).
- **Tampilan medan** (08-10-2026, perintah work owner; dikerjakan sesi ASIS CLAIM PROP, ditinjau dan di-commit sesi
  CLAIM PROP): Consultant / Adjuster dropdown memilih dan menampilkan nama (`Tata.Tampilan`, saringan NAME saja), ID
  tetap disimpan; isian angka hanya angka dengan separator Indonesia, paling banyak 4 desimal (`ketikAngka.ts`,
  `InputAngka.tsx`, `nilai.tampilAngka`); isian tanggal diketik `dd-mm-yyyy` (+ `hh:mm`) dengan tombol kalender, aksi
  server hanya saat lengkap dan sah (`ketikTanggal.ts`, `InputTanggal.tsx`); label Type Estimation List dari
  screenshot (`LabelKode["EstimationType"]`); nilai tampil di sel tabel tidak lagi terpecah satu huruf per baris.
- **Reporter Address dari CLIENT_ADDRESS** (08-10-2026, "ubah jangan dari json, ambil dari client address"; dikerjakan
  sesi ASIS CLAIM PROP): `GetAddressCeding` membaca tabel datar CLIENT_ADDRESS (ASMADDRESS, RWNAME, DISTRICTNAME,
  CITYNAME) klien milik agen ceding, bukan `M_CLIENT.JSONDATA` - baris Kantor (tipe 2) dulu, tanpa Email (tipe 7).
  Nama ceding / SOB tetap dari JSON master treaty.
- **Perbaikan Claim Information** (08-10-2026, daftar work owner): Date of Loss / Received Date tanggal saja
  (KTanggal, katalog kTgl); Reporter Phone Number kendali `telepon` (hanya angka, nol di depan tetap, kolom teks);
  label kode `models.LabelKode` - Report Type 1 Direct / 2 Via Email / 3 Via Fax / 4 via Postal Mail/Courier / 5 Via
  Telephone, Reporter Status 1 Ceding Co Name / 2 SOB Name / 3 Others (nilai tersimpan tetap kode);
  PeriodPolicyTBA berlabel "Policy Period TBA ?"; Consultant / Adjuster / Province = dropdown yang dapat dicari
  (PilihSaring inti, saringan server). Catastrophe dan Edit RNM Share diperbaiki: penanda mode `EditCatastrope` /
  `IsEditRNMShare` tidak punya kolom sehingga hilang saat halaman dibaca ulang (aksi 409) - kini `Layar.mode` dikirim
  server dan dikembalikan layar di setiap aksi (`models.ModeLayar`, daftar putih nilai true/false).
- **Tombol "+" Consultant / Adjuster** (08-10-2026, pilihan "Popup di Claim Prop"): popup `TambahAdjuster`
  (label section MstAdjusterConsultant) menyimpan master lewat API modul Adjuster Consultant - rute pinjaman
  `"POST /api/adjuster-consultant": {"claimprop"}` di `cmd/api/rakit.go` (dijaga `TestPanggilanLintasModulTerdaftar`).
  Modul itu membuat ID dan menolak nama kembar (422, pesannya tampil). Sesudah tersimpan layar menjalankan
  SetConsultant / SetAdjsuter dengan ID baru. Tombol nonaktif selama ID hanya-baca (IsOutstanding = 1).
- **Layar kasus berkulit Kelola User** (08-10-2026, "SAMAIN DENGAN MENU KELOLA USER SKIN NYA"): setiap bagian
  berlabel = kartu `panel` + `panel__title`, isi `form-grid`, medan berlabel di atas kotak `field__input`, hanya-baca
  tetap berkotak (`field__input--readonly`). Pengelompokan di `frontend/components/susun.ts` (diuji): label
  pendamping `Q` / `/` / `U/Y` menjadi label medannya, `%` menjadi satuan, tombol sesudah medan menempel di kanan
  kotaknya, tombol akhir = baris aksi kanan (Save / Submit = `btn--primary`). Isi tata server tidak berubah.
- **Kulit Kelola User dari inbox sampai isi kasus** (08-10-2026, "SEMUA DONG DARI LUAR SAMPE DALAM"): nilai tema
  `.kelola-user` inti disalin ke token `--cp-*` di akar `.claimprop__akar` (terang + gelap), pola Marketing Officer -
  tidak menumpang ke kelas inti. Akar dipasang di halaman awal dan layar kasus (popup / Modal ikut karena tanpa
  portal); semua tabel `claimprop__tabel` (kartu timbul berkepala navy), tab berjalur cekung, tombol kapsul, tombol
  utama merah bergradasi, isian cekung, kartu bagian timbul. Warna hanya di blok token (dijaga gaya.test.ts).
- **Layout lama Pega dirapikan** (08-10-2026, "berikut layout lama, ikuti dan rapihkan" + "font juga kecilin"): server
  mengirim format layout XML di `Tata.Letak` - `dua` (Inline grid double), `sebaris` (Inline / Inline labels left),
  `tab` (layout group Tab), `judul` (kepala, label di sel tengah) - serta `Ikon` (pi-plus / pi-trash / pi-pencil /
  pi-check / pyWorkActionsAddWork) dan `PerHalaman` (Claim History paging 5). Medan = Stacked with labels left (label
  kiri 150px, hanya-baca = teks). Claim Treaty 8|5, Claim Information kiri polis / kanan pelapor, baris
  "Quarter/Year" (label dari XML), Interest / Estimation / Spreading = tab, RNM Share di atas tab. Huruf modul 13px,
  isian dan tombol 32px. "View Master" di screenshot Pega hidup tidak ada di XML - tidak dibangun.
- **Pesan pra-proses tidak tampil saat kasus dibuat / dibuka** (08-10-2026, "BARU BUAT UDAH ADA WARNING
  ERROR"): `BukaKasus` membersihkan pesan `CheeckNoRNM_Act` (termasuk ProteksiData langkah 12) sesudah pra-proses;
  bendera Protect / IsError tetap dihitung. Pesan ProteksiData tampil pada Save to issue RNM / Submit.

## Pemuat data lama

`backend/alat/pemuatlama` (prompt §6 butir 11; AC 9–12, 123, 132). Sumber baca-saja: `OS_AKSEPTASI_KLAIM`
(CASEID `ASM-FW-GCNMFW-WORK CLMP-%`, baris berlaku per kasus menurut AC 123) dan `JSON_KLAIM` (halaman `.ClaimData`
bila ada). Tujuan: `T_WORK_CLAIM` + `T_GENERAL_CLAIM` (SUMBER `PEGA`, ID = pyID Pega `CLMP-n`) dan tabel `T_CLAIM_*`
lewat `SimpanHalaman` — jalur yang sama dengan aplikasi.

```
go run ./modul/claimprop/backend/alat/pemuatlama -keluaran <folder>             # uji-kering (bawaan, baca saja)
go run ./modul/claimprop/backend/alat/pemuatlama -keluaran <folder> -jalankan   # tulis - HANYA work owner
```

- `-jalankan` ditolak bila `IS_PEGA_PROD=true`; satu transaksi per kasus; kasus yang ID-nya sudah ada dilewati.
- Berkas keluaran: `claimprop-arsip-medan-*.csv` (medan yang tidak masuk kolom + sebab) dan `claimprop-galat-*.csv`.
  Keduanya memuat data kasus — simpan di luar repositori.
- Uji-kering DEV 07-10-2026: 2.451 kasus, 6.596 baris OS (4.145 riwayat), 6 kasus berhalaman JSON, **2.032 siap**,
  **419 ditunda** (baris berlaku `STS_REJECT = 1`, menunggu modul Komite Claim Prop), **0 gagal**; keluar 0.
- Urutan resmi: work owner menjalankan `-migrate` (520–534, 980) → uji-kering ulang → keputusan OQ-CP-18 →
  `-jalankan` oleh work owner / DBA. Tidak pernah dijalankan agen.

## Uji SQL di DEV (baca saja)

`go test -tags ujidev -run TestSQLDiDEV -v ./modul/claimprop/backend/repository/` — setiap SELECT dijalankan lewat
metode aslinya dengan masukan `UJI-*`, setiap INSERT / UPDATE / DELETE hanya diurai `DBMS_SQL.PARSE`. 07-10-2026:
58 ok, 0 gagal, 38 "objek belum ada" (35 menunggu migrasi 520–533; 3 skema luar OQ-CP-11).
