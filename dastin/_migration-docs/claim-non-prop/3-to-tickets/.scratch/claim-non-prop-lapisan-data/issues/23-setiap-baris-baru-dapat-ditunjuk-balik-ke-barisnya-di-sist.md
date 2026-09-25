---
status: selesai
---

# 23: Setiap baris baru dapat ditunjuk balik ke barisnya di sistem lama


> **SELESAI 19 September 2026 — view-nya ditulis ulang dari satu tabel menjadi sembilan belas.**
>
> Versi sebelumnya hanya menggabungkan `AKSEPTASI`. Akibatnya dua, dan **keduanya diam**:
>
> 1. Baris korelasi untuk **18 tabel lain** tetap muncul, tetapi `KEADAAN_BARIS`-nya `NULL` — bukan karena barisnya belum lengkap, melainkan **karena view-nya tidak melihat tabelnya**.
> 2. Penyaring `KEADAAN_BARIS = 'LENGKAP'` yang dijanjikan butir 2 akan **membuang seluruh 18 tabel itu tanpa sepatah pesan**.
>
> Itu persis cacat yang ADR-0019 larang: **`NULL` diperlakukan sebagai nilai** — dan di alat paritas, baris yang hilang adalah **selisih yang tidak pernah terhitung**.
>
> **Tiga keadaan dibedakan sekarang, dan tidak satu pun `NULL`:**
>
> | Nilai | Artinya |
> |---|---|
> | `LENGKAP` / `MENUNGGU_KURS` / `GAGAL_URAI` | tabelnya punya kolom keadaan, dan inilah isinya — 11 tabel |
> | `TIDAK_BERLAKU` | tabelnya **tidak** punya kolom keadaan — 8 tabel acuan dan catatan. Berbeda dari "belum lengkap" |
> | `TABEL_TIDAK_DIKENALI` | baris korelasi menunjuk tabel di luar daftar. **Tanpa cabang ini ia hilang dari view** |
>
> Kolom `KEBERADAAN` memisahkan pertanyaan kedua dari yang pertama: `'HILANG'` berarti baris kanoniknya tidak ada, apa pun keadaannya — itu **kegagalan migrasi**, bukan baris yang belum lengkap.
>
> **Tiga tabel sengaja di luar**: `MIGRASI_KORELASI`, `MIGRASI_PENDARATAN`, `MIGRASI_NILAI_DITOLAK`. Ketiganya jembatan, bukan sasaran — baris korelasi tidak pernah menunjuk tabel korelasi. Bila ternyata ada, cabang terakhir menangkapnya, dan itu memang yang seharusnya terjadi.

*Asal: `T-19` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Perbandingan baris per baris mungkin dilakukan.

`V_PARITAS_SHADOW`.

**Tidak termasuk:** Pembandingnya sendiri dan toleransinya — AK-3.

**Blocked by:**

- ~~`15`~~ *(selesai)* — Jembatan ke sistem lama berdiri sebagai tabel terpisah


**Dasar:** DECIDED(ADR-0005). EVIDENCED: `BLUEPRINT.md` §8.4 — `IndexObject` adalah **posisi numerik, bukan surrogate key**, sehingga jembatannya hilang begitu baris berpindah ke kunci sendiri.

- [x] Setiap baris korelasi muncul **tepat sekali** — cabangnya disaring `TABEL_TUJUAN`, ke-19 nilainya saling lepas, dan cabang terakhir menampung sisanya.
- [x] Baris ber-`KEADAAN_BARIS` bukan `LENGKAP` dapat **dikeluarkan** lewat satu penyaring, bukan dihitung sebagai selisih — dan **tanpa ikut membuang 18 tabel** yang tidak punya kolom keadaan.
- [x] Baris kanonik yang **hilang** terbaca sebagai `KEBERADAAN = 'HILANG'`, terpisah dari soal keadaan.

**Ketidakpastian:** Baseline pembandingnya sendiri belum tentu benar: `TotalUR` selalu nol dan dua rumus premi pemulihan bekerja atas nilai berbeda. AK-2b menempatkan ketiganya sebagai pengecualian bernama.
