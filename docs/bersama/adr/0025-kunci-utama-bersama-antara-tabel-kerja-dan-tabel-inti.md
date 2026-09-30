---
status: accepted
tanggal: 2026-09-23
sumber: spec Claim Life dan spec penyimpanan Treaty In, keputusan work owner
---

# Tabel inti berbagi kunci utama dengan tabel kerja

Hubungan antara **tabel kerja** dan **tabel inti** di bawahnya **tidak memakai kolom penyambung**.
Keduanya memakai **kunci utama yang sama persis**.

## Bentuknya

| | |
| --- | --- |
| Tabel kerja | satu baris per objek kerja, memegang identitas |
| Tabel inti | satu baris per objek kerja, berbagi kunci utama yang sama |
| Kolom penyambung | **tidak ada** |

`[keputusan work owner]` Pola ini dipakai pada Claim Life dan pada Treaty In. Ia bukan kebetulan
dua modul, melainkan bentuk yang sama untuk persoalan yang sama.

## Kenapa

1. Hubungannya **satu lawan satu dan wajib**. Kolom penyambung terpisah hanya menambah tempat
   untuk tidak sinkron, tanpa menambah keterangan apa pun.
2. Penggabungan menjadi murah dan tidak mungkin salah pasang.
3. Menghapus pertanyaan "kunci mana yang dipakai" pada setiap tabel anak di bawahnya - seluruhnya
   memakai kunci yang sama.

## Akibat

1. Tabel inti tidak punya sequence sendiri. Identitas lahir di tabel kerja.
2. Baris tabel inti tidak dapat ada tanpa baris tabel kerja. Urutan penulisan mengikutinya.
3. Test yang menemukan kolom penyambung terpisah antara keduanya **gagal**.

## Yang catatan ini TIDAK putuskan

Hubungan tabel inti ke tabel **anak** di bawahnya. Itu kunci tamu biasa, bukan kunci bersama.
