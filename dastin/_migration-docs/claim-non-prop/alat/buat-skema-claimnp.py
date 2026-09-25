# -*- coding: utf-8 -*-
"""Rancangan tabel Claim Non Prop - hanya PK dan FK. HTML + Excel.

HANYA MEMBACA. Nol berkas sumber berubah.

Bentuknya mengikuti Diagram-Skema-Tabel-NusantaraRe.xlsx yang sudah ada:
akar T_WORK_CLAIM, T_GENERAL_CLAIM 1:1 SHARED PK, lalu tabel lini berawalan
T_CLAIM_ - sekelas T_CLAIMP_ pada sheet Claim Prop.

Sumber:
  struktur-claimdata-lama.md / datar-claimdata-lama.csv   pohon .ClaimData
  datar-claimdata-salinan.csv                             salinan beku
  ../../MEMORI_PEMAHAMAN.MD                               arti bisnis
  pengetahuan/ddl/, pengetahuan/SCHEMA-ACTUAL.csv         PK/UNIQUE lama
  ddl-usulan/*.sql                                        22 tabel yang sudah diputuskan
  Diagram-Skema-Tabel-NusantaraRe.xlsx                    konvensi penamaan dan bentuk

Keluaran:
  ERD-CLAIM-NON-PROP.html
  Diagram-Skema-Tabel-ClaimNonProp.xlsx

Jalankan:  python alat/buat-skema-claimnp.py
"""
from __future__ import print_function
import io, os, re, sys, csv, glob, collections

ALAT = os.path.dirname(os.path.abspath(__file__))
AKAR = os.path.normpath(os.path.join(ALAT, '..'))
ROOT = os.path.normpath(os.path.join(AKAR, '..', '..'))
TANGGAL = '20 September 2026'

try:
    import openpyxl
    from openpyxl.styles import Font, PatternFill, Alignment, Border, Side
except ImportError:
    print('GAGAL: openpyxl tidak ada. pip install openpyxl')
    sys.exit(1)


def baca(p):
    return io.open(p, encoding='utf-8').read()


def baca_csv(nama):
    with io.open(os.path.join(AKAR, nama), encoding='utf-8-sig', newline='') as f:
        return list(csv.DictReader(f))


# =====================================================================
# 1. Bukti dari sumber
# =====================================================================
# Sumbernya struktur-claimdata-lama.md LANGSUNG - keempat CSV pendampingnya
# sudah dipindahkan ke sampah/, jadi ia tidak lagi sah dipakai sebagai sumber.
ERD = os.path.join(AKAR, '4-erd-dan-tabel-datar')
SPEC = os.path.join(AKAR, '2-to-spec')
STRUKTUR = baca(os.path.join(ERD, 'struktur-claimdata-lama.md'))

# Nama simpul di pohon bagian 4 dan 5 - dipakai untuk MEMERIKSA bahwa tiap
# tabel rancangan memang berasal dari simpul yang ada, bukan dari ingatan.
SIMPUL_POHON = set()
POLA_BLOK = re.compile('```' + chr(10) + '(.*?)' + chr(10) + '```', re.S)
for blok in POLA_BLOK.findall(STRUKTUR):
    for b_ in blok.split(chr(10)):
        for m in re.finditer(r'([A-Z][A-Za-z0-9_]{2,})(\[\]|\{\})', b_):
            SIMPUL_POHON.add(m.group(1))
        m = re.match(r'^[\s│├└─●•|+-]*([A-Z][A-Za-z0-9_]{2,})\s*$', b_)
        if m:
            SIMPUL_POHON.add(m.group(1))
if len(SIMPUL_POHON) < 20:
    print('GAGAL: pohon di struktur-claimdata-lama.md tidak terurai (%d simpul)'
          % len(SIMPUL_POHON))
    sys.exit(1)

# Angka ringkasan bagian 2 - dikutip, dan dibandingkan dengan yang kita pakai.
def angka(pola, bawaan=0):
    m = re.search(pola, STRUKTUR)
    return int(m.group(1)) if m else bawaan


N_SIMPUL = angka(r'Simpul dalam pohon `\.ClaimData`\s*\|\s*\*\*(\d+)\*\*')
N_SKALAR_AKAR = angka(r'Properti skalar langsung di `\.ClaimData`\s*\|\s*\*\*(\d+)\*\*')
N_WADAH_AKAR = angka(r'Kontainer \(Page/Page List\) langsung di `\.ClaimData`\s*\|\s*\*\*(\d+)\*\*')

# tabel yang sudah diputuskan di ddl-usulan/ - dipakai sebagai PEMBANDING,
# supaya rancangan ini tidak menjadi rancangan kedua yang berdiri sendiri.
tabel_klaimnp = []
for p in sorted(glob.glob(os.path.join(SPEC, 'ddl-usulan', '[0-9]*.sql'))):
    m = re.search(r'CREATE TABLE KLAIMNP\.(\w+)\s*\(', baca(p))
    if m:
        tabel_klaimnp.append(m.group(1))

# PK/UNIQUE di DDL lama - dari pengetahuan/
lama_pk = collections.defaultdict(lambda: {'pk': False, 'uq': False})
with io.open(os.path.join(AKAR, 'pengetahuan', 'SCHEMA-ACTUAL.csv'), encoding='utf-8-sig') as f:
    for b in csv.DictReader(f, delimiter=';'):
        if b.get('is_pk') == 'Y':
            lama_pk[b['object_name']]['pk'] = True
        if b.get('is_unique') == 'Y':
            lama_pk[b['object_name']]['uq'] = True

# Salinan beku, dari tabel bagian 6 berkas itu sendiri.
SALINAN_BEKU = {}
for m in re.finditer(r'\|\s*`AdjustmentList\[\]\.(\w+)`\s*\|\s*`([^`]+)`', STRUKTUR):
    SALINAN_BEKU[m.group(1)] = m.group(2)

# =====================================================================
# 2. Rancangan tabel - keputusan, satu per satu, dengan sebabnya
# =====================================================================
# Tiap entri: (NAMA, KARDINALITAS, INDUK, KUNCI_TAMU, CATATAN_FK, ASAL_PEGA,
#              PADANAN_KLAIMNP, KELOMPOK)
# KELOMPOK: 'akar' lintas-lini · 'np' Claim Non Prop · 'bersama' lintas-lini
#           · 'komite' · 'migrasi'
T = []


def tabel(nama, kard, induk, kunci, catatan, asal, padanan, kelompok, tambahan=()):
    T.append({'nama': nama, 'kard': kard, 'induk': induk, 'kunci': kunci,
              'catatan': catatan, 'asal': asal, 'padanan': padanan,
              'kelompok': kelompok, 'tambahan': list(tambahan)})


# --- akar, sudah ditetapkan workbook yang ada ------------------------
tabel('T_WORK_CLAIM', u'akar', None, None, u'',
      u'pyWorkPage — satu baris per work object',
      u'—', 'akar',
      [u'PK   ID            teks berformat  CLMNP-xxxxxx (klaim) / TKMT-xxxxxx (komite)',
       u'FK   COVER_KEY  →  T_WORK_CLAIM.ID      self · nullable',
       u'LINI = NONPROP pada kedua baris · baris komite lahir saat penyerahan',
       u'CHECK:  ID LIKE \'CLMNP-%\' ∧ COVER_KEY IS NULL  ∨  ID LIKE \'TKMT-%\' ∧ COVER_KEY NOT NULL'])

tabel('T_GENERAL_CLAIM', u'1:1', 'T_WORK_CLAIM', u'—  SHARED PK (ID = ID)',
      u'SHARED PK — tidak ada kolom FK terpisah',
      u'ClaimData{} · 91 properti skalar',
      u'KLAIM', 'bersama',
      [u'kolom khas lini WAJIB nullable — tabel ini dipakai semua lini klaim'])

# --- tabel Claim Non Prop, urut mengikuti pohon ----------------------
tabel('T_CLAIM_CLAIM_AMOUNT', u'1:N', 'T_GENERAL_CLAIM', 'CLAIM_ID', u'CASCADE',
      u'ClaimData.ListClaimAmount[] · 15 field · nilai kerugian per mata uang',
      u'NILAI_KLAIM_MATA_UANG', 'np')

tabel('T_CLAIM_SPREADING_RISK', u'1:N', 'T_GENERAL_CLAIM', 'CLAIM_ID', u'CASCADE',
      u'ClaimData.SpreadingRisk[] · 29 field · ★ inti alokasi XOL, n=215',
      u'ALOKASI_LAYER', 'np',
      [u'⚠ BERBEDA DARI CLAIM PROP. Di sheet Claim Prop, SpreadingRisk BUKAN tabel',
       u'   (PAGE-REMOVE lalu RDB-LIST tiap kali — bacaan master treaty).',
       u'   Di Non Prop ia HASIL ALOKASI dan disalin beku ke AdjustmentList[] — jadi data.'])

tabel('T_CLAIM_RETENTION_CEDANT', u'1:N', 'T_GENERAL_CLAIM', 'CLAIM_ID', u'CASCADE',
      u'baris "UR" di dalam ClaimData.SpreadingRisk[] — DIPISAH jadi tabel sendiri',
      u'RETENSI_CEDANT', 'np',
      [u'Di sistem lama ia BARIS di daftar alokasi yang sama, jadi tiap loop ikut',
       u'melihatnya kecuali menyaring sendiri. Dipisah supaya penjumlahan tidak',
       u'pernah mencampurnya diam-diam. [keputusan ADR-0010]'])

tabel('T_CLAIM_SPREAD_LOSS', u'1:N', 'T_GENERAL_CLAIM', 'CLAIM_ID', u'CASCADE',
      u'ClaimData.CNPSpreadLoss[] · 8 field · pembagian kerugian per COB',
      u'PEMBAGIAN_KERUGIAN', 'np')

tabel('T_CLAIM_SPREADING', u'1:N', 'T_GENERAL_CLAIM', 'CLAIM_ID', u'CASCADE',
      u'ClaimData.SpreadingClaim[] · penyebaran ke retrosesi',
      u'PENYEBARAN', 'np',
      [u'≡ sama dengan T_CLAIM_SPREADING di Claim Prop dan Claim Fac In'])

tabel('T_CLAIM_BREAK_QS', u'1:N', 'T_GENERAL_CLAIM', 'CLAIM_ID', u'CASCADE',
      u'ClaimData.SpreadingBreakQS[] · n=28 · penyebaran quota share',
      u'PENYEBARAN · bagian QS', 'np',
      [u'SEJAJAR, bukan anak _SPREADING — sama seperti di Claim Prop.',
       u'Dulu digabung ke T_CLAIM_SPREADING dengan pembeda JENIS_REASURANSI;',
       u'dipisah 20-09-2026 mengikuti nama yang berlaku di berkas rujukan.'])

tabel('T_CLAIM_REINSTATEMENT', u'1:N', 'T_GENERAL_CLAIM', 'CLAIM_ID', u'CASCADE',
      u'ClaimData.ReinstatementList[] · 12 field · premi pemulihan',
      u'PREMI_PEMULIHAN', 'np')

tabel('T_CLAIM_ESTIMATION', u'1:N', 'T_GENERAL_CLAIM', 'CLAIM_ID', u'CASCADE',
      u'ClaimData.EstimationList[] · 7 field · estimasi awal',
      u'ESTIMASI_AWAL', 'np')

tabel('T_CLAIM_INTEREST', u'1:N', 'T_GENERAL_CLAIM', 'CLAIM_ID', u'CASCADE',
      u'ClaimData.InterestList[] · n=39 · kepentingan yang dipertanggungkan',
      u'OBJEK_PERTANGGUNGAN · bagian interest', 'np',
      [u'≡ sama dengan T_CLAIM_INTEREST di Claim Prop'])

tabel('T_CLAIM_OBJECT', u'1:N', 'T_GENERAL_CLAIM', 'CLAIM_ID', u'CASCADE',
      u'ClaimData.ObjectList[] · objek pertanggungan',
      u'OBJEK_PERTANGGUNGAN', 'np',
      [u'≡ sama dengan T_CLAIM_OBJECT di Claim Fac In',
       u'⚠ Di Claim Prop ObjectList ditandai "NOL penulis · residu save-as Non Prop"',
       u'   — di Non Prop dan Fac In-lah ia berasal.'])

tabel('T_CLAIM_OBJECT_ITEM', u'1:N', 'T_CLAIM_OBJECT', 'OBJECT_ID', u'CASCADE',
      u'ClaimData.ObjectList[].ObjectItemList[] · rincian tiap objek',
      u'OBJEK_PERTANGGUNGAN · bagian item', 'np',
      [u'≡ sama dengan T_CLAIM_OBJECT_ITEM di Claim Fac In',
       u'Dulu digabung ke T_CLAIM_OBJECT; dipisah 20-09-2026 mengikuti berkas rujukan',
       u'dan struktur lama, yang memang menaruhnya sebagai halaman bersarang.'])

tabel('T_CLAIM_RECEIVER', u'1:N', 'T_GENERAL_CLAIM', 'CLAIM_ID', u'CASCADE',
      u'ClaimData.ReceiverClaim[] · 13 field · sepasang rekening (utama + varian …2)',
      u'REKENING_PENERIMA', 'np')

tabel('T_VIEW_SUGGEST', u'1:N', 'T_GENERAL_CLAIM', 'CLAIM_ID', u'CASCADE',
      u'ClaimData.SuggestList[] · kronologi / audit trail',
      u'KRONOLOGI_KLAIM', 'bersama',
      [u'BERSAMA dengan PremiumList — punya DUA induk, PREMIUM_LIST_ID dan CLAIM_ID,',
       u'keduanya nullable, dengan CHECK ( (PREMIUM_LIST_ID IS NULL) <> (CLAIM_ID IS NULL) ).'])

tabel('DOCUMENT_CLAIM', u'1:N', 'T_GENERAL_CLAIM', 'CLAIM_ID', u'⚠ [data DBA]',
      u'ClaimData.Attachment[] + AttachmentPaid{} · dokumen klaim',
      u'DOKUMEN_KLAIM', 'bersama',
      [u'⚠ LINTAS-LINI. Daftar kolom TIDAK terbaca dari korpus — [data DBA].',
       u'   Induknya untuk Non Prop mengikuti pola Prop: CLAIM_ID → T_GENERAL_CLAIM.ID.'])

tabel('T_CLAIM_ADJUSTMENT', u'1:N', 'T_GENERAL_CLAIM', 'CLAIM_ID', u'CASCADE',
      u'ClaimData.AdjustmentList[] · 47 field skalar · ★ unit transaksi pembayaran',
      u'AKSEPTASI + ADJUSTMENT', 'np',
      [u'FK   KOMITE_ID  →  T_WORK_CLAIM.ID     1:1 · nullable · index UNIK',
       u'⚠ Di KLAIMNP ia DUA tabel: AKSEPTASI (kunci alami klaim×layer×mata uang)',
       u'   dan ADJUSTMENT (baris penyesuaian). Di sini digambar satu kotak karena',
       u'   satu halaman Pega; pemisahannya keputusan ADR-0024. [terbuka] mana yang dipakai.'])

tabel('T_CLAIM_ADJ_SPREADING', u'1:N', 'T_CLAIM_ADJUSTMENT', 'ADJUSTMENT_ID', u'CASCADE',
      u'AdjustmentList[].SpreadingAdjustment[] ═ salinan beku dari ClaimData.SpreadingClaim',
      u'PENYEBARAN (baris beku)', 'np',
      [u'⚠ NAMA TUJUAN ≠ NAMA SUMBER. SpreadingClaim di tingkat klaim menjadi',
       u'   SpreadingAdjustment di dalam Adjustment. Nama yang sama berarti DUA hal.'])

tabel('T_CLAIM_ADJ_QUOTA_SHARE', u'1:N', 'T_CLAIM_ADJUSTMENT', 'ADJUSTMENT_ID', u'CASCADE',
      u'AdjustmentList[].SpreadingQuotaShare[] ═ salinan beku dari ClaimData.SpreadingBreakQS',
      u'PENYEBARAN (baris beku)', 'np')

tabel('T_CLAIM_ADJ_LOSS_ALLOCATION', u'1:N', 'T_CLAIM_ADJUSTMENT', 'ADJUSTMENT_ID', u'CASCADE',
      u'AdjustmentList[].LossAllocation[] + SpreadingRisk[] ═ salinan beku dari ClaimData.SpreadingRisk',
      u'ALOKASI_LAYER (baris beku)', 'np',
      [u'Tiap Adjustment membawa POTRET alokasi saat ia dibuat. Perhitungan berikutnya',
       u'di tingkat klaim TIDAK mengubah potret itu. Perilaku ini harus dipertahankan.'])

tabel('T_CLAIM_ADJ_LAYER', u'1:N', 'T_CLAIM_ADJUSTMENT', 'ADJUSTMENT_ID', u'CASCADE',
      u'AdjustmentList[].CNPLayerList[] · XOL · XOLID   [bukti: identitas-kelas]',
      u'⚠ TIDAK ADA di KLAIMNP', 'np',
      [u'⚠ BELUM ADA PADANANNYA di 22 tabel ddl-usulan/. Jalur penuhnya TIDAK PERNAH',
       u'   ditulis di 279 berkas; keberadaannya dipastikan dari identitas kelas (bagian 5).',
       u'[terbuka] apakah ia tabel sendiri, atau kolom di T_CLAIM_ADJ_LAYER_CURRENCY.'])

tabel('T_CLAIM_ADJ_LAYER_CURRENCY', u'1:N', 'T_CLAIM_ADJ_LAYER', 'LAYER_ID', u'CASCADE',
      u'AdjustmentList[].CNPLayerList[].CNPCurrencyList[] · 21 field · KEDALAMAN 4',
      u'⚠ TIDAK ADA di KLAIMNP', 'np',
      [u'⚠ Data-Adjustment BERSARANG DI DALAM DIRINYA SENDIRI TIGA TINGKAT:',
       u'   Adjustment → CNPLayerList → CNPCurrencyList, dengan arti transaksi, layer,',
       u'   dan mata uang. Satu tabel untuk ketiganya akan salah.'])

tabel('T_CLAIM_RETRO', u'1:N', 'T_CLAIM_SPREAD_LOSS', 'SPREAD_LOSS_ID', u'CASCADE',
      u'ClaimData.CNPSpreadLoss[].RetroList[] · TREATYTYPENAME · COMMISION · OVR_COMM …',
      u'⚠ TIDAK ADA di KLAIMNP', 'np',
      [u'⚠ Kelasnya GISFW-Int-RETROCESSIONLIFE — kelas INTEGRASI, bukan kelas data.',
       u'   RetroList juga menempel di ClaimData.SpreadingRisk[]; induknya karena itu DUA.',
       u'[terbuka] satu tabel dengan dua induk nullable, atau dua tabel.'])

# --- komite, mengikuti workbook yang ada -----------------------------
tabel('T_GENERAL_KOMITE', u'1:1', 'T_WORK_CLAIM', u'—  SHARED PK (ID = ID)',
      u'SHARED PK · baris TKMT- untuk lini NONPROP',
      u'AdjustmentList[].ComiteeClaim[] + ClaimData.ClaimComitee[]',
      u'— (di luar 22 tabel)', 'komite',
      [
       u'FK   ADJUSTMENT_ID  →  T_CLAIM_ADJUSTMENT.ID    NOT NULL · index UNIK',
       u'⚠ Folder Komite Claim Non Prop TIDAK dibuka. Bentuk tabel ini diambil dari',
       u'   sheet Komite Claim Prop yang sudah ada — lintas-lini, tujuh kolom.'])

tabel('T_KOMITE_KOMITELIST', u'1:N', 'T_GENERAL_KOMITE', 'DATA_KOMITE_ID', u'CASCADE',
      u'KomiteList · satu baris = satu penyetuju · sembilan kolom',
      u'— (di luar 22 tabel)', 'komite')

# --- tabel yang tidak lahir dari pohon, tetapi ada di KLAIMNP --------
tabel('T_CLAIM_CLOSING_DATE', u'bebas', None, None, u'',
      u'TIDAK DITEMUKAN di pohon — angka 25 di sistem lama punya dua sumber',
      u'TUTUP_BUKU', 'np',
      [u'Data bertanggal berlaku, lingkup global. Tidak menempel pada satu klaim,',
       u'jadi tidak punya FK. [keputusan ADR-0025]'])

tabel('T_CLAIM_RATE', u'bebas', None, None, u'',
      u'TIDAK DITEMUKAN di pohon — brokerage dan pajak dari konstanta di rule',
      u'TARIF_BERLAKU', 'np',
      [u'Sama: tarif bertanggal berlaku, tanpa induk.'])

tabel('T_CLAIM_CORRECTION', u'bebas', None, None, u'',
      u'TIDAK DITEMUKAN di pohon — pengganti 29 langkah tambalan di dalam activity',
      u'KOREKSI_NILAI', 'np',
      [u'Menunjuk baris mana pun lewat TABEL_SASARAN + ID_SASARAN.',
       u'FOREIGN KEY tidak dapat menunjuk banyak tabel, jadi ia TANPA FK.'])

tabel('T_CLAIM_DATE_GUARD', u'1:N', 'T_GENERAL_CLAIM', 'CLAIM_ID', u'CASCADE',
      u'TIDAK DITEMUKAN di pohon — penjaga batas tanggal',
      u'KLAIM_PENJAGA_TANGGAL', 'np')

tabel('T_CLAIM_OUTBOUND_ARCHIVE', u'1:N', 'T_GENERAL_CLAIM', 'CLAIM_ID', u'CASCADE',
      u'InputParamOs / InputParamOsCNP → OS_AKSEPTASI_KLAIM (bagian 7.1)',
      u'ARSIP_MUATAN_KELUAR', 'np',
      [u'Arsip apa yang BENAR-BENAR dikirim. Barisnya TAMBAHAN, bukan keadaan —',
       u'penulisannya SELISIH: InputParamOs.X = X − OutOSAcc.pxResults(1).X.',
       u'Hanya SELECT dan INSERT; tanpa UPDATE, tanpa DELETE.'])

tabel('T_CLAIM_MIG_CORRELATION', u'bebas', None, None, u'',
      u'jembatan ke sistem lama — pzInsKey · pyID · CASEID · IndexObject',
      u'MIGRASI_KORELASI', 'migrasi',
      [u'Menunjuk balik lewat TABEL_TUJUAN + ID_TUJUAN, bukan lewat FK. Berumur terbatas.',
       u'⚠ OS_AKSEPTASI_KLAIM di sistem lama TIDAK punya primary key maupun unique',
       u'   constraint — itu sebab REQ-018 ada, dan sebab korelasi ini perlu.'])

tabel('T_CLAIM_MIG_LANDING', u'bebas', None, None, u'',
      u'tempat bentuk lama mendarat utuh — DATA_JSON apa adanya',
      u'MIGRASI_PENDARATAN', 'migrasi',
      [u'⚠ DATA_JSON lama MEMUAT yang tidak boleh diimpor: GetPageJSONString()',
       u'   menyerialisasi seluruh halaman. Buang seluruh kunci px* py* pz*.',
       u'   SEMUA nilai di JSON lama bertipe string; tanggal DUA format berdampingan.'])

tabel('T_CLAIM_MIG_REJECTED', u'1:N', 'T_CLAIM_MIG_LANDING', 'LANDING_ID', u'nullable',
      u'nilai yang tidak dapat diurai — tercatat, tidak dibulatkan, tidak dibuang',
      u'MIGRASI_NILAI_DITOLAK', 'migrasi',
      [u'Kolomnya boleh kosong: sebagian nilai ditolak berasal langsung dari kolom',
       u'tabel lama, bukan dari muatan. Kosong = bukan dari muatan, BUKAN tidak diketahui.'])

# --- yang BUKAN tabel ------------------------------------------------
BUKAN_TABEL = [
    (u'ClaimData.ListTotalEstimation[]', u'total per mata uang', u'SELECT SUM — view V_TOTAL_ESTIMASI'),
    (u'ClaimData.TotalInterestInsured[]', u'total nilai pertanggungan', u'SELECT SUM — view V_TOTAL_NILAI_PERTANGGUNGAN'),
    (u'ClaimData.SpreadingAdjustment[]', u'salinan dari AdjustmentList(n)', u'variabel; agregat lewat view'),
    (u'ClaimData.SpreadingAdjustmentQS[]', u'pola sama', u'variabel; agregat lewat view'),
    (u'AdjustmentList[].ListClaimAcceptation[]', u'rekap akseptasi', u'view V_REKAP_AKSEPTASI'),
    (u'AdjustmentList[].CNPSpreadLoss[]', u'satu field CNPFlagOuts saja', u'kolom di T_CLAIM_ADJ_*'),
    (u'AdjustmentList[].CurencyAdjustment[]', u'hasil query kurs (Kurs.pxResults)', u'kurs dibaca saat hitung'),
    (u'AdjustmentList[].DataCommitteeTreaty{}', u'5 field circum/remarks', u'kolom KOMITE_* di T_CLAIM_ADJUSTMENT'),
    (u'ClaimData.PolicyData{} · QuotationData{} · MarketingData{}', u'bacaan master polis/treaty', u'query master, bukan data klaim'),
    (u'ObjectList[].ObjectItemList[].Adjustment{}', u'satu field CurrencyID · kedalaman 4', u'residu; tidak dimigrasi'),
]

# --- relasi datar -----------------------------------------------------
RELASI = []
for t in T:
    if t['induk'] and t['kunci']:
        RELASI.append((t['induk'], t['nama'], t['kunci'], t['kard']))
RELASI.insert(0, ('T_WORK_CLAIM', 'T_WORK_CLAIM', 'COVER_KEY', '1:N'))
RELASI.append(('T_WORK_CLAIM', 'T_CLAIM_ADJUSTMENT', 'KOMITE_ID', '1:1'))
RELASI.append(('T_CLAIM_ADJUSTMENT', 'T_GENERAL_KOMITE', 'ADJUSTMENT_ID', '1:1'))

# --- pemeriksaan: tiap simpul yang dikutip memang ADA di pohonnya ------
# Rancangan yang menyebut halaman yang tidak ada di sumber adalah rancangan
# yang ditulis dari ingatan. Ini yang mencegahnya.
dikutip, tak_ada = set(), []
for t_ in T:
    # nama halaman yang dikutip di kolom asal, mis. ClaimData.SpreadingRisk[]
    for nm in re.findall('([A-Z][A-Za-z0-9_]{3,})', t_['asal']):
        if nm in ('ClaimData', 'TIDAK', 'DITEMUKAN', 'DATA_JSON', 'Pega',
                  'OS_AKSEPTASI_KLAIM', 'GetPageJSONString', 'NONPROP',
                  # ada di sumbernya, tetapi di LUAR pohon bagian 4/5:
                  'IndexObject',        # field skalar, bukan wadah
                  'InsKey', 'WorkPage', # potongan pzInsKey / pyWorkPage
                  'InputParamOs', 'InputParamOsCNP',   # wadah keluar, bagian 7.1
                  'KomiteList'):        # bagian 3: menempel di work page,
                                        # DI LUAR pohon .ClaimData - disebut di sana
            continue
        dikutip.add(nm)
for nm in sorted(dikutip):
    if nm not in SIMPUL_POHON and nm.upper() != nm:
        tak_ada.append(nm)

PETA = {t['nama']: t for t in T}
ANAK = collections.defaultdict(list)
for t in T:
    if t['induk']:
        ANAK[t['induk']].append(t['nama'])

# =====================================================================
# 3. HTML
# =====================================================================
def esc(s):
    return (u'%s' % s).replace('&', '&amp;').replace('<', '&lt;') \
                      .replace('>', '&gt;').replace('"', '&quot;')


WARNA = {'akar': 'w-akar', 'bersama': 'w-bersama', 'np': 'w-np',
         'komite': 'w-komite', 'migrasi': 'w-migrasi'}

LEBAR, T_KEPALA, T_BARIS = 300, 26, 15.2
JARAK_X, JARAK_Y, MX, MY = 340, 22, 30, 92


def baris_kotak(t):
    b = []
    if t['nama'] == 'T_WORK_CLAIM':
        b += t['tambahan']
        return b
    if t['kunci'] and u'SHARED PK' in t['kunci']:
        b.append(u'PK   ID  =  %s.ID' % t['induk'])
    else:
        b.append(u'PK   ID')
        if t['induk']:
            b.append(u'FK   %s  →  %s.ID   %s' % (t['kunci'], t['induk'], t['catatan']))
    b += t['tambahan']
    b.append(u'←  %s' % t['asal'])
    b.append(u'≡  %s' % t['padanan'])
    return b


def tinggi(t):
    return T_KEPALA + len(baris_kotak(t)) * T_BARIS + 8


# tata letak: kolom menurut kedalaman pohon
def kedalaman(nama, lihat=None):
    lihat = lihat or set()
    t = PETA[nama]
    if not t['induk'] or nama in lihat:
        return 0
    return 1 + kedalaman(t['induk'], lihat | {nama})


kol = collections.defaultdict(list)
for t in T:
    kol[kedalaman(t['nama'])].append(t['nama'])
# sebar kolom dalam supaya tidak terlalu tinggi
N_KOL = 5
kolom_isi = [[] for _ in range(N_KOL)]
kolom_t = [0.0] * N_KOL
for d in sorted(kol):
    for nama in kol[d]:
        i = min(range(min(d, N_KOL - 1), N_KOL), key=lambda j: kolom_t[j])
        kolom_isi[i].append(nama)
        kolom_t[i] += tinggi(PETA[nama]) + JARAK_Y
pos = {}
for ix, daftar in enumerate(kolom_isi):
    y = MY
    for nama in daftar:
        pos[nama] = (MX + ix * JARAK_X, y)
        y += tinggi(PETA[nama]) + JARAK_Y
W = MX * 2 + (N_KOL - 1) * JARAK_X + LEBAR
H = MY + max(kolom_t) + 24

svg = [u'<svg viewBox="0 0 %d %d" width="100%%" class="erd" role="img" '
       u'aria-label="ERD Claim Non Prop">' % (W, H), u'<defs>']
for ident, warna in (('pohon', 'var(--r-pohon)'), ('silang', 'var(--r-silang)')):
    svg.append(u'<marker id="m-%s" viewBox="0 0 12 12" refX="11" refY="6" markerWidth="10" '
               u'markerHeight="10" orient="auto-start-reverse"><path d="M11,6 L0,0 M11,6 L0,6 '
               u'M11,6 L0,12" fill="none" stroke="currentColor" stroke-width="1.4"/></marker>'
               % ident)
    svg.append(u'<marker id="s-%s" viewBox="0 0 12 12" refX="11" refY="6" markerWidth="10" '
               u'markerHeight="10" orient="auto-start-reverse"><path d="M3,1 L3,11" fill="none" '
               u'stroke="currentColor" stroke-width="1.6"/></marker>' % ident)
svg.append(u'</defs>')

for induk, anak, kunci, kard in RELASI:
    if induk not in pos or anak not in pos or induk == anak:
        continue
    silang = kard == '1:1' and PETA.get(anak, {}).get('induk') != induk
    xa, ya = pos[induk]
    xb, yb = pos[anak]
    y1, y2 = ya + tinggi(PETA[induk]) / 2.0, yb + tinggi(PETA[anak]) / 2.0
    x1, x2 = (xa + LEBAR, xb) if xb >= xa else (xa, xb + LEBAR)
    dx = max(30, abs(x2 - x1) * 0.4)
    arah = 1 if x2 >= x1 else -1
    svg.append(u'<path class="rel %s" d="M%.1f,%.1f C%.1f,%.1f %.1f,%.1f %.1f,%.1f" '
               u'marker-start="url(#s-%s)" marker-end="url(#%s-%s)"><title>%s → %s · %s · %s</title></path>'
               % ('r-silang' if silang else 'r-pohon', x1, y1, x1 + dx * arah, y1,
                  x2 - dx * arah, y2, x2, y2,
                  'silang' if silang else 'pohon',
                  's' if kard == '1:1' else 'm', 'silang' if silang else 'pohon',
                  esc(induk), esc(anak), esc(kunci), esc(kard)))

for nama, (x, y) in pos.items():
    t = PETA[nama]
    b = baris_kotak(t)
    svg.append(u'<g class="ent %s">' % WARNA[t['kelompok']])
    svg.append(u'<rect class="badan" x="%d" y="%.1f" width="%d" height="%.1f" rx="3"/>'
               % (x, y, LEBAR, tinggi(t)))
    svg.append(u'<path class="kepala" d="M%d,%.1f h%d v%d h-%d z"/>'
               % (x, y + 3, LEBAR, T_KEPALA - 3, LEBAR))
    svg.append(u'<text class="t-nama" x="%d" y="%.1f">%s</text>' % (x + 9, y + 18, esc(nama)))
    svg.append(u'<text class="t-kard" x="%d" y="%.1f">%s</text>'
               % (x + LEBAR - 9, y + 18, esc(t['kard'])))
    for i, teks in enumerate(b):
        yy = y + T_KEPALA + i * T_BARIS + 11
        kls = ('b-pk' if teks.startswith('PK') else 'b-fk' if teks.startswith('FK')
               else 'b-awas' if teks.lstrip().startswith(('⚠', '[terbuka]')) else
               'b-asal' if teks.startswith(('←', '≡')) else 'b-ket')
        svg.append(u'<text class="%s" x="%d" y="%.1f">%s</text>'
                   % (kls, x + 9, yy, esc(teks if len(teks) <= 56 else teks[:55] + u'…')))
    svg.append(u'</g>')
svg.append(u'</svg>')
SVG = u'\n'.join(svg)

baris_relasi = u''.join(
    u'<tr><td class="num">%d</td><td>%s</td><td>%s</td><td><code>%s</code></td><td>%s</td></tr>'
    % (i, esc(a), esc(b), esc(k), esc(kd))
    for i, (a, b, k, kd) in enumerate(RELASI, 1))
baris_bukan = u''.join(u'<tr><td><code>%s</code></td><td>%s</td><td>%s</td></tr>'
                       % (esc(a), esc(b), esc(c)) for a, b, c in BUKAN_TABEL)

HTML = u"""<!DOCTYPE html><html lang="id"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Rancangan tabel Claim Non Prop</title><style>
:root{--tinta:#16181d;--tinta-2:#525a68;--tinta-3:#7b8494;--latar:#fff;--latar-2:#f7f8fa;
--garis:#d8dce4;--badan:#fdfdfe;--r-pohon:#1f4fd8;--r-silang:#c0262b;
--akar:#1e3a8a;--akar-t:#fff;--bersama:#f6b93b;--bersama-t:#3d2c00;
--np:#bcd7f7;--np-t:#10233f;--komite:#ddc6f2;--komite-t:#2e1b45;--migrasi:#d5d9e0;--migrasi-t:#2b3039;
--pk:#b45309;--fk:#1d4ed8;--awas:#c0262b}
@media(prefers-color-scheme:dark){:root:not([data-theme="light"]){
--tinta:#e9ebf0;--tinta-2:#aeb6c3;--tinta-3:#8c94a3;--latar:#14161a;--latar-2:#1b1e24;
--garis:#2f343d;--badan:#1b1e24;--r-pohon:#7fa2ff;--r-silang:#ff8a8a;
--akar:#2a4fa8;--akar-t:#fff;--bersama:#8a6a1c;--bersama-t:#ffe9b8;
--np:#22375c;--np-t:#cfe0fb;--komite:#3d2b57;--komite-t:#e4d3f7;--migrasi:#2b303a;--migrasi-t:#cdd3dc;
--pk:#e0a352;--fk:#8fabff;--awas:#ff8a8a}}
*{box-sizing:border-box}body{margin:0;background:var(--latar);color:var(--tinta);
font:15px/1.55 ui-sans-serif,-apple-system,"Segoe UI",Roboto,Helvetica,Arial,sans-serif}
header{border-bottom:1px solid var(--garis);padding:22px 18px 16px}.bingkai{padding:0 18px 56px}
h1{font-size:21px;margin:0 0 8px}h2{font-size:16px;margin:26px 0 8px}
p.ket{color:var(--tinta-2);font-size:13px;max-width:92ch;margin:0 0 8px}
.awas{display:inline-block;margin-top:8px;padding:6px 10px;border-radius:4px;background:var(--latar-2);
border-left:3px solid var(--awas);font-size:13px}
code{font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;font-size:.92em}
.legenda{display:flex;flex-wrap:wrap;gap:7px 18px;margin:10px 0;font-size:12.4px;color:var(--tinta-2)}
.legenda div{display:flex;align-items:center;gap:6px}
.kot{width:14px;height:14px;border-radius:3px;display:inline-block}
.gar{width:30px;height:0;border-top-width:2px;border-top-style:solid;display:inline-block}
.bungkus{overflow:auto;border:1px solid var(--garis);border-radius:6px;background:var(--latar-2);padding:8px}
.erd{display:block;min-width:1740px}
.erd .badan{fill:var(--badan);stroke-width:1.2}.erd .kepala{stroke:none}
.erd .t-nama{font-size:12.4px;font-weight:700}.erd .t-kard{font-size:10.4px;text-anchor:end;opacity:.85}
.erd text{font-size:10.2px}
.erd .b-pk{fill:var(--pk);font-weight:700}.erd .b-fk{fill:var(--fk);font-weight:600}
.erd .b-awas{fill:var(--awas)}.erd .b-asal{fill:var(--tinta-3);font-style:italic}
.erd .b-ket{fill:var(--tinta-2)}
.erd .rel{fill:none;stroke-width:1.3;opacity:.85}.erd .rel:hover{stroke-width:2.8;opacity:1}
.erd .r-pohon{stroke:var(--r-pohon);color:var(--r-pohon)}
.erd .r-silang{stroke:var(--r-silang);color:var(--r-silang);stroke-dasharray:6 4}
.w-akar .badan{stroke:var(--akar)}.w-akar .kepala{fill:var(--akar)}
.w-akar .t-nama,.w-akar .t-kard{fill:var(--akar-t)}
.w-bersama .badan{stroke:#b8860b}.w-bersama .kepala{fill:var(--bersama)}
.w-bersama .t-nama,.w-bersama .t-kard{fill:var(--bersama-t)}
.w-np .badan{stroke:#3f79c4}.w-np .kepala{fill:var(--np)}.w-np .t-nama,.w-np .t-kard{fill:var(--np-t)}
.w-komite .badan{stroke:#8e5cc0}.w-komite .kepala{fill:var(--komite)}
.w-komite .t-nama,.w-komite .t-kard{fill:var(--komite-t)}
.w-migrasi .badan{stroke:#69717f;stroke-dasharray:5 3}.w-migrasi .kepala{fill:var(--migrasi)}
.w-migrasi .t-nama,.w-migrasi .t-kard{fill:var(--migrasi-t)}
table{border-collapse:collapse;width:100%;font-size:12.6px}
th{position:sticky;top:0;background:var(--latar-2);text-align:left;font-weight:600;padding:6px 10px;
border-bottom:1px solid var(--garis);color:var(--tinta-2)}
td{padding:4px 10px;border-bottom:1px solid var(--garis);vertical-align:top}
td.num,th.num{text-align:right;font-variant-numeric:tabular-nums}
.tbl{max-height:60vh;overflow:auto;border:1px solid var(--garis);border-radius:6px}
footer{color:var(--tinta-3);font-size:12px;margin-top:26px;border-top:1px solid var(--garis);padding-top:12px}
</style></head><body>
<header><h1>Rancangan tabel <b>Claim Non Prop</b> &mdash; hanya PK dan FK</h1>
<p class="ket"><b>STEMPEL ASAL.</b> Akar <code>T_WORK_CLAIM</code> &middot; __N_TABEL__ tabel
&middot; __N_RELASI__ relasi. Bentuk pohon dari <code>struktur-claimdata-lama.md</code> /
<code>datar-claimdata-lama.csv</code> (489 simpul, 279 berkas XML); arti bisnis dari
<code>MEMORI_PEMAHAMAN.MD</code> bagian 4.2&ndash;4.3; PK/UNIQUE lama dari
<code>pengetahuan/SCHEMA-ACTUAL.csv</code>; konvensi penamaan dan bentuk dari
<code>Diagram-Skema-Tabel-NusantaraRe.xlsx</code> sheet <i>Claim Prop</i>.
Dibangkitkan <code>alat/buat-skema-claimnp.py</code> pada <b>__TANGGAL__</b>.</p>
<p class="awas"><b>Ini rancangan, bukan bentuk yang sudah berdiri.</b> Kolom selain PK dan FK
sengaja tidak digambar &mdash; sesuai konvensi workbook. Baris <code>&equiv;</code> di tiap kotak
menunjuk padanannya di <code>ddl-usulan/</code>, supaya rancangan ini <b>tidak menjadi rancangan
kedua yang berdiri sendiri</b>. Folder <code>Komite Claim Non Prop</code> tidak dibuka.</p>
</header><div class="bingkai">
<div class="legenda">
<div><span class="kot" style="background:var(--akar)"></span> akar lintas-lini</div>
<div><span class="kot" style="background:var(--np)"></span> tabel Claim Non Prop <code>T_CLAIM_</code></div>
<div><span class="kot" style="background:var(--bersama)"></span> lintas-lini (dipakai lini lain juga)</div>
<div><span class="kot" style="background:var(--komite)"></span> tabel Komite</div>
<div><span class="kot" style="background:var(--migrasi)"></span> tabel migrasi, berumur terbatas</div>
<div><span class="gar" style="border-color:var(--r-pohon)"></span> induk &rarr; anak (pohon)</div>
<div><span class="gar" style="border-color:var(--r-silang);border-top-style:dashed"></span> penunjuk silang 1:1</div>
</div>
<div class="bungkus">__SVG__</div>
<p class="ket">Baris <code>&larr;</code> = halaman Pega asalnya. Baris <code>&equiv;</code> =
padanannya di <code>ddl-usulan/</code>. Baris <code>&#9888;</code> = hal yang membedakan rancangan
ini dari sheet <i>Claim Prop</i>, atau yang masih terbuka.</p>

<h2>Daftar relasi &mdash; __N_RELASI__</h2>
<div class="tbl"><table><thead><tr><th class="num">#</th><th>Induk</th><th>Anak</th>
<th>Kunci tamu</th><th>Kard.</th></tr></thead><tbody>__RELASI__</tbody></table></div>

<h2>Sepuluh yang bukan tabel</h2>
<div class="tbl"><table><thead><tr><th>Halaman Pega</th><th>Sebab</th><th>Di Go</th></tr></thead>
<tbody>__BUKAN__</tbody></table></div>

<footer><b>TURUNAN.</b> Bila halaman ini berbeda dari sumbernya, <b>alatnya yang salah</b>.
Bentuk Excel-nya: <code>Diagram-Skema-Tabel-ClaimNonProp.xlsx</code>.</footer>
</div></body></html>"""

HTML = (HTML.replace('__SVG__', SVG).replace('__TANGGAL__', TANGGAL)
        .replace('__N_TABEL__', str(len(T))).replace('__N_RELASI__', str(len(RELASI)))
        .replace('__RELASI__', baris_relasi).replace('__BUKAN__', baris_bukan))
io.open(os.path.join(ERD, 'ERD-CLAIM-NON-PROP.html'), 'w', encoding='utf-8',
        newline='\n').write(HTML)

# =====================================================================
# 4. Excel - kotak ERD, teknik persis Diagram-Skema-Tabel-NusantaraRe.xlsx
# =====================================================================
# Bukan satu sel per baris. Tekniknya: seluruh kolom selebar 2,3 dan tiap
# baris kotak adalah SATU MERGE melintasi LEBAR_KOTAK kolom sempit, dengan
# border medium di kiri dan kanan, medium di atas pada kepala, dan medium di
# bawah pada baris terakhir. Itulah yang menggambar garis kotaknya.
LEBAR_KOTAK = 32                 # kolom sempit per kotak, sepadan D13:AG13
LEBAR_KOLOM = 2.3
TINGGI_BARIS = 15.0

# Palet diambil apa adanya dari workbook contoh.
ISI_KEPALA = {'akar': '1F4E79', 'bersama': 'C55A11', 'np': '2E75B6',
              'komite': '7030A0', 'migrasi': '808080'}
ISI_PK = 'FFF2CC'
ISI_BADAN = 'EDF3F9'
TINTA_REDUP = '7F7F7F'
TINTA_AWAS = 'C00000'
TINTA_FK = '1F4E79'

MEDIUM = Side(style='medium', color='000000')
TEPI = Border(*[Side(style='thin', color='9AA2B1')] * 4)


def sisi(kiri=False, kanan=False, atas=False, bawah=False):
    return Border(left=MEDIUM if kiri else None, right=MEDIUM if kanan else None,
                  top=MEDIUM if atas else None, bottom=MEDIUM if bawah else None)


wb = openpyxl.Workbook()
ws = wb.active
ws.title = 'Claim Non Prop'
ws.sheet_view.showGridLines = False

MAKS_KOL = 100
for i in range(1, MAKS_KOL + 1):
    ws.column_dimensions[openpyxl.utils.get_column_letter(i)].width = LEBAR_KOLOM


def merge_tulis(r, c, teks, isi, tinta, ukuran, tebal, rata, kiri, kanan, atas, bawah,
                lebar=LEBAR_KOTAK):
    ws.merge_cells(start_row=r, start_column=c, end_row=r, end_column=c + lebar - 1)
    sel = ws.cell(r, c, teks)
    sel.font = Font(name='Arial', size=ukuran, bold=tebal, color=tinta)
    if isi:
        sel.fill = PatternFill('solid', fgColor=isi)
    sel.alignment = Alignment(horizontal=rata, vertical='center')
    for k in range(c, c + lebar):
        ws.cell(r, k).border = sisi(kiri=(k == c and kiri), kanan=(k == c + lebar - 1 and kanan),
                                    atas=atas, bawah=bawah)
    ws.row_dimensions[r].height = TINGGI_BARIS
    return r + 1


def tulis_kotak(nama, r0, c0):
    """Gambar satu kotak entitas. Kembalikan baris pertama SESUDAH kotak."""
    t = PETA[nama]
    merge_tulis(r0, c0, u'%s        %s' % (nama, t['kard']),
                ISI_KEPALA[t['kelompok']], 'FFFFFF', 9, True, 'center',
                True, True, True, False)
    r = r0 + 1
    isi_baris = baris_kotak(t)
    for i, teks in enumerate(isi_baris):
        akhir = (i == len(isi_baris) - 1)
        if teks.startswith('PK'):
            isi, tinta, tebal = ISI_PK, '000000', True
        elif teks.startswith('FK'):
            isi, tinta, tebal = ISI_BADAN, TINTA_FK, True
        elif teks.lstrip().startswith((u'⚠', '[terbuka]')):
            isi, tinta, tebal = ISI_BADAN, TINTA_AWAS, False
        elif teks.startswith((u'←', u'≡')):
            isi, tinta, tebal = ISI_BADAN, TINTA_REDUP, False
        else:
            isi, tinta, tebal = ISI_BADAN, '000000', False
        r = merge_tulis(r, c0, u'  ' + teks, isi, tinta, 8, tebal, 'left',
                        True, True, False, akhir)
    return r


def panah(r, c0, label):
    sel = ws.cell(r, c0, u'▼  %s' % label)
    sel.font = Font(name='Arial', size=8, bold=True, color=TINTA_FK)
    ws.row_dimensions[r].height = TINGGI_BARIS
    return r + 1


# --- judul -----------------------------------------------------------
merge_tulis(1, 2, u'CLAIM NON PROP — rancangan tabel · hanya PK dan FK',
            None, '000000', 14, True, 'left', False, False, False, False, lebar=60)
merge_tulis(2, 2,
            u'Akar T_WORK_CLAIM · %d tabel · %d relasi · sumber: '
            u'struktur-claimdata-lama.md · MEMORI_PEMAHAMAN.MD 4.2–4.3 · '
            u'pengetahuan/ · ddl-usulan/ · dibangkitkan %s'
            % (len(T), len(RELASI), TANGGAL),
            None, '7F7F7F', 9, False, 'left', False, False, False, False, lebar=76)

# --- tata letak tangga, sama dengan sheet Claim Prop -----------------
KOL_TINGKAT = {0: 2, 1: 4, 2: 6, 3: 12, 4: 18, 5: 24}
def akar_dari(n):
    """Leluhur teratas sebuah tabel. Yang tidak sampai ke T_WORK_CLAIM tidak
    masuk tangga - ia turun ke bagian DI LUAR POHON bersama induknya."""
    lihat = set()
    while PETA[n]['induk'] and n not in lihat:
        lihat.add(n)
        n = PETA[n]['induk']
    return n


tingkat = collections.defaultdict(list)
for nama in PETA:
    if akar_dari(nama) == 'T_WORK_CLAIM':
        tingkat[min(kedalaman(nama), 5)].append(nama)
for d in tingkat:
    tingkat[d].sort(key=lambda n: (0 if n == 'T_WORK_CLAIM' else 1, n))

r = 4
for d in sorted(tingkat):
    c0 = KOL_TINGKAT[d]
    for i, nama in enumerate(tingkat[d]):
        t = PETA[nama]
        if t['induk']:
            r = panah(r, c0, u'%s   ← %s' % (t['kunci'], t['induk']))
        r = tulis_kotak(nama, r, c0)
        r += 1
    r += 1

# --- tabel tanpa induk, seperti bagian DI LUAR POHON -----------------
# tabel tanpa induk, BESERTA keturunannya - supaya panahnya tidak menunjuk
# kotak yang berdiri jauh di bagian lain.
bebas = [n for n in sorted(PETA) if akar_dari(n) != 'T_WORK_CLAIM']
bebas.sort(key=lambda n: (kedalaman(n), n))
if bebas:
    r += 1
    merge_tulis(r, 2, u'DI LUAR POHON — tanpa induk, tanpa foreign key',
                None, '000000', 11, True, 'left', False, False, False, False, lebar=60)
    r += 2
    for nama in bebas:
        t_ = PETA[nama]
        c0 = 2 + 2 * kedalaman(nama)
        if t_['induk']:
            r = panah(r, c0, u'%s   ← %s' % (t_['kunci'], t_['induk']))
        r = tulis_kotak(nama, r, c0) + 1

# --- sepuluh yang bukan tabel ----------------------------------------
r += 2
merge_tulis(r, 2, u'SEPULUH YANG BUKAN TABEL', None, '000000', 11, True, 'left',
            False, False, False, False, lebar=60)
r += 1
for a_, b_, c_ in BUKAN_TABEL:
    merge_tulis(r, 2, a_, None, '000000', 9, False, 'left', False, False, False, False, lebar=26)
    merge_tulis(r, 29, b_, None, '555555', 9, False, 'left', False, False, False, False, lebar=24)
    merge_tulis(r, 54, c_, None, '555555', 9, False, 'left', False, False, False, False, lebar=24)
    r += 1

# --- legenda ---------------------------------------------------------
r += 2
merge_tulis(r, 2, u'LEGENDA', None, '000000', 11, True, 'left',
            False, False, False, False, lebar=20)
r += 1
for warna, lbl, ket in (
        (ISI_KEPALA['akar'], u'Kotak biru tua', u'akar lintas-lini T_WORK_CLAIM'),
        (ISI_KEPALA['np'], u'Kotak biru', u'tabel Claim Non Prop — awalan T_CLAIM_'),
        (ISI_KEPALA['bersama'], u'Kotak jingga', u'tabel lintas-lini, dipakai lini lain juga'),
        (ISI_KEPALA['komite'], u'Kotak ungu', u'tabel Komite'),
        (ISI_KEPALA['migrasi'], u'Kotak abu', u'tabel migrasi — berumur terbatas'),
        (ISI_PK, u'Baris kuning', u'kunci utama (PK)'),
        (None, u'▼', u'hubungan induk → anak, beserta nama kunci tamunya'),
        (None, u'←', u'halaman Pega asalnya'),
        (None, u'≡', u'padanannya di ddl-usulan/'),
        (None, u'⚠', u'beda dari sheet Claim Prop, atau masih terbuka')):
    merge_tulis(r, 2, u'  ' + lbl, warna, 'FFFFFF' if warna and warna != ISI_PK else '000000',
                9, True, 'left', bool(warna), bool(warna), bool(warna), bool(warna), lebar=12)
    merge_tulis(r, 15, ket, None, '555555', 9, False, 'left',
                False, False, False, False, lebar=40)
    r += 1


# --- sheet Daftar Relasi ---------------------------------------------
ws2 = wb.create_sheet('Daftar Relasi')
ws2.sheet_view.showGridLines = False
ws2['A1'] = u'DAFTAR RELASI — Claim Non Prop · PK dan FK saja'
ws2['A1'].font = Font(bold=True, size=12)
for j, h in enumerate([u'#', u'Skema', u'Induk', u'Anak', u'Kunci tamu', u'Kard.'], 1):
    c = ws2.cell(3, j, h)
    c.font = Font(bold=True, size=10)
    c.fill = PatternFill('solid', fgColor='E8ECF2')
    c.border = TEPI
for i, (a, b, k, kd) in enumerate(RELASI, 1):
    ws2.cell(i + 3, 1, i).font = Font(size=9)
    ws2.cell(i + 3, 2, u'Claim Non Prop').font = Font(size=9)
    ws2.cell(i + 3, 3, a).font = Font(size=9)
    ws2.cell(i + 3, 4, b).font = Font(size=9)
    ws2.cell(i + 3, 5, k).font = Font(size=9)
    ws2.cell(i + 3, 6, kd).font = Font(size=9)
r = len(RELASI) + 5
for cat in (u'Seluruh kunci tamu WAJIB ber-index. Yang bertanda UNIK wajib index UNIK.',
            u'SHARED PK: T_GENERAL_CLAIM.ID dan T_GENERAL_KOMITE.ID = T_WORK_CLAIM.ID — tidak ada kolom FK terpisah.',
            u'T_VIEW_SUGGEST punya DUA induk — PREMIUM_LIST_ID dan CLAIM_ID, keduanya nullable, dengan CHECK tepat satu terisi.',
            u'⚠ T_CLAIM_ADJ_LAYER dan T_CLAIM_ADJ_LAYER_CURRENCY BELUM punya padanan di 22 tabel ddl-usulan/.',
            u'⚠ T_CLAIM_RETRO punya DUA induk mungkin — CNPSpreadLoss[] dan SpreadingRisk[]. [terbuka]',
            u'⚠ Folder Komite Claim Non Prop TIDAK dibuka; bentuk tabel komite diambil dari sheet Komite Claim Prop.',
            u'⚠ OS_AKSEPTASI_KLAIM lama tidak punya primary key maupun unique constraint — REQ-018.'):
    ws2.cell(r, 3, cat).font = Font(size=9, color='555555')
    r += 1
for c_, w_ in (('A', 5), ('B', 16), ('C', 32), ('D', 34), ('E', 26), ('F', 8)):
    ws2.column_dimensions[c_].width = w_

# --- sheet Peta Sumber -----------------------------------------------
ws3 = wb.create_sheet('Peta Sumber')
ws3.sheet_view.showGridLines = False
ws3['A1'] = u'PETA SUMBER — tabel baru ← halaman Pega ← padanan ddl-usulan/'
ws3['A1'].font = Font(bold=True, size=12)
for j, h in enumerate([u'Tabel baru', u'Kelompok', u'Kard.', u'Induk', u'Kunci tamu',
                       u'Asal di Pega', u'Padanan ddl-usulan/'], 1):
    c = ws3.cell(3, j, h)
    c.font = Font(bold=True, size=10)
    c.fill = PatternFill('solid', fgColor='E8ECF2')
    c.border = TEPI
for i, t in enumerate(T, 1):
    for j, v in enumerate([t['nama'], t['kelompok'], t['kard'], t['induk'] or u'—',
                           t['kunci'] or u'—', t['asal'], t['padanan']], 1):
        cc = ws3.cell(i + 3, j, v)
        cc.font = Font(size=9, color='C0262B' if u'⚠' in u'%s' % v else '000000')
        cc.alignment = Alignment(vertical='top', wrap_text=(j >= 6))
for c_, w_ in (('A', 34), ('B', 11), ('C', 8), ('D', 24), ('E', 24), ('F', 62), ('G', 26)):
    ws3.column_dimensions[c_].width = w_

# tabel ddl-usulan/ yang TIDAK punya kotak di sini, dan sebaliknya
padanan = set()
for t in T:
    for x in re.split(r'[+·,]', t['padanan']):
        x = x.strip()
        if x and x in tabel_klaimnp:
            padanan.add(x)
r = len(T) + 5
ws3.cell(r, 1, u'REKONSILIASI terhadap ddl-usulan/').font = Font(bold=True, size=11)
r += 1
ws3.cell(r, 1, u'22 tabel ddl-usulan/, %d terpadankan di sini' % len(padanan)).font = Font(size=9)
r += 1
sisa = [x for x in tabel_klaimnp if x not in padanan]
ws3.cell(r, 1, u'belum terpadankan: %s' % (u', '.join(sisa) if sisa else u'nihil')).font = \
    Font(size=9, color='C0262B' if sisa else '2E7D32')
r += 1
ws3.cell(r, 1, u'tabel baru TANPA padanan: %s'
         % u', '.join(t['nama'] for t in T if u'TIDAK ADA' in t['padanan'])).font = \
    Font(size=9, color='C0262B')

TUJUAN = os.path.join(ERD, 'Diagram-Skema-Tabel-ClaimNonProp.xlsx')
try:
    wb.save(TUJUAN)
except PermissionError:
    # Excel memegang berkasnya. Hasil kerja TIDAK dibuang dan berkas lama
    # TIDAK dibiarkan seolah-olah mutakhir - keduanya akan menyesatkan.
    TUJUAN = TUJUAN[:-5] + '.BARU.xlsx'
    wb.save(TUJUAN)
    print('')
    print('!! Diagram-Skema-Tabel-ClaimNonProp.xlsx SEDANG TERBUKA di Excel.')
    print('!! Hasil baru ditulis ke  %s' % os.path.basename(TUJUAN))
    print('!! Tutup Excel, hapus yang lama, ganti nama yang .BARU - atau jalankan ulang.')
    print('')

# =====================================================================
print('ERD-CLAIM-NON-PROP.html            %.0f KB'
      % (os.path.getsize(os.path.join(ERD, 'ERD-CLAIM-NON-PROP.html')) / 1024.0))
print('Diagram-Skema-Tabel-ClaimNonProp.xlsx  %.0f KB'
      % (os.path.getsize(TUJUAN) / 1024.0))
print('')
print('tabel                 : %d' % len(T))
for k in ('akar', 'bersama', 'np', 'komite', 'migrasi'):
    n = sum(1 for t in T if t['kelompok'] == k)
    print('  %-10s          %d' % (k, n))
print('relasi                : %d' % len(RELASI))
print('bukan tabel           : %d' % len(BUKAN_TABEL))
# --- peta nama, terbaca mesin: dipakai buat-erd-excel.py --------------
with io.open(os.path.join(ERD, 'peta-nama-tabel.tsv'), 'w', encoding='utf-8') as _f:
    _f.write(u'NAMA_T\tPADANAN_DDL\tKELOMPOK\tINDUK\tKUNCI\tASAL_PEGA\n')
    for _t in T:
        _f.write(u'%s\t%s\t%s\t%s\t%s\t%s\n'
                 % (_t['nama'], _t['padanan'], _t['kelompok'],
                    _t['induk'] or u'', _t['kunci'] or u'', _t['asal']))
print('peta-nama-tabel.tsv                %d baris' % len(T))

print('')
print('REKONSILIASI terhadap struktur-claimdata-lama.md')
print('  simpul di pohon      : %d' % len(SIMPUL_POHON))
print('  dikutip rancangan    : %d' % len(dikutip))
print('  dikutip tapi TAK ADA : %s' % (', '.join(tak_ada) if tak_ada else 'nihil'))
print('  angka bagian 2       : %d simpul · %d skalar akar · %d wadah akar'
      % (N_SIMPUL, N_SKALAR_AKAR, N_WADAH_AKAR))
print('')
print('REKONSILIASI terhadap ddl-usulan/')
print('  22 tabel, terpadankan: %d' % len(padanan))
print('  belum terpadankan    : %s' % (', '.join(sisa) if sisa else 'nihil'))
print('  baru tanpa padanan   : %s'
      % ', '.join(t['nama'] for t in T if u'TIDAK ADA' in t['padanan']))

# --- rekonsiliasi ketiga: terhadap nama yang BERLAKU di berkas rujukan --
# Diagram-Skema-Tabel-NusantaraRe.xlsx diubah work owner 20-09-2026:
# awalan T_CLAIMP_ diganti T_CLAIM_, dan sheet "Claim Fac In" ditambahkan.
# Rancangan ini WAJIB memakai nama yang sama persis untuk tabel yang sama.
import openpyxl as _ox
_RUJ = os.path.join(ERD, 'Diagram-Skema-Tabel-NusantaraRe.xlsx')
_w = _ox.load_workbook(_RUJ, read_only=True, data_only=True)
_ada, _induk = set(), {}
for _sh in _w:
    for _row in _sh.iter_rows(values_only=True):
        for _v in _row:
            if isinstance(_v, str):
                for _m in re.findall(r'T_CLAIM_[A-Z0-9_]+', _v):
                    _ada.add(_m)
_ws = _w['Daftar Relasi']
for _row in _ws.iter_rows(min_row=3, values_only=True):
    if _row and isinstance(_row[0], (int, float)) and _row[2] and _row[3]:
        if str(_row[3]).startswith('T_CLAIM_'):
            _induk.setdefault(str(_row[3]), set()).add((str(_row[2]), str(_row[4] or '')))
_w.close()

_milik = set(t['nama'] for t in T if t['nama'].startswith('T_CLAIM_'))
_bersama = sorted(_milik & _ada)
_khas = sorted(_milik - _ada)
_hanya_rujukan = sorted(_ada - _milik)
_beda_induk = []
for _n in _bersama:
    _t = next(x for x in T if x['nama'] == _n)
    _pasangan = (_t['induk'], _t['kunci'])
    if _n in _induk and _pasangan not in _induk[_n]:
        _beda_induk.append('%s: di sini %s.%s, di rujukan %s'
                           % (_n, _t['induk'], _t['kunci'],
                              ' / '.join('%s.%s' % x for x in sorted(_induk[_n]))))
print('')
print('REKONSILIASI terhadap Diagram-Skema-Tabel-NusantaraRe.xlsx (keadaan 20-09-2026)')
print('  T_CLAIM_ di rujukan  : %d  %s' % (len(_ada), ', '.join(sorted(_ada))))
print('  T_CLAIM_ di sini     : %d' % len(_milik))
print('  DIPAKAI BERSAMA      : %d  %s' % (len(_bersama), ', '.join(_bersama)))
print('  khas Non Prop        : %d  %s' % (len(_khas), ', '.join(_khas)))
print('  ada di rujukan saja  : %d  %s' % (len(_hanya_rujukan),
                                           ', '.join(_hanya_rujukan) or 'nihil'))
print('  induk BERBEDA        : %s' % ('; '.join(_beda_induk) if _beda_induk else 'nihil'))
sys.exit(0)
