---
status: accepted
label: DECIDED
---

# Perhitungan dijalankan sebagai turunan dari data, bukan sebagai akibat penekanan tombol

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Setiap nilai hasil perhitungan adalah **turunan** dari masukannya. Begitu masukan berubah, turunannya ikut berubah. Tidak ada tombol "hitung", tidak ada keadaan setengah jadi, dan tidak ada urutan tindakan yang harus diingat pengguna.

## Consequences

Sistem lama melakukan sebaliknya, dan itu terbukti dari XML: **tidak ada satu pun activity yang memanggil `CountClaimTNP_Act`** — ia dipicu dari Section lewat `<pyActivity>`, 42 rujukan di tiga berkas. Urutan jalannya `CountClaimTNP_Act`, `CountLossAllocation_act`, `GenerateCFS_act`, dan `SaveToOS` tidak ditetapkan di kode mana pun; ia ditentukan kontrol mana yang ditekan.

**Sapuan `MEMORI_PEMAHAMAN.MD` atas urutan kerja, SOP, dan langkah petugas menghasilkan nihil.** Urutan baku itu **tidak terdokumentasi di sumber mana pun** — tidak di XML, tidak di memori. Itu sendiri temuan: urutan yang menentukan hasil perhitungan hanya hidup di kepala petugas, dan akan hilang bersama orangnya.

Keputusan ini menutup jalur itu, bukan memperbaikinya. Konsekuensinya pada rancangan: perhitungan tidak boleh ditempatkan di lapisan tampilan; ia menjadi fungsi murni atas data, dipanggil ulang kapan pun masukannya berubah, dan hasilnya tidak disimpan sebagai keadaan yang dapat menyimpang dari masukannya.

**Keputusan ini mengikat bersama ADR-0011 dan ADR-0012, bukan berdiri sendiri.** Uraian rantai sebab-akibatnya ada di `BLUEPRINT.md` bagian 12: karena hasil bergantung urutan klik, selalu ada klaim yang keluar jalur; karena selalu ada yang keluar jalur, perbaikan termurah adalah menambal klaim itu satu per satu; dan begitulah 29 langkah tambalan menumpuk selama delapan tahun. Mengambil satu dari tiga keputusan tanpa dua lainnya akan mengembalikan pola yang sama dalam bentuk baru.
