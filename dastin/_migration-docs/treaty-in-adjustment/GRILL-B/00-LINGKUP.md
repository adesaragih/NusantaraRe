> Modul  : Treaty In Adjustment · Ronde B · 2026-09-25
> Peran  : interogator
> Masukan: 379 XML `Treaty In Adjustment/` + 329 XML `Treaty In/` · `PENGETAHUAN.md` · `GRILL-A/` · ADR induk 0040 & 0055 · `SPEC-INVARIAN.md` · `SPEC-MODEL-DATA.md` §10.2
> Status : TERBUKA
> Sifat  : TAMBAH-SAJA

# 00 · LINGKUP RONDE B — identitas, penomoran, dan rantai addendum

## 1. Apa yang diadili ronde ini

Bagaimana sebuah versi kontrak **dikenali**: pengenalnya, nomornya, dan hubungannya dengan versi
lain. Termasuk di dalamnya nasib pengenal warisan `‹kontrak 7 karakter›/Rnn`.

**Tidak** diadili di sini: definisi materialitas (C), mekanisme pembekuan (E), migrasi data warisan
(I). Temuan yang menyentuhnya **dititipkan**, tidak diputuskan.

## 2. Sumber yang SUDAH menjawab sebagian cabang ini

Didaftar lebih dulu, sebelum satu pertanyaan pun diajukan — aturan yang lahir dari `MA-06`
(`TA-05`). Semuanya **dibaca**, bukan dirujuk dari ingatan.

| Sumber | Yang sudah diputuskannya | Sisa untuk cabang B |
|---|---|---|
| **ADR-0040** butir 1 | `OLDID` pecah menjadi dua hubungan; "addendum atas" berhenti menjadi kolom dan **menjadi hubungan versi** | konfirmasi |
| **ADR-0040** butir 2 | kunci alami = cedant + source of business + periode + `ProportionType`; ia **memperingatkan, tidak melarang** | konfirmasi |
| **ADR-0040** butir 4 | perpanjangan di tengah periode **bukan** addendum — operasi tersendiri | konfirmasi |
| **ADR-0040** konsekuensi | deteksi duplikat adalah kemampuan **BARU**; jangan dihitung gratis | konfirmasi |
| **ADR-0055** §4.1 dan Konteks | revisi-di-tempat melahirkan **versi baru yang mulai dari `DRAFT`**; keadaan disimpan **pada versi**, bukan pada kontrak | konfirmasi |
| **ADR-0055** perubahan 24 Sep | nomor versi **tidak dipakai ulang**; lompatan penomoran **jujur**, dan itu akibat syarat (a) di bawah INV-04 — bukan mekanisme baru | konfirmasi |
| **INV-04** | `NOMOR_URUT_VERSI` unik di dalam satu `KONTRAK` | **belum punya kriteria yang dapat gagal** — lihat §4 |
| **INV-25** | paling banyak **satu** versi per kontrak berada di keadaan tak-terminal | konfirmasi; menutup sendiri titipan picker NA-13 |
| **SPEC-MODEL-DATA §10.2** | `VERSI_KONTRAK` punya `ID_VERSI_KONTRAK`, `ID_KONTRAK`, `NOMOR_URUT_VERSI` | **tidak ada atribut pengenal tampilan** — lubang sesungguhnya |
| **GRL-01** | addendum adalah `VERSI_KONTRAK` milik spesifikasi induk | konfirmasi |
| **GRL-04** | penyesuaian bukan entitas tersendiri | konfirmasi |

**Kesimpulan pendaftaran ini:** dari tujuh pertanyaan yang semula tampak milik cabang B, **enam
tinggal konfirmasi**. Yang benar-benar belum diputuskan **satu**, dan itulah `B1`.

## 3. Bahan yang dibawa masuk dari ronde A

| Butir | Isi |
|---|---|
| **TDA-01** | nomor revisi dihitung dari baris yang **dipilih**, penjaga duplikatnya mati, tabrakan berakhir sebagai `UPDATE` yang melapor berhasil |
| **TDA-12** | pengenal diurai dengan **offset meleset satu**; penomorannya patah pada revisi kesepuluh, dengan dua kemungkinan akhir |
| **NA-02** | sebabnya offset tunggal, bukan dua pengurai |
| **TDA-11 / NA-13** | daftar pilihan menyatukan kontrak dan addendum tanpa saringan apa pun |
| **`UA-1`** | bentuk pengenal yang benar-benar terjadi di produksi — tiga kueri |
| **`DB-2`** | "rujukan lama sebuah addendum menunjuk versi yang disesuaikan" — pertanyaan fakta ke orang |

## 4. Bahan to-spec yang lahir sebelum pertanyaan pertama

**INV-04 belum "memperbaiki dengan sendirinya" apa pun sampai ada kriteria yang dapat gagal**
(`METODE` §5.2, §5.3). Sebuah invarian yang hanya berbunyi "unik" tidak menyatakan apa yang terjadi
saat ia dilanggar, dan justru itulah yang salah di sistem lama.

> **Bahan to-spec B-1.** *Menyimpan versi dengan nomor urut yang sudah dipakai di dalam kontrak yang
> sama **ditolak**, dan pesannya **menyebut nomor yang bentrok**.*
>
> **Berubah dari:** tabrakan pengenal masuk ke cabang `UPDATE`, menimpa baris yang sudah ada, lalu
> melaporkan berhasil (**TDA-01**, §6.2). Label: **PERUBAHAN**.
>
> Uji negatifnya jelas: simpan dua versi dengan nomor sama, harapkan penolakan beserta nomornya.

## 5. Aturan bukti yang diwarisi ronde A

Berlaku penuh di ronde ini:

* **`00-LINGKUP` ronde A §5a** — kemunculan di dalam `pyIncludedRuleXML` / `pyRuleVersionsList`
  tidak dihitung dan tidak jadi bukti; setiap hitungan menyebut badan dan salinan terpisah.
* **`TA-04`** — sapuan bernilai nol tidak dilaporkan sebelum penyapunya terbukti menemukan kasus
  positif yang diketahui. Perkakasnya `tools/tulis.py`, delapan jenis penulis.
* **`TA-05`** — setiap ADR yang masih terbuka di daftar rekonsiliasi cabang dibaca **sebelum**
  temuan ditulis. §2 di atas adalah penerapannya.
* **`TA-06`** — klaim keterlihatan menyebut hasil perkalian seluruh tingkat.

## 6. Aturan berhenti

Ronde B ditutup ketika bentuk pengenal versi diputuskan dan keenam butir konfirmasi §2 ditandai.
Yang tersisa sesudah itu bukan pertanyaan cabang B dan dititipkan ke cabangnya.

## 7. Anggaran

Sisa **32** pertanyaan untuk sepuluh cabang (B s.d. K), batas atas **40** untuk seluruh grilling.
Ronde B direncanakan **satu sampai dua** pertanyaan, sebab enam dari tujuh perkaranya sudah dijawab
sumber induk. Bila ternyata butuh lebih, tambahannya **diambil dari cabang lain** dan disebutkan —
bukan menaikkan batas atas.
