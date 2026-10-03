# 09: Deductible — satu rumus, satu makna

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 05 (Insured Interest / TSI) · 08 (baris adjustment)
**Menutup:** AC 50 · 51 · 52 · 53 *(4 AC)* — US 33–35

## Hasil & nilai pengguna

Deductible dihitung dengan **satu** rumus yang disepakati, atas basis yang dipilih eksplisit, dan
dikurangkan pada titik yang sama setiap kali. Claim Admin berhenti menebak versi mana yang sedang
berlaku.

⚠️ Korpus memuat **dua rumus deductible yang bertentangan**. Tiket ini menutup pertentangan itu
dengan memilih versi 2024.

## Area codebase

Perhitungan deductible · pemilihan basis · urutan pengurangan terhadap share ceding.

## Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/CountValueADJTreaty_Act.xml` | 6 | penyelarasan mata uang sebelum deductible dari TSI |
| `Activity/CountValueADJTreaty_Act.xml` | 4 · 5 · 7 · 8 · 9 | tiga jenis basis risiko individual dan cabang per tipe |

⚠️ `[terverifikasi]` Properti bernama sama berarti **berbeda** di dua rule lama — sumber langsung
pertentangan rumus.

## ADR terkait

**ADR-0003** (uang non-float).

## Acceptance criteria

- [ ] `[keputusan work owner]` Deductible = **MAX(persentase × basis, nilai flat)** *(AC 50 spec)*
- [ ] `[terverifikasi]` Basis dipilih antara **TSI** dan **nilai klaim** *(AC 51 spec)*
- [ ] `[keputusan work owner]` Rumus yang berlaku adalah versi **2024** — deductible dikurangkan **sesudah** share ceding diterapkan. Versi 2022 **ditinggalkan** *(AC 52 spec)*
- [ ] ⚠️ Properti bernama sama diberi **satu makna saja**. **Alasan menyimpang:** di dua rule lama nama yang sama berarti dua hal berbeda, sehingga pembaca tidak dapat tahu rumus mana yang sedang berjalan *(AC 53 spec)*

## Perintah verifikasi

```
jalankan test "persentase x basis > flat -> deductible = hasil persentase"
jalankan test "flat > persentase x basis -> deductible = flat"
jalankan test "basis TSI dan basis nilai klaim menghasilkan angka berbeda sesuai pilihan"
jalankan test "deductible dikurangkan SESUDAH share ceding (versi 2024)"
cari implementasi rumus deductible kedua                  -> nihil
```
