---
status: aktif
golongan: pelestarian
---

# 31: Layer dicatat beserta limit, deductible, MDP, dan ketentuan pemulihan limitnya

*Asal: `DAFTAR-PEKERJAAN.md` `P-16` · `SPEC-MODEL-DATA.md` §10.3, §10.3a, §10.3b · `SPEC-INVARIAN.md` `INV-05`.*

**What to build:** **PK** mencatat layer sebuah kontrak non-proporsional — nomor layer dan bagiannya,
limit, deductible, tarif penyesuaian, minimum deposit premium, dan **ketentuan pemulihan limitnya** —
lalu membacanya kembali utuh.

Artefak: entitas `LAYER`, `PEMULIHAN_LIMIT`, dan kedua tabel anak paket uang `NILAI_MDP` dan
`NILAI_MDP_MINIMUM`; kunci asingnya; `INV-05`.

**PEMBUAT PERTAMA** untuk `LAYER`, `PEMULIHAN_LIMIT`, `NILAI_MDP`, `NILAI_MDP_MINIMUM`.

**Kenapa begini:** `LAYER` **diselesaikan atas sumber yang kurang**, dan itu tercatat: pohon 985
simpul tidak pernah melihat **enam belas properti** kelas `Data-TreatyInLimits` (`L-8`). Dari yang
kemudian diadili lahir `LIMIT_AGREGAT` — batas total sepanjang periode, terpisah dari limit per
kejadian — yang **hilang tanpa disengaja** karena kembaran mata uang keduanya dibuang lebih dulu
sementara kembaran pertamanya tidak pernah ikut masuk. **Menghilangkannya mengubah arti kontrak.**
Dan `MDP` bukan satu angka: ia **daftar per mata uang** di sistem lama, sehingga ia tabel anak —
golongan A menurut `P-8`.

**Persyaratan:** `INV-05` (`NOMOR_LAYER` + `BAGIAN_LAYER` unik di dalam satu versi) · `INV-36` ·
`INV-41` · `INV-46` (tidak ada pasangan kolom kembar untuk dua mata uang) · `INV-49` (pemulihan
adalah **besaran berulang**, dikecualikan dari partisi `INV-47`) · `ADR-0035` — `DEDUCTIBLE_KEDUA`
tersimpan **sebagai teks** di sistem lama, dan baris yang gagal dikonversi **tidak disamarkan menjadi
nol**

**Tidak termasuk:** **Syarat yang berbeda antar pemulihan** — irisan `32`. Di sini pemulihan
tersimpan sebagai daftar; **nilainya boleh berbeda** adalah kemampuan tersendiri dan bergolongan
**BARU**.
**Bagian NuRe atas layer** — irisan `33`.
**Nasib `PERSEN_ROL` dan premi diperoleh** — **`F-15`**. `PERSEN_ROL` berdiri di §10.3 atas putusan
`TDA-17`, dan `F-15` membantahnya dengan rumus `DetailCalculationROL`: `premi ÷ limit × 100`. Tabel
anak `NILAI_PREMI_DIPEROLEH` **tidak dibuat**. Keputusannya pemilik proses; irisan ini **tidak
mendahuluinya ke arah mana pun**.
**Tingkat pencatatan `LIMIT_AGREGAT`** — `T-3`.

**Jalur gagal:** Dua layer bernomor **dan** berbagian sama pada satu versi -> **ditolak** `INV-05` ·
`DEDUCTIBLE_KEDUA` warisan yang tidak dapat dikonversi dari teks -> **kegagalan bernama**, bukan nol
(`ADR-0035`) · Baris MDP tanpa mata uang -> ditolak · `PERSEN_ROL` bernilai `9989998` -> ditolak
`INV-42`.

**Uji:** **Negatif:** layer kembar pada sumbu penuh; MDP tanpa mata uang; sisipkan `9989998`;
konversi `DEDUCTIBLE_KEDUA` dari teks yang bukan angka.
**Positif — dan ia yang menangkap kunci yang terlalu sempit:** satu versi dengan **layer 1 bagian A
dan layer 1 bagian B** -> **diterima**. `UNIQUE` atas nomor layer saja lulus uji negatif di atas dan
menolak susunan berlapis yang justru pokok kontrak non-proporsional.
**Positif kedua:** satu layer dengan **dua baris MDP bermata uang berbeda** -> diterima.
**Positif ketiga:** satu layer dengan **dua baris pemulihan limit** -> diterima, dan `INV-47`
**tidak** menuntut keduanya berjumlah utuh — itu yang `INV-49` kecualikan.

**Menggantikan:** `P-16` melestarikan daftar `Limits[]`. Yang bergeser: `LIMIT_AGREGAT`,
`MDP_MINIMUM`, `PERSEN_MDP_MINIMUM`, `MDP_DIGABUNG`, dan `TANPA_HITUNG_PREMI_PEMULIHAN`
**ditambahkan** — kelimanya ada di sistem lama dan **tidak pernah terlihat** oleh pohon (`L-8`,
§10.3a). Dan `PEMULIHAN_LIMIT` berdiri sebagai entitas tersendiri (§10.23c butir 1).

**Blocked by:** 14 · 15

**Dasar:**
```
EVIDENCED(SPEC-MODEL-DATA.md §10.3a - 16 properti Data-TreatyInLimits, diadili satu per satu)
        EVIDENCED(SetReinstatementPct@ekspor-2026-09 - ReinstatementPct="100", AdditionalPct="100" sebagai tetapan)
        DECIDED(INV-05, INV-49, ADR-0035, KTV-A, KTV-C)
        DIASUMSIKAN-CLEAR(KTV-A)
        DIASUMSIKAN-CLEAR(T-6)
        DIASUMSIKAN-CLEAR(T-3)
        DIASUMSIKAN-CLEAR(F-15)
```

- [ ] `LAYER`, `PEMULIHAN_LIMIT`, `NILAI_MDP`, `NILAI_MDP_MINIMUM` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `INV-05` terpasang atas `NOMOR_LAYER` + `BAGIAN_LAYER` di dalam satu versi
- [ ] ketiga uji positif lulus — layer berbagian, MDP dua mata uang, dua baris pemulihan
- [ ] konversi `DEDUCTIBLE_KEDUA` dari teks **melaporkan kegagalannya**; tidak ada baris yang diam-diam menjadi nol
- [ ] `INV-49` tertulis sebagai pengecualian bernama terhadap `INV-47` — bukan dibiarkan terbaca sebagai pelanggaran
- [ ] penanda `F-15` pada `PERSEN_ROL` **masih berbunyi** di `SPEC-MODEL-DATA.md` §10.3 saat tiket ini dinyatakan selesai, atau keputusannya sudah turun dan tiket ini disesuaikan
- [ ] `KTV-A`, `T-6`, `T-3`, `F-15` tercatat di `ASUMSI-CLEAR.md`
