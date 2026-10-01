# PROMPT — GELOMBANG 3, TIGA MODUL KLAIM SEKALIGUS: **Claim Prop · Komite Claim Prop · Claim Fac In** *(sesi baru, folder `D:\XML\RNM_BRD\OUTPUT_HASIL_RNM`, cabang **`dev`**)*

> Permintaan work owner 01-10-2026: *"kerjakan 3 modul besar sekaligus, jangan hanya 1 tiket, jangan banyak yang di-skip, baca ulang XML, ambil logic,
> button, dan method yang benar dari XML, jangan membuat menu yang tidak ada di Pega; XML patokan dasar; tiket = hasil grilling; tiket yang salah
> diperbaiki di dokumennya."*
>
> **Mengapa tiga ini:** ketiganya `ready-for-agent` dengan spec, 13–16 tiket, dan 5–8 ronde grilling. Polanya sudah ada di aplikasi — Claim Life +
> Komite Claim Life — jadi bentuk layanan, tahap, kotak kerja, dan penyerahan ke Komite **ditiru**, bukan dirancang ulang. Claim Fac In memakai skema
> lini PROP, jadi Claim Prop dikerjakan lebih dulu di setiap gelombang.
>
> ⛔ **Syarat mulai:** sesi gelombang 2 *(Endorsement)* dan sesi seragam kolom `T_WORK_POLIS` sudah selesai melapor, dan `git status --porcelain` kosong.

## 0. ATURAN BACA XML

1. **Setiap pembacaan activity mencetak `pyStepsBlockName`**. `//` = langkah ter-remark, tidak pernah jalan, jangan dibangun.
2. Nomor baris **hanya** dari `sed -e 's/></>\n</g' <berkas> | grep -n …`; tulis perintahnya di kepala PARITAS. *(Dua sesi terakhir memakai penomoran yang
   bergeser 5–65 baris dari perintah ini.)*
3. WHEN dari urutan aksi. Precondition `WhenTrue 2` = lanjut, `3` = lewati. Visibilitas sel `pyVisible OTHER` + `pyCondition 1=2` = **tombol mati**.
4. Cari **penulis** kolom, bukan hanya pembacanya.
5. Setiap tombol = tombol di section XML yang memanggil activity yang sama, label **VERBATIM**. Nol menu, layar, tombol, aksi tanpa bukti XML.
6. Aksi XML yang perilakunya jelas tetapi tidak ada di tiket = **dibangun** dan dicatat sebagai tambahan tiket. Yang tidak jelas = OQ dengan bukti.
7. Tiket yang berbeda dari XML atau katalog DEV diperbaiki di dokumennya *(catatan bertanggal, kalimat lama dikutip)*. `grilling-ronde-*` tersegel.

## 1. BATAS KERJA

| Hal | Isi |
| --- | --- |
| Cabang | **`dev`**. `git commit -o -- <jalur>`, pesan berawalan nama modul. **Nol `git push`, nol `git pull`, nol `git pull --rebase`** |
| Folder boleh disunting | `APP_RNM/modul/claimprop/**`, `APP_RNM/modul/komiteclaimprop/**`, `APP_RNM/modul/claimfacin/**`, ketiga `inti/backend/daftar/modul_<nama>_gen.go` |
| Pengecualian, paket 0 *(satu commit)* | lihat §2 — hanya berkas yang disebut di sana |
| Dilarang | `inti/` selain berkas bangkitan, modul lain di luar §2, `package.json`, `go.mod`, konfigurasi root. Pola Claim Life dan Komite **ditiru**, **tidak diimpor**; lintas modul hanya lewat `inti/backend/kontrak` yang sudah ada — kontrak baru = catat, jangan buat |
| Nomor *(ralat 01-10-2026, commit `5fa9e48`)* | Claim Prop **520–559** / slot **980** · Claim Fac In **560–599** / slot **978** · Komite Claim Prop **680–719** / slot **986**. Migrasi Fac In yang merujuk tabel Prop otomatis berjalan sesudahnya |
| Oracle | skema `POOLDATA`. `-migrate` oleh **work owner**. Uji `db` tanpa `ORACLE_DSN` = SKIP, dijalankan dengan `-p 1`; POOLDATA bukan target uji `db`. Nol penulisan ke DEV |
| Prosedur | tidak dipanggil, logikanya ditiru di Go. Nol `COMMIT` di teks SQL |
| Layanan luar | `ConnectREST` *(Claim Prop 6, Komite Claim Prop 3, Claim Fac In 8)*, `LinkService`, kasir, Arasapas, email, PDF: **stub outbox**; alamat nyata tidak dipanggil dan tidak ditulis |
| Rahasia | nol nama orang, nomor polis, kredensial, host/IP, URL `M_LINK_SERVICE`, baris `EMAILKOMITE`. Fixture `UJI-` |
| Tabel kerja | setiap `T_WORK_*` baru mengikuti aturan kolom wajib *(`ID`, `COVER_KEY`, `LINI`, `STATUS_WORK`, `CREATE_OP`, `CREATE_OP_NAME`, `TGL_CREATE`, `TGL_UPDATE`)* |
| CSS | kelas berawalan modul, di berkas CSS folder frontend modul itu |

## 2. PENGECUALIAN — paket 0, satu commit, sebelum yang lain

| # | Masalah *(diukur asisten 01-10-2026)* | Berkas | Ubah |
| ---: | --- | --- | --- |
| X1 | `T_WORK_CLAIM` dipakai bersama lini *(kolom `LINI`, tiket 00 Claim Fac In: tabel dipakai bersama, keputusan K1·K2)*, tetapi **kotak masuk Claim Life tidak menyaring `LINI`** *(`modul/claimlife/backend/repository/inbox.go`: `WHERE NVL(w.TAHAP, …) = :tahap [AND w.CREATE_OP = :akun]`)*. Kasus Prop atau Fac yang ditulis ke tabel itu akan **muncul di inbox Claim Life** | `modul/claimlife/backend/repository/inbox.go` + ujinya | tambah `AND w.LINI = :lini` dengan nilai `inti.LiniLife` di sqlInbox dan sqlCacahInbox. Perilaku Claim Life tidak berubah *(semua barisnya `LIFE`)*. Uji: baris `LINI` lain tidak tampil |
| X2 | Kotak kerja Komite Claim Life membaca `T_WORK_CLAIM`/daftar komite tanpa saringan lini | `modul/komiteclaimlife/backend/repository/komite_inbox.go` + ujinya | saringan lini yang sama, nilai lini Komite Claim Life dari XML-nya; uji serupa |
| X3 | penjaga `TestSetiapPemanggilBukaMemeriksaBolehDilewati` *(`modul/claimlife/backend/repository/batasanpemakaian_test.go`)* menghitung pemanggil `skemauji.Buka()` **di seluruh `APP_RNM`**, sehingga setiap modul yang menambah uji `db` harus menyunting berkas Claim Life *(merah 01-10-2026: 15 vs 14 karena uji baru PremiumList)* | berkas itu | persempit pindaiannya ke `modul/claimlife/` saja, seperti `e13ad9e`; aturannya tidak berubah. Penjaga yang sama per modul ditulis di folder modul masing-masing |
| X4 | `T_GENERAL_CLAIM` **tidak** punya `LINI` | — | **jangan** `ALTER`. Pembeda lini dibaca lewat `T_WORK_CLAIM.LINI` *(join `ID`)*. Bila XML menuntut lebih, OQ |

Nilai `LINI` untuk Prop dan Fac **dari XML** *(mis. `Param.STS_KLAIM` seperti Claim Life `"LIFE"` di `services/komite.go:42`)* — sebut buktinya.
Konstanta lini baru ditulis di folder modulnya sendiri, bukan di `inti/backend/lini.go`.

## 3. FAKTA DEV *(asisten 01-10-2026, agregat, baca-saja)*

| Tabel warisan | Baris | Dipakai |
| --- | ---: | --- |
| `FACINPRODUCTION` | ±1,88 juta | Claim Fac In |
| `OS_AKSEPTASI_KLAIM` | ±39,7 ribu | Claim Fac In *(28 rujukan XML)*, cermin klaim |
| `TREATYINPRODUCTION` | ±41,9 ribu | Claim Prop |
| `TREATYINDETAILEDM` | ±11,8 ribu | Claim Prop |
| `TREATYBUSINESS`, `TREATYREINSURER`, `TREATYYEAR`, `PROPORTIONALARRG` | 4.553 / 430 / 178 / 3.188 | Claim Prop, Komite Claim Prop |
| `MONITORING_KLAIM_LOG` | ±24,4 ribu | Komite Claim Prop |
| `JSON_POLIS` | 148 | Claim Prop, Claim Fac In |
| `T_STORAGE_IMAGE`, `T_FOLDER_IMAGE` | 517 / 1 | lampiran ketiganya |

Tabel aplikasi bersama: `T_WORK_CLAIM` *(punya `LINI`, kini 0 baris)*, `T_GENERAL_CLAIM` *(tanpa `LINI`)*.

**Aturan tabel besar:** setiap kueri ke tabel di atas memakai kunci ber-index *(sebut indeksnya)*, **nol pemindaian penuh**, uji tidak menyentuhnya.

## 4. MODUL A — CLAIM PROP *(329 XML: Activity 112, Section 36, RDBList 58, RD 19, Harness 11, FlowAction 12, When 61, ConnectREST 6)*

Dokumen: `modul/claimprop/docs/` — `spec.md` *(ready-for-agent, 18-09-2026)*, `grilling-ronde-1`–`6`, `STRUKTUR-TABEL-CLAIM-PROP.md`,
`RELASI-TABEL-CLAIM-PROP.md`, `periksa-ulang-aturan-baru.md`, `utang-lintas-modul.md`, `issues/00`–`15`.

| Butir | Kerjakan |
| --- | --- |
| tiket 00 skema | tabel aplikasi baru lini PROP **boleh** *(pola Claim Life `T_CLAIMLF_*`)*, di rentang 520–559, nama tabel dari `STRUKTUR-TABEL-CLAIM-PROP.md`. Tabel **warisan** tidak di-DDL. Kolom `FLAG_ON_GOING_COMMITTEE` di `T_GENERAL_CLAIM` **tidak dibuat** *(keputusan ① tiket 00)*. Enam properti spreading = tabel berbeda bila tiket menetapkannya |
| tiket 01–15 | semuanya, urut `urutan-tiket` bila ada; setiap tiket dicocokkan dengan PARITAS lebih dulu |
| kasus | baris `T_WORK_CLAIM` dengan `LINI` Prop dari XML; pengenal kasus dari XML *(awalan + sequence modul sendiri)* |
| hulu | polis treaty dibaca langsung dari tabel warisan *(NB Treaty In belum dibangun)* — baca-saja, kolom yang dibaca XML saja |

## 5. MODUL B — KOMITE CLAIM PROP *(80 XML: Activity 35, RDBList 23, RD 4, Section 2, FlowAction 2, When 6, ConnectREST 3, nol Harness)*

Dokumen: `modul/komiteclaimprop/docs/` — `spec.md`, `grilling-ronde-1`–`8`, `STRUKTUR-TABEL-KOMITE-CLAIM-PROP.md`, `RELASI-TABEL-KOMITE-CLAIM-PROP.md`,
`utang-lintas-modul.md`, `issues/00`–`13`, `urutan-tiket.md`.

| Butir | Kerjakan |
| --- | --- |
| tiket 00 skema | pola Komite Claim Life *(baris `T_WORK_CLAIM` berawalan komite + `COVER_KEY` ke klaim induk)*; tabel komite baru hanya bila `STRUKTUR-TABEL-KOMITE-CLAIM-PROP.md` menetapkannya, di 680–719 |
| tiket 01–12 | penyerahan dari Claim Prop, kotak kerja dan giliran, layar tiga wajah, keputusan, persetujuan bersyarat, tangga, penolakan, nomor akseptasi, uang, efek keluar DB, PDF, kasir/Arasapas/email *(stub)* |
| tiket 13 | migrasi data lama — **di luar lingkup** kecuali tiket menyatakan sebaliknya; tetap terbuka |
| penyerahan | Komite Claim Life menerima klaim lewat kontrak **`kontrak.KlaimKomite`** yang sudah ada di `inti/backend/kontrak`, disediakan Claim Life. Bila antarmuka itu cukup umum untuk lini lain, **Claim Prop menyediakan implementasinya sendiri** dan Komite Claim Prop memakainya — tanpa saling impor. Bila antarmukanya harus diubah atau kontrak baru dibutuhkan, **berhenti di tiket itu**, catat OQ, lanjutkan tiket lain |
| menu | **tidak ada harness** di korpus: kotak kerja Komite dibuka lewat aksi yang XML tunjukkan *(cari pemanggilnya)*; tombol menu hanya bila ada jalan masuk portal berbukti |

## 6. MODUL C — CLAIM FAC IN *(482 XML: Activity 179, Section 57, RDBList 63, RD 27, Harness 16, FlowAction 33, DataPage 8, When 60, ConnectREST 8)*

Dokumen: `modul/claimfacin/docs/` — `spec.md`, `grilling-ronde-1`–`5` *(+ ulang)*, `STRUKTUR-TABEL-CLAIM-FACIN.md`, `RELASI-TABEL-CLAIM-FACIN.md`,
`issues/00`–`14`, `urutan-tiket.md`.

| Butir | Kerjakan |
| --- | --- |
| tiket 00 skema | dua tabel baru *(objek pertanggungan dan item objek, khas FAC)* + kolom induk kedua, **expand–contract** seperti tiket 00 menetapkannya; di 560–599, **sesudah** tabel Prop ada. Bila tiket menuntut `ALTER` pada tabel milik Claim Prop atau Claim Life, migrasinya di folder Fac In dengan nomor 560-an dan diumumkan di laporan |
| tiket 01–13 | semuanya |
| tiket 14 | migrasi paritas dan penamaan — kerjakan bagian penamaan; pemindahan data lama di luar lingkup |
| hulu | polis fakultatif dibaca langsung dari `FACINPRODUCTION`, `JSON_POLIS` *(NB FacIn belum dibangun)*, baca-saja |
| cermin | `OS_AKSEPTASI_KLAIM` ditulis **persis** seperti XML *(pola N13 Claim Life `OS_AKSEPTASI_KLAIM_LIFE`)*; sebut setiap penulisnya |

## 7. MENU

Satu tombol per modul dari baris `M_NAV_MENU` yang sudah ada di `KLAIM`: **Claim Prop**, **Komite Claim Prop**, **Claim Fac In**. Halaman awal = layar
pintu masuk portal menurut XML, dengan bukti. Slot hanya `UPDATE … SET DIMIGRASI = '1' … WHERE KODE = '<modul>'` + `_down`. Nol `INSERT`. Komite Claim
Prop tanpa jalan masuk portal berbukti → `DIMIGRASI` tetap `'0'`, OQ.

## 8. URUTAN — tiga modul bergantian, satu commit per paket, vertikal

| Gelombang | Claim Prop | Komite Claim Prop | Claim Fac In |
| ---: | --- | --- | --- |
| 1 | **pengecualian §2 (X1–X3)** + PARITAS + RALAT + OQ | PARITAS + RALAT + OQ | PARITAS + RALAT + OQ |
| 2 | tiket 00 skema + kerangka + baca polis treaty | tiket 00 skema + kerangka | *(menunggu tabel Prop)* kerangka + baca polis fac |
| 3 | register klaim, tahap awal | penerimaan penyerahan, kotak kerja, giliran | tiket 00 skema + register |
| 4 | outstanding, uang, spreading | layar tiga wajah, keputusan | tahap dan outstanding |
| 5 | penyerahan ke Komite, tutup klaim | tangga, penolakan, nomor akseptasi, uang | objek pertanggungan, uang |
| 6 | efek keluar *(stub)*, frontend lengkap, menu | efek keluar *(stub)*, frontend, menu | cermin `OS_AKSEPTASI_KLAIM`, efek keluar, frontend, menu |
| 7 | dokumen: status tiket, PARITAS, OQ, `MODUL.md` | sama | sama |

Setiap commit: `go build`, `go vet` dan `go test` *(dengan dan tanpa tag `db`; `db` dengan `-p 1`)*, `gofmt`, `tsc`, `vitest`, `vite build`, di pohon kerja
utama. Satu tiket yang terhenti **tidak** menghentikan modul lain. Sesudah gelombang 6: backend + `npm run dev` *(port **5174**)*, buka ketiga layar,
cocokkan tombol dengan PARITAS.

## 9. BERHENTI HANYA BILA

- Aturan hanya bisa dibangun dengan DDL pada **tabel warisan** — OQ, jangan tulis migrasi.
- Butuh kontrak lintas modul baru di `inti/backend/kontrak` atau menyunting modul lain di luar §2 — OQ, lanjutkan tiket lain.
- Perilaku aksi XML tidak dapat dipastikan — OQ dengan bukti `bNNN`, lanjutkan tiket lain.

## 10. LAPORAN

Satu pesan:

1. Tabel **modul → paket → commit → tiket → bukti XML → angka uji**.
2. Pengecualian X1–X3 beserta uji gigitnya.
3. Nilai `LINI` Prop dan Fac beserta buktinya.
4. Ralat dan tiket yang diperbaiki dokumennya.
5. Tombol dan aksi XML di luar tiket beserta statusnya.
6. OQ baru per modul.
7. Halaman awal tiap modul beserta buktinya.
8. Daftar migrasi baru per modul, dan pengingat `-migrate` untuk work owner.
9. Persentase kesiapan per modul dan total migrasi *(modul dimigrasi dari 20)*.
10. Bab **TELEMETRI EKSEKUSI**.

---

*Disusun 1 Oktober 2026 dari katalog DEV (tabel warisan dan cacah barisnya, kolom `LINI` di `T_WORK_CLAIM`/`T_GENERAL_CLAIM`), kode kotak masuk Claim Life
dan Komite Claim Life, penjaga `batasanpemakaian_test.go`, tiket 00 Claim Prop dan Claim Fac In, daftar rule tiga folder korpus, dan penukaran rentang
`5fa9e48`.*
