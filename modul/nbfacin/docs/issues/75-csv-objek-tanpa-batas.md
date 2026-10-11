# 75: Unggah CSV objek TANPA BATAS baris (1 juta baris harus bisa)

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 10-10-2026 "INI KENAPA KOSONG? PADAHAL
> CSV NYA ADA 1JT BARIS", lalu "UNGGAHAN CSV TIDAK USAH ADA BATASANNYA, JADI 1 JUTA JUGA HARUS BISA". Ini mengganti batas 3000 baris
> (tiket 66 RALAT D7).

**Status:** ready-for-agent. Kontrak dikirim ke sesi c3 10-10-2026; frontend dikerjakan sesi 0f.

## Kenapa rancangan lama tidak bisa `[terverifikasi]`

Alur tiket 66:

1. Layar membaca berkas utuh di peramban (`File.text()`), memecah semua baris, lalu mengirim satu JSON
   (`POST …/objek/csv`, `handlers.batasBadanCsvObjek` = 16 MiB).
2. Server memvalidasi dan menyimpan seluruhnya dalam **satu transaksi**. Untuk 3000 baris itu sudah ± 1–3 menit dan case
   terkunci (tiket 66 "Risiko transaksi panjang").

Pada 1 juta baris:

- peramban gagal membuat teks sebesar itu (bug "Berkas CSV kosong.", tiket 66);
- JSON berukuran GB;
- satu transaksi berisi jutaan statement berjalan berjam-jam;
- tab Object memuat semua objek sekaligus.

## Keputusan work owner 10-10-2026

- **K75-1 Rancangan:** berkas dikirim UTUH ke server (streaming), lalu dibaca, divalidasi, dan disimpan server di LATAR per
  kelompok. Layar menampilkan status / kemajuan. Tabel penampung unggahan + status pekerjaan = **migrasi 257** (rentang nbfacin
  255–259), ditulis agent; work owner / DBA yang menjalankan.
- **K75-2 Bila gagal di tengah:** tandai gagal ("gagal di baris N") + tombol **Ulangi** yang melanjutkan dari kelompok yang gagal.
  Selama proses, Save lain di case itu ditolak.
- **K75-3 Tab Object per halaman dari server** untuk case besar. Tab Coverage / Spreading yang memuat semua objek menyusul dengan
  pola yang sama.

## Rancangan (usulan agent, mengikat kontrak)

- **Tanpa batas baris.** Tidak ada batas badan 16 MiB untuk berkas. Server membaca aliran CSV baris demi baris; memori tidak
  bergantung pada jumlah baris.
- **Aturan validasi per baris SAMA dengan tiket 66** (judul 61 / 55 kolom, isian wajib, master, alamat, dst.), hanya dijalankan
  mengalir. Master di-cache per pekerjaan.
- **Semantik "Save Data" tetap MENGGANTI semua objek case** (tiket 66):
  - kelompok pertama menghapus objek lama, kelompok berikut menambah;
  - hapus `T_CURRENCYLIST` / `T_LISTINSTALLMENT` + hitung ringkasan dilakukan SEKALI di akhir.
- **Kemajuan tahan restart:** status, baris diproses, dan kelompok terakhir yang sukses disimpan di tabel. Pekerjaan yang terputus
  karena server mati dapat dilanjutkan dengan Ulangi.
- **Kunci case:** selama status memvalidasi / menyimpan, `PUT …/objek` dan Save Data lain untuk case itu mendapat 409 berpesan.
- **Mode lama dipertahankan untuk case kecil** ("jangan ubah yang sudah berhasil"):
  - tab Object beralih ke mode per halaman hanya bila jumlah lokasi > ambang (usulan 500);
  - di bawah ambang, `GET`/`PUT …/objek` utuh seperti sekarang.

## Fase

| Fase | Isi |
| --- | --- |
| F1 | unggah berkas (stream) + validasi latar + status + View Upload per halaman dari server |
| F2 | Save Data latar per kelompok + Ulangi + kunci case |
| F3 | tab Object per halaman (case > ambang) + simpan perubahan saja |
| F4 | Coverage / Spreading / tab lain untuk case besar; CSV Aneka / Marine Cargo (tiket 73 P6 / 74 P6) memakai infrastruktur yang sama |

## Backend F1 + F2 — dibangun (sesi c3, 10-10-2026)

**Status:** dibangun, uji hijau. Migrasi 257 DITULIS (belum dijalankan; work owner / DBA).

- **Migrasi 257** (`migrations/257_csv_objek_unggah.sql` + `_down`):
  - `T_NB_CSV_UNGGAH` (pekerjaan: STATUS, FASE validasi / simpan, kemajuan, header langkah 5.3, GALAT_PROSES / GALAT_DI_BARIS,
    DIUBAH = detak);
  - `T_NB_CSV_BARIS` (penampung: SATU KOLOM PER KOLOM CSV, `C_<kolom>` VARCHAR2(4000) × 61 — dipilih supaya View Upload per
    halaman dan validasi membaca kolom bernama tanpa mengurai ulang);
  - `T_NB_CSV_PESAN` (galat / pesan per baris).
  - ⚠️ **Temuan c3:** `T_LOCATIONLIST.SEQ_NO` NUMBER(5) (186) hanya memuat 99 999 lokasi; 1 juta baris = 1 juta lokasi →
    257 memperlebar ke NUMBER(10). Jalur mundur tidak mempersempit kembali (ORA-01440).
  - Pembersihan penampung lama TIDAK dibuat (usulan menunggu keputusan work owner).
- **Endpoint:**
  - `POST /api/nbfacin/kasus/{caseId}/objek/csv/berkas` — multipart/form-data medan `berkas`, dibaca MENGALIR
    (`r.MultipartReader`, tanpa ParseMultipartForm / batas badan) → 202 status. Judul salah → 400 tanpa pekerjaan.
  - `GET …/objek/csv/aktif` (null bila belum pernah), `GET …/objek/csv/{idUnggah}`,
    `GET …/objek/csv/{idUnggah}/baris?halaman=&ukuran=` (ukuran 1..1000).
  - `POST …/objek/csv/{idUnggah}/simpan` (hanya dari `valid`) dan `…/ulangi` (hanya dari `gagal`) → 202.
  - Status: `{idUnggah, status, fase, totalBaris, barisDiproses, barisTersimpan, galat[≤100], jumlahGalat, pesan[≤100],
    jumlahPesan, galatProses, namaBerkas, dibuat, adaObjek}` — `fase`, `jumlahPesan`, `galatProses` tambahan c3.
  - 409: pekerjaan berjalan (PUT `…/objek`, Save Data lama, unggah baru) / langkah tidak sah untuk status; 404 pekerjaan case lain.
- **Dialek** = pembaca layar tiket 66 (pemisah ";", kutip, BOM, baris kosong dilewati, judul 61 / 55 kolom).
- **Aturan per baris SAMA dengan tiket 66.** `uraiBarisCsv` dan `rencanakanBarisCsv` dipecah dari jalur JSON; jalur lama tidak
  berubah dan ujinya tetap hijau. Galat ketat: satu galat → `gagal_validasi`.
- **Simpan per kelompok 1000 baris:**
  - kelompok 1 = General + hapus CurrencyList / cicilan + GANTI objek (spreading dibuang); kelompok 2.. = TAMBAH (SEQ_NO lokasi
    lanjut, `ObjekOracle.TambahObjek`);
  - kemajuan pekerjaan ditulis DI TRANSAKSI objek kelompoknya, jadi Ulangi melanjutkan tepat dari kelompok yang belum COMMIT
    (tanpa baris ganda);
  - master baru (risk address / akumulasi) tetap COMMIT sendiri (D3).
- **Tahan restart:** pekerjaan berjalan tanpa detak > 15 menit (dihitung di SQL, bebas zona waktu) dianggap terputus → status
  `gagal`, kunci lepas, Ulangi melanjutkan.
- **Uji:**
  - pembaca mengalir: 1 000 000 baris sintetis + validasi per baris, puncak heap 3,5 MiB, ± 35 detik; dilewati di `-short`;
  - alur 2 500 baris (3 kelompok), gagal validasi, ulangi sesudah putus di kelompok 3, kunci / basi;
  - kabel HTTP (multipart, rute, 400 / 503);
  - SQL (bind, kolom dua arah terhadap 257).

## Backend F3 — dibangun (sesi c3, 10-10-2026)

**Status:** dibangun, uji hijau (`go vet`, `go test -short ./...`). Tanpa migrasi baru. Belum dicoba di DEV.

- **`GET …/kasus/{caseId}`** memuat `jumlahLokasi` (`number | null`; null bila tak terbaca). Layar memakainya untuk memilih
  mode; usulan ambang 500 ada di `services.AmbangHalamanObjek`, dan keputusan tampilannya di layar.
  - `api.ts` `KasusNB` mendapat `jumlahLokasi?: number | null` (aditif, opsional), karena uji kontrak Go ↔ TS
    `TestKunciKasusSamaDenganApiTs` mewajibkannya.
- **`GET …/objek/halaman?halaman=&ukuran=&cari=`** → `{total, halaman, ukuran, lokasi: [{id, …objek}]}`.
  - `halaman` ≥ 1; `ukuran` 1..500 (batas ekspresi IN Oracle, ORA-01795). Di luar itu → 400.
  - `id` = `T_LOCATIONLIST.ID`. Medan objek SAMA dengan `GET …/objek`, termasuk olahan bacanya (tanggal kerugian, total per
    mata uang).
  - Urut `SEQ_NO, ID`.
  - Cari `[keputusan agent]`: memuat, tanpa beda huruf, di Object No / Object Name / alamat risiko / kode pos. Wildcard
    pengguna di-escape.
  - Pohon anak dibaca dengan kueri baca objek yang SAMA, dibatasi `AND l.ID IN (…)` halaman itu (`lingkupBaca`), jadi GET
    utuh tidak berubah.
- **`PATCH …/objek`** `{ubah: [{id, …objek}], tambah: [objek], hapus: [id]}` → `{jumlahLokasi}`; halaman dibaca ulang layar.
  - Urutan pemeriksaan sama dengan `PUT …/objek`: identitas 401, hak ubah 409, isian 400, basis data 503, kunci unggah 409
    (K75-2), mata uang 400. Penyiapan per baris juga sama (loss ratio, coverage).
  - Batas sekali simpan `[keputusan agent]`: `ubah` + `hapus` ≤ 500, `tambah` ≤ 500, badan ≤ 8 MiB. ID wajib angka dan tidak
    ganda.
  - Satu transaksi:
    1. sentuh case dan pastikan General;
    2. ID `ubah` / `hapus` wajib milik case — bila tidak, **409** "Lokasi objek sudah berubah sejak halaman dimuat";
    3. pohon `ubah` + `hapus` dihapus, dengan urutan tabel SAMA dengan `PUT`;
    4. `ubah` disisip ulang di `SEQ_NO`-nya, dengan ID baru;
    5. `tambah` disisip sesudah `SEQ_NO` terbesar.
- **`[keputusan agent]`, dilaporkan untuk dikonfirmasi:**
  - **K75-4 — `ubah` atas lokasi yang sudah ber-spreading → 409.** Pesannya: "…perubahan per halaman untuk lokasi
    ber-spreading belum tersedia (tiket 75 F4)". Alasannya: spreading per lokasi belum bisa dipertahankan tanpa membaca case
    utuh. `hapus` tetap boleh, dan spreading lokasi itu ikut terhapus, sama seperti `PUT`.
  - **K75-5 — `SEQ_NO` tidak dirapatkan sesudah `hapus`.** Urutan tetap terjaga, tetapi celah boleh ada; pembaca FIRE
    memakai `ORDER BY`. Merapatkan 1 juta baris tanpa MERGE berbiaya kuadrat, dan `MERGE … UPDATE SET` ditolak penjaga
    ADR-U-0033 (pola regex). `PUT` dan Save Data CSV tetap menulis 1..n.
- **Uji:**
  - SQL: lingkup dibatasi tepat satu pola per kueri baca; bind berurutan; hapus = tabel dan urutan `sqlHapusObjek`;
    SEQ_NO per baris;
  - service: isian, 503, 404 / 409, kunci unggah, olahan baca;
  - kabel HTTP: bentuk JSON, 400, 404, 409, `jumlahLokasi` / null.
  - Logika transaksi `PerbaruiObjek` belum teruji atas Oracle sungguhan (sesi ini tanpa akses DEV).
- **Ikut diperbaiki (regresi penjaga dari sesi yang sama):**
  - `inti/backend/outbox/penyimpan.go` `sqlPungutEfek` sekarang memakai `%s` dua kali, bukan `%[1]s`, sehingga penjaga
    nama-tabel-telanjang bisa membacanya;
  - STRUKTUR `T_NB_CSV_UNGGAH.DIBUAT` / `DIUBAH` bertipe "timestamp", sesuai kosakata penjaga golongan tipe.

## Frontend F1 + F2 (sesi 0f, 10-10-2026)

**Status:** dibangun, uji hijau. Belum bisa dicoba sebelum migrasi 257 dijalankan di DEV.

- **`BilahCsvObjek.tsx` ditulis ulang.** Berkas TIDAK dibaca utuh di peramban:
  1. Judul diperiksa dari 64 KB awal (`csvObjek.ts` `periksaJudulBerkas` / `judulCsvSah`, aturan judul tiket 66), sehingga salah
     template langsung ketahuan.
  2. Berkas di-stream ke `POST …/objek/csv/berkas` (`api.ts` `unggahCsvObjek`, multipart `berkas`, batas waktu 30 menit lewat
     parameter opsional baru `mintaFormulir(…, batasWaktuMs)` di `inti/frontend/klien.ts` — aditif, bawaan tetap).
  3. Status dibaca tiap 2 detik selama memvalidasi / menyimpan, dengan pesan kemajuan ("350.000 baris diperiksa", "120.000 /
     1.000.000 baris (12%)").
- **Galat per baris:** paling banyak 100 pesan, ditambah "… dan N galat lain".
- **View Upload:** dibaca per halaman dari server.
- **Save Data:** hanya dari status `valid`; konfirmasi ganti bila `adaObjek`.
- **Gagal:** tombol **Ulangi** (K75-2), plus `galatProses` bila ada.
- **Status dipulihkan** sesudah muat ulang (`…/objek/csv/aktif`). Tab detail dimuat ulang SEKALI saat `menyimpan` → `selesai`.
- **Batas 3000 baris tidak lagi dipakai jalur layar.** `hitungBarisBerisi` / `BATAS_BARIS_CSV_OBJEK` tetap ada untuk jalur JSON lama.
- **Uji:** `BilahCsvObjek.test.tsx` (judul dari potongan awal, tanpa `File.text()`, pesan status, Ulangi, Save hanya dari valid).

## Frontend F3 — tab Object per halaman (sesi 0f, 10-10-2026)

**Status:** dibangun, uji hijau.

- **Pemilihan mode** (`pages/InwardFacultative.tsx`, tab Object FIRE): `jumlahLokasi` > `AMBANG_HALAMAN_OBJEK` (500, sama dengan
  `services.AmbangHalamanObjek`) memakai **`TabObjectHalaman`**. Selain itu `TabObject` tetap, tanpa perubahan alur; hanya
  `ObjectAddress` kini diekspor supaya dipakai ulang.
- **`components/TabObjectHalaman.tsx`:**
  - grid 10 baris per halaman dari `GET …/objek/halaman` + kotak cari server;
  - sub-tab detail dan aturan Save SAMA dengan `TabObject`;
  - perubahan dicatat per `id` dan terbawa saat pindah halaman;
  - lokasi baru di daftar terpisah di atas grid;
  - hapus ditandai dulu (baris redup bercoret, "Batal hapus").
- **Save** mengirim HANYA perubahan (`PATCH …/objek`, `api.ts` `perbaruiObjek`), lalu halaman dibaca ulang. Batas per Save mengikuti
  backend: ubah + hapus ≤ 500, tambah ≤ 500, dengan pesan di layar.
- **409 server** tampil apa adanya: lokasi berubah sejak dimuat, K75-4, unggahan berjalan.
- **Inti (aditif):** `OpsiMinta.metode` + `'PATCH'` (`inti/frontend/klien.ts`).
- **Uji:** `components/TabObjectHalaman.test.tsx`. Dua uji teks-sumber lama (FIRE → TabObject) dilonggarkan untuk cabang ukuran;
  yang diuji tetap FIRE → TabObject.

### Keputusan work owner 10-10-2026 (F3)

- **K75-4** Disetujui (sementara): di mode per halaman, mengubah lokasi yang sudah ber-spreading ditolak 409 berpesan. Hapus tetap
  boleh, dan spreading-nya ikut terhapus. Dibuka di F4 bersama Coverage / Spreading per halaman.
- **K75-5** Disetujui: SEQ_NO tidak dirapatkan sesudah hapus per halaman (boleh ada celah, urutan tetap). PUT tab Object dan Save
  Data CSV tetap 1..n.

## Bug DEV 11-10-2026 — CSV FIRE masuk ke case MARINE CARGO (sesi c3)

**Gejala:** unggahan CSV FIRE ±1 juta baris tersimpan di NB-151032 (MARINE CARGO); seharusnya di NB-151031.

**Akar** `[terverifikasi]` (dibaca sesi 0f): `NB FacIn\Activity\ChooseActivityMappingCSV_Act.xml` menjalankan
`InsertUploadFire_act` (L439) HANYA bila When IsFire (prakondisi L499). Marine Cargo, Aneka, dan Briguna masing-masing punya
activity sendiri. Port tiket 66 / 75 memproses template FIRE untuk SEMUA Class of Business.

**Perbaikan backend** (`services/csvfire.go` `wajibCsvFire`):
- Predikatnya sama dengan `lini.fire` GET kasus (`LiniKasus().Fire`); lini tak terbaca = bukan FIRE.
- Dipasang di Save Data JSON lama, `…/objek/csv/berkas`, `…/simpan`, `…/ulangi` (pekerjaan lama ikut ditolak), dan jalur latar
  `simpanUnggah` (pekerjaan ditandai gagal tanpa menulis objek).
- Respons 409: "Upload CSV template FIRE hanya untuk Class of Business FIRE (Pega ChooseActivityMappingCSV_Act). CSV untuk Class
  of Business ini belum tersedia di sistem baru."
- Uji: `TestCsvFireHanyaKasusFire` (Marine, Aneka, lini tak terbaca), `TestGalatCsvBukanFire`.

**Pemindahan data:** keputusan work owner 11-10-2026 "Pindahkan langsung". Skripnya `docs/alat/pindah_objek_151032_ke_151031.sql`
— DITULIS saja, belum dijalankan, untuk DBA. Isinya:
- cek sebelum dan sesudah;
- satu blok tanpa COMMIT dengan penjaga;
- daftar data yang tidak pulih, beserta opsi Flashback Query `[dugaan]` yang harus dibaca segera.
