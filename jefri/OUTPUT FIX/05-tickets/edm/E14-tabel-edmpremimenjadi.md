# E14: Tabel rumus premi-menjadi per lini, jenis endorsement, dan jalur

**What to build:** Properti "premi menjadi" **tidak dihitung dengan satu rumus**. Ia punya belasan
bentuk berbeda yang dipilih menurut **lini bisnis × jenis endorsement × jalur**.

`[terverifikasi]` **Dua puluh dua penugasan** di enam berkas menghasilkan **sebelas rumus berbeda
secara semantik**.

⛔ **JANGAN diseragamkan.** Modelkan sebagai **tabel keputusan**, bukan satu fungsi dengan cabang
tersembunyi. Menyeragamkannya mengubah angka pada lini tertentu.

📌 Pembagian "pakai nilai pembayaran lama vs premi lama" **bukan sifat lini** — ia berlaku hanya pada
satu cabang tertentu, dan di sana hanya **dua lini** yang memakai bentuk berbeda.

**Asal (Pega).** `Activity\CountEndorsementData` · `CountDataEDMElse` · `CountPaymentEdm_Act` ·
`CountPaymentEdmTSIObj_Act` · `ReCountPremiLifeEDM` · `CopyAllObj_ACT` ·
`DataTransform\CountPremiEDM_DT`

**Keputusan.** **K-046** · K-029 · K-048

**Blocked by:** E13

**Status:** blocked

- [ ] Terimplementasi sebagai **tabel** (lini × jenis endorsement × jalur), bukan satu rumus
- [ ] Kedua puluh dua penugasan terwakili; **sebelas bentuk semantik** dapat dibedakan
- [ ] Cabang "batal": **dua lini** memakai nilai pembayaran lama, empat lainnya premi lama
- [ ] Cabang "batal sejak semula": ketujuh lini memakai bentuk yang **sama**
- [ ] **K-046** `K046_EDMPremiMenjadi_TidakSeragam_PerLiniDanEdmType`
- [ ] **K-046** `K046_EdmType1_HanyaLife_LiniLainTidakDihitungUlang` — hanya satu lini yang cabang perhitungannya berjalan untuk jenis endorsement pertama
- [ ] **K-046** `K046_LabelMBU_GerbangIsMarineCargo_IkutiKondisi` — label langkah menyebut satu lini, gerbangnya lini lain; **ikuti gerbang**
- [ ] **K-046** `K046_LabelBatalSejakSemula_Gerbang2_BukanSatu` — label dan gerbang bertentangan di lima lini
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
