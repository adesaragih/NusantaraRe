---
status: aktif
golongan: pelestarian
---

# 22: Retensi cedant dicatat per kelompok treaty per mata uang

*Asal: `DAFTAR-PEKERJAAN.md` `P-07` · `SPEC-MODEL-DATA.md` §10.11 · `SPEC-INVARIAN.md` `INV-08`.*

**What to build:** **PK** mencatat berapa yang **ditahan cedant sendiri**, dirinci **per kelompok
treaty** dan **per mata uang**. Baris kembar pada sumbu itu ditolak.

Artefak: entitas `RETENSI_CEDANT`, kunci asingnya ke `VERSI_KONTRAK` dan `KELOMPOK_TREATY`, kolom
`KODE_MATA_UANG`, dan `INV-08`.

**PEMBUAT PERTAMA** untuk `RETENSI_CEDANT`.

**Kenapa begini:** Kunci alaminya **menyebut mata uang**, dan itu bukan kerapian. Sebuah kontrak dapat
menahan jumlah berbeda dalam dua mata uang untuk kelompok treaty yang sama; menulis `UNIQUE` tanpa
mata uang **mengubah artinya** menjadi *"dilarang dua baris bermata uang berbeda"* dan **menolak data
yang sah**. `INV-08` sempat berstatus **klaim, bukan penegakan**, justru karena kolom mata uangnya
belum ada — `P-8` golongan B yang membebaskannya dengan menaruh mata uang **di baris yang sama**,
sebagaimana bentuknya memang di sistem lama.

**Persyaratan:** `INV-08` (kelompok treaty + mata uang unik di dalam satu versi) · `INV-36` (nilai uang
terisi ⟹ mata uangnya terisi) · `INV-44` · `INV-47` — retensi masuk rekonsiliasi bersama penyerahan
(`INV-51`), dan **`SUMBU_REKONSILIASI`**-nya dinyatakan §10.11

**Tidak termasuk:** **`INV-51`** — *"retensi + penyerahan sama dengan 100 persen"* — bergolongan
**CONSTRAINT BELUM DIBUKTIKAN** (`SPEC-INVARIAN.md` §5.1) dan ditegakkan lewat *materialized view*
bersama `INV-47` dan `INV-50`. Ia berdiri di irisan `38`, **beserta pemantau kebasiannya**.
**Tingkat pencatatan** besaran retensi — `T-2`…`T-4` belum kembali dari teknik treaty.

**Jalur gagal:** Dua baris berkelompok treaty **dan** mata uang sama pada satu versi -> **ditolak**
`INV-08` · Nilai retensi terisi dengan mata uang kosong -> ditolak `INV-36` · Kelompok treaty yang
tidak ada di tabel acuan -> ditolak.

**Uji:** **Negatif:** sisipkan baris kembar pada sumbu penuh; isi nilai tanpa mata uang.
**Positif — dan ia yang menangkap lingkup kunci yang terlalu sempit:** satu versi dengan **dua baris
berkelompok treaty sama tetapi mata uang berbeda** -> **diterima**. Inilah uji yang memisahkan
`UNIQUE` yang benar dari `UNIQUE` tanpa mata uang; keduanya lulus uji negatif di atas.

**Menggantikan:** `P-07` melestarikan daftar `Retention` di sistem lama. Yang bergeser hanya
**kunci alaminya menjadi dapat dikompilasi** — di sistem lama tidak ada constraint apa pun atasnya.

**Blocked by:** 14 · 15

**Dasar:**
```
EVIDENCED(TreatyIn.Retention@ekspor-2026-09 - daftar per kelompok treaty; Currency pada baris yang sama)
        DECIDED(INV-08, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `RETENSI_CEDANT` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`, dengan `KODE_MATA_UANG` pada barisnya
- [ ] `INV-08` terpasang sebagai `UNIQUE` **bertiga kolom**, dan mata uangnya ada di dalamnya
- [ ] uji positif lulus: kelompok treaty sama, mata uang berbeda, **diterima**
- [ ] `INV-51` dinyatakan **di luar irisan ini** dengan irisan `38` sebagai pemiliknya
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`
