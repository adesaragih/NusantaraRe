# PROMPT — GILIRAN 16 *(folder `OUTPUT_HASIL_RNM`, cabang `main` @ `5138cfb` atau lebih baru)*: **OQ-N12 diputuskan (a) — `CLAIM_GROSS` = `CLAIM_AMOUNT` sementara**

> Hanya konteks **Claim Life, PremiumList Life, Komite Claim Life**. GILIRAN-3 *(A/B/C)*, 4–15 tetap rujukan. **Mulai hanya sesudah sesi
> Treaty lanjutan 3 selesai dan pohon `main` bersih** kecuali `App.tsx` milik work owner. Berkas `tco_*` tidak disentuh.

## 0. KEPUTUSAN WORK OWNER *(29-09-2026, kutipan: "A")* — **N12 (a)**

Catatan 7 tiket 03 **dipertahankan**: `CLAIM_GROSS` dibaca sama dengan `CLAIM_AMOUNT` baris adjustment, bertanda
`[sementara — menunggu OQ-N11 pemilik ekspor]`. Dasar, diperiksa asisten pada `Activity/SaveOutStandingLife_Act.xml`:

| Fakta | Bukti |
| --- | --- |
| gerbang Save to RNM menolak bila `CLAIM_GROSS` kosong | langkah **11.17.1** `Page-Set-Messages` b5936, syarat `.CLAIM_GROSS==""` b6043 |
| langkah itu **hidup** | tidak ada `pyStepsBlockName` `//` di dalam langkah 11.17.1; `//` terdekat b6178 milik langkah **12** *(ref b6167)* |
| `CLAIM_GROSS` tanpa penulis di korpus | pembaca: `CountClaimAmountLife_Act`, `RejectOSClaimLife_Act`, `SaveOutStandingLife_Act`, `SpreadingClaimLife_Act`; penulis `PropertiesName .CLAIM_GROSS`: **0** |
| kesimpulan | bila `CLAIM_GROSS` benar-benar kosong di produksi, **tidak satu klaim pun** dapat melewati Save to RNM di Pega — mustahil untuk sistem yang dipakai; nilainya pasti diisi rule yang tidak diekspor *(dugaan: Declare Expression)* |

## 1. KERJAKAN — satu commit

- Tiket 03: blok bertanggal *"Keputusan work owner 29-09-2026 — OQ-N12 (a)"* dengan tabel bukti §0; catatan 7 diberi label sementara.
- `OQ-untuk-tim.md`: OQ-N12 **ditutup (a)**; OQ-N11 di daftar pemilik ekspor ditambah kalimat: *"mohon ekspor Declare Expression atau
  rule lain yang mengisi `CLAIM_GROSS` pada kelas `Int-LIFE_PREMIUM_DETAIL`/`AdjustmentList`"*.
- Kode: satu komentar di fungsi yang menyamakan `CLAIM_GROSS` dengan `CLAIM_AMOUNT` merujuk OQ-N12/N11; **perilaku tidak berubah**; uji
  yang menguji gerbang 11.17.1 tetap hijau.
- `058_seq_work_polis_mulai_ulang.sql`: kepala berkas ditambah perintah audit angka 22373 *(asisten, 28-09-2026)*:
  `SELECT COUNT(*), MAX(TO_NUMBER(REGEXP_SUBSTR(IDPEGA,'[0-9]+$'))) FROM POOLDATA.JSON_POLIS WHERE IDPEGA LIKE '%NBLF-%'` → 33 baris,
  maksimum 22373. **Berkas migrasi yang sudah dijalankan tidak boleh berubah pernyataannya** — hanya komentar; bila 058 sudah terpasang
  di DEV, tulis perintah audit itu di tiket 00 PremiumList saja.

Commit `docs: OQ-N12 (a) — CLAIM_GROSS = CLAIM_AMOUNT sementara; audit angka 22373`.

## 2. LAPORAN

Satu pesan pendek: commit, angka uji dengan dan tanpa tag `db`, OQ ditutup.

---

*Disusun 29 September 2026 dari jawaban work owner "A" dan pembacaan batas langkah 11.17.1 / 12 di `SaveOutStandingLife_Act.xml`
(b5936, b6043, b6167, b6178).*
