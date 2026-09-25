# Riwayat sesi — 19 September 2026

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09, rule termutakhir `pxUpdateDateTime = 2026-08-30`); `pengetahuan/ddl/` 49 berkas.
> Berkas ini mencatat **urutan keputusan dan sebabnya**, bukan isinya. Isi setiap keputusan hidup di berkas yang memakainya; bila keduanya berbeda, **berkas itu yang berlaku**.

**Transkrip mentah** sesi ini ada di `C:\Users\Administrator\.claude\projects\d--XML-NURE\a207e2e1-7a13-4c51-84ad-283831eb8a81.jsonl` (20 MB). Ia tidak disalin ke folder ini: 20 MB JSONL bukan artefak yang dapat dibaca orang, dan seluruh isinya yang berguna sudah turun menjadi keputusan di bawah.

---

## Kenapa berkas ini ada

Tiga hal yang tidak tersimpan di tempat lain:

1. **Urutan** — keputusan mana membuka keputusan mana. ADR menyimpan isi putusan, bukan rantainya.
2. **Sebab sebuah premis gugur** — beberapa keputusan hari ini membatalkan pekerjaan yang sudah dikerjakan. Tanpa catatan sebabnya, orang berikutnya akan mengerjakannya lagi.
3. **Koreksi atas diri sendiri** — lima angka dan satu klaim dicabut hari ini, dan **tujuh cacat ditemukan di DDL kami sendiri** saat tiketnya dikerjakan (bagian 10). Yang dicabut lebih mudah hilang daripada yang ditambahkan.

---

## 1. Register pertanyaan ditutup seluruhnya

**Sebelum**: aktif 42 · selesai 17. **Sesudah**: aktif **0** · selesai **57** · dihapus **2**.

Empat puluh butir ditutup, dua dihapus. **Sebagian besar tidak ditutup dengan jawaban**, dan itu pokoknya:

| Cara tutup | Berapa | Artinya |
|---|---|---|
| TIDAK BERLAKU | 5 | pertanyaannya mengandaikan sesuatu yang tidak ikut pindah |
| DI LUAR KEPEMILIKAN | 3 | perilaku modul lain; tidak dirancang, tidak ditanyakan |
| DITUTUP OLEH BENTUK | 8 | rancangannya benar di setiap jawaban yang mungkin |
| DIPUTUSKAN | 22 | keputusan diambil dan dicatat di berkas yang memakainya |
| DIHAPUSKAN | 1 | A13 — pengecualiannya dihapus, jadi tidak ada yang perlu berwenang |
| DIHAPUS dari register | 2 | A10b, A11b — jawabannya tidak mengubah satu kolom pun |

**Pola yang paling banyak dipakai: menutup dengan bentuk.** Sebuah pertanyaan berhenti menahan bukan ketika dijawab, melainkan ketika rancangannya benar di setiap jawaban yang mungkin. Itu yang menutup B4, C6, F, G1, H, dan I.

Rinciannya: `_selesai/OPEN-QUESTIONS.md`, bagian *Penutupan register*.

## 2. Keputusan pokok, berurut

| # | Keputusan | Membuka |
|---|---|---|
| 1 | **HRD keluar dari modul ini.** Pelaku disimpan sebagai **potret**, bukan rujukan hidup | K1 (jadi K1a), A10b, A11b, ADR-0016 ditulis ulang beserta judulnya |
| 2 | **ADR-0029 berlaku seragam**, kedua cabangnya dicabut | A16, C6, E14, ASK-AKUNTANSI butir 7, syarat batas ke-3 |
| 3 | **Presisi adalah keputusan maju**, bukan temuan yang menunggu data | B1, B2, E1, E2, E12; REQ-001 turun jadi verifikasi |
| 4 | **F, H, I ditutup sebagai hipotesis permanen** — status pertanyaannya, bukan cabangnya | C8, B8, tiket `17`, `18` |
| 5 | **Pengecualian pembekuan tambalan dihapuskan** | A13; ADR-0012 |
| 6 | **REQ-032 turun jadi verifikasi**; batas 30 byte jadi **aturan tetap** | gerbang penamaan dicabut; tiket `02` mati, `03` `04` lepas |
| 7 | **Gerbang `IS JSON` dicabut** dari tabel pendaratan | tiket `13` + `14` sebagai satu pekerjaan |

## 3. Premis yang gugur — dan pekerjaan yang ikut dibatalkan

### Tabel singkatan tertutup — tiket `02` mati

Premisnya: nama pengenal menabrak batas 30 byte, jadi singkatan diperlukan. **Salah.** Yang menabrak batas bukan bahasa melainkan **satu kebiasaan** — mencantumkan daftar kolom di nama constraint. `UQ_AKSEPTASI_KLAIM_LAYER_MATA_UANG` panjang karena ia **mengeja isinya**, padahal isi itu sudah ada di katalog Oracle.

Begitu nama constraint berhenti mengeja kolom, tekanannya hilang — dan **nama tabel serta kolom sendiri tidak pernah menabrak batas**.

Buktinya satu baris:

| Bentuk | Nama | Byte | Terbaca? |
|---|---|---:|---|
| mengeja kolom | `UQ_AKSEPTASI_KLAIM_LAYER_MATA_UANG` | 34 | ya, tetapi lewat batas |
| bersingkatan *(dibatalkan)* | `UQ_AKS_KLM_LYR_MTU` | 18 | tidak — perlu daftar untuk membacanya |
| **berlaku** | `UQ_AKSEPTASI_1` | **14** | ya |

Bentuk yang berlaku **lebih pendek sekaligus lebih terbaca** daripada bentuk bersingkatan. Itu yang membuat daftar singkatan bukan hanya tidak perlu, melainkan lebih buruk.

### Gerbang `IS JSON` di tabel pendaratan — kriteria tiket `13` dicabut

Premisnya: muatan yang bukan JSON sah ditolak di pintu. **Bertabrakan dengan dua fakta**:

- **D43** — `adoptJSONObject` mengadopsi teks `HASIL1` tanpa memeriksa bentuk. Bentuk yang **tidak dibatasi** tidak dapat dijaga dengan CHECK; ia hanya dapat ditampung.
- **D41** — sistem lama memang menghasilkan JSON rusak: `Nett`, `Deductible`, `KaliDeduct` dikirim tanpa kutip, dan nilai kosong menghasilkan `"Nett":,`.

Gerbang itu akan menolak **persis muatan yang paling perlu tercatat**, dan menolaknya melanggar AK-4 (tidak ada baris yang dibuang karena "cuma sedikit").

Penggantinya: muatan mendarat utuh, dan **hasil penguraian dicatat sebagai fakta** — `BENTUK_TERURAI`, `SEBAB_GAGAL_URAI`.

## 4. Angka yang dicabut — koreksi atas diri sendiri

Lima, dan tidak satu pun datang dari luar.

| Angka | Bunyi lama | Yang benar | Sebab |
|---|---|---|---|
| Pertanyaan terbuka | "aktif 41" | **42** | E17 masuk sesudah angka itu ditulis, tidak ikut dinaikkan |
| Pertanyaan selesai | "selesai 26" | **17** | tidak dapat direproduksi dengan perintah mana pun; daftar SELESAI saat itu memuat 17 butir |
| Kolom di kamus | "368" / "391", lalu **"452"** | **396** kolom tabel (+ 56 kolom view) | 368 dan 391 tidak dapat direproduksi. **452 ikut dicabut** kemudian hari itu juga oleh rekonsiliasi tabel datar: ia mencacah kolom tabel **dan** kolom view sekaligus, dan perintah yang menyertainya menghasilkan 457. Sekerabat dengan kelas keenam — bukan salah tulis, melainkan **pencacahnya menghitung populasi yang lain dari yang disebut namanya**. `sampah/keluaran/RINGKASAN-TABEL-DATAR.md` bagian 2.1 |
| Pengenal terpanjang | constraint `UQ_NILAI_KLAIM_MATA_UANG_1`, 26 byte | **view `V_TOTAL_NILAI_PERTANGGUNGAN`, 27 byte** | pemeriksaan pertama hanya mencacah tabel dan constraint, bukan view |
| Kolom terpanjang | `PENYEBAB_KERUGIAN_ID`, 20 | **`PREMI_PEMULIHAN_PORSI_IDR`, 25** | ditulis dari ingatan, lalu dihitung |
| Cacah pengenal DDL | "316", lalu "333" | **355** | **alat ukurnya yang cacat** — lihat di bawah |

**Kelas keenam, dan ia berbeda dari kelima yang di atasnya.** Lima yang pertama salah tulis atau salah ingat — keduanya dapat ditangkap dengan membaca ulang. Yang keenam **tidak dapat**: angkanya keluar dari sebuah alat, alat itu dijalankan, dan hasilnya disalin dengan benar. Yang cacat **alat ukurnya**.

`alat/periksa-penamaan.py` versi pertama hanya mencacah **tabel dan constraint**, dan melewatkan view, sequence, kolom, akun, peran, tablespace, serta kebijakan pengawasan. Ia melaporkan 316, lalu 333, dan **keduanya tampak terverifikasi** karena ada perintahnya.

Itu sebabnya "terpanjang 26 byte" juga keliru: pengenal terpanjang di skema ini **nama view**, dan alat itu tidak pernah melihat view.

> **Perintah yang menyertai sebuah angka membuktikan angkanya dapat direproduksi — bukan bahwa ia mengukur hal yang benar.** Alat ukur ikut diperiksa, bukan hanya keluarannya.

**Setiap angka di berkas hasil kini membawa perintah yang menghasilkannya.** Angka yang tidak dapat direproduksi dicabut, bukan diwariskan.

## 5. Pekerjaan yang selesai

| Tiket | Keluaran |
|---|---|
| `01` | ditinjau ulang terhadap aturan penamaan final — **nol nama berubah**; dua nama bukan-kata-utuh (`KLAIMNP`, `_APP`) dicatat terbuka beserta alasan tidak diganti |
| `03` | naik dari `selesai-sebagian` ke **`selesai`** — empat aturan penamaan final; 105 objek bernama ditulis ulang |
| `04` | naik ke **`selesai`** — `IX_ADJ_REKENING` → `IX_ADJUSTMENT_1` |
| `12` | `01_KLAIM.sql` ditulis ulang. Dua kolom baru: **`LINI_USAHA`** (§21.3) dan **`PENUTUPAN_LAMA`** (§21.6); empat aturan dinyatakan jatuh ke aplikasi |
| `13` + `14` | dikerjakan sebagai satu pekerjaan. Gerbang `IS JSON` dicabut; **`ID_PENDARATAN`** menjadi seam antara keduanya |
| `02` | **dibatalkan** ke `_mati/` — tiket mati pertama di papan ini |

**Penamaan**: satu tabel diganti nama karena constraint-nya melewati batas — `DAFTAR_KLAIM_PENJAGA_TANGGAL` (28 byte, constraint 33) → **`KLAIM_PENJAGA_TANGGAL`** (21, constraint 26). Kata `DAFTAR_` dibuang karena sebuah tabel memang sudah daftar. Kolom `ID_DAFTAR` ikut → `ID_PENJAGA_TANGGAL`.

~~**Diperiksa, bukan diperkirakan**: 316 pengenal unik di seluruh `ddl-usulan/`~~ — **DICABUT.** Angka 316 (dan 333 pada bagian kedua) adalah keluaran **alat yang cacat**: ia hanya mencacah tabel dan constraint, melewatkan view, sequence, kolom, akun, dan kebijakan. Yang benar, dari alat yang sudah diperbaiki: **355 pengenal unik**, nol di atas 30 byte, nol sisa singkatan, nol nama constraint ganda. Petanya `alat/peta-nama-2026-09-19.tsv`.

## 6. Keadaan pada pertengahan sesi — **digantikan bagian 13**

> Angka di bawah adalah keadaan **sebelum** REQ-032 turun derajat dan sebelas tiket berikutnya ditutup. Dipertahankan sebagai catatan perjalanan; **yang berlaku bagian 13**.

| Register | Keadaan |
|---|---|
| Pertanyaan | aktif **0** · selesai **57** · dihapus **2** |
| REQ | aktif **35** · selesai 1 · **mati 1** — tidak ada yang menahan keputusan rancangan |
| Tiket | aktif **15** · menunggu instance **8** · tertahan **1** · selesai **14** · mati **1** |

**Yang menahan pelaksanaan** — empat, dan hanya pelaksanaan:

| REQ | Menahan |
|---|---|
| REQ-018 | tiket `15`; penanganan migrasi di `12`, `18`, `29`, `34` |
| REQ-033 | pemetaan view kompatibilitas — `29`, `32` |
| REQ-021 | `35`, `39` |
| REQ-037 | klaim ADR-0017 di `35` |

Tiket `11` tertahan **lingkup**, bukan pengetahuan: prasyaratnya tinggal satu — batch aplikasi dimulai.

## 7. Batas yang tetap berdiri

**Register yang kosong bukan pengetahuan yang lengkap.** Tiga hal tetap tidak diketahui, dan ketiganya tercatat sebagai batas:

1. **Empat hipotesis permanen, delapan cabang** — F, H, I, G1. Cabangnya tertulis utuh di lampiran `_selesai/OPEN-QUESTIONS.md`; tidak satu pun dipilih. *Angka "tujuh" dicabut: ia mencampur satuan hipotesis dengan satuan cabang.*
2. **Tiga puluh lima permintaan belum dijawab** — 34 di berkas kedua, REQ-032 di berkas pertama.
3. **Folder `Komite Claim Non Prop` tidak dibuka** — tertutup untuk pekerjaan ini, bukan ditunda.

**Yang dapat membuka register kembali**: jawaban DBA yang menggugurkan sebuah ramalan, atau sesi Komite bila kelak terjadi.

## 8. Aturan yang berlaku sepanjang sisa pekerjaan

- Setiap kesimpulan lama yang gugur **dicabut di tempatnya**, bukan hanya dilaporkan.
- Setiap presisi merujuk **tabel domain**, bukan ditulis per kolom.
- Folder Komite **tertutup**.
- DDL yang ditulis harus **dapat dijalankan** — seam-nya uji SQL.
- Setiap angka membawa **perintah yang menghasilkannya**. Angka tanpa perintah tidak dihitung sebagai jawaban.
- Register menyimpan **keadaan**; laporan menyimpan **bukti**. Laporan tidak pernah menggantikan register.

---

# Bagian kedua — REQ-032 turun, dan seluruh papan lapisan data ditutup

## 9. REQ-032 ditutup sebagai penahan

**Batas 30 byte menjadi aturan tetap**, bukan pengamanan sementara. 30 byte sah di setiap versi Oracle; 128 hanya di sebagian. Pertanyaan yang jawabannya tidak mengubah apa pun bukan pertanyaan yang perlu ditunggu.

**Tabel singkatan tertutup dibatalkan** — tiket `02` mati, tiket mati pertama di papan ini. Premisnya gugur: yang menabrak batas bukan bahasa melainkan kebiasaan mencantumkan daftar kolom di nama constraint.

105 objek bernama ditulis ulang; satu tabel diganti nama. ~~316 → 333 pengenal diperiksa~~ — **dicabut**, keduanya keluaran alat yang cacat; yang benar **355**, lihat bagian 17.

## 10. Sebelas tiket ditutup, dan tujuh cacat ditemukan saat mengerjakannya

Papan lapisan data: **aktif 0 · menunggu instance 8 · tertahan 1 · selesai 29 · mati 1.**


### Cacah tiket — satu angka, satu perintah

Pesan sebelumnya menulis tiga angka berbeda untuk satu hal. Yang berlaku angka di bawah, dan **perintahnya ikut ditulis supaya dapat diperiksa ulang** — aturan yang berlaku untuk seluruh angka di proyek ini, dan yang belum sempat berlaku ke berkas riwayat ini sendiri.

**Tiket ditutup 19 September 2026: 24.**

```
grep -lE '^> \*\*(SELESAI|SELESAI PENUH|DIBATALKAN)( [A-Z]+)? 19 September 2026' \
  *.md _selesai/*.md _tertahan/*.md _mati/*.md | wc -l
```

`01` `02` `03` `04` `05` `06` `12` `13` `14` `15` `16` `17` `18` `19` `20` `21` `22` `23` `25` `26` `29` `30` `31` `35` — dua puluh tiga ditutup selesai, satu (`02`) dibatalkan.

**Cacah papan: 39, dan jumlahnya cocok.**

```
grep -h '^status:' *.md _selesai/*.md _tertahan/*.md _mati/*.md | sort | uniq -c
```

| Status | Cacah |
|---|---:|
| `selesai` | 29 |
| `menunggu-instance` | 8 |
| `tertahan` | 1 |
| `mati` | 1 |
| **Jumlah** | **39** |

Dan 39 sama dengan cacah berkas tiket: `ls *.md _selesai/*.md _tertahan/*.md _mati/*.md | grep -v README | wc -l`.

Yang penting bukan jumlahnya melainkan apa yang ditemukan:

| # | Cacat | Akibatnya bila lolos |
|---|---|---|
| 1 | **23 constraint pasangan IDR–kurs berarah terbalik** | baris ber-kurs terisi tapi IDR kosong **ditolak** — melanggar kriteria tiketnya sendiri |
| 2 | **Constraint yang sama saling mengunci** | mengisi **satu** nilai IDR memaksa **tiga lainnya** terisi. Tidak ada yang pernah memutuskan itu; ia akibat bentuk |
| 3 | **Asal-usul kurs tidak pernah ditegakkan** | ADR-0029 butir kelima; dan FINDING-006 — kurs `1` yang berarti "tidak ditemukan" tak dapat dibedakan tanpa asal-usulnya |
| 4 | **`DISUNTING_OLEH` tak punya constraint** | kriteria tiket `17` butir 5 ada sejak awal dan tak pernah ditegakkan: suntingan tanpa pelaku |
| 5 | **`KEPUTUSAN_KOMITE NOT NULL`** | keadaan **keempat** mustahil disimpan — dan keadaan keempat justru yang lolos penjaga `CloseClaimMD` (D2) |
| 6 | **`V_AKSEPTASI_KOMPATIBEL` tak menyentuh arsip** | ia mengeluarkan **nilai mutlak** ke jalur yang menjumlahkan **tambahan** — melipatgandakan setiap nilai pada pengiriman kedua |
| 7 | **44 `REVOKE` menggugurkan pemasangan bersih** | `ORA-01927`. **Kelas cacat yang sama** dengan sembilan `REVOKE` yang sudah diperbaiki di `01` — terulang karena tidak disapu ke seluruh berkas saat pertama ketahuan |

Cacat 7 adalah pelajarannya sendiri: **sekali sebuah kelas cacat ketahuan, ia disapu ke seluruh berkas, bukan diperbaiki di satu tempat.**

Tiga cacat lain yang lebih sunyi: `V_AKSEPTASI_DITOLAK` mengelompokkan berbeda dari pasangannya (sehingga ada kelompok yang jatuh dari keduanya — hilang, dan tidak terlihat hilang); `V_PARITAS_SHADOW` hanya melihat satu dari 19 tabel (sehingga penyaring keadaan membuang 18 tabel tanpa pesan); dan `ARSIP_MUATAN_KELUAR` memberi akun aplikasi `UPDATE`/`DELETE` padahal tiketnya menuntut tulis-sekali.

## 11. Satu berkas DML, dan ia sengaja dipisahkan

`ddl-usulan/Z00_ISIAN_AWAL.sql` — **satu-satunya berkas DML** di seluruh `ddl-usulan/`. Empat baris: tutup buku `25`, brokerage 2,5%, PPh 2%, PPN 2,2%.

Dipisahkan justru supaya kekecualiannya **terlihat**. Ke-22 berkas tabel, kesembilan view, dan berkas skema tetap menyatakan *"TIDAK ADA DML DI BERKAS INI"*, dan pernyataan itu tetap benar — diperiksa dengan `grep -lE "^(INSERT|UPDATE|DELETE|MERGE)" *.sql`, yang mengembalikan berkas ini saja.

## 12. Dua penyimpangan dari tiket yang dicatat, bukan ditutupi

1. **`ESTIMASI_AWAL` tidak diberi `NOMOR_URUT`** meski tiket `19` menyebut "ketiga tabel". Ia punya kunci alami yang sesungguhnya, dan kunci alami **lebih kuat**: nomor urut hanya melarang dua baris bernomor sama; kunci alami melarang dua baris **berarti sama**. Menambahkannya akan melemahkan tabel.
2. **Kriteria "muatan bukan JSON sah ditolak" dicabut** dari tiket `13`. Bentuk muatan lama tidak dibatasi (D43) dan sistem lama memang menghasilkan JSON rusak (D41); gerbang itu akan menolak persis muatan yang paling perlu tercatat.

## 13. Keadaan akhir

| | | Perintahnya |
|---|---|---|
| Pertanyaan | aktif **0** · selesai 57 · dihapus 2 | kepala `_selesai/OPEN-QUESTIONS.md` |
| REQ | aktif 35 · selesai 1 · mati 1 | kepala `ORACLE-REQUESTS.md` |
| Tiket | **aktif 0** · menunggu instance 8 · tertahan 1 · selesai 29 · mati 1 = **39** | `grep -h '^status:' … \| sort \| uniq -c` |
| DDL | 36 berkas · **355 pengenal** · **nol** di atas 30 byte · **nol** singkatan | `alat/periksa-penamaan.py` |

Tidak ada REQ yang menahan keputusan rancangan; empat menahan pelaksanaan — REQ-018, REQ-033, REQ-021, REQ-037.

**Yang tersisa, dan kenapa:**

- **Delapan tiket uji** — menunggu instance nyata. Rancangannya lengkap; yang kurang mesinnya.
- **Tiket `11`** — menunggu batch aplikasi. Tertahan lingkup, bukan pengetahuan.
- **Satu keluaran spec** — pemetaan migrasi kolom demi kolom, menunggu REQ-012. Pemetaan yang sah hanya datang dari `Data-Admin-DB-Table`; menuliskannya sekarang berarti mengabadikan kecocokan nama sebagai kebenaran.

---

# Bagian ketiga — empat perbaikan sebelum case ditutup

## 14. Dua kolom dinamai menurut isinya

`TANGGAL_TUTUP` memuat **hari dalam bulan**, bukan tanggal. `NILAI_TARIF` memuat **persen**, dan satuannya hanya tersirat — terbukti dari komentar kami sendiri, yang menurunkan faktor 102,2 sebagai `100 + tarif`, dan itu hanya benar bila nilainya persen.

Proyek ini menuntut kurs membawa asal-usulnya dan uang membawa mata uangnya. Tarif tidak boleh jadi kekecualian.

`TANGGAL_TUTUP` → **`HARI_TUTUP_BUKU`** · `NILAI_TARIF` → **`NILAI_TARIF_PERSEN`**. `CHECK` 1–31 sudah ada sebagai `CK_TUTUP_BUKU_1`.

## 15. Tidak ada pembangkit pengenal — di seluruh skema

Pertanyaannya tentang empat `INSERT` di `Z00`; jawabannya jauh lebih besar dari itu.

**Kedua puluh dua tabel berkunci primer `NUMBER(19)`, dan tidak satu pun menyatakan dari mana nilainya datang.** Nol `IDENTITY`, nol `SEQUENCE`, nol `DEFAULT` di 36 berkas. Hak `CREATE SEQUENCE` sudah diberikan di `00_SKEMA_DAN_AKUN.sql` — **niatnya ada, objeknya tidak pernah dibuat.**

Jadi kedua pilihan yang ditawarkan — "hilangkan ID dan biarkan identity mengisinya" atau "majukan sequence sesudah COMMIT" — keduanya mengandaikan sesuatu yang tidak ada.

**Dipilih sequence, satu per tabel, `SQ_<tabel>`.** Alasannya **persis alasan yang menetapkan batas 30 byte**: pilih bentuk yang sah di setiap versi. `IDENTITY` baru ada sejak 12.1; sequence sah jauh sebelumnya. REQ-032 diturunkan jadi verifikasi dengan janji bahwa jawabannya **tidak mengubah apa pun** — memakai `IDENTITY` akan membatalkan janji itu.

`Z00` memakai `NEXTVAL`. Dengan itu tidak ada nomor yang diandaikan **dan** tidak ada sequence yang tertinggal — keduanya sekaligus, tanpa perlu memilih di antara keduanya.

## 16. Sapuan kelas cacat — nol

Aturan 3 (*sekali kelas cacat ketahuan, ia disapu*) diterapkan **sebelum** ia terulang, bukan sesudah:

```
grep -E 'VALUES\s*\(\s*[0-9]+\s*,|\bID_[A-Z_]*\s*=\s*[0-9]+' ddl-usulan/*.sql
```

Sebelum perbaikan: **4** (keempatnya di `Z00`). Sesudah: **nol** di 36 berkas.

## 17. Cacah tiket — satu angka, dan alatnya diperbaiki

Tiga angka berbeda untuk satu hal sudah diganti satu angka berperintah di §10 dan §13.

Dan `alat/periksa-penamaan.py` ditulis ulang: versi sebelumnya **ikut mengubah berkas DDL** saat dijalankan, dan itu berbahaya diserahkan kepada orang yang belum tahu. Sekarang ia **hanya membaca**, mencacah kedelapan jenis pengenal (sebelumnya hanya tabel dan constraint — itu sebabnya angka 316 dan "terpanjang 26 byte" keduanya keliru), dan **keluar dengan kode 1** bila ada yang melanggar.

Keluarannya sekarang: **355 pengenal unik · nol di atas 30 byte · nol sisa singkatan · nol nama constraint ganda.**

## 18. Case ditutup

`KEADAAN-AKHIR.md` ditulis untuk orang yang datang tanpa sesi ini.
