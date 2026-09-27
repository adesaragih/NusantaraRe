# PROMPT — lanjutan 9: **F0.6 form login ditiadakan → Register dirampungkan menurut XML → Outstanding → Detail & Tutup**, satu giliran; keputusan **av** *(menunggu modul PremiumList Life)* dan **aw** *(Send Back to Register)*

> Dibaca sesudah brief modul dan lanjutan 1–8; semuanya tetap berlaku kecuali yang diralat di sini.
> Mekanisme giliran = lanjutan 8 §1 **tanpa perubahan** *(ia bekerja: satu pesan sesudah A3 Register)*.
> Perintah work owner 27 September 2026, tambahan hari ini: **"FORM LOGIN DITIADAKAN DULU"** dan, atas
> pertanyaan sumber sepuluh medan polis, **"YA MENUNGGU MODULNYA! UNTUK XML NYA SUDAH ADA DI FOLDER
> REFERENSI!"** — XML modul itu ada di `D:\XML\RNM_BRD\PremiumList Life\` *(READ-ONLY, untuk modulnya
> kelak, bukan untuk dibaca datanya dari Claim Life sekarang)*.

---

## 0. KEADAAN SESUDAH `82c61dd` — F0.3 → A3 REGISTER DIVERIFIKASI ULANG

| Klaim laporan | Diperiksa ulang | Hasil |
| --- | --- | --- |
| 7 commit sejak `5b49954` | `git rev-list --count` = 7 *(6 kode + 1 docs)* | ✅ |
| Go 274 · 0 · 34; JS 136; build 46 modul | dijalankan ulang: **274 PASS · 0 FAIL · 34 SKIP**; vitest 9 berkas **136**; **46** modul | ✅ |
| migrasi 016 `TAHAP` + `TGL_CREATE`; ralat tiket 08 & 14 | `016_kolom_tahap_dan_tgl_create.sql` + `_down`; tiket 08 baris 122, tiket 14 baris 1659 | ✅ |
| menu dua butir berbukti, penjaga menjatuhkan penambahan | `MENU` = `Claim Life` *(270)*, `Inbox Claim Life` `[tidak ada di korpus]`, `Register` *(155)* | ✅ |
| ralat angka commit `6b1618c` dicatat di `LAPORAN-GILIRAN-F0.md` | baris 67–69 | ✅ |
| nol kebocoran | grep seluruh berkas yang berubah: nol | ✅ |

**Dua penyimpangan dari XML yang ditemukan saat membaca ulang `InputRegisterClaimLife.xml` — diralat
di §2**: *(a)* ketiga medan `Ceding`, `Class of Business`, `Marketing Officer` di Pega **read-only**
*(`pyEditOptions Read-only`, `pyReadOnly true`, format `pxTextInput`, `pyLabelFor` `CedingCoName` /
`BusinessName` / `MarketingName`; blok 11500–11720, 9890, 12171)* dan terikat `.PolicyDataLife.*`,
sedangkan `repository/rujukan.go` + rute `GET /api/rujukan/{jenis}` + `api.ts cariRujukan` dibangun
sebagai "tiga isian ber-autocomplete" — dan **nol** berkas `.tsx` memakainya *(kode mati di dua
lapis)*; *(b)* PARITAS baris 226 menulis `SearchPolicyHolder_act` = "rute + kontrol lewat
`/api/rujukan/ceding` dan dropdown-nya" — padahal activity itu hanya `Property-Set
SearchPolicyHolder.CARI1 = @toUpperCase(...)`: teks pencarian popup `SearchPolicy_Harness`, bukan ceding.

---

## 1. KEPUTUSAN **av** `[DIPUTUSKAN work owner — 27 September 2026]`: sepuluh medan polis **menunggu modul PremiumList Life**

| Butir | Ketentuan |
| --- | --- |
| Sumber di Pega | `PolicyDataLife` diisi dari **kasus PremiumList Life** *(class `ASM-FW-GISFW-Work-LIFE`; popup `SearchPolicy_Harness` memakai `InboxPremiumList_Claim` atas class itu)*. Mirror Oracle-nya *(`JSON_POLIS`, view `POLICYJSONLIFE`, `SEARCH_POLIS`)* ditulis modul PremiumList Life *(`PremiumList Life\Activity\InsertJsonPolisLife_Act.xml` → prosedur `INSERTJSONPOLISLIFE`)* |
| Keputusan | Claim Life **tidak** membaca `JSON_POLIS` / `POLICYJSONLIFE` / `SEARCH_POLIS` / `pc_*`. Medan polis diisi kelak dari tabel relasional modul PremiumList Life *(tiket 02: "dibaca hidup dari tabel polis, `T_PREMIUM_LIST` dkk")* |
| Di layar sekarang | `PanelDataPolis` **tetap** menampilkan seluruh medan VERBATIM dengan penanda; teks penandanya menjadi **"menunggu modul PremiumList Life"** *(bukan "belum bersumber" tanpa sebab)*; uji menghitung cacahnya |
| Popup `Choose Policy No` | *(3776/3827 → `SearchPolicy_Harness`, daftar kasus PremiumList Life)* → tombol ada, membuka `BelumTersedia` yang menyebut `SearchPolicy_Harness` dan modul yang ditunggu. `Find Insured` / `GET /api/peserta-life` atas `M_LIFE_PREMIUM_DETAIL` **tetap** *(keputusan tiket 02)* |
| Aksi yang ikut menunggu | `PreCaimLife_Act` *(BUSINESS by `OLDID` = `PolicyDataLife.BusinessCode` → `BusinessID`)* dan `SetMOClaim_Act` *(`MarketingData.*` ← `PolicyDataLife.MOID/MarketingCode/MarketingName/TeamGroup/BranchCode/BranchName`, fallback `Obj-Browse` `DataMarketing`)* di Outstanding bergantung `PolicyDataLife` → ditandai **menunggu modul** di PARITAS dengan bukti langkahnya, bukan dihilangkan |
| Yang dibuang | rute `GET /api/rujukan/{jenis}`, `handlers/rujukan.go`, `api.ts` `cariRujukan`/`JENIS_RUJUKAN`/`BarisRujukan` — kode mati. `repository/rujukan.go` **dihapus** juga, kecuali dipakai ulang giliran ini sebagai pencarian `BUSINESS` ber-`OLDID` untuk `PreCaimLife_Act` *(dan itu pun menunggu modul)* — pilih satu, nyatakan |
| Tiket | ralat **tiket 02** *(blok bertanggal: medan polis read-only dari `PolicyDataLife`; sumber menunggu modul PremiumList Life; bukti baris)*; PARITAS §5.3 baris 226 diralat; `OQ-untuk-tim.md`: urutan modul — PremiumList Life dibutuhkan Register Claim Life *(keputusan work owner)* |

⛔ `setDetailClaim_act` *(tombol `Choose` popup, `SearchPolicy_Section.xml` 3337–3356, 3468)* adalah
**residu uji pengembang**, bukan logika bisnis: `Obj-Open-By-Handle` pada **satu instance handle
literal** *(pengenal kasus tertulis mati di rule; jangan dikutip)*, `Property-Set` tanggal literal
Januari–Februari 2026, prasyarat yang
membandingkan `.NAME_OF_INSURED` dengan **satu nama literal** *(jangan dikutip)*, lalu `Obj-Save`.
**Tidak ditiru**; PARITAS baris 99 *(no. 28, "A3 — Tutup & lihat, MENULIS")* diralat menjadi
"tidak ditiru — residu uji, bukti langkah"; dilaporkan di `OQ-untuk-tim.md` untuk pengembang Pega.

---

## 2. F0.6 — FORM LOGIN DITIADAKAN `[perintah work owner 27-09-2026]`

| Bagian | Ketentuan |
| --- | --- |
| Halaman `Masuk` | **dibuang** beserta ujinya. Aplikasi membuka **Inbox** langsung |
| Identitas stub | dibentuk saat aplikasi menyala dari env Vite: `VITE_STUB_PELAKU` *(bawaan `UJI-ADMIN`)* dan `VITE_STUB_PERAN` *(dipisah koma; bawaan ketiga peran ADR-U-0002)*, **hanya** bila `VITE_AUTH_STUB=true`; tanpa itu tidak ada identitas dan setiap permintaan ditolak backend — gagal tertutup seperti sekarang. `store/sesi` tetap bentuknya *(tiket 07 / ADR-U-0030 kelak mengisinya dari IAM)*; `sessionStorage` tidak lagi dipakai untuk identitas stub |
| Topbar | tetap menampilkan akun + peran *(dari env)*; tombol `Keluar` **dibuang** *(tanpa masuk tidak ada keluar)*; pita "mode stub" tetap |
| Uji | `Shell.test`/`App`: tanpa `Masuk`, Inbox dirender langsung; header `X-Pelaku`/`X-Peran` terkirim dari env; `VITE_AUTH_STUB` bukan `true` → layar menyatakan identitas tidak ada, bukan pecah |
| Dokumen | `.env.example` frontend *(dua variabel baru)*, `PANDUAN-MENJALANKAN.txt` bab 4, `LAPORAN-GILIRAN-F0.md` |

Commit: `frontend: kerangka — F0.6 form login ditiadakan, identitas stub dari env`.

---

## 3. REGISTER DIRAMPUNGKAN MENURUT XML *(commit `claim-life: A3 — Register (2)`)*

| Aksi | Bukti langkah *(dibaca ulang hari ini)* | Yang dibangun |
| --- | --- | --- |
| Medan read-only | §0 *(a)* | ketiga medan + `Type` tampil **read-only** dari `PolicyDataLife` → "menunggu modul"; dropdown dan rutenya dibuang |
| `SaveInsuredClaim_Act` *(tombol `Select Insured` 20008/20060)* | `SearchInput.CARI3=""`; tiap baris `TempDetail1.pxResults` → `TempDetail.pxResults(<APPEND>)`: `GROSS_PREMIUM`, `NET_PREMIUM`, `SUM_INSURED`, `CEDING_RETENTION` = `@divide(@toDecimal(@replaceAll(x,",",".")),1,4)`; `SHARE_NUSANTARA_RE` = itu, **fallback** `SHARE_NUSANTARA_RE_GROSS` bila `"0"` atau `""`; `AGE` = `AGE`, else `ENTRY_AGE`, else `CURRENT_AGE`; `BEGIN_DATE`, `CERTIFICATE_NO`, `DOB`, … disalin; `local.RETROBEGINDATE/RETROEXPIREDDATE` | memindahkan peserta terpilih ke daftar klaim **di layar** *(sebelum `POST`)*; aturan fallback share dan umur **diuji**; bandingkan dengan `kolomSalin` tiket 02/03 — bila backend sudah menyalin kolom yang sama saat `POST`, layar tidak menghitung ulang uang *(ADR-U-0003)*, hanya memilih |
| `DeletePesertaClaimLife` | `Local.IndexPremium = .pxListSubscript`; `.AdjustmentList.IndexPremiumList` disesuaikan; `Obj-Save` pyWorkPage | sebelum `POST` = hapus dari pilihan di layar; **sesudah klaim ada** *(tahap Input Register)* = rute `DELETE /api/klaim-life/{id}/peserta/{pesertaId}` bergerbang tahap + pemegang *(services)*, ikut menghapus baris adjustment peserta itu *(bukti `IndexPremiumList`)*, jejak audit; uji |
| `LoadDataPesertaSpesifik_Act` *(`Find Insured` 7057/7111 dengan nama)* | `Page-Remove TempDetail1`; `SearchInput.CARI1=""`, `CARI3=1`; `SearchPolicyHolder.CARI3` di-uppercase; `RDB-List` `GetPesertaClaim_sql1` *(class Work-ClaimLife)*; lalu tiap tanggal `@addCalendar(x,0,0,0,0,7,0,0)` | `GET /api/peserta-life?pl=&nama=` *(nama di-uppercase, `LIKE` **berbatas**, hanya bersama `pl` — indeks tabel 66,8 juta baris)*; SQL dibaca dari `GetPesertaClaim_sql1.xml`; **+7 jam** = koreksi zona Pega *(GMT → WIB)*, **tidak ditiru** *(Oracle `DATE` kita tanpa zona)*, dicatat |
| `SelectAllClaimLife_act` *(`Select All` 24489)* | `Select.CARI1` bergantian `"true"/"false"`; setiap `PremiumListDetail.IsAccept = Select.CARI1` | tombol memilih/melepas seluruh baris hasil pencarian di layar; uji |
| `ValidasiClaimReceived_Act`, `ValidasiSTNC_Act` | sheet baris **115, 118**: dipanggil dari **`ClaimLifeDetailGCNM`**, bukan Register → PARITAS kelompoknya diralat ke **Detail & Tutup**. Isi: `RDB-List GetProductName` *(`PRODUCTINWARD_LIFE`)* → `MAXEXPIREDCLAIM` / `MAXDATARECEIVE`; selisih hari `@DateTimeDifference(DOL, CLAIM_RECEIVED_DATE,"D")` resp. `(EFFECTIVE_DATE, RECEIVED_DATE)`; bila melebihi → `.MAXCLAIM_RECEIVED` / `.STNC` = tanggal terformat, else `""`; `local.errmsg` `"Max Claim invalid"` / `"STNC invalid"` | dikerjakan di §5 *(Detail)*, kolom `STNC_CLAIM` sudah ada di `T_GENERAL_CLAIM` |
| `UploadCSV_ClaimLife` *(flow action 12)* | `UploadCSVClaimLife_Act`, `pxUploadCSVResults` | dibaca sebagai pohon; bila menulis peserta dari CSV → rute unggah berbatas + validasi kolom; bila hanya UI Pega → dinyatakan |

---

## 4. OUTSTANDING *(commit `claim-life: A3 — Outstanding`)* — jawaban pohon dan keputusan **aw**

### 4.1 Peta penyambung `Register_Flow.xml` *(12 `Embed-Rule-Obj-Flow-FromToTask`, dibaca hari ini)*

| Penyambung | Syarat | `pyPosition` yang ditulis |
| --- | --- | --- |
| keluar Assignment2 lewat flow action `InputRegisterClaimLife` | STATUS | `ReasLifeAdmin` |
| keluar Assignment1 lewat `OSClaimLife` *(`Transition4`)* | STATUS | `ReasLifeMedicalAdvisor` *(sebelum `Decision3`)* |
| `IsSendtoAdmin` *(tiga cabang: Decision3 `Transition11`, Decision1, Decision2)* | WHEN | `ReasLifeAdmin` |
| `IsSendtoMedical` *(Decision2)* | WHEN | `ReasLifeMedicalAdvisor` |
| keluar Assignment3 lewat `MedicalCheck` | STATUS | `ReasLifeSPV` |
| keluar Assignment4 lewat `AkseptasiClaimLife` | STATUS | — |
| tiga `ELSE` + satu `ALWAYS` | — | — |

Jadi **di Outstanding `pyPosition` = `ReasLifeAdmin`**. Tombol `Send Back to Register`
*(`InputOSClaimLife.xml` 21404)* → `pyLocalAction SendtoAdmin` *(21433)* → `SendtoAdmin_Section`
tombol `Submit` → `SendtoAdmin_Act`: `Property-Set pyWorkPage.SendtoAdmin="1"` **hanya bila**
`pyPosition=="ReasLifeMedicalAdvisor"` *(WhenTrue=2, WhenFalse=3 = lewati)*. Tombol `Send to Medical
Check` *(21839)* → `refresh` + `SendtoAdmin_Act1` *(`SendtoAdmin="0"`, prasyarat **sama**)* +
`finishAssignment`. Penulis `SendtoAdmin`/`SendtoMedical` di seluruh modul hanya **tiga** activity
*(`SendtoAdmin_Act`, `SendtoAdmin_Act1`, `SendtoMedical_Act`)* dan **tidak satu pun** berjalan pada
posisi Admin.

**Akibat menurut XML apa adanya:** kasus baru tidak pernah dapat dikembalikan ke Input Register dari
Outstanding; kasus yang pulang dari Medical/SPV membawa `SendtoAdmin=1` yang **tidak pernah direset**,
sehingga `Decision3` mengirimnya ke Input Register **tombol apa pun** yang ditekan, lalu Register →
Outstanding → Register lagi. Itu **cacat rule warisan** *(prasyarat terbalik atau salah tempel)*, bukan
maksud bisnis: label tombol, penyambung `Decision3 → Assignment2`, dan ADR-U-0002 menyebut jalur
balik ini sebagai fitur.

### 4.2 Keputusan **aw** `[DIPUTUSKAN — 27 September 2026, dari maksud XML yang terang; cacat dilaporkan; veto work owner]`

`Send Back to Register` = `Pindah` ke **Input Register** dengan `SENDTO_ADMIN="1"`, jejak;
`Send to Medical Check` = `Pindah` ke **Medical Check** dengan `SENDTO_ADMIN` dikosongkan *(maksud
`SendtoAdmin_Act1` = `"0"`)*; keduanya bergerbang pemegang tahap Outstanding. Ralat **tiket 08**
*(blok bertanggal: bukti 4.1, cacat, keputusan)*; catat cacatnya di `OQ-untuk-tim.md` untuk pengembang
Pega dengan baris buktinya.

### 4.3 Cakupan Outstanding *(brief 6 §3 kelompok 2; PARITAS baris 20, flow action 2 & 4)*

`OSClaimLife` + `InputOSClaimLife` *(label VERBATIM: `Claim No` 951, `Type` 6898, `Marketing Officer`
7686, `Ceding` 9407, `Policy Holder` 9602, `Class of Business` 10036, `Date Received Email` 11250,
`Response Date` 11456, `Confirmation Date` 11662, `Status` 11868, `Updated Status` 12271, `Realization
Date` 12468, tombol `Send to Medical Check` 21349/21839, `Send Back to Register` 21404, `Close Claim`
22750/22808, `Select All` 24489)*; `PreCaimLife_Act`/`SetMOClaim_Act` *(menunggu modul — §1)*;
`SetClaimXOL_Act`; `RetroClaimLife` + `RetroDetailClaimLife`; `ViewClaimDetailLifeGCNM` *(membuka
panel Detail)*; `CloseClaim` *(`CloseClaim_Section`, `ProtectCloseClaim_act`)*; `SendtoAdmin` lokal.
Aksi layar OS *(`pyAction` unik)*: `refresh` 18, `postValue` 6, `localAction` 4, `setFocus` 3, `save` 2,
`finishAssignment` 2, `editItem` 2, `deleteRow` 2, `load` 1 — tiap `save`/`finishAssignment`/`deleteRow`
ditelusuri ke activity-nya. Baris Inbox tab *Outstanding Claim* kini membuka layar ini.

---

## 5. DETAIL & TUTUP *(commit `claim-life: A3 — Detail & Tutup`)* — selama giliran masih hidup

`ClaimLifeDetailGCNM` *(pohon; `Save Adjustment` 22590; `Diagnose_Harness` popup 5080/5224 —
`BelumTersedia` bernama sampai kelompok Medis)*, `ShowEditClaimLife` + `EditDateClaimLife_Section`,
`ValidasiDOL_Act`/`ValidasiSTNC_Act`/`ValidasiClaimReceived_Act` *(§3)*, `SetSTS_Reject`,
`SetIndexAdjustmentList`, `DocumentLife` *(tampil daftar; unggah/hapus = kelompok Dokumen)*. Rute
yang sudah ada *(akseptasi, tolak, putaran, tanggal kejadian)* mendapat kontrolnya.

---

## 6. LANGKAH 0 · 7. LAPORAN · 8. TELEMETRI

**Langkah 0**: `git add PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE-LANJUTAN-9.md` → commit `docs: brief
lanjutan 9 — login ditiadakan, av menunggu PremiumList Life, aw Send Back to Register` → `git status
--porcelain` kosong → uji hijau *(274 · 34 SKIP · 136 JS · 46 modul)*.

**Laporan akhir** persis lanjutan 8 §6, ditambah: daftar kode yang **dibuang** *(rujukan, Masuk)* dengan
sebab XML-nya; tabel **tombol Outstanding → aksi Pega → rute/kontrol kita**; daftar ralat tiket
*(02, 08, PARITAS)* dengan barisnya; cacah medan "menunggu modul PremiumList Life".

**Telemetri** persis lanjutan 8 §7 *(per paket di `LAPORAN-GILIRAN-F0.md`, ringkasan di laporan akhir,
taksiran token disebut taksiran)*.

---

*Disusun 27 September 2026 sesudah verifikasi `82c61dd` (build, vitest, `go test -tags=db`
dijalankan ulang; 7 commit dibaca statistiknya; migrasi 016, tiket 08/14, `LAPORAN-GILIRAN-F0.md`
dibaca), pembacaan ulang `InputRegisterClaimLife.xml` (kontrol read-only), `InputOSClaimLife.xml`
(tombol dan `pyAction`), `Register_Flow.xml` (12 penyambung, DOM), `SendtoAdmin_Act/_Act1/
SendtoMedical_Act`, `PreCaimLife_Act`, `SetMOClaim_Act`, `setDetailClaim_act`, `SearchPolicyHolder_act`,
enam activity Register, `PremiumList Life` (`InsertJsonPolisLife_Act`, tiga RDBList), dan katalog DEV
(`JSON_POLIS`, `POLICYJSONLIFE`, `SEARCH_POLIS`, `LIFE_PREMIUM_SUMMARY` — hanya definisi dan cacah).*
