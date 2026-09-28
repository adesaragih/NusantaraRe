// Menyimpan Blob sebagai berkas di mesin pemakai.
//
// ⛔ Satu-satunya tempat `a.href` dipasang untuk unduhan berkas dari backend,
// dan nilainya SELALU objek URL lokal (`URL.createObjectURL`), tidak pernah
// alamat backend: pranala ke backend tidak membawa header identitas
// (`services/unduhdokumen.test.ts`). Pola yang sama dengan `unduhXlsx`.

/** simpanBlob memicu unduhan `blob` dengan nama berkas aslinya. */
export function simpanBlob(blob: Blob, namaBerkas: string): void {
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = namaBerkas
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}
