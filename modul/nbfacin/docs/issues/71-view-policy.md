# 71: Menu View Policy (Facultative) — cari polis jadi, detail baca-saja (sesuai sistem baru)

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 09-10-2026 "jadi untuk cari case
> yang udh selesai darimana? dipega ada menu view policy, fungsinya untuk cari data polis yang udh jadi". XML diserahkan work
> owner ke `D:\migrasi\RNM\DDL\` (ViewOpenPolicyMarineInbox, GetListViewPolicy_act, GridOpenPolicyMarine,
> ViewEndorsementData_Act, BackViewPolicy_Dt, ResetFilterValue_DT). Riset: dua subagen sesi 0f, klaim inti diverifikasi ulang.

**Status:** kontrak backend dikirim ke sesi c3 09-10-2026; frontend sesi 0f.

## Bukti `[terverifikasi]`

- `DDL\ViewOpenPolicyMarineInbox.xml` (`DATA-PORTAL!VIEWOPENPOLICYMARINEINBOX`): mode daftar bertab Facultative (L1224) / Treaty
  (L17980) / Life (L35445); tab Facultative blok "List Polis": "View Option" `.FilterForViewRetro` ("Choose", wajib) + "Filter"
  `.FilterTermForPolicy` (placeholder "Policy No / ID"; Enter / tombol Filter → `GetListViewPolicy_act` SEARCHON "POLIFAC",
  L5691/L5969) + grid `GridOpenPolicyMarine` (L8594). Mode detail (container `IsShowOpenPolicyMarine`): Back, Refresh, "View
  Data" (versi; tampil bila >1), `Periode_IsUW` (L47714) + `ViewInwardFacultativeDtl_IsUW` baca-saja, View Facretro, riwayat
  akseptasi, Attachment File.
- `DDL\GetListViewPolicy_act.xml` blok 8 "POLIS FACIN": sama dengan `GetListViewFollowing_act` (tiket 70) KECUALI opsi 5 / 6 / 7
  cocok AWALAN: `INSUREDNAME LIKE 'x%'` (L2419), `SOB LIKE 'x%'` (L2581), `UPPER(QQNAME) LIKE 'x%'` (L2743); tanda kutip dibuang
  dari filter; tanpa cadangan json_following. SQL `BrowseDataViewPolis` (L2934) TIDAK ada di korpus — `BrowseDataViewPolis1`
  (FACINPRODUCTION) adalah salinan Save-As-nya `[dugaan kuat]`.
- `DDL\GridOpenPolicyMarine.xml`: `TempListPolicy.pxResults`, 20 per halaman, kolom ID · Policy Number · Reff. Number · Insured
  Name · Source of Business · QQ Name · Start Period · End Period · link "View Details" (sama dengan 9 kolom pertama
  GridViewFollowingNB); View Details → `SetDetailViewOpenPolicyMarine_PreDT` (TIDAK ada) → `GetPolicyData_ACT(prodke "0")` →
  `BrowsePolicyJSONMarine_Act`.
- "View Data": `GetProdKeOldData_SQL` = `select PRODKE … from pooldata.json_polis where NOPOLIS = … order by TGL_INPUT desc`;
  label = `substring(IDPEGA, 18)` (`GetPolicyData_ACT` L1262) = ID case versi itu; tampil bila versi > 1. Isi detail tiap versi
  di Pega = `JSON_POLIS.DATA_JSON` (adoptJSONObject) — JSON.
- Polis lahir di Pega TIDAK punya baris di tabel case sistem baru (`T_*`); datanya di `FACINPRODUCTION` / `CEDING_FACINPRODUCTION` /
  `JSON_POLIS`.

## Keputusan work owner 09-10-2026

- **K71-1** Detail = **campuran**: polis/versi yang case-nya ada di sistem baru → layar case baca-saja penuh (semua tab, riwayat,
  lampiran); polis lama Pega → ringkasan dari `FACINPRODUCTION` + `CEDING_FACINPRODUCTION`. Tanpa DATA_JSON (K70-1).
- **K71-2** "View Data": "tampilkan saja, kalau memang tidak ada endorsmentnya kan, brarti listnya cuma 1, tidak ada dropdown" —
  daftar versi dari kolom relasional `JSON_POLIS` (NOPOLIS / PRODKE / IDPEGA / TGL_INPUT; bukan DATA_JSON); dropdown hanya bila
  versi > 1.
- **K71-3** Akses: semua pengguna NB FacIn (baca saja).
- Lanjutan K70-3: View Option 3 No Slip (ARASAPAS) tetap ditunda.

## Keputusan agent (sesi 0f)

- **A71-1** Cakupan: tab Facultative "List Polis" saja; Treaty / Life (modul lain) dan "List Offer" tidak.
- **A71-2** Opsi 7 QQ Name: `UPPER(QQNAME) LIKE UPPER(filter)%` — Pega hanya meng-UPPER kolom (input huruf kecil tak pernah cocok);
  maksudnya tak peka huruf. Opsi 5 / 6 awalan apa adanya seperti Pega.
- **A71-3** Ringkasan polis lama: kepala (nomor polis, endorsement, ID case, tertanggung, QQ, SOB, ceding, periode, policy type,
  reff number, group, status bisnis, master polis) + baris FACINPRODUCTION (objek, lokasi risiko, item, coverage, akumulasi,
  mata uang, TSI 100 %, TSI RNM, rate, premi RNM, treaty, % share, RI commission) + daftar ceding; teks kolom sistem baru.
- **A71-4** Detail case sistem baru = halaman case yang sudah ada (case Resolved otomatis baca-saja), tombol Back kembali ke
  View Policy.
- **A71-5** View Facretro tidak dibangun (FacRetro = tiket 65).

## Backend sesi c3 (09-10-2026)

Berkas: `services/polis.go`, `repository/polis.go` (+ `CariPolisQQ` di `repository/following.go`), `models/polis.go`,
`handlers/polis.go`; layanan `DenganPolis(NewPolisOracle)`.

| Endpoint | Jawaban |
| --- | --- |
| `GET /api/nbfacin/polis?opsi=&filter=&halaman=` | `{baris, total, halaman, ukuran: 20}` (baris = BarisFollowing tiket 70) |
| `GET /api/nbfacin/polis/versi?nopolis=` | `{versi: [{prodke, caseId, sistemBaru}]}` |
| `GET /api/nbfacin/polis/ringkasan?nopolis=&prodke=` | `{kepala (14), baris (13), ceding [{cedingName, pctShare}]}` |

- **K71-B1 Cari** `[terverifikasi]` `DDL\GetListViewPolicy_act.xml`:
  - langkah 6 membuang `"` dan `'` dari filter, langkah 7 filter kosong -> keluar (lalu pangkas spasi);
  - blok 8 = klausa tiket 70 berurutan, KECUALI 8.7 `INSUREDNAME LIKE 'x%'`, 8.8 `SOB LIKE 'x%'` (awalan) dan 8.9 opsi 7
    `UPPER(QQNAME) LIKE UPPER(:1)` + `x%` (A71-2);
  - opsi 4 tetap `%x%`; opsi 3 / lain 400; bind;
  - SQL `BrowseDataViewPolis` tidak ada di korpus -> `BrowseDataViewPolis1` (FACINPRODUCTION) `[dugaan kuat]`;
  - tanpa penjaga case: hanya login (K71-3).
- **K71-B2 Versi** (GetProdKeOldData_SQL): kolom relasional JSON_POLIS NOPOLIS / PRODKE / IDPEGA, urut TGL_INPUT turun (lalu
  PRODKE turun, `[keputusan agent]`). caseId = SUBSTR(IDPEGA, 19) (= Java `substring(IDPEGA, 18)` GetPolicyData_ACT L1262).
  sistemBaru = baris T_WORK_POLIS ber-ID itu, LINI FAC. Polis tanpa baris JSON_POLIS -> satu versi "0" dari FACINPRODUCTION
  (MIN(IDPEGA), NOENDORS kosong `[keputusan agent]`). Tidak ada -> 404 "Polis tidak ditemukan.".
- **K71-B3 Ringkasan** (A71-3):
  - IDPEGA versi = JSON_POLIS (NOPOLIS, PRODKE; ganda -> TGL_INPUT terbaru); polis tanpa JSON_POLIS + prodke "0" ->
    FACINPRODUCTION MIN(IDPEGA);
  - baris = FACINPRODUCTION ber-IDPEGA + NOPOLIS itu, urut OBJECTNO, IDX_OBJECT_ITEM, COVERAGENOTE, JN_REAS `[keputusan agent]`;
  - kepala dari baris pertama; reffNumber = RISLIPCEDING terakhir menurut SPREAD_DATE se-NOPOLIS (rumus BrowseDataViewPolis1);
  - policyType = kode (kodePolicyType, label lama ikut);
  - RATE (kolom VARCHAR2; TO_NUMBER implisit dapat menulis ".15") dinormalkan ke desimal bertitik bila terurai, selain itu apa
    adanya;
  - angka lewat TO_CHAR TM9 -> apd -> teks; tanggal DD-MM-YYYY; NULL -> "";
  - ceding = CEDING_FACINPRODUCTION ber-IDPEGA itu, urut CEDINGID, CEDINGNAME;
  - nol DATA_JSON; versi / baris tidak ada -> 404; nopolis / prodke kosong -> 400.
- **K71-B4 Detail case sistem baru**: GET kasus / submit (riwayat) / lampiran tidak punya penjaga workbasket -> akun NB FacIn tanpa
  workbasket posisi membaca tanpa 403. GET submit memberi `bolehSubmit` false + alasan, bukan 403. Lampiran: hanya login.
- **K71-B5 Menu**: menu modul DATAR (satu baris `M_NAV_MENU` per modul; 962 sudah menyalakan `DIMIGRASI`). Halaman
  `nbfacin-view-policy` cukup di `HALAMAN_NBFACIN` frontend; tanpa migrasi, tanpa baris menu. Penjaga backend
  (`inti/backend/penjaga`) TIDAK membaca daftar halaman frontend.
- MODUL.md baris `JSON_POLIS` (tabel warisan) diperbarui: ditulis 69a-2, dibaca kolom relasional, DATA_JSON tidak dibaca.

**Uji**:
- services: `TestJenisCariPolis` (10 kasus + kutip), `TestCariPolis` (akun tanpa peran), `TestVersiDanRingkasanPolis`;
- repository: `TestSQLPolis` (27 kolom dihitung dari SQL = 27 indeks Scan, bind, nol DATA_JSON);
- handlers: `TestPolisHandler` (VersiPolis / RingkasanPolis + blok `kepala` 14 / `baris` 13 = api.ts dua arah + urutan; 400 / 401 /
  404).

**Mutasi 11 / 11 mati.** ⚠️ `PolisOracle` (urutan baca: JSON_POLIS -> cadangan FACINPRODUCTION) hanya teruji lewat teks SQL.

## Menu (09-10-2026)

- **K71-4** "View Policy dibuatin di menu baru. dibawah menu home" (gambar layar sidebar): tombol sidebar **View Policy** tepat
  di bawah Home, di luar golongan. Dibangun generik di inti: `MenuModul.menuAtas` (aditif, pola `antreanBeranda`) →
  `frontend/daftar.ts` `ENTRI_MENU` (`atas: true`) → `susunMenu` (`MenuTersusun.atas`; tampil bila tombol modul pemiliknya
  tampil; halamannya tidak menyalakan tombol NB FacIn) → `Shell` (ikon `IkonDokumen`) + palet Ctrl+K (sesudah Home).
  Pengecualian sadar atas "satu modul satu menu" (30-09-2026). Tanpa baris `M_NAV_MENU` / migrasi. Tombol sementara di portal
  NB FacIn dicabut.
- Uji inti diperbarui: `frontend/daftar.sinkron.test.ts` (palet = Home + tombol atas + tombol modul), `frontend/Shell.test.ts`
  (tombol atas bukan tombol modul). Kegagalan lama di luar tiket ini tidak berubah (premiumlistlife, treatycontractout,
  nbtreatyin, masterproductnamelife, Shell "menu profil").

## Tab Treaty / Life + List Offer Facultative (work owner 11-10-2026)

- Permintaan: "menu view policy ini buatin untuk life dan treaty juga, dipega ada kaya gitu … user friendly dan tidak halu".
  Sumber: `DDL\ViewOpenPolicyMarineInbox.xml` (tab Facultative / Treaty / Life) + `DDL\GetListViewPolicy_act.xml`
  (SEARCHON OFFERFAC / TREATYPOL / POLISLIFE / OFFERLIFE) + rule kiriman work owner: `BrowseDataViewOfferF`,
  `BrowsePolisTreaty`, `BrowsePolisLife`, `BrowseOfferLife`, grid `GridOpenoffer` / `GridOpenPolicyMarineTreaty` /
  `GridOpenPolicyLife` / `GridOpenOfferLife`, dan `FilterForViewPolicy`.
- Jawaban work owner atas View Option: Offer Facultative = sama dengan polis Facultative (tanpa QQ Name, karena FACINOFFER
  tanpa kolom QQ); Treaty = `FilterForViewPolicy`; Life = tanpa View Option.
- **List Offer Treaty TIDAK dibangun** (work owner: "list offer treaty tidak usah"). Grid `GridOpenOfferTreaty` juga tidak
  dikirim, dan activity tidak punya SEARCHON-nya.
- View Details Treaty / Life belum diport (rule detail belum dikirim). View Details Offer Facultative = layar case baca-saja.
- Frontend: `pages/ViewPolicy.tsx` (tab), `components/DaftarCariPolis.tsx` (daftar generik). Backend: kontrak c3, empat
  endpoint `/api/nbfacin/polis/{offer,treaty,life,life-offer}`.

## Backend sesi c3 — tab Offer / Treaty / Life (11-10-2026)

Berkas: `repository/viewpolis.go`, `services/viewpolis.go`, `handlers/viewpolis.go`, `models/viewpolis.go`; layanan
`DenganViewPolis(NewViewPolisOracle)`. Semua jawaban `{baris, total, halaman, ukuran: 20}`; filter kosong (atau hanya tanda
kutip / spasi) -> `baris: []`, `total: 0`, tanpa kueri; hanya login (K71-3); halaman bukan angka / < 1 -> 400; tanpa basis
data -> 503.

| Endpoint | SEARCHON / SQL | `baris[]` |
| --- | --- | --- |
| `GET /api/nbfacin/polis/offer?opsi=&filter=&halaman=` | OFFERFAC / `BrowseDataViewOfferF` (FACINOFFER, DISTINCT) | `caseId, noOffer, insuredName, sourceOfBusiness, qqName (""), startPeriod, endPeriod` |
| `GET /api/nbfacin/polis/treaty?opsi=&filter=&halaman=` | TREATYPOL / `BrowsePolisTreaty` (TREATYINPRODUCTION, DISTINCT) | `caseId, policyNo, insuredName, sourceOfBusiness, startPeriod, endPeriod` |
| `GET /api/nbfacin/polis/life?filter=&halaman=` | POLISLIFE / `BrowsePolisLife` (LIFEINPRODUCTION) | `caseId, policyNo, endorsNo, typePolicy, sourceOfBusiness, cedingCoName, qq, riSlipCeding` |
| `GET /api/nbfacin/polis/life-offer?filter=&halaman=` | OFFERLIFE / `BrowseOfferLife` (JSON_OFFER_LIFE versi terakhir) | `caseId, noOffer, typePolicy, businessName, cedingCoName, qq, lastStatus` |

- **K71-B6 Klausa** `[terverifikasi]` `DDL\GetListViewPolicy_act.xml` blok 9-12 + `pyBrowseSQL` keempat rule: teks CARI5 Pega
  dipertahankan, `{ASIS}` / literal `'…'` diganti BIND (nol nilai filter di teks SQL). Kurung WHERE mengikuti rule: Offer /
  Treaty `WHERE (…)`, Offer Life `AND (…)`, Life `WHERE …` TANPA kurung. Kutip dibuang (langkah 6). IDPEGA dua bentuk
  (`IN (:1, :2)` polos / berawalan; caseId = ekspresi dua bentuk).
- **Offer** (View Option, kontrak 9d: "samain kaya polis facultative", tanpa QQ): 1 Case ID hanya bila memuat `NB-` / `RNW-`
  (prakondisi 9.2; selain itu daftar kosong), 2 `NOOFFER = x`, 4 `RISLIPCEDING LIKE %x%`, 5 `INSUREDNAME LIKE x%`, 6
  `SOB LIKE x%`. Opsi 3 / 7 / lainnya -> 400 `[keputusan agent]` (No Slip ditunda K70-3; FACINOFFER tanpa QQ). ⚠️ Pega blok 9
  aslinya MENDETEKSI ISI (RNM -> NOOFFER; NB- / RNW- -> IDPEGA polos; selain itu `SOB LIKE %x% OR RISLIPCEDING LIKE %x%`) -
  View Option di sini keputusan layar baru. caseId Offer TIDAK dipotong 10 karakter (rule memilih IDPEGA utuh; tiga tab lain
  `substr(IDPEGA, 19, 10)`).
- **Treaty** 10.1-10.7 berurutan (`FilterForViewPolicy` 1-5): 2 `NOPOLIS = x`, ditimpa 10.2 `NOENDORS = x` (subkueri) bila
  panjang TEPAT 28 karakter; 1 + `NB-` IDPEGA, 1 + `EDMT-` subkueri IDPEGA (opsi 1 tanpa keduanya -> kosong; `RNW-` tidak
  dikenal 10.3); 3 NOOFFER; 4 INSUREDNAME `%x%`; 5 SOB `%x%`; semua `and PC0.NOENDORS is null`. Opsi lain 400.
- **Life** 11.1-11.6 berurutan (tanpa View Option): RNM & <= 22 -> NOPOLIS =; panjang 5 -> `NOPOLIS like %x`; `NBLF-` ->
  IDPEGA; > 22 & RNM -> subkueri NOENDORS; `EDML-` -> subkueri IDPEGA; tanpa RNM / NBLF- / EDML- DAN panjang bukan 5 ->
  `(SOBNAME LIKE OR CEDINGCONAME LIKE OR POLICYHOLDERNAME LIKE) AND NOENDORS IS NULL` `%x%` (tiga bind, BERKURUNG). Prakondisi
  11.2 `@length(filter)==5 T=(kosong) F=3`: benar -> lanjut (jalankan), salah -> lewati - syaratnya hanya panjang, bukan isi
  ("ABC12" pun 11.2). 11.4 / 11.5 `T=` kosong juga dibaca lanjut.
- **Offer Life** 12.1-12.3: `NBLF-` -> IDPEGA; `OFR-` (menimpa) -> `A.ID = SUBSTR(:1, 5)`; selainnya CEDING / POLICYHOLDER
  `%x%` (dua bind). Subkueri `C` Pega tanpa skema - di sini berkualifikasi seperti `A` / `b`.
- **Urutan halaman** `[keputusan agent]` (rule tanpa ORDER BY): Offer NOOFFER, CASEID; Treaty NOPOLIS, CASEID; Life NOPOLIS,
  NOENDORS NULLS FIRST, CASEID; Offer Life TGL_INPUT turun, CASEID. Tanggal `TO_CHAR(…, 'DD-MM-YYYY')`.
- **Tabel POOLDATA milik modul lain, BACA SAJA, tanpa impor modul**: TREATYINPRODUCTION (dinyatakan nbtreatyin / edmtreatyin),
  LIFEINPRODUCTION (endorsementlife), JSON_OFFER_LIFE (tak dinyatakan modul mana pun; ditulis Pega). Tidak ditambahkan ke
  tabel warisan MODUL.md (satu tabel warisan = satu modul).
- **DDL** `[belum terverifikasi lebar / tipe]` (tidak ada di `D:\migrasi\RNM\DDL`; semua dibaca sebagai teks kecuali tanggal):
  - LIFEINPRODUCTION: IDPEGA, NOPOLIS, NOENDORS, CEDINGCONAME, POLICYHOLDERNAME, TYPE, SOBNAME, RISLIPRNM - nama kolom ada di
    daftar INSERT `SaveLifeinProduction_SQL` (dikutip `endorsementlife/backend/repository/edm_putusan.go`);
  - JSON_OFFER_LIFE: ID, IDPEGA, TGL_INPUT, OLDID, CEDINGCONAME, POLICYHOLDERNAME, TYPECEDING, BUSINESSNAME, STATUS - nama
    kolom hanya dari `BrowseOfferLife` (OLDID / STATUS juga `GetOfferLife_sql`, premiumlistlife issue 00); TGL_INPUT diurutkan
    apa adanya (tipe belum terverifikasi).
  - FACINOFFER / TREATYINPRODUCTION `[terverifikasi]` `DDL\FACINOFFER.txt` / `DDL\TREATYINPRODUCTION.txt`: BEGINDATE / ENDDATE
    DATE; IDPEGA / NOPOLIS / NOENDORS / NOOFFER VARCHAR2(100); SOB / INSUREDNAME VARCHAR2(1000); RISLIPCEDING VARCHAR2(4000).

**Penyimpangan dari Pega (dicatat, bukan diam-diam)**:
1. Filter kosong Life / Offer Life: langkah 7 Pega HANYA keluar untuk POLIFAC / TREATYPOL / OFFERFAC; POLISLIFE / OFFERLIFE
   lanjut dengan `LIKE '%%'` (= semua baris). Di sini kosong -> daftar kosong (kontrak 9d).
2. `TRANS 1==1 -> SKP` (10.1, 10.3, 11.1-11.5; 7 lompatan): label SKP TIDAK ADA di activity (nol `pyStepsLabel`). `[dugaan]`
   (keputusan sesi 9d) langkah berikutnya tetap dievaluasi dan menimpa. **K71-5 `[keputusan work owner 11-10-2026]`**
   penyimpangan sadar dari bacaan itu: filter TEPAT 5 karakter yang tidak memuat RNM / NBLF- / EDML- -> 11.2 (digit akhir
   NOPOLIS) MENANG, bukan 11.6 (`00012`, `ABC12` -> `NOPOLIS like '%x'`). Filter 5 karakter yang memuat NBLF- / EDML- tetap
   ditimpa 11.3 / 11.5 (urutan langkah).
3. **K71-6 `[keputusan work owner 11-10-2026]`** 11.6 DIBERI KURUNG: `(SOBNAME LIKE … OR CEDINGCONAME LIKE … OR
   POLICYHOLDERNAME LIKE …) AND NOENDORS IS NULL` -> hanya polis induk. (Pega `A OR B OR C AND NOENDORS IS NULL` tanpa kurung
   mengikat `NOENDORS IS NULL` ke POLICYHOLDERNAME saja.)
4. Treaty DISTINCT tanpa `pzInsKey` (IDPEGA utuh): IDPEGA bentuk polos / berawalan untuk case yang sama menjadi satu baris.

**Uji**:
- repository `TestSQLViewPolis` (21 jenis: teks klausa persis, nol kutip, nilai bind, nomor bind jumlah / halaman),
  `TestSQLDasarViewPolis` (SQL dasar keempat rule persis, halaman + urutan);
- services `TestJenisCariOfferFac`, `TestJenisCariTreaty` (28 karakter rune, RNW- opsi 1), `TestJenisCariLife` (K71-5: `00012` /
  `ABC12` / 5 rune non-ASCII -> 11.2, 4 / 6 karakter -> 11.6, RNM 5 karakter -> 11.2, `NBLF-` / `EDML-` -> 11.3 / 11.5),
  `TestJenisCariOfferLife`, `TestCariViewPolis` (kutip, offset 20, kosong tanpa kueri, 400 / 401 / 503);
- handlers `TestViewPolisHandler` (JSON keempat tab persis - kunci + urutan untuk frontend 9d; kosong; 400 opsi / halaman;
  401; 503).

**Mutasi 16 / 16 mati; K71-5 / K71-6 3 / 3 mati.** ⚠️ `ViewPolisOracle` (Scan sel teks) hanya teruji lewat teks SQL - belum dijalankan ke Oracle DEV.

## Tombol View tab Treaty (work owner 11-10-2026: "tombol view treaty nya mana")

- `GET /api/nbfacin/polis/treaty/berkas?nopolis=` -> `{modul: "nbtreatyin" | "edmtreatyin", kasus}`; frontend membuka berkas lewat
  `PropsRute.onLihatBerkas` (kontrak 9d).
- Pola Claim Non Prop / Claim Prop View polis (OQ-CNP-13, `modul/claimnonprop/backend/repository/acuan.go` `sqlBerkasPolis` +
  `models.ModulBerkasPolis`) DISALIN, tanpa impor modul lain:
  `SELECT g.ID, g.PRODKE FROM T_GENERAL_POLIS_TREATY g JOIN T_WORK_POLIS w ON w.ID = g.ID WHERE g.NOPOLIS = :1 ORDER BY g.PRODKE DESC
  FETCH FIRST 1 ROWS ONLY`. PRODKE `NUMBER(10) DEFAULT 0 NOT NULL` `[terverifikasi]` migrasi nbtreatyin 320 -> dibaca bilangan;
  0 -> nbtreatyin, >= 1 -> edmtreatyin.
- Tanpa baris -> 404 "Polis ini belum punya berkas NB / EDM Treaty In di sistem baru." (polis Pega lama); nomor polis kosong
  (sesudah pangkas) -> 400; hanya login; tanpa basis data -> 503. T_GENERAL_POLIS_TREATY milik nbtreatyin, baca saja.
- Berkas: `repository/viewpolis.go` (`TabelGeneralPolisTreaty`, `sqlBerkasPolisTreaty`, `ViewPolisOracle.BerkasPolisTreaty`),
  `models/viewpolis.go` (`BerkasPolisTreaty`, `ModulBerkasPolisTreaty`), `services/viewpolis.go` (`ErrBerkasTreatyTidakAda`),
  `handlers/viewpolis.go`, pemetaan 404 di `handlers/handlers.go`.
- Uji: `TestSQLBerkasPolisTreaty` (teks SQL persis, 1 bind), `TestBerkasPolisTreaty` (PRODKE 0 / 1 / 7 / -1, pangkas, 404, 400
  tanpa membaca, 401, 503), `TestBerkasPolisTreatyHandler` (JSON persis, 404, 400, 401, 503). Mutasi 6 / 6 mati.
- **View Details Life ditunda** (work owner 11-10-2026: "yauda tidak jadi, biarin aja dlu"). Case `NBLF-…` dipegang modul
  `premiumlistlife`, yang rutenya belum menerima `bukaKasus` / `onBeranda` dari modul lain. Bila dilanjutkan: tambahan aditif di
  `modul/premiumlistlife/frontend/rute.tsx` (pola NB Treaty In), lalu `onLihat` di List Polis Life memanggil `onBukaModul`.
