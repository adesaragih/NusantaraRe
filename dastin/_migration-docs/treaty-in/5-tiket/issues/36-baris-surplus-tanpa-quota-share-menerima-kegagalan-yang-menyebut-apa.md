---
status: aktif
golongan: baru
---

# 36: Baris surplus tanpa baris quota share pada versi yang sama menghasilkan kegagalan yang menyebutkan apa yang kurang

*Asal: `DAFTAR-PEKERJAAN.md` `P-23` · `SPEC-INVARIAN.md` `INV-35` · `ADR-0035`.*

**What to build:** **PK** yang mencatat baris `SURPLUS` sementara versi itu **tidak punya** baris
`QUOTA_SHARE` menerima **kegagalan yang menyebutkan apa yang kurang** — bukan galat tanpa isi, dan
bukan hasil perhitungan bernilai nol.

**Kenapa begini:** Kapasitas surplus dihitung **retensi × jumlah lines**, dan retensinya datang dari
susunan quota share. Tanpa baris quota share, **penyebutnya tidak ada** — dan sistem lama menyelesaikan
keadaan itu dengan menghasilkan **nol**. Nol yang sebenarnya *"tidak dapat dihitung"* **tidak dapat
dibedakan lagi setelah tersimpan**, dan itu persis larangan `ADR-0035`. `INV-35` menyebutnya dengan
kata-kata itu: *"ketiadaannya kegagalan yang dilaporkan, bukan nol"*.

**Persyaratan:** **`INV-35`** (baris `SURPLUS` menuntut adanya baris `QUOTA_SHARE` pada **versi yang
sama**), bergolongan **APLIKASI** · `ADR-0035` · `INV-42` (nilai uang tidak pernah memuat angka
penanda kegagalan) · `INV-48` (kapasitas surplus adalah **kelipatan**, pengecualian bernama terhadap
`INV-47`)

**Tidak termasuk:** **Aturan tepat-satu pada satu baris** — irisan `35`. Yang di sini **satu baris
terhadap saudaranya di versi yang sama**.
**`KAPASITAS_SURPLUS` sebagai kolom tersimpan** — tidak ada; turunan.
**Kenapa `INV-35` di lapisan aplikasi dan bukan constraint** — ia melintasi baris di dalam versi lewat
`LAYER`, dan bentuk itu tidak dapat ditulis sebagai `CHECK`. Letaknya **keputusan**, bukan kelalaian.

**Jalur gagal:** Simpan baris `SURPLUS` pada versi tanpa baris `QUOTA_SHARE` -> **kegagalan bernama**
yang menyebut *"tidak ada baris quota share pada versi ini"* · Kapasitas surplus dihitung dan
menghasilkan **nol** -> **kriteria selesai tidak terpenuhi**; itu perilaku yang irisan ini hapus ·
Galat tanpa isi -> idem.

**Uji:** **Negatif:** simpan surplus tanpa quota share pada versi kosong; dan pada versi yang punya
layer lain tetapi tanpa quota share sama sekali.
**Positif — dan ia yang menangkap pemeriksaan yang dijalankan pada lingkup yang salah:** baris
`SURPLUS` pada versi yang punya baris `QUOTA_SHARE` **di layer yang berbeda** -> **diterima**.
`INV-35` berbunyi *"pada versi yang sama"*, bukan *"pada layer yang sama"* — pemeriksaan yang
dijalankan per layer menolak susunan yang sah dan lulus kedua uji negatif di atas.

**Menggantikan:** perilaku sistem lama yang menghasilkan **nol** ketika penyebutnya tidak ada.
**Golongannya BARU** — pemeriksaannya memang belum pernah ada — sehingga `CARA MENYALAKANNYA` wajib,
dan hasilnya **tidak dapat diuji terhadap data lama**: data lama memuat nol yang artinya bukan nol.

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-6, CO-7, CO-8** |
| 2 | **siapa membaca, seberapa sering** — **pemilik proses Treaty In, mingguan**. Angka pertamanya menjawab berapa versi warisan memuat surplus tanpa quota share |
| 3 | **ambang berangka** — kuasa memblokir dinyalakan ketika pelanggaran pada versi **yang lahir di sistem baru** nol selama **4 minggu berturut-turut**; versi warisan dihitung terpisah |
| 4 | **siapa boleh menyalakan** — pemilik proses Treaty In |

**Blocked by:** 34

**Dasar:**
```
EVIDENCED(SPEC-INVARIAN.md INV-35 - "ketiadaannya kegagalan yang dilaporkan, bukan nol")
        DECIDED(ADR-0035, INV-35, INV-48)
```

- [ ] `INV-35` terpasang **di lapisan aplikasi**, dan letaknya tertulis sebagai keputusan yang menyebut sebabnya
- [ ] kegagalannya **menyebut apa yang kurang**, dan keterangannya **tersimpan**, bukan hanya muncul di layar
- [ ] tidak ada jalur yang menghasilkan kapasitas surplus bernilai nol saat penyebutnya tidak ada — diperiksa, bukan diandaikan
- [ ] uji positif lulus: quota share di layer berbeda pada versi yang sama **diterima**
- [ ] keempat butir **CARA MENYALAKANNYA** terisi
- [ ] hitungan versi warisan yang melanggar dilaporkan terpisah
