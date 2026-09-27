import { beforeEach, describe, expect, it, vi } from 'vitest'

import {
  ApiFailure,
  bolehTutupDiLayar,
  kasusTertutup,
  penghalangDariGalat,
  tutupKlaim,
  STATUS_WORK_SELESAI,
  TAHAP_PENAWAR_TUTUP,
  type Klaim,
} from './api'

// Uji sisi klien penutupan kasus — butir bb.
//
// Yang dijaga: tombolnya hanya ditawarkan pada tahap yang di Pega
// menawarkannya, penutupannya memakai POST, dan 409 membawa SELURUH
// penghalang beserta kalimat rule-nya.

/** Klaim minimal; hanya medan yang diuji berkas ini yang berarti. */
function klaimUji(tahap: string, statusWork = ''): Klaim {
  return {
    id: 'UJI-KLAIM-1',
    nomorKlaim: 'UJI-CLM-1',
    nomorPolis: 'UJI-POL-1',
    namaBisnis: '',
    kodeStatus: '',
    statusTurunan: '',
    claimRetro: { amount: '0', currency: 'IDR' },
    peserta: [],
    cacahBaris: 0,
    tahap,
    statusWork,
  } as unknown as Klaim
}

describe('bolehTutupDiLayar', () => {
  it('hanya kedua tahap yang layarnya menawarkan Close Claim', () => {
    // ⛔ Cacah berkas, bukan pilihan kami: `pyLocalAction>CloseClaim` ada di
    // TEPAT DUA section — InputOSClaimLife (b22837, b22988) dan
    // InputAkseptasiClaimLife (b21457, b21602).
    expect(bolehTutupDiLayar(klaimUji('Outstanding Claim'))).toBe(true)
    expect(bolehTutupDiLayar(klaimUji('Claim Analis'))).toBe(true)
    expect(bolehTutupDiLayar(klaimUji('Medical Check'))).toBe(false)
    expect(bolehTutupDiLayar(klaimUji('Input Register'))).toBe(false)
  })

  it('tahap KOSONG tidak menawarkannya — gagal TERTUTUP', () => {
    // Backend tidak pernah menyisipkan ke `T_WORK_CLAIM`, jadi klaim tanpa
    // baris work adalah keadaan nyata. Menawarkan tombolnya "kalau-kalau"
    // berarti menawarkan penutupan atas kasus yang tahapnya tidak diketahui.
    expect(bolehTutupDiLayar(klaimUji(''))).toBe(false)
    expect(bolehTutupDiLayar(null)).toBe(false)
  })

  it('kasus yang sudah tertutup tidak menawarkannya lagi', () => {
    const tutup = klaimUji('Outstanding Claim', STATUS_WORK_SELESAI)
    expect(kasusTertutup(tutup)).toBe(true)
    expect(bolehTutupDiLayar(tutup)).toBe(false)
  })

  it('status kerja dibandingkan PERSIS', () => {
    // Sama seperti gerbangnya di Go: pembanding yang lebih longgar daripada
    // aslinya akan memperlakukan kasus yang masih terbuka sebagai tertutup,
    // lalu menolak setiap perubahan tanpa jalan keluar yang terlihat.
    for (const mirip of [
      ' Resolved-Completed',
      'resolved-completed',
      'Resolved',
      'Resolved-Complete',
    ]) {
      expect(kasusTertutup(klaimUji('Outstanding Claim', mirip))).toBe(false)
    }
  })

  it('kedua nama tahap VERBATIM', () => {
    expect(TAHAP_PENAWAR_TUTUP).toEqual(['Outstanding Claim', 'Claim Analis'])
  })
})

describe('penghalangDariGalat', () => {
  it('membawa SELURUH penghalang, bukan yang pertama', () => {
    // Pega memasang pesannya di DALAM loop, sekali per peserta yang
    // tertandai. Melaporkan satu saja memaksa pemakai menutup berulang kali
    // dan menemukan satu penghalang baru setiap kali.
    const galat = new ApiFailure(409, {
      code: 'DITOLAK_BACKEND',
      message: 'klaim belum boleh ditutup: ada peserta yang belum diaksep',
      penghalang: [
        { urutan: 1, nomorSertifikat: 'UJI-001', pesan: 'UJI-001 is not approved yet, on list 1' },
        { urutan: 3, nomorSertifikat: 'UJI-003', pesan: 'UJI-003 is not approved yet, on list 3' },
      ],
    })
    const halangan = penghalangDariGalat(galat)
    expect(halangan).toHaveLength(2)
    expect(halangan[1]?.urutan).toBe(3)
  })

  it('galat lain memberi daftar KOSONG, bukan penghalang karangan', () => {
    // 403 dan 503 bukan penolakan gerbang. Memperlakukannya sebagai
    // penghalang akan menuduh peserta yang tidak bersalah.
    expect(penghalangDariGalat(new ApiFailure(403, { code: 'DITOLAK_BACKEND' }))).toEqual([])
    expect(penghalangDariGalat(new Error('putus'))).toEqual([])
  })
})

describe('tutupKlaim', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('memakai POST pada rute /tutup', async () => {
    // ⛔ POST, bukan GET. Penutupan yang dapat dipicu pranala, prefetch,
    // atau tombol kembali adalah penutupan yang terjadi tanpa ada yang
    // menekannya.
    const palsu = vi
      .spyOn(globalThis, 'fetch')
      .mockResolvedValue(new Response(null, { status: 204 }))

    await tutupKlaim('UJI-KLAIM-1')

    expect(palsu).toHaveBeenCalledTimes(1)
    const panggilan = palsu.mock.calls[0]
    expect(panggilan).toBeDefined()
    expect(String(panggilan?.[0])).toContain('/api/klaim-life/UJI-KLAIM-1/tutup')
    expect(panggilan?.[1]?.method).toBe('POST')
  })

  it('409 dilempar sebagai ApiFailure yang membawa penghalangnya', async () => {
    // ⚠️ mockImplementation, BUKAN mockResolvedValue: badan Response hanya
    // dapat dibaca sekali, dan uji ini memanggil dua kali. Yang kedua akan
    // menerima badan yang sudah habis - dan gagalnya akan terbaca seolah
    // kodenya tidak membaca penghalang.
    vi.spyOn(globalThis, 'fetch').mockImplementation(() =>
      Promise.resolve(
        new Response(
          JSON.stringify({
            galat: 'klaim belum boleh ditutup: ada peserta yang belum diaksep',
            penghalang: [
              {
                urutan: 2,
                nomorSertifikat: 'UJI-002',
                pesan: 'UJI-002 is not approved yet, on list 2',
              },
            ],
          }),
          { status: 409, headers: { 'Content-Type': 'application/json' } },
        ),
      ),
    )

    await expect(tutupKlaim('UJI-KLAIM-1')).rejects.toBeInstanceOf(ApiFailure)
    try {
      await tutupKlaim('UJI-KLAIM-1')
      expect.unreachable('penutupan yang ditolak harus melempar')
    } catch (e) {
      // ⛔ Kunci amplopnya `galat`, dan penghalangnya ikut di amplop yang
      // SAMA. Cacat `galat` vs `error` lahir dari dua bentuk amplop.
      const halangan = penghalangDariGalat(e)
      expect(halangan).toHaveLength(1)
      expect(halangan[0]?.pesan).toBe('UJI-002 is not approved yet, on list 2')
    }
  })
})
