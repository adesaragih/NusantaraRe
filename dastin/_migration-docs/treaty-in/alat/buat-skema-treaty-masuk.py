# -*- coding: utf-8 -*-
"""
Membangkitkan dua artefak bentuk data Treaty In + Treaty In Adjustment:

    4-erd-dan-tabel-datar/Diagram-Skema-Tabel-TreatyMasuk.xlsx
    4-erd-dan-tabel-datar/ERD-TREATY-MASUK.html

Bentuknya meniru contoh dari modul claim-non-prop (Diagram-Skema-Tabel-Gabungan.xlsx
dan ERD-GABUNGAN.html). ISINYA tidak disalin dari sana — seluruhnya dari
SPEC-MODEL-DATA.md, ERD.md, STRUKTUR-DATA.md, dan KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md.

KELUARAN INI TURUNAN. Bila salah satu berbeda dari sumbernya, ALATNYA yang salah.
Jangan sunting hasilnya dengan tangan.

Jalankan:  PYTHONIOENCODING=utf-8 python alat/buat-skema-treaty-masuk.py
"""
import os, openpyxl
from openpyxl.styles import Font, PatternFill, Alignment
from openpyxl.utils import get_column_letter

KELUARAN = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', '4-erd-dan-tabel-datar')

# ── warna lajur ────────────────────────────────────────────────────────────────
INTI   = '1F4E79'   # tulang punggung
ANAK   = 'C55A11'   # anak versi, dipakai kedua cabang
CABANG = '2E75B6'   # layer dan kedua cabangnya
SPREAD = '31859C'   # potongan dan penyebaran
SEKAT  = 'C00000'   # Adjustment — melintasi sekat
ACUAN  = '7030A0'   # tabel acuan bersama
LUAR   = '808080'   # di luar gelombang ini

KREM   = 'FFF2CC'   # baris PK
PUCAT  = 'EDF3F9'   # baris FK dan catatan
PITA_N = 'FCE4D6'   # kotak catatan bergaris putus

# ── data: satu sumber untuk kedua keluaran ─────────────────────────────────────
# ('tabel', kardinalitas, lajur_dalam, [baris...])
# baris: ('pk'|'fk'|'x'|'n', teks)   x = atribut pembeda cabang, n = catatan

LAJUR = [
 (INTI, 'INTI — TULANG PUNGGUNG · kontrak dan versinya', None, [

  ('KONTRAK', 'akar', 0, [
    ('pk', 'PK   ID_KONTRAK'),
    ('fk', 'FK   ID_KONTRAK_DISALIN_DARI  →  KONTRAK.ID_KONTRAK      [hapus: putus]'),
    ('fk', 'FK   ID_CEDANT               →  CEDANT [luar]            [hapus: tolak]'),
    ('fk', 'FK   ID_ASAL_BISNIS          →  ASAL_BISNIS [luar]       [hapus: tolak]'),
    ('x',  '▸    SIFAT_PROPORSI   PEMBEDA CABANG — menentukan anak mana yang berlaku'),
    ('n',  'kunci alami: cedant + asal bisnis + periode + sifat proporsi'),
    ('n',  'kunci alami MEMPERINGATKAN, tidak melarang — ADR-0040 §2'),
    ('n',  'hanya lapisan yang tidak pernah berubah sepanjang hidup kontrak'),
  ]),

  ('VERSI_KONTRAK', '1:N', 1, [
    ('pk', 'PK   ID_VERSI_KONTRAK'),
    ('fk', 'FK   ID_KONTRAK              →  KONTRAK.ID_KONTRAK       [hapus: tolak]'),
    ('fk', 'FK   ID_VERSI_KONTRAK_DASAR  →  VERSI_KONTRAK            [hapus: tolak]'),
    ('fk', 'FK   ID_REASURADUR_PEMIMPIN  →  REASURADUR [luar]        [hapus: tolak]'),
    ('fk', 'FK   KELAS_BISNIS_KONTRAK    →  KELAS_BISNIS             [hapus: tolak]'),
    ('n',  'kunci alami: NOMOR_URUT_VERSI di dalam kontraknya'),
    ('n',  '⭐ terisinya ID_VERSI_KONTRAK_DASAR = versi ini hasil PENYESUAIAN'),
    ('n',  '   itu satu-satunya pembeda — tidak ada kolom jenis, tidak ada entitas terpisah'),
    ('n',  'penunjuknya TUNGGAL: sisi "baru" adalah versi induk barisnya sendiri'),
    ('n',  'memikul seluruh isi kontrak + keadaan persetujuannya sendiri'),
  ]),
 ]),

 (ANAK, 'ANAK VERSI — dipakai KEDUA cabang · seluruhnya menggantung pada VERSI_KONTRAK', '⇦  induknya VERSI_KONTRAK di lajur INTI — tidak digambar dua kali', [
  ('MATA_UANG_KONTRAK', '1:N', 2, [
    ('pk','PK   ID_MATA_UANG_KONTRAK'),
    ('fk','FK   ID_VERSI_KONTRAK  →  VERSI_KONTRAK        [hapus: ikut hapus]'),
    ('fk','FK   KODE_MATA_UANG    →  MATA_UANG            [hapus: tolak]'),
    ('n','kunci alami: kode mata uang di dalam versinya'),
  ]),
  ('RETENSI_CEDANT', '1:N', 2, [
    ('pk','PK   ID_RETENSI_CEDANT'),
    ('fk','FK   ID_VERSI_KONTRAK    →  VERSI_KONTRAK      [hapus: ikut hapus]'),
    ('fk','FK   ID_KELOMPOK_TREATY  →  KELOMPOK_TREATY    [hapus: tolak]'),
    ('n','kunci alami: kelompok treaty + mata uang'),
  ]),
  ('EGNPI', '1:N', 2, [
    ('pk','PK   ID_EGNPI'),
    ('fk','FK   ID_VERSI_KONTRAK    →  VERSI_KONTRAK      [hapus: ikut hapus]'),
    ('fk','FK   ID_KELOMPOK_TREATY  →  KELOMPOK_TREATY    [hapus: tolak]'),
    ('n','kunci alami: kelompok treaty + mata uang'),
  ]),
  ('PORTOFOLIO', '1:N', 2, [
    ('pk','PK   ID_PORTOFOLIO'),
    ('fk','FK   ID_VERSI_KONTRAK  →  VERSI_KONTRAK        [hapus: ikut hapus]'),
  ]),
  ('PERIODE_PELAPORAN', '1:N', 2, [
    ('pk','PK   ID_PERIODE_PELAPORAN'),
    ('fk','FK   ID_VERSI_KONTRAK  →  VERSI_KONTRAK        [hapus: ikut hapus]'),
    ('n','menyimpan tanggal jatuh tempo — TIDAK menyimpan bendera terlambat'),
  ]),
  ('PERIODE_AKUMULASI', '1:N', 2, [
    ('pk','PK   ID_PERIODE_AKUMULASI'),
    ('fk','FK   ID_VERSI_KONTRAK  →  VERSI_KONTRAK        [hapus: ikut hapus]'),
  ]),
  ('TERMIN', '1:N', 2, [
    ('pk','PK   ID_TERMIN'),
    ('fk','FK   ID_VERSI_KONTRAK  →  VERSI_KONTRAK        [hapus: ikut hapus]'),
    ('n','kunci alami: nomor termin di dalam versinya'),
  ]),
  ('SKALA_KOASURANSI', '1:N', 2, [
    ('pk','PK   ID_SKALA_KOASURANSI'),
    ('fk','FK   ID_VERSI_KONTRAK  →  VERSI_KONTRAK        [hapus: ikut hapus]'),
    ('n','⚠ di sistem lama tersimpan sebagai SKALAR — ia daftar'),
  ]),
  ('BATAS_PER_BAHAYA', '1:N', 2, [
    ('pk','PK   ID_BATAS_PER_BAHAYA'),
    ('fk','FK   ID_VERSI_KONTRAK  →  VERSI_KONTRAK        [hapus: ikut hapus]'),
    ('fk','FK   ID_BAHAYA         →  BAHAYA               [hapus: tolak]'),
    ('n','⭐ menggantikan 8 kolom bernama bahaya — bahayanya NILAI, bukan nama kolom'),
  ]),
  ('CATATAN_PERSETUJUAN', '1:N', 2, [
    ('pk','PK   ID_CATATAN_PERSETUJUAN'),
    ('fk','FK   ID_VERSI_KONTRAK  →  VERSI_KONTRAK        [hapus: TOLAK]'),
    ('n','⚠ TOLAK disengaja — jejak yang dapat dihapus bersama bendanya bukan jejak'),
  ]),
  ('DOKUMEN_KONTRAK', '1:N', 2, [
    ('pk','PK   ID_DOKUMEN_KONTRAK'),
    ('fk','FK   ID_VERSI_KONTRAK  →  VERSI_KONTRAK        [hapus: ikut hapus]'),
    ('fk','FK   ID_DOKUMEN        →  DOKUMEN [luar]       [hapus: tolak]'),
    ('n','dokumen DIRUJUK, tidak dimiliki — ADR-0027'),
  ]),
  ('JEJAK_PERUBAHAN', '1:N', 2, [
    ('pk','PK   ID_JEJAK_PERUBAHAN'),
    ('fk','FK   ID_VERSI_KONTRAK  →  VERSI_KONTRAK        [hapus: TOLAK]'),
    ('n','⚠ TOLAK disengaja — sama dengan CATATAN_PERSETUJUAN'),
    ('n','tidak ada padanannya di sistem lama — ADR-0045'),
  ]),
 ]),

 (CABANG, 'CABANG — satu-satunya tempat kedua sisi berpisah', None, [
  ('LAYER', '1:N', 2, [
    ('pk','PK   ID_LAYER'),
    ('fk','FK   ID_VERSI_KONTRAK  →  VERSI_KONTRAK        [hapus: ikut hapus]'),
    ('fk','FK   KELAS_BISNIS      →  KELAS_BISNIS         [hapus: tolak]'),
    ('n','kunci alami: nomor layer + bagian layer di dalam versinya'),
    ('n','ada di KEDUA cabang — layer (non-prop) atau kelompok limit (prop)'),
  ]),
  ('DETAIL_PROPORSIONAL', '1:N', 3, [
    ('pk','PK   ID_DETAIL_PROPORSIONAL'),
    ('fk','FK   ID_LAYER            →  LAYER              [hapus: ikut hapus]'),
    ('fk','FK   ID_KELOMPOK_TREATY  →  KELOMPOK_TREATY    [hapus: tolak]'),
    ('x','▸    hanya bila KONTRAK.SIFAT_PROPORSI = PROPORSIONAL   — INV-32'),
    ('n','kunci alami: kelompok treaty di dalam layernya'),
  ]),
  # --- ditambahkan 24 September 2026: dua entitas yang lahir sesudah berkas ini ditulis ---
  ('PEMULIHAN_LIMIT', '1:N', 3, [
    ('pk','PK   ID_PEMULIHAN_LIMIT'),
    ('fk','FK   ID_LAYER  →  LAYER                        [hapus: ikut hapus]'),
    ('n','satu baris per ketentuan pemulihan, bukan dua skalar pada layer'),
    ('n','kunci alami BELUM BERNOMOR — SPEC-INVARIAN.md 7.3'),
    ('n','kemampuan BARU P-59: sistem lama menulis kedua persennya tetap "100" di kode'),
  ]),
  ('PERISTIWA_KONTRAK', '1:N', 2, [
    ('pk','PK   ID_PERISTIWA_KONTRAK'),
    ('fk','FK   ID_VERSI_KONTRAK  →  VERSI_KONTRAK       [hapus: ikut hapus]'),
    ('n','TIDAK punya kunci alami, dan itu keputusan — peristiwa yang sama'),
    ('n','dapat terjadi dua kali pada versi yang sama (10.21a)'),
  ]),
  ('BAGIAN', '1:o1', 3, [
    ('pk','PK   ID_BAGIAN'),
    ('fk','FK   ID_LAYER  →  LAYER                        [hapus: ikut hapus]'),
    ('x','▸    hanya bila KONTRAK.SIFAT_PROPORSI = NON_PROPORSIONAL — INV-33'),
    ('n','paling banyak SATU per layer, bukan banyak'),
    ('n','bagian NuRe pada cabang proporsional adalah ATRIBUT, bukan baris'),
  ]),
 ]),

 (SPREAD, 'POTONGAN DAN PENYEBARAN — satu konsep, DUA pelekatan', '⇦  induknya BAGIAN atau DETAIL_PROPORSIONAL — tepat satu terisi tiap baris', [
  ('POTONGAN', '1:N', 4, [
    ('pk','PK   ID_POTONGAN'),
    ('fk','FK   ID_INDUK_POTONGAN  →  BAGIAN  ATAU  DETAIL_PROPORSIONAL  [ikut hapus]'),
    ('fk','FK   ID_JENIS_POTONGAN  →  JENIS_POTONGAN     [hapus: tolak]'),
    ('n','⚠ TEPAT SATU induk terisi — invarian, bukan pilihan bebas'),
    ('n','aturannya TIDAK BOLEH tertulis dua kali — INV-63'),
    ('n','nilai potongannya TURUNAN, tidak disimpan'),
  ]),
  ('PENYEBARAN', '1:N', 4, [
    ('pk','PK   ID_PENYEBARAN'),
    ('fk','FK   ID_INDUK_PENYEBARAN  →  BAGIAN  ATAU  DETAIL_PROPORSIONAL [ikut hapus]'),
    ('fk','FK   ID_JENIS_REASURANSI  →  JENIS_REASURANSI  [hapus: tolak]'),
    ('n','kunci alami: jenis reasuransi di dalam induknya'),
  ]),
  ('RINCIAN_PENYEBARAN', '1:N', 5, [
    ('pk','PK   ID_RINCIAN_PENYEBARAN'),
    ('fk','FK   ID_PENYEBARAN  →  PENYEBARAN             [hapus: ikut hapus]'),
  ]),
  ('NILAI_PENYEBARAN', '1:N', 6, [
    ('pk','PK   ID_NILAI_PENYEBARAN'),
    ('fk','FK   ID_RINCIAN_PENYEBARAN  →  RINCIAN_PENYEBARAN [hapus: ikut hapus]'),
    ('fk','FK   KODE_MATA_UANG         →  MATA_UANG          [hapus: tolak]'),
    ('n','⭐ berbaris per mata uang — menggantikan kolom kembar Rp/Usd'),
  ]),
 ]),

 (SEKAT, 'TREATY IN ADJUSTMENT — MELINTASI SEKAT · dua relasi, tidak lebih', '⇦  seluruh sambungan antarmodul ada di dua relasi di bawah ini', [
  ('NILAI_SELISIH', '1:N', 2, [
    ('pk','PK   ID_NILAI_SELISIH'),
    ('fk','FK   ID_VERSI_KONTRAK  →  VERSI_KONTRAK   ⟦SEKAT⟧  [hapus: ikut hapus]'),
    ('fk','FK   ID_BESARAN        →  BESARAN_DAPAT_DISESUAIKAN ⟦SEKAT⟧ [tolak]'),
    ('n','satu baris = SATU besaran yang berubah pada sebuah versi'),
    ('n','⭐ sisi "lama" TIDAK disimpan — dibaca lewat VERSI_KONTRAK_DASAR induknya'),
    ('n','bentuk SEMPIT — seluruh besaran yang dapat disesuaikan adalah uang atau persen'),
    ('n','kunci alami: besaran + kunci padanan baris'),
  ]),
  ('BESARAN_DAPAT_DISESUAIKAN', 'acuan', 1, [
    ('pk','PK   ID_BESARAN'),
    ('x','▸    SATUAN_BESARAN   UANG | PERSENTASE — mengikat bentuk baris selisih'),
    ('n','⚠ diisi 31 sebagai AWAL, BUKAN sebagai batas'),
    ('n','   ketiadaan sebuah besaran bukan larangan sampai dikonfirmasi bisnis'),
    ('n','dipakai BERSAMA kedua modul — tidak boleh diubah sepihak'),
  ]),
 ]),

 (ACUAN, 'TABEL ACUAN BERSAMA — wajib tabel, bukan CHECK · ADR-0038, INV-62', None, [
  ('MATA_UANG',        'acuan', 1, [('pk','PK   KODE_MATA_UANG'), ('n','dirujuk MATA_UANG_KONTRAK dan NILAI_PENYEBARAN')]),
  ('JENIS_POTONGAN',   'acuan', 1, [('pk','PK   ID_JENIS_POTONGAN'), ('n','⭐ berasal dari Comment di sistem lama — itu NAMA potongannya')]),
  ('JENIS_REASURANSI', 'acuan', 1, [('pk','PK   ID_JENIS_REASURANSI')]),
  ('BAHAYA',           'acuan', 1, [('pk','PK   ID_BAHAYA'), ('n','gempa · banjir Jabodetabek · banjir nasional · RSMD(?)'), ('n','⚠ kepanjangan RSMD belum dipastikan — Uji Y')]),
  ('KELOMPOK_TREATY',  'acuan', 1, [('pk','PK   ID_KELOMPOK_TREATY')]),
  ('KELAS_BISNIS',     'acuan', 1, [('pk','PK   ID_KELAS_BISNIS')]),
 ]),

 (LUAR, 'DI LUAR GELOMBANG INI — dimodelkan, tidak dibangun sekarang', None, [
  ('RETRO_KELUAR', '1:N  [G2]', 2, [
    ('pk','PK   ID_RETRO_KELUAR'),
    ('fk','FK   ID_VERSI_KONTRAK  →  VERSI_KONTRAK        [hapus: ikut hapus]'),
    ('n','⚠ arah KELUAR — bukan fakultatif masuk; ADR-0020 TIDAK berlaku'),
  ]),
  ('PENCAPAIAN', '1:N  [G3]', 1, [
    ('pk','PK   ID_PENCAPAIAN'),
    ('fk','FK   ID_KONTRAK  →  KONTRAK                    [hapus: tolak]'),
  ]),
 ]),
]

# batas yang belum diketahui — bukan entitas
BATAS = ('BATAS PENERBITAN KE LUAR — BELUM DIKETAHUI ISINYA',
         'TREATYINOFFER tidak punya satu pun penulis yang terjangkau: penulisnya dua, '
         'pemanggilnya satu, dan kedua pemanggil pemanggilnya tertutup (satu ber-blok mati, '
         'satu di balik tombol bersyarat NEVER). Ditutup Uji AB + satu pertanyaan bisnis.')

# ── daftar relasi, untuk sheet dan tabel HTML ──────────────────────────────────
RELASI = [
 ('inti','KONTRAK','VERSI_KONTRAK','ID_KONTRAK','1:N','tolak','kontrak yang punya versi tidak boleh hilang'),
 ('inti','KONTRAK','KONTRAK','ID_KONTRAK_DISALIN_DARI','1:o N','putus','pecahan pertama OLDID — ADR-0040 §1'),
 ('inti','VERSI_KONTRAK','VERSI_KONTRAK','ID_VERSI_KONTRAK_DASAR','1:o N','tolak','⭐ dasar penyesuaian — inti sambungan Adjustment'),
 ('anak','VERSI_KONTRAK','MATA_UANG_KONTRAK','ID_VERSI_KONTRAK','1:N','ikut hapus',''),
 ('anak','VERSI_KONTRAK','RETENSI_CEDANT','ID_VERSI_KONTRAK','1:o N','ikut hapus',''),
 ('anak','VERSI_KONTRAK','EGNPI','ID_VERSI_KONTRAK','1:o N','ikut hapus',''),
 ('anak','VERSI_KONTRAK','PORTOFOLIO','ID_VERSI_KONTRAK','1:o N','ikut hapus',''),
 ('anak','VERSI_KONTRAK','PERIODE_PELAPORAN','ID_VERSI_KONTRAK','1:o N','ikut hapus',''),
 ('anak','VERSI_KONTRAK','PERIODE_AKUMULASI','ID_VERSI_KONTRAK','1:o N','ikut hapus',''),
 ('anak','VERSI_KONTRAK','TERMIN','ID_VERSI_KONTRAK','1:o N','ikut hapus',''),
 ('anak','VERSI_KONTRAK','SKALA_KOASURANSI','ID_VERSI_KONTRAK','1:o N','ikut hapus','daftar, bukan skalar'),
 ('anak','VERSI_KONTRAK','BATAS_PER_BAHAYA','ID_VERSI_KONTRAK','1:o N','ikut hapus','menggantikan 8 kolom bernama bahaya'),
 ('anak','VERSI_KONTRAK','CATATAN_PERSETUJUAN','ID_VERSI_KONTRAK','1:N','TOLAK','⚠ jejak tidak ikut terhapus'),
 ('anak','VERSI_KONTRAK','DOKUMEN_KONTRAK','ID_VERSI_KONTRAK','1:o N','ikut hapus','dokumen dirujuk, tidak dimiliki'),
 ('anak','VERSI_KONTRAK','JEJAK_PERUBAHAN','ID_VERSI_KONTRAK','1:N','TOLAK','⚠ jejak tidak ikut terhapus'),
 ('cabang','VERSI_KONTRAK','LAYER','ID_VERSI_KONTRAK','1:N','ikut hapus',''),
 ('cabang','LAYER','DETAIL_PROPORSIONAL','ID_LAYER','1:N','ikut hapus','hanya bila SIFAT_PROPORSI = PROPORSIONAL — INV-32'),
 ('cabang','LAYER','BAGIAN','ID_LAYER','1:o1','ikut hapus','hanya bila SIFAT_PROPORSI = NON_PROPORSIONAL — INV-33'),
 ('sebar','BAGIAN','POTONGAN','ID_INDUK_POTONGAN','1:o N','ikut hapus','induk pilihan — tepat satu terisi'),
 ('sebar','DETAIL_PROPORSIONAL','POTONGAN','ID_INDUK_POTONGAN','1:o N','ikut hapus','induk pilihan — tepat satu terisi'),
 ('sebar','BAGIAN','PENYEBARAN','ID_INDUK_PENYEBARAN','1:o N','ikut hapus','induk pilihan — tepat satu terisi'),
 ('sebar','DETAIL_PROPORSIONAL','PENYEBARAN','ID_INDUK_PENYEBARAN','1:o N','ikut hapus','induk pilihan — tepat satu terisi'),
 ('sebar','PENYEBARAN','RINCIAN_PENYEBARAN','ID_PENYEBARAN','1:N','ikut hapus',''),
 ('sebar','RINCIAN_PENYEBARAN','NILAI_PENYEBARAN','ID_RINCIAN_PENYEBARAN','1:N','ikut hapus','berbaris per mata uang'),
 ('SEKAT','VERSI_KONTRAK','NILAI_SELISIH','ID_VERSI_KONTRAK','1:o N','ikut hapus','⟦SEKAT⟧ Treaty In → Adjustment'),
 ('SEKAT','BESARAN_DAPAT_DISESUAIKAN','NILAI_SELISIH','ID_BESARAN','1:N','tolak','⟦SEKAT⟧ acuan bersama → Adjustment'),
 ('acuan','MATA_UANG','MATA_UANG_KONTRAK','KODE_MATA_UANG','1:N','tolak',''),
 ('acuan','MATA_UANG','NILAI_PENYEBARAN','KODE_MATA_UANG','1:N','tolak',''),
 ('acuan','JENIS_POTONGAN','POTONGAN','ID_JENIS_POTONGAN','1:N','tolak',''),
 ('acuan','JENIS_REASURANSI','PENYEBARAN','ID_JENIS_REASURANSI','1:N','tolak',''),
 ('acuan','BAHAYA','BATAS_PER_BAHAYA','ID_BAHAYA','1:N','tolak',''),
 ('acuan','KELOMPOK_TREATY','DETAIL_PROPORSIONAL','ID_KELOMPOK_TREATY','1:N','tolak',''),
 ('acuan','KELOMPOK_TREATY','RETENSI_CEDANT','ID_KELOMPOK_TREATY','1:N','tolak',''),
 ('acuan','KELOMPOK_TREATY','EGNPI','ID_KELOMPOK_TREATY','1:N','tolak',''),
 ('acuan','KELAS_BISNIS','LAYER','KELAS_BISNIS','1:N','tolak',''),
 ('acuan','KELAS_BISNIS','VERSI_KONTRAK','KELAS_BISNIS_KONTRAK','1:N','tolak',''),
 ('luar','CEDANT [luar]','KONTRAK','ID_CEDANT','1:N','tolak','nama TIDAK disalin — ADR-0023'),
 ('luar','ASAL_BISNIS [luar]','KONTRAK','ID_ASAL_BISNIS','1:N','tolak','nama TIDAK disalin'),
 ('luar','REASURADUR [luar]','VERSI_KONTRAK','ID_REASURADUR_PEMIMPIN','1:N','tolak','nama TIDAK disalin'),
 ('luar','DOKUMEN [luar]','DOKUMEN_KONTRAK','ID_DOKUMEN','1:N','tolak','dirujuk, tidak dimiliki'),
 ('cabang','LAYER','PEMULIHAN_LIMIT','ID_LAYER','1:o N','ikut hapus','ditambahkan 24 Sep 2026'),
 ('anak','VERSI_KONTRAK','PERISTIWA_KONTRAK','ID_VERSI_KONTRAK','1:o N','ikut hapus','ditambahkan 24 Sep 2026'),
 ('G2/G3','VERSI_KONTRAK','RETRO_KELUAR','ID_VERSI_KONTRAK','1:o N','ikut hapus','[G2] arah KELUAR'),
 ('G2/G3','KONTRAK','PENCAPAIAN','ID_KONTRAK','1:o N','tolak','[G3]'),
]

SAMBUNGAN = [
 ('1','relasi lintas-sekat','VERSI_KONTRAK','NILAI_SELISIH','ID_VERSI_KONTRAK','1:o N',
  'Treaty In','Adjustment','selisih menggantung pada versi yang dihasilkan penyesuaian'),
 ('2','relasi lintas-sekat','BESARAN_DAPAT_DISESUAIKAN','NILAI_SELISIH','ID_BESARAN','1:N',
  'bersama','Adjustment','kosakata besaran TERTUTUP — hapus ditolak selama dipakai'),
 ('3','entitas bersama','BESARAN_DAPAT_DISESUAIKAN','—','—','acuan',
  'bersama','bersama','diisi 31 sebagai AWAL, bukan sebagai batas'),
 ('4','penunjuk, BUKAN relasi baru','VERSI_KONTRAK','VERSI_KONTRAK','ID_VERSI_KONTRAK_DASAR','1:o N',
  'Treaty In','Treaty In','sisi "lama" sebuah selisih — di-SELECT, tidak disalin'),
 ('5','entitas yang TIDAK dibuat','PENYESUAIAN','—','—','—',
  '—','—','⭐ dibuang: 19 properti kelas addendum lama seluruhnya sudah punya rumah'),
]


def _rows_of(t):
    return t[3]


def tulis_xlsx():
    wb = openpyxl.Workbook()

    # ══ sheet 1: Peta Skema ══════════════════════════════════════════════════
    ws = wb.active; ws.title = 'Peta Skema'
    ws.sheet_view.showGridLines = False
    for c in range(1, 86):
        ws.column_dimensions[get_column_letter(c)].width = 2.3

    ws['B1'] = 'PETA SKEMA — TREATY MASUK · Treaty In + Treaty In Adjustment'
    ws['B1'].font = Font(bold=True, size=14)
    ws.merge_cells('B1:BY1')
    n_tabel = sum(len(l[3]) for l in LAJUR)
    ws['B2'] = ('%d tabel · %d relasi · 2 relasi melintasi sekat antarmodul · '
                'skema TREATY_MASUK · hanya PK dan FK' % (n_tabel, len(RELASI)))
    ws['B2'].font = Font(size=9); ws.merge_cells('B2:BY2')

    r = 4
    for warna, judul, nota, tabel in LAJUR:
        col = 17                                   # kolom Q
        cl = get_column_letter(col)
        ws.cell(r, col, '  ' + judul).font = Font(bold=True, size=11, color='FFFFFF')
        ws.cell(r, col).fill = PatternFill('solid', start_color=warna)
        ws.merge_cells('%s%d:%s%d' % (cl, r, 'BJ', r)); r += 2

        if nota:
            ws.cell(r, col, '  ' + nota).font = Font(bold=True, size=9, color=warna)
            ws.cell(r, col).fill = PatternFill('solid', start_color=PITA_N)
            ws.merge_cells('%s%d:%s%d' % (cl, r, 'BJ', r)); r += 2

        for nama, kard, dalam, baris in tabel:
            c0 = 17 + dalam * 3
            cl0 = get_column_letter(c0)
            ws.cell(r, c0, '%s        %s' % (nama, kard)).font = Font(bold=True, size=9, color='FFFFFF')
            ws.cell(r, c0).fill = PatternFill('solid', start_color=warna)
            ws.merge_cells('%s%d:%s%d' % (cl0, r, 'BJ', r)); r += 1
            for jenis, teks in baris:
                sel = ws.cell(r, c0, '  ' + teks)
                if jenis == 'pk':
                    sel.font = Font(bold=True, size=8); sel.fill = PatternFill('solid', start_color=KREM)
                elif jenis == 'fk':
                    sel.font = Font(bold=True, size=8); sel.fill = PatternFill('solid', start_color=PUCAT)
                elif jenis == 'x':
                    sel.font = Font(bold=True, size=8, color=warna); sel.fill = PatternFill('solid', start_color=PITA_N)
                else:
                    sel.font = Font(size=8); sel.fill = PatternFill('solid', start_color=PUCAT)
                ws.merge_cells('%s%d:%s%d' % (cl0, r, 'BJ', r)); r += 1
            r += 1
        r += 1

    # kotak batas yang belum diketahui
    ws.cell(r, 17, '  ⌧  ' + BATAS[0]).font = Font(bold=True, size=10, color='C00000')
    ws.cell(r, 17).fill = PatternFill('solid', start_color=PITA_N)
    ws.merge_cells('Q%d:BJ%d' % (r, r)); r += 1
    ws.cell(r, 17, '  ' + BATAS[1]).font = Font(size=8, color='C00000')
    ws.cell(r, 17).fill = PatternFill('solid', start_color=PITA_N)
    ws.merge_cells('Q%d:BJ%d' % (r, r))

    # ══ sheet 2: Sambungan Antarmodul ════════════════════════════════════════
    ws2 = wb.create_sheet('Sambungan Antarmodul')
    ws2['A1'] = 'SAMBUNGAN ANTARMODUL — seluruh titik temu Treaty In ⟷ Treaty In Adjustment'
    ws2['A1'].font = Font(bold=True, size=12)
    ws2['A2'] = ('Dua relasi, satu entitas bersama, satu penunjuk, dan satu entitas yang '
                 'sengaja TIDAK dibuat. Tidak ada yang lain.')
    ws2['A2'].font = Font(size=9)
    hdr = ['#','Jenis','Dari','Ke','Kunci','Kard.','Sisi Treaty In','Sisi Adjustment','Keterangan']
    for i, h in enumerate(hdr, 1):
        c = ws2.cell(4, i, h); c.font = Font(bold=True, size=9, color='FFFFFF')
        c.fill = PatternFill('solid', start_color=INTI)
    for j, row in enumerate(SAMBUNGAN, 5):
        for i, v in enumerate(row, 1):
            c = ws2.cell(j, i, v); c.font = Font(size=9, color=SEKAT if row[1].startswith('relasi') else '000000')
            c.alignment = Alignment(vertical='top', wrap_text=(i == 9))
    for i, w in enumerate([5,24,28,22,26,9,14,16,62], 1):
        ws2.column_dimensions[get_column_letter(i)].width = w

    # ══ sheet 3: Daftar Relasi ═══════════════════════════════════════════════
    ws3 = wb.create_sheet('Daftar Relasi')
    ws3['A1'] = 'DAFTAR RELASI — %d relasi · setiap baris menyatakan kardinalitas, kunci, dan perilaku hapus' % len(RELASI)
    ws3['A1'].font = Font(bold=True, size=12)
    hdr = ['#','Lajur','Induk','Anak','Kunci tamu','Kard.','Perilaku hapus','Catatan']
    for i, h in enumerate(hdr, 1):
        c = ws3.cell(3, i, h); c.font = Font(bold=True, size=9, color='FFFFFF')
        c.fill = PatternFill('solid', start_color=INTI)
    for j, row in enumerate(RELASI, 4):
        for i, v in enumerate((j-3,) + row, 1):
            c = ws3.cell(j, i, v); c.font = Font(size=9, color=SEKAT if row[0] == 'SEKAT' else '000000')
            c.alignment = Alignment(vertical='top', wrap_text=(i == 8))
    for i, w in enumerate([5,9,28,28,26,8,14,56], 1):
        ws3.column_dimensions[get_column_letter(i)].width = w

    # ══ sheet 4: Peta Sumber ═════════════════════════════════════════════════
    ws4 = wb.create_sheet('Peta Sumber')
    ws4['A1'] = 'PETA SUMBER — tiap tabel, lajurnya, induknya, dan baris PK/FK-nya'
    ws4['A1'].font = Font(bold=True, size=12)
    hdr = ['#','Tabel','Lajur','Kard.','Milik','Baris PK / FK']
    for i, h in enumerate(hdr, 1):
        c = ws4.cell(3, i, h); c.font = Font(bold=True, size=9, color='FFFFFF')
        c.fill = PatternFill('solid', start_color=INTI)
    j = 4; n = 0
    for warna, judul, nota, tabel in LAJUR:
        lajur = judul.split('—')[0].strip()
        milik = ('Adjustment' if warna == SEKAT else
                 'bersama' if warna == ACUAN else
                 'di luar gelombang' if warna == LUAR else 'Treaty In')
        for nama, kard, dalam, baris in tabel:
            n += 1
            kunci = ' · '.join(t for jn, t in baris if jn in ('pk', 'fk'))
            for i, v in enumerate((n, nama, lajur, kard, milik, kunci), 1):
                c = ws4.cell(j, i, v); c.font = Font(size=9)
                c.alignment = Alignment(vertical='top', wrap_text=(i == 6))
            j += 1
    for i, w in enumerate([5,30,22,10,18,96], 1):
        ws4.column_dimensions[get_column_letter(i)].width = w

    p = os.path.join(KELUARAN, 'Diagram-Skema-Tabel-TreatyMasuk.xlsx')
    wb.save(p)
    return p, n_tabel


def tulis_html(n_tabel):
    X0, W, BARIS, PITA = 298.0, 800.0, 20.0, 20.0
    y = 86.0
    lines, boxes, pitas, notas = [], [], [], []
    posisi = {}                                     # nama -> (x, y_tengah, y_atas, y_bawah)

    for warna, judul, nota, tabel in LAJUR:
        pitas.append((X0, y, W, '#' + warna, judul)); y += PITA + 20
        if nota:
            notas.append((X0, y, W, nota)); y += PITA + 20
        for nama, kard, dalam, baris in tabel:
            x = X0 + dalam * 17
            w = W - dalam * 17
            h = PITA + BARIS * len(baris)
            boxes.append((x, y, w, h, '#' + warna, '%s      %s' % (nama, kard), baris))
            posisi[nama] = (x, y + h / 2, y, y + h)
            y += h + 20
        y += 16

    y_batas = y + 10
    tinggi = y_batas + 70

    def esc(s):
        return (s.replace('&', '&amp;').replace('<', '&lt;').replace('>', '&gt;'))

    out = []
    out.append('<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg" '
               'font-family="Arial, Helvetica, sans-serif">' % (1180, int(tinggi), 1180, int(tinggi)))
    out.append('<style>.kp{font-size:11px;font-weight:700;fill:#fff;text-anchor:middle}'
               '.br{font-family:ui-monospace,Consolas,Menlo,monospace;font-size:10.5px}'
               '.pita{font-size:11.5px;font-weight:700;fill:#fff}'
               '.pnb{font-size:10px;font-weight:700;fill:#C55A11}'
               '.bts{font-size:10px;font-weight:700;fill:#C00000}</style>')

    # garis induk → anak
    for lajur_nama, induk, anak, kunci, kard, hapus, cat in RELASI:
        if induk not in posisi or anak not in posisi:
            continue
        if lajur_nama in ('acuan', 'luar'):
            continue
        xi, yi, _, _ = posisi[induk]
        xa, ya, _, _ = posisi[anak]
        if xa <= xi:
            continue
        bengkok = xa - 9
        warna = '#C00000' if lajur_nama == 'SEKAT' else '#1F4E79'
        lebar = '2.6' if lajur_nama == 'SEKAT' else '2.0'
        dash = ' stroke-dasharray="6,3"' if lajur_nama == 'SEKAT' else ''
        out.append('<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="%s"%s/>'
                   % (bengkok, yi, bengkok, ya, warna, lebar, dash))
        out.append('<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="%s"%s/>'
                   % (bengkok, ya, xa, ya, warna, lebar, dash))

    for x, yy, w, warna, teks in pitas:
        out.append('<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s"/>' % (x, yy, w, PITA, warna))
        out.append('<text class="pita" x="%.1f" y="%.1f">%s</text>' % (x + 8, yy + 14.4, esc(teks)))
    for x, yy, w, teks in notas:
        out.append('<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#FCE4D6" stroke="#C55A11" '
                   'stroke-width="1.6" stroke-dasharray="5,3"/>' % (x, yy, w, PITA))
        out.append('<text class="pnb" x="%.1f" y="%.1f">%s</text>' % (x + 8, yy + 14.4, esc(teks)))

    for x, yy, w, h, warna, judul, baris in boxes:
        out.append('<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#EDF3F9" stroke="%s" stroke-width="1.8"/>'
                   % (x, yy, w, h, warna))
        out.append('<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s"/>' % (x, yy, w, PITA, warna))
        out.append('<text class="kp" x="%.1f" y="%.1f">%s</text>' % (x + w / 2, yy + 14.4, esc(judul)))
        yb = yy + PITA
        for jenis, teks in baris:
            if jenis == 'pk':
                out.append('<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#FFF2CC"/>' % (x + 1.4, yb, w - 2.8, BARIS))
            elif jenis == 'x':
                out.append('<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#FCE4D6"/>' % (x + 1.4, yb, w - 2.8, BARIS))
            berat = ' font-weight="700"' if jenis in ('pk', 'fk', 'x') else ''
            warna_t = ' fill="%s"' % ('#' + warna.lstrip('#')) if jenis == 'x' else ''
            out.append('<text class="br" x="%.1f" y="%.1f"%s%s>%s</text>' % (x + 10, yb + 14.2, berat, warna_t, esc(teks)))
            yb += BARIS

    out.append('<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#FFF5F5" stroke="#C00000" '
               'stroke-width="2" stroke-dasharray="7,4"/>' % (X0, y_batas, W, 46.0))
    out.append('<text class="bts" x="%.1f" y="%.1f">⌧  %s</text>' % (X0 + 10, y_batas + 17, esc(BATAS[0])))
    out.append('<text class="br" x="%.1f" y="%.1f" fill="#C00000">%s</text>'
               % (X0 + 10, y_batas + 35, esc(BATAS[1][:118])))
    out.append('</svg>')
    svg = '\n'.join(out)

    baris_rel = '\n'.join(
        '<tr%s><td>%d</td><td>%s</td><td><code>%s</code></td><td><code>%s</code></td>'
        '<td><code>%s</code></td><td>%s</td><td>%s</td><td>%s</td></tr>'
        % (' class="sekat"' if r[0] == 'SEKAT' else '', i, r[0], r[1], r[2], r[3], r[4], r[5], esc(r[6]))
        for i, r in enumerate(RELASI, 1))

    baris_sam = '\n'.join(
        '<tr class="sekat"><td>%s</td><td>%s</td><td><code>%s</code></td><td><code>%s</code></td>'
        '<td><code>%s</code></td><td>%s</td><td>%s</td></tr>'
        % (r[0], r[1], r[2], r[3], r[4], r[5], esc(r[8]))
        for r in SAMBUNGAN)

    html = HTML_KERANGKA % dict(n_tabel=n_tabel, n_rel=len(RELASI), svg=svg,
                                baris_rel=baris_rel, baris_sam=baris_sam)
    p = os.path.join(KELUARAN, 'ERD-TREATY-MASUK.html')
    open(p, 'w', encoding='utf-8').write(html)
    return p


HTML_KERANGKA = """<!DOCTYPE html>
<html lang="id"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Peta Skema Treaty Masuk</title>
<style>
:root{--tinta:#14243a;--redup:#5b6b80;--garis:#d6dfea;--latar:#f5f7fa;--kertas:#fff;
      --sekat:#C00000;--pohon:#1F4E79;}
@media (prefers-color-scheme: dark){:root:not([data-theme="light"]){
  --tinta:#e8edf4;--redup:#9fb0c4;--garis:#2b3644;--latar:#10151c;--kertas:#161c25;}}
:root[data-theme="dark"]{--tinta:#e8edf4;--redup:#9fb0c4;--garis:#2b3644;--latar:#10151c;--kertas:#161c25;}
*{box-sizing:border-box}
body{margin:0;background:var(--latar);color:var(--tinta);
     font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Arial,sans-serif;line-height:1.5;}
.bungkus{max-width:1400px;margin:0 auto;padding:28px 16px 64px;}
h1{font-size:1.5rem;margin:0 0 6px;letter-spacing:-.01em}
h2{font-size:1.05rem;margin:38px 0 10px;padding-bottom:6px;border-bottom:1px solid var(--garis)}
.sub{color:var(--redup);font-size:.86rem;margin:0 0 22px}
.papan{background:var(--kertas);border:1px solid var(--garis);border-radius:10px;padding:10px;overflow:auto;}
svg{display:block}
table{border-collapse:collapse;width:100%%;font-size:.82rem;background:var(--kertas);
      border:1px solid var(--garis);border-radius:10px;overflow:hidden}
th{background:#1F4E79;color:#fff;text-align:left;padding:7px 9px;font-weight:600;font-size:.76rem;letter-spacing:.02em}
td{padding:6px 9px;border-top:1px solid var(--garis);vertical-align:top}
tr.sekat td{color:var(--sekat);font-weight:600}
code{font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;font-size:.95em}
.lgd{display:flex;flex-wrap:wrap;gap:8px 20px;margin:12px 0 0;font-size:.8rem;color:var(--redup)}
.lgd span{display:flex;align-items:center;gap:7px}
.sw{width:26px;height:12px;border-radius:2px;display:inline-block}
.ln{width:30px;height:0;display:inline-block}
.kartu{display:flex;flex-wrap:wrap;gap:10px;margin:16px 0 0}
.kartu div{background:var(--kertas);border:1px solid var(--garis);border-radius:8px;
           padding:8px 12px;font-size:.8rem;color:var(--redup)}
.kartu b{color:var(--tinta);font-size:1.05rem}
.awas{background:#FFF5F5;border:1px solid #C00000;border-radius:8px;padding:10px 14px;
      margin:16px 0 0;font-size:.84rem;color:#C00000}
@media (max-width:640px){.bungkus{padding:20px 16px 48px}h1{font-size:1.2rem}}
</style></head><body><div class="bungkus">

<h1>Peta Skema Treaty Masuk &mdash; Nusantara Re</h1>
<p class="sub">Menyatukan <b>Treaty In</b> dan <b>Treaty In Adjustment</b> menjadi satu peta,
skema <code>TREATY_MASUK</code>. Hanya PK dan FK, ditambah satu pembeda cabang.
Dibangkitkan oleh <code>alat/buat-skema-treaty-masuk.py</code> dari
<code>SPEC-MODEL-DATA.md</code>, <code>ERD.md</code>, <code>STRUKTUR-DATA.md</code>, dan
<code>KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md</code>. <b>Turunan</b> &mdash; bila ia berbeda dari
sumbernya, alatnya yang salah.</p>

<div class="kartu">
  <div><b>%(n_tabel)d</b><br>tabel</div>
  <div><b>%(n_rel)d</b><br>relasi</div>
  <div><b>2</b><br>relasi melintasi sekat</div>
  <div><b>1</b><br>entitas milik Adjustment</div>
  <div><b>1</b><br>batas belum diketahui</div>
</div>

<div class="awas"><b>Sekat antarmodul tipis, dan itu hasil gerbang G1.</b>
Sebuah penyesuaian <b>adalah</b> sebuah <code>VERSI_KONTRAK</code> &mdash; bukan dokumen tersendiri.
Adjustment karena itu hanya menyumbang satu entitas: catatan selisihnya.</div>

<h2>Peta</h2>
<div class="papan">
%(svg)s
</div>
<div class="lgd">
  <span><i class="sw" style="background:#1F4E79"></i>inti &mdash; kontrak dan versi</span>
  <span><i class="sw" style="background:#C55A11"></i>anak versi, kedua cabang</span>
  <span><i class="sw" style="background:#2E75B6"></i>cabang</span>
  <span><i class="sw" style="background:#31859C"></i>potongan &amp; penyebaran</span>
  <span><i class="sw" style="background:#C00000"></i>Adjustment &mdash; melintasi sekat</span>
  <span><i class="sw" style="background:#7030A0"></i>acuan bersama</span>
  <span><i class="sw" style="background:#808080"></i>di luar gelombang ini</span>
  <span><i class="sw" style="background:#FFF2CC;border:1px solid #d9c98a"></i>baris PK</span>
  <span><i class="sw" style="background:#FCE4D6;border:1px solid #C55A11"></i>pembeda cabang</span>
  <span><i class="ln" style="border-top:3px solid #1F4E79"></i>induk &rarr; anak</span>
  <span><i class="ln" style="border-top:3px dashed #C00000"></i>relasi lintas-sekat</span>
</div>

<h2>Sambungan antarmodul &mdash; seluruhnya</h2>
<table><thead><tr><th>#</th><th>Jenis</th><th>Dari</th><th>Ke</th><th>Kunci</th><th>Kard.</th>
<th>Keterangan</th></tr></thead><tbody>
%(baris_sam)s
</tbody></table>

<h2>Daftar relasi</h2>
<table><thead><tr><th>#</th><th>Lajur</th><th>Induk</th><th>Anak</th><th>Kunci tamu</th>
<th>Kard.</th><th>Perilaku hapus</th><th>Catatan</th></tr></thead><tbody>
%(baris_rel)s
</tbody></table>

</div></body></html>
"""


if __name__ == '__main__':
    p1, n = tulis_xlsx()
    p2 = tulis_html(n)
    print('xlsx :', os.path.normpath(p1))
    print('html :', os.path.normpath(p2))
    print('tabel:', n, '| relasi:', len(RELASI))
