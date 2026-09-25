# E03: Tracer — alur masuk endorsement sampai kasus berfase Policy

**What to build:** Seorang petugas memilih **nomor polis + jenis endorsement + tanggal**, dan sebuah
kasus endorsement tercipta — atau **ditolak dengan alasan spesifik**.

Ini irisan **tertipis yang lengkap** untuk alur masuk. Begitu hijau, premisnya terbukti: endorsement
lahir **langsung di fase polis**, tanpa fase penawaran maupun binding.

`[terverifikasi]` Kasus lahir dari kelas portal menuju kelas kerja berprefiks `EDM-`, dengan **tautan
tiga arah** antara kasus portal, kasus endorsement, dan handle assignment. Tidak ada satu pun
konektor di flow yang menetapkan fase penawaran atau binding.

⛔ **Penolakan adalah keluaran bisnis yang sah**, bukan kegagalan teknis — keduanya harus dapat
dibedakan.

**Asal (Pega).** `Activity\SetValueToEDMWork` langkah 7–13 · `DataTransform\DataToEDM` ·
`Flow\InputAddendumFacultativeIn` konektor `Start2 → Assignment7` · `Activity\SetEdmType`

**Keputusan.** K-049 (Seam 5) · K-029 · K-044

**Blocked by:** `..\08-registry-rules-eval.md` · E01
⛔ Tiket NB itu **belum dikerjakan** — lihat §0 indeks.

**Status:** blocked

- [ ] **Seam 5 `services/endorsement.OpenCase` hidup**
- [ ] Kasus lahir berkelas kerja endorsement, berprefiks **`EDM-`**
- [ ] **Tautan tiga arah** terisi lengkap — kasus portal ↔ kasus endorsement ↔ handle assignment
- [ ] Nomor polis lama dan alasan endorsement terisi
- [ ] Kasus **berfase Policy**; kasus uji negatif: **tidak pernah** melewati fase penawaran atau binding
- [ ] Ditempatkan di workbasket marketing dengan tiket admin polis
- [ ] Jenis endorsement dari domain terkunci (K-029); kode usang **tetap diport**, cabangnya tidak dihapus
- [ ] Penolakan **terbedakan** dari kegagalan sistem, dengan alasan spesifik per gerbang
- [ ] Tiap transisi menyebut rule Pega asalnya dalam komentar (§4.6)
