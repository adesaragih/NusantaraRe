# ADR-0056 — Tidak ada stored procedure di sistem baru

**Status:** diterima, 24 September 2026 — arahan pemilik proses (K-4)
**Berlaku untuk:** seluruh skema `TREATY_MASUK`, sejak gelombang pertama
**Menggantikan:** tidak ada · **Bersandar pada:** ADR-0023, ADR-0035, ADR-0041, ADR-0051

## Konteks

Sistem lama menyimpan Treaty In lewat lima *stored procedure* di `POOLDATA`. Kelimanya dibedah
baris demi baris; hasilnya `2-to-spec/PEMETAAN-PROCEDURE.md`. Yang terbaca **bukan** pilihan
arsitektur melainkan satu cetakan yang disalin empat kali — dan cacatnya ikut tersalin empat kali.

Tiga bentuk cacat yang muncul berulang, dan ketiganya **tidak mungkin ada** kalau logikanya berada
di lapisan aplikasi yang dapat diuji:

1. **Pekerjaan yang tidak dilakukan, dilaporkan berhasil.** `PEGA_TREATY_IN` dipanggil dari layar
   addendum menghasilkan `UPDATE … WHERE ID = 'XXXXXXX/Rnn'` yang mengenai **nol baris**; Oracle
   tidak menimbulkan galat, `COMMIT` berhasil, `StsSave := 1`, dan pesannya berbunyi *"Data Sudah
   Disimpan Dengan ID : …"*. Kedua procedure rinci **tidak punya cabang `ELSE`** sama sekali, dan
   tetap menyetel `StsSave := 1` — jalur perbarui **tidak melakukan apa pun**, dengan `IDPegaOut`
   keluar `NULL`.
2. **`EXCEPTION WHEN OTHERS` tanpa `RAISE`, di tiga tingkat, di kelima procedure.** Setiap galat
   ditelan. Satu di antaranya (`PEGA_M_TREATY_IN_DETAIL`) tidak pernah memberi nilai kepada
   `ErrMsg` sama sekali, walaupun mendeklarasikannya `OUT`.
3. **Penyulihan teks buta atas dokumen JSON** — `replace(DataPega,'UnknownId',id)` mengganti setiap
   kemunculan kata itu di mana pun di dalam dokumen, tanpa dibatasi pada ruas pengenal.

Ketiganya bertahan bertahun-tahun karena **tidak ada satu pun yang dapat dijalankan sebagai uji**.
Procedure tidak punya seam: ia menulis, meng-`COMMIT`, dan menelan galatnya sendiri di dalam satu
benda yang hanya dapat diperiksa dengan menjalankannya terhadap basis data sungguhan.

## Keputusan

> **Aplikasi Golang tidak memanggil stored procedure, dan skema `TREATY_MASUK` tidak memuat
> satu pun `CREATE PROCEDURE`, `CREATE FUNCTION`, maupun trigger yang membawa aturan bisnis.**

Oracle tetap dipakai sebagai basis data. Yang dipindahkan adalah **tempat aturannya tinggal**:

| Yang boleh tinggal di Oracle | Yang tidak |
|---|---|
| tipe, `NOT NULL`, `CHECK` domain, kunci utama dan asing, keunikan | keputusan cabang, penerbitan pengenal, pengambilan nilai induk, pemilihan jalur sisip-versus-perbarui |
| *materialized view* untuk invarian lintas baris (ADR sebelumnya) | pesan untuk pengguna, penanganan galat, batas transaksi |

**Batas transaksi dimiliki pemanggil.** Tidak ada `COMMIT` di dalam objek basis data.

**Kegagalan adalah kegagalan** (ADR-0035): perintah yang tidak mengubah baris mana pun ketika ia
seharusnya mengubah satu baris **adalah galat**, bukan keberhasilan yang sunyi. Lapisan penyimpan
memeriksa jumlah baris terpengaruh dan menolak bila bukan yang diharapkan.

## Alternatif yang ditimbang

| Pilihan | Kenapa tidak |
|---|---|
| **Pertahankan procedure, perbaiki cacatnya** | memperbaiki ketiga cacat menuntut uji, dan uji menuntut seam. Membangun seam di sekitar PL/SQL berarti membangun kerangka uji basis data untuk satu modul — ongkos yang sama dengan memindahkannya, tanpa memperoleh keterujian di tempat lain |
| **Pindahkan sekarang, biarkan procedure lama tetap ada untuk Pega** | tidak ditolak, dan **memang itu yang terjadi selama masa berdampingan**. ADR-0041 sudah mengaturnya: hak tulis Pega dicabut per gelombang. ADR ini tentang **skema baru**, bukan tentang mencabut procedure lama lebih cepat |
| **Pakai trigger untuk integritas** | trigger dapat berhenti bekerja tanpa gagal, dan penegakan semacam itu menuntut pemantau kebasian tersendiri. Constraint tidak dapat berhenti diam-diam |

## Konsekuensi

- **Kelima cacat §PEMETAAN-PROCEDURE lenyap karena bentuknya, bukan karena ditambal.** Tidak ada
  tempat bagi "berhasil tanpa melakukan apa pun" ketika jumlah baris terpengaruh diperiksa dan
  cabangnya ditulis sebagai kode yang punya uji.
- **Dua proyeksi datar tidak dibawa.** `TREATYINDETAIL` dan `TREATYINDETAILEDM` adalah turunan dari
  dokumen JSON; di skema baru bentuk kanoniknya relasional, sehingga proyeksinya tidak perlu
  ditulis ulang oleh siapa pun. Ketimpangan tujuh kolom di antara keduanya ikut lenyap.
- **Konsumen hilir yang membaca kedua tabel datar itu harus dilayani** — lewat *view* atas bentuk
  kanonik baru, bukan lewat procedure yang menyalin (ADR-0041, ADR-0051). **View itu tidak
  mengembalikan ketimpangan tujuh kolom**: ia memaparkan besaran yang sama untuk kontrak maupun
  addendum.
- **Penerbitan pengenal pindah ke aplikasi**, dan bersamanya kesempatan memperbaiki `lpad(...,6,'0')`
  yang membuat skema pengenal lama pecah. Bentuk penggantinya bukan urusan ADR ini.
- **Satu hal yang TIDAK diputuskan di sini:** apa yang terjadi pada kelima procedure lama. Mereka
  tetap melayani Pega sampai gelombangnya pindah, dan pencabutannya diatur ADR-0041 langkah 1.
- **Satu temuan diangkat keluar dari lingkup Treaty In.** `GET_TOKEN_STORAGE` menerbitkan token
  akses dari `STANDARD_HASH('ASMAPP' || timestamp,'MD5')` — **dapat ditebak sepenuhnya, tanpa unsur
  acak**. Itu bukan utang migrasi; ia terbuka sekarang, di sistem yang berjalan. Ia naik ke
  `DAFTAR-ESKALASI-MANAJEMEN.md` dan **tidak** ditambal oleh sesi ini.

## Syarat pembalikan

Keputusan ini dibalik bila salah satu terbukti:

1. ada konsumen yang **memanggil** procedure ini secara langsung — bukan membaca tabelnya — dan
   tidak dapat diubah; maka yang dibutuhkan **pembungkus**, dan bentuknya diputuskan tersendiri;
2. ada kewajiban kepatuhan yang menuntut aturan tertentu ditegakkan di dalam basis data, di luar apa
   yang dapat dinyatakan sebagai constraint.

**Tidak satu pun ditandai terverifikasi.** Keduanya belum diperiksa, dan pemeriksaannya menunggu
daftar konsumen hilir yang sampai sekarang tidak ada (ADR-0051).
