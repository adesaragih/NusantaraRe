---
status: aktif
---

# 10: Pindahkan — nomor urut warisan diisi menurut kronologi, per batch kontrak

*Asal: `DAFTAR-PEKERJAAN.md` `P-65` bagian 2 · `GRL-17` · `ADR-0043`.*

**What to build:** Batch kontrak warisan memperoleh `NOMOR_URUT_VERSI` yang **diberikan ulang menurut
kronologi** — sumber utama tanggal komentar pembuatan, cadangan `EDMDATE`. Pengenal `/Rnn`
**dilestarikan apa adanya**.

Batch demi batch, dan **hijau di antara batch**: bentuk lama masih berdiri.

**Persyaratan:** `GRL-17` · `ADR-0042` · `ADR-0043` · `GRL-11` (turunan versi berlaku bersandar pada nomor urut)

**Tidak termasuk:** **Pencabutan pembacaan dari pengenal** — irisan 12.
**Pembersihan nomor `/Rnn` yang bentrok** — dilarang; `ADR-0042` melarang membersihkan sejarah.

**Jalur gagal:** Dua versi satu kontrak memperoleh nomor urut **sama** -> kontraknya masuk **pengecualian
migrasi bernomor**, dilaporkan per kontrak, dan **tidak** memperoleh turunan "versi berlaku" sampai
seseorang memutuskannya · Kronologi tidak dapat ditentukan dari sumber mana pun -> idem, bukan
ditebak.

**Uji:** **Negatif:** batch yang menghasilkan nomor urut kembar **harus** memunculkan pengecualian,
bukan lolos diam-diam.
**Positif:** kontrak yang kronologinya jelas memperoleh urutan yang **sama dengan urutan
`/Rnn`-nya** — bila keduanya berbeda tanpa sebab, yang salah adalah pengisian ini, bukan datanya.

**Menggantikan:** `TDA-01` dan `TDA-12` — *penjaga duplikat mati, tabrakan berakhir sebagai `UPDATE` yang
menimpa dan melapor berhasil*, dan *offset pengurai nomor revisi meleset satu*. Sistem lama
menghitung nomor dari `@substring(ID,10,12)` — memilih `/R01` ketika `R02` sudah ada menghasilkan
`R02` lagi, dan barisnya **tertimpa tanpa galat**. **Nomor yang dapat dipakai ulang bukan
urutan.**

**Blocked by:** 05

**Dasar:**
```
EVIDENCED(TreatyInRevisi_post@ekspor-2026-09, TreatyInEdmCheckDuplicate@ekspor-2026-09 - MATI)
        DECIDED(GRL-17, ADR-0042, ADR-0043)
```

- [ ] batch mengisi nomor urut dari tanggal komentar pembuatan, cadangan `EDMDATE`
- [ ] baris yang tidak dapat diurutkan menjadi **pengecualian migrasi bernomor**, dilaporkan **per kontrak**
- [ ] pengenal `/Rnn` tidak berubah pada satu baris pun
- [ ] **`NOMOR_URUT_VERSI` dikecualikan dari kriteria identik `ADR-0043`**, dan sebabnya tertulis di daftar harapan uji paritas: ia mekanisme yang **diberikan**, bukan angka yang **dipindahkan**
- [ ] `UA-21` dijalankan atau tercatat belum dapat dijalankan: berapa kontrak bernomor `/Rnn` ganda, berapa addendum tertimpa
