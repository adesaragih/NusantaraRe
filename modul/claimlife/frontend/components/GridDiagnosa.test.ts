import { describe, expect, it } from 'vitest'

import { DIAGNOSA } from '../labels'
import { bolehUbahDiagnosa, type Diagnosa, type Peserta } from '../api'
import { PILIHAN_KELOMPOK_DIAGNOSA, pilihanKelompok, ringkasanDiagnosa, selDiagnosa } from './GridDiagnosa'

// Uji grid diagnosa — butir bd.

/** Peserta seadanya; hanya medan yang diuji yang berarti. */
function peserta(kodeStatus: string, diagnosa: Diagnosa[] = []): Peserta {
  return {
    id: 'P-1',
    nomorPremiList: '',
    nomorPolis: '',
    nomorSertifikat: '006',
    mataUang: 'IDR',
    isCheck: 'true',
    tanggalKejadian: '',
    tanggalTerimaKlaim: '',
    tanggalDokumenLengkap: '',
    tanggalKonfirmasi: '',
    penandaTerimaKlaim: '',
    penandaTerimaKlaimAlasan: '',
    baris: [],
    total: {} as Peserta['total'],
    dokumen: [],
    diagnosa,
    kodeStatus,
  }
}

function baris(id: number, nama = '', kodeIcd = ''): Diagnosa {
  return {
    id,
    pesertaId: 'P-1',
    urutan: id,
    kodeIcd,
    nama,
    groupDiagnose: '',
    kodeStatus: '',
  }
}

describe('bolehUbahDiagnosa', () => {
  it('KEDUA keputusan mengunci, bukan hanya yang ditolak', () => {
    // ⛔ Inti uji ini. Gerbang `pyDisabledWhen` b4682 menyebut DUA nilai
    // dengan `||`. Layar yang hanya mengenal "2" akan membiarkan diagnosa
    // peserta yang sudah DIAKSEP tetap dapat disunting — separuh gerbang
    // hilang, dan tidak ada satu pun uji satu-kasus yang akan tahu.
    expect(bolehUbahDiagnosa(peserta('1'))).toBe(false)
    expect(bolehUbahDiagnosa(peserta('2'))).toBe(false)
  })

  it('belum diputus berarti masih boleh', () => {
    expect(bolehUbahDiagnosa(peserta(''))).toBe(true)
    expect(bolehUbahDiagnosa(peserta('0'))).toBe(true)
  })

  it('kode warisan 4 TIDAK mengunci', () => {
    // "4" ada di data warisan dan belum diputuskan work owner. Ia bukan
    // keputusan, jadi ia tidak mengunci — dan yang menebaknya mengunci akan
    // membekukan peserta yang sebenarnya masih berjalan.
    expect(bolehUbahDiagnosa(peserta('4'))).toBe(true)
  })
})

describe('selDiagnosa', () => {
  it('baris kosong ditandai, bukan dibiarkan kosong', () => {
    // ⛔ `Add` b4700 menambah baris KOSONG — keadaan yang WAJAR. Sel kosong
    // tidak dapat dibedakan dari kolom yang gagal terbaca (ADR-U-0027).
    expect(selDiagnosa('')).toBe('—')
    expect(selDiagnosa('   ')).toBe('—')
  })

  it('isi ditampilkan apa adanya', () => {
    expect(selDiagnosa('Diabetes mellitus')).toBe('Diabetes mellitus')
  })
})

describe('ringkasanDiagnosa (11-10-2026: English, ringkas)', () => {
  it('kosong dan boleh: menyebut + Add', () => {
    expect(ringkasanDiagnosa([], true)).toBe('No diagnosis yet. Click + Add to add a row.')
  })

  it('kosong dan baca saja: tanpa ajakan Add', () => {
    expect(ringkasanDiagnosa([], false)).toBe('No diagnosis recorded.')
  })

  it('ada baris: tanpa kalimat (grid yang berbicara)', () => {
    expect(ringkasanDiagnosa([baris(1, 'UJI PENYAKIT', 'A00')], true)).toBeNull()
  })
})

describe('pilihanKelompok', () => {
  it('nilai tersimpan di luar daftar tetap tampil di depan; kosong tidak ditambah', () => {
    expect(pilihanKelompok(['UJI-G1', 'UJI-G2'], 'UJI-LAMA')).toEqual(['UJI-LAMA', 'UJI-G1', 'UJI-G2'])
    expect(pilihanKelompok(['UJI-G1'], 'UJI-G1')).toEqual(['UJI-G1'])
    expect(pilihanKelompok([], ' ')).toEqual([])
  })

  it('daftar bawaan KOSONG — tidak dikarang (OQ-L)', () => {
    expect(PILIHAN_KELOMPOK_DIAGNOSA).toEqual([])
  })
})

describe('label grid diagnosa', () => {
  it('VERBATIM dari korpus, termasuk huruf besarnya', () => {
    // ⛔ Kepala kolom di `ClaimLifeDetailGCNM.xml` ditulis HURUF BESAR.
    // Merapikannya menjadi "Diagnose" berarti layar ini menyebut kolom
    // dengan kata yang tidak ada di sistem lama — dan orang yang beralih
    // membaca layar, bukan kode.
    expect(DIAGNOSA.kolomNama).toBe('DIAGNOSE')
    expect(DIAGNOSA.kolomIcd).toBe('ICD CODE')
    expect(DIAGNOSA.kolomKelompok).toBe('GROUP DIAGNOSE')
    expect(DIAGNOSA.tambah).toBe('Add')
    expect(DIAGNOSA.hapus).toBe('Delete')
    expect(DIAGNOSA.pilih).toBe('Choose')
  })
})
