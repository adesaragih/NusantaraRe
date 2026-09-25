> Modul  : Treaty In Adjustment · peta sumber modul induk
> Dibuat : 2026-09-26
> Guna   : dicari **sebelum** setiap butir GRILL dan sebelum menggolongkan apa pun sebagai REKOMENDASI
> Dasar  : pelajaran `MA-06` dan `TA-05` — sumber induk dibaca lebih dulu, bukan diputuskan ulang

# Peta sumber modul induk `_migration-docs/treaty-in`

## 0. Urutan wewenang — tidak berubah

| Jenis pernyataan | Yang menang |
|---|---|
| **fakta tentang sistem lama** | ekspor XML dan `PENGETAHUAN.md` Adjustment — **di atas** dokumen induk (`METODE` §3.7) |
| **rancangan** | ADR dan spesifikasi induk — **kecuali** ada usulan `REV` yang sudah diputuskan |

Fakta di dokumen induk **tidak diwarisi begitu saja**: kutip dulu, lalu periksa ke ekspor.

**Berkas induk tidak disunting.** Usulan perubahan masuk `USULAN-REVISI-ADR.md`; perubahan artefak
induk ditunjukkan sebagai diff lebih dulu.

## 1. Tanda ⚠ — ditulis sebelum `PENGETAHUAN.md` Adjustment

`PENGETAHUAN.md` Adjustment mulai ditulis **23 September 2026** dan dikoreksi terus sampai 26
September. Berkas induk bertanda ⚠ ditulis **sebelum atau bersamaan** dengan itu dan **belum**
memperhitungkan koreksinya — terutama pencabutan G2 (§2.1: `OLDDATA` potret, bukan rujukan), TDA-07,
dan hitungan yang ikut menghitung salinan terbungkus. Berkas bertanda ⚠ **dikutip, lalu diperiksa**.

## 2. Akar folder

| Berkas | Isinya, satu kalimat | Status | Tanggal | |
|---|---|---|---|---|
| `CONTEXT.md` | glosarium dan batas modul induk; §1.1 dan §4 sudah dikoreksi sesi ini | hidup | 24 Sep | |
| `PENGETAHUAN.md` | AS-IS modul Treaty In dari ekspor Pega | hidup | 23 Sep | ⚠ |
| `SPEC-MODEL-DATA.md` | model data to-spec induk | **BENTUK AWAL — langkah 0**; langkah 8–10 belum jalan | 24 Sep | |
| `SPEC-INVARIAN.md` | daftar invarian `INV-nn` beserta bentuk penegakannya | hidup | 24 Sep | |
| `UJI-NEGATIF-INVARIAN.md` | uji yang harus gagal untuk tiap invarian | hidup | 23 Sep | |
| `PETA-TELUSUR-JSON.md` | telusur 667 jalur `JSONDATA` ke atribut | hidup | 23 Sep | ⚠ |
| `INVENTARIS-STRUKTUR-DATA.md` | inventaris struktur data sistem lama | hidup | 23 Sep | ⚠ |
| `ERD-TREATY-MASUK.md` | ERD naratif; **§2 teks mengikat, §5 gambar alat bantu** | hidup | 23 Sep | ⚠ |
| `RINGKASAN-GRILLING.md` | ringkasan grilling modul induk | arsip | 23 Sep | ⚠ |
| `KEPUTUSAN-TANPA-VERIFIKASI.md` | keputusan yang diambil tanpa uji data, beserta syarat pembalikannya | hidup | 23 Sep | ⚠ |
| `TETAPAN-DI-KODE.md` | nilai yang ditanam di kode sistem lama | hidup | 24 Sep | |
| `DAFTAR-ESKALASI-MANAJEMEN.md` | butir untuk manajemen; butir 1 dikoreksi sesi ini | hidup | 24 Sep | |
| `SEAM-ADJUSTMENT.md` | sambungan ke modul Adjustment | **sudah diberi banner koreksi** sesi ini | 24 Sep | ⚠ |
| `PERMINTAAN-DBA-1-UJI-A-SAMPAI-H.sql` | kueri profil data yang sudah dikirim ke DBA | hidup | — | |

## 3. `docs/adr/` — 22 ADR, 0034 s.d. 0055

Seluruhnya 23 September kecuali **0055** (24 September). Yang paling sering dipakai sesi ini:

| ADR | Isinya | Dipakai di |
|---|---|---|
| **0034** | arsip JSON sistem lama tidak punya jalur baca | F-bentuk-simpan → KONFIRMASI |
| **0036** | beku saat disetujui, hitung saat dibaca | GRL-02, `REV-1` |
| **0037** | masukan versus turunan; *"'Boleh diubah' = turunan dari keadaan"* | GRL-11 |
| **0038** | aturan yang bisa berubah disimpan sebagai data | H-sumber-editabilitas → KONFIRMASI |
| **0040** | identitas kontrak dan addendum; kunci alami **memperingatkan** (butir 2); **lapisan beku melarang** (butir 3) | GRL-05, GRL-09, `R-D2` dicabut |
| **0041** | satu fakta satu penulis | GRL-11 |
| **0042** | sejarah pindah apa adanya | GRL-06, GRL-09 |
| **0044** | wewenang: keadaan dan peran terpisah; *"Nama orang tidak pernah muncul di dalam aturan"* | G → KONFIRMASI |
| **0045** | jejak perubahan sebagai fakta mesin | D-jejak → KONFIRMASI |
| **0048** | seam Treaty In Adjustment; butir 2 = SELECT, tidak menyalin | GRL-03, `REV-1` |
| **0049** | jenis addendum satu sumbu | GRL-06 |
| **0051** | anggap ada konsumen hilir | GRL-09, `DB-11` |
| **0052** | satu rantai persetujuan tanpa pengecualian | GRL-07, `REV-2` |
| **0053** | mata uang daftar; `CurrencyRelation` pasif | E-mata-uang → KONFIRMASI |
| **0054** | keadaan warisan tak terpetakan | `UA-12`, cabang I |
| **0055** | daftar keadaan dan perpindahan; perubahan 24 Sep menambah `DIBATALKAN` | GRL-08, `REV-3` |

Sisanya — 0035, 0039, 0043, 0046, 0047, 0050 — belum terpakai di ronde A atau B, dan **belum
diperiksa butir demi butir**. Itu dicatat supaya tidak dibaca sebagai "sudah diperiksa dan tidak
relevan".

## 4. `4-erd-dan-tabel-datar/`

| Berkas | Isinya | Status | Tanggal | |
|---|---|---|---|---|
| `STRUKTUR-DATA.md` | **13 entitas tulang punggung + anak + entitas luar**; sumber `ID_VERSI_KONTRAK_DASAR`, `CEDANT [luar]`, `DOKUMEN_KONTRAK` | hidup — **paling sering menjawab** | 24 Sep | |
| `ERD.md` | relasi antarentitas; §2.6 menyatakan `NILAI_SELISIH` tidak punya relasi ke "versi lama" | hidup | 23 Sep | ⚠ |
| `TABEL-DATAR.md` | tabel datar turunan dan asal kolomnya | hidup | 23 Sep | ⚠ |
| `ISI-FOLDER.md` | penjelasan isi folder dan konvensi penamaan | hidup | 23 Sep | ⚠ |
| `PETA-NAMA-TABEL-TREATYIN.md` | padanan nama tabel lama dan baru | hidup | 23 Sep | ⚠ |
| `struktur-treatyin-lama.md` | struktur sistem lama | hidup | 23 Sep | ⚠ |
| `KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md` | keputusan sambungan ke Adjustment | **sudah diberi banner koreksi** sesi ini | 24 Sep | ⚠ |
| `TEMUAN-ADJUSTMENT-DITUNDA.md` | temuan Adjustment yang ditunda sesi induk | sebagian sudah diadili `GRILL-A/02-SIDANG` | 23 Sep | ⚠ |
| `AUDIT-PENYEBUT-POHON.md` · `TITIK-BUTA-POHON.md` | audit penyisir pohon clipboard dan titik butanya | hidup | 24 Sep | |
| `BENTUK-PENULIS-PROPERTI.md` · `SAPUAN-DAN-NAMA-TAGNYA.md` | bentuk penulis properti dan nama tag sapuan | hidup — **bersinggungan langsung dengan `TA-04`** | 24 Sep | |
| `JENIS-ATURAN-TAK-TEREKSPOR.md` | jenis aturan yang tidak ikut ter-ekspor | hidup — **bersinggungan dengan `EXP-1`** | 24 Sep | |
| `PRA-PEMISAHAN-ANGKA-DAN-TABRAKAN.md` · `PERIKSA-BUTIR-CONTOH.md` | pemisahan angka dan pemeriksaan butir contoh | hidup | 24 Sep | |

> **Dua berkas yang harus saya baca lebih awal daripada ini, dan belum:**
> `BENTUK-PENULIS-PROPERTI.md` dan `SAPUAN-DAN-NAMA-TAGNYA.md`. Keduanya bertanggal 24 September —
> hari yang sama dengan `MA-05`, ketika saya menyapu penulis dengan nama tag yang salah. Sangat
> mungkin jawabannya sudah ada di sana. **Dibaca sebelum butir GRILL berikutnya**, dan hasilnya
> dilaporkan apa adanya termasuk bila ternyata `TA-04` hanya menemukan ulang isinya.

## 5. `5-tiket/`

| Berkas | Isinya | Status | Tanggal |
|---|---|---|---|
| `DAFTAR-PEKERJAAN.md` | daftar hal yang harus ada di sistem baru beserta asal-usulnya; **bukan tiket, bukan taksiran** | hidup | 24 Sep |
| `BENTUK-TIKET.md` | bentuk baku sebuah tiket | hidup | 24 Sep |
| `KEPUTUSAN-PEMBAGIAN-TIKET.md` | dasar pembagian tiket | hidup | 24 Sep |
| `LUBANG-SPESIFIKASI.md` | tempat yang menuntut kembali ke ekspor; **dilaporkan, tidak ditambal** | hidup | 24 Sep |

## 6. Cara memakai peta ini

1. **Sebelum setiap butir GRILL**, dan sebelum menggolongkan apa pun sebagai REKOMENDASI: cari
   jawabannya lewat tabel di atas.
2. **Sudah dijawab** → jadikan **KONFIRMASI**, dengan **kutipan dan barisnya**, tanpa pertanyaan.
3. **Dijawab sebagian** → tanyakan **sisanya saja**, seperti pada cabang B.
4. **Berkas bertanda ⚠** → kutip, lalu periksa ke ekspor sebelum dipakai sebagai fakta.
