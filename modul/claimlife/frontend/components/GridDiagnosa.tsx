// Grid diagnosa SATU PESERTA — butir bd, A3 kelompok Medis; isi kotak MEDICAL ADVISOR (RincianPeserta).
//
// Meniru `Section/ClaimLifeDetailGCNM.xml` b3923 `.DiagnoseList`, dibaca
// 27-09-2026 bersama `Activity/SetDisease.xml` dan `Activity/SetSTS_Reject.xml`.
//
// ⛔ BANYAK per peserta, dan buktinya TOMBOLNYA — grid ber-`Add` b4690 DAN
// ber-`Delete` b6160 (lihat CariDiagnosa untuk riwayat OQ-K.2).
//
// ⭐ DIRAPIKAN 11-10-2026 (work owner): urutan kolom
//
//   Find Disease (tombol per baris, b5061) | DIAGNOSE b4188 | ICD CODE b4337 |
//   GROUP DIAGNOSE b4490 (dropdown b5863) | Delete b6160
//
// dengan `+ Add` di kepala kotak; teks berbahasa Inggris ringkas. Baca saja: hanya DIAGNOSE | ICD CODE | GROUP DIAGNOSE.
//
// ⛔ DIAGNOSE dan ICD CODE READ-ONLY (b5374, b5566): nilainya dari katalog lewat
// `Find Disease` -> `Choose` b2509 -> `SetDisease`. Mengetiknya bebas berarti
// diagnosa di luar katalog dapat masuk.
//
// ⛔ GROUP DIAGNOSE — butir **bf**, OQ-L. Dropdown b5863 ber-`pyListSource
// associated`: daftar pilihannya hidup pada rule properti `.GROUPDIAGNOSE`
// kelas `Data-DiagnoseLife`, yang TIDAK ada di ekspor korpus. Kolomnya ADA
// (migrasi 018); daftar pilihannya TIDAK DIKARANG — `PILIHAN_KELOMPOK_DIAGNOSA`
// kosong sampai sumbernya diketahui; nilai tersimpan tetap ditampilkan.
//
// ⛔ GERBANG (11-10-2026): dapat diubah HANYA di tahap Medical Check oleh
// Medical Advisor (b4908, b6401 pyPosition) dan selama peserta belum diputus
// (`STS_REJECT` b4682, b5059, b5870, b6152). Pemanggil menghitungnya
// (`boleh`); server menegakkan yang sama (409/403).

import { useState } from 'react'

import { DETAIL, DIAGNOSA } from '../labels'
import { hapusDiagnosa, tambahDiagnosa, ubahDiagnosa, type Diagnosa, type Penyakit, type Peserta } from '../api'
import { pesanGalat } from '../../../../inti/frontend/klien'
import { CariDiagnosa } from './CariDiagnosa'

/**
 * Daftar pilihan GROUP DIAGNOSE (b5863 `pyListSource associated`, OQ-L). KOSONG: sumbernya tidak ada di korpus dan
 * belum diketahui tabelnya di DEV — tidak dikarang. Diisi bila work owner menunjuk tabel/daftarnya.
 */
export const PILIHAN_KELOMPOK_DIAGNOSA: readonly string[] = []

/** Keterangan dropdown GROUP DIAGNOSE yang daftarnya belum ada. */
export const KELOMPOK_BELUM_ADA = 'Group Diagnose list is not configured yet.'

/** Sel yang belum diisi (`Add` menambah baris KOSONG — wajar, bukan kerusakan; ADR-U-0027). */
export function selDiagnosa(nilai: string): string {
  return nilai.trim() === '' ? '—' : nilai
}

/** Kalimat grid kosong; null bila ada baris. */
export function ringkasanDiagnosa(daftar: Diagnosa[], boleh: boolean): string | null {
  if (daftar.length > 0) return null
  return boleh ? 'No diagnosis yet. Click + Add to add a row.' : 'No diagnosis recorded.'
}

/** Pilihan dropdown satu baris: daftar + nilai tersimpan (bila di luar daftar, tetap ditampilkan). */
export function pilihanKelompok(daftar: readonly string[], tersimpan: string): string[] {
  const t = tersimpan.trim()
  return t !== '' && !daftar.includes(t) ? [t, ...daftar] : [...daftar]
}

export function GridDiagnosa({
  klaimID,
  peserta,
  onBerubah,
  boleh,
  pilihan = PILIHAN_KELOMPOK_DIAGNOSA,
}: {
  klaimID: string
  peserta: Peserta
  /** Dipanggil sesudah setiap perubahan supaya pemanggil membaca ulang. */
  onBerubah: () => void | Promise<void>
  /** Dapat diubah: tahap Medical Check, Medical Advisor, peserta belum diputus (dihitung pemanggil). */
  boleh: boolean
  /** Daftar GROUP DIAGNOSE (bawaan: PILIHAN_KELOMPOK_DIAGNOSA). */
  pilihan?: readonly string[]
}) {
  const [sibuk, setSibuk] = useState(false)
  const [galat, setGalat] = useState<string | null>(null)
  const kosong = ringkasanDiagnosa(peserta.diagnosa, boleh)

  async function jalankan(kerja: () => Promise<void>): Promise<void> {
    if (sibuk) return
    setSibuk(true)
    setGalat(null)
    try {
      await kerja()
      await onBerubah()
    } catch (e) {
      setGalat(pesanGalat(e) ?? 'The diagnosis could not be changed.')
    } finally {
      setSibuk(false)
    }
  }

  // `Choose` b2509 -> `SetDisease`: DUA kolom dari katalog, kelompok tetap.
  async function pilihUntuk(d: Diagnosa, p: Penyakit): Promise<void> {
    await jalankan(() =>
      ubahDiagnosa(klaimID, peserta.id, d.id, { kodeIcd: p.kodeIcd, nama: p.nama, groupDiagnose: d.groupDiagnose }),
    )
  }

  async function pilihKelompok(d: Diagnosa, kelompok: string): Promise<void> {
    await jalankan(() =>
      ubahDiagnosa(klaimID, peserta.id, d.id, { kodeIcd: d.kodeIcd, nama: d.nama, groupDiagnose: kelompok }),
    )
  }

  const tambah = (
    <button
      type="button"
      className="btn btn--ghost btn--sm"
      disabled={sibuk}
      onClick={() => {
        void jalankan(async () => {
          await tambahDiagnosa(klaimID, peserta.id)
        })
      }}
    >
      + {DIAGNOSA.tambah}
    </button>
  )

  return (
    <section className="diagnosa-grid" aria-label={DIAGNOSA.kolomNama}>
      {boleh && <div className="diagnosa-grid__kepala">{tambah}</div>}
      {galat !== null && (
        <p role="alert" className="alert alert--error">
          {galat}
        </p>
      )}
      {kosong !== null && <p className="diagnosa-grid__kosong">{kosong}</p>}

      {peserta.diagnosa.length > 0 && (
        <div className="inbox__gulir">
          <table className="diagnosa-grid__tabel">
            <thead>
              <tr>
                {boleh && <th scope="col">{DETAIL.cariPenyakit}</th>}
                <th scope="col">{DIAGNOSA.kolomNama}</th>
                <th scope="col">{DIAGNOSA.kolomIcd}</th>
                <th scope="col">{DIAGNOSA.kolomKelompok}</th>
                {boleh && <th scope="col">{DIAGNOSA.hapus}</th>}
              </tr>
            </thead>
            <tbody>
              {peserta.diagnosa.map((d) => {
                const daftar = pilihanKelompok(pilihan, d.groupDiagnose)
                return (
                  <tr key={d.id}>
                    {boleh && (
                      <td>
                        {/* `Find Disease` b5061 berdiri DI DALAM baris (sel 37): barisnya konteks `Choose`. */}
                        <CariDiagnosa onPilih={(p) => pilihUntuk(d, p)} sibuk={sibuk} />
                      </td>
                    )}
                    {/* ⛔ Teks, bukan input: `Read-only` b5374 dan b5566. */}
                    <td>{selDiagnosa(d.nama)}</td>
                    <td>{selDiagnosa(d.kodeIcd)}</td>
                    <td>
                      {boleh ? (
                        <select
                          className="field__input"
                          aria-label={`${DIAGNOSA.kolomKelompok} ${d.urutan}`}
                          value={d.groupDiagnose.trim()}
                          disabled={sibuk || pilihan.length === 0}
                          title={pilihan.length === 0 ? KELOMPOK_BELUM_ADA : undefined}
                          onChange={(e) => {
                            void pilihKelompok(d, e.target.value)
                          }}
                        >
                          <option value="">—</option>
                          {daftar.map((k) => (
                            <option key={k} value={k}>
                              {k}
                            </option>
                          ))}
                        </select>
                      ) : (
                        selDiagnosa(d.groupDiagnose)
                      )}
                    </td>
                    {boleh && (
                      <td>
                        <button
                          type="button"
                          className="btn btn--ghost btn--sm"
                          disabled={sibuk}
                          onClick={() => {
                            void jalankan(() => hapusDiagnosa(klaimID, peserta.id, d.id))
                          }}
                        >
                          {DIAGNOSA.hapus}
                        </button>
                      </td>
                    )}
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      )}

      {boleh && peserta.diagnosa.length > 0 && <div className="diagnosa-grid__kaki">{tambah}</div>}
      {boleh && pilihan.length === 0 && peserta.diagnosa.length > 0 && (
        <p className="diagnosa-grid__catatan">{KELOMPOK_BELUM_ADA}</p>
      )}
    </section>
  )
}
