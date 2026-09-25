# -*- coding: utf-8 -*-
r"""
buat-peta-nama-tabel-treatyin.py — tabel datar modul Treaty In, penamaan T_<NAMA>.

Bentuknya meniru `_migration-docs/claim-non-prop/4-erd-dan-tabel-datar/peta-nama-tabel.tsv`:
kolom NAMA_T memegang nama TABEL DATAR (bentuk pipih turunan pohon clipboard Pega),
kolom PADANAN_DDL memegang nama OBJEK ORACLE yang tunduk pada §16. Keduanya berlaku di
lapisan yang berbeda dan tidak saling menggantikan.

Sumber isinya: `4-erd-dan-tabel-datar/struktur-treatyin-lama.md` §4-§7 dan
`datar-treatyin-lama.csv`, keduanya turunan dari `alat/buat-pohon-treatyin.py`.

Keluaran (TURUNAN; bila beda dari definisinya, ALAT INI yang salah):
  4-erd-dan-tabel-datar/peta-nama-tabel-treatyin.tsv
  4-erd-dan-tabel-datar/PETA-NAMA-TABEL-TREATYIN.md

Jalankan:  PYTHONIOENCODING=utf-8 python alat/buat-peta-nama-tabel-treatyin.py
"""
import os, io, csv, collections

KELUARAN = os.path.join(r'D:\XML_NURE', '_migration-docs', 'treaty-in',
                        '4-erd-dan-tabel-datar')

# NAMA_T, PADANAN_DDL, KELOMPOK, INDUK, KUNCI, ASAL_PEGA
TABEL = [
    # ---------------------------------------------------------------- akar
    # KEPALA MEMAKAI NAMA YANG SUDAH ADA: `TREATY_IN` adalah tabel Oracle yang hidup
    # sekarang, bukan nama baru. Karena itu ia TIDAK berawalan T_ \u2014 awalan itu
    # menandai tabel datar yang kita turunkan, dan tabel ini bukan turunan kita.
    ('TREATY_IN', 'KONTRAK + VERSI_KONTRAK', 'akar', '', '',
     'TreatyIn{} \u00b7 99 skalar akar \u00b7 kelas ASM-FW-GISFW-Int-TREATY_IN \u00b7 '
     'TABEL SUDAH ADA \u2014 20 kolom bisnis nyata, 79 sisanya hanya di M_TREATY_IN.JSONDATA'),
    ('T_TREATY_REVISION', 'VERSI_KONTRAK \u00b7 bagian addendum', 'edm',
     'TREATY_IN', 'TREATY_IN_ID',
     'TreatyIn.OLDID \u00b7 EDMState \u00b7 EDMEffective \u00b7 EDMMaterialType '
     '\u00b7 RevisionState \u00b7 RevisionDate \u00b7 AddendumPremi \u00b7 '
     'padanan lamanya TREATY_IN_EDM / M_TREATY_IN_EDM \u2014 TIDAK dipakai sebagai nama '
     'karena tabel itu bentuk SALINAN, bukan bagian addendum'),

    # ---------------------------------------------------------------- layer
    ('T_TREATY_LIMITS', 'LAYER', 'limit', 'TREATY_IN', 'TREATY_IN_ID',
     'TreatyIn.Limits[] \u00b7 18 field \u00b7 \u2605 inti layer, n=307'),
    ('T_TREATY_LIMIT_DETAIL', 'DETAIL_PROPORSIONAL', 'limit', 'T_TREATY_LIMITS',
     'LIMIT_ID',
     'TreatyIn.Limits[].Detail[] \u00b7 14 field \u00b7 \u2605 inti proporsional, n=104'),
    ('T_TREATY_LIMIT_COB', 'KELAS_BISNIS_LAYER', 'limit', 'T_TREATY_LIMIT_DETAIL',
     'LIMIT_DETAIL_ID',
     'Limits[].Detail[].COBList[] \u00b7 ClassOfBusiness \u00b7 ClassOfBusinessID'),
    ('T_TREATY_LIMIT_AMOUNT', 'NILAI_LAYER', 'limit', 'T_TREATY_LIMIT_DETAIL',
     'LIMIT_DETAIL_ID',
     'Detail[] \u00b7 EPIList[] RetentionList[] CashLossList[] ClaimCoopList[] PLAList[] '
     'IOOLimitList[] \u2014 SATU tabel berkolom PERAN, bukan enam tabel'),
    ('T_TREATY_LIMIT_ACHIEVEMENT', 'PENCAPAIAN', 'limit', 'T_TREATY_LIMIT_DETAIL',
     'LIMIT_DETAIL_ID',
     'Limits[].Detail[].AchievementLists[] \u00b7 n=1 \u00b7 '
     '\u26a0 kelas TIDAK dideklarasikan di mana pun'),
    ('T_TREATY_LIMIT_MEASURE', 'BESARAN_LAYER', 'limit', 'T_TREATY_LIMITS', 'LIMIT_ID',
     'Limits[] \u00b7 MDPList[] EgnpiTotalList[] PremiumEarnedList[] '
     '\u2014 SATU tabel berkolom PERAN'),
    ('T_TREATY_REINSTATEMENT', 'PEMULIHAN_LIMIT', 'limit', 'T_TREATY_LIMITS', 'LIMIT_ID',
     'Limits[].Reinstatement_List[] \u00b7 ReinstatementPct \u00b7 AdditionalAmount1/2 '
     '\u00b7 nama ber-garis-bawah di sumbernya'),
    ('T_TREATY_LIMIT_GROUP', 'KELOMPOK_LAYER', 'limit', 'T_TREATY_LIMITS', 'LIMIT_ID',
     'Limits[].TreatyGroupList[] \u00b7 TreatyGroup \u00b7 TreatyGroupID'),
    ('T_TREATY_LIMIT_GROUP_COB', 'KELAS_BISNIS_KELOMPOK', 'limit', 'T_TREATY_LIMIT_GROUP',
     'LIMIT_GROUP_ID',
     'Limits[].TreatyGroupList[].ClassOfBusinessList[] \u25c4 KEDALAMAN 4'),

    # ---------------------------------------------------------------- share
    ('T_TREATY_SHARE', 'BAGIAN', 'share', 'TREATY_IN', 'TREATY_IN_ID',
     'TreatyIn.Share[] \u00b7 10 field \u00b7 \u2605 inti penyebaran, n=237'),
    ('T_TREATY_SHARE_SPREADING', 'PENYEBARAN_XOL', 'share', 'T_TREATY_SHARE', 'SHARE_ID',
     'Share[].SpreadingListXOL[] \u00b7 Pct \u00b7 ReinsTypeName \u00b7 ReinsTypeID, n=45'),
    ('T_TREATY_SHARE_SPREAD_AMOUNT', 'NILAI_TERSEBAR', 'share',
     'T_TREATY_SHARE_SPREADING', 'SHARE_SPREADING_ID',
     '11 daftar RNMSpreadedList\u2026[] di bawah SpreadingListXOL \u25c4 KEDALAMAN 4 '
     '\u2014 SATU tabel berkolom PERAN'),
    ('T_TREATY_SHARE_AMOUNT', 'NILAI_BAGIAN', 'share', 'T_TREATY_SHARE', 'SHARE_ID',
     '17 daftar uang langsung di bawah Share[] \u2014 GrossPremiumList NetPremiumList '
     'RnmLimitList BrokerageList RNMSpreadedList\u2026 \u2014 SATU tabel berkolom PERAN'),
    ('T_TREATY_SHARE_DEDUCTION', 'POTONGAN', 'share', 'T_TREATY_SHARE', 'SHARE_ID',
     'Share[].DeductionList[] \u00b7 Deduction \u00b7 DeductionPct \u00b7 Comment'),

    # ---------------------------------------------------------------- fakultatif
    ('T_TREATY_FAC_SHARE', 'BAGIAN_FAKULTATIF', 'fac', 'TREATY_IN',
     'TREATY_IN_ID',
     'TreatyIn.FacultativeShareList[] \u00b7 9 field \u00b7 kelas SAMA dengan Share[], n=58'),
    ('T_TREATY_FAC_SHARE_AMOUNT', 'NILAI_BAGIAN_FAKULTATIF', 'fac', 'T_TREATY_FAC_SHARE',
     'FAC_SHARE_ID',
     'FacultativeShareList[] \u00b7 GrossPremiumList[] NetPremiumList[] RnmLimitList[] '
     'RnmGrossPremiDisplay[] RnmLimitListDisplay[] \u2014 SATU tabel berkolom PERAN'),
    ('T_TREATY_FAC_SHARE_DEDUCTION', 'POTONGAN_FAKULTATIF', 'fac', 'T_TREATY_FAC_SHARE',
     'FAC_SHARE_ID', 'FacultativeShareList[].DeductionList[]'),
    ('T_TREATY_FAC_REINSURER', 'REASURADUR_FAKULTATIF', 'fac', 'TREATY_IN',
     'TREATY_IN_ID',
     'TreatyIn.ShareFacultativeReinsurers[] (n=47) DAN '
     'FacultativeShareList[].ShareFacultativeReinsurers[] (n=6) \u2014 kelas sama, dua induk'),
    ('T_TREATY_FAC_LIMITS', 'LAYER_FAKULTATIF', 'fac', 'T_TREATY_FAC_REINSURER',
     'FAC_REINSURER_ID',
     'ShareFacultativeReinsurers[].FacultativeLimits[] \u00b7 TreatyType \u00b7 TreatyTypeID'),
    ('T_TREATY_FAC_LIMIT_DETAIL', 'DETAIL_PROPORSIONAL_FAKULTATIF', 'fac',
     'T_TREATY_FAC_LIMITS', 'FAC_LIMIT_ID',
     'FacultativeLimits[].Detail[] \u00b7 11 field \u25c4 KEDALAMAN 4 \u00b7 '
     'kelas SAMA dengan Limits[].Detail[]'),

    # ---------------------------------------------------------------- kepala
    ('T_TREATY_CURRENCY', 'MATA_UANG_KONTRAK', 'kepala', 'TREATY_IN',
     'TREATY_IN_ID',
     'TreatyIn.CurrencyList[] \u00b7 Currency \u00b7 CurrencyID \u00b7 Conversion '
     '\u00b7 PeriodStart \u00b7 PeriodEnd'),
    ('T_TREATY_EGNPI', 'EGNPI', 'kepala', 'TREATY_IN', 'TREATY_IN_ID',
     'TreatyIn.EGNPI[] \u00b7 9 field \u00b7 Amount \u00b7 AmountIDR \u00b7 Proportion'),
    ('T_TREATY_RETENTION', 'RETENSI', 'kepala', 'TREATY_IN', 'TREATY_IN_ID',
     'TreatyIn.Retention[] \u00b7 6 field \u00b7 '
     '\u26a0 memakai kelas TreatyInEGNPI, bukan kelas sendiri'),
    ('T_TREATY_LIMIT_SUMMARY', 'RINGKASAN_LIMIT', 'kepala', 'TREATY_IN',
     'TREATY_IN_ID',
     'LimitSummaryList[] LimitShareSummaryList[] LimitFacShareSummaryList[] '
     'MDPSummaryList[] \u2014 SATU tabel berkolom LINGKUP'),
    ('T_TREATY_INSTALLMENT', 'ANGSURAN', 'kepala', 'TREATY_IN', 'TREATY_IN_ID',
     'TreatyIn.Installment[] \u00b7 Currency \u00b7 ID'),
    ('T_TREATY_INSTALLMENT_ITEM', 'RINCIAN_ANGSURAN', 'kepala', 'T_TREATY_INSTALLMENT',
     'INSTALLMENT_ID',
     'Installment[].InstallmentList[] \u00b7 DueDate \u00b7 InstallmentPct \u00b7 '
     'PaymentDate \u00b7 WPC   [bukti: identitas-kelas \u2014 bentuknya hanya terbaca '
     'lewat salinan OLDDATA/ValueDifference]'),
    ('T_TREATY_REPORTING_PERIOD', 'PERIODE_PELAPORAN', 'kepala', 'TREATY_IN',
     'TREATY_IN_ID',
     'TreatyIn.ReportingPeriodList[] \u00b7 Period \u00b7 InitialDate \u00b7 '
     'SubmissionDue \u00b7 ConfirmationDue \u00b7 SettlementDue'),
    ('T_TREATY_ACCUMULATION', 'AKUMULASI', 'kepala', 'TREATY_IN',
     'TREATY_IN_ID', 'TreatyIn.AccumulationList[] \u00b7 Period \u00b7 ReportDate'),
    ('T_TREATY_PORTFOLIO', 'PORTOFOLIO', 'kepala', 'TREATY_IN', 'TREATY_IN_ID',
     'TreatyIn.Portfolio[] \u00b7 Type \u00b7 TypePortfolio \u00b7 Description'),
    ('T_TREATY_HAZARD_LIMIT', 'BATAS_BAHAYA', 'kepala', 'TREATY_IN',
     'TREATY_IN_ID',
     'Earthquake FloodJab FloodNation RSMDLimit MaxCoGroup MaxCoNonGroup + '
     'Currency\u2026 \u2014 10 skalar berpasangan DIPUTAR jadi daftar'),
    ('T_TREATY_TOTAL', 'REKAP_KONTRAK', 'kepala', 'TREATY_IN', 'TREATY_IN_ID',
     '20 daftar Total\u2026[] tingkat akar + TotalEgnpiAmountNP + TotalRetentionAmountNP '
     '\u2014 SATU tabel berkolom PERAN'),
    ('T_TREATY_RETRO_SHARE', 'BAGIAN_RETRO', 'kepala', 'TREATY_IN',
     'TREATY_IN_ID',
     'TreatyIn.ShareReins[] + RetroList \u00b7 '
     '\u26a0 kelas TIDAK dideklarasikan di mana pun'),
    ('T_VIEW_COMMENT', 'CATATAN_PERSETUJUAN', 'bersama', 'TREATY_IN',
     'TREATY_IN_ID',
     'TreatyIn.CommentList[] \u00b7 OperatorName \u00b7 Date \u00b7 Suggest \u00b7 '
     'IsApproved \u00b7 kelas pinjaman SuggestList'),
    ('DOCUMENT_TREATY_IN', 'DOKUMEN_KONTRAK', 'bersama', 'TREATY_IN',
     'TREATY_IN_ID',
     'M_ATTACHMENTTREATY_2 + T_STORAGE_IMAGE + CATEGORY_ATTACH_REAS \u2014 '
     'lewat RDB List, bukan lewat pohon clipboard'),

    # ---------------------------------------------------------------- selisih
    ('T_TREATY_VALUE_DIFFERENCE', 'NILAI_SELISIH', 'edm', 'T_TREATY_REVISION',
     'REVISION_ID',
     'TreatyIn.ValueDifference{} \u00b7 217 simpul \u00b7 nilai sekarang \u2212 nilai lama '
     '\u00b7 ditulis TreatyEDMCalculateDifference (11 langkah, SELURUHNYA hidup)'),
    ('T_TREATY_VALUE_BEFORE_PRORATE', 'NILAI_SEBELUM_PRO_RATE', 'edm',
     'TREATY_IN', 'TREATY_IN_ID',
     'TreatyIn.ValueBeforeProrate{} \u00b7 15 simpul \u00b7 kelas sama dengan akar'),
    ('\u2014 TIDAK DIBUAT \u2014', 'OLDDATA tidak menjadi tabel', 'edm', '', '',
     'TreatyIn.OLDDATA{} \u00b7 142 simpul \u00b7 '
     'ATURAN BISNIS: data lama TIDAK disimpan ulang \u2014 ia DISELECT'),
    ('\u2014 TIDAK DIBUAT \u2014', 'ActualValue tidak menjadi tabel', 'edm', '', '',
     'TreatyIn.ActualValue{} \u00b7 152 simpul \u00b7 '
     '\u26a0 namanya "nilai aktual" tetapi 8 isian menerima SELISIH \u2014 '
     'lihat struktur-treatyin-lama.md \u00a75.3'),

    # ---------------------------------------------------------------- migrasi
    ('T_TREATY_MIG_CORRELATION', 'MIGRASI_KORELASI', 'migrasi', '', '',
     'jembatan ke sistem lama \u2014 pzInsKey \u00b7 TreatyIn.ID \u00b7 TreatyIn.OLDID'),
    ('T_TREATY_MIG_LANDING', 'MIGRASI_PENDARATAN', 'migrasi', '', '',
     'tempat bentuk lama mendarat utuh \u2014 M_TREATY_IN.JSONDATA apa adanya'),
    ('T_TREATY_MIG_REJECTED', 'MIGRASI_NILAI_DITOLAK', 'migrasi', 'T_TREATY_MIG_LANDING',
     'LANDING_ID',
     'nilai yang tidak dapat diurai \u2014 tercatat, tidak dibulatkan, tidak dibuang'),
    ('T_TREATY_OUTBOUND_ARCHIVE', 'ARSIP_MUATAN_KELUAR', 'migrasi',
     'TREATY_IN', 'TREATY_IN_ID',
     '175 argumen ke PEGA_TREATY_IN / PEGA_M_TREATY_IN_EDM / '
     'PEGA_M_TREATY_IN_DETAIL / \u2026_DETAIL_EDM'),

    # ---------------------------------------------------------------- tak ditemukan
    ('T_TREATY_OFFER', '\u26a0 TIDAK ADA PENULIS HIDUP', 'lubang', '', '',
     'TREATYINOFFER \u2014 nol penulis terjangkau (sapuan dua tingkat). '
     'BUKAN sumber migrasi, BUKAN dasar rekonsiliasi'),
    ('T_TREATY_AUTHORITY_LIMIT', 'BATAS_WEWENANG', 'lubang', '', '',
     'TIDAK DITEMUKAN di pohon \u2014 batas wewenang persetujuan, menunggu dokumen'),
]

# Tambatan tiap tabel datar ke simpul pohon `TreatyIn`. Inilah yang membuat ERD dapat
# dibangkitkan DARI STRUKTUR, bukan dari ingatan: nama fieldnya diambil dari
# datar-treatyin-lama.csv lewat jalur ini. Nilai berupa daftar = tabel itu menggabungkan
# beberapa daftar sebentuk menjadi satu, dengan kolom PERAN.
JALUR_PEGA = {
    'TREATY_IN': ['TreatyIn'],
    'T_TREATY_REVISION': ['TreatyIn'],

    'T_TREATY_LIMITS': ['TreatyIn.Limits'],
    'T_TREATY_LIMIT_DETAIL': ['TreatyIn.Limits.Detail'],
    'T_TREATY_LIMIT_COB': ['TreatyIn.Limits.Detail.COBList'],
    'T_TREATY_LIMIT_AMOUNT': ['TreatyIn.Limits.Detail.EPIList',
                              'TreatyIn.Limits.Detail.RetentionList',
                              'TreatyIn.Limits.Detail.CashLossList',
                              'TreatyIn.Limits.Detail.ClaimCoopList',
                              'TreatyIn.Limits.Detail.PLAList',
                              'TreatyIn.Limits.Detail.IOOLimitList'],
    'T_TREATY_LIMIT_ACHIEVEMENT': ['TreatyIn.Limits.Detail.AchievementLists'],
    'T_TREATY_LIMIT_MEASURE': ['TreatyIn.Limits.MDPList',
                               'TreatyIn.Limits.EgnpiTotalList',
                               'TreatyIn.Limits.PremiumEarnedList'],
    'T_TREATY_REINSTATEMENT': ['TreatyIn.Limits.Reinstatement_List'],
    'T_TREATY_LIMIT_GROUP': ['TreatyIn.Limits.TreatyGroupList'],
    'T_TREATY_LIMIT_GROUP_COB': ['TreatyIn.Limits.TreatyGroupList.ClassOfBusinessList'],

    'T_TREATY_SHARE': ['TreatyIn.Share'],
    'T_TREATY_SHARE_SPREADING': ['TreatyIn.Share.SpreadingListXOL'],
    'T_TREATY_SHARE_SPREAD_AMOUNT': ['TreatyIn.Share.SpreadingListXOL.RNMSpreadedListXOL',
                                     'TreatyIn.Share.SpreadingListXOL.GrossPremiumList'],
    'T_TREATY_SHARE_AMOUNT': ['TreatyIn.Share.GrossPremiumList',
                              'TreatyIn.Share.NetPremiumList',
                              'TreatyIn.Share.RnmLimitList',
                              'TreatyIn.Share.BrokerageList',
                              'TreatyIn.Share.RNMSpreadedListNetXOL'],
    'T_TREATY_SHARE_DEDUCTION': ['TreatyIn.Share.DeductionList'],

    'T_TREATY_FAC_SHARE': ['TreatyIn.FacultativeShareList'],
    'T_TREATY_FAC_SHARE_AMOUNT': ['TreatyIn.FacultativeShareList.GrossPremiumList',
                                  'TreatyIn.FacultativeShareList.RnmLimitList'],
    'T_TREATY_FAC_SHARE_DEDUCTION': ['TreatyIn.FacultativeShareList.DeductionList'],
    'T_TREATY_FAC_REINSURER': ['TreatyIn.ShareFacultativeReinsurers',
                               'TreatyIn.FacultativeShareList.ShareFacultativeReinsurers'],
    'T_TREATY_FAC_LIMITS': ['TreatyIn.ShareFacultativeReinsurers.FacultativeLimits'],
    'T_TREATY_FAC_LIMIT_DETAIL':
        ['TreatyIn.ShareFacultativeReinsurers.FacultativeLimits.Detail'],

    'T_TREATY_CURRENCY': ['TreatyIn.CurrencyList'],
    'T_TREATY_EGNPI': ['TreatyIn.EGNPI'],
    'T_TREATY_RETENTION': ['TreatyIn.Retention'],
    'T_TREATY_LIMIT_SUMMARY': ['TreatyIn.LimitSummaryList',
                               'TreatyIn.LimitShareSummaryList',
                               'TreatyIn.LimitFacShareSummaryList',
                               'TreatyIn.MDPSummaryList'],
    'T_TREATY_INSTALLMENT': ['TreatyIn.Installment'],
    # bentuknya HANYA terbaca lewat salinan — lihat struktur-treatyin-lama.md §5.2
    'T_TREATY_INSTALLMENT_ITEM': ['TreatyIn.ValueDifference.Installment.InstallmentList'],
    'T_TREATY_REPORTING_PERIOD': ['TreatyIn.ReportingPeriodList'],
    'T_TREATY_ACCUMULATION': ['TreatyIn.AccumulationList'],
    'T_TREATY_PORTFOLIO': ['TreatyIn.Portfolio'],
    'T_TREATY_HAZARD_LIMIT': ['TreatyIn'],
    'T_TREATY_TOTAL': ['TreatyIn.TotalShareNetNP', 'TreatyIn.TotalShareGrossNP',
                       'TreatyIn.TotalEgnpiAmountNP', 'TreatyIn.TotalRetentionAmountNP'],
    'T_TREATY_RETRO_SHARE': ['TreatyIn.ShareReins'],
    'T_VIEW_COMMENT': ['TreatyIn.CommentList'],
    'DOCUMENT_TREATY_IN': [],

    'T_TREATY_VALUE_DIFFERENCE': ['TreatyIn.ValueDifference'],
    'T_TREATY_VALUE_BEFORE_PRORATE': ['TreatyIn.ValueBeforeProrate'],

    'T_TREATY_MIG_CORRELATION': [],
    'T_TREATY_MIG_LANDING': [],
    'T_TREATY_MIG_REJECTED': [],
    'T_TREATY_OUTBOUND_ARCHIVE': [],
    'T_TREATY_OFFER': [],
    'T_TREATY_AUTHORITY_LIMIT': [],
}

KELOMPOK = [
    ('akar', 'Akar — tabel TREATY_IN yang sudah ada'),
    ('bersama', 'Dipakai bersama modul lain'),
    ('limit', 'Layer dan rincian proporsional'),
    ('share', 'Bagian dan penyebaran'),
    ('fac', 'Fakultatif'),
    ('kepala', 'Daftar tingkat kontrak'),
    ('edm', 'Addendum dan selisih'),
    ('migrasi', 'Jembatan migrasi'),
    ('lubang', 'Belum ada isinya'),
]
KEPALA = ['NAMA_T', 'PADANAN_DDL', 'KELOMPOK', 'INDUK', 'KUNCI', 'ASAL_PEGA', 'JALUR_PEGA']

if not os.path.isdir(KELUARAN):
    os.makedirs(KELUARAN)


# ------------------------------------------------------------------ penjaga §16
# Oracle menolak pengenal >30 bita, dan penolakannya baru terjadi saat DDL dijalankan.
# Pemeriksaan ini menahannya di sini supaya tidak lolos ke skrip DDL.
def periksa_nama():
    salah, panjang = [], []
    for baris in TABEL:
        for kol, nm in (('NAMA_T', baris[0]), ('PADANAN_DDL', baris[1])):
            if not nm or nm.startswith('—') or ' ' in nm:
                continue          # sel keterangan, bukan pengenal
            b = len(nm.encode('utf-8'))
            if b > 30:
                salah.append((kol, nm, b))
            elif b == 30:
                panjang.append((kol, nm))
    nama = [t[0] for t in TABEL if not t[0].startswith('—')]
    ganda = sorted(set(n for n in nama if nama.count(n) > 1))
    if salah:
        raise SystemExit('§16 DILANGGAR — pengenal >30 bita:\n' +
                         '\n'.join('  %s  %s = %d bita' % s for s in salah))
    if ganda:
        raise SystemExit('NAMA_T ganda: %s' % ', '.join(ganda))
    return panjang


def periksa_tambatan():
    """Tiap tabel wajib punya entri JALUR_PEGA, dan tiap jalur wajib ADA di hasil
    sapuan. Tanpa penjaga ini, ERD bisa menggambar simpul yang tidak pernah ada."""
    nama = [t[0] for t in TABEL if not t[0].startswith('—')]
    hilang = [n for n in nama if n not in JALUR_PEGA]
    if hilang:
        raise SystemExit('tanpa tambatan JALUR_PEGA: %s' % ', '.join(hilang))
    lebih = [n for n in JALUR_PEGA if n not in nama]
    if lebih:
        raise SystemExit('JALUR_PEGA menyebut tabel yang tidak ada: %s' % ', '.join(lebih))

    p = os.path.join(KELUARAN, 'datar-treatyin-lama.csv')
    if not os.path.exists(p):
        print('⚠ datar-treatyin-lama.csv belum ada — jalankan '
              'buat-pohon-treatyin.py lebih dulu; tambatan TIDAK diperiksa')
        return 0
    ada = set()
    with io.open(p, encoding='utf-8-sig') as h:
        for r in csv.DictReader(h):
            ada.add(r['JALUR'])
    palsu = sorted(set(j for v in JALUR_PEGA.values() for j in v) - ada)
    if palsu:
        raise SystemExit('JALUR_PEGA menunjuk simpul yang TIDAK ADA di pohon:\n  ' +
                         '\n  '.join(palsu))
    return sum(len(v) for v in JALUR_PEGA.values())


mepet = periksa_nama()
n_tambat = periksa_tambatan()

# ------------------------------------------------------------------ TSV
p_tsv = os.path.join(KELUARAN, 'peta-nama-tabel-treatyin.tsv')
with io.open(p_tsv, 'w', encoding='utf-8', newline='') as h:
    w = csv.writer(h, delimiter='\t', lineterminator='\n')
    w.writerow(KEPALA)
    for t in TABEL:
        w.writerow(list(t) + [' | '.join(JALUR_PEGA.get(t[0], []))])

# ------------------------------------------------------------------ Markdown
dibuat = [t for t in TABEL if not t[0].startswith('\u2014')]
tidak = [t for t in TABEL if t[0].startswith('\u2014')]
per_kel = collections.Counter(t[2] for t in dibuat)

L = []
A = L.append
A('# Peta nama tabel datar — Treaty In')
A('')
A('**Tanggal:** 23 September 2026 · **Skema:** `TREATY_MASUK`')
A('')
A('> **TURUNAN.** Dibangkitkan [`alat/buat-peta-nama-tabel-treatyin.py`](../alat/'
  'buat-peta-nama-tabel-treatyin.py) dan **tidak pernah disunting tangan**. Bila berbeda '
  'dari definisinya, **alatnya yang salah**.')
A('>')
A('> Isinya diturunkan dari [`struktur-treatyin-lama.md`](struktur-treatyin-lama.md) '
  '— pohon `TreatyIn` hasil sapuan **329 berkas XML**.')
A('')
A('Pendamping mesin-baca: [`peta-nama-tabel-treatyin.tsv`](peta-nama-tabel-treatyin.tsv).')
A('')
A('---')
A('')
A('## 0. Dua nama, dua lapisan')
A('')
A('| Kolom | Lapisan | Tunduk pada |')
A('|---|---|---|')
A('| `NAMA_T` | **tabel datar** — bentuk pipih yang memetakan pohon clipboard Pega '
  'satu lawan satu | konvensi `T_` modul claim-non-prop |')
A('| `PADANAN_DDL` | **objek Oracle** di skema `TREATY_MASUK` | §16 — kata utuh bahasa '
  'Indonesia, ≤30 bita |')
A('')
A('Keduanya berlaku bersamaan dan **tidak saling menggantikan**. Awalan `T_` justru '
  'menandai bahwa barisnya adalah tabel datar, bukan entitas rancangan.')
A('')
A('### Satu pengecualian: kepalanya memakai nama yang sudah ada')
A('')
A('Kepala struktur ini **`TREATY_IN`**, tanpa awalan `T_`, karena ia **tabel Oracle yang '
  'sudah hidup sekarang** — bukan tabel datar yang kita turunkan. Awalan `T_` menandai '
  'sesuatu yang kita buat; menempelkannya pada tabel yang sudah ada akan mengaburkan '
  'justru hal yang paling perlu terlihat, yaitu bahwa tabel itu **bukan** milik kita.')
A('')
A('Seluruh anaknya mengikuti struktur pohon dan memakai `T_…`, dengan kunci tamu '
  '`TREATY_IN_ID`.')
A('')
A('> **20 lawan 99.** `TREATY_IN` yang sudah ada memuat **20 kolom bisnis** — terbaca dari '
  'argumen `POOLDATA.PEGA_TREATY_IN`, bukan ditebak. Halaman `TreatyIn` punya **99 skalar '
  'akar**. Selisih **79** itu tidak punya kolom di mana pun; ia hanya ada di dalam '
  '`M_TREATY_IN.JSONDATA`. Memakai nama tabel yang sudah ada **tidak** berarti memakai '
  'bentuk kolomnya yang sekarang.')
A('')
A('---')
A('')
A('## 1. Ringkasan')
A('')
A('| | |')
A('|---|---:|')
A('| Tabel datar dibuat | **%d** |' % len(dibuat))
A('| Halaman yang **sengaja tidak** menjadi tabel | **%d** |' % len(tidak))
A('| Kedalaman maksimum induk | **4** |')
A('')
A('| Kelompok | Tabel |')
A('|---|---:|')
for k, ket in KELOMPOK:
    if per_kel.get(k):
        A('| %s | %d |' % (ket, per_kel[k]))
A('')
A('---')
A('')
A('## 2. Daftar tabel')
A('')
for k, ket in KELOMPOK:
    baris = [t for t in TABEL if t[2] == k]
    if not baris:
        continue
    A('### %s' % ket)
    A('')
    A('| `NAMA_T` | Padanan DDL | Induk | Kunci tamu | Asal di pohon Pega |')
    A('|---|---|---|---|---|')
    for n, d, _, ind, ku, asal in baris:
        nn = n if n.startswith('\u2014') else '`%s`' % n
        ii = '`%s`' % ind if ind else '—'
        kk = '`%s`' % ku if ku and not ku.startswith('\u2014') else (ku or '—')
        A('| %s | %s | %s | %s | %s |' % (nn, d, ii, kk, asal))
    A('')
A('---')
A('')
A('## 3. Empat keputusan bentuk yang perlu dibaca, bukan disimpulkan')
A('')
A('**(a) Daftar uang berbentuk sama digabung, dengan kolom PERAN.** Kelas '
  '`ASM-FW-GISFW-Data-TreatyInLimitsSpreading` muncul **120 kali** dan hampir seluruhnya '
  'berisi pasangan `Currency` + `Value`. Ia bukan entitas; ia bentuk pembawa satu nilai '
  'uang bermata uang. Membuat 120 tabel akan menyalin satu bentuk seratus dua puluh kali; '
  'membuat satu tabel tanpa kolom peran akan menumpuk besaran yang berbeda. Karena itu: '
  'satu tabel per **induk**, dengan kolom `PERAN` yang memegang nama daftar aslinya.')
A('')
A('**(b) `OLDDATA` tidak menjadi tabel.** Aturan bisnisnya sudah ditetapkan: **data lama '
  'tidak disimpan ulang — ia DISELECT** dari versi yang disesuaikan. 142 simpul di bawah '
  '`TreatyIn.OLDDATA` karena itu tidak melahirkan satu kolom pun.')
A('')
A('**(c) `ActualValue` juga tidak menjadi tabel — dan sebabnya berbeda.** Bukan karena '
  'aturan bisnis, melainkan karena **namanya tidak cocok dengan isinya**: delapan isian di '
  '`TreatyEDMDifferenceShare` menulis hasil pengurangan ke dalamnya. Membuat tabel bernama '
  '"nilai aktual" yang berisi selisih akan mewariskan kekeliruan itu. Adjudikasinya '
  'ditangguhkan ke [`TEMUAN-ADJUSTMENT-DITUNDA.md`](TEMUAN-ADJUSTMENT-DITUNDA.md).')
A('')
A('**(d) Sepuluh skalar batas bahaya diputar menjadi daftar.** `Earthquake`, `FloodJab`, '
  '`FloodNation`, `RSMDLimit`, `MaxCoGroup`, `MaxCoNonGroup` masing-masing berpasangan '
  'dengan kolom mata uangnya sendiri. Bentuk lamanya memaksa satu kolom baru setiap kali '
  'ada jenis bahaya baru; `T_TREATY_HAZARD_LIMIT` membuatnya satu baris baru.')
A('')
A('---')
A('')
A('## 4. Yang TIDAK menjadi tabel, dan sebabnya')
A('')
A('| Halaman / nama | Simpul | Sebab |')
A('|---|---:|---|')
A('| `TreatyIn.OLDDATA{}` | 142 | aturan bisnis — data lama **diselect**, tidak disalin |')
A('| `TreatyIn.ActualValue{}` | 152 | nama tidak cocok dengan isinya (§3c) |')
A('| `TREATYINOFFER` | — | **nol penulis terjangkau** di sistem berjalan |')
A('')
A('Dua yang pertama **nyata dan hidup** di sistem lama; yang tidak dibawa adalah '
  'bentuk penyimpanannya, bukan datanya.')
L.append('')

p_md = os.path.join(KELUARAN, 'PETA-NAMA-TABEL-TREATYIN.md')
with io.open(p_md, 'w', encoding='utf-8') as h:
    h.write('\n'.join(L))

print('tabel datar dibuat      : %d' % len(dibuat))
print('sengaja tidak dibuat    : %d' % len(tidak))
print('§16 pengenal >30 bita   : 0  (diperiksa, bukan diandaikan)')
print('tambatan ke pohon       : %d jalur, seluruhnya ADA di hasil sapuan' % n_tambat)
for kol, nm in mepet:
    print('   [!] %s tepat 30 bita, tanpa sisa: %s' % (kol, nm))
for k, ket in KELOMPOK:
    if per_kel.get(k):
        print('   %-34s %d' % (ket, per_kel[k]))
print('---')
print(p_tsv)
print(p_md)
