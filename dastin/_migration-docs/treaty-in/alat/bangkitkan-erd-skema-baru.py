# -*- coding: utf-8 -*-
r"""
bangkitkan-erd-skema-baru.py -- ERD skema baru, DIBANGKITKAN dari ddl-usulan/.

    masukan : 2-to-spec/ddl-usulan/*.sql   <- relasi, kardinalitas, PK, CHECK
              2-to-spec/KAMUS-KOLOM.md     <- nama dan tipe kolom (MENGIKAT, butir 5)
    keluaran: 4-erd-dan-tabel-datar/ERD-SKEMA-BARU.xlsx        (tujuh lembar)
              4-erd-dan-tabel-datar/erd-skema-baru.json        (bahan ERD.md)

KENAPA DIBANGKITKAN, BUKAN DIGAMBAR
    Alasannya sama dengan butir 5 urutan wewenang: berkas yang DITULIS TANGAN
    dapat menyimpang dari spesifikasi; yang DIBANGKITKAN tidak bisa. Bila
    keluarannya menyimpang dari ddl-usulan/, ALAT INI yang salah.

SUMBER TIAP BESARAN -- dinyatakan supaya tidak tertukar
    relasi dan ON DELETE  <- FOREIGN KEY di ddl-usulan/, apa adanya
    kardinalitas          <- UNIQUE di Z00_KUNCI_ALAMI.sql. FK yang kolomnya
                             SENDIRIAN unik -> 1:1; selain itu -> 1:N.
                             TIDAK ditebak dari bentuk pohon clipboard lama.
    nama dan tipe kolom   <- CREATE TABLE, yang dibangkitkan dari definisi yang
                             sama dengan KAMUS-KOLOM.md

APA YANG TIDAK DAPAT DIBUKTIKAN PERKAKAS INI
    Ia TIDAK dapat menyatakan himpunan entitasnya LENGKAP. Semesta sec 10 masih
    kurang 340 properti titik buta (L-8, ditagih M-4), dan gambar yang dibangkitkan
    dari himpunan yang bolong akan TERLIHAT LENGKAP. Kalimat itu ikut dicetak ke
    dalam lembar "Catatan & Batas".

PENOLAKAN dilaporkan satu baris per sebab.

Pakai:  PYTHONIOENCODING=utf-8 python alat/bangkitkan-erd-skema-baru.py
"""
import io, os, re, json, sys
from collections import Counter, defaultdict, OrderedDict

from openpyxl import Workbook
from openpyxl.styles import Font, PatternFill, Alignment, Border, Side

AKAR = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
DDL = os.path.join(AKAR, '2-to-spec', 'ddl-usulan')
KELUAR = os.path.join(AKAR, '4-erd-dan-tabel-datar')
XLSX = os.path.join(KELUAR, 'ERD-SKEMA-BARU.xlsx')

TANGGAL = '25 September 2026'
tolak = Counter()

# ══ 1. BACA ddl-usulan/ ═════════════════════════════════════════════════════
tabel = OrderedDict()       # nama -> dict(pk, kolom[], fk[], cek[], berkas)
for f in sorted(os.listdir(DDL)):
    if not f.lower().endswith('.sql'):
        tolak['berkas di ddl-usulan/ yang bukan .sql'] += 1
        continue
    s = io.open(os.path.join(DDL, f), encoding='utf-8').read()
    m = re.search(r'CREATE TABLE\s+\S+\.(\w+)\s*\(\n(.*?)\n\);', s, re.S)
    if not m:
        tolak['berkas .sql tanpa CREATE TABLE (kerangka, bukan tabel)'] += 1
        continue
    nm, badan = m.group(1), m.group(2)
    kolom = []
    for l in badan.split('\n'):
        l = l.strip().rstrip(',')
        if not l:
            continue
        bag = l.split()
        kolom.append(dict(nama=bag[0], tipe=bag[1],
                          wajib=('NOT NULL' in l)))
    pk = re.search(r'ADD CONSTRAINT\s+\w+\s+PRIMARY KEY\s*\((\w+)\)', s)
    fk = []
    for fm in re.finditer(
            r'ADD CONSTRAINT\s+(\w+)\s+FOREIGN KEY\s*\((\w+)\)\s*\n?\s*'
            r'REFERENCES\s+\S+\.(\w+)\s*\((\w+)\)\s*([^;]*);', s):
        ekor = fm.group(5) or ''
        od = re.search(r'ON DELETE\s+([A-Z ]+)', ekor)
        fk.append(dict(kolom=fm.group(2), induk=fm.group(3),
                       kolom_induk=fm.group(4),
                       on_delete=(od.group(1).strip() if od else None),
                       constraint=fm.group(1)))
    cek = [c.strip() for c in re.findall(r'ADD CONSTRAINT\s+\w+\s+\n?\s*CHECK\s*\(', s)]
    punya_cek = 'ADD CONSTRAINT CK_' in s
    tabel[nm] = dict(nama=nm, pk=(pk.group(1) if pk else None), kolom=kolom,
                     fk=fk, cek=punya_cek, berkas=f)

# ══ 2. UNIQUE dari Z00 -> kardinalitas ══════════════════════════════════════
uq = defaultdict(list)
z = os.path.join(DDL, 'Z00_KUNCI_ALAMI.sql')
if os.path.exists(z):
    s = io.open(z, encoding='utf-8').read()
    for m in re.finditer(r'ALTER TABLE\s+\S+\.(\w+)\s+ADD CONSTRAINT\s+\w+\s+UNIQUE\s*\(([^)]*)\)', s):
        uq[m.group(1)].append([c.strip() for c in m.group(2).split(',')])
else:
    tolak['Z00_KUNCI_ALAMI.sql tidak ada -- kardinalitas tidak dapat dibaca'] += 1


def kardinalitas(anak, kolom_fk):
    """1:1 hanya bila kolom FK itu SENDIRIAN unik. Dibaca dari UNIQUE, bukan ditebak."""
    for cols in uq.get(anak, []):
        if cols == [kolom_fk]:
            return '1:1'
    return '1:N'


# ══ 3. PEMBAGIAN LEMBAR -- daftar BERNAMA, bukan aturan yang menebak ════════
# Sumbu pemisahnya BUKAN prop / non-prop: di skema baru itu KOLOM
# (DETAIL_PROPORSIONAL.JENIS_TREATY), bukan tabel. Memisahkan lembar menurut
# sumbu itu menduplikasi hampir seluruh entitas -- persis yang terjadi di v2,
# di mana 34 dari 34 bertanda "bersama".
LEMBAR = OrderedDict([
    ('Kontrak dan versi', ['KONTRAK', 'VERSI_KONTRAK', 'DOKUMEN_ADDENDUM',
                           'JEJAK_PERUBAHAN', 'CATATAN_PERSETUJUAN',
                           'PERISTIWA_KONTRAK', 'DOKUMEN_KONTRAK',
                           # --- empat di bawah TIDAK ada di pembagian yang
                           # diterima; ditempatkan di sini dan DILAPORKAN.
                           'PERIODE_PELAPORAN', 'PERIODE_AKUMULASI',
                           'SKALA_KOASURANSI', 'PORTOFOLIO']),
    ('Layer dan proporsional', ['LAYER', 'DETAIL_PROPORSIONAL', 'PEMULIHAN_LIMIT',
                                'BAGIAN', 'POTONGAN', 'BATAS_PER_BAHAYA']),
    ('Uang per mata uang', ['MATA_UANG_KONTRAK', 'EGNPI', 'RETENSI_CEDANT', 'TERMIN',
                            'NILAI_MDP', 'NILAI_MDP_MINIMUM', 'NILAI_CADANGAN_PREMI',
                            'NILAI_PREMI_BRUTO', 'NILAI_PREMI_BRUTO_MINIMUM',
                            'NILAI_PENYEBARAN']),
    ('Penyebaran', ['PENYEBARAN', 'RINCIAN_PENYEBARAN', 'NILAI_PENYEBARAN']),
    ('Acuan', ['MATA_UANG', 'JENIS_POTONGAN', 'JENIS_REASURANSI', 'BAHAYA',
               'KELOMPOK_TREATY', 'KELAS_BISNIS']),
])
# empat yang pembagian aslinya lewatkan -- dicetak, bukan didiamkan
DITAMBAHKAN = ['PERIODE_PELAPORAN', 'PERIODE_AKUMULASI', 'SKALA_KOASURANSI', 'PORTOFOLIO']

di_lembar = defaultdict(list)
for lb, isi in LEMBAR.items():
    for e in isi:
        di_lembar[e].append(lb)

tak_berlembar = sorted(set(tabel) - set(di_lembar))
for e in tak_berlembar:
    tolak['tabel yang tidak masuk lembar mana pun'] += 1
berlembar_tanpa_tabel = sorted(set(di_lembar) - set(tabel))
for e in berlembar_tanpa_tabel:
    tolak['nama di pembagian lembar yang tidak punya tabel'] += 1

# ══ 4. ENTITAS LUAR -- dirujuk, tidak dimiliki ══════════════════════════════
LUAR_SKEMA = {
    'CEDANT': 'master cedant -- dirujuk lewat ID_CEDANT, TIDAK DIMILIKI skema ini',
    'ASAL_BISNIS': 'master source of business -- dirujuk lewat ID_ASAL_BISNIS, TIDAK DIMILIKI',
    'REASURADUR': 'master reasuradur -- dirujuk lewat ID_REASURADUR_PEMIMPIN, TIDAK DIMILIKI',
    'DAFTAR_RETRO': 'acuan pembagian kapasitas -- PEMILIKNYA BELUM DITETAPKAN (eskalasi butir 7)',
    'DOKUMEN': 'penyimpanan berkas di luar basis data -- DOKUMEN_KONTRAK menyimpan RUJUKAN',
}

# ══ 5. GAYA -- disalin dari contooh.xlsx, konvensi pemilik proses ═══════════
AR13 = Font(name='Arial', size=13, bold=True)
AR9 = Font(name='Arial', size=9, bold=True)
AR8 = Font(name='Arial', size=8)
AR8B = Font(name='Arial', size=8, bold=True)
F_AKAR = PatternFill('solid', fgColor='1F4E79')
F_ANAK = PatternFill('solid', fgColor='2E75B6')
F_PK = PatternFill('solid', fgColor='FFF2CC')
F_SUMBER = PatternFill('solid', fgColor='EDF3F9')
PUTIH = Font(name='Arial', size=9, bold=True, color='FFFFFF')
M = Side(style='medium')
KIRI = Alignment(horizontal='left', vertical='center')
TENGAH = Alignment(horizontal='center', vertical='center')
LEBAR_KOTAK = 34


def kotak(ws, r, c, teks, fill, font, align, atas=False, bawah=False, lebar=LEBAR_KOTAK):
    ws.merge_cells(start_row=r, start_column=c, end_row=r, end_column=c + lebar - 1)
    sel = ws.cell(r, c, teks)
    sel.fill = fill
    sel.font = font
    sel.alignment = align
    for i in range(lebar):
        s2 = ws.cell(r, c + i)
        s2.border = Border(left=(M if i == 0 else None), right=(M if i == lebar - 1 else None),
                           top=(M if atas else None), bottom=(M if bawah else None))
    return r + 1


def panah(ws, r, c):
    sel = ws.cell(r, c, '▼')
    sel.font = AR8B
    sel.alignment = TENGAH
    sel.border = Border(left=M)
    return r + 1


# ══ 6. SUSUN POHON per lembar ══════════════════════════════════════════════
induk_dari = {}
for nm, t in tabel.items():
    p = [x['induk'] for x in t['fk'] if x['induk'] != nm]
    induk_dari[nm] = p

wb = Workbook()
wb.remove(wb.active)
rekap_lembar = []

for lb, isi in LEMBAR.items():
    isi = [e for e in isi if e in tabel]
    ws = wb.create_sheet(lb)
    ws.sheet_view.showGridLines = False
    for i in range(1, 80):
        ws.column_dimensions[ws.cell(1, i).column_letter].width = 2.3

    n_rel_dalam = 0
    r = 1
    r = kotak(ws, r, 2, 'SKEMA TREATY MASUK — %s' % lb.upper(), PatternFill(), AR13, KIRI, lebar=50)
    r = kotak(ws, r, 2,
              '%d tabel · dibangkitkan dari 2-to-spec/ddl-usulan/ · %s        '
              'BUKAN digambar tangan — bila menyimpang dari DDL, alatnya yang salah'
              % (len(isi), TANGGAL), PatternFill(), AR8, KIRI, lebar=70)
    r += 1

    dalam = set(isi)
    anak_dari = defaultdict(list)
    akar = []
    for e in isi:
        p_dalam = [p for p in induk_dari[e] if p in dalam and p != e]
        if p_dalam:
            anak_dari[p_dalam[0]].append(e)
        else:
            akar.append(e)

    digambar = set()

    def gambar(e, tingkat):
        global r, n_rel_dalam
        if e in digambar:
            return
        digambar.add(e)
        c = 2 + 2 * tingkat
        t = tabel[e]
        akar_p = (tingkat == 0)

        # baris nama
        lain = [x for x in di_lembar[e] if x != lb]
        judul = e
        if not akar_p:
            fkp = [x for x in t['fk'] if x['induk'] in dalam]
            if fkp:
                judul += '        ' + kardinalitas(e, fkp[0]['kolom'])
        if lain:
            judul += '        BERSAMA → ' + ' · '.join(lain)
        if akar_p:
            luar_p = [p for p in induk_dari[e] if p not in dalam]
            if luar_p:
                judul += '   (akar di lembar ini · induknya %s → lembar %s)' % (
                    luar_p[0], ', '.join(di_lembar.get(luar_p[0], ['—'])))
        r = kotak(ws, r, c, judul, F_AKAR if akar_p else F_ANAK, PUTIH, TENGAH, atas=True)

        # baris PK
        pk = t['pk'] or '(TIDAK ADA PK — melanggar INV-01)'
        tipe_pk = next((k['tipe'] for k in t['kolom'] if k['nama'] == t['pk']), '')
        r = kotak(ws, r, c, 'PK   %-28s %s' % (pk, tipe_pk), F_PK, AR8B, KIRI)

        # baris FK
        if not t['fk']:
            r = kotak(ws, r, c, '(tidak punya kunci asing)', F_SUMBER, AR8, KIRI)
        for i, x in enumerate(t['fk']):
            od = x['on_delete'] or '(ON DELETE tidak dinyatakan — INV-18)'
            r = kotak(ws, r, c, 'FK   %-26s →  %s.%s   %s   %s'
                      % (x['kolom'], x['induk'], x['kolom_induk'],
                         kardinalitas(e, x['kolom']), od),
                      F_SUMBER, AR8, KIRI)
            if x['induk'] in dalam:
                n_rel_dalam += 1

        # DIRUJUK OLEH -- sisi induk dari relasi yang TIDAK digambar sebagai
        # cabang pohon. Tanpa baris ini, sebuah entitas yang menjadi INDUK lewat
        # kolom di anaknya (DOKUMEN_ADDENDUM <- VERSI_KONTRAK) tampak menggantung
        # sendirian, dan pembacanya menyimpulkan ia tidak berelasi apa-apa.
        anak_pohon = set(anak_dari.get(e, []))
        dirujuk = [(a, x['kolom']) for a in isi for x in tabel[a]['fk']
                   if x['induk'] == e and a != e and a not in anak_pohon]
        for a, kol in sorted(dirujuk):
            r = kotak(ws, r, c, 'DIRUJUK OLEH  %s.%s   %s   — digambar di atas, bukan sebagai cabang'
                      % (a, kol, kardinalitas(a, kol)), F_SUMBER, AR8, KIRI)

        # kunci alami
        for cols in uq.get(e, []):
            r = kotak(ws, r, c, 'UNIQUE  (%s)' % ', '.join(cols), F_SUMBER, AR8, KIRI)
        if t['cek']:
            r = kotak(ws, r, c,
                      'CHECK   tepat satu kolom induk terisi — KTV-B, induk polimorfik',
                      F_SUMBER, AR8, KIRI)
        # tutup kotak
        for i in range(LEBAR_KOTAK):
            ws.cell(r - 1, c + i).border = Border(
                left=(M if i == 0 else None), right=(M if i == LEBAR_KOTAK - 1 else None),
                bottom=M)
        r += 1

        for a in sorted(anak_dari.get(e, [])):
            r = panah(ws, r, c + 2 + 2)
            gambar(a, tingkat + 1)

    for a in akar:
        gambar(a, 0)

    rekap_lembar.append((lb, len(isi), n_rel_dalam))

# ══ 6b. SENSUS: kolom rujukan yang TIDAK punya kunci asing ═════════════
# Pembangkit DDL menurunkan kunci asing dari POLA NAMA `ID_<ENTITAS>`. Setiap
# kolom rujukan yang dinamai lain karena itu TIDAK memperoleh kunci asing, dan
# ketiadaannya TIDAK BERBUNYI di mana pun -- tabelnya tetap berdiri, DDL-nya
# tetap sah, dan hanya pembacaan dua arah yang menemukannya.
LUAR_KOLOM = {
    'ID_CEDANT': 'CEDANT -- master, DI LUAR skema (ERD.md sec 2.8)',
    'ID_ASAL_BISNIS': 'ASAL_BISNIS -- master, DI LUAR skema (sec 2.8)',
    'ID_REASURADUR_PEMIMPIN': 'REASURADUR -- master, DI LUAR skema (sec 2.8)',
    'ID_KETUA_TREATY': 'REASURADUR -- master, DI LUAR skema',
    'ID_DOKUMEN': 'DOKUMEN -- penyimpanan berkas DI LUAR basis data (sec 2.8)',
    'ID_DAFTAR_RETRO': 'acuan kapasitas -- PEMILIKNYA BELUM DITETAPKAN (eskalasi butir 7)',
    'ID_SUSUNAN_RETRO': 'susunan retro -- acuan kapasitas, pemiliknya belum ditetapkan',
}
punya_fk = {(t, x['kolom']) for t in tabel for x in tabel[t]['fk']}
tanpa_fk = []
for t in sorted(tabel):
    for k in tabel[t]['kolom']:
        nmk = k['nama']
        if (t, nmk) in punya_fk or nmk == 'ID_' + t:
            continue
        if nmk in LUAR_KOLOM:
            tanpa_fk.append((t, nmk, 'DI LUAR SKEMA', LUAR_KOLOM[nmk]))
        elif nmk.startswith('KODE_MATA_UANG') or nmk.startswith('MATA_UANG_'):
            tanpa_fk.append((t, nmk, 'SEHARUSNYA ADA', 'MATA_UANG.KODE -- INV-44 menuntutnya'))
        elif 'KELAS_BISNIS' in nmk:
            tanpa_fk.append((t, nmk, 'SEHARUSNYA ADA', 'KELAS_BISNIS.ID_KELAS_BISNIS'))
        elif nmk == 'ID_KONTRAK_DISALIN_DARI':
            tanpa_fk.append((t, nmk, 'SEHARUSNYA ADA', 'KONTRAK -- rujukan ke tabelnya sendiri, ERD.md sec 2.1'))
        elif nmk == 'ID_INDUK' and t == 'JENIS_REASURANSI':
            tanpa_fk.append((t, nmk, 'SEHARUSNYA ADA', 'JENIS_REASURANSI -- jenis bersusun, rujukan ke tabelnya sendiri'))
        elif nmk == 'ID_JENIS_REASURANSI_INDUK':
            tanpa_fk.append((t, nmk, 'SEHARUSNYA ADA', 'JENIS_REASURANSI.ID_JENIS_REASURANSI'))
        elif nmk.startswith('ID_') and nmk[3:] not in tabel:
            tanpa_fk.append((t, nmk, 'TIDAK TERBACA', 'namanya tidak cocok tabel mana pun, dan bukan daftar luar'))
n_harus = sum(1 for x in tanpa_fk if x[2] == 'SEHARUSNYA ADA')

# ══ 7. LEMBAR "Daftar Relasi" ══════════════════════════════════════════════
ws = wb.create_sheet('Daftar Relasi')
ws.sheet_view.showGridLines = False
for k, w in zip('ABCDEFGHIJ', [4, 26, 30, 26, 8, 30, 26, 26, 34]):
    ws.column_dimensions[k].width = w
ws.cell(1, 1, 'DAFTAR RELASI — seluruh kunci asing di 2-to-spec/ddl-usulan/ · %s' % TANGGAL).font = AR13
ws.cell(2, 1, 'Dibangkitkan. Sumber relasi: FOREIGN KEY. Sumber kardinalitas: UNIQUE di '
              'Z00_KUNCI_ALAMI.sql. ON DELETE diambil APA ADANYA, tidak ditebak.').font = AR8
kepala = ['#', 'Induk', 'Anak', 'Kolom kunci tamu', 'Kard.', 'ON DELETE',
          'Lembar induk', 'Lembar anak', 'Catatan']
for i, h in enumerate(kepala, 1):
    s = ws.cell(4, i, h)
    s.font = AR8B
    s.fill = F_PK
    s.border = Border(top=M, bottom=M)
baris_rel = []
rr = 5
n = 0
for nm in sorted(tabel):
    for x in tabel[nm]['fk']:
        n += 1
        lin = ', '.join(di_lembar.get(x['induk'], ['— TIDAK DIGAMBAR']))
        lan = ', '.join(di_lembar.get(nm, ['— TIDAK DIGAMBAR']))
        cat = []
        if x['induk'] == nm:
            cat.append('rujukan ke tabelnya sendiri')
        if tabel[nm]['cek'] and x['induk'] in ('BAGIAN', 'DETAIL_PROPORSIONAL'):
            cat.append('induk polimorfik — KTV-B, tepat satu terisi')
        if lin != lan:
            cat.append('relasi LINTAS LEMBAR')
        od = x['on_delete'] or '(tidak dinyatakan)'
        nilai = [n, x['induk'], nm, x['kolom'], kardinalitas(nm, x['kolom']),
                 od, lin, lan, ' · '.join(cat)]
        for i, v in enumerate(nilai, 1):
            s = ws.cell(rr, i, v)
            s.font = AR8
            s.alignment = KIRI
            if od == '(tidak dinyatakan)' and i == 6:
                s.fill = PatternFill('solid', fgColor='FCE4D6')
        baris_rel.append(dict(no=n, induk=x['induk'], anak=nm, kolom=x['kolom'],
                              kard=kardinalitas(nm, x['kolom']), on_delete=od,
                              lembar_induk=lin, lembar_anak=lan))
        rr += 1
ws.cell(rr + 1, 1, 'JUMLAH').font = AR8B
ws.cell(rr + 1, 2, '%d kunci asing' % n).font = AR8B
ws.cell(rr + 2, 1, 'ON DELETE tidak dinyatakan').font = AR8B
ws.cell(rr + 2, 2, '%d dari %d — INV-18 menuntutnya DITETAPKAN SADAR' % (
    sum(1 for b in baris_rel if b['on_delete'] == '(tidak dinyatakan)'), n)).font = AR8B

# -- sensus kolom rujukan TANPA kunci asing ---------------------------------
rr += 5
ws.cell(rr, 1, 'KOLOM RUJUKAN YANG TIDAK PUNYA KUNCI ASING — %d, di antaranya %d SEHARUSNYA ADA'
        % (len(tanpa_fk), n_harus)).font = AR13
rr += 1
ws.cell(rr, 1, 'Pembangkit DDL menurunkan kunci asing dari POLA NAMA ID_<ENTITAS>. Kolom rujukan '
        'yang dinamai lain TIDAK memperolehnya, dan ketiadaannya TIDAK BERBUNYI — tabelnya '
        'tetap berdiri dan DDL-nya tetap sah.').font = AR8
rr += 2
for i, h in enumerate(['#', 'Tabel', 'Kolom', 'Kedudukan', 'Sasaran yang dimaksud'], 1):
    c_ = ws.cell(rr, i, h); c_.font = AR8B; c_.fill = F_PK
    c_.border = Border(top=M, bottom=M)
rr += 1
for i, (t_, k_, ked, sas) in enumerate(tanpa_fk, 1):
    for j, v in enumerate([i, t_, k_, ked, sas], 1):
        c_ = ws.cell(rr, j, v); c_.font = AR8; c_.alignment = KIRI
        if ked == 'SEHARUSNYA ADA' and j == 4:
            c_.fill = PatternFill('solid', fgColor='FCE4D6')
    rr += 1

# ══ 8. LEMBAR "Catatan & Batas" ════════════════════════════════════════════
ws = wb.create_sheet('Catatan & Batas')
ws.sheet_view.showGridLines = False
for k, w in zip('ABC', [46, 18, 96]):
    ws.column_dimensions[k].width = w
ws.cell(1, 1, 'CATATAN & BATAS — sumber tiap angka, dan apa yang berkas ini TIDAK dapat buktikan').font = AR13
r = 3


def bar(a, b='', c='', bold=False, fill=None):
    global r
    for i, v in enumerate((a, b, c), 1):
        s = ws.cell(r, i, v)
        s.font = AR8B if bold else AR8
        s.alignment = Alignment(horizontal='left', vertical='top', wrap_text=True)
        if fill:
            s.fill = fill
    r += 1


bar('SUMBER TIAP BESARAN', 'Diambil dari', 'Catatan', bold=True, fill=F_PK)
bar('daftar entitas', 'ddl-usulan/', 'CREATE TABLE. Dicocokkan dua arah dengan KAMUS-KOLOM.md dan STRUKTUR-DATA.md sebelum lembar ini dibangkitkan')
bar('relasi', 'ddl-usulan/', 'FOREIGN KEY apa adanya')
bar('kardinalitas', 'Z00_KUNCI_ALAMI.sql', 'UNIQUE. FK yang kolomnya SENDIRIAN unik -> 1:1; selain itu 1:N. TIDAK ditebak dari bentuk pohon clipboard lama')
bar('ON DELETE', 'ddl-usulan/', 'apa adanya. Yang tidak dinyatakan ditulis "(tidak dinyatakan)", BUKAN ditebak')
bar('nama dan tipe kolom', 'KAMUS-KOLOM.md', 'mengikat — urutan wewenang butir 5')
r += 1
bar('YANG BERKAS INI TIDAK DAPAT BUKTIKAN', '', '', bold=True, fill=F_PK)
bar('bahwa himpunan entitasnya LENGKAP', 'L-8 · M-4',
    'Semesta sec 10 masih kurang 340 properti titik buta, dan BELUM DIPERIKSA SIAPA PUN. Gambar yang dibangkitkan dari himpunan yang bolong akan TERLIHAT LENGKAP — itu bentuk kerusakan yang paling berbahaya di proyek ini')
bar('bahwa constraint-nya BERJALAN', 'L-3',
    'Tidak ada instans Oracle yang terjangkau. Setiap baris di ddl-usulan/ adalah BACAAN, bukan hasil uji')
bar('bahwa presisinya benar', 'KTV-A',
    'Presisi diputuskan TANPA VERIFIKASI, dan BERTENGGAT: boleh dipersempit hanya SEBELUM data dimuat (tiket 44)')
bar('bahwa ON DELETE sudah diputuskan', 'INV-18',
    'NOL dari %d kunci asing menyatakannya di DDL. Keputusannya ADA — ERD.md sec 2 memutuskannya untuk KETIGA PULUH ENAM relasi dalam-skema — dan HILANG di jalan menuju berkas yang akan dibangun' % n)
bar('bahwa seluruh rujukan terjaga', 'INV-44 dkk',
    '%d kolom rujukan TIDAK punya kunci asing; %d di antaranya SEHARUSNYA ADA. Delapan belas menunjuk MATA_UANG.KODE, dan INV-44 menuntutnya. Daftarnya di lembar "Daftar Relasi"' % (len(tanpa_fk), n_harus))
bar('bahwa ID_VERSI_KONTRAK_DASAR ada', 'F-19',
    'Kolom itu TIDAK ADA di sec 10, TIDAK ADA di KAMUS-KOLOM.md, dan TIDAK ADA di ddl-usulan/ — padahal ERD.md sec 2.2 menyebutnya inti sambungan Adjustment dan tiket 01 menamainya sebagai artefaknya')
r += 1
bar('ENTITAS LUAR — DIRUJUK, TIDAK DIMILIKI', '', '', bold=True, fill=F_PK)
for k, v in LUAR_SKEMA.items():
    bar(k, '', v)
r += 1
bar('SENGAJA DI LUAR PENYERAHAN PERTAMA', '', '', bold=True, fill=F_PK)
bar('RETRO_KELUAR', 'GEL-2', 'arah keluar — SPEC-MODEL-DATA sec 2.3')
bar('PENCAPAIAN', 'GEL-3', 'dimodelkan sekarang, tidak dibangun sekarang')
bar('NILAI_SELISIH', 'Adjustment', 'sec 11.3 — BELUM PUNYA DDL. Tidak digambar dari kamus; dilaporkan sebagai lubang')
bar('BESARAN_DAPAT_DISESUAIKAN', 'Adjustment', 'sec 11.3 — idem')
bar('KUNCI_PADANAN', 'Adjustment', 'disebut PETA-TELUSUR-JSON sec 176 sebagai kunci NILAI_SELISIH. BUKAN entitas di sec 10, dan TIDAK punya DDL')
r += 1
bar('PEMBAGIAN LEMBAR', '', '', bold=True, fill=F_PK)
bar('sumbu pemisahnya BUKAN prop / non-prop', '',
    'Di sistem lama prop dan non-prop memisahkan LAYAR. Di skema baru ia KOLOM (DETAIL_PROPORSIONAL.JENIS_TREATY), bukan tabel. Memisahkan lembar menurut sumbu itu menduplikasi hampir seluruh entitas — persis yang terjadi di v2, di mana 34 dari 34 bertanda "bersama"')
bar('empat entitas DITAMBAHKAN ke lembar 1', '',
    'PERIODE_PELAPORAN, PERIODE_AKUMULASI, SKALA_KOASURANSI, PORTOFOLIO tidak ada di pembagian lembar yang diterima. Keempatnya anak langsung VERSI_KONTRAK. DILAPORKAN, bukan didiamkan')
for lb, nn, nrel in rekap_lembar:
    bar(lb, '%d tabel' % nn, '%d relasi digambar di dalam lembar ini' % nrel)

os.makedirs(KELUAR, exist_ok=True)
try:
    wb.save(XLSX)
    simpan = 'DITULIS'
except Exception as e:
    simpan = 'GAGAL DITULIS -- %s' % e

# ══ 9. laporan ═════════════════════════════════════════════════════════════
print('=== bangkitkan-erd-skema-baru.py ===')
print('tabel dibaca dari ddl-usulan/ : %d' % len(tabel))
print('kunci asing                   : %d' % n)
print('lembar                        : %d  (%d isi + Daftar Relasi + Catatan & Batas)'
      % (len(rekap_lembar) + 2, len(rekap_lembar)))
print()
print('DITOLAK / DILAPORKAN -- satu baris per sebab:')
for k, v in tolak.most_common():
    print('   %-58s %d' % (k, v))
if tak_berlembar:
    print('      tabel tanpa lembar, apa adanya: %s' % ', '.join(tak_berlembar))
if berlembar_tanpa_tabel:
    print('      nama di pembagian tanpa tabel : %s' % ', '.join(berlembar_tanpa_tabel))
if not tolak:
    print('   (tidak ada)')
print()
print('ON DELETE tidak dinyatakan    : %d dari %d  -- INV-18'
      % (sum(1 for b in baris_rel if b['on_delete'] == '(tidak dinyatakan)'), n))
print('kolom rujukan TANPA kunci asing: %d  (%d SEHARUSNYA ADA)' % (len(tanpa_fk), n_harus))
print('relasi LINTAS LEMBAR          : %d'
      % sum(1 for b in baris_rel if b['lembar_induk'] != b['lembar_anak']))
print('entitas di lebih dari satu lembar: %s'
      % ', '.join('%s (%s)' % (k, ', '.join(v)) for k, v in di_lembar.items() if len(v) > 1))
print()
for lb, nn, nrel in rekap_lembar:
    print('   %-26s %2d tabel  %2d relasi dalam lembar' % (lb, nn, nrel))
print()
print('%s : %s' % (simpan, XLSX))
print()
print('-- BATAS PERKAKAS INI --')
print('   Ia TIDAK dapat menyatakan himpunan entitasnya LENGKAP (L-8: 340 properti')
print('   titik buta belum diperiksa siapa pun), dan TIDAK dapat menyatakan satu pun')
print('   constraint BERJALAN (L-3: tidak ada instans yang terjangkau).')

json.dump({'tanpa_fk': tanpa_fk, 'tabel': {k: dict(v, kolom=v['kolom']) for k, v in tabel.items()},
           'relasi': baris_rel, 'lembar': {k: v for k, v in LEMBAR.items()},
           'di_lembar': {k: v for k, v in di_lembar.items()},
           'ditambahkan_ke_lembar1': DITAMBAHKAN,
           'luar_skema': LUAR_SKEMA, 'uq': {k: v for k, v in uq.items()},
           'tolak': dict(tolak), 'tanggal': TANGGAL},
          io.open(os.path.join(KELUAR, 'erd-skema-baru.json'), 'w', encoding='utf-8'),
          ensure_ascii=False, indent=1)
