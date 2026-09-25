# -*- coding: utf-8 -*-
r"""
saring-kemampuan-ke-tiket.py -- menyaring daftar kemampuan menjadi bahan to-ticket.

APA YANG DAPAT DIBUKTIKAN PERKAKAS INI
    Bahwa setiap baris kemampuan di DAFTAR-PEKERJAAN.md terbaca, tergolongkan, dan
    TIDAK ADA YANG HILANG -- jumlah golongan selalu menutup ke jumlah baris.

APA YANG TIDAK
    Ia TIDAK mengadili. Penggolongan G2 dan batch REV-3 di bawah adalah SAPUAN KATA
    atas kolom Kemampuan / Asal / Entitas, dan sapuan kata TIDAK DAPAT membaca arti.
    Hasilnya CALON, bukan putusan; setiap calon diperiksa dengan mata dan yang
    berbeda dicatat sebagai koreksi bernama. Batasnya dicetak di dalam keluarannya.

PENOLAKAN dilaporkan satu baris per sebab.

Pakai:  PYTHONIOENCODING=utf-8 python alat/saring-kemampuan-ke-tiket.py
"""
import io, os, re, json
from collections import Counter, OrderedDict

AKAR = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
DAFTAR = os.path.join(AKAR, '5-tiket', 'DAFTAR-PEKERJAAN.md')
ISSUES = os.path.join(AKAR, '5-tiket', 'issues')

teks = io.open(DAFTAR, encoding='utf-8').read()
baris = teks.split('\n')

# -- 1. baca baris kemampuan ------------------------------------------------
# bentuk: | P-nn | kemampuan | gol | asal | entitas | keadaan |
ROW = re.compile(r'^\|\s*(~~)?(P-\d+)(~~)?\s*\|(.*)$')
tolak = Counter()
kem = OrderedDict()
bagian = ''
for l in baris:
    if l.startswith('### '):
        bagian = l[4:].strip()
    m = ROW.match(l)
    if not m:
        continue
    sel = [c.strip() for c in l.strip().strip('|').split('|')]
    if len(sel) < 6:
        tolak['baris P-nn berkolom kurang dari enam'] += 1
        continue
    kode = m.group(2)
    if kode in kem:
        tolak['baris P-nn kembar (baris kedua diabaikan)'] += 1
        continue
    kem[kode] = dict(kode=kode, bagian=bagian, dicoret=bool(m.group(1)),
                     kemampuan=sel[1], gol=sel[2], asal=sel[3],
                     entitas=sel[4], keadaan=sel[5])

# -- 2. lingkup ronde ini: P-01..P-59, yang tidak dicoret --------------------
def nomor(k):
    return int(k.split('-')[1])

lingkup, luar = [], Counter()
for k, v in kem.items():
    if v['dicoret']:
        luar['DICORET sebagai SIFAT (U-6)'] += 1
        continue
    if nomor(k) > 59:
        luar['kemampuan Adjustment P-60..P-66 -- sudah ditiketkan ronde lalu'] += 1
        continue
    lingkup.append(v)

# -- 3. saringan G2 -- CALON, bukan putusan ---------------------------------
# Tepi yang KEPUTUSAN-PEMBAGIAN-TIKET.md sec 2 "Yang TEGAS TIDAK masuk" sebutkan.
G2_TEPI = [
    ('GEL-2 RETRO_KELUAR',      lambda v: 'RETRO_KELUAR' in v['entitas'] or 'retro keluar' in v['kemampuan'].lower()),
    ('GEL-3 PENCAPAIAN',        lambda v: 'PENCAPAIAN' in v['entitas']),
    ('penerbitan ke luar',      lambda v: 'TREATYINOFFER' in (v['kemampuan'] + v['asal'] + v['entitas']).upper()),
    ('hitung ulang / perbaikan data lama (ADR-0043)',
                                lambda v: 'ADR-0043' in v['asal'] and 'hitung ulang' in v['kemampuan'].lower()),
    ('invarian AccountingMode -- menunggu Uji X-2',
                                lambda v: 'X-2' in v['keadaan'] or 'AccountingMode' in v['kemampuan']),
]

# -- 4. batch REV-3 -- CALON. Kata yang menandai sentuhan ke DAFTAR KEADAAN. --
KATA_KEADAAN = ('ADR-0055', 'ADR-0044', 'ADR-0052', 'keadaan', 'KEADAAN',
                'setuju', 'Setuju', 'DISETUJUI', 'persetujuan', 'menyetujui',
                'tolak', 'Tolak', 'DITOLAK', 'menolak versi', 'batal', 'Batal',
                'DIBATALKAN', 'DRAFT', 'AJUKAN', 'ajukan', 'mengajukan',
                'siklus hidup')

def sentuh_keadaan(v):
    ladang = ' '.join((v['kemampuan'], v['asal'], v['entitas']))
    kena = [k for k in KATA_KEADAAN if k in ladang]
    return kena

for v in lingkup:
    v['g2_tepi'] = [n for n, f in G2_TEPI if f(v)]
    v['kata_keadaan'] = sentuh_keadaan(v)

# -- 5. tiket yang sudah ada -- kemampuan mana yang dirujuknya ---------------
sudah = Counter()
peta_tiket = {}
if os.path.isdir(ISSUES):
    for f in sorted(os.listdir(ISSUES)):
        if not f.endswith('.md') or f == 'README.md':
            continue
        t = io.open(os.path.join(ISSUES, f), encoding='utf-8').read()
        for k in set(re.findall(r'`(P-\d+)`', t)):
            sudah[k] += 1
            peta_tiket.setdefault(k, []).append(f[:2])
else:
    tolak['folder issues tidak ada'] += 1

# -- 6. cetak ---------------------------------------------------------------
print('=== saring-kemampuan-ke-tiket.py ===')
print('baris kemampuan terbaca : %d' % len(kem))
print()
print('DITOLAK / DIKELUARKAN -- satu baris per sebab:')
for k, n in list(luar.items()) + list(tolak.items()):
    print('   %-62s %d' % (k, n))
print('   %-62s %d' % ('SISA -- lingkup ronde ini', len(lingkup)))
print()

kena_g2 = [v for v in lingkup if v['g2_tepi']]
print('SARINGAN G2 -- CALON tepi "tegas tidak masuk" : %d' % len(kena_g2))
for v in kena_g2:
    print('   %-6s %-46s <- %s' % (v['kode'], v['kemampuan'][:46], ', '.join(v['g2_tepi'])))
if not kena_g2:
    print('   (nol)')
print()

b2 = [v for v in lingkup if v['kata_keadaan']]
b1 = [v for v in lingkup if not v['kata_keadaan']]
print('CALON BATCH -- sapuan kata, BUKAN adjudikasi:')
print('   batch 2 (menyentuh keadaan/persetujuan/penolakan/pembatalan) : %d' % len(b2))
print('   batch 1 (tidak menyentuh)                                    : %d' % len(b1))
print('   menutup ke lingkup : %s' % ('YA' if len(b1) + len(b2) == len(lingkup) else 'TIDAK'))
print()
print('SUDAH DIRUJUK tiket yang ada:')
for k in sorted(sudah, key=nomor):
    if nomor(k) <= 59:
        print('   %-6s tiket %s' % (k, ', '.join(peta_tiket[k])))
print()
print('-- BATAS PERKAKAS INI, dan ia bagian dari keluarannya --')
print('   Kedua penggolongan di atas SAPUAN KATA. Ia tidak dapat membaca arti, dan')
print('   akan salah pada dua arah: kemampuan yang MENYEBUT keadaan tanpa bergantung')
print('   padanya, dan kemampuan yang BERGANTUNG pada keadaan tanpa menyebutnya.')
print('   Angka di atas CALON. Yang berlaku adalah adjudikasi bermata, dan tiap')
print('   koreksi terhadap angka ini ditulis bernama di USULAN-IRISAN-*.md.')

json.dump({'lingkup': lingkup, 'luar': dict(luar), 'tolak': dict(tolak),
           'sudah': dict(sudah), 'peta_tiket': peta_tiket},
          io.open(os.path.join(AKAR, 'alat', 'saring-kemampuan.json'), 'w', encoding='utf-8'),
          ensure_ascii=False, indent=1)
