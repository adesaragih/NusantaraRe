---
status: tertahan
---

# 07: Sirkulasi menutup dan menolak klaim, pemegangnya diselesaikan dari peran

*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Analis mengirim usulan menutup atau menolak klaim; sirkulasinya tetap satu
jenjang, dan pemegangnya diselesaikan **dari roster berdasarkan peran** — bukan dari nama
yang ditulis di dalam aturan. Bila tidak ada baris roster aktif yang memegang peran itu,
sirkulasi **ditolak saat dibuat**.

Jenis sirkulasi menutup dan menolak klaim pada `BentukSirkulasi`; baris fungsi akibat untuk
keduanya.

**Persyaratan:** `S-008`.

**Tidak termasuk:** —

**Jalur gagal:** nol baris roster aktif memegang peran itu → `422` menyebut **peran yang
kosong**, nol sirkulasi lahir, dan maksud tidak ditulis ke klaim.

**Uji:** P5-04 (jalur tutup klaim) · P5-05 (jalur tolak klaim) · **PS-01** (jenis dinyatakan
pemanggil) · **PS-02** (pemegang dari peran) · **PS-03** (penolakan saat peran kosong).
Ketiga uji `PS-xx` **dirancang gagal** terhadap sistem lama — deviasi 23, 24, 25 pada
`REGISTER-DEVIASI.md` bagian 6, seluruhnya **belum diratifikasi**.

**Menggantikan:** `CreateChildKomiteCloseNP_Act`·8, 9 — satu jenjang dengan pemegang dari nama
orang, surel, inisial, dan **dua sebutan jabatan berbeda untuk satu kedudukan** (`F-20`) ·
`CreateChildKomiteCloseNP_Act`·12 — jenis dibedakan dari nilai kolom analisis.

**Blocked by:**

- `05` — Akibat pada klaim dihitung di satu tempat dan ditulis dalam transaksi yang sama


**Dasar:** DECIDED(`K6-4`, keputusan beku no. 5). EVIDENCED: `F-20`.
`K6-4` menggantikan `H-2`, yang tetap tercatat `BUNYI HILANG`.

- [ ] Jenis sirkulasi datang dari pemanggil; nol pembacaan kolom analisis untuk menentukannya.
- [ ] Pemegang diselesaikan dari roster **berdasarkan peran**; nol nama orang di dalam aturan.
- [ ] Peran yang kosong **menolak pembentukan**, dan galatnya menyebut peran itu.
- [ ] Ketiga uji `PS-xx` **gagal** terhadap sistem lama. Uji `PS-xx` yang **lulus** paritas
      adalah tanda deviasinya tidak terpasang, dan itu kegagalan tiket ini.
- [ ] Deret uji tetap dieja `PS-xx`; ia **tidak** dinormalkan menjadi `P-xx`, karena
      penamaannya sengaja tidak terbaca sebagai ronde tujuh.

**Ketidakpastian:** Tidak ada.
