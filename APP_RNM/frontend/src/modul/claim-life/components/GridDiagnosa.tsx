// Grid diagnosa SATU PESERTA — butir bd, A3 kelompok Medis.
//
// Meniru `Section/ClaimLifeDetailGCNM.xml` b3923 `.DiagnoseList`, dibaca
// 27-09-2026 bersama `Activity/SetDisease.xml` dan `Activity/SetSTS_Reject.xml`.
//
// ⛔ BANYAK per peserta, dan buktinya TOMBOLNYA — bukan bentuk gridnya. Grid
// dapat berarti tampilan satu baris; grid ber-`Add` b4690 DAN ber-`Delete`
// b6160 tidak dapat. Ronde sebelumnya berhenti pada gridnya dan menyimpulkan
// "satu lawan banyak, tidak dapat diputuskan" (OQ-K.2); jawabannya dua baris
// di bawah tempat pembacaan itu berhenti.
//
// ⛔ TIGA kolom, dan kepalanya VERBATIM — termasuk huruf besarnya:
//
//   b4188 `DIAGNOSE`        <- `.DISEASE` b5422, `Read-only` b5374
//   b4337 `ICD CODE`        <- `.ICDCODE` b5616, `Read-only` b5566
//   b4490 `GROUP DIAGNOSE`  <- `.GROUPDIAGNOSE` b5860, `pxDropdown` b5863
//
// ⛔ Kedua kolom pertama READ-ONLY, dan itu bukan kekurangan: nilainya datang
// dari baris katalog yang dipilih lewat `Find Disease` b5061 -> `Choose`
// b2509 -> `SetDisease` b260/b307. Membuatnya dapat diketik berarti nama
// diagnosa yang tidak ada di katalog dapat masuk, dan tidak ada satu pun yang
// akan membandingkannya lagi.
//
// ⛔ `GROUP DIAGNOSE` DINYATAKAN belum tersedia — butir **bf**, OQ-L. Dropdown
// b5863 ber-`pyListSource associated`: daftar pilihannya hidup pada rule
// properti `.GROUPDIAGNOSE` kelas `Data-DiagnoseLife`, yang tidak ada di
// ekspor dan tidak ada di katalog DEV. Kolomnya ADA (migrasi 018); nilainya
// TIDAK DIKARANG. Menawarkan daftar tebakan berarti memaksa orang memilih
// kelompok yang salah, dan salahnya baru terlihat berbulan-bulan kemudian.
//
// ⛔ GERBANGNYA `STS_REJECT` PESERTA, satu kalimat di empat tempat di dalam
// grid (b4682, b5059, b5870, b6152). Layar bukan penjaga — backend menolak
// dengan 409 — tetapi tombol yang tetap hidup padahal pasti ditolak adalah
// tombol yang mengajari orang mengabaikan galat.

import { useState } from 'react'

import { DIAGNOSA } from '../labels'
import {
  bolehUbahDiagnosa,
  hapusDiagnosa,
  tambahDiagnosa,
  ubahDiagnosa,
  type Diagnosa,
  type Penyakit,
  type Peserta,
} from '../api'
import { pesanGalat } from '../../../inti/klien'
import { CariDiagnosa } from './CariDiagnosa'
import { BelumTersedia } from '../../../inti/components/ui/dasar'

/**
 * Menyusun kalimat untuk sel yang belum diisi.
 *
 * ⛔ Bukan sel kosong. `Add` b4700 menambah baris KOSONG — itu keadaan yang
 * WAJAR, bukan kerusakan — dan sel kosong tidak dapat dibedakan dari kolom
 * yang gagal terbaca. ADR-U-0027: kosong berarti sesuatu, dan yang berarti
 * sesuatu harus terlihat.
 */
export function selDiagnosa(nilai: string): string {
  return nilai.trim() === '' ? '—' : nilai
}

/**
 * Menyusun kalimat ringkas isi grid.
 *
 * Dipisah supaya dapat diuji tanpa DOM.
 */
export function ringkasanDiagnosa(daftar: Diagnosa[], boleh: boolean): string {
  if (daftar.length === 0) {
    return boleh
      ? 'Belum ada diagnosa. Tekan Add untuk menambah baris.'
      : 'Belum ada diagnosa, dan peserta ini sudah diputus.'
  }
  const kosong = daftar.filter((d) => d.nama.trim() === '').length
  if (kosong === 0) return `${daftar.length} diagnosa.`
  return `${daftar.length} diagnosa, ${kosong} di antaranya belum diisi.`
}

export function GridDiagnosa({
  klaimID,
  peserta,
  onBerubah,
}: {
  klaimID: string
  peserta: Peserta
  /** Dipanggil sesudah setiap perubahan supaya pemanggil membaca ulang. */
  onBerubah: () => void | Promise<void>
}) {
  const boleh = bolehUbahDiagnosa(peserta)
  const [sibuk, setSibuk] = useState(false)
  const [galat, setGalat] = useState<string | null>(null)

  async function jalankan(kerja: () => Promise<void>): Promise<void> {
    if (sibuk) return
    setSibuk(true)
    setGalat(null)
    try {
      await kerja()
      await onBerubah()
    } catch (e) {
      setGalat(pesanGalat(e) ?? 'Perubahan diagnosa gagal.')
    } finally {
      setSibuk(false)
    }
  }

  // `Choose` b2509 -> `SetDisease`: DUA kolom dari katalog, kelompok tetap.
  async function pilihUntuk(d: Diagnosa, p: Penyakit): Promise<void> {
    await jalankan(() =>
      ubahDiagnosa(klaimID, peserta.id, d.id, {
        kodeIcd: p.kodeIcd,
        nama: p.nama,
        groupDiagnose: d.groupDiagnose,
      }),
    )
  }

  return (
    <section className="diagnosa-grid">
      <h4 className="diagnosa-grid__judul">{DIAGNOSA.kolomNama}</h4>
      <p role="status">{ringkasanDiagnosa(peserta.diagnosa, boleh)}</p>
      {galat !== null && <p role="alert">{galat}</p>}

      {peserta.diagnosa.length > 0 && (
        <table className="diagnosa-grid__tabel">
          <thead>
            <tr>
              <th scope="col">{DIAGNOSA.kolomNama}</th>
              <th scope="col">{DIAGNOSA.kolomIcd}</th>
              <th scope="col">{DIAGNOSA.kolomKelompok}</th>
              <th scope="col">Aksi</th>
            </tr>
          </thead>
          <tbody>
            {peserta.diagnosa.map((d) => (
              <tr key={d.id}>
                {/* ⛔ Teks, bukan input: `Read-only` b5374 dan b5566. */}
                <td>{selDiagnosa(d.nama)}</td>
                <td>{selDiagnosa(d.kodeIcd)}</td>
                <td>
                  {/* ⛔ Butir bf — daftar pilihannya belum ada, OQ-L.
                      Nilai yang SUDAH tersimpan tetap ditampilkan; yang
                      belum ada hanyalah cara memilihnya. */}
                  {d.groupDiagnose.trim() !== '' ? (
                    d.groupDiagnose
                  ) : (
                    <BelumTersedia apa={DIAGNOSA.kolomKelompok} />
                  )}
                </td>
                <td>
                  {boleh ? (
                    <>
                      {/* `Find Disease` b5061 berdiri DI DALAM baris grid
                          (sel 37), bukan di luarnya — karena itu barisnya
                          yang menjadi konteks `Choose`. */}
                      <CariDiagnosa
                        onPilih={(p) => pilihUntuk(d, p)}
                        sibuk={sibuk}
                      />{' '}
                      <button
                        type="button"
                        disabled={sibuk}
                        onClick={() => {
                          void jalankan(() =>
                            hapusDiagnosa(klaimID, peserta.id, d.id),
                          )
                        }}
                      >
                        {DIAGNOSA.hapus}
                      </button>
                    </>
                  ) : (
                    <span>Terkunci</span>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {boleh ? (
        <p className="diagnosa-grid__aksi">
          <button
            type="button"
            disabled={sibuk}
            onClick={() => {
              void jalankan(async () => {
                await tambahDiagnosa(klaimID, peserta.id)
              })
            }}
          >
            {DIAGNOSA.tambah}
          </button>
        </p>
      ) : (
        <p className="diagnosa-grid__terkunci">
          Peserta ini sudah diputus; diagnosanya tidak dapat diubah lagi
          (gerbang <code>.STS_REJECT</code> b4682).
        </p>
      )}
    </section>
  )
}
