# Fixture kasus — hasil de-identifikasi (tiket 15)

Lima berkas kasus nyata (`DDL\P-5 *.txt` korpus) yang sudah melewati
`backend/alat/deidentifikasi`. Disimpan di sini atas keputusan work owner 1 Oktober 2026
(`docs/KEPUTUSAN-30-09-2026.md` butir 36–38).

| Berkas | Siklus | Lini |
| --- | --- | --- |
| `nb-fire-1.json` | New Business | FIRE |
| `nb-kredit-1.json` | New Business | Asuransi Kredit |
| `nb-marinecargo-1.json` | New Business | MARINE CARGO |
| `rnw-fire-1.json` | Renewal | FIRE |
| `edm-fire-1.json` | Endorsement | FIRE |

- **Yang dikosongkan:** nilai kunci `DaftarBuang` (identitas, nomor polis/kasus), `DaftarKosongkan`
  (teks bebas, bagian alamat), dan jalur `JalurBuang` (`ID` tingkat dokumen). Kuncinya tetap ada.
- **Yang utuh:** setiap daun lain, termasuk setiap angka (seluruhnya disimpan sebagai teks).
- **Nama berkas netral** — padanan dengan nomor kasus asli **sengaja tidak disimpan** di repositori.
- ⛔ Jangan menyalin berkas mentah ke sini. `TestFixtureKasusSudahDibersihkan` gagal bila ada daun
  BUANG yang berisi atau nama berkas tidak berbentuk `siklus-lini-urut.json`.
- Menambah fixture: jalankan alat dengan `-keluar` ke folder di luar repositori, pastikan
  kebocoran 0, beri nama netral, lalu salin ke sini dan ubah jumlah di uji penjaga.

## Sensus coverage per lini (audit `TestBerkasKasusP5`, 01-10-2026)

Jendela: jalur pemilih rumus tiap lini (`rekonsiliasi.jalurCoverage`; CountGrossPremi_Act 5.1.2 / 5.2.1,
CountGPWMarinePAMbu_Act 1.1.1) — **bukan** seluruh berkas: grep mentah memungut salinan coverage di wadah lain.
Cara pertama = `telusuri` Go (tesnya sendiri); cara kedua = urai Python independen, dijalankan dari folder ini:

```
PYTHONIOENCODING=utf-8 py -c "
import json,glob
J={'FIRE':['LocationList','Property','PropertyItemList','CoverageList'],'ANEKA':['LocationList','Property','RiskLocation','OccupationList','AnekaList','CoverageList'],'MARINE CARGO':['CargoList','CoverageList']}
def n(m,k):
    if not k: return 1
    v=m.get(k[0]) if isinstance(m,dict) else None
    return sum(n(x,k[1:]) for x in v) if isinstance(v,list) else (n(v,k[1:]) if isinstance(v,dict) else 0)
for b in sorted(glob.glob('*.json')):
    a=json.load(open(b,encoding='utf-8')); print(b,{l:n(a,j) for l,j in J.items() if n(a,j)})"
```

Hasil kedua cara sama: edm-fire-1 FIRE 45 · nb-fire-1 FIRE 5 · rnw-fire-1 FIRE 4 · nb-kredit-1 ANEKA 1 ·
nb-marinecargo-1 MARINE CARGO 104. Bentuk nilai di jalur FIRE (urai, 54 coverage): `NetRate` `"0"` 54 ·
`PctAdjustment` `"0.0"` 49, `"75"` 5 · `LostLimit` `"100"` 54 · `CoverageBasis` `"1"` 54. Master policy: akar
`QuotationData.PolicyType` `"0"` di keempat kasus non-marine, `"2"` di nb-marinecargo-1; `IsMOP` tidak ada di kelimanya.
