# 66: Bilah CSV objek General (Download Template / Upload CSV / View Upload / Save Data)

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 08-10-2026 (gambar layar bilah
> "Please upload file with .csv format") "button insert accumulation cuma muncul di it, dan 4 button sisanya perbaiki, biar
> bisa di klik, 4 itu harusnya tidak muncul di uw". Riset korpus subagen: scratchpad sesi 0f `riset-D-csv-objek.md`.

**Status:** frontend selesai 08-10-2026 (uji hijau); backend selesai sesi c3 08-10-2026 (lihat "Backend sesi c3").

## Keputusan work owner 08-10-2026

- **K66-1** Insert Accumulation hanya akun IT (Pega: When `IsIT`); tombol tetap nonaktif (fungsinya di luar permintaan).
- **K66-2** Empat tombol CSV tidak tampil di UW / tahap baca-saja (Pega: `Periode_IsUW` tanpa tombol ini `[terverifikasi]`).
- **K66-3** Template = berkas asli `D:\migrasi\RNM\DDL\templatecsvfire.csv` (55 kolom, pemisah `;`) — "Saya taruh file-nya".
- **K66-4** Save Data MENGGANTI semua objek (seperti Pega) + konfirmasi bila case sudah punya objek — "Ganti semua + konfirmasi".

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- Tombol di `Section\Periode.xml` (`ASM-FW-GISFW-DATA-OFFERFACIN!PERIODE`) sel 93-97; bilah tampil `(IsFire || …) &&
  .ConfirmBinding != 1 && IsB2B != 'ASM'`.
- Upload: `UploadCSV_DataPolis` / `UploadCSVDataPolis_Act` → buffer `TempWorkPage.ListB2BHostUpload` (belum menyimpan).
- View Upload: harness `ViewCSVResult_B2BHost` (43 kolom format B2B lama, hanya 21 sama dengan parser).
- Save Data: `ChooseActivityMappingCSV_Act` → `InsertUploadFire_act` (+ `MappingJsonH2HNusare2_act`, `MappingCoverage`,
  `AddPaymentCurency_ACT`); pesan "Please Upload CSV" (b.492), "Success Insert Data" (b.7577), "KodePos … tidak ditemukan!".
  Master risk address / accumulation lewat prosedur (`InsertUpdateRISKADDRESS`, `RDBMASTERACCUMULATION`) — ADR-0043: port Go
  yang sudah ada (tiket 37 / 46).

## Keputusan agent

- **K66-5** View Upload menampilkan 55 kolom template (bukan 43 kolom grid Pega yang tidak sinkron dengan parsernya).
- **K66-6** Download Template berisi baris judul saja; baris contoh berkas Pega berisi nama tertanggung / marketing / ceding,
  tidak disalin.
- **K66-7** Pembacaan CSV di layar (judul wajib sama dengan template, tak peka huruf); normalisasi angka dan pemetaan di backend.

## Frontend (sesi 0f)

- `components/csvObjek.ts` (`KOLOM_CSV_OBJEK`, `TEMPLATE_CSV_OBJEK`, `bacaCsvObjek`), `components/BilahCsvObjek.tsx`,
  `api.ts` `simpanCsvObjek` (`POST …/objek/csv`), `labels.ts` `TEKS_CSV_OBJEK`, `pages/InwardFacultative.tsx` (bilah + muat
  ulang tab detail sesudah Save Data), `nbfacin.css` `.nbf-csvobj__*`. Uji `BilahCsvObjek.test.tsx`, `InwardFacultative.test.tsx`.

## Backend (diminta ke sesi c3, 08-10-2026)

Lihat pesan kontrak tiket 66.

## Keputusan backend 08-10-2026 (dikumpulkan sesi 9d; work owner lewat AskUserQuestion kecuali ditandai)

- **R1** pesan = daftar pesan KodePos, ATAU ["Success Insert Data"] bila tanpa pesan (`[terverifikasi]` InsertUploadFire_act
  langkah 6 b.7548: When `TempError.ERRMSG == ""`). R2-R4 diterima (precond=false = selalu jalan; template tanpa RNM_SHARE_PCT /
  COMM_PCT / COUNTRY / NOTES / COVERAGE_CONDITIONS; TSI_OBJECT_ITEM_6 tidak dinormalisasi Pega).
- **D1** "Timpa yang terisi saja": NO_OFFER_SLIP, BEGIN/END_DATE, OfferingDate, TypeFacultative; % Share, RI Comm, IsProRate tidak
  disentuh; beda antar baris → baris terakhir menang + log; tanggal dd/MM/yyyy `[dugaan]` → teks Pega WIB seperti General, tak
  terurai → 400.
- **D2** "Cocokkan benar": Object Type tak peka huruf ke 9 nilai baku → ejaan baku; selain itu "Others".
- **D3** `[keputusan agent]` master (risk address + akumulasi, cari-atau-tambah) lebih dulu; objek + header + currency SATU
  transaksi. ⚠️ Kontrak "gagal → tidak ada yang tertulis" TIDAK berlaku untuk master baru (sama dengan Pega: prosedur men-COMMIT
  per panggilan).
- **D4** `[keputusan agent]` ProvinceID dari nama provinsi RW; CZone dari CZoneZip; tidak ditemukan → akumulasi tidak dibuat,
  ACCUMULATION_CODE tetap nilai CSV, tanpa pesan baru (log). GetAccumulation_SQL (kelas Int-ACCUMULATION) tidak ada di korpus →
  `belum terverifikasi`.
- **D5** "Ikut Pega + premi dihitung": coverage lewat siapkanCoverage; spreading dibuang; T_CURRENCYLIST + T_LISTINSTALLMENT dihapus
  tanpa dibangun ulang.
- **D6 + D7** "Ketat": mata uang tak ditemukan → 400 (baris, CURRENCY); ItemTypeID tak ditemukan → kosong; angka wajib terurai
  (termasuk TSI_OBJECT_ITEM_6); tanggal wajib terurai; tanpa baris → 400 "Please Upload CSV"; maks 500 baris; kunci di luar 55 →
  400; semua galat dalam satu 400.
- **D8** `[keputusan agent]` literal ID coverage di-port apa adanya; keaktifan di master COVERAGE DEV = pemeriksaan DBA; ID tak
  ada di master saat runtime → 400 jelas.
- **D9** sesuai usulan c3 (okupasi = OLDID, risk address persis, KodePos BrowseRW2_SQL, TOP_RISK "Y", ADJUSTABLE_PCT ≠ "" dan ≠
  "100", penjaga = PUT …/objek, tanpa migrasi).

## Backend sesi c3 (08-10-2026)

`POST /api/nbfacin/kasus/{caseId}/objek/csv` — badan `{baris: {kolom: teks}[]}` (55 kolom `services.KolomCsvObjek` = judul
`DDL\templatecsvfire.csv` = `KOLOM_CSV_OBJEK` frontend; diperiksa tiga arah 55 = 55 = 55, selisih kosong), jawab 200 `{baris:
ObjekFire[], pesan: string[]}` (`HasilCsvObjekServer`). Tanpa migrasi, tanpa CALL.

- **K66-B1 Urutan** (`services/csvobjek.go` `SimpanCsvObjek`): 401 → penjaga PUT …/objek (403 / 409) → 503 → 404 → 400 isian
  CSV (kunci / jumlah / angka / tanggal, terkumpul) → header General (periksaGeneral, 400) → rencana per baris (pencarian master
  baca-saja; mata uang dan ID coverage master, 400 terkumpul) → periksaObjek (400) → master baru (D3) → SATU transaksi:
  SimpanGeneral, hapus CurrencyList + cicilan, objek diganti (spreading TIDAK dipertahankan) → baca ulang.
- **K66-B2 Per baris** (InsertUploadFire_act 5.2-5.14): normalisasi angka ("." dibuang, "," → "."), huruf besar TYPE_ADDRESS /
  ADDRESS / OBJECT_TYPE / KELURAHAN, ACCUMULATION_NAME = ADDRESS; satu lokasi (ObjectNo = nomor baris, RiskLocation = TYPE_ADDRESS
  + " " + ADDRESS, ObjectName = OBJECT_TYPE, ObjectType D2); RW BrowseRW2_SQL (zip + kelurahan tanpa spasi) bila keduanya terisi —
  tanpa hasil (termasuk zip / kelurahan kosong) → pesan verbatim "KodePos … and kelurahan … pada lokasi ke - n tidak ditemukan!"
  (tanpa "\n" Pega, satu butir per pesan); okupasi satu per lokasi (nama kosong → OCCUPATION TYPE FIRE ber-OLDID, urut Name);
  ListCauseOfLoss Detail "-".
- **K66-B3 Item / coverage** (MappingJsonH2HNusare2_act / MappingCoverage): slot ITEM_TYPE, _1.._6 yang terisi; Currency =
  CURRENCY master (GetDataCurrencyByName_SQL); ItemTypeID = OBJECTITEMTYPE (TYPE FIRE, nama tanpa spasi, ISACTIVE '1') — ⚠️ domain
  ID ini (OBJECTITEMTYPE.ID, sesuai Pega CSV) berbeda dengan pemilih item layar (V_JN_OBJ_ITEM.MJOI_KODE, tiket 39)
  `[pertanyaan terbuka]`; Note = INTEREST; coverage per peril RATE_x > 0 (FLEXAS 100815, 4.1A CC 100828, 4.3 100825, 4.2 PRGBI
  100829, OTHERS 100840), CoverageBasis "1", IndemnityPercentage 100, LoL = LOL, akumulasi = kode master / CSV.
- **K66-B4 Master** (D3 / D4): risk address `ADDRESS + POSTALCODE` persis → ID; tidak ada dan RW ditemukan → port tiket 37
  (`SisipRiskAddress`, transaksi sendiri). Akumulasi `NOTE = ACCUMULATION_NAME` persis → ID; tidak ada dan RW ditemukan → port tiket
  46 (`TambahAkumulasi`: Accumulation "10000001", tipe "FLEXAS", ScopeArea "AREA", Note = ADDRESS, Keyword menurut Object Type baku).
- **K66-B5 Header** (D1 b): dibangun dari General tersimpan lewat kabel + `periksaGeneral` (konversi tanggal / lebar sama dengan
  Save for later); Ceding Co dipertahankan; TypeFacultative ditulis label layar "Facultative In" (butir 78.3: data FacultativeIn =
  layar Facultative In) `[keputusan agent]`; Prorate coverage memakai periode BARU.
- **K66-B6 Keputusan agent kecil**: kode okupasi kosong tidak dicari (RD Pega membuang filter kosong lalu memungut okupasi
  pertama — tidak ditiru); TOP_RISK dibandingkan persis "Y" (ikut Pega); perbandingan ObjectType Pega `[dugaan]` peka huruf.
- Repository `repository/csvobjek.go` (`PembacaCsvObjek`): BrowseRW2_SQL, RISKADDRESS / ACCUMULATION / NATION / PROVINCE persis,
  GetDataCurrencyByName_SQL, GetDataObjectItem, BrowseOccupationFacInFIRE_RD (filter D TYPE + F OldID; C dibuang karena kosong),
  COVERAGE IN (…), hapus T_LISTINSTALLMENT + T_CURRENCYLIST case. `services/objek.go` `gantiObjek` (opsi periode / buang spreading /
  tulisan sebelum) dipakai PUT …/objek dan CSV.

**Uji** — `TestSimpanCsvObjek`, `TestCsvObjekPesanDanMaster`, `TestCsvObjekItemDanHeader`, `TestCsvObjekDitolak`,
`TestJenisObjekBaku`, `TestKolomCsvObjekSamaDenganFrontend` (services); `TestSQLCsvObjek` (repository); `TestCsvObjekHandler`
(kunci dua arah dengan api.ts `HasilCsvObjekServer`, 401 / 400 / 503). gofmt / vet bersih; `go test ./modul/nbfacin/...` hijau;
penjaga hanya kegagalan lama premiumlistlife; index git 0. **Mutasi 29 / 30 mati**; hidup: tanpa cek provinsi sebelum
TambahAkumulasi — ekuivalen (TambahAkumulasi sendiri menolak provinceId kosong → log, kode CSV tetap).

## Terbuka

- `[pertanyaan terbuka — DBA]` kelima ID coverage (100815 / 100828 / 100825 / 100829 / 100840) aktif di master COVERAGE DEV.
- `[pertanyaan terbuka — Product]` domain ItemTypeID (OBJECTITEMTYPE Pega CSV vs V_JN_OBJ_ITEM layar).
- `[pertanyaan terbuka — pemilik export Pega]` `ASM-FW-GISFW-INT-ACCUMULATION!…!GETACCUMULATION_SQL` (penamaan ulang nama akumulasi)
  tidak ada di korpus.

## Lanjutan 08-10-2026 — template 61 kolom (sesi c3)

Work owner menambah 6 kolom di belakang `DDL	emplatecsvfire.csv` (bukan kolom Pega): KODE_POS, KECAMATAN, KABUPATEN_KOTA,
PROVINSI, SUMBER_KODE_POS, TINGKAT_KEYAKINAN. Keputusan work owner (AskUserQuestion, lewat 9d): "Bagian template, dipakai";
kode pos "POSTAL_CODE dulu"; wilayah "Bantu cari RW"; SUMBER_KODE_POS / TINGKAT_KEYAKINAN "Hanya ditampilkan".

- **K66-B7** `services.KolomCsvObjek` = 61 (1-55 parser Pega tetap); diperiksa tiga arah: judul berkas = Go = `KOLOM_CSV_OBJEK`
  frontend (61 = 61 = 61, urutan sama, selisih dua arah kosong). Kunci di luar 61 → 400. SUMBER_KODE_POS / TINGKAT_KEYAKINAN
  diterima, tidak dipakai.
- **K66-B8 Kode pos** (`cariRWCsv`): BrowseRW2_SQL dengan POSTAL_CODE dulu; tanpa hasil dan KODE_POS terisi (berbeda) → coba
  KODE_POS. Kode pos yang MENEMUKAN RW dipakai lokasi (ZipCode), pencarian / penambahan risk address dan akumulasi `[keputusan
  agent: ZipCode lokasi ikut]`. Keduanya gagal → pesan verbatim memakai POSTAL_CODE (seperti Pega). Kelurahan kosong → tidak
  dicari (Pega langkah 5.6 dilewati).
- **K66-B9 Wilayah** (`pilihRWCsv`, aturan `[keputusan agent]`): master RW memuat `DISTRICTNAME` / `CITYNAME` / `PROVINCENAME`
  (`[terverifikasi]` DDL\RW.txt; sudah dipilih BrowseRW2_SQL). > 1 baris RW → baris pertama yang cocok dengan SETIAP
  KECAMATAN / KABUPATEN_KOTA / PROVINSI yang terisi (tak peka huruf, dipangkas); tidak ada → baris pertama (perilaku Pega) + log;
  satu baris → dipakai (beda wilayah → log).

**Uji** — `TestCsvObjekKodePosDanWilayah` (POSTAL_CODE ketemu; gagal + KODE_POS ketemu termasuk risk address yang sudah ada;
keduanya gagal → pesan; KODE_POS = POSTAL_CODE satu pencarian; banyak RW dipilih menurut wilayah / semua kolom terisi wajib
cocok), `TestKolomCsvObjekSamaDenganFrontend` 61. Mutasi aturan baru 8 / 8 mati. go test nbfacin + config + cmd hijau; penjaga
hanya kegagalan lama premiumlistlife; index git 0.

## RALAT D7 08-10-2026 — batas 3000 baris + cache per permintaan (sesi c3)

> ⚠️ RALAT D7 "maks 500 baris": work owner (gambar layar galat "paling banyak 500 baris CSV" untuk berkas 1128 baris) — "perbaiki,
> jangan dibatasin 500. 3000 aja buat batas nya". Batas kini **3000** (`services.BatasBarisCsvObjek`, satu tempat; pesan tetap
> "paling banyak 3000 baris CSV").

- **K66-B10 Badan** `handlers.batasBadanCsvObjek` = **16 MiB**. Dasar: overhead kunci JSON 1146 byte / baris (Σ panjang 61 nama
  kolom + kutip / titik dua / koma); 3000 baris × rata-rata 60 byte / sel = 13,75 MiB — dihitung dua cara (rumus dan JSON
  rakitan, sama). 16 MiB menampung rata-rata ≤ 72 byte / sel; di atasnya 400.
- **K66-B11 Cache per permintaan** (`denganCacheCsv`, salinan Service satu permintaan): RW (kode pos + kelurahan), risk address
  (alamat + kode pos), akumulasi (nama), negara, provinsi, mata uang, ItemTypeID, okupasi, CZone. Galat tidak di-cache.
- **K66-B12 Perbaikan cacat** (ditemukan saat menelusur, ada sejak K66-B4): pencarian master dilakukan untuk semua baris sebelum
  penulisan, sehingga dua baris beralamat sama yang belum ada di master MENAMBAH risk address / akumulasi DUA kali (akumulasi kedua
  ditolak "sudah ada" → kode CSV). Kini sebelum menambah dicari ulang, dan master baru dicatat ke cache → baris berikutnya
  memakainya (seperti Pega: Obj-Browse tiap baris sesudah prosedur baris sebelumnya men-COMMIT).
- **Uji** — `TestCsvObjekCacheSamaHasil` (40 baris / 5 alamat: hasil, objek tersimpan, jumlah master baru SAMA dengan / tanpa
  cache; master baru 5 + 5; alamat sama kode pos beda → dua risk address), `TestCsvObjekBatas3000` (3000 lolos; 3001 ditolak di
  `TestCsvObjekDitolak`), `TestCsvObjekHandler` (± 8,5 MiB diterima, > 16 MiB → 400). Mutasi 8 / 8 mati.

**Jumlah query pencarian master** (DB tiruan menghitung panggilan; termasuk cari ulang sebelum menambah):

| Baris | Alamat berbeda | Dengan cache | Tanpa cache | Master baru (risk + akumulasi) |
| ---: | ---: | ---: | ---: | --- |
| 1128 | 1128 | 2263 | 11281 | 1128 + 1128 |
| 1128 | 200 | 407 | 9425 | 200 + 200 |
| 3000 | 3000 | 6007 | 30001 | 3000 + 3000 |
| 3000 | 200 | 407 | 24401 | 200 + 200 |

`[dugaan]` **Perkiraan waktu** (asumsi 1-3 ms per round-trip Oracle; belum diukur di DEV):
- Pencarian master (cache): 0,4-18 s.
- Master baru: tiap alamat baru ± 2 statement risk address + beberapa statement mesin master akumulasi, masing-masing transaksi
  sendiri → terburuk 3000 alamat ± 18k statement ≈ 20-60 s.
- Transaksi objek (tidak diubah): `ObjekOracle.GantiObjek` ± 20 round-trip per lokasi (5 NEXTVAL + 5 INSERT lokasi / property /
  risk / bangunan / sekitar, + item, coverage, okupasi + table of limit, kerugian) → 1128 baris ± 22k, 3000 baris ± 60k statement
  dalam SATU transaksi ≈ 1-3 menit.

⚠️ **Risiko transaksi panjang — DILAPORKAN, strategi transaksi TIDAK diubah:** baris T_WORK_POLIS case dan baris objek case
terkunci selama transaksi (hanya case itu; Save lain atas case yang sama menunggu); server Go tanpa WriteTimeout (hanya
ReadHeaderTimeout 10 s), tetapi proxy / peramban bisa memutus permintaan > 60 s sementara transaksi tetap jalan sampai selesai.
Opsi bila perlu (keputusan): (a) array binding / INSERT ber-NEXTVAL di repository objek (round-trip turun ± 10×, menyentuh jalur
PUT …/objek); (b) proses latar + status; (c) batas lebih rendah.

## Butir 4 08-10-2026 — penulisan objek berkelompok (sesi c3)

`[keputusan work owner]` opsi (a) "Percepat penyimpanan". Berlaku untuk PUT …/objek DAN Save Data CSV (repository objek yang sama).

- **K66-B13 Traversal satu, penulis dua** — fungsi sisip (`objek.go` / `item.go` / `coverage.go` / `deductible.go` / `okupasi.go` /
  `fea.go` / `kerugian.go`) kini metode `tulisObjek` yang menulis lewat `penulisBaris` (nomor = NEXTVAL, sisip = satu baris).
  Logika, urutan DFS, SEQ_NO, PARENT_ID, uid, dan argumen TIDAK berubah.
- **K66-B14 Jalur berkelompok** (`repository/tulisberkelompok.go`): (1) lintasan hitung - jumlah ID per sequence; (2) ID diambil
  berkelompok DARI sequence (bukan ditebak): sequence ber-jumlah sama dalam SATU kueri `SELECT s1.NEXTVAL, s2.NEXTVAL, … FROM
  DUAL CONNECT BY LEVEL <= :1`; (3) lintasan tulis - ID dipakai berurutan, baris disangga per teks SQL (urutan kemunculan
  pertama = induk sebelum anak); (4) INSERT array-bind go-ora (v2.8.19, `[terverifikasi]` sumber driver: argumen slice →
  arrayBindCount, elemen nil = NULL) per potongan 500 baris; baris terpengaruh wajib = jumlah baris. Kelompok satu baris memakai
  bind skalar biasa (sama dengan jalur lama). Pengaman: kumpulan ID kurang / lebih dari yang dipakai → galat sebelum eksekusi.
- **Kenapa ID diambil lebih dulu**: anak (property, item, coverage, deductible, okupasi → table of limit, kerugian → coins) butuh
  ID induk SEBELUM induk ditulis; INSERT … RETURNING ber-array tidak dipakai (perilaku driver tidak terverifikasi).
- **Kesetaraan** (`tulisberkelompok_test.go`, tanpa basis data): jalur lama (satu NEXTVAL + satu INSERT per baris) vs berkelompok
  atas objek kecil (semua jenis anak, dua varian SQL deductible) dan CSV 1128 baris sintetis → baris per teks SQL IDENTIK
  (urutan, ID sequence, PARENT_ID, SEQ_NO, uid, nilai kolom), urutan teks SQL sama, induk dieksekusi sebelum anak. Semua uji PUT
  …/objek tetap hijau.
- **Round-trip penyisipan** (sentuh / hapus sama di kedua jalur; pengambilan ID dihitung termasuk fetch prefetch 25 baris
  bawaan go-ora): 1128 baris 24 816 → 79 (± 314×); 3000 baris 66 000 → 186 (± 355×). `[dugaan]` waktu transaksi objek 3000
  baris turun dari ± 1-3 menit ke beberapa detik (belum diukur di DEV).
- **Mutasi** 8 / 9 mati; hidup: penyangga berbagi slice argumen - ekuivalen (setiap pemanggil membuat slice baru per baris).
- ⚠️ **Belum diuji ke Oracle** (tidak ada basis data di sesi ini): array-bind go-ora dan `CONNECT BY` NEXTVAL perlu dicoba di DEV
  (Save Object satu lokasi + CSV kecil + CSV besar) sebelum dipakai luas.
- WriteTimeout server / proxy TIDAK diubah: permintaan CSV besar bisa ± 1 menit (frontend sudah memasang batas waktu 10 menit
  untuk simpanCsvObjek).


## Bug DEV 10-10-2026 — berkas CSV 1 juta baris tampil "Berkas CSV kosong." (sesi 0f)

- **Laporan work owner:** "INI KENAPA KOSONG? PADAHAL CSV NYA ADA 1JT BARIS".
- **Sebab:** `BilahCsvObjek` membaca berkas utuh lewat `File.text()`. Untuk berkas ratusan MB, peramban gagal membuat teks sebesar
  itu (teks kosong), sehingga pesan yang tampil "kosong". Pesan yang benar adalah batas 3000 baris (RALAT D7, keputusan work owner
  08-10-2026).
- **Perbaikan:**
  - `csvObjek.ts` `hitungBarisBerisi` menghitung baris berisi secara mengalir per potongan byte (`Blob.stream()`) dan berhenti
    begitu melewati batas. Baris kosong atau yang hanya berisi pemisah tidak dihitung, sama dengan pembaca CSV.
  - `BilahCsvObjek` memeriksanya SEBELUM `File.text()`, dengan pesan "Berkas CSV berisi lebih dari 3000 baris; paling banyak 3000
    baris per unggahan.".
  - Berkas berukuran > 0 yang teksnya kosong mendapat pesan sendiri ("tidak dapat dibaca").
- **Batas 3000 tidak diubah.** Unggahan berukuran jutaan baris butuh rancangan lain (unggah berkas ke server + proses latar);
  itu keputusan work owner.
