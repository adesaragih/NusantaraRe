---
status: aktif
golongan: baru
---

# 43: Pelaksana migrasi mematikan penegakan trigger selama pemindahan dan menyalakannya kembali, dan keadaan sakelarnya terlihat

*Asal: `DAFTAR-PEKERJAAN.md` `P-53` · `SPEC-INVARIAN.md` §4.3b, CO-1…CO-8 · `ADR-0042`.*

**What to build:** **PM** mematikan penegakan trigger selama pemindahan data, menyalakannya kembali
sesudahnya, dan **keadaan sakelar itu terlihat tanpa membuka basis data**. Menyalakan kembali
**memeriksa ulang** baris yang masuk selama ia mati.

**Kenapa begini:** `ADR-0042` memerintahkan sejarah **pindah apa adanya**, dan sebagian data lama
**memang melanggar** aturan yang sistem baru tegakkan — kunci alami yang berubah di tengah jalan,
keadaan yang tidak punya padanan, nilai nol yang sebenarnya kegagalan. Trigger yang hidup akan
menolaknya, dan pemindahan **berhenti pada baris pertama yang menyimpang**. Mematikannya adalah jalan
yang benar; yang berbahaya adalah **lupa menyalakannya kembali**.

> Sakelar yang keadaannya hanya terbaca dengan membuka basis data akan dibiarkan mati, dan **tidak
> ada yang tahu sampai data yang melanggar sudah masuk lewat jalur biasa**. Itu bentuk kegagalan yang
> sama dengan *materialized view* yang berhenti me-refresh — penegakan yang berhenti **tanpa satu
> galat pun**.

**Persyaratan:** `SPEC-INVARIAN.md` §4.3b dan **CO-1…CO-8** · `ADR-0042` · `INV-26`
(`WARISAN_TAK_TERPETAKAN` hanya dapat dimasuki lewat migrasi, **tidak pernah oleh sistem berjalan**) —
sakelar ini **yang membuat larangan itu dapat ditegakkan**

**Tidak termasuk:** **Pemindahan datanya sendiri** — irisan `44` untuk tabel acuan, dan sisanya batch
2. Irisan ini membangun **sakelarnya**, dan mengujinya dengan data uji.
**Daftar trigger yang dimatikan** — daftarnya **dikumpulkan dari tiket yang memasangnya**, dan tiap
tiket bertrigger wajib mendaftarkan diri: irisan `18` (`INV-19`) dan irisan-irisan batch 2 (`INV-28`,
`INV-54`). Yang dibangun di sini **mekanismenya**, bukan daftar isinya.
**Mematikan constraint yang bukan trigger** — `CHECK` dan `UNIQUE` **tidak dimatikan**. Membiarkan
kunci ganda masuk merusak hal yang tidak dapat diperbaiki tanpa memilih baris mana yang dibuang.

**Jalur gagal:** Menyalakan kembali tanpa memeriksa ulang baris yang masuk -> **kriteria selesai tidak
terpenuhi** · Keadaan sakelar tidak terlihat di luar basis data -> idem · Pemeriksaan ulang menemukan
pelanggaran -> **dilaporkan per baris**, dan pemindahannya **tidak dinyatakan selesai** · Sistem
berjalan menulis baris saat sakelar mati -> **dilarang**; sakelar mati hanya sah selama jendela
pemindahan, dan jendelanya tercatat.

**Uji:** **Negatif:** nyalakan kembali dengan sengaja menyisipkan baris melanggar saat mati —
pemeriksaan ulang **harus menemukannya**; matikan sakelar lalu tulis lewat jalur aplikasi biasa —
**harus ditolak**.
**Positif — dan ia yang menangkap sakelar yang terlalu lebar:** saat sakelar mati, `UNIQUE` dan
`CHECK` **tetap menolak** baris yang melanggar. Sebuah sakelar yang mematikan **segalanya** lulus uji
negatif pertama dan **membiarkan kerusakan yang tidak dapat diperbaiki**.
**Positif kedua:** keadaan sakelar terbaca **benar** pada kedua posisinya, diperiksa dari luar basis
data.

**Menggantikan:** tidak ada. **Golongannya BARU** — sistem lama tidak punya trigger yang perlu
dimatikan, sebab ia hampir tidak punya penegakan di basis data sama sekali.

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-1…CO-8**, dan §4.3b |
| 2 | **siapa membaca, seberapa sering** — **PM, pada setiap pembukaan dan penutupan jendela pemindahan**; dan **pemilik proses Treaty In, harian** selama gelombang pemindahan berjalan. Keadaan sakelar adalah hal yang **harus dilihat tiap hari**, bukan tiap minggu |
| 3 | **ambang berangka** — jendela pemindahan dinyatakan tertutup ketika sakelar **menyala** dan pemeriksaan ulang mengembalikan **nol pelanggaran yang belum diadili**. Pelanggaran yang sudah diadili dan diputuskan tetap dibawa **dihitung terpisah** dan tidak menahan penutupan |
| 4 | **siapa boleh menyalakan atau mematikan** — **PM**, dan **hanya** di dalam jendela pemindahan yang pemilik proses buka |

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(SPEC-INVARIAN.md §4.3b - urutan penyalaan penegakan; CO-1..CO-8)
        DECIDED(ADR-0042, INV-26)
        DIASUMSIKAN-CLEAR(L-3)
```

- [ ] sakelar berdiri, dan keadaannya **terbaca dari luar basis data** pada kedua posisinya
- [ ] menyalakan kembali **memeriksa ulang** baris yang masuk selama mati, dan melaporkan pelanggaran **per baris**
- [ ] uji positif lulus: `UNIQUE` dan `CHECK` **tetap menolak** saat sakelar mati
- [ ] jalur aplikasi biasa **ditolak** selama sakelar mati — diuji, bukan diandaikan
- [ ] **daftar trigger** yang tunduk pada sakelar ini dibangkitkan dari tiket yang memasangnya, bukan diketik ulang
- [ ] jendela pemindahan **tercatat** — kapan dibuka, kapan ditutup, oleh siapa
- [ ] keempat butir **CARA MENYALAKANNYA** terisi, termasuk ambang berangka
- [ ] `L-3` tercatat di `ASUMSI-CLEAR.md`: ujinya menuntut instans yang terjangkau
