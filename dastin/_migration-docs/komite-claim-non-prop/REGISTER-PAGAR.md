> Modul  : Komite Claim Non Prop · Ronde 06 · 2026-09-20
> Peran  : auditor
> Masukan: `INVENTARIS-BUKTI.md` (2026-09-20) · `KETETAPAN.md` · `PUTUSAN-01.md` §6 dan §7 · `GRILL-05/05-GERBANG.md` · `GRILL-06/01-PEMBACAAN.md` · penyuntingan penutup ronde 6
> Status : TERBUKA
> Sifat  : HIDUP

# REGISTER PAGAR

Satu baris per pagar. **Pagar tanpa baris inventaris tidak sah**, dan siapa pun yang
membacanya berhak mengabaikannya sampai rujukannya dilengkapi. Baris yang ditunjuk boleh
berupa berkas yang belum dipegang **atau data produksi yang belum diambil**
(`INVENTARIS-BUKTI.md` §2.5). Bila rujukan itu ternyata menunjuk sesuatu yang ada di repo,
pagar gugur sendiri — tanpa ronde, tanpa keputusan.

---

## 1. Pagar yang berlaku

| # | Aliran | Apa yang ditahan | Baris inventaris | Pembuka | Dipasang | Status |
|---|---|---|---|---|---|---|
| PG-03 | A-3 | menetapkan periode buku sebagai aturan tetap (mekanikanya boleh ditulis) | §2.5 baris 4 — isi badan `PROC_GENERATE_SEQUENCE_NUMBER` pada basis data berjalan belum pernah dibaca | B5-7 · DBA | ronde 1.5, C-01 | **BERLAKU** |
| PG-04 | A-4 | seluruh muatan kiriman ke Kasir, penyaring jenisnya, perangkaian JSON-nya, **dan jumlah kiriman per akseptasi** (ditambah ronde 6) | §2.5 baris 1 — nol baris `POOLDATA.DIRECTTOKASIR_LOG` | satu `SELECT` baca-saja atas log kiriman | ronde 1.5, `PUTUSAN-01.md` §6 | **BERLAKU — diperiksa ronde 6, tidak menyusut** |
| PG-05 | A-5 | penulisan baris akseptasi, rincian layer, cakupan pembalikan | §2.1 — `PEGA_JSON_OS_AKSEP_KLAIM`, `PEGA_JSON_OS_AKSEP_SUBJECTIVITY`, `XOL2_AKSEP_KLAIM` tidak dipegang | source ketiga prosedur itu | ronde 1.5, `PUTUSAN-01.md` §6 | **BERLAKU** |
| PG-06 | A-6 | **isi** dokumen dan surat; bukan kapan terbit dan siapa penerimanya | §2.3 — `PostEmailKomiteCNP` dipanggil, berkasnya tidak ada | B5-4 · admin Pega | ronde 1.5, `PUTUSAN-01.md` §6 | **BERLAKU (ringan)** |

---

## 2. Pagar yang gugur

| # | Aliran | Bunyi pagar | Mengapa gugur | Tanggal |
|---|---|---|---|---|
| PG-01 | A-1b | "A-1b beku, dibuka oleh ekspor dua rule" | Keempat rule pemagarnya sudah ada di repo sebelum pagarnya dipasang: `CreateChildKomiteCNP_Act`, `CreateChildKomiteCloseNP_Act`, `ProteksiSendKomiteCNP_Act`, `FilterEmailKomiteWithLimit` — `INVENTARIS-BUKTI.md` §2.4. Pagar gugur sendiri, tanpa keputusan | gugur 2026-09-20, ronde 4 |
| PG-07 | A-1a | "menunggu parameter metode dua activity inti (G-01/Q-9)" | `KomiteRouter.xml` seluruhnya `Property-Set`, yang memang terekspor; pagar ini tidak pernah menyentuh apa yang ditahannya | gugur 2026-09-20, ronde 5, `GRILL-05/05-GERBANG.md` |

Kedua pagar yang gugur **tidak dipasang kembali**. Keduanya adalah kejadian keempat dan
kelima dari pola yang sama: pernyataan tentang ketiadaan bukti dibuat tanpa memeriksa
inventaris, lalu diwarisi berkas berikutnya sebagai premis
(`INVENTARIS-BUKTI.md` §3).

---

## 3. Pagar yang diteruskan dengan dasar berbeda

| # | Aliran | Bunyi lama | Dasar lama | Dasar baru | Status |
|---|---|---|---|---|---|
| PG-02 | A-1b | E-3: aksi pengalihan tidak diekspos **sampai A-1b cair** | A-1b beku | — | **DIPINDAH ronde 6 ke bagian 4 — bukan pagar**; baris ini tidak dihapus |

A-1b sudah cair sejak ronde 4, sehingga dasar lama PG-02 habis. `H-4` meneruskannya
("pengalihan sesaat tetap tidak diekspos") tanpa menyebut dasar penggantinya. Aturan berkas
hidup ini berbunyi: **pagar yang diteruskan dengan dasar baru wajib menyatakan dasar
barunya.** Sampai dinyatakan, PG-02 dicatat sebagai pagar yang bertahan atas kekuatan
ketetapan saja, bukan atas kekuatan bukti yang belum dipegang — dan karena itu ia **bukan**
pagar dalam arti berkas ini, melainkan keputusan rancangan.

Perlakuan yang disarankan, bukan diputuskan di sini: pindahkan `H-4` dari register pagar ke
`REGISTER-DEVIASI.md`, karena tidak mengeksposnya adalah pilihan rancangan yang sah dan
tidak menunggu bukti apa pun.

---

## 4. Yang bukan pagar, meski sering disebut begitu

| Hal | Mengapa bukan pagar |
|---|---|
| Sisa `G-01` (tiga langkah tanpa pengikatan parameter) | pengambilan bukti biasa, `03-LUBANG.md` B5-5; tidak menahan aliran mana pun |
| `Q-2`, `Q-3`, `Q-6`, `Q-7` | seluruhnya ditutup ronde 5 dengan membaca berkas yang sudah dipegang |
| `N-13` (pintu masuk lain `KomitePostAdjustmentCWP`) | penyelidikan dihentikan di `06-PUTUSAN.md` §2; bukan ditahan, melainkan dinyatakan tidak berpengaruh |
| `PG-02` — pengalihan sesaat tidak diekspos *(dipindah ke sini ronde 6)* | `H-4` kini punya bunyi penuh (`KETETAPAN.md` bagian 4), dan ia **ketetapan**, bukan pagar: ia tidak menunggu bukti apa pun. Dasarnya bukan lagi "A-1b beku" melainkan bahwa kebutuhan nyatanya sudah terpenuhi oleh delegasi tetap (`H-3`). Saran ronde 5 untuk memindahkannya ke `REGISTER-DEVIASI.md` **ditolak**: bagian 3 register deviasi sudah benar menyatakan `H-4` bukan deviasi, sebab sistem lama pun tidak mengeksposnya. Ketetapan tinggal di `KETETAPAN.md`; tidak setiap keputusan butuh register |

---

## 5. Aturan yang berlaku bagi register ini

1. Setiap pagar baru **wajib** menyebut nomor bagian dan nama objek pada
   `INVENTARIS-BUKTI.md`, bukan kalimat "belum ada".
2. Pagar berlaku bagi **alasan**, bukan hanya bagi kesimpulan. Argumen yang bersandar pada
   isi aliran beku tidak sah meski kesimpulannya berada di aliran terbuka.
3. Pagar yang rujukannya menunjuk sesuatu yang ada di repo gugur seketika, dan
   kegugurannya dicatat di bagian 2 — bukan dihapus.
4. Register ini dimutakhirkan setiap kali satu berkas bukti masuk atau satu kueri pulang.
   Register yang tidak dimutakhirkan berhenti menjadi dasar yang sah bagi pagar mana pun.

---

## 6. Pemutakhiran ronde 6 — 2026-09-20

**Bagian 3 kini kosong.** Tidak ada lagi pagar yang bertahan tanpa dasar yang dinyatakan.
Baris `PG-02` tidak dihapus dari bagian 3; ia diberi status dan isinya pindah ke bagian 4.

**`PG-04` diperiksa dan tidak berubah.** Pembacaan ronde 6 mengandaikan `PG-04` akan menyusut
menjadi pagar atas isi muatan saja bila `G-03` runtuh. **`G-03` tidak runtuh**
(`GRILL-06/01-PEMBACAAN.md` P6-1), sehingga penyusutan itu tidak terjadi.

Satu fakta struktural yang terbaca ronde 6 justru memperluas apa yang pagar ini tahan:
`Connect-REST` berada di langkah **10.7**, di dalam loop yang sama dengan langkah 10.3, dan
pra-syaratnya hanya `IsPEGAPROD` — bukan `.TreatyName=="UR"`. Penulisan log pada langkah
10.9 tidak berpra-syarat sama sekali. Berapa kiriman yang benar-benar terjadi per nomor
akseptasi karena itu ikut berada di balik pagar ini, dan ia dijawab oleh `SELECT` yang sama.
Fakta strukturnya dicatat; **artinya tidak disimpulkan**.

`PG-03`, `PG-05`, `PG-06` tidak berubah. `PG-01` dan `PG-07` tetap tercatat gugur.
