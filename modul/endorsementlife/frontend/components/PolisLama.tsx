// Popup `View Old Policy` - `ViewOldPolicy_EDM` (QR), `_QP`, `_TP`, `_TR` (tiket 03).
//
// Tombolnya di layar yatim `ShowLifePremiumSummary_EDM` (b64965/b65522/b66083/b66640) dipindah ke kepala
// `InputEDMLife` (RALAT R02, OQ-EDM-004): satu tombol, popup sesuai `.Type` kasus. Varian QP/TP/TR di Pega
// kosong karena wadah luarnya bersyarat `.Type = 'QR'` (R03) - di sini grid dalamnya diisi (spec AC 9).

import { useEffect, useState } from 'react'

import { Gagal, Halaman, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambilPolisLama, type PolisLamaEDM } from '../api'
import { POLIS_LAMA_PESERTA, POLIS_LAMA_REKAP, type KolomKorpus, type TipeEDM } from '../kolomKorpus'
import { POLIS_LAMA_EDM, UMUM_EDM } from '../labels'
import { UKURAN_HALAMAN_EDM } from '../tampilan'
import { TabelKorpus } from './TabelKorpus'

/**
 * Kolom korpus yang DISEMBUNYIKAN dari popup View Old Policy (tabel peserta dan rekap) - keputusan work owner
 * 08-10-2026, tampilan saja: daftar korpus lengkap tetap di `kolomKorpus.ts` sebagai bukti, data backend tidak
 * berubah. `CLAIM_AMOUNT` ada di peserta QR/QP/TR dan rekap QP (tidak di TP).
 */
export const KOLOM_POLIS_LAMA_TERSEMBUNYI: readonly string[] = ['CLAIM_AMOUNT']

/** Kolom popup yang tampil: daftar korpus dikurangi `KOLOM_POLIS_LAMA_TERSEMBUNYI`, urutan tetap. */
export function kolomPolisLamaTampil(kolom: readonly KolomKorpus[]): readonly KolomKorpus[] {
  return kolom.filter((k) => !KOLOM_POLIS_LAMA_TERSEMBUNYI.includes(k.kolom))
}

/**
 * Kolom "Endorsement No" tabel rekap View Old Policy - kolom TAMPILAN bukan korpus (keputusan work owner 10-10-2026,
 * opsi a): tanpa kolom ini baris NB dan baris EDM /01 tampak kembar (COB, PL NUMBER, CURRENCY sama). Isinya dari
 * `PL_NUMBER_EDM` yang sudah ada di respons (`teksEndorsementNo`); daftar korpus `POLIS_LAMA_REKAP` tidak diubah.
 */
export const KOLOM_ENDORSEMENT_NO: KolomKorpus = {
  label: POLIS_LAMA_EDM.endorsementNo,
  kolom: 'ENDORSEMENT_NO_TAMPIL',
  jenis: 't',
  berkas: '',
  baris: 0,
  tag: 'pyValue',
}

/**
 * Teks Endorsement No satu baris rekap: `PL_NUMBER_EDM` apa adanya bila berakhiran "/<angka>" dan berbeda dari
 * `PL_NUMBER`; selain itu (kosong - baris NB PremiumList menyalin `T_PREMIUM_LIST.PL_NUMBER_EDM` NB yang tidak
 * pernah diisi; sama dengan PL NUMBER; atau tanpa akhiran) "New Business".
 */
export function teksEndorsementNo(b: Record<string, string>): string {
  const edm = (b.PL_NUMBER_EDM ?? '').trim()
  if (edm === '' || edm === (b.PL_NUMBER ?? '').trim() || !/\/\d+$/.test(edm)) return POLIS_LAMA_EDM.newBusiness
  return edm
}

/** Kolom rekap yang tampil: kolom korpus (tanpa yang disembunyikan) + Endorsement No tepat sesudah PL NUMBER. */
export function kolomRekapPolisLama(kolom: readonly KolomKorpus[]): readonly KolomKorpus[] {
  const tampil = [...kolomPolisLamaTampil(kolom)]
  const i = tampil.findIndex((k) => k.kolom === 'PL_NUMBER')
  tampil.splice(i < 0 ? tampil.length : i + 1, 0, KOLOM_ENDORSEMENT_NO)
  return tampil
}

/** Baris rekap dengan kolom tampilan Endorsement No. */
export function barisRekapPolisLama(rekap: ReadonlyArray<Record<string, string>>): Array<Record<string, string>> {
  return rekap.map((b) => ({ ...b, [KOLOM_ENDORSEMENT_NO.kolom]: teksEndorsementNo(b) }))
}

/** Nomor baris pertama tabel peserta halaman ini - berlanjut antar halaman (07-10-2026). */
export function nomorAwalPeserta(halaman: number, ukuran: number = UKURAN_HALAMAN_EDM): number {
  return (Math.max(halaman, 1) - 1) * ukuran + 1
}

/** Tipe yang dikenal popup; tipe lain jatuh ke varian induk `QR` (seperti tombol `.Type=='QR'`). */
export function tipePopup(t: string): TipeEDM {
  return t === 'QP' || t === 'TP' || t === 'TR' ? t : 'QR'
}

export default function PolisLama({ kasusId, tipe, onTutup }: { kasusId: string; tipe: string; onTutup: () => void }) {
  const [halaman, setHalaman] = useState(1)
  const [isi, setIsi] = useState<PolisLamaEDM | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const t = tipePopup(tipe)

  useEffect(() => {
    let batal = false
    ambilPolisLama(kasusId, halaman)
      .then((p) => {
        if (!batal) setIsi(p)
      })
      .catch((e: unknown) => {
        if (!batal) setGalat(e)
      })
    return () => {
      batal = true
    }
  }, [kasusId, halaman])

  return (
    <Modal judul={POLIS_LAMA_EDM.viewOldPolicy} onTutup={onTutup} labelBatal={UMUM_EDM.tutup} penuh>
      {galat !== null && <Gagal galat={galat} />}
      {isi === null && galat === null && <Memuat pesan={UMUM_EDM.memuat} />}
      {isi !== null && isi.peserta.length === 0 && <Kosong pesan={UMUM_EDM.kosong} />}
      {isi !== null && isi.peserta.length > 0 && (
        <>
          {/* Kolom "No" paling kiri, berlanjut antar halaman (keputusan work owner 07-10-2026); rekap tanpa nomor. */}
          <TabelKorpus
            kolom={kolomPolisLamaTampil(POLIS_LAMA_PESERTA[t])}
            baris={isi.peserta.map((p) => p.nilai)}
            kunci={(_, i) => isi.peserta[i]?.id ?? String(i)}
            nomorAwal={nomorAwalPeserta(halaman)}
          />
          <Halaman halaman={halaman} ukuran={UKURAN_HALAMAN_EDM} total={isi.total} onPindah={setHalaman} />
        </>
      )}
      {isi !== null && isi.rekap.length > 0 && (
        // Kolom Endorsement No sesudah PL NUMBER (keputusan work owner 10-10-2026).
        <TabelKorpus
          kolom={kolomRekapPolisLama(POLIS_LAMA_REKAP[t])}
          baris={barisRekapPolisLama(isi.rekap)}
          kunci={(b, i) => `${b.PL_NUMBER_EDM ?? ''}/${b.CURRENCY ?? ''}/${i}`}
        />
      )}
    </Modal>
  )
}
