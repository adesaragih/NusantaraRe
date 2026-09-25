---
status: tertahan
---

# 64: Tabel acuan pembagian kapasitas punya pemilik yang bertanggung jawab atas isinya

*Asal: `DAFTAR-PEKERJAAN.md` `P-49` · `ADR-0050` · eskalasi manajemen butir 7.*

**What to build:** Tabel acuan pembagian kapasitas membawa **pemiliknya**, dan setiap perubahan isinya
berjejak beserta siapa dan kapan.

**Persyaratan:** `ADR-0050` (sumber luar yang tidak dipercaya) · `ADR-0045`

**Tidak termasuk:** **Isi tabelnya** — ia data, bukan bentuk; pemiliknya yang mengisinya.

**Jalur gagal:** Perubahan isi tabel acuan tanpa jejak -> mustahil · Tabel tanpa pemilik -> tiket ini tidak dapat dinyatakan selesai.

**Uji:** **Negatif:** ubah isi tanpa pelaku. **Positif:** perubahan oleh pemilik yang sah tercatat lengkap.

**Menggantikan:** **`PROPORTIONALARRG`** — master susunan retro sistem lama: **seluruh kolom
`VARCHAR2(1000)`, tanpa kunci utama maupun indeks**, dan tanpa pemilik yang tercatat di mana pun.

**Blocked by:** `15` · **penetapan pemilik**

**Dasar:**
```
EVIDENCED(PROPORTIONALARRG DDL@ekspor-2026-09 - seluruh kolom VARCHAR2(1000), nol PK, nol indeks)
        DECIDED(ADR-0050, ADR-0045)
```

- [ ] tabel acuan kapasitas membawa pemiliknya
- [ ] perubahan isinya berjejak
- [ ] **PENGHALANG:** penetapan pemilik · **siapa menjawab:** manajemen, eskalasi butir 7 · **yang berubah:** siapa yang bertanggung jawab, bukan bentuk tabelnya
