# -*- coding: utf-8 -*-
u"""
buat-skema-gabungan.py — satu peta untuk DUA berkas yang sudah ada.

Masukan (DIBACA SAJA, tidak pernah ditulis):
    Diagram-Skema-Tabel-NusantaraRe.xlsx     Klaim Life · PremiumList · Claim Prop · Komite Prop
    Diagram-Skema-Tabel-ClaimNonProp.xlsx    Claim Non Prop

Keluaran:
    Diagram-Skema-Tabel-Gabungan.xlsx        4 sheet
    ERD-GABUNGAN.html                        satu halaman, tanpa rujukan luar

Yang dijamin alat ini
---------------------
1. Garis penghubung DIGAMBAR, bukan disiratkan. Di Excel dengan border sel
   (vertikal = border kiri, horizontal = border atas); di HTML dengan <path>.
2. TIDAK ADA garis yang tertutup kotak. Jaminannya geometris, bukan kebetulan:
   - kotak hanya menempati kolom >= KOL_LAJUR; koridor tali (kolom 2..KOL_LAJUR-2)
     tidak pernah ditempati kotak;
   - di dalam lajur, batang pohon sebuah induk berada di kolom induk+1, dan
     seluruh keturunannya digambar berurutan (DFS) mulai kolom induk+3, sehingga
     kolom induk+1 kosong sepanjang rentang baris batang itu;
   - tiap tali panjang mendapat kolom jalurnya sendiri lewat pewarnaan graf
     selang, jadi dua tali tidak pernah berbagi kolom pada baris yang sama.
3. Tiap angka yang dilaporkan dihitung dari daftar di berkas ini, bukan diingat.
"""

import collections
import datetime
import io
import os
import re
import sys

import openpyxl
from openpyxl.styles import Alignment, Border, Font, PatternFill, Side

BASIS = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
ERD = os.path.join(BASIS, '4-erd-dan-tabel-datar')
RUJUK_NURE = os.path.join(ERD, 'Diagram-Skema-Tabel-NusantaraRe.xlsx')
RUJUK_CNP = os.path.join(ERD, 'Diagram-Skema-Tabel-ClaimNonProp.xlsx')
KELUAR_XLSX = os.path.join(ERD, 'Diagram-Skema-Tabel-Gabungan.xlsx')
KELUAR_HTML = os.path.join(ERD, 'ERD-GABUNGAN.html')
TANGGAL = datetime.date.today().isoformat()

LAPOR = []


def lapor(baris):
    LAPOR.append(baris)
    print(baris)


# =====================================================================
# 1. Daftar tabel — gabungan kedua berkas, satu baris per tabel
# =====================================================================
# lini    : 'inti' 'komite' 'nonprop' 'prop' 'life' 'polis' 'migrasi'
# berkas  : 'NURE' 'CNP' 'KEDUANYA'   — di berkas mana tabel ini sudah ada
# induk   : nama induk di dalam lajur yang sama (None = akar lajur)
T = []
URUT = {}


def tabel(nama, lini, berkas, kard, induk, lajur, baris, catatan=u''):
    if nama in URUT:
        raise SystemExit(u'tabel ganda: %s' % nama)
    URUT[nama] = len(T)
    T.append({'nama': nama, 'lini': lini, 'berkas': berkas, 'kard': kard,
              'induk': induk, 'lajur': lajur, 'baris': list(baris),
              'catatan': catatan})


# --- INTI — tulang punggung, dipakai SEMUA lini klaim ----------------
tabel('T_WORK_CLAIM', 'inti', 'KEDUANYA', u'akar', None, 'inti', [
    u'PK   ID           CLM-/KMT- (FAC)   CLMP-/TKMT- (PROP)',
    u'                  CLMNP-/KMTNP- (NONPROP)   CLMLF-/KMTLF- (LIFE)',
    u'FK   COVER_KEY  →  T_WORK_CLAIM.ID     self · nullable',
    u'     LINI menandai lini pemilik baris — satu tabel, empat lini',
], u'akar lintas-lini · satu baris per work object')

tabel('T_GENERAL_CLAIM', 'inti', 'KEDUANYA', u'1:1', 'T_WORK_CLAIM', 'inti', [
    u'PK   ID   =   T_WORK_CLAIM.ID       SHARED PK',
    u'     tidak ada kolom FK terpisah',
    u'     kolom khas satu lini WAJIB nullable',
], u'induk bersama LIFE, PROP, FAC IN dan NONPROP')

tabel('T_VIEW_SUGGEST', 'inti', 'KEDUANYA', u'1:N', 'T_GENERAL_CLAIM', 'inti', [
    u'PK   ID',
    u'FK   CLAIM_ID         →  T_GENERAL_CLAIM.ID    nullable · CASCADE',
    u'FK   PREMIUM_LIST_ID  →  T_PREMIUM_LIST.ID     nullable · CASCADE',
    u'CHECK ( (PREMIUM_LIST_ID IS NULL) <> (CLAIM_ID IS NULL) )',
], u'⭐ DUA INDUK — satu-satunya tabel yang menyambung lajur KLAIM dan lajur POLIS')

tabel('DOCUMENT_CLAIM', 'inti', 'KEDUANYA', u'1:N', 'T_GENERAL_CLAIM', 'inti', [
    u'PK   ID',
    u'FK   CLAIM_ID  →  T_GENERAL_CLAIM.ID        pola PROP, FAC IN, NONPROP',
    u'FK   PREMIUM_LIST_DETAIL_ID  →  T_CLAIMLF_PREMIUMLIST_DETAIL.ID   pola LIFE',
    u'⚠ induknya BERBEDA menurut lini — lihat sheet Tali Penghubung',
    u'⚠ daftar kolom TIDAK terbaca dari korpus — [data DBA]',
], u'kelas Pega ASM-FW-GCNMFW-Int-DOCUMENT_CLAIM')

tabel('T_GENERAL_KOMITE', 'komite', 'KEDUANYA', u'1:1', 'T_WORK_CLAIM', 'inti', [
    u'PK   ID   =   T_WORK_CLAIM.ID  (baris komite)     SHARED PK',
    u'FK   ADJUSTMENT_ID  →  adjustment lini ybs   NOT NULL · UNIK',
    u'⚠ TANPA REFERENCES — DUA sasaran sesudah 20-09-2026:',
    u'     T_CLAIM_ADJUSTMENT (PROP · FAC IN · NONPROP)  atau  T_CLAIMLF_ADJUSTMENT (LIFE)',
    u'     7 kolom · CLOSE_FILE dan RESERVED_CLAIM hanya diisi jalur PROP',
], u'lintas-lini · tujuh kolom sesudah 19-09-2026')

tabel('T_KOMITE_KOMITELIST', 'komite', 'KEDUANYA', u'1:N', 'T_GENERAL_KOMITE', 'inti', [
    u'PK   ID',
    u'FK   DATA_KOMITE_ID  →  T_GENERAL_KOMITE.ID      CASCADE',
    u'     9 kolom · satu baris = satu penyetuju',
], u'dipakai bersama seluruh lini komite')

# --- T_CLAIM_ yang DIPAKAI BERSAMA lebih dari satu lini ---------------
# Sesudah [keputusan work owner 20-09-2026] awalan T_CLAIMP_ diganti T_CLAIM_,
# dan tabel-tabel ini tidak lagi milik satu lini.
BERSAMA = [
    ('T_CLAIM_CLAIM_AMOUNT', 'T_GENERAL_CLAIM', 'CLAIM_ID', None, None,
     u'PROP · NONPROP', u'← ClaimData.ListClaimAmount[] · 15 field', ()),
    ('T_CLAIM_INTEREST', 'T_GENERAL_CLAIM', 'CLAIM_ID', None, None,
     u'PROP · NONPROP', u'← ClaimData.InterestList[] · n=39', ()),
    ('T_CLAIM_ESTIMATION', 'T_GENERAL_CLAIM', 'CLAIM_ID',
     'T_CLAIM_OBJECT_ITEM', 'OBJECT_ITEM_ID',
     u'PROP · FAC IN · NONPROP', u'← ClaimData.EstimationList[] · 7 field', ()),
    ('T_CLAIM_SPREADING', 'T_GENERAL_CLAIM', 'CLAIM_ID',
     'T_CLAIM_OBJECT_ITEM', 'OBJECT_ITEM_ID',
     u'PROP · FAC IN · NONPROP', u'← ClaimData.SpreadingClaim[] · + TREATY_ID', ()),
    ('T_CLAIM_BREAK_QS', 'T_GENERAL_CLAIM', 'CLAIM_ID',
     'T_CLAIM_OBJECT_ITEM', 'OBJECT_ITEM_ID',
     u'PROP · FAC IN · NONPROP',
     u'← ClaimData.SpreadingBreakQS[] · n=28 · SEJAJAR, bukan anak _SPREADING', ()),
    ('T_CLAIM_OBJECT', 'T_GENERAL_CLAIM', 'CLAIM_ID', None, None,
     u'FAC IN · NONPROP', u'← ClaimData.ObjectList[]',
     (u'⚠ di PROP ditandai "NOL penulis · residu save-as" — asalnya di FAC dan NONPROP',)),
    ('T_CLAIM_OBJECT_ITEM', 'T_CLAIM_OBJECT', 'OBJECT_ID', None, None,
     u'FAC IN · NONPROP', u'← ObjectList[].ObjectItemList[]', ()),
    ('T_CLAIM_ADJUSTMENT', 'T_GENERAL_CLAIM', 'CLAIM_ID',
     'T_CLAIM_OBJECT_ITEM', 'OBJECT_ITEM_ID',
     u'PROP · FAC IN · NONPROP',
     u'← ClaimData.AdjustmentList[] · 47 field · ★ unit transaksi',
     (u'FK   KOMITE_ID  →  T_WORK_CLAIM.ID     1:1 · nullable · UNIK',)),
    ('T_CLAIM_ADJ_SPREADING', 'T_CLAIM_ADJUSTMENT', 'ADJUSTMENT_ID', None, None,
     u'PROP · FAC IN · NONPROP',
     u'← AdjustmentList[].SpreadingAdjustment[] ═ salinan beku', ()),
    ('T_CLAIM_ADJ_QUOTA_SHARE', 'T_CLAIM_ADJUSTMENT', 'ADJUSTMENT_ID', None, None,
     u'PROP · FAC IN · NONPROP',
     u'← AdjustmentList[].SpreadingQuotaShare[] ═ salinan beku', ()),
    ('T_CLAIM_ADJ_LOSS_ALLOCATION', 'T_CLAIM_ADJUSTMENT', 'ADJUSTMENT_ID', None, None,
     u'PROP · NONPROP', u'← AdjustmentList[].LossAllocation[] ═ potret alokasi', ()),
    ('T_CLAIM_FAC_RETRO', 'T_GENERAL_CLAIM', 'CLAIM_ID',
     'T_CLAIM_ADJUSTMENT', 'ADJUSTMENT_ID',
     u'PROP · FAC IN', u'← ClaimData.FacRetroList',
     (u'⚠ TIDAK dipakai NONPROP — di sana retro menempel pada CNPSpreadLoss[],',
      u'   tabelnya T_CLAIM_RETRO di lajur bawah. Nama mirip, halaman asal BERBEDA.')),
]
for nm, ind, kn, ind2, kn2, pakai, asal, tam in BERSAMA:
    _brk = 'NURE' if u'NONPROP' not in pakai else 'KEDUANYA'
    isi = [u'PK   ID', u'FK   %s  →  %s.ID     CASCADE' % (kn, ind)]
    if ind2:
        isi.append(u'FK   %s  →  %s.ID   nullable · jalur FAC IN' % (kn2, ind2))
        isi.append(u'CHECK  tepat satu dari %s / %s terisi' % (kn, kn2))
    isi.extend(tam)
    isi.append(asal)
    tabel(nm, 'bersama', _brk, u'1:N', ind, 'bersama', isi,
          u'dipakai lini: %s' % pakai)

# --- KHAS CLAIM NON PROP ---------------------------------------------
NP = [
    ('T_CLAIM_SPREADING_RISK', 'T_GENERAL_CLAIM', 'CLAIM_ID',
     u'← ClaimData.SpreadingRisk[] · 29 field · ★ inti alokasi XOL',
     (u'⚠ di PROP SpreadingRisk BUKAN tabel — bacaan master treaty, bukan data',)),
    ('T_CLAIM_RETENTION_CEDANT', 'T_GENERAL_CLAIM', 'CLAIM_ID',
     u'← baris "UR" di SpreadingRisk[] — dipisah [ADR-0010]', ()),
    ('T_CLAIM_SPREAD_LOSS', 'T_GENERAL_CLAIM', 'CLAIM_ID',
     u'← ClaimData.CNPSpreadLoss[] · 8 field · pembagian kerugian per COB', ()),
    ('T_CLAIM_RETRO', 'T_CLAIM_SPREAD_LOSS', 'SPREAD_LOSS_ID',
     u'← CNPSpreadLoss[].RetroList[]',
     (u'⚠ BUKAN T_CLAIM_FAC_RETRO — halaman asalnya berbeda',
      u'⚠ belum ada padanannya di 22 tabel ddl-usulan/')),
    ('T_CLAIM_REINSTATEMENT', 'T_GENERAL_CLAIM', 'CLAIM_ID',
     u'← ClaimData.ReinstatementList[] · 12 field · premi pemulihan', ()),
    ('T_CLAIM_RECEIVER', 'T_GENERAL_CLAIM', 'CLAIM_ID',
     u'← ClaimData.ReceiverClaim[] · 13 field · sepasang rekening', ()),
    ('T_CLAIM_DATE_GUARD', 'T_GENERAL_CLAIM', 'CLAIM_ID',
     u'TIDAK DITEMUKAN di pohon — penjaga batas tanggal', ()),
    ('T_CLAIM_OUTBOUND_ARCHIVE', 'T_GENERAL_CLAIM', 'CLAIM_ID',
     u'← InputParamOs / InputParamOsCNP → OS_AKSEPTASI_KLAIM',
     (u'hanya SELECT dan INSERT · tanpa UPDATE, tanpa DELETE',)),
    ('T_CLAIM_ADJ_LAYER', 'T_CLAIM_ADJUSTMENT', 'ADJUSTMENT_ID',
     u'← AdjustmentList[].CNPLayerList[] · XOL · XOLID',
     (u'⚠ induknya di lajur T_CLAIM_ BERSAMA — lihat tali penghubung',
      u'⚠ belum ada padanannya di 22 tabel ddl-usulan/')),
    ('T_CLAIM_ADJ_LAYER_CURRENCY', 'T_CLAIM_ADJ_LAYER', 'LAYER_ID',
     u'← CNPLayerList[].CNPCurrencyList[] · 21 field · KEDALAMAN 4',
     (u'⚠ belum ada padanannya di 22 tabel ddl-usulan/',)),
]
for nm, ind, kn, asal, tam in NP:
    isi = [u'PK   ID', u'FK   %s  →  %s.ID     CASCADE' % (kn, ind)]
    isi.extend(tam)
    isi.append(asal)
    tabel(nm, 'np', 'CNP', u'1:N', ind, 'np', isi)

for nm, pad, ket in [
        ('T_CLAIM_CLOSING_DATE', u'TUTUP_BUKU', u'data bertanggal berlaku, lingkup global'),
        ('T_CLAIM_RATE', u'TARIF_BERLAKU', u'tarif bertanggal berlaku, tanpa induk'),
        ('T_CLAIM_CORRECTION', u'KOREKSI_NILAI', u'menunjuk lewat TABEL_SASARAN + ID_SASARAN')]:
    tabel(nm, 'np', 'CNP', u'bebas', None, 'np', [
        u'PK   ID', u'     TANPA FK — tidak menempel pada satu klaim',
        u'≡ %s' % pad, ket])

tabel('T_CLAIM_MIG_CORRELATION', 'migrasi', 'CNP', u'bebas', None, 'np', [
    u'PK   ID',
    u'     TABEL_TUJUAN + ID_TUJUAN — menunjuk balik TANPA FK',
    u'jembatan ke sistem lama · pzInsKey · pyID · CASEID · berumur terbatas'])
tabel('T_CLAIM_MIG_LANDING', 'migrasi', 'CNP', u'bebas', None, 'np', [
    u'PK   ID',
    u'     DATA_JSON apa adanya — buang seluruh kunci px* py* pz*'])
tabel('T_CLAIM_MIG_REJECTED', 'migrasi', 'CNP', u'1:N', 'T_CLAIM_MIG_LANDING', 'np', [
    u'PK   ID',
    u'FK   LANDING_ID  →  T_CLAIM_MIG_LANDING.ID     nullable',
    u'nilai yang tidak dapat diurai — tercatat, tidak dibulatkan'])

# --- CLAIM LIFE — T_CLAIMLF_, TIDAK disentuh keputusan 20-09-2026 -----
LIFE = [
    ('T_CLAIMLF_PREMIUMLIST_DETAIL', 'T_GENERAL_CLAIM', 'CLAIM_ID',
     u'induk langsung seluruh pohon klaim LIFE', ()),
    ('T_CLAIMLF_ADJUSTMENT', 'T_CLAIMLF_PREMIUMLIST_DETAIL', 'PREMIUM_LIST_DETAIL_ID',
     u'⚠ nama kunci tidak cocok induk — dilaporkan, tidak diubah',
     (u'FK   KOMITE_ID  →  T_WORK_CLAIM.ID     1:1 · nullable · UNIK',)),
    ('T_CLAIMLF_ADJUSTMENT_SPREADING', 'T_CLAIMLF_ADJUSTMENT', 'ADJUSTMENT_ID', u'', ()),
    ('T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO', 'T_CLAIMLF_ADJUSTMENT_SPREADING',
     'SPREADING_ID', u'', ()),
]
for nm, ind, kn, asal, tam in LIFE:
    isi = [u'PK   ID', u'FK   %s  →  %s.ID   CASCADE' % (kn, ind)]
    isi.extend(tam)
    if asal:
        isi.append(asal)
    tabel(nm, 'life', 'NURE', u'1:N', ind, 'life', isi)

# --- POLIS / PREMIUMLIST ---------------------------------------------
tabel('T_WORK_POLIS', 'polis', 'NURE', u'akar', None, 'polis', [
    u'PK   ID',
    u'     baris NB dan baris EDM sama-sama di sini — TIDAK saling menunjuk',
], u'akar lajur POLIS — sejajar T_WORK_CLAIM, bukan anaknya')
tabel('T_PREMIUM_LIST', 'polis', 'NURE', u'1:1', 'T_WORK_POLIS', 'polis', [
    u'PK   ID   =   T_WORK_POLIS.ID      SHARED PK',
    u'     WORK_POLIS_ID DIBUANG 18-09-2026',
    u'     versi berjalan = PRODKE terbesar · NB = 1 · EDM = 2,3,…',
])
tabel('T_PREMIUM_LIST_SUMMARY', 'polis', 'NURE', u'1:N', 'T_PREMIUM_LIST', 'polis', [
    u'PK   ID', u'FK   PREMIUM_LIST_ID  →  T_PREMIUM_LIST.ID    CASCADE',
    u'     TIDAK disalin antar versi'])
tabel('T_PREMIUM_LIST_DETAIL', 'polis', 'NURE', u'1:N', 'T_PREMIUM_LIST', 'polis', [
    u'PK   ID', u'FK   PREMIUM_LIST_ID  →  T_PREMIUM_LIST.ID    CASCADE',
    u'FK   PARENT_ID  →  T_PREMIUM_LIST_DETAIL.ID   self · nullable'])
tabel('T_PREMIUM_LIST_SPREADING', 'polis', 'NURE', u'1:N', 'T_PREMIUM_LIST_DETAIL', 'polis', [
    u'PK   ID', u'FK   DETAIL_ID  →  T_PREMIUM_LIST_DETAIL.ID    CASCADE'])
tabel('T_PREMIUM_LIST_SPREADING_RETRO', 'polis', 'NURE', u'1:N',
      'T_PREMIUM_LIST_SPREADING', 'polis', [
          u'PK   ID', u'FK   SPREADING_ID  →  T_PREMIUM_LIST_SPREADING.ID   CASCADE'])
tabel('DOCUMENT_POLIS', 'polis', 'NURE', u'1:N', 'T_PREMIUM_LIST', 'polis', [
    u'PK   ID', u'FK   ⚠ [data DBA] — kolom penyambung belum terbaca dari korpus',
    u'kelas Pega  ASM-FW-GISFW-Int-DOCUMENT_POLIS'],
      u'lintas-lini')

PETA = dict((t['nama'], t) for t in T)

# =====================================================================
# 2. Tali penghubung — yang MENYAMBUNG lajur satu ke lajur lain
# =====================================================================
TALI = [
    ('T_GENERAL_CLAIM', '@bersama', 'tali', u'CLAIM_ID',
     u'12 tabel T_CLAIM_ bersama bergantung pada induk yang sama'),
    ('T_GENERAL_CLAIM', '@np', 'tali', u'CLAIM_ID',
     u'8 tabel khas Non Prop bergantung pada induk yang sama'),
    ('T_GENERAL_CLAIM', '@life', 'tali', u'CLAIM_ID',
     u'pintu masuk pohon LIFE — lewat T_CLAIMLF_PREMIUMLIST_DETAIL'),
    ('T_CLAIM_ADJUSTMENT', '@np', 'tali', u'ADJUSTMENT_ID',
     u'T_CLAIM_ADJ_LAYER khas Non Prop menggantung pada adjustment bersama'),
    ('T_PREMIUM_LIST', 'T_VIEW_SUGGEST', 'tali', u'PREMIUM_LIST_ID',
     u'⭐ satu-satunya kolom yang menyambung lajur POLIS ke lajur KLAIM'),
    ('T_CLAIM_ADJUSTMENT', 'T_GENERAL_KOMITE', 'silang', u'ADJUSTMENT_ID',
     u'penutup lingkar PROP + FAC IN + NONPROP — satu tabel sesudah 20-09-2026'),
    ('T_CLAIMLF_ADJUSTMENT', 'T_GENERAL_KOMITE', 'silang', u'ADJUSTMENT_ID',
     u'penutup lingkar LIFE · sasaran KEDUA dari kolom yang sama'),
    ('T_CLAIM_ADJUSTMENT', 'T_WORK_CLAIM', 'silang', u'KOMITE_ID',
     u'pintasan tampilan PROP + FAC IN + NONPROP · nullable · UNIK'),
    ('T_CLAIMLF_ADJUSTMENT', 'T_WORK_CLAIM', 'silang', u'KOMITE_ID',
     u'pintasan tampilan LIFE · nullable · UNIK'),
    ('T_CLAIMLF_PREMIUMLIST_DETAIL', 'DOCUMENT_CLAIM', 'silang', u'PREMIUM_LIST_DETAIL_ID',
     u'⚠ induk DOCUMENT_CLAIM di LIFE BERBEDA dari lini lain'),
]

# lajur: (kunci, judul, teks penambat ke lajur di atasnya)
LAJUR = [
    ('inti', u'INTI — TULANG PUNGGUNG BERSAMA · dipakai SEMUA lini', None),
    ('bersama', u'T_CLAIM_ DIPAKAI BERSAMA — PROP · FAC IN · NON PROP '
                u'[keputusan work owner 20-09-2026]',
     u'⇦  T_GENERAL_CLAIM   · kotaknya ada di lajur INTI di atas — tidak digambar dua kali'),
    ('np', u'KHAS CLAIM NON PROP — tidak ada padanannya di lini lain',
     u'⇦  induknya ada di lajur di atas — lihat baris FK tiap kotak'),
    ('life', u'CLAIM LIFE — awalan T_CLAIMLF_ TIDAK ikut diganti',
     u'⇦  T_GENERAL_CLAIM   · kotaknya ada di lajur INTI di atas — tidak digambar dua kali'),
    ('polis', u'POLIS / PREMIUMLIST — sheet "PremiumList NB + EDM"', None),
]

# =====================================================================
# 3. Tata letak — satu kali, dipakai Excel DAN HTML
# =====================================================================
LEBAR_KOTAK = 46          # kolom sempit per kotak
INDEN = 3                 # kolom per tingkat kedalaman
KOL_LAJUR = 17            # kolom pertama yang boleh ditempati kotak
KOR_AWAL, KOR_AKHIR = 2, 15   # koridor tali penghubung — TIDAK PERNAH ditempati kotak
SELA = 1                  # baris kosong antar kotak

KOTAK = {}                # nama -> dict(r0, c0, tinggi, lebar)
PENAMBAT = {}             # '@lajur' -> dict(r0, c0, tinggi, lebar)
PITA = []                 # (r, teks, warna)
SEGMEN = []               # (jenis, 'V'|'H', r1, c1, r2, c2)
TEKS_KORIDOR = []         # (r, c, teks, warna)


def anak_dari(induk, lajur):
    return [t['nama'] for t in T if t['lajur'] == lajur and t['induk'] == induk]


def gambar_subpohon(nama, r, tingkat):
    """DFS pre-order. Kembalikan (baris_sesudah, daftar_baris_kepala_anak)."""
    c0 = KOL_LAJUR + tingkat * INDEN
    t = PETA[nama]
    tinggi = 1 + len(t['baris']) + (1 if t['catatan'] else 0)
    KOTAK[nama] = {'r0': r, 'c0': c0, 'tinggi': tinggi, 'lebar': LEBAR_KOTAK - tingkat * INDEN,
                   'tingkat': tingkat}
    r += tinggi + SELA
    anak = anak_dari(nama, t['lajur'])
    baris_anak = []
    for a in anak:
        baris_anak.append(r)
        r, _ = gambar_subpohon(a, r, tingkat + 1)
    if anak:
        # batang pohon: kolom induk+1, dari bawah kotak induk sampai anak terakhir
        batang = c0 + 1
        atas = KOTAK[nama]['r0'] + tinggi
        bawah = baris_anak[-1] + 1
        SEGMEN.append(('pohon', 'V', atas, batang, bawah, batang))
        for ra in baris_anak:
            SEGMEN.append(('pohon', 'H', ra + 1, batang, ra + 1, c0 + INDEN))
    return r, baris_anak


baris = 4
for kunci, judul, penambat_ke in LAJUR:
    PITA.append((baris, judul, kunci))
    baris += 2
    if penambat_ke:
        lb = LEBAR_KOTAK
        PENAMBAT['@' + kunci] = {'r0': baris, 'c0': KOL_LAJUR, 'tinggi': 1, 'lebar': lb,
                                 'teks': penambat_ke, 'tingkat': 0}
        baris += 1 + SELA
        akar_lajur = [t['nama'] for t in T
                      if t['lajur'] == kunci and (t['induk'] is None or t['induk'] not in
                                                  [x['nama'] for x in T if x['lajur'] == kunci])]
        baris_anak = []
        for nm in akar_lajur:
            baris_anak.append(baris)
            baris, _ = gambar_subpohon(nm, baris, 1)
        batang = KOL_LAJUR + 1
        SEGMEN.append(('pohon', 'V', PENAMBAT['@' + kunci]['r0'] + 1, batang,
                       baris_anak[-1] + 1, batang))
        for ra in baris_anak:
            SEGMEN.append(('pohon', 'H', ra + 1, batang, ra + 1, KOL_LAJUR + INDEN))
    else:
        akar_lajur = [t['nama'] for t in T
                      if t['lajur'] == kunci and t['induk'] is None]
        for nm in akar_lajur:
            baris, _ = gambar_subpohon(nm, baris, 0)
    baris += 2

BARIS_AKHIR_PETA = baris


# --- tali panjang: satu kolom jalur sendiri per tali ------------------
def kotak_atau_penambat(nama):
    return PENAMBAT[nama] if nama.startswith('@') else KOTAK[nama]


pakai_sumber = collections.Counter()
pakai_tujuan = collections.Counter()
rencana = []
for sm, tj, jenis, label, ket in TALI:
    ks, kt = kotak_atau_penambat(sm), kotak_atau_penambat(tj)
    os_ = min(pakai_sumber[sm], max(0, ks['tinggi'] - 1))
    ot_ = min(pakai_tujuan[tj], max(0, kt['tinggi'] - 1))
    pakai_sumber[sm] += 1
    pakai_tujuan[tj] += 1
    rs, rt = ks['r0'] + os_ + (1 if ks['tinggi'] > 1 else 0), \
             kt['r0'] + ot_ + (1 if kt['tinggi'] > 1 else 0)
    rencana.append({'sm': sm, 'tj': tj, 'jenis': jenis, 'label': label, 'ket': ket,
                    'rs': rs, 'rt': rt, 'cs': ks['c0'], 'ct': kt['c0'],
                    'atas': min(rs, rt), 'bawah': max(rs, rt)})

# pewarnaan graf selang -> kolom jalur; jalur 0 paling dekat lajur
jalur_akhir = []
for p in sorted(rencana, key=lambda x: (x['atas'], x['bawah'])):
    for i, akhir in enumerate(jalur_akhir):
        if akhir < p['atas']:
            p['jalur'] = i
            jalur_akhir[i] = p['bawah']
            break
    else:
        p['jalur'] = len(jalur_akhir)
        jalur_akhir.append(p['bawah'])
JUM_JALUR = len(jalur_akhir)
if KOR_AKHIR - JUM_JALUR + 1 < KOR_AWAL:
    raise SystemExit(u'koridor kurang lebar: perlu %d jalur' % JUM_JALUR)

for p in rencana:
    tc = KOR_AKHIR - p['jalur']
    p['tc'] = tc
    SEGMEN.append((p['jenis'], 'H', p['rs'], tc, p['rs'], p['cs']))
    SEGMEN.append((p['jenis'], 'V', min(p['rs'], p['rt']), tc, max(p['rs'], p['rt']), tc))
    SEGMEN.append((p['jenis'], 'H', p['rt'], tc, p['rt'], p['ct']))
    TEKS_KORIDOR.append((p['rt'], tc + 1, u'▶ ' + p['label'], p['jenis']))

# --- periksa jaminan: tidak ada segmen yang jatuh di dalam kotak ------
terpakai = set()
for nm, k in list(KOTAK.items()) + list(PENAMBAT.items()):
    for r in range(k['r0'], k['r0'] + k['tinggi']):
        for c in range(k['c0'], k['c0'] + k['lebar']):
            terpakai.add((r, c))

bentrok = []
for jenis, arah, r1, c1, r2, c2 in SEGMEN:
    if arah == 'V':
        titik = [(r, c1) for r in range(r1, r2)]
    else:
        titik = [(r1, c) for c in range(min(c1, c2), max(c1, c2))]
    for tt in titik:
        # sebuah titik (r,c) menggambar garis di TEPI KIRI / TEPI ATAS sel itu,
        # jadi yang menutupinya hanyalah sel (r,c) sendiri bila ia bagian kotak
        if tt in terpakai:
            bentrok.append((jenis, arah, tt))

# =====================================================================
# 4. Excel
# =====================================================================
LEBAR_KOLOM = 2.3
TINGGI_BARIS = 15.0
ISI = {'inti': '1F4E79', 'komite': '7030A0', 'bersama': 'C55A11', 'np': '2E75B6',
       'life': 'BF8F00', 'polis': '31859C', 'migrasi': '808080'}
ISI_PK = 'FFF2CC'
ISI_BADAN = 'EDF3F9'
ISI_PENAMBAT = 'FCE4D6'
TINTA_REDUP = '7F7F7F'
TINTA_AWAS = 'C00000'
TINTA_FK = '1F4E79'
WARNA_GARIS = {'pohon': ('medium', '1F4E79'), 'tali': ('thick', 'C55A11'),
               'silang': ('mediumDashed', 'C00000')}

wb = openpyxl.Workbook()
ws = wb.active
ws.title = 'Peta Gabungan'
ws.sheet_view.showGridLines = False
MAKS_KOL = KOL_LAJUR + LEBAR_KOTAK + 4
for i in range(1, MAKS_KOL + 1):
    ws.column_dimensions[openpyxl.utils.get_column_letter(i)].width = LEBAR_KOLOM

BATAS = {}


def sisi(r, c, **ganti):
    kini = BATAS.setdefault((r, c), {})
    kini.update(ganti)
    b = {}
    for k, v in kini.items():
        gaya, warna = v
        b[k] = Side(style=gaya, color=warna)
    ws.cell(r, c).border = Border(**b)


def merge_tulis(r, c, teks, isi, tinta, ukuran, tebal, rata, lebar,
                kiri=False, kanan=False, atas=False, bawah=False, warna_tepi='1F4E79'):
    if lebar > 1:
        ws.merge_cells(start_row=r, start_column=c, end_row=r, end_column=c + lebar - 1)
    sel = ws.cell(r, c, teks)
    sel.font = Font(name='Arial', size=ukuran, bold=tebal, color=tinta)
    if isi:
        sel.fill = PatternFill('solid', fgColor=isi)
    sel.alignment = Alignment(horizontal=rata, vertical='center')
    for k in range(c, c + lebar):
        g = {}
        if kiri and k == c:
            g['left'] = ('medium', warna_tepi)
        if kanan and k == c + lebar - 1:
            g['right'] = ('medium', warna_tepi)
        if atas:
            g['top'] = ('medium', warna_tepi)
        if bawah:
            g['bottom'] = ('medium', warna_tepi)
        if g:
            sisi(r, k, **g)
    ws.row_dimensions[r].height = TINGGI_BARIS
    return r + 1


def gaya_baris(teks):
    s = teks.lstrip()
    if s.startswith('PK'):
        return ISI_PK, '000000', True
    if s.startswith('FK'):
        return ISI_BADAN, TINTA_FK, True
    if s.startswith((u'⚠', u'⭐', 'CHECK')):
        return ISI_BADAN, TINTA_AWAS, False
    if s.startswith((u'←', u'≡')):
        return ISI_BADAN, TINTA_REDUP, False
    return ISI_BADAN, '000000', False


merge_tulis(1, 2, u'PETA GABUNGAN — CLAIM NON PROP · CLAIM PROP · CLAIM LIFE · '
                  u'KOMITE · POLIS / PREMIUMLIST', None, '000000', 14, True, 'left', 60)
merge_tulis(2, 2, u'%d tabel · %d relasi pohon · %d tali penghubung · hanya PK dan FK · '
                  u'menyatukan Diagram-Skema-Tabel-NusantaraRe.xlsx dan '
                  u'Diagram-Skema-Tabel-ClaimNonProp.xlsx · dibangkitkan %s'
            % (len(T), sum(1 for t in T if t['induk']), len(TALI), TANGGAL),
            None, '7F7F7F', 9, False, 'left', 76)

for r, judul, kunci in PITA:
    merge_tulis(r, KOL_LAJUR, u'  ' + judul, ISI[kunci], 'FFFFFF', 11, True, 'left',
                LEBAR_KOTAK, True, True, True, True, warna_tepi=ISI[kunci])

for nm, k in KOTAK.items():
    t = PETA[nm]
    kepala = u'%s        %s' % (nm, t['kard'])
    if t['berkas'] == 'KEDUANYA':
        kepala += u'        ⬚ ada di KEDUA berkas'
    merge_tulis(k['r0'], k['c0'], kepala, ISI[t['lini']], 'FFFFFF', 9, True, 'center',
                k['lebar'], True, True, True, False)
    isi_baris = list(t['baris']) + ([t['catatan']] if t['catatan'] else [])
    rr = k['r0'] + 1
    for i, teks in enumerate(isi_baris):
        f, tt, tb = gaya_baris(teks)
        akhir = (i == len(isi_baris) - 1)
        merge_tulis(rr, k['c0'], u'  ' + teks, f, tt, 8, tb, 'left', k['lebar'],
                    True, True, False, akhir)
        rr += 1

for nm, k in PENAMBAT.items():
    merge_tulis(k['r0'], k['c0'], u'  ' + k['teks'], ISI_PENAMBAT, 'C55A11', 9, True, 'left',
                k['lebar'], True, True, True, True, warna_tepi='C55A11')

for jenis, arah, r1, c1, r2, c2 in SEGMEN:
    gaya, warna = WARNA_GARIS[jenis]
    if arah == 'V':
        for r in range(r1, r2):
            sisi(r, c1, left=(gaya, warna))
            ws.row_dimensions[r].height = TINGGI_BARIS
    else:
        for c in range(min(c1, c2), max(c1, c2)):
            sisi(r1, c, top=(gaya, warna))

for r, c, teks, jenis in TEKS_KORIDOR:
    sel = ws.cell(r, c, teks)
    sel.font = Font(name='Arial', size=7, bold=True,
                    color='C55A11' if jenis == 'tali' else 'C00000')

# --- legenda ---------------------------------------------------------
rl = BARIS_AKHIR_PETA + 1
rl = merge_tulis(rl, KOL_LAJUR, u'LEGENDA', None, '000000', 11, True, 'left', 20)
for warna, lbl, ket in (
        (ISI['inti'], u'Kotak biru tua', u'INTI — dipakai seluruh lini klaim'),
        (ISI['bersama'], u'Kotak jingga', u'T_CLAIM_ dipakai bersama PROP · FAC IN · NON PROP'),
        (ISI['np'], u'Kotak biru', u'khas Claim Non Prop — tanpa padanan di lini lain'),
        (ISI['life'], u'Kotak emas', u'Claim Life — T_CLAIMLF_'),
        (ISI['polis'], u'Kotak toska', u'Polis / PremiumList — T_PREMIUM_LIST'),
        (ISI['komite'], u'Kotak ungu', u'Komite — lintas-lini'),
        (ISI['migrasi'], u'Kotak abu', u'tabel migrasi — berumur terbatas'),
        (ISI_PK, u'Baris kuning', u'kunci utama (PK)'),
        (ISI_PENAMBAT, u'Pita jingga', u'penambat — tabel INTI, digambar sekali saja di atas'),
        (None, u'Garis biru penuh', u'hubungan induk → anak di dalam satu lajur'),
        (None, u'Garis jingga tebal', u'TALI PENGHUBUNG — menyambung lajur, menyambung kedua berkas'),
        (None, u'Garis merah putus', u'penunjuk 1:1 — bukan hubungan induk-anak')):
    merge_tulis(rl, KOL_LAJUR, u'  ' + lbl, warna,
                'FFFFFF' if warna and warna not in (ISI_PK, ISI_PENAMBAT) else '000000',
                9, True, 'left', 14, bool(warna), bool(warna), bool(warna), bool(warna))
    merge_tulis(rl, KOL_LAJUR + 15, ket, None, '555555', 9, False, 'left', 40)
    rl += 1

for r in range(1, rl + 4):
    ws.row_dimensions[r].height = TINGGI_BARIS

# --- sheet Tali Penghubung -------------------------------------------
ws2 = wb.create_sheet('Tali Penghubung')
KEPALA2 = [u'#', u'Jenis', u'Dari', u'Ke', u'Kunci', u'Kard.',
           u'Di NusantaraRe.xlsx', u'Di ClaimNonProp.xlsx', u'Keterangan']
for j, h in enumerate(KEPALA2, start=1):
    s = ws2.cell(3, j, h)
    s.font = Font(name='Arial', size=9, bold=True, color='FFFFFF')
    s.fill = PatternFill('solid', fgColor='1F4E79')
s = ws2.cell(1, 1, u'TALI PENGHUBUNG — yang membuat kedua berkas menjadi satu skema')
s.font = Font(name='Arial', size=13, bold=True)
s = ws2.cell(2, 1, u'Bagian A: tabel yang ADA DI KEDUA berkas. '
                   u'Bagian B: relasi yang melintasi batas kedua berkas.')
s.font = Font(name='Arial', size=9, color='7F7F7F')

r2 = 4
s = ws2.cell(r2, 1, u'A · TABEL YANG DIPAKAI BERSAMA — titik temu kedua berkas')
s.font = Font(name='Arial', size=10, bold=True, color='C55A11')
r2 += 1
SHEET_NURE = {'T_WORK_CLAIM': u'keempat sheet', 'T_GENERAL_CLAIM': u'Klaim Life · Claim Prop',
              'T_GENERAL_KOMITE': u'Klaim Life · Claim Prop · Komite Claim Prop',
              'T_KOMITE_KOMITELIST': u'Klaim Life · Claim Prop · Komite Claim Prop',
              'T_VIEW_SUGGEST': u'PremiumList NB + EDM · Claim Prop',
              'DOCUMENT_CLAIM': u'Klaim Life + Komite'}
n = 0
for t in T:
    if t['berkas'] != 'KEDUANYA':
        continue
    n += 1
    for j, v in enumerate([n, u'tabel bersama', t['nama'], u'—', u'—', t['kard'],
                           SHEET_NURE.get(t['nama'], u'—'), u'Claim Non Prop',
                           t['catatan'] or u''], start=1):
        c_ = ws2.cell(r2, j, v)
        c_.font = Font(name='Arial', size=9)
        c_.fill = PatternFill('solid', fgColor='FFF2CC')
    r2 += 1

r2 += 1
s = ws2.cell(r2, 1, u'B · RELASI YANG MELINTASI BATAS — digambar sebagai garis di sheet Peta Gabungan')
s.font = Font(name='Arial', size=10, bold=True, color='C55A11')
r2 += 1
for i, p in enumerate(rencana, start=1):
    sm = p['sm'].lstrip('@')
    tj = p['tj'].lstrip('@')
    asal_sm = PETA[p['sm']]['berkas'] if p['sm'] in PETA else u'lajur'
    asal_tj = PETA[p['tj']]['berkas'] if p['tj'] in PETA else u'lajur'
    for j, v in enumerate([i,
                           u'tali penghubung' if p['jenis'] == 'tali' else u'penunjuk 1:1 silang',
                           sm, tj, p['label'],
                           u'1:N' if p['jenis'] == 'tali' else u'1:1',
                           u'ya' if asal_sm in ('NURE', 'KEDUANYA', 'lajur') else u'—',
                           u'ya' if asal_sm in ('CNP', 'KEDUANYA', 'lajur') else u'—',
                           p['ket']], start=1):
        c_ = ws2.cell(r2, j, v)
        c_.font = Font(name='Arial', size=9,
                       color='C55A11' if p['jenis'] == 'tali' else 'C00000')
    r2 += 1

for j, w in enumerate([4, 20, 32, 30, 24, 7, 32, 22, 70], start=1):
    ws2.column_dimensions[openpyxl.utils.get_column_letter(j)].width = w

# --- sheet Daftar Relasi Gabungan ------------------------------------
ws3 = wb.create_sheet('Daftar Relasi Gabungan')
KEPALA3 = [u'#', u'Lajur', u'Induk', u'Anak', u'Kunci tamu', u'Kard.',
           u'Berkas asal', u'Catatan']
s = ws3.cell(1, 1, u'DAFTAR RELASI GABUNGAN — hanya PK dan FK')
s.font = Font(name='Arial', size=13, bold=True)
for j, h in enumerate(KEPALA3, start=1):
    s = ws3.cell(3, j, h)
    s.font = Font(name='Arial', size=9, bold=True, color='FFFFFF')
    s.fill = PatternFill('solid', fgColor='1F4E79')
r3 = 4
BERKAS_NAMA = {'NURE': u'NusantaraRe.xlsx', 'CNP': u'ClaimNonProp.xlsx',
               'KEDUANYA': u'keduanya'}
nrel = 0
for t in T:
    if not t['induk']:
        continue
    nrel += 1
    kunci = u'—  SHARED PK (ID = ID)' if t['kard'] == '1:1' else u''
    if not kunci:
        for b in t['baris']:
            m = re.match(r'FK\s+(\w+)', b.strip())
            if m:
                kunci = m.group(1)
                break
    for j, v in enumerate([nrel, t['lajur'], t['induk'], t['nama'], kunci, t['kard'],
                           BERKAS_NAMA[t['berkas']], t['catatan'] or u''], start=1):
        ws3.cell(r3, j, v).font = Font(name='Arial', size=9)
    r3 += 1
for p in rencana:
    nrel += 1
    for j, v in enumerate([nrel, u'LINTAS-LAJUR', p['sm'].lstrip('@'), p['tj'].lstrip('@'),
                           p['label'], u'1:N' if p['jenis'] == 'tali' else u'1:1',
                           u'keduanya', p['ket']], start=1):
        ws3.cell(r3, j, v).font = Font(name='Arial', size=9, bold=True,
                                       color='C55A11' if p['jenis'] == 'tali' else 'C00000')
    r3 += 1
for j, w in enumerate([4, 14, 32, 34, 26, 7, 20, 74], start=1):
    ws3.column_dimensions[openpyxl.utils.get_column_letter(j)].width = w

# --- sheet Peta Sumber -----------------------------------------------
ws4 = wb.create_sheet('Peta Sumber')
s = ws4.cell(1, 1, u'PETA SUMBER — tiap tabel, lajurnya, dan berkas tempat ia sudah ditetapkan')
s.font = Font(name='Arial', size=13, bold=True)
for j, h in enumerate([u'#', u'Tabel', u'Lajur', u'Kard.', u'Induk',
                       u'Ditetapkan di berkas', u'Baris PK / FK'], start=1):
    s = ws4.cell(3, j, h)
    s.font = Font(name='Arial', size=9, bold=True, color='FFFFFF')
    s.fill = PatternFill('solid', fgColor='1F4E79')
r4 = 4
for i, t in enumerate(T, start=1):
    kunci_baris = u' · '.join(b.strip() for b in t['baris']
                                  if b.strip().startswith(('PK', 'FK')))
    for j, v in enumerate([i, t['nama'], t['lajur'], t['kard'], t['induk'] or u'—',
                           BERKAS_NAMA[t['berkas']], kunci_baris], start=1):
        c_ = ws4.cell(r4, j, v)
        c_.font = Font(name='Arial', size=9,
                       bold=(t['berkas'] == 'KEDUANYA'))
        if t['berkas'] == 'KEDUANYA':
            c_.fill = PatternFill('solid', fgColor='FFF2CC')
    r4 += 1
for j, w in enumerate([4, 36, 12, 8, 32, 22, 110], start=1):
    ws4.column_dimensions[openpyxl.utils.get_column_letter(j)].width = w

wb.save(KELUAR_XLSX)

# =====================================================================
# 5. HTML — SVG, geometri SAMA dengan Excel
# =====================================================================
KX, KY = 17.0, 20.0
PAD = 26


def px(r, c):
    return (PAD + (c - 1) * KX, PAD + (r - 1) * KY)


LEBAR_SVG = PAD * 2 + (MAKS_KOL + 2) * KX
TINGGI_SVG = PAD * 2 + (rl + 4) * KY


def esc(s):
    return (s.replace('&', '&amp;').replace('<', '&lt;').replace('>', '&gt;'))


bag = []
# lapis 1: garis. Digambar LEBIH DULU, tetapi rutenya memang tidak pernah
# melewati kotak, jadi tidak ada yang menutupinya apa pun urutannya.
GAYA_SVG = {'pohon': ('#1F4E79', '2.2', ''),
            'tali': ('#C55A11', '3.4', ''),
            'silang': ('#C00000', '2.4', '7,4')}
for jenis, arah, r1, c1, r2, c2 in SEGMEN:
    warna, tebal, putus = GAYA_SVG[jenis]
    x1, y1 = px(r1, c1)
    x2, y2 = px(r2, c2)
    bag.append(u'<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" '
               u'stroke-width="%s"%s/>'
               % (x1, y1, x2, y2, warna, tebal,
                  (u' stroke-dasharray="%s"' % putus) if putus else u''))
for p in rencana:
    x, y = px(p['rt'], p['ct'])
    warna = '#C55A11' if p['jenis'] == 'tali' else '#C00000'
    bag.append(u'<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="%s"/>'
               % (x, y, x - 9, y - 4.5, x - 9, y + 4.5, warna))
for r, c, teks, jenis in TEKS_KORIDOR:
    x, y = px(r, c)
    warna = '#C55A11' if jenis == 'tali' else '#C00000'
    bag.append(u'<text class="tl" x="%.1f" y="%.1f" fill="%s">%s</text>'
               % (x + 2, y - 4, warna, esc(teks)))

# lapis 2: pita lajur
for r, judul, kunci in PITA:
    x, y = px(r, KOL_LAJUR)
    bag.append(u'<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#%s"/>'
               % (x, y, LEBAR_KOTAK * KX, KY, ISI[kunci]))
    bag.append(u'<text class="pita" x="%.1f" y="%.1f">%s</text>'
               % (x + 8, y + KY * 0.72, esc(judul)))

# lapis 3: penambat
for nm, k in PENAMBAT.items():
    x, y = px(k['r0'], k['c0'])
    bag.append(u'<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#FCE4D6" '
               u'stroke="#C55A11" stroke-width="1.6" stroke-dasharray="5,3"/>'
               % (x, y, k['lebar'] * KX, KY))
    bag.append(u'<text class="pnb" x="%.1f" y="%.1f">%s</text>'
               % (x + 8, y + KY * 0.72, esc(k['teks'])))

# lapis 4: kotak
for nm, k in KOTAK.items():
    t = PETA[nm]
    x, y = px(k['r0'], k['c0'])
    w = k['lebar'] * KX
    isi_baris = list(t['baris']) + ([t['catatan']] if t['catatan'] else [])
    h = (1 + len(isi_baris)) * KY
    bag.append(u'<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#EDF3F9" '
               u'stroke="#1F4E79" stroke-width="1.8"/>' % (x, y, w, h))
    bag.append(u'<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#%s"/>'
               % (x, y, w, KY, ISI[t['lini']]))
    kepala = nm + u'      ' + t['kard']
    if t['berkas'] == 'KEDUANYA':
        kepala += u'      ⬚ ada di KEDUA berkas'
    bag.append(u'<text class="kp" x="%.1f" y="%.1f">%s</text>'
               % (x + w / 2.0, y + KY * 0.72, esc(kepala)))
    yy = y + KY
    for teks in isi_baris:
        f, tt, tb = gaya_baris(teks)
        if f == ISI_PK:
            bag.append(u'<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#FFF2CC"/>'
                       % (x + 1.4, yy, w - 2.8, KY))
        bag.append(u'<text class="br" x="%.1f" y="%.1f" fill="#%s"%s>%s</text>'
                   % (x + 8, yy + KY * 0.72, tt,
                      u' font-weight="700"' if tb else u'', esc(teks)))
        yy += KY
    bag.append(u'<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="none" '
               u'stroke="#1F4E79" stroke-width="1.8"/>' % (x, y, w, h))

SVG = u'\n'.join(bag)

BARIS_TALI = u'\n'.join(
    u'<tr class="%s"><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td><code>%s</code></td>'
    u'<td>%s</td><td>%s</td></tr>'
    % (p['jenis'], i, u'tali penghubung' if p['jenis'] == 'tali' else u'penunjuk 1:1 silang',
       esc(p['sm'].lstrip('@')), esc(p['tj'].lstrip('@')), esc(p['label']),
       u'1:N' if p['jenis'] == 'tali' else u'1:1', esc(p['ket']))
    for i, p in enumerate(rencana, start=1))

BARIS_BERSAMA = u'\n'.join(
    u'<tr><td>%d</td><td><code>%s</code></td><td>%s</td><td>%s</td><td>%s</td></tr>'
    % (i, esc(t['nama']), esc(t['kard']), esc(SHEET_NURE.get(t['nama'], u'—')),
       esc(t['catatan'] or u''))
    for i, t in enumerate([x for x in T if x['berkas'] == 'KEDUANYA'], start=1))

CACAH = collections.Counter(t['lajur'] for t in T)
BARIS_CACAH = u' · '.join(u'%s <b>%d</b>' % (k, CACAH[k]) for k, _, _ in LAJUR)

HTML = u"""<!DOCTYPE html>
<html lang="id"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Peta Gabungan Skema</title>
<style>
:root{--tinta:#14243a;--redup:#5b6b80;--garis:#d6dfea;--latar:#f5f7fa;--kertas:#fff;
      --tali:#C55A11;--silang:#C00000;--pohon:#1F4E79;}
:root:not([data-theme="light"]){}
@media (prefers-color-scheme: dark){:root:not([data-theme="light"]){
  --tinta:#e8edf4;--redup:#9fb0c4;--garis:#2b3purple;--latar:#10151c;--kertas:#161c25;}}
:root[data-theme="dark"]{--tinta:#e8edf4;--redup:#9fb0c4;--garis:#2b3644;--latar:#10151c;--kertas:#161c25;}
*{box-sizing:border-box}
body{margin:0;background:var(--latar);color:var(--tinta);
     font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Arial,sans-serif;
     line-height:1.5;}
.bungkus{max-width:1400px;margin:0 auto;padding:28px 16px 64px;}
h1{font-size:1.5rem;margin:0 0 6px;letter-spacing:-.01em}
h2{font-size:1.05rem;margin:38px 0 10px;padding-bottom:6px;border-bottom:1px solid var(--garis)}
.sub{color:var(--redup);font-size:.86rem;margin:0 0 22px}
.papan{background:var(--kertas);border:1px solid var(--garis);border-radius:10px;
       padding:10px;overflow:auto;}
svg{display:block}
table{border-collapse:collapse;width:100%;font-size:.82rem;background:var(--kertas);
      border:1px solid var(--garis);border-radius:10px;overflow:hidden}
th{background:#1F4E79;color:#fff;text-align:left;padding:7px 9px;font-weight:600;
   font-size:.76rem;letter-spacing:.02em}
td{padding:6px 9px;border-top:1px solid var(--garis);vertical-align:top}
tr.tali td{color:var(--tali)} tr.silang td{color:var(--silang)}
code{font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;font-size:.95em}
.lgd{display:flex;flex-wrap:wrap;gap:8px 20px;margin:12px 0 0;font-size:.8rem;
     color:var(--redup)}
.lgd span{display:flex;align-items:center;gap:7px}
.sw{width:26px;height:12px;border-radius:2px;display:inline-block}
.ln{width:30px;height:0;display:inline-block}
.kartu{display:flex;flex-wrap:wrap;gap:10px;margin:16px 0 0}
.kartu div{background:var(--kertas);border:1px solid var(--garis);border-radius:8px;
           padding:8px 12px;font-size:.8rem;color:var(--redup)}
.kartu b{color:var(--tinta);font-size:1.05rem}
@media (max-width:640px){.bungkus{padding:20px 16px 48px}h1{font-size:1.2rem}}
</style></head><body><div class="bungkus">

<h1>Peta Gabungan Skema &mdash; Nusantara Re</h1>
<p class="sub">Menyatukan <b>Diagram-Skema-Tabel-NusantaraRe.xlsx</b> (Claim Life &middot;
PremiumList &middot; Claim Prop &middot; Komite Claim Prop) dan
<b>Diagram-Skema-Tabel-ClaimNonProp.xlsx</b> menjadi satu peta.
Hanya PK dan FK. Dibangkitkan __TANGGAL__ oleh <code>alat/buat-skema-gabungan.py</code>.</p>

<div class="kartu">
  <div><b>__NTABEL__</b><br>tabel</div>
  <div><b>__NPOHON__</b><br>relasi pohon</div>
  <div><b>__NTALI__</b><br>tali penghubung</div>
  <div><b>__NBERSAMA__</b><br>tabel dipakai KEDUA berkas</div>
</div>

<h2>Peta</h2>
<div class="papan">
<svg width="__W__" height="__H__" viewBox="0 0 __W__ __H__"
     xmlns="http://www.w3.org/2000/svg" font-family="Arial, Helvetica, sans-serif">
<style>
.kp{font-size:11px;font-weight:700;fill:#fff;text-anchor:middle}
.br{font-family:ui-monospace,Consolas,Menlo,monospace;font-size:10.5px}
.pita{font-size:11.5px;font-weight:700;fill:#fff}
.pnb{font-size:10px;font-weight:700;fill:#C55A11}
.tl{font-family:ui-monospace,Consolas,Menlo,monospace;font-size:9px;font-weight:700}
</style>
__SVG__
</svg>
</div>
<div class="lgd">
  <span><i class="sw" style="background:#1F4E79"></i>INTI &mdash; seluruh lini klaim</span>
  <span><i class="sw" style="background:#C55A11"></i>T_CLAIM_ bersama &mdash; Prop &middot; Fac In &middot; Non Prop</span>
  <span><i class="sw" style="background:#2E75B6"></i>khas Claim Non Prop</span>
  <span><i class="sw" style="background:#BF8F00"></i>Claim Life</span>
  <span><i class="sw" style="background:#31859C"></i>Polis / PremiumList</span>
  <span><i class="sw" style="background:#7030A0"></i>Komite</span>
  <span><i class="sw" style="background:#808080"></i>Migrasi</span>
  <span><i class="sw" style="background:#FFF2CC;border:1px solid #d9c98a"></i>baris PK</span>
  <span><i class="ln" style="border-top:3px solid #1F4E79"></i>induk &rarr; anak</span>
  <span><i class="ln" style="border-top:4px solid #C55A11"></i>tali penghubung antar-lajur</span>
  <span><i class="ln" style="border-top:3px dashed #C00000"></i>penunjuk 1:1 silang</span>
</div>

<h2>A &middot; Tabel yang dipakai KEDUA berkas</h2>
<table><thead><tr><th>#</th><th>Tabel</th><th>Kard.</th>
<th>Sheet di NusantaraRe.xlsx</th><th>Keterangan</th></tr></thead>
<tbody>
__BERSAMA__
</tbody></table>

<h2>B &middot; Tali penghubung &mdash; relasi yang melintasi batas kedua berkas</h2>
<table><thead><tr><th>#</th><th>Jenis</th><th>Dari</th><th>Ke</th><th>Kunci</th>
<th>Kard.</th><th>Keterangan</th></tr></thead>
<tbody>
__TALI__
</tbody></table>

<h2>Cara garis digambar &mdash; dan mengapa ia tidak pernah tertutup</h2>
<p class="sub" style="margin-top:0">
Kolom 2&ndash;15 adalah <b>koridor tali</b>; tidak ada satu pun kotak yang boleh
menempatinya. Di dalam tiap lajur, batang pohon sebuah induk berada di kolom
induk+1 sedangkan seluruh keturunannya digambar berurutan mulai kolom induk+3,
sehingga kolom batang itu kosong sepanjang rentangnya. Tiap tali panjang
mendapat kolom jalurnya sendiri lewat pewarnaan graf selang, jadi dua tali tidak
pernah berbagi kolom pada baris yang sama. Pemeriksaan otomatis di
<code>alat/buat-skema-gabungan.py</code> menghitung ulang seluruh titik garis
terhadap seluruh sel yang ditempati kotak &mdash; hasilnya <b>__BENTROK__</b> bentrokan.
Di Excel garis yang sama digambar dengan border sel: ruas tegak sebagai border
kiri, ruas datar sebagai border atas.</p>

</div></body></html>
"""

HTML = (HTML.replace('__SVG__', SVG)
            .replace('__W__', '%d' % int(LEBAR_SVG))
            .replace('__H__', '%d' % int(TINGGI_SVG))
            .replace('__TANGGAL__', TANGGAL)
            .replace('__NTABEL__', str(len(T)))
            .replace('__NPOHON__', str(sum(1 for t in T if t['induk'])))
            .replace('__NTALI__', str(len(TALI)))
            .replace('__NBERSAMA__', str(sum(1 for t in T if t['berkas'] == 'KEDUANYA')))
            .replace('__BERSAMA__', BARIS_BERSAMA)
            .replace('__TALI__', BARIS_TALI)
            .replace('__BENTROK__', str(len(bentrok))))
HTML = HTML.replace('--garis:#2b3purple;', '--garis:#2b3644;')

with io.open(KELUAR_HTML, 'w', encoding='utf-8') as f:
    f.write(HTML)

# =====================================================================
# 6. Rekonsiliasi terhadap KEDUA berkas masukan
# =====================================================================
POLA_NAMA = re.compile(r'\b(T_[A-Z0-9]+(?:_[A-Z0-9]+)*|DOCUMENT_CLAIM|DOCUMENT_POLIS)\b')


def nama_dalam(path):
    w = openpyxl.load_workbook(path, read_only=True, data_only=True)
    ada = set()
    for sh in w:
        for row in sh.iter_rows(values_only=True):
            for v in row:
                if isinstance(v, str):
                    ada.update(POLA_NAMA.findall(v))
    w.close()
    return ada


ADA_NURE = nama_dalam(RUJUK_NURE)
ADA_CNP = nama_dalam(RUJUK_CNP)
MILIK = set(PETA)

lapor(u'')
lapor(u'=== REKONSILIASI =================================================')
lapor(u'tabel di peta gabungan                     : %d' % len(T))
lapor(u'nama tabel terbaca di NusantaraRe.xlsx     : %d' % len(ADA_NURE))
lapor(u'nama tabel terbaca di ClaimNonProp.xlsx    : %d' % len(ADA_CNP))
hilang_n = sorted(ADA_NURE - MILIK)
hilang_c = sorted(ADA_CNP - MILIK)
lapor(u'ada di NusantaraRe tetapi TIDAK di peta    : %d  %s'
      % (len(hilang_n), u', '.join(hilang_n) or u'nihil'))
lapor(u'ada di ClaimNonProp tetapi TIDAK di peta   : %d  %s'
      % (len(hilang_c), u', '.join(hilang_c) or u'nihil'))
lebih = sorted(MILIK - ADA_NURE - ADA_CNP)
lapor(u'ada di peta tetapi TIDAK di kedua berkas   : %d  %s'
      % (len(lebih), u', '.join(lebih) or u'nihil'))
bersama_hitung = sorted(n for n in MILIK if n in ADA_NURE and n in ADA_CNP)
bersama_tandai = sorted(t['nama'] for t in T if t['berkas'] == 'KEDUANYA')
lapor(u'irisan nyata kedua berkas                  : %d  %s'
      % (len(bersama_hitung), u', '.join(bersama_hitung)))
lapor(u'ditandai KEDUANYA di peta                  : %d  %s'
      % (len(bersama_tandai), u', '.join(bersama_tandai)))
lapor(u'selisih penandaan                          : %s'
      % (u', '.join(sorted(set(bersama_hitung) ^ set(bersama_tandai))) or u'nihil'))
# --- rekonsiliasi 3: tiap relasi di sheet "Daftar Relasi" berkas contoh ----
PASANG = set()
_FK = re.compile(r'FK\s+\w+\s+→\s+(T_[A-Z0-9_]+|DOCUMENT_[A-Z]+)\.ID')
for t in T:
    if t['induk']:
        PASANG.add((t['induk'], t['nama']))
    # kunci asing kedua (jalur FAC IN, pintasan komite) terbaca dari baris kotaknya
    for _b in t['baris']:
        for _m in _FK.findall(_b):
            PASANG.add((_m, t['nama']))
for p_ in rencana:
    PASANG.add((p_['sm'].lstrip('@'), p_['tj'].lstrip('@')))
    PASANG.add((p_['tj'].lstrip('@'), p_['sm'].lstrip('@')))
PASANG.add(('T_WORK_CLAIM', 'T_WORK_CLAIM'))
PASANG.add(('T_PREMIUM_LIST_DETAIL', 'T_PREMIUM_LIST_DETAIL'))

w = openpyxl.load_workbook(RUJUK_NURE, read_only=True, data_only=True)
sh = w['Daftar Relasi']
rel_rujuk, tak_tertutup = [], []
for row in sh.iter_rows(min_row=4, values_only=True):
    if not row or not row[0] or not isinstance(row[0], (int, float)):
        continue
    induk, anak = row[2], row[3]
    if not induk or not anak:
        continue
    rel_rujuk.append((induk, anak))
    if (induk, anak) not in PASANG:
        tak_tertutup.append(u'%s → %s' % (induk, anak))
w.close()
lapor(u'')
lapor(u'=== REKONSILIASI 3 — terhadap sheet "Daftar Relasi" NusantaraRe.xlsx ===')
lapor(u'relasi tercatat di berkas contoh           : %d' % len(rel_rujuk))
lapor(u'yang TIDAK tertutup peta gabungan          : %d  %s'
      % (len(tak_tertutup), u', '.join(tak_tertutup) or u'nihil'))
lapor(u'relasi pohon di peta gabungan              : %d' % sum(1 for t in T if t['induk']))
lapor(u'tali penghubung di peta gabungan           : %d' % len(rencana))

lapor(u'')
lapor(u'=== GEOMETRI GARIS ===============================================')
lapor(u'ruas garis digambar                        : %d' % len(SEGMEN))
lapor(u'  pohon  : %d' % sum(1 for s in SEGMEN if s[0] == 'pohon'))
lapor(u'  tali   : %d' % sum(1 for s in SEGMEN if s[0] == 'tali'))
lapor(u'  silang : %d' % sum(1 for s in SEGMEN if s[0] == 'silang'))
lapor(u'kolom jalur koridor yang terpakai          : %d dari %d tersedia'
      % (JUM_JALUR, KOR_AKHIR - KOR_AWAL + 1))
lapor(u'titik garis yang jatuh di dalam kotak      : %d  %s'
      % (len(bentrok), u'' if not bentrok else str(bentrok[:10])))
lapor(u'')
lapor(u'keluaran: %s  (%d bait)' % (os.path.basename(KELUAR_XLSX),
                                    os.path.getsize(KELUAR_XLSX)))
lapor(u'keluaran: %s  (%d bait)' % (os.path.basename(KELUAR_HTML),
                                    os.path.getsize(KELUAR_HTML)))
