# 05: Insured Interest — objek pertanggungan dan TSI

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 01 (registrasi klaim dan nomor polis) · 04 (klasifikasi lini bisnis)
**Menutup:** AC 41 · 42 · 43 *(3 AC)* — US 10–13

## Hasil & nilai pengguna

Claim Admin mendaftar objek pertanggungan beserta nilai TSI-nya, masing-masing dengan mata uang dan
kursnya sendiri, dan melihat total per mata uang beserta totalnya dalam IDR. Dengan itu dasar
perhitungan klaim jelas dan plafonnya terlihat sebelum estimasi disusun.

⚠️ **Baca namanya baik-baik.** "Interest" di modul ini adalah **objek pertanggungan dengan TSI**,
**bukan bunga finansial**. Salah baca di sini akan merusak seluruh perhitungan di hilir.

## Area codebase

Entitas objek pertanggungan · subtotal TSI per mata uang · template teks per kelas lini bisnis ·
layar Insured Interest.

## Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/CountTotalInsterest_Act.xml` | — | menghitung total TSI; **nol rate, nol jumlah hari, nol basis 360/365** — bukan bunga |
| `Activity/SetValueToClaim_Act.xml` | 9 · 10 · 11 · 12 | template teks objek pertanggungan menurut kelas lini bisnis (lihat tiket 04) |
| `Activity/ProteksiData_act.xml` | — | gerbang wajib — daftar objek pertanggungan harus terisi sebelum simpan |
| `Section/InputAcceptation_Est.xml` | — | tampilan daftar objek dan subtotalnya |

## ADR terkait

**ADR-0003** (uang non-float — nilai TSI dan kurs).

## Acceptance criteria

- [ ] ⚠️ `[terverifikasi]` "Interest" di modul ini adalah **Insured Interest (objek pertanggungan dengan TSI)**, bukan bunga finansial — nol rate, nol jumlah hari, nol basis 360/365 di seluruh rule perhitungannya *(AC 41 spec)*
- [ ] `[terverifikasi]` Daftar objek pertanggungan **wajib terisi** sebelum simpan *(AC 42 spec)*
- [ ] `[terverifikasi]` Total TSI **per mata uang** dan total dalam IDR tersedia *(AC 43 spec)*

## Perintah verifikasi

```
jalankan test "simpan tanpa objek pertanggungan -> DITOLAK"
jalankan test "dua objek beda mata uang -> subtotal per mata uang benar, total IDR benar"
jalankan test "kurs per objek dipakai, bukan kurs klaim"
cari rate/jumlah hari/basis 360/365 di perhitungan interest -> nihil
```
