# 07: Estimasi klaim per treaty

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 05 (Insured Interest / TSI) · 06 (loss allocation dan spreading)
**Menutup:** AC 34 · 35 · 36 · 37 · 38 · 39 · 40 *(7 AC)* — US 14–19

## Hasil & nilai pengguna

Claim Admin melihat estimasi terbagi per treaty sesuai penempatan risiko, masing-masing dengan mata
uang dan kursnya sendiri. Ia diperingatkan bila total estimasi melampaui TSI atau plafon cash call,
dan menghapus baris membuat seluruh subtotal ikut menyesuaikan — jadi angkanya tidak pernah
menggantung.

## Area codebase

Entitas baris estimasi · subtotal estimasi per mata uang · validasi tanggal · peringatan plafon.

## Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/AddEstimation_Act.xml` | — | baris estimasi dibentuk di clipboard; lookup kurs |
| `Activity/CountEstimation_Act.xml` | — | bercabang eksplisit satu-mata-uang versus multi-mata-uang |
| `Activity/CurencyEstimation_Act.xml` | — | mata uang ditetapkan **per baris** |
| `Activity/DeleteEstimation_Act.xml` | — | penghapusan baris dan penyesuaian subtotal |
| `Activity/CheckDateDOL_Act.xml` | — | batas bawah tanggal estimasi = Date of Loss |
| `Activity/AddAdjustment_Act.xml` | 5 | ⚠️ menulis **potret** nilai estimasi ke baris adjustment; **nol rule lain memperbaruinya** |
| `Activity/CountValueADJTreaty_Act.xml` | 16 · 17 | satu-satunya pembaca potret itu |

## ADR terkait

**ADR-0003** (uang non-float).

## Acceptance criteria

- [ ] `[terverifikasi]` Baris estimasi dibentuk **per treaty**, bersumber dari hasil loss allocation *(AC 34 spec)*
- [ ] `[terverifikasi]` Tanggal estimasi wajib **antara Date of Loss dan hari ini**, inklusif di kedua ujung, zona `Asia/Jakarta` *(AC 35 spec)*
- [ ] Total estimasi melebihi TSI → peringatan *(AC 36 spec)*
- [ ] Total estimasi melebihi plafon cash call → peringatan *(AC 37 spec)*
- [ ] Menghapus baris estimasi menyesuaikan seluruh subtotal *(AC 38 spec)*
- [ ] ⚠️ Dua nama kolom Pega **tidak dibawa apa adanya** ke nama kolom baru. **Alasan menyimpang:** satu kolom bernama persen berisi **nilai uang**, satu kolom bernama jenis kerugian berisi **identitas treaty** — nama yang berbohong membuat kode tidak terbaca *(AC 39 spec)*
- [ ] ⚠️ **Pagar nilai mengikat estimasi TERKINI, bukan potret.** Potret tetap disimpan untuk jejak audit tetapi **tidak** dipakai sebagai pembanding. **Alasan menyimpang:** di Pega nilainya potret yang ditulis sekali saat baris adjustment dibuat dan **nol rule lain memperbaruinya**, sehingga pagar dapat mengikat angka yang sudah basi *(AC 40 spec)*

## Perintah verifikasi

```
jalankan test "estimasi bertanggal sebelum Date of Loss -> DITOLAK"
jalankan test "estimasi bertanggal besok -> DITOLAK"
jalankan test "estimasi bertanggal tepat Date of Loss -> diterima"
jalankan test "total estimasi > TSI -> peringatan tampil"
jalankan test "total estimasi > plafon cash call -> peringatan tampil"
jalankan test "hapus satu baris -> seluruh subtotal menyesuaikan"
jalankan test "estimasi berubah sesudah baris adjustment dibuat -> pagar memakai nilai TERKINI"
jalankan test "potret estimasi tetap tersimpan di jejak audit"
```
