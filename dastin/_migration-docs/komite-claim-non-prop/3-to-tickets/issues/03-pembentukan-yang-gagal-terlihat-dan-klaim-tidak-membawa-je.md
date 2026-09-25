---
status: tertahan
---

# 03: Pembentukan yang gagal terlihat, dan klaim tidak membawa jejaknya

*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Analis yang pembentukannya gagal **diberi tahu alasannya** dan klaimnya
tidak berubah — tidak menerima pengenal sirkulasi, tidak berpindah keadaan, tidak menerima
akibat apa pun. Maksud pengajuan tercatat pada kolom tersendiri sejak pengajuan dan tetap
tercatat meski sirkulasinya gagal lahir.

Kolom maksud pengajuan pada tabel klaim; batas transaksi pembentukan; jalur galat yang
terlihat di permukaan.

**Persyaratan:** `S-009`, `S-010`, `S-011`, `S-012`, `S-015`.

**Tidak termasuk:** kolom maksud ini menyentuh **tabel milik sisi Klaim**. Penambahannya
perlu persetujuan pemilik tabelnya; itu di luar papan ini.

**Jalur gagal:** gangguan di tengah transaksi → seluruh transaksi batal, `503` yang dapat
dibedakan dari penolakan, klaim tanpa pengenal sirkulasi dan tanpa niat pemberitahuan
tercatat · validasi tidak terpenuhi → `422` menyebut medan yang gagal, nol berkas, nol
roster, nol tulisan ke klaim · maksud tidak dapat ditulis → pembentukan batal seluruhnya,
klaim tidak berubah sama sekali.

**Uji:** P5-07 (pembentukan gagal karena pesan validasi) · P5-08 (alokasi penyebaran kosong) ·
BARU untuk ketiadaan jalur keluar senyap.

**Menggantikan:** `CreateChildKomiteCNP_Act`·29 — penjaga yang hanya melewati satu langkah
sementara langkah sesudahnya menulis pengenal milik sirkulasi lain lalu menyimpan ·
`CreateChildKomiteCloseNP_Act`·6 — penanda maksud ditulis ke induk pada langkah yang membuat
berkas anak, sebelum komite memutus apa pun · `ProteksiSendKomiteCNP_Act` — gerbang yang
menandai galat dan tidak menghentikan (`F-21`) · `IsFlagError` — penanda yang tidak dibaca
satu rule pun, **dibuang**.

**Blocked by:**

- `02` — Sirkulasi usulan pembayaran lahir lengkap dengan jenjangnya


**Dasar:** DECIDED(`K5-2`, `K5-3`, `K5-7`, `H-5`, ADR-0031). EVIDENCED: `F-21`, `F-22`;
GRILL-05 N-07, N-08.

- [ ] Tiap jalur pembentukan yang berakhir tanpa sirkulasi berakhir dengan galat yang
      **terlihat pengaju** dan tercatat; pemanggil tidak pernah menerima jawaban berhasil
      tanpa sirkulasi.
- [ ] Klaim yang pembentukannya gagal **tidak berubah selain maksudnya** — diperiksa dengan
      membandingkan seluruh barisnya sebelum dan sesudah.
- [ ] Nol klaim membawa pengenal sirkulasi tanpa sirkulasi yang bersesuaian. Ini yang
      mewujudkan invarian `I-11`.
- [ ] Maksud dan akibat menempati **kolom yang berbeda**; maksud tidak pernah menjadi akibat.
- [ ] Validasi **menolak** kasus-guna; nol jalur yang memasang pesan lalu melanjutkan.

**Ketidakpastian:** Penambahan kolom maksud pada tabel klaim menunggu persetujuan pemilik
tabel sisi Klaim. Menunggui, tidak menahan: bila ditolak, yang berubah adalah **tempat maksud
disimpan**, bukan bahwa maksud dan akibat terpisah.
