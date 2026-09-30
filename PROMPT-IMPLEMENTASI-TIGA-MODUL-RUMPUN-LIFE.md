# PROMPT — TIGA MODUL SEKALIGUS, RUMPUN LIFE: **Master Contract Retro Life · Master Product Name Life · Endorsement Life** *(sesi baru, folder `OUTPUT_HASIL_RNM`, cabang **`dev`**)*

> ⛔ **DITAHAN 30-09-2026** *(keputusan work owner: "aku mau hanya ini, Treaty Contract Retro Life")*. Kerjakan dulu **hanya**
> Master Contract Retro Life lewat **`PROMPT-IMPLEMENTASI-MODUL-MASTER-CONTRACT-RETRO-LIFE.md`**. Brief ini dipakai lagi untuk Master Product
> Name Life dan Endorsement Life sesudahnya; bab 2 di bawah sudah tercakup brief tunggal itu.

> Permintaan work owner 30-09-2026: *"kerjakan 3 modul sekaligus, jangan hanya 1 tiket, jangan banyak yang di-skip, baca ulang XML, ambil logic,
> button, dan method yang benar dari XML, jangan membuat menu yang tidak ada di Pega. XML patokan dasar; tiket = hasil grilling; tiket yang salah
> diperbaiki di dokumennya."*
>
> Ketiganya adalah sisa **rumpun Life** *(`PROMPT-EKSEKUSI-HULU-HILIR.md` §1: Master Contract Retro → Master Product Name → PremiumList → Endorsement)*.
> PremiumList Life dan Claim Life **sudah dibangun**. Ketiga modul punya spec `ready-for-agent`, tiket, dan grilling; korpusnya paling kecil
> *(66, 114, dan 75 berkas XML)*.

## 0. ATURAN BACA XML — berlaku di setiap tiket

1. **Setiap pembacaan activity mencetak `pyStepsBlockName`**. `//` = langkah ter-remark, **tidak pernah jalan** — jangan dibangun *(pelajaran GILIRAN-11)*.
2. Pecah XML dengan `sed -e 's/></>\n</g'` dan kutip nomor baris `bNNN` di setiap bukti.
3. WHEN dibaca **dari urutan aksi**, bukan dari daftar `pyPropertiesName`/`pyPropertiesValue` yang diratakan. Precondition `WhenTrue 2` = lanjut, `3` = lewati.
4. Cari **penulis** sebuah kolom, bukan hanya pembacanya.
5. Setiap tombol di layar baru = satu tombol di section XML *(nama section + `bNNN`)*, memanggil activity yang sama. **Nol tombol, menu, layar, atau aksi
   tanpa bukti XML.**
6. Tiket yang berbeda dari XML atau katalog DEV **diperbaiki di dokumennya**: catatan bertanggal 30-09-2026, kalimat lama dikutip, bukti XML atau katalog
   disebut. Berkas `grilling-ronde-*` tersegel dan **tidak disunting**.

## 1. BATAS KERJA

| Hal | Isi |
| --- | --- |
| Cabang | **`dev`**. Commit dengan jalur eksplisit `git commit -o -- <jalur>`, awalan pesan = nama modul *(`mastercontractretrolife: …`)*. **Nol `git push`, nol `git pull --rebase`** |
| Folder boleh disunting | `APP_RNM/modul/mastercontractretrolife/**`, `APP_RNM/modul/masterproductnamelife/**`, `APP_RNM/modul/endorsementlife/**`, dan ketiga `inti/backend/daftar/modul_<nama>_gen.go` hasil `go generate` |
| Dilarang | `inti/`, modul lain *(termasuk `premiumlistlife` dan `claimlife`)*, konfigurasi root. Pola dari modul lain **boleh ditiru**, **tidak boleh diimpor**. Butuh fungsi bersama atau kontrak baru → catat di laporan |
| Nomor | Retro Life **100–139**, slot **958–959** · Product Name **140–179**, slot **960–961** · Endorsement **480–519**, slot **976–977** *(`MODUL.md` masing-masing)* |
| Oracle | skema `POOLDATA`. `-migrate` dijalankan **work owner**. Uji tag `db` tanpa `ORACLE_DSN` = SKIP; **POOLDATA bukan target uji `db`**. Nol penulisan ke DEV oleh sesi ini |
| Prosedur Oracle | **tidak dipanggil**; logikanya ditiru di Go *(keputusan o)*. Nol `COMMIT` di teks SQL |
| Layanan luar | `ConnectREST` di Product Name dan Endorsement: **stub** di DEV lewat pola outbox yang sudah ada di modul lain, alamat nyata tidak dipanggil |
| Rahasia | nol nama orang, nomor polis, kredensial, host/IP, alias tnsnames di berkas apa pun. Fixture berawalan `UJI-` |
| CSS | kelas khusus modul berawalan nama modul, di berkas CSS di folder frontend modul itu, diimpor dari `rute.tsx`-nya |

## 2. MODUL A — MASTER CONTRACT RETRO LIFE

**Seluruh isi `PROMPT-IMPLEMENTASI-MODUL-MASTER-CONTRACT-RETRO-LIFE.md` berlaku**, termasuk ralat **K1–K8**, paket 0–11, dan OQ-MCRL-01/02.
Ringkasnya: lima tabel warisan `*_LIFE`, **nol PK/FK di DEV**, jadi **nol DDL**; kaskade dan keunikan di Go; `ID` = `'1'||LPAD(seq.NEXTVAL,6,'0')`;
gerbang tahun ikut Pega *(mati)* sampai OQ-MCRL-01 dijawab.

## 3. MODUL B — MASTER PRODUCT NAME LIFE

### 3.1 Fakta DEV *(diukur asisten 30-09-2026, agregat, baca-saja)*

| Hal | Isi |
| --- | --- |
| Tabel | `M_PRODUCT_LIFE` 196 baris, `M_PRODUCTINWARD_LIFE` 197 baris — keduanya `JSONDATA` + satu constraint `IS JSON`, **nol PK** |
| Pembaca | view **`PRODUCT_LIFE`**, **`PRODUCTINWARD_LIFE`**, **`DOCUMENTCLAIM_LIFE`**; prosedur `PEGA_M_PRODUCT_LIFE`, `PEGA_M_PRODUCT_INWARD_LIFE`; kode **Claim Life** membaca `PRODUCTINWARD_LIFE` *(ambang produk, kategori dokumen, rate, validasi tanggal)* |
| Definisi view | **`modul/masterproductnamelife/docs/dba-view-produk-life.md`** *(asisten, dari `ALL_VIEWS`)* — daftar kunci JSON yang dibaca, peka huruf besar-kecil |
| Sequence | `M_PRODUCT_LIFE_SEQ`, `M_PRODUCT_INWARD_LIFE_SEQ`, `M_PRODUCT_TYPE_LIFE_SEQ` |
| Lampiran | `M_ATTACHMENTPRODUCTNAME` 2 baris, `T_STORAGE_IMAGE` 517, `T_FOLDER_IMAGE` 1 — tabel warisan |

### 3.2 Ralat

| # | Klaim di dokumen | Fakta | Menjadi |
| ---: | --- | --- | --- |
| P1 | tiket **01** dan spec penyimpangan 2: *"skema relasional penuh"*, satu `product_life` + lima tabel anak, dua tabel lama digabung, JSON dibuang | tiga view, dua prosedur, dan Claim Life membaca `JSONDATA` kedua tabel lama. Produk yang hanya ditulis ke tabel baru **tidak terlihat** oleh semuanya | **tiket 01 ditangguhkan → OQ-MPNL-01**. Bawaan: **ikut Pega** — simpan ke `M_PRODUCT_LIFE` dan `M_PRODUCTINWARD_LIFE` dengan bentuk `JSONDATA` persis yang ditulis activity Pega *(`SetProductName*`, `Save*`)*, termasuk empat kolom hasil flatten `RIRISKID`, `RIRISK`, `PRODUCTNAME`, `BEGIN_DATE` seperti prosedur `PEGA_M_PRODUCT_LIFE`. **Nol tabel baru, nol DDL**. Model Go tetap bertipe *(struct bernama)*; JSON hanya di lapisan repository |
| P2 | batasan AC 38 *"tidak mengurai `JSONDATA` m_product_life"* | batasan itu milik **Claim Life** sebagai pembaca | modul ini **penulis** tabelnya, jadi boleh membaca dan menulis `JSONDATA` miliknya sendiri. Claim Life tetap membaca lewat view |
| P3 | `OutwardList` bukan tabel, jalur OR mati | view `PRODUCT_LIFE` membaca `OutwardList[0].OVR_COMM` | kunci `OutwardList` **tetap ditulis** dengan bentuk yang sama seperti Pega menulisnya *(cari penulisnya di XML)*; jangan dibuang dari JSON |
| P4 | simpan atomik dua sisi *(tiket 03)* | tetap benar | kedua tabel ditulis dalam **satu transaksi**; gagal satu = gagal semua |
| P5 | lampiran dan efek keluar *(tiket 08, 09)* | tabel warisan ada; `ConnectREST` ke layanan luar | tulis ke tabel warisan seperti XML; kirim berkas lewat stub outbox; nol panggilan ke alamat nyata |

Uji wajib: satu produk `UJI-` yang disimpan layanan, dibaca kembali dengan **teks SQL ketiga view** di skema uji tiruan, menghasilkan nilai yang sama
dengan masukannya *(uji `db` SKIP tanpa DSN, plus uji murni atas bentuk JSON yang dihasilkan terhadap daftar kunci di `dba-view-produk-life.md`)*.

### 3.3 Tiket

| Tiket | Kerjakan |
| --- | --- |
| 01 | **ditangguhkan** *(P1)*; ralat bertanggal di tiket dan spec |
| 02 | produk sisi umum, identitas dari sequence, gagal terang-terangan |
| 03 | sisi inward, simpan atomik *(P4)* |
| 04 | tujuh pemilih master — sumber data dari RD/RDB XML masing-masing |
| 05 | validasi wajib isi seragam di kedua sisi — hanya aturan yang ada di XML atau keputusan tertulis di tiket |
| 06 | daftar plan dan batas underwriting |
| 07 | daftar komentar, dokumen, underwriting finansial |
| 08, 09 | lampiran opsional, efek keluar penyimpanan berkas *(P5)* |

Satu harness di korpus = halaman awal modul.

## 4. MODUL C — ENDORSEMENT LIFE

### 4.1 Fakta DEV dan kode

| Hal | Isi |
| --- | --- |
| Tabel aplikasi milik PremiumList | `T_PREMIUM_LIST`, `T_PREMIUM_LIST_DETAIL` *(migrasi 051, 052)* — ada di DEV, **0 baris** |
| Tabel warisan | `JSON_POLIS` 148 baris, `M_TEMPUPLOADLIFE` 1.350, **`M_LIFE_PREMIUM_DETAIL` ±66,8 juta baris** |
| Penulis `M_LIFE_PREMIUM_DETAIL` di kode | PremiumList *(`polis_warisan.go`, cermin `SaveMasterLPDet`)* |
| Pembaca peserta | Claim Life membaca `M_LIFE_PREMIUM_DETAIL` dan **sudah** menyaring `EDMSTATUS` `Batal`/`Delete` *(`pesertapolis.go`)* |
| Harness | `EditDetail_Harness`, `EndorsmentLife_harnes`, `InboxEndorsementLife`, `ViewCSVResult_LifeEDM`, `ViewOldPolicy_EDM`, `_QP`, `_TP`, `_TR` |

### 4.2 Ralat

| # | Klaim | Fakta | Menjadi |
| ---: | --- | --- | --- |
| E1 | tiket **00** `ALTER` dua tabel, blok oleh PremiumList tiket 00 | kedua tabel sudah ada; blok terpenuhi | tiket 00 **dikerjakan** di rentang **480–519**, di folder migrasi Endorsement. `PRODKE NUMBER DEFAULT 1` supaya kode PremiumList **tidak perlu diubah** dan baris new business otomatis versi 1. Semua kolom EDM lain nullable. `PARENT_ID` self-FK nullable ber-index |
| E2 | tiket **11** kontrak hilir ke Claim Life | Claim Life sudah membaca `M_LIFE_PREMIUM_DETAIL` dengan penyaring `EDMSTATUS` | Endorsement **menulis baris `M_LIFE_PREMIUM_DETAIL`** dengan `EDMSTATUS` `Old`/`New`/`Delete`/`Batal` persis `Endorsement Life/RDBList/SaveMasterLPDet.xml`, dan summary seperti `InsertPLSummary`. **Nol perubahan di kode Claim Life dan PremiumList**. Uji: baris `Delete` dan `Batal` tidak lolos penyaring Claim Life *(tiru penyaringnya di uji, jangan impor)* |
| E3 | tiket **12** migrasi skema/data lama, `needs-info` OQ-001 | memindah data warisan `JSON_POLIS` ke tabel baru belum diputuskan | **di luar lingkup sesi ini**; tetap `needs-info` |
| E4 | `M_LIFE_PREMIUM_DETAIL` | ±66,8 juta baris | setiap kueri memakai kunci ber-index *(`PL_NUMBER`, sebut indeksnya dari katalog uji atau XML)*; **nol pemindaian penuh**; uji tidak pernah menyentuhnya |
| E5 | efek keluar dan alarm *(tiket 10)*, `ConnectREST` | layanan luar | stub outbox; alarm = pesan di layar + log, tanpa alamat nyata |

### 4.3 Tiket

00 *(E1)*, 01 pilih polis dan gerbang kelayakan, 02 buat case dan muat polis lama, 03 popup lihat polis lama per Type *(empat harness `ViewOldPolicy_EDM*`)*,
04 penomoran `PL_NUMBER_EDM` *(`BAHAN-tiket-pl-number-edm.md`)*, 05 maksud endorsement dan `EDMSTATUS`, 06 jurnal balik delete dan batal, 07 unggah CSV
*(`ViewCSVResult_LifeEDM`)*, 08 alur keputusan, kunci field, audit, 09 simpan transaksi campuran dan anti-dobel, 10 *(E5)*, 11 *(E2)*. Tiket 12 *(E3)* tidak.

Halaman awal = `InboxEndorsementLife` bila XML menunjukkan harness itu pintu masuk portal *(sebut buktinya)*.

## 5. MENU — satu modul, satu menu

Menu sudah datar *(migrasi 901 di `dev`)*. Setiap modul mendapat **satu** tombol dari baris `M_NAV_MENU` yang sudah ada: **Master Contract Retro Life** dan
**Master Product Name Life** di `MASTER`, **Endorsement Life** di `TREATY`. Slot menu tiap modul hanya:

```sql
UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1', TGL_UBAH = SYSDATE WHERE KODE = '<modul>'
/
```

ditambah `_down` ke `'0'`. **Nol `INSERT`**. `menu.ts` menyatakan `HALAMAN_AWAL`. Layar lain dibuka dari dalam layar, sesuai aksi di XML.

## 6. PARITAS LEBIH DULU — per modul

`modul/<nama>/docs/PARITAS-LAYAR-DAN-AKSI.md`: setiap harness, section, tombol, dan aksi → activity *(`bNNN` + `pyStepsBlockName`)* → RDB/RD → tabel dan
kolom → rute API → komponen React → status. Cocokkan tiket dengan tabel paritas **sebelum** menulis kode. Aksi yang ada di XML tetapi tidak dicakup tiket
= **dibangun** bila perilakunya jelas dari XML *(jangan di-skip)*, dicatat sebagai tambahan tiket; bila tidak jelas = OQ.

## 7. URUTAN — satu commit per paket, uji hijau di setiap commit

Kerjakan **vertikal** *(backend + frontend + uji per paket)*. Tiga modul berjalan **bergantian per gelombang**, bukan satu modul sampai selesai lalu modul
berikutnya.

| Gelombang | Retro Life | Product Name | Endorsement |
| ---: | --- | --- | --- |
| 1 | paket 0: PARITAS + ralat K1–K8 + OQ | PARITAS + ralat P1–P5 + OQ | PARITAS + ralat E1–E5 + OQ |
| 2 | paket 1 kerangka + baca | kerangka + baca produk lewat JSON + tujuh pemilih *(04)* | tiket 00 migrasi + kerangka + baca polis |
| 3 | paket 2–3 tahun, kontrak, batas | 02, 03 simpan dua sisi atomik | 01, 02, 03 pilih polis, case, popup polis lama |
| 4 | paket 4–5 reinsurer, security | 05, 06, 07 validasi dan daftar | 04, 05 penomoran, maksud endorsement |
| 5 | paket 6–8 business, rate, terapkan ke semua, hapus berjenjang, galat | 08, 09 lampiran dan efek keluar | 06, 07, 09 jurnal balik, CSV, simpan campuran |
| 6 | paket 9–10 frontend lengkap + menu | frontend lengkap + menu | 08, 10, 11 + frontend lengkap + menu |
| 7 | dokumen: status tiket, PARITAS, OQ, `MODUL.md` Status → dimigrasi | sama | sama |

Setiap commit: `go build`, `go vet` dan `go test` *(dengan dan tanpa tag `db`)*, `gofmt`, `tsc`, `vitest`, `vite build`, di pohon kerja utama. Sesudah
gelombang 6: jalankan backend dan `npm run dev`, buka ketiga layar di peramban, pastikan tombol yang tampil = tombol di tabel PARITAS.

## 8. BERHENTI HANYA BILA

- Sebuah aturan hanya bisa dibangun dengan DDL pada **tabel warisan** — OQ, jangan tulis migrasi; lanjutkan tiket lain.
- Perlu menyunting `inti/`, `claimlife`, atau `premiumlistlife` — catat, lanjutkan tiket lain.
- Perilaku sebuah aksi XML tidak dapat dipastikan — OQ dengan bukti `bNNN`, lanjutkan tiket lain.

Satu tiket yang terhenti **tidak** menghentikan dua modul lainnya.

## 9. LAPORAN

Satu pesan:

1. Tabel **modul → paket → commit → tiket → bukti XML → angka uji**.
2. Ralat yang diterapkan *(K, P, E)* dan tiket yang diperbaiki dokumennya.
3. Tombol dan aksi XML yang ditemukan di luar tiket, dan statusnya.
4. OQ baru per modul.
5. Halaman awal tiap modul beserta buktinya.
6. Pengingat `-migrate` untuk work owner: Endorsement 480-an, slot 958, 960, 976.
7. Persentase kesiapan per modul dan total migrasi *(modul dimigrasi dari 20)*.
8. Bab **TELEMETRI EKSEKUSI**.

## 10. OQ UNTUK WORK OWNER

| OQ | Pertanyaan | Bawaan sampai dijawab |
| --- | --- | --- |
| OQ-MCRL-01 | Gerbang tahun Retro Life ternyata mati di Pega. Tetap ditegakkan, atau ikut Pega? | ikut Pega |
| OQ-MCRL-02 | PK/FK lima tabel Retro Life tidak ada di DEV walau dokumen DBA menyebutnya sudah dipasang. Dipasang di lingkungan lain? | nol DDL |
| OQ-MPNL-01 | Produk tetap disimpan sebagai JSON di dua tabel lama seperti Pega, atau dipindah ke tabel relasional baru? Tabel baru membuat tiga view, dua prosedur, dan Claim Life tidak melihat produk baru, kecuali kedua bentuk ditulis bersamaan | JSON seperti Pega |
| OQ-EDM-001 | Data endorsement lama di `JSON_POLIS` dipindah atau tidak *(tiket 12)* | tidak dikerjakan |

---

*Disusun 30 September 2026 dari katalog DEV (`ALL_OBJECTS`, `ALL_CONSTRAINTS`, `ALL_DEPENDENCIES`, `ALL_VIEWS`, `ALL_SEQUENCES`, cacah baris), daftar rule
tiga folder korpus, spec dan tiket ketiga modul, serta kode `claimlife` dan `premiumlistlife` yang membaca atau menulis tabel yang sama.*
