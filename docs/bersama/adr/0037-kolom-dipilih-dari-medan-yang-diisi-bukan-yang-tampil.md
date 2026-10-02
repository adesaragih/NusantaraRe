---
status: accepted
tanggal: 2026-09-23
sumber: spec Master Product Name Life, keputusan work owner
---

# Kolom dipilih dari medan yang benar-benar diisi, bukan dari yang sekadar tampil

`[keputusan work owner]` Daftar kolom disusun dari **medan yang di-set oleh aturan penyimpan**.
Itulah data yang benar-benar terisi.

Medan yang **hanya** muncul di layar dengan **kondisi tampil yang selalu salah** - sisa rancangan
yang tidak pernah hidup - **tidak menjadi kolom**.

## Hubungannya dengan ADR-0017

ADR-0017 menetapkan daftar medan disusun dari sapuan aturan, mencakup Activity **dan** Section.
Catatan ini menyaring hasil sapuan itu, dan keduanya sejalan:

| Langkah | Sumber | Gunanya |
| ---: | --- | --- |
| 1 | sapuan Activity **dan** Section | menemukan **calon** medan - jangan ada yang terlewat |
| 2 | penyaringan di catatan ini | menentukan calon mana yang **menjadi kolom** |

Sapuan yang luas dan penyaringan yang ketat bukan hal yang bertentangan. Melewatkan langkah 1
membuat medan hilang; melewatkan langkah 2 membuat tabel penuh kolom yang tidak pernah terisi.

## Akibat

1. Medan yang gugur di langkah 2 **dicatat beserta alasannya**, bukan dihilangkan diam-diam -
   supaya keputusan itu dapat diperiksa ulang bila kemudian ternyata keliru.
2. Medan yang terbukti diisi aturan penyimpan menjadi kolom, meskipun tidak tampil di layar mana
   pun.
3. Test membaca dan menulis kembali dokumen nyata, dan memastikan nol medan terisi yang hilang.
