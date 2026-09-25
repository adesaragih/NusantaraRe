# -*- coding: utf-8 -*-
"""Tabel datar dan ERD modul Komite Claim Non Prop - bentuk SISTEM LAMA.

HANYA MEMBACA. Nol berkas sumber berubah.

Ini pendataran pohon halaman kerja Komite menjadi tabel. Ia BUKAN rancangan
skema baru; rancangan itu hidup di SPEC-KOMITE-01.md bagian Model data
(tujuh objek KLAIMNP). Yang digambar di sini adalah bentuk yang ADA.

Penamaan mengikuti Diagram-Skema-Tabel-NusantaraRe.xlsx lewat
claim-non-prop/alat/buat-skema-claimnp.py: akar T_WORK_CLAIM, lalu
T_GENERAL_KOMITE 1:1 SHARED PK, lalu tabel lini berawalan T_KOMITE_.

Sumber:
  4-erd-dan-tabel-datar/struktur-komite-lama.md   pohon halaman kerja
  4-erd-dan-tabel-datar/datar-komite-lama.csv     cacah ref per simpul
  claim-non-prop/.../buat-skema-claimnp.py        konvensi nama dan bentuk

Keluaran:
  datar-komite-tabel.csv        tabel datar - satu baris per kolom
  RELASI-KOMITE.csv             satu baris per relasi
  ERD-KOMITE-LAMA.html
  Diagram-Skema-Tabel-Komite.xlsx

Jalankan:  python alat/buat-skema-komite.py
"""
from __future__ import print_function
import io, os, sys, csv, collections

ALAT = os.path.dirname(os.path.abspath(__file__))
AKAR = os.path.normpath(os.path.join(ALAT, '..'))
ERD = os.path.join(AKAR, '4-erd-dan-tabel-datar')
TANGGAL = u'22 September 2026'

try:
    import openpyxl
    from openpyxl.styles import Font, PatternFill, Alignment, Border, Side
except ImportError:
    print('GAGAL: openpyxl tidak ada. pip install openpyxl')
    sys.exit(1)

# ---------------------------------------------------------------- sumber
POHON = {}
with io.open(os.path.join(ERD, 'datar-komite-lama.csv'), encoding='utf-8-sig', newline='') as f:
    for r in csv.DictReader(f):
        POHON[r['JALUR']] = r

# Sebagian jalur bukan properti bisnis dan karena itu tidak masuk pohon -
# pxSubscript salah satunya. Cacahnya dihitung LANGSUNG dari XML supaya tetap
# turunan, bukan diketik.
SUMBER_XML = os.path.normpath(os.path.join(AKAR, '..', '..', 'Komite Claim Non Prop'))
_cache_grep = {}

def grep_xml(nama):
    if nama in _cache_grep:
        return _cache_grep[nama]
    import re as _re
    pola = _re.compile(r'\b%s\b' % _re.escape(nama))
    tot, brk = 0, 0
    for r, _, fs in os.walk(SUMBER_XML):
        for f in fs:
            if not f.lower().endswith('.xml'):
                continue
            s = io.open(os.path.join(r, f), encoding='utf-8', errors='replace').read()
            c = len(pola.findall(s))
            if c:
                tot += c
                brk += 1
    _cache_grep[nama] = (tot, brk)
    return _cache_grep[nama]

def ref(jalur):
    r = POHON.get(jalur)
    if r:
        return (int(r['REF']), int(r['CACAH_BERKAS']))
    if jalur and not jalur.startswith(u'('):
        return grep_xml(jalur.split('.')[-1])
    return (0, 0)

# ---------------------------------------------------------------- tabel
T = []
def tabel(nama, kard, induk, kunci, asal, kelompok, kolom, catatan=(),
          pk_jalur=u'(surogat pendataran)', pk_yakin=u'konvensi',
          fk_jalur=None, fk_yakin=u'konvensi'):
    """pk_jalur / fk_jalur: dari mana kunci itu berasal DI SISTEM LAMA.

    Pega tidak punya kunci baris pada page list - identitas sebuah baris
    adalah POSISInya (pxSubscript, 695 rujukan di 59 berkas). PK dan FK di
    bawah karena itu sebagian surogat yang lahir dari pendataran, dan itu
    dinyatakan per baris alih-alih disamarkan.
    """
    if fk_jalur is None:
        # Bedanya nyata: page list punya baris berposisi, page bersarang tidak
        # punya baris sama sekali - ia SATU halaman yang menempel pada induknya.
        fk_jalur = (u'(posisi baris di daftar induk - pxSubscript)' if kard == u'1:N'
                    else u'(halaman bersarang; tidak ada kolom penghubung)')
    T.append({'nama': nama, 'kard': kard, 'induk': induk, 'kunci': kunci,
              'asal': asal, 'kelompok': kelompok, 'kolom': list(kolom),
              'catatan': list(catatan), 'pk_jalur': pk_jalur, 'pk_yakin': pk_yakin,
              'fk_jalur': fk_jalur, 'fk_yakin': fk_yakin})

# kolom: (NAMA_KOLOM, JALUR_PEGA, KEYAKINAN)
K = lambda n, j, y=u'jalur-penuh': (n, j, y)
# Baris page list dikenali POSISInya. Kolom ini yang menjaga urutannya
# saat didatarkan; tanpanya urutan giliran hilang.
SUB = K(u'SUBSCRIPT', u'pxSubscript')

tabel('T_WORK_CLAIM', u'akar', None, None,
      u'pyWorkPage - satu baris per work object', 'akar',
      [K(u'ID', u'pyWorkPage.pyID'),
       K(u'COVER_KEY', u'pyWorkPage.pxCoverInsKey'),
       K(u'LINI', u'- (ditetapkan workbook NusantaraRe)', u'konvensi')],
      [u'Baris komite ber-ID TKMT-xxxxxx; baris klaim CLMNP-xxxxxx.',
       u'COVER_KEY adalah SATU-SATUNYA tali berkas sirkulasi ke klaim induknya.'])

tabel('T_GENERAL_KOMITE', u'1:1', 'T_WORK_CLAIM', u'-  SHARED PK (ID = ID)',
      u'pyWorkPage{} - 18 properti skalar', 'komite',
      [K(u'ACCEPT_STATUS', u'pyWorkPage.AcceptStatus'),
       K(u'COMMENT', u'pyWorkPage.Comment'),
       K(u'IS_SUBJECTIVITY', u'pyWorkPage.IsSubjectivity'),
       K(u'SUBJECTIVITY_NOTE', u'pyWorkPage.SubjectivityNote'),
       K(u'KOMITE_COUNT', u'pyWorkPage.KomiteCount'),
       K(u'KOMITE_LOOP', u'pyWorkPage.KomiteLoop'),
       K(u'IS_CLOSE_FILE', u'pyWorkPage.IsCloseFile'),
       K(u'IS_REJECT', u'pyWorkPage.IsReject'),
       K(u'IS_PREVIOUS', u'pyWorkPage.IsPrevious'),
       K(u'CLM_NO', u'pyWorkPage.CLMNO'),
       K(u'CREATED_AT', u'pyWorkPage.pxCreateDateTime'),
       K(u'CREATED_BY', u'pyWorkPage.pxCreateOperator'),
       K(u'CREATED_BY_NAME', u'pyWorkPage.pxCreateOpName'),
       K(u'LOCK_HANDLE', u'pyWorkPage.pxLockHandle'),
       K(u'WORK_ID_PREFIX', u'pyWorkPage.pyWorkIDPrefix'),
       K(u'INS_KEY', u'pyWorkPage.pzInsKey')],
      [u'FK   ADJUSTMENT_ID  ->  T_CLAIM_ADJUSTMENT.ID   NOT NULL - index UNIK',
       u'   Di sistem lama tali ini POSISIONAL: Adjustment.IndexObject adalah',
       u'   indeks balik ke baris AdjustmentList, dan KomitePostAdjustment',
       u'   memakainya lewat Local.IdxParent / Local.IdxAdjustment (24 + 19 ref).',
       u'PK = pzInsKey, satu-satunya pengenal baris yang benar-benar ada (683 ref).',
       u'Tali ke KLAIM adalah pxCoverInsKey di T_WORK_CLAIM, bukan kolom di sini.',
       u'KOMITE_COUNT dan KOMITE_LOOP adalah DUA penyimpan untuk satu fakta.',
       u'   Keduanya digambar apa adanya; K5-1 menghapus keduanya di skema baru.'],
      pk_jalur=u'pyWorkPage.pzInsKey', pk_yakin=u'jalur-penuh',
      fk_jalur=u'pyWorkPage.pxCoverInsKey', fk_yakin=u'jalur-penuh')

tabel('T_KOMITE_KOMITELIST', u'1:N', 'T_GENERAL_KOMITE', 'DATA_KOMITE_ID',
      u'KomiteList[] - satu baris = satu jenjang', 'komite',
      [K(u'KOMITE_ID', u'pyWorkPage.KomiteList.KomiteID'),
       K(u'KOMITE_APROVAL', u'pyWorkPage.KomiteList.KomiteAproval'),
       K(u'KOMITE_COMMENT', u'pyWorkPage.KomiteList.KomiteComment'),
       K(u'DATE_APPROVAL', u'pyWorkPage.KomiteList.DateApproval'),
       K(u'DATE_APPROVE', u'pyWorkPage.KomiteList.DateApprove'),
       K(u'KOMITE_EMAIL', u'KomiteList(n).KomiteEmail', u'rujukan-relatif'),
       K(u'ID_KOMITE', u'KomiteList(n).IDKomite', u'rujukan-relatif'),
       K(u'KOMITE_POST', u'KomiteList(n).KomitePost', u'rujukan-relatif'),
       K(u'INITIAL', u'KomiteList(n).Initial', u'rujukan-relatif')],
      [u'Sembilan kolom - cocok dengan tebakan workbook NusantaraRe. Lihat bagian 3.',
       u'DATE_APPROVAL dan DATE_APPROVE adalah DUA kolom untuk satu fakta.',
       u'NOL kolom urutan. DEGREE hidup di EMAILKOMITE, di luar folder ini.'])

tabel('T_KOMITE_LOSS_CONTEXT', u'1:1', 'T_GENERAL_KOMITE', 'DATA_KOMITE_ID',
      u'Komite{} -> GCNMFW-Data-Comitee', 'komite',
      [K(u'CIRCUM_CAUSE_OF_LOSS', u'pyWorkPage.Komite.CircumtansesCouseOfLoss'),
       K(u'EXTENT_OF_LOSS', u'pyWorkPage.Komite.ExtentOfLoss'),
       K(u'LEGAL_LIABILITY', u'pyWorkPage.Komite.LegalLiability'),
       K(u'OCCUPATION', u'pyWorkPage.Komite.Occupation'),
       K(u'REMARKS', u'pyWorkPage.Komite.Remarks')],
      [u'AWAS. Kelas GCNMFW-Data-Comitee dideklarasikan HANYA di sini, dan isinya',
       u'   konteks kerugian - BUKAN daftar jenjang. Kelas yang sama menaungi',
       u'   KomiteList/ComiteeClaim/ClaimComitee tanpa dideklarasikan.',
       u'   Satu kelas, dua arti yang tidak berhubungan.'])

tabel('T_KOMITE_ADJ_SNAPSHOT', u'1:1', 'T_GENERAL_KOMITE', 'DATA_KOMITE_ID',
      u'Adjustment{} -> GCNMFW-Data-Adjustment - SALINAN usulan', 'salinan',
      [K(u'ACCEPTED_NO', u'pyWorkPage.Adjustment.AcceptedNo'),
       K(u'PAYMENT_TYPE', u'pyWorkPage.Adjustment.PaymentType'),
       K(u'INDEX_OBJECT', u'pyWorkPage.Adjustment.IndexObject'),
       K(u'CNP_INDEX_INTERIM', u'pyWorkPage.Adjustment.CNPIndexInterim'),
       K(u'XOL_ID', u'pyWorkPage.Adjustment.XOLID'),
       K(u'CNP_ACC_NO_ADJUST_F', u'pyWorkPage.Adjustment.CNPAccNoAdjustF'),
       K(u'CNP_ACC_NO_OTHER_F', u'pyWorkPage.Adjustment.CNPAccNoOtherF'),
       K(u'CNP_ACC_NO_SALVAGE', u'pyWorkPage.Adjustment.CNPAccNoSalvage'),
       K(u'CNP_ACC_NO_REINSTATE', u'pyWorkPage.Adjustment.CNPAccNoReinstate'),
       K(u'IS_PROPOSE_CLOSE', u'pyWorkPage.Adjustment.IsProposeClose')],
      [u'SALINAN. Rumahnya T_CLAIM_ADJUSTMENT di sisi Klaim; ini potret beku.',
       u'CNPLayerList, SpreadingRisk, AlokasiXOLPaid disentuh sebagai LIST UTUH -',
       u'   isinya tidak pernah dirujuk dari modul ini, jadi tanpa tabel anak.',
       u'Empat kolom CNP_ACC_NO_* adalah penanda per komponen: empat penanda, satu nomor.'])

tabel('T_KOMITE_CLAIM_SNAPSHOT', u'1:1', 'T_GENERAL_KOMITE', 'DATA_KOMITE_ID',
      u'ClaimData{} -> GCNMFW-Data-ClaimData - SALINAN klaim', 'salinan',
      [K(u'NO_CLAIM', u'pyWorkPage.ClaimData.NoClaim'),
       K(u'INSURED_NAME', u'pyWorkPage.ClaimData.InsuredName'),
       K(u'INSURED_RELATIONSHIP', u'pyWorkPage.ClaimData.InsuredRelationship'),
       K(u'PAYABLE', u'pyWorkPage.ClaimData.Payable'),
       K(u'TOTAL_LIST_CLAIM_AMOUNT', u'pyWorkPage.ClaimData.TotalListClaimAmount'),
       K(u'TOTAL_LIST_CLAIM_AMOUNT_IDR', u'pyWorkPage.ClaimData.TotalListClaimAmountIDR')],
      [u'SALINAN KETIGA. Klaim yang sama juga dibaca lewat pyWorkCover.ClaimData',
       u'   (64 simpul) dan ditulis lewat TempMainWork.ClaimData (39 simpul).',
       u'   Satu klaim, tiga arah, tiga saat. Lihat struktur-komite-lama.md bagian 5.',
       u'InterestList disentuh sebagai LIST UTUH - isinya tidak dirujuk.'])

tabel('T_KOMITE_CLAIM_POLICY', u'1:1', 'T_KOMITE_CLAIM_SNAPSHOT', 'CLAIM_SNAPSHOT_ID',
      u'ClaimData.PolicyData{}', 'tipis',
      [K(u'END_DATE_TIME', u'pyWorkPage.ClaimData.PolicyData.EndDateTime')])

tabel('T_KOMITE_CLAIM_ADJ', u'1:N', 'T_KOMITE_CLAIM_SNAPSHOT', 'CLAIM_SNAPSHOT_ID',
      u'ClaimData.AdjustmentList[]', 'tipis', [],
      [u'NOL kolom sendiri terbaca. Ia dilewati untuk mencapai CurencyAdjustment.'])

tabel('T_KOMITE_ADJ_CURRENCY', u'1:N', 'T_KOMITE_CLAIM_ADJ', 'CLAIM_ADJ_ID',
      u'ClaimData.AdjustmentList[].CurencyAdjustment[]  <- KEDALAMAN 4', 'tipis',
      [K(u'CURRENCY_ID', u'pyWorkPage.ClaimData.AdjustmentList.CurencyAdjustment.CurrencyID')],
      [u'Salah eja CurencyAdjustment ada di sumbernya.'])

tabel('T_KOMITE_CLAIM_OBJECT', u'1:N', 'T_KOMITE_CLAIM_SNAPSHOT', 'CLAIM_SNAPSHOT_ID',
      u'ClaimData.ObjectList[]', 'tipis', [],
      [u'NOL kolom sendiri terbaca.'])

tabel('T_KOMITE_OBJECT_ITEM', u'1:N', 'T_KOMITE_CLAIM_OBJECT', 'CLAIM_OBJECT_ID',
      u'ClaimData.ObjectList[].ObjectItemList[]', 'tipis', [],
      [u'NOL kolom sendiri terbaca.'])

tabel('T_KOMITE_OBJECT_ITEM_ADJ', u'1:N', 'T_KOMITE_OBJECT_ITEM', 'OBJECT_ITEM_ID',
      u'...ObjectItemList[].Adjustment[]  <- KEDALAMAN 5', 'tipis',
      [K(u'CURRENCY_ID', u'pyWorkPage.ClaimData.ObjectList.ObjectItemList.Adjustment.CurrencyID')],
      [u'Simpul TERDALAM modul ini. Bernama Adjustment tetapi BUKAN usulan -',
       u'   ia rincian mata uang per item objek. Nama yang sama, arti berbeda.'])

tabel('T_KOMITE_QUOTATION', u'1:1', 'T_GENERAL_KOMITE', 'DATA_KOMITE_ID',
      u'QuotationData{}', 'tipis',
      [K(u'INSURED_NAME', u'pyWorkPage.QuotationData.InsuredName')])

tabel('T_KOMITE_TREATY_MASTER', u'1:1', 'T_GENERAL_KOMITE', 'DATA_KOMITE_ID',
      u'TreatyInMaster{}', 'tipis',
      [K(u'REPORTING_START', u'pyWorkPage.TreatyInMaster.ReportingStart')])

# Tiap tabel 1:N berasal dari page list, dan baris page list dikenali POSISInya.
# SUBSCRIPT dipasang di semuanya supaya urutan tidak hilang saat didatarkan -
# KomiteRouter langkah 6 beriterasi menurut urutan itu, bukan menurut kolom lain.
for t in T:
    if t['kard'] == u'1:N' and not any(k[0] == u'SUBSCRIPT' for k in t['kolom']):
        t['kolom'].append(SUB)

PETA = dict((t['nama'], t) for t in T)

RELASI = [(t['induk'], t['nama'], t['kunci'], t['kard']) for t in T if t['induk']]
RELASI.append(('T_GENERAL_KOMITE', 'T_CLAIM_ADJUSTMENT', 'ADJUSTMENT_ID', u'1:1 silang'))

# ------------------------------------------------------------ tabel datar
baris_datar = []
for t in T:
    n = t['nama']
    # Kolom bernama ID yang sudah ada di pohon TIDAK diduplikasi sebagai baris PK
    # generik; ia yang ditandai PK, lengkap dengan jalur Pega asalnya.
    id_sendiri = [k for k in t['kolom'] if k[0] == u'ID']
    rpk, bpk = ref(t['pk_jalur'])
    if t['kunci'] and u'SHARED PK' in t['kunci']:
        baris_datar.append([n, u'ID', u'PK', u'SHARED PK = %s.ID' % t['induk'],
                            t['pk_jalur'], rpk or u'', bpk or u'', t['pk_yakin']])
    elif not id_sendiri:
        baris_datar.append([n, u'ID', u'PK', u'kunci primer',
                            t['pk_jalur'], rpk or u'', bpk or u'', t['pk_yakin']])
    if t['induk'] and not (t['kunci'] and u'SHARED PK' in t['kunci']):
        rfk, bfk = ref(t['fk_jalur'])
        baris_datar.append([n, t['kunci'], u'FK', u'-> %s.ID' % t['induk'],
                            t['fk_jalur'], rfk or u'', bfk or u'', t['fk_yakin']])
    for (kol, jalur, yakin) in t['kolom']:
        r, b = ref(jalur)
        kunci = u'PK' if (kol == u'ID' and id_sendiri) else u''
        ket = u'kunci primer' if kunci else u''
        baris_datar.append([n, kol, kunci, ket, jalur, r or u'', b or u'', yakin])

with io.open(os.path.join(ERD, 'datar-komite-tabel.csv'), 'w', encoding='utf-8', newline='') as f:
    w = csv.writer(f)
    w.writerow(['TABEL', 'KOLOM', 'KUNCI', 'KETERANGAN_KUNCI', 'JALUR_PEGA', 'REF', 'CACAH_BERKAS', 'KEYAKINAN'])
    for b in baris_datar:
        w.writerow(b)

with io.open(os.path.join(ERD, 'RELASI-KOMITE.csv'), 'w', encoding='utf-8', newline='') as f:
    w = csv.writer(f)
    w.writerow(['INDUK', 'ANAK', 'KOLOM_FK', 'KARDINALITAS'])
    for r in RELASI:
        w.writerow(list(r))

# ------------------------------------------------------------------ HTML
def esc(s):
    return (u'%s' % s).replace('&', '&amp;').replace('<', '&lt;').replace('>', '&gt;')

WARNA = {'akar': 'w-akar', 'komite': 'w-komite', 'salinan': 'w-salinan', 'tipis': 'w-tipis'}
LEBAR, T_KEPALA, T_BARIS = 430, 26, 15.2
JARAK_X, JARAK_Y, MX, MY = 470, 24, 30, 96
# Lebar kotak 430px, teks monospace 10.2px ~ 6.13px/karakter, sisa 20px padding.
MUAT = 66

def bungkus(teks, lebar=MUAT, gantung=u'     '):
    """Potong baris panjang di batas kata supaya tidak melimpah keluar kotak."""
    kata, baris, kini = teks.split(), [], u''
    for k in kata:
        calon = (kini + u' ' + k).strip() if kini else k
        if len(calon) > lebar and kini:
            baris.append(kini)
            kini = gantung + k
        else:
            kini = calon
    if kini:
        baris.append(kini)
    return baris or [u'']

def baris_kotak(t):
    b = []
    if t['kunci'] and u'SHARED PK' in t['kunci']:
        b.append(u'PK   ID  =  %s.ID' % t['induk'])
    else:
        b.append(u'PK   ID')
        if t['induk']:
            b.append(u'FK   %s  ->  %s.ID' % (t['kunci'], t['induk']))
    for (kol, jalur, yakin) in t['kolom']:
        r, _ = ref(jalur)
        tanda = u'*' if yakin == u'rujukan-relatif' else u' '
        b.append(u'%s %-28s n=%s' % (tanda, kol, r))
    for c in t['catatan']:
        b += bungkus(c)
    b += bungkus(u'<-  %s' % t['asal'])
    return b

def tinggi(t):
    return T_KEPALA + len(baris_kotak(t)) * T_BARIS + 8

def kedalaman(n, lihat=None):
    lihat = lihat or set()
    if n in lihat:
        return 0
    lihat.add(n)
    t = PETA[n]
    return 0 if not t['induk'] or t['induk'] not in PETA else 1 + kedalaman(t['induk'], lihat)

tingkat = collections.defaultdict(list)
for n in PETA:
    tingkat[kedalaman(n)].append(n)
for d in tingkat:
    tingkat[d].sort()

POS, x = {}, MX
maks_y = 0
for d in sorted(tingkat):
    y = MY
    for n in tingkat[d]:
        POS[n] = (x, y)
        y += tinggi(PETA[n]) + JARAK_Y
    maks_y = max(maks_y, y)
    x += JARAK_X
LEBAR_SVG, TINGGI_SVG = x + MX, maks_y + 40

svg = []
for (induk, anak, kol, kard) in RELASI:
    if induk not in POS or anak not in POS:
        continue
    x1, y1 = POS[induk]
    x2, y2 = POS[anak]
    x1 += LEBAR
    y1 += tinggi(PETA[induk]) / 2.0
    y2 += tinggi(PETA[anak]) / 2.0
    mid = (x1 + x2) / 2.0
    svg.append(u'<path class="rel r-pohon" d="M%.1f %.1f C%.1f %.1f %.1f %.1f %.1f %.1f"><title>%s -> %s (%s, %s)</title></path>'
               % (x1, y1, mid, y1, mid, y2, x2, y2, esc(induk), esc(anak), esc(kol), esc(kard)))

for n, (x0, y0) in POS.items():
    t = PETA[n]
    h = tinggi(t)
    svg.append(u'<g class="%s">' % WARNA[t['kelompok']])
    svg.append(u'<rect class="badan" x="%d" y="%d" width="%d" height="%.1f" rx="5"/>' % (x0, y0, LEBAR, h))
    svg.append(u'<path class="kepala" d="M%d %d h%d a5 5 0 0 1 5 5 v%d h-%d v-%d a5 5 0 0 1 5 -5 z"/>'
               % (x0 + 5, y0, LEBAR - 10, T_KEPALA - 5, LEBAR, T_KEPALA - 5))
    svg.append(u'<text class="t-nama" x="%d" y="%d">%s</text>' % (x0 + 10, y0 + 18, esc(n)))
    svg.append(u'<text class="t-kard" x="%d" y="%d">%s</text>' % (x0 + LEBAR - 10, y0 + 18, esc(t['kard'])))
    yy = y0 + T_KEPALA + 12
    for b in baris_kotak(t):
        k = 'b-ket'
        if b.startswith('PK'):
            k = 'b-pk'
        elif b.startswith('FK'):
            k = 'b-fk'
        elif b.startswith('<-'):
            k = 'b-asal'
        elif b.startswith('AWAS') or b.startswith('SALINAN') or b.startswith('NOL'):
            k = 'b-awas'
        svg.append(u'<text class="%s" x="%d" y="%.1f" xml:space="preserve">%s</text>' % (k, x0 + 10, yy, esc(b)))
        yy += T_BARIS
    svg.append(u'</g>')

BARIS_TBL = u''.join(
    u'<tr><td><code>%s</code></td><td><code>%s</code></td><td>%s</td><td>%s</td><td><code>%s</code></td><td class="num">%s</td><td>%s</td></tr>'
    % (esc(b[0]), esc(b[1]), esc(b[2]), esc(b[3]), esc(b[4]), b[5], esc(b[7]))
    for b in baris_datar)

HTML = u"""<!DOCTYPE html><html lang="id"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Tabel datar dan ERD - Komite Claim Non Prop (sistem lama)</title><style>
:root{--tinta:#16181d;--tinta-2:#525a68;--tinta-3:#7b8494;--latar:#fff;--latar-2:#f7f8fa;
--garis:#d8dce4;--badan:#fdfdfe;--r-pohon:#1f4fd8;
--akar:#1e3a8a;--akar-t:#fff;--komite:#ddc6f2;--komite-t:#2e1b45;
--salinan:#f6b93b;--salinan-t:#3d2c00;--tipis:#d5d9e0;--tipis-t:#2b3039;
--pk:#b45309;--fk:#1d4ed8;--awas:#c0262b}
@media(prefers-color-scheme:dark){:root:not([data-theme="light"]){
--tinta:#e9ebf0;--tinta-2:#aeb6c3;--tinta-3:#8c94a3;--latar:#14161a;--latar-2:#1b1e24;
--garis:#2f343d;--badan:#1b1e24;--r-pohon:#7fa2ff;
--akar:#2a4fa8;--akar-t:#fff;--komite:#3d2b57;--komite-t:#e4d3f7;
--salinan:#8a6a1c;--salinan-t:#ffe9b8;--tipis:#2b303a;--tipis-t:#cdd3dc;
--pk:#e0a352;--fk:#8fabff;--awas:#ff8a8a}}
*{box-sizing:border-box}body{margin:0;background:var(--latar);color:var(--tinta);
font:15px/1.55 ui-sans-serif,-apple-system,"Segoe UI",Roboto,Helvetica,Arial,sans-serif}
header{border-bottom:1px solid var(--garis);padding:22px 18px 16px}.bingkai{padding:0 18px 56px}
h1{font-size:21px;margin:0 0 8px}h2{font-size:16px;margin:26px 0 8px}
p.ket{color:var(--tinta-2);font-size:13px;max-width:92ch;margin:0 0 8px}
.awas{display:inline-block;margin-top:8px;padding:6px 10px;border-radius:4px;background:var(--latar-2);
border-left:3px solid var(--awas);font-size:13px;max-width:92ch}
code{font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;font-size:.92em}
.legenda{display:flex;flex-wrap:wrap;gap:7px 18px;margin:10px 0;font-size:12.4px;color:var(--tinta-2)}
.legenda div{display:flex;align-items:center;gap:6px}
.kot{width:14px;height:14px;border-radius:3px;display:inline-block}
.bungkus{overflow:auto;border:1px solid var(--garis);border-radius:6px;background:var(--latar-2);padding:8px}
.erd{display:block}
.erd .badan{fill:var(--badan);stroke-width:1.2}.erd .kepala{stroke:none}
.erd .t-nama{font-size:12.4px;font-weight:700}.erd .t-kard{font-size:10.4px;text-anchor:end;opacity:.85}
.erd text{font-size:10.2px;font-family:ui-monospace,Menlo,Consolas,monospace}
.erd .b-pk{fill:var(--pk);font-weight:700}.erd .b-fk{fill:var(--fk);font-weight:600}
.erd .b-awas{fill:var(--awas)}.erd .b-asal{fill:var(--tinta-3);font-style:italic}
.erd .b-ket{fill:var(--tinta-2)}
.erd .rel{fill:none;stroke-width:1.3;opacity:.85}.erd .rel:hover{stroke-width:2.8;opacity:1}
.erd .r-pohon{stroke:var(--r-pohon)}
.w-akar .badan{stroke:var(--akar)}.w-akar .kepala{fill:var(--akar)}
.w-akar .t-nama,.w-akar .t-kard{fill:var(--akar-t)}
.w-komite .badan{stroke:#8e5cc0}.w-komite .kepala{fill:var(--komite)}
.w-komite .t-nama,.w-komite .t-kard{fill:var(--komite-t)}
.w-salinan .badan{stroke:#b8860b}.w-salinan .kepala{fill:var(--salinan)}
.w-salinan .t-nama,.w-salinan .t-kard{fill:var(--salinan-t)}
.w-tipis .badan{stroke:#69717f;stroke-dasharray:5 3}.w-tipis .kepala{fill:var(--tipis)}
.w-tipis .t-nama,.w-tipis .t-kard{fill:var(--tipis-t)}
table{border-collapse:collapse;width:100%;font-size:12.6px}
th{position:sticky;top:0;background:var(--latar-2);text-align:left;font-weight:600;padding:6px 10px;
border-bottom:1px solid var(--garis);color:var(--tinta-2)}
td{padding:4px 10px;border-bottom:1px solid var(--garis);vertical-align:top}
td.num,th.num{text-align:right;font-variant-numeric:tabular-nums}
.tbl{max-height:60vh;overflow:auto;border:1px solid var(--garis);border-radius:6px}
footer{color:var(--tinta-3);font-size:12px;margin-top:26px;border-top:1px solid var(--garis);padding-top:12px}
</style></head><body>
<header><h1>Tabel datar dan ERD &mdash; <b>Komite Claim Non Prop</b>, sistem lama</h1>
<p class="ket"><b>STEMPEL ASAL.</b> Akar <code>T_WORK_CLAIM</code> &middot; __N_TABEL__ tabel
&middot; __N_RELASI__ relasi &middot; __N_KOLOM__ kolom. Pohonnya dari
<code>struktur-komite-lama.md</code> / <code>datar-komite-lama.csv</code> &mdash; 688 simpul, sapuan
<b>59 berkas XML</b> seluruhnya. Konvensi penamaan dari
<code>Diagram-Skema-Tabel-NusantaraRe.xlsx</code> lewat
<code>claim-non-prop/alat/buat-skema-claimnp.py</code>.
Dibangkitkan <code>alat/buat-skema-komite.py</code> pada <b>__TANGGAL__</b>.</p>
<p class="awas"><b>Ini pendataran SISTEM LAMA, bukan rancangan skema baru.</b>
Rancangan yang berlaku hidup di <code>SPEC-KOMITE-01.md</code> bagian <i>Model data</i> &mdash;
<b>tujuh objek</b> <code>KLAIMNP</code>, 64 kolom. Empat belas kotak di bawah ini memerikan
bentuk yang <b>ada</b>, termasuk cacatnya: dua penyimpan untuk satu fakta, dua kolom tanggal
untuk satu fakta, dan enam tabel salinan yang <b>tidak boleh</b> ikut dibuat.</p>
<div class="legenda">
<div><span class="kot" style="background:var(--akar)"></span>akar lintas-lini</div>
<div><span class="kot" style="background:var(--komite)"></span>milik modul Komite</div>
<div><span class="kot" style="background:var(--salinan)"></span>SALINAN &mdash; rumahnya di sisi Klaim</div>
<div><span class="kot" style="background:var(--tipis)"></span>tipis &mdash; nol atau satu kolom; jangan dibuat</div>
<div><code>*</code> kolom bertingkat <i>rujukan-relatif</i></div>
<div><code>n=</code> cacah rujukan di 59 berkas</div>
</div></header>
<div class="bingkai">
<h2>ERD</h2>
<div class="bungkus"><svg class="erd" viewBox="0 0 __W__ __H__" width="__W__" height="__H__"
xmlns="http://www.w3.org/2000/svg" role="img" aria-label="ERD tabel Komite Claim Non Prop">
__SVG__
</svg></div>
<h2>Tabel datar &mdash; satu baris per kolom</h2>
<p class="ket">Kolom <code>REF</code> adalah cacah rujukan mentah di 59 berkas XML; ia penanda
seberapa hidup sebuah kolom, bukan penanda penting. Kolom berkeyakinan
<i>rujukan-relatif</i> tidak pernah ditulis dengan jalur penuh.</p>
<div class="tbl"><table><thead><tr><th>TABEL</th><th>KOLOM</th><th>KUNCI</th>
<th>KETERANGAN</th><th>JALUR PEGA</th><th class="num">REF</th><th>KEYAKINAN</th></tr></thead>
<tbody>__BARIS__</tbody></table></div>
<footer>TURUNAN. Dibangkitkan <code>alat/buat-skema-komite.py</code>.
Bila berkas ini berbeda dari <code>datar-komite-lama.csv</code>, <b>alatnya yang salah</b>.
Nol DDL dijalankan. Nol berkas sumber berubah.</footer>
</div></body></html>"""

n_kolom = sum(len(t['kolom']) for t in T)
HTML = (HTML.replace('__SVG__', u'\n'.join(svg))
            .replace('__W__', str(int(LEBAR_SVG))).replace('__H__', str(int(TINGGI_SVG)))
            .replace('__BARIS__', BARIS_TBL).replace('__TANGGAL__', TANGGAL)
            .replace('__N_TABEL__', str(len(T))).replace('__N_RELASI__', str(len(RELASI)))
            .replace('__N_KOLOM__', str(n_kolom)))
# SVG yang cacat tergambar sebagai kotak KOSONG - tanpa pesan galat, tanpa
# tanda apa pun. Satu tag tak tertutup pernah menelan seluruh teks di dalam
# kotaknya. Karena itu SVG diurai dulu; keluaran yang tidak sah dihentikan
# di sini alih-alih dikirim dan ditemukan mata.
import xml.etree.ElementTree as ET
_svg = u'<svg xmlns="http://www.w3.org/2000/svg">%s</svg>' % u'\n'.join(svg)
try:
    _akar = ET.fromstring(_svg.encode('utf-8'))
except ET.ParseError as e:
    print(u'GAGAL: SVG tidak sah - %s' % e)
    sys.exit(1)
_n_text = len(_akar.findall('.//{http://www.w3.org/2000/svg}text'))
_harus = sum(len(baris_kotak(t)) + 2 for t in T)
if _n_text != _harus:
    print(u'GAGAL: cacah <text> %d, seharusnya %d' % (_n_text, _harus))
    sys.exit(1)

with io.open(os.path.join(ERD, 'ERD-KOMITE-LAMA.html'), 'w', encoding='utf-8') as f:
    f.write(HTML)

# ----------------------------------------------------------------- Excel
wb = openpyxl.Workbook()
tipis = Side(style='thin', color='D8DCE4')
BR = Border(left=tipis, right=tipis, top=tipis, bottom=tipis)
KEPALA = Font(bold=True, color='FFFFFF', size=10)
ISI_KEPALA = PatternFill('solid', fgColor='1E3A8A')

ws = wb.active
ws.title = 'BACA-DULU'
for i, t in enumerate([
        u'Tabel datar dan ERD - Komite Claim Non Prop, SISTEM LAMA',
        u'',
        u'TURUNAN. Dibangkitkan alat/buat-skema-komite.py pada %s.' % TANGGAL,
        u'Sumber: 4-erd-dan-tabel-datar/datar-komite-lama.csv (688 simpul, 59 berkas XML).',
        u'Konvensi nama: Diagram-Skema-Tabel-NusantaraRe.xlsx lewat buat-skema-claimnp.py.',
        u'',
        u'INI BUKAN RANCANGAN SKEMA BARU.',
        u'Rancangan yang berlaku: SPEC-KOMITE-01.md bagian Model data - tujuh objek KLAIMNP.',
        u'',
        u'Tabel   : %d' % len(T),
        u'Relasi  : %d' % len(RELASI),
        u'Kolom   : %d' % n_kolom,
        u'',
        u'Nol DDL dijalankan. Nol berkas sumber berubah.'], start=1):
    ws.cell(row=i, column=1, value=t)
ws.column_dimensions['A'].width = 100

ws = wb.create_sheet('TABEL-DATAR')
hdr = ['TABEL', 'KOLOM', 'KUNCI', 'KETERANGAN_KUNCI', 'JALUR_PEGA', 'REF', 'CACAH_BERKAS', 'KEYAKINAN']
for c, h in enumerate(hdr, start=1):
    s = ws.cell(row=4, column=c, value=h)
    s.font, s.fill, s.border = KEPALA, ISI_KEPALA, BR
for r, b in enumerate(baris_datar, start=5):
    for c, v in enumerate(b, start=1):
        s = ws.cell(row=r, column=c, value=v)
        s.border = BR
for col, w in zip('ABCDEFGH', (26, 30, 7, 30, 62, 7, 7, 17)):
    ws.column_dimensions[col].width = w
ws.auto_filter.ref = 'A4:H%d' % (4 + len(baris_datar))
ws.freeze_panes = 'A5'

ws = wb.create_sheet('RELASI')
for c, h in enumerate(['INDUK', 'ANAK', 'KOLOM_FK', 'KARDINALITAS'], start=1):
    s = ws.cell(row=4, column=c, value=h)
    s.font, s.fill, s.border = KEPALA, ISI_KEPALA, BR
for r, rel in enumerate(RELASI, start=5):
    for c, v in enumerate(rel, start=1):
        s = ws.cell(row=r, column=c, value=v)
        s.border = BR
for col, w in zip('ABCD', (28, 30, 22, 14)):
    ws.column_dimensions[col].width = w
ws.freeze_panes = 'A5'

ws = wb.create_sheet('DAFTAR-TABEL')
for c, h in enumerate(['TABEL', 'KARDINALITAS', 'INDUK', 'KUNCI', 'CACAH_KOLOM', 'KELOMPOK', 'ASAL_PEGA'], start=1):
    s = ws.cell(row=4, column=c, value=h)
    s.font, s.fill, s.border = KEPALA, ISI_KEPALA, BR
for r, t in enumerate(T, start=5):
    for c, v in enumerate([t['nama'], t['kard'], t['induk'] or u'-', t['kunci'] or u'-',
                           len(t['kolom']), t['kelompok'], t['asal']], start=1):
        s = ws.cell(row=r, column=c, value=v)
        s.border = BR
for col, w in zip('ABCDEFG', (28, 12, 26, 26, 12, 10, 56)):
    ws.column_dimensions[col].width = w
ws.freeze_panes = 'A5'

# Workbook sering sedang dibuka Excel saat alat ini dijalankan ulang. Bila
# terkunci, simpan ke nama berdampingan alih-alih mati setelah CSV dan HTML
# sudah tertulis - keluaran separuh yang diam lebih mahal daripada pesan ini.
XLSX = os.path.join(ERD, 'Diagram-Skema-Tabel-Komite.xlsx')
try:
    wb.save(XLSX)
except (IOError, OSError) as e:
    XLSX = os.path.join(ERD, 'Diagram-Skema-Tabel-Komite.BARU.xlsx')
    wb.save(XLSX)
    print(u'AWAS: workbook lama terkunci (%s). Disimpan sebagai %s.'
          % (e.__class__.__name__, os.path.basename(XLSX)))

# -------------------------------------------------------------- ringkasan
print(u'tabel  : %d' % len(T))
print(u'relasi : %d' % len(RELASI))
print(u'kolom  : %d' % n_kolom)
print(u'baris tabel datar : %d' % len(baris_datar))
for k in ('akar', 'komite', 'salinan', 'tipis'):
    n = [t['nama'] for t in T if t['kelompok'] == k]
    print(u'  %-8s %d  %s' % (k, len(n), ' '.join(n)))
kosong = [t['nama'] for t in T if not t['kolom']]
print(u'tabel tanpa kolom sendiri : %d  %s' % (len(kosong), ' '.join(kosong)))
