# Usulan untuk tim inti — dua hal dari modul `nbfacin` (tiket NB-01)

**Untuk:** tim inti (pemilik `inti/backend/`) · **Dari:** pemilik modul `nbfacin` · **Tanggal:** 1 Oktober 2026

Dua bagian yang sekarang hidup di `modul/nbfacin/backend/services/premium/` sebenarnya milik `inti`.
Keduanya **tidak** menunggu jawaban usulan ini — modul tetap jalan dengan salinan lokal — tetapi
setelah diterima, salinan lokal dihapus.

## 1. Parser koma desimal K-027 → `inti/backend/utils`

**Masalahnya:** ADR-U-0034 meminta satu fungsi konversi teks → desimal untuk seluruh pemanggil,
karena format pemisah desimal yang berbeda-beda paling mudah merusak angka uang. Repo kini punya tiga
aturan berbeda:

| Tempat | Koma | Titik | Keduanya |
| --- | --- | --- | --- |
| `modul/treatycontractout/backend/models/tco_kurs.go` (`UraiNilaiKursTCO`) | diterima, diganti titik | diterima | ditolak |
| `modul/premiumlistlife/backend/models/polis_unggah.go` (`UangCSV`) | **ditolak** (`ErrUangBerkoma`) | diterima | ditolak |
| `modul/nbfacin/backend/services/premium/premium.go` | ditolak | diterima (`utils.ParseDecimal`) | ditolak |
| *pemuat `FACINOFFER.RATE` nbfacin — belum ada* | **wajib diterima** (K-027) | ambigu (`70.000`) | — |

Ketiganya masing-masing berdasar: kolom yang berbeda menyimpan teks dengan konvensi berbeda. Justru
karena itu konvensinya perlu dinamai per fungsi di satu tempat, bukan ditulis ulang per modul.

📌 **Ralat 1 Oktober 2026:** versi pertama usulan ini mencatat `premium.go` menerima koma. Itu sudah
tidak berlaku: data kerja Pega (berkas kasus) bertitik desimal, dan koma K-027 ternyata format **kolom
Oracle**. Karena itu `premium.go` kini menolak koma, dan pemuat `repository` NB kelak menjadi pemakai
pertama fungsi K-027 di `inti` (keputusan work owner, `KEPUTUSAN-30-09-2026.md` butir 15).

**Dasar K-027** (`modul/nbfacin/docs/00-KEPUTUSAN-WORK-OWNER.md`): `RATE` dan `PCTLIMIT` produksi
tersimpan sebagai teks berkoma (`0,0244`, `70,000 ` dengan spasi). `70,000` harus 70, bukan 70 ribu;
spasi di ujung dipangkas; nol di belakang dipertahankan.

**Usulan:** satu fungsi di `inti/backend/utils` di samping `ParseDecimal` — misalnya
`ParseDecimalKoma` — dengan konvensi yang ditulis eksplisit. Lalu ketiga pemanggil di atas ditinjau
apakah memang satu konvensi.

## 2. Uji "uang + rasio gagal kompilasi" untuk tipe `inti` → `inti/backend/uang`

**Masalahnya:** `TestUangTambahRasioGagalKompilasi` beserta `testdata/kontrol/` dan
`testdata/uangtambahrasio/` membuktikan sifat tipe **milik `inti`** (`uang.Money` tidak dapat
dijumlahkan dengan, atau diganti, `uang.Ratio` — ADR-F-0004), tetapi tinggal di modul.

**Usulan:** pindahkan kasus `uang.Money dan uang.Ratio` (berikut kedua folder `testdata`) ke
`inti/backend/uang`. Kasus tipe lokal `rasio` (`kompilasi_gagal_rasio.go`) tetap di `nbfacin`.

**Cara ujinya bekerja:** uji menjalankan `go build` atas berkas yang wajib gagal, dan atas berkas
kontrol yang wajib terkompilasi. Tanpa kontrol, pembangun yang rusak akan terbaca sebagai bukti.
Bila perintah `go` tidak ada, uji di-`Skip`.
