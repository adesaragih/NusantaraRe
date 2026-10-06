# Pertanyaan terbuka — apakah tabel pendaratan `T_TREATY_*` ikut membawa nilai `DITURUNKAN`?

**Untuk pemilik proses. Diajukan 6 Oktober 2026.**
**Satu pertanyaan, dan jawabannya memutuskan lima tabel sekaligus.**

---

## Pertanyaannya

> **Apakah tabel pendaratan `T_TREATY_*` ikut membawa nilai `DITURUNKAN`,
> atau hanya yang `DIPETAKAN`?**

| Jawab | Akibatnya |
| --- | --- |
| **"ikut"** | Kelima tabel di bawah **dibangun**, kolomnya diturunkan dari dokumen seperti ketiga belas tabel migrasi `437`. |
| **"tidak"** | Kelimanya **dicoret** dari daftar ke-36 — bukan dibiarkan menggantung sebagai pekerjaan yang belum selesai. |

⛔ Keduanya berbiaya, dan itu sebabnya pertanyaan ini tidak dijawab sendiri: membangun lima tabel
yang tidak boleh ada sama mahalnya dengan tidak membangun lima yang diperlukan.

---

## Kelima tabelnya

Jalur sumber masing-masing **100%** bertanda `DITURUNKAN — tidak disimpan (ADR-0037, §4)` di
`2-to-spec/PENELUSURAN-JSON-KE-KOLOM.md`:

| Tabel | Jalur bertanda `DITURUNKAN` |
| --- | ---: |
| `T_TREATY_LIMIT_MEASURE` | 7 dari 7 |
| `T_TREATY_LIMIT_SUMMARY` | 4 dari 4 |
| `T_TREATY_REINSTATEMENT` | 2 dari 2 |
| `T_TREATY_SHARE_SPREAD_AMOUNT` | 20 dari 20 |
| `T_TREATY_TOTAL` | 16 dari 16 |

---

## Mengapa ini pertanyaan, bukan hal yang sudah terjawab

**`ADR-0037` berbicara tentang MODEL RELASIONAL BARU** — `KONTRAK`, `VERSI_KONTRAK`, `LAYER`, dan
saudaranya. Di sana aturannya masuk akal: nilai yang dapat dihitung ulang tidak diberi kolom.

**Tabel `T_TREATY_*` adalah hal yang BERBEDA.** Ia **pendaratan** — ia mendaratkan dokumen Pega
apa adanya supaya layar dapat menampilkan apa yang sistem lama tampilkan. Bagi pendaratan,
"dapat dihitung ulang" bukan alasan untuk tidak menyimpan: yang menghitungnya di sistem lama
adalah Pega, dan Pega tidak akan ada.

### ⚠️ Dan keduanya SUDAH TERBUKTI BERSELISIH, di layar yang sedang berjalan

`Limits[].MDPList` dan `Limits[].PremiumEarnedList` ditandai **`DITURUNKAN`** di dokumen itu —
**tetapi tab Limits menampilkan keduanya hari ini.**

⇒ Jadi `DITURUNKAN` **tidak dapat** dibaca sebagai "jangan didaratkan" tanpa seseorang
menyatakannya. Kalau ia dibaca begitu, kedua medan yang sedang tampil di layar itu seharusnya
tidak pernah ada.

---

## Yang TIDAK dilakukan sambil menunggu

- ⛔ Kelima tabel **tidak dibangun**.
- ⛔ Kelimanya **tidak dicoret**.
- ⛔ Nol tebakan ke arah mana pun.

Ketiga tabel yang jalurnya `DIPETAKAN` dan cacah elemennya terukur lebih besar dari nol —
`T_TREATY_LIMIT_AMOUNT`, `T_TREATY_SHARE_AMOUNT`, `T_TREATY_FAC_SHARE_AMOUNT` — **dibangun lebih
dulu** oleh migrasi `438`, sebab keduanya tidak bergantung pada jawaban ini.
