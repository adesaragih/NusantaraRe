// Popup `View Upload` b9340 → `Section/ViewCSVResult_LifeEDM.xml`: atas isi berkas terunggah, nilainya teks mentah
// seperti daftar sementara Pega; `Generate Data Detail` b14322 mengunduhnya kembali berkepala `CSVPropHeaders` b276
// (23 kolom, tidak berubah). Kolom tabel = kolom grid Participants yang tampil (`kolomGrid`, urutan dan label sama) -
// keputusan work owner 10-10-2026 "Sama dengan grid peserta", menyimpang dari grid 38 kolom b1082 (`HASIL_CSV`).

import { useMemo, useState } from 'react'

import { Halaman, Kosong, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { simpanBlob } from '../../../../inti/frontend/lib/simpanBlob'
import { KEPALA_UNDUH_CSV, namaBerkasUnduh, rakitCSV } from '../csv'
import type { KolomKorpus } from '../kolomKorpus'
import { UMUM_EDM, UNGGAH_EDM } from '../labels'
import { kolomGrid, UKURAN_HALAMAN_EDM } from '../tampilan'
import { TabelKorpus } from './TabelKorpus'

/**
 * Kolom grid Participants sebagai teks mentah: nilai unggahan belum diurai (`10,5`, `02/01/1990` tampil apa adanya).
 * Dihitung tiap render supaya ikut bila `KOLOM_GRID_DI_RINCIAN` berubah. `berkas`/`baris` hanya untuk kunci React.
 */
export function kolomViewUpload(): KolomKorpus[] {
  return kolomGrid().map(([label, kolom], i) => ({ label, kolom, jenis: 't', berkas: '', baris: i, tag: 'pyValue' }))
}

export default function HasilCSV({ baris, onTutup }: { baris: ReadonlyArray<Record<string, string>>; onTutup: () => void }) {
  const [halaman, setHalaman] = useState(1)
  const potong = useMemo(
    () => baris.slice((halaman - 1) * UKURAN_HALAMAN_EDM, halaman * UKURAN_HALAMAN_EDM),
    [baris, halaman],
  )
  const unduh = () => {
    simpanBlob(new Blob([rakitCSV(KEPALA_UNDUH_CSV, baris)], { type: 'text/csv;charset=utf-8' }), namaBerkasUnduh(new Date()))
  }
  return (
    <Modal
      judul={UNGGAH_EDM.viewUpload}
      onTutup={onTutup}
      labelBatal={UMUM_EDM.tutup}
      penuh
      aksi={
        <button type="button" className="btn btn--sm" onClick={unduh}>
          {UNGGAH_EDM.generateDataDetail}
        </button>
      }
    >
      {baris.length === 0 ? (
        <Kosong pesan={UMUM_EDM.kosong} />
      ) : (
        <>
          <TabelKorpus kolom={kolomViewUpload()} baris={potong} kunci={(_, i) => String((halaman - 1) * UKURAN_HALAMAN_EDM + i)} />
          <Halaman halaman={halaman} ukuran={UKURAN_HALAMAN_EDM} total={baris.length} onPindah={setHalaman} />
        </>
      )}
    </Modal>
  )
}
