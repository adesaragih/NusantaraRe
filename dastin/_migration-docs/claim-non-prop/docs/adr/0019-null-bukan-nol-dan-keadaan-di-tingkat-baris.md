---
status: accepted
label: DECIDED
---

# NULL bukan nol, dan keadaan perhitungan disimpan di tingkat baris

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.

Tiga lapis, disusun dari yang termurah. Lapis berikutnya hanya ditambahkan di tempat lapis sebelumnya tidak cukup.

## Lapis 1 — `NULL` berarti belum dihitung, nol berarti dihitung dan hasilnya nol

Keduanya **tidak boleh pernah disamakan**, di mana pun, oleh kode mana pun. Tidak ada `NVL(x,0)` yang diam-diam menghapus perbedaan itu, tidak ada kolom nilai ber-`DEFAULT 0`, dan tidak ada layar yang menampilkan kosong sebagai `0`.

Biayanya nol — hanya disiplin — dan ia menangkap sebagian besar manfaat dari seluruh keputusan ini.

## Lapis 2 — keadaan di tingkat baris, bukan tingkat kolom

Setiap baris hasil perhitungan membawa tiga kolom: `dihitung_pada`, `status_perhitungan`, dan `versi_aturan` yang dipakai. Satu perhitungan gagal maka statusnya terbaca di baris itu, dan kolom yang gagal bernilai `NULL`.

Tiga kolom, bukan ratusan.

## Lapis 3 — penanda per-nilai, hanya di satu tempat

Hanya untuk **nilai uang hasil konversi mata uang**, karena di sana satu mata uang dapat gagal sementara mata uang lain berhasil **di dalam baris yang sama** — satu-satunya kasus yang tidak tertangkap lapis 2.

## Consequences

Usulan pertama saya — satu kolom penanda untuk setiap nilai turunan — ditolak, dan alasan penolakannya benar: ratusan kolom tambahan akan diisi asal saat implementasi, lalu tidak ada yang percaya isinya. **Penanda yang tidak dipercaya lebih buruk daripada tidak ada penanda.**

Urutannya juga dibalik dari usulan saya. Bukan *"kalau terlalu mahal, batasi ke uang"*, melainkan *"mulai dari yang murah dan menyeluruh, lalu perdalam hanya di tempat yang paling mahal kalau salah"*.

**Masalah yang diselesaikan** — tiga mekanisme di sistem lama, semuanya berakar pada tidak adanya cara membedakan nol dari gagal:

| Mekanisme | Kegagalan | Terlihat sebagai |
|---|---|---|
| `GETCURRENCYSTANDARD` | kurs tidak ada | kurs bernilai 1 (`FINDING-006`) |
| Database link HRD | tautan putus | daftar pengguna kosong (`ADR-0016`) |
| Precondition tidak cocok | kondisi gagal | langkah dilewati tanpa pesan (`FINDING-002` bagian 2) |

Ketiganya menghasilkan jawaban yang bentuknya benar dan isinya salah.
