---
status: aktif
---

# 01: Roster dan tabel seleksi berdiri sebagai data yang dikelola

*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Pemilik proses menyemai dan menyunting roster jenjang dan tabel seleksi
tanpa menyentuh kode; konfigurasi seleksi yang meninggalkan satu pun kombinasi masukan tanpa
aturan **ditolak saat disimpan**, bukan ditemukan oleh klaim pertama yang jatuh ke dalamnya.

`ROSTER_JENJANG` dan `ATURAN_SELEKSI` beserta sequence dan constraint-nya; operasi
`KelolaRoster` dan `KelolaTabelSeleksi`; jalur baca dan tulis keduanya; layar pengelolaan.

**Persyaratan:** `S-050`, `S-051`, `S-053`, `S-054`, `S-055`, `S-056`, `S-057`, `S-036`,
`S-069`.

**Tidak termasuk:** `S-052` — pergantian pemegang tercatat sebagai peristiwa. Ia menunggu
objek `PERISTIWA_ROSTER` ditambahkan ke model data; `PERISTIWA_SIRKULASI` tidak dapat
menampungnya karena pengenal sirkulasinya wajib. **Nama** peran administratif juga tidak
ditetapkan di sini — gerbangnya dibangun, namanya menunggu pemilik proses.

**Jalur gagal:** baris yang menduplikasi kombinasi masukan yang sudah ada → `422`, tabel yang
berlaku tidak berubah · penyimpanan yang meninggalkan kombinasi tak tercakup → `422` dengan
galat yang **menyebut kombinasinya** · delegasi menunjuk orang yang tidak aktif → `422`, baris
roster tidak berubah · dua baris roster aktif berderajat sama dalam satu kelas → ditolak
constraint · pemanggil tanpa peran administratif → `403`.

**Uji:** P5-09, P5-10, P5-11 (kelas dari kombinasi masukan) · P5-16 (hak diperiksa) · BARU
untuk penolakan cakupan tidak lengkap dan untuk kosong-berarti-tidak-membatasi.

**Menggantikan:** `FilterEmailKomiteWithLimit` — penyaring roster; `.LIMIT_TOP` **dibuang**,
ia diambil tetapi tidak pernah menyaring (`F-14`) · `CreateChildKomiteCNP_Act`·12, 13 — dua
tetapan kelas di dalam langkah · `CreateChildKomiteCNP_Act`·26.8.3 — substitusi satu orang ke
orang lain, **dibuang** (`F-19`) · `SetKomiteList_Act` — pemetaan nama orang ke jabatan
(`F-18`) · `EMAILKOMITE` sebagai sumber runtime.

**Blocked by:**

- None (can start immediately)


**Dasar:** DECIDED(`K5-5`, `K6-1`, `K6-3`, `D-5`, `H-1`, `H-3`, `H-7`, keputusan beku no. 5,
ADR-0006, ADR-0016). EVIDENCED: `F-14`, `F-16`, `F-18`, `F-19`.

- [ ] Aturan seleksi disunting sebagai data; perubahan berlaku tanpa rilis.
- [ ] Penyimpanan yang meninggalkan satu pun kombinasi masukan tanpa aturan **ditolak**, dan
      galatnya menyebut kombinasi yang tidak tercakup.
- [ ] Keempat kolom batas dan penanda bersyarat pada `ATURAN_SELEKSI` boleh kosong, dan
      **kosong berarti "tidak membatasi", bukan nol** — dinyatakan eksplisit di komentar
      kolom **dan** di kaki berkas DDL. Tanpa pernyataan itu baris semai akan "dirapikan"
      jadi nol oleh pembaca berikutnya, dan seluruh isinya berubah arti (`ADR-0019`, `K5-6`).
- [ ] Isi awal disemai dari nilai yang **berpenulis** di sistem lama, meninggalkan batas
      nilainya kosong; nol angka ambang di dalam kode, dan angka yang hanya ada pada
      deskripsi langkah **tidak** dipakai (`K6-1`).
- [ ] Baris roster yang dinonaktifkan **tetap tersimpan**; mencabut seseorang tidak menghapus
      jejaknya.
- [ ] Delegasi tetap adalah atribut roster, bukan penukaran nama di dalam aturan.
- [ ] Nol cabang berdasarkan identitas orang di jalur mana pun — diuji dengan pencarian atas
      kode, bukan dengan memanggil operasi.
- [ ] Nama tabel, kolom, constraint, dan sequence di bawah 30 byte, berbentuk
      `<peran>_<tabel>[_n]` yang tidak mengeja kolom; pemeriksa penamaan keluar dengan kode
      gagal bila satu pun melewati batas.

**Ketidakpastian:** **Nama peran administratif** belum ditetapkan — keputusan pemilik proses.
Gerbangnya dibangun sekarang; bila jawabannya lain, yang berubah adalah **nama peran yang
dicocokkan**, bukan adanya gerbang. Menunggui, tidak menahan.
