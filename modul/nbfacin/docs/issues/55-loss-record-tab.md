# 55: Tab utama Loss Record (kasus FIRE)

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 06-10-2026 "lanjut ke loss
> record" (tab ini dua kali dilewati karena di Pega hanya tampil untuk `IsIT`). Survei korpus oleh subagen (pengurai XML);
> butir kunci diperiksa ulang sesi 0f: `When\IsIT.xml`, prasyarat langkah `Activity\SaveLossRecord.xml` (btnsave, 1 / 3 / 5
> tahun, kurs), daftar properti loss ratio.
>
> **Berbeda dari tiket 42** (sub-tab Loss Record per LOKASI di tab Object, `.Property.ListCauseOfLoss`, `T_LISTCAUSEOFLOSS`).
> Tab ini memakai daftar tingkat CASE `.OfferFacIn.ListCauseOfLoss` (kelas `ASM-FW-GISFW-Data-CauseOfLoss`).

**Status:** frontend selesai 06-10-2026 (uji hijau); backend + migrasi tabel baru diserahkan ke sesi c3.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `Section\InputInwardFacultativeDtl.xml` kontainer "Loss Record" (S88): `pyContainerVisibleWhen=IsIT`, aktif
  `IsTABActive_LossRecord`; isi `ObjectList_LossRecord` (bila `IsFire`) + `SummaryLossRecord_Section`; tombol "Save"
  (`SaveLossRecord`, `btnsave="save"`).
- `When\IsIT.xml` (`DATA-PORTAL!ISIT`): `OperatorID.pyPosition = "IT Developer"` ATAU `OperatorID.pyOrgDivision = "IT"`.
- `Section\ObjectList_LossRecord.xml` (`ASM-FW-GISFW-WORK!OBJECTLIST_LOSSRECORD`, memo "add premi par"):
  - bar "Please upload file with .csv format": Download Template CSV (`DownloadMasterLossRecord` → berkas
    `TemplateCSVLossRecord`, isi berkas TIDAK ada di korpus), Upload CSV (`UploadCSV_DataPolis`), Save Data
    (`SaveLossRecord`, btnsave kosong);
  - Premi PAR / Premi MB / Premi PL / Premi Other COB (pxNumber 4 desimal, dapat diisi); Premi PAR default
    `GetPremiPAR_ACT` (Σ Premium × `.KursIDR` dari `TempTotal`, hanya bila 0 / kosong — `.KursIDR` tidak pernah diisi di
    korpus `[dugaan]` hasilnya kosong);
  - grid `.OfferFacIn.ListCauseOfLoss`: Uw Year · Class Of Business · Date of Loss (dd/MM/yyyy) · Cause Of Loss · Currency ·
    Loss Amount · Loss Amount ( IDR ) (baca-saja, 4 desimal), tambah / hapus baris, ubah baris lewat flow action
    `InputLossRecord` (`Section\InputLossRecord_Sec.xml`).
- `Activity\SaveLossRecord.xml`: CSV (CURRENCY, DATE_OF_LOSS dd/MM/yyyy, CLASS_OF_BUSINESS, CAUSE_OF_LOSS, LOSS_AMOUNT format
  Indonesia) → daftar diganti utuh; kurs `RDBList\GetKursLimitByName_SQL.xml` (`pooldata.treatyexchangeyearly`, TOIDR, per
  Begin date polis); Uw Year / tahun-ke dari `RDBList\Get5YearAgo.xml` (5 jendela 12 bulan mundur dari periode polis);
  Total Loss per Class Of Business + Currency; loss ratio:
  - Gross = Premi PAR + MB + PL + Other; Deduction = RI Commission / 100;
  - Loss Ratio N tahun = Σ Loss IDR / (Gross × N) × 100; Net = Σ Loss IDR / (Gross × (1 − Deduction) × N) × 100;
  - disimpan `.OfferFacIn.LossRatio1/3/5YearPercent`, `.NetLossRatio1/3/5Year`; tampil di `SummaryLossRecord_Section`.
- Sampel work object di `DDL\P-5 *.txt` dan `DDL\CONTOH\*.xml`: `ListCauseOfLoss` kosong, `LossRatio1YearPercent=0` —
  `[dugaan]` fitur belum dipakai produksi (selaras dengan tab yang hanya untuk IT).

## Cacat Pega `[terverifikasi]`

1. Ember 3 tahun `.TahunKe>1 || .TahunKe<=3` dan 5 tahun `.TahunKe>=1 || .TahunKe<=5` — selalu benar: SEMUA kerugian masuk.
2. Dua tombol simpan berbeda: "Save Data" (btnsave kosong) wajib CSV dan menghitung rasio; "Save" tab (btnsave "save")
   melewati hitung rasio.
3. `.DOL` diisi teks mentah CSV; `Local.DOL` (untuk Uw Year) dibangun terpisah.
4. Kurs: `{Search.CARI9} BETWEEN startdate AND enddate` atas kolom VARCHAR2 (`DDL\TREATYEXCHANGEYEARLY.txt`).
5. Premi PAR default memakai `.KursIDR` yang tidak pernah diisi.

## Batas pengetahuan

- Tidak ada tabel rancangan untuk daftar kerugian tingkat case, Premi PAR/MB/PL/Other, dan loss ratio 3/5 tahun / net
  (`services/loader/skema_gen.go`: hanya jalur `LocationList/Property/ListCauseOfLoss`) → butuh tabel baru (migrasi).
- Isi berkas template CSV, tabel enumerasi `matauang`, `BrowseBusiness_RD` / `BrowseCurrency_RD` belum dibaca / tidak ada.

## Keputusan work owner 06-10-2026 ("1. untuk sementara untuk POSITION=IT saja. 2. dibetulkan 3. sudah ditambah dengan nama templatecsvlossrecord 4. setuju 5. gbr nya itu")

1. **Akses** — tab hanya untuk akun login ber-`M_LOGIN_GO.JOB_POSITION` = "IT" (backend: `KasusNB.aksesLossRecord`, 403 di
   endpoint loss-record). Tanpa nilai → tab disembunyikan.
2. **Cacat Pega dibetulkan** — ember 3 / 5 tahun = tahun ke-1..3 / 1..5; satu Save yang selalu menghitung ulang; Premi PAR
   bawaan = Σ premi spread RNM × kurs Begin date; kurs lewat pembaca TREATYEXCHANGEYEARLY Cedant (bukan BETWEEN VARCHAR);
   Loss Amount ( IDR ) selalu terisi; Total Loss dihitung (Pega: halaman clipboard, "No items" sesudah muat ulang).
3. **Template CSV** = `D:\migrasi\RNM\DDL\templatecsvlossrecord.csv` (CRLF, pemisah `;`, kolom CLASS_OF_BUSINESS ·
   DATE_OF_LOSS · CAUSE_OF_LOSS · CURRENCY · LOSS_AMOUNT + catatan di kolom ke-10) → `labels.ts`
   `TEMPLATE_CSV_LOSS_RECORD`, diuji sama persis dengan berkas itu.
4. **Tabel baru** (migrasi ditulis c3, dijalankan work owner).
5. Gambar layar Pega diterima — tata letak mengikutinya.

## Frontend (sesi 0f)

- `components/TabLossRecord.tsx` (+ `TabLossRecord.test.tsx`), `components/csvLossRecord.ts` (baca CSV di layar: kolom
  menurut nama kepala, `;` atau `,`, sel berkutip, angka format Indonesia, tanggal DD/MM/YYYY, Cause of Loss huruf besar);
  `api.ts` kontrak; `labels.ts` `LOSS_RECORD` / `TEKS_LOSS_RECORD` / `TEMPLATE_CSV_LOSS_RECORD`; `nbfacin.css` `.nbf-lr*`;
  `pages/InwardFacultative.tsx` (tab disaring `aksesLossRecord`). Diperiksa visual lewat Chrome headless (data sintetis).

## Keputusan agent

- **L55-1** Uw Year dan Loss Amount ( IDR ) baca-saja di ubah baris (Pega dapat diketik tetapi dihitung ulang saat Date of
  Loss / Loss Amount berubah).
- **L55-2** Currency dan Cause of Loss dapat diubah (Pega: baca-saja di `InputLossRecord_Sec`, sehingga baris tambahan manual
  tidak dapat diberi mata uang).
- **L55-3** CSV bergalat (kolom hilang, tanggal / jumlah salah, tanpa baris) tidak mengganti grid sama sekali; CSV sah
  mengganti grid utuh (seperti `SaveLossRecord` langkah 4) dan baru tersimpan lewat Save Data.
- **L55-5** (work owner 06-10-2026 "tombol save loss recordnya blm ada") Tombol **Save** tab di kaki ditambahkan
  (`InputInwardFacultativeDtl.xml` kontainer Loss Record, `<pyLabel>Save</pyLabel>` → `SaveLossRecord`); aksinya sama
  dengan Save Data (keputusan 2: satu simpan yang selalu menghitung ulang).
- **L55-6** (work owner 06-10-2026 "isi loss record nya di rapiin dong biar lebih user friendly") Panel ubah baris: info
  baca-saja di atas (Insured Name, Uw Year, Loss Amount ( IDR )), isian 3 kolom berlabel atas rata kiri (Date of Loss ·
  Currency · Loss Amount / Cause of Loss · Class Of Business), tinggi seragam, kartu beraksen kiri. Penyebab label rata kanan
  sebelumnya: aturan kolom tombol `td:last-child` ikut mengenai sel colSpan panel - kini dibatasi ke baris data.
- **L55-7** (work owner 06-10-2026 "Class Of Business nya kalau yang diketik tidak ada dari list, langsung dihps aja
  ketikannya") `IsianSaran` opsi `wajibDariDaftar`: saat kotak ditinggal sesudah diketik, teks dicari ulang ke
  `…/class-of-business/saran`; label yang sama persis (tanpa beda huruf besar/kecil) dipilih, selain itu isian dikosongkan.
  Tambahan (work owner 06-10-2026 "iya boleh itu ditambah juga"): Upload CSV juga wajib dari daftar -
  `csvLossRecord.ts` `cocokkanKelasBisnis` (tiap nilai berbeda dicari sekali; ejaan diganti ke NOTE daftar; kosong /
  tidak ada -> galat per nomor baris berkas, grid tidak diganti). Backend PUT memvalidasi ulang (diminta ke c3); `hitung`
  tidak divalidasi (dipanggil sambil mengetik).
- **L55-4** Loss ratio tampil 4 desimal + "%" (Pega 2 desimal) — seragam dengan angka lain modul ini.

## Backend (sesi c3, 06-10-2026)

**Status:** backend selesai; migrasi 206 DITULIS, belum dijalankan. Kontrak = `frontend/api.ts` `TampilanLossRecord` /
`KasusNB.aksesLossRecord`. Rule dibaca ulang sendiri (pengurai XML): `Activity\SaveLossRecord.xml`, `RDBList\Get5YearAgo.xml`,
`RDBList\GetKursLimitByName_SQL.xml`, `Activity\GetPremiPAR_ACT.xml`, `ReportDefinition\BrowseBusiness_RD.xml` +
`Section\InputLossRecord_Sec.xml`.

- **K55-1 Akses** — `aksesLossRecord` (GET kasus dan jawaban PUT General) = `M_LOGIN_GO.JOB_POSITION` akun pelaku
  (`inti.PelakuDari` → AkunID = LOGIN_ID), dipangkas, tanpa beda huruf = "IT". LUNAK: tanpa sesi / akun tak ada / kolom
  kosong / galat baca → false (galat dicatat log) — GET kasus tidak pernah gagal karena tab ini. Ketiga endpoint tab
  DITEGAKKAN di layanan (ADR-0030 butir 1): tanpa identitas **401** (`inti.ErrTanpaIdentitas`, pola inti — kontrak menyebut
  403), bukan "IT" **403** (`inti.ErrTanpaWewenang`, kini dipetakan handler nbfacin). ⚠️ Ketegangan ADR-0030 butir 3
  ("modul memakai aturan peran yang sudah ada; peran baru ditambahkan ke aturan bersama"): sumber peran bersama hari ini =
  workbasket login (`Pelaku.Peran`) dan hak menu `M_LOGIN_GO_MENU.HAK`; JOB_POSITION adalah sumber kedua. Sesuai keputusan
  work owner ("untuk sementara") — `[pertanyaan terbuka]` bila kelak dipindah ke aturan peran bersama (mis. hak menu).
- **K55-2 Non-FIRE** — TIDAK ditolak: backend belum punya predikat FIRE tingkat case (`IsFire` belum diport); tab FIRE
  lain (Object, Coverage, Spreading, Scoring) juga tidak menolaknya. Layar hanya menampilkan tab untuk kasus FIRE.
- **K55-3 Tabel** — migrasi `206_loss_record_kasus.sql`: `T_LOSSRECORD` (satu per kasus, UNIQUE, induk T_GENERAL_POLIS) +
  `T_LOSSRECORDLIST` (baris). Bukan amandemen loader (jalur case tidak ada di workbook; data sampel kosong) — jalur
  `.OfferFacIn.ListCauseOfLoss` dan Premi PAR / loss ratio tingkat case tetap ke penampung loader `[pertanyaan terbuka]`
  bila data lama kelak memuatnya.
- **K55-4 Hitung** (GET, hitung, dan PUT — selalu menghitung ulang; satu Save, cacat 2 dibetulkan):
  - DOL kosong → Begin date (WIB) − 1 tahun (Calendar.add, 29-Feb → 28-Feb); CauseOfLoss huruf besar (Pega); Currency
    dipangkas + huruf besar (kontrak; Pega TIDAK mengubah currency).
  - Jendela L = 1..5: [ADD_MONTHS(Begin, −12L), ADD_MONTHS(End, −12L)] inklusif, tanggal WIB, aturan ADD_MONTHS Oracle (hari
    terakhir bulan → hari terakhir bulan); Uw Year = tahun ADD_MONTHS(Begin, −12L). DOL di beberapa jendela (mis. tepat
    Begin − 1 tahun dengan End = Begin + 1 tahun) → jendela TERAKHIR (L terbesar) seperti loop Pega. Di luar jendela → Uw
    Year "" dan tidak masuk rasio. Begin / End kosong → Uw Year kosong + pesan.
  - Loss Amount IDR: currency IDR atau kosong → = Loss Amount (Pega melewati konversi untuk keduanya); lain → × kurs
    TREATYEXCHANGEYEARLY pada Begin date (pembaca Cedant / Scoring: CURRENCY.ID → KursPada; cacat 4 dibetulkan), 8 desimal;
    kurs tidak ada → kosong + pesan (sekali per mata uang).
  - Total per Class Of Business + Currency urut kemunculan; Loss Amount IDR kosong dilewati.
  - Rasio (cacat 1 dibetulkan): Y1 = TahunKe 1, Y3 = 1..3, Y5 = 1..5; Gross = PAR + MB + PL + Other (kosong = 0); Deduction
    = RI Commission % (T_GENERAL_POLIS.PAY_RI_COMMISION, sumber tab Deduction) / 100 EKSAK; skala 20 lalu 8 desimal. Gross 0
    → rasio kosong + pesan; RI Commission 100 → net kosong + pesan.
  - Premi PAR bawaan (cacat 5 dibetulkan): kosong / 0 → Σ PremiumSpreaded per mata uang (total spreading per mata uang,
    `[dugaan]` = TempTotal Pega) × kurs Begin date, dibulatkan 0 desimal setengah ke atas `[dugaan]`; kurs tidak ada → pesan,
    Premi PAR tetap.
  - hitung: urutan baris dipertahankan. PUT: urut DOL terbaru dulu (stabil, tanpa DOL di akhir), ganti utuh, premi + 6 rasio
    + per baris Uw Year / TahunKe / Loss Amount IDR / kurs disimpan.
- **K55-5 Masukan** — badan ketat (kunci asing → 400); angka teks desimal ≤ 8 desimal; dateOfLoss DD-MM-YYYY atau kosong;
  ≤ 500 baris; classOfBusiness ≤ 4000, causeOfLoss ≤ 500, currency ≤ 10 bita.
- **K55-6 Class Of Business** — `GET /api/nbfacin/class-of-business/saran?q=`: BrowseBusiness_RD dengan ketiga saringan
  (ID, Note, Group) dikirim KOSONG oleh section → semua grup; NOTE mengandung q (tidak peka huruf), NOTE kosong dibuang, urut
  ID DESC (RD), batas 50 (RD 500); id = BUSINESS.ID, label = NOTE.
- **K55-7 Urutan deploy** — 206 cukup sebelum tab Loss Record dipakai: GET kasus TIDAK membaca tabel 206 (hanya
  M_LOGIN_GO, sudah ada sejak inti 902); tanpa 206 hanya endpoint loss-record yang gagal ORA-00942. Urutan aman: 202 → 203 →
  204 → 205 → 206 → biner (205 tetap WAJIB sebelum biner karena GET kasus membaca T_SCORINGRISK).
- **K55-8** (c3, perbaikan sesudah telaah dua sumbu)
  - Hasil yang disimpan wajib muat NUMBER(38,8): premi, 6 rasio, Loss Amount IDR melebihi 30 digit bulat → 400 ber-jalur
    (pola coverage), bukan ORA-01438 (mis. Gross premi sangat kecil).
  - RI Commission ≥ 100 % → net loss ratio kosong + pesan (sebelumnya hanya = 100; > 100 memberi rasio negatif).
  - PUT: Total per COB dihitung atas baris TERURUT, sama dengan GET berikutnya.
  - Akses dicek handler SEBELUM badan diurai: bukan IT dengan badan rusak tetap 403 (bukan 400).
  - Pembulatan bersama (`bulatSetengahAtas`, juga Scoring) tidak lagi menelan galat: di luar presisi → nilai asli.
- **K55-9** (c3, catatan telaah spec — perilaku = Pega, dicatat)
  - Premi PAR 0 yang diketik pengguna diganti bawaan bila ada premi spread (Pega memeriksa nilai halaman saat itu, sama).
  - Premi PAR bawaan: satu mata uang tanpa kurs → seluruh bawaan dilewati (+ pesan); Pega membulatkan jumlah berjalan 4
    desimal tiap langkah, di sini sekali 0 desimal di akhir.
  - DOL bawaan memakai Calendar −1 tahun (Begin 28-02-2025 → 28-02-2024), jendela memakai ADD_MONTHS (L1 mulai 29-02-2024)
    → DOL bawaan itu jatuh di luar L1 (seperti Pega). Pega memformat Begin tanpa zona waktu untuk DOL bawaan; di sini WIB
    `[dugaan]`.
  - Penolakan non-FIRE dan 401-vs-403 tanpa sesi: lihat K55-1 / K55-2 — menunggu penerimaan work owner.
- **K55-10** (c3, keputusan work owner 06-10-2026 "Class Of Business WAJIB dari daftar BUSINESS.NOTE") — PUT
  `…/loss-record`: tiap baris `classOfBusiness` wajib terisi DAN sama persis (dipangkas, tanpa beda huruf besar-kecil) dengan
  satu `BUSINESS.NOTE` (`repository/bisnis.go` `NoteKelasBisnis`: satu kueri `UPPER(TRIM(NOTE)) IN (…)`, beberapa baris
  berkunci sama → ID terkecil); disimpan dengan ejaan NOTE tabel. Gagal → 400 `ErrMasukanLossRecord`, pesan per baris
  bernomor mulai 1 ("baris 2: Class Of Business wajib diisi", "baris 3: Class Of Business \"X\" tidak ada di daftar
  BUSINESS"), tanpa simpan. POST `…/hitung` TIDAK memeriksanya (layar menghitung sambil mengetik). Data lama tidak
  dimigrasi — GET menampilkan COB tersimpan apa adanya; Save berikutnya baru menuntutnya dari daftar.
- **K55-11** (c3, permintaan sesi 9d) `GET /api/nbfacin/class-of-business/saran?q=`: dengan q, NOTE yang SAMA PERSIS (UPPER(TRIM(NOTE))
  = UPPER(TRIM(q))) diurutkan lebih dulu, lalu ID DESC, batas 50 tetap — layar memutuskan "ada di daftar" dari jawaban saran,
  jadi NOTE yang persis sama tidak boleh terpotong batas bila > 50 NOTE mengandung q. Tanpa q: tetap ID DESC.

## Ditegaskan work owner 07-10-2026

"Tab Loss Record dan Correspondence itu cuma untuk IT aja , selain itu gausa munculin" - tetap: tab hanya tampil dan hanya
dapat diubah oleh akun ber-JOB_POSITION "IT"; akun lain tidak melihat tabnya. Saat case keluar dari tahap marketing, tab ini
ikut baca-saja (syarat TAHAP saja, tanpa workbasket marketing - keputusan c3 K57-20 dipertahankan).
