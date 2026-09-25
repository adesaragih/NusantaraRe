# -*- coding: utf-8 -*-
r"""
urai-erd-v2.py -- mengurai Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx APA ADANYA
menjadi bahan mesin-baca, tanpa menambah dan tanpa mengurangi satu baris pun.

    masukan : 4-erd-dan-tabel-datar/Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx
    keluaran: alat/erd-v2.json

KENAPA ADA
    Pemilik proses menetapkan berkas itu sebagai SUMBER struktur data, "PERSIS,
    TIDAK ADA PERUBAHAN". Maka struktur data, tabel datar, dan ERD HTML
    ketiga-tiganya harus dibangkitkan DARI berkas itu -- bukan ditulis ulang
    dengan tangan, sebab yang ditulis tangan menyimpang dan yang dibangkitkan
    tidak.

APA YANG DAPAT DIBUKTIKAN PERKAKAS INI
    Bahwa tiap kotak, tiap kunci, tiap baris asal, dan tiap baris Daftar Relasi
    di dalam berkas itu TERBACA dan TERSALIN. Cacahnya dicetak per lembar.

APA YANG TIDAK
    Ia TIDAK menilai isinya. Berkas sumbernya BERBANNER
    "POTRET SISTEM LAMA -- BUKAN RANCANGAN", dan perkakas ini MENYALIN banner itu
    apa adanya ke dalam keluarannya, bukan membuangnya.

    Ia juga TIDAK mendamaikannya dengan 2-to-spec/. Berkas itu memuat 44 tabel
    berawalan T_; ddl-usulan/ memuat 35 entitas bernama lain. Keduanya dibiarkan
    berdiri masing-masing, dan selisihnya bukan urusan perkakas ini.

BENTUK YANG DIBACA -- diukur dari berkasnya, bukan diandaikan
    kotak      = TIGA baris berurutan, tiap baris satu sel terlebur 46 kolom
    kedalaman  = (kolom_awal - 2) / 4
    baris 1    = "NAMA        <kardinalitas> . <sekat> . <catatan>"
    baris 2    = "PK   KOL"  atau  "FK   KOL -> INDUK.KOL     ON DELETE"
    baris 3    = "<- <jalur Pega>   [<padanan model baru>]   bukti: <...>"

PENOLAKAN dilaporkan satu baris per sebab.

Pakai:  PYTHONIOENCODING=utf-8 python alat/urai-erd-v2.py
"""
import io, os, re, json, sys
from collections import Counter, OrderedDict

from openpyxl import load_workbook

AKAR = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
ERD = os.path.join(AKAR, '4-erd-dan-tabel-datar')
SUMBER = os.path.join(ERD, 'Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx')
KELUAR = os.path.join(AKAR, 'alat', 'erd-v2.json')

LEMBAR_GAMBAR = ['Treaty In Prop', 'Treaty In Non Prop',
                 'Treaty In EDM Prop', 'Treaty In EDM Non Prop']
tolak = Counter()
contoh = []

wb = load_workbook(SUMBER)
hilang = [x for x in LEMBAR_GAMBAR if x not in wb.sheetnames]
if hilang:
    print('BERHENTI: lembar tidak ada -- %s' % ', '.join(hilang))
    sys.exit(2)

banner = [str(wb[LEMBAR_GAMBAR[0]].cell(1, 1).value or ''),
          str(wb[LEMBAR_GAMBAR[0]].cell(2, 1).value or '')]


def sel(ws, r, c):
    v = ws.cell(r, c).value
    return '' if v is None else str(v).strip()


def awal_kolom(ws, r):
    """Kolom pertama yang berisi di baris ini -- ia menentukan kedalaman kotak."""
    for c in range(1, 60):
        if sel(ws, r, c):
            return c
    return 0


lembar = OrderedDict()
semua_kotak = {}

for nama in LEMBAR_GAMBAR:
    ws = wb[nama]
    judul = sel(ws, 3, 2)
    subjudul = sel(ws, 4, 2)
    kotak, r = [], 5
    while r <= ws.max_row:
        c = awal_kolom(ws, r)
        if not c:
            r += 1
            continue
        b1 = sel(ws, r, c)
        b2 = sel(ws, r + 1, awal_kolom(ws, r + 1) or c)
        b3 = sel(ws, r + 2, awal_kolom(ws, r + 2) or c)
        if not (b2.startswith('PK') or b2.startswith('FK')):
            tolak['baris yang bukan kepala kotak (tidak diikuti PK/FK)'] += 1
            contoh.append(('bukan kepala kotak', '%s r%d: %s' % (nama, r, b1[:60])))
            r += 1
            continue
        if not b3.startswith('←'):
            tolak['kotak tanpa baris asal (tidak diawali panah kiri)'] += 1
            contoh.append(('tanpa baris asal', '%s r%d: %s' % (nama, r, b1[:60])))

        dalam = (c - 2) // 4
        if (c - 2) % 4:
            tolak['kotak yang kolom awalnya bukan kelipatan 4 dari kolom B'] += 1

        # baris 1 -> nama + penanda
        bagian = re.split(r'\s{4,}', b1, 1)
        nm = bagian[0].strip()
        meta = bagian[1].strip() if len(bagian) > 1 else ''
        kard = ''
        mk = re.match(r'^(1:1|1:N|N:1|N:M)\b', meta)
        if mk:
            kard = mk.group(1)
        bersama = ''
        mb = re.search(r'BERSAMA\s*→\s*([^·]+)', meta)
        if mb:
            bersama = mb.group(1).strip()

        # baris 2 -> kunci
        pk, fk = '', None
        if b2.startswith('PK'):
            pk = b2[2:].strip()
        else:
            mf = re.match(r'^FK\s+(\S+)\s*→\s*(\S+?)\.(\S+)\s+(.*)$', b2)
            if mf:
                fk = dict(kolom=mf.group(1), induk=mf.group(2),
                          kolom_induk=mf.group(3), on_delete=mf.group(4).strip())
            else:
                tolak['baris FK yang bentuknya tidak terbaca'] += 1
                contoh.append(('FK tak terbaca', '%s r%d: %s' % (nama, r + 1, b2[:70])))

        # baris 3 -> asal Pega, padanan model baru, bukti
        jalur, padanan, bukti = '', '', ''
        if b3.startswith('←'):
            sisa = b3[1:].strip()
            mp = re.search(r'\[([^\]]*)\]', sisa)
            if mp:
                padanan = mp.group(1).strip()
                sisa = sisa[:mp.start()] + sisa[mp.end():]
            mbk = re.search(r'bukti:\s*(.*)$', sisa)
            if mbk:
                bukti = mbk.group(1).strip()
                sisa = sisa[:mbk.start()]
            jalur = re.sub(r'\s{2,}', ' ', sisa).strip()

        k = dict(nama=nm, dalam=dalam, kolom_awal=c, baris=r, kard=kard,
                 bersama=bersama, meta=meta, pk=pk, fk=fk,
                 jalur=jalur, padanan=padanan, bukti=bukti,
                 teks=[b1, b2, b3])
        kotak.append(k)
        semua_kotak.setdefault(nm, k)
        r += 3

    lembar[nama] = dict(judul=judul, subjudul=subjudul, kotak=kotak)

# ══ Daftar Relasi -- disalin apa adanya ════════════════════════════════════
relasi = []
if 'Daftar Relasi' in wb.sheetnames:
    ws = wb['Daftar Relasi']
    kepala, r0 = None, None
    for r in range(1, ws.max_row + 1):
        baris = [sel(ws, r, c) for c in range(1, 60)]
        isi = [x for x in baris if x]
        if not isi:
            continue
        if kepala is None and isi[0] == '#':
            kepala, r0 = isi, r
            continue
        if kepala is not None:
            relasi.append(OrderedDict(zip(kepala, isi + [''] * (len(kepala) - len(isi)))))
else:
    tolak['lembar Daftar Relasi tidak ada'] += 1

catatan = []
if 'Catatan & Batas' in wb.sheetnames:
    ws = wb['Catatan & Batas']
    for r in range(1, ws.max_row + 1):
        isi = [sel(ws, r, c) for c in range(1, 60)]
        isi = [x for x in isi if x]
        if len(isi) >= 2:
            catatan.append([isi[0], ' '.join(isi[1:])])
        elif len(isi) == 1 and r > 2:
            catatan.append([isi[0], ''])
else:
    tolak['lembar Catatan & Batas tidak ada'] += 1

json.dump(dict(sumber=os.path.basename(SUMBER), banner=banner,
               lembar=lembar, relasi=relasi, catatan=catatan,
               tolak=dict(tolak)),
          io.open(KELUAR, 'w', encoding='utf-8'), ensure_ascii=False, indent=1)

# ══════════════════════════ CETAK ══════════════════════════════════════════
P = print
P('=== urai-erd-v2.py ===')
P()
P('SUMBER : %s' % os.path.basename(SUMBER))
P('BANNER yang DISALIN, bukan dibuang:')
for b in banner:
    P('   %s' % b[:96])
P()
P('KOTAK terbaca per lembar:')
for nama, L in lembar.items():
    dal = Counter(k['dalam'] for k in L['kotak'])
    P('   %-24s %3d kotak   kedalaman %s'
      % (nama, len(L['kotak']), ' '.join('%d:%d' % (d, dal[d]) for d in sorted(dal))))
P('   %-24s %3d baris' % ('Daftar Relasi', len(relasi)))
P('   %-24s %3d baris' % ('Catatan & Batas', len(catatan)))
P()
P('TABEL UNIK di keempat lembar : %d' % len(semua_kotak))
berpadanan = sum(1 for k in semua_kotak.values() if k['padanan'])
P('   berpadanan model baru      : %d' % berpadanan)
P('   TANPA padanan model baru   : %d' % (len(semua_kotak) - berpadanan))
for nm, k in sorted(semua_kotak.items()):
    if not k['padanan']:
        P('        - %s' % nm)
P()
P('DITOLAK / TIDAK TERBACA -- satu baris per sebab:')
for k, v in tolak.most_common():
    P('   %-64s %d' % (k, v))
if not tolak:
    P('   (tidak ada)')
for s, c in contoh[:12]:
    P('      [%s] %s' % (s, c))
P()
P('-- BATAS PERKAKAS INI --')
P('   Ia MENYALIN, ia tidak menilai. Sumbernya berbanner POTRET SISTEM LAMA,')
P('   dan banner itu ikut tersalin ke keluarannya.')
P('   Ia TIDAK mendamaikan 44 tabel T_ ini dengan 35 entitas di ddl-usulan/.')
