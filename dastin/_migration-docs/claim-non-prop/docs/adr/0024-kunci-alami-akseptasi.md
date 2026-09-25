---
status: accepted
label: DECIDED
---

# Satu akseptasi per klaim per layer per mata uang

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas** (2026-09-08/09); dan `pengetahuan/DDL_Script_ClaimNonProp.xls` **versi 2026-09-18 10:36, 48 objek**.

Kunci alami baris akseptasi adalah **(klaim, layer, mata uang)**, ditegakkan sebagai `UNIQUE`, dengan primary key surrogate.

Sumbernya bukan tebakan: blok yang **dikomentari** di dalam `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP` memakai persis kombinasi itu —

```sql
SELECT count(1) INTO id_count FROM OS_AKSEPTASI_KLAIM a
 WHERE CASEID = PegaID AND a.data_json.TypeLoss = LayerT AND a.data_json.Currency = Currency;
```

Orang yang menulis blok itu sedang menjawab pertanyaan yang sama dan sampai pada kesimpulan yang sama. Kodenya dimatikan; alasannya tidak.

## Penguat dari bentuk penyimpanan — `CLAIMXOL`

Dipindahkan ke sini dari `BLUEPRINT.md` §13.2 pada 18 September 2026, karena tempatnya di sini.

View `POOLDATA.CLAIMXOL` membongkar JSON akseptasi dengan `JSON_TABLE`:

```sql
json_table (data_json, '$.CNPLayerList[*]'
  columns( XOL varchar2 path '$.XOL',
    nested path '$.CNPCurrencyList[*]' columns(
      TotalXOLGross varchar2 path '$.TotalXOLGross', ... )))
```

**`$.CNPLayerList[*]` membungkus `$.CNPCurrencyList[*]`** — layer di tingkat luar, mata uang di tingkat dalam. Jadi bentuk penyimpanan lama sendiri sudah bersusun dua tingkat persis menurut dua dimensi kunci ini, dan urutannya menetapkan arah sarangnya: satu layer memuat banyak mata uang, bukan sebaliknya.

Ini penguat, bukan dasar. Dasarnya tetap blok yang dikomentari di atas, karena di sanalah kombinasinya dinyatakan sebagai **kunci**, bukan sekadar sebagai susunan.

## Consequences

Sistem lama tidak menegakkannya sama sekali: `OS_AKSEPTASI_KLAIM` **tidak punya primary key maupun unique constraint**, dan prosedurnya **selalu `INSERT`, tidak pernah `UPDATE`** karena logika upsert-nya dikomentari seluruhnya. Baris ganda bukan kemungkinan teoretis.

**Karena itu pelanggaran kunci akan ditemukan, dan itu bukan bukti kuncinya salah.** Saat REQ-018 kembali, hasilnya dipisah tiga:

| Kelompok | Perlakuan |
|---|---|
| Baris **identik di semua nilai** | duplikat murni — digabung |
| Baris **berbeda, dengan pola waktu** | ambil yang terakhir, **catat yang dibuang** |
| Baris **berbeda, tanpa pola waktu** | bawa contohnya untuk ditinjau — berarti ada dimensi pembeda yang belum tertangkap |

Kelompok ketiga adalah satu-satunya yang dapat menggugurkan kunci ini. Dua kelompok pertama adalah akibat bug selalu-`INSERT`, bukan bukti tentang model data.
