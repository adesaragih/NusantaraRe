# -*- coding: utf-8 -*-
"""
telusur-jalur-ke-kolom.py — LANGKAH 1: satu baris per jalur JSONDATA.

MASUKAN
    PETA-TELUSUR-JSON.md §6   -- 667 jalur, sudah beradjudikasi empat golongan
    alat/definisi-skema-treaty-masuk.json -- kolom skema baru beserta kolom *Asal*-nya

APA YANG DAPAT DIBUKTIKAN PERKAKAS INI
    Bahwa sebuah jalur JSON PUNYA baris di tabel penelusuran, dan nasibnya tertulis.
    Untuk jalur DIPETAKAN, ia dapat menunjukkan kolom skema baru yang kolom *Asal*-nya
    menyebut ruas terakhir jalur itu.

APA YANG TIDAK
    Ia TIDAK dapat membuktikan pemetaan itu BENAR. Pencocokannya leksikal -- nama ruas
    terakhir -- dan dua ruas bernama sama di kelas berbeda akan tercocokkan ke kolom
    yang sama. Setiap pemetaan bertanda COCOK adalah CALON yang masih harus dibaca.
    Ia juga TIDAK dapat menyatakan semesta LENGKAP: pohon sumbernya berlubang 340
    properti (L-8), dan PETA-TELUSUR-JSON sendiri berlabel "semestanya kurang".

PENOLAKAN dilaporkan satu baris per sebab.
"""
import io, os, re, json
from collections import Counter, defaultdict

AKAR = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
PETA = os.path.join(AKAR, 'PETA-TELUSUR-JSON.md')
DEF  = os.path.join(AKAR, 'alat', 'definisi-skema-treaty-masuk.json')

# Pohon cermin: milik modul Adjustment, DI LUAR lingkup prompt ini.
CERMIN = ('ActualValue', 'ValueDifference', 'OLDDATA', 'ValueBeforeProrate')

# ── 1. baca keempat daftar §6 ────────────────────────────────────────────────
teks = io.open(PETA, encoding='utf-8').read()
GOLONGAN = [('DIPETAKAN', 'DIPETAKAN — '), ('DITURUNKAN', 'DITURUNKAN — '),
            ('DIBUANG', 'DIBUANG — '), ('DITUNDA', 'DITUNDA — ')]
jalur = []          # (nasib, jalur)
tolak = Counter()
for i, (nama, kepala) in enumerate(GOLONGAN):
    m = re.search(r'^### ' + kepala + r'.*?$', teks, re.M)
    if not m:
        tolak['bagian §6 tidak ditemukan: ' + nama] += 1
        continue
    mulai = m.end()
    n = re.search(r'^### ', teks[mulai:], re.M)
    seg = teks[mulai:mulai + n.start()] if n else teks[mulai:]
    for baris in seg.split('\n'):
        b = baris.strip().strip('`')
        if not b or b.startswith(('```', '#', '|', '>', '*', '-')):
            continue
        if not re.match(r'^[A-Za-z][A-Za-z0-9_\[\]\.]*$', b):
            tolak['baris di §6 yang bukan jalur'] += 1
            continue
        jalur.append((nama, b))

# ── 2. pisahkan lingkup ──────────────────────────────────────────────────────
def cermin(j):
    return j.split('[')[0].split('.')[0] in CERMIN or any(('.' + c) in j for c in CERMIN)

dalam, luar = [], []
for nasib, j in jalur:
    (luar if cermin(j) else dalam).append((nasib, j))

# ── 3. peta ruas-terakhir -> (tabel, kolom) dari definisi skema ──────────────
D = json.load(io.open(DEF, encoding='utf-8'))
asal_ke = defaultdict(list)
for e, v in D['kolom'].items():
    for c in v:
        for t in re.findall(r'[A-Za-z][A-Za-z0-9_]*', c['asal'] or ''):
            if t in ('baru', 'dan', 'atau', 'pada', 'baris', 'yang', 'sama', 'tabel',
                     'acuan', 'ADR', 'lampiran', 'salah', 'eja', 'ada', 'di', 'sumbernya'):
                continue
            asal_ke[t].append((e, c['kolom'], c['tipe'], c['kosong']))

def ruas(j):
    return re.split(r'[\.\[]', j.rstrip('[]'))[-1] or j

hasil = []
cocok = Counter()
for nasib, j in dalam:
    r = ruas(j)
    kand = asal_ke.get(r, [])
    if nasib != 'DIPETAKAN':
        hasil.append((j, nasib, '', '', '', '', ''))
        cocok[nasib] += 1
    elif len(kand) == 1:
        e, k, tp, ks = kand[0]
        hasil.append((j, nasib, e, k, tp, ks, 'COCOK-TUNGGAL'))
        cocok['DIPETAKAN·cocok tunggal'] += 1
    elif len(kand) > 1:
        e, k, tp, ks = kand[0]
        hasil.append((j, nasib, e, k, tp, ks, 'COCOK-GANDA(%d)' % len(kand)))
        cocok['DIPETAKAN·cocok ganda'] += 1
    else:
        hasil.append((j, nasib, '', '', '', '', 'TIDAK-TERCOCOKKAN'))
        cocok['DIPETAKAN·tidak tercocokkan'] += 1

print('=== telusur-jalur-ke-kolom.py ===')
print('jalur dibaca dari §6           : %d' % len(jalur))
print('  di dalam LINGKUP prompt ini  : %d' % len(dalam))
print('  DIKECUALIKAN (pohon cermin)  : %d  -> milik modul Adjustment' % len(luar))
print()
print('DITOLAK -- satu baris per sebab:')
for k, v in tolak.most_common():
    print('   %-46s %d' % (k, v))
if not tolak:
    print('   (tidak ada)')
print()
print('NASIB, di dalam lingkup:')
for k, v in cocok.most_common():
    print('   %-34s %d' % (k, v))
print('   %-34s %d' % ('JUMLAH', sum(cocok.values())))

json.dump({'dalam': hasil, 'luar': luar, 'tolak': dict(tolak)},
          io.open(os.path.join(AKAR, 'alat', 'telusur-jalur.json'), 'w', encoding='utf-8'),
          ensure_ascii=False, indent=1)


# ══════════════════════════════════════════════════════════════════════════════
# 4. GOLONGKAN yang tidak tercocokkan -- masing-masing BERSEBAB, bukan didaftar
# ══════════════════════════════════════════════════════════════════════════════
import re as _re

TABEL_DARI_DAFTAR = {
    'AccumulationList': 'PERIODE_AKUMULASI', 'CurrencyList': 'MATA_UANG_KONTRAK',
    'EGNPI': 'EGNPI', 'Portfolio': 'PORTOFOLIO', 'ReportingPeriodList': 'PERIODE_PELAPORAN',
    'Retention': 'RETENSI_CEDANT', 'CoInScale': 'SKALA_KOASURANSI',
}
SUDAH_DIBUANG = {
    'CedingStatusActive': '§14.5 — bendera status master, bukan fakta kontrak',
    'SourceStatusActive': '§14.5 — idem',
    'PositionUsername':   '§12.5 — pemegang sekarang adalah turunan, bukan fakta',
    'BrokeragePercentP':  '§13 — akhiran P = proportional; MELEBUR ke PERSEN_BROKERAGE',
}
# KTV-C koreksi 1, 24 September 2026: keempat mata uang bahaya ini BUKAN golongan C.
# Sec 10.18 menulis asal NILAI_BATAS sendiri -- "+ mata uangnya" -- jadi ia golongan B,
# dan mata uangnya kini berumah di BATAS_PER_BAHAYA.KODE_MATA_UANG. Golongan GOL-C
# dikosongkan, BUKAN dihapus: mengosongkannya mencetak nol, menghapusnya mencetak diam.
GOL_C = set()
MELEBUR_B = {'CurrencyEarthquake', 'CurrencyFloodJab', 'CurrencyFloodNat', 'CurrencyRSMD'}
# Turunan yang diadili 24 September 2026 saat menutup F-2. Tanpa baris ini keduanya
# terbaca YATIM -- "tidak ada rumah, DAN TIDAK ADA ALASAN TERTULIS" -- padahal
# alasannya sudah tertulis, dan itu dua pernyataan berbeda tentang keadaan yang sama.
TURUNAN_F2 = {
    'TreatyYear': 'diturunkan `@substring(Commencement,0,4)` (`TreatyInSetTreatyYear`); '
                  'sisa pertanyaannya dipersempit ke `T-7` — `COCOK-SILANG-CACAH-ATRIBUT.md` §4.3',
    'RevisionDate': 'jejak mesin — satu-satunya penulisnya `@CurrentDateTime()` '
                    '(`TreatyInEDMSetValue`), sekeluarga dengan `EDMDate` (ADR-0006)',
    'CoInScale': 'penampung layar tak terdeklarasi — **nol penulis**; skala ko-asuransi '
                 'yang sebenarnya daftar `CoInScale[]` → `SKALA_KOASURANSI`',
}
AGREGAT_PASANGAN = {'CurrencyEgnpiAmount': 'TotalEgnpiAmount',
                    'CurrencyInstallmentAmount': 'TotalInstallmentAmount'}
NAMA_MASTER = ('Ceding', 'LeadingReinsName', 'LeadingReinsSource')

def golong(j):
    dasar = j.rstrip('[]')
    akhir = _re.split(r'[\.\[]', dasar)[-1] or dasar
    if dasar.startswith(('ShareFacultativeReinsurers', 'FacultativeShareList', 'ShareReins')):
        return ('GEL-2', 'cabang fakultatif keluar — `RETRO_KELUAR`, di luar gelombang 1 (§2.3)')
    if dasar in TABEL_DARI_DAFTAR:
        return ('TABEL', 'jalur ini menamai **daftarnya**, bukan sebuah ruas — ia menjadi **tabel** `%s`' % TABEL_DARI_DAFTAR[dasar])
    if dasar in SUDAH_DIBUANG:
        return ('DIBUANG-ULANG', SUDAH_DIBUANG[dasar])
    if dasar in GOL_C:
        return ('GOL-C', 'paket uang **golongan C**, menunggu `T-6` (P-8)')
    if dasar in MELEBUR_B:
        return ('MELEBUR-B', 'mata uang batas per bahaya — **melebur** ke '
                             '`BATAS_PER_BAHAYA.KODE_MATA_UANG` (`KTV-C` koreksi 1)')
    if dasar in TURUNAN_F2:
        return ('TURUNAN-F2', TURUNAN_F2[dasar])
    if dasar in AGREGAT_PASANGAN:
        return ('AGREGAT', 'berpasangan dengan `%s` yang **agregat** dan tidak disimpan (§4.2)' % AGREGAT_PASANGAN[dasar])
    if akhir in ('TreatyGroup', 'ReinsName', 'BrokerName') or dasar in NAMA_MASTER:
        return ('NAMA-MASTER', 'nama dibaca dari master; hanya pengenalnya disimpan (ADR-0041)')
    if akhir in ('ClassOfBusinessID',) or dasar.endswith(('TreatyGroupList', 'ClassofBusinessList', 'ClassOfBusinessList')):
        return ('BERSUSUN', 'daftar kelas bisnis / kelompok treaty bersusun di bawah induknya — **belum punya rumah bernama**')
    return ('YATIM', '**tidak ada rumah, dan tidak ada alasan tertulis**')

tt = [h[0] for h in hasil if h[6] == 'TIDAK-TERCOCOKKAN']
gol = Counter(); rinci = defaultdict(list)
for j in tt:
    g, alasan = golong(j)
    gol[g] += 1; rinci[g].append((j, alasan))

print()
print('GOLONGAN dari %d jalur DIPETAKAN yang tidak tercocokkan:' % len(tt))
for g, n in gol.most_common():
    print('   %-16s %d' % (g, n))
for g in ('YATIM', 'BERSUSUN'):
    for j, a in rinci.get(g, []):
        print('      [%s] %s' % (g, j))

json.dump({'gol': {k: v for k, v in gol.items()},
           'rinci': {k: v for k, v in rinci.items()}},
          io.open(os.path.join(AKAR, 'alat', 'telusur-golongan.json'), 'w', encoding='utf-8'),
          ensure_ascii=False, indent=1)
