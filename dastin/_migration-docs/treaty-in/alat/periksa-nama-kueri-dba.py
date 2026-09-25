# -*- coding: utf-8 -*-
"""
periksa-nama-kueri-dba.py — 24 September 2026

KENAPA PERKAKAS INI ADA
    Satu putaran dengan DBA diukur dalam HARI, bukan menit. Satu nama kolom yang
    keliru menghasilkan ORA-00904, kuerinya gagal, berkasnya dikembalikan -- dan
    bila kelirunya pada beberapa kueri, putarannya berulang.

    Seluruh nama tabel dan kolom di berkas permintaan berasal dari ekspor Pega,
    dari inventaris yang kita susun sendiri, atau dari pembacaan kita. TIDAK SATU
    PUN pernah diuji terhadap basis datanya.

KEBENARAN DASAR yang kita punya, dan batasnya:
    Table/*.txt      DDL asli tujuh tabel  -> kolomnya DAPAT diperiksa
    PROCEDURE/*.txt  badan prosedur        -> menyebut kolom, dipakai sebagai penguat
    selebihnya       TIDAK ADA DDL-nya     -> hanya dapat ditandai "belum terkonfirmasi"

    Perkakas ini TIDAK dapat menyatakan sebuah nama benar. Ia hanya dapat
    menyatakan sebuah nama SALAH, atau BELUM TERPERIKSA. Uji A-0 di dalam berkas
    permintaanlah yang menutup sisanya, dan ia dijalankan DBA.

ATURAN PERKAKAS: yang ditolak ikut dilaporkan beserta jumlahnya.
"""
import collections
import io
import os
import re

AKAR = r'D:\XML_NURE\_migration-docs\treaty-in'
SQL = os.path.join(AKAR, 'PERMINTAAN-DBA-1-UJI-A-SAMPAI-H.sql')

RE_TABEL_ALIAS = re.compile(
    r'\b(?:FROM|JOIN)\s+POOLDATA\.([A-Z_0-9]+)\s+([a-zA-Z][a-zA-Z0-9_]*)\b', re.I)
RE_TABEL_POLOS = re.compile(
    r'\b(?:FROM|JOIN)\s+POOLDATA\.([A-Z_0-9]+)\s*(?:,|\n|$)', re.I)
# alias milik JSON_TABLE dan turunannya: kolomnya MAYA, bukan kolom basis data
RE_JSON_ALIAS = re.compile(r'\)\s*\)\s*([a-zA-Z][a-zA-Z0-9_]*)\b')
RE_KOLOM_PATH = re.compile(r'\b([A-Z][A-Z_0-9]{1,})\s+(?:VARCHAR2|NUMBER|DATE|CLOB)')
RE_QUALIFIED = re.compile(r'\b([a-zA-Z][a-zA-Z0-9_]*)\.([A-Z][A-Z_0-9]{1,})\b')

KATA_SQL = {'SELECT', 'FROM', 'WHERE', 'GROUP', 'ORDER', 'BY', 'AND', 'OR', 'NOT',
            'IN', 'IS', 'NULL', 'AS', 'ON', 'CASE', 'WHEN', 'THEN', 'ELSE', 'END',
            'COUNT', 'SUM', 'MIN', 'MAX', 'AVG', 'DISTINCT', 'HAVING', 'UNION',
            'PATH', 'COLUMNS', 'NESTED', 'VARCHAR2', 'NUMBER', 'DATE', 'CLOB'}


def kolom_dari_ddl():
    """Kolom nyata per tabel, dari DDL asli di folder Table/."""
    peta = {}
    d = os.path.join(AKAR, 'Table')
    if not os.path.isdir(d):
        return peta
    for f in os.listdir(d):
        if not f.lower().endswith('.txt'):
            continue
        s = io.open(os.path.join(d, f), encoding='utf-8', errors='replace').read()
        m = re.search(r'CREATE TABLE\s+"([^"]+)"\."([^"]+)"', s)
        if not m:
            continue
        tabel = m.group(2).upper()
        kol = set(re.findall(r'"\s*([A-Z_0-9]+)\s*"\s+(?:VARCHAR2|NUMBER|DATE|CLOB|TIMESTAMP)', s))
        peta[tabel] = kol
    return peta


def kolom_dari_prosedur():
    """Nama yang disebut badan prosedur -- penguat, bukan kebenaran dasar."""
    peta = collections.defaultdict(set)
    d = os.path.join(AKAR, 'PROCEDURE')
    if not os.path.isdir(d):
        return peta
    for f in os.listdir(d):
        s = io.open(os.path.join(d, f), encoding='utf-8', errors='replace').read()
        for t in set(re.findall(r'POOLDATA\.([A-Z_0-9]+)', s)):
            peta[t] |= set(re.findall(r'\b([A-Z][A-Z_0-9]{2,})\b', s))
    return peta


def main():
    s = io.open(SQL, encoding='utf-8').read()
    ddl = kolom_dari_ddl()
    prosedur = kolom_dari_prosedur()

    salah = []          # tabel ada DDL-nya, kolomnya TIDAK ada -> pasti salah
    belum = collections.Counter()   # tabel tanpa DDL -> belum terperiksa
    tabel_dipakai = collections.Counter()
    n_maya = 0
    n_stmt = 0

    for idx, bagian in enumerate(s.split(';')):
        if 'SELECT' not in bagian.upper():
            continue
        n_stmt += 1
        alias_nyata = {}
        for t, a in RE_TABEL_ALIAS.findall(bagian):
            alias_nyata[a] = t.upper()
            tabel_dipakai[t.upper()] += 1
        for t in RE_TABEL_POLOS.findall(bagian):
            tabel_dipakai[t.upper()] += 1
        alias_maya = set(RE_JSON_ALIAS.findall(bagian))
        # kolom hasil JSON_TABLE juga maya
        maya_kol = set(RE_KOLOM_PATH.findall(bagian))
        n_maya += len(maya_kol)

        for a, k in RE_QUALIFIED.findall(bagian):
            if k in KATA_SQL or k in maya_kol:
                continue
            if a in alias_maya or a not in alias_nyata:
                continue
            t = alias_nyata[a]
            if t in ddl:
                if k not in ddl[t]:
                    salah.append((idx, t, a, k))
            else:
                belum[t] += 1

    print('pernyataan SELECT diperiksa :', n_stmt)
    print('tabel dengan DDL asli       :', len(ddl), '->', ', '.join(sorted(ddl)))
    print('kolom MAYA (JSON_TABLE) dilewati :', n_maya)
    print()
    print('=== NAMA KOLOM YANG PASTI SALAH (tabel punya DDL) ===')
    if not salah:
        print('   tidak ada')
    for idx, t, a, k in salah:
        print('   pernyataan %-4d %s.%s  (alias %s)' % (idx, t, k, a))
    print()
    print('=== BELUM TERPERIKSA: tabel tanpa DDL di folder Table/ ===')
    for t, n in belum.most_common():
        print('   %-26s %d rujukan kolom' % (t, n))
    print()
    print('=== seluruh tabel yang dirujuk berkas ini ===')
    for t, n in tabel_dipakai.most_common():
        tanda = 'DDL ada ' if t in ddl else 'DDL TIDAK ADA'
        print('   %-26s %3d kemunculan   %s' % (t, n, tanda))


if __name__ == '__main__':
    main()
