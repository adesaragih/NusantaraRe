# Pertanyaan terbuka — apakah daftar pilihan tab Portfolio hanya empat nilai itu?

**Untuk pemilik proses. Diajukan 6 Oktober 2026.**
**Tabnya SUDAH dibangun dan berjalan. Yang ditanya di sini hanya kelengkapan daftarnya.**

---

## Pertanyaannya

> **`Portfolio Type` dan `Premium / Loss Type` di Pega — apakah daftar pilihannya
> benar-benar hanya dua nilai masing-masing, atau aturan propertinya mengizinkan
> nilai lain yang kebetulan belum pernah dipakai?**

| Jawab | Akibatnya |
| --- | --- |
| **"hanya itu"** | Tidak ada yang berubah. Daftar yang sekarang sudah lengkap. |
| **"ada lagi: …"** | Nilai yang disebut ditambahkan ke `PORTOFOLIO.opsiArah` / `.opsiJenis` di `frontend/labels.ts`. Satu baris per nilai, nol perubahan lain. |

---

## Mengapa ini tidak dapat dijawab dari ekspor

Kedua sel di `Section/TreatyInTabsProportional.xml` berbunyi:

```
pyControlDisplayTitle = "Control inherited from property"
pyListDataSource      = (kosong)
```

Artinya daftar pilihannya **tidak ada di seksi itu** — ia hidup di `Rule-Obj-Property`
kelas `ASM-FW-GISFW-Data-TreatyInPortfolio`.

⛔ Dan **aturan properti tidak ikut diekspor**. Isi `D:\XML_NURE\Treaty In`:
`Activity` · `ConnectREST` · `DataTransform` · `DecisionTable` · `FlowAction` ·
`Harness` · `RDBList` · `ReportDefinition` · `Section` · `SystemSettings` · `When`.
Nol direktori `Property`. Pencarian `TypePortfolio` ke seluruh korpus mengembalikan
hanya seksi, harness, Activity, dan dokumen migrasi — nol definisi properti.

---

## Yang dipakai sebagai gantinya: PENGUKURAN, bukan tebakan

Sapuan 6 Oktober 2026 atas **seluruh 1.855 baris** `POOLDATA.M_TREATY_IN`
(nol `ROWNUM`, nol sampel; 1.855 dokumen terurai, **nol gagal urai**):
**1.925 elemen `Portfolio[]` di 845 dokumen.**

| Kunci | Nilai | Cacah | Kosong |
| --- | --- | ---: | ---: |
| `TypePortfolio` | `Withdrawal` | 1.578 | 0 |
| | `Assumption` | 347 | |
| `Type` | `Premium` | 1.032 | 0 |
| | `Loss` | 893 | |

Keempat kombinasinya terpakai: `Withdrawal/Premium` 845 · `Withdrawal/Loss` 733 ·
`Assumption/Premium` 187 · `Assumption/Loss` 160.

⚠️ **Yang pengukuran ini TIDAK dapat katakan:** nilai yang **sah tetapi belum
pernah dipakai** tidak meninggalkan jejak di data. Itulah seluruh isi pertanyaan ini.

---

## Yang sudah dikerjakan sambil menunggu, supaya jawabannya murah

- ⭐ Tabnya **berjalan** dengan keempat nilai terukur itu.
- ⭐ Nilai **di luar daftar tidak dijatuhkan**: `TabPortofolio.tsx` menambahkan
  nilai asing ke bagian bawah daftar pilihannya (`const asing = …`). Jadi kontrak
  yang membawa nilai kelima akan **tetap menampilkannya**, bukan diam-diam
  mengosongkannya — dan nilai itu akan terlihat oleh siapa pun yang membuka
  kontraknya.
- ⭐ Pilihan **kosong** disediakan, sebab baris baru memang lahir kosong
  (`Activity/TreatyInPropAdd.xml` hanya menetapkan `.Description = ""`).

Jadi jawaban "ada lagi" berbiaya **satu baris per nilai**, dan sampai jawabannya
tiba, nol data yang hilang dari layar.

---

## ⚠️ Temuan menyamping — BUKAN bagian pertanyaan ini, dan bukan milik repo ini

`D:\XML_NURE\_migration-docs\treaty-in\2-to-spec\PENELUSURAN-JSON-KE-KOLOM.md`
baris 267–268 memetakan:

```
Portfolio[].Type           -> PORTOFOLIO.ARAH_PORTOFOLIO
Portfolio[].TypePortfolio  -> PORTOFOLIO.JENIS_PORTOFOLIO
```

Pengukuran mengatakan sebaliknya: `Type` bernilai `Premium`/`Loss` (sebuah **jenis**),
`TypePortfolio` bernilai `Withdrawal`/`Assumption` (sebuah **arah**). Uji
`TestTiket24PortofolioEmpatKombinasiDiterima` memakai `ARAH ∈ {MASUK, KELUAR}` dan
`JENIS ∈ {PREMI, KLAIM}`, yang sejalan dengan pengukuran dan **tidak** sejalan dengan
kedua baris pemetaan itu.

⛔ **Tidak diubah di sini.** Berkas itu ada di luar repo, tabel `PORTOFOLIO` milik
model acuan (tiket 24) dan **belum dimuat satu baris pun**, dan menukar pemetaannya
adalah keputusan pemilik proses. Dicatat supaya tidak terbaca sebagai sudah beres.
