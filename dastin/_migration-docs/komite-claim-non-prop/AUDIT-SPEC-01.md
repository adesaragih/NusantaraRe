> Modul  : Komite Claim Non Prop · Audit spesifikasi · 2026-09-21
> Peran  : auditor
> Masukan: SPEC-KOMITE-01.md (2026-09-21) · KETETAPAN.md · CONTEXT.md · REGISTER-PAGAR.md · REGISTER-DEVIASI.md · INVENTARIS-BUKTI.md — kelimanya Ronde 06, 2026-09-20
> Status : DITUTUP 2026-09-21
> Sifat  : TAMBAH-SAJA

# AUDIT-SPEC-01

Audit atas `SPEC-KOMITE-01.md` sebelum `to-tickets`. **Nol baris spesifikasi diubah.**
Temuan masuk daftar perbaikan di bagian akhir, tidak diterapkan sendiri.

Gerbang lolos: keenam berkas masukan ada di disk. `AUDIT-SPEC-01.md` belum ada sebelum
audit ini.

---

## 0. Temuan pokok, disebut lebih dulu

`SPEC-KOMITE-01.md` ditulis mengikuti template naratif — Problem Statement, Solution, User
Stories, Implementation Decisions, Testing Decisions, Out of Scope, Further Notes. Skill
`to-spec` memerintahkan **enam belas bagian bernomor** yang berpusat pada **persyaratan
`S-xxx`**, masing-masing dengan `Aturan`, `Bila gagal`, `Tanda`, `Beda AS-IS`, dan `Uji`.

Diukur ke disk:

| Yang dicari | Cacah di `SPEC-KOMITE-01.md` |
|---|---:|
| Persyaratan `S-xxx` | **0** |
| Blok `Bila gagal` | **0** |
| Baris `Tanda :` | **0** |
| Baris `Beda AS-IS` | **0** |
| Bagian `Ketertelusuran` | **0** |
| Rujukan uji `P-xx` / `P5-xx` | **1** — satu kalimat di *Prior art*, bukan pemetaan per persyaratan |
| Rujukan deviasi dengan nomor | **0** — register dirujuk sebagai "dua puluh dua baris", tanpa satu nomor pun |
| Tanda bukti (`EVIDENCED` / `DIPULIHKAN` / `TAFSIR`) | **2**, keduanya `TAFSIR (menunggu C-01)` |

**Akibatnya langsung pada tahap berikutnya:** tiket lahir dari persyaratan. Tidak ada
persyaratan, sehingga langkah 3 dan langkah 4 audit ini tidak dapat dijalankan sebagaimana
dirancang — bukan karena persyaratannya buruk, melainkan karena lapisnya tidak ada.

Ini **bukan** berarti isinya hilang. Sebagian besar substansi keenam belas bagian memang
mendarat, dalam bentuk paragraf keputusan. Yang hilang adalah **bentuk yang dapat dipecah**:
satu aturan, satu jalur gagal, satu tanda, satu uji, satu nomor.

---

## 1. Peta serapan enam belas bagian

| # | Bagian `to-spec` | Mendarat di | Utuh? | Yang hilang |
|---|---|---|---|---|
| 1 | Cakupan & yang sengaja tidak ditulis | *Out of Scope*; kalimat "tidak memutuskan apa pun tentang uang" ada di kepala | **SEBAGIAN** | Tabel aliran (A-1a/A-1b/A-2/A-3 dengan statusnya) tidak ada · cacah persyaratan tidak ada karena persyaratannya tidak ada · cacah pagar tidak dinyatakan (ada 8, tidak dicacah) |
| 2 | Pembentukan sirkulasi (A-1b) | *Implementation Decisions* → Tabel seleksi dan roster | **SEBAGIAN** | Seluruh butir substantif ada: tiga masukan, setiap kombinasi punya aturan, galat konfigurasi, jumlah jenjang turunan, derajat dari roster, nilai penentu kelas, jalur `K6-4`, roster idempoten, sirkulasi tanpa jenjang ditolak. Yang hilang: bentuk `S-xxx` dan jalur gagal per aturan |
| 3 | Validasi & kegagalan pembentukan | *Implementation Decisions* → Batas modul dan transaksi; *Solution* | **SEBAGIAN** | Satu transaksi, kegagalan terlihat, tali dari sisi sirkulasi — ketiganya dinyatakan. Yang hilang: galat mana dan pesan apa per jalur gagal; `H-5` (validasi menolak, bukan menandai) tidak dirujuk nomornya di mana pun |
| 4 | Maksud dan akibat | *Implementation Decisions* → Maksud dan akibat | **YA** | — |
| 5 | Mesin keadaan | *Implementation Decisions* → Mesin keadaan | **YA** | `stateDiagram-v2` ada; 11 invarian ada, masing-masing dengan cara membuktikannya dilanggar |
| 6 | Fungsi akibat | *Implementation Decisions* → Maksud dan akibat, tabel tujuh baris | **YA** | Satu tempat, tiga nilai tampil, laporan tidak membaca hasil telanjang — ketiganya ada |
| 7 | Wewenang & identitas | *Implementation Decisions* → Wewenang | **SEBAGIAN** | Gerbang menolak, kesetaraan penuh, kedudukan vs orang, hak baca — ada. Yang hilang: `P6-3` tidak dirujuk, sehingga pembaca tidak tahu bahwa sistem lama **tidak punya gerbang sama sekali**, dan perbedaan itu yang menjadikan ini deviasi, bukan penguatan |
| 8 | Pencatatan keputusan | *Implementation Decisions* → Pencatatan keputusan | **SEBAGIAN** | Catatan tetap, satu kolom waktu, peran disalin, versi usulan, dua peristiwa, penjaga di batas tulis, `K5-6` — seluruhnya ada. Yang hilang: `TIDAK_SAMPAI` dibahas di *Solution* dan invarian I-5, tetapi tidak di bagian pencatatan; **nomor urut sirkulasi per usulan (`F-11`) hanya muncul di tabel skema**, tidak sebagai aturan |
| 9 | Tulis-balik ke induk | *Implementation Decisions* → Tulis-balik ke klaim | **YA** | Sepuluh sasaran dipetakan; `ComiteeClaim` vs `ClaimComitee` dibedakan; empat sasaran tanpa kolom pendaratan diputus |
| 10 | Model data | *Implementation Decisions* → Skema | **SEBAGIAN** | Enam objek disebut dengan isi, asal-usul, aturan penamaan, presisi uang, sequence, zona waktu. Yang hilang: **kolom, tipe, kunci, dan constraint per tabel** — skill meminta "tiap tabel: kolom, tipe, kunci, constraint, satu kalimat asal-usul". Nol DDL |
| 11 | Kontrak API | *Implementation Decisions* → Kontrak API | **YA** | Metode, jalur, masukan, keluaran, galat, idempotensi lengkap; lima port hilir dinamai dan kosong, masing-masing berpagar |
| 12 | Layar (A-2) | *Implementation Decisions* → Layar | **SEBAGIAN** | Pernyataan cakupan ada; kunci dua kali ada; blok `1=2`/`1=3`/`NEVER` dinyatakan dibuang. Yang hilang: **tabel medan** (sumber lama · kendali · wajib · baca-saja bila · tampil bila · validasi klien · validasi server). Tanpanya 296 kendali tidak dapat dipecah menjadi tiket |
| 13 | Penomoran (A-3) | *Implementation Decisions* → Penomoran akseptasi | **YA** | Syarat terbit, idempotensi constraint, terbit sekali, galat gangguan dapat dibedakan, panjang dijaga hilir, ambang parameter bernama berpagar |
| 14 | Register pagar | *Out of Scope*, tabel delapan baris | **SEBAGIAN** | Delapan pagar ada dengan register dan baris inventaris. Yang hilang: **kolom `Titik sentuh`** (persyaratan mana yang berhenti di sana) — tidak dapat diisi karena persyaratannya tidak ada |
| 15 | Deviasi | *Testing Decisions* → Prior art | **TIDAK** | Register dirujuk sebagai "dua puluh dua baris" dan "dua puluh belum diratifikasi", **tanpa satu nomor pun**. Skill memerintahkan rujukan **dengan nomor**; tanpa itu tidak ada satu pun uji deviasi yang dapat ditarik ke sebuah tiket. Status belum-diratifikasi disebut — itu bagian yang benar |
| 16 | Ketertelusuran | — | **TIDAK** | Tidak ada rantai rule lama → `S-xxx` → `P-xx` di mana pun. Sebelas rule inti sistem lama (`KomiteRouter`, `KomitePostAdjustment`, `KomitePostAdjustmentCWP`, `CreateChildKomiteCNP_Act`, `CreateChildKomiteCloseNP_Act`, `FilterEmailKomiteWithLimit`, `IsKomiteLoop`, `ShowTransfer`, `KomiteTreaty_Flow`, `InsertXOLKlaimCNP`, `HitServiceToKasirKMT_Act`) tidak disebut satu pun di dalam spesifikasi |

**Cacah:** `YA` = 6 · `SEBAGIAN` = 8 · `TIDAK` = 2.

Dua yang paling mudah hilang di bentuk naratif memang hilang: **ketertelusuran** (`TIDAK`)
dan **deviasi dirujuk dengan nomor** (`TIDAK`). Keduanya persis yang skill peringatkan.

---

## 2. Uji kesiapan tiket

| Uji | Cacah lolos | Cacah gagal | Nomor yang gagal |
|---|---:|---:|---|
| T-1 · Dapat dikerjakan sendiri | — | — | **tidak dapat dijalankan** |
| T-2 · Punya "Bila gagal" | 0 | — | **tidak dapat dijalankan** |
| T-3 · Punya uji | 0 | — | **tidak dapat dijalankan** |
| T-4 · Punya tanda bukti yang sah | 0 | — | **tidak dapat dijalankan** |
| T-5 · Merujuk ketetapan dengan nomor | — | — | **tidak dapat dijalankan** |
| T-6 · Tidak menyeberang pagar | — | — | **tidak dapat dijalankan** |

Keenam uji berlaku atas `S-xxx`. Cacah `S-xxx` adalah **nol**, sehingga tidak ada objek
yang diuji. Itu dilaporkan sebagai ketidakmampuan menjalankan uji, **bukan** sebagai nol
kegagalan — nol kegagalan atas nol objek adalah kalimat kosong.

**Yang dapat diperiksa tanpa `S-xxx`**, diperiksa dan hasilnya:

| Hal | Hasil |
|---|---|
| Ketetapan dirujuk dengan nomor, bukan ditulis ulang | **LOLOS.** Seluruh rujukan berbentuk nomor (`D-1`…`D-5`, `E-1`…`E-5`, `J-1`…`J-4`, `H-3`, `H-4`, `H-6`, `H-7`, `K5-1`…`K5-7`, `K6-1`…`K6-4`, keputusan beku no. 1–9, `ADR-0003/0006/0016/0030/0031/0032`, `F-1`, `F-6`, `F-7`, `F-9`, `F-11`, `F-14`, `F-17`, `C-01`, `K-07`) |
| `H-2` tidak dipakai | **LOLOS.** Disebut hanya untuk menyatakan ia tidak mengikat, dan `K6-4` yang menggantikannya |
| Nol nama orang | **LOLOS.** Jalur khusus satu operator disebut tanpa namanya |
| Nol angka ambang telanjang | **LOLOS.** Isi awal tabel seleksi dinyatakan disemai, angkanya tidak ditulis |
| Nol angka yang asalnya hanya deskripsi langkah | **LOLOS.** Angka yang dilarang `K6-1` tidak muncul |
| Nol kata terlarang `CONTEXT.md` bagian 7 sebagai nama | **LOLOS** untuk nama baru. Satu kolom milik sisi Klaim memakai kata itu; spesifikasi menyatakan tidak menamainya ulang dan tidak membuat nama baru dengannya — penanganan yang benar |
| Nol kode Go/JSX | **LOLOS** |
| Nol persyaratan untuk A-4, A-5, isi A-6 | **LOLOS.** Ketiganya hanya muncul sebagai pagar dan port kosong |
| Alasan tunduk pada pagar | **LOLOS.** Tidak ada argumen yang bersandar pada isi aliran beku; port hilir dinamai tanpa menyatakan isinya |

---

## 3. Pasangan cerita–persyaratan

Pemasangan formal tidak dapat dijalankan: satu sisi pasangan tidak ada. Sebagai pengganti
terdekat, tiap cerita diperiksa terhadap **paragraf keputusan** yang menopangnya.

| Keadaan | Cacah | Nomor |
|---|---:|---|
| Bertopang keputusan | 43 | 1–41, 43, 44 |
| **Cerita tanpa topangan** | 1 | **42** |
| **Keputusan tanpa cerita** | 4 | lihat di bawah |

**Cerita 42** — "sirkulasi lama terbaca di sistem baru dengan riwayat keputusannya" —
bertentangan dengan *Out of Scope*, yang menyatakan migrasi data di luar cakupan dan
bentuknya bergantung pada satu cacah yang belum diambil (`INVENTARIS-BUKTI.md` §2.5 baris 6).
Cerita itu menjanjikan yang tidak dispesifikasikan.

**Empat keputusan tanpa cerita:**

1. Blok layar bersyarat `1=2`, `1=3`, `NEVER` dibuang dan pembuangannya dicatat.
2. Tiap kunci layar ditegakkan dua kali — klien dan server.
3. Lima port hilir dinamai dan dibiarkan kosong.
4. Penanda "sedang disirkulasikan" tidak dimigrasikan sebagai kolom; ia turunan (`D-2`).

Keempatnya sah sebagai keputusan, tetapi tidak ada cerita yang memintanya — sehingga bila
dipecah menjadi tiket, tiketnya tidak dapat menyebut siapa yang diuntungkan.

---

## 4. Delapan pagar

| PAGAR | Aliran | Register | Baris inventaris yang ditunjuk | Baris itu ada? | Benar belum dipegang? | Sah? |
|---|---|---|---|---|---|---|
| PAGAR-01 | A-3 | `PG-03` | §2.5 baris 4 — isi badan prosedur penerbit pada basis data berjalan | **ya** | ya — "nol", yang dipegang hasil reverse-engineer | **SAH** |
| PAGAR-02 | A-5 | `PG-05` | §2.1 — "dua prosedur akseptasi" | ya | ya | **SAH, penamaan kurang** |
| PAGAR-03 | A-5 | `PG-05` | §2.1 — "idem" | ya | ya | **SAH, penamaan kurang** |
| PAGAR-04 | A-5 | `PG-05` | §2.1 — "sasaran rincian layer" | ya | ya | **SAH, penamaan kurang** |
| PAGAR-05 | A-4 | `PG-04` | §2.5 baris 1 — "nol baris log kiriman" | ya | ya | **SAH, penamaan kurang** |
| PAGAR-06 | A-4 | `PG-04` | §2.5 baris 1 — "idem" | ya | ya | **SAH, penamaan kurang** |
| PAGAR-07 | A-6 | `PG-06` | §2.3 — "rule penyusun surat jenjang berikutnya" | ya | ya | **SAH, penamaan kurang** |
| PAGAR-08 | A-6 | `PG-06` | §2.3 — "rule pengirim berlampiran" | ya | ya | **SAH, penamaan kurang** |

**Nol pagar gugur.** Kedelapan rujukan menunjuk baris yang benar-benar ada dan benar-benar
belum dipegang; tidak satu pun menunjuk sesuatu yang ternyata ada di repo. Pola yang sudah
sembilan kali terjadi tidak terulang di sini.

**Tetapi enam dari delapan tidak menamai objeknya.** Aturan kerja `B` pada `KETETAPAN.md`
bagian 7 menuntut pagar menyebut "nomor bagian **dan nama objeknya**, bukan kalimat 'belum
ada'". `PAGAR-02`, `-03`, `-04` menunjuk §2.1 tanpa menyebut `PEGA_JSON_OS_AKSEP_KLAIM`,
`PEGA_JSON_OS_AKSEP_SUBJECTIVITY`, atau `XOL2_AKSEP_KLAIM`. `PAGAR-05` dan `-06` menunjuk
§2.5 baris 1 tanpa menyebut `POOLDATA.DIRECTTOKASIR_LOG`. `PAGAR-07` dan `-08` menunjuk §2.3
tanpa menyebut `PostEmailKomiteCNP` dan `SendEmailWithAttachments`. Dua di antaranya
("idem") bahkan menunjuk baris sebelumnya, bukan inventaris.

Itu tidak membatalkan pagarnya — objeknya dapat ditemukan — tetapi ia melanggar aturan yang
justru dipasang agar pembaca berikutnya tidak perlu menebak.

**Kolom `Sementara ini` — diperiksa satu per satu.** Tujuh dari delapan menamai lubang
(`PORT_… kosong`, "parameter bernama tanpa nilai", "tidak ada persyaratan; lubang
dinyatakan"). Satu, `PAGAR-06`, berbunyi "Niat dicatat sebelum panggilan dengan kunci
idempotensi per efek; penjaganya tidak dirancang". Klausa pertamanya adalah **ketetapan
yang sudah ada** (keputusan beku no. 6), bukan pendekatan sementara yang dikarang, dan
klausa keduanya menamai lubang. **Lolos**, dengan catatan bahwa bentuknya paling dekat ke
batas di antara kedelapannya.

**Satu ketidakcocokan di dalam bagian itu sendiri:** kalimat pembuka *Out of Scope* berbunyi
"Tiga aliran berpagar", sedangkan tabelnya memuat **empat** — A-3, A-4, A-5, A-6.

---

## 5. Seam — temuan, bukan keputusan

Usulan yang diaudit: batas aplikasi modul (`BentukSirkulasi`, `CatatKeputusan`,
`BacaSirkulasi`), diuji dalam proses terhadap skema Oracle nyata; HTTP dan layar sebagai
adapter tipis.

**Apakah ketiganya menutupi seluruh aliran terbuka?** Tidak seluruhnya. Dua operasi
administratif berdiri di luar ketiganya:

- Pengelolaan **roster** (`D-5`, `H-3`) — menopang cerita 35, 36, 39, 40.
- Pengelolaan **tabel seleksi** (`K5-5`, `K6-1`, `K6-3`) — menopang cerita 37, 38, dan
  penolakan konfigurasi yang cakupan kombinasinya tidak lengkap.

Keduanya muncul di *Kontrak API* sebagai `PUT /roster-jenjang` dan `PUT /aturan-seleksi`,
tetapi tidak muncul di seam. Seam yang diusulkan karena itu **tiga operasi, bukan satu**,
dan sebenarnya **lima** bila kedua operasi administratif ikut diuji lewat batas yang sama.

**Apakah ada yang hanya dapat diuji lewat HTTP atau layar?** Ya, satu keluarga:

- Aturan layar A-2 — kendali baca-saja, medan wajib, kunci jenjang pertama, blok yang
  dibuang. Separuhnya (sisi server) jatuh ke dalam seam; separuh klien tidak. Klaim "adapter
  tipis" berlaku untuk HTTP, **tidak** berlaku untuk layar: spesifikasi sendiri menuntut tiap
  kunci ditegakkan **dua kali**, dan penegakan kedua itu hidup di luar seam.

**Apakah invarian dapat ditegakkan di dalam seam?** Delapan dari sebelas ya. Tiga tidak:

- **I-2** (derajat unik) ditegakkan constraint basis data, bukan kode modul.
- **I-8** (jumlah jenjang tidak pernah disimpan sebagai kolom) hanya dapat dibuktikan dengan
  memeriksa **katalog skema**, bukan dengan memanggil operasi mana pun.
- **I-10** (satu kolom waktu per jenjang) sama — pemeriksaan bentuk, bukan pemeriksaan
  perilaku.

Ketiganya tetap dapat diuji **selama uji berjalan terhadap skema nyata**, seperti yang
usulan seam nyatakan — tetapi bentuk ujinya adalah pemeriksaan skema, dan itu perlu
dinyatakan agar tidak dikira uji perilaku yang tertinggal.

Keputusan seam **tidak diambil di sini**. Milik pemilik proses.

---

## 6. Putusan

**BELUM SIAP.**

Sebabnya satu, dan bukan soal mutu isi: **lapisan persyaratan tidak ada.** Tiket lahir dari
persyaratan, dan cacah `S-xxx` adalah nol. Enam uji kesiapan tiket tidak dapat dijalankan,
kolom `Titik sentuh` pada delapan pagar tidak dapat diisi, dan ketertelusuran tidak dapat
dibangun karena tidak ada yang dapat ditunjuk di tengah rantai.

Yang **tidak** menjadi sebab: substansi. Enam bagian terserap utuh, delapan sebagian dengan
butir yang dapat disebut satu per satu, kedelapan pagar sah dan nol yang gugur, seluruh
ketetapan dirujuk dengan nomor, dan sembilan pemeriksaan disiplin lolos. Pekerjaannya bukan
menulis ulang, melainkan **menambahkan lapisan yang hilang di atas isi yang sudah ada**.

---

## 7. Yang harus diperbaiki sebelum `to-tickets`

Urut dari yang memblokir.

1. **Tulis lapisan `S-xxx`.** Tiap keputusan yang sudah ada dipecah menjadi persyaratan
   berbentuk `Aturan` · `Bila gagal` · `Tanda` · `Beda AS-IS` · `Uji`. Ini memblokir seluruh
   tahap berikutnya: tanpa persyaratan tidak ada tiket, tidak ada kriteria selesai untuk
   jalur gagal, dan tidak ada yang dapat ditunjuk oleh pagar maupun ketertelusuran.
2. **Bangun bagian Ketertelusuran.** Rantai rule lama → `S-xxx` → `P-xx`, memuat sebelas rule
   inti yang saat ini tidak disebut satu pun. Rule yang jatuh ke aliran beku diberi nomor
   pagar, bukan dikosongkan. Tanpa ini shadow-run kehilangan petanya dan tiket tidak dapat
   menyebut apa yang digantikannya.
3. **Rujuk deviasi dengan nomor.** Dua puluh dua baris register ditarik ke persyaratan yang
   bersangkutan lewat nomornya, bukan disalin dan bukan dirujuk sebagai jumlah. Yang belum
   diratifikasi tetap bertanda demikian.
4. **Lengkapi model data.** Kolom, tipe, kunci, dan constraint per tabel untuk keenam objek,
   plus kolom maksud pada tabel klaim. Tanpa ini tiket skema tidak dapat ditulis.
5. **Lengkapi tabel medan layar.** Sumber lama · kendali · wajib · baca-saja bila · tampil
   bila · validasi klien · validasi server. Tanpa ini 296 kendali tidak dapat dipecah.
6. **Namai objek pada enam pagar.** `PAGAR-02`, `-03`, `-04` menyebut ketiga objek §2.1;
   `PAGAR-05`, `-06` menyebut `POOLDATA.DIRECTTOKASIR_LOG`; `PAGAR-07`, `-08` menyebut kedua
   rule §2.3. Ganti dua "idem" dengan rujukan penuh. Aturan kerja `B`.
7. **Isi kolom `Titik sentuh`** pada kedelapan pagar begitu `S-xxx` ada.
8. **Tambahkan tabel aliran dan cacah** pada bagian cakupan: status A-1a, A-1b, A-2, A-3,
   cacah persyaratan, cacah pagar.
9. **Perbaiki "Tiga aliran berpagar"** menjadi empat, atau pisahkan `PAGAR-01` dari ketiga
   aliran beku — A-3 terbatas, bukan beku.
10. **Selesaikan cerita 42.** Migrasi data ada di *Out of Scope*, sehingga cerita itu
    menjanjikan yang tidak dispesifikasikan: cabut ceritanya, atau tarik migrasi baca ke
    dalam cakupan.
11. **Beri cerita pada empat keputusan yang tidak punya**, atau nyatakan keempatnya sebagai
    keputusan teknis yang memang tidak berasal dari permintaan siapa pun.
12. **Nyatakan dua operasi administratif di dalam seam** — roster dan tabel seleksi — atau
    nyatakan tegas bahwa keduanya diuji di seam lain, beserta alasannya.
13. **Nyatakan bahwa klaim "adapter tipis" tidak berlaku untuk layar**, dan sebutkan bagaimana
    penegakan kunci sisi klien diuji.
14. **Tandai I-2, I-8, I-10 sebagai pemeriksaan skema**, bukan pemeriksaan perilaku, agar
    bentuk ujinya tidak dikira tertinggal.

Butir 1 sampai 3 memblokir `to-tickets`. Butir 4 dan 5 memblokir sebagian tiket, bukan
seluruhnya. Butir 6 sampai 14 tidak memblokir pemecahan tiket, tetapi seluruhnya harus
tertutup sebelum spesifikasi ini dipakai sebagai rujukan tetap.
