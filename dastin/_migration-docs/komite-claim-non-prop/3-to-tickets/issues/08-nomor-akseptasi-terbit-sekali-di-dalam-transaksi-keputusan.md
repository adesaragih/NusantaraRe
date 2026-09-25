---
status: tertahan
---

# 08: Nomor akseptasi terbit sekali, di dalam transaksi keputusan

*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Persetujuan jenjang terakhir atas usulan yang **tidak bersyarat**
menerbitkan nomor akseptasi, di dalam transaksi keputusan yang sama. Persetujuan **bersyarat**
menutup sirkulasi tanpa menerbitkan nomor. Gagal menerbitkan membatalkan seluruh keputusan —
tidak pernah ada keputusan tersimpan tanpa nomor ketika syaratnya terpenuhi.

Pembungkus pencatat di atas sumber urutan Oracle yang ada; constraint keunikan pada kolom
nomor akseptasi yang sudah ada di sisi Klaim.

**Persyaratan:** `S-038`, `S-039`, `S-040`, `S-041`, `S-042`, `S-043`, `S-066`.

**Tidak termasuk:** **nilai ambang periode buku** — **PAGAR-01**. Mekanika penomoran ditulis
penuh; nilainya adalah parameter bernama tanpa isi, dan pembukanya adalah isi badan
`PROC_GENERATE_SEQUENCE_NUMBER` pada basis data berjalan (`INVENTARIS-BUKTI.md` §2.5 baris 4).
Isi baris akseptasi dan bentuk simpannya — **PAGAR-02**, tiket `11`.

**Jalur gagal:** sumber urutan tidak menjawab → `503`, seluruh transaksi batal, nol keputusan
tersimpan · penerbitan kedua atas usulan yang sama → ditolak constraint, pemanggil menerima
hasil penerbitan pertama · permintaan menyunting nomor → `409` · parameter periode buku tidak
terisi → `503`; **nol nomor terbit dengan periode terkaan**.

**Uji:** BARU — syarat terbit, idempotensi dengan mengulang perintah ber-kunci sama,
pembatalan transaksi saat penerbitan gagal · BARU — uji DDL untuk keunikan nomor.

**Menggantikan:** `GetSequenceNumber_SQL` dan `PROC_GENERATE_SEQUENCE_NUMBER` — prosedur lama
**dipakai apa adanya** pada fase 1, dibungkus pencatat (`D-4`) · `GetKodeProdNonLife_SQL` —
kode produksi, **berpagar** PAGAR-01 · idempotensi lewat baca-lalu-tulis — diganti constraint.

**Blocked by:**

- `05` — Akibat pada klaim dihitung di satu tempat dan ditulis dalam transaksi yang sama


**Dasar:** DECIDED(`D-4`, `J-3`, `E-2`, keputusan beku no. 7). Dasar faktualnya bertanda
**TAFSIR (menunggu C-01)** — berkas DDL yang dipegang adalah hasil reverse-engineer, bukan
ekspor langsung.

- [ ] Nomor terbit **bila dan hanya bila** jenjang terakhir menyetujui dan usulan tidak
      bersyarat.
- [ ] Persetujuan bersyarat menutup sirkulasi, **tidak** menerbitkan nomor, dan nilai
      tampilnya berbunyi disetujui bersyarat.
- [ ] Idempotensi ditegakkan **constraint**, bukan baca-lalu-tulis.
- [ ] Percobaan ulang ber-kunci sama mengembalikan **nomor yang sama**, bukan nomor kedua.
- [ ] Nomor tidak pernah disunting sesudah terbit.
- [ ] Nol keputusan final tidak bersyarat tersimpan tanpa nomor.
- [ ] Periode buku ditentukan di **satu lapis**, Oracle; sisi aplikasi tidak menggeser ulang
      dan tidak menambal. Ambangnya parameter bernama, bukan angka di dalam kode.

**Ketidakpastian:** **C-01** — isi badan prosedur penerbit pada basis data berjalan belum
pernah dibandingkan. Menunggui, tidak menahan: mekanika sudah ditetapkan `D-4`; bila isinya
berbeda, yang berubah adalah **pembungkusnya**, bukan bahwa nomor terbit di dalam transaksi.
