# -*- coding: utf-8 -*-
r"""
cocok-enam-sumber-struktur.py -- GERBANG 0 DIPERLEBAR: enam sumber struktur data
diadu satu sama lain, dan selisihnya dicetak PER NAMA.

KENAPA ADA, dan kenapa ia BUKAN cocok-tiga-sumber-entitas.py (kini di alat/_arsip/)

    Pendahulunya mengadu TIGA sumber dan mencetak "LULUS". Ia lulus karena
    NILAI_SELISIH dan BESARAN_DAPAT_DISESUAIKAN dikeluarkan lebih dulu ke daftar
    LUAR dengan sebab "modul Adjustment -- embargo, sec 11.3".

    Embargo itu SUDAH DICABUT. Modul Adjustment sudah digrilling (GRILL-A..TDA),
    sudah to-spec (USULAN-DIFF-KE-INDUK.md, enam diff DITERAPKAN), dan sudah
    to-ticket (13 tiket; tiket 06 adalah PEMBUAT PERTAMA NILAI_SELISIH).
    Sebabnya hilang, penagihnya tidak dicabut -- dan pemeriksa yang mengecualikan
    apa yang seharusnya diperiksanya menjawab "LULUS" untuk kedua kemungkinan.

    Maka di sini keduanya DIPERIKSA, bukan dikecualikan, dan daftar LUAR menyusut
    menjadi dua: GEL-2 dan GEL-3.

ENAM SUMBER, dan masing-masing MENGIKAT atas hal yang berbeda

    S1  4-erd-dan-tabel-datar/STRUKTUR-DATA.md        daftar entitas  (MENGIKAT, butir 4)
    S2  2-to-spec/KAMUS-KOLOM.md                      nama+tipe kolom (MENGIKAT, butir 5)
    S3  2-to-spec/ddl-usulan/*.sql                    tabel, PK, FK   (yang dibangun)
    S4  4-erd-dan-tabel-datar/peta-nama-tabel-treatyin.tsv  peta nama tabel datar
    S5  4-erd-dan-tabel-datar/datar-treatyin-lama.csv  985 simpul pohon Pega (AS-IS)
        -- pendamping mesin-baca pohon-treatyin.txt, yang bentuk RINGKASnya
    S6  PETA-TELUSUR-JSON.md sec 6                    nasib 667 jalur JSON
    S7  4-erd-dan-tabel-datar/KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md  G1..G4 Adjustment

APA YANG DAPAT DIBUKTIKAN PERKAKAS INI
    Bahwa keenam sumber menyebut nama yang sama, bahwa tipe kolom kunci asing
    sama dengan tipe kunci utama yang dirujuknya, bahwa tiap padanan di peta nama
    masih ada di DDL, dan bahwa tiap simpul pohon punya nasib di PETA-TELUSUR.

APA YANG TIDAK
    Ia TIDAK dapat menyatakan himpunannya LENGKAP. Semesta sec 10 masih kurang
    340 properti titik buta (L-8). Gambar yang dibangkitkan dari himpunan bolong
    akan TERLIHAT lengkap.
    Ia juga TIDAK membaca isi pasal: sebuah entitas yang namanya cocok di enam
    tempat dengan ARTI berbeda tetap LULUS di sini.
    Dan pohon-treatyin.txt adalah RINGKASAN -- daftar skalarnya dipotong "(+91)".
    Karena itu cacah simpul diambil dari datar-treatyin-lama.csv, bukan darinya.

PENOLAKAN dilaporkan satu baris per sebab, beserta jumlahnya.

Pakai:  PYTHONIOENCODING=utf-8 python alat/cocok-enam-sumber-struktur.py
"""
import io, os, re, csv, json, sys
from collections import Counter, OrderedDict

AKAR = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
ERD = os.path.join(AKAR, '4-erd-dan-tabel-datar')
S1 = os.path.join(ERD, 'STRUKTUR-DATA.md')
S2 = os.path.join(AKAR, '2-to-spec', 'KAMUS-KOLOM.md')
S3 = os.path.join(AKAR, '2-to-spec', 'ddl-usulan')
S4 = os.path.join(ERD, 'peta-nama-tabel-treatyin.tsv')
S5 = os.path.join(ERD, 'datar-treatyin-lama.csv')
S5TXT = os.path.join(ERD, 'pohon-treatyin.txt')
S6 = os.path.join(AKAR, 'PETA-TELUSUR-JSON.md')
S7 = os.path.join(ERD, 'KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md')

tolak = Counter()
catat = []          # (sebab, contoh apa adanya)

# ══ S1. STRUKTUR-DATA.md -- daftar entitas ══════════════════════════════════
t1 = io.open(S1, encoding='utf-8').read()
stru = set(re.findall(r'^>?\s*\|\s*\*{0,2}`([A-Z][A-Z0-9_]+)`\*{0,2}\s*\|\s*\d+\s*\|', t1, re.M))

# ══ S2. KAMUS-KOLOM.md -- nama entitas + kolomnya ═══════════════════════════
t2 = io.open(S2, encoding='utf-8').read()
kamus, kamus_kol, kini = set(), {}, None
for l in t2.split('\n'):
    m = re.match(r'^## `([A-Z][A-Z0-9_]*)`\s*$', l)
    if m:
        kini = m.group(1); kamus.add(kini); kamus_kol[kini] = []
        continue
    if l.startswith('## '):
        kini = None
        tolak['judul "## " di KAMUS yang BUKAN nama entitas'] += 1
        catat.append(('judul "## " bukan entitas', l[3:].strip()[:70]))
        continue
    if kini:
        mc = re.match(r'^\|\s*`([A-Z][A-Z0-9_]*)`\s*\|', l)
        if mc:
            kamus_kol[kini].append(mc.group(1))

# ══ S3. ddl-usulan/*.sql -- tabel, kolom+tipe, PK, FK ═══════════════════════
ddl = OrderedDict()     # nama -> dict(berkas, kolom{nama:tipe}, pk, fk[])
tanpa_create = []
for f in sorted(os.listdir(S3)):
    if not f.lower().endswith('.sql'):
        tolak['berkas di ddl-usulan/ yang bukan .sql'] += 1
        continue
    s = io.open(os.path.join(S3, f), encoding='utf-8').read()
    m = re.search(r'CREATE TABLE\s+\S+\.(\w+)\s*\((.*?)\n\)\s*;', s, re.S)
    if not m:
        tanpa_create.append(f)
        tolak['berkas .sql tanpa CREATE TABLE (berkas kerangka, bukan tabel)'] += 1
        continue
    nm, badan = m.group(1), m.group(2)
    kol = OrderedDict()
    for bl in badan.split('\n'):
        mk = re.match(r'\s*(\w+)\s+(NUMBER\([\d,]+\)|VARCHAR2\(\d+ ?CHAR\)|DATE|CLOB|TIMESTAMP\S*)', bl)
        if mk:
            kol[mk.group(1)] = mk.group(2)
    pk = None
    mp = re.search(r'ADD CONSTRAINT\s+\w+\s+PRIMARY KEY\s*\((\w+)\)', s)
    if mp:
        pk = mp.group(1)
    else:
        tolak['tabel tanpa PRIMARY KEY yang terbaca'] += 1
        catat.append(('tabel tanpa PK terbaca', nm))
    fks = []
    for fm in re.finditer(
            r'ADD CONSTRAINT\s+(\w+)\s+FOREIGN KEY\s*\((\w+)\)\s*\n?\s*REFERENCES\s+\S+\.(\w+)\s*\((\w+)\)\s*([^;]*);', s):
        od = re.search(r'ON DELETE\s+([A-Z ]+)', fm.group(5) or '')
        fks.append(dict(kolom=fm.group(2), induk=fm.group(3), kolom_induk=fm.group(4),
                        on_delete=(od.group(1).strip() if od else '(tidak dinyatakan)')))
    ddl[nm] = dict(berkas=f, kolom=kol, pk=pk, fk=fks)

# ══ S4. peta-nama-tabel-treatyin.tsv ════════════════════════════════════════
peta = []
with io.open(S4, encoding='utf-8') as fh:
    for r in csv.DictReader(fh, delimiter='\t'):
        nt = (r.get('NAMA_T') or '').strip()
        if not nt or nt.startswith('⚠') or 'TIDAK DIBUAT' in nt:
            tolak['baris peta nama bertanda TIDAK DIBUAT (bukan tabel)'] += 1
            continue
        peta.append(r)

# ══ S5. datar-treatyin-lama.csv -- 985 simpul pohon Pega ════════════════════
simpul = []
with io.open(S5, encoding='utf-8-sig') as fh:
    for r in csv.DictReader(fh):
        simpul.append(r)
akar_pohon = set()
for r in simpul:
    j = r['JALUR']
    if j.startswith('TreatyIn.'):
        akar_pohon.add(j[len('TreatyIn.'):].split('.')[0])
# bentuk RINGKAS, dipakai hanya untuk menyatakan batasnya
txt_baris = len(io.open(S5TXT, encoding='utf-8').read().split('\n'))
txt_dipotong = len(re.findall(r'\(\+\d+\)', io.open(S5TXT, encoding='utf-8').read()))

# ══ S6. PETA-TELUSUR-JSON.md sec 6 -- nasib per jalur ═══════════════════════
t6 = io.open(S6, encoding='utf-8').read()
nasib = {}
for gol in ('DIPETAKAN', 'DITURUNKAN', 'DIBUANG', 'DITUNDA'):
    mg = re.search(r'^### %s[^\n]*\n+```\n(.*?)\n```' % gol, t6, re.S | re.M)
    if not mg:
        tolak['golongan nasib yang blok kodenya tidak terbaca di PETA-TELUSUR'] += 1
        catat.append(('golongan nasib tak terbaca', gol))
        continue
    for l in mg.group(1).split('\n'):
        l = l.strip()
        if l:
            nasib[l] = gol


def polos(j):
    """Samakan bentuk jalur: buang awalan TreatyIn. dan tanda []"""
    if j.startswith('TreatyIn.'):
        j = j[len('TreatyIn.'):]
    return j.replace('[]', '')


nasib_polos = {}
for k, v in nasib.items():
    nasib_polos.setdefault(polos(k), v)

# ══ S7. KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md -- entitas dari G1..G4 ════════════
t7 = io.open(S7, encoding='utf-8').read()
adj = set(re.findall(r'^\|\s*`?([A-Z][A-Z0-9_]+)`?\s*\|\s*\d+\s*\|\s*(?:Adjustment|\*\*bersama\*\*)',
                     t7, re.M))

# ══ DAFTAR LUAR -- menyusut dari empat menjadi DUA, dan sebabnya disebut ════
LUAR = {
    'RETRO_KELUAR': 'GEL-2 -- arah keluar, SPEC-MODEL-DATA sec 2.3',
    'PENCAPAIAN':   'GEL-3 -- dimodelkan, tidak dibangun sekarang',
}
DICABUT = {
    'NILAI_SELISIH':
        'embargo sec 11.3 DICABUT -- Adjustment sudah to-spec dan to-ticket; '
        'tiket 06 PEMBUAT PERTAMA',
    'BESARAN_DAPAT_DISESUAIKAN':
        'embargo sec 11.3 DICABUT -- tabel acuan G3, dirujuk NILAI_SELISIH.ID_BESARAN',
}

dnama = set(ddl)
sel_kamus_tanpa_ddl = sorted(kamus - dnama)
sel_ddl_tanpa_kamus = sorted(dnama - kamus)
sel_stru_tanpa_kamus = sorted(x for x in (stru - kamus) if x not in LUAR)
sel_kamus_tanpa_stru = sorted(kamus - stru)
stru_luar = sorted(x for x in (stru - kamus) if x in LUAR)

# ══ PEMERIKSAAN B -- tipe kunci asing versus tipe kunci utama ═══════════════
tipe_beda, fk_yatim, fk_semua = [], [], []
for nm, t in ddl.items():
    for f in t['fk']:
        f2 = dict(f, anak=nm, berkas=t['berkas'])
        fk_semua.append(f2)
        if f['induk'] not in ddl:
            fk_yatim.append(f2)
            continue
        ta = t['kolom'].get(f['kolom'])
        ti = ddl[f['induk']]['kolom'].get(f['kolom_induk'])
        if ta and ti and ta != ti:
            tipe_beda.append(dict(f2, tipe_anak=ta, tipe_induk=ti))

# ══ PEMERIKSAAN C -- padanan di peta nama versus DDL ════════════════════════
# Yang tidak ada di DDL BELUM TENTU basi. Tiga sebab dipisahkan dengan DAFTAR
# BERNAMA -- bukan ditebak dari bentuk namanya.
LUAR_GELOMBANG = {                      # sengaja di luar penyerahan pertama
    'PENCAPAIAN': 'GEL-3, STRUKTUR-DATA sec 1.5',
    'BAGIAN_RETRO': 'GEL-2 RETRO_KELUAR, arah keluar',
    'BAGIAN_FAKULTATIF': 'fakultatif -- GEL-2, SPEC-MODEL-DATA sec 14.3',
    'NILAI_BAGIAN_FAKULTATIF': 'fakultatif -- GEL-2',
    'POTONGAN_FAKULTATIF': 'fakultatif -- GEL-2',
    'REASURADUR_FAKULTATIF': 'fakultatif -- GEL-2',
    'LAYER_FAKULTATIF': 'fakultatif -- GEL-2',
    'DETAIL_PROPORSIONAL_FAKULTATIF': 'fakultatif -- GEL-2',
}
PENUNJANG = {                           # bukan entitas model; penunjang peralihan
    'MIGRASI_KORELASI': 'penunjang migrasi, bukan entitas model',
    'MIGRASI_PENDARATAN': 'penunjang migrasi',
    'MIGRASI_NILAI_DITOLAK': 'penunjang migrasi',
    'ARSIP_MUATAN_KELUAR': 'arsip muatan keluar, bukan entitas model',
}
AKAN_MENDARAT = {                       # akan ada begitu embargo sec 11.3 selesai
    'NILAI_SELISIH': 'entitas Adjustment -- DIDAFTAR di STRUKTUR-DATA, belum ber-DDL',
}
peta_cocok, peta_luar, peta_penunjang, peta_mendarat, peta_belum = [], [], [], [], []
for r in peta:
    for cal in re.split(r'[+·]', r.get('PADANAN_DDL') or ''):
        cal = cal.strip().strip('`')
        if not re.match(r'^([A-Z][A-Z0-9_]+)$', cal):
            continue
        p = (r['NAMA_T'], cal)
        if cal in ddl:
            peta_cocok.append(p)
        elif cal in LUAR_GELOMBANG:
            peta_luar.append(p + (LUAR_GELOMBANG[cal],))
        elif cal in PENUNJANG:
            peta_penunjang.append(p + (PENUNJANG[cal],))
        elif cal in AKAN_MENDARAT:
            peta_mendarat.append(p + (AKAN_MENDARAT[cal],))
        else:
            peta_belum.append(p)
peta_basi = peta_belum

# ══ PEMERIKSAAN D -- akar jalur Pega di peta nama versus pohon ══════════════
peta_akar_asing = []
for r in peta:
    for j in (r.get('JALUR_PEGA') or '').split('|'):
        j = j.strip()
        if not j or j == 'TreatyIn':
            continue
        a = polos(j).split('.')[0]
        if a and a not in akar_pohon:
            peta_akar_asing.append((r['NAMA_T'], j))

# ══ PEMERIKSAAN E -- tiap simpul pohon punya nasib ══════════════════════════
tanpa_nasib = []
for r in simpul:
    j = r['JALUR']
    if j == 'TreatyIn':
        continue
    if polos(j) not in nasib_polos:
        tanpa_nasib.append((j, r['JENIS'], r['BUKTI']))

# ══ PEMERIKSAAN F -- kolom bernama ID_* yang TIDAK berkunci asing ═══════════
tanpa_fk = []
for nm, t in ddl.items():
    berfk = {f['kolom'] for f in t['fk']}
    for k in t['kolom']:
        if k.startswith('ID_') and k != t['pk'] and k not in berfk:
            tanpa_fk.append((nm, k))

# ══════════════════════════ CETAK ══════════════════════════════════════════
P = print
P('=== cocok-enam-sumber-struktur.py -- GERBANG 0 DIPERLEBAR ===')
P()
P('CACAH per sumber:')
P('   %-56s %d' % ('S1 STRUKTUR-DATA.md      entitas (MENGIKAT butir 4)', len(stru)))
P('   %-56s %d' % ('S2 KAMUS-KOLOM.md        entitas (MENGIKAT butir 5)', len(kamus)))
P('   %-56s %d' % ('S3 ddl-usulan/           CREATE TABLE', len(ddl)))
P('   %-56s %d' % ('S4 peta-nama-tabel...tsv baris tabel datar', len(peta)))
P('   %-56s %d' % ('S5 datar-treatyin-lama.csv  simpul pohon Pega', len(simpul)))
P('   %-56s %d' % ('   -- di antaranya ADA_DI_ADJUSTMENT',
                   sum(1 for r in simpul if (r.get('ADA_DI_ADJUSTMENT') or '').strip() == 'ya')))
P('   %-56s %d' % ('S6 PETA-TELUSUR-JSON sec 6  jalur bernasib', len(nasib)))
for g in ('DIPETAKAN', 'DITURUNKAN', 'DIBUANG', 'DITUNDA'):
    P('   %-56s %d' % ('   -- ' + g, sum(1 for v in nasib.values() if v == g)))
P('   %-56s %d' % ('S7 KEPUTUSAN-SAMBUNGAN   entitas Adjustment', len(adj)))
P()
P('DITOLAK / TIDAK DIHITUNG -- satu baris per sebab:')
for k, v in tolak.most_common():
    P('   %-64s %d' % (k, v))
if not tolak:
    P('   (tidak ada)')
for sebab, contoh in catat[:20]:
    P('      [%s] %s' % (sebab, contoh))
P()
P('DAFTAR LUAR -- menyusut dari EMPAT menjadi DUA')
for k, v in sorted(LUAR.items()):
    P('   TETAP DI LUAR  %-28s %s' % (k, v))
for k, v in sorted(DICABUT.items()):
    P('   MASUK LINGKUP  %-28s %s' % (k, v))
P()
P('-- A. HIMPUNAN ENTITAS ------------------------------------------------')


def cetak(judul, daftar, arti):
    P('   %-34s : %d' % (judul, len(daftar)))
    if daftar:
        P('        ARTINYA: %s' % arti)
        for x in daftar:
            P('          - %s' % (x if isinstance(x, str) else ' -> '.join(map(str, x))))


cetak('di KAMUS, tidak ada DDL-nya', sel_kamus_tanpa_ddl, 'tabel TIDAK AKAN TERBANGUN')
cetak('ada DDL, tidak ada di KAMUS', sel_ddl_tanpa_kamus, 'kolomnya TIDAK BERNAMA')
cetak('di KAMUS, tidak di STRUKTUR', sel_kamus_tanpa_stru, 'berkas MENGIKAT kehilangan entitas')
cetak('di STRUKTUR, tidak di KAMUS', sel_stru_tanpa_kamus,
      'DIDAFTAR tetapi tidak berkolom dan tidak ber-DDL -- ia TIDAK AKAN DIBANGUN')
P('   sengaja di luar penyerahan pertama : %d  (%s)' % (len(stru_luar), ', '.join(stru_luar)))
P()
P('-- B. TIPE KUNCI ASING versus TIPE KUNCI UTAMA -------------------------')
P('   kunci asing dibaca dari DDL        : %d' % len(fk_semua))
P('   menggantung (induknya tiada)       : %d' % len(fk_yatim))
for x in fk_yatim:
    P('          - %s.%s -> %s' % (x['anak'], x['kolom'], x['induk']))
P('   TIPE BERBEDA dari yang dirujuk     : %d' % len(tipe_beda))
for x in tipe_beda:
    P('          - %-26s %-22s %-14s  ->  %s.%s %s'
      % (x['anak'], x['kolom'], x['tipe_anak'], x['induk'], x['kolom_induk'], x['tipe_induk']))
P()
P('-- C. PETA NAMA TABEL DATAR versus DDL ---------------------------------')
P('   padanan yang MASIH ADA di DDL      : %d' % len(peta_cocok))
P('   sengaja di luar gelombang ini      : %d' % len(peta_luar))
for nt, cal, seb in peta_luar:
    P('          - %-30s -> %-32s %s' % (nt, cal, seb))
P('   penunjang, bukan entitas model     : %d' % len(peta_penunjang))
for nt, cal, seb in peta_penunjang:
    P('          - %-30s -> %-32s %s' % (nt, cal, seb))
P('   akan mendarat (Adjustment)         : %d' % len(peta_mendarat))
for nt, cal, seb in peta_mendarat:
    P('          - %-30s -> %-32s %s' % (nt, cal, seb))
P('   BELUM TERPADANKAN, sebab tak tahu  : %d   <- ini yang BASI' % len(peta_basi))
for nt, cal in peta_basi:
    P('          - %-30s -> %s' % (nt, cal))
P()
P('-- D. AKAR JALUR PEGA di peta nama versus POHON ------------------------')
P('   akar yang tidak ada di pohon       : %d' % len(peta_akar_asing))
for nt, j in peta_akar_asing[:20]:
    P('          - %-30s %s' % (nt, j))
P()
P('-- E. CAKUPAN NASIB: tiap simpul pohon punya nasib ---------------------')
P('   simpul diperiksa                   : %d' % (len(simpul) - 1))
P('   simpul TANPA nasib di PETA-TELUSUR : %d' % len(tanpa_nasib))
jns = Counter(x[1] for x in tanpa_nasib)
for k, v in jns.most_common():
    P('          menurut JENIS: %-10s %d' % (k, v))
for x in tanpa_nasib[:15]:
    P('          - %s  [%s]' % (x[0], x[1]))
if len(tanpa_nasib) > 15:
    P('          ... dan %d lagi' % (len(tanpa_nasib) - 15))
P()
P('-- F. KOLOM ID_* TANPA KUNCI ASING -------------------------------------')
P('   rujukan yang SEHARUSNYA ADA        : %d' % len(tanpa_fk))
for nm, k in tanpa_fk[:12]:
    P('          - %s.%s' % (nm, k))
if len(tanpa_fk) > 12:
    P('          ... dan %d lagi' % (len(tanpa_fk) - 12))
P()
lulus = not (sel_kamus_tanpa_ddl or sel_ddl_tanpa_kamus or sel_kamus_tanpa_stru
             or sel_stru_tanpa_kamus or fk_yatim or tipe_beda or peta_basi)
P('GERBANG 0 DIPERLEBAR : %s' % ('LULUS' if lulus else 'GAGAL -- lihat A..F di atas'))
P()
P('-- BATAS PERKAKAS INI, dan ia bagian dari keluarannya --')
P('   1. Ia mencocokkan NAMA, bukan ARTI. Entitas yang namanya cocok di enam')
P('      tempat dengan arti berbeda tetap LULUS di sini.')
P('   2. "Cocok" BUKAN "lengkap": semesta sec 10 masih kurang 340 properti')
P('      titik buta (L-8). Gambar dari himpunan bolong TERLIHAT lengkap.')
P('   3. pohon-treatyin.txt adalah RINGKASAN -- %d baris, %d daftar skalar' % (txt_baris, txt_dipotong))
P('      DIPOTONG dengan "(+n)". Cacah simpul karena itu diambil dari')
P('      datar-treatyin-lama.csv (%d simpul), bukan darinya.' % len(simpul))
P('   4. Pemeriksaan E menyamakan jalur dengan membuang "[]" dan awalan')
P('      "TreatyIn."; jalur yang berbeda HANYA pada tanda itu terbaca sama.')

json.dump(dict(
    stru=sorted(stru), kamus=sorted(kamus), ddl=sorted(ddl),
    ddl_rinci={k: dict(berkas=v['berkas'], pk=v['pk'], kolom=v['kolom'],
                       fk=v['fk']) for k, v in ddl.items()},
    kamus_kolom=kamus_kol, adj=sorted(adj),
    luar=LUAR, dicabut=DICABUT,
    selisih=dict(kamus_tanpa_ddl=sel_kamus_tanpa_ddl, ddl_tanpa_kamus=sel_ddl_tanpa_kamus,
                 kamus_tanpa_struktur=sel_kamus_tanpa_stru,
                 struktur_tanpa_kamus=sel_stru_tanpa_kamus),
    fk=fk_semua, fk_yatim=fk_yatim, tipe_beda=tipe_beda,
    peta_basi=peta_basi, peta_akar_asing=peta_akar_asing,
    tanpa_nasib=tanpa_nasib, tanpa_fk=tanpa_fk,
    cacah=dict(simpul=len(simpul), nasib=len(nasib), peta=len(peta)),
    lulus=lulus),
    io.open(os.path.join(AKAR, 'alat', 'himpunan-struktur.json'), 'w', encoding='utf-8'),
    ensure_ascii=False, indent=1)
sys.exit(0 if lulus else 1)
